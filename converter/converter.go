package converter

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

func RGBToGraycale(r uint32, g uint32,b uint32) float32{
	result := float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114
	return result
}
func PixelToChar(gray uint8) rune{
	
	const asciiRamp = "@%#*+=-:. "
	scale := float32(gray) /255

	index := int(scale * float32(len(asciiRamp)-1))
	return rune(asciiRamp[index])

}

func ImageToAscii(img image.Image,height int,width int,aspectRatio float64) [][]rune {
	result := make([][]rune,height)
	for i := 0; i!= height; i++{
		result[i] = make([]rune,width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcY := int(float64(y) / aspectRatio)
			r, g, b, _ := img.At(x, srcY).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			gray := uint8(RGBToGraycale(uint32(r8), uint32(g8), uint32(b8)))
			char := PixelToChar(gray)
			result[y][x] = char
		}
	}
	return result
}

func PrintAsciiImage(img [][]rune, height int,width int) {
	
	for _,row := range img {
		fmt.Println(string(row))
	}
}
