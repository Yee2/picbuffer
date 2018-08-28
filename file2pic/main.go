package main

import (
	"os"
	"flag"
	"fmt"
	"image"
	"github.com/Yee2/picture"
	"path/filepath"
	"strings"
	"image/png"
	"golang.org/x/image/bmp"
	"image/jpeg"
	"image/gif"
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
		case ".jpg", ".jpeg":
			jpeg.Encode(out, p.Image(), &jpeg.Options{Quality: 100})
		case ".git":
			gif.Encode(out, p.Image(), &gif.Options{NumColors: 256})
		default:
			fmt.Println("不支持格式！")
			return
		}
	} else {
		var buffer [100]byte
		_, err := p.Read(buffer[0:2])
		if err != nil {
			fmt.Printf("读取失败:%s\n", *src)
			return
		}
		length = (int(buffer[0]) << 8) + int(buffer[1])
		for {
			n, err := p.Read(buffer[:])
			if n > 0 {
				if length-n < 0 {
					fmt.Printf("%s\n", buffer[:length])
					break
				}else{
					fmt.Printf("%s", buffer[:n])
				}
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
