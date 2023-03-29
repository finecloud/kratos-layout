package server

import (
	"github.com/gorilla/handlers"
	"google.golang.org/protobuf/proto"
	http2 "net/http"

	"encoding/json"
	"github.com/author_name/project_name/internal/conf"
	"github.com/author_name/project_name/internal/data"
	"github.com/author_name/project_name/pkg/casdoor_auth"
	"github.com/author_name/project_name/pkg/log_id"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewHTTPServer new HTTP server.
func NewHTTPServer(c *conf.Server, auth *conf.Data, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			log_id.MetadataServer(),
			selector.Server(casdoor_auth.Server()).Match(
				casdoor_auth.NewWhiteListMatcher(auth.Casdoor.IgnoreUrls),
			).Build(),
			logging.Server(logger),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(data.MetricReqDurationHistogram)),
				metrics.WithRequests(prom.NewCounter(data.MetricReqTotalCounter))),
		),
		http.ResponseEncoder(func(w http2.ResponseWriter, request *http2.Request, v interface{}) error {
			pbJson := func(v interface{}) ([]byte, error) {
				mOption := protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
					UseEnumNumbers:  true,
				}
				if m, ok := v.(proto.Message); ok {
					return mOption.Marshal(m)
				} else {
					return json.Marshal(v)
				}
			}
			data, err := pbJson(struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			}{
				Code:    0,
				Message: "success",
				Data:    v,
			})
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(data)
			if err != nil {
				return err
			}
			return nil
		}),
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
	opts = append(opts, http.Filter(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Card", "GroupId", "X-Requested-With", "Access-Control-Allow-Credentials", "User-Agent", "Content-Length", "Authorization"}),
	)))
	srv := http.NewServer(opts...)

	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	srv.Handle("/metrics", promhttp.Handler())

	return srv
}
