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
