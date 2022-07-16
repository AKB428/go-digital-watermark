package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {

	var filePath string
	flag.StringVar(&filePath, "f", "", "filePath")
	flag.Parse()

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
	// https://pkg.go.dev/image#Paletted.At
	bounds := img.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		r, g, b, a := img.At(x, 0).RGBA()
		fmt.Printf("%6d %6d %6d %6d %6d\n", x, r>>8, g>>8, b>>8, a>>8)
	}
}
