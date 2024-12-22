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

package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/json"
	"github.com/horm-database/common/types"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// SendEmailCode 发送邮箱验证码
func SendEmailCode(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.SendEmailCodeRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.ErrServerDecode,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if req.Account == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "account can`t be empty")
	}

	return nil, logic.SendEmailCode(ctx, &req)
}

// Register 注册
func Register(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.RegisterRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.ErrServerDecode,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if req.Account == "" || req.Code == "" || req.Password == "" || req.Nickname == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "account/code/password/nickname can`t be empty")
	}

	return nil, logic.Register(ctx, &req)
}

func Login(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.LoginRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.ErrServerDecode,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if req.Account == "" || req.Password == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "account/password can`t be empty")
	}

	return logic.Login(ctx, head, &req)
}

func ResetPassword(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ResetPasswordRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.ErrServerDecode,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if req.Account == "" || req.Password == "" || req.Code == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "account/code/password can`t be empty")
	}

	return nil, logic.ResetPassword(ctx, &req)
}

func FindUser(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.FindUserRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	return logic.FindUser(ctx, req.Keyword)
}

func FindUserByID(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.FindUserByIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	return logic.FindUserByID(ctx, req.Ids)
}
