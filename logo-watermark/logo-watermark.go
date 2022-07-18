package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

func main() {

	var oImageFilePath string
	var logoImagePath string
	var outFolderPath string
	flag.StringVar(&oImageFilePath, "f", "", "originalImageFilePath")
	flag.StringVar(&logoImagePath, "l", "", "logoImagePath")
	flag.StringVar(&outFolderPath, "o", "", "outFolderPath")
	flag.Parse()

	originFile, err := os.Open(oImageFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer originFile.Close()

	logoFile, err := os.Open(logoImagePath)
	if err != nil {
		fmt.Println(err)
	}
	defer logoFile.Close()

	originImg, format, err := image.Decode(originFile)
	if err != nil {
		log.Fatalf("failed to decode image: %s", err.Error())
	}
	logoImg, _, err := image.Decode(logoFile)
	if err != nil {
		log.Fatalf("failed to decode image: %s", err.Error())
	}

	startPointLogo := image.Point{0, 0}

	logoRectangle := image.Rectangle{startPointLogo, startPointLogo.Add(logoImg.Bounds().Size())}
	originRectangle := image.Rectangle{image.Point{0, 0}, originImg.Bounds().Size()}

	rgba := image.NewRGBA(originRectangle)
	draw.Draw(rgba, originRectangle, originImg, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, logoRectangle, logoImg, image.Point{0, 0}, draw.Over)

	if format == "jpeg" {
		out, err := os.Create("logo-watermark.jpg")
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 80

		jpeg.Encode(out, rgba, &opt)
	} else {
		out, err := os.Create("logo-watermark.png")
		if err != nil {
			fmt.Println(err)
		}
		png.Encode(out, rgba)
	}
}
