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
	"sync"

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
func webSafePalette() []color.Color {
    pale := make([]color.Color, 0, 256)
    steps := []uint8{0x00, 0x33, 0x66, 0x99, 0xCC, 0xFF}
    for _, r := range steps {
        for _, g := range steps {
            for _, b := range steps {
                pale = append(pale, color.RGBA{r, g, b, 255})
            }
        }
    }
    // pad to 256 if needed
    for len(pale) < 256 {
        pale = append(pale, color.Black)
    }
    return pale
}

func AsciiToGifSlow(imgs []Ascii_t, height int,width int,delays []int, isColored bool){
	anim := gif.GIF{
		LoopCount: 0,
	}
	pale := webSafePalette()
	file, err := os.Create("ascii.gif")

	if err != nil {
		panic("Could not create image")
	}
	defer file.Close()
	frames := make([]*image.Paletted,len(imgs))
	var wg sync.WaitGroup
	for i,chars := range imgs {
		wg.Add(1)
		go func(i int ,chars Ascii_t){
			defer wg.Done()
			rgba := image.NewRGBA(image.Rect(0,0,width*7,height*13))
			face := basicfont.Face7x13
			drawer := &font.Drawer{
				Dst: rgba,
				Src: image.NewUniform(color.White),
				Face: face,
				Dot: fixed.Point26_6{X : fixed.I(20),Y : fixed.I(50)},
			}
			lineHeight := drawer.Face.Metrics().Height.Ceil()
			charWidth := face.Advance
			for y := 0; y < len(chars.AsciiChars); y++{
				drawer.Dot.X = fixed.I(0)
				drawer.Dot.Y = fixed.I((y+1)* lineHeight)

				for x := 0; x < len(chars.AsciiChars[y]) ; x++{
					if isColored {
						drawer.Src =image.NewUniform(color.RGBA{
							R : uint8(chars.RgbColors[y][x].R  ),
							G : uint8(chars.RgbColors[y][x].G ),
							B : uint8(chars.RgbColors[y][x].B ),
							A: 255,
						}) 
					}
					char := chars.AsciiChars[y][x]
					drawer.DrawString(string(char))
					drawer.Dot.X = fixed.I((x+1) * charWidth)
				}
			}
			paletted := image.NewPaletted(rgba.Bounds(),pale)
			draw.FloydSteinberg.Draw(paletted,rgba.Bounds(),rgba,image.Point{})
			frames[i] = paletted
			fmt.Println(i)
		}(i,chars)
		anim.Image = frames
	}
	wg.Wait()
	anim.Delay = append(anim.Delay, delays...)
	gif.EncodeAll(file,&anim)
}

func AsciiToGifFast(imgs []Ascii_t, height int,width int,delays []int, isColored bool){
	anim := gif.GIF{
		LoopCount: 0,
	}
	pale := webSafePalette()
	file, err := os.Create("ascii.gif")

	if err != nil {
		panic("Could not create image")
	}
	defer file.Close()
	frames := make([]*image.Paletted,len(imgs))
	var wg sync.WaitGroup
	for i,chars := range imgs {
		wg.Add(1)
		go func(i int ,chars Ascii_t){
			defer wg.Done()
			paletted := image.NewPaletted(image.Rect(0,0,width*7,height*13),pale)
			face := basicfont.Face7x13
			drawer := &font.Drawer{
				Dst: paletted,
				Src: image.NewUniform(color.White),
				Face: face,
				Dot: fixed.Point26_6{X : fixed.I(20),Y : fixed.I(50)},
			}
			lineHeight := drawer.Face.Metrics().Height.Ceil()
			charWidth := face.Advance
			for y := 0; y < len(chars.AsciiChars); y++{
				drawer.Dot.X = fixed.I(0)
				drawer.Dot.Y = fixed.I((y+1)* lineHeight)

				for x := 0; x < len(chars.AsciiChars[y]) ; x++{
					if isColored {
						drawer.Src =image.NewUniform(color.RGBA{
							R : uint8(chars.RgbColors[y][x].R  ),
							G : uint8(chars.RgbColors[y][x].G ),
							B : uint8(chars.RgbColors[y][x].B ),
							A: 255,
						}) 
					}
					char := chars.AsciiChars[y][x]
					drawer.DrawString(string(char))
					drawer.Dot.X = fixed.I((x+1) * charWidth)
				}
			}
			frames[i] = paletted
			fmt.Println(i)
		}(i,chars)
		anim.Image = frames
	}
	wg.Wait()
	anim.Delay = append(anim.Delay, delays...)
	gif.EncodeAll(file,&anim)
}
















