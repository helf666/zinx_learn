package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	//告知当前的链接是否已经退出
	ExitChan chan bool //由reader告知 writer退出的信号

	//无缓冲d管道，用于读写goroutin
	msgChan chan []byte

	//消息管理MsgID
	MsgHandler ziface.IMsgHandler
}

// 初始化连接的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     ConnID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("ConnID=", c.ConnID, "reader is exit,remote addr is", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512字节

		dp := NewDataPack()
		//读取客户端的Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg err", err)
			break
		}

		//拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error")
			break
		}

		//根据datalen再次读取data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
			}
		}
		msg.SetData(data)

		//将Msg交给路由器
		req := Request{
			conn: c,
			msg:  msg,
		}

		//根据绑定的conn对应router调用
		//根据绑定好的ID找到对应的API
		go c.MsgHandler.DoMsgHandler(&req)

	}
}

// 提供一个sendmsg方法 将要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	//将data进行封包  msgdatalen msgid  data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgDataPackage(msgId, data)) //序列化好的msg
	if err != nil {
		fmt.Println("package err", err)
	}
	//将数据发送到客户端
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn writer exit")

	//不断的阻塞等待channel的消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err", err, "conn writer exit")
				return
			}

		case <-c.ExitChan:
			//Reader退出 writer也要退出
			return
		}

	}
}

func (c *Connection) Start() {
	fmt.Println("conn start connid = ", c.ConnID)
	//启动当前连接读数据的业务\
	go c.StartReader()
	//启动当前连接写数据的业务
	go c.StartWriter()
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
	//告知 writer关闭
	c.ExitChan <- true
	//关闭管道
	close(c.ExitChan)
	close(c.msgChan)
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
