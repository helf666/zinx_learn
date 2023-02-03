package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("PingRouter Handle")

	//读取客户端的数据，然后再写业务
	fmt.Println("recv from client : msgid", request.GetMsgId(), " data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping"))
	if err != nil {
		fmt.Println(err)
	}

}
func main() {
	s := znet.NewServer("zinx v0.3")

	//添加router
	s.AddRouter(&PingRouter{})

	s.Server()
}
