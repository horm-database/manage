package table

import (
	"context"
	"time"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/api/pb"
	"github.com/samber/lo"
)

func GetUserBaseFromUser(user *TblUser) *pb.UsersBase {
	return &pb.UsersBase{
		UserID:    user.Id,
		Account:   user.Account,
		Nickname:  user.Nickname,
		AvatarUrl: user.AvatarUrl,
		Gender:    user.Gender,
	}
}

func InsertUser(ctx context.Context, user *TblUser) error {
	_, err := GetTableORM("tbl_user").Insert(user).Exec(ctx)
	return err
}

func UpdateUserByID(ctx context.Context, id uint64, updateInfo horm.Map) error {
	_, err := GetTableORM("tbl_user").Update(updateInfo).Eq("id", id).Exec(ctx)
	return err
}

func GetUserByAccount(ctx context.Context, account string) (bool, *TblUser, error) {
	user := TblUser{}

	notFind, err := GetTableORM("tbl_user").FindBy("account", account).Exec(ctx, &user)

	return notFind, &user, err
}

func GetUserByID(ctx context.Context, id uint64) (bool, *TblUser, error) {
	user := TblUser{}

	notFind, err := GetTableORM("tbl_user").FindBy("id", id).Exec(ctx, &user)

	return notFind, &user, err
}

func GetUsersByKeyword(ctx context.Context, keyword string) ([]*TblUser, error) {
	users := []*TblUser{}

	where := horm.Where{
		"OR": horm.Where{
			"account ~":  "%" + keyword + "%",
			"nickname ~": "%" + keyword + "%",
		},
	}
	_, err := GetTableORM("tbl_user").FindAll(where).Exec(ctx, &users)

	return users, err
}

func GetUserBasesByIds(ctx context.Context, userIds []uint64) ([]*pb.UsersBase, error) {
	if len(userIds) == 0 {
		return []*pb.UsersBase{}, nil
	}

	userIds = lo.Uniq(userIds)

	users := []*TblUser{}

	_, err := GetTableORM("tbl_user").FindAllBy("id", userIds).Exec(ctx, &users)
	if err != nil {
		return nil, err
	}

	ret := []*pb.UsersBase{}

	for _, v := range users {
		ret = append(ret, GetUserBaseFromUser(v))
	}

	return ret, err
}

func GetUserBasesMapByIds(ctx context.Context, userIds []uint64) (map[uint64]*pb.UsersBase, error) {
	ret := map[uint64]*pb.UsersBase{}

	users, err := GetUserBasesByIds(ctx, userIds)
	if len(users) == 0 || err != nil {
		return ret, err
	}

	for _, v := range users {
		ret[v.UserID] = v
	}

	return ret, err
}

func GetSequence(ctx context.Context) (uint64, error) {
	ret := proto.ModRet{}
	_, err := GetTableORM("tbl_sequence").Insert(horm.Map{"created_at": time.Now()}).Exec(ctx, &ret)

	if err != nil {
		return 0, err
	}

	return ret.ID.Uint64(), nil
}
