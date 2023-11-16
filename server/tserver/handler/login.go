package handler

import (
	"fmt"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/global"
	"proxy/server/utils/pack"
	"time"
)

// LoginRouter Struct
type LoginRouter struct {
	BaseHandler
}

func (h *LoginRouter) Handle(request ziface.IRequest) error {
	//开始时间
	start := time.Now()

	//data := request.GetData()
	//dp := pack.NewDataPackLogin()
	//msg, err := dp.UnPack(data)
	//if err != nil {
	//	request.GetConnection().Stop()
	//	return fmt.Errorf(`[login handler] unpack error: %v`, err)
	//}
	//if msg.GetServer() == nil {
	//	request.GetConnection().Stop()
	//	return errors.New("[login handler] msg server id is nil")
	//}

	//调用gameServer验证登录
	//url := fmt.Sprintf(`%v%v`, library.NewPoll().Get(), config.HttpConfig.GameServerDoLoginAPI)
	//httpClient := httpclient.NewClient(&httpclient.Config{})
	//resp, err := httpClient.NewRequest("POST", url, request.GetData()).SetHeader(map[string]any{
	//	"Content-Type": "application/octet-stream",
	//	"proxy_id":     request.GetConnection().GetTCPServer().GetID(),
	//	"server_id":    0,
	//	"user_id":      0,
	//	"client_ip":    request.GetConnection().GetRemoteIP(),
	//	"trace_id":     request.GetTraceId(),
	//}).Do()
	//if err != nil {
	//	return fmt.Errorf(`login call http server is error: %v`, err)
	//}
	//request.SetAargs("gameserver_id", resp.GetDataFromHeader("gameserver_id")) //game server id
	//if resp.GetHeaderCode() != 200 {
	//	return fmt.Errorf(`login call http server is code:%v`, resp.Response.StatusCode)
	//}
	//
	////处理返回数据
	//userId := cast.ToUint64(resp.GetDataFromHeader("user_id"))     //玩家ID
	//serverId := cast.ToUint32(resp.GetDataFromHeader("server_id")) //区服ID
	//if uin == 0 || userId == 0 || serverId == 0 {
	//	logger.Error(request, "[login handler] login fail.")
	//	return err
	//}

	//设置conn属性
	conn := request.GetConnection()
	//conn.SetProxyId(conn.GetTCPServer().GetID()) //网关ID
	//conn.SetServerId(serverId)                   //区服ID
	//conn.SetUserId(userId)                       //玩家ID

	//添加登录后的玩家到connManager
	userId := uint64(111)
	conn.SetUserId(userId)
	fmt.Println("login===========")
	if connection, _ := conn.GetTCPServer().GetConnMgr().GetConnByUserId(userId); connection != nil {
		fmt.Println("login2222222===========")
		//踢掉原connection,并推送消息给客户端
		downMsg := pack.NewMessageDown(global.CMD_DOWN_KICK_OUT, 0, []byte("kick out"))
		dp := pack.NewDataPackDown()
		if downData, err := dp.Pack(downMsg); err != nil {
			logger.Errorf(request, `[login handler] connection.SendByteMsg pack. error:%v`, err)
		} else {
			if err := connection.SendByteMsg(downData); err != nil {
				logger.Errorf(request, `[login handler] connection.SendByteMsg. error: %v`, err)
				return err
			}
		}
		connection.GetTCPServer().GetConnMgr().Remove(connection)
		connection.SetKickOut(1)
		connection.Stop()
	}

	//把conn添加到players
	if err := conn.GetTCPServer().GetConnMgr().AddConnByUserId(conn); err != nil {
		logger.Errorf(request, `[login handler] GetConnMgr.AddConnByUserId. error: %v`, err)
		return err
	}

	//结束时间(毫秒)
	end := time.Since(start).Milliseconds()
	logger.Infof(request, `[Status Code:%d | MsgID:%d | Execution Time:%d(ms)]`, 200, request.GetMsgID(), end)

	//return
	return PushMessage(request, global.CMD_DOWN_LOGIN, 0, []byte("login successful!"))
}
