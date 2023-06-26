package iface

type PacketKickOut interface {
	GetHeadLen() uint32 //获取包头长度方法

	Pack(IMessageKickOut) ([]byte, error)
	UnPack([]byte) (IMessageKickOut, error)
}
