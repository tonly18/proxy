package iface

type PacketPing interface {
	GetHeadLen() uint32 //获取包头长度方法

	Pack(IMessagePing) ([]byte, error)
	UnPack([]byte) (IMessagePing, error)
}
