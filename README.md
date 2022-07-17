# go-digital-watermark 

## Golangで画像に対して電子透かしを実行するアプローチ検証コードです


## steganography (ステガノグラフィー)

指定されたPNGファイルのRGB値を変更し、指定された文字列のバイトデータを埋め込みます。
デフォルトでは x,y座標の先頭からy=0を固定で1行目を書き換えます。
たとえば ABCであれば 1,0の座標データを書き換えます (1ピクセルはRGBAの4バイトのためx=1以内)。
また、RGBAのAはデータ埋め込みに使用しないため1座標に埋め込めるのは３バイトになります。
このためx軸(横の長さ[width])*3バイトまでの文字列を埋め込むことが可能です。

ABC -> x,y=[1,0]を書き換え (3byte)

ABCD -> x,y=[1,0][2,0]を書き換え (4byte)

ABCDE -> x,y=[1,0][2,0]を書き換え (5byte)

ABCDEFG -> x,y=[1,0][2,0][3,0]を書き換え (7byte)

### 実行ファイルビルド

```
:go build -o ./bin/steganography steganography/steganography.go
```

### 埋め込み

image/red.png に"ABCDE"文字列を埋め込む

```
./bin/steganography  -f image/red.png -t ABCDE
```

image/red.png に"あいう漢字"文字列を埋め込む

```
./bin/steganography -f image/red.png -t あいう漢字
```

### デコード

```
:./bin/steganography -f sg.png -d
```

