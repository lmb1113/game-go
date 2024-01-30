package main

import (
	"encoding/json"
	"fmt"
	"game/global"
	"game/msg"
	"game/pack"
	"time"
)

func (r *RoomService) HandleMove(req *msg.MoveReq) {
	room, ok := r.GetGameRoom(req.RoomId)
	if !ok {
		return
	}
	var resp = &msg.MoveResp{
		RoomId:    req.RoomId,
		UserId:    req.UserId,
		X:         req.X,
		Y:         req.Y,
		Direction: req.Direction,
	}
	data, _ := json.Marshal(resp)
	if room.UserA != nil && room.UserA.UserId == req.UserId {
		room.UserA.X = req.X
		room.UserA.Y = req.Y
		room.UserA.Direction = req.Direction
		if conn, has := getConn(room.UserB.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
	if room.UserB != nil && room.UserB.UserId == req.UserId {
		room.UserB.X = req.X
		room.UserB.Y = req.Y
		room.UserB.Direction = req.Direction
		if conn, has := getConn(room.UserA.UserId); has {
			pack.Send(conn, msg.MsgMoveResp, data)
		}
	}
}

func (r *RoomService) HandleBlood(req *msg.BloodReq) {
	room, ok := r.GetGameRoom(req.RoomId)
	if !ok {
		return
	}
	var resp = &msg.BloodResp{
		Blood:  req.Blood,
		RoomId: req.RoomId,
	}
	if req.Blood <= 0 {
		if !room.TryLock() {
			return
		}
		room.Unlock()
		room.PlayStatus = 4
		if req.UserId == room.UserA.UserId {
			room.Result = 1
			room.Record[0]++
		} else {
			room.Result = 2
			room.Record[1]++
		}
		// 只打三局
		r.SendGameStatus(req.RoomId)
		room.Round++
		if room.Round >= 3 {
			// 判断总结果
			if room.Record[0] > room.Record[1] {
				room.FinalResult = 1
			} else {
				room.FinalResult = 1
			}
			room.PlayStatus = 5
			r.SendGameStatus(req.RoomId)
		} else {
			go func() {
				time.Sleep(3 * time.Second)
				room.PlayStatus = 3
				r.InitPlayData(req.RoomId)
				r.SendGameStatus(req.RoomId)
			}()
		}

	}
	respJson, _ := json.Marshal(resp)
	if room.UserA != nil && room.UserB != nil {
		if room.UserA.UserId == req.UserId {
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
		UserId:    room.UserA.UserId,
		UserName:  "玩家A",
		X:         100,
		Y:         global.ScreenHeight / 2,
		Blood:     100,
		Direction: 2,
	}
	room.UserB = &msg.ModelInfo{
		UserId:    room.UserB.UserId,
		UserName:  "玩家B",
		X:         global.ScreenWidth / 4 * 3,
		Y:         global.ScreenHeight / 2,
		Blood:     100,
		Direction: 1,
	}
	room.Status = 2
	room.PlayStatus = 3 // todo 开始倒计时

	var userAData = &msg.MoveResp{
		RoomId:    roomId,
		UserId:    room.UserA.UserId,
		X:         room.UserA.X,
		Y:         room.UserA.Y,
		Direction: room.UserA.Direction,
	}
	dataA, _ := json.Marshal(userAData)

	var userBData = &msg.MoveResp{
		RoomId:    roomId,
		UserId:    room.UserB.UserId,
		X:         room.UserB.X,
		Y:         room.UserB.Y,
		Direction: room.UserB.Direction,
	}
	dataB, _ := json.Marshal(userBData)
	var resp msg.BloodResp
	resp.Blood = 100
	bloodData, _ := json.Marshal(resp)
	r.SendGameStatus(roomId)
	fmt.Println(string(dataA), string(dataB))
	if room.UserA != nil {
		if conn, has := getConn(room.UserA.UserId); has {
			pack.Send(conn, msg.MsgBloodResp, bloodData)
			pack.Send(conn, msg.MsgMoveResp, dataB)
		}
	}
	if room.UserB != nil {
		if conn, has := getConn(room.UserB.UserId); has {
			pack.Send(conn, msg.MsgBloodResp, bloodData)
			pack.Send(conn, msg.MsgMoveResp, dataA)
		}
	}
}
