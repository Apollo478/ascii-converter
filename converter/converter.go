package converter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func RGBToGraycale(r uint32, g uint32,b uint32) float32{
	result := float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114
	return result
}
func PixelToChar(gray uint8) rune{
	
	const asciiRamp  = "@%#*+=-:. "
	const revRamp = " .:-=+*#%@"
	scale := float32(gray) /255

	index := int(scale * float32(len(revRamp)-1))
	return rune(revRamp[index])

}

func ImageToAscii(img image.Image,height int,width int,aspectRatio float64) Ascii_t {
	var res Ascii_t 
	res.AsciiChars = make([][]rune,height)
	res.RgbColors = make([][]Rgb,height)
	for i := 0; i!= height; i++{
		res.AsciiChars[i] = make([]rune,width)
		res.RgbColors[i] = make([]Rgb,width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcY := int(float64(y) / aspectRatio)
			r, g, b, _ := img.At(x, srcY).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			gray := uint8(RGBToGraycale(uint32(r8), uint32(g8), uint32(b8)))
			res.RgbColors[y][x] = Rgb{
				uint32(r8),
				uint32(g8), 
				uint32(b8),
			}
			char := PixelToChar(gray)
			res.AsciiChars[y][x] = char
		}
	}
	return res
}

func PrintAsciiImage(img [][]rune, height int,width int) {
	
	for _,row := range img {
		fmt.Println(string(row))
	}
}
func AsciiToImage(ascii Ascii_t, height int, width int,isColored bool){
	img := image.NewRGBA(image.Rect(0,0,width*7,height*13))
	draw.Draw(img,img.Bounds(),image.NewUniform(color.Black),image.Point{},draw.Src)
	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.White),
		Face: face,
		Dot: fixed.Point26_6{X : fixed.I(20),Y : fixed.I(50)},
	}
	lineHeight := drawer.Face.Metrics().Height.Ceil()
	charWidth := face.Advance
	for y := 0; y < height; y++{
		drawer.Dot.X = fixed.I(0)
		drawer.Dot.Y = fixed.I((y+1)* lineHeight)

		for x := 0; x < width ; x++{
			char := ascii.AsciiChars[y][x]
			if isColored {
				drawer.Src =image.NewUniform(color.RGBA{
					R : uint8(ascii.RgbColors[y][x].R  ),
					G : uint8(ascii.RgbColors[y][x].G ),
					B : uint8(ascii.RgbColors[y][x].B ),
					A: 255,
				}) 
			}
			drawer.DrawString(string(char))
			drawer.Dot.X = fixed.I((x+1) * charWidth)
		}
	}
	file, err := os.Create("ascii.png")

	if err != nil {
		panic("Could not create image")
	}
	defer file.Close()
	png.Encode(file,img)
}
func AsciiToGif(imgs [][][]rune, height int,width int,delays []int){
	palette := []color.Color{
		color.Black,
		color.White,
		color.RGBA{0,255,0,255},
		color.RGBA{255,0,0,255},
	}
	anim := gif.GIF{
		LoopCount: 0,
	}
	file, err := os.Create("ascii.gif")

	if err != nil {
		panic("Could not create image")
	}
	defer file.Close()
	for _,chars := range imgs {
		
		frame := image.NewPaletted(image.Rect(0,0,width*7,height*13),palette)
		face := basicfont.Face7x13
		drawer := &font.Drawer{
			Dst: frame,
			Src: image.NewUniform(color.White),
			Face: face,
			Dot: fixed.Point26_6{X : fixed.I(20),Y : fixed.I(50)},
		}
		lineHeight := drawer.Face.Metrics().Height.Ceil()
		charWidth := face.Advance
		for y := 0; y < len(chars); y++{
			drawer.Dot.X = fixed.I(0)
			drawer.Dot.Y = fixed.I((y+1)* lineHeight)

			for x := 0; x < len(chars[y]) ; x++{
				char := chars[y][x]
				drawer.DrawString(string(char))
				drawer.Dot.X = fixed.I((x+1) * charWidth)
			}
		}
		anim.Image = append(anim.Image, frame)
	}
	anim.Delay = append(anim.Delay, delays...)
	gif.EncodeAll(file,&anim)
}


















