package utils

import "proxy/library/logger"

//程序退出前的清理工作
func FinishClear() {
	logger.FinishClear() //日志清理
}
