// Copyright (c) 2024 The horm-database Authors (such as CaoHao <18500482693@163.com>). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
