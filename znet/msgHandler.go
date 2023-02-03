package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
