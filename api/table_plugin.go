package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
	"github.com/horm-database/server/consts"
)

// AddTablePlugin 新增表插件
func AddTablePlugin(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddTablePluginRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableId == 0 || req.PluginID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id/plugin_id can`t be empty")
	}

	if req.Type != consts.PrePlugin && req.Type != consts.PostPlugin && req.Type != consts.DeferPlugin {
		return nil, errs.Newf(errs.RetWebParamEmpty, "input param [type] is invalid")
	}

	return logic.AddTablePlugin(ctx, head.Userid, &req)
}

// UpdateTablePlugin 修改表插件
func UpdateTablePlugin(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateTablePluginRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	return nil, logic.UpdateTablePlugin(ctx, head.Userid, &req)
}

// DelTablePlugin 删除表插件
func DelTablePlugin(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DelTablePluginRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Id == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "id can`t be empty")
	}

	return nil, logic.DelTablePlugin(ctx, head.Userid, req.Id)
}

// TablePlugins 表插件
func TablePlugins(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.TableIDRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.TableID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "table_id can`t be empty")
	}

	return logic.TablePlugins(ctx, head.Userid, req.TableID)
}