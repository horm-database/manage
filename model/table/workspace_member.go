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
	pageResult := proto.Detail{}

	members := []*TblWorkspaceMember{}

	where := horm.Where{}
	where["workspace_id"] = workspaceID

	_, err := GetTableORM("tbl_workspace_member").
		FindAll(where).
		Order("status", "-updated_at").
		Page(page, size).
		Exec(ctx, &pageResult, &members)

	return &pageResult, members, err
}

func GetWorkspaceMembersJoined(ctx context.Context,
	workspaceID, page, size int) (*proto.Detail, []*TblWorkspaceMember, error) {
	pageResult := proto.Detail{}

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
		Exec(ctx, &pageResult, &members)

	return &pageResult, members, err
}
