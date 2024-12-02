package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/horm-database/common/crypto"
	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/types"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/manage/model/cache"
	"github.com/horm-database/manage/model/mail"
	"github.com/horm-database/manage/model/table"
	"github.com/horm-database/manage/srv/transport/web/head"
	"github.com/samber/lo"
)

// SendEmailCode 发送邮箱验证码
func SendEmailCode(ctx context.Context, req *pb.SendEmailCodeRequest) error {
	key := fmt.Sprintf("%s%s", consts.CachePreEmailCode, req.Account)

	ttl, err := cache.TTLCacheByKey(ctx, key)
	if err != nil {
		return err
	}

	if ttl > consts.CacheEmailCodeExpire-consts.CacheSendEmailFrequently { // 邮件发送过于频繁
		return errs.New(errs.RetWebEmailSendFrequently, "send email frequently")
	}

	code := 1000 + rand.Intn(8888)

	var subject, body string
	if req.Type == 0 {
		subject = "聚码数据—邮箱身份验证"

		now := time.Now()
		body = fmt.Sprintf(`亲爱的用户：<br><br>
	您好，感谢使用聚码服务，您正在进行邮箱验证，<br><br>
	本次请求的验证码为 <font size="4" style="color:#FFA500;"><b>%d</b></font><font style="color:#989898;">（为了保证您的账号安全，请在 5 分钟内完成验证）</font>，如非本人操作请忽略，切勿将此验证码泄露给他人，以免给您账号下的数据带来损失。<br><br>
	聚码数据团队<br>%d年%02d月%02d日<br>`, code, now.Year(), now.Month(), now.Day())
	}

	err = mail.SendMail([]string{req.Account}, nil, []byte(subject), []byte(body))
	if err != nil {
		return err
	}

	err = cache.SetCacheByKey(ctx, key, code, consts.CacheEmailCodeExpire)
	if err != nil {
		return err
	}

	return nil
}

func Register(ctx context.Context, req *pb.RegisterRequest) error {
	key := fmt.Sprintf("%s%s", consts.CachePreEmailCode, req.Account)

	var saveCode string
	_, err := cache.GetCacheByKey(ctx, key, &saveCode)
	if err != nil {
		return err
	}

	if saveCode != req.Code {
		return errs.New(errs.RetWebCodeIncorrectly, "code verify incorrectly")
	}

	notFind, tblUser, err := table.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return err
	}

	if !notFind {
		return errs.Newf(errs.RetWebAccountExists, "account is registered")
	}

	tblUser.Id, err = GenerateUserid(ctx)
	if err != nil {
		return err
	}

	tblUser.Account = req.Account
	tblUser.Nickname = req.Nickname
	tblUser.Password = crypto.MD5Str(req.Password)
	err = table.InsertUser(ctx, tblUser)

	// 验证码被用掉之后不可重复利用
	_ = cache.DelCacheByKey(ctx, key)

	return err
}

func Login(ctx context.Context, head *head.WebReqHeader, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	notFind, tblUser, err := table.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}

	if notFind {
		return nil, errs.Newf(errs.RetWebAccountNotExists, "account is not exists")
	}

	if tblUser.Password != crypto.MD5Str(req.Password) {
		return nil, errs.Newf(errs.RetWebPasswordIncorrect, "password verification failed")
	}

	loginToken := crypto.MD5Str(fmt.Sprintf("%d%d%d", tblUser.Id, time.Now().Unix(), rand.Intn(10000)))

	update := horm.Map{}
	update["last_login_time"] = time.Now().Unix()
	update["last_login_ip"] = head.Ip
	update["token"] = loginToken

	err = table.UpdateUserByID(ctx, tblUser.Id, update)
	if err != nil {
		return nil, err
	}

	ret := pb.LoginResponse{
		UserID:     tblUser.Id,
		Account:    tblUser.Account,
		Mobile:     tblUser.Mobile,
		Nickname:   tblUser.Nickname,
		Token:      loginToken,
		Gender:     tblUser.Gender,
		Company:    tblUser.Company,
		Department: tblUser.Department,
		City:       tblUser.City,
		Province:   tblUser.Province,
		Country:    tblUser.Country,
	}

	return &ret, nil
}

func ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) error {
	key := fmt.Sprintf("%s%s", consts.CachePreEmailCode, req.Account)

	var saveCode string
	_, err := cache.GetCacheByKey(ctx, key, &saveCode)
	if err != nil {
		return err
	}

	if saveCode != req.Code {
		return errs.New(errs.RetWebCodeIncorrectly, "code verify incorrectly")
	}

	notFind, tblUser, err := table.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return err
	}

	if notFind {
		return errs.Newf(errs.RetWebAccountExists, "not find user")
	}

	update := horm.Map{}
	update["password"] = crypto.MD5Str(req.Password)

	err = table.UpdateUserByID(ctx, tblUser.Id, update)

	// 验证码被用掉之后不可重复利用
	_ = cache.DelCacheByKey(ctx, key)

	return err
}

func FindUser(ctx context.Context, keyword string) (*pb.FindUserResponse, error) {
	ret := pb.FindUserResponse{Users: []*pb.UsersBase{}}
	if keyword == "" {
		return &ret, nil
	}

	users, err := table.GetUsersByKeyword(ctx, keyword)
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		ret.Users = append(ret.Users, table.GetUserBaseFromUser(v))
	}

	return &ret, nil
}

func FindUserByID(ctx context.Context, ids []uint64) (*pb.FindUserResponse, error) {
	ret := pb.FindUserResponse{Users: []*pb.UsersBase{}}
	if len(ids) == 0 {
		return &ret, nil
	}

	users, err := table.GetUserBasesByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	ret.Users = users

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func GenerateUserid(ctx context.Context) (uint64, error) {
	id, err := table.GetSequence(ctx)
	if err != nil {
		return 0, err
	}

	id = id % 100000000

	idStr := fmt.Sprintf("1%07d%04d", id, rand.Intn(9999))
	userid, err := strconv.ParseUint(idStr, 10, 64)
	return userid, err
}

// IsManager 是否管理员 true-是 false-否
func IsManager(userid uint64, managers string) bool {
	if managers == "" {
		return false
	}

	managerArr := GetUserIds(managers)
	if len(managers) == 0 {
		return false
	}

	for _, manager := range managerArr {
		if manager == userid {
			return true
		}
	}

	return false
}

func GetExpireTime(expireTime int64, expireType int8) int64 {
	start := time.Now()
	if expireTime > time.Now().Unix() { // 如还未过期，从过期时间开始算
		start = time.Unix(expireTime, 0)
	}

	switch expireType {
	case consts.ExpireType1month:
		return start.AddDate(0, 1, 0).Unix()
	case consts.ExpireType3month:
		return start.AddDate(0, 3, 0).Unix()
	case consts.ExpireTypeHalfYear:
		return start.AddDate(0, 6, 0).Unix()
	case consts.ExpireTypeYear:
		return start.AddDate(1, 0, 0).Unix()
	}

	return 0
}

func GetUsersFromMap(userMap map[uint64]*pb.UsersBase, userIds []uint64) []*pb.UsersBase {
	ret := []*pb.UsersBase{}

	for _, userid := range userIds {
		u, ok := userMap[userid]
		if ok {
			ret = append(ret, u)
		}
	}

	return ret
}

func GetUserIds(userIds ...interface{}) []uint64 {
	ret := []uint64{}

	if len(userIds) == 0 {
		return ret
	}

	for _, userid := range userIds {
		switch u := userid.(type) {
		case uint64:
			if u > 0 {
				ret = append(ret, u)
			}
		case []uint64:
			ret = append(ret, u...)
		case string:
			if u != "" {
				us := types.SplitUint64(u, ",")
				ret = append(ret, us...)
			}
		case []byte:
			if len(u) > 0 {
				us := types.SplitUint64(string(u), ",")
				ret = append(ret, us...)
			}
		}
	}

	return lo.Uniq(ret)
}
