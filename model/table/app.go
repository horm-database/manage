package table

import (
	"context"
	"fmt"

	"github.com/horm-database/common/proto"
	"github.com/horm-database/go-horm/horm"
	"github.com/horm-database/manage/consts"
	"github.com/horm-database/server/model/table"
)

func AddApp(ctx context.Context, appInfo *table.TblAppInfo) error {
	_, err := GetTableORM("tbl_app_info").Insert(appInfo).Exec(ctx)
	return err
}

func UpdateAppByID(ctx context.Context, appid uint64, update horm.Map) error {
	_, err := GetTableORM("tbl_app_info").Eq("appid", appid).Update(update).Exec(ctx)
	return err
}

func GetAppList(ctx context.Context, userid uint64, page, size int) (*proto.Detail, []*table.TblAppInfo, error) {
	pageResult := proto.Detail{}

	apps := []*table.TblAppInfo{}

	// 在线的应用和我管理的应用。
	where := horm.Where{
		"OR": horm.Where{
			"status": consts.StatusOnline,
			"AND": horm.Where{
				"status":    consts.StatusOffline,
				"manager ~": "%" + fmt.Sprint(userid) + "%",
			},
		},
	}

	_, err := GetTableORM("tbl_app_info").FindAll(where).Order("-appid").Page(page, size).Exec(ctx, &pageResult, &apps)

	return &pageResult, apps, err
}

func GetAppListByAppids(ctx context.Context, appids []uint64) ([]*table.TblAppInfo, error) {
	apps := []*table.TblAppInfo{}

	_, err := GetTableORM("tbl_app_info").FindAllBy("appid", appids).Exec(ctx, &apps)

	return apps, err
}

func GetMyAppListByKeyword(ctx context.Context, userid uint64,
	keyword string, status int8) ([]*table.TblAppInfo, error) {
	apps := []*table.TblAppInfo{}

	where := horm.Where{
		"manager ~": "%" + fmt.Sprint(userid) + "%",
	}

	if status != 0 {
		where["status"] = status
	}

	if keyword != "" {
		where["OR"] = horm.Where{
			"appid ~": "%" + fmt.Sprint(keyword) + "%",
			"name ~":  "%" + fmt.Sprint(keyword) + "%",
		}
	}

	_, err := GetTableORM("tbl_app_info").FindAll(where).Order("-appid").Exec(ctx, &apps)

	return apps, err
}

func GetAppDetail(ctx context.Context, appid uint64) (bool, *table.TblAppInfo, error) {
	app := table.TblAppInfo{}

	isNil, err := GetTableORM("tbl_app_info").
		FindBy("appid", appid).
		Exec(ctx, &app)

	return isNil, &app, err
}
