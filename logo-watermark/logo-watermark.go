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

/*TODO
1. ロゴの位置が指定できる 4箇所
2. 出力フォルダーが指定できる
3. ファイル名を加工する(オリジナルファイル名+"文字列")にする
*/

func main() {

	var oImageFilePath string
	var logoImagePath string
	var outFolderPath string
	var logoPosition string
	flag.StringVar(&oImageFilePath, "f", "", "originalImageFilePath")
	flag.StringVar(&logoImagePath, "l", "", "logoImagePath")
	flag.StringVar(&outFolderPath, "o", "", "outFolderPath")
	flag.StringVar(&logoPosition, "p", "", "TopLeft | TopRight | BottomLeft | BottomRight ")
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

	oPoint := originImg.Bounds().Size()
	lPoint := logoImg.Bounds().Size()

	x, y := positionInt(logoPosition, oPoint.X, oPoint.Y, lPoint.X, lPoint.Y)
	fmt.Printf("%d, %d, %d ,%d\n", oPoint.X, oPoint.Y, lPoint.X, lPoint.Y)
	fmt.Printf("%d, %d\n", x, y)

	startPointLogo := image.Point{x, y}

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

func positionInt(position string, w int, h int, lw int, lh int) (int, int) {

	const xLeft = 0
	const yTop = 0
	xRight := w - lw
	ybottom := h - lh

	switch position {
	case "TopLeft":
		return xLeft, yTop
	case "TopRight":
		return xRight, yTop
	case "BottomLeft":
		return xLeft, ybottom
	case "BottomRight":
		return xRight, ybottom
	default:
		// TopLeftをデフォルトとする
		return xLeft, yTop
	}
}
