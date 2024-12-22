// Copyright (c) 2024 The horm-database Authors. All rights reserved.
// This file Author:  CaoHao <18500482693@163.com> .
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

package mail

import (
	"github.com/horm-database/manage/util"
)

// SendMail 发送邮箱
// to 收件人邮箱
// cc 抄送人邮箱
// subject 邮件主题
// msg 邮件内容
func SendMail(to, cc []string, subject, msg []byte) error {
	var user = "聚码数据"
	var from = "649947921@qq.com"
	var password = "lumwyaqngqycbffh"
	var mailType = "Content-Type: text/html; charset=UTF-8"

	return util.SendMail(from, user, password, to, cc, mailType, subject, msg)
}
