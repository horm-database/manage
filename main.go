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
