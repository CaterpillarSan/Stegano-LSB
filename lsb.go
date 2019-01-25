package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const BUFSIZE = 2048
const INPUT_TEXTFILE = "input.txt"
const OUTPUT_FILEHEAD = "stgn_"

func main() {

	var (
		toStegano bool   // true : encryption / false : decription
		filename  string // 使用するファイル名
	)

	flag.BoolVar(&toStegano, "enc", false, "true: encrypt, false: decrypt")
	flag.StringVar(&filename, "file", "example.png", "Image filename")

	flag.Parse()

	// @Note 該当ファイルはmainと同じディレクトリに置く
	img := inputImage(filename)

	if toStegano {
		// Encryption
		encImg := encrypt(img)
		outputImage(filename, encImg)
		fmt.Println("==== Encryption OK ====")
	} else {
		// Decription
		text := decrypt(img)
		fmt.Println("==== Decription ====\n" + text)
	}
}

// mainと同じディレクトリにある INPUT_TEXTFILE の文字列をBUFSIZEぶん読み出す
func inputTextBuf() []byte {
	file, err := os.Open(INPUT_TEXTFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, BUFSIZE)
	if _, err = file.Read(buf); err != nil {
		log.Fatal(err)
	}

	return buf
}

// filenameで指定されたpngファイルを読み出す
func inputImage(filename string) *image.RGBA {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return imageToRGBA(img)
}

// img構造体をファイルに書き出す
func outputImage(inputFilename string, img *image.RGBA) {
	// 書き出し用ファイル準備
	outfile, _ := os.Create(OUTPUT_FILEHEAD + inputFilename)
	defer outfile.Close()

	// 書き出し
	// 圧縮入るとまずそうなので念のため
	var encoder = &png.Encoder{CompressionLevel: -1}
	encoder.Encode(outfile, img)
}

// 画像に文字を入れ込む
func encrypt(img *image.RGBA) *image.RGBA {
	// pixel := img.Pix
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	vec := inputTextBuf()
	tsize := len(vec)

	// 文字列に対して画像が小さすぎない前提
	for i := 0; i < tsize*2; i++ {
		c := vec[i/2]

		// @todo 正方形だと死ぬ
		xx := i % dx
		yy := i % dy
		rr, gg, bb, aa := img.At(xx, yy).RGBA()

		r, g, b, a := uint8(rr), uint8(gg), uint8(bb), uint8(aa) //つらい

		if i%2 == 0 {
			g = (g & 0xF8) | ((c >> 0) & 3)
			b = (b & 0xF8) | ((c >> 2) & 3)
		} else {
			g = (g & 0xF8) | ((c >> 4) & 3)
			b = (b & 0xF8) | ((c >> 6) & 3)

		}

		img.SetRGBA(xx, yy, color.RGBA{r, g, b, a})
	}

	return img
}

// 画像から文字列取得
func decrypt(img *image.RGBA) string {
	var vec [BUFSIZE]byte // めんどいので文字数固定
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	tsize := len(vec)
	for i := 0; i < tsize*2; i++ {
		xx := i % dx
		yy := i % dy
		rr, gg, bb, aa := img.At(xx, yy).RGBA()

		_, g, b, _ := uint8(rr), uint8(gg), uint8(bb), uint8(aa) //つらい

		if i%2 == 0 {
			vec[i/2] = vec[i/2] | (g&3)<<0 | (b&3)<<2
		} else {
			vec[i/2] = vec[i/2] | (g&3)<<4 | (b&3)<<6
		}
	}

	return string(vec[:])
}

// convert given image to RGBA image
func imageToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()

	var m *image.RGBA
	var width, height int

	width = b.Dx()
	height = b.Dy()

	m = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	return m
}
