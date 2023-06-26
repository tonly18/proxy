package handler

import (
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/znet"
)

//BaseHandler Struct
type BaseHandler struct {
	znet.BaseRouter
}

func (bh *BaseHandler) PreHandle(request ziface.IRequest) error {
	return nil
}

func (bh *BaseHandler) Handle(request ziface.IRequest) error {
	return nil
}

func (bh *BaseHandler) PostHandle(request ziface.IRequest) error {
	return nil
}
