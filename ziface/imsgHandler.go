package ziface

type IMsgHandler interface {
	//执行对应的router
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
