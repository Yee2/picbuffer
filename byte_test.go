package picture

import (
	"testing"
	"image/color"
)

func TestColor2Byte(t *testing.T) {
	var c1, c2 color.RGBA
	for i := 'a'; i <= 'z'; i++ {
		h, l := byte2color(c1, c2, byte(i))
		if color2byte(h, l) != byte(i) {
			hr, hg, hb, ha := h.RGBA()
			lr, lg, lb, la := l.RGBA()
			t.Logf("%b%b%b%b%b%b%b%b\n", hr&1, hg&1, hb&1, ha&1, lr&1, lg&1, lb&1, la&1)
			t.Fatalf("转换失败:%c(%08b),结果:%08b", byte(i), i, color2byte(byte2color(c1, c2, byte(i))))
		}
		if color2byte(byte2color(c1, c2, byte(i+('A'-'a')))) != byte(i+('A'-'a')) {
			t.Fatalf("转换失败:%c(%b),结果:%b", byte(i), i, color2byte(byte2color(c1, c2, byte(i))))
		}
	}
	for i := '0'; i <= '9'; i++ {
		if color2byte(byte2color(c1, c2, byte(i))) != byte(i) {
			t.Fatalf("转换失败:%c(%b)", byte(i), i)
		}
	}
}
