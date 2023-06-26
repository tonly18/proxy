package iface

type PacketUp interface {
	GetHeadLen() uint32 //获取包头长度方法

	Pack(IMessageUp) ([]byte, error)
	UnPack([]byte) (IMessageUp, error)
}
