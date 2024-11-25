package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func AddFilter(ctx context.Context, filter *table.TblFilter) (int, error) {
	modResult := proto.ModResult{}
	_, err := GetTableORM("tbl_filter").Insert(filter).Exec(ctx, &modResult)
	if err != nil {
		return 0, err
	}

	return modResult.ID.Int(), nil
}

func UpdateFilterByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_filter").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetFilterList(ctx context.Context, page, size int) (*proto.Detail, []*table.TblFilter, error) {
	pageResult := proto.Detail{}

	filters := []*table.TblFilter{}

	_, err := GetTableORM("tbl_filter").
		FindAll().
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &filters)

	return &pageResult, filters, err
}

func GetFilterByIDs(ctx context.Context, filterIDs []int) ([]*table.TblFilter, error) {
	filters := []*table.TblFilter{}

	_, err := GetTableORM("tbl_filter").
		FindAllBy("id", filterIDs).
		Exec(ctx, &filters)

	return filters, err
}

func GetFilterByID(ctx context.Context, filterID int) (bool, *table.TblFilter, error) {
	filter := table.TblFilter{}

	isNil, err := GetTableORM("tbl_filter").FindBy("id", filterID).Exec(ctx, &filter)

	return isNil, &filter, err
}
