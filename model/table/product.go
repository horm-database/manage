package table

import (
	"context"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
)

func AddProduct(ctx context.Context, product *TblProduct) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_product").Insert(product).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func UpdateProductByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_product").Eq("id", id).Update(update).Exec(ctx)
	return err
}

func GetProductByID(ctx context.Context, id int) (bool, *TblProduct, error) {
	product := TblProduct{}

	isNil, err := GetTableORM("tbl_product").FindBy("id", id).Exec(ctx, &product)

	return isNil, &product, err
}

func GetProductList(ctx context.Context, status int8, page, size int) (*proto.Detail, []*TblProduct, error) {
	pageResult := proto.Detail{}

	products := []*TblProduct{}

	where := horm.Where{}
	if status > 0 {
		where["status"] = status
	}

	_, err := GetTableORM("tbl_product").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageResult, &products)

	return &pageResult, products, err
}
