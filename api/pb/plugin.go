package pb

type PluginListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type PluginListResponse struct {
	Total     uint64        `json:"total"`      // 总数
	TotalPage uint32        `json:"total_page"` // 总页数
	Page      int           `json:"page"`       // 分页
	Size      int           `json:"size"`       // 每页大小
	Plugins   []*PluginBase `json:"plugins"`    // 插件列表
}

// PluginBase 插件基础信息
type PluginBase struct {
	Id           int          `json:"id"`            // id
	Name         string       `json:"name"`          // 产品名称
	Intro        string       `json:"intro"`         // 中文简介
	Version      []int        `json:"version"`       // 所有支持的插件版本，逗号分开
	Func         string       `json:"func"`          // 插件注册函数
	SupportTypes []int8       `json:"support_types"` // 支持的过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器，多个逗号分隔，空串为全部支持
	Online       int8         `json:"online"`        // 状态 1-上线 2-下线
	Source       int8         `json:"source"`        // 来源：1-官方插件 2-第三方插件 3-个人插件
	Desc         string       `json:"desc"`          // 详细介绍
	Creator      *UsersBase   `json:"creator"`       // creator
	Manager      []*UsersBase `json:"manager"`       // 管理员
	CreatedAt    int64        `json:"create_time"`   // 记录创建时间
	UpdatedAt    int64        `json:"update_time"`   // 记录最后修改时间
}

type PluginConfigsRequest struct {
	PluginID      int `json:"plugin_id"`      // 插件 id
	PluginVersion int `json:"plugin_version"` // 插件版本
}

type PluginConfigsResponse struct {
	Configs []*PluginConfig `json:"configs"`
}

// AddPluginRequest 新增插件
type AddPluginRequest struct {
	Name         string `json:"name"`          // 过滤器名称
	Intro        string `json:"intro"`         // 中文简介
	Func         string `json:"func"`          // 过滤器注册函数名
	SupportTypes []int8 `json:"support_types"` // 支持的过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Desc         string `json:"desc"`          // 插件介绍
}

// UpdatePluginRequest 更新插件
type UpdatePluginRequest struct {
	PluginID     int    `json:"plugin_id"`     // 插件 id
	Name         string `json:"name"`          // 过滤器名称
	Intro        string `json:"intro"`         // 中文简介
	SupportTypes []int8 `json:"support_types"` // 支持的过滤器类型 1-前置过滤器 2-后置过滤器 3-defer 过滤器
	Desc         string `json:"desc"`          // 插件介绍
}

type AddPluginResponse struct {
	ID int `json:"id"`
}

// ReplacePluginConfigRequest 新增/修改插件配置
type ReplacePluginConfigRequest struct {
	PluginID      int    `json:"plugin_id"`      // 插件 id
	PluginVersion int    `json:"plugin_version"` // 插件版本
	Key           string `json:"key"`            // 插件配置 key
	Name          string `json:"name"`           // 插件配置名
	Type          int8   `json:"type"`           // 配置类型 1-bool、2-string、3-int、4-uint、5-float、6-枚举、7-时间、8-array、9-map、10-multi-conf
	NotNull       int8   `json:"not_null"`       // 是否必输 1-是 2-否
	MoreInfo      string `json:"more_info"`      // 更多细节
	Default       string `json:"default"`        // 默认值，仅用于预填充配置值。
	Desc          string `json:"desc"`           // 配置描述
}

// DelPluginConfigRequest 删除插件配置
type DelPluginConfigRequest struct {
	PluginID      int    `json:"plugin_id"`      // 插件 id
	PluginVersion int    `json:"plugin_version"` // 插件版本
	Key           string `json:"key"`            // 插件配置 key
}

// PluginConfig 插件配置
type PluginConfig struct {
	Id       int    `json:"id"`
	Key      string `json:"key"`       // 插件配置 key
	Name     string `json:"name"`      // 插件配置名
	Type     int8   `json:"type"`      // 配置类型 1-bool、2-string、3-int、4-uint、5-float、6-枚举 7-时间、8-array、9-map、10-multi-conf
	NotNull  int8   `json:"not_null"`  // 是否必输 1-是 2-否
	MoreInfo string `json:"more_info"` // 更多细节
	Default  string `json:"default"`   // 默认值，仅用于预填充配置值。
	Desc     string `json:"desc"`      // 配置描述
}

// MinMaxMoreInfo int、uint、float 类型的最小最大值 more_info
type MinMaxMoreInfo struct {
	Min interface{} `json:"min,omitempty"` // 最小值，如果为 nil 则无最小值，可根据 type 可以转化为 int、uint、float64
	Max interface{} `json:"max,omitempty"` // 最大值，如果为 nil 则无最大值，可根据 type 可以转化为 int、uint、float64
}

// EnumMoreInfo 枚举类型的 more_info
type EnumMoreInfo struct {
	Multiple bool          `json:"multiple"` // 是否多选 true-多选 false-单选
	Options  []*EnumOption `json:"options"`  // 选项
}

type EnumOption struct {
	Key  string `json:"key,omitempty"`  // 选项 key
	Name string `json:"name,omitempty"` // 选项名
}

// TimeMoreInfo 时间类型的 more_info
type TimeMoreInfo struct {
	Type string `json:"type"` //包含 time、date、time_interval、date_interval
}

// TypeMoreInfo array、map 类型的 more_info
type TypeMoreInfo struct {
	Type     int8   `json:"type"`                // 配置类型 1-bool、2-string、3-int、4-uint、5-float、6-枚举（array 不包含） 7-时间
	MoreInfo string `json:"more_info,omitempty"` // 更多细节
}

// MultiConfMoreInfo multi-conf 类型的 more_info 为 []*MultiConfMoreInfo
type MultiConfMoreInfo struct {
	Key      string `json:"key,omitempty"`       // 插件配置 key
	Name     string `json:"name,omitempty"`      // 插件配置名
	Type     int8   `json:"type,omitempty"`      // 配置类型 1-bool、2-string、3-int、4-uint、5-float、6-枚举、7-时间
	MoreInfo string `json:"more_info,omitempty"` // 更多细节
	Default  string `json:"default,omitempty"`   // 默认值，仅用于预填充配置值。
	Desc     string `json:"desc,omitempty"`      // 配置描述
}
