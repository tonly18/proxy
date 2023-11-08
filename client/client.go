package main

import (
	"fmt"
	"io"
	"net"
	"proxy/server/global"
	"proxy/server/utils/pack"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		panic(fmt.Sprintf("net dial error: %v", err))
	}

	testServer(conn)
	//loginServer(conn)

	time.Sleep(3 * time.Second)
}

// 返回值处理
func ReadFromConn(conn net.Conn) {
	for {
		//读取包头
		pb := pack.NewDataPackDown()
		headerBuffer := make([]byte, pb.GetHeadLen())
		if _, err := io.ReadFull(conn, headerBuffer); err != nil {
			fmt.Println("io read full error: ", err)
			break
		}
		downHeader, err := pb.UnPack(headerBuffer)
		if err != nil {
			fmt.Println("unpack error----------: ", err)
			break
		}
		fmt.Println("login------downMsg.head.len:::::::::::::::", downHeader.GetMsgLen())
		fmt.Println("login------downMsg.head.cmd:::::::::::::::", downHeader.GetCmd())
		fmt.Println("login------downMsg.head.code::::::::::::::", downHeader.GetCode())

		//读取包体
		downData := make([]byte, downHeader.GetMsgLen()-pb.GetHeadLen())
		if _, err := io.ReadFull(conn, downData); err != nil {
			fmt.Println("io read full error: ", err)
			break
		}
		fmt.Println("len(downData)-------", len(downData))
		fmt.Println("string(downData)----", string(downData))

		fmt.Println("------------------------------------------")
		break
	}
}

// test
func testServer(conn net.Conn) {
	go func() {
		ReadFromConn(conn)
	}()

	upMsg := pack.NewMessageUp(0, []byte("1234567890"))
	fmt.Println("upMsg.MsgLen:::::::::::::::::", upMsg.GetMsgLen())
	fmt.Println("upMsg.cmd::::::::::::::::::::", upMsg.GetCmd())
	fmt.Println("upMsg.data:::::::::::::::::::", string(upMsg.GetData()))
	fmt.Println("---------------------------------------")

	dp := pack.NewDataPackUp()
	data, err := dp.Pack(upMsg)
	fmt.Println("dp.pack-err:::::::::::::::::", err)
	fmt.Println("dp.pack-len(data):::::::::::", len(data))
	fmt.Println("---------------------------------------")

	//send message
	n, err := conn.Write(data)
	fmt.Println("conn.write-err(testServer)::::::::", n, err)
	fmt.Println("--------------------------------------- \n")

	time.Sleep(5 * time.Second)
	conn.Close()
}

// login
func loginServer(conn net.Conn) {
	go func() {
		ReadFromConn(conn)
	}()

	upMsg := pack.NewMessageUp(global.CMD_UP_LOGIN, []byte("do-login"))
	fmt.Println("upMsg.MsgLen:::::::::::::::::", upMsg.GetMsgLen())
	fmt.Println("upMsg.cmd::::::::::::::::::::", upMsg.GetCmd())
	fmt.Println("upMsg.data:::::::::::::::::::", string(upMsg.GetData()))
	fmt.Println("---------------------------------------")

	dp := pack.NewDataPackUp()
	data, err := dp.Pack(upMsg)
	fmt.Println("dp.pack-err:::::::::::::::::", err)
	fmt.Println("dp.pack-len(data):::::::::::", len(data))
	fmt.Println("---------------------------------------")

	//send message
	n, err := conn.Write(data)
	fmt.Println("conn.write-err(testServer)::::::::", n, err)
	fmt.Println("--------------------------------------- \n")

	time.Sleep(5 * time.Second)
	conn.Close()
}
