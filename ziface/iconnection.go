package ziface

import "net"

type IConnection interface {
	//启动连接 让当前连接准备开始工作
	Start()
	//停止连接 结束当前连接的工作
	Stop()
	//获取当前连接的 socket conn
	GetTcpConnection() *net.TCPConn
	//获取当前连接的id
	GetConnId() uint32
	//获取远程客户端的 tcp 状态 ip port
	GetRemoterAddr() net.Addr
	//发送数据
	// send(data []byte) error
}

// 定义一个抽象的函数类型,用于处理业务
type HandleFunc func(*net.TCPConn, []byte, int) error
