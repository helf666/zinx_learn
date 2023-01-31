package znet

import "zinx/ziface"

type BaseRouter struct {
}

// 抽象层都实现，实际中自己可能实现要实现的东西
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (br *BaseRouter) Handle(request ziface.IRequest) {

}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
