package head

const (
	Version     = "web-version"      // 客户端版本
	RequestID   = "web-request-id"   // 请求唯一id
	Timestamp   = "web-timestamp"    // 请求时间戳（精确到毫秒）
	Timeout     = "web-timeout"      // 请求超时时间，单位 ms
	UserID      = "web-userid"       // userid
	WorkspaceID = "web-workspace-id" // workspace-id
	Caller      = "web-caller"       // caller
	AuthRand    = "web-auth-rand"    // 随机生成 0-99999999 的数字，相同 timestamp 不允许出现同样的 ip、auth_rand。为了避免碰撞，0-99999999，单机理论最大支持 1000 亿/秒的并发。
	Sign        = "web-sign"         // 签名 md5(workspace_id+userid+token+version+request_id+timestamp+timeout+caller+auth_rand)
)
