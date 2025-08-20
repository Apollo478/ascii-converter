package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

)

func RGBToGraycale(r uint32, g uint32,b uint32) float32{
	result := float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114
	return result
}
func pixelToChar(gray uint8) rune{
	
	const asciiRamp = "@%#*+=-:. "
	scale := float32(gray) /255

	index := int(scale * float32(len(asciiRamp)-1))
	return rune(asciiRamp[index])

}
func main(){
	file,err := os.Open("images/smol.png")
	if err != nil {
		panic("Could not open file")
	}
	defer file.Close()
	img,_,err := image.Decode(file)
	if err != nil {
		panic("Could not decode image")
	}
	bounds := img.Bounds()
	width  := bounds.Dx()
	height  := bounds.Dy()
	grayScaleMatrix := make([][]uint8,height)
	for i := 0; i!= height; i++{
		grayScaleMatrix[i] = make([]uint8,width)
	}
for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
        r, g, b, _ := img.At(x, y).RGBA()
        r8 := uint8(r >> 8)
        g8 := uint8(g >> 8)
        b8 := uint8(b >> 8)

        gray := uint8(RGBToGraycale(uint32(r8), uint32(g8), uint32(b8)))
        char := pixelToChar(gray)
        fmt.Printf("%c", char)
    }
    fmt.Println()
}}
