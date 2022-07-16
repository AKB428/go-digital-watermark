package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	flag.Parse()

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	ftBin, err := ioutil.ReadFile("./font/ipaexg.ttf")
	if err != nil {
		log.Fatalf("failed to load font: %s", err.Error())
	}
	ft, err := truetype.Parse(ftBin)
	if err != nil {
		log.Fatalf("failed to parse font: %s", err.Error())
	}

	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %s", err.Error())
	}
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	text := "てすと"

	col := color.RGBA{33, 33, 33, 1}

	opt := truetype.Options{
		Size: 20,
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
