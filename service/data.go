package main

import (
	"encoding/json"
	"game/msg"
	"game/pack"
	"sync"
)

type GameRoom struct {
	RoomId uint64     `json:"room_id"`
	UserA  *ModelInfo `json:"user_a"`
	UserB  *ModelInfo `json:"user_b"`
}

type ModelInfo struct {
	UserId   string  `json:"user_id"`
	UserName string  `json:"user_name"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Blood    float32 `json:"blood"`
}

var GameRoomMar sync.Map

func GetGameRoom(roomId uint64) (*GameRoom, bool) {
	room, has := GameRoomMar.Load(roomId)
	if has {
		return room.(*GameRoom), true
	}
	return nil, false
}

func SetGameRoom(roomId uint64, room *GameRoom) {
	GameRoomMar.Store(roomId, room)
}

func (g *GameRoom) HandleMove(req *msg.MoveReq) {
	if g.UserA != nil && g.UserA.UserId == req.Id {
		g.UserA.X = req.X
		g.UserA.Y = req.Y
		g.UserA.Blood = req.Blood
	}
	if g.UserB != nil && g.UserB.UserId == req.Id {
		g.UserB.X = req.X
		g.UserB.Y = req.Y
		g.UserB.Blood = req.Blood
	}
	if g.UserA != nil {
		dataB, _ := json.Marshal(g.UserB)
		if conn, has := getConn(g.UserA.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserA.UserId, dataB)
		}
	}
	if g.UserB != nil {
		dataA, _ := json.Marshal(g.UserA)
		if conn, has := getConn(g.UserB.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserB.UserId, dataA)
		}
	}
}
