package ImageTools

import "image/color"

type ColorReader interface {
	At(x, y int) color.Color
}
