package Display

import (
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"sort"
)

type XTermColor struct {
	xtermCode uint16
	hexColor  uint32
}

var BaseLUT = []XTermColor{
	// Primary 3-bit (8 colors). Unique representation!
	{xtermCode: 0, hexColor: 0x000000},
	{xtermCode: 1, hexColor: 0x800000},
	{xtermCode: 2, hexColor: 0x008000},
	{xtermCode: 3, hexColor: 0x808000},
	{xtermCode: 4, hexColor: 0x000080},
	{xtermCode: 5, hexColor: 0x800080},
	{xtermCode: 6, hexColor: 0x008080},
	{xtermCode: 7, hexColor: 0xc0c0c0},

	// Equivalent "bright" versions of original 8 colors.
	{xtermCode: 8, hexColor: 0x808080},
	{xtermCode: 9, hexColor: 0xff0000},
	{xtermCode: 10, hexColor: 0x00ff00},
	{xtermCode: 11, hexColor: 0xffff00},
	{xtermCode: 12, hexColor: 0x0000ff},
	{xtermCode: 13, hexColor: 0xff00ff},
	{xtermCode: 14, hexColor: 0x00ffff},
	{xtermCode: 15, hexColor: 0xffffff},

	// Strictly ascending.
	{xtermCode: 16, hexColor: 0x000000},
	{xtermCode: 17, hexColor: 0x00005f},
	{xtermCode: 18, hexColor: 0x000087},
	{xtermCode: 19, hexColor: 0x0000af},
	{xtermCode: 20, hexColor: 0x0000d7},
	{xtermCode: 21, hexColor: 0x0000ff},
	{xtermCode: 22, hexColor: 0x005f00},
	{xtermCode: 23, hexColor: 0x005f5f},
	{xtermCode: 24, hexColor: 0x005f87},
	{xtermCode: 25, hexColor: 0x005faf},
	{xtermCode: 26, hexColor: 0x005fd7},
	{xtermCode: 27, hexColor: 0x005fff},
	{xtermCode: 28, hexColor: 0x008700},
	{xtermCode: 29, hexColor: 0x00875f},
	{xtermCode: 30, hexColor: 0x008787},
	{xtermCode: 31, hexColor: 0x0087af},
	{xtermCode: 32, hexColor: 0x0087d7},
	{xtermCode: 33, hexColor: 0x0087ff},
	{xtermCode: 34, hexColor: 0x00af00},
	{xtermCode: 35, hexColor: 0x00af5f},
	{xtermCode: 36, hexColor: 0x00af87},
	{xtermCode: 37, hexColor: 0x00afaf},
	{xtermCode: 38, hexColor: 0x00afd7},
	{xtermCode: 39, hexColor: 0x00afff},
	{xtermCode: 40, hexColor: 0x00d700},
	{xtermCode: 41, hexColor: 0x00d75f},
	{xtermCode: 42, hexColor: 0x00d787},
	{xtermCode: 43, hexColor: 0x00d7af},
	{xtermCode: 44, hexColor: 0x00d7d7},
	{xtermCode: 45, hexColor: 0x00d7ff},
	{xtermCode: 46, hexColor: 0x00ff00},
	{xtermCode: 47, hexColor: 0x00ff5f},
	{xtermCode: 48, hexColor: 0x00ff87},
	{xtermCode: 49, hexColor: 0x00ffaf},
	{xtermCode: 50, hexColor: 0x00ffd7},
	{xtermCode: 51, hexColor: 0x00ffff},
	{xtermCode: 52, hexColor: 0x5f0000},
	{xtermCode: 53, hexColor: 0x5f005f},
	{xtermCode: 54, hexColor: 0x5f0087},
	{xtermCode: 55, hexColor: 0x5f00af},
	{xtermCode: 56, hexColor: 0x5f00d7},
	{xtermCode: 57, hexColor: 0x5f00ff},
	{xtermCode: 58, hexColor: 0x5f5f00},
	{xtermCode: 59, hexColor: 0x5f5f5f},
	{xtermCode: 60, hexColor: 0x5f5f87},
	{xtermCode: 61, hexColor: 0x5f5faf},
	{xtermCode: 62, hexColor: 0x5f5fd7},
	{xtermCode: 63, hexColor: 0x5f5fff},
	{xtermCode: 64, hexColor: 0x5f8700},
	{xtermCode: 65, hexColor: 0x5f875f},
	{xtermCode: 66, hexColor: 0x5f8787},
	{xtermCode: 67, hexColor: 0x5f87af},
	{xtermCode: 68, hexColor: 0x5f87d7},
	{xtermCode: 69, hexColor: 0x5f87ff},
	{xtermCode: 70, hexColor: 0x5faf00},
	{xtermCode: 71, hexColor: 0x5faf5f},
	{xtermCode: 72, hexColor: 0x5faf87},
	{xtermCode: 73, hexColor: 0x5fafaf},
	{xtermCode: 74, hexColor: 0x5fafd7},
	{xtermCode: 75, hexColor: 0x5fafff},
	{xtermCode: 76, hexColor: 0x5fd700},
	{xtermCode: 77, hexColor: 0x5fd75f},
	{xtermCode: 78, hexColor: 0x5fd787},
	{xtermCode: 79, hexColor: 0x5fd7af},
	{xtermCode: 80, hexColor: 0x5fd7d7},
	{xtermCode: 81, hexColor: 0x5fd7ff},
	{xtermCode: 82, hexColor: 0x5fff00},
	{xtermCode: 83, hexColor: 0x5fff5f},
	{xtermCode: 84, hexColor: 0x5fff87},
	{xtermCode: 85, hexColor: 0x5fffaf},
	{xtermCode: 86, hexColor: 0x5fffd7},
	{xtermCode: 87, hexColor: 0x5fffff},
	{xtermCode: 88, hexColor: 0x870000},
	{xtermCode: 89, hexColor: 0x87005f},
	{xtermCode: 90, hexColor: 0x870087},
	{xtermCode: 91, hexColor: 0x8700af},
	{xtermCode: 92, hexColor: 0x8700d7},
	{xtermCode: 93, hexColor: 0x8700ff},
	{xtermCode: 94, hexColor: 0x875f00},
	{xtermCode: 95, hexColor: 0x875f5f},
	{xtermCode: 96, hexColor: 0x875f87},
	{xtermCode: 97, hexColor: 0x875faf},
	{xtermCode: 98, hexColor: 0x875fd7},
	{xtermCode: 99, hexColor: 0x875fff},
	{xtermCode: 100, hexColor: 0x878700},
	{xtermCode: 101, hexColor: 0x87875f},
	{xtermCode: 102, hexColor: 0x878787},
	{xtermCode: 103, hexColor: 0x8787af},
	{xtermCode: 104, hexColor: 0x8787d7},
	{xtermCode: 105, hexColor: 0x8787ff},
	{xtermCode: 106, hexColor: 0x87af00},
	{xtermCode: 107, hexColor: 0x87af5f},
	{xtermCode: 108, hexColor: 0x87af87},
	{xtermCode: 109, hexColor: 0x87afaf},
	{xtermCode: 110, hexColor: 0x87afd7},
	{xtermCode: 111, hexColor: 0x87afff},
	{xtermCode: 112, hexColor: 0x87d700},
	{xtermCode: 113, hexColor: 0x87d75f},
	{xtermCode: 114, hexColor: 0x87d787},
	{xtermCode: 115, hexColor: 0x87d7af},
	{xtermCode: 116, hexColor: 0x87d7d7},
	{xtermCode: 117, hexColor: 0x87d7ff},
	{xtermCode: 118, hexColor: 0x87ff00},
	{xtermCode: 119, hexColor: 0x87ff5f},
	{xtermCode: 120, hexColor: 0x87ff87},
	{xtermCode: 121, hexColor: 0x87ffaf},
	{xtermCode: 122, hexColor: 0x87ffd7},
	{xtermCode: 123, hexColor: 0x87ffff},
	{xtermCode: 124, hexColor: 0xaf0000},
	{xtermCode: 125, hexColor: 0xaf005f},
	{xtermCode: 126, hexColor: 0xaf0087},
	{xtermCode: 127, hexColor: 0xaf00af},
	{xtermCode: 128, hexColor: 0xaf00d7},
	{xtermCode: 129, hexColor: 0xaf00ff},
	{xtermCode: 130, hexColor: 0xaf5f00},
	{xtermCode: 131, hexColor: 0xaf5f5f},
	{xtermCode: 132, hexColor: 0xaf5f87},
	{xtermCode: 133, hexColor: 0xaf5faf},
	{xtermCode: 134, hexColor: 0xaf5fd7},
	{xtermCode: 135, hexColor: 0xaf5fff},
	{xtermCode: 136, hexColor: 0xaf8700},
	{xtermCode: 137, hexColor: 0xaf875f},
	{xtermCode: 138, hexColor: 0xaf8787},
	{xtermCode: 139, hexColor: 0xaf87af},
	{xtermCode: 140, hexColor: 0xaf87d7},
	{xtermCode: 141, hexColor: 0xaf87ff},
	{xtermCode: 142, hexColor: 0xafaf00},
	{xtermCode: 143, hexColor: 0xafaf5f},
	{xtermCode: 144, hexColor: 0xafaf87},
	{xtermCode: 145, hexColor: 0xafafaf},
	{xtermCode: 146, hexColor: 0xafafd7},
	{xtermCode: 147, hexColor: 0xafafff},
	{xtermCode: 148, hexColor: 0xafd700},
	{xtermCode: 149, hexColor: 0xafd75f},
	{xtermCode: 150, hexColor: 0xafd787},
	{xtermCode: 151, hexColor: 0xafd7af},
	{xtermCode: 152, hexColor: 0xafd7d7},
	{xtermCode: 153, hexColor: 0xafd7ff},
	{xtermCode: 154, hexColor: 0xafff00},
	{xtermCode: 155, hexColor: 0xafff5f},
	{xtermCode: 156, hexColor: 0xafff87},
	{xtermCode: 157, hexColor: 0xafffaf},
	{xtermCode: 158, hexColor: 0xafffd7},
	{xtermCode: 159, hexColor: 0xafffff},
	{xtermCode: 160, hexColor: 0xd70000},
	{xtermCode: 161, hexColor: 0xd7005f},
	{xtermCode: 162, hexColor: 0xd70087},
	{xtermCode: 163, hexColor: 0xd700af},
	{xtermCode: 164, hexColor: 0xd700d7},
	{xtermCode: 165, hexColor: 0xd700ff},
	{xtermCode: 166, hexColor: 0xd75f00},
	{xtermCode: 167, hexColor: 0xd75f5f},
	{xtermCode: 168, hexColor: 0xd75f87},
	{xtermCode: 169, hexColor: 0xd75faf},
	{xtermCode: 170, hexColor: 0xd75fd7},
	{xtermCode: 171, hexColor: 0xd75fff},
	{xtermCode: 172, hexColor: 0xd78700},
	{xtermCode: 173, hexColor: 0xd7875f},
	{xtermCode: 174, hexColor: 0xd78787},
	{xtermCode: 175, hexColor: 0xd787af},
	{xtermCode: 176, hexColor: 0xd787d7},
	{xtermCode: 177, hexColor: 0xd787ff},
	{xtermCode: 178, hexColor: 0xd7af00},
	{xtermCode: 179, hexColor: 0xd7af5f},
	{xtermCode: 180, hexColor: 0xd7af87},
	{xtermCode: 181, hexColor: 0xd7afaf},
	{xtermCode: 182, hexColor: 0xd7afd7},
	{xtermCode: 183, hexColor: 0xd7afff},
	{xtermCode: 184, hexColor: 0xd7d700},
	{xtermCode: 185, hexColor: 0xd7d75f},
	{xtermCode: 186, hexColor: 0xd7d787},
	{xtermCode: 187, hexColor: 0xd7d7af},
	{xtermCode: 188, hexColor: 0xd7d7d7},
	{xtermCode: 189, hexColor: 0xd7d7ff},
	{xtermCode: 190, hexColor: 0xd7ff00},
	{xtermCode: 191, hexColor: 0xd7ff5f},
	{xtermCode: 192, hexColor: 0xd7ff87},
	{xtermCode: 193, hexColor: 0xd7ffaf},
	{xtermCode: 194, hexColor: 0xd7ffd7},
	{xtermCode: 195, hexColor: 0xd7ffff},
	{xtermCode: 196, hexColor: 0xff0000},
	{xtermCode: 197, hexColor: 0xff005f},
	{xtermCode: 198, hexColor: 0xff0087},
	{xtermCode: 199, hexColor: 0xff00af},
	{xtermCode: 200, hexColor: 0xff00d7},
	{xtermCode: 201, hexColor: 0xff00ff},
	{xtermCode: 202, hexColor: 0xff5f00},
	{xtermCode: 203, hexColor: 0xff5f5f},
	{xtermCode: 204, hexColor: 0xff5f87},
	{xtermCode: 205, hexColor: 0xff5faf},
	{xtermCode: 206, hexColor: 0xff5fd7},
	{xtermCode: 207, hexColor: 0xff5fff},
	{xtermCode: 208, hexColor: 0xff8700},
	{xtermCode: 209, hexColor: 0xff875f},
	{xtermCode: 210, hexColor: 0xff8787},
	{xtermCode: 211, hexColor: 0xff87af},
	{xtermCode: 212, hexColor: 0xff87d7},
	{xtermCode: 213, hexColor: 0xff87ff},
	{xtermCode: 214, hexColor: 0xffaf00},
	{xtermCode: 215, hexColor: 0xffaf5f},
	{xtermCode: 216, hexColor: 0xffaf87},
	{xtermCode: 217, hexColor: 0xffafaf},
	{xtermCode: 218, hexColor: 0xffafd7},
	{xtermCode: 219, hexColor: 0xffafff},
	{xtermCode: 220, hexColor: 0xffd700},
	{xtermCode: 221, hexColor: 0xffd75f},
	{xtermCode: 222, hexColor: 0xffd787},
	{xtermCode: 223, hexColor: 0xffd7af},
	{xtermCode: 224, hexColor: 0xffd7d7},
	{xtermCode: 225, hexColor: 0xffd7ff},
	{xtermCode: 226, hexColor: 0xffff00},
	{xtermCode: 227, hexColor: 0xffff5f},
	{xtermCode: 228, hexColor: 0xffff87},
	{xtermCode: 229, hexColor: 0xffffaf},
	{xtermCode: 230, hexColor: 0xffffd7},
	{xtermCode: 231, hexColor: 0xffffff},

	// Gray-scale range.
	{xtermCode: 232, hexColor: 0x080808},
	{xtermCode: 233, hexColor: 0x121212},
	{xtermCode: 234, hexColor: 0x1c1c1c},
	{xtermCode: 235, hexColor: 0x262626},
	{xtermCode: 236, hexColor: 0x303030},
	{xtermCode: 237, hexColor: 0x3a3a3a},
	{xtermCode: 238, hexColor: 0x444444},
	{xtermCode: 239, hexColor: 0x4e4e4e},
	{xtermCode: 240, hexColor: 0x585858},
	{xtermCode: 241, hexColor: 0x626262},
	{xtermCode: 242, hexColor: 0x6c6c6c},
	{xtermCode: 243, hexColor: 0x767676},
	{xtermCode: 244, hexColor: 0x808080},
	{xtermCode: 245, hexColor: 0x8a8a8a},
	{xtermCode: 246, hexColor: 0x949494},
	{xtermCode: 247, hexColor: 0x9e9e9e},
	{xtermCode: 248, hexColor: 0xa8a8a8},
	{xtermCode: 249, hexColor: 0xb2b2b2},
	{xtermCode: 250, hexColor: 0xbcbcbc},
	{xtermCode: 251, hexColor: 0xc6c6c6},
	{xtermCode: 252, hexColor: 0xd0d0d0},
	{xtermCode: 253, hexColor: 0xdadada},
	{xtermCode: 254, hexColor: 0xe4e4e4},
	{xtermCode: 255, hexColor: 0xeeeeee},
}

var hexToPos map[uint32]int
var hexValues []uint32

var HVals []uint32

var HValToHex map[uint32]uint32

func hsv2rgb(H float32, S float32, V float32) (float32, float32, float32) {
	var r, g, b float32
	if S == 0 {
		r = V
		g = V
		b = V
	} else {
		h := H * 6
		if h == 6 {
			h = 0
		}
		i := float32(math.Floor(float64(h)))
		v1 := V * (1 - S)
		v2 := V * (1 - S*(h-i))
		v3 := V * (1 - S*(1-(h-i)))

		if i == 0 {
			r = V
			g = v3
			b = v1
		} else if i == 1 {
			r = v2
			g = V
			b = v1
		} else if i == 2 {
			r = v1
			g = V
			b = v3
		} else if i == 3 {
			r = v1
			g = v2
			b = V
		} else if i == 4 {
			r = v3
			g = v1
			b = V
		} else {
			r = V
			g = v1
			b = v2
		}
	}
	return r, g, b
}

func getApproxHex(v uint32) uint32 {
	var lastValue uint32 = 0
	for i := 0; i < len(hexValues); i++ {
		iv := hexValues[i]
		if (iv > v && lastValue < v) || iv == v {
			// Got!
			return iv
		}
		lastValue = iv
	}

	return hexValues[len(hexValues)-1]
}
func getApproxH(v uint32) uint32 {
	var lastValue uint32 = 0
	for i := 0; i < len(HVals); i++ {
		iv := HVals[i]
		if (iv > v && lastValue < v) || iv == v {
			// Got!
			return iv
		}
		lastValue = iv
	}

	return HVals[len(hexValues)-1]
}

func GetXTermColor(rgb uint32) uint16 {
	v := getApproxHex(rgb)
	pos := hexToPos[v]
	return BaseLUT[pos].xtermCode
}

func GetXTermColorHVal(h uint32) uint16 {
	v := getApproxH(h)
	z := HValToHex[v]
	pos := hexToPos[z]
	return BaseLUT[pos].xtermCode
}

func InitLut() {
	HVals = make([]uint32, 360)
	hexValues = make([]uint32, len(BaseLUT))
	hexToPos = make(map[uint32]int)
	HValToHex = make(map[uint32]uint32)
	for i := 0; i < len(BaseLUT); i++ {
		item := BaseLUT[i]
		hexToPos[item.hexColor] = i
		hexValues[i] = item.hexColor
		c := colorful.Color{
			float64((item.hexColor&0xFF0000)>>16) / 255,
			float64((item.hexColor&0xFF00)>>8) / 255,
			float64(item.hexColor&0xFF) / 255,
		}

		h, _, v := c.Hsv()
		if math.Round(v) == 1 {
			HVals[i] = uint32(h)
			HValToHex[uint32(h)] = item.hexColor
		}
	}

	sort.Slice(hexValues, func(i, j int) bool {
		return hexValues[i] < hexValues[j]
	})

	sort.Slice(HVals, func(i, j int) bool {
		return HVals[i] < HVals[j]
	})
}
