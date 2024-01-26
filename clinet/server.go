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

func CreateRoom(conn net.Conn, id string) {
	req := msg.CreateRoomReq{
		Id: id,
	}
	reqJson, _ := json.Marshal(req)
	pack.Send(conn, msg.MsgCreateRoom, reqJson)
	return
}

func JoinRoom(conn net.Conn, id string, roomId uint64) {
	req := msg.JoinRoomReq{
		Id:     id,
		RoomId: roomId,
	}
	reqJson, _ := json.Marshal(req)
	pack.Send(conn, msg.MsgJoinRoom, reqJson)
	return
}
