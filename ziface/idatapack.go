package ziface

type IDataPack interface {
	//获取头长度的方法
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
