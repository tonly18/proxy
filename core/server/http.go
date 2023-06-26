package server

import (
	"context"
	"fmt"
	"net/http"
	"proxy/server/config"
	"time"
)

type httpServer struct {
	Http http.Server
}

type HttpServerConfig struct {
	IP      string
	Port    int
	Handler http.Handler
}

func NewHttpServer(config *HttpServerConfig) *httpServer {
	return &httpServer{
		Http: http.Server{
			Addr:           fmt.Sprintf(`%s:%d`, config.IP, config.Port),
			Handler:        config.Handler,
			ReadTimeout:    5 * time.Second,  //从链接被接受开始,到request body完全读取为止。
			WriteTimeout:   10 * time.Second, //http:从request head读取结束开始到response write完成为止。
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *httpServer) Start() error {
	fmt.Printf("[HTTP SERVER START] SERVER IS STARTING AT HOST: %s, PORT: %d\n", config.HttpConfig.HttpHost, config.HttpConfig.HttpPort)

	return s.Http.ListenAndServe()
}

func (s *httpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Http.Shutdown(ctx)
}
