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
	"time"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/consts"
)

func InsertWorkspaceMember(ctx context.Context, member *TblWorkspaceMember) error {
	_, err := GetTableORM("tbl_workspace_member").Insert(member).Exec(ctx)
	return err
}

func ReplaceWorkspaceMember(ctx context.Context, member horm.Map) error {
	_, err := GetTableORM("tbl_workspace_member").Replace(member).Exec(ctx)
	return err
}

func UpdateWorkspaceMemberByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_workspace_member").Update(update).Eq("id", id).Exec(ctx)
	return err
}

func GetWorkspaceMemberByUser(ctx context.Context, workspaceID int, userid uint64) (bool, *TblWorkspaceMember, error) {
	member := TblWorkspaceMember{}

	where := horm.Where{
		"userid":       userid,
		"workspace_id": workspaceID,
	}

	isNil, err := GetTableORM("tbl_workspace_member").Find(where).Exec(ctx, &member)

	return isNil, &member, err
}

func GetWorkspaceMemberByUsers(ctx context.Context, workspaceID int, userIds []uint64) ([]*TblWorkspaceMember, error) {
	members := []*TblWorkspaceMember{}

	where := horm.Where{
		"workspace_id": workspaceID,
		"userid":       userIds,
	}

	_, err := GetTableORM("tbl_workspace_member").FindAll(where).Exec(ctx, &members)

	return members, err
}

func GetWorkspaceMembersAll(ctx context.Context,
	workspaceID, page, size int) (*proto.Detail, []*TblWorkspaceMember, error) {
	pageRet := proto.Detail{}

	members := []*TblWorkspaceMember{}

	where := horm.Where{}
	where["workspace_id"] = workspaceID

	_, err := GetTableORM("tbl_workspace_member").
		FindAll(where).
		Order("status", "-updated_at").
		Page(page, size).
		Exec(ctx, &pageRet, &members)

	return &pageRet, members, err
}

func GetWorkspaceMembersJoined(ctx context.Context,
	workspaceID, page, size int) (*proto.Detail, []*TblWorkspaceMember, error) {
	pageRet := proto.Detail{}

	members := []*TblWorkspaceMember{}

	where := horm.Where{}
	where["workspace_id"] = workspaceID
	where["status"] = []int8{consts.WorkspaceMemberStatusRenewal, consts.WorkspaceMemberStatusJoined}
	where["OR"] = horm.OR{
		"expire_time":   0,
		"expire_time >": time.Now().Unix(),
	}

	_, err := GetTableORM("tbl_workspace_member").
		FindAll(where).
		Order("-updated_at").
		Page(page, size).
		Exec(ctx, &pageRet, &members)

	return &pageRet, members, err
}
