package global

import (
	"os"
	"proxy/core/zinx/ziface"
	"strings"
)

//运行环境: local、dev、test、prod
var PROXY_SERVER_ENV string

//当前工作目录
var PROXY_SERVER_WORK_PATH_ENV string

//TCP包命令
const (
	CMD_UP_DOLIGON = 1
	CMD_UP_PING    = 34
	CMD_UP_PROTO   = 35

	CMD_DOWN_KICK_OUT      = 25
	CMD_DOWN_KICK_OUT_SELF = 26
	CMD_DOWN_PONG          = 46
)

func init() {
	PROXY_SERVER_ENV = strings.ToLower(os.Getenv("ZINX_ENV"))
	PROXY_SERVER_WORK_PATH_ENV, _ = os.Getwd()
}

//TCP Server
var globalTcpServer ziface.IServer

//SetTCPServer 获取tcpServer
func SetTCPServer(tcpServer ziface.IServer) ziface.IServer {
	globalTcpServer = tcpServer

	return globalTcpServer
}

//GetTCPServer 获取tcpServer
func GetTCPServer() ziface.IServer {
	return globalTcpServer
}
