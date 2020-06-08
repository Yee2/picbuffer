package picture

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"sync"
)

var (
	InvError = errors.New("invalid argument")
)

type Per struct {
	img     image.NRGBA
	offset  int64
	offsetR int64
	sync.Locker
}

func (p *Per) Image() image.Image {
	return &p.img
}
func (p *Per) Size() int64 {
	return int64(len(p.img.Pix) / 8)
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
		abs = int64(len(p.img.Pix)/8) - offset
	default:
		return 0, InvError
	}
	if abs < 0 || abs > p.Size() {
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
	if offset >= p.Size() {
		return InvError
	}

	for i := 0; i < 8; i++ {
		if b&(1<<i) > 0 {
			p.img.Pix[int(offset)*8+i] |= 1
		} else {
			p.img.Pix[int(offset)*8+i] &= 0b1111_1110
		}
	}
	return nil
}
func (p *Per) at(offset int64) (b byte, e error) {
	if offset > p.Size() {
		return 0, io.EOF
	}
	for i := 0; i < 8; i++ {
		b |= (p.img.Pix[int(offset)*8+i] & 1) << i
	}
	return b, nil
}

func NewPer(img image.Image) *Per {
	pad := image.NewNRGBA(img.Bounds())
	if i, ok := img.(*image.NRGBA); ok {
		pad = i
	} else {
		draw.Draw(pad, img.Bounds(), img, image.Point{}, draw.Src)
	}
	return &Per{img: *pad, offset: 0, Locker: &sync.Mutex{}}
}
