package handler

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/tonly18/httpclient"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/config"
	"proxy/server/global"
	"proxy/server/library"
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
	//if msg.GetUin() == nil {
	//	request.GetConnection().Stop()
	//	return errors.New("[login handler] msg uin is nil")
	//}
	//if msg.GetServer() == nil {
	//	request.GetConnection().Stop()
	//	return errors.New("[login handler] msg server id is nil")
	//}

	//调用gameServer验证登录
	url := fmt.Sprintf(`%v%v`, library.NewPoll().Get(), config.HttpConfig.GameServerDoLoginAPI)
	client := httpclient.NewHttpClient(&httpclient.Config{})
	resp, err := client.NewRequest("POST", url, request.GetData()).SetHeader(map[string]any{
		"Content-Type": "application/octet-stream",
		"proxy_id":     request.GetConnection().GetTCPServer().GetID(),
		"server_id":    0,
		"user_id":      0,
		"client_ip":    request.GetConnection().GetRemoteIP(),
		"trace_id":     request.GetTraceId(),
	}).Do()
	if err != nil {
		return fmt.Errorf(`login call http server is error: %v`, err)
	}
	request.SetAargs("gameserver_id", resp.GetDataFromHeader("gameserver_id")) //game server id
	if resp.GetHeaderCode() != 200 {
		return fmt.Errorf(`login call http server is code:%v`, resp.Response.StatusCode)
	}

	//处理返回数据
	uin := cast.ToUint64(resp.GetDataFromHeader("uin"))            //帐号ID(一个帐号ID对应多个玩家ID)
	userId := cast.ToUint64(resp.GetDataFromHeader("user_id"))     //玩家ID
	serverId := cast.ToUint32(resp.GetDataFromHeader("server_id")) //区服ID
	if uin == 0 || userId == 0 || serverId == 0 {
		logger.Error(request, "[login handler] login fail.")
		return err
	}

	//设置conn属性
	conn := request.GetConnection()
	conn.SetProxyId(conn.GetTCPServer().GetID()) //网关ID
	conn.SetServerId(serverId)                   //区服ID
	conn.SetUIN(uin)                             //帐号ID
	conn.SetUserId(userId)                       //玩家ID

	//添加登录后的玩家到connManager
	if connection, _ := conn.GetTCPServer().GetConnMgr().GetConnByUserId(userId); connection != nil {
		//踢掉原connection,并推送消息给客户端
		downMsg := pack.NewMessageKickOut(global.CMD_DOWN_KICK_OUT, 5)
		pb := pack.NewDataPackKickOut()
		if downData, err := pb.Pack(downMsg); err != nil {
			logger.Errorf(request, `[login handler] connection.SendByteMsg pack. error:%v`, err)
		} else {
			if err := connection.SendByteMsg(downData); err != nil {
				logger.Errorf(request, `[login handler] connection.SendByteMsg. error: %v`, err)
				return err
			}
		}
		connection.GetTCPServer().GetConnMgr().Remove(connection)
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
	return nil
}
