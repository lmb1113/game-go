package main

import (
	"errors"
	"game/global"
	"game/msg"
	"game/utils/pkg/flake"
	"sync"
)

type RoomService struct{}

var GameRoomMar sync.Map

func GetGameRoom(roomId uint64) (*msg.GameRoom, bool) {
	room, has := GameRoomMar.Load(roomId)
	if has {
		return room.(*msg.GameRoom), true
	}
	return nil, false
}

func (r *RoomService) GetGameRoom(roomId uint64) (*msg.GameRoom, bool) {
	room, has := GameRoomMar.Load(roomId)
	if has {
		return room.(*msg.GameRoom), true
	}
	return nil, false
}

func (r *RoomService) Create(id string, nickname string) uint64 {
	roomId, _ := flake.GetID()
	room := &msg.GameRoom{
		RoomId: roomId,
		UserA: &msg.ModelInfo{
			UserId: id,
			X:      100,
			Y:      global.ScreenHeight / 2,
			Blood:  100,
		},
		RoomName: nickname + "创建的房间",
		Number:   1,
	}
	GameRoomMar.Store(roomId, room)
	return roomId
}

func (r *RoomService) Delete() {

}

func (r *RoomService) List() []*msg.GameRoom {
	var resp []*msg.GameRoom
	GameRoomMar.Range(func(key, value any) bool {
		resp = append(resp, value.(*msg.GameRoom))
		return true
	})
	return resp
}

func (r *RoomService) Join(id string, roomId uint64) error {
	room, ok := r.GetGameRoom(roomId)
	if !ok {
		return errors.New("房间不存在")
	}
	if room.Number == 2 {
		return errors.New("房间已满")
	}
	room.Number++
	room.UserB = &msg.ModelInfo{
		UserId: id,
		Blood:  100,
		X:      global.ScreenWidth / 4 * 3,
		Y:      global.ScreenHeight / 2,
	}
	return nil
}
