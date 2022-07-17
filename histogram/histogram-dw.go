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

/*
 */

func main() {

	var filePath string
	var embeddedText string
	var detailView bool
	var decode bool
	// 検証用 ファイルをコピーするだけ
	var copy bool
	flag.BoolVar(&copy, "c", false, "copy")
	flag.BoolVar(&decode, "d", false, "decode")
	flag.StringVar(&filePath, "f", "", "filePath")
	flag.StringVar(&embeddedText, "t", "", "embedded Text 100文字まで")
	flag.BoolVar(&detailView, "v", false, "detailView")

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

	if decode {
		decodeSteganography(img, config)
		return
	}

	if detailView {
		// https://pkg.go.dev/image#Paletted.At
		bounds := img.Bounds()

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, 0).RGBA()
			fmt.Printf("%6d %6d %6d %6d %6d\n", x, r>>8, g>>8, b>>8, a>>8)
		}
		return
	}

	if copy {
		copyFile(img)
		return
	}

	drawSteganography(img, embeddedText)
}

func drawSteganography(oimg image.Image, text string) {

	bounds := oimg.Bounds()
	img := image.NewNRGBA(bounds)

	b := []byte(text)
	fmt.Println(b)
	fmt.Println(len(b))

	bc := len(b)
	var counter int

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {

			//fmt.Printf("%6d %6d %6d %6d\n", r>>8, g>>8, b>>8, a>>8)
			if y == 0 {
				// このときに文字列を埋め込む
				// 文字列分だけ書き込んでbreakで抜ける
				// 最後は0x00 0x00 0x00 0x00で終端する
				if bc >= counter {
					img.Set(x, y, color.RGBA{
						R: alignment(b, bc, counter),
						G: alignment(b, bc, counter+1),
						B: alignment(b, bc, counter+2),
						A: alignment(b, bc, counter+3),
					})
					counter += 4
				} else {
					img.Set(x, y, oimg.At(x, y))
				}
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

func alignment(b []byte, len int, counter int) uint8 {

	if counter < len {
		return uint8(b[counter])
	} else {
		return 0
	}
}

func decodeSteganography(img image.Image, config image.Config) {
	for x := 0; x < config.Width; x++ {
		r, g, b, a := img.At(x, 0).RGBA()
		//fmt.Printf("%6d %6d %6d %6d %6d\n", x, r>>8, g>>8, b>>8, a>>8)

		//var data []byte

		r8 := r >> 8
		g8 := g >> 8
		b8 := b >> 8
		a8 := a >> 8

		rb := i32tob(r)
		gb := i32tob(g)
		bb := i32tob(b)
		ab := i32tob(a)

		fmt.Println(rb[1])
		fmt.Println(gb[1])
		fmt.Println(bb[1])
		fmt.Println(ab[1])

		slice := []byte{rb[1], gb[1], bb[1], ab[1]}

		fmt.Println(string(slice))
		// 処理終端
		if r8 == 0 && g8 == 0 && b8 == 0 && a8 == 00 {
			break
		}
	}
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
