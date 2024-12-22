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

// WorkspaceBaseInfoRequest 空间基础信息
type WorkspaceBaseInfoRequest struct {
	Workspace string `json:"workspace"` // workspace
}

// WorkspaceBaseInfoResponse 空间基础信息
type WorkspaceBaseInfoResponse struct {
	Id         int          `json:"workspace_id"` // 空间 id
	Workspace  string       `json:"workspace"`    // workspace
	Name       string       `json:"name"`         // 空间名
	Intro      string       `json:"intro"`        // 简介
	Company    string       `json:"company"`      // 公司
	Department string       `json:"department"`   // 部门
	Creator    *UsersBase   `json:"creator"`      // creator
	Manager    []*UsersBase `json:"manager"`      // 管理员
	CreateTime int64        `json:"create_time"`  // 创建时间戳
	Role       int8         `json:"role"`         // 我的空间成员角色 0-非空间成员 1-空间成员 2-空间管理员 3-权限已过期
	Status     int8         `json:"status"`       // 我的空间成员状态 1-待审批 2-续期审批 3-暂未申请 4-已加入 5-审批拒绝  6-已退出
	ExpireTime int          `json:"expire_time"`  // 我的空间成员过期时间
	OutTime    int          `json:"out_time"`     // 我的空间成员退出时间
}

// WorkspaceMemberListRequest 工作空间成员列表
type WorkspaceMemberListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

// WorkspaceMemberListResponse 工作空间成员列表
type WorkspaceMemberListResponse struct {
	Total     uint64             `json:"total"`      // 总数
	TotalPage uint32             `json:"total_page"` // 总页数
	Page      int                `json:"page"`       // 分页
	Size      int                `json:"size"`       // 每页大小
	IsManager bool               `json:"is_manager"` // 是否管理员
	Members   []*WorkspaceMember `json:"members"`    // 用户列表
}

type WorkspaceMember struct {
	MemberID   int    `json:"member_id"`   // member id
	Userid     uint64 `json:"userid"`      // userid
	Account    string `json:"account"`     // 账号
	Nickname   string `json:"nickname"`    // 昵称
	Role       int8   `json:"role"`        // 角色 0:- 1:普通成员 2:管理员
	Status     int8   `json:"status"`      // 状态 1-待审批 2-续期审批 3-未加入 4-正常 5-审批拒绝  6-已退出 9-已过期
	JoinTime   int    `json:"join_time"`   // 申请/加入时间
	ExpireType int8   `json:"expire_type"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	ExpireTime int    `json:"expire_time"` // 过期时间
	OutTime    int    `json:"out_time"`    // 退出时间
}

type WorkspaceMemberInviteRequest struct {
	Userid     uint64 `json:"userid"`      // userid
	ExpireType int8   `json:"expire_type"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	Reason     string `json:"reason"`      // 邀请理由
}

type WorkspaceJoinApplyRequest struct {
	ExpireType int8   `json:"expire_type"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	Reason     string `json:"reason"`      // 申请理由
}

type WorkspaceApprovalRequest struct {
	Userid uint64 `json:"userid"` // 待审批 userid
	Status int8   `json:"status"` // 1-审批通过 2-审批拒绝
	Reason string `json:"reason"` // 拒绝理由（ status=2 时输入）
}

type WorkspaceMemberRemoveRequest struct {
	Userid uint64 `json:"userid"` // userid
	Reason string `json:"reason"` // 移除理由
}

type MaintainWorkspaceManagerRequest struct {
	Manager []uint64 `json:"manager"` // 管理员
}
