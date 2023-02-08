package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//负责worker池的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 执行对应的router
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//从request找到msgid
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api! msgID", request.GetMsgId(), "is NOT FOUND need register")
	}
	//根据msgid调度找到对应的router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1. 判断当前msg绑定的API是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated" + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID", msgID, "success")
}

// 启动一个工作池 （只能发生一次）
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//1个worker启动

		//当前worker对应的channel消息队列 开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker 阻塞等待消息
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个工作流
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID ", workerID, "is started")
	//不断的阻塞等待对应的消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)

		}
	}
}

// 将消息交给taskqueue , 由 worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker  分布式
	workerID := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	fmt.Println("Add conID = ", request.GetConnection().GetConnId(),
		"request msgID = ", request.GetMsgId(),
		"to workerID = ", workerID)

	mh.TaskQueue[workerID] <- request
}
