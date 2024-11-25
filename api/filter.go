package api

import (
	"context"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/logic"
	"github.com/horm-database/manage/srv/transport/web/head"
)

// AddFilter 新增插件
func AddFilter(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.AddFilterRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.Name == "" || len(req.SupportTypes) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "name/support_types can`t be empty")
	}

	return logic.AddFilter(ctx, head.Userid, &req)
}

// UpdateFilter 更新插件
func UpdateFilter(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.UpdateFilterRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.FilterID == 0 || req.Name == "" || len(req.SupportTypes) == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "filter_id/name/support_types can`t be empty")
	}

	return nil, logic.UpdateFilter(ctx, head.Userid, &req)
}

// ReplaceFilterConfig 新增/修改插件配置
func ReplaceFilterConfig(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.ReplaceFilterConfigRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.FilterID == 0 || req.Key == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "filter_id/key can`t be empty")
	}

	return nil, logic.ReplaceFilterConfig(ctx, head.Userid, &req)
}

// DelFilterConfig 删除插件配置
func DelFilterConfig(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.DelFilterConfigRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.FilterID == 0 || req.Key == "" {
		return nil, errs.Newf(errs.RetWebParamEmpty, "filter_id/key can`t be empty")
	}

	return nil, logic.DelFilterConfig(ctx, head.Userid, &req)
}

// FilterList 插件列表
func FilterList(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.FilterListRequest{}
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

	return logic.FilterList(ctx, &req)
}

// FilterConfigs 插件配置列表
func FilterConfigs(ctx context.Context, head *head.WebReqHeader, reqBuf []byte) (interface{}, error) {
	req := pb.FilterConfigsRequest{}
	err := DecodeAndAuth(ctx, head, reqBuf, &req)
	if err != nil {
		return nil, err
	}

	if req.FilterID == 0 {
		return nil, errs.Newf(errs.RetWebParamEmpty, "filter id can`t be empty")
	}

	return logic.FilterConfigs(ctx, &req)
}
