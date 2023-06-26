package tcpserver

import (
	"fmt"
	"proxy/core/zinx/znet"
	"proxy/server/global"
	"proxy/server/tserver/core"
	"proxy/server/tserver/router"
	"proxy/utils"
)

//StartTCPServer new tcp server
func StartTCPServer(sig *utils.Signal) {
	//0 tcp
	tcpServer := global.SetTCPServer(znet.NewServer())

	//1 信号
	go func() {
		select {
		case <-sig.GetCtx().Done():
			fmt.Println("[TCP SERVER StartTCPServer STOP] TCP SERVER IS STOPPING!")
			tcpServer.Stop()
		}
	}()

	//2 添加router
	router.InitRouter(tcpServer)

	//3 Hook
	tcpServer.SetOnConnStart(core.OnConnStartFunc)
	tcpServer.SetOnConnStop(core.OnConnStopFunc)

	//4 启动zinx服务
	tcpServer.Serve()

	fmt.Println("[TCP SERVER StartTCPServer STOP] TCP SERVER IS STOPPED!")
}
