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

func main() {
	var toStegano bool
	var filename string
	var secretText string

	flag.BoolVar(&toStegano, "enc", false, "True: encrypt, False: decrypt")
	flag.StringVar(&filename, "file", "example.png", "Image filename")
	flag.StringVar(&secretText, "text", "This is sample secret text.", "Text you wanto hide")

	flag.Parse()

	img := inputImage(filename)

	if toStegano {
		encrypt(img)
		fmt.Println("==== Encryption OK ====")
	} else {
		text := decrypt(img)
		fmt.Println("==== Decription ====\n" + text)
	}
}

func inputTextBuf() []byte {
	file, err := os.Open(`input.txt`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, BUFSIZE)
	_, err = file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf
}

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

func outputImage(img *image.RGBA) {
	// 書き出し用ファイル準備
	outfile, _ := os.Create("out.png")
	defer outfile.Close()
	// 書き出し

	var encoder = &png.Encoder{CompressionLevel: -1}
	encoder.Encode(outfile, img)

}

func encrypt(img *image.RGBA) {
	// pixel := img.Pix
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	vec := inputTextBuf()
	tsize := len(vec)

	// 文字列が長すぎない前提

	for i := 0; i < tsize*2; i++ {
		c := vec[i/2]
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

	// ファイル保存
	outputImage(img)
}

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
