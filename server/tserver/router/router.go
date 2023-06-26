package router

import (
	"proxy/core/zinx/ziface"
	"proxy/server/global"
	"proxy/server/tserver/handler"
)

func InitRouter(s ziface.IServer) {
	//test handler
	s.AddRouter(0, &handler.TestRouter{})

	//proto handler
	s.AddRouter(global.CMD_UP_PROTO, &handler.PublicRouter{})

	//login handler
	s.AddRouter(global.CMD_UP_DOLIGON, &handler.LoginRouter{})

	//ping handler
	s.AddRouter(global.CMD_UP_PING, &handler.PingRouter{})
}
