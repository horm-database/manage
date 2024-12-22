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
	"fmt"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/orm/obj"
	"github.com/samber/lo"
)

// AddDB 新增数据库
func AddDB(ctx context.Context, userid uint64, req *pb.AddDBRequest) (*pb.AddDBResponse, error) {
	role, _, err := GetUserProductRole(ctx, userid, req.ProductID)
	if err != nil {
		return nil, err
	}

	if role != consts.ProductRoleManager && role != consts.ProductRoleDeveloper {
		return nil, errs.New(errs.RetWebCantCreateDB, "user id not product manager or developer, can`t create table")
	}

	data := obj.TblDB{
		Name:            req.Name,
		Intro:           req.Intro,
		Desc:            req.Desc,
		ProductID:       req.ProductID,
		Creator:         userid,
		Manager:         fmt.Sprint(userid),
		Status:          consts.StatusOnline,
		WriteTimeoutTmp: req.WriteTimeout,
		ReadTimeoutTmp:  req.ReadTimeout,
		WarnTimeoutTmp:  req.WarnTimeout,
		OmitErrorTmp:    req.OmitError,
		DebugTmp:        req.Debug,
		Type:            req.Type,
		Version:         req.Version,
		Network:         req.Network,
		Address:         req.Address,
		BakAddress:      req.BakAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id, err := table.AddDB(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &pb.AddDBResponse{ID: id}, nil
}

func UpdateDBBase(ctx context.Context, userid uint64, req *pb.UpdateDBBaseRequest) error {
	_, err := IsDBManager(ctx, userid, req.DBId)
	if err != nil {
		return err
	}

	update := horm.Map{
		"name":  req.Name,
		"intro": req.Intro,
		"desc":  req.Desc,
	}

	return table.UpdateDBByID(ctx, req.DBId, update)
}

func MaintainDBManager(ctx context.Context, userid uint64, req *pb.MaintainDBManagerRequest) error {
	db, err := IsDBManager(ctx, userid, req.DBId)
	if err != nil {
		return err
	}

	managerUids := lo.Uniq(req.Manager)

	members, err := table.GetProductMemberByUsers(ctx, db.ProductID, managerUids)
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

	return table.UpdateDBByID(ctx, req.DBId, update)
}

func UpdateDBStatus(ctx context.Context, userid uint64, req *pb.UpdateDBStatusRequest) error {
	_, err := IsDBManager(ctx, userid, req.DBId)
	if err != nil {
		return err
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateDBByID(ctx, req.DBId, update)
}

func UpdateDBNetwork(ctx context.Context, userid uint64, req *pb.UpdateDBNetworkRequest) error {
	_, err := IsDBManager(ctx, userid, req.DBId)
	if err != nil {
		return err
	}

	update := horm.Map{
		"type":          req.Type,
		"version":       req.Version,
		"network":       req.Network,
		"address":       req.Address,
		"bak_address":   req.BakAddress,
		"write_timeout": req.WriteTimeout,
		"read_timeout":  req.ReadTimeout,
		"warn_timeout":  req.WarnTimeout,
		"omit_error":    req.OmitError,
		"debug":         req.Debug,
	}

	return table.UpdateDBByID(ctx, req.DBId, update)
}

func DBBase(ctx context.Context, userid uint64, dbID int) (*pb.DBBaseResponse, error) {
	db, dbManagerUids, err := GetDBAndManagers(ctx, dbID)
	if err != nil {
		return nil, err
	}

	// 获取所属产品及其管理员（拥有 db 管理员权限）
	myRole, product, productManagers, err := GetProductAndManagers(ctx, userid, db.ProductID)
	if err != nil {
		return nil, err
	}

	productManagerUids := GetUseridFromProductMember(productManagers)
	userBases, err := table.GetUserBasesMapByIds(ctx, GetUserIds(db.Creator, dbManagerUids, productManagerUids))
	if err != nil {
		return nil, err
	}

	ret := pb.DBBaseResponse{
		Info:    GetDBBase(db),
		Creator: userBases[db.Creator],
		Manager: GetUsersFromMap(userBases, dbManagerUids),
		ProductInfo: &pb.ProductBase{
			Id:        product.Id,
			Name:      product.Name,
			Intro:     product.Intro,
			Status:    product.Status,
			CreatedAt: product.CreatedAt.Unix(),
		},
		ProductManager: GetUsersFromMap(userBases, productManagerUids),
		Tables:         []*pb.TableBase{},
	}

	if myRole == consts.ProductRoleManager || (myRole != consts.ProductRoleNotJoin &&
		myRole != consts.ProductRoleExpired && lo.IndexOf(dbManagerUids, userid) != -1) {
		ret.IsManager = true
	}

	tables, err := table.GetDBTables(ctx, dbID)
	if err != nil {
		return nil, err
	}

	for _, v := range tables {
		ret.Tables = append(ret.Tables, GetTableBase(v))
	}

	return &ret, nil
}

func DBNetworkDetail(ctx context.Context, userid uint64, dbID int) (*pb.DBNetworkInfoResponse, error) {
	db, err := IsDBManager(ctx, userid, dbID)
	if err != nil {
		return nil, err
	}

	ret := pb.DBNetworkInfoResponse{
		Type:         db.Type,
		Version:      db.Version,
		Network:      db.Network,
		Address:      db.Address,
		BakAddress:   db.BakAddress,
		WriteTimeout: db.WriteTimeoutTmp,
		ReadTimeout:  db.ReadTimeoutTmp,
		WarnTimeout:  db.WarnTimeoutTmp,
		OmitError:    db.OmitErrorTmp,
		Debug:        db.DebugTmp,
	}

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func GetDBAndManagers(ctx context.Context, dbID int) (*obj.TblDB, []uint64, error) {
	isNil, db, err := table.GetDBByID(ctx, dbID)
	if err != nil {
		return nil, nil, err
	}

	if isNil {
		return nil, nil, errs.New(errs.RetWebNotFindDB, "not find db")
	}

	return db, GetUserIds(db.Manager), nil
}

func IsDBManager(ctx context.Context, userid uint64, dbID int) (*obj.TblDB, error) {
	db, dbManagerUids, err := GetDBAndManagers(ctx, dbID)
	if err != nil {
		return nil, err
	}

	productRole, _, err := GetUserProductRole(ctx, userid, db.ProductID)
	if err != nil {
		return db, err
	}

	if productRole == consts.ProductRoleManager {
		return db, nil
	}

	// user is db manager
	if lo.IndexOf(dbManagerUids, userid) != -1 {
		return db, nil
	}

	return db, errs.New(errs.RetWebNotDBManager, "user is not manager of db")
}

func GetDBBase(db *obj.TblDB) *pb.DBBase {
	if db == nil {
		return nil
	}

	return &pb.DBBase{
		Id:         db.Id,
		Name:       db.Name,
		Intro:      db.Intro,
		Desc:       db.Desc,
		ProductID:  db.ProductID,
		Type:       db.Type,
		Version:    db.Version,
		Status:     db.Status,
		CreateTime: db.CreatedAt.Unix(),
	}
}

func GetDBByID(dbs []*obj.TblDB, id int) *obj.TblDB {
	if len(dbs) == 0 {
		return nil
	}

	for _, v := range dbs {
		if v.Id == id {
			return v
		}
	}

	return nil
}
