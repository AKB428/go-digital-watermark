package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"image/draw"
	"image/jpeg"
	"image/png"
)

/*
TODO
1.ロゴの透明値を指定できる
2.ポジションにセンターを追加
*/
func main() {

	var oImageFilePath string
	var logoImagePath string
	var outFolderPath string
	var logoPosition string
	var useOriginalFilename bool
	flag.StringVar(&oImageFilePath, "f", "", "originalImageFilePath")
	flag.StringVar(&logoImagePath, "l", "", "logoImageFilePath")
	flag.StringVar(&outFolderPath, "o", "", "outFolderPath")
	flag.StringVar(&logoPosition, "p", "", "TopLeft | TopRight | BottomLeft | BottomRight | Center")
	flag.BoolVar(&useOriginalFilename, "u", false, "出力ファイル名にオリジナルファイル名を使う")
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
		out, err := os.Create(outFilePath(oImageFilePath, "jpg", outFolderPath, useOriginalFilename))
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 80

		err = jpeg.Encode(out, rgba, &opt)
		if err != nil {
			log.Fatalf("failed to encode image: %s", err.Error())
		}

	} else {
		out, err := os.Create(outFilePath(oImageFilePath, "png", oImageFilePath, useOriginalFilename))
		if err != nil {
			fmt.Println(err)
		}
		err = png.Encode(out, rgba)
		if err != nil {
			log.Fatalf("failed to encode image: %s", err.Error())
		}
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
	case "Center":
		cx, cy := w/2, h/2
		clx, cly := lw/2, lh/2
		return cx - clx, cy - cly
	default:
		// TopLeftをデフォルトとする
		return xLeft, yTop
	}
}

func outFilePath(oFilePath string, ext string, outFolder string, useOfn bool) string {
	const defaultFileName = "logo-watermark"
	const addFileName = "-lw"
	var filaneme string

	if useOfn {
		ext := filepath.Ext(oFilePath)
		pf := strings.TrimSuffix(filepath.Base(oFilePath), ext)

		filaneme = pf + addFileName
	} else {
		filaneme = defaultFileName
	}

	if outFolder == "" {
		return filaneme + "." + ext
	} else {
		return filepath.Join(outFolder, filaneme+"."+ext)
	}

}
