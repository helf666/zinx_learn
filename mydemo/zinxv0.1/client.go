package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start..")

	time.Sleep(1 * time.Second)
	//1 直接连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,", err)
		return
	}
	for {
		//2 连接调用write方法写数据
		conn.Write([]byte("hello zinxv0.1"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err", err)
			return

		}
		fmt.Printf("server call back %s,cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
