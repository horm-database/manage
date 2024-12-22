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
package logic

import (
	"context"
	"fmt"

	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	sc "github.com/horm-database/server/consts"
	st "github.com/horm-database/server/model/table"
	"github.com/horm-database/server/plugin/conf"
)

// AddPlugin 新增插件
func AddPlugin(ctx context.Context, userid uint64, req *pb.AddPluginRequest) (*pb.AddPluginResponse, error) {
	data := st.TblPlugin{
		Name:    req.Name,
		Intro:   req.Intro,
		Version: "",
		Func:    req.Func,
		Online:  consts.StatusOnline,
		Source:  consts.PluginSourcePrivate,
		Desc:    req.Desc,
		Creator: userid,
		Manager: fmt.Sprint(userid),
	}

	if len(req.SupportTypes) > 0 {
		data.SupportTypes = types.JoinInt8(req.SupportTypes, ",")
	}

	id, err := table.AddPlugin(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &pb.AddPluginResponse{ID: id}, nil
}

func UpdatePlugin(ctx context.Context, userid uint64, req *pb.UpdatePluginRequest) error {
	update := horm.Map{
		"name":          req.Name,
		"intro":         req.Intro,
		"support_types": types.JoinInt8(req.SupportTypes, ","),
		"desc":          req.Desc,
	}

	return table.UpdatePluginByID(ctx, req.PluginID, update)
}

// ReplacePluginConfig 新增/修改插件配置
func ReplacePluginConfig(ctx context.Context, userid uint64, req *pb.ReplacePluginConfigRequest) error {
	data := st.TblPluginConfig{
		PluginID:      req.PluginID,
		PluginVersion: req.PluginVersion,
		Key:           req.Key,
		Name:          req.Name,
		Type:          req.Type,
		NotNull:       req.NotNull,
		MoreInfo:      req.MoreInfo,
		Default:       req.Default,
		Desc:          req.Desc,
	}

	return table.ReplacePluginConfig(ctx, &data)
}

// DelPluginConfig 删除插件配置
func DelPluginConfig(ctx context.Context, userid uint64, req *pb.DelPluginConfigRequest) error {
	return table.DelPluginConfigByKey(ctx, req.PluginID, req.PluginVersion, req.Key)
}

func PluginList(ctx context.Context, req *pb.PluginListRequest) (*pb.PluginListResponse, error) {
	pageInfo, plugins, err := table.GetPluginList(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.PluginListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      req.Page,
		Size:      req.Size,
		Plugins:   make([]*pb.PluginBase, len(plugins)),
	}

	var userIds []uint64
	for _, v := range plugins {
		userIds = append(userIds, GetUserIds(v.Creator, v.Manager)...)
	}

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for k, v := range plugins {
		ret.Plugins[k] = &pb.PluginBase{
			Id:           v.Id,
			Name:         v.Name,
			Intro:        v.Intro,
			Version:      types.SplitInt(v.Version, ","),
			Func:         v.Func,
			SupportTypes: PluginTypes(v.SupportTypes),
			Online:       v.Online,
			Source:       v.Source,
			Desc:         v.Desc,
			Creator:      userMaps[v.Creator],
			Manager:      GetUsersFromMap(userMaps, GetUserIds(v.Manager)),
			CreatedAt:    v.CreatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
		}
	}

	return &ret, nil
}

// PluginConfigs 插件配置列表
func PluginConfigs(ctx context.Context, req *pb.PluginConfigsRequest) (*pb.PluginConfigsResponse, error) {
	pluginConfigs, err := table.GetPluginConfigs(ctx, req.PluginID, req.PluginVersion)
	if err != nil {
		return nil, err
	}

	ret := pb.PluginConfigsResponse{
		Configs: make([]*pb.PluginConfig, len(pluginConfigs)),
	}

	for k, v := range pluginConfigs {
		ret.Configs[k] = &pb.PluginConfig{
			Id:       v.Id,
			Key:      v.Key,
			Name:     v.Name,
			Type:     v.Type,
			NotNull:  v.NotNull,
			MoreInfo: v.MoreInfo,
			Default:  v.Default,
			Desc:     v.Desc,
		}
	}

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func PluginsToMap(plugins []*st.TblPlugin) map[int]*st.TblPlugin {
	ret := map[int]*st.TblPlugin{}
	for _, v := range plugins {
		ret[v.Id] = v
	}
	return ret
}

func PluginConfigsToMap(pluginConfigs []*st.TblPluginConfig) map[string][]*st.TblPluginConfig {
	ret := map[string][]*st.TblPluginConfig{}
	for _, v := range pluginConfigs {
		key := fmt.Sprintf("%d_%d", v.PluginID, v.PluginVersion)
		ret[key] = append(ret[key], v)
	}
	return ret
}

func GetDefaultScheduleConfig() *conf.ScheduleConfig {
	return &conf.ScheduleConfig{
		Async:         false,
		SkipError:     false,
		Timeout:       1000,
		RequestSource: []string{"api", "web"},
		OpType:        []string{"read", "mod", "del"},
		GrayScale:     100,
		AppRule: &conf.AppRule{
			ActType: sc.ActionTypeExec,
			AppIDs:  []uint64{},
		},
		CustomRule: &conf.CustomRule{
			ActType:  sc.ActionTypeExec,
			RuleType: sc.CondTypeAny,
			Rules:    []*conf.Rule{},
		},
	}
}

func PluginTypes(typ string) []int8 {
	if typ == "" {
		return []int8{sc.PrePlugin, sc.PostPlugin, sc.DeferPlugin}
	}

	return types.SplitInt8(typ, ",")
}

func PluginTypeDesc(typ int8) string {
	switch typ {
	case sc.PrePlugin:
		return "pre-plugin"
	case sc.PostPlugin:
		return "post-plugin"
	default:
		return "defer-plugin"
	}
}
