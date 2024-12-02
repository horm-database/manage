package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/json"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	cc "github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/server/consts"
	"github.com/horm-database/server/filter/conf"
	st "github.com/horm-database/server/model/table"
	"github.com/samber/lo"
)

// TableFilters 表插件
func TableFilters(ctx context.Context, userid uint64, tableID int) (*pb.TableFiltersResponse, error) {
	_, _, err := IsTableManager(ctx, userid, tableID)
	if err != nil {
		return nil, err
	}

	tableFilters, err := table.GetTableFilters(ctx, tableID)
	if err != nil {
		return nil, err
	}

	ret := pb.TableFiltersResponse{
		PreFilters:   []*pb.TableFilter{},
		PostFilters:  []*pb.TableFilter{},
		DeferFilters: []*pb.TableFilter{},
	}

	if len(tableFilters) == 0 {
		return &ret, nil
	}

	var filterIDs []int
	for _, tableFilter := range tableFilters {
		filterIDs = append(filterIDs, tableFilter.FilterId)
	}

	filterInfos, err := table.GetFilterByIDs(ctx, filterIDs)
	if err != nil {
		return nil, err
	}

	filterInfoMaps := FiltersToMap(filterInfos)

	var idVersions []int
	for _, tableFilter := range tableFilters {
		idVersions = append(idVersions, tableFilter.FilterId)
		idVersions = append(idVersions, tableFilter.FilterVersion)
	}

	filterConfigs, err := table.GetConfigsByFilterIDVersions(ctx, idVersions...)
	if err != nil {
		return nil, err
	}

	filterConfigMaps := FilterConfigsToMap(filterConfigs)

	var userIds []uint64
	for _, filterInfo := range filterInfos {
		userIds = append(userIds, GetUserIds(filterInfo.Creator, filterInfo.Manager)...)
	}

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, tableFilter := range tableFilters {
		filter := pb.TableFilter{
			Id:             tableFilter.Id,
			TableId:        tableFilter.TableId,
			FilterId:       tableFilter.FilterId,
			FilterVersion:  tableFilter.FilterVersion,
			Type:           tableFilter.Type,
			Front:          tableFilter.Front,
			Desc:           tableFilter.Desc,
			Status:         tableFilter.Status,
			CreatedAt:      tableFilter.CreatedAt.Unix(),
			UpdatedAt:      tableFilter.UpdatedAt.Unix(),
			FilterConfigs:  []*pb.TableFilterConfig{},
			ScheduleConfig: &conf.ScheduleConfig{},
		}

		filterInfo := filterInfoMaps[tableFilter.FilterId]
		if filterInfo != nil {
			filter.FilterInfo = &pb.FilterBase{
				Id:           filterInfo.Id,
				Name:         filterInfo.Name,
				Intro:        filterInfo.Intro,
				Version:      types.SplitInt(filterInfo.Version, ","),
				Func:         filterInfo.Func,
				SupportTypes: FilterTypes(filterInfo.SupportTypes),
				Online:       filterInfo.Online,
				Source:       filterInfo.Source,
				Desc:         filterInfo.Desc,
				Creator:      userMaps[filterInfo.Creator],
				Manager:      GetUsersFromMap(userMaps, GetUserIds(filterInfo.Manager)),
				CreatedAt:    filterInfo.CreatedAt.Unix(),
				UpdatedAt:    filterInfo.UpdatedAt.Unix(),
			}
		}

		if tableFilter.ScheduleConfig == "" {
			filter.ScheduleConfig = GetDefaultScheduleConfig()
		} else {
			_ = json.Api.Unmarshal([]byte(tableFilter.ScheduleConfig), &filter.ScheduleConfig)
		}

		tableFilterConfigValues := map[string]interface{}{}
		if tableFilter.Config != "" {
			_ = json.Api.Unmarshal([]byte(tableFilter.Config), &tableFilterConfigValues)
		}

		configs := filterConfigMaps[fmt.Sprintf("%d_%d", tableFilter.FilterId, tableFilter.FilterVersion)]
		if len(configs) > 0 {
			for _, conf := range configs {
				tableFilterConfig := pb.TableFilterConfig{
					Config: &pb.FilterConfig{
						Id:       conf.Id,
						Key:      conf.Key,
						Name:     conf.Name,
						Type:     conf.Type,
						NotNull:  conf.NotNull,
						MoreInfo: conf.MoreInfo,
						Default:  conf.Default,
						Desc:     conf.Desc,
					},
				}

				configValue, ok := tableFilterConfigValues[conf.Key]

				if ok {
					tableFilterConfig.IsSet = true
					tableFilterConfig.Value = configValue
				}

				filter.FilterConfigs = append(filter.FilterConfigs, &tableFilterConfig)
			}
		}

		switch tableFilter.Type {
		case consts.PreFilter:
			ret.PreFilters = append(ret.PreFilters, &filter)
		case consts.PostFilter:
			ret.PostFilters = append(ret.PostFilters, &filter)
		case consts.DeferFilter:
			ret.DeferFilters = append(ret.DeferFilters, &filter)
		}
	}

	ret.PreFilters, err = sortTableFilters("table pre-filter", ret.PreFilters)
	if err != nil {
		return nil, err
	}

	ret.PostFilters, err = sortTableFilters("table post-filter", ret.PostFilters)
	if err != nil {
		return nil, err
	}

	ret.DeferFilters, err = sortTableFilters("table defer-filter", ret.DeferFilters)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// AddTableFilter 新增表插件
func AddTableFilter(ctx context.Context, userid uint64,
	req *pb.AddTableFilterRequest) (*pb.AddTableFilterResponse, error) {
	_, _, err := IsTableManager(ctx, userid, req.TableId)
	if err != nil {
		return nil, err
	}

	isNil, filter, err := table.GetFilterByID(ctx, req.FilterId)
	if err != nil {
		return nil, err
	}

	if isNil {
		return nil, errs.Newf(errs.RetWebNotFindFilter, "not find table %d", req.FilterId)
	}

	supportVersions := types.SplitInt(filter.Version, ",")
	if lo.IndexOf(supportVersions, req.FilterVersion) == -1 {
		return nil, errs.Newf(errs.RetWebNotFindFilter,
			"filter %s not support version %d", filter.Name, req.FilterVersion)
	}

	if filter.Online != cc.StatusOnline {
		return nil, errs.Newf(errs.RetWebNotFindFilter, "filter %s is not online", filter.Name)
	}

	supportTypes := FilterTypes(filter.SupportTypes)
	if lo.IndexOf(supportTypes, req.Type) == -1 {
		return nil, errs.Newf(errs.RetWebNotFindFilter,
			"filter %s not support %s", filter.Name, FilterTypeDesc(req.Type))
	}

	tableFilters, err := table.GetTableFilters(ctx, req.TableId, req.Type)
	if err != nil {
		return nil, err
	}

	if len(tableFilters) == 0 && req.Front != 0 {
		return nil, errs.Newf(errs.RetWebIsFirstFilter, "this is first filter, front must be zero")
	}

	insertTableFilter := st.TblTableFilter{
		TableId:        req.TableId,
		FilterId:       req.FilterId,
		FilterVersion:  req.FilterVersion,
		Type:           req.Type,
		Front:          req.Front,
		ScheduleConfig: json.MarshalToString(req.ScheduleConfig),
		Config:         json.MarshalToString(req.FilterConfigs),
		Desc:           req.Desc,
		Status:         cc.StatusOnline,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	id, err := table.InsertTableFilter(ctx, &insertTableFilter)
	if err != nil {
		return nil, err
	}

	if len(tableFilters) > 0 {
		for _, tableFilter := range tableFilters {
			if tableFilter.Front == req.Front {
				_ = table.UpdateTableFilterByID(ctx, tableFilter.Id, horm.Map{"front": id})
			}
		}
	}

	return &pb.AddTableFilterResponse{ID: id}, nil
}

// UpdateTableFilter 更新表插件
func UpdateTableFilter(ctx context.Context, userid uint64, req *pb.UpdateTableFilterRequest) error {
	isNil, tableFilter, err := table.GetTableFilterByID(ctx, req.Id)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebNotFindTableFilter, "not find table filter [%d]", req.Id)
	}

	_, _, err = IsTableManager(ctx, userid, tableFilter.TableId)
	if err != nil {
		return err
	}

	isNil, filter, err := table.GetFilterByID(ctx, tableFilter.FilterId)
	if err != nil {
		return err
	}

	if isNil {
		return errs.Newf(errs.RetWebNotFindFilter, "not find table %d", tableFilter.FilterId)
	}

	supportVersions := types.SplitInt(filter.Version, ",")
	if lo.IndexOf(supportVersions, req.FilterVersion) == -1 {
		return errs.Newf(errs.RetWebNotFindFilter,
			"filter %s not support version %d", filter.Name, req.FilterVersion)
	}

	supportTypes := FilterTypes(filter.SupportTypes)
	if lo.IndexOf(supportTypes, req.Type) == -1 {
		return errs.Newf(errs.RetWebNotFindFilter,
			"filter %s not support %s", filter.Name, FilterTypeDesc(req.Type))
	}

	if req.Type != tableFilter.Type || req.Front != tableFilter.Front {
		backTableFilters, err := table.GetTableBackFilter(ctx, tableFilter.TableId, tableFilter.Type, tableFilter.Id)
		if err != nil {
			return err
		}

		tableTypeFilters, err := table.GetTableFilters(ctx, tableFilter.TableId, req.Type)
		if err != nil {
			return err
		}

		// 先剔除自身
		if len(backTableFilters) > 0 {
			var ids = getTableFiltersID(backTableFilters)
			err = table.UpdateTableFilterByIDs(ctx, ids, horm.Map{"front": tableFilter.Front})
			if err != nil {
				return err
			}
		}

		// 再插入新位置
		if len(tableTypeFilters) > 0 {
			for _, tf := range tableTypeFilters {
				if tf.Id != req.Id && tf.Front == req.Front {
					err = table.UpdateTableFilterByID(ctx, tf.Id, horm.Map{"front": req.Id})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	updateTableFilter := horm.Map{
		"filter_version":  req.FilterVersion,
		"type":            req.Type,
		"front":           req.Front,
		"schedule_config": json.MarshalToString(req.ScheduleConfig),
		"config":          json.MarshalToString(req.FilterConfigs),
		"desc":            req.Desc,
	}

	err = table.UpdateTableFilterByID(ctx, req.Id, updateTableFilter)
	if err != nil {
		return err
	}

	return nil
}

// DelTableFilter 删除表插件
func DelTableFilter(ctx context.Context, userid uint64, id int) error {
	isNil, tableFilter, err := table.GetTableFilterByID(ctx, id)
	if err != nil {
		return err
	}
	if isNil {
		return errs.Newf(errs.RetWebNotFindTableFilter, "not find table filter [%d]", id)
	}

	_, _, err = IsTableManager(ctx, userid, tableFilter.TableId)
	if err != nil {
		return err
	}

	backTableFilters, err := table.GetTableBackFilter(ctx, tableFilter.TableId, tableFilter.Type, tableFilter.Id)
	if err != nil {
		return err
	}

	if len(backTableFilters) > 0 {
		var ids = getTableFiltersID(backTableFilters)
		err = table.UpdateTableFilterByIDs(ctx, ids, horm.Map{"front": tableFilter.Front})
		if err != nil {
			return err
		}
	}

	_ = table.DelTableFilter(ctx, id)

	return nil
}

///////////////////////////////// function /////////////////////////////////////////

func sortTableFilters(typ string, tableFilters []*pb.TableFilter) ([]*pb.TableFilter, error) {
	if len(tableFilters) == 0 {
		return []*pb.TableFilter{}, nil
	}

	var head *pb.TableFilter

	for _, tableFilter := range tableFilters {
		if tableFilter.Front == 0 {
			head = tableFilter
			break
		}
	}

	if head == nil {
		return nil, errs.Newf(errs.RetFilterFrontNotFind,
			"table_id %d not find head of %s", tableFilters[0].TableId, typ)
	}

	ret := []*pb.TableFilter{}
	ret = append(ret, head)

	currentTableFilter := head
	for i := 0; i < len(tableFilters); i++ {
		backTableFilter := findBackTableFilter(currentTableFilter, tableFilters)
		if backTableFilter == nil { // 最后一个
			break
		}

		ret = append(ret, backTableFilter)
		currentTableFilter = backTableFilter
	}

	return ret, nil
}

func findBackTableFilter(currentTableFilter *pb.TableFilter, tableFilters []*pb.TableFilter) *pb.TableFilter {
	for _, tableFilter := range tableFilters {
		if tableFilter.Front == currentTableFilter.Id {
			return tableFilter
		}
	}
	return nil
}

func getTableFiltersID(tableFilters []*st.TblTableFilter) []int {
	var ids []int
	for _, v := range tableFilters {
		ids = append(ids, v.Id)
	}

	return ids
}
