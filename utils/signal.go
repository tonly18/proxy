package utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"proxy/core/zinx/zconf"
	"proxy/server/config"
	"syscall"
	"time"
)

type Signal struct {
	sigChan    chan os.Signal
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewSignal() *Signal {
	s := &Signal{
		sigChan:    make(chan os.Signal),
		ctx:        nil,
		cancelFunc: nil,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())

	return s
}

func (s *Signal) notify() {
	signal.Notify(s.sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
}

func (s *Signal) Waiter() error {
	s.notify()

	for sig := range s.sigChan {
		switch sig {
		case syscall.SIGINT:
			//ctrl + c
			//fmt.Println("control signal int:", s)
			return nil
		case syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT:
			//fmt.Println("control signal(hup|term|quit):", s)
			return nil
		case syscall.SIGUSR1:
			//fmt.Println("control signal usr1:", s)
			return nil
		case syscall.SIGUSR2:
			//重新加载proxy配置文件
			if err := config.LoadProxyConfig(); err != nil {
				fmt.Println("reload proxy config file error:", err)
			} else {
				fmt.Println("reload proxy config file successfully")
			}
			continue
		default:
			//fmt.Println("other signal:", s)
		}
	}

	return nil
}

func (s *Signal) Cannel() {
	s.cancelFunc()
	if zconf.ZINX_ENV == "prod" {
		time.Sleep(5 * time.Second)
	} else {
		time.Sleep(5 * time.Second)
	}
}

func (s *Signal) GetCtx() context.Context {
	return s.ctx
}
