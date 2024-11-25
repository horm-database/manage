package table

import (
	"time"

	"github.com/horm-database/manage/consts"
	"github.com/horm-database/orm"
)

// GetTableORM 获取经校验后的 ORM
func GetTableORM(table string) *orm.ORM {
	cli := orm.NewORM(consts.DBConfigName)
	cli.Name(table)
	return cli
}

type TblUser struct {
	Id            uint64    `orm:"id,uint64,omitempty" json:"id,omitempty"`                        // 用户id
	Account       string    `orm:"account,string,omitempty" json:"account,omitempty"`              // 账号，可以是 admin，邮箱等。。。
	Password      string    `orm:"password,string,omitempty" json:"password,omitempty"`            // 密码
	Nickname      string    `orm:"nickname,string,omitempty" json:"nickname,omitempty"`            // 昵称
	Mobile        string    `orm:"mobile,string,omitempty" json:"mobile,omitempty"`                // 手机号
	Token         string    `orm:"token,string,omitempty" json:"token,omitempty"`                  // token
	AvatarUrl     string    `orm:"avatar_url,string,omitempty" json:"avatar_url,omitempty"`        // 头像
	Gender        int8      `orm:"gender,int,omitempty" json:"gender,omitempty"`                   // 性别 1-男  2-女
	Company       string    `orm:"company,string,omitempty" json:"company,omitempty"`              // 公司
	Department    string    `orm:"department,string,omitempty" json:"department,omitempty"`        // 部门
	City          string    `orm:"city,string,omitempty" json:"city,omitempty"`                    // 城市
	Province      string    `orm:"province,string,omitempty" json:"province,omitempty"`            // 省份
	Country       string    `orm:"country,string,omitempty" json:"country,omitempty"`              // 国家
	LastLoginTime int       `orm:"last_login_time,int,omitempty" json:"last_login_time,omitempty"` // 上次登录时间
	LastLoginIP   string    `orm:"last_login_ip,string,omitempty" json:"last_login_ip,omitempty"`  // 上次登录ip
	CreatedAt     time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`                // 记录创建时间
	UpdatedAt     time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`                // 记录最后修改时间
}

type TblWorkspaceMember struct {
	Id          int       `orm:"id,int,omitempty" json:"id,omitempty"`                     // member id
	WorkspaceID int       `orm:"workspace_id,int,omitempty" json:"workspace_id,omitempty"` // workspace id
	UserID      uint64    `orm:"userid,uint64,omitempty" json:"userid,omitempty"`          // 用户id
	Status      int8      `orm:"status,int8,omitempty" json:"status,omitempty"`            // 1-待审批 2-续期审批 3-暂未申请 4-已加入 5-审批拒绝  6-已退出
	JoinTime    int64     `orm:"join_time,int,omitempty" json:"join_time,omitempty"`       // 加入时间
	ExpireType  int8      `orm:"expire_type,int8,omitempty" json:"expire_type,omitempty"`  // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	ExpireTime  int       `orm:"expire_time,int,omitempty" json:"expire_time,omitempty"`   // 过期时间
	OutTime     int       `orm:"out_time,int,omitempty" json:"out_time,omitempty"`         // 退出时间
	CreatedAt   time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`          // 记录创建时间
	UpdatedAt   time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`          // 记录最后修改时间
}

type TblCollectTable struct {
	Id        int       `orm:"id,int,omitempty" json:"id,omitempty"`             // id
	UserID    uint64    `orm:"userid,uint64,omitempty" json:"userid,omitempty"`  // 用户id
	TableID   int       `orm:"table_id,int,omitempty" json:"table_id,omitempty"` // 表id
	CreatedAt time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`  // 记录创建时间
	UpdatedAt time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`  // 记录最后修改时间
}

type TblProduct struct {
	Id        int       `orm:"id,int,omitempty" json:"id,omitempty"`              // id
	Name      string    `orm:"name,string,omitempty" json:"name,omitempty"`       // 产品名称
	Intro     string    `orm:"intro,string,omitempty" json:"intro,omitempty"`     // 简介
	Creator   uint64    `orm:"creator,uint64,omitempty" json:"creator,omitempty"` // Creator
	Manager   string    `orm:"manager,string,omitempty" json:"manager,omitempty"` // 管理员，多个逗号分隔
	Status    int8      `orm:"status,int8,omitempty" json:"status,omitempty"`     // 1-正常 2-下线
	CreatedAt time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`   // 记录创建时间
	UpdatedAt time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`   // 记录最后修改时间
}

type TblProductMember struct {
	Id         int       `orm:"id,int,omitempty" json:"id,omitempty"`                    // member id
	ProductID  int       `orm:"product_id,int,omitempty" json:"product_id,omitempty"`    // product id
	UserID     uint64    `orm:"userid,uint64,omitempty" json:"userid,omitempty"`         // 用户id
	Role       int8      `orm:"role,int8,omitempty" json:"role,omitempty"`               // 1-管理员 2-开发者 3-运营者
	Status     int8      `orm:"status,int8,omitempty" json:"status,omitempty"`           // 1-待审批 2-续期审批 3-角色变更审批 4-已加入 5-审批拒绝  6-已退出
	JoinTime   int64     `orm:"join_time,int,omitempty" json:"join_time,omitempty"`      // 加入时间
	ExpireType int8      `orm:"expire_type,int8,omitempty" json:"expire_type,omitempty"` // 0: 永久 1: 一个月 2: 三个月 3: 半年 4: 一年
	ExpireTime int       `orm:"expire_time,int,omitempty" json:"expire_time,omitempty"`  // 过期时间
	OutTime    int       `orm:"out_time,int,omitempty" json:"out_time,omitempty"`        // 退出时间
	ChangeRole int8      `orm:"change_role,int8,omitempty" json:"change_role,omitempty"` // 变更为目标角色 0-无 2-开发者 3-运营者
	CreatedAt  time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`         // 记录创建时间
	UpdatedAt  time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`         // 记录最后修改时间
}

type TblSearchKeyword struct {
	Id        int       `orm:"id,int,omitempty" json:"id,omitempty"`                // id
	Type      int8      `orm:"type,int8,omitempty" json:"type,omitempty"`           // 1-product 2-db 3-table
	Sid       int       `orm:"sid,int,omitempty" json:"sid,omitempty"`              // 检索id
	SName     string    `orm:"sname,string,omitempty" json:"sname,omitempty"`       // 检索名
	Field     string    `orm:"field,string,omitempty" json:"field,omitempty"`       // 字段
	SKey      string    `orm:"skey,string,omitempty" json:"skey,omitempty"`         // 检索key
	SContent  string    `orm:"scontent,string,omitempty" json:"scontent,omitempty"` // 检索内容
	CreatedAt time.Time `orm:"created_at,datetime,omitempty" json:"created_at"`     // 记录创建时间
	UpdatedAt time.Time `orm:"updated_at,datetime,omitempty" json:"updated_at"`     // 记录最后修改时间
}
