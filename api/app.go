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
package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// AddApp 新增应用
func AddApp(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddAppRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "app name can`t be empty")
	}

	return logic.AddApp(ctx, head.Userid, &req)
}

// UpdateApp 应用基础信息更新
func UpdateApp(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateAppRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.Name == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/name can`t be empty")
	}

	return nil, logic.UpdateApp(ctx, head.Userid, &req)
}

// ResetAppSecret 重置应用秘钥
func ResetAppSecret(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid can`t be empty")
	}

	return logic.ResetAppSecret(ctx, head.Userid, req.Appid)
}

// UpdateAppStatus 应用状态更新
func UpdateAppStatus(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateAppStatusRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.UpdateAppStatus(ctx, head.Userid, &req)
}

// MaintainAppManager 应用管理员维护
func MaintainAppManager(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.MaintainAppManagerRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.AppID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid can`t be empty")
	}

	if len(req.Manager) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "app manager can`t be empty")
	}

	return nil, logic.MaintainAppManager(ctx, head.Userid, &req)
}

// AppList 应用列表
func AppList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 20
	}

	return logic.AppList(ctx, head.Userid, &req)
}

// AppDetail 应用详情
func AppDetail(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid can`t be empty")
	}

	return logic.AppDetail(ctx, head.Userid, req.Appid)
}
