package codec

import (
	"context"

	cc "github.com/horm-database/common/codec"
	"github.com/horm-database/common/log/logger"
)

var (
	GCtx context.Context
)

// InitGlobalContext 获取初始化的 context、msg
func InitGlobalContext(env, container, callee string) {
	ctx, msg := cc.NewMessage(context.Background())

	msg.WithEnv(env)
	msg.WithCalleeServiceName(callee)

	msg.WithLogger(logger.DefaultLogger.With(
		logger.Field{"callee", callee},
		logger.Field{"container", container}))

	GCtx = ctx
}
