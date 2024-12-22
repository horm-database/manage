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

func InsertTablePlugin(ctx context.Context, tablePlugin *table.TblTablePlugin) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_table_plugin").Insert(tablePlugin).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateTablePluginByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_table_plugin").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func UpdateTablePluginByIDs(ctx context.Context, id []int, update horm.Map) error {
	_, err := GetTableORM("tbl_table_plugin").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func DelTablePlugin(ctx context.Context, id int) error {
	_, err := GetTableORM("tbl_table_plugin").DeleteBy("id", id).Exec(ctx)
	return err
}

func GetTablePluginByID(ctx context.Context, id int) (bool, *table.TblTablePlugin, error) {
	tablePlugin := table.TblTablePlugin{}

	isNil, err := GetTableORM("tbl_table_plugin").FindBy("id", id).Exec(ctx, &tablePlugin)

	return isNil, &tablePlugin, err
}

func GetTablePlugins(ctx context.Context, tableID int, typ ...int8) ([]*table.TblTablePlugin, error) {
	tablePlugins := []*table.TblTablePlugin{}

	where := horm.Where{}
	where["table_id"] = tableID

	if len(typ) == 1 && typ[0] != 0 {
		where["type"] = typ[0]
	} else if len(typ) > 1 {
		where["type"] = typ
	}

	_, err := GetTableORM("tbl_table_plugin").FindAll(where).Exec(ctx, &tablePlugins)

	return tablePlugins, err
}

func GetTableBackPlugin(ctx context.Context, tableID int, typ int8, frontID int) ([]*table.TblTablePlugin, error) {
	tablePlugins := []*table.TblTablePlugin{}

	where := horm.Where{
		"table_id": tableID,
		"type":     typ,
		"front":    frontID,
	}

	_, err := GetTableORM("tbl_table_plugin").FindAll(where).Exec(ctx, &tablePlugins)

	return tablePlugins, err
}
