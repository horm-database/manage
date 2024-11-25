package logic

import (
	"context"
	"fmt"

	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	sc "github.com/horm-database/server/consts"
	"github.com/horm-database/server/filter/conf"
	st "github.com/horm-database/server/model/table"
)

// AddFilter 新增插件
func AddFilter(ctx context.Context, userid uint64, req *pb.AddFilterRequest) (*pb.AddFilterResponse, error) {
	data := st.TblFilter{
		Name:    req.Name,
		Intro:   req.Intro,
		Version: "",
		Func:    req.Func,
		Online:  consts.StatusOnline,
		Source:  consts.FilterSourcePrivate,
		Desc:    req.Desc,
		Creator: userid,
		Manager: fmt.Sprint(userid),
	}

	if len(req.SupportTypes) > 0 {
		data.SupportTypes = types.JoinInt8(req.SupportTypes, ",")
	}

	id, err := table.AddFilter(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &pb.AddFilterResponse{ID: id}, nil
}

func UpdateFilter(ctx context.Context, userid uint64, req *pb.UpdateFilterRequest) error {
	update := horm.Map{
		"name":          req.Name,
		"intro":         req.Intro,
		"support_types": types.JoinInt8(req.SupportTypes, ","),
		"desc":          req.Desc,
	}

	return table.UpdateFilterByID(ctx, req.FilterID, update)
}

// ReplaceFilterConfig 新增/修改插件配置
func ReplaceFilterConfig(ctx context.Context, userid uint64, req *pb.ReplaceFilterConfigRequest) error {
	data := st.TblFilterConfig{
		FilterID:      req.FilterID,
		FilterVersion: req.FilterVersion,
		Key:           req.Key,
		Name:          req.Name,
		Type:          req.Type,
		NotNull:       req.NotNull,
		MoreInfo:      req.MoreInfo,
		Default:       req.Default,
		Desc:          req.Desc,
	}

	return table.ReplaceFilterConfig(ctx, &data)
}

// DelFilterConfig 删除插件配置
func DelFilterConfig(ctx context.Context, userid uint64, req *pb.DelFilterConfigRequest) error {
	return table.DelFilterConfigByKey(ctx, req.FilterID, req.FilterVersion, req.Key)
}

func FilterList(ctx context.Context, req *pb.FilterListRequest) (*pb.FilterListResponse, error) {
	pageInfo, filters, err := table.GetFilterList(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.FilterListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      req.Page,
		Size:      req.Size,
		Filters:   make([]*pb.FilterBase, len(filters)),
	}

	var userIds []uint64
	for _, v := range filters {
		userIds = append(userIds, GetUserIds(v.Creator, v.Manager)...)
	}

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for k, v := range filters {
		ret.Filters[k] = &pb.FilterBase{
			Id:           v.Id,
			Name:         v.Name,
			Intro:        v.Intro,
			Version:      types.SplitInt(v.Version, ","),
			Func:         v.Func,
			SupportTypes: FilterTypes(v.SupportTypes),
			Online:       v.Online,
			Source:       v.Source,
			Desc:         v.Desc,
			Creator:      userMaps[v.Creator],
			Manager:      GetUsersFromMap(userMaps, GetUserIds(v.Manager)),
			CreatedAt:    v.CreatedAt.Unix(),
			UpdatedAt:    v.UpdatedAt.Unix(),
		}
	}

	return &ret, nil
}

// FilterConfigs 插件配置列表
func FilterConfigs(ctx context.Context, req *pb.FilterConfigsRequest) (*pb.FilterConfigsResponse, error) {
	filterConfigs, err := table.GetFilterConfigs(ctx, req.FilterID, req.FilterVersion)
	if err != nil {
		return nil, err
	}

	ret := pb.FilterConfigsResponse{
		Configs: make([]*pb.FilterConfig, len(filterConfigs)),
	}

	for k, v := range filterConfigs {
		ret.Configs[k] = &pb.FilterConfig{
			Id:       v.Id,
			Key:      v.Key,
			Name:     v.Name,
			Type:     v.Type,
			NotNull:  v.NotNull,
			MoreInfo: v.MoreInfo,
			Default:  v.Default,
			Desc:     v.Desc,
		}
	}

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func FiltersToMap(filters []*st.TblFilter) map[int]*st.TblFilter {
	ret := map[int]*st.TblFilter{}
	for _, v := range filters {
		ret[v.Id] = v
	}
	return ret
}

func FilterConfigsToMap(filterConfigs []*st.TblFilterConfig) map[string][]*st.TblFilterConfig {
	ret := map[string][]*st.TblFilterConfig{}
	for _, v := range filterConfigs {
		key := fmt.Sprintf("%d_%d", v.FilterID, v.FilterVersion)
		ret[key] = append(ret[key], v)
	}
	return ret
}

func GetDefaultScheduleConfig() *conf.ScheduleConfig {
	return &conf.ScheduleConfig{
		Async:         false,
		SkipError:     false,
		Timeout:       1000,
		RequestSource: []string{"api", "web"},
		OpType:        []string{"read", "mod", "del"},
		GrayScale:     100,
		AppRule: &conf.AppRule{
			ActType: sc.ActionTypeExec,
			AppIDs:  []uint64{},
		},
		CustomRule: &conf.CustomRule{
			ActType:  sc.ActionTypeExec,
			RuleType: sc.CondTypeAny,
			Rules:    []*conf.Rule{},
		},
	}
}

func FilterTypes(typ string) []int8 {
	if typ == "" {
		return []int8{sc.PreFilter, sc.PostFilter, sc.DeferFilter}
	}

	return types.SplitInt8(typ, ",")
}

func FilterTypeDesc(typ int8) string {
	switch typ {
	case sc.PreFilter:
		return "pre-filter"
	case sc.PostFilter:
		return "post-filter"
	default:
		return "defer-filter"
	}
}
