package global

import (
	"github.com/joho/godotenv"
	"os"
	"proxy/core/zinx/ziface"
)

var (
	// 运行环境: local、dev、test、prod
	PROXY_SERVER_ENV string
	// 当前工作目录
	PROXY_SERVER_PATH string
	// 日志目录
	PROXY_SERVER_LOG_PATH string
	// 配置文件目录
	PROXY_SERVER_CONFIG_PATH string
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err = godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	PROXY_SERVER_ENV = os.Getenv("ZINX_ENV")
	PROXY_SERVER_PATH = os.Getenv("PROXY_SERVER_PATH")
	PROXY_SERVER_LOG_PATH = os.Getenv("PROXY_SERVER_LOG_PATH")
	PROXY_SERVER_CONFIG_PATH = os.Getenv("PROXY_SERVER_CONFIG_PATH")
}

// 全局 TCP Server
var globalTcpServer ziface.IServer

// SetTCPServer 获取tcpServer
func SetTCPServer(tcpServer ziface.IServer) ziface.IServer {
	globalTcpServer = tcpServer

	return globalTcpServer
}

// GetTCPServer 获取tcpServer
func GetTCPServer() ziface.IServer {
	return globalTcpServer
}
