package znet

import (
	"fmt"
	"net"
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
	//当前server添加一个router,server的注册的链接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[strat]Server Listen at IP:%s,Port:%d、n", s.IP, s.Port)

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
			dealConn := NewConnection(conn, cid, s.Router)
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

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router success")
}

func NewServer(name string) ziface.IServer { //调用接口
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8999,
		Router:    nil,
	}
	return s
}