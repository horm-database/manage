package consts

const (
	CachePreEmailCode = "PreEmailCode_"
)

const (
	CacheEmailCodeExpire     = 5 * 60
	CacheSendEmailFrequently = 120
)

const (
	WorkspaceMemberNotJoin = 0 // 非空间成员
	WorkspaceMember        = 1 // 空间成员
	WorkspaceMemberManager = 2 // 空间管理员
	WorkspaceMemberExpired = 3 // 权限已过期
)

const (
	WorkspaceMemberStatusApproval = 1 // 待审批
	WorkspaceMemberStatusRenewal  = 2 // 续期审批
	WorkspaceMemberStatusNotApply = 3 // 暂未申请
	WorkspaceMemberStatusJoined   = 4 // 已加入
	WorkspaceMemberStatusReject   = 5 // 审批拒绝
	WorkspaceMemberStatusQuit     = 6 // 已退出
	WorkspaceMemberStatusExpired  = 9 // 已过期
)

const (
	ExpireTypePermanent = 0 // 永久
	ExpireType1month    = 1 // 一个月
	ExpireType3month    = 2 // 三个月
	ExpireTypeHalfYear  = 3 // 半年
	ExpireTypeYear      = 4 // 一年
)

const (
	StatusOnline  = 1 // 正常
	StatusOffline = 2 // 下线
)

const (
	SearchTypeProduct = 1 // product
	SearchTypeDB      = 2 // db
	SearchTypeTable   = 3 // table
)

const (
	ProductRoleNotJoin   = 0 // 非产品成员
	ProductRoleManager   = 1 // 产品管理员
	ProductRoleDeveloper = 2 // 产品开发者
	ProductRoleOperator  = 3 // 产品运营者
	ProductRoleExpired   = 4 // 成员权限已过期
)

const (
	ProductMemberStatusNotApply   = 0 // 未加入
	ProductMemberStatusApproval   = 1 // 待审批
	ProductMemberStatusRenewal    = 2 // 续期审批
	ProductMemberStatusChangeRole = 3 // 角色变更审批
	ProductMemberStatusJoined     = 4 // 已加入
	ProductMemberStatusReject     = 5 // 审批拒绝
	ProductMemberStatusQuit       = 6 // 已退出
	ProductMemberStatusExpired    = 7 // 已过期
)

const (
	FilterSourceOfficial = 1
	FilterSourceThird    = 2
	FilterSourcePrivate  = 3
)