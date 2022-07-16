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
	"unicode/utf8"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
TODO
1. ラベルを左下、右下、左上、右上　どれかが選べる
<ok> 2. ラベル文字列は引数指定できる
3. フォントの色を指定できる(RGBA)
4. フォントのサイズを指定できる
5. 出力ファイル名を指定できる
6. 出力ファイル拡張子を指定できる
[エラー処理]
1. 文字数、フォントサイズ　から許容できる画像サイズ(横縦)をチェックする
---
[AWS]
7. 出力先をS3にしてOUT URLを返却する
----
[lambda] 別プログラム
1. AWS lambda化する
2. lambda edge化する
---
[ラベルでない透かし] 別プログラム
pngの拡張データを使う
jpeg exifを使う
jepg exif以外？を使う
--
[別プログラム]
ロゴ(画像)ラベル埋め込み
*/
func main() {
	const DEFAULT_LABEL_TEXT = "電子透かしテスト"
	fontSize := 30
	var filePath string
	var labelPosition string
	var labelText string
	flag.StringVar(&filePath, "f", "", "filePath")
	flag.StringVar(&labelPosition, "p", "", "Label Position")
	flag.StringVar(&labelText, "t", "", "labelText")
	flag.Parse()

	if labelText == "" {
		labelText = DEFAULT_LABEL_TEXT
	}

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

	wordCount := utf8.RuneCountInString(labelText)
	col := color.RGBA{00, 00, 00, 120}

	opt := truetype.Options{
		Size: float64(fontSize),
	}
	face := truetype.NewFace(ft, &opt)

	x, y := labelPositionInt(labelPosition, config.Width, config.Height, fontSize, wordCount)
	dot := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	fmt.Println(dot)

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  dot,
	}
	d.DrawString(labelText)

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

func labelPositionInt(labelPosition string, width int, height int, fontSize int, wordCount int) (int, int) {

	bottomHeightPadding := 10
	upperHeight := fontSize
	bottomHeight := height - bottomHeightPadding

	leftWidht := 0
	rightWidht := width - (wordCount * fontSize)

	switch labelPosition {
	case "UpperLeft":
		return leftWidht, upperHeight
	case "UpperRight":
		return rightWidht, upperHeight
	case "BottomLeft":
		return leftWidht, bottomHeight
	case "BottomRight":
		return rightWidht, bottomHeight
	default:
		return leftWidht, upperHeight
	}
}
