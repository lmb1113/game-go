package clinet

import (
	"encoding/json"
	"game/msg"
	"game/pack"
	"net"
)

func GetRoomList(conn net.Conn) {
	pack.Send(conn, msg.MsgRoomList, nil)
	return
}

func CreateRoom(conn net.Conn, userId uint64) {
	req := msg.CreateRoomReq{
		UserId: userId,
	}
	reqJson, _ := json.Marshal(req)
	pack.Send(conn, msg.MsgCreateRoom, reqJson)
	return
}

func JoinRoom(conn net.Conn, userId uint64, roomId uint64) {
	req := msg.JoinRoomReq{
		UserId: userId,
		RoomId: roomId,
	}
	reqJson, _ := json.Marshal(req)
	pack.Send(conn, msg.MsgJoinRoom, reqJson)
	return
}
