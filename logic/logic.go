package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game/clinet"
	"game/global"
	"game/msg"
	"game/pack"
	"game/resources"
	"game/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	tilesImage      *ebiten.Image
	ckxImage        *ebiten.Image
	atmImage        *ebiten.Image
	basketballImage *ebiten.Image
	chickenImage    *ebiten.Image
	gameFont60      font.Face
	gameFont24      font.Face
	gameFontB24     font.Face
)

func reverseImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	reversedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取原始像素点的颜色值
			originalColor := img.At(x, y)
			// 计算反转后的像素点位置
			reversedX := width - x - 1
			reversedY := y
			// 设置反转后的像素点颜色值
			reversedImg.Set(reversedX, reversedY, originalColor)
		}
	}
	return reversedImg
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, _ := image.Decode(bytes.NewReader(resources.Tiles_png))
	tilesImage = ebiten.NewImageFromImage(img)

	img2, _, _ := image.Decode(bytes.NewReader(resources.T2_png))
	atmImage = ebiten.NewImageFromImage(reverseImage(img2))
	img3, _, _ := image.Decode(bytes.NewReader(resources.T3_png))
	ckxImage = ebiten.NewImageFromImage(img3)
	img5, _, _ := image.Decode(bytes.NewReader(resources.T5_png))
	basketballImage = ebiten.NewImageFromImage(img5)
	img6, _, _ := image.Decode(bytes.NewReader(resources.T6_png))
	chickenImage = ebiten.NewImageFromImage(reverseImage(img6))
	whiteImage.Fill(color.White)
	tt1, err := opentype.Parse(resources.Font1)
	if err != nil {
		log.Fatal(err)
	}
	tt2, err := opentype.Parse(resources.Font2)
	if err != nil {
		log.Fatal(err)
	}
	gameFont60, err = opentype.NewFace(tt1, &opentype.FaceOptions{
		Size:    60,
		DPI:     global.Dpi,
		Hinting: font.HintingVertical,
	})
	gameFont24, err = opentype.NewFace(tt1, &opentype.FaceOptions{
		Size:    24,
		DPI:     global.Dpi,
		Hinting: font.HintingVertical,
	})
	gameFontB24, err = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    24,
		DPI:     global.Dpi,
		Hinting: font.HintingVertical,
	})
}

type Game struct {
	layers      [][]int
	stars       [starsCount]Star
	keys        []ebiten.Key
	UserA       *ModelInfo
	UserB       *ModelInfo
	IsA         bool
	Status      int   // 1 待开始 2 战斗中 3 战斗结束 4 暂停
	Page        uint8 // 主页面 2战斗页面 3 服务器页面
	ServiceList *ui.Server
	RoomId      uint64
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
	X        float64
	Y        float64
	blood    float32
	Skill    []Skill
	lastTime time.Time
	duration time.Duration
}

func (g *Game) Update() error {
	if g.Page == 3 {
		g.ServiceList.Update()
	} else {
		g.keys = inpututil.AppendPressedKeys(g.keys[:0])
		g.HandleCtrl(g.keys)
		g.HandleRemoteCtrl()
		x, y := ebiten.CursorPosition()
		for i := 0; i < starsCount; i++ {
			g.stars[i].Update(float32(x*scale), float32(y*scale))
		}
		if g.Status == 2 {
			if g.UserA.blood == 0 {
				g.Status = 3
			}
			if g.UserB.blood == 0 {
				g.Status = 3
			}
		}
	}
	return nil
}

func (g *Game) RoomJoin() {
	time.Sleep(time.Second)
	for data := range clinet.RoomChannel {
		g.RoomId = data.RoomId
		g.Page = 1
		g.Status = 2
		g.IsA = data.IsA
	}
}

func (g *Game) HandleCtrl(keys []ebiten.Key) {
	for _, key := range keys {
		switch key {
		case ebiten.KeyEnter:
			if g.Status != 2 {
				g.Status = 2
				g.UserA.blood = 100
				g.UserB.blood = 100
			}
		default:
			if g.Status != 2 {
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
	if clinet.GameRoomInfo.RoomId != 0 {
		if g.IsA {
			//g.GetMeObj().X = clinet.GameRoomInfo.UserA.X
			g.GetRivalObj().X = clinet.GameRoomInfo.UserB.X
			//g.GetMeObj().Y = clinet.GameRoomInfo.UserA.Y
			g.GetRivalObj().Y = clinet.GameRoomInfo.UserB.Y
			g.GetMeObj().blood = clinet.BloodResp.Blood
		} else {
			//g.GetMeObj().X = clinet.GameRoomInfo.UserB.X
			g.GetRivalObj().X = clinet.GameRoomInfo.UserA.X
			//g.GetMeObj().Y = clinet.GameRoomInfo.UserB.Y
			g.GetRivalObj().Y = clinet.GameRoomInfo.UserA.Y
			g.GetMeObj().blood = clinet.BloodResp.Blood
		}
	}
}

func (g *Game) handleKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyA:
		if g.GetMeObj().X > 0 {
			g.GetMeObj().X -= 3
		}
	case ebiten.KeyS:
		if g.GetMeObj().Y+global.ModelHeight < global.ScreenHeight {
			g.GetMeObj().Y += 3
		}
	case ebiten.KeyD:
		if g.GetMeObj().X+global.ModelWidth < global.ScreenWidth {
			g.GetMeObj().X += 3
		}
	case ebiten.KeyW:
		if g.GetMeObj().Y > 0 {
			g.GetMeObj().Y -= 3
		}
	}
	msgData, _ := json.Marshal(msg.MoveReq{
		RoomId: g.RoomId,
		Id:     clinet.Uid,
		X:      g.GetMeObj().X,
		Y:      g.GetMeObj().Y,
	})
	pack.Send(clinet.GetConn(), msg.MsgMove, msgData)
}

func (g *Game) handleSkill(key ebiten.Key) {
	switch key {
	case ebiten.KeyJ:
		msgData, _ := json.Marshal(msg.SkillReq{
			Id:   clinet.Uid,
			X:    g.GetMeObj().X,
			Y:    g.GetMeObj().Y,
			Type: 1,
		})
		pack.Send(clinet.GetConn(), msg.MsgSkill, msgData)
		g.GetMeObj().Skill = append(g.GetMeObj().Skill, &SkillA{
			SkillBase{
				Type: 1,
			},
		})
	case ebiten.KeyK:
		msgData, _ := json.Marshal(msg.SkillReq{
			Id:   clinet.Uid,
			X:    g.GetMeObj().X,
			Y:    g.GetMeObj().Y,
			Type: 2,
		})
		pack.Send(clinet.GetConn(), msg.MsgSkill, msgData)
		g.GetMeObj().Skill = append(g.GetMeObj().Skill, &SkillA{
			SkillBase{
				Type: 2,
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
	if g.Page == 3 {
		g.ServiceList.Draw(screen)
	} else if g.Page == 1 {
		if g.Status == 1 {
			text.Draw(screen, "开始游戏", gameFont60, global.ScreenWidth/2-60*2, global.ScreenHeight/2, color.White)
			text.Draw(screen, "Enter", gameFont24, global.ScreenWidth/2-24*2, global.ScreenHeight/2+80, color.White)
		}
		if g.Status == 3 {
			text.Draw(screen, "游戏结束", gameFont60, global.ScreenWidth/2-60*2, global.ScreenHeight/2, color.White)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.1, 0.1)
		op.GeoM.Translate(g.UserA.X, g.UserA.Y)
		text.Draw(screen, "CXK", gameFont24, global.ScreenWidth/2-60-72, 50, color.White)
		text.Draw(screen, "VS", gameFont60, global.ScreenWidth/2-60, 70, color.RGBA{
			R: 0xFF,
			G: 0,
			B: 0,
			A: 1,
		})
		text.Draw(screen, "ATM", gameFont24, global.ScreenWidth/2+60, 50, color.White)
		if !(g.Status == 3 || g.Status == 4) {
			for i := 0; i < starsCount; i++ {
				g.stars[i].Draw(screen)
			}
		}
		screen.DrawImage(ckxImage, op)
		op1 := &ebiten.DrawImageOptions{}
		op1.GeoM.Scale(0.1, 0.1)
		op1.GeoM.Translate(g.UserB.X, g.UserB.Y)
		screen.DrawImage(atmImage, op1)
		g.DrawBlood(screen)
		{
			for index, skill := range g.GetMeObj().Skill {
				if skill.Handle(screen, g.GetMeObj().X, g.GetMeObj().Y, g.GetRivalObj(), g.RoomId) {
					g.GetMeObj().Skill = append(g.GetMeObj().Skill[:index], g.GetMeObj().Skill[index+1:]...)
				}
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
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
		vs[i].ColorA = 1
	}
	op1 := &ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vs, is, whiteSubImage, op1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return global.ScreenWidth, global.ScreenHeight
}

func Handle() {
	ebiten.SetWindowSize(global.ScreenWidth, global.ScreenHeight)
	ebiten.SetWindowTitle("测试")
	go clinet.Init()
	GameModel = NewGame()
	if err := ebiten.RunGameWithOptions(GameModel, &ebiten.RunGameOptions{
		InitUnfocused: true,
	}); err != nil {
		log.Fatal(err)
	}
}

var GameModel *Game

const (
	scale      = 64
	starsCount = 1024
)

type Star struct {
	fromx, fromy, tox, toy, brightness float32
}

func (s *Star) Init() {
	s.tox = rand.Float32() * global.ScreenWidth * scale
	s.fromx = s.tox
	s.toy = rand.Float32() * global.ScreenHeight * scale
	s.fromy = s.toy
	s.brightness = rand.Float32() * 0xff
}

func (s *Star) Update(x, y float32) {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - x) / 32
	s.toy += (s.toy - y) / 32
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}
	if s.fromx < 0 || global.ScreenWidth*scale < s.fromx || s.fromy < 0 || global.ScreenHeight*scale < s.fromy {
		s.Init()
	}
}

func (s *Star) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(screen, s.fromx/scale, s.fromy/scale, s.tox/scale, s.toy/scale, 1, c, true)
}

func NewGame() *Game {
	g := &Game{
		Status: 1,
		Page:   3,
	}
	for i := 0; i < starsCount; i++ {
		g.stars[i].Init()
	}
	g.UserA = &ModelInfo{
		X:        100,
		Y:        global.ScreenHeight / 2,
		blood:    100,
		duration: 150 * time.Millisecond,
	}

	g.UserB = &ModelInfo{
		X:        global.ScreenWidth / 4 * 3,
		Y:        global.ScreenHeight / 2,
		blood:    100,
		duration: 150 * time.Millisecond,
	}
	g.ServiceList = ui.NewService()
	g.ServiceList.Back = func() {
		g.Page = 1
	}
	g.ServiceList.Join = func(roomId uint64, isA bool) {
		g.RoomId = roomId
		g.IsA = isA
	}
	go g.RoomJoin()
	return g
}

type Skill interface {
	Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo, roomId uint64) bool
}

type SkillBase struct {
	Type      int     // 1 扩散 2 射线 3 todo
	Name      string  // 名称
	Damage    int     // 伤害值
	CD        int     // 技能CD
	Releasing bool    // 技能是否释放中
	Direction uint    // 1 左边 2 右边
	X         float64 // x坐标
	Y         float64 // y坐标
	sync.Mutex
	once sync.Once
}

type SkillA struct {
	SkillBase
}

func (s *SkillA) Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo, roomId uint64) bool {
	s.once.Do(func() {
		s.X = x
		s.Y = y
	})
	fmt.Println("")
	if s.X <= -40 || s.X > global.ScreenWidth+40 {
		return true
	}

	if s.Y <= -40 || s.Y > global.ScreenHeight+40 {
		return true
	}
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.1, 0.1)
	op1.GeoM.Translate(s.X, s.Y)
	s.X += 5
	fmt.Printf(" 人物坐标 [%f,%f] 技能坐标 [%f,%f] \n", obj.X, obj.Y, s.X, s.Y)
	if math.Abs(s.X-obj.X) <= 5 && math.Abs(s.X-obj.X) <= 50 {
		obj.blood--
		msgData, _ := json.Marshal(msg.BloodReq{
			RoomId: roomId,
			Id:     clinet.Uid,
			Blood:  obj.blood,
		})
		pack.Send(clinet.GetConn(), msg.MsgBlood, msgData)
		return true
	}
	switch s.Type {
	case 1:
		screen.DrawImage(basketballImage, op1)
	case 2:
		screen.DrawImage(chickenImage, op1)
	}
	return false
}
