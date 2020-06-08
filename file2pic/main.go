package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"picture"
	"strings"
)

func main() {
	src := flag.String("i", "", "源图片")
	dst := flag.String("o", "out.png", "输出图片")
	text := flag.String("c", "", "需要编码的内容")
	flag.Parse()
	f, err := os.Open(*src)
	if err != nil {
		fmt.Printf("无法读取文件:%s\n", *src)
		flag.Usage()
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("目标文件不是图片:%s\n", *src)
		flag.Usage()
		return
	}
	p := picture.NewPer(img)

	if length := len(*text); length > 0 {
		_, err = p.Write([]byte{byte(length >> 8), byte(length)})
		if err != nil {
			fmt.Printf("写入失败:%s\n", *src)
			return
		}
		_, err = p.Write([]byte(*text))
		if err != nil {
			fmt.Printf("写入失败:%s\n", *src)
			return
		}

		out, err := os.OpenFile(*dst, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("无法创建文件:%s\n", *dst)
			flag.Usage()
			return
		}
		defer out.Close()
		switch strings.ToLower(filepath.Ext(*dst)) {
		case ".png":
			png.Encode(out, p.Image())
		case ".bmp":
			bmp.Encode(out, p.Image())
		default:
			fmt.Println("不支持格式！")
			return
		}
	} else {
		var buffer [0xffff]byte
		_, err := p.Read(buffer[0:2])
		if err != nil {
			fmt.Printf("读取失败:%s\n", *src)
			return
		}
		length = (int(buffer[0]) << 8) + int(buffer[1])
		n, err := p.Read(buffer[:length])
		fmt.Println(string(buffer[:n]))
		if err != nil {
			fmt.Println(err)
		}
	}
}
