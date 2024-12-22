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
package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/orm/obj"
)

func AddTable(ctx context.Context, tb *obj.TblTable) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_table").Insert(tb).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateTableByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_table").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetIndexTable(ctx context.Context, page, size int) (*proto.Detail, []*obj.TblTable, error) {
	pageRet := proto.Detail{}

	tables := []*obj.TblTable{}

	where := horm.Where{}
	where["status"] = consts.StatusOnline

	_, err := GetTableORM("tbl_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &tables)

	return &pageRet, tables, err
}

func GetTableByID(ctx context.Context, id int) (bool, *obj.TblTable, error) {
	table := obj.TblTable{}
	isNil, err := GetTableORM("tbl_table").FindBy("id", id).Exec(ctx, &table)
	return isNil, &table, err
}

func GetTableByIds(ctx context.Context, ids []int) ([]*obj.TblTable, error) {
	tables := []*obj.TblTable{}

	if len(ids) == 0 {
		return tables, nil
	}

	where := horm.Where{}
	where["id"] = ids

	_, err := GetTableORM("tbl_table").
		FindAll(where).
		Order("-id").
		Exec(ctx, &tables)

	return tables, err
}

func GetDBTables(ctx context.Context, dbID int) ([]*obj.TblTable, error) {
	tables := []*obj.TblTable{}

	_, err := GetTableORM("tbl_table").
		FindAllBy("db", dbID).
		Order("-id").
		Exec(ctx, &tables)

	return tables, err
}

///////////////////////////////// function /////////////////////////////////////////

func TablesToMap(tables []*obj.TblTable) map[int]*obj.TblTable {
	ret := map[int]*obj.TblTable{}

	for _, v := range tables {
		ret[v.Id] = v
	}

	return ret
}
