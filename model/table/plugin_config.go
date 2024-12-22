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

	"github.com/horm-database/common/errs"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/server/model/table"
)

func ReplacePluginConfig(ctx context.Context, pluginConfig *table.TblPluginConfig) error {
	_, err := GetTableORM("tbl_plugin_config").Replace(pluginConfig).Exec(ctx)
	return err
}

func DelPluginConfigByKey(ctx context.Context, pluginID, version int, key string) error {
	where := horm.Where{
		"plugin_id":      pluginID,
		"plugin_version": version,
		"key":            key,
	}

	_, err := GetTableORM("tbl_plugin_config").Delete(where).Exec(ctx)
	return err
}

func GetPluginConfigs(ctx context.Context, pluginID, version int) ([]*table.TblPluginConfig, error) {
	pluginConfigs := []*table.TblPluginConfig{}

	where := horm.Where{
		"plugin_id":      pluginID,
		"plugin_version": version,
	}

	_, err := GetTableORM("tbl_plugin_config").FindAll(where).Order("id").Exec(ctx, &pluginConfigs)

	return pluginConfigs, err
}

func GetConfigsByPluginIDVersions(ctx context.Context, idVersions ...int) ([]*table.TblPluginConfig, error) {
	pluginConfigs := []*table.TblPluginConfig{}

	if len(idVersions) < 2 {
		return pluginConfigs, nil
	}

	if len(idVersions)%2 != 0 {
		return nil, errs.New(errs.ErrSystem, "GetConfigsByPluginIDVersions input idVersion is invalid")
	}

	tmp := horm.Where{}
	idVersionWhere := []horm.Where{}

	for k, v := range idVersions {
		if k%2 == 0 {
			tmp["plugin_id"] = v
		} else {
			tmp["plugin_version"] = v
			idVersionWhere = append(idVersionWhere, tmp)
			tmp = horm.Where{}
		}
	}

	where := horm.Where{
		"OR": idVersionWhere,
	}

	_, err := GetTableORM("tbl_plugin_config").FindAll(where).Order("id").Exec(ctx, &pluginConfigs)

	return pluginConfigs, err
}
