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

// DBBase 数据库基础信息
type DBBase struct {
	Id         int    `json:"db_id"`
	Name       string `json:"name"`        // 数据库名称
	Intro      string `json:"intro"`       // 中文简介
	Desc       string `json:"desc"`        // 详细描述
	ProductID  int    `json:"product_id"`  // 产品id
	Type       int    `json:"type"`        // 数据库类型 0-nil（仅执行插件） 1-elastic 2-mongo 3-redis 10-mysql 11-postgresql 12-clickhouse 13-oracle 14-DB2 15-sqlite
	Version    string `json:"version"`     // 数据库版本，比如elastic v6，v7
	Status     int8   `json:"status"`      // 状态 1-正常 2-下线
	CreateTime int64  `json:"create_time"` // 记录创建时间
}

type AddDBRequest struct {
	Name      string `json:"name"`       // 数据库名称
	Intro     string `json:"intro"`      // 简介
	Desc      string `json:"desc"`       // 详细介绍
	ProductID int    `json:"product_id"` // 产品id

	// db params
	WriteTimeout int  `json:"write_timeout"` // 写超时（毫秒）
	ReadTimeout  int  `json:"read_timeout"`  // 读超时（毫秒）
	WarnTimeout  int  `json:"warn_timeout"`  // 告警超时（ms），如果请求耗时超过这个时间，就会打 warning 日志
	OmitError    int8 `json:"omit_error"`    // 是否忽略 error 日志，0-否 1-是
	Debug        int8 `json:"debug"`         // 是否开启 debug 日志，正常的数据库请求也会被打印到日志，0-否 1-是，会造成海量日志，慎重开启

	// db address
	Type       int    `json:"type"`        // 数据库类型 0-nil（仅执行插件） 1-elastic 2-mongo 3-redis 10-mysql 11-postgresql 12-clickhouse 13-oracle 14-DB2 15-sqlite
	Version    string `json:"version"`     // 数据库版本，比如elastic v6，v7
	Network    string `json:"network"`     // network
	Address    string `json:"address"`     // address
	BakAddress string `json:"bak_address"` // backup address
}

type AddDBResponse struct {
	ID int `json:"id"` // 仓库 id
}

type MaintainDBManagerRequest struct {
	DBId    int      `json:"db_id"`   // db id
	Manager []uint64 `json:"manager"` // 管理员
}

type UpdateDBBaseRequest struct {
	DBId  int    `json:"db_id"` // 仓库id
	Name  string `json:"name"`  // 数据库名称
	Intro string `json:"intro"` // 中文简介
	Desc  string `json:"desc"`  // 详细描述
}

type UpdateDBStatusRequest struct {
	DBId   int  `json:"db_id"`  // db id
	Status int8 `json:"status"` // 状态 1-在线 2-离线
}

type UpdateDBNetworkRequest struct {
	DBId int `json:"db_id"` // db id

	// db address
	Type       int    `json:"type"`        // 数据库类型 0-nil（仅执行插件） 1-elastic 2-mongo 3-redis 10-mysql 11-postgresql 12-clickhouse 13-oracle 14-DB2 15-sqlite
	Version    string `json:"version"`     // 数据库版本，比如elastic v6，v7
	Network    string `json:"network"`     // network
	Address    string `json:"address"`     // address
	BakAddress string `json:"bak_address"` // backup address

	// db params
	WriteTimeout int `json:"write_timeout"` // 写超时（毫秒）
	ReadTimeout  int `json:"read_timeout"`  // 读超时（毫秒）

	// 日志
	WarnTimeout int  `json:"warn_timeout"` // 告警超时（ms），如果请求耗时超过这个时间，就会打 warning 日志
	OmitError   int8 `json:"omit_error"`   // 是否忽略 error 日志，0-否 1-是
	Debug       int8 `json:"debug"`        // 是否开启 debug 日志，正常的数据库请求也会被打印到日志，0-否 1-是，会造成海量日志，慎重开启
}

type DBIdRequest struct {
	DbID int `json:"db_id"`
}

type DBBaseResponse struct {
	Info           *DBBase      `json:"info"`            // 数据库基础信息
	IsManager      bool         `json:"is_manager"`      // 是否数据库管理员
	Creator        *UsersBase   `json:"creator"`         // creator
	Manager        []*UsersBase `json:"manager"`         // 管理员
	ProductInfo    *ProductBase `json:"product_info"`    // 所属产品基础信息
	ProductManager []*UsersBase `json:"product_manager"` // 所属产品管理员（拥有库管理员所有权限）
	Tables         []*TableBase `json:"tables"`          // 数据库下所有表
}

type DBNetworkInfoResponse struct {
	// db address
	Type       int    `json:"type"`        // 数据库类型 0-nil（仅执行插件） 1-elastic 2-mongo 3-redis 10-mysql 11-postgresql 12-clickhouse 13-oracle 14-DB2 15-sqlite
	Version    string `json:"version"`     // 数据库版本，比如elastic v6，v7
	Network    string `json:"network"`     // network
	Address    string `json:"address"`     // address
	BakAddress string `json:"bak_address"` // backup address

	// db params
	WriteTimeout int `json:"write_timeout"` // 写超时（毫秒）
	ReadTimeout  int `json:"read_timeout"`  // 读超时（毫秒）

	// 日志
	WarnTimeout int  `json:"warn_timeout"` // 告警超时（ms），如果请求耗时超过这个时间，就会打 warning 日志
	OmitError   int8 `json:"omit_error"`   // 是否忽略 error 日志，0-否 1-是
	Debug       int8 `json:"debug"`        // 是否开启 debug 日志，正常的数据库请求也会被打印到日志，0-否 1-是，会造成海量日志，慎重开启
}
