package log_id

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	mwMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
)

func MetadataClient(ops ...mwMetadata.Option) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if logId := ctx.Value("K_LOGID"); logId != nil {
				ctx = metadata.AppendToClientContext(ctx, "x-md-global-log_id", fmt.Sprint(logId))
			}
			return mwMetadata.Client(ops...)(handler)(ctx, req)
		}
	}
}

// MetadataServer Server is middleware server-side metadata.
func MetadataServer(ops ...mwMetadata.Option) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		afterHandler := func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			md := metadata.Metadata{}
			var ok bool
			var logId string
			if md, ok = metadata.FromServerContext(ctx); ok {
				logId = md.Get("x-md-global-log_id")
			}
			if !ok || len(logId) == 0 {
				nLogId, _ := LocalGen()
				md.Set("x-md-global-log_id", fmt.Sprint(nLogId))
				ctx = metadata.NewServerContext(ctx, md)
			}

			return handler(ctx, req)
		}
		return mwMetadata.Server(ops...)(afterHandler)
	}
}

// LogID 使用方法：注册到kratos logger中
func LogID() log.Valuer {
	return func(ctx context.Context) interface{} {
		if md, ok := metadata.FromServerContext(ctx); ok {
			logId := md.Get("x-md-global-log_id")
			if len(logId) != 0 {
				return logId
			}
		}
		return "unknown_log_id"
	}
}
