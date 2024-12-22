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

// AppCanAccessTableRequest 我的能接入指定仓库的所有应用
type AppCanAccessTableRequest struct {
	TableID int    `json:"table_id"` // 数据库
	Keyword string `json:"keyword"`  // 过滤关键词
}

type AppCanAccessTableResponse struct {
	Apps []*AppCanAccessTable `json:"apps"` // 我的能接入仓库的应用
}

type AppCanAccessTable struct {
	Appid        uint64           `json:"appid"`         // 应用appid
	AppName      string           `json:"app_name"`      // 应用名称
	Intro        string           `json:"intro"`         // 简介
	AccessStatus int8             `json:"access_status"` // 接入状态：0-未接入 1-正常 2-下线 3-审核中 4-审核撤回 5-拒绝 11-拥有库数据权限（无需申请表权限）
	AccessInfo   *TableAccessInfo `json:"access_info"`   // 接入信息
}

type TableAccessInfo struct {
	AccessID       int      `json:"access_id"`        // 接入 ID
	AccessQueryAll int8     `json:"access_query_all"` // 是否支持所有的 query 语句，1-true 2-false
	AccessOp       []string `json:"access_op"`        // 支持的操作
	AccessReason   string   `json:"access_reason"`    // 接入原因
}

// AppApplyAccessTableRequest 应用申请接入表数据
type AppApplyAccessTableRequest struct {
	Appid    uint64   `json:"appid"`     // 应用appid
	TableID  int      `json:"table_id"`  // 表ID
	QueryAll int8     `json:"query_all"` // 是否支持所有的 query 语句，1-true 2-false
	Op       []string `json:"op"`        // 支持的操作
	Reason   string   `json:"reason"`    // 接入原因
}

// AppAccessTableApprovalRequest 应用接入表数据审批
type AppAccessTableApprovalRequest struct {
	Appid   uint64 `json:"appid"`    // 应用appid
	TableID int    `json:"table_id"` // 表ID
	Status  int8   `json:"status"`   // 1-审批通过 2-审批拒绝
	Reason  string `json:"reason"`   // 拒绝理由（ status=2 时输入）
}

// AppAccessTableWithdrawRequest 应用接入表数据撤销申请
type AppAccessTableWithdrawRequest struct {
	Appid   uint64 `json:"appid"`    // 应用appid
	TableID int    `json:"table_id"` // 表ID
	Reason  string `json:"reason"`   // 撤销理由
}

// AppAccessTableUpdateRequest 编辑表数据访问权限
type AppAccessTableUpdateRequest struct {
	Appid    uint64   `json:"appid"`     // 应用appid
	TableID  int      `json:"table_id"`  // 表ID
	QueryAll int      `json:"query_all"` // 是否支持所有的 query 语句，1-true 2-false
	Op       []string `json:"op"`        // 支持的操作
	Reason   string   `json:"reason"`    // 编辑原因
}

// AppAccessTableOnOffRequest 表数据访问权限上/下线
type AppAccessTableOnOffRequest struct {
	Appid   uint64 `json:"appid"`    // 应用appid
	TableID int    `json:"table_id"` // 表ID
	Status  int8   `json:"status"`   // 状态：1-上线 2-下线
	Reason  string `json:"reason"`   // 上/下线原因
}

// TablesAllAppAccessListRequest 访问该表的应用列表
type TablesAllAppAccessListRequest struct {
	TableID int `json:"table_id"` // 表ID
	Page    int `json:"page"`     // 分页
	Size    int `json:"size"`     // 每页大小
}

type TablesAllAppAccessListResponse struct {
	Total           uint64            `json:"total"`             // 总数
	TotalPage       uint32            `json:"total_page"`        // 总页数
	Page            int               `json:"page"`              // 分页
	Size            int               `json:"size"`              // 每页大小
	IsManager       bool              `json:"is_manager"`        // 是否表管理员
	AppAccessTables []*AppAccessTable `json:"app_access_tables"` // 访问列表
}

// AppsAllTableAccessListRequest 该应用访问的表列表
type AppsAllTableAccessListRequest struct {
	Appid uint64 `json:"appid"` // 应用id
	Page  int    `json:"page"`  // 分页
	Size  int    `json:"size"`  // 每页大小
}

type AppsAllTableAccessListResponse struct {
	Total           uint64            `json:"total"`             // 总数
	TotalPage       uint32            `json:"total_page"`        // 总页数
	Page            int               `json:"page"`              // 分页
	Size            int               `json:"size"`              // 每页大小
	AppAccessTables []*AppAccessTable `json:"app_access_tables"` // 访问列表
}

type AppAccessTable struct {
	Id        int        `json:"id"`
	App       *AppBase   `json:"app,omitempty"`   // 应用信息
	Table     *TableBase `json:"table,omitempty"` // 表信息
	QueryAll  int8       `json:"query_all"`       // 是否支持所有的 query 语句，1-true 2-false
	Op        []string   `json:"op"`              // 支持的操作
	Status    int8       `json:"status"`          // 状态：1-正常 2-下线 3-审核中 4-审核撤回 5-拒绝
	ApplyUser *UsersBase `json:"apply_user"`      // 申请者
	Reason    string     `json:"reason"`          // 接入原因
	CreatedAt int64      `json:"create_time"`     // 记录创建时间
	UpdatedAt int64      `json:"update_time"`     // 最后更新时间
}
