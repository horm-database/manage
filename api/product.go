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

// AddProduct 新增产品
func AddProduct(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddProductRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product name can`t be empty")
	}

	return logic.AddProduct(ctx, head.Userid, &req)
}

// UpdateProduct 产品基础信息更新
func UpdateProduct(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateProductRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 || req.Name == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product id/name can`t be empty")
	}

	return nil, logic.UpdateProduct(ctx, head.Userid, &req)
}

// UpdateProductStatus 产品状态更新
func UpdateProductStatus(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateProductStatusRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product id can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.UpdateProductStatus(ctx, head.Userid, &req)
}

// MaintainProductManager 产品管理员维护
func MaintainProductManager(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.MaintainProductManagerRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product id can`t be empty")
	}

	if len(req.Manager) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product manager can`t be empty")
	}

	return nil, logic.MaintainProductManager(ctx, head.Userid, &req)
}

// ProductList 产品列表
func ProductList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductListRequest{}
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

	return logic.ProductList(ctx, &req)
}

// ProductDetail 产品详情
func ProductDetail(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductDetailRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product_id can`t be empty")
	}

	return logic.ProductDetail(ctx, head.Userid, &req)
}

// ProductMemberList 产品成员列表
func ProductMemberList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductMemberListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product_id can`t be empty")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 20
	}

	return logic.ProductMemberList(ctx, head.Userid, &req)
}

// ProductJoinApply 申请加入产品/续期
func ProductJoinApply(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductJoinApplyRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product_id can`t be empty")
	}

	if req.Role != 0 && req.Role != consts.ProductRoleDeveloper && req.Role != consts.ProductRoleOperator {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [role] is invalid")
	}

	if req.ExpireType > consts.ExpireTypeYear || req.ExpireType < consts.ExpireTypePermanent {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [expire_type] is invalid")
	}

	return nil, logic.ProductJoinApply(ctx, head.Userid, &req)
}

// ProductApproval 产品权限审批
func ProductApproval(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductApprovalRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 || req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid/product_id can`t be empty")
	}

	if req.Status != consts.ApprovalAccess && req.Status != consts.ApprovalReject {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.ProductApproval(ctx, head.Userid, &req)
}

// ProductChangeRoleApply 申请变更角色
func ProductChangeRoleApply(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductChangeRoleApplyRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "product_id can`t be empty")
	}

	if req.Role != consts.ProductRoleDeveloper && req.Role != consts.ProductRoleOperator {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [role] is invalid")
	}

	return nil, logic.ProductChangeRoleApply(ctx, head.Userid, &req)
}

// ProductChangeRoleApproval 产品角色变更审批
func ProductChangeRoleApproval(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductApprovalRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 || req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid/product_id can`t be empty")
	}

	if req.Status != consts.ApprovalAccess && req.Status != consts.ApprovalReject {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.ProductChangeRoleApproval(ctx, head.Userid, &req)
}

// ProductMemberRemove 将指定用户移出产品
func ProductMemberRemove(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ProductMemberRemoveRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.ProductID == 0 || req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid can`t be empty")
	}

	return nil, logic.ProductMemberRemove(ctx, head.Userid, &req)
}
