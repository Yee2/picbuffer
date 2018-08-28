package picture

import (
	"testing"
	"os"
	"image"
	"image/png"
)

const (
	filename = "out.bmp"
	data     = "hello world!"
	length   = len("hello world!")
)

func TestDecode(t *testing.T) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("读取图片失败:%s", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("解析图片失败:%s", err)
	}
	var raw [100]byte

	buffer := NewPer(img)
	n, err := buffer.Read(raw[0:length])
	if err != nil {
		t.Fatalf("%s", err)
	}
	if string(raw[0:n]) != data{
		t.Fatalf("无法读取字符串:%s", raw[:n])
	}
}
func TestEncode(t *testing.T) {
	f, err := os.Open("test.png")
	if err != nil {
		t.Skipf("未找到测试文件，跳过测试")
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		t.Skipf("文件损坏，跳过测试")
	}
	buffer := NewPer(img)
	n, err := buffer.Write([]byte(data))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if n != length {
		t.Fatalf("写入数据长度错误")
	}

	out, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		t.Skipf("无法创建文件，跳过测试")
	}
	out.Truncate(0)

	err = png.Encode(out, buffer.Image())
	if err != nil {
		t.Fatalf("保存图片失败:%s", err)
	}
	out.Close()
	t.Run("", TestDecode)
}


func TestSet(t *testing.T) {
	buffer := NewPer(image.NewRGBA(image.Rect(0, 0, 20, 20)))
	buffer.set(0, 'a')
	buffer.test()
	b, err := buffer.at(0)
	if err != nil {
		t.Fatalf("error:%s", err)
	}
	if b != 'a' {
		t.Fatalf("error")
	}
}
