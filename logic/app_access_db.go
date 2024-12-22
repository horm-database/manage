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
	"strings"

	"github.com/horm-database/common/consts"
	"github.com/horm-database/common/errs"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	mc "github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	sc "github.com/horm-database/server/consts"
	st "github.com/horm-database/server/model/table"
)

func DBSupportOps(ctx context.Context, userid uint64, dbID int) (*pb.SupportOpsResponse, error) {
	isNil, db, err := table.GetDBByID(ctx, dbID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindDB, "not find db")
	}

	var ret = pb.SupportOpsResponse{
		DBType:     db.Type,
		SupportOps: []string{},
	}

	if db.Type != consts.DBTypeRedis {
		ret.SupportOps = []string{consts.OpInsert, consts.OpReplace, consts.OpUpdate, consts.OpDelete,
			consts.OpFind, consts.OpFindAll, consts.OpCreate, consts.OpDrop}
	} else {
		ret.SupportOps = []string{consts.OpExpire, consts.OpTTL, consts.OpExists, consts.OpDel, consts.OpSet,
			consts.OpSetEx, consts.OpSetNX, consts.OpMSet, consts.OpGet, consts.OpMGet, consts.OpGetSet, consts.OpIncr,
			consts.OpDecr, consts.OpIncrBy, consts.OpSetBit, consts.OpGetBit, consts.OpBitCount, consts.OpHSet,
			consts.OpHSetNx, consts.OpHmSet, consts.OpHIncrBy, consts.OpHIncrByFloat, consts.OpHDel, consts.OpHGet,
			consts.OpHMGet, consts.OpHGetAll, consts.OpHKeys, consts.OpHVals, consts.OpHExists, consts.OpHLen,
			consts.OpHStrLen, consts.OpLPush, consts.OpRPush, consts.OpLPop, consts.OpRPop, consts.OpLLen,
			consts.OpSAdd, consts.OpSMove, consts.OpSPop, consts.OpSRem, consts.OpSCard, consts.OpSMembers,
			consts.OpSIsMember, consts.OpSRandMember, consts.OpZAdd, consts.OpZRem, consts.OpZRemRangeByScore,
			consts.OpZRemRangeByRank, consts.OpZIncrBy, consts.OpZPopMin, consts.OpZPopMax, consts.OpZCard,
			consts.OpZScore, consts.OpZRank, consts.OpZRevRank, consts.OpZCount, consts.OpZRange,
			consts.OpZRangeByScore, consts.OpZRevRange, consts.OpZRevRangeByScore}
	}

	return &ret, nil
}

func AppCanAccessDB(ctx context.Context, userid uint64,
	req *pb.AppCanAccessDBRequest) (*pb.AppCanAccessDBResponse, error) {
	isNil, _, err := table.GetDBByID(ctx, req.DbID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindDB, "not find db")
	}

	var ret = pb.AppCanAccessDBResponse{Apps: []*pb.AppCanAccessDB{}}

	apps, err := table.GetMyAppListByKeyword(ctx, userid, req.Keyword, consts.StatusOnline)
	if err != nil {
		return nil, err
	}

	if len(apps) == 0 {
		return &ret, nil
	}

	accessDBs, err := table.GetAppAccessDBs(ctx, GetAppidFromApps(apps), req.DbID)
	if err != nil {
		return nil, err
	}

	// 未接入
	var isAppends = map[uint64]bool{}
	for _, app := range apps {
		tmp := pb.AppCanAccessDB{
			Appid:   app.Appid,
			AppName: app.Name,
			Intro:   app.Intro,
		}

		accessDB := GetAccessDBByAppidDBId(accessDBs, app.Appid, req.DbID)
		if accessDB == nil {
			tmp.AccessStatus = 0
			ret.Apps = append(ret.Apps, &tmp)
			isAppends[app.Appid] = true
		} else if accessDB.Status == sc.AuthStatusOffline ||
			accessDB.Status == sc.AuthStatusCancel || accessDB.Status == sc.AuthStatusReject {
			tmp.AccessStatus = accessDB.Status
			tmp.AccessInfo = &pb.DBAccessInfo{
				AccessID:     accessDB.Id,
				AccessRoot:   accessDB.Root,
				AccessOp:     strings.Split(accessDB.Op, ","),
				AccessReason: accessDB.Reason,
			}
		}
	}

	// 已接入
	for _, app := range apps {
		isAppend, _ := isAppends[app.Appid]
		if isAppend {
			continue
		}

		tmp := pb.AppCanAccessDB{
			Appid:   app.Appid,
			AppName: app.Name,
			Intro:   app.Intro,
		}

		accessDB := GetAccessDBByAppidDBId(accessDBs, app.Appid, req.DbID)
		tmp.AccessStatus = accessDB.Status
		tmp.AccessInfo = &pb.DBAccessInfo{
			AccessID:     accessDB.Id,
			AccessRoot:   accessDB.Root,
			AccessOp:     strings.Split(accessDB.Op, ","),
			AccessReason: accessDB.Reason,
		}
		ret.Apps = append(ret.Apps, &tmp)
	}

	return &ret, nil
}

func AppApplyAccessDB(ctx context.Context, userid uint64,
	req *pb.AppApplyAccessDBRequest) (*pb.AppApplyAccessResponse, error) {
	isNil, _, err := table.GetDBByID(ctx, req.DbID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindDB, "not find db")
	}

	_, err = IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return nil, err
	}

	isNil, accessDB, err := table.GetAppAccessDB(ctx, req.Appid, req.DbID)
	if err != nil {
		return nil, err
	}

	var ret = pb.AppApplyAccessResponse{}

	if isNil {
		data := st.TblAccessDB{
			Appid:     req.Appid,
			DB:        req.DbID,
			Root:      req.Root,
			Op:        strings.Join(req.Op, ","),
			ApplyUser: userid,
			Status:    sc.AuthStatusChecking,
			Reason:    req.Reason,
		}
		ret.AccessID, err = table.InsertAccessDB(ctx, &data)
		if err != nil {
			return nil, err
		}
	} else {
		ret.AccessID = accessDB.Id
		if accessDB.Status == sc.AuthStatusNormal {
			return nil, errs.New(errs.RetWebAccessStatusNormal, "app already has access permission of db")
		} else if accessDB.Status == sc.AuthStatusChecking {
			return nil, errs.New(errs.RetWebAccessStatusChecking,
				"the application's access to the database is under review")
		}

		update := horm.Map{
			"root":       req.Root,
			"op":         strings.Join(req.Op, ","),
			"status":     sc.AuthStatusChecking,
			"apply_user": userid,
			"reason":     req.Reason,
		}

		err = table.UpdateAccessDBByID(ctx, accessDB.Id, update)
		if err != nil {
			return nil, err
		}
	}

	return &ret, nil
}

// AppAccessDBApproval 申请权限审批
func AppAccessDBApproval(ctx context.Context, userid uint64, req *pb.AppAccessDBApprovalRequest) error {
	_, err := IsDBManager(ctx, userid, req.DbID)
	if err != nil {
		return err
	}

	isNil, accessDB, err := table.GetAppAccessDB(ctx, req.Appid, req.DbID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access apply")
	}

	if accessDB.Status != sc.AuthStatusChecking {
		return errs.New(errs.RetWebAccessStatusNotChecking,
			"the status of application access to database is not under review")
	}

	var update horm.Map

	if req.Status == mc.ApprovalAccess {
		update = horm.Map{
			"status": sc.AuthStatusNormal,
		}
	} else {
		update = horm.Map{
			"status": sc.AuthStatusReject,
		}
	}

	return table.UpdateAccessDBByID(ctx, accessDB.Id, update)
}

// AppAccessDBWithdraw 应用接入仓库撤销申请
func AppAccessDBWithdraw(ctx context.Context, userid uint64, req *pb.AppAccessDBWithdrawRequest) error {
	isNil, accessDB, err := table.GetAppAccessDB(ctx, req.Appid, req.DbID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access apply")
	}

	if accessDB.ApplyUser != userid {
		return errs.New(errs.RetWebAccessPermissionDeny, "is not my access, can`t withdraw")
	}

	if accessDB.Status != sc.AuthStatusChecking {
		return errs.New(errs.RetWebAccessStatusNotChecking,
			"the status of application access to database is not under review")
	}

	update := horm.Map{
		"status": sc.AuthStatusCancel,
	}

	return table.UpdateAccessDBByID(ctx, accessDB.Id, update)
}

// AppAccessDBUpdate 编辑仓库访问权限
func AppAccessDBUpdate(ctx context.Context, userid uint64, req *pb.AppAccessDBUpdateRequest) error {
	_, err := IsDBManager(ctx, userid, req.DbID)
	if err != nil {
		return err
	}

	isNil, accessDB, err := table.GetAppAccessDB(ctx, req.Appid, req.DbID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access info")
	}

	update := horm.Map{
		"root": req.Root,
		"op":   strings.Join(req.Op, ","),
	}

	return table.UpdateAccessDBByID(ctx, accessDB.Id, update)
}

// AppAccessDBOnOff 仓库访问权限上/下线
func AppAccessDBOnOff(ctx context.Context, userid uint64, req *pb.AppAccessDBOnOffRequest) error {
	_, err := IsDBManager(ctx, userid, req.DbID)
	if err != nil {
		return err
	}

	isNil, accessDB, err := table.GetAppAccessDB(ctx, req.Appid, req.DbID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access info")
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateAccessDBByID(ctx, accessDB.Id, update)
}

func DBsAllAppAccessList(ctx context.Context, userid uint64,
	req *pb.DBsAllAppAccessListRequest) (ret *pb.DBsAllAppAccessListResponse, err error) {
	_, err = IsDBManager(ctx, userid, req.DbID)
	if err == nil { // 作为仓库管理员，返回访问该仓库的所有应用
		pageInfo, accessDBs, err := table.GetAppAccessDBListByDBID(ctx, req.DbID, req.Page, req.Size)
		if err != nil {
			return nil, err
		}

		ret = &pb.DBsAllAppAccessListResponse{
			Total:        pageInfo.Total,
			TotalPage:    pageInfo.TotalPage,
			Page:         req.Page,
			Size:         req.Size,
			IsDBManager:  true,
			AppAccessDBs: []*pb.AppAccessDB{},
		}

		if len(accessDBs) > 0 {
			var userIds, appids []uint64

			for _, v := range accessDBs {
				appids = append(appids, v.Appid)
				userIds = append(userIds, v.ApplyUser)
			}

			apps, err := table.GetAppListByAppids(ctx, appids)
			if err != nil {
				return nil, err
			}

			for _, app := range apps {
				userIds = append(userIds, GetUserIds(app.Creator, app.Manager)...)
			}

			userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
			if err != nil {
				return nil, err
			}

			for _, v := range accessDBs {
				ret.AppAccessDBs = append(ret.AppAccessDBs, &pb.AppAccessDB{
					Id:        v.Id,
					App:       GetAppBaseFromApp(userid, GetAppByAppid(apps, v.Appid), userMaps),
					DB:        nil,
					Root:      v.Root,
					Op:        strings.Split(v.Op, ","),
					Status:    v.Status,
					ApplyUser: userMaps[v.ApplyUser],
					Reason:    v.Reason,
					CreatedAt: v.CreatedAt.Unix(),
					UpdatedAt: v.UpdatedAt.Unix(),
				})
			}
		}
	} else { // 不是仓库管理员，查询自己应用访问该仓库
		apps, err := table.GetMyAppListByKeyword(ctx, userid, "", 0)
		if err != nil {
			return nil, err
		}

		if len(apps) == 0 {
			return &pb.DBsAllAppAccessListResponse{
				Page:         req.Page,
				Size:         req.Size,
				IsDBManager:  false,
				AppAccessDBs: []*pb.AppAccessDB{},
			}, nil
		}

		pageInfo, accessDBs, err := table.GetAppAccessDBsPages(ctx,
			GetAppidFromApps(apps), req.DbID, req.Page, req.Size)
		if err != nil {
			return nil, err
		}

		ret = &pb.DBsAllAppAccessListResponse{
			Total:        pageInfo.Total,
			TotalPage:    pageInfo.TotalPage,
			Page:         req.Page,
			Size:         req.Size,
			IsDBManager:  false,
			AppAccessDBs: []*pb.AppAccessDB{},
		}

		if len(accessDBs) > 0 {
			var userIds, appids []uint64

			for _, v := range accessDBs {
				appids = append(appids, v.Appid)
				userIds = append(userIds, v.ApplyUser)
			}

			for _, app := range apps {
				userIds = append(userIds, GetUserIds(app.Creator, app.Manager)...)
			}

			userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
			if err != nil {
				return nil, err
			}

			for _, v := range accessDBs {
				ret.AppAccessDBs = append(ret.AppAccessDBs, &pb.AppAccessDB{
					Id:        v.Id,
					App:       GetAppBaseFromApp(userid, GetAppByAppid(apps, v.Appid), userMaps),
					DB:        nil,
					Root:      v.Root,
					Op:        strings.Split(v.Op, ","),
					Status:    v.Status,
					ApplyUser: userMaps[v.ApplyUser],
					Reason:    v.Reason,
					CreatedAt: v.CreatedAt.Unix(),
					UpdatedAt: v.UpdatedAt.Unix(),
				})
			}
		}
	}

	return ret, nil
}

func AppsAllDBAccessList(ctx context.Context, userid uint64,
	req *pb.AppsAllDBAccessListRequest) (*pb.AppsAllDBAccessListResponse, error) {
	_, err := IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return nil, err
	}

	pageInfo, accessDBs, err := table.GetAppAccessDBListByAppid(ctx, req.Appid, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.AppsAllDBAccessListResponse{
		Total:        pageInfo.Total,
		TotalPage:    pageInfo.TotalPage,
		Page:         req.Page,
		Size:         req.Size,
		AppAccessDBs: []*pb.AppAccessDB{},
	}

	if len(accessDBs) > 0 {
		var userIds []uint64
		var dbids []int

		for _, v := range accessDBs {
			dbids = append(dbids, v.DB)
			userIds = append(userIds, v.ApplyUser)
		}

		dbs, err := table.GetDBByIds(ctx, dbids)
		if err != nil {
			return nil, err
		}

		userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
		if err != nil {
			return nil, err
		}

		for _, v := range accessDBs {

			ret.AppAccessDBs = append(ret.AppAccessDBs, &pb.AppAccessDB{
				Id:        v.Id,
				App:       nil,
				DB:        GetDBBase(GetDBByID(dbs, v.DB)),
				Root:      v.Root,
				Op:        strings.Split(v.Op, ","),
				Status:    v.Status,
				ApplyUser: userMaps[v.ApplyUser],
				Reason:    v.Reason,
				CreatedAt: v.CreatedAt.Unix(),
				UpdatedAt: v.UpdatedAt.Unix(),
			})
		}
	}

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func GetAccessDBByAppidDBId(accessDBs []*st.TblAccessDB, appid uint64, dbId int) *st.TblAccessDB {
	if len(accessDBs) == 0 {
		return nil
	}

	for _, v := range accessDBs {
		if v.Appid == appid && v.DB == dbId {
			return v
		}
	}

	return nil
}
