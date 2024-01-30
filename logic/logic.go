package logic

import (
	"encoding/json"
	"fmt"
	"game/clinet"
	"game/global"
	"game/msg"
	"game/pack"
	"game/resources"
	"game/ui"
	"game/utils/pkg/flake"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"sync"
	"time"
)

type Game struct {
	layers      [][]int
	keys        []ebiten.Key
	UserA       *ModelInfo
	UserB       *ModelInfo
	IsA         bool
	ServiceList *ui.Server
	RoomId      uint64
	Status      uint8 // 1 主页面 2 未加入服务器 3 连接到服务器 4 已加入游戏
	sync.Mutex
}

func (g *Game) GetMeObj() *ModelInfo {
	if g.IsA {
		return g.UserA
	}
	return g.UserB
}

func (g *Game) GetRivalObj() *ModelInfo {
	if !g.IsA {
		return g.UserA
	}
	return g.UserB
}

type ModelInfo struct {
	X         float64
	Y         float64
	blood     float32
	Skill     sync.Map
	lastTime  time.Time
	duration  time.Duration
	Direction int `json:"direction"` // 1 左 2右
}

func (g *Game) Update() error {
	if g.Status == 2 {
		g.ServiceList.Update()
	} else {
		g.keys = inpututil.AppendPressedKeys(g.keys[:0])
		g.HandleCtrl(g.keys)
		g.HandleRemoteCtrl()
	}
	if clinet.GameRoomInfo.PlayStatus == 4 || clinet.GameRoomInfo.PlayStatus == 5 {
		g.UserA = &ModelInfo{
			X:         100,
			Y:         global.ScreenHeight / 2,
			blood:     100,
			duration:  150 * time.Millisecond,
			Direction: 2,
		}

		g.UserB = &ModelInfo{
			X:         global.ScreenWidth / 4 * 3,
			Y:         global.ScreenHeight / 2,
			blood:     100,
			duration:  150 * time.Millisecond,
			Direction: 1,
		}
	}
	return nil
}

func (g *Game) RoomJoin() {
	time.Sleep(time.Second)
	for data := range clinet.RoomChannel {
		g.RoomId = data.RoomId
		g.Status = 4
		g.IsA = data.IsA
	}
}

func (g *Game) RemoteSkill() {
	time.Sleep(time.Second)
	for data := range clinet.SkillChannel {
		var skill SkillA
		err := json.Unmarshal(data, &skill)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("对方释放一个技能")
		if g.IsA {
			g.UserB.Skill.Store(skill.SkillId, &skill)
		} else {
			g.UserA.Skill.Store(skill.SkillId, &skill)
		}
	}
}

func (g *Game) HandleCtrl(keys []ebiten.Key) {
	for _, key := range keys {
		switch key {
		case ebiten.KeyEnter:
			if g.Status == 1 {
				g.Status = 2

			}
			if clinet.GameRoomInfo.PlayStatus == 5 {
				msgData, _ := json.Marshal(msg.RoomReq{
					RoomId: g.RoomId,
					UserId: clinet.Uid,
				})
				pack.Send(clinet.GetConn(), msg.MsgExitRoom, msgData)
				g.Status = 2
			}
		default:
			if g.Status != 4 {
				continue
			}
			g.handleKey(key)
			go g.UserA.Do(func() {
				g.handleSkill(key)
			})
		}
	}
}

func (g *Game) HandleRemoteCtrl() {
	if clinet.GameRoomInfo.PlayStatus != 3 {
		return
	}
	g.GetRivalObj().X = clinet.MoveResp.X
	g.GetRivalObj().Y = clinet.MoveResp.Y
	g.GetMeObj().blood = clinet.BloodResp.Blood
	g.GetRivalObj().Direction = clinet.MoveResp.Direction
}

func (g *Game) handleKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyA:
		if g.GetMeObj().X > 0 {
			g.GetMeObj().X -= 3
		}
		g.GetMeObj().Direction = 1
	case ebiten.KeyS:
		if g.GetMeObj().Y+global.ModelHeight < global.ScreenHeight {
			g.GetMeObj().Y += 3
		}
	case ebiten.KeyD:
		if g.GetMeObj().X+global.ModelWidth < global.ScreenWidth {
			g.GetMeObj().X += 3
		}
		g.GetMeObj().Direction = 2
	case ebiten.KeyW:
		if g.GetMeObj().Y > 0 {
			g.GetMeObj().Y -= 3
		}
	default:
		return
	}
	msgData, _ := json.Marshal(msg.MoveReq{
		RoomId:    g.RoomId,
		UserId:    clinet.Uid,
		X:         g.GetMeObj().X,
		Y:         g.GetMeObj().Y,
		Direction: g.GetMeObj().Direction,
	})
	pack.Send(clinet.GetConn(), msg.MsgMove, msgData)
}

func (g *Game) handleSkill(key ebiten.Key) {
	id, _ := flake.GetID()
	switch key {
	case ebiten.KeyJ:
		g.GetMeObj().Skill.Store(id, &SkillA{
			SkillBase{
				Type:      1,
				SkillId:   id,
				Direction: g.GetMeObj().Direction,
			},
		})
	case ebiten.KeyK:
		g.GetMeObj().Skill.Store(id, &SkillA{
			SkillBase{
				Type:      2,
				SkillId:   id,
				Direction: g.GetMeObj().Direction,
			},
		})
	case ebiten.KeyL:
		g.GetMeObj().Skill.Store(id, &SkillA{
			SkillBase{
				Type:      3,
				SkillId:   id,
				Direction: g.GetMeObj().Direction,
			},
		})
	case ebiten.KeyU:
		g.GetMeObj().Skill.Store(id, &SkillA{
			SkillBase{
				Type:      4,
				SkillId:   id,
				Direction: g.GetMeObj().Direction,
			},
		})
	case ebiten.KeyI:
		g.GetMeObj().Skill.Store(id, &SkillA{
			SkillBase{
				Type:      5,
				SkillId:   id,
				Direction: g.GetMeObj().Direction,
			},
		})
	}
}

func (m *ModelInfo) Do(fn func()) {
	currentTime := time.Now()
	if currentTime.Sub(m.lastTime) >= m.duration {
		fn()
		m.lastTime = currentTime
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Status == 1 {
		g.DrawBg(screen)
		ui.Text(screen, "开始游戏", ui.FaceA60, ui.AlignCenter)
		text.Draw(screen, "Enter", resources.GameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.Black)
	}

	if g.Status == 2 {
		g.ServiceList.Draw(screen)
	}

	// 加入游戏
	if g.Status == 4 {
		g.DrawBg(screen)
		if clinet.GameRoomInfo.Status == 1 {
			ui.Text(screen, clinet.GameRoomInfo.GetStatusText(), ui.FaceA60, ui.AlignCenter)
		}
		if clinet.GameRoomInfo.Status == 3 {
			ui.Text(screen, clinet.GameRoomInfo.GetStatusText(), ui.FaceA60, ui.AlignCenter)
		}

		if clinet.GameRoomInfo.Status == 2 {
			g.DrawBg(screen)
			if clinet.GameRoomInfo.PlayStatus == 4 {
				var content string
				if (g.IsA && clinet.GameRoomInfo.Result == 2) || (!g.IsA && clinet.GameRoomInfo.Result == 1) {
					content = "加油呀，您快要输了！"
					text.Draw(screen, "马上开始下一回合", resources.GameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.Black)
				} else {
					content = "再接再厉，马上胜利了！"
					text.Draw(screen, "马上开始下一回合", resources.GameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.Black)
				}
				ui.Text(screen, content, ui.FaceA60, ui.AlignCenter)
			} else if clinet.GameRoomInfo.PlayStatus == 5 {
				var content string
				if (g.IsA && clinet.GameRoomInfo.Result == 2) || (!g.IsA && clinet.GameRoomInfo.Result == 1) {
					content = "很遗憾，您输了！"
					text.Draw(screen, "回车重新开始比赛", resources.GameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.Black)
				} else {
					content = "恭喜您，胜利了！"
					text.Draw(screen, "回车重新开始比赛", resources.GameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.Black)
				}
				ui.Text(screen, content, ui.FaceA60, ui.AlignCenter)
			} else {
				switch clinet.GameRoomInfo.PlayStatus {
				case 1:
					ui.Text(screen, clinet.GameRoomInfo.GetStatusText(), ui.FaceA60, ui.AlignCenter)
				case 2:
					ui.Text(screen, clinet.GameRoomInfo.GetStatusText(), ui.FaceA60, ui.AlignCenter)
				case 3:
				}
				text.Draw(screen, clinet.GameRoomInfo.GetRecord(), resources.GameFont24, global.ScreenWidth/2-24*2, 110, color.RGBA{
					R: 0xFF,
					G: 0,
					B: 0,
					A: 0xff,
				})
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(0.1, 0.1)
				op.GeoM.Translate(g.UserA.X, g.UserA.Y)
				text.Draw(screen, "CXK", resources.GameFont24, global.ScreenWidth/2-60-72, 50, color.Black)
				text.Draw(screen, "VS", resources.GameFont60, global.ScreenWidth/2-60, 70, color.RGBA{
					R: 0xFF,
					G: 0,
					B: 0,
					A: 0xff,
				})
				text.Draw(screen, "ATM", resources.GameFont24, global.ScreenWidth/2+60, 50, color.Black)
				screen.DrawImage(resources.CkxImage, op)
				op1 := &ebiten.DrawImageOptions{}
				op1.GeoM.Scale(0.1, 0.1)
				op1.GeoM.Translate(g.UserB.X, g.UserB.Y)
				screen.DrawImage(resources.AtmImage, op1)
				g.DrawBlood(screen)
				{
					g.GetMeObj().Skill.Range(func(key, value any) bool {
						if value.(*SkillA).Handle(screen, g.GetMeObj().X, g.GetMeObj().Y, g.GetRivalObj(), g.RoomId) {
							g.GetMeObj().Skill.Delete(key)
						}
						return true
					})

					g.GetRivalObj().Skill.Range(func(key, value any) bool {
						if value.(*SkillA).HandleRemote(screen, value.(*SkillA), g.RoomId) {
							g.GetRivalObj().Skill.Delete(key)
						}
						return true
					})
				}
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) DrawBg(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.GeoM.Scale(0.6667, 0.6667)
	screen.DrawImage(resources.BgImage, op)
}

func (g *Game) DrawBlood(screen *ebiten.Image) {
	// 血条
	var path vector.Path
	path.MoveTo(10, 100)
	path.LineTo(4*g.UserA.blood, 100)
	path.LineTo(4*g.UserA.blood, 120)
	path.LineTo(10, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", g.UserA.blood), int(4*g.UserA.blood)+30, 100)
	path.Close()

	path.MoveTo(global.ScreenWidth-10, 100)
	path.LineTo(global.ScreenWidth-4*g.UserB.blood, 100)
	path.LineTo(global.ScreenWidth-4*g.UserB.blood, 120)
	path.LineTo(global.ScreenWidth-10, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", g.UserB.blood), int(global.ScreenWidth-4*g.UserB.blood)-30, 100)
	path.Close()
	var vs []ebiten.Vertex
	var is []uint16
	op := &vector.StrokeOptions{}
	op.Width = 5
	op.LineJoin = vector.LineJoinRound
	vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0xFF
		vs[i].ColorG = 0x00
		vs[i].ColorB = 0x00
		vs[i].ColorA = 0xff
	}
	op1 := &ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vs, is, resources.WhiteSubImage, op1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return global.ScreenWidth, global.ScreenHeight
}

func Handle() {
	ebiten.SetWindowSize(global.ScreenWidth, global.ScreenHeight)
	ebiten.SetWindowTitle("偶像练习生")
	resources.Init()
	go clinet.Init()
	GameModel = NewGame()
	if err := ebiten.RunGameWithOptions(GameModel, &ebiten.RunGameOptions{
		InitUnfocused: true,
	}); err != nil {
		log.Fatal(err)
	}
}

var GameModel *Game

func NewGame() *Game {
	g := &Game{
		Status: 1,
	}
	g.UserA = &ModelInfo{
		X:         100,
		Y:         global.ScreenHeight / 2,
		blood:     100,
		duration:  150 * time.Millisecond,
		Direction: 2,
	}

	g.UserB = &ModelInfo{
		X:         global.ScreenWidth / 4 * 3,
		Y:         global.ScreenHeight / 2,
		blood:     100,
		duration:  150 * time.Millisecond,
		Direction: 1,
	}
	g.ServiceList = ui.NewService()
	g.ServiceList.Back = func() {
		g.Status = 1
	}
	g.ServiceList.Join = func(roomId uint64, isA bool) {
		g.RoomId = roomId
		g.IsA = isA
		g.Status = 4
	}
	g.ServiceList.CreateCall = func() {
	}
	go g.RoomJoin()
	go g.RemoteSkill()
	return g
}

type Skill interface {
	Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo, roomId uint64) bool
	HandleRemote(screen *ebiten.Image, skill *SkillA, roomId uint64) bool
}

type SkillBase struct {
	Type      int     // 1 扩散 2 射线 3 todo
	Name      string  // 名称
	Damage    int     // 伤害值
	CD        int     // 技能CD
	Releasing bool    // 技能是否释放中
	Direction int     // 1 左边 2 右边
	X         float64 // x坐标
	Y         float64 // y坐标
	SkillId   uint64  `json:"skill_id"`
	once      sync.Once
}

type SkillA struct {
	SkillBase
}

func (s *SkillA) Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo, roomId uint64) bool {
	s.once.Do(func() {
		s.X = x
		s.Y = y
		msgData, _ := json.Marshal(msg.SkillReq{
			UserId:    clinet.Uid,
			RoomId:    roomId,
			X:         s.X,
			Y:         s.Y,
			Type:      s.Type,
			SkillId:   s.SkillId,
			Direction: s.Direction,
		})
		fmt.Println("释放技能")
		pack.Send(clinet.GetConn(), msg.MsgSkill, msgData)
	})
	if s.X <= -40 || s.X > global.ScreenWidth+40 {
		return true
	}

	if s.Y <= -40 || s.Y > global.ScreenHeight+40 {
		return true
	}
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.1, 0.1)
	op1.GeoM.Translate(s.X, s.Y)
	if s.Direction == 2 {
		s.X += 8
	} else {
		s.X -= 8
	}
	//fmt.Printf(" 人物坐标 [%f,%f] 技能坐标 [%f,%f] \n", obj.X, obj.Y, s.X, s.Y)
	if math.Abs(s.X-obj.X) <= 5 && math.Abs(s.Y-obj.Y) <= 50 {
		obj.blood--
		msgData, _ := json.Marshal(msg.BloodReq{
			RoomId: roomId,
			UserId: clinet.Uid,
			Blood:  obj.blood,
		})
		if obj.blood <= 0 {
			clinet.GameRoomInfo.PlayStatus = 4
			// 最低只能到0
			obj.blood = 0
		}
		fmt.Println("技能命中")
		pack.Send(clinet.GetConn(), msg.MsgBlood, msgData)
		skillData, _ := json.Marshal(msg.SkillReq{
			UserId:    clinet.Uid,
			RoomId:    roomId,
			X:         s.X,
			Y:         s.Y,
			Type:      s.Type,
			SkillId:   s.SkillId,
			Direction: s.Direction,
			Releasing: true,
		})
		pack.Send(clinet.GetConn(), msg.MsgSkill, skillData)
		return true
	}
	switch s.Type {
	case 1:
		if s.Direction == 1 {
			screen.DrawImage(resources.BasketballImage2, op1)
		} else {
			screen.DrawImage(resources.BasketballImage, op1)
		}
	case 2:
		if s.Direction == 1 {
			screen.DrawImage(resources.ChickenImage2, op1)
		} else {
			screen.DrawImage(resources.ChickenImage, op1)
		}
	case 3:
		text.Draw(screen, "律师函", resources.GameFont24, int(s.X), int(s.Y), color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 0xff,
		})
	case 4:
		if s.Direction == 1 {
			screen.DrawImage(resources.Skill4A, op1)
		} else {
			screen.DrawImage(resources.Skill4B, op1)
		}
	case 5:
		if s.Direction == 1 {
			screen.DrawImage(resources.Skill5A, op1)
		} else {
			screen.DrawImage(resources.Skill5B, op1)
		}
	}
	return false
}

func (s *SkillA) HandleRemote(screen *ebiten.Image, skill *SkillA, roomId uint64) bool {
	s.Y = skill.Y
	s.X = skill.X
	s.Type = skill.Type
	s.Name = skill.Name
	s.Damage = skill.Damage
	s.Releasing = skill.Releasing
	if s.Releasing {
		return true
	}
	if s.X <= -40 || s.X > global.ScreenWidth+40 {
		s.X = 99999
		return true
	}

	if s.Y <= -40 || s.Y > global.ScreenHeight+40 {
		s.Y = 99999
		return true
	}
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.1, 0.1)
	op1.GeoM.Translate(s.X, s.Y)
	if s.Direction == 2 {
		s.X += 8
	} else {
		s.X -= 8
	}
	switch s.Type {
	case 1:
		if s.Direction == 1 {
			screen.DrawImage(resources.BasketballImage2, op1)
		} else {
			screen.DrawImage(resources.BasketballImage, op1)
		}
	case 2:
		if s.Direction == 1 {
			screen.DrawImage(resources.ChickenImage2, op1)
		} else {
			screen.DrawImage(resources.ChickenImage, op1)
		}
	case 3:
		text.Draw(screen, "律师函", resources.GameFont24, int(s.X), int(s.Y), color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 0xff,
		})
	case 4:
		if s.Direction == 1 {
			screen.DrawImage(resources.Skill4A, op1)
		} else {
			screen.DrawImage(resources.Skill4B, op1)
		}
	case 5:
		if s.Direction == 1 {
			screen.DrawImage(resources.Skill5A, op1)
		} else {
			screen.DrawImage(resources.Skill5B, op1)
		}
	}
	return false
}
