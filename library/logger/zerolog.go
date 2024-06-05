package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"proxy/core/zinx/zconf"
	"proxy/library/logger/hook"
	"proxy/utils"
	"runtime"
	"time"
)

// logger
var logger zerolog.Logger

// init
func init() {
	fs, err := os.OpenFile(fmt.Sprintf(`%v/proxy.log`, zconf.GlobalObject.LogDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf(`open log file is error: %v`, err))
	}
	runtime.SetFinalizer(fs, func(f *os.File) {
		f.Close()
	})
	runtime.KeepAlive(fs)

	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger = zerolog.New(fs).With().Timestamp().Logger().Hook(&hook.ZeroLogHook{})
}

// Debug
func Debug(ctx context.Context, args ...any) {
	logger.Debug().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msg(fmt.Sprint(args...))
}
func Debugf(ctx context.Context, format string, args ...any) {
	logger.Debug().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msgf(format, args...)
}

// Info
func Info(ctx context.Context, args ...any) {
	logger.Info().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msg(fmt.Sprint(args...))
}
func Infof(ctx context.Context, format string, args ...any) {
	logger.Info().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msgf(format, args...)
}

// Warning
func Warning(ctx context.Context, args ...any) {
	logger.Warn().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msg(fmt.Sprint(args...))

}
func Warningf(ctx context.Context, format string, args ...any) {
	logger.Warn().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msgf(format, args...)
}

// Error
func Error(ctx context.Context, args ...any) {
	logger.Error().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msg(fmt.Sprint(args...))
}
func Errorf(ctx context.Context, format string, args ...any) {
	logger.Error().Fields(map[string]any{
		"proxy_id":      ctx.Value(utils.ProxyID),
		"server_id":     ctx.Value(utils.ServerID),
		"user_id":       ctx.Value(utils.UserID),
		"client_ip":     ctx.Value(utils.ClientIP),
		"trace_id":      ctx.Value(utils.TraceID),
		"gameserver_id": ctx.Value(utils.GameServerID),
	}).Msgf(format, args...)
}
