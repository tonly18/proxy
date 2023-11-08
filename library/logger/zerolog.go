package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"proxy/core/zinx/zconf"
	"proxy/library/logger/hook"
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

	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger = zerolog.New(fs).With().Timestamp().Logger().Hook(&hook.ZeroLogHook{})
}

// Debug
func Debug(ctx context.Context, args ...any) {
	logger.Debug().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msg(fmt.Sprint(args...))
}
func Debugf(ctx context.Context, format string, args ...any) {
	logger.Debug().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msgf(format, args...)
}

// Info
func Info(ctx context.Context, args ...any) {
	logger.Info().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msg(fmt.Sprint(args...))
}
func Infof(ctx context.Context, format string, args ...any) {
	logger.Info().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msgf(format, args...)
}

// Warning
func Warning(ctx context.Context, args ...any) {
	logger.Warn().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msg(fmt.Sprint(args...))

}
func Warningf(ctx context.Context, format string, args ...any) {
	logger.Warn().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msgf(format, args...)
}

// Error
func Error(ctx context.Context, args ...any) {
	logger.Error().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msg(fmt.Sprint(args...))
}
func Errorf(ctx context.Context, format string, args ...any) {
	logger.Error().Fields(map[string]any{
		"gameserver_id": ctx.Value("gameserver_id"),
		"proxy_id":      ctx.Value("proxy_id"),
		"server_id":     ctx.Value("server_id"),
		"user_id":       ctx.Value("user_id"),
		"client_ip":     ctx.Value("client_ip"),
		"trace_id":      ctx.Value("trace_id"),
	}).Msgf(format, args...)
}
