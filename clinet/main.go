package clinet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"game/msg"
	"game/pack"
	"net"
	"time"
)

var conn net.Conn
var Uid string

func GetConn() net.Conn {
	return conn
}

type ModelInfo struct {
	UserId   string  `json:"user_id"`
	UserName string  `json:"user_name"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Blood    float32 `json:"blood"`
}

type GameRoom struct {
	RoomId uint64     `json:"room_id"`
	UserA  *ModelInfo `json:"user_a"`
	UserB  *ModelInfo `json:"user_b"`
}

var GameRoomInfo GameRoom

var LoginResp msg.LoginMsgResp

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
	//uuid, _ := uuid.NewUUID()git remote add origin git@github.com:lmb1113/game-go.git
	Uid = fmt.Sprintf("%d", time.Now().Unix())
	pack.Send(conn, msg.MsgLogin, Uid, nil)
	go handleConnection(conn)
	fmt.Println("已连接到服务器")
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
			json.Unmarshal(scannedPack.Msg, &GameRoomInfo)
		case msg.MsgLoginResp:
			json.Unmarshal(scannedPack.Msg, &LoginResp)
		}
	}
}
