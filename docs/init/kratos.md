# Kratos
> [Kratos](https://go-kratos.dev/en/docs/getting-started/start) 作为B站的Golang 脚手架 有一定的封装和管理思想，因此采用此框架进行开发

## Install Kratos

```
go get -u github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## Create a service

```
# Create a template project
kratos new -r https://github.com/finecloud/kratos-layout.git github.com/finecloud/server 

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

make all
make run
```

## Generate other auxiliary files by Makefile

```
# Download and update dependencies
make init
# Generate API swagger json files by proto file
make swagger
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```

## Automated Initialization (wire)

```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker

```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 9000:9000 <your-docker-image-name>
``
