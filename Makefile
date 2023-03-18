GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
MICRO_MOD=$(shell cat go.mod|sed -n '1p'|awk '{print $$2}')
GO_OPTS=$(shell echo $(API_PROTO_FILES) | awk '{for(i=1;i<=NF;i++){print "--go_opt=M"$$i"=$(MICRO_MOD)/"$$i}}' \
						| awk -F/ 'OFS="/"{$$NF="";print}'	\
						| sed 's/.$$//g')

GO_INTERNAL_API_OPTS=$(shell echo $(INTERNAL_PROTO_FILES) | awk '{for(i=1;i<=NF;i++){print "--go_opt=M"$$i"=$(MICRO_MOD)/"$$i}}' \
						| awk -F/ 'OFS="/"{$$NF="";print}'	\
						| sed 's/.$$//g')

.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	go get -u github.com/google/wire/cmd/wire
	go get github.com/swaggo/swag/cmd/swag@v1.7.8
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go mod download

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
               --proto_path=./third_party \
               --go_out=paths=source_relative:. \
               --go-errors_out=paths=source_relative:. \
               $(API_PROTO_FILES)

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=. \
	       --proto_path=./third_party \
	       $(GO_INTERNAL_API_OPTS)	\
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=. \
		   --proto_path=./third_party \
		   $(GO_OPTS)		\
 	       --go_out=paths=source_relative:. \
 	       --go-http_out=paths=source_relative:. \
 	       --go-grpc_out=paths=source_relative:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)  -X main.Name=$(NAME)" -o ./bin/ ./...

.PHONY: docker
# 构建docker中的exec
docker:
	mkdir -p bin/ \
	&& GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s -extldflags "-static" -X main.Version=$(VERSION) -X main.Name=$(NAME)' -o ./bin/ ./...

.PHONY: run
# go run
run:
	go run -ldflags "-X main.Version=$(VERSION) -X main.Name=$(NAME)" ./cmd/$(NAME)

.PHONY: generate
# generate
generate:
	go generate ./...

.PHONY: swagger
# 编译生成swagger.json文件
swagger:
	protoc --proto_path=. \
		   --proto_path=./third_party \
		   --openapiv2_out . \
		   --openapiv2_opt logtostderr=true \
           --openapiv2_opt json_names_for_fields=false \
           --openapiv2_opt enums_as_ints=true	\
           $(API_PROTO_FILES)

.PHONY: all
# generate all
all:
	make init
	make errors;
	make config;
	make generate;
	make api;
	make swagger;
	go mod tidy;

.PHONY: service
service:
	kratos proto server api/$(model)/$(name).proto -t internal/service

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
