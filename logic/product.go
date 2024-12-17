package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/proto"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	"github.com/samber/lo"
)

func AddProduct(ctx context.Context, userid uint64, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	if lo.IndexOf(req.Manager, userid) == -1 {
		req.Manager = append(req.Manager, userid)
	}

	userMap, err := table.GetUserBasesMapByIds(ctx, req.Manager)
	if err != nil {
		return nil, err
	}

	product := table.TblProduct{
		Name:    req.Name,
		Intro:   req.Intro,
		Creator: userid,
		Manager: types.JoinUint64(req.Manager, ","),
		Status:  consts.StatusOnline,
	}

	id, err := table.AddProduct(ctx, &product)

	if errs.Code(err) == 1062 {
		return nil, errs.Newf(errs.RetWebDuplicateProductName, "product name is duplicated")
	}

	if err != nil {
		return nil, err
	}

	member := table.TblProductMember{
		ProductID:  id,
		UserID:     userid,
		Role:       consts.ProductRoleDeveloper,
		Status:     consts.ProductMemberStatusJoined,
		JoinTime:   time.Now().Unix(),
		ExpireType: consts.ExpireTypePermanent,
		ExpireTime: 0,
	}

	_, err = table.InsertProductMember(ctx, &member)
	if err != nil {
		return nil, err
	}

	// 写入检索信息
	searchs := []*table.TblSearchKeyword{}
	searchs = append(searchs, &table.TblSearchKeyword{
		Type:     consts.SearchTypeProduct,
		Sid:      id,
		SName:    req.Name,
		Field:    "name",
		SContent: req.Name})

	if req.Intro != "" {
		searchs = append(searchs, &table.TblSearchKeyword{
			Type:     consts.SearchTypeProduct,
			Sid:      id,
			SName:    req.Name,
			Field:    "intro",
			SContent: req.Intro})
	}

	u, ok := userMap[userid]
	if ok {
		name := fmt.Sprintf("%s(%s)", u.Nickname, u.Account)
		searchs = append(searchs, &table.TblSearchKeyword{
			Type:     consts.SearchTypeProduct,
			Sid:      id,
			SName:    req.Name,
			Field:    "creator",
			SKey:     fmt.Sprint(userid),
			SContent: name})
	}

	for _, uid := range req.Manager {
		u, ok = userMap[uid]
		if ok {
			name := fmt.Sprintf("%s(%s)", u.Nickname, u.Account)
			searchs = append(searchs, &table.TblSearchKeyword{
				Type:     consts.SearchTypeProduct,
				Sid:      id,
				SName:    req.Name,
				Field:    "manager",
				SKey:     fmt.Sprint(userid),
				SContent: name})
		}
	}

	table.AddSearchKeywords(ctx, searchs)

	return &pb.AddProductResponse{ID: id}, nil
}

func UpdateProduct(ctx context.Context, userid uint64, req *pb.UpdateProductRequest) error {
	role, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if role != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	update := horm.Map{
		"name":  req.Name,
		"intro": req.Intro,
	}

	return table.UpdateProductByID(ctx, req.ProductID, update)
}

func MaintainProductManager(ctx context.Context, userid uint64, req *pb.MaintainProductManagerRequest) error {
	myRole, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if myRole != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	managerUids := lo.Uniq(req.Manager)

	members, err := table.GetProductMemberByUsers(ctx, req.ProductID, managerUids)
	if err != nil {
		return err
	}

	roleMap := map[uint64]int8{}
	for _, member := range members {
		roleMap[member.UserID] = GetProductRole(member)
	}

	for _, uid := range managerUids {
		role := roleMap[uid]
		if role == consts.ProductRoleNotJoin {
			return errs.Newf(errs.RetWebIsNotMember, "user [%d] is not member of product", uid)
		}
	}

	update := horm.Map{
		"manager": types.JoinUint64(managerUids, ","),
	}

	err = table.UpdateProductByID(ctx, req.ProductID, update)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProductStatus(ctx context.Context, userid uint64, req *pb.UpdateProductStatusRequest) error {
	role, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if role != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateProductByID(ctx, req.ProductID, update)
}

func ProductList(ctx context.Context, req *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	pageInfo, products, err := table.GetProductList(ctx, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.ProductListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      req.Page,
		Size:      req.Size,
		Products:  make([]*pb.ProductBase, len(products)),
	}

	for k, v := range products {
		ret.Products[k] = &pb.ProductBase{
			Id:        v.Id,
			Name:      v.Name,
			Intro:     v.Intro,
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Unix(),
		}
	}

	return &ret, nil
}

func ProductDetail(ctx context.Context, userid uint64,
	req *pb.ProductDetailRequest) (*pb.ProductDetailResponse, error) {
	isNil, product, err := table.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindProduct, "not find product")
	}

	managerUids := GetUserIds(product.Manager)
	userBases, err := table.GetUserBasesMapByIds(ctx, GetUserIds(product.Creator, managerUids))
	if err != nil {
		return nil, err
	}

	_, member, err := table.GetProductMemberByUser(ctx, req.ProductID, userid)
	if err != nil {
		return nil, err
	}

	ret := pb.ProductDetailResponse{
		Info: &pb.ProductBase{
			Id:        product.Id,
			Name:      product.Name,
			Intro:     product.Intro,
			Status:    product.Status,
			CreatedAt: product.CreatedAt.Unix(),
		},
		Creator:    userBases[product.Creator],
		Manager:    GetUsersFromMap(userBases, managerUids),
		Role:       GetProductRole(member, product),
		Status:     member.Status,
		ChangeRole: member.ChangeRole,
		ExpireTime: member.ExpireTime,
		OutTime:    member.OutTime,
		DBs:        []*pb.DBBase{},
	}

	dbs, err := table.GetProductDBs(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	for _, v := range dbs {
		ret.DBs = append(ret.DBs, GetDBBase(v))
	}

	return &ret, nil
}

func ProductMemberList(ctx context.Context, userid uint64,
	req *pb.ProductMemberListRequest) (*pb.ProductMemberListResponse, error) {
	_, myMember, err := table.GetProductMemberByUser(ctx, req.ProductID, userid)
	if err != nil {
		return nil, err
	}

	isNil, product, err := table.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindProduct, "not find product")
	}

	ret := pb.ProductMemberListResponse{
		Total:      0,
		TotalPage:  0,
		Page:       req.Page,
		Size:       req.Size,
		Role:       GetProductRole(myMember, product),
		Status:     myMember.Status,
		ChangeRole: myMember.ChangeRole,
		Members:    []*pb.ProductMember{},
	}

	var pageRet *proto.Detail
	var members []*table.TblProductMember

	switch ret.Role {
	case consts.ProductRoleManager:
		pageRet, members, err = table.GetProductMembersAll(ctx, req.ProductID, req.Page, req.Size)
	case consts.ProductRoleDeveloper, consts.ProductRoleOperator:
		pageRet, members, err = table.GetProductMembersJoined(ctx, req.ProductID, req.Page, req.Size)
	}

	if err != nil {
		return nil, err
	}

	if pageRet != nil {
		ret.Total = pageRet.Total
		ret.TotalPage = pageRet.TotalPage
	}

	if len(members) > 0 {
		userIds := GetUseridFromProductMember(members)
		userBases, err := table.GetUserBasesMapByIds(ctx, userIds)
		if err != nil {
			return nil, err
		}

		for _, v := range members {
			userBase, ok := userBases[v.UserID]
			if !ok {
				continue
			}

			member := pb.ProductMember{
				MemberID:   v.Id,
				Userid:     v.UserID,
				Account:    userBase.Account,
				Nickname:   userBase.Nickname,
				JoinTime:   int(v.JoinTime),
				ExpireType: v.ExpireType,
				ExpireTime: v.ExpireTime,
				OutTime:    v.OutTime,
				ChangeRole: v.ChangeRole,
			}

			member.Role, member.Status = GetProductRealRoleStatus(v, product)
			ret.Members = append(ret.Members, &member)
		}
	}

	return &ret, nil
}

// ProductJoinApply 申请加入产品 / 续期
func ProductJoinApply(ctx context.Context, userid uint64, req *pb.ProductJoinApplyRequest) error {
	isNil, member, err := table.GetProductMemberByUser(ctx, req.ProductID, userid)
	if err != nil {
		return err
	}

	if isNil { // 新成员申请加入产品
		if req.Role != consts.ProductRoleDeveloper && req.Role != consts.ProductRoleOperator {
			return errs.Newf(errs.RetWebParamEmpty, "input param [role] is invalid")
		}

		newMember := table.TblProductMember{
			ProductID:  req.ProductID,
			UserID:     userid,
			Role:       req.Role,
			Status:     consts.ProductMemberStatusApproval,
			JoinTime:   time.Now().Unix(),
			ExpireType: req.ExpireType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = table.InsertProductMember(ctx, &newMember)
		return err
	} else if GetProductRole(member) == consts.ProductRoleNotJoin { // 重新申请加入产品
		if member.Status == consts.ProductMemberStatusApproval ||
			member.Status == consts.ProductMemberStatusRenewal {
			return errs.Newf(errs.RetWebMemberUnderApproval, "under approval, please do not apply repeatedly")
		}

		if req.Role != consts.ProductRoleDeveloper && req.Role != consts.ProductRoleOperator {
			return errs.Newf(errs.RetWebParamEmpty, "input param [role] is invalid")
		}

		replace := horm.Map{
			"id":          member.Id,
			"product_id":  member.ProductID,
			"userid":      member.UserID,
			"role":        req.Role,
			"status":      consts.ProductMemberStatusApproval,
			"join_time":   time.Now().Unix(),
			"expire_type": req.ExpireType,
			"expire_time": 0,
			"out_time":    0,
			"updated_at":  time.Now(),
		}
		return table.ReplaceProductMember(ctx, replace)
	} else { // 申请续期
		if member.ExpireTime == 0 || int64(member.ExpireTime)-time.Now().Unix() > 7*86400 { // 只有7天内过期的用户才允许续期
			return errs.Newf(errs.RetWebIsMember, "user is already member of product")
		} else {
			if member.Status == consts.ProductMemberStatusApproval ||
				member.Status == consts.ProductMemberStatusRenewal {
				return errs.Newf(errs.RetWebMemberUnderApproval, "under approval, please do not apply repeatedly")
			}

			if req.Role != member.Role { // 续期不能改变角色
				return errs.Newf(errs.RetWebParamEmpty, "renewal input param [role] must be empty")
			}

			update := horm.Map{
				"status":      consts.ProductMemberStatusRenewal,
				"expire_type": req.ExpireType,
				"out_time":    0,
			}

			return table.UpdateProductMemberByID(ctx, member.Id, update)
		}
	}
}

// ProductApproval 产品权限审批
func ProductApproval(ctx context.Context, userid uint64, req *pb.ProductApprovalRequest) error {
	myRole, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if myRole != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	isNil, member, err := table.GetProductMemberByUser(ctx, req.ProductID, req.Userid)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebIsNotApply, "user has not applied for product permissions")
	}

	if member.Status != consts.ProductMemberStatusApproval && member.Status != consts.ProductMemberStatusRenewal {
		return errs.Newf(errs.RetWebMemberNotUnderApproval, "user is not in approval status")
	}

	var update horm.Map
	if req.Status == consts.ApprovalAccess {
		update = horm.Map{
			"status":      consts.ProductMemberStatusJoined,
			"expire_time": GetExpireTime(int64(member.ExpireTime), member.ExpireType),
		}

		if member.Status == consts.ProductMemberStatusApproval {
			update["join_time"] = time.Now().Unix()
		}
	} else {
		update = horm.Map{
			"status": consts.ProductMemberStatusReject,
		}
	}

	return table.UpdateProductMemberByID(ctx, member.Id, update)
}

// ProductChangeRoleApply 申请变更角色
func ProductChangeRoleApply(ctx context.Context, userid uint64, req *pb.ProductChangeRoleApplyRequest) error {
	_, member, err := table.GetProductMemberByUser(ctx, req.ProductID, userid)
	if err != nil {
		return err
	}

	role := GetProductRole(member)
	if role == consts.ProductRoleNotJoin {
		return errs.Newf(errs.RetWebIsNotMember, "user is not member of product")
	}

	if role == consts.ProductRoleExpired {
		return errs.Newf(errs.RetWebMemberExpired, "product member permission has expired, please renewal first")
	}

	if member.Status == consts.ProductMemberStatusRenewal {
		return errs.Newf(errs.RetWebMemberUnderApproval, "please complete the renewal approval first")
	}

	if member.Status != consts.ProductMemberStatusJoined {
		return errs.Newf(errs.RetWebIsNotMember, "user is not member of product")
	}

	if req.Role == member.Role {
		return nil
	}

	update := horm.Map{
		"status":      consts.ProductMemberStatusChangeRole,
		"change_role": req.Role,
	}

	return table.UpdateProductMemberByID(ctx, member.Id, update)
}

// ProductChangeRoleApproval 产品角色变更审批
func ProductChangeRoleApproval(ctx context.Context, userid uint64, req *pb.ProductApprovalRequest) error {
	myRole, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if myRole != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	_, member, err := table.GetProductMemberByUser(ctx, req.ProductID, req.Userid)
	if err != nil {
		return err
	}

	role := GetProductRole(member)
	if role == consts.ProductRoleNotJoin {
		return errs.Newf(errs.RetWebIsNotMember, "user is not member of product")
	}

	if member.Status != consts.ProductMemberStatusChangeRole {
		return errs.Newf(errs.RetWebMemberNotUnderApproval, "user is not in role change approval status")
	}

	var update horm.Map
	update["status"] = consts.ProductMemberStatusJoined
	update["change_role"] = 0

	if req.Status == consts.ApprovalAccess {
		update["role"] = member.ChangeRole
	}

	return table.UpdateProductMemberByID(ctx, member.Id, update)
}

// ProductMemberRemove 将指定用户移出产品
func ProductMemberRemove(ctx context.Context, userid uint64, req *pb.ProductMemberRemoveRequest) error {
	myRole, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return err
	}

	if myRole != consts.ProductRoleManager {
		return errs.New(errs.RetWebMemberNotManager, "not product manager")
	}

	_, member, err := table.GetProductMemberByUser(ctx, req.ProductID, req.Userid)
	if err != nil {
		return err
	}

	if GetProductRole(member) == consts.ProductRoleNotJoin {
		return errs.Newf(errs.RetWebIsNotMember, "user is already not member of the product")
	}

	update := horm.Map{
		"status":   consts.ProductMemberStatusQuit,
		"out_time": time.Now().Unix(),
	}

	return table.UpdateProductMemberByID(ctx, member.Id, update)
}

///////////////////////////////// function /////////////////////////////////////////

func GetUserProductRole(ctx context.Context, userid uint64, productID int) (int8, *table.TblProduct, error) {
	isNil, product, err := table.GetProductByID(ctx, productID)
	if err != nil {
		return 0, product, err
	}

	if isNil {
		return 0, product, errs.New(errs.RetWebNotFindProduct, "not find product")
	}

	_, member, err := table.GetProductMemberByUser(ctx, productID, userid)
	if err != nil {
		return 0, product, err
	}

	role := GetProductRole(member, product)
	if role == consts.ProductRoleNotJoin {
		return role, product, errs.New(errs.RetWebIsNotMember, "user is not member of product")
	}

	if role == consts.ProductRoleExpired {
		return role, product, errs.New(errs.RetWebMemberExpired, "product member permission has expired")
	}

	return role, product, nil
}

// GetProductRole 获取产品用户角色 0-非产品成员 1-管理员（仅当 product 不为空时判断）2-开发者 3-运营者 4-成员权限已过期
func GetProductRole(member *table.TblProductMember, product ...*table.TblProduct) int8 {
	if member == nil || member.Id == 0 {
		return consts.ProductRoleNotJoin
	}

	if member.Status != consts.ProductMemberStatusJoined &&
		member.Status != consts.ProductMemberStatusRenewal && // 续期审批，可以在产品权限正常使用情况下申请
		member.Status != consts.ProductMemberStatusChangeRole { // 角色变更申请肯定是在产品权限正常情况下申请的
		return consts.ProductRoleNotJoin
	}

	// 已过期
	if member.ExpireTime != 0 && int64(member.ExpireTime) < time.Now().Unix() {
		return consts.ProductRoleExpired
	}

	// 是否需要产品管理员判断
	if len(product) > 0 && product[0] != nil {
		if IsManager(member.UserID, product[0].Manager) {
			return consts.ProductRoleManager
		}
	}

	return member.Role
}

// GetProductRealRoleStatus 获取产品实际的角色和状态 role 和 status
// role 0:- 1:管理员 2:开发者 3:运营者
// status 0-未加入 1-待审批 2-续期审批 3-角色变更审批 4-已加入 5-审批拒绝 6-已退出 9-已过期
func GetProductRealRoleStatus(member *table.TblProductMember, product ...*table.TblProduct) (int8, int8) {
	if member == nil || member.Id == 0 {
		return consts.ProductRoleNotJoin, consts.ProductMemberStatusNotApply
	}

	role := member.Role

	// 是否需要产品管理员判断
	if len(product) > 0 && product[0] != nil {
		if IsManager(member.UserID, product[0].Manager) {
			role = consts.ProductRoleManager
		}
	}

	r := GetProductRole(member, product...)
	switch r {
	case consts.ProductRoleNotJoin:
		switch member.Status {
		case consts.ProductMemberStatusNotApply, consts.ProductMemberStatusQuit:
			return consts.ProductRoleNotJoin, consts.ProductMemberStatusNotApply
		default:
			return role, member.Status
		}

	case consts.ProductRoleExpired:
		switch member.Status {
		case consts.ProductMemberStatusRenewal:
			return role, consts.ProductMemberStatusRenewal
		default:
			return role, consts.ProductMemberStatusExpired
		}
	default:
		return r, member.Status
	}
}

func GetProductAndManagers(ctx context.Context, userid uint64,
	productID int) (int8, *table.TblProduct, []*table.TblProductMember, error) {
	isNil, product, err := table.GetProductByID(ctx, productID)
	if err != nil {
		return 0, nil, nil, err
	}

	if isNil {
		return 0, nil, nil, errs.New(errs.RetWebNotFindProduct, "not find product")
	}

	members, err := table.GetProductMemberByUsers(ctx, productID, GetUserIds(userid, product.Manager))
	if err != nil {
		return 0, nil, nil, err
	}

	var myRole int8
	var managers = []*table.TblProductMember{}

	if len(members) == 0 {
		return consts.ProductRoleNotJoin, product, managers, nil
	}

	for _, member := range members {
		role := GetProductRole(member, product)
		if role == consts.ProductRoleManager {
			managers = append(managers, member)
		}

		if member.UserID == userid {
			myRole = role
		}
	}

	return myRole, product, managers, nil
}

func GetUseridFromProductMember(members []*table.TblProductMember) []uint64 {
	ret := []uint64{}
	for _, member := range members {
		ret = append(ret, member.UserID)
	}

	return ret
}
