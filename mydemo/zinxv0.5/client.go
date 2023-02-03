package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
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
		//封包消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgDataPackage(0, []byte("zinx v0.5")))
		if err != nil {
			fmt.Println("pack err", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write err", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {

			fmt.Println("read head err", err)
			break
		}
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack err", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Println("msg id", msg.Id, " len=", msg.DataLen, " data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
