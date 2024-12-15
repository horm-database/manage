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
	pageResult := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID

	_, err := GetTableORM("tbl_access_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &accessTables)

	return &pageResult, accessTables, err
}

func GetAppAccessTablesPages(ctx context.Context, appids []uint64,
	tableID int, page, size int) (*proto.Detail, []*table.TblAccessTable, error) {
	pageResult := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["table_id"] = tableID
	where["appid"] = appids

	_, err := GetTableORM("tbl_access_table").FindAll(where).Order("-id").
		Page(page, size).Exec(ctx, &pageResult, &accessTables)

	return &pageResult, accessTables, err
}

func GetAppAccessTableListByAppid(ctx context.Context,
	appid uint64, page, size int) (*proto.Detail, []*table.TblAccessTable, error) {
	pageResult := proto.Detail{}

	accessTables := []*table.TblAccessTable{}

	where := horm.Where{}
	where["appid"] = appid

	_, err := GetTableORM("tbl_access_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &accessTables)

	return &pageResult, accessTables, err
}
