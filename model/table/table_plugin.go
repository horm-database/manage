package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func InsertTablePlugin(ctx context.Context, tablePlugin *table.TblTablePlugin) (int, error) {
	modResult := proto.ModResult{}
	_, err := GetTableORM("tbl_table_plugin").Insert(tablePlugin).Exec(ctx, &modResult)
	if err != nil {
		return 0, err
	}

	return modResult.ID.Int(), nil
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
