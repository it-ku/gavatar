package gavatar

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

type AvatarGenerate struct {
	W        int
	H        int
	fontFile string
	fontSize float64
	bg       color.Color
	fg       color.Color
	ctx      *freetype.Context
}

func NewAvatarGenerate(fontFile string) *AvatarGenerate {
	return &AvatarGenerate{
		fontFile: fontFile,
		bg:       color.White,
		fg:       color.Black,
		W:        64,
		H:        64,
		fontSize: 64,
	}
}

// SetBackgroundColor 设置背景颜色
func (ag *AvatarGenerate) SetBackgroundColor(c color.Color) {
	ag.bg = c
}

// SetBackgroundColorHex 设置背景颜色(十六进制)
//
// []uint32{ 0xff6200, 0x42c58e, 0x5a8de1, 0x785fe0 }
func (ag *AvatarGenerate) SetBackgroundColorHex(hex uint32) {
	ag.bg = ag.hexToRGBA(hex)
}

// SetFrontColor 设置字体颜色
func (ag *AvatarGenerate) SetFrontColor(c color.Color) {
	ag.fg = c
}

// SetFrontColorHex 设置字体颜色(十六进制)
func (ag *AvatarGenerate) SetFrontColorHex(hex uint32) {
	ag.fg = ag.hexToRGBA(hex)
}

// SetAvatarSize 设置头像大小，默认64
func (ag *AvatarGenerate) SetAvatarSize(width, height int) {
	ag.W = width
	ag.H = height
}

// SetFontSize 设置字体大小，默认字号64
func (ag *AvatarGenerate) SetFontSize(size float64) {
	ag.fontSize = size
}

func (ag *AvatarGenerate) GenerateImage(s string, outFileName string) error {
	bs, err := ag.GenerateImageContent(s)
	if err != nil {
		return err
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("create file error: %w", err)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	if _, err := b.Write(bs); err != nil {
		return fmt.Errorf("write bytes to file: %w", err)
	}
	if err = b.Flush(); err != nil {
		return fmt.Errorf("flush image: %w", err)
	}

	return nil
}

func (ag *AvatarGenerate) GenerateImageContent(s string) ([]byte, error) {
	rgba := ag.CreateBackgroundImage()
	if ag.ctx == nil {
		if err := ag.DrawContext(rgba); err != nil {
			return nil, err
		}
	}
	x := ag.W / 2 - (ag.GetFontWidth() *3/5)/2
	y := ag.H / 2 + ag.GetFontWidth() *4/11

	pt := freetype.Pt(x, y)
	if _, err := ag.ctx.DrawString(s, pt); err != nil {
		return nil, fmt.Errorf("draw string error: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, rgba); err != nil {
		return nil, fmt.Errorf("png encode error: %w", err)
	}

	return buf.Bytes(), nil

}

func (ag *AvatarGenerate) CreateBackgroundImage() *image.RGBA {
	bg := image.NewUniform(ag.bg)
	rgba := image.NewRGBA(image.Rect(0, 0, ag.W, ag.H))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	return rgba
}

func (ag *AvatarGenerate) DrawContext(rgba *image.RGBA) error {
	fontBytes, err := ioutil.ReadFile(ag.fontFile)
	if err != nil {
		return fmt.Errorf("error when open font file: %w", err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return fmt.Errorf("error when parse font file: %w", err)
	}
	c := freetype.NewContext()
	c.SetDPI(72)                      // 设置像素密度
	c.SetFont(f)                      // 指定字体
	c.SetFontSize(ag.fontSize)        // 指定字体大小
	c.SetClip(rgba.Bounds())          // 指定画布绘制范围
	c.SetDst(rgba)                    // 指定画布对象
	c.SetSrc(image.NewUniform(ag.fg)) // 字体颜色
	c.SetHinting(font.HintingNone)    //

	ag.ctx = c
	return nil
}

func (ag *AvatarGenerate) hexToRGBA(h uint32) *color.RGBA {
	rgba := &color.RGBA{
		R: uint8(h >> 16),
		G: uint8((h & 0x00ff00) >> 8),
		B: uint8(h & 0x0000ff),
		A: 255,
	}
	return rgba
}

func (ag *AvatarGenerate) GetFontWidth() int {
	return int(ag.ctx.PointToFixed(ag.fontSize) >> 6)
}
