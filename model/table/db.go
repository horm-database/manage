package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/orm/obj"
)

func AddDB(ctx context.Context, db *obj.TblDB) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_db").Insert(db).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateDBByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_db").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetDBByID(ctx context.Context, id int) (bool, *obj.TblDB, error) {
	db := obj.TblDB{}

	isNil, err := GetTableORM("tbl_db").FindBy("id", id).Exec(ctx, &db)

	return isNil, &db, err
}

func GetDBByIds(ctx context.Context, ids []int) ([]*obj.TblDB, error) {
	dbs := []*obj.TblDB{}

	_, err := GetTableORM("tbl_db").FindAllBy("id", ids).Exec(ctx, &dbs)

	return dbs, err
}

func GetProductDBs(ctx context.Context, productID int) ([]*obj.TblDB, error) {
	dbs := []*obj.TblDB{}

	_, err := GetTableORM("tbl_db").FindAllBy("product_id", productID).Order("-id").Exec(ctx, &dbs)

	return dbs, err
}
