package msg

const (
	Msg          uint16 = 0xFF55
	MsgLogin     uint8  = 101
	MsgLoginResp uint8  = 102
	MsgPlay      uint8  = 103
	MsgPlayResp  uint8  = 104
	MsgSkill     uint8  = 105
	MsgSkillResp uint8  = 106
	MsgMove      uint8  = 107
	MsgMoveResp  uint8  = 108
	MsgBlood     uint8  = 109
	MsgBloodResp uint8  = 110
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

type LoginMsg struct {
	Name string `json:"name"`
}

type LoginMsgResp struct {
	BaseResp
	IsA bool `json:"is_a"`
}

type MoveReq struct {
	Id    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Blood float32 `json:"blood"`
}

type MoveResp struct {
	BaseResp
}

type SkillReq struct {
	Id        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Type      int     `json:"type"`      // 1 篮球 2 鸡
	Direction int     `json:"direction"` // 1 左 2右
}

type MsgBloodReq struct {
	Id    string  `json:"id"`
	Blood float32 `json:"blood"`
}
