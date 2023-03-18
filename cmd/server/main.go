package main

import (
	"flag"
	"fmt"
	"github.com/author_name/project_name/internal/conf"
	"github.com/author_name/project_name/pkg/casdoor_auth"
	"github.com/author_name/project_name/pkg/color"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	//// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

const _UI = `
 ______ _____  ______   ______    ______  _        ______   _    _   _____   
| |      | |  | |  \ \ | |       | |     | |      / |  | \ | |  | | | | \ \  
| |----  | |  | |  | | | |----   | |     | |   _  | |  | | | |  | | | |  | | 
|_|     _|_|_ |_|  |_| |_|____   |_|____ |_|__|_| \_|__|_/ \_|__|_| |_|_/_/  

service : %s
version : %s
`

func init() {
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, http *http.Server, rg registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(GetApplicationName()),
		kratos.Version(GetVersion()),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			http,
		),
		kratos.Registrar(rg),
	)
}

func main() {
	flag.Parse()
	InitConfigs(flagconf)
	fmt.Println(color.Blue(fmt.Sprintf(_UI, GetApplicationName(), GetVersion())))

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"Service.id", id,
		"Service.name", GetApplicationName(),
		"Service.version", GetVersion(),
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	c := config.New(
		config.WithLogger(logger),
		config.WithSource(
			NewNacosConfigSource(),
		),
		config.WithDecoder(func(value *config.KeyValue, m map[string]interface{}) error {
			return yaml.Unmarshal(value.Value, m)
		}),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	casdoor := bc.Data.Casdoor
	casdoor_auth.Client(
		casdoor.Endpoint,
		casdoor.ClientId,
		casdoor.ClientSecret,
		casdoor.JwtSecret,
		casdoor.OrganizationName,
		casdoor.ApplicationName,
	)

	app, cleanup, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
