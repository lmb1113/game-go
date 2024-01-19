package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"game/msg"
	"game/pack"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close() // 关闭连接
	fmt.Println(conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'V' {
			if len(data) > 4 {
				length := int16(0)
				binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &length)
				if int(length)+4 <= len(data) {
					return int(length) + 4, data[:int(length)+4], nil
				}
			}
		}
		return
	})
	for scanner.Scan() {
		scannedPack := new(pack.Package)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		fmt.Printf("%+v", scannedPack)
		switch scannedPack.MsgType {
		case msg.MsgLogin:
			setConn(string(scannedPack.Hostname), conn)
			room, has := GetGameRoom(12345)
			if !has {
				SetGameRoom(12345, &GameRoom{
					RoomId: 12345,
					UserA: &ModelInfo{
						UserId:   string(scannedPack.Hostname),
						UserName: string(scannedPack.Hostname),
						Blood:    100,
					},
				})
			} else {
				room.UserB = &ModelInfo{
					UserId:   string(scannedPack.Hostname),
					UserName: string(scannedPack.Hostname),
					Blood:    100,
				}
				SetGameRoom(12345, room)
			}
			resp := &msg.LoginMsgResp{
				BaseResp: msg.BaseResp{
					Code: msg.CodeOk,
				},
				IsA: !has,
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgLoginResp, "server", respData)
		case msg.MsgMove:
			room, has := GetGameRoom(12345)
			if !has {
				continue
			}
			fmt.Println("移动成功")
			var msgData msg.MoveReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			room.HandleMove(&msgData)
			resp := &msg.LoginMsgResp{
				BaseResp: msg.BaseResp{
					Code: msg.CodeOk,
				},
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgMoveResp, "server", respData)
		case msg.MsgBlood:
			room, has := GetGameRoom(12345)
			if !has {
				continue
			}
			fmt.Println("移动成功")
			var msgData msg.MsgBloodReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			room.HandleBlood(&msgData)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Println("监听失败：", err)
		return
	}
	defer listener.Close()
	fmt.Println("服务器已启动，监听地址：0.0.0.0:9090")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败：", err)
			continue
		}
		go handleConnection(conn) // 开启一个新的协程处理连接
	}
}
