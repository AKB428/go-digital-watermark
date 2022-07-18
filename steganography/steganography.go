package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

/*
 指定した画像ファイルにテキスト文字列を透かしで入れるプログラム
 RGB値を変更するため元画像は多少劣化する。ただし透かしの判別は難しい。
*/
/*
TODO
mode1A スタート位置を可変にする
mode1B スタート位置をランダムにしてスタートポジションはファイル名等に書き出す
mode2  ビット操作にし、RGBの末尾の1bitを変更し文字を書き込む->1バイト書き込むのに8バイトのRGBが必要
*/

func main() {

	var filePath string
	var embeddedText string
	var detailView bool
	var decode bool
	var copy bool
	flag.BoolVar(&copy, "c", false, "プログラム検証用　rgb copy")
	flag.BoolVar(&decode, "d", false, "decode")
	flag.StringVar(&filePath, "f", "", "filePath")
	flag.StringVar(&embeddedText, "t", "", "embedded Text [max = image width * 3byte]")
	flag.BoolVar(&detailView, "v", false, "detailView")

	flag.Parse()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %s", err.Error())
	}
	//フォーマット名表示
	fmt.Println("画像フォーマット：" + format)
	//サイズ表示
	point := img.Bounds().Size()
	fmt.Println("横幅=" + strconv.Itoa(point.X) + ", 縦幅=" + strconv.Itoa(point.Y))

	if format != "png" {
		log.Fatalf("対応している画像フォーマットではありません")
	}

	if decode {
		decodeText := decodeSteganography(img, point.X)
		fmt.Println(decodeText)
		return
	}

	if detailView {
		fmt.Printf("%8s %6s %6s %6s %6s\n", "[x,y=0]", "R", "G", "B", "A")
		for x := 0; x < point.X; x++ {
			r, g, b, a := img.At(x, 0).RGBA()
			fmt.Printf("%8d %6d %6d %6d %6d\n", x+1, r>>8, g>>8, b>>8, a>>8)
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
	fmt.Printf("input text byte %d\n", len(b))

	bc := len(b)
	var counter int
	fmt.Println("Original Data")
	fmt.Printf("%8s %6s %6s %6s %6s\n", "[x,y=0]", "R", "G", "B", "A")
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {

			if y == 0 {
				// y座標が0の場合に文字列を埋め込む、指定された文字数分だけ処理
				// 最後は0x00 0x00 0x00で終端する
				if bc >= counter || bc%3 != 0 && bc >= counter-3 {
					// fmt.Printf("%d : %d\n", bc, counter)

					or, og, ob, oa := oimg.At(x, y).RGBA()
					fmt.Printf("%8d %6d %6d %6d %6d\n", x+1, or>>8, og>>8, ob>>8, oa>>8)

					color := color.RGBA{
						R: alignment(b, bc, counter),
						G: alignment(b, bc, counter+1),
						B: alignment(b, bc, counter+2),
						A: 255, //A値は255で固定、RGB値に依存するため任意の値を入れるとRGB値が壊れる
					}
					img.Set(x, y, color)
					counter += 3

				} else {
					img.Set(x, y, oimg.At(x, y))
				}
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

func decodeSteganography(img image.Image, width int) string {
	textb := []byte{}
	fmt.Printf("%8s %6s %6s %6s %6s\n", "[x,y=0]", "R", "G", "B", "A")

	for x := 0; x < width; x++ {
		r, g, b, a := img.At(x, 0).RGBA()

		r8 := r >> 8
		g8 := g >> 8
		b8 := b >> 8

		rb := i32tob(r)
		gb := i32tob(g)
		bb := i32tob(b)
		ab := i32tob(a)

		fmt.Printf("%8d %6d %6d %6d %6d\n", x+1, rb[1], gb[1], bb[1], ab[1])

		textb = append(textb, rb[1], gb[1], bb[1])

		// 処理終端
		if r8 == 0 && g8 == 0 && b8 == 0 {
			break
		}
	}
	return string(textb)
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
