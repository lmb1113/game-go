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
	defer func(addr string) {
		fmt.Println("=========用户离线=========", addr)
		conn.Close() // 关闭连接
	}(conn.RemoteAddr().String())
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
			resp := &msg.LoginMsgResp{
				BaseResp: msg.BaseResp{
					Code: msg.CodeOk,
				},
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgLoginResp, respData)
		case msg.MsgMove:
			fmt.Println("移动成功")
			var msgData msg.MoveReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			new(RoomService).HandleMove(&msgData)
			resp := &msg.LoginMsgResp{
				BaseResp: msg.BaseResp{
					Code: msg.CodeOk,
				},
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgMoveResp, respData)
		case msg.MsgBlood:
			fmt.Println("血量上报", scannedPack.Msg)
			var msgData msg.BloodReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			new(RoomService).HandleBlood(&msgData)
		case msg.MsgCreateRoom:
			fmt.Println("创建房间", scannedPack.Msg)
			var msgData msg.CreateRoomReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			roomId := new(RoomService).Create(msgData.Id, string(scannedPack.Hostname))
			resp := &msg.CreateRoomResp{
				RoomId: roomId,
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgCreateRoomResp, respData)
		case msg.MsgRoomList:
			fmt.Println("获取房间列表", scannedPack.Msg)
			list := new(RoomService).List()
			resp := &msg.GetRoomResp{
				RoomList: list,
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgRoomListResp, respData)
		case msg.MsgJoinRoom:
			fmt.Println("加入房间", scannedPack.Msg)
			var msgData msg.JoinRoomReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			err := new(RoomService).Join(msgData.Id, msgData.RoomId)
			if err != nil {
				fmt.Println(err)
			}
			resp := &msg.JoinRoomResp{
				RoomId: msgData.RoomId,
			}
			respData, _ := json.Marshal(resp)
			pack.Send(conn, msg.MsgJoinRoomResp, respData)
			new(RoomService).InitPlayData(msgData.RoomId)
		case msg.MsgSkill:
			fmt.Println("使用技能", scannedPack.Msg)
			var msgData msg.SkillReq
			json.Unmarshal(scannedPack.Msg, &msgData)
			err := new(RoomService).HandleSkill(msgData)
			if err != nil {
				fmt.Println(err)
			}
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
