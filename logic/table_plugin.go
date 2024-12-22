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

package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/json"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	cc "github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/server/consts"
	st "github.com/horm-database/server/model/table"
	"github.com/horm-database/server/plugin/conf"
	"github.com/samber/lo"
)

// TablePlugins 表插件
func TablePlugins(ctx context.Context, userid uint64, tableID int) (*pb.TablePluginsResponse, error) {
	_, _, err := IsTableManager(ctx, userid, tableID)
	if err != nil {
		return nil, err
	}

	tablePlugins, err := table.GetTablePlugins(ctx, tableID)
	if err != nil {
		return nil, err
	}

	ret := pb.TablePluginsResponse{
		PrePlugins:   []*pb.TablePlugin{},
		PostPlugins:  []*pb.TablePlugin{},
		DeferPlugins: []*pb.TablePlugin{},
	}

	if len(tablePlugins) == 0 {
		return &ret, nil
	}

	var pluginIDs []int
	for _, tablePlugin := range tablePlugins {
		pluginIDs = append(pluginIDs, tablePlugin.PluginID)
	}

	pluginInfos, err := table.GetPluginByIDs(ctx, pluginIDs)
	if err != nil {
		return nil, err
	}

	pluginInfoMaps := PluginsToMap(pluginInfos)

	var idVersions []int
	for _, tablePlugin := range tablePlugins {
		idVersions = append(idVersions, tablePlugin.PluginID)
		idVersions = append(idVersions, tablePlugin.PluginVersion)
	}

	pluginConfigs, err := table.GetConfigsByPluginIDVersions(ctx, idVersions...)
	if err != nil {
		return nil, err
	}

	pluginConfigMaps := PluginConfigsToMap(pluginConfigs)

	var userIds []uint64
	for _, pluginInfo := range pluginInfos {
		userIds = append(userIds, GetUserIds(pluginInfo.Creator, pluginInfo.Manager)...)
	}

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, tablePlugin := range tablePlugins {
		plugin := pb.TablePlugin{
			Id:             tablePlugin.Id,
			TableId:        tablePlugin.TableId,
			PluginID:       tablePlugin.PluginID,
			PluginVersion:  tablePlugin.PluginVersion,
			Type:           tablePlugin.Type,
			Front:          tablePlugin.Front,
			Desc:           tablePlugin.Desc,
			Status:         tablePlugin.Status,
			CreatedAt:      tablePlugin.CreatedAt.Unix(),
			UpdatedAt:      tablePlugin.UpdatedAt.Unix(),
			PluginConfigs:  []*pb.TablePluginConfig{},
			ScheduleConfig: &conf.ScheduleConfig{},
		}

		pluginInfo := pluginInfoMaps[tablePlugin.PluginID]
		if pluginInfo != nil {
			plugin.PluginInfo = &pb.PluginBase{
				Id:           pluginInfo.Id,
				Name:         pluginInfo.Name,
				Intro:        pluginInfo.Intro,
				Version:      types.SplitInt(pluginInfo.Version, ","),
				Func:         pluginInfo.Func,
				SupportTypes: PluginTypes(pluginInfo.SupportTypes),
				Online:       pluginInfo.Online,
				Source:       pluginInfo.Source,
				Desc:         pluginInfo.Desc,
				Creator:      userMaps[pluginInfo.Creator],
				Manager:      GetUsersFromMap(userMaps, GetUserIds(pluginInfo.Manager)),
				CreatedAt:    pluginInfo.CreatedAt.Unix(),
				UpdatedAt:    pluginInfo.UpdatedAt.Unix(),
			}
		}

		if tablePlugin.ScheduleConfig == "" {
			plugin.ScheduleConfig = GetDefaultScheduleConfig()
		} else {
			_ = json.Api.Unmarshal([]byte(tablePlugin.ScheduleConfig), &plugin.ScheduleConfig)
		}

		tablePluginConfigValues := map[string]interface{}{}
		if tablePlugin.Config != "" {
			_ = json.Api.Unmarshal([]byte(tablePlugin.Config), &tablePluginConfigValues)
		}

		configs := pluginConfigMaps[fmt.Sprintf("%d_%d", tablePlugin.PluginID, tablePlugin.PluginVersion)]
		if len(configs) > 0 {
			for _, conf := range configs {
				tablePluginConfig := pb.TablePluginConfig{
					Config: &pb.PluginConfig{
						Id:       conf.Id,
						Key:      conf.Key,
						Name:     conf.Name,
						Type:     conf.Type,
						NotNull:  conf.NotNull,
						MoreInfo: conf.MoreInfo,
						Default:  conf.Default,
						Desc:     conf.Desc,
					},
				}

				configValue, ok := tablePluginConfigValues[conf.Key]

				if ok {
					tablePluginConfig.IsSet = true
					tablePluginConfig.Value = configValue
				}

				plugin.PluginConfigs = append(plugin.PluginConfigs, &tablePluginConfig)
			}
		}

		switch tablePlugin.Type {
		case consts.PrePlugin:
			ret.PrePlugins = append(ret.PrePlugins, &plugin)
		case consts.PostPlugin:
			ret.PostPlugins = append(ret.PostPlugins, &plugin)
		case consts.DeferPlugin:
			ret.DeferPlugins = append(ret.DeferPlugins, &plugin)
		}
	}

	ret.PrePlugins, err = sortTablePlugins("table pre-plugin", ret.PrePlugins)
	if err != nil {
		return nil, err
	}

	ret.PostPlugins, err = sortTablePlugins("table post-plugin", ret.PostPlugins)
	if err != nil {
		return nil, err
	}

	ret.DeferPlugins, err = sortTablePlugins("table defer-plugin", ret.DeferPlugins)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// AddTablePlugin 新增表插件
func AddTablePlugin(ctx context.Context, userid uint64,
	req *pb.AddTablePluginRequest) (*pb.AddTablePluginResponse, error) {
	_, _, err := IsTableManager(ctx, userid, req.TableId)
	if err != nil {
		return nil, err
	}

	isNil, plugin, err := table.GetPluginByID(ctx, req.PluginID)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.Newf(errs.RetWebNotFindPlugin, "not find table %d", req.PluginID)
	}

	supportVersions := types.SplitInt(plugin.Version, ",")
	if lo.IndexOf(supportVersions, req.PluginVersion) == -1 {
		return nil, errs.Newf(errs.RetWebNotFindPlugin,
			"plugin %s not support version %d", plugin.Name, req.PluginVersion)
	}

	if plugin.Online != cc.StatusOnline {
		return nil, errs.Newf(errs.RetWebNotFindPlugin, "plugin %s is not online", plugin.Name)
	}

	supportTypes := PluginTypes(plugin.SupportTypes)
	if lo.IndexOf(supportTypes, req.Type) == -1 {
		return nil, errs.Newf(errs.RetWebNotFindPlugin,
			"plugin %s not support %s", plugin.Name, PluginTypeDesc(req.Type))
	}

	tablePlugins, err := table.GetTablePlugins(ctx, req.TableId, req.Type)
	if err != nil {
		return nil, err
	}

	if len(tablePlugins) == 0 && req.Front != 0 {
		return nil, errs.Newf(errs.RetWebIsFirstPlugin, "this is first plugin, front must be zero")
	}

	insertTablePlugin := st.TblTablePlugin{
		TableId:        req.TableId,
		PluginID:       req.PluginID,
		PluginVersion:  req.PluginVersion,
		Type:           req.Type,
		Front:          req.Front,
		ScheduleConfig: json.MarshalToString(req.ScheduleConfig),
		Config:         json.MarshalToString(req.PluginConfigs),
		Desc:           req.Desc,
		Status:         cc.StatusOnline,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	id, err := table.InsertTablePlugin(ctx, &insertTablePlugin)
	if err != nil {
		return nil, err
	}

	if len(tablePlugins) > 0 {
		for _, tablePlugin := range tablePlugins {
			if tablePlugin.Front == req.Front {
				_ = table.UpdateTablePluginByID(ctx, tablePlugin.Id, horm.Map{"front": id})
			}
		}
	}

	return &pb.AddTablePluginResponse{ID: id}, nil
}

// UpdateTablePlugin 更新表插件
func UpdateTablePlugin(ctx context.Context, userid uint64, req *pb.UpdateTablePluginRequest) error {
	isNil, tablePlugin, err := table.GetTablePluginByID(ctx, req.Id)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebNotFindTablePlugin, "not find table plugin [%d]", req.Id)
	}

	_, _, err = IsTableManager(ctx, userid, tablePlugin.TableId)
	if err != nil {
		return err
	}

	isNil, plugin, err := table.GetPluginByID(ctx, tablePlugin.PluginID)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebNotFindPlugin, "not find table %d", tablePlugin.PluginID)
	}

	supportVersions := types.SplitInt(plugin.Version, ",")
	if lo.IndexOf(supportVersions, req.PluginVersion) == -1 {
		return errs.Newf(errs.RetWebNotFindPlugin,
			"plugin %s not support version %d", plugin.Name, req.PluginVersion)
	}

	supportTypes := PluginTypes(plugin.SupportTypes)
	if lo.IndexOf(supportTypes, req.Type) == -1 {
		return errs.Newf(errs.RetWebNotFindPlugin,
			"plugin %s not support %s", plugin.Name, PluginTypeDesc(req.Type))
	}

	if req.Type != tablePlugin.Type || req.Front != tablePlugin.Front {
		backTablePlugins, err := table.GetTableBackPlugin(ctx, tablePlugin.TableId, tablePlugin.Type, tablePlugin.Id)
		if err != nil {
			return err
		}

		tableTypePlugins, err := table.GetTablePlugins(ctx, tablePlugin.TableId, req.Type)
		if err != nil {
			return err
		}

		// 先剔除自身
		if len(backTablePlugins) > 0 {
			var ids = getTablePluginsID(backTablePlugins)
			err = table.UpdateTablePluginByIDs(ctx, ids, horm.Map{"front": tablePlugin.Front})
			if err != nil {
				return err
			}
		}

		// 再插入新位置
		if len(tableTypePlugins) > 0 {
			for _, tf := range tableTypePlugins {
				if tf.Id != req.Id && tf.Front == req.Front {
					err = table.UpdateTablePluginByID(ctx, tf.Id, horm.Map{"front": req.Id})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	updateTablePlugin := horm.Map{
		"plugin_version":  req.PluginVersion,
		"type":            req.Type,
		"front":           req.Front,
		"schedule_config": json.MarshalToString(req.ScheduleConfig),
		"config":          json.MarshalToString(req.PluginConfigs),
		"desc":            req.Desc,
	}

	err = table.UpdateTablePluginByID(ctx, req.Id, updateTablePlugin)
	if err != nil {
		return err
	}

	return nil
}

// DelTablePlugin 删除表插件
func DelTablePlugin(ctx context.Context, userid uint64, id int) error {
	isNil, tablePlugin, err := table.GetTablePluginByID(ctx, id)
	if err != nil {
		return err
	}
	if isNil {
		return errs.Newf(errs.RetWebNotFindTablePlugin, "not find table plugin [%d]", id)
	}

	_, _, err = IsTableManager(ctx, userid, tablePlugin.TableId)
	if err != nil {
		return err
	}

	backTablePlugins, err := table.GetTableBackPlugin(ctx, tablePlugin.TableId, tablePlugin.Type, tablePlugin.Id)
	if err != nil {
		return err
	}

	if len(backTablePlugins) > 0 {
		var ids = getTablePluginsID(backTablePlugins)
		err = table.UpdateTablePluginByIDs(ctx, ids, horm.Map{"front": tablePlugin.Front})
		if err != nil {
			return err
		}
	}

	_ = table.DelTablePlugin(ctx, id)

	return nil
}

///////////////////////////////// function /////////////////////////////////////////

func sortTablePlugins(typ string, tablePlugins []*pb.TablePlugin) ([]*pb.TablePlugin, error) {
	if len(tablePlugins) == 0 {
		return []*pb.TablePlugin{}, nil
	}

	var head *pb.TablePlugin

	for _, tablePlugin := range tablePlugins {
		if tablePlugin.Front == 0 {
			head = tablePlugin
			break
		}
	}

	if head == nil {
		return nil, errs.Newf(errs.ErrPrefixPluginNotFount,
			"table_id %d not find head of %s", tablePlugins[0].TableId, typ)
	}

	ret := []*pb.TablePlugin{}
	ret = append(ret, head)

	currentTablePlugin := head
	for i := 0; i < len(tablePlugins); i++ {
		backTablePlugin := findBackTablePlugin(currentTablePlugin, tablePlugins)
		if backTablePlugin == nil { // 最后一个
			break
		}

		ret = append(ret, backTablePlugin)
		currentTablePlugin = backTablePlugin
	}

	return ret, nil
}

func findBackTablePlugin(currentTablePlugin *pb.TablePlugin, tablePlugins []*pb.TablePlugin) *pb.TablePlugin {
	for _, tablePlugin := range tablePlugins {
		if tablePlugin.Front == currentTablePlugin.Id {
			return tablePlugin
		}
	}
	return nil
}

func getTablePluginsID(tablePlugins []*st.TblTablePlugin) []int {
	var ids []int
	for _, v := range tablePlugins {
		ids = append(ids, v.Id)
	}

	return ids
}
