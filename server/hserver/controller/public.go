package controller

import (
	"io/ioutil"
	"proxy/core/server"
	"proxy/library/logger"
	"proxy/server/global"
)

//PublicController gameServer向proxy发送消息
func PublicController(req *server.Request) *server.Response {
	//需要推送消息的玩家ID
	playerIds := req.GetPlayerID()
	if len(playerIds) == 0 {
		return &server.Response{
			Code: 100000,
			Type: 2,
		}
	}

	//data
	data, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return &server.Response{
			Code: 100005,
			Msg:  err,
			Type: 2,
		}
	}
	defer req.Request.Body.Close()
	if data != nil && len(data) == 0 {
		return &server.Response{
			Code: 100010,
			Type: 2,
		}
	}

	//向客户端发送消息
	for _, uid := range playerIds {
		conn, err := global.GetTCPServer().GetConnMgr().GetConnByUserId(uint64(uid))
		if err != nil { //玩家不存在
			logger.Errorf(req, `[http server public] player is not exist. user id:%v, error:%v`, uid, err)
			continue
		}
		if err := conn.SendByteMsg(data); err != nil {
			logger.Errorf(req, `[http server public] player.SendByteMsg. user id:%v, error:%v`, uid, err)
			continue
		}
	}

	//return
	return &server.Response{
		Code: 0,
	}
}
