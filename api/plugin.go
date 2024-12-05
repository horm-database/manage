package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// AddPlugin 新增插件
func AddPlugin(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddPluginRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" || len(req.SupportTypes) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "name/support_types can`t be empty")
	}

	return logic.AddPlugin(ctx, head.Userid, &req)
}

// UpdatePlugin 更新插件
func UpdatePlugin(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdatePluginRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.PluginID == 0 || req.Name == "" || len(req.SupportTypes) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "plugin_id/name/support_types can`t be empty")
	}

	return nil, logic.UpdatePlugin(ctx, head.Userid, &req)
}

// ReplacePluginConfig 新增/修改插件配置
func ReplacePluginConfig(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ReplacePluginConfigRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.PluginID == 0 || req.Key == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "plugin_id/key can`t be empty")
	}

	return nil, logic.ReplacePluginConfig(ctx, head.Userid, &req)
}

// DelPluginConfig 删除插件配置
func DelPluginConfig(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DelPluginConfigRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.PluginID == 0 || req.Key == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "plugin_id/key can`t be empty")
	}

	return nil, logic.DelPluginConfig(ctx, head.Userid, &req)
}

// PluginList 插件列表
func PluginList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.PluginListRequest{}
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

	return logic.PluginList(ctx, &req)
}

// PluginConfigs 插件配置列表
func PluginConfigs(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.PluginConfigsRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.PluginID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "plugin id can`t be empty")
	}

	return logic.PluginConfigs(ctx, &req)
}
