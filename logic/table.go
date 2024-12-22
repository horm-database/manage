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
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	cc "github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/orm/obj"
	"github.com/samber/lo"
)

// AddTable 新增表
func AddTable(ctx context.Context, userid uint64, req *pb.AddTableRequest) (*pb.AddTableResponse, error) {
	_, err := IsDBManager(ctx, userid, req.DB)
	if err != nil {
		return nil, err
	}

	data := obj.TblTable{
		Name:        req.Name,
		Intro:       req.Intro,
		Desc:        req.Desc,
		TableVerify: req.TableVerify,
		DB:          req.DB,
		Status:      cc.StatusOnline,
		Creator:     userid,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := table.AddTable(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &pb.AddTableResponse{ID: id}, nil
}

func UpdateTableBase(ctx context.Context, userid uint64, req *pb.UpdateTableBaseRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	update := horm.Map{
		"intro": req.Intro,
		"desc":  req.Desc,
	}

	return table.UpdateTableByID(ctx, req.TableID, update)
}

func UpdateTableStatus(ctx context.Context, userid uint64, req *pb.UpdateTableStatusRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateTableByID(ctx, req.TableID, update)
}

func UpdateTableAdvance(ctx context.Context, userid uint64, req *pb.UpdateTableAdvanceRequest) error {
	_, _, err := IsTableManager(ctx, userid, req.TableID)
	if err != nil {
		return err
	}

	update := horm.Map{
		"table_verify": req.TableVerify,
	}

	return table.UpdateTableByID(ctx, req.TableID, update)
}

func TableDetail(ctx context.Context, userid uint64, tableID int) (*pb.TableDetailResponse, error) {
	isNil, tableInfo, err := table.GetTableByID(ctx, tableID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.Newf(errs.RetWebNotFindTable, "not find table [%d]", tableID)
	}

	db, dbManagerUids, err := GetDBAndManagers(ctx, tableInfo.DB)
	if err != nil {
		return nil, err
	}

	// 获取所属产品及其管理员（拥有 db 管理员权限）
	myRole, product, productManagers, err := GetProductAndManagers(ctx, userid, db.ProductID)
	if err != nil {
		return nil, err
	}

	productManagerUids := GetUseridFromProductMember(productManagers)

	userBases, err := table.GetUserBasesMapByIds(ctx, GetUserIds(tableInfo.Creator, dbManagerUids, productManagerUids))
	if err != nil {
		return nil, err
	}

	ret := pb.TableDetailResponse{
		Info:      GetTableBase(tableInfo),
		DBInfo:    GetDBBase(db),
		Creator:   userBases[tableInfo.Creator],
		DBManager: GetUsersFromMap(userBases, dbManagerUids),
		ProductInfo: &pb.ProductBase{
			Id:        product.Id,
			Name:      product.Name,
			Intro:     product.Intro,
			Status:    product.Status,
			CreatedAt: product.CreatedAt.Unix(),
		},
		ProductManager: GetUsersFromMap(userBases, productManagerUids),
	}

	if myRole == cc.ProductRoleManager || (myRole != cc.ProductRoleNotJoin &&
		myRole != cc.ProductRoleExpired && lo.IndexOf(dbManagerUids, userid) != -1) {
		ret.IsManager = true
	}

	ret.TableFields = []*pb.TableField{
		{
			Field:     "id",
			Type:      3,
			Len:       "",
			Empty:     false,
			Status:    1,
			Default:   "",
			IsPrimary: true,
			IsIndex:   true,
			Comment:   "primary id",
			More:      "{}",
		},
		{
			Field:     "name",
			Type:      2,
			Len:       "64",
			Empty:     false,
			Status:    1,
			Default:   "",
			IsPrimary: false,
			IsIndex:   false,
			Comment:   "用户姓名",
			More:      "{}",
		},
		{
			Field:     "gender",
			Type:      16,
			Len:       "",
			Empty:     false,
			Status:    1,
			Default:   "1",
			IsPrimary: false,
			IsIndex:   false,
			Comment:   "用户性别 1-male 2-female",
			More:      `[{"id":1, "name":"male"},{"id":2, "name":"female"}]`,
		},
		{
			Field:     "readme",
			Type:      15,
			Len:       "64KB",
			Empty:     true,
			Status:    1,
			Default:   "1",
			IsPrimary: false,
			IsIndex:   false,
			Comment:   "自我介绍",
			More:      ``,
		},
		{
			Field:     "created_at",
			Type:      19,
			Len:       "",
			Empty:     false,
			Status:    1,
			Default:   "CURRENT_TIMESTAMP",
			IsPrimary: false,
			IsIndex:   false,
			Comment:   "记录创建时间",
			More:      ``,
		},
		{
			Field:     "updated_at",
			Type:      19,
			Len:       "",
			Empty:     false,
			Status:    1,
			Default:   "CURRENT_TIMESTAMP",
			IsPrimary: false,
			IsIndex:   false,
			Comment:   "记录最后修改时间",
			More:      ``,
		},
	}

	ret.TableIndexs = []*pb.TableIndex{
		{
			Name:   "PRIMARY",
			Type:   "PRIMARYKEY",
			Fields: "id",
		},
		{
			Name:   "name",
			Type:   "UNIQUEKEY",
			Fields: "name",
		},
	}

	ret.LangStructs = map[string]*pb.LangStruct{
		"go": {
			Language: "go",
			Struct:   "type Student struct {\n\tId        int       `orm:\"id,int,omitempty\" json:\"id,omitempty\"`           // primary id\n\tName      string    `orm:\"name,string,omitempty\" json:\"name,omitempty\"`    // 用户姓名\n\tGender    int8      `orm:\"gender,int8,omitempty\" json:\"gender,omitempty\"`  // 用户性别 1-male 2-female\n\tReadme    []byte    `orm:\"readme,bytes,omitempty\" json:\"readme,omitempty\"` // 自我介绍\n\tCreatedAt time.Time `orm:\"created_at,time,omitempty\" json:\"created_at\"`    // 记录创建时间\n\tUpdatedAt time.Time `orm:\"updated_at,time,omitempty\" json:\"updated_at\"`    // 记录最后修改时间\n}",
		},
	}

	return &ret, nil
}

func TableAdvanceConfig(ctx context.Context, userid uint64, tableID int) (*pb.TableAdvanceConfigResponse, error) {
	tableInfo, _, err := IsTableManager(ctx, userid, tableID)
	if err != nil {
		return nil, err
	}

	ret := pb.TableAdvanceConfigResponse{}
	ret.Definition = tableInfo.Definition
	ret.TableVerify = tableInfo.TableVerify

	ret.Definition = "CREATE TABLE `tbl_user` (\n  `id` bigint NOT NULL COMMENT '用户id',\n  `account` varchar(128) NOT NULL DEFAULT '' COMMENT '账号，可以是 admin，邮箱等。。。',\n  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',\n  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '密码',\n  `mobile` varchar(64) NOT NULL DEFAULT '' COMMENT '手机号',\n  `token` varchar(64) NOT NULL DEFAULT '' COMMENT 'token',\n  `avatar_url` varchar(512) NOT NULL DEFAULT '' COMMENT '头像',\n  `gender` tinyint NOT NULL DEFAULT '1' COMMENT '性别 1-男  2-女',\n  `company` varchar(128) NOT NULL DEFAULT '' COMMENT '公司',\n  `department` varchar(128) NOT NULL DEFAULT '' COMMENT '部门',\n  `city` varchar(128) NOT NULL DEFAULT '' COMMENT '城市',\n  `province` varchar(128) NOT NULL DEFAULT '' COMMENT '省份',\n  `country` varchar(128) NOT NULL DEFAULT '' COMMENT '国家',\n  `last_login_time` int NOT NULL DEFAULT '0' COMMENT '上次登录时间',\n  `last_login_ip` varchar(32) NOT NULL DEFAULT '' COMMENT '上次登录ip',\n  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',\n  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后修改时间',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `account` (`account`,`mobile`),\n  KEY `mobile` (`mobile`),\n  KEY `nickname` (`nickname`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号信息'"
	ret.TableVerify = "tbl_user_*"
	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func IsTableManager(ctx context.Context, userid uint64, tableID int) (*obj.TblTable, *obj.TblDB, error) {
	isNil, tableInfo, err := table.GetTableByID(ctx, tableID)
	if err != nil {
		return nil, nil, err
	}

	if isNil {
		return tableInfo, nil, errs.Newf(errs.RetWebNotFindTable, "not find table [%d]", tableID)
	}

	db, err := IsDBManager(ctx, userid, tableInfo.DB)
	return tableInfo, db, err
}

func GetTableAndDBByTableID(ctx context.Context, tableID int) (*obj.TblTable, *obj.TblDB, error) {
	isNil, tableInfo, err := table.GetTableByID(ctx, tableID)
	if err != nil {
		return nil, nil, err
	}

	if isNil {
		return nil, nil, errs.Newf(errs.RetWebNotFindTable, "not find table [%d]", tableID)
	}

	isNil, dbInfo, err := table.GetDBByID(ctx, tableInfo.DB)
	if err != nil {
		return tableInfo, nil, err
	}

	if isNil {
		return tableInfo, nil, errs.Newf(errs.RetWebNotFindDB, "not find db [%d]", tableInfo.DB)
	}

	return tableInfo, dbInfo, nil
}

func GetTableBase(v *obj.TblTable) *pb.TableBase {
	if v == nil {
		return nil
	}

	return &pb.TableBase{
		Id:         v.Id,
		Name:       v.Name,
		Intro:      v.Intro,
		Desc:       v.Desc,
		Status:     v.Status,
		CreateTime: v.CreatedAt.Unix(),
	}
}

func GetTableByID(tables []*obj.TblTable, id int) *obj.TblTable {
	if len(tables) == 0 {
		return nil
	}

	for _, v := range tables {
		if v.Id == id {
			return v
		}
	}

	return nil
}
