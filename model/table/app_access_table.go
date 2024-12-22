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

package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func GetAppAccessTables(ctx context.Context, appids []uint64, tableID int) ([]*table.TblAccessTable, error) {
	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID
	where["appid"] = appids

	_, err := GetTableORM("tbl_access_table").FindAll(where).Order("-id").Exec(ctx, &accessTables)

	return accessTables, err
}

func GetAppAccessTable(ctx context.Context, appid uint64, tableID int) (bool, *table.TblAccessTable, error) {
	accessTable := table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID
	where["appid"] = appid

	isNil, err := GetTableORM("tbl_access_table").Find(where).Exec(ctx, &accessTable)

	return isNil, &accessTable, err
}

func InsertAccessTable(ctx context.Context, accessTable *table.TblAccessTable) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_access_table").Insert(accessTable).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateAccessTableByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_access_table").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetAppAccessTableListByTableID(ctx context.Context,
	tableID, page, size int) (*proto.Detail, []*table.TblAccessTable, error) {
	pageRet := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID

	_, err := GetTableORM("tbl_access_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &accessTables)

	return &pageRet, accessTables, err
}

func GetAppAccessTablesPages(ctx context.Context, appids []uint64,
	tableID int, page, size int) (*proto.Detail, []*table.TblAccessTable, error) {
	pageRet := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID
	where["appid"] = appids

	_, err := GetTableORM("tbl_access_table").FindAll(where).Order("-id").
		Page(page, size).Exec(ctx, &pageRet, &accessTables)

	return &pageRet, accessTables, err
}

func GetAppAccessTableListByAppid(ctx context.Context,
	appid uint64, page, size int) (*proto.Detail, []*table.TblAccessTable, error) {
	pageRet := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["appid"] = appid

	_, err := GetTableORM("tbl_access_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &accessTables)

	return &pageRet, accessTables, err
}
