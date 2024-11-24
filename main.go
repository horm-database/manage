package main

import (
	"time"

	"github.com/horm/common/errs"
	"github.com/horm/common/log"
	"github.com/horm/manage/api"
	"github.com/horm/manage/auth"
	"github.com/horm/manage/srv"
	"github.com/horm/manage/srv/codec"

	_ "go.uber.org/automaxprocs"
)

func main() {
	server := srv.NewServer(api.ServerDesc)

	go func() {
		time.Sleep(10 * time.Millisecond)
	}()

	err := auth.InitWorkspaceID(codec.GCtx)
	if err != nil {
		panic(errs.Newf(errs.RetSystem, "init workspace id error: %v", err))
	}

	if err := server.Serve(); err != nil {
		log.Fatal(codec.GCtx, err)
	}
}
