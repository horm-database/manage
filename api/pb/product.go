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

type AddProductRequest struct {
	Name    string   `json:"name"`    // 产品名称
	Intro   string   `json:"intro"`   // 简介
	Manager []uint64 `json:"manager"` // 管理员
}

type AddProductResponse struct {
	ID int `json:"id"` // product id
}

type UpdateProductRequest struct {
	ProductID int    `json:"product_id"` // 产品id
	Name      string `json:"name"`       // 产品名称
	Intro     string `json:"intro"`      // 简介
}

type UpdateProductStatusRequest struct {
	ProductID int  `json:"product_id"` // 产品id
	Status    int8 `json:"status"`     // 状态 1-在线 2-离线
}

type MaintainProductManagerRequest struct {
	ProductID int      `json:"product_id"` // 产品id
	Manager   []uint64 `json:"manager"`    // 管理员
}

type ProductListRequest struct {
	Status int8 `json:"status"` // 状态 0-全部 1-在线 2-离线
	Page   int  `json:"page"`   // 分页
	Size   int  `json:"size"`   // 每页大小
}

type ProductListResponse struct {
	Total     uint64         `json:"total"`      // 总数
	TotalPage uint32         `json:"total_page"` // 总页数
	Page      int            `json:"page"`       // 分页
	Size      int            `json:"size"`       // 每页大小
	Products  []*ProductBase `json:"products"`   // 产品列表
}

// ProductBase 产品基础信息
type ProductBase struct {
	Id        int    `json:"id"`          // id
	Name      string `json:"name"`        // 产品名称
	Intro     string `json:"intro"`       // 简介
	Status    int8   `json:"status"`      // 1-正常 2-下线
	CreatedAt int64  `json:"create_time"` // 记录创建时间
}

type ProductDetailRequest struct {
	ProductID int `json:"product_id"` // 产品 id
}

type ProductDetailResponse struct {
	Info       *ProductBase `json:"info"`        // 基础信息
	Creator    *UsersBase   `json:"creator"`     // creator
	Manager    []*UsersBase `json:"manager"`     // 管理员
	Role       int8         `json:"role"`        // 我的成员角色: 0-非产品成员 1-管理员 2-开发者 3-运营者 4-成员权限已过期
	Status     int8         `json:"status"`      // 我的成员状态 1-待审批 2-续期审批 3-角色变更审批 4-已加入 5-审批拒绝  6-已退出
	ExpireTime int          `json:"expire_time"` // 我的成员过期时间
	OutTime    int          `json:"out_time"`    // 我的成员退出时间
	ChangeRole int8         `json:"change_role"` // 变更为目标角色 0-无 2-开发者 3-运营者
	DBs        []*DBBase    `json:"dbs"`         // 产品下数据库
}

type ProductMemberListRequest struct {
	ProductID int `json:"product_id"` // 产品 id
	Page      int `json:"page"`       // 分页
	Size      int `json:"size"`       // 每页大小
}

type ProductMemberListResponse struct {
	Total      uint64           `json:"total"`       // 总数
	TotalPage  uint32           `json:"total_page"`  // 总页数
	Page       int              `json:"page"`        // 分页
	Size       int              `json:"size"`        // 每页大小
	Role       int8             `json:"role"`        // 我的成员角色: 0-非产品成员 1-管理员 2-开发者 3-运营者 4-成员权限已过期
	Status     int8             `json:"status"`      // 我的成员状态 1-待审批 2-续期审批 3-角色变更审批 4-正常 5-审批拒绝  6-已退出
	ChangeRole int8             `json:"change_role"` // 变更为目标角色 0-无 2-开发者 3-运营者
	Members    []*ProductMember `json:"members"`     // 成员列表
}

type ProductMember struct {
	MemberID   int    `json:"member_id"`   // member id
	Userid     uint64 `json:"userid"`      // userid
	Account    string `json:"account"`     // 账号
	Nickname   string `json:"nickname"`    // 昵称
	Role       int8   `json:"role"`        // 角色: 0:- 1:管理员 2:开发者 3:运营者
	Status     int8   `json:"status"`      // 状态 0-未加入 1-待审批 2-续期审批 3-角色变更审批 4-正常 5-审批拒绝 6-已退出 9-已过期
	JoinTime   int    `json:"join_time"`   // 申请/加入时间
	ExpireType int8   `json:"expire_type"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	ExpireTime int    `json:"expire_time"` // 过期时间
	OutTime    int    `json:"out_time"`    // 退出时间
	ChangeRole int8   `json:"change_role"` // 变更为目标角色 0-无 2-开发者 3-运营者
}

type ProductJoinApplyRequest struct {
	ProductID  int    `json:"product_id"`  // 产品 id
	Role       int8   `json:"role"`        // 0-续期（续期不能变更角色）2-开发者 3-运营者
	ExpireType int8   `json:"expire_type"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	Reason     string `json:"reason"`      // 申请理由
}

type ProductApprovalRequest struct {
	ProductID int    `json:"product_id"` // 产品 id
	Userid    uint64 `json:"userid"`     // 待审批 userid
	Status    int8   `json:"status"`     // 1-审批通过 2-审批拒绝
	Reason    string `json:"reason"`     // 拒绝理由（ status=2 时输入）
}

type ProductChangeRoleApplyRequest struct {
	ProductID int    `json:"product_id"` // 产品 id
	Role      int8   `json:"role"`       // 2-开发者 3-运营者
	Reason    string `json:"reason"`     // 变更理由
}

type ProductMemberRemoveRequest struct {
	ProductID int    `json:"product_id"` // 产品 id
	Userid    uint64 `json:"userid"`     // userid
	Reason    string `json:"reason"`     // 移除理由
}
