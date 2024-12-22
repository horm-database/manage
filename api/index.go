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

// IndexTableList 首页所有的表
func IndexTableList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.IndexTableListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 30
	}

	return logic.IndexTableList(ctx, &req)
}

// CollectTableList 收藏的表
func CollectTableList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.CollectTableListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 30
	}

	return logic.CollectTableList(ctx, head.Userid, &req)
}

// CollectTable 收藏表
func CollectTable(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.CollectTableRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Status != consts.CollectTableAdd && req.Status != consts.CollectTableDel {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.CollectTable(ctx, head.Userid, &req)
}
