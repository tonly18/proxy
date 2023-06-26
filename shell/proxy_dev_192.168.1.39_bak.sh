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
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build1072476409=/tmp/go-build -gno-record-gcc-switches"

#proxy server
export ZINX_ENV=dev

#golang
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn
export GOROOT=/usr/local/go
export GOPATH=/root/item
export GOBIN=/root/item/bin


#path
export PATH=/usr/local/php74/bin:$PATH:$GOROOT/bin:$GOBIN


##git
git pull origin master:master


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
    exit 1
fi



##synchronize
echo "proxy is rsync......"
TARGET="/sdadata/item/proxyserver"
if [ -f "${SOURCE}/${SERVER_NAME}" ]; then
    rsync -a "${SOURCE}/${SERVER_NAME}" "${TARGET}/${SERVER_NAME}"
    rsync -a "${SOURCE}/conf/" "${TARGET}/conf/"
else
    echo "${SOURCE}/${SERVER_NAME} not exist!!!"
    exit 1
fi
if [ $? -ne 0 ]; then
    echo "proxy rsync failed!!!"
    exit 1
fi


##start
echo "proxy is starting......"
echo "cd ${TARGET}"
cd ${TARGET}
PID=`ps aux | grep "proxyserver/proxy" | grep -v grep | awk '{printf $2}'`
if [ -n "$PID" ]; then
    echo "kill -SIGINT ${PID}"
    kill -SIGINT $PID
    /bin/sleep 5
fi


##start
/sdadata/item/proxyserver/proxy >> /sdadata/item/proxyserver/log/server.log 2>&1 &

if [ $? -eq 0 ]; then
    echo "proxy start is successfully!!!"
    echo "----------------------- `date '+%Y-%m-%d %H:%M:%S'` -----------------------" >> /sdadata/item/proxyserver/log/server.log
    exit 0
else
    echo "proxy start failed!!!"
    exit 1
fi
