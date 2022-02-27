package main

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	nConfig "github.com/go-kratos/nacos/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

const (
	EnvNacosAddr          = "NacosAddr"
	EnvNacosPort          = "NacosPort"
	EnvNacosNamespaceID   = "NacosNamespaceID"
	EnvNacosLogLevel      = "NacosLogLevel"
	EnvNacosGroupName     = "NacosGroupName"
	EnvNacosLogMaxAge     = "NacosLogMaxAge"
	EnvNacosLogRotateTime = "NacosLogRotateTime"
)

type Config struct {
	Service struct {
		Name    string
		Version string
	}

	Data struct {
		Nacos struct {
			Addr          string `mapstructure:"addr"`
			Port          uint64 `mapstructure:"port"`
			NameSpaceId   string `mapstructure:"namespace_id"`
			LogRotateTime string `mapstructure:"log_rotate_time"`
			LogMaxAge     int64  `mapstructure:"log_max_age"`
			LogLevel      string `mapstructure:"log_level"`
			ClusterName   string `mapstructure:"cluster_name"`
			GroupName     string `mapstructure:"group_name"`
			Weight        string `mapstructure:"weight"`
		}
	}
}

var ConfigStruct Config

func NewNacosConfigSource() config.Source {
	applicationConfig := ConfigStruct.Service
	name := applicationConfig.Name
	fileType := "yaml"

	cc := constant.NewClientConfig(
		constant.WithNamespaceId(getNacosNamespaceID()),
		constant.WithMaxAge(getNacosLogMaxAge()),
		constant.WithRotateTime(getNacosLogRotateTime()),
		constant.WithLogLevel(getNacosLogLevel()))

	scs := []constant.ServerConfig{
		*constant.NewServerConfig(getNacosAddr(), getNacosPort()),
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: scs,
		})

	if err != nil {
		panic("new nacosConfig config client failed:" + err.Error())
	}
	return nConfig.NewConfigSource(configClient,
		nConfig.DataID(fmt.Sprintf("%s.%s", name, fileType)), nConfig.Group(getNacosGroupName()))

}

func InitConfigs(path string) {
	defaultName := "config"
	configType := "yml"
	defaultPath := path
	defaultViper := viper.New()
	defaultViper.SetConfigName(defaultName)
	defaultViper.AddConfigPath(defaultPath)
	defaultViper.SetConfigType(configType)
	err := defaultViper.ReadInConfig()
	if err != nil {
		panic("read " + defaultName + " config failed:" + err.Error())
	}
	for k, v := range defaultViper.AllSettings() {
		viper.SetDefault(k, v)
	}
	localConfigStruct := Config{}
	_ = viper.Unmarshal(&localConfigStruct)
	ConfigStruct = localConfigStruct

	viper.AutomaticEnv()
}

func GetApplicationName() string {
	applicationConfig := ConfigStruct.Service
	return applicationConfig.Name
}

func GetVersion() string {
	applicationConfig := ConfigStruct.Service
	return applicationConfig.Version
}

func getNacosNamespaceID() string {
	nacosNameSpaceId := viper.GetString(EnvNacosNamespaceID)
	if nacosNameSpaceId != "" {
		return nacosNameSpaceId
	}
	return ConfigStruct.Data.Nacos.NameSpaceId
}

func getNacosAddr() string {
	nacosAddr := viper.GetString(EnvNacosAddr)
	if nacosAddr != "" {
		return nacosAddr
	}
	return ConfigStruct.Data.Nacos.Addr
}

func getNacosPort() uint64 {
	nacosPort := viper.GetUint64(EnvNacosPort)
	if nacosPort != 0 {
		return nacosPort
	}
	return ConfigStruct.Data.Nacos.Port
}

func getNacosGroupName() string {
	nacosGroupName := viper.GetString(EnvNacosGroupName)
	if nacosGroupName != "" {
		return nacosGroupName
	}
	return ConfigStruct.Data.Nacos.GroupName
}

func getNacosLogLevel() string {
	logLevel := viper.GetString(EnvNacosLogLevel)
	if logLevel != "" {
		return logLevel
	}
	return ConfigStruct.Data.Nacos.LogLevel
}

func getNacosLogRotateTime() string {
	nacosLogRotateTime := viper.GetString(EnvNacosLogRotateTime)
	if nacosLogRotateTime != "" {
		return nacosLogRotateTime
	}
	return ConfigStruct.Data.Nacos.LogRotateTime
}

func getNacosLogMaxAge() int64 {
	nacosLogMaxAge := viper.GetInt64(EnvNacosLogMaxAge)
	if nacosLogMaxAge != -1 {
		return nacosLogMaxAge
	}
	return ConfigStruct.Data.Nacos.LogMaxAge
}
