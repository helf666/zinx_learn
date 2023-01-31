package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	//告知当前的链接是否已经退出
	ExitChan chan bool //通道
	//改连接处理的router
	Router ziface.IRouter
}

// 初始化连接的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   ConnID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("ConnID=", c.ConnID, "reader is exit,remote addr is", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err", err)
			continue
		}
		//得到当前的request数据
		req := Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running")
}

func (c *Connection) Start() {
	fmt.Println("conn start connid = ", c.ConnID)
	//启动当前连接读数据的业务\
	go c.StartReader()
	//启动当前连接写数据的业务
}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Stop,connid= ", c.ConnID)

	//如果当前连接已经被关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	//关闭socket
	c.Conn.Close()
	//关闭管道
	close(c.ExitChan)
}

// 获取当前连接的 socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接的id
func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

// 获取远程客户端的 tcp 状态 ip port
func (c *Connection) GetRemoterAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
// func (c *Connection) send(data []byte) error {
// 	return nil
// }
