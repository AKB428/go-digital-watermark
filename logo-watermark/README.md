# logo watermark (透かし画像 ロゴ透かし)

指定したオリジナル画像に、ロゴ画像を重ねる「透かし画像」作成プログラムです

## ビルド

```
go build -o ./bin/logo-watermark logo-watermark/logo-watermark.go
```

## 実行

### 元画像がJPEG (ロゴはpng)

```
./bin/logo-watermark -f ./image/photo-jpeg.jpg -l ./image/logo-sample2.png
```

### 元画像がPNG (ロゴはpng)

```
./bin/logo-watermark -f ./image/photo-hiziribashi.png -l ./image/logo-sample2.png
```

![logo-watermark-TL](../image-embedded-sample/logo-watermark-TL.png "image-embedded-sample/logo-watermark-TL.png")


### ロゴの位置を指定 (BottomRight=右下の場合)

```
./bin/logo-watermark -f ./image/photo-jpeg.jpg -l ./image/logo-sample2.png -p BottomRight
```

![logo-watermark-BR](../image-embedded-sample/logo-watermark-BR.jpg "image-embedded-sample/logo-watermark-BR.jpg")

### 出力フォルダの指定 -o オプション

```
./bin/logo-watermark -f ./image/photo-jpeg.jpg -l ./image/logo-sample2.png -p BottomRight -o ./private
```

### 出力ファイル名にオリジナルのファイル名を使用する -u オプション

```
./bin/logo-watermark -f ./image/photo-jpeg.jpg -l ./image/logo-sample2.png -p BottomLeft -u
```

photo-jpeg-lw.jpg がカレントディレクトリに出力される。オリジナルのファイル名の末尾に"-lw"をつける。

### -u と -o の組み合わせ

```
 ./bin/logo-watermark -f ./image/photo-jpeg.jpg -l ./image/logo-sample2.png -p BottomLeft -u -o ./private
```