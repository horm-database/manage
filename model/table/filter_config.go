package table

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func ReplaceFilterConfig(ctx context.Context, filterConfig *table.TblFilterConfig) error {
	_, err := GetTableORM("tbl_filter_config").Replace(filterConfig).Exec(ctx)
	return err
}

func DelFilterConfigByKey(ctx context.Context, filterID, version int, key string) error {
	where := horm.Where{
		"filter_id":      filterID,
		"filter_version": version,
		"key":            key,
	}

	_, err := GetTableORM("tbl_filter_config").Delete(where).Exec(ctx)
	return err
}

func GetFilterConfigs(ctx context.Context, filterID, version int) ([]*table.TblFilterConfig, error) {
	filterConfigs := []*table.TblFilterConfig{}

	where := horm.Where{
		"filter_id":      filterID,
		"filter_version": version,
	}

	_, err := GetTableORM("tbl_filter_config").FindAll(where).Order("id").Exec(ctx, &filterConfigs)

	return filterConfigs, err
}

func GetConfigsByFilterIDVersions(ctx context.Context, idVersions ...int) ([]*table.TblFilterConfig, error) {
	filterConfigs := []*table.TblFilterConfig{}

	if len(idVersions) < 2 {
		return filterConfigs, nil
	}

	if len(idVersions)%2 != 0 {
		return nil, errs.New(errs.RetSystem, "GetConfigsByFilterIDVersions input idVersion is invalid")
	}

	tmp := horm.Where{}
	idVersionWhere := []horm.Where{}

	for k, v := range idVersions {
		if k%2 == 0 {
			tmp["filter_id"] = v
		} else {
			tmp["filter_version"] = v
			idVersionWhere = append(idVersionWhere, tmp)
			tmp = horm.Where{}
		}
	}

	where := horm.Where{
		"OR": idVersionWhere,
	}

	_, err := GetTableORM("tbl_filter_config").FindAll(where).Order("id").Exec(ctx, &filterConfigs)

	return filterConfigs, err
}
