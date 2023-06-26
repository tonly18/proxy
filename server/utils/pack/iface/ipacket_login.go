package iface

type PacketLogin interface {
	GetHeadLen() uint32 //获取包头长度方法

	Pack(IMessageLogin) ([]byte, error)
	UnPack([]byte) (IMessageLogin, error)
}
