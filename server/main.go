package main

import (
	"fmt"
	httpserver "proxy/server/hserver"
	tcpserver "proxy/server/tserver"
	"proxy/utils"
)

func main() {
	//创建监听退出chan,监听指定信号 ctrl+c kill
	sig := utils.NewSignal()

	/**********************************************************
	 * TCP SERVER
	 **********************************************************/
	go func() {
		fmt.Println("[TCP SERVER IS START]")
		tcpserver.StartTCPServer(sig)
	}()

	/**********************************************************
	 * HTTP SERVER
	 **********************************************************/
	go func() {
		fmt.Println("[HTTP SERVER IS START]")
		httpserver.StartHttpServer(sig)
	}()

	//Block: wait for signal
	if err := sig.Waiter(); err == nil {
		sig.Cannel()
	}

	//clear
	utils.FinishClear()

	//Finish
	fmt.Println("[All Services Are Stopped!!!]")
}
