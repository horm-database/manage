package pb

type AppIDRequest struct {
	Appid uint64 `json:"appid"` // appid
}

type AddAppRequest struct {
	Name    string   `json:"name"`    // 应用名称
	Intro   string   `json:"intro"`   // 简介
	Manager []uint64 `json:"manager"` // 管理员
}

type AddAppResponse struct {
	Appid uint64 `json:"appid"` // 应用appid
}

type UpdateAppRequest struct {
	Appid uint64 `json:"appid"` // 应用appid
	Name  string `json:"name"`  // 应用名称
	Intro string `json:"intro"` // 简介
}

type ResetAppSecretResponse struct {
	Secret string `json:"secret"` // 新秘钥
}

type UpdateAppStatusRequest struct {
	Appid  uint64 `json:"appid"`  // 应用appid
	Status int8   `json:"status"` // 状态 1-在线 2-离线
}

type MaintainAppManagerRequest struct {
	AppID   uint64   `json:"appid"`   // 应用appid
	Manager []uint64 `json:"manager"` // 管理员
}

type AppListRequest struct {
	Page int `json:"page"` // 分页
	Size int `json:"size"` // 每页大小
}

type AppListResponse struct {
	Total     uint64     `json:"total"`      // 总数
	TotalPage uint32     `json:"total_page"` // 总页数
	Page      int        `json:"page"`       // 分页
	Size      int        `json:"size"`       // 每页大小
	Apps      []*AppBase `json:"apps"`       // 应用列表
}

// AppBase 应用基础信息
type AppBase struct {
	Appid     uint64       `json:"appid"`      // 应用appid
	Name      string       `json:"name"`       // 应用名称
	Intro     string       `json:"intro"`      // 简介
	IsManager bool         `json:"is_manager"` // 是否管理员
	Creator   *UsersBase   `json:"creator"`    // Creator
	Manager   []*UsersBase `json:"manager"`    // 管理员，多个逗号分隔
	Status    int8         `json:"status"`     // 1-正常 2-下线
	CreatedAt int64        `json:"created_at"` // 创建时间
	UpdatedAt int64        `json:"updated_at"` // 最后修改时间
}

type AppDetailResponse struct {
	AppInfo *AppBase `json:"app_info"` // app 基础信息
	Secret  string   `json:"secret"`   // app 秘钥
}
