package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	//服务器的名称
	Name string
	//ip版本号
	IPVersion string
	//监听的ip
	IP string
	//监听的端口
	Port int
	//当前server的消息管理模块，用来绑定msgid和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[zinx] server name is %s ,IP is %s, port is %d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx] version is %s,max packagesize is %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxPackageSize)
	//非阻塞
	go func() {
		//1获取一个tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve TCP error:", err)
			return
		}
		//2监听一个服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s, s.IPVersion, "err", err)
			return
		}
		fmt.Println("start zinx server succ", s.Name, "success,listening")
		var cid uint32 = 0

		//3阻塞等待客户端进行连接，处理客户端连接的业务
		for {
			//如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err,", err)
				return
			}
			//处理新链接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Server() {
	s.Start()

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router success")
}

func NewServer(name string) ziface.IServer { //调用接口
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,    //"127.0.0.1",
		Port:       utils.GlobalObject.TcpPort, //8999,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
