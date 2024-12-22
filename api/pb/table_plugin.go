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
package pb

import (
	"github.com/horm-database/server/plugin/conf"
)

type TablePluginsResponse struct {
	PrePlugins   []*TablePlugin `json:"pre_plugins"`   // 前置插件
	PostPlugins  []*TablePlugin `json:"post_plugins"`  // 后置插件
	DeferPlugins []*TablePlugin `json:"defer_plugins"` // 延迟插件
}

type AddTablePluginRequest struct {
	TableId        int                    `json:"table_id"`        // 表id
	PluginID       int                    `json:"plugin_id"`       // 插件id
	PluginVersion  int                    `json:"plugin_version"`  // plugin 版本
	Type           int8                   `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                    `json:"front"`           // plugin execute front of me
	Desc           string                 `json:"desc"`            // 描述
	ScheduleConfig *conf.ScheduleConfig   `json:"schedule_config"` // 插件调度配置
	PluginConfigs  map[string]interface{} `json:"plugin_configs"`  // 插件配置
}

type AddTablePluginResponse struct {
	ID int `json:"id"`
}

type UpdateTablePluginRequest struct {
	Id             int                    `json:"id"`
	PluginVersion  int                    `json:"plugin_version"`  // plugin 版本
	Type           int8                   `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                    `json:"front"`           // plugin execute front of me
	Desc           string                 `json:"desc"`            // 描述
	ScheduleConfig *conf.ScheduleConfig   `json:"schedule_config"` // 插件调度配置
	PluginConfigs  map[string]interface{} `json:"plugin_configs"`  // 插件配置
}

type DelTablePluginRequest struct {
	Id int `json:"id"`
}

type TablePlugin struct {
	Id             int                  `json:"id"`              // id
	TableId        int                  `json:"table_id"`        // 表id
	PluginID       int                  `json:"plugin_id"`       // 插件id
	PluginVersion  int                  `json:"plugin_version"`  // plugin 版本
	Type           int8                 `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                  `json:"front"`           // plugin execute front of me
	Desc           string               `json:"desc"`            // 描述
	Status         int8                 `json:"status"`          // 状态 1-启用 2-停用
	CreatedAt      int64                `json:"create_time"`     // 添加时间
	UpdatedAt      int64                `json:"update_time"`     // 最后修改时间
	ScheduleConfig *conf.ScheduleConfig `json:"schedule_config"` // 插件调度配置
	PluginInfo     *PluginBase          `json:"plugin_info"`     // 插件信息
	PluginConfigs  []*TablePluginConfig `json:"plugin_configs"`  // 插件配置
}

type TablePluginConfig struct {
	Config *PluginConfig `json:"config"` // 插件配置
	IsSet  bool          `json:"is_set"` // 是否已经设置
	Value  interface{}   `json:"value"`  // 设置值
}
