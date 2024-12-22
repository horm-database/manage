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

// UsersBase 用户基础信息
type UsersBase struct {
	UserID    uint64 `json:"userid"`     // 用户id
	Account   string `json:"account"`    // 账号
	Nickname  string `json:"nickname"`   // 昵称
	AvatarUrl string `json:"avatar_url"` // 头像
	Gender    int8   `json:"gender"`     // 性别 1-男  2-女
}

// SendEmailCodeRequest 发送验证码到邮箱请求
type SendEmailCodeRequest struct {
	Account string `json:"account"` // 账号
	Type    int    `json:"type"`    // 邮件类型 0-账号验证
}

// RegisterRequest 注册
type RegisterRequest struct {
	Account  string `json:"account"`  // 账号
	Nickname string `json:"nickname"` // 昵称
	Code     string `json:"code"`     // 验证码
	Password string `json:"password"` // 密码
}

// LoginRequest 登录
type LoginRequest struct {
	Account  string `json:"account"`  // 账号
	Password string `json:"password"` // 密码
}

// LoginResponse 登录
type LoginResponse struct {
	UserID     uint64 `json:"userid"`     // 用户id
	Account    string `json:"account"`    // 账号
	Mobile     string `json:"mobile"`     // 手机号
	Nickname   string `json:"nickname"`   // 昵称
	Token      string `json:"token"`      // token
	Gender     int8   `json:"gender"`     // 性别 1-男  2-女
	Company    string `json:"company"`    // 公司
	Department string `json:"department"` // 部门
	City       string `json:"city"`       // 城市
	Province   string `json:"province"`   // 省份
	Country    string `json:"country"`    // 国家
}

// ResetPasswordRequest 重设密码
type ResetPasswordRequest struct {
	Account  string `json:"account"`  // 账号
	Code     string `json:"code"`     // 验证码
	Password string `json:"password"` // 密码
}

// FindUserRequest 检索用户
type FindUserRequest struct {
	Keyword string `json:"keyword"` // 关键词
}

// FindUserResponse 检索用户
type FindUserResponse struct {
	Users []*UsersBase `json:"users"`
}

// FindUserByIDRequest 通过 userid 查找用户
type FindUserByIDRequest struct {
	Ids []uint64 `json:"ids"` // user ids
}
