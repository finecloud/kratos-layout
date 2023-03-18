package server

import (
	"github.com/author_name/project_name/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	nRegistry "github.com/go-kratos/nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosRegistry(logger log.Logger, confData *conf.Data) registry.Registrar {
	helper := log.NewHelper(logger)
	cc := constant.NewClientConfig(
		constant.WithNamespaceId(confData.GetNacos().GetNamespaceId()),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithMaxAge(confData.GetNacos().GetLogMaxAge()),
		constant.WithRotateTime(confData.GetNacos().GetLogRotateTime()),
		constant.WithLogLevel(confData.GetNacos().GetLogLevel()))

	scs := []constant.ServerConfig{
		*constant.NewServerConfig(
			confData.GetNacos().GetAddr(),
			confData.GetNacos().GetPort()),
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: scs,
		})

	if err != nil {
		helper.Fatalf("new nacos registry client failed: %v", err)
	}

	return nRegistry.New(namingClient,
		nRegistry.WithCluster(confData.GetNacos().GetClusterName()),
		nRegistry.WithGroup(confData.GetNacos().GetGroupName()),
		nRegistry.WithWeight(confData.GetNacos().GetWeight()))
}
