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

	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

var CurrentWorkspace *table.TblWorkspace

func GetCurrentWorkspace(ctx context.Context) (*table.TblWorkspace, error) {
	if CurrentWorkspace == nil {
		_, err := GetTableORM("tbl_workspace").Find().Exec(ctx, &CurrentWorkspace)
		return CurrentWorkspace, err
	}

	return CurrentWorkspace, nil
}

func GetWorkspace(ctx context.Context, workspace string) (*table.TblWorkspace, error) {
	workspaceInfo := table.TblWorkspace{}

	_, err := GetTableORM("tbl_workspace").
		Eq("workspace", workspace).
		Find().Exec(ctx, &workspaceInfo)

	return &workspaceInfo, err
}

func GetWorkspaceByID(ctx context.Context, id int) (*table.TblWorkspace, error) {
	if CurrentWorkspace != nil && CurrentWorkspace.Id == id {
		return CurrentWorkspace, nil
	}

	workspaceInfo := table.TblWorkspace{}

	_, err := GetTableORM("tbl_workspace").FindBy("id", id).Exec(ctx, &workspaceInfo)

	return &workspaceInfo, err
}

func UpdateWorkspaceByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_workspace").Eq("id", id).Update(update).Exec(ctx)
	return err
}
