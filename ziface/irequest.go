package ziface

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	//得到请求
	GetData() []byte
	//得到消息的Id
	GetMsgId() uint32
}
