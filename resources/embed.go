package resources

import (
	_ "embed"
)

var (
	//go:embed images/22.png
	Tile_png []byte
	//go:embed images/Angle/dirt_E.png
	Tiles_png []byte
	//go:embed images/2.png
	T2_png []byte
	//go:embed images/3.png
	T3_png []byte
	//go:embed images/5.png
	T5_png []byte
	//go:embed images/6.png
	T6_png []byte
	//go:embed 1.ttf
	Font1 []byte
	//go:embed 2.ttf
	Font2 []byte
	//go:embed 3.ttf
	Font3 []byte
)
