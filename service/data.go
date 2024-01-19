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
	}
	if g.UserB != nil && g.UserB.UserId == req.Id {
		g.UserB.X = req.X
		g.UserB.Y = req.Y
	}
	data, _ := json.Marshal(g)
	if g.UserA != nil {
		if conn, has := getConn(g.UserA.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserA.UserId, data)
		}
	}
	if g.UserB != nil {
		if conn, has := getConn(g.UserB.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserB.UserId, data)
		}
	}
}

func (g *GameRoom) HandleBlood(req *msg.MsgBloodReq) {
	if g.UserA != nil && g.UserB != nil {
		if g.UserA.UserId == req.Id {
			g.UserB.Blood = req.Blood
		} else {
			g.UserA.Blood = req.Blood
		}
	}
	data, _ := json.Marshal(g)
	if g.UserA != nil {
		if conn, has := getConn(g.UserA.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserA.UserId, data)
		}
	}
	if g.UserB != nil {
		if conn, has := getConn(g.UserB.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, g.UserB.UserId, data)
		}
	}
}
