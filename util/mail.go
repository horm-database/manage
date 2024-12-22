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
package util

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/horm-database/common/errs"
)

// SendMail 发送邮箱
// from 邮箱账号
// password 邮箱密码
// to 收件人邮箱
// cc 抄送人邮箱
// subject 邮件主题
// msg 邮件内容
func SendMail(from, user, password string, to, cc []string, mailType string, subject, body []byte) error {
	// 服务器认证信息
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	var msg []byte

	msg = append(msg, fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\n", emailsArr2Str(to), user, from)...)

	if len(cc) > 0 {
		msg = append(msg, fmt.Sprintf("%sCC:%s\r\n", msg, emailsArr2Str(cc))...)
	}

	msg = append(msg, "Subject:"...)
	msg = append(msg, subject...)
	msg = append(msg, fmt.Sprintf("\r\n%s\r\n\r\n", mailType)...)
	msg = append(msg, body...)

	// 发送邮件
	err := smtp.SendMail("smtp.qq.com:587", auth, from, to, msg)
	if err != nil {
		return errs.Newf(errs.RetWebEmailSendFailed, "send email error: %v", err)
	}

	return nil
}

func emailsArr2Str(emails []string) string {
	var str strings.Builder
	for k, email := range emails {
		if k == 0 {
			str.WriteString(email)
		} else {
			str.WriteString(";")
			str.WriteString(email)
		}
	}

	return str.String()
}
