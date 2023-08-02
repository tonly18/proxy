package core

import (
	"context"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/config"
	"proxy/server/library/httpclient"
)

//OnConnStartFunc 上线
func OnConnStartFunc(conn ziface.IConnection) {
	//玩家在线数
	userOnLineNumber := conn.GetTCPServer().GetConnMgr().Len()

	//call http server
	client := httpclient.NewHttpClient(&httpclient.Config{})
	resp, err := client.Get(config.HttpConfig.GameServerOnOffLineAPI, map[string]any{
		"online": userOnLineNumber,
		"status": 1,
	}).SetHeader(map[string]any{
		"proxy_id":  conn.GetTCPServer().GetID(),
		"server_id": conn.GetServerId(),
		"user_id":   conn.GetUserId(),
		"client_ip": conn.GetRemoteIP(),
		"socket_id": conn.GetConnID(),
	}).Do()
	if err != nil {
		logger.Error(context.Background(), "[OnConnStartFunc] player on line error: ", err, ", api: ", config.HttpConfig.GameServerOnOffLineAPI)
	}
	if resp.GetHeaderCode() != 200 {
		logger.Error(context.Background(), "[OnConnStartFunc] player on line error, code: ", resp.GetHeaderCode(), ", api: ", config.HttpConfig.GameServerOnOffLineAPI)
	}

	logger.Info(conn, "[OnConnStartFunc] player on line, conn id: ", conn.GetConnID(), ", userId:", conn.GetUserId())
}

//OnConnStopFunc 下线
func OnConnStopFunc(conn ziface.IConnection) {
	userOnLineNumber := conn.GetTCPServer().GetConnMgr().Len() //玩家在线数
	proxyId := conn.GetTCPServer().GetID()
	serverId := conn.GetServerId()
	userId := conn.GetUserId()
	clientId := conn.GetRemoteIP()
	socketId := conn.GetConnID()

	//删除conn
	conn.Stop()

	//call http server
	client := httpclient.NewHttpClient(&httpclient.Config{})
	resp, err := client.Get(config.HttpConfig.GameServerOnOffLineAPI, map[string]any{
		"online": userOnLineNumber,
		"status": 0,
	}).SetHeader(map[string]any{
		"proxy_id":  proxyId,
		"server_id": serverId,
		"user_id":   userId,
		"client_ip": clientId,
		"socket_id": socketId,
	}).Do()
	if err != nil {
		logger.Error(context.Background(), "[OnConnStopFunc] player off line error: ", err, ", api: ", config.HttpConfig.GameServerOnOffLineAPI)
	}
	if resp.GetHeaderCode() != 200 {
		logger.Error(context.Background(), "[OnConnStopFunc] player off line error, code: ", resp.GetHeaderCode(), ", api: ", config.HttpConfig.GameServerOnOffLineAPI)
	}

	logger.Info(conn, "[OnConnStopFunc] player off line, conn id: ", socketId, ", userId:", userId)
}
