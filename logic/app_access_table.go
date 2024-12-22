// Copyright (c) 2024 The horm-database Authors. All rights reserved.
// This file Author:  CaoHao <18500482693@163.com> .
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

func TableSupportOps(ctx context.Context, userid uint64, tableID int) (*pb.SupportOpsResponse, error) {
	_, db, err := GetTableAndDBByTableID(ctx, tableID)
	if err != nil {
		return nil, err
	}

	var ret = pb.SupportOpsResponse{
		DBType:     db.Type,
		SupportOps: []string{},
	}

	if db.Type != consts.DBTypeRedis {
		ret.SupportOps = []string{consts.OpInsert, consts.OpReplace,
			consts.OpUpdate, consts.OpDelete, consts.OpFind, consts.OpFindAll}
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

func AppCanAccessTable(ctx context.Context, userid uint64,
	req *pb.AppCanAccessTableRequest) (*pb.AppCanAccessTableResponse, error) {
	tableInfo, _, err := GetTableAndDBByTableID(ctx, req.TableID)
	if err != nil {
		return nil, err
	}

	var ret = pb.AppCanAccessTableResponse{Apps: []*pb.AppCanAccessTable{}}

	apps, err := table.GetMyAppListByKeyword(ctx, userid, req.Keyword, consts.StatusOnline)
	if err != nil {
		return nil, err
	}

	if len(apps) == 0 {
		return &ret, nil
	}

	accessTables, err := table.GetAppAccessTables(ctx, GetAppidFromApps(apps), req.TableID)
	if err != nil {
		return nil, err
	}

	accessDBs, err := table.GetAppAccessDBs(ctx, GetAppidFromApps(apps), tableInfo.DB)
	if err != nil {
		return nil, err
	}

	// 未接入
	var isAppends = map[uint64]bool{}
	for _, app := range apps {
		tmp := pb.AppCanAccessTable{
			Appid:   app.Appid,
			AppName: app.Name,
			Intro:   app.Intro,
		}

		accessTable := GetAccessTableByAppidTableId(accessTables, app.Appid, req.TableID)
		accessDB := GetAccessDBByAppidDBId(accessDBs, app.Appid, tableInfo.DB)

		// 拥有表所属仓库数据权限，无需申请表权限
		if accessDB != nil && accessDB.Status == sc.AuthStatusNormal &&
			(accessDB.Root == sc.DBRootAll || accessDB.Root == sc.DBRootTableData) {
			continue
		}

		if accessTable == nil {
			tmp.AccessStatus = 0
			ret.Apps = append(ret.Apps, &tmp)
			isAppends[app.Appid] = true
		} else if accessTable.Status == sc.AuthStatusOffline ||
			accessTable.Status == sc.AuthStatusCancel || accessTable.Status == sc.AuthStatusReject {
			tmp.AccessStatus = accessTable.Status
			tmp.AccessInfo = &pb.TableAccessInfo{
				AccessID:       accessTable.Id,
				AccessQueryAll: accessTable.QueryAll,
				AccessOp:       strings.Split(accessTable.Op, ","),
				AccessReason:   accessTable.Reason,
			}
			ret.Apps = append(ret.Apps, &tmp)
			isAppends[app.Appid] = true
		}
	}

	// 已经拥有所属仓库权限，可接入/也可以不接入
	for _, app := range apps {
		isAppend, _ := isAppends[app.Appid]
		if isAppend {
			continue
		}

		tmp := pb.AppCanAccessTable{
			Appid:   app.Appid,
			AppName: app.Name,
			Intro:   app.Intro,
		}

		accessDB := GetAccessDBByAppidDBId(accessDBs, app.Appid, tableInfo.DB)

		if accessDB != nil && accessDB.Status == sc.AuthStatusNormal &&
			(accessDB.Root == sc.DBRootAll || accessDB.Root == sc.DBRootTableData) {
			tmp.AccessStatus = 11
			ret.Apps = append(ret.Apps, &tmp)
			isAppends[app.Appid] = true
		}
	}

	// 已接入
	for _, app := range apps {
		isAppend, _ := isAppends[app.Appid]
		if isAppend {
			continue
		}

		tmp := pb.AppCanAccessTable{
			Appid:   app.Appid,
			AppName: app.Name,
			Intro:   app.Intro,
		}

		accessTable := GetAccessTableByAppidTableId(accessTables, app.Appid, req.TableID)

		tmp.AccessStatus = accessTable.Status
		tmp.AccessInfo = &pb.TableAccessInfo{
			AccessID:       accessTable.Id,
			AccessQueryAll: accessTable.QueryAll,
			AccessOp:       strings.Split(accessTable.Op, ","),
			AccessReason:   accessTable.Reason,
		}
		ret.Apps = append(ret.Apps, &tmp)
	}

	return &ret, nil
}

func AppApplyAccessTable(ctx context.Context, userid uint64,
	req *pb.AppApplyAccessTableRequest) (*pb.AppApplyAccessResponse, error) {
	isNil, _, err := table.GetTableByID(ctx, req.TableID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.New(errs.RetWebNotFindTable, "not find table")
	}

	_, err = IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return nil, err
	}

	isNil, accessTable, err := table.GetAppAccessTable(ctx, req.Appid, req.TableID)
	if err != nil {
		return nil, err
	}

	var ret = pb.AppApplyAccessResponse{}

	if isNil {
		data := st.TblAccessTable{
			Appid:     req.Appid,
			TableId:   req.TableID,
			QueryAll:  req.QueryAll,
			Op:        strings.Join(req.Op, ","),
			Status:    sc.AuthStatusChecking,
			ApplyUser: userid,
			Reason:    req.Reason,
		}

		ret.AccessID, err = table.InsertAccessTable(ctx, &data)
		if err != nil {
			return nil, err
		}
	} else {
		ret.AccessID = accessTable.Id
		if accessTable.Status == sc.AuthStatusNormal {
			return nil, errs.New(errs.RetWebAccessStatusNormal, "app already has access permission of table")
		} else if accessTable.Status == sc.AuthStatusChecking {
			return nil, errs.New(errs.RetWebAccessStatusChecking,
				"the application's access to the database is under review")
		}

		update := horm.Map{
			"query_all":  req.QueryAll,
			"op":         strings.Join(req.Op, ","),
			"status":     sc.AuthStatusChecking,
			"apply_user": userid,
			"reason":     req.Reason,
		}

		err = table.UpdateAccessTableByID(ctx, accessTable.Id, update)
		if err != nil {
			return nil, err
		}
	}

	return &ret, nil
}

// AppAccessTableApproval 申请权限审批
func AppAccessTableApproval(ctx context.Context, userid uint64, req *pb.AppAccessTableApprovalRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	isNil, accessTable, err := table.GetAppAccessTable(ctx, req.Appid, req.TableID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access apply")
	}

	if accessTable.Status != sc.AuthStatusChecking {
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

	return table.UpdateAccessTableByID(ctx, accessTable.Id, update)
}

// AppAccessTableWithdraw 应用接入表数据撤销申请
func AppAccessTableWithdraw(ctx context.Context, userid uint64, req *pb.AppAccessTableWithdrawRequest) error {
	isNil, accessTable, err := table.GetAppAccessTable(ctx, req.Appid, req.TableID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access apply")
	}

	if accessTable.ApplyUser != userid {
		return errs.New(errs.RetWebAccessPermissionDeny, "is not my access, can`t withdraw")
	}

	if accessTable.Status != sc.AuthStatusChecking {
		return errs.New(errs.RetWebAccessStatusNotChecking,
			"the status of application access to table is not under review")
	}

	update := horm.Map{
		"status": sc.AuthStatusCancel,
	}

	return table.UpdateAccessTableByID(ctx, accessTable.Id, update)
}

// AppAccessTableUpdate 编辑表数据访问权限
func AppAccessTableUpdate(ctx context.Context, userid uint64, req *pb.AppAccessTableUpdateRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	isNil, accessTable, err := table.GetAppAccessTable(ctx, req.Appid, req.TableID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access info")
	}

	update := horm.Map{
		"query_all": req.QueryAll,
		"op":        strings.Join(req.Op, ","),
	}

	return table.UpdateAccessTableByID(ctx, accessTable.Id, update)
}

// AppAccessTableOnOff 仓库访问权限上/下线
func AppAccessTableOnOff(ctx context.Context, userid uint64, req *pb.AppAccessTableOnOffRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	isNil, accessTable, err := table.GetAppAccessTable(ctx, req.Appid, req.TableID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.New(errs.RetWebNotFindAccessInfo, "not find access info")
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateAccessTableByID(ctx, accessTable.Id, update)
}

func TablesAllAppAccessList(ctx context.Context, userid uint64,
	req *pb.TablesAllAppAccessListRequest) (ret *pb.TablesAllAppAccessListResponse, err error) {
	_, _, err = IsTableManager(ctx, userid, req.TableID)
	if err == nil { // 作为表管理员，返回访问该表的所有应用
		pageInfo, accessTables, err := table.GetAppAccessTableListByTableID(ctx, req.TableID, req.Page, req.Size)
		if err != nil {
			return nil, err
		}

		ret = &pb.TablesAllAppAccessListResponse{
			Total:           pageInfo.Total,
			TotalPage:       pageInfo.TotalPage,
			Page:            req.Page,
			Size:            req.Size,
			IsManager:       true,
			AppAccessTables: []*pb.AppAccessTable{},
		}

		if len(accessTables) > 0 {
			var userIds, appids []uint64

			for _, v := range accessTables {
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

			for _, v := range accessTables {
				ret.AppAccessTables = append(ret.AppAccessTables, &pb.AppAccessTable{
					Id:        v.Id,
					App:       GetAppBaseFromApp(userid, GetAppByAppid(apps, v.Appid), userMaps),
					Table:     nil,
					QueryAll:  v.QueryAll,
					Op:        strings.Split(v.Op, ","),
					Status:    v.Status,
					ApplyUser: userMaps[v.ApplyUser],
					Reason:    v.Reason,
					CreatedAt: v.CreatedAt.Unix(),
					UpdatedAt: v.UpdatedAt.Unix(),
				})
			}
		}
	} else { // 不是表管理员，查询自己应用访问该表
		apps, err := table.GetMyAppListByKeyword(ctx, userid, "", 0)
		if err != nil {
			return nil, err
		}

		if len(apps) == 0 {
			return &pb.TablesAllAppAccessListResponse{
				Page:            req.Page,
				Size:            req.Size,
				IsManager:       false,
				AppAccessTables: []*pb.AppAccessTable{},
			}, nil
		}

		pageInfo, accessTables, err := table.GetAppAccessTablesPages(ctx,
			GetAppidFromApps(apps), req.TableID, req.Page, req.Size)
		if err != nil {
			return nil, err
		}

		ret = &pb.TablesAllAppAccessListResponse{
			Total:           pageInfo.Total,
			TotalPage:       pageInfo.TotalPage,
			Page:            req.Page,
			Size:            req.Size,
			IsManager:       false,
			AppAccessTables: []*pb.AppAccessTable{},
		}

		if len(accessTables) > 0 {
			var userIds, appids []uint64

			for _, v := range accessTables {
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

			for _, v := range accessTables {
				ret.AppAccessTables = append(ret.AppAccessTables, &pb.AppAccessTable{
					Id:        v.Id,
					App:       GetAppBaseFromApp(userid, GetAppByAppid(apps, v.Appid), userMaps),
					Table:     nil,
					QueryAll:  v.QueryAll,
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

func AppsAllTableAccessList(ctx context.Context, userid uint64,
	req *pb.AppsAllTableAccessListRequest) (*pb.AppsAllTableAccessListResponse, error) {
	_, err := IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return nil, err
	}

	pageInfo, accessTables, err := table.GetAppAccessTableListByAppid(ctx, req.Appid, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.AppsAllTableAccessListResponse{
		Total:           pageInfo.Total,
		TotalPage:       pageInfo.TotalPage,
		Page:            req.Page,
		Size:            req.Size,
		AppAccessTables: []*pb.AppAccessTable{},
	}

	if len(accessTables) > 0 {
		var userIds []uint64
		var tableIds []int

		for _, v := range accessTables {
			tableIds = append(tableIds, v.TableId)
			userIds = append(userIds, v.ApplyUser)
		}

		tables, err := table.GetTableByIds(ctx, tableIds)
		if err != nil {
			return nil, err
		}

		userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
		if err != nil {
			return nil, err
		}

		for _, v := range accessTables {

			ret.AppAccessTables = append(ret.AppAccessTables, &pb.AppAccessTable{
				Id:        v.Id,
				App:       nil,
				Table:     GetTableBase(GetTableByID(tables, v.TableId)),
				QueryAll:  v.QueryAll,
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

func GetAccessTableByAppidTableId(accessTables []*st.TblAccessTable, appid uint64, tableID int) *st.TblAccessTable {
	if len(accessTables) == 0 {
		return nil
	}

	for _, v := range accessTables {
		if v.Appid == appid && v.TableId == tableID {
			return v
		}
	}

	return nil
}
