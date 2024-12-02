package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func InsertTableFilter(ctx context.Context, tableFilter *table.TblTableFilter) (int, error) {
	modResult := proto.ModResult{}
	_, err := GetTableORM("tbl_table_filter").Insert(tableFilter).Exec(ctx, &modResult)
	if err != nil {
		return 0, err
	}

	return modResult.ID.Int(), nil
}

func UpdateTableFilterByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_table_filter").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func UpdateTableFilterByIDs(ctx context.Context, id []int, update horm.Map) error {
	_, err := GetTableORM("tbl_table_filter").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func DelTableFilter(ctx context.Context, id int) error {
	_, err := GetTableORM("tbl_table_filter").DeleteBy("id", id).Exec(ctx)
	return err
}

func GetTableFilterByID(ctx context.Context, id int) (bool, *table.TblTableFilter, error) {
	tableFilter := table.TblTableFilter{}

	isNil, err := GetTableORM("tbl_table_filter").FindBy("id", id).Exec(ctx, &tableFilter)

	return isNil, &tableFilter, err
}

func GetTableFilters(ctx context.Context, tableID int, typ ...int8) ([]*table.TblTableFilter, error) {
	tableFilters := []*table.TblTableFilter{}

	where := horm.Where{}
	where["table_id"] = tableID

	if len(typ) == 1 && typ[0] != 0 {
		where["type"] = typ[0]
	} else if len(typ) > 1 {
		where["type"] = typ
	}

	_, err := GetTableORM("tbl_table_filter").FindAll(where).Exec(ctx, &tableFilters)

	return tableFilters, err
}

func GetTableBackFilter(ctx context.Context, tableID int, typ int8, frontID int) ([]*table.TblTableFilter, error) {
	tableFilters := []*table.TblTableFilter{}

	where := horm.Where{
		"table_id": tableID,
		"type":     typ,
		"front":    frontID,
	}

	_, err := GetTableORM("tbl_table_filter").FindAll(where).Exec(ctx, &tableFilters)

	return tableFilters, err
}
