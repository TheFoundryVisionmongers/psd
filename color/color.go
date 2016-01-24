package color

import "image/color"

// Gray1 represents an 1-bit monochrome bitmap color.
type Gray1 struct {
	Y bool
}

// RGBA implements color.Color interface's method.
func (c Gray1) RGBA() (r, g, b, a uint32) {
	if c.Y {
		return 0xffff, 0xffff, 0xffff, 0xffff
	}
	return 0, 0, 0, 0xffff
}

func gray1Model(c color.Color) color.Color {
	if _, ok := c.(Gray1); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	return Gray1{y >= 0x8000}
}

// Gray32 represents a 32-bit float grayscale color.
type Gray32 struct {
	Y float32
}

// RGBA implements color.Color interface's method.
func (c Gray32) RGBA() (r, g, b, a uint32) {
	// FIXME: How should I convert to RGBA?
	var y uint32
	switch {
	case c.Y >= 0x10001:
		y = 0xffffffff
	case c.Y <= 0:
		y = 0
	default:
		y = uint32(c.Y * 0xffff)
	}
	return y, y, y, 0xffff
}

func gray32Model(c color.Color) color.Color {
	if _, ok := c.(Gray32); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	y := (299*r + 587*g + 114*b + 500) / 1000
	return Gray32{float32(y) / 0xffff}
}

type NGrayA struct {
	Y uint8
	A uint8
}

func (c NGrayA) RGBA() (uint32, uint32, uint32, uint32) {
	y := uint32(c.Y) * 0x101
	if c.A == 0xff {
		return y, y, y, 0xffff
	}
	if c.A == 0 {
		return 0, 0, 0, 0
	}
	a := uint32(c.A) * 0x101
	y = y * a / 0xffff
	return y, y, y, a
}

func nGrayAModel(c color.Color) color.Color {
	if _, ok := c.(NGrayA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return NGrayA{0, 0}
	}
	y := (299*r + 587*g + 114*b + 500) / 1000
	if a == 0xffff {
		return NGrayA{uint8(y >> 8), 0xff}
	}
	y = (y * 0xffff) / a
	return NGrayA{uint8(y >> 8), uint8(a >> 8)}
}

type NGrayA32 struct {
	Y uint16
	A uint16
}

func (c NGrayA32) RGBA() (uint32, uint32, uint32, uint32) {
	y := uint32(c.Y)
	if c.A == 0xffff {
		return y, y, y, 0xffff
	}
	if c.A == 0 {
		return 0, 0, 0, 0
	}
	a := uint32(c.A)
	y = y * a / 0xffff
	return y, y, y, a
}

func nGrayA32Model(c color.Color) color.Color {
	if _, ok := c.(NGrayA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return NGrayA32{0, 0}
	}
	y := (299*r + 587*g + 114*b + 500) / 1000
	if a == 0xffff {
		return NGrayA32{uint16(y), 0xffff}
	}
	y = (y * 0xffff) / a
	return NGrayA32{uint16(y), uint16(a)}
}

type NGrayA64 struct {
	Y float32
	A float32
}

func (c NGrayA64) RGBA() (uint32, uint32, uint32, uint32) {
	var y uint32
	switch {
	case c.Y >= 0x10001:
		y = 0xffffffff
	case c.Y <= 0:
		y = 0
	default:
		y = uint32(c.Y * 0xffff)
	}
	switch {
	case c.A >= 1:
		return y, y, y, 0xffff
	case c.A <= 0:
		return 0, 0, 0, 0
	}
	a := uint32(c.A * 0xffff)
	y = y * a / 0xffff
	return y, y, y, a
}

func nGrayA64Model(c color.Color) color.Color {
	if _, ok := c.(NGrayA64); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return NGrayA64{0, 0}
	}
	y := (299*r + 587*g + 114*b + 500) / 1000
	if a == 0xffff {
		return NGrayA64{float32(y) / 0xffff, 1}
	}
	y = (y * 0xffff) / a
	return NGrayA64{float32(y) / 0xffff, float32(a) / 0xffff}
}

// NCMYKA represents a non-alpha-premultiplied CMYK color, having 8 bits for each of cyan,
// magenta, yellow, black and alpha.
// NCMYKA is different from color.CMYK, CMYK is inverted value.
//
// It is not associated with any particular color profile.
type NCMYKA struct {
	C, M, Y, K, A uint8
}

// RGBA implements color.Color interface's method.
func (c NCMYKA) RGBA() (uint32, uint32, uint32, uint32) {
	w := uint32(c.K) * 0x10201
	r := uint32(c.C) * w / 0xffff
	g := uint32(c.M) * w / 0xffff
	b := uint32(c.Y) * w / 0xffff
	if c.A == 0xff {
		return r, g, b, 0xffff
	}
	if c.A == 0 {
		return 0, 0, 0, 0
	}

	a := uint32(c.A) * 0x101
	r = r * a / 0xffff
	g = g * a / 0xffff
	b = b * a / 0xffff
	return r, g, b, a
}

func nCMYKAModel(c color.Color) color.Color {
	if _, ok := c.(NCMYKA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	cc, mm, yy, kk := color.RGBToCMYK(uint8(r>>8), uint8(g>>8), uint8(b>>8))
	cc = uint8((uint32(cc) * 0xffff) / a)
	mm = uint8((uint32(mm) * 0xffff) / a)
	yy = uint8((uint32(yy) * 0xffff) / a)
	kk = uint8((uint32(kk) * 0xffff) / a)
	return NCMYKA{255 - cc, 255 - mm, 255 - yy, 255 - kk, uint8(a >> 8)}
}

type NCMYKA80 struct {
	C, M, Y, K, A uint16
}

// RGBA implements color.Color interface's method.
func (c NCMYKA80) RGBA() (uint32, uint32, uint32, uint32) {
	w := uint32(c.K)
	r := uint32(c.C) * w / 0xffff
	g := uint32(c.M) * w / 0xffff
	b := uint32(c.Y) * w / 0xffff
	if c.A == 0xffff {
		return r, g, b, 0xffff
	}
	if c.A == 0 {
		return 0, 0, 0, 0
	}

	a := uint32(c.A)
	r = r * a / 0xffff
	g = g * a / 0xffff
	b = b * a / 0xffff
	return r, g, b, a
}

func nCMYKA80Model(c color.Color) color.Color {
	if _, ok := c.(NCMYKA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return NCMYKA80{0xffff, 0xffff, 0xffff, 0xffff, 0}
	}
	w := r
	if w < g {
		w = g
	}
	if w < b {
		w = b
	}
	if w == 0 {
		return NCMYKA80{0xffff, 0xffff, 0xffff, 0xffff, uint16(a)}
	}
	cc := (w - r) * 0xffff / w
	mm := (w - g) * 0xffff / w
	yy := (w - b) * 0xffff / w
	kk := 0xffff - w
	if a == 0xffff {
		return NCMYKA80{uint16(0xffff - cc), uint16(0xffff - mm), uint16(0xffff - yy), uint16(0xffff - kk), 0xffff}
	}
	cc = (cc * 0xffff) / a
	mm = (mm * 0xffff) / a
	yy = (yy * 0xffff) / a
	kk = (kk * 0xffff) / a
	return NCMYKA80{uint16(0xffff - cc), uint16(0xffff - mm), uint16(0xffff - yy), uint16(0xffff - kk), uint16(a)}
}

// These are color model.
var (
	Gray1Model    = color.ModelFunc(gray1Model)
	NGrayAModel   = color.ModelFunc(nGrayAModel)
	Gray32Model   = color.ModelFunc(gray32Model)
	NGrayA32Model = color.ModelFunc(nGrayA32Model)
	NGrayA64Model = color.ModelFunc(nGrayA64Model)
	NCMYKAModel   = color.ModelFunc(nCMYKAModel)
	NCMYKA80Model = color.ModelFunc(nCMYKA80Model)
)
