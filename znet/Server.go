package znet

import (
	"fmt"
	"net"
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
		//3阻塞等待客户端进行连接，处理客户端连接的业务
		for {
			//如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err,", err)
				return
			}
			//已经有客户端连接，做业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err,", err)
					}
					fmt.Printf("server recv %s,cnt = %d\n", buf, cnt)

					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()

	//阻塞状态
	select {}
}

func NewServer(name string) *Server { //ziface.IServer
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8999,
	}
	return s
}
