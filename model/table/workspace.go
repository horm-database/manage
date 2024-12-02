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
