package vcode

import (
	"image"
	"image/color"
)


type PuzzleMask struct {
	img image.Image
	flag bool // true(default): color puzzle image, false: color base image
}

func NewPMask(puzzle image.Image) *PuzzleMask {
	return &PuzzleMask{img: puzzle, flag: true}
}

func (p *PuzzleMask)Switch() {
	p.flag = !p.flag
}

func (p *PuzzleMask)ColorModel() color.Model {
	return color.AlphaModel
}

func (p *PuzzleMask)Bounds() image.Rectangle {
	return p.img.Bounds()
}

func (p *PuzzleMask)At(x, y int) color.Color {
	r, _, _, _ := p.img.At(x, y).RGBA()
	if (r == 0 && p.flag) || (r != 0 &&  ! p.flag) {
		return color.Opaque
	}

	return color.Transparent
}
