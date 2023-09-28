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
	"proxy/server/utils"
	"strings"
	"time"
)

//PublicRouter Struct
type PublicRouter struct {
	BaseHandler
}

func (h *PublicRouter) Handle(request ziface.IRequest) error {
	//开始时间
	start := time.Now()

	//判断玩家是否登录
	userid := request.GetConnection().GetUserId() //当前玩家ID
	if userid == 0 {
		return errors.New("[Public Handle] player not login")
	}

	//获取玩家
	conn, err := global.GetTCPServer().GetConnMgr().GetConnByUserId(userid)
	if err != nil {
		return fmt.Errorf(`[Public Handle] player not exist. error: %w`, err)
	}

	//http client
	client := httpclient.NewHttpClient(&httpclient.Config{})
	resp, err := client.NewRequest("POST", config.HttpConfig.GameServerCommandAPI, request.GetData()).SetHeader(map[string]any{
		"Content-Type": "application/octet-stream",
		"proxy_id":     conn.GetTCPServer().GetID(),
		"server_id":    conn.GetProperty("server_id"),
		"user_id":      userid,
		"trace_id":     request.GetTraceId(),
		"client_ip":    conn.GetRemoteIP(),
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

	/** 调用PHP GameServer **/
	if downRawData == nil || len(downRawData) < 16 {
		resp, err = client.NewRequest("POST", config.HttpConfig.GameServerPHPCommandAPI, request.GetData()).SetHeader(map[string]any{
			"Content-Type": "application/octet-stream",
			"proxy_id":     conn.GetTCPServer().GetID(),
			"server_id":    conn.GetProperty("server_id"),
			"user_id":      userid,
			"trace_id":     request.GetTraceId(),
			"client_ip":    conn.GetRemoteIP(),
			"socket_id":    conn.GetConnID(),
		}).Do()
		if err != nil {
			return fmt.Errorf(`[Public Handle] call http(php) server is error: %v`, err)
		}
		if resp.GetHeaderCode() != 200 {
			return fmt.Errorf(`[Public Handle] call http(php) server is code:%v`, resp.Response.StatusCode)
		}
		//处理返回值
		downRawData, err = resp.GetData()
		if err != nil {
			return fmt.Errorf(`[Public Handle] resp.GetData(php) error:%v`, err)
		}
	}

	if downRawData != nil && len(downRawData) < 16 {
		return fmt.Errorf(`[Public Handle] resp.GetData rawDownData is empty. downRawData: %v`, downRawData)
	}

	//json
	resultJson := utils.ResultJSON{}
	if err := json.Unmarshal(downRawData, &resultJson); err == nil {
		return fmt.Errorf(`[Public Handle] resp.GetData rawDownData: %v`, resultJson)
	}

	//给客户端(玩家)推送消息
	playerId := resp.GetDataFromHeader("player_id")
	playerIds := make([]uint64, 0, 20)
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
