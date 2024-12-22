// Copyright (c) 2024 The horm-database Authors. All rights reserved.
// This file Author:  CaoHao <18500482693@163.com> .
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
	pageRet := proto.Detail{}

	products := []*TblProduct{}

	where := horm.Where{}
	if status > 0 {
		where["status"] = status
	}

	_, err := GetTableORM("tbl_product").
		FindAll(where).
		Order("-id").
		Page(page, size).
		Exec(ctx, &pageRet, &products)

	return &pageRet, products, err
}
