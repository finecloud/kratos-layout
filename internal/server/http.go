package server

import (
	v1 "github.com/LikeRainDay/kratos-layout/api/helloworld/v1"
	"github.com/LikeRainDay/kratos-layout/internal/conf"
	"github.com/LikeRainDay/kratos-layout/internal/service"
	"github.com/LikeRainDay/kratos-layout/pkg/casdoor_auth"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new HTTP server.
func NewHTTPServer(c *conf.Server, auth *conf.Data, logger log.Logger, greeter *service.GreeterService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			selector.Server(casdoor_auth.Server()).Match(
				casdoor_auth.NewWhiteListMatcher(auth.Casdoor.IgnoreUrls),
			).Build(),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)

	v1.RegisterGreeterHTTPServer(srv, greeter)

	return srv
}
