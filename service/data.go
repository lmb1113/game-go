package main

import (
	"encoding/json"
	"game/global"
	"game/msg"
	"game/pack"
)

func (r *RoomService) HandleMove(req *msg.MoveReq) {
	room, ok := r.GetGameRoom(req.RoomId)
	if !ok {
		return
	}
	if room.UserA != nil && room.UserA.UserId == req.Id {
		room.UserA.X = req.X
		room.UserA.Y = req.Y
	}
	if room.UserB != nil && room.UserB.UserId == req.Id {
		room.UserB.X = req.X
		room.UserB.Y = req.Y
	}
	data, _ := json.Marshal(room)
	if room.UserA != nil {
		if conn, has := getConn(room.UserA.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
	if room.UserB != nil {
		if conn, has := getConn(room.UserB.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
}

func (r *RoomService) HandleBlood(req *msg.BloodReq) {
	room, ok := r.GetGameRoom(req.RoomId)
	if !ok {
		return
	}
	var resp msg.BloodResp
	resp.Blood = req.Blood
	respJson, _ := json.Marshal(resp)
	if room.UserA != nil && room.UserB != nil {
		if room.UserA.UserId == req.Id {
			room.UserB.Blood = req.Blood
			if conn, has := getConn(room.UserB.UserId); has {
				pack.Send(conn, msg.MsgBloodResp, respJson)
			}
		} else {
			room.UserA.Blood = req.Blood
			if conn, has := getConn(room.UserA.UserId); has {
				pack.Send(conn, msg.MsgBloodResp, respJson)
			}
		}
	}
}

func (r *RoomService) InitPlayData(roomId uint64) {
	room, ok := r.GetGameRoom(roomId)
	if !ok {
		return
	}
	room.UserA = &msg.ModelInfo{
		UserId:   room.UserA.UserId,
		UserName: "玩家A",
		X:        100,
		Y:        global.ScreenHeight / 2,
		Blood:    100,
	}
	room.UserB = &msg.ModelInfo{
		UserId:   room.UserB.UserId,
		UserName: "玩家B",
		X:        global.ScreenWidth / 4 * 3,
		Y:        global.ScreenHeight / 2,
		Blood:    100,
	}
	data, _ := json.Marshal(room)
	var resp msg.BloodResp
	resp.Blood = 100
	bloodData, _ := json.Marshal(resp)
	if room.UserA != nil {
		if conn, has := getConn(room.UserA.UserId); has {
			pack.Send(conn, msg.MsgBloodResp, bloodData)
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
	if room.UserB != nil {
		if conn, has := getConn(room.UserB.UserId); has {
			pack.Send(conn, msg.MsgBloodResp, bloodData)
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
}
