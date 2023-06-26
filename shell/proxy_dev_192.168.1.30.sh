#!/bin/bash


##go env
GO111MODULE="auto"
GOARCH="amd64"
GOBIN="/root/item/bin"
GOCACHE="/root/.cache/go-build"
GOENV="/root/.config/go/env"
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
GOVERSION="go1.20"
GCCGO="gccgo"
GOAMD64="v1"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="0"
GOMOD=""
GOWORK=""
CGO_CFLAGS="-O2 -g"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-O2 -g"
CGO_FFLAGS="-O2 -g"
CGO_LDFLAGS="-O2 -g"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -fno-caret-diagnostics -Qunused-arguments -Wl,--no-gc-sections -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build2085342253=/tmp/go-build -gno-record-gcc-switches"

#proxy server
export ZINX_ENV=dev

#golang
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn
export GOROOT=/usr/local/go
export GOPATH=/root/item
export GOBIN=/root/item/bin

#path
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/usr/local/go/bin:/root/item/bin:/root/bin:/usr/local/go/bin:/root/item/bin


##work path
TARGET="/sdadata/item/proxyserver"
echo "cd ${TARGET}"
cd ${TARGET}


##compile
echo -e "\n"
echo "compile proxy(develop)..."
ssh -t wangkebiao@192.168.1.39 "sudo /data/git/proxy/shell/proxy_dev_192.168.1.39.sh"
if [ $? -ne 0 ]; then
    echo "proxy compile failed!!!"
    exit 1
fi


##synchronize
echo -e "\n"
echo "proxy is rsync..."
rsync -av wangkebiao@192.168.1.39:/data/git/proxy/server/conf/*_dev* /sdadata/item/proxyserver/conf
if [ $? -ne 0 ]; then
    echo "proxy rsync config failed!!!"
    exit 1
fi
rsync -av wangkebiao@192.168.1.39:/data/git/proxy/server/proxy /sdadata/item/proxyserver
if [ $? -ne 0 ]; then
    echo "proxy rsync proxy failed!!!"
    exit 1
fi



##start
echo -e "\n"
echo "proxy is starting......"
PID=`ps aux | grep "proxyserver/proxy" | grep -v grep | awk '{printf $2}'`
if [ -n "$PID" ]; then
    echo "kill -SIGINT ${PID}"
    sudo kill -SIGINT $PID
    /bin/sleep 5
fi

/sdadata/item/proxyserver/proxy >> /sdadata/item/proxyserver/log/server.log 2>&1 &


if [ $? -eq 0 ]; then
    echo "proxy start is successfully!!!"
    echo "----------------------- `date '+%Y-%m-%d %H:%M:%S'` -----------------------" >> /sdadata/item/proxyserver/log/server.log
    exit 0
else
    echo "proxy start failed!!!"
    exit 1
fi


