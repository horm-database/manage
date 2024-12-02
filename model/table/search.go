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
