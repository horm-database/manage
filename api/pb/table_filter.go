package pb

import (
	"github.com/horm-database/server/filter/conf"
)

type TableFiltersResponse struct {
	PreFilters   []*TableFilter `json:"pre_filters"`   // 前置插件
	PostFilters  []*TableFilter `json:"post_filters"`  // 后置插件
	DeferFilters []*TableFilter `json:"defer_filters"` // 延迟插件
}

type AddTableFilterRequest struct {
	TableId        int                    `json:"table_id"`        // 表id
	FilterId       int                    `json:"filter_id"`       // 插件id
	FilterVersion  int                    `json:"filter_version"`  // filter 版本
	Type           int8                   `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                    `json:"front"`           // filter execute front of me
	Desc           string                 `json:"desc"`            // 描述
	ScheduleConfig *conf.ScheduleConfig   `json:"schedule_config"` // 插件调度配置
	FilterConfigs  map[string]interface{} `json:"filter_configs"`  // 插件配置
}

type AddTableFilterResponse struct {
	ID int `json:"id"`
}

type UpdateTableFilterRequest struct {
	Id             int                    `json:"id"`
	FilterVersion  int                    `json:"filter_version"`  // filter 版本
	Type           int8                   `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                    `json:"front"`           // filter execute front of me
	Desc           string                 `json:"desc"`            // 描述
	ScheduleConfig *conf.ScheduleConfig   `json:"schedule_config"` // 插件调度配置
	FilterConfigs  map[string]interface{} `json:"filter_configs"`  // 插件配置
}

type DelTableFilterRequest struct {
	Id int `json:"id"`
}

type TableFilter struct {
	Id             int                  `json:"id"`              // id
	TableId        int                  `json:"table_id"`        // 表id
	FilterId       int                  `json:"filter_id"`       // 插件id
	FilterVersion  int                  `json:"filter_version"`  // filter 版本
	Type           int8                 `json:"type"`            // 过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Front          int                  `json:"front"`           // filter execute front of me
	Desc           string               `json:"desc"`            // 描述
	Status         int8                 `json:"status"`          // 状态 1-启用 2-停用
	CreatedAt      int64                `json:"create_time"`     // 添加时间
	UpdatedAt      int64                `json:"update_time"`     // 最后修改时间
	ScheduleConfig *conf.ScheduleConfig `json:"schedule_config"` // 插件调度配置
	FilterInfo     *FilterBase          `json:"filter_info"`     // 插件信息
	FilterConfigs  []*TableFilterConfig `json:"filter_configs"`  // 插件配置
}

type TableFilterConfig struct {
	Config *FilterConfig `json:"config"` // 插件配置
	IsSet  bool          `json:"is_set"` // 是否已经设置
	Value  interface{}   `json:"value"`  // 设置值
}
