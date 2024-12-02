package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
	"github.com/horm-database/server/consts"
)

// AddTableFilter 新增表插件
func AddTableFilter(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddTableFilterRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableId == 0 || req.FilterId == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id/filter_id can`t be empty")
	}

	if req.Type != consts.PreFilter && req.Type != consts.PostFilter && req.Type != consts.DeferFilter {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [type] is invalid")
	}

	return logic.AddTableFilter(ctx, head.Userid, &req)
}

// UpdateTableFilter 修改表插件
func UpdateTableFilter(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateTableFilterRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	return nil, logic.UpdateTableFilter(ctx, head.Userid, &req)
}

// DelTableFilter 删除表插件
func DelTableFilter(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DelTableFilterRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Id == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "id can`t be empty")
	}

	return nil, logic.DelTableFilter(ctx, head.Userid, req.Id)
}

// TableFilters 表插件
func TableFilters(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TableIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id can`t be empty")
	}

	return logic.TableFilters(ctx, head.Userid, req.TableID)
}
