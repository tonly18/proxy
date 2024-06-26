package core

import (
	"context"
	"fmt"
	"github.com/tonly18/httpclient"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/config"
	"proxy/server/library"
	"proxy/utils"
)

// OnConnStartFunc 上线
func OnConnStartFunc(conn ziface.IConnection) {
	//玩家在线数
	userOnLineNumber, _ := conn.GetConnMgr().Len()

	logger.Info(context.Background(), "OnConnStartFunc. userOnLineNumber:", userOnLineNumber)
	return

	//call http server
	serverId := conn.GetProperty(utils.ServerID)
	userId := conn.GetProperty(utils.UserID)

	url := fmt.Sprintf(`%v%v`, library.NewPoll().Get(), config.HttpConfig.GameServerOnOffLineAPI)
	httpClient := httpclient.NewClient(&httpclient.Config{})
	resp, err := httpClient.Get(url, map[string]any{
		"online": userOnLineNumber,
		"status": 1,
	}).SetHeader(map[string]any{
		"proxy_id":  conn.GetTCPServer().GetID(),
		"server_id": serverId,
		"user_id":   userId,
		"client_ip": conn.GetRemoteIP(),
		"socket_id": conn.GetConnID(),
	}).Do()
	if err != nil {
		logger.Error(conn, "[OnConnStartFunc] player on line fail error: ", err, ", api: ", url)
		return
	}
	if resp.GetHeaderCode() != 200 {
		logger.Error(conn, "[OnConnStartFunc] player on line fail, code: ", resp.GetHeaderCode(), ", api: ", url)
	}

	logger.Info(conn, "[OnConnStartFunc] player on line success, conn id: ", conn.GetConnID(), ", userId:", userId)
}

// OnConnStopFunc 下线
func OnConnStopFunc(conn ziface.IConnection) {
	userOnLineNumber, _ := conn.GetConnMgr().Len() //玩家在线数
	proxyId := conn.GetTCPServer().GetID()
	serverId := conn.GetProperty(utils.ServerID)
	userId := conn.GetProperty(utils.UserID)
	clientId := conn.GetRemoteIP()
	socketId := conn.GetConnID()
	kick := conn.GetProperty(utils.Kick)

	logger.Info(context.Background(), "OnConnStopFunc. userOnLineNumber:", userOnLineNumber)
	return

	//call http server
	url := fmt.Sprintf(`%v%v`, library.NewPoll().Get(), config.HttpConfig.GameServerOnOffLineAPI)
	httpClient := httpclient.NewClient(&httpclient.Config{})
	resp, err := httpClient.Get(url, map[string]any{
		"online": userOnLineNumber,
		"status": 0,
		"kick":   kick,
	}).SetHeader(map[string]any{
		"proxy_id":  proxyId,
		"server_id": serverId,
		"user_id":   userId,
		"client_ip": clientId,
		"socket_id": socketId,
	}).Do()
	if err != nil {
		logger.Error(conn, "[OnConnStopFunc] player off line fail error: ", err, ", api: ", url)
		return
	}
	if resp.GetHeaderCode() != 200 {
		logger.Error(conn, "[OnConnStopFunc] player off line fail, code: ", resp.GetHeaderCode(), ", api: ", url)
	}

	logger.Info(conn, "[OnConnStopFunc] player off line success, conn id: ", socketId, ", userId:", userId)
}
