package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/json"
	"github.com/horm-database/common/types"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/auth"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// WorkspaceBaseInfo 工作空间基础信息
func WorkspaceBaseInfo(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceBaseInfoRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.RetServerDecodeFail,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if req.Workspace == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "workspace can`t be empty")
	}

	return logic.WorkspaceBaseInfo(ctx, head.Userid, &req)
}

// WorkspaceJoinApply 申请加入工作空间/续期
func WorkspaceJoinApply(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceJoinApplyRequest{}
	err := json.Api.Unmarshal(reqBuf, &req)
	if err != nil {
		return nil, errs.Newf(errs.RetServerDecodeFail,
			"decode request error: %v, request:[%s]", err, types.QuickReplaceLFCR2Space(reqBuf))
	}

	if head.Userid == 0 {
		return nil, errs.New(errs.RetWebNotLogin, "please login first")
	}

	notFind, user, err := table.GetUserByID(ctx, head.Userid)
	if err != nil {
		return nil, err
	}

	if notFind {
		return nil, errs.New(errs.RetWebNotFindUser, "not find user")

	}

	if !auth.SignSuccess(head, user.Token) {
		//return errs.Newf(errs.RetServerAuthFail, "signature failed")
	}

	if head.Userid == 0 {
		return nil, errs.New(errs.RetWebNotLogin, "please login first")
	}

	if req.ExpireType > consts.ExpireTypeYear || req.ExpireType < consts.ExpireTypePermanent {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [expire_type] is invalid")
	}

	return nil, logic.WorkspaceJoinApply(ctx, head.Userid, int(head.WorkspaceId), &req)
}

// WorkspaceApproval 空间权限审批
func WorkspaceApproval(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceApprovalRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid can`t be empty")
	}

	if req.Status != consts.ApprovalAccess && req.Status != consts.ApprovalReject {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.WorkspaceApproval(ctx, head.Userid, int(head.WorkspaceId), &req)
}

// WorkspaceMemberInvite 管理员邀请用户加入空间
func WorkspaceMemberInvite(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceMemberInviteRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid can`t be empty")
	}

	if req.ExpireType > consts.ExpireTypeYear || req.ExpireType < consts.ExpireTypePermanent {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [expire_type] is invalid")
	}

	return nil, logic.WorkspaceMemberInvite(ctx, head.Userid, int(head.WorkspaceId), &req)
}

// WorkspaceMemberRemove 将指定用户移出空间
func WorkspaceMemberRemove(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceMemberRemoveRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Userid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "userid can`t be empty")
	}

	return nil, logic.WorkspaceMemberRemove(ctx, head.Userid, int(head.WorkspaceId), &req)
}

// WorkspaceMemberList 工作空间成员列表
func WorkspaceMemberList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.WorkspaceMemberListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	return logic.WorkspaceMemberList(ctx, head.Userid, int(head.WorkspaceId), &req)
}

// MaintainWorkspaceManager 空间管理员维护
func MaintainWorkspaceManager(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.MaintainWorkspaceManagerRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if len(req.Manager) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "workspace manager can`t be empty")
	}

	return nil, logic.MaintainWorkspaceManager(ctx, head.Userid, int(head.WorkspaceId), req.Manager)
}
