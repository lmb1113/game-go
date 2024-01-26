package ui

import (
	"bytes"
	"game/clinet"
	"game/global"
	"game/resources"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

const VScrollBarWidth = 16

type VScrollBar struct {
	X      int
	Y      int
	Height int

	thumbRate           float64
	thumbOffset         int
	dragging            bool
	draggingStartOffset int
	draggingStartY      int
	contentOffset       int
}

func (v *VScrollBar) thumbSize() int {
	const minThumbSize = VScrollBarWidth

	r := v.thumbRate
	if r > 1 {
		r = 1
	}
	s := int(float64(v.Height) * r)
	if s < minThumbSize {
		return minThumbSize
	}
	return s
}

func (v *VScrollBar) thumbRect() image.Rectangle {
	if v.thumbRate >= 1 {
		return image.Rectangle{}
	}

	s := v.thumbSize()
	return image.Rect(v.X, v.Y+v.thumbOffset, v.X+VScrollBarWidth, v.Y+v.thumbOffset+s)
}

func (v *VScrollBar) maxThumbOffset() int {
	return v.Height - v.thumbSize()
}

func (v *VScrollBar) ContentOffset() int {
	return v.contentOffset
}

func (v *VScrollBar) Update(contentHeight int) {
	v.thumbRate = float64(v.Height) / float64(contentHeight)

	if !v.dragging && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tr := v.thumbRect()
		if tr.Min.X <= x && x < tr.Max.X && tr.Min.Y <= y && y < tr.Max.Y {
			v.dragging = true
			v.draggingStartOffset = v.thumbOffset
			v.draggingStartY = y
		}
	}
	if v.dragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			_, y := ebiten.CursorPosition()
			v.thumbOffset = v.draggingStartOffset + (y - v.draggingStartY)
			if v.thumbOffset < 0 {
				v.thumbOffset = 0
			}
			if v.thumbOffset > v.maxThumbOffset() {
				v.thumbOffset = v.maxThumbOffset()
			}
		} else {
			v.dragging = false
		}
	}

	v.contentOffset = 0
	if v.thumbRate < 1 {
		v.contentOffset = int(float64(contentHeight) * float64(v.thumbOffset) / float64(v.Height))
	}
}

func (v *VScrollBar) Draw(dst *ebiten.Image) {
	sd := image.Rect(v.X, v.Y, v.X+VScrollBarWidth, v.Y+v.Height)
	drawNinePatches(dst, sd, imageSrcRects[imageTypeVScrollBarBack])

	if v.thumbRate < 1 {
		drawNinePatches(dst, v.thumbRect(), imageSrcRects[imageTypeVScrollBarFront])
	}
}

const (
	textBoxPaddingLeft = 8
	textBoxPaddingTop  = 4
)

type TextBox struct {
	Rect       image.Rectangle
	Text       string
	vScrollBar *VScrollBar
	offsetX    int
	offsetY    int
	Button     *Button
}

func (t *TextBox) AppendLine(line string) {
	if t.Text == "" {
		t.Text = line
	} else {
		t.Text += "\n" + line
	}
}

func (t *TextBox) Update() {
	if t.vScrollBar == nil {
		t.vScrollBar = &VScrollBar{}
	}
	t.vScrollBar.X = t.Rect.Max.X - VScrollBarWidth
	t.vScrollBar.Y = t.Rect.Min.Y
	t.vScrollBar.Height = t.Rect.Dy()

	_, h := t.contentSize()
	t.vScrollBar.Update(h)

	t.offsetX = 0
	t.offsetY = t.vScrollBar.ContentOffset()
}

func (t *TextBox) contentSize() (int, int) {
	h := len(strings.Split(t.Text, "\n"))*lineSpacingInPixels + textBoxPaddingTop
	return t.Rect.Dx(), h
}

func (t *TextBox) viewSize() (int, int) {
	return t.Rect.Dx() - VScrollBarWidth - textBoxPaddingLeft, t.Rect.Dy()
}

func (t *TextBox) contentOffset() (int, int) {
	return t.offsetX, t.offsetY
}

func (t *TextBox) Draw(dst *ebiten.Image) {
	drawNinePatches(dst, t.Rect, imageSrcRects[imageTypeTextBox])

	textOp := &text.DrawOptions{}
	x := -float64(t.offsetX) + textBoxPaddingLeft
	y := -float64(t.offsetY) + textBoxPaddingTop
	textOp.GeoM.Translate(x, y)
	textOp.GeoM.Translate(float64(t.Rect.Min.X), float64(t.Rect.Min.Y))
	textOp.ColorScale.ScaleWithColor(color.Black)
	textOp.LineSpacing = lineSpacingInPixels
	text.Draw(dst.SubImage(t.Rect).(*ebiten.Image), t.Text, &text.GoTextFace{
		Source: uiFaceSource,
		Size:   uiFontSize,
	}, textOp)
	t.vScrollBar.Draw(dst)
}

type Server struct {
	button1        *Button
	button2        *Button
	button3        *Button
	textBoxLog     *TextBox
	buttonList     []*Button
	buttonListSize int
	Dy             float64
	Back           func()
	Join           func(roomId uint64, isA bool)
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
	screen.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	s.button1.Draw(screen)
	s.button2.Draw(screen)
	s.button3.Draw(screen)
	for _, button := range s.buttonList {
		button.Draw(screen)
	}
}
