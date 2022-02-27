# 坏境配置

## Go环境

1. 你需要[Download Golang](https://golang.org/)
2. 国内仓库环境配置`export GOPROXY=https://goproxy.cn`
3. 安装依赖，项目根目录执行`make all`
4. [protoc 安装](https://grpc.io/docs/protoc-installation/)
5. 查看go 环境变量是否正常`go env`



## 异常

### protoc 找不到

> 如果出现protoc-gen-go-http未发现等，可以尝试以下方法。

1. 将GO对应的`$GOPATH/bin`增加到`~/.bash_profile`中
2. 或者直接将`$GOPATH/bin`下的文件cp 到 `/usr/local/bin/`目录下。

### 替换Https 更改了git方式拉取

git config --global url."git@github.com:".insteadOf "https://github.com/"

### 无法从 `https://golang.org/x/tools`下获取内容
> 最快速的方式[squid server](https://hub.docker.com/r/ubuntu/squid)代理服务。


1. [最有效](https://blog.csdn.net/yxf771hotmail/article/details/88233857)

```shell
# 在 $GOPATH/src/github.com 下执行
# mkdir -p $GOPATH/src/github.com && cd $GOPATH/src/github.com 
git clone https://github.com/mdempsky/gocode
git clone https://github.com/uudashr/gopkgs/cmd/gopkgs
git clone https://github.com/ramya-rao-a/go-outline
git clone https://github.com/acroca/go-symbols
git clone https://github.com/go-delve/delve/cmd/dlv
git clone https://github.com/stamblerre/gocode
git clone https://github.com/rogpeppe/godef
git clone https://github.com/sqs/goreturns

# 在 $GOPATH/src/golang.org/x 下执行
# mkdir -p $GOPATH/src/golang.org/x && cd $GOPATH/src/golang.org/x
git clone git@github.com:golang/tools.git

# 在 $GOPATH 下执行
# cd $GOPATH 
go install github.com/mdempsky/gocode
go install github.com/uudashr/gopkgs/cmd/gopkgs
go install github.com/ramya-rao-a/go-outline
go install github.com/acroca/go-symbols
go install golang.org/x/tools/cmd/guru
go install golang.org/x/tools/cmd/gorename
go install github.com/go-delve/delve/cmd/dlv
go install github.com/stamblerre/gocode
go install github.com/rogpeppe/godef
go install github.com/sqs/goreturns
go install golang.org/x/lint/golint

```

2. [手动下载方式](https://blog.csdn.net/ckx178/article/details/89156585)
3. [解决方法link](https://www.cnblogs.com/shockerli/p/go-get-golang-org-x-solution.html)

# 引用

[Go Pkg](https://pkg.go.dev/search?q=protoc-gen-go-http&m=package)