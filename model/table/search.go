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
	"time"

	"github.com/horm-database/common/codec"
	"github.com/horm-database/go-horm/horm"
)

func AddSearchKeyword(ctx context.Context, sk *TblSearchKeyword) {
	go func(ctx context.Context) {
		_, err := GetTableORM("tbl_search_keyword").Replace(sk).Exec(ctx)
		if err != nil { //重试
			time.Sleep(time.Second)
			_, err = GetTableORM("tbl_search_keyword").Replace(sk).Exec(ctx)
			if err != nil {
				time.Sleep(time.Second)
				_, _ = GetTableORM("tbl_search_keyword").Replace(sk).Exec(ctx)
			}
		}
	}(codec.CloneContext(ctx))
}

func AddSearchKeywords(ctx context.Context, sks []*TblSearchKeyword) {
	go func(ctx context.Context) {
		_, err := GetTableORM("tbl_search_keyword").Replace(sks).Exec(ctx)
		if err != nil { //重试
			time.Sleep(time.Second)
			_, err = GetTableORM("tbl_search_keyword").Replace(sks).Exec(ctx)
			if err != nil { //重试
				time.Sleep(time.Second)
				_, _ = GetTableORM("tbl_search_keyword").Replace(sks).Exec(ctx)
			}
		}
	}(codec.CloneContext(ctx))
}

func DelSearchKeywords(ctx context.Context, typ int8, sid int, field, skey string) {
	where := horm.Where{
		"type":  typ,
		"sid":   sid,
		"field": field,
		"skey":  skey,
	}

	go func(ctx context.Context) {
		_, err := GetTableORM("tbl_search_keyword").Delete(where).Exec(ctx)
		if err != nil { //重试
			time.Sleep(time.Second)
			_, err = GetTableORM("tbl_search_keyword").Delete(where).Exec(ctx)
			if err != nil { //重试
				time.Sleep(time.Second)
				_, err = GetTableORM("tbl_search_keyword").Delete(where).Exec(ctx)
			}
		}
	}(codec.CloneContext(ctx))
}
