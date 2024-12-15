package table

import (
	"context"
	"time"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/consts"
)

func InsertProductMember(ctx context.Context, member *TblProductMember) (int, error) {
	modRet := proto.ModRet{}
	_, err := GetTableORM("tbl_product_member").Insert(member).Exec(ctx, &modRet)
	if err != nil {
		return 0, err
	}

	return modRet.ID.Int(), nil
}

func ReplaceProductMember(ctx context.Context, member horm.Map) error {
	_, err := GetTableORM("tbl_product_member").Replace(member).Exec(ctx)
	return err
}

func UpdateProductMemberByID(ctx context.Context, id int, update horm.Map) error {
	_, err := GetTableORM("tbl_product_member").Update(update).Eq("id", id).Exec(ctx)
	return err
}

func GetProductMemberByUser(ctx context.Context, productID int, userid uint64) (bool, *TblProductMember, error) {
	member := TblProductMember{}

	where := horm.Where{
		"product_id": productID,
		"userid":     userid,
	}

	isNil, err := GetTableORM("tbl_product_member").Find(where).Exec(ctx, &member)

	return isNil, &member, err
}

func GetProductMemberByUsers(ctx context.Context, productID int, userIds []uint64) ([]*TblProductMember, error) {
	members := []*TblProductMember{}

	where := horm.Where{
		"product_id": productID,
		"userid":     userIds,
	}

	_, err := GetTableORM("tbl_product_member").FindAll(where).Exec(ctx, &members)

	return members, err
}

func GetProductMembersAll(ctx context.Context,
	productID, page, size int) (*proto.Detail, []*TblProductMember, error) {
	pageResult := proto.Detail{}

	members := []*TblProductMember{}

	where := horm.Where{}
	where["product_id"] = productID

	_, err := GetTableORM("tbl_product_member").
		FindAll(where).
		Order("status", "-updated_at").
		Page(page, size).
		Exec(ctx, &pageResult, &members)

	return &pageResult, members, err
}

func GetProductMembersJoined(ctx context.Context,
	productID, page, size int) (*proto.Detail, []*TblProductMember, error) {
	pageResult := proto.Detail{}

	members := []*TblProductMember{}

	where := horm.Where{}
	where["product_id"] = productID
	where["status"] = []int8{consts.ProductMemberStatusRenewal, consts.ProductMemberStatusChangeRole, consts.ProductMemberStatusJoined}
	where["OR"] = horm.OR{
		"expire_time":   0,
		"expire_time >": time.Now().Unix(),
	}

	_, err := GetTableORM("tbl_product_member").
		FindAll(where).
		Order("-updated_at").
		Page(page, size).
		Exec(ctx, &pageResult, &members)

	return &pageResult, members, err
}
