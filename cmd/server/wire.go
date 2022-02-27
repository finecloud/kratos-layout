// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/LikeRainDay/kratos-layout/internal/biz"
	"github.com/LikeRainDay/kratos-layout/internal/conf"
	"github.com/LikeRainDay/kratos-layout/internal/data"
	"github.com/LikeRainDay/kratos-layout/internal/server"
	"github.com/LikeRainDay/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/gorm"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger, *gorm.Config, webMultiTenancyOption *http.WebMultiTenancyOption) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
