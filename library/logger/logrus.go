package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"proxy/core/zinx/zconf"
	"runtime"
)

//init
func init() {
	initLogrus() //Logrus日志
}

//init logrus
func initLogrus() {
	//设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	//记录调用者
	logrus.SetReportCaller(false)
	//设置日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//日志输出
	if fileLog, err := os.OpenFile(zconf.GlobalObject.LogDir+"/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		runtime.SetFinalizer(fileLog, FinishClear)
		logrus.SetOutput(fileLog)
	} else {
		logrus.SetOutput(os.Stdout)
	}
}

//Error
func Error(ctx context.Context, args ...any) {
	logrus.WithFields(logrus.Fields{
		"TraceId":  ctx.Value("trace_id"),
		"ProxyId":  ctx.Value("proxy_id"),
		"ServerId": ctx.Value("server_id"),
		"UserId":   ctx.Value("user_id"),
		"IP":       ctx.Value("client_ip"),
	}).Error(args...)
}
func Errorf(ctx context.Context, format string, args ...any) {
	logrus.WithFields(logrus.Fields{
		"TraceId":  ctx.Value("trace_id"),
		"ProxyId":  ctx.Value("proxy_id"),
		"ServerId": ctx.Value("server_id"),
		"UserId":   ctx.Value("user_id"),
		"IP":       ctx.Value("client_ip"),
	}).Errorf(format, args...)
}

//Info
func Info(ctx context.Context, args ...any) {
	logrus.WithFields(logrus.Fields{
		"TraceId":  ctx.Value("trace_id"),
		"ProxyId":  ctx.Value("proxy_id"),
		"ServerId": ctx.Value("server_id"),
		"UserId":   ctx.Value("user_id"),
		"IP":       ctx.Value("client_ip"),
	}).Info(args...)
}
func Infof(ctx context.Context, format string, args ...any) {
	logrus.WithFields(logrus.Fields{
		"TraceId":  ctx.Value("trace_id"),
		"ProxyId":  ctx.Value("proxy_id"),
		"ServerId": ctx.Value("server_id"),
		"UserId":   ctx.Value("user_id"),
		"IP":       ctx.Value("client_ip"),
	}).Infof(format, args...)
}

//FinishClear 关闭文件句柄
func FinishClear(f *os.File) {
	f.Close()
}
