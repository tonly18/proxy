package zconf

import (
	"os"
	"strings"
)

//运行环境: local、dev、test、prod
var ZINX_ENV string

//当前工作目录
var ZINX_WORK_PATH_ENV string

func init() {
	ZINX_ENV = strings.ToLower(os.Getenv("ZINX_ENV"))
	ZINX_WORK_PATH_ENV, _ = os.Getwd()
}
