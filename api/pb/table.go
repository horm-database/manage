package pb

// TableBase 表基础信息
type TableBase struct {
	Id         int    `json:"table_id"`
	Name       string `json:"name"`        // 数据名称
	Intro      string `json:"intro"`       // 简介
	Desc       string `json:"desc"`        // 详细描述
	Status     int8   `json:"status"`      // 状态 1-正常 2-下线
	CreateTime int64  `json:"create_time"` // 记录创建时间
}

type AddTableRequest struct {
	Name        string `json:"name"`         // 数据名称（执行单元名）
	Intro       string `json:"intro"`        // 中文简介
	Desc        string `json:"desc"`         // 详细描述
	TableVerify string `json:"table_verify"` // 表校验，为空时不校验，默认同 name，即只允许访问 name 表/索引
	DB          int    `json:"db"`           // 所属数据库
}

type UpdateTableBaseRequest struct {
	TableID int    `json:"table_id"`
	Intro   string `json:"intro"` // 中文简介
	Desc    string `json:"desc"`  // 详细描述
}

type AddTableResponse struct {
	ID int `json:"id"` // table id
}

type UpdateTableStatusRequest struct {
	TableID int  `json:"table_id"` // table id
	Status  int8 `json:"status"`   // 状态 1-在线 2-离线
}

type UpdateTableAdvanceRequest struct {
	TableID     int    `json:"table_id"`     // table id
	TableVerify string `json:"table_verify"` // 表校验，为空时不校验，默认同 name，即只允许访问 name 表/索引
}

type TableIDRequest struct {
	TableID int `json:"table_id"`
}

type TableDetailResponse struct {
	Info           *TableBase             `json:"info"`            // 表基础信息
	IsManager      bool                   `json:"is_manager"`      // 是否表管理员（即库管理员或产品管理员）
	Creator        *UsersBase             `json:"creator"`         // creator
	DBInfo         *DBBase                `json:"db_info"`         // 所属数据库基础信息
	DBManager      []*UsersBase           `json:"db_manager"`      // 所属库管理员
	ProductInfo    *ProductBase           `json:"product_info"`    // 所属产品基础信息
	ProductManager []*UsersBase           `json:"product_manager"` // 所属产品管理员（拥有库管理员所有权限）
	TableFields    []*TableField          `json:"table_fields"`    // 表字段
	TableIndexs    []*TableIndex          `json:"table_indexs"`    // 表索引
	LangStructs    map[string]*LangStruct `json:"lang_structs"`    // 各语言结构体
	ReadMe         string                 `json:"readme"`          // 表操作文档
}

type TableAdvanceConfigResponse struct {
	Definition  string `json:"definition"`   // 表定义
	TableVerify string `json:"table_verify"` // 表校验
}

// TableField 表字段
type TableField struct {
	Field     string `json:"field"`      // 字段
	Type      int8   `json:"type"`       // 类型 1-bool、2-string、3-int、4-int8、5-int16、6-int32、7-int64、8-uint、9-uint8、10-uint16、11-uint32、12-uint64、13-float、14-float64、15-blob(bytes)、16-enum、17-json 18-date（2006-01-02）、19-datetime（2023-10-21T13:40:21+08:00）
	Len       string `json:"len"`        // 长度
	Empty     bool   `json:"empty"`      // 是否可空 true-是 false-否
	Status    int8   `json:"status"`     // 状态 1-正常 2-隐藏（对非表管理员）
	Default   string `json:"default"`    // 默认值
	IsPrimary bool   `json:"is_primary"` // 是否主键
	IsIndex   bool   `json:"is_index"`   // 是否索引
	Comment   string `json:"comment"`    // 注释
	More      string `json:"more"`       // 其他，例如枚举值。
}

// TableIndex 表索引
type TableIndex struct {
	Name   string `json:"name"`   // 索引名
	Type   string `json:"type"`   // 索引类型
	Fields string `json:"fields"` // 包含列
}

// LangStruct 各语言对应结构体
type LangStruct struct {
	Language string `json:"language"` // 语言：go、java、c++、node、php
	Struct   string `json:"struct"`
}
