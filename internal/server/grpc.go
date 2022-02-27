package server

import (
	v1 "github.com/LikeRainDay/kratos-layout/api/helloworld/v1"
	"github.com/LikeRainDay/kratos-layout/internal/conf"
	"github.com/LikeRainDay/kratos-layout/internal/data"
	"github.com/LikeRainDay/kratos-layout/internal/service"
	"github.com/LikeRainDay/kratos-layout/pkg/casdoor_auth"
	"github.com/LikeRainDay/kratos-layout/pkg/log_id"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, logger log.Logger, greeter *service.GreeterService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			log_id.MetadataServer(),
			casdoor_auth.Server(),
			logging.Server(logger),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(data.MetricReqDurationHistogram)),
				metrics.WithRequests(prom.NewCounter(data.MetricReqTotalCounter))),
		),
		grpc.Logger(logger),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
