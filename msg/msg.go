package msg

import (
	"strconv"
	"sync"
)

const (
	Msg               uint16 = 0xFF55
	MsgLogin          uint8  = 101
	MsgLoginResp      uint8  = 102
	MsgPlay           uint8  = 103
	MsgPlayResp       uint8  = 104
	MsgSkill          uint8  = 105
	MsgSkillResp      uint8  = 106
	MsgMove           uint8  = 107
	MsgMoveResp       uint8  = 108
	MsgBlood          uint8  = 109
	MsgBloodResp      uint8  = 110
	MsgRoomList       uint8  = 111
	MsgRoomListResp   uint8  = 112
	MsgCreateRoom     uint8  = 113
	MsgCreateRoomResp uint8  = 114
	MsgJoinRoom       uint8  = 115
	MsgJoinRoomResp   uint8  = 116
	MsgGameStatusReq  uint8  = 117
	MsgGameStatusResp uint8  = 118
	MsgExitRoom       uint8  = 119
	MsgExitRoomResp   uint8  = 120
)

const (
	CodeOk  = 200
	CodeErr = 500
)

type BaseMsg struct {
	Header  uint16 `json:"header"`
	MsgType uint8  `json:"msg_type"`
	MsgLen  uint32 `json:"msg_len"`
	Data    []byte `json:"data"`
}

type BaseResp struct {
	Code int `json:"code"`
}

type BaseReq struct {
	Id string `json:"id"`
}

type LoginReq struct {
	UserId uint64 `json:"id"`
}

type LoginResp struct {
	BaseResp
}

type LoginMsgResp struct {
	BaseResp
	IsA bool `json:"is_a"`
}

type MoveReq struct {
	UserId    uint64  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	RoomId    uint64  `json:"room_id"`
	Direction int     `json:"direction"` // 1 左 2右
}

type MoveResp struct {
	RoomId    uint64  `json:"room_id"`
	UserId    uint64  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction int     `json:"direction"` // 1 左 2右
}

type SkillReq struct {
	SkillId   uint64  `json:"skill_id"`
	UserId    uint64  `json:"user_id"`
	RoomId    uint64  `json:"room_id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Type      int     `json:"type"`      // 1 篮球 2 鸡
	Direction int     `json:"direction"` // 1 左 2右
	Releasing bool    `json:"releasing"`
}

type BloodReq struct {
	RoomId uint64  `json:"room_id"`
	UserId uint64  `json:"user_id"`
	Blood  float32 `json:"blood"`
}

type BloodResp struct {
	RoomId uint64  `json:"room_id"`
	Id     string  `json:"id"`
	Blood  float32 `json:"blood"`
}

type GetRoomReq struct {
	Id string `json:"id"`
}

type GetRoomResp struct {
	RoomList []*GameRoom `json:"blood"`
}

type CreateRoomReq struct {
	UserId uint64 `json:"user_id"`
}

type CreateRoomResp struct {
	RoomId uint64 `json:"room_id"`
	IsA    bool   `json:"is_a"`
}

type GameRoom struct {
	RoomId      uint64     `json:"room_id"`
	RoomName    string     `json:"room_name"`
	UserA       *ModelInfo `json:"user_a"`
	UserB       *ModelInfo `json:"user_b"`
	UserId      string     `json:"user_id"`
	Number      int        `json:"number"`
	Status      uint8      `json:"status"`       //  1 等待玩家加入 2 游戏中 3 玩家退出
	PlayStatus  uint8      `json:"play_status"`  // status=2时候的状态 1等待玩家进入 2等待开始 3 开始  4 回合结束 5 总结束
	Round       uint8      `json:"round"`        // 回合
	Result      uint8      `json:"result"`       // 回合结果1 玩家A胜利 2 玩家B胜利
	FinalResult uint8      `json:"final_result"` // 最终结果1 玩家A胜利 2 玩家B胜利
	Record      [2]int8    `json:"record"`       // 比赛记录
	sync.Mutex
}

func (g *GameRoom) GetStatusText() string {
	switch g.Status {
	case 1:
		return "等待玩家加入"
	case 2:
		return "游戏中"
	case 3:
		return "玩家退出,请重新开始战局"
	default:
		return ""
	}
}

func (g *GameRoom) GetRecord() string {
	return strconv.Itoa(int(g.Record[0])) + "/" + strconv.Itoa(int(g.Record[1]))
}

func (g *GameRoom) GetPlayStatusText() string {
	switch g.PlayStatus {
	case 1:
		return "等待玩家进入"
	case 2:
		return "等待开始"
	case 3:
		return "开始中"
	case 4:
		return "游戏结束"
	default:
		return ""
	}
}

type ModelInfo struct {
	UserId    uint64  `json:"user_id"`
	UserName  string  `json:"user_name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Blood     float32 `json:"blood"`
	Direction int     `json:"direction"` // 1 左 2右
}

type JoinRoomReq struct {
	UserId uint64 `json:"user_id"`
	RoomId uint64 `json:"room_id"`
}

type JoinRoomResp struct {
	RoomId uint64 `json:"room_id"`
}

type GameStatusReq struct {
	RoomId uint64 `json:"room_id"`
}

type GameStatusResp struct {
	RoomId         uint64 `json:"room_id"`
	Status         uint8  `json:"status"`
	PlayStatus     uint8  `json:"play_status"`
	StatusText     string `json:"status_text"`
	PlayStatusText string `json:"play_status_text"`
}

type RoomReq struct {
	RoomId uint64 `json:"room_id"`
	UserId uint64 `json:"user_id"`
}
