package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"proxy/proto/pb/down"
	"proxy/proto/pb/up"
	"proxy/server/utils/pack"
	"time"
)

//cpack "proxy/client/library/pack"

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		fmt.Println("net dial error: ", err)
	}

	//go doLoginProxy(conn)
	//time.Sleep(1 * time.Second)

	//go getGuildList(conn)
	//time.Sleep(1 * time.Second)

	//loginGameServer(conn)
	testServer(conn)
	time.Sleep(1 * time.Second)

	//ping
	go func() {
		//	for true {
		//		pingProxy(conn)
		//		time.Sleep(3 * time.Second)
		//	}
	}()

	//go setAvatar(conn)
	//time.Sleep(1 * time.Second)

	//发送
	go func() {
		/*
			for {
				//发送消息
				//var repeat uint32 = 1
				//var serverId int32 = 101
				//var deviceId string = "deviceId"
				//var seq int32 = int32(time.Now().Second())
				//
				//upmsg := up.UpMsg{
				//	XRepeat: &repeat,
				//	XLogin: &up.Login{
				//		XServerId: &serverId,
				//		XDeviceId: &deviceId,
				//		XSeq:      &seq,
				//	},
				//}
				//dataBuffer, err := proto.Marshal(&upmsg)
				//if err != nil {
				//	fmt.Println("proto marshal error: ", err)
				//	continue
				//}
				//
				////pack
				//pb := znet.NewMessagePack()
				//msgPackage, err := pb.Pack(znet.NewMessage(1000, dataBuffer))
				//if err != nil {
				//	fmt.Println("db.pack error: ", err)
				//	break
				//}
				//if _, err := conn.Write(msgPackage); err != nil {
				//	fmt.Println("conn write error: ", err)
				//	break
				//}

				//time.Sleep(2 * time.Second)

				//msg2 := znet.NewMsgPackage(201, []byte("this is a test message from the client-201!!!"))
				//pb2 := znet.NewDataPack()
				//msgBuf2, err := pb2.Pack(msg2)
				//if err != nil {
				//	fmt.Println("db.pack error: ", err)
				//	break
				//}
				//if _, err := conn.Write(msgBuf2); err != nil {
				//	fmt.Println("conn write error: ", err)
				//	break
				//}
				//
				////time.Sleep(2 * time.Second)
				//
				//msgBuf = append(msgBuf, msgBuf2...)
				//if _, err := conn.Write(msgBuf); err != nil {
				//	fmt.Println("conn write error: ", err)
				//	break
				//}

				//收到的消息
				//buf := make([]byte, 512)
				//if _, err := conn.Read(buf); err != nil {
				//	fmt.Println("conn read error: ", err)
				//	continue
				//}
				//fmt.Println("server send message: ", string(buf))

				//sleep
				time.Sleep(1 * time.Second)

				//close
				//conn.Close()
				//break
			}
		*/
	}()

	//收取
	go func() {
		/*
			for {
				dp := pack.NewDataPackDown()

				//读取包头
				headBuffer := make([]byte, dp.GetHeadLen())
				if _, err := io.ReadFull(conn, headBuffer); err != nil {
					fmt.Println("io read full error: ", err)
					break
				}
				msg, err := dp.UnPack(headBuffer)
				if err != nil {
					fmt.Println("unpack error+++++001: ", err)
					break
				}
				msgPb := msg.(*pack.MessageDown)
				fmt.Println("len(headBuffer):::", string(headBuffer))
				fmt.Println("msg.head.len:::::::::::::::", msg.GetMsgLen())
				fmt.Println("msg.head.cmd:::::::::::::::", msg.GetCmd())
				fmt.Println("msg.head.errorcode:::::::::", msg.GetErrorCode())
				fmt.Println("msg.head.pblen:::::::::::::", msg.GetPbLen())

				//读取包体
				msgPb.Pb = make([]byte, msg.GetPbLen())
				if _, err := io.ReadFull(conn, msgPb.Pb); err != nil {
					fmt.Println("io read full error: ", err)
					break
				}

				//合包
				//msg, err = dp.UnPack(append(headBuffer, msgPb.Pb...))

				//解析pb
				//fmt.Println("msg.id++++++++++++: ", msg.GetCmd())
				//fmt.Println("++++++++++++msg.data: ", string(msg.GetData()))

				donwMsg := down.DownMsg{}
				proto.Unmarshal(msg.GetPb(), &donwMsg)
				//fmt.Println("donwMsg::::::", donwMsg)

				guildReply := donwMsg.GetXGuildReply()
				fmt.Println("guildReply.result::::::", guildReply.GetXList().GetXResult())
				fmt.Println("guildReply.serverId::::", guildReply.GetXList().GetXServerId())
				fmt.Println("guildReply.guilds::::::", guildReply.GetXList().GetXGuilds())

				//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatar())
				//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatar())
				//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXItemId())
				//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXGold())
				//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatarUrl())

				fmt.Println("\n")
			}
		*/
	}()

	//阻塞
	//select {}
	//time.Sleep(30 * time.Second)
}

func doLoginProxy(conn net.Conn) {
	//loginMsg := pack.NewMessageLogin(1, []byte("229"), []byte("1"), []byte("session_key-202305"), []byte("version-202305"))
	loginMsg := pack.NewMessageLogin(1, []byte("867"), []byte("1"), []byte("78100f10719bd2a7"), []byte("version-202305"))
	dp := pack.NewDataPackLogin()
	data, _ := dp.Pack(loginMsg)

	fmt.Println("loginMsg.MsgLen: ", loginMsg.GetMsgLen())
	fmt.Println("loginMsg.Cmd: ", loginMsg.GetCmd())
	fmt.Println("loginMsg.Uin: ", loginMsg.GetUin())
	fmt.Println("loginMsg.server: ", loginMsg.GetServer())
	fmt.Println("loginMsg.seesionKey: ", string(loginMsg.GetSessionKey()))
	fmt.Println("loginMsg.version: ", string(loginMsg.GetVersion()))

	length, err := conn.Write(data)
	if err != nil {
		fmt.Println("conn write error: ", err)
	}
	fmt.Println("length::::::", length)
	fmt.Println("------------------do login------------------")

	//sleep
	//time.Sleep(1 * time.Second)
	//conn.Close()
}

func pingProxy(conn net.Conn) {
	pingMsg := pack.NewMessagePing(34)
	dp := pack.NewDataPackPing()
	data, _ := dp.Pack(pingMsg)

	fmt.Println("pingMsg.MsgLen: ", pingMsg.GetMsgLen())
	fmt.Println("pingMsg.Cmd: ", pingMsg.GetCmd())

	if _, err := conn.Write(data); err != nil {
		fmt.Println("conn write error: ", err)
	}

	//close
	//conn.Close()
	//break

	//sleep
	//time.Sleep(10 * time.Second)
}

func setAvatar(conn net.Conn) {
	fmt.Println("--------------------setAvatar--------------------")
	//var result uint32 = 0
	//var avatar int32 = 111
	//var itemId uint32 = 222
	//var gold int64 = 333
	//upMsg := &up.UpMsg{
	//	XRepeat: &result,
	//	XSetAvatar: &up.SetAvatar{
	//		XAvatar: &avatar,
	//		XItemId: &itemId,
	//		XGold:   &gold,
	//	},
	//}
	//upDataPb, _ := proto.Marshal(upMsg)
	//
	//msg := cpack.NewMessageUp(35, upDataPb, []byte("uin-01"), []byte("key-01"), []uint64{100, 101}, []uint32{200, 201})
	//dp := cpack.NewDataPackUp()
	//data, err := dp.Pack(msg)
	//fmt.Println("data-err:", err)
	//
	//n, err := conn.Write(data)
	//fmt.Println("n-err:", n, err)
}

//公会 - 列表
func getGuildList(conn net.Conn) {
	doLoginProxy(conn)

	//userId := uint32(867) //当前游戏玩家ID
	var repeat uint32 = 0
	var serverId int32 = 1
	var aid int32 = 1
	upMsg := &up.UpMsg{
		XRepeat: &repeat,
		XGuild: &up.Guild{
			XList: &up.GuildList{
				XAreaId:   &aid,
				XServerId: &serverId,
			},
		},
	}

	//pack
	upPb, _ := proto.Marshal(upMsg)
	//fmt.Println("pb.len---len(upPb)::::::::::::::::::", len(upPb))
	//fmt.Println("---------------------------------------")

	msg := pack.NewMessageUp(35, upPb, []byte("uin-001"), []byte("key-202305"), []uint64{100, 101}, []uint32{200, 201})
	//fmt.Println("msg.MsgLen:::::::::::::::::", msg.GetMsgLen())
	//fmt.Println("msg.len(pb)::::::::::::::::", len(msg.GetPb()))
	//fmt.Println("---------------------------------------")

	dp := pack.NewDataPackUp()
	data, err := dp.Pack(msg)
	fmt.Println("dp.pack-err:::::::::::::::::", err)
	fmt.Println("dp.pack-len(data):::::::::::", len(data))
	fmt.Println("---------------------------------------")

	//unpack
	//res, err := dp.UnPack(data)
	//fmt.Println("dp.UnPack.err:::::::::::::::", err)
	//fmt.Println("res.datalen:::::::::::::::::", res.GetMsgLen())
	//fmt.Println("res.len(res.GetPb())::::::::", len(res.GetPb()))
	//fmt.Println("res.troop::::::::", res.GetTroop())
	//fmt.Println("---------------------------------------")

	//upMessage := up.UpMsg{}
	//err = proto.Unmarshal(res.GetPb(), &upMessage)
	//fmt.Println("unmarshal.err:::::::::::::::", err)
	//fmt.Println("upMessage.repeat::::::::::::", upMessage.GetXRepeat())
	//fmt.Println("upMessage.guild:::::::::::::", upMessage.GetXGuild())
	//fmt.Println("upMessage.areaId::::::::::::", upMessage.GetXGuild().GetXList().GetXAreaId())
	//fmt.Println("upMessage.serverId::::::::::", upMessage.GetXGuild().GetXList().GetXServerId())

	//send message
	n, err := conn.Write(data)
	fmt.Println("conn.write-err::::::::", n, err)
	time.Sleep(1 * time.Second)
	fmt.Println("--------------------------------------- \n")

	//处理返回数据
	go func() {
		for {
			dp := pack.NewDataPackDown()

			//读取包头
			headBuffer := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(conn, headBuffer); err != nil {
				fmt.Println("io read full error: ", err)
				break
			}
			msg, err := dp.UnPack(headBuffer)
			if err != nil {
				fmt.Println("unpack error+++++001: ", err)
				break
			}
			msgPb := msg.(*pack.MessageDown)
			fmt.Println("len(headBuffer):::", string(headBuffer))
			fmt.Println("msg.head.len:::::::::::::::", msg.GetMsgLen())
			fmt.Println("msg.head.cmd:::::::::::::::", msg.GetCmd())
			fmt.Println("msg.head.errorcode:::::::::", msg.GetErrorCode())
			fmt.Println("msg.head.pblen:::::::::::::", msg.GetPbLen())

			//读取包体
			msgPb.Pb = make([]byte, msg.GetPbLen())
			if _, err := io.ReadFull(conn, msgPb.Pb); err != nil {
				fmt.Println("io read full error: ", err)
				break
			}

			//合包
			//msg, err = dp.UnPack(append(headBuffer, msgPb.Pb...))

			//解析pb
			//fmt.Println("msg.id++++++++++++: ", msg.GetCmd())
			//fmt.Println("++++++++++++msg.data: ", string(msg.GetData()))

			donwMsg := down.DownMsg{}
			proto.Unmarshal(msg.GetPb(), &donwMsg)
			//fmt.Println("donwMsg::::::", donwMsg)

			guildReply := donwMsg.GetXGuildReply()
			fmt.Println("guildReply.result::::::", guildReply.GetXList().GetXResult())
			fmt.Println("guildReply.serverId::::", guildReply.GetXList().GetXServerId())
			fmt.Println("guildReply.guilds::::::", guildReply.GetXList().GetXGuilds())

			//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatar())
			//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatar())
			//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXItemId())
			//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXGold())
			//fmt.Println("downmsg++++++++++++: ", downmsg.GetXSetAvatarReply().GetXAvatarUrl())

			fmt.Println("\n")
			break
		}
	}()

	time.Sleep(5 * time.Second)
	conn.Close()
}

//初始始化数据
func loginGameServer(conn net.Conn) {
	doLoginProxy(conn)

	//userId := uint32(867) //当前游戏玩家ID
	var repeat uint32 = 0
	var serverId int32 = 1
	upMsg := &up.UpMsg{
		XRepeat: &repeat,
		XLogin: &up.Login{
			XServerId: proto.Int32(serverId),
			XDeviceId: proto.String("abc-def-92345"),
			XSeq:      proto.Int32(1),
		},
	}

	//pack
	upPb, _ := proto.Marshal(upMsg)
	fmt.Println("pb.len---len(upPb)::::::::::::::::::", len(upPb))
	fmt.Println("---------------------------------------")

	msg := pack.NewMessageUp(35, upPb, []byte("867"), []byte("key-20230531"), []uint64{100, 101}, []uint32{200, 201})
	//fmt.Println("msg.MsgLen:::::::::::::::::", msg.GetMsgLen())
	//fmt.Println("msg.len(pb)::::::::::::::::", len(msg.GetPb()))
	//fmt.Println("---------------------------------------")

	dp := pack.NewDataPackUp()
	data, err := dp.Pack(msg)
	fmt.Println("dp.pack-err:::::::::::::::::", err)
	fmt.Println("dp.pack-len(data):::::::::::", len(data))
	fmt.Println("---------------------------------------")

	//unpack
	//res, err := dp.UnPack(data)
	//fmt.Println("dp.UnPack.err:::::::::::::::", err)
	//fmt.Println("res.datalen:::::::::::::::::", res.GetMsgLen())
	//fmt.Println("res.len(res.GetPb())::::::::", len(res.GetPb()))
	//fmt.Println("res.troop::::::::", res.GetTroop())
	//fmt.Println("---------------------------------------")

	//upMessage := up.UpMsg{}
	//err = proto.Unmarshal(res.GetPb(), &upMessage)
	//login := upMessage.GetXLogin()
	//fmt.Println("upMessage.err:::::::::::::::", err)
	//fmt.Println("upMessage.repeat::::::::::::", upMessage.GetXRepeat())
	//fmt.Println("upMessage.login.serverId::::", *(login.XServerId))
	//fmt.Println("upMessage.login.seq:::::::::", *(login.XSeq))
	//fmt.Println("upMessage.login.deviceId::::", *(login.XDeviceId))

	//send message
	n, err := conn.Write(data)
	fmt.Println("conn.write-err(loginGameServer)::::::::", n, err)
	fmt.Println("--------------------------------------- \n")

	time.Sleep(100 * time.Millisecond)

	//返回值处理
	go func() {
		for {
			pb := pack.NewDataPackDown()
			//读取包头
			headerBuffer := make([]byte, pb.GetHeadLen())
			if _, err := io.ReadFull(conn, headerBuffer); err != nil {
				fmt.Println("io read full error: ", err)
				break
			}
			downMsg, err := pb.UnPack(headerBuffer)
			if err != nil {
				fmt.Println("unpack error----------: ", err)
				break
			}

			fmt.Println("login------downMsg.head.len:::::::::::::::", downMsg.GetMsgLen())
			fmt.Println("login------downMsg.head.cmd:::::::::::::::", downMsg.GetCmd())
			fmt.Println("login------downMsg.head.errorcode:::::::::", downMsg.GetErrorCode())
			fmt.Println("login------downMsg.head.pblen:::::::::::::", downMsg.GetPbLen())

			//读取包体
			downMsgPB := make([]byte, downMsg.GetPbLen())
			if _, err := io.ReadFull(conn, downMsgPB); err != nil {
				fmt.Println("io read full error: ", err)
				break
			}
			//fmt.Println("len(downMsgPB)-------", len(downMsgPB))

			//返序列化proto
			donwMsg := down.DownMsg{}
			if err := proto.Unmarshal(downMsgPB, &donwMsg); err != nil {
				fmt.Println("err::::::::::", err)
			}
			//fmt.Println("donwMsg::::::", donwMsg)

			//读取Key长度
			fmt.Println("read key......")
			keyLenBuffer := make([]byte, 4)
			if _, err := io.ReadFull(conn, keyLenBuffer); err != nil {
				fmt.Println("io read full keyLenBuffer error: ", err)
				break
			}
			keyLen := uint32(0)
			dataBuffer := bytes.NewReader(keyLenBuffer)
			binary.Read(dataBuffer, binary.LittleEndian, &keyLen)
			fmt.Println("keyLen:::::::::", keyLen)

			//读取Key
			keyBuffer := make([]byte, keyLen-1)
			if _, err := io.ReadFull(conn, keyBuffer); err != nil {
				fmt.Println("io read full keyBuffer error: ", err)
				break
			}
			fmt.Println("keyBuffer:::::::::", string(keyBuffer))

			//读取0
			zeroBuffer := make([]byte, 1)
			if _, err := io.ReadFull(conn, zeroBuffer); err != nil {
				fmt.Println("io read full zeroBuffer error: ", err)
				break
			}
			zero := int8(-1)
			dataBuffer = bytes.NewReader(zeroBuffer)
			binary.Read(dataBuffer, binary.LittleEndian, &zero)
			fmt.Println("zero:::::::::", zero)

			loginReply := donwMsg.GetXLoginReply()
			fmt.Println("loginReply.result::::::", loginReply.GetXResult())
			fmt.Println("loginReply.userData::::", loginReply.GetXUserData())

			fmt.Println("------------------------------------------")
			break
		}
	}()

	time.Sleep(60 * time.Second)
	conn.Close()
}

//test
func testServer(conn net.Conn) {
	doLoginProxy(conn)

	//userId := uint32(867) //当前游戏玩家ID
	var repeat uint32 = 0
	var serverId int32 = 1
	upMsg := &up.UpMsg{
		XRepeat: &repeat,
		XLogin: &up.Login{
			XServerId: proto.Int32(serverId),
			XDeviceId: proto.String("test-server-867"),
			XSeq:      proto.Int32(1),
		},
	}

	//pack
	upPb, _ := proto.Marshal(upMsg)
	fmt.Println("pb.len---len(upPb)::::::::::::::::::", len(upPb))
	fmt.Println("---------------------------------------")

	msg := pack.NewMessageUp(0, upPb, []byte("867"), []byte("key-20230609"), []uint64{100, 101}, []uint32{200, 201})
	//fmt.Println("msg.MsgLen:::::::::::::::::", msg.GetMsgLen())
	//fmt.Println("msg.len(pb)::::::::::::::::", len(msg.GetPb()))
	//fmt.Println("---------------------------------------")

	dp := pack.NewDataPackUp()
	data, err := dp.Pack(msg)
	fmt.Println("dp.pack-err:::::::::::::::::", err)
	fmt.Println("dp.pack-len(data):::::::::::", len(data))
	fmt.Println("---------------------------------------")

	//send message
	n, err := conn.Write(data)
	fmt.Println("conn.write-err(testServer)::::::::", n, err)
	fmt.Println("--------------------------------------- \n")

	time.Sleep(60 * time.Second)
	conn.Close()
}
