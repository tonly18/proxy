package utils

import (
	"context"
	"os"
	"os/signal"
	"proxy/core/zinx/zconf"
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

	for s := range s.sigChan {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			//fmt.Println("Program Exit...", s)
			return nil
		case syscall.SIGUSR1:
			//fmt.Println("usr1 signal", s)
			return nil
		case syscall.SIGUSR2:
			//fmt.Println("usr2 signal", s)
			return nil
		default:
			//fmt.Println("other signal", s)
		}
	}

	return nil
}

func (s *Signal) Cannel() {
	s.cancelFunc()
	if zconf.ZINX_ENV == "prod" {
		time.Sleep(10 * time.Second)
	} else {
		time.Sleep(5 * time.Second)
	}
}

func (s *Signal) GetCtx() context.Context {
	return s.ctx
}
