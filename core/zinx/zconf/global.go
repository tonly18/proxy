package zconf

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	//运行环境: local、dev、test、prod
	ZINX_ENV string
	//当前工作目录
	ZINX_WORK_PATH string
	//日志目录
	ZINX_LOG_PATH string
	//配置文件目录
	ZINX_CONFIG_PATH string
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err = godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	ZINX_ENV = os.Getenv("ZINX_ENV")
	ZINX_WORK_PATH = os.Getenv("PROXY_SERVER_PATH")
	ZINX_LOG_PATH = os.Getenv("PROXY_SERVER_LOG_PATH")
	ZINX_CONFIG_PATH = os.Getenv("PROXY_SERVER_CONFIG_PATH")
}
