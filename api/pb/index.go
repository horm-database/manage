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

package pb

type IndexTableListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type IndexTableListResponse struct {
	Total     uint64       `json:"total"`      // 总数
	TotalPage uint32       `json:"total_page"` // 总页数
	Page      int          `json:"page"`       // 分页
	Size      int          `json:"size"`       // 每页大小
	Tables    []*TableBase `json:"tables"`     // 用户列表
}

type CollectTableListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type CollectTableRequest struct {
	TableID int `json:"table_id"` // 表id
	Status  int `json:"status"`   // 1-收藏 2-取关
}
