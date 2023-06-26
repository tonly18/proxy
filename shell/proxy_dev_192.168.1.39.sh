#!/bin/bash


##go env
GO111MODULE="auto"
GOARCH="amd64"
GOBIN="/root/item/bin"
GOCACHE="/home/wangkebiao/.cache/go-build"
GOENV="/home/wangkebiao/.config/go/env"
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOINSECURE=""
GOMODCACHE="/root/item/pkg/mod"
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
GOPATH="/root/item"
GOPRIVATE=""
GOPROXY="https://goproxy.cn"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
GOVCS=""
GOVERSION="go1.18.3"
GCCGO="gccgo"
GOAMD64="v1"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
GOMOD="/data/git/proxy/go.mod"
GOWORK=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build457313608=/tmp/go-build -gno-record-gcc-switches"

#proxy server
export ZINX_ENV=dev

#golang
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn
export GOROOT=/usr/local/go
export GOPATH=/root/item
export GOBIN=/root/item/bin


#path
export PATH=/usr/local/php74/bin:/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/go/bin:/root/item/bin:/root/bin



echo "[start] execute shell(192.168.1.39)"

##git
echo "cd /data/git/proxy"
cd /data/git/proxy
git pull origin develop:develop


##compile proxy
echo "start compiling proxy."
SERVER_NAME="proxy"
SOURCE="/data/git/proxy/server"
echo "cd ${SOURCE}"
cd ${SOURCE}
if [ -f "${SOURCE}/${SERVER_NAME}" ]; then
    echo "rm -f ${SOURCE}/${SERVER_NAME}"
    rm -f "${SOURCE}/${SERVER_NAME}"
fi

go mod download
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${SERVER_NAME} "${SOURCE}/main.go"
if [ $? -ne 0 ]; then
    echo "compiling proxy failed!!!"
    echo "[end]execute shell(192.168.1.39)"
    exit 1
else
    echo "compiling proxy successfully!!!"
    echo "[end]execute shell(192.168.1.39)"
    exit 0
fi
