// Copyright (c) 2024 The horm-database Authors (such as CaoHao <18500482693@163.com>). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package logic

import (
	"context"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/proto"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	tb "github.com/horm-database/server/model/table"
	"github.com/samber/lo"
)

// WorkspaceBaseInfo 工作空间基础信息
func WorkspaceBaseInfo(ctx context.Context, userid uint64,
	req *pb.WorkspaceBaseInfoRequest) (*pb.WorkspaceBaseInfoResponse, error) {
	tblWorkspace, err := table.GetWorkspace(ctx, req.Workspace)
	if err != nil {
		return nil, err
	}

	if tblWorkspace == nil || tblWorkspace.Id == 0 {
		return nil, errs.New(errs.RetWebWorkspaceNotExists, "workspace not exists")
	}

	managerUids := GetUserIds(tblWorkspace.Manager)
	userMaps, err := table.GetUserBasesMapByIds(ctx, GetUserIds(tblWorkspace.Creator, managerUids))
	if err != nil {
		return nil, err
	}

	ret := pb.WorkspaceBaseInfoResponse{
		Id:         tblWorkspace.Id,
		Workspace:  tblWorkspace.Workspace,
		Name:       tblWorkspace.Name,
		Intro:      tblWorkspace.Intro,
		Company:    tblWorkspace.Company,
		Department: tblWorkspace.Department,
		Creator:    userMaps[tblWorkspace.Creator],
		Manager:    GetUsersFromMap(userMaps, managerUids),
		CreateTime: tblWorkspace.CreatedAt.Unix(),
	}

	ret.Role = consts.WorkspaceMemberNotJoin
	ret.Status = consts.WorkspaceMemberStatusNotApply
	ret.ExpireTime = 0
	ret.OutTime = 0

	if userid != 0 {
		_, member, err := table.GetWorkspaceMemberByUser(ctx, tblWorkspace.Id, userid)
		if err != nil {
			return nil, err
		}

		ret.Role = GetWorkspaceRole(member, tblWorkspace)
		ret.Status = member.Status
		ret.ExpireTime = member.ExpireTime
		ret.OutTime = member.OutTime
	}

	return &ret, nil
}

// WorkspaceJoinApply 申请加入空间 / 续期
func WorkspaceJoinApply(ctx context.Context, userid uint64, workspaceID int, req *pb.WorkspaceJoinApplyRequest) error {
	isNil, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, userid)
	if err != nil {
		return err
	}

	if isNil { // 新成员申请加入空间
		newMember := table.TblWorkspaceMember{
			WorkspaceID: workspaceID,
			UserID:      userid,
			Status:      consts.WorkspaceMemberStatusApproval,
			JoinTime:    time.Now().Unix(),
			ExpireType:  req.ExpireType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		return table.InsertWorkspaceMember(ctx, &newMember)
	} else if GetWorkspaceRole(member) == consts.WorkspaceMemberNotJoin { // 重新申请加入空间
		if member.Status == consts.WorkspaceMemberStatusApproval ||
			member.Status == consts.WorkspaceMemberStatusRenewal {
			return errs.Newf(errs.RetWebMemberUnderApproval, "under approval, please do not apply repeatedly")
		}

		replace := horm.Map{
			"id":           member.Id,
			"workspace_id": member.WorkspaceID,
			"userid":       member.UserID,
			"status":       consts.WorkspaceMemberStatusApproval,
			"join_time":    time.Now().Unix(),
			"expire_type":  req.ExpireType,
			"expire_time":  0,
			"out_time":     0,
			"updated_at":   time.Now(),
		}
		return table.ReplaceWorkspaceMember(ctx, replace)
	} else { // 申请续期
		if member.ExpireTime == 0 || int64(member.ExpireTime)-time.Now().Unix() > 7*86400 { // 只有7天内过期的用户才允许续期
			return errs.Newf(errs.RetWebIsMember, "user is already member of workspace")
		} else {
			if member.Status == consts.WorkspaceMemberStatusApproval ||
				member.Status == consts.WorkspaceMemberStatusRenewal {
				return errs.Newf(errs.RetWebMemberUnderApproval, "under approval, please do not apply repeatedly")
			}

			update := horm.Map{
				"status":      consts.WorkspaceMemberStatusRenewal,
				"expire_type": req.ExpireType,
				"out_time":    0,
			}
			return table.UpdateWorkspaceMemberByID(ctx, member.Id, update)
		}
	}
}

// WorkspaceApproval 空间权限审批
func WorkspaceApproval(ctx context.Context, userid uint64, workspaceID int, req *pb.WorkspaceApprovalRequest) error {
	myRole, _, err := GetUserWorkspaceRole(ctx, userid, workspaceID)
	if err != nil {
		return err
	}

	if myRole != consts.WorkspaceMemberManager {
		return errs.New(errs.RetWebMemberNotManager, "not workspace manager")
	}

	isNil, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, req.Userid)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebIsNotApply, "user has not applied for workspace permissions")
	}

	if member.Status != consts.WorkspaceMemberStatusApproval && member.Status != consts.WorkspaceMemberStatusRenewal {
		return errs.Newf(errs.RetWebMemberNotUnderApproval, "user is not in approval status")
	}

	var update horm.Map
	if req.Status == consts.ApprovalAccess {
		update = horm.Map{
			"status":      consts.WorkspaceMemberStatusJoined,
			"expire_time": GetExpireTime(int64(member.ExpireTime), member.ExpireType),
		}

		if member.Status == consts.WorkspaceMemberStatusApproval {
			update["join_time"] = time.Now().Unix()
		}
	} else {
		update = horm.Map{
			"status": consts.WorkspaceMemberStatusReject,
		}
	}

	return table.UpdateWorkspaceMemberByID(ctx, member.Id, update)
}

// WorkspaceMemberInvite 管理员邀请用户加入空间
func WorkspaceMemberInvite(ctx context.Context, userid uint64,
	workspaceID int, req *pb.WorkspaceMemberInviteRequest) error {
	myRole, _, err := GetUserWorkspaceRole(ctx, userid, workspaceID)
	if err != nil {
		return err
	}

	if myRole != consts.WorkspaceMemberManager {
		return errs.New(errs.RetWebMemberNotManager, "not workspace manager")
	}

	isNil, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, req.Userid)
	if err != nil {
		return err
	}

	if isNil { // 邀请新成员加入空间
		newMember := table.TblWorkspaceMember{
			WorkspaceID: workspaceID,
			UserID:      req.Userid,
			Status:      consts.WorkspaceMemberStatusJoined,
			JoinTime:    time.Now().Unix(),
			ExpireType:  req.ExpireType,
			ExpireTime:  int(GetExpireTime(0, req.ExpireType)),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		return table.InsertWorkspaceMember(ctx, &newMember)
	} else if GetWorkspaceRole(member) == consts.WorkspaceMemberNotJoin { // 邀请重新加入空间
		replace := horm.Map{
			"id":           member.Id,
			"workspace_id": member.WorkspaceID,
			"userid":       member.UserID,
			"status":       consts.WorkspaceMemberStatusJoined,
			"join_time":    time.Now().Unix(),
			"expire_type":  req.ExpireType,
			"expire_time":  GetExpireTime(0, req.ExpireType),
			"out_time":     0,
			"updated_at":   time.Now(),
		}
		return table.ReplaceWorkspaceMember(ctx, replace)
	} else { // 邀请续期
		if member.ExpireTime == 0 || int64(member.ExpireTime)-time.Now().Unix() > 7*86400 { // 只有7天内过期的用户才允许续期
			return errs.Newf(errs.RetWebIsMember, "user is already member of workspace")
		} else {
			update := horm.Map{
				"status":      consts.WorkspaceMemberStatusJoined,
				"expire_type": req.ExpireType,
				"expire_time": GetExpireTime(int64(member.ExpireTime), req.ExpireType),
				"out_time":    0,
			}
			return table.UpdateWorkspaceMemberByID(ctx, member.Id, update)
		}
	}
}

// WorkspaceMemberRemove 将指定用户移出空间
func WorkspaceMemberRemove(ctx context.Context, userid uint64,
	workspaceID int, req *pb.WorkspaceMemberRemoveRequest) error {
	myRole, _, err := GetUserWorkspaceRole(ctx, userid, workspaceID)
	if err != nil {
		return err
	}

	if myRole != consts.WorkspaceMemberManager {
		return errs.New(errs.RetWebMemberNotManager, "not workspace manager")
	}

	_, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, req.Userid)
	if err != nil {
		return err
	}

	if GetWorkspaceRole(member) == consts.WorkspaceMemberNotJoin {
		return errs.Newf(errs.RetWebNotWorkspaceMember, "user is already not member of the workspace")
	}

	update := horm.Map{
		"status":   consts.WorkspaceMemberStatusQuit,
		"out_time": time.Now().Unix(),
	}

	return table.UpdateWorkspaceMemberByID(ctx, member.Id, update)
}

// WorkspaceMemberList 工作空间成员列表
func WorkspaceMemberList(ctx context.Context, userid uint64, workspaceID int,
	req *pb.WorkspaceMemberListRequest) (*pb.WorkspaceMemberListResponse, error) {
	myRole, workspace, err := GetUserWorkspaceRole(ctx, userid, workspaceID)
	if err != nil {
		return nil, err
	}

	var pageRet *proto.Detail
	var members []*table.TblWorkspaceMember

	if myRole == consts.WorkspaceMemberManager {
		pageRet, members, err = table.GetWorkspaceMembersAll(ctx, workspaceID, req.Page, req.Size)
	} else {
		pageRet, members, err = table.GetWorkspaceMembersJoined(ctx, workspaceID, req.Page, req.Size)
	}

	if err != nil {
		return nil, err
	}

	var ret = pb.WorkspaceMemberListResponse{
		Total:     pageRet.Total,
		TotalPage: pageRet.TotalPage,
		Page:      req.Page,
		Size:      req.Size,
		IsManager: myRole == consts.WorkspaceMemberManager,
		Members:   []*pb.WorkspaceMember{},
	}

	userIds := GetUseridFromWorkspaceMember(members)
	userBases, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, v := range members {
		userBase := userBases[v.UserID]
		if userBase == nil {
			continue
		}

		member := pb.WorkspaceMember{
			MemberID:   v.Id,
			Userid:     v.UserID,
			Account:    userBase.Account,
			Nickname:   userBase.Nickname,
			JoinTime:   int(v.JoinTime),
			ExpireType: v.ExpireType,
			ExpireTime: v.ExpireTime,
			OutTime:    v.OutTime,
		}

		member.Role, member.Status = GetWorkspaceRealRoleStatus(v, workspace)
		ret.Members = append(ret.Members, &member)
	}

	return &ret, nil
}

func MaintainWorkspaceManager(ctx context.Context, userid uint64, workspaceID int, manager []uint64) error {
	myRole, _, err := GetUserWorkspaceRole(ctx, userid, workspaceID)
	if err != nil {
		return err
	}

	if myRole != consts.WorkspaceMemberManager {
		return errs.New(errs.RetWebMemberNotManager, "not workspace manager")
	}

	managerUids := lo.Uniq(manager)

	members, err := table.GetWorkspaceMemberByUsers(ctx, workspaceID, managerUids)
	if err != nil {
		return err
	}

	roleMap := map[uint64]int8{}
	for _, member := range members {
		roleMap[member.UserID] = GetWorkspaceRole(member)
	}

	for _, uid := range managerUids {
		role := roleMap[uid]
		if role == consts.ProductRoleNotJoin {
			return errs.Newf(errs.RetWebIsNotMember, "user [%d] is not member of workspace", uid)
		}
	}

	update := horm.Map{
		"manager": types.JoinUint64(managerUids, ","),
	}

	err = table.UpdateWorkspaceByID(ctx, workspaceID, update)
	if err != nil {
		return err
	}

	return nil
}

///////////////////////////////// function /////////////////////////////////////////

func GetUserWorkspaceRole(ctx context.Context, userid uint64, workspaceID int) (int8, *tb.TblWorkspace, error) {
	workspace, err := table.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return 0, workspace, err
	}

	_, member, err := table.GetWorkspaceMemberByUser(ctx, workspaceID, userid)
	if err != nil {
		return 0, workspace, err
	}

	role := GetWorkspaceRole(member, workspace)
	if role == consts.WorkspaceMemberNotJoin {
		return 0, workspace, errs.New(errs.RetWebIsNotMember, "user is not member of workspace")
	}

	if role == consts.WorkspaceMemberExpired {
		return 0, workspace, errs.New(errs.RetWebMemberExpired, "workspace member permission has expired")
	}

	return role, workspace, nil
}

// GetWorkspaceRole 获取空间用户角色 0-非空间成员 1-空间成员 2-空间管理员（仅当 workspace 不为空时判断） 3-权限已过期
func GetWorkspaceRole(member *table.TblWorkspaceMember, workspace ...*tb.TblWorkspace) int8 {
	if member == nil || member.Id == 0 {
		return consts.WorkspaceMemberNotJoin
	}

	if member.Status != consts.WorkspaceMemberStatusJoined &&
		member.Status != consts.WorkspaceMemberStatusRenewal { // 续期审批，可以在空间权限正常使用情况下申请
		return consts.WorkspaceMemberNotJoin
	}

	// 已过期
	if member.ExpireTime != 0 && int64(member.ExpireTime) < time.Now().Unix() {
		return consts.WorkspaceMemberExpired
	}

	// 是否需要空间管理员判断
	if len(workspace) > 0 && workspace[0] != nil {
		if IsManager(member.UserID, workspace[0].Manager) {
			return consts.WorkspaceMemberManager
		}
	}

	return consts.WorkspaceMember
}

// GetWorkspaceRealRoleStatus 获取空间实际的角色和状态 role 和 status
// role 0:- 1:普通成员 2:管理员
// status 1-待审批 2-续期审批 3-未加入 4-正常 5-审批拒绝  6-已退出 9-已过期
func GetWorkspaceRealRoleStatus(member *table.TblWorkspaceMember, workspace ...*tb.TblWorkspace) (int8, int8) {
	if member == nil || member.Id == 0 {
		return consts.WorkspaceMemberNotJoin, consts.WorkspaceMemberStatusNotApply
	}

	var role int8 = consts.WorkspaceMember

	// 是否需要产品管理员判断
	if len(workspace) > 0 && workspace[0] != nil {
		if IsManager(member.UserID, workspace[0].Manager) {
			role = consts.WorkspaceMemberManager
		}
	}

	r := GetWorkspaceRole(member, workspace...)

	switch r {
	case consts.WorkspaceMemberNotJoin:
		switch member.Status {
		case consts.WorkspaceMemberStatusNotApply, consts.WorkspaceMemberStatusQuit:
			return consts.WorkspaceMemberNotJoin, consts.WorkspaceMemberStatusNotApply
		default:
			return role, member.Status
		}

	case consts.ProductRoleExpired:
		switch member.Status {
		case consts.WorkspaceMemberStatusRenewal:
			return role, consts.WorkspaceMemberStatusRenewal
		default:
			return role, consts.ProductMemberStatusExpired
		}

	default:
		return r, member.Status
	}
}

func GetUseridFromWorkspaceMember(members []*table.TblWorkspaceMember) []uint64 {
	ret := []uint64{}
	for _, member := range members {
		ret = append(ret, member.UserID)
	}

	return ret
}
