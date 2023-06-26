package router

import (
	"net/http"
	"proxy/server/hserver/controller"
)

func InitRouter() *http.ServeMux {
	handler := http.NewServeMux()

	//public
	handler.HandleFunc("/v1/public", controller.WrapHandle(controller.PublicController))

	//monitor
	handler.HandleFunc("/v1/m/goroutine", controller.MonitorGroutionController)
	handler.HandleFunc("/v1/m/memory", controller.MonitorMemoryController)
	//test
	handler.HandleFunc("/v1/test", controller.TestController)

	//return
	return handler
}
