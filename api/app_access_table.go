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

// TableSupportOps 表所支持的所有操作
func TableSupportOps(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TableIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	return logic.TableSupportOps(ctx, head.Userid, req.TableID)
}

// AppCanAccessTable 我的能接入指定表的所有应用
func AppCanAccessTable(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppCanAccessTableRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	return logic.AppCanAccessTable(ctx, head.Userid, &req)
}

// AppApplyAccessTable 应用申请接入表数据
func AppApplyAccessTable(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppApplyAccessTableRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/table_id can`t be empty")
	}

	if req.QueryAll != sc.TableQueryAllTrue && req.QueryAll != sc.TableQueryAllFalse {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [query_all] is invalid")
	}

	return logic.AppApplyAccessTable(ctx, head.Userid, &req)
}

// AppAccessTableApproval 应用接入表数据审批
func AppAccessTableApproval(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessTableApprovalRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/table_id can`t be empty")
	}

	if req.Status != consts.ApprovalAccess && req.Status != consts.ApprovalReject {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.AppAccessTableApproval(ctx, head.Userid, &req)
}

// AppAccessTableWithdraw 应用接入表数据撤销申请
func AppAccessTableWithdraw(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessTableWithdrawRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/table_id can`t be empty")
	}

	return nil, logic.AppAccessTableWithdraw(ctx, head.Userid, &req)
}

// AppAccessTableUpdate 编辑表数据访问权限
func AppAccessTableUpdate(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessTableUpdateRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/table_id can`t be empty")
	}

	return nil, logic.AppAccessTableUpdate(ctx, head.Userid, &req)
}

// AppAccessTableOnOff 表数据访问权限上/下线
func AppAccessTableOnOff(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppAccessTableOnOffRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Appid == 0 || req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "appid/table_id can`t be empty")
	}

	if req.Status != consts.StatusOnline && req.Status != consts.StatusOffline {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [status] is invalid")
	}

	return nil, logic.AppAccessTableOnOff(ctx, head.Userid, &req)
}

// TablesAllAppAccessList 访问该表的应用列表
func TablesAllAppAccessList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TablesAllAppAccessListRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table id can`t be empty")
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 20
	}

	return logic.TablesAllAppAccessList(ctx, head.Userid, &req)
}

// AppsAllTableAccessList 该应用访问的表列表
func AppsAllTableAccessList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AppsAllTableAccessListRequest{}
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

	return logic.AppsAllTableAccessList(ctx, head.Userid, &req)
}
