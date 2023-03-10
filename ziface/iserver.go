package ziface

//定义一个服务器接口

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Server()
	//注册路由
	AddRouter(msgID uint32, router IRouter)
}
