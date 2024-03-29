package main

import (
	"bytes"
	"fmt"
	"game/resources"
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

const (
	screenWidth  = 1080
	screenHeight = 600
	modelWidth   = 50
	modelHeight  = 50
	dpi          = 72
)

const (
	tileSize = 8
)

var (
	tilesImage      *ebiten.Image
	ckxImage        *ebiten.Image
	atmImage        *ebiten.Image
	basketballImage *ebiten.Image
	chickenImage    *ebiten.Image
	gameFont60      font.Face
	gameFont24      font.Face
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
	tt, err := opentype.Parse(resources.Font1)
	if err != nil {
		log.Fatal(err)
	}
	gameFont60, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    60,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	gameFont24, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
}

type Game struct {
	layers [][]int
	stars  [starsCount]Star
	keys   []ebiten.Key
	Cxk    *ModelInfo
	Atm    *ModelInfo
	Status int // 1 待开始 2 战斗中 3 战斗结束 4 暂停]
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
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.HandleCtrl(g.keys)
	x, y := ebiten.CursorPosition()
	for i := 0; i < starsCount; i++ {
		g.stars[i].Update(float32(x*scale), float32(y*scale))
	}

	if g.Status == 2 {
		if g.Cxk.blood == 0 {
			g.Status = 3
		}
		if g.Atm.blood == 0 {
			g.Status = 3
		}
	}
	return nil
}

func (g *Game) HandleCtrl(keys []ebiten.Key) {
	for _, key := range keys {
		switch key {
		case ebiten.KeyEnter:
			if g.Status != 2 {
				g.Status = 2
				g.Cxk.blood = 100
				g.Atm.blood = 100
			}
		default:
			if g.Status != 2 {
				continue
			}
			g.handleKey(key)
			go g.Cxk.Do(func() {
				g.handleSkill(key)
			})
			go g.Atm.Do(func() {
				g.handleSkill(key)
			})
		}
	}
}

func (g *Game) handleKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyA:
		if g.Cxk.X > 0 {
			g.Cxk.X -= 3
		}
	case ebiten.KeyS:
		if g.Cxk.Y+modelHeight < screenHeight {
			g.Cxk.Y += 3
		}
	case ebiten.KeyD:
		if g.Cxk.X+modelWidth < screenWidth {
			g.Cxk.X += 3
		}
	case ebiten.KeyW:
		if g.Cxk.Y > 0 {
			g.Cxk.Y -= 3
		}
	case ebiten.KeyNumpad4:
		if g.Atm.X > 0 {
			g.Atm.X -= 3
		}
	case ebiten.KeyNumpad2:
		if g.Atm.Y+modelHeight < screenHeight {
			g.Atm.Y += 3
		}
	case ebiten.KeyNumpad6:
		if g.Atm.X+modelWidth < screenWidth {
			g.Atm.X += 3
		}
	case ebiten.KeyNumpad8:
		if g.Atm.Y > 0 {
			g.Atm.Y -= 3
		}
	}
}

func (g *Game) handleSkill(key ebiten.Key) {
	switch key {
	case ebiten.KeyJ:
		g.Cxk.Skill = append(g.Cxk.Skill, &SkillA{
			SkillBase{
				Type: 1,
			},
		})
	case ebiten.KeyK:
		g.Cxk.Skill = append(g.Cxk.Skill, &SkillA{
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
	if g.Status == 1 {
		text.Draw(screen, "开始游戏", gameFont60, screenWidth/2-60*2, screenHeight/2, color.White)
		text.Draw(screen, "Enter", gameFont24, screenWidth/2-24*2, screenHeight/2+80, color.White)
	}
	if g.Status == 3 {
		text.Draw(screen, "游戏结束", gameFont60, screenWidth/2-60*2, screenHeight/2, color.White)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.1, 0.1)
	op.GeoM.Translate(g.Cxk.X, g.Cxk.Y)
	text.Draw(screen, "CXK", gameFont24, screenWidth/2-60-72, 50, color.White)
	text.Draw(screen, "VS", gameFont60, screenWidth/2-60, 70, color.RGBA{
		R: 0xFF,
		G: 0,
		B: 0,
		A: 1,
	})
	text.Draw(screen, "ATM", gameFont24, screenWidth/2+60, 50, color.White)
	if !(g.Status == 3 || g.Status == 4) {
		for i := 0; i < starsCount; i++ {
			g.stars[i].Draw(screen)
		}
	}
	screen.DrawImage(ckxImage, op)
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.1, 0.1)
	op1.GeoM.Translate(g.Atm.X, g.Atm.Y)
	screen.DrawImage(atmImage, op1)
	g.DrawBlood(screen)
	{
		for _, skill := range g.Cxk.Skill {
			skill.Handle(screen, g.Cxk.X, g.Cxk.Y, g.Atm)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) DrawBlood(screen *ebiten.Image) {
	var path vector.Path
	path.MoveTo(10, 100)
	path.LineTo(4*g.Cxk.blood, 100)
	path.LineTo(4*g.Cxk.blood, 120)
	path.LineTo(10, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", g.Cxk.blood), int(4*g.Cxk.blood)+30, 100)
	path.Close()

	path.MoveTo(screenWidth-10, 100)
	path.LineTo(screenWidth-4*g.Atm.blood, 100)
	path.LineTo(screenWidth-4*g.Atm.blood, 120)
	path.LineTo(screenWidth-10, 120)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", g.Atm.blood), int(screenWidth-4*g.Atm.blood)-30, 100)
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
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("测试")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

const (
	scale      = 64
	starsCount = 1024
)

type Star struct {
	fromx, fromy, tox, toy, brightness float32
}

func (s *Star) Init() {
	s.tox = rand.Float32() * screenWidth * scale
	s.fromx = s.tox
	s.toy = rand.Float32() * screenHeight * scale
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
	if s.fromx < 0 || screenWidth*scale < s.fromx || s.fromy < 0 || screenHeight*scale < s.fromy {
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
	}
	for i := 0; i < starsCount; i++ {
		g.stars[i].Init()
	}
	g.Cxk = &ModelInfo{
		X:        100,
		Y:        screenHeight / 2,
		blood:    0,
		duration: 150 * time.Millisecond,
	}

	g.Atm = &ModelInfo{
		X:        screenWidth / 4 * 3,
		Y:        screenHeight / 2,
		blood:    0,
		duration: 150 * time.Millisecond,
	}
	return g
}

type Skill interface {
	Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo)
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

func (s *SkillA) Handle(screen *ebiten.Image, x float64, y float64, obj *ModelInfo) {
	s.once.Do(func() {
		s.X = x
		s.Y = y
	})
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.1, 0.1)
	op1.GeoM.Translate(s.X, s.Y)
	s.X += 5
	if s.X-obj.X == 0 && math.Abs(s.Y-obj.Y) <= 50 {
		obj.blood--
	}
	switch s.Type {
	case 1:
		screen.DrawImage(basketballImage, op1)
	case 2:
		screen.DrawImage(chickenImage, op1)
	}
}
