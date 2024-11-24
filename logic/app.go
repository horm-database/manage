package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/horm/common/crypto"
	"github.com/horm/common/errs"
	"github.com/horm/common/types"
	"github.com/horm/go-horm/horm"
	"github.com/horm/manage/api/pb"
	"github.com/horm/manage/consts"
	"github.com/horm/manage/model/table"
	st "github.com/horm/server/model/table"
)

func AddApp(ctx context.Context, userid uint64, req *pb.AddAppRequest) (*pb.AddAppResponse, error) {
	if !types.InArrayUint64(req.Manager, userid) {
		req.Manager = append(req.Manager, userid)
	}

	appid, err := GenerateAppID(ctx)
	if err != nil {
		return nil, err
	}

	appInfo := st.TblAppInfo{
		Appid:   appid,
		Name:    req.Name,
		Secret:  GenerateAppSecret(),
		Intro:   req.Intro,
		Creator: userid,
		Manager: types.JoinUint64(req.Manager, ","),
		Status:  consts.StatusOnline,
	}

	err = table.AddApp(ctx, &appInfo)
	if err != nil {
		return nil, err
	}

	return &pb.AddAppResponse{Appid: appid}, nil
}

func UpdateApp(ctx context.Context, userid uint64, req *pb.UpdateAppRequest) error {
	_, err := IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return err
	}

	update := horm.Map{
		"name":  req.Name,
		"intro": req.Intro,
	}

	return table.UpdateAppByID(ctx, req.Appid, update)
}

func ResetAppSecret(ctx context.Context, userid, appid uint64) (*pb.ResetAppSecretResponse, error) {
	_, err := IsAppManager(ctx, userid, appid)
	if err != nil {
		return nil, err
	}

	newSecret := GenerateAppSecret()
	update := horm.Map{
		"secret": newSecret,
	}

	err = table.UpdateAppByID(ctx, appid, update)
	if err != nil {
		return nil, err
	}

	return &pb.ResetAppSecretResponse{Secret: newSecret}, nil
}

func UpdateAppStatus(ctx context.Context, userid uint64, req *pb.UpdateAppStatusRequest) error {
	_, err := IsAppManager(ctx, userid, req.Appid)
	if err != nil {
		return err
	}

	update := horm.Map{
		"status": req.Status,
	}

	return table.UpdateAppByID(ctx, req.Appid, update)
}

func MaintainAppManager(ctx context.Context, userid uint64, req *pb.MaintainAppManagerRequest) error {
	_, err := IsAppManager(ctx, userid, req.AppID)
	if err != nil {
		return err
	}

	managerUids := types.UniqUint64(req.Manager)

	update := horm.Map{
		"manager": types.JoinUint64(managerUids, ","),
	}

	return table.UpdateAppByID(ctx, req.AppID, update)
}

func AppList(ctx context.Context, userid uint64, req *pb.AppListRequest) (*pb.AppListResponse, error) {
	pageInfo, apps, err := table.GetAppList(ctx, userid, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	ret := pb.AppListResponse{
		Total:     pageInfo.Total,
		TotalPage: pageInfo.TotalPage,
		Page:      req.Page,
		Size:      req.Size,
		Apps:      []*pb.AppBase{},
	}

	var userIds []uint64
	for _, app := range apps {
		userIds = append(userIds, GetUserIds(app.Creator, app.Manager)...)
	}

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		ret.Apps = append(ret.Apps, GetAppBaseFromApp(userid, app, userMaps))
	}

	return &ret, nil
}

func AppDetail(ctx context.Context, userid, appid uint64) (*pb.AppDetailResponse, error) {
	app, err := IsAppManager(ctx, userid, appid)
	if err != nil {
		return nil, err
	}

	var userIds []uint64
	userIds = append(userIds, GetUserIds(app.Creator, app.Manager)...)

	userMaps, err := table.GetUserBasesMapByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	ret := pb.AppDetailResponse{
		AppInfo: &pb.AppBase{
			Appid:     app.Appid,
			Name:      app.Name,
			Intro:     app.Intro,
			IsManager: IsManager(userid, app.Manager),
			Creator:   userMaps[app.Creator],
			Manager:   GetUsersFromMap(userMaps, GetUserIds(app.Manager)),
			Status:    app.Status,
			CreatedAt: app.CreatedAt.Unix(),
			UpdatedAt: app.UpdatedAt.Unix(),
		},
		Secret: app.Secret,
	}

	return &ret, nil
}

///////////////////////////////// function /////////////////////////////////////////

func GenerateAppID(ctx context.Context) (uint64, error) {
	id, err := table.GetSequence(ctx)
	if err != nil {
		return 0, err
	}

	id = id % 1000000

	idStr := fmt.Sprintf("1%06d%02d", id, rand.Intn(99))
	userid, err := strconv.ParseUint(idStr, 10, 64)
	return userid, err
}

func GenerateAppSecret() string {
	return crypto.MD5Str(fmt.Sprintf("%d_%d", time.Now().UnixMilli(), rand.Intn(999999999)))
}

func IsAppManager(ctx context.Context, userid, appid uint64) (*st.TblAppInfo, error) {
	isNil, app, err := table.GetAppDetail(ctx, appid)
	if err != nil {
		return nil, err
	}

	if isNil {
		return app, errs.Newf(errs.RetWebNotFindApp, "not find app [%d]", appid)
	}

	if !types.InArrayUint64(GetUserIds(app.Manager), userid) {
		return app, errs.Newf(errs.RetWebMemberNotManager, "user is not manager of app [%s]", app.Name)
	}

	return app, nil
}

func GetAppidFromApps(apps []*st.TblAppInfo) []uint64 {
	ret := []uint64{}
	for _, app := range apps {
		ret = append(ret, app.Appid)
	}

	return ret
}

func GetAppBaseFromApp(userid uint64, app *st.TblAppInfo, userMaps map[uint64]*pb.UsersBase) *pb.AppBase {
	if app == nil || userMaps == nil {
		return nil
	}

	return &pb.AppBase{
		Appid:     app.Appid,
		Name:      app.Name,
		Intro:     app.Intro,
		IsManager: IsManager(userid, app.Manager),
		Creator:   userMaps[app.Creator],
		Manager:   GetUsersFromMap(userMaps, GetUserIds(app.Manager)),
		Status:    app.Status,
		CreatedAt: app.CreatedAt.Unix(),
		UpdatedAt: app.UpdatedAt.Unix(),
	}
}

func GetAppByAppid(apps []*st.TblAppInfo, appid uint64) *st.TblAppInfo {
	if len(apps) == 0 {
		return nil
	}

	for _, v := range apps {
		if v.Appid == appid {
			return v
		}
	}

	return nil
}
