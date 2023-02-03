package main

import (
	"fmt"
	"io"
	"net"
	"zinx/znet"
)

func main() {
	//1.创建socket
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err", err)
		return
	}
	//2.从客户端读取数据
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept err", err)
				return
			}
			go func(conn net.Conn) {
				dp := znet.NewDataPack()
				for {
					//第一次读吧包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的
						msg := msgHead.(*znet.Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据len长度再从io流进行读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack", err)
							return
						}

						fmt.Println("server recv", msg.Id, "datalen:", msg.DataLen, string(msg.Data))

					}
				}

			}(conn)
		}
	}()

	fmt.Println("run client")
	//
	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}
	dp := znet.NewDataPack()
	//封装第一个msg
	msg1 := &znet.Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
		return
	}
	//封装第二个msg
	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println(err)
		return
	}
	//两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}
