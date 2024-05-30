package zlog

import (
	"fmt"
	"log/slog"
	"os"
	"proxy/library/command"
	"runtime"
	"time"
)

// SetLogFile sets the log file of StdZinxLog
func SetLogFile2(fileDir string, fileName string) {
	fs, err := os.OpenFile(fmt.Sprintf(`%v/%v`, fileDir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	runtime.SetFinalizer(fs, func(f *os.File) {
		f.Close()
	})
	runtime.KeepAlive(fs)

	opt := &slog.HandlerOptions{
		//AddSource: true,
		//Level:     slog.LevelError,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				itime := a.Value.Time()
				a.Value = slog.StringValue(itime.Format(time.DateTime))
			}
			return a
		},
	}
	slogHandler := slog.NewTextHandler(fs, opt)
	slog.SetDefault(slog.New(slogHandler))
}

func Infof2(format string, v ...any) {
	slog.Info(fmt.Sprintf(format, v...))
}

func Info2(v string, arg ...any) {
	slog.Info(fmt.Sprintf(`%v%v`, v, command.SliceJoinString(arg, " ")))
}

func Warnf2(format string, v ...interface{}) {
	slog.Warn(fmt.Sprintf(format, v...))
}

func Warn2(v string, arg ...any) {
	slog.Warn(fmt.Sprintf(`%v%v`, v, command.SliceJoinString(arg, " ")))
}

func Errorf2(format string, v ...interface{}) {
	slog.Error(fmt.Sprintf(format, v...))
}

func Error2(v string, arg ...any) {
	slog.Error(fmt.Sprintf(`%v%v`, v, command.SliceJoinString(arg, " ")))
}
