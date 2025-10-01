package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

// Type Image  since we've defined all image methods so
// we can use it as an image.Image
type Foo struct{}

func (i Foo) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Foo) Bounds() image.Rectangle {
	return image.Rect(0, 0, 255, 255)
}

func (i Foo) At(x, y int) color.Color {
	return color.RGBA{uint8(x), uint8(y), uint8(x + y), 0xff}

}

func main() {
	m := Foo{}
	pic.ShowImage(m)
}
