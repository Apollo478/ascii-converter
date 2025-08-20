package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

)
func RGBToGraycale(r uint32, g uint32,b uint32) float32{
	result := float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114
	return result
}
func main(){
	file,err := os.Open("dd.jpg")
	if err != nil {
		panic("could not open file")
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
	for i := 0 ;i!= height; i++ {
		
		for j := 0; j!= height; j++ {
			r,g,b,_ := img.At(i,j).RGBA()
			grayScaleMatrix[i][j] = uint8(RGBToGraycale(r,g,b))
			fmt.Printf("%3d ",grayScaleMatrix[i][j] )
		}
		fmt.Println()
	}
}
