package ziface

/*
将打包封装到packet中，定义抽象层接口
*/
type Packet interface {
	GetHeadLen() uint32
	Pack(IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
