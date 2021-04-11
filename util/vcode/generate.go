package vcode

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"labsystem/configs"
	"labsystem/util"
	"os"
)

type VImageTyp int
const (
	VImageWidth  = 320
	VImageHeight = 180
	PuzzleWidth  = 50
	PuzzleHeight = 50
	VImagePad    = 20
)

const (
	_ VImageTyp = iota
	Puzzle
	Image
)

var (
	puzzle01 = image.Point{X: 0, Y: 0}
	puzzle02 = image.Point{X: PuzzleWidth, Y: 0}
	puzzle03 = image.Point{X: 0, Y: PuzzleHeight}
	puzzle04 = image.Point{X: PuzzleWidth, Y: PuzzleHeight}

	VImagePath    = configs.CurProjectPath() + "/static/"
	VSrcImagePath = VImagePath + "src/"
	PuzzleSuffix  = "puzzles"
)

var ErrInvalidVCode = errors.New("invalid verify code")

func getPuzzlePos(no int) *image.Point {
	switch no {
	case 2:
		return &puzzle02
	case 3:
		return &puzzle03
	case 4:
		return &puzzle04
	default:
		return &puzzle01
	}
}

func CheckVCode(vcode int) bool {
	if vcode < PuzzleWidth || vcode > VImageWidth - VImagePad* 2 - PuzzleWidth {
		return false
	}

	return true
}

func GetVCode() int {
	return PuzzleWidth + (util.RandIntN(VImageWidth - VImagePad* 2 - 2 * PuzzleWidth) / 10) * 10
}

func GenerateRandomVImage(puzzleNo ,vcode int, imgFilename string) (puz io.ReadWriter, base io.ReadWriter, err error) {
	if !CheckVCode(vcode) {
		return nil, nil, ErrInvalidVCode
	}
	var bgImgF, puzAllImgF *os.File
	// get puzzle image
	puzPos := getPuzzlePos(puzzleNo)
	puzAllImgF, err = os.Open(VSrcImagePath + PuzzleSuffix + ".png")
	if err != nil {
		return nil, nil, err
	}
	puzModelImg := image.NewAlpha(image.Rect(0, 0, VImageWidth, VImageHeight))
	// background render white
	for i := 0;i < 300;i++ {
		for j := 0;j < 400;j++ {
			puzModelImg.Set(j, i, color.White)
		}
	}
	puzAllModelImg, err := png.Decode(puzAllImgF)
	if err != nil {
		return nil, nil, err
	}
	draw.Draw(puzModelImg, image.Rect(VImagePad + vcode, (VImageHeight - PuzzleHeight) / 2, VImagePad+ vcode + PuzzleWidth, (VImageHeight - PuzzleHeight) / 2 + PuzzleHeight), puzAllModelImg, *puzPos, draw.Src)
	// get base image
	bgImgF, err = os.Open(imgFilename)
	if err != nil {
		return nil, nil, err
	}
	bgImg, err := jpeg.Decode(bgImgF)
	if err != nil {
		return nil, nil, err
	}
	// get puzzle mask
	mask := NewPMask(puzModelImg)
	puzImg, baseImg := image.NewRGBA(image.Rect(0, 0, VImageWidth, VImageHeight)), image.NewRGBA(image.Rect(0, 0, VImageWidth, VImageHeight))
	draw.DrawMask(puzImg, image.Rect(-1 * vcode, 0, VImageWidth - vcode, VImageHeight), bgImg, image.Point{}, mask, image.Point{}, draw.Src)
	mask.Switch()
	draw.DrawMask(baseImg, baseImg.Bounds(), bgImg, image.Point{}, mask, image.Point{}, draw.Over)
	var puzBuf, baseBuf []byte
	puz, base = bytes.NewBuffer(puzBuf), bytes.NewBuffer(baseBuf)
	err = png.Encode(puz, puzImg)
	if err != nil {
		return nil, nil, err
	}
	err = png.Encode(base, baseImg)
	if err != nil {
		return nil, nil, err
	}

	return
}
