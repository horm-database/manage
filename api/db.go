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

	cc "github.com/horm-database/common/consts"
	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// AddDB 新增数据库
func AddDB(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddDBRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" || req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "name/product_id can`t be empty")
	}

	return logic.AddDB(ctx, head.Userid, &req)
}

// UpdateDBBase 仓库基础信息更新
func UpdateDBBase(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateDBBaseRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DBId == 0 || req.Name == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "name/db_id can`t be empty")
	}

	return nil, logic.UpdateDBBase(ctx, head.Userid, &req)
}

// UpdateDBStatus 仓库状态更新
func UpdateDBStatus(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateDBStatusRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DBId == 1 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.UpdateDBStatus(ctx, head.Userid, &req)
}

// MaintainDBManager 仓库管理员维护
func MaintainDBManager(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.MaintainDBManagerRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DBId == 1 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	if len(req.Manager) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db manager can`t be empty")
	}

	return nil, logic.MaintainDBManager(ctx, head.Userid, &req)
}

// UpdateDBNetwork 数据库网络信息更新
func UpdateDBNetwork(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateDBNetworkRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DBId == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	if req.OmitError != cc.FALSE && req.OmitError != cc.TRUE {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [omit_error] is invalid")
	}

	if req.Debug != cc.FALSE && req.Debug != cc.TRUE {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [debug] is invalid")
	}

	return nil, logic.UpdateDBNetwork(ctx, head.Userid, &req)
}

// DBBase 数据库基础信息
func DBBase(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DBIdRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	return logic.DBBase(ctx, head.Userid, req.DbID)
}

// DBNetworkDetail 数据库网络配置
func DBNetworkDetail(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DBIdRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	return logic.DBNetworkDetail(ctx, head.Userid, req.DbID)
}
