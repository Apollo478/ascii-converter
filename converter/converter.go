package converter

import (
	"fmt"
	"image"
	"image/color"
	"time"

	//"image/draw"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"sync"

	"os"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/term"
)
var RevRamp string = ""
const (
	SimpleRamp = ".-+*=%@#&WMN$"
	FullRamp = ".'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)                                
func ReverseRamp(ramp string) string {
	runes := []rune(ramp)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func RGBToGraycale(r uint32, g uint32,b uint32) float32{
	result :=(float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114)
	return result
}
func PixelToChar(gray uint8) rune{
	
	scale := float32(gray ) /255

	index := int(scale * float32(len(RevRamp)-1))
	return rune( RevRamp[index] )

}
func ImageToGrayScale(img image.Image,opts Options,prevFrame image.Image) [][]uint8{
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()
	height = int(float32(height) * float32(opts.AspectRatio))
	grayScale := make([][]uint8,height)	
	for i := 0; i!= height; i++{
		grayScale[i] = make([]uint8,width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcY := int(float64(y) / opts.AspectRatio)
			if srcY >= img.Bounds().Dy() {
				srcY = img.Bounds().Dy() - 1
			}
			r, g, b ,a:= img.At(x, srcY).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)
			if a8 == 0 && opts.BlendPrev   && prevFrame != nil && prevFrame.Bounds().Dx() >  x{
				srcY := int(float64(y) / opts.AspectRatio)
				if   srcY >= prevFrame.Bounds().Dy() {
					srcY = prevFrame.Bounds().Dy() - 1
				}
				pr, pg, pb ,_ := prevFrame.At(x, srcY).RGBA()
				pr8 := uint8(pr >> 8)
				pg8 := uint8(pg >> 8)
				pb8 := uint8(pb >> 8)
				grayScale[y][x] = uint8(RGBToGraycale(uint32(pr8), uint32(pg8), uint32(pb8)))
			} else {
				grayScale[y][x] = uint8(RGBToGraycale(uint32(r8), uint32(g8), uint32(b8)))
			}
		}
	}
	return grayScale
	
}
func CompressGrayScale(gray [][]uint8,opts Options) [][]uint8{
	
	if opts.Compression == 0 {
		return	nil
	}
		
	if len(gray) == 0 {
		return nil
	}
	height := len(gray) / opts.Compression
	width := len(gray[0]) / opts.Compression
	grayScale := make([][]uint8,height)	
	for i := 0; i!= height; i++{
		grayScale[i] = make([]uint8,width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++{
			sum := 0	
		 for dy := 0; dy < opts.Compression; dy++ {
                for dx := 0; dx < opts.Compression; dx++ {
                    sum += int(gray[y*opts.Compression+dy][x*opts.Compression+dx])
                }
            }
            grayScale[y][x] = uint8(sum / (opts.Compression*opts.Compression))
		}
	}
	return grayScale
}

func ImageToRgb(img image.Image,opts Options,prevFrame image.Image) [][]Rgb {
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()
	height = int(float32(height) * float32(opts.AspectRatio))
	rgbScale := make([][]Rgb,height)	
	for i := 0; i!= height; i++{
		rgbScale[i] = make([]Rgb,width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcY := int(float64(y) / opts.AspectRatio)
			if srcY >= img.Bounds().Dy() {
				srcY = img.Bounds().Dy() - 1
			}
			r, g, b, a := img.At(x, srcY).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)
			
			if a8 == 0 && opts.BlendPrev && prevFrame != nil && prevFrame.Bounds().Dx() >  x  {
				if srcY >= prevFrame.Bounds().Dy() {
					srcY = prevFrame.Bounds().Dy() - 1
				}
				pr, pg, pb ,pa := prevFrame.At(x, srcY).RGBA()
				pr8 := uint8(pr >> 8)
				pg8 := uint8(pg >> 8)
				pb8 := uint8(pb >> 8)
				pa8 := uint8(pa >> 8)
				rgbScale[y][x] = Rgb{
					uint32(pr8),
					uint32(pg8), 
					uint32(pb8),
					255,
				}
				if opts.UseAlpha {
					rgbScale[y][x].A = uint32(pa8)
				}
			} else {
				rgbScale[y][x] = Rgb{
					uint32(r8),
					uint32(g8), 
					uint32(b8),
					255,
				}
				if opts.UseAlpha {
					rgbScale[y][x].A = uint32(a8)
				}
			}
		}
	}
	return rgbScale
}

func CompressRgb(rgb [][]Rgb,opts Options) [][]Rgb {
	if opts.Compression == 0 {
		return	nil
	}
		
	if len(rgb) == 0 {
		return nil
	}
	height := len(rgb) / opts.Compression
	width := len(rgb[0]) / opts.Compression
	rgbScale := make([][]Rgb,height)	
	for i := 0; i!= height; i++{
		rgbScale[i] = make([]Rgb,width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++{
			Rsum := 0	
			Gsum := 0	
			Bsum := 0	
			Asum := 0	
		 for dy := 0; dy < opts.Compression; dy++ {
                for dx := 0; dx < opts.Compression; dx++ {
                    Rsum += int(rgb[y*opts.Compression+dy][x*opts.Compression+dx].R)
                    Gsum += int(rgb[y*opts.Compression+dy][x*opts.Compression+dx].G)
                    Bsum += int(rgb[y*opts.Compression+dy][x*opts.Compression+dx].B)
                    Asum += int(rgb[y*opts.Compression+dy][x*opts.Compression+dx].A)
                }
            }
            rgbScale[y][x] = Rgb{
				R:	uint32(Rsum / (opts.Compression*opts.Compression)),
				G:	uint32(Gsum / (opts.Compression*opts.Compression)),
				B:	uint32(Bsum / (opts.Compression*opts.Compression)),
				A:	uint32(Asum / (opts.Compression*opts.Compression)),
			}

		}
	}
	return rgbScale
}

func ImageToAscii(img image.Image,opts Options,prevFrame image.Image) Ascii_t {
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()

	grayScale := ImageToGrayScale(img,opts,prevFrame)
	rgbScale := ImageToRgb(img,opts,prevFrame)
	if len(grayScale) == 0 {
		return Ascii_t{}
	}
	if opts.Compression != 0 {
		grayScale = CompressGrayScale(grayScale,opts)
		rgbScale = CompressRgb(rgbScale,opts)
		height = len(grayScale)
		width = len(grayScale[0])
	}
	var res Ascii_t 
	res.AsciiChars = make([][]rune,height)
	res.RgbColors = make([][]Rgb,height)
	for i := 0; i!= height; i++{
		res.AsciiChars[i] = make([]rune,width)
		res.RgbColors[i] = make([]Rgb,width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			char := PixelToChar(grayScale[y][x])
			res.AsciiChars[y][x] = char
		}
	}
	res.RgbColors = rgbScale
	return res
}

func PrintAsciiImage(ascii Ascii_t, opts Options) {
	if opts.ClearScreen {
		fmt.Print("\033c")
	}
	for y := 0; y < len(ascii.AsciiChars); y++ {
		for x := 0; x < len(ascii.AsciiChars[y]); x++ {
			char := ascii.AsciiChars[y][x]
			color := ascii.RgbColors[y][x]

			if opts.UseColor {
				fmt.Printf("\x1b[38;2;%d;%d;%dm%c",
					color.R, color.G, color.B, char)
			} else {
				fmt.Printf("%c", char)
			}
		}
		fmt.Print("\x1b[0m\n") 
	}
}

func PrintAsciiGif(asciis []Ascii_t, opts Options,delays []int) {
	
		for {
			for i,ascii := range asciis {
				PrintAsciiImage(ascii,opts)
				time.Sleep(time.Duration(delays[i]*10) * time.Millisecond)
			}
		}
}

func AsciiToImage(ascii Ascii_t, opts Options){
	if len(ascii.AsciiChars) == 0 {
		return
	}
	height :=len(ascii.AsciiChars) 
	width := len(ascii.AsciiChars[0])

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
			if opts.UseColor {
				drawer.Src =image.NewUniform(color.RGBA{
					R : uint8(ascii.RgbColors[y][x].R  ),
					G : uint8(ascii.RgbColors[y][x].G ),
					B : uint8(ascii.RgbColors[y][x].B ),
					A: uint8(ascii.RgbColors[y][x].A ),
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
    for len(pale) < 256 {
        pale = append(pale, color.Black)
    }
    return pale
}

func ResizePaletted(img *image.Paletted,opts Options) *image.Paletted {
	if opts.FitTerminal {
		fd := int(os.Stdout.Fd())
		if term.IsTerminal(fd) {
			width, height, err := term.GetSize(fd)	
			if err != nil {
				panic("Could not get terminal size")
			}
			opts.Height = height * 2 - 1
			opts.Width = width
		}
	}
	if img.Bounds().Dx() == opts.Width && img.Bounds().Dy() == opts.Height {
		return img
	}
	dst := image.NewPaletted(image.Rect(0,0,opts.Width,opts.Height),img.Palette)
	
	xRatio := float64(img.Bounds().Dx()) / float64(opts.Width)
	yRatio := float64(img.Bounds().Dy()) / float64(opts.Height)

	for y := 0; y < opts.Height; y++ {
		for x := 0; x < opts.Width; x++ {
			origX := int(float64(x) * xRatio)
			origY := int(float64(y) * yRatio)

			c := img.At(origX, origY)

			index := dst.Palette.Index(c)
			dst.SetColorIndex(x, y, uint8(index))
		}
	}
	return dst
}
func ResizeImg(img image.Image,opts Options) image.Image {
	if opts.FitTerminal {
		fd := int(os.Stdout.Fd())
		if term.IsTerminal(fd) {
			width, height, err := term.GetSize(fd)	
			if err != nil {
				panic("Could not get terminal size")
			}
			opts.Height = height * 2 - 1
			opts.Width = width
		}
	}
	dst := image.NewRGBA(image.Rect(0,0,opts.Width,opts.Height))
	draw.BiLinear.Scale(dst,dst.Bounds(),img,img.Bounds(),draw.Over,nil)
	return dst
}

func AsciiToGifSlow(imgs []Ascii_t,opts Options,delays []int,disposal []byte,palets []color.Palette){
	anim := gif.GIF{
		LoopCount: 0,
	}
	file, err := os.Create("ascii.gif")

	if err != nil {
		panic("Could not create image")
	}
	defer file.Close()
	frames := make([]*image.Paletted,len(imgs))
	PrintProgress(0,len(imgs))
	var wg sync.WaitGroup
	framesDone := 0
	for i,chars := range imgs {
		if opts.Parallel {
			wg.Add(1)
			go func(i int ,chars Ascii_t){
				defer wg.Done()
				frames[i] = AsciiToPalleted(chars,opts,palets[i])
				if opts.ShowProgress {
					PrintProgress(framesDone,len(imgs))
					framesDone++
				}
			}(i,chars)
		} else {
			frames[i] = AsciiToPalleted(chars,opts,palets[i])
			if opts.ShowProgress {
				PrintProgress(framesDone,len(imgs))
				framesDone++
			}
		}
		anim.Image = frames
	}
	if opts.Parallel {
		wg.Wait()
	}
	anim.Delay = append(anim.Delay, delays...)
	anim.Disposal = append(anim.Disposal, disposal...)
	gif.EncodeAll(file,&anim)
}

func PrintProgress(curr int,max int) {
	
fd := int(os.Stdout.Fd())
	width := 0 // fallback
	if term.IsTerminal(fd) {
		w, _, err := term.GetSize(fd)
		if err == nil {
			width = w
		}
	}

	barWidth := width - 10 
	progress := float64(curr) / float64(max)
	filled := int(progress * float64(barWidth))

	fmt.Print("\r")

	fmt.Print("[")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			fmt.Print("█")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %3d%%", int(progress*100))

	if curr == max {
		fmt.Println()
	}}

func AsciiToPalleted(chars Ascii_t,opts Options,pale []color.Color) *image.Paletted {
	
	var height int = 0
	var width int = 0
			if len(chars.AsciiChars) != 0{
				height = len(chars.AsciiChars)
				width = len(chars.AsciiChars[0])
			}
			rgba := image.NewRGBA(image.Rect(0,0,width*7,height*13))
			draw.Draw(rgba,rgba.Bounds(),image.NewUniform(color.Black),image.Point{},draw.Src)
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
					if opts.UseColor && uint8(chars.RgbColors[y][x].A )!= 0 {
						drawer.Src =image.NewUniform(color.RGBA{
							R : uint8(chars.RgbColors[y][x].R  ),
							G : uint8(chars.RgbColors[y][x].G ),
							B : uint8(chars.RgbColors[y][x].B ),
							A: uint8(chars.RgbColors[y][x].A ),
						}) 
					}
					char := chars.AsciiChars[y][x]
					drawer.DrawString(string(char))
					drawer.Dot.X = fixed.I((x+1) * charWidth)
				}
			}
			paletted := image.NewPaletted(rgba.Bounds(),pale)
			draw.FloydSteinberg.Draw(paletted,rgba.Bounds(),rgba,image.Point{})
			return paletted
}
// func AsciiToGifFast(imgs []Ascii_t, height int,width int,delays []int,disposal []byte, isColored bool){
// 	anim := gif.GIF{
// 		LoopCount: 0,
// 	}
// 	pale := webSafePalette()
// 	file, err := os.Create("ascii.gif")
//
// 	if err != nil {
// 		panic("Could not create image")
// 	}
// 	defer file.Close()
// 	frames := make([]*image.Paletted,len(imgs))
// 	var wg sync.WaitGroup
// 	for i,chars := range imgs {
// 		wg.Add(1)
// 		go func(i int ,chars Ascii_t){
// 			defer wg.Done()
// 			if len(chars.AsciiChars) != 0{
// 				height = len(chars.AsciiChars)
// 				width = len(chars.AsciiChars[0])
// 			}
// 			paletted := image.NewPaletted(image.Rect(0,0,width*7,height*13),pale)
// 			face := basicfont.Face7x13
// 			drawer := &font.Drawer{
// 				Dst: paletted,
// 				Src: image.NewUniform(color.White),
// 				Face: face,
// 				Dot: fixed.Point26_6{X : fixed.I(20),Y : fixed.I(50)},
// 			}
// 			lineHeight := drawer.Face.Metrics().Height.Ceil()
// 			charWidth := face.Advance
// 			for y := 0; y < len(chars.AsciiChars); y++{
// 				drawer.Dot.X = fixed.I(0)
// 				drawer.Dot.Y = fixed.I((y+1)* lineHeight)
//
// 				for x := 0; x < len(chars.AsciiChars[y]) ; x++{
// 					if isColored {
// 						drawer.Src =image.NewUniform(color.RGBA{
// 							R : uint8(chars.RgbColors[y][x].R  ),
// 							G : uint8(chars.RgbColors[y][x].G ),
// 							B : uint8(chars.RgbColors[y][x].B ),
// 							A: uint8(chars.RgbColors[y][x].A ),
// 						}) 
// 					}
// 					char := chars.AsciiChars[y][x]
// 					drawer.DrawString(string(char))
// 					drawer.Dot.X = fixed.I((x+1) * charWidth)
// 				}
// 			}
// 			frames[i] = paletted
// 			fmt.Println(i)
// 		}(i,chars)
// 		anim.Image = frames
// 	}
// 	wg.Wait()
// 	anim.Delay = append(anim.Delay, delays...)
// 	anim.Disposal = append(anim.Disposal, disposal...)
// 	gif.EncodeAll(file,&anim)
// }















