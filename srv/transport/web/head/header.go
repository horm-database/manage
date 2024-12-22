// Copyright (c) 2024 The horm-database Authors (such as CaoHao <18500482693@163.com>). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
