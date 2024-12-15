package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func InsertAccessDB(ctx context.Context, accessDB *table.TblAccessDB) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_access_db").Insert(accessDB).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateAccessDBByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_access_db").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetAppAccessDBs(ctx context.Context, appids []uint64, db int) ([]*table.TblAccessDB, error) {
	accessDBs := []*table.TblAccessDB{}

	where := horm.Where{}
	where["db"] = db
	where["appid"] = appids

	_, err := GetTableORM("tbl_access_db").FindAll(where).Order("-id").Exec(ctx, &accessDBs)

	return accessDBs, err
}

func GetAppAccessDBsPages(ctx context.Context, appids []uint64,
	db int, page, size int) (*proto.Detail, []*table.TblAccessDB, error) {
	pageResult := proto.Detail{}

	accessDBs := []*table.TblAccessDB{}

	where := horm.Where{}
	where["db"] = db
	where["appid"] = appids

	_, err := GetTableORM("tbl_access_db").FindAll(where).Order("-id").
		Page(page, size).Exec(ctx, &pageResult, &accessDBs)

	return &pageResult, accessDBs, err
}

func GetAppAccessDB(ctx context.Context, appid uint64, db int) (bool, *table.TblAccessDB, error) {
	accessDB := table.TblAccessDB{}

	where := horm.Where{}
	where["db"] = db
	where["appid"] = appid

	isNil, err := GetTableORM("tbl_access_db").Find(where).Exec(ctx, &accessDB)

	return isNil, &accessDB, err
}

func GetAppAccessDBListByDBID(ctx context.Context, db, page, size int) (*proto.Detail, []*table.TblAccessDB, error) {
	pageResult := proto.Detail{}

	accessDBs := []*table.TblAccessDB{}

	where := horm.Where{}
	where["db"] = db

	_, err := GetTableORM("tbl_access_db").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &accessDBs)

	return &pageResult, accessDBs, err
}

func GetAppAccessDBListByAppid(ctx context.Context, appid uint64, page, size int) (*proto.Detail, []*table.TblAccessDB, error) {
	pageResult := proto.Detail{}

	accessDBs := []*table.TblAccessDB{}

	where := horm.Where{}
	where["appid"] = appid

	_, err := GetTableORM("tbl_access_db").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &accessDBs)

	return &pageResult, accessDBs, err
}
