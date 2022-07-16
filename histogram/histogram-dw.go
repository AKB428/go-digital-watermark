package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {

	var filePath string
	var detailView bool
	// 検証用 ファイルをコピーするだけ
	var copy bool
	flag.StringVar(&filePath, "f", "", "filePath")
	flag.BoolVar(&detailView, "d", false, "detailView")
	flag.BoolVar(&copy, "c", false, "copy")
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

	if detailView {
		return
	}

	if copy {
		copyFile(img)
		return
	}

	drawSteganography(img)
}

func drawSteganography(oimg image.Image) {

	bounds := oimg.Bounds()
	img := image.NewNRGBA(bounds)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {

			//fmt.Printf("%6d %6d %6d %6d\n", r>>8, g>>8, b>>8, a>>8)
			if y == 0 {
				// このときに文字列を埋め込む
				// 文字列分だけ書き込んでbreakで抜ける

				img.Set(x, y, color.RGBA{
					R: 12,
					G: 21,
					B: 31,
					A: 255,
				})
				//fmt.Printf("%6d %6d %6d %6d %6d %6d\n", x, y, uint8(r), uint8(g), uint8(b), uint8(a))
			} else {
				img.Set(x, y, oimg.At(x, y))
			}
		}
	}

	f, err := os.Create("sg.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func copyFile(oimg image.Image) {

	bounds := oimg.Bounds()
	img := image.NewNRGBA(bounds)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {

			img.Set(x, y, oimg.At(x, y))

		}
	}

	f, err := os.Create("sg.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
