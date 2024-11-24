package pb

// SupportOpsResponse 数据支持的所有操作
type SupportOpsResponse struct {
	DBType     int      `json:"db_type"`     // 数据库类型 0-nil（仅执行插件） 1-elastic 2-mongo 3-redis 10-mysql 11-postgresql 12-clickhouse 13-oracle 14-DB2 15-sqlite
	SupportOps []string `json:"support_ops"` // 数据支持的所有操作
}

// AppCanAccessDBRequest 我的能接入指定仓库的所有应用
type AppCanAccessDBRequest struct {
	DbID    int    `json:"db_id"`   // 数据库
	Keyword string `json:"keyword"` // 过滤关键词
}

type AppCanAccessDBResponse struct {
	Apps []*AppCanAccessDB `json:"apps"` // 我的能接入仓库的应用
}

type AppCanAccessDB struct {
	Appid        uint64        `json:"appid"`         // 应用appid
	AppName      string        `json:"app_name"`      // 应用名称
	Intro        string        `json:"intro"`         // 简介
	AccessStatus int8          `json:"access_status"` // 接入状态：0-未接入 1-正常 2-下线 3-审核中 4-审核撤回 5-拒绝
	AccessInfo   *DBAccessInfo `json:"access_info"`   // 接入信息
}

type DBAccessInfo struct {
	AccessID     int      `json:"access_id"`     // 接入 ID
	AccessRoot   int8     `json:"access_root"`   // 权限 1-超级权限（所有权限，包含DDL）  2-表数据权限（库下表的所有增删改查权限，不包含 DDL）  3-无
	AccessOp     []string `json:"access_op"`     // 支持的操作
	AccessReason string   `json:"access_reason"` // 接入原因
}

// AppApplyAccessDBRequest 应用申请接入仓库
type AppApplyAccessDBRequest struct {
	Appid  uint64   `json:"appid"`  // 应用appid
	DbID   int      `json:"db_id"`  // 数据库
	Root   int8     `json:"root"`   // 权限 1-超级权限（所有权限，包含DDL）  2-表数据权限（库下表的所有增删改查权限，不包含 DDL）  3-无
	Op     []string `json:"op"`     // 支持的操作
	Reason string   `json:"reason"` // 接入原因
}

type AppApplyAccessResponse struct {
	AccessID int `json:"access_id"` // 申请ID
}

// AppAccessDBApprovalRequest 应用接入仓库审批
type AppAccessDBApprovalRequest struct {
	Appid  uint64 `json:"appid"`  // 应用appid
	DbID   int    `json:"db_id"`  // 数据库
	Status int8   `json:"status"` // 1-审批通过 2-审批拒绝
	Reason string `json:"reason"` // 拒绝理由（ status=2 时输入）
}

// AppAccessDBWithdrawRequest 应用接入仓库撤销申请
type AppAccessDBWithdrawRequest struct {
	Appid  uint64 `json:"appid"`  // 应用appid
	DbID   int    `json:"db_id"`  // 数据库
	Reason string `json:"reason"` // 撤销理由
}

// AppAccessDBUpdateRequest 编辑仓库访问权限
type AppAccessDBUpdateRequest struct {
	Appid  uint64   `json:"appid"`  // 应用appid
	DbID   int      `json:"db_id"`  // 数据库
	Root   int8     `json:"root"`   // 权限 1-超级权限（所有权限，包含DDL）  2-表数据权限（库下表的所有增删改查权限，不包含 DDL）  3-无
	Op     []string `json:"op"`     // 支持的操作
	Reason string   `json:"reason"` // 编辑原因
}

// AppAccessDBOnOffRequest 仓库访问权限上/下线
type AppAccessDBOnOffRequest struct {
	Appid  uint64 `json:"appid"`  // 应用appid
	DbID   int    `json:"db_id"`  // 数据库
	Status int8   `json:"status"` // 状态：1-上线 2-下线
	Reason string `json:"reason"` // 上/下线原因
}

// DBsAllAppAccessListRequest 访问该仓库的应用列表
type DBsAllAppAccessListRequest struct {
	DbID int `json:"db_id"` // 数据库
	Page int `json:"page"`  // 分页
	Size int `json:"size"`  // 每页大小
}

type DBsAllAppAccessListResponse struct {
	Total        uint64         `json:"total"`          // 总数
	TotalPage    uint32         `json:"total_page"`     // 总页数
	Page         int            `json:"page"`           // 分页
	Size         int            `json:"size"`           // 每页大小
	IsDBManager  bool           `json:"is_db_manager"`  // 是否仓库管理员
	AppAccessDBs []*AppAccessDB `json:"app_access_dbs"` // 访问列表
}

// AppsAllDBAccessListRequest 该应用访问的仓库列表
type AppsAllDBAccessListRequest struct {
	Appid uint64 `json:"appid"` // 应用id
	Page  int    `json:"page"`  // 分页
	Size  int    `json:"size"`  // 每页大小
}

type AppsAllDBAccessListResponse struct {
	Total        uint64         `json:"total"`          // 总数
	TotalPage    uint32         `json:"total_page"`     // 总页数
	Page         int            `json:"page"`           // 分页
	Size         int            `json:"size"`           // 每页大小
	AppAccessDBs []*AppAccessDB `json:"app_access_dbs"` // 访问列表
}

type AppAccessDB struct {
	Id        int        `json:"id"`
	App       *AppBase   `json:"app,omitempty"` // 应用信息
	DB        *DBBase    `json:"db,omitempty"`  // 库信息
	Root      int8       `json:"root"`          // 超级权限 1-超级权限（所有权限，包含DDL）  2-表数据权限（库下表的所有增删改查权限，不包含 DDL）  3-无
	Op        []string   `json:"op"`            // 支持的操作
	Status    int8       `json:"status"`        // 状态：1-正常 2-下线 3-审核中 4-审核撤回 5-拒绝
	ApplyUser *UsersBase `json:"apply_user"`    // 申请者
	Reason    string     `json:"reason"`        // 接入原因
	CreatedAt int64      `json:"create_time"`   // 记录创建时间
	UpdatedAt int64      `json:"update_time"`   // 最后更新时间
}
