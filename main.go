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

package main

import (
	"github.com/horm-database/common/errs"
	"github.com/horm-database/common/log"
	"github.com/horm-database/manage/api"
	"github.com/horm-database/manage/auth"
	"github.com/horm-database/manage/srv"
	"github.com/horm-database/manage/srv/codec"
	_ "go.uber.org/automaxprocs"
)

func main() {
	server := srv.NewServer(api.ServerDesc)

	err := auth.InitWorkspaceID(codec.GCtx)
	if err != nil {
		panic(errs.Newf(errs.ErrSystem, "init workspace id error: %v", err))
	}

	if err := server.Serve(); err != nil {
		log.Fatal(codec.GCtx, err)
	}
}
