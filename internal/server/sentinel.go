package server

import (
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
	"github.com/finecloud/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewSentinelDataSource(logger log.Logger, cd *conf.Data) {
	newFlowRuleDataSource(logger, cd)
}

func newFlowRuleDataSource(logger log.Logger, confData *conf.Data) {
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

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: scs,
		})

	if err != nil {
		helper.Fatalf("new sentinel nacos config client failed: %v", err)
	}

	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	nds, err := nacos.NewNacosDataSource(configClient, confData.GetSentinel().GetGroupName(),
		confData.GetSentinel().GetDataIdFlow(), h)
	if err != nil {
		helper.Fatalf("new sentinel nacos data_source failed: %v", err)
	}

	err = nds.Initialize()
	if err != nil {
		helper.Fatalf("sentinel datasource initial failed: %v", err)
	}
}

func newCircuitBreakerDataSource(logger log.Logger, confData *conf.Data) {
	log := log.NewHelper(logger)
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

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: scs,
		})

	if err != nil {
		log.Fatalf("new sentinel nacos config client failed: %v", err)
	}

	h := datasource.NewCircuitBreakerRulesHandler(datasource.CircuitBreakerRuleJsonArrayParser)
	nds, err := nacos.NewNacosDataSource(configClient, confData.GetSentinel().GetGroupName(),
		confData.GetSentinel().GetDataIdCb(), h)
	if err != nil {
		log.Fatalf("new sentinel nacos data_source failed: %v", err)
	}

	err = nds.Initialize()
	if err != nil {
		log.Fatalf("sentinel datasource initial failed: %v", err)
	}
}
