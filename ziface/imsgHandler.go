package ziface

type IMsgHandler interface {
	//执行对应的router
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	//将消息发送到消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
	//开启工作池
	StartWorkerPool()
	//开启一个工作
	StartOneWorker(workerID int, taskQueue chan IRequest)
}
