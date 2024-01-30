package ui

import (
	"bytes"
	"game/clinet"
	"game/global"
	"game/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"
)

const (
	uiFontSize          = 14
	lineSpacingInPixels = 16
)

var (
	uiImage      *ebiten.Image
	uiFaceSource *text.GoTextFaceSource
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.UI_png))
	if err != nil {
		log.Fatal(err)
	}
	uiImage = ebiten.NewImageFromImage(img)
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(resources.Font3))
	if err != nil {
		log.Fatal(err)
	}
	uiFaceSource = s
}

type imageType int

const (
	imageTypeButton imageType = iota
	imageTypeButtonPressed
	imageTypeTextBox
	imageTypeVScrollBarBack
	imageTypeVScrollBarFront
	imageTypeCheckBox
	imageTypeCheckBoxPressed
	imageTypeCheckBoxMark
)

var imageSrcRects = map[imageType]image.Rectangle{
	imageTypeButton:          image.Rect(0, 0, 16, 16),
	imageTypeButtonPressed:   image.Rect(16, 0, 32, 16),
	imageTypeTextBox:         image.Rect(0, 16, 16, 32),
	imageTypeVScrollBarBack:  image.Rect(16, 16, 24, 32),
	imageTypeVScrollBarFront: image.Rect(24, 16, 32, 32),
	imageTypeCheckBox:        image.Rect(0, 32, 16, 48),
	imageTypeCheckBoxPressed: image.Rect(16, 32, 32, 48),
	imageTypeCheckBoxMark:    image.Rect(32, 32, 48, 48),
}

func drawNinePatches(dst *ebiten.Image, dstRect image.Rectangle, srcRect image.Rectangle) {
	srcX := srcRect.Min.X
	srcY := srcRect.Min.Y
	srcW := srcRect.Dx()
	srcH := srcRect.Dy()

	dstX := dstRect.Min.X
	dstY := dstRect.Min.Y
	dstW := dstRect.Dx()
	dstH := dstRect.Dy()

	op := &ebiten.DrawImageOptions{}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			op.GeoM.Reset()

			sx := srcX
			sy := srcY
			sw := srcW / 4
			sh := srcH / 4
			dx := 0
			dy := 0
			dw := sw
			dh := sh
			switch i {
			case 1:
				sx = srcX + srcW/4
				sw = srcW / 2
				dx = srcW / 4
				dw = dstW - 2*srcW/4
			case 2:
				sx = srcX + 3*srcW/4
				dx = dstW - srcW/4
			}
			switch j {
			case 1:
				sy = srcY + srcH/4
				sh = srcH / 2
				dy = srcH / 4
				dh = dstH - 2*srcH/4
			case 2:
				sy = srcY + 3*srcH/4
				dy = dstH - srcH/4
			}

			op.GeoM.Scale(float64(dw)/float64(sw), float64(dh)/float64(sh))
			op.GeoM.Translate(float64(dx), float64(dy))
			op.GeoM.Translate(float64(dstX), float64(dstY))
			dst.DrawImage(uiImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image), op)
		}
	}
}

type Button struct {
	Rect      image.Rectangle
	Text      string
	Type      string
	mouseDown bool
	Dy        float64
	Y0        int
	Y1        int
	onPressed func(b *Button)
}

func (b *Button) Update() {
	b.Y1 += int(b.Dy)
	b.Y0 += int(b.Dy)
	b.Rect = image.Rect(b.Rect.Min.X, b.Y0, b.Rect.Max.X, b.Y1)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if b.Rect.Min.X <= x && x < b.Rect.Max.X && b.Rect.Min.Y <= y && y < b.Rect.Max.Y {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(dst *ebiten.Image) {
	t := imageTypeButton
	if b.mouseDown {
		t = imageTypeButtonPressed
	}
	drawNinePatches(dst, b.Rect, imageSrcRects[t])

	op := &text.DrawOptions{}
	// 按钮文字居中
	op.GeoM.Translate(float64(b.Rect.Min.X+b.Rect.Max.X)/2, float64(b.Rect.Min.Y+b.Rect.Max.Y)/2)
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = lineSpacingInPixels
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	text.Draw(dst, b.Text, &text.GoTextFace{
		Source: uiFaceSource,
		Size:   uiFontSize,
	}, op)
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
}

type Server struct {
	button1        *Button
	button2        *Button
	button3        *Button
	buttonList     []*Button
	buttonListSize int
	Dy             float64
	Back           func()
	Join           func(roomId uint64, isA bool)
	CreateCall     func()
}

func NewService() *Server {
	s := &Server{}
	s.button1 = &Button{
		Rect: image.Rect(16, 16, 116, 48),
		Text: "返回",
		Y0:   16,
		Y1:   48,
	}
	s.button1.SetOnPressed(func(b *Button) {
		s.Back()
	})

	s.button2 = &Button{
		Rect: image.Rect(120, 16, 220, 48),
		Text: "刷新",
		Y0:   16,
		Y1:   48,
	}

	s.button2.SetOnPressed(func(b *Button) {
		clinet.GetRoomList(clinet.GetConn())
		s.buttonListSize = len(clinet.RoomResp.RoomList)*60 - 100
		var buttonList []*Button
		for i, room := range clinet.RoomResp.RoomList {
			button := &Button{
				Rect: image.Rect(16, 100+50*i+10, global.ScreenWidth-16, 100+50*i+50),
				Text: "[" + strconv.FormatUint(room.RoomId, 10) + "] " + room.RoomName + "         人数[" + strconv.Itoa(room.Number) + "/2]",
				Y0:   100 + 50*i + 10,
				Y1:   100 + 50*i + 50,
			}
			button.SetOnPressed(func(b *Button) {
				clinet.JoinRoom(clinet.GetConn(), clinet.Uid, room.RoomId)
				s.Join(room.RoomId, false)
			})
			buttonList = append(buttonList, button)
		}
		s.Dy = 0
		s.buttonList = buttonList
	})
	s.button3 = &Button{
		Rect: image.Rect(224, 16, 324, 48),
		Text: "创建",
		Y0:   16,
		Y1:   48,
	}
	s.button3.SetOnPressed(func(b *Button) {
		clinet.CreateRoom(clinet.GetConn(), clinet.Uid)
		s.CreateCall()
	})
	return s
}

func (s *Server) Update() error {

	// 鼠标滚动服务器
	_, dy := ebiten.Wheel()
	if dy < 0 && int(s.Dy+dy) < -(s.buttonListSize-global.ScreenHeight) {
		dy = 0
	}

	if dy > 0 && int(s.Dy+dy) > 0 {
		dy = 0
	}
	s.button1.Update()
	s.button2.Update()
	s.button3.Update()
	s.Dy += dy * 10
	for _, button := range s.buttonList {
		button.Dy = dy * 10
		button.Update()
	}
	return nil
}

func (s *Server) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.GeoM.Scale(0.6667, 0.6667)
	screen.DrawImage(resources.BgImage, op)
	s.button1.Draw(screen)
	s.button2.Draw(screen)
	s.button3.Draw(screen)
	for _, button := range s.buttonList {
		button.Draw(screen)
	}
}
