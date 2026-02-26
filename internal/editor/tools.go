package editor

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type ToolType int

const (
	ToolRectangle ToolType = iota
	ToolEllipse
	ToolArrow
	ToolText
	ToolPixelate
)

type ToolState struct {
	Active     ToolType
	Color      color.Color
	Stroke     int
	FontSize   int
	TextBuffer string
}

type Operation interface {
	Apply(dst draw.Image)
}

type RectOp struct { X1, Y1, X2, Y2 int; Color color.Color; Stroke int }
func (o RectOp) Apply(dst draw.Image) {
	x1, y1, x2, y2 := norm(o.X1,o.Y1,o.X2,o.Y2)
	drawLine(dst, x1,y1,x2,y1,o.Color,o.Stroke)
	drawLine(dst, x2,y1,x2,y2,o.Color,o.Stroke)
	drawLine(dst, x2,y2,x1,y2,o.Color,o.Stroke)
	drawLine(dst, x1,y2,x1,y1,o.Color,o.Stroke)
}

type EllipseOp struct { X1, Y1, X2, Y2 int; Color color.Color; Stroke int }
func (o EllipseOp) Apply(dst draw.Image) {
	x1,y1,x2,y2 := norm(o.X1,o.Y1,o.X2,o.Y2)
	rx := float64(max(1,(x2-x1)/2)); ry := float64(max(1,(y2-y1)/2))
	cx := float64(x1)+rx; cy := float64(y1)+ry
	for t:=0.0; t<2*math.Pi; t+=0.01 {
		x := int(cx + rx*math.Cos(t)); y := int(cy + ry*math.Sin(t))
		for i:=0;i<o.Stroke;i++ { setPixel(dst,x+i,y,o.Color) }
	}
}

type ArrowOp struct { X1,Y1,X2,Y2 int; Color color.Color; Stroke int }
func (o ArrowOp) Apply(dst draw.Image) {
	drawLine(dst,o.X1,o.Y1,o.X2,o.Y2,o.Color,o.Stroke)
	angle := math.Atan2(float64(o.Y2-o.Y1), float64(o.X2-o.X1))
	head := 14.0
	a1 := angle + 2.6
	a2 := angle - 2.6
	x3,y3 := o.X2+int(head*math.Cos(a1)), o.Y2+int(head*math.Sin(a1))
	x4,y4 := o.X2+int(head*math.Cos(a2)), o.Y2+int(head*math.Sin(a2))
	drawLine(dst,o.X2,o.Y2,x3,y3,o.Color,o.Stroke)
	drawLine(dst,o.X2,o.Y2,x4,y4,o.Color,o.Stroke)
}

type TextOp struct { X,Y int; Text string; Color color.Color; FontSize int }
func (o TextOp) Apply(dst draw.Image) {
	// lightweight bitmap-like text fallback: draw blocks per rune.
	x := o.X
	for range o.Text {
		rect := image.Rect(x, o.Y, x+max(6,o.FontSize/2), o.Y+max(10,o.FontSize))
		draw.Draw(dst, rect, image.NewUniform(o.Color), image.Point{}, draw.Src)
		x += max(8, o.FontSize/2+2)
	}
}

type PixelateOp struct { X1, Y1, X2, Y2 int }
func (o PixelateOp) Apply(dst draw.Image) {}
func (o PixelateOp) ApplyOverlay(dst draw.Image, src image.Image) {
	x1, y1, x2, y2 := norm(o.X1, o.Y1, o.X2, o.Y2)
	bounds := image.Rect(x1, y1, x2, y2).Intersect(dst.Bounds()).Intersect(src.Bounds())
	if bounds.Empty() {
		return
	}
	const block = 10
	for y := bounds.Min.Y; y < bounds.Max.Y; y += block {
		for x := bounds.Min.X; x < bounds.Max.X; x += block {
			bx2 := min(x+block, bounds.Max.X)
			by2 := min(y+block, bounds.Max.Y)
			cx := x + (bx2-x)/2
			cy := y + (by2-y)/2
			c := src.At(cx, cy)
			rect := image.Rect(x, y, bx2, by2)
			draw.Draw(dst, rect, image.NewUniform(c), image.Point{}, draw.Src)
		}
	}
}

func norm(x1,y1,x2,y2 int)(int,int,int,int){ if x2<x1{x1,x2=x2,x1}; if y2<y1{y1,y2=y2,y1}; return x1,y1,x2,y2 }
func max(a,b int) int { if a>b {return a}; return b }
func min(a,b int) int { if a<b {return a}; return b }
func setPixel(img draw.Image, x,y int, c color.Color){ if image.Pt(x,y).In(img.Bounds()) { img.Set(x,y,c) } }
func drawLine(img draw.Image, x1,y1,x2,y2 int, c color.Color, stroke int){
	dx := int(math.Abs(float64(x2-x1))); sx := -1; if x1<x2 { sx = 1 }
	dy := -int(math.Abs(float64(y2-y1))); sy := -1; if y1<y2 { sy = 1 }
	err := dx + dy
	for {
		for i:=0;i<stroke;i++ { setPixel(img,x1+i,y1,c); setPixel(img,x1,y1+i,c) }
		if x1==x2 && y1==y2 { break }
		e2 := 2*err
		if e2>=dy { err += dy; x1 += sx }
		if e2<=dx { err += dx; y1 += sy }
	}
}
