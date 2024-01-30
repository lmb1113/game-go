package clinet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"game/msg"
	"game/pack"
	"game/utils/pkg/flake"
	"net"
)

var conn net.Conn
var Uid uint64

func GetConn() net.Conn {
	return conn
}

type ModelInfo struct {
	UserId    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Blood     float32 `json:"blood"`
	Direction int     `json:"direction"` // 1 左 2右
}

var GameRoomInfo msg.GameRoom
var MoveResp msg.MoveResp

var LoginResp msg.LoginMsgResp
var RoomResp msg.GetRoomResp
var BloodResp msg.BloodResp
var createRoomResp msg.CreateRoomResp
var RoomChannel chan msg.CreateRoomResp
var SkillChannel chan []byte

func Init() {
	serverAddr, err := net.ResolveTCPAddr("tcp", "192.168.31.245:9090")
	if err != nil {
		fmt.Println("解析服务器地址失败：", err)
		return
	}

	conn, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		fmt.Println("连接服务器失败：", err)
		return
	}

	defer conn.Close()
	Uid, _ = flake.GetID()
	var loginReq msg.LoginReq
	loginReq.UserId = Uid
	loginReqJson, _ := json.Marshal(loginReq)
	pack.Send(conn, msg.MsgLogin, loginReqJson)
	go handleConnection(conn)
	fmt.Println("已连接到服务器")
	RoomChannel = make(chan msg.CreateRoomResp, 100)
	SkillChannel = make(chan []byte, 100)
	select {}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // 关闭连接
	fmt.Println(conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'V' {
			if len(data) > 4 {
				length := int16(0)
				binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &length)
				if int(length)+4 <= len(data) {
					return int(length) + 4, data[:int(length)+4], nil
				}
				fmt.Println("格式错误")
			}
		}
		return
	})
	for scanner.Scan() {
		scannedPack := new(pack.Package)
		err := scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		if err != nil {
			fmt.Println("拆包错误", err.Error())
			continue
		}
		fmt.Printf("%+v", scannedPack)
		switch scannedPack.MsgType {
		case msg.MsgMoveResp:
			json.Unmarshal(scannedPack.Msg, &MoveResp)
		case msg.MsgLoginResp:
			json.Unmarshal(scannedPack.Msg, &LoginResp)
		case msg.MsgRoomListResp:
			json.Unmarshal(scannedPack.Msg, &RoomResp)
		case msg.MsgBloodResp:
			json.Unmarshal(scannedPack.Msg, &BloodResp)
		case msg.MsgJoinRoomResp:
			json.Unmarshal(scannedPack.Msg, &createRoomResp)
			createRoomResp.IsA = false
			RoomChannel <- createRoomResp
			fmt.Println("加入房间成功")
		case msg.MsgCreateRoomResp:
			json.Unmarshal(scannedPack.Msg, &createRoomResp)
			createRoomResp.IsA = true
			RoomChannel <- createRoomResp
		case msg.MsgSkillResp:
			fmt.Println("对方释放技能")
			SkillChannel <- scannedPack.Msg
		case msg.MsgGameStatusResp:
			fmt.Println("收到房间状态响应")
			json.Unmarshal(scannedPack.Msg, &GameRoomInfo)
		}
	}
}
