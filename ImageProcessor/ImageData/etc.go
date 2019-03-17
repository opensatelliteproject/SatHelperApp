package ImageData

import (
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

func init() {
	fontBytes := MustAsset("FreeMono.ttf")
	loadedFont, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	draw2d.RegisterFont(draw2d.FontData{
		Name:   "FreeMono",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	}, loadedFont)
}
