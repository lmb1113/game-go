package resources

import (
	"bytes"
	_ "embed"
	"game/global"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
)

var (
	//go:embed images/bg.png
	Bg_png []byte
	//go:embed images/2.png
	T2_png []byte
	//go:embed images/3.png
	T3_png []byte
	//go:embed images/5.png
	T5_png []byte
	//go:embed images/6.png
	T6_png []byte
	//go:embed images/7.png
	T7_png []byte

	//go:embed images/8.png
	T8_png []byte
	//go:embed images/9.png
	T9_png []byte
	//go:embed 1.ttf
	Font1 []byte
	//go:embed 2.ttf
	Font2 []byte
	//go:embed 3.ttf
	Font3 []byte
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	WhiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
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

func Init() {
	// Decode an image from the image file's byte slice.
	bgImg, _, _ := image.Decode(bytes.NewReader(Bg_png))
	BgImage = ebiten.NewImageFromImage(bgImg)

	img2, _, _ := image.Decode(bytes.NewReader(T2_png))
	AtmImage = ebiten.NewImageFromImage(reverseImage(img2))
	AtmImage2 = ebiten.NewImageFromImage(img2)

	img3, _, _ := image.Decode(bytes.NewReader(T3_png))
	CkxImage = ebiten.NewImageFromImage(img3)
	CkxImage2 = ebiten.NewImageFromImage(reverseImage(img3))

	img5, _, _ := image.Decode(bytes.NewReader(T5_png))
	BasketballImage = ebiten.NewImageFromImage(img5)
	BasketballImage2 = ebiten.NewImageFromImage(reverseImage(img5))

	img6, _, _ := image.Decode(bytes.NewReader(T6_png))
	ChickenImage = ebiten.NewImageFromImage(reverseImage(img6))
	ChickenImage2 = ebiten.NewImageFromImage(img6)

	img7, _, _ := image.Decode(bytes.NewReader(T7_png))
	Skill3A = ebiten.NewImageFromImage(img7)
	Skill3B = ebiten.NewImageFromImage(reverseImage(img7))

	img8, _, _ := image.Decode(bytes.NewReader(T8_png))
	Skill4A = ebiten.NewImageFromImage(img8)
	Skill4B = ebiten.NewImageFromImage(reverseImage(img8))

	img9, _, _ := image.Decode(bytes.NewReader(T9_png))
	Skill5A = ebiten.NewImageFromImage(img9)
	Skill5B = ebiten.NewImageFromImage(reverseImage(img9))
	whiteImage.Fill(color.White)
	tt1, err := opentype.Parse(Font1)
	if err != nil {
		log.Fatal(err)
	}
	tt2, err := opentype.Parse(Font2)
	if err != nil {
		log.Fatal(err)
	}
	GameFont60, err = opentype.NewFace(tt1, &opentype.FaceOptions{
		Size:    60,
		DPI:     global.Dpi,
		Hinting: font.HintingVertical,
	})
	GameFont24, err = opentype.NewFace(tt1, &opentype.FaceOptions{
		Size:    24,
		DPI:     global.Dpi,
		Hinting: font.HintingVertical,
	})
	GameFontB24, err = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
}

var (
	TilesImage       *ebiten.Image
	BgImage          *ebiten.Image
	CkxImage         *ebiten.Image
	CkxImage2        *ebiten.Image
	AtmImage         *ebiten.Image
	AtmImage2        *ebiten.Image
	BasketballImage  *ebiten.Image
	BasketballImage2 *ebiten.Image
	ChickenImage     *ebiten.Image
	ChickenImage2    *ebiten.Image
	Skill3A          *ebiten.Image
	Skill3B          *ebiten.Image
	Skill4A          *ebiten.Image
	Skill4B          *ebiten.Image
	Skill5A          *ebiten.Image
	Skill5B          *ebiten.Image
	GameFont60       font.Face
	GameFont24       font.Face
	GameFontB24      font.Face
)
