package draw

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Context provides simple methods for drawing over an image.
type Context struct {
	rgba        *image.RGBA
	penColor    color.Color
	fillColor   color.Color
	textColor   color.Color
	fontDrawer  *font.Drawer
	font        *truetype.Font
	fontOptions *truetype.Options
}

// NewContext creates a new context for drawing over image.
func NewContext(rgba *image.RGBA) *Context {
	return &Context{
		rgba:      rgba,
		penColor:  color.Black,
		fillColor: color.Transparent,
		textColor: color.Black,
		fontDrawer: &font.Drawer{
			Dst:  rgba,
			Src:  image.NewUniform(color.Black),
			Face: basicfont.Face7x13,
		},
	}
}

// SetPen changes the pen color (outline color).
func (c *Context) SetPen(clr color.Color) {
	c.penColor = clr
}

// SetFill changes the fill color.
func (c *Context) SetFill(clr color.Color) {
	c.fillColor = clr
}

// SetFontFace changes the font face and font options.
func (c *Context) SetFontFace(font *truetype.Font, options *truetype.Options) {
	c.font = font
	*c.fontOptions = *options
	c.fontDrawer.Face = truetype.NewFace(font, c.fontOptions)
}

// SetFontSize changes the font size only.
func (c *Context) SetFontSize(size float64) {
	c.fontOptions.Size = size
	c.SetFontFace(c.font, c.fontOptions)
}

// SetTextColor changes the text color.
func (c *Context) SetTextColor(clr color.Color) {
	c.textColor = clr
}

// Dot draw a single dot at x,y coordinates.
func (c *Context) Dot(x, y int) {
	c.rgba.Set(x, y, c.penColor)
}

// Dots draws a sequence of dots.
func (c *Context) Dots(points []image.Point) {
	for _, point := range points {
		c.Dot(point.X, point.Y)
	}
}

// Line draws an approximation of a straight line between two points using Bresenham's algorithm.
func (c *Context) Line(x0, y0, x1, y1 int) {
	swap0and1 := false
	swapXandY := math.Abs(float64(y1-y0)) >= math.Abs(float64(x1-x0))
	if swapXandY && y0 > y1 {
		swap0and1 = true
	} else if !swapXandY && x0 > x1 {
		swap0and1 = true
	}

	if swap0and1 {
		x0, x1 = x1, x0
	}

	if swapXandY {
		x0, y0, x1, y1 = y0, x0, y1, x1
	}

	dx := x1 - x0
	dy := y1 - y0

	yi := 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}

	D := 2*dy - dx
	y := y0

	for x := x0; x < x1; x++ {
		if swapXandY {
			c.Dot(y, x)
		} else {
			c.Dot(x, y)
		}

		if D > 0 {
			y = y + yi
			D = D - 2*dx
		}

		D = D + 2*dy
	}

}

// Rect draws a rectangle with pen's color.
func (c *Context) Rect(x0, y0, x1, y1 int) {
	if c.penColor != color.Transparent {
		c.Line(x0, y0, x1, y0)
		c.Line(x1, y0, x1, y1)
		c.Line(x1, y1, x0, y1)
		c.Line(x0, y0, x0, y1)
	}
	if c.fillColor != color.Transparent {
		draw.Draw(c.rgba, image.Rect(x0, y0, x1, y0), &image.Uniform{c.fillColor}, image.ZP, draw.Src)
	}
}

// Cross draws a cross centered at x,y.
func (c *Context) Cross(x, y, size int) {
	c.Line(x, y-size, x, y+size)
	c.Line(x-size, y, x+size, y)
}

// Path draws a sequence of points, connected by lines.
func (c *Context) Path(points []image.Point) {
	var last image.Point
	for i, point := range points {
		if i > 0 {
			c.Line(last.X, last.Y, point.X, point.Y)
		}
		last = point
	}
}

// Parabola draws a parabola arc with the specified coefficients a, b and c.
func (c *Context) Parabola(a1, b1, c1 float64) {
	for x := c.rgba.Bounds().Min.X; x < c.rgba.Bounds().Max.X; x++ {
		y := int(a1*math.Pow(float64(x), 2) + b1*float64(x) + c1 + 0.5)
		if image.Rect(x, y, x, y).In(c.rgba.Bounds()) {
			c.Dot(x, y)
		}
	}
}

// Text draws the given text at x,y with the font chosen in context.
// The default font is golang.org/x/image/font/basicfont.
func (c *Context) Text(x, y int, text string) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	c.fontDrawer.Src = image.NewUniform(c.textColor)
	c.fontDrawer.Dot = point
	c.fontDrawer.DrawString(text)
}
