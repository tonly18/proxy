package znet

import (
	"proxy/core/zinx/ziface"
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

	return req
}

func PutRequest(request ziface.IRequest) {
	request.Reset()
	requestPool.Put(request)
}
