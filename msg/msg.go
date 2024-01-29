package msg

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
	Blood     float32 `json:"blood"`
	RoomId    uint64  `json:"room_id"`
	Direction int     `json:"direction"` // 1 左 2右
}

type MoveResp struct {
	BaseResp
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
	RoomId   uint64     `json:"room_id"`
	RoomName string     `json:"room_name"`
	UserA    *ModelInfo `json:"user_a"`
	UserB    *ModelInfo `json:"user_b"`
	UserId   string     `json:"user_id"`
	Number   int        `json:"number"`
	Status   int        `json:"status"` // 1 空闲 2满
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
