package picture

import (
	"image/draw"
	"image/color"
	"io"
	"errors"
	"sync"
	"image"
	"fmt"
)

var (
	InvError = errors.New("invalid argument")
)

type Per struct {
	img    draw.Image
	offset int64
	length int64
	sync.Locker
}

func (p *Per) Image() image.Image {
	return p.img
}
func (p *Per) Size() int64 {
	return p.length
}

func (p *Per) Write(bs []byte) (n int, err error) {
	p.Lock()
	defer p.Unlock()
	for _, b := range bs {
		err := p.set(p.offset, b)
		if err != nil {
			return n, err
		}
		n++
		p.offset++
	}
	return n, nil
}
func (p *Per) Seek(offset int64, whence int) (int64, error) {
	p.Lock()
	defer p.Unlock()
	var abs int64
	switch whence {
	case io.SeekCurrent:
		abs = p.offset + offset
	case io.SeekStart:
		abs = offset
	case io.SeekEnd:
		abs = p.length - offset
	default:
		return 0, InvError
	}
	if abs < 0 || abs > p.length {
		return 0, InvError
	}
	p.offset = abs
	return abs, nil
}
func (p *Per) Read(bs []byte) (n int, err error) {
	p.Lock()
	defer p.Unlock()
	var char byte
	for i := range bs {
		char, err = p.at(p.offset)
		if err != nil {
			return n, err
		}
		bs[i] = char
		p.offset++
		n++
	}
	return
}

func (p *Per) set(offset int64, b byte) error {
	if offset >= p.length {
		return InvError
	}
	p1, p2 := int(offset*2), int(offset*2+1)
	y1 := p1 / p.img.Bounds().Max.X
	y2 := p2 / p.img.Bounds().Max.X
	x1 := p1 % p.img.Bounds().Max.X
	x2 := p2 % p.img.Bounds().Max.X
	c1, c2 := byte2color(p.img.At(x1, y1), p.img.At(x2, y2), b)
	p.img.Set(x1, y1, c1)
	p.img.Set(x2, y2, c2)
	if y1 != 0{
		fmt.Printf("%d,%d",x1,y1)
	}
	return nil
}
func (p *Per) at(offset int64) (byte, error) {
	if offset == p.length {
		return 0, io.EOF
	}
	if offset > p.length {
		return 0, InvError
	}
	p1, p2 := int(offset*2), int(offset*2+1)
	y1 := p1 / p.img.Bounds().Max.X
	y2 := p2 / p.img.Bounds().Max.X
	x1 := p1 % p.img.Bounds().Max.X
	x2 := p2 % p.img.Bounds().Max.X
	return color2byte(p.img.At(x1, y1), p.img.At(x2, y2)), nil
}

func color2byte(h color.Color, l color.Color) (byte) {
	var char uint8 = 0
	data := make([]uint32, 0, 8)
	hr, hg, hb, ha := h.RGBA()
	data = append(data, hr, hg, hb, ha)
	lr, lg, lb, la := l.RGBA()
	data = append(data, lr, lg, lb, la)
	for i := 0; i < 8; i++ {
		if data[7-i]&1 == 1 {
			char = char | (1 << uint(i))
		}
	}
	return byte(char)
}
func byte2color(h color.Color, l color.Color, char byte) (color.Color, color.Color) {
	data := make([]uint8, 0, 8)
	r, g, b, a := h.RGBA()
	data = append(data, uint8(r), uint8(g), uint8(b), uint8(a))
	r, g, b, a = l.RGBA()
	data = append(data, uint8(r), uint8(g), uint8(b), uint8(a))
	for i := 0; i < 8; i++ {
		if char&(1<<uint(i)) == (1 << uint(i)) {
			// 最后一位设置为 1
			data[7-i] = data[7-i] | 1
		} else {
			// 最后一位设置为 0
			data[7-i] = data[7-i] & (^uint8(1))
		}
	}
	// TODO: 不知道为什么，这边一定要减去 2,不然无法正常工作
	return color.RGBA{R: data[0]-2, G: data[1]-2, B: data[2]-2, A: data[3]},
		color.RGBA{R: data[4]-2, G: data[5]-2, B: data[6]-2, A: data[7]}
}
func NewPer(img image.Image) (*Per) {
	pad := image.NewRGBA(img.Bounds())
	draw.Draw(pad, img.Bounds(), img, image.ZP, draw.Src)
	Pixels := (img.Bounds().Max.X - img.Bounds().Min.X) * (img.Bounds().Max.Y - img.Bounds().Min.Y)
	return &Per{img: pad, offset: 0, length: int64(Pixels / 2), Locker: &sync.Mutex{}}
}
