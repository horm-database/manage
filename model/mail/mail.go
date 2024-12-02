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
