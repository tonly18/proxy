package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tonly18/httpclient"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/config"
	"proxy/server/global"
	"proxy/server/library"
	"proxy/server/utils"
	"strings"
	"time"
)

// PublicRouter Struct
type PublicRouter struct {
	BaseHandler
}

func (t *PublicRouter) Handle(request ziface.IRequest) error {
	//开始时间
	start := time.Now()

	//判断玩家是否登录
	userId := request.GetConnection().GetUserId() //当前玩家ID
	if userId == 0 {
		return errors.New("[Public Handle] player not login")
	}

	//获取玩家connection
	conn, err := global.GetTCPServer().GetConnMgr().GetConnByUserId(userId)
	if err != nil {
		return fmt.Errorf(`[Public Handle] player not exist. error: %w`, err)
	}

	//http client
	url := fmt.Sprintf(`%v%v`, library.NewPoll().Get(), config.HttpConfig.GameServerCommandAPI)
	client := httpclient.NewHttpClient(&httpclient.Config{})
	resp, err := client.NewRequest("POST", url, request.GetData()).SetHeader(map[string]any{
		"Content-Type": "application/octet-stream",
		"proxy_id":     conn.GetTCPServer().GetID(),
		"server_id":    conn.GetProperty("server_id"),
		"user_id":      userId,
		"client_ip":    conn.GetRemoteIP(),
		"trace_id":     request.GetTraceId(),
	}).Do()
	if err != nil {
		return fmt.Errorf(`[Public Handle] call http server is error: %w`, err)
	}
	if resp.GetHeaderCode() != 200 {
		return fmt.Errorf(`[Public Handle] call http server is code: %v`, resp.Response.StatusCode)
	}

	//处理返回值
	downRawData, err := resp.GetData()
	if err != nil {
		return fmt.Errorf(`[Public Handle] resp.GetData error: %w`, err)
	}
	request.SetAargs("gameserver_id", resp.GetDataFromHeader("gameserver_id")) //game server id
	if downRawData != nil && len(downRawData) < 16 {
		return fmt.Errorf(`[Public Handle] resp.GetData rawDownData is empty. downRawData: %v`, downRawData)
	}

	//json
	if downRawData != nil {
		resultJson := utils.ResultJSON{}
		if err := json.Unmarshal(downRawData, &resultJson); err == nil {
			return fmt.Errorf(`[Public Handle] resp.GetData rawDownData: %v`, resultJson)
		}
	}

	//给客户端(玩家)推送消息
	playerId := resp.GetDataFromHeader("player_id") //玩家ID
	playerIds := make([]uint64, 0, 10)
	if playerId != "" {
		for _, v := range strings.Split(playerId, ",") {
			playerIds = append(playerIds, cast.ToUint64(v))
		}
	}
	for _, v := range playerIds {
		connection, err := global.GetTCPServer().GetConnMgr().GetConnByUserId(v)
		if err != nil { //玩家不存在
			logger.Errorf(request, `[Public Handle] player not exist. userId:%v, length:%d, error:%v`, v, len(downRawData), err)
			continue
		}
		if err := connection.SendByteMsg(downRawData); err != nil {
			logger.Errorf(request, `[Public Handle] connection.SendByteMsg. userId:%v, length:%d, error:%v`, v, len(downRawData), err)
			continue
		}
	}

	//结束时间(毫秒)
	end := time.Since(start).Milliseconds()
	logger.Infof(request, `[Status Code:%d | MsgID:%d | Length:%d | Execution Time:%d(ms)]`, resp.GetHeaderCode(), request.GetMsgID(), len(downRawData), end)

	//return
	return nil
}
