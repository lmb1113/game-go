package ui

import (
	"game/global"
	"game/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"unicode/utf8"
)

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)
const (
	FaceA24 = iota
	FaceA60
	FaceB24
)

func Text(dst *ebiten.Image, content string, faceType int, align int) {
	var face font.Face
	var faceLen = 12
	switch faceType {
	case FaceA24:
		face = resources.GameFont24
		faceLen = 24
	case FaceA60:
		face = resources.GameFont60
		faceLen = 60
	case FaceB24:
		face = resources.GameFontB24
		faceLen = 24
	default:
		return
	}
	var x, y int
	if align == AlignCenter {
		textLen := utf8.RuneCountInString(content)
		x, y = global.ScreenWidth/2-faceLen*textLen/2, global.ScreenHeight/2
	}
	if align == AlignLeft {
		x, y = 0, global.ScreenHeight/2
	}
	if align == AlignRight {
		x, y = global.ScreenWidth, global.ScreenHeight/2
	}
	text.Draw(dst, content, face, x, y, color.Black)
}
