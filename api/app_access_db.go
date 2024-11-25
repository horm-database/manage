package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
	sc "github.com/horm-database/server/consts"
)

// DBSupportOps 数据库所支持的所有操作
func DBSupportOps(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DBIdRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	return logic.DBSupportOps(ctx, head.Userid, req.DbID)
}

// AppCanAccessDB 我的能接入指定仓库的所有应用
func AppCanAccessDB(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppCanAccessDBRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	return logic.AppCanAccessDB(ctx, head.Userid, &req)
}

// AppApplyAccessDB 应用申请接入仓库
func AppApplyAccessDB(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppApplyAccessDBRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/db_id can`t be empty")
	}

	if req.Root != sc.DBRootAll && req.Root != sc.DBRootTableData && req.Root != sc.DBRootNone {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [root] is invalid")
	}

	return logic.AppApplyAccessDB(ctx, head.Userid, &req)
}

// AppAccessDBApproval 应用接入仓库审批
func AppAccessDBApproval(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessDBApprovalRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/db_id can`t be empty")
	}

	if req.Status != consts.ApprovalAccess && req.Status != consts.ApprovalReject {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.AppAccessDBApproval(ctx, head.Userid, &req)
}

// AppAccessDBWithdraw 应用接入仓库撤销申请
func AppAccessDBWithdraw(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessDBWithdrawRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/db_id can`t be empty")
	}

	return nil, logic.AppAccessDBWithdraw(ctx, head.Userid, &req)
}

// AppAccessDBUpdate 编辑仓库访问权限
func AppAccessDBUpdate(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessDBUpdateRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/db_id can`t be empty")
	}

	return nil, logic.AppAccessDBUpdate(ctx, head.Userid, &req)
}

// AppAccessDBOnOff 仓库访问权限上/下线
func AppAccessDBOnOff(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessDBOnOffRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/db_id can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.AppAccessDBOnOff(ctx, head.Userid, &req)
}

// DBsAllAppAccessList 访问该仓库的应用列表
func DBsAllAppAccessList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DBsAllAppAccessListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.DbID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "db id can`t be empty")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 20
	}

	return logic.DBsAllAppAccessList(ctx, head.Userid, &req)
}

// AppsAllDBAccessList 该应用访问的仓库列表
func AppsAllDBAccessList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppsAllDBAccessListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid can`t be empty")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 20
	}

	return logic.AppsAllDBAccessList(ctx, head.Userid, &req)
}
