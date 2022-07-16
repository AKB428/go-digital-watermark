package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
TODO
1. ラベルを左下、右下、左上、右上　どれかが選べる
2. ラベルは引数指定できる
3. フォントの色を指定できる(RGBA)
4. フォントのサイズを指定できる
5. 出力ファイル名を指定できる
6. 出力ファイル拡張子を指定できる
---
[AWS]
7. 出力先をS3にしてOUT URLを返却する
----
[lambda] 別プログラム
1. AWS lambda化する
2. lambda edge化する
*/
func main() {
	const DEFAULT_LABEL = "電子透かしテスト"

	flag.Parse()

	filePath := flag.Args()[0]

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file.Close()
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	//フォーマット名表示
	fmt.Println("画像フォーマット：" + format)
	//サイズ表示
	fmt.Println("横幅=" + strconv.Itoa(config.Width) + ", 縦幅=" + strconv.Itoa(config.Height))

	ftBin, err := ioutil.ReadFile("./font/ipaexg.ttf")
	if err != nil {
		log.Fatalf("failed to load font: %s", err.Error())
	}
	ft, err := truetype.Parse(ftBin)
	if err != nil {
		log.Fatalf("failed to parse font: %s", err.Error())
	}

	file2, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file2.Close()

	var img image.Image
	if format == "png" {
		img, err = png.Decode(file2)
		if err != nil {
			log.Fatalf("failed to decode image: %s", err.Error())
		}
	} else if format == "jpeg" {
		img, err = jpeg.Decode(file2)
		if err != nil {
			log.Fatalf("failed to decode image: %s", err.Error())
		}
	} else {
		log.Fatalf("failed to decode image: %s", err.Error())
	}

	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	text := DEFAULT_LABEL

	col := color.RGBA{00, 00, 00, 120}

	opt := truetype.Options{
		Size: 30,
	}
	face := truetype.NewFace(ft, &opt)

	x, y := 100, 100
	dot := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 26)}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  dot,
	}
	d.DrawString(text)

	newFile, err := os.Create("out.png")
	if err != nil {
		log.Fatalf("failed to create file: %s", err.Error())
	}
	defer newFile.Close()

	b := bufio.NewWriter(newFile)
	if err := png.Encode(b, dst); err != nil {
		log.Fatalf("failed to encode image: %s", err.Error())
	}
}
