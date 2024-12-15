package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func AddPlugin(ctx context.Context, plugin *table.TblPlugin) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_plugin").Insert(plugin).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdatePluginByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_plugin").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetPluginList(ctx context.Context, page, size int) (*proto.Detail, []*table.TblPlugin, error) {
	pageResult := proto.Detail{}

	plugins := []*table.TblPlugin{}

	_, err := GetTableORM("tbl_plugin").
		FindAll().
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &plugins)

	return &pageResult, plugins, err
}

func GetPluginByIDs(ctx context.Context, pluginIDs []int) ([]*table.TblPlugin, error) {
	plugins := []*table.TblPlugin{}

	_, err := GetTableORM("tbl_plugin").
		FindAllBy("id", pluginIDs).
		Exec(ctx, &plugins)

	return plugins, err
}

func GetPluginByID(ctx context.Context, pluginID int) (bool, *table.TblPlugin, error) {
	plugin := table.TblPlugin{}

	isNil, err := GetTableORM("tbl_plugin").FindBy("id", pluginID).Exec(ctx, &plugin)

	return isNil, &plugin, err
}
