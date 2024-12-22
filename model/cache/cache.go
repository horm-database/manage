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
package cache

import (
	"context"

	"github.com/horm-database/manage/consts"
	"github.com/horm-database/orm"
)

// GetCacheORM 获取经校验后的 ORM
func GetCacheORM() *orm.ORM {
	return orm.NewORM(consts.CacheConfigName)
}

// GetCacheByKey 获取缓存
func GetCacheByKey(ctx context.Context, key string, ret interface{}) (bool, error) {
	return GetCacheORM().Get(key).Exec(ctx, ret)
}

// SetCacheByKey 设置缓存信息
func SetCacheByKey(ctx context.Context, key string, value interface{}, expire int) error {
	_, err := GetCacheORM().SetEX(key, value, expire).Exec(ctx)
	return err
}

// DelCacheByKey 删除缓存信息
func DelCacheByKey(ctx context.Context, key string) error {
	_, err := GetCacheORM().Del(key).Exec(ctx)
	return err
}

// TTLCacheByKey 查看缓存到期时间
func TTLCacheByKey(ctx context.Context, key string) (int, error) {
	var ttl int
	_, err := GetCacheORM().TTL(key).Exec(ctx, &ttl)

	return ttl, err
}
