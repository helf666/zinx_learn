package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("PingRouter PreHandle")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("Hello1\n"))
	if err != nil {
		fmt.Println("call back before ping err")
	}
}
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("Hello2\n"))
	if err != nil {
		fmt.Println("call back in ping err")
	}
}
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("PingRouter PostHandle")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("Hello3\n"))
	if err != nil {
		fmt.Println("call back after ping err")
	}
}
func main() {
	s := znet.NewServer("zinx v0.3")

	//添加router
	s.AddRouter(&PingRouter{})

	s.Server()
}
