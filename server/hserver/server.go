package httpserver

import (
	"fmt"
	"net/http"
	"proxy/core/server"
	"proxy/server/config"
	"proxy/server/hserver/router"
	"proxy/utils"
)

func StartHttpServer(sig *utils.Signal) {
	//create http server
	config := server.HttpServerConfig{
		IP:      config.HttpConfig.HttpHost,
		Port:    config.HttpConfig.HttpPort,
		Handler: router.InitRouter(),
	}
	httpserver := server.NewHttpServer(&config)

	//signal
	go func() {
		select {
		case <-sig.GetCtx().Done():
			fmt.Println("[HTTP SERVER STOP] SERVER IS STOPPING!")
			httpserver.Stop()
		}
	}()

	//listen
	if err := httpserver.Start(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf(`"[HTTP SERVER ERROR] START ERROR: %v`, err))
	}

	fmt.Println("[HTTP SERVER STOP] SERVER IS STOPPED!")
}
