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
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// AddTable 新增表
func AddTable(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddTableRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" || req.DB == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "name/db_id can`t be empty")
	}

	return logic.AddTable(ctx, head.Userid, &req)
}

// UpdateTableBase 表基础信息更新
func UpdateTableBase(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateTableBaseRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	return nil, logic.UpdateTableBase(ctx, head.Userid, &req)
}

// UpdateTableStatus 表状态更新
func UpdateTableStatus(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateTableStatusRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 1 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.UpdateTableStatus(ctx, head.Userid, &req)
}

// UpdateTableAdvance 表高级配置更新
func UpdateTableAdvance(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateTableAdvanceRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 1 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	return nil, logic.UpdateTableAdvance(ctx, head.Userid, &req)
}

// TableDetail 表详情
func TableDetail(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TableIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id can`t be empty")
	}

	return logic.TableDetail(ctx, head.Userid, req.TableID)
}

// TableAdvanceConfig 表高级配置
func TableAdvanceConfig(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TableIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id can`t be empty")
	}

	return logic.TableAdvanceConfig(ctx, head.Userid, req.TableID)
}
