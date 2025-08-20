package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/Apollo478/ascii-converter/converter"
)

func main(){
	file,err := os.Open("images/riri3.jpg")
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
	aspectRatio := 0.5
	height = int(float32(height) * float32(aspectRatio))
	fmt.Println(height,width)
	ascii := converter.ImageToAscii(img,height,width,aspectRatio)
	converter.PrintAsciiImage(ascii,height,width)
	
}
