package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/orm/obj"
)

func AddTable(ctx context.Context, tb *obj.TblTable) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_table").Insert(tb).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateTableByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_table").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetIndexTable(ctx context.Context, page, size int) (*proto.Detail, []*obj.TblTable, error) {
	pageRet := proto.Detail{}

	tables := []*obj.TblTable{}

	where := horm.Where{}
	where["status"] = consts.StatusOnline

	_, err := GetTableORM("tbl_table").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &tables)

	return &pageRet, tables, err
}

func GetTableByID(ctx context.Context, id int) (bool, *obj.TblTable, error) {
	table := obj.TblTable{}
	isNil, err := GetTableORM("tbl_table").FindBy("id", id).Exec(ctx, &table)
	return isNil, &table, err
}

func GetTableByIds(ctx context.Context, ids []int) ([]*obj.TblTable, error) {
	tables := []*obj.TblTable{}

	if len(ids) == 0 {
		return tables, nil
	}

	where := horm.Where{}
	where["id"] = ids

	_, err := GetTableORM("tbl_table").
		FindAll(where).
		Order("-id").
		Exec(ctx, &tables)

	return tables, err
}

func GetDBTables(ctx context.Context, dbID int) ([]*obj.TblTable, error) {
	tables := []*obj.TblTable{}

	_, err := GetTableORM("tbl_table").
		FindAllBy("db", dbID).
		Order("-id").
		Exec(ctx, &tables)

	return tables, err
}

///////////////////////////////// function /////////////////////////////////////////

func TablesToMap(tables []*obj.TblTable) map[int]*obj.TblTable {
	ret := map[int]*obj.TblTable{}

	for _, v := range tables {
		ret[v.Id] = v
	}

	return ret
}
