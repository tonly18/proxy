package iface

type PacketDown interface {
	GetHeadLen() uint32 //获取包头长度方法

	Pack(IMessageDown) ([]byte, error)
	UnPack([]byte) (IMessageDown, error)
}
