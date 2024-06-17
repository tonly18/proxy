package znet

import (
	"proxy/core/zinx/ziface"
	"proxy/library/command"
	"sync"
)

var requestPool = &sync.Pool{
	New: func() any {
		return allocateRequest()
	},
}

func allocateRequest() ziface.IRequest {
	return NewRequest(nil, nil)
}

func GetRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := requestPool.Get().(*Request)
	req.conn = conn
	req.msg = msg
	if req.router != nil {
		req.router = nil
	}
	if req.steps != PRE_HANDLE {
		req.steps = PRE_HANDLE
	}
	req.traceId = command.GenTraceID()

	return req
}

func PutRequest(request ziface.IRequest) {
	requestPool.Put(request)
}
