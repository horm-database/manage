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
	"github.com/horm-database/manage/model/table"
	st "github.com/horm-database/server/model/table"
	"github.com/samber/lo"
)

func AddApp(ctx context.Context, userid uint64, req *pb.AddAppRequest) (*pb.AddAppResponse, error) {
	if lo.IndexOf(req.Manager, userid) == -1 {
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

	managerUids := lo.Uniq(req.Manager)

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

	if lo.IndexOf(GetUserIds(app.Manager), userid) == -1 {
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
