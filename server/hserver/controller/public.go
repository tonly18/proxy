package controller

import (
	"github.com/spf13/cast"
	"io"
	"proxy/core/server"
	"proxy/library/logger"
	"proxy/server/global"
	"strings"
)

// PublicController gameServer向proxy发送消息
func PublicController(req *server.Request) *server.Response {
	//需要推送消息的玩家ID
	toMsgPlayerId := req.GetData("player_id")
	toMsgPlayerIds := strings.Split(toMsgPlayerId.(string), ",")
	playerIds := make([]int, 0, len(toMsgPlayerIds)/2)
	for _, v := range toMsgPlayerIds {
		if v != "" {
			playerIds = append(playerIds, cast.ToInt(v))
		}
	}
	if len(playerIds) == 0 {
		return &server.Response{
			Code: 100000,
			Type: 2,
		}
	}

	//data
	data, err := io.ReadAll(req.Request.Body)
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
		conn, err := global.GetTCPServer().GetConnMgr().GetByUserId(uint64(uid))
		if err != nil { //玩家不存在
			logger.Errorf(req, `[http server public] player is not exist. user id:%v, error:%v`, uid, err)
			continue
		}
		if err := conn.SendBuffMsg(data); err != nil {
			logger.Errorf(req, `[http server public] player.SendBuffMsg. user id:%v, error:%v`, uid, err)
			continue
		}
	}

	//return
	return &server.Response{
		Code: 0,
	}
}
