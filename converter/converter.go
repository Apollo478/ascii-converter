package converter

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"strings"

	_ "image/png"
	"os"
	"sync"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/term"
)
var RevRamp string = ""
var asciiBuffer []byte
var prevChars [][]rune
var prevColors [][]Rgb
const (
	SimpleRamp = " .-+*=%@#&WMN$"
	FullRamp = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	DefaultWidth = 140
	DefaultHeight = 120

)                                
var RosePine = []Rgb{
    {197, 199, 198,255}, 
    {235, 188, 186,255}, 
    {235, 188, 186,255}, 
    {245, 224, 220,255}, 
    {156, 207, 216,255}, 
    {163, 190, 140,255}, 
    {224, 222, 244,255}, 
}
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

func ImageToGrayScale(img image.Image,opts Options) [][]uint8{
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
			r, g, b ,_:= img.At(x, srcY).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			grayScale[y][x] = uint8(RGBToGraycale(uint32(r8), uint32(g8), uint32(b8)))
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

func ImageToRgb(img image.Image,opts Options) [][]Rgb {
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
				A:255,
			}
			if opts.UseAlpha {
				rgbScale[y][x].A = uint32(Asum / (opts.Compression*opts.Compression))
			}

		}
	}
	return rgbScale
}

func RgbBufferToAscii(buffer []byte, opts Options) Ascii_t {
	rawWidth := opts.Width
	rawHeight := opts.Height

	displayWidth := rawWidth
	displayHeight := int(float64(rawHeight) * opts.AspectRatio)

	grayScale := make([][]uint8, displayHeight)
	rgb := make([][]Rgb, displayHeight)
	for y := 0; y < displayHeight; y++ {
		grayScale[y] = make([]uint8, displayWidth)
		rgb[y] = make([]Rgb, displayWidth)
	}

	i := 0
	for y := 0; y < rawHeight; y++ {
		for x := 0; x < rawWidth; x++ {
			r := buffer[i]
			g := buffer[i+1]
			b := buffer[i+2]

			displayY := int(float64(y) * opts.AspectRatio)
			if displayY >= displayHeight {
				displayY = displayHeight - 1
			}

			gray := uint8(RGBToGraycale(uint32(r), uint32(g), uint32(b)))
			grayScale[displayY][x] = gray
			rgb[displayY][x] = Rgb{uint32(r), uint32(g), uint32(b), 255}

			i += 3
		}
	}

	if opts.Compression != 0 {
		grayScale = CompressGrayScale(grayScale, opts)
		rgb = CompressRgb(rgb, opts)
		displayHeight = len(grayScale)
		displayWidth = len(grayScale[0])
	}

	res := Ascii_t{
		AsciiChars: make([][]rune, displayHeight),
		RgbColors:  make([][]Rgb, displayHeight),
	}
	for y := 0; y < displayHeight; y++ {
		res.AsciiChars[y] = make([]rune, displayWidth)
		res.RgbColors[y] = make([]Rgb, displayWidth)
		for x := 0; x < displayWidth; x++ {
			res.AsciiChars[y][x] = PixelToChar(grayScale[y][x])
			res.RgbColors[y][x] = rgb[y][x]
		}
	}
	return res
}

func ImageToAscii(img image.Image,opts Options) (Ascii_t,error) {
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()

	grayScale := ImageToGrayScale(img,opts)
	rgbScale := ImageToRgb(img,opts)
	if len(grayScale) == 0 {
		return Ascii_t{},errors.New("Empty image")
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
	return res,nil
}
func intToBytes(i int) []byte {
	if i == 0 {
		return []byte{'0'}
	}
	var buf [20]byte 
	n := len(buf)
	for i > 0 {
		n--
		buf[n] = byte('0' + i%10)
		i /= 10
	}
	return buf[n:]
}
func PrintAsciiSlow(ascii Ascii_t,opts Options) {
	for y := 0; y < len(ascii.AsciiChars); y++ {
		for x := 0; x < len(ascii.AsciiChars[y]); x++ {
			ch := ascii.AsciiChars[y][x]
			color := ascii.RgbColors[y][x]

			if opts.UseColor {
				fmt.Printf("\x1b[38;2;%d;%d;%dm%s",
				color.R, color.G, color.B, string(ch))
			}else {
				fmt.Printf("%c",ch)
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\x1b[0m")
}
func PrintAsciiImage(ascii Ascii_t, opts Options) {
	os.Stdout.WriteString("\033[?25l")
	defer os.Stdout.WriteString("\033[?25h")

	if len(prevChars) != len(ascii.AsciiChars) || len(prevChars[0]) != len(ascii.AsciiChars[0]) {
		prevChars = make([][]rune, len(ascii.AsciiChars))
		prevColors = make([][]Rgb, len(ascii.RgbColors))
		for y := range ascii.AsciiChars {
			prevChars[y] = make([]rune, len(ascii.AsciiChars[y]))
			prevColors[y] = make([]Rgb, len(ascii.RgbColors[y]))
		}
		if opts.ClearScreen {
			os.Stdout.WriteString("\033[2J\033[H") 
		}
	}

	asciiBuffer = asciiBuffer[:0]

	for y := 0; y < len(ascii.AsciiChars); y++ {
		for x := 0; x < len(ascii.AsciiChars[y]); x++ {
			char := ascii.AsciiChars[y][x]
			color := ascii.RgbColors[y][x]

			if prevChars[y][x] == char &&
				(!opts.UseColor || (prevColors[y][x] == color)) {
				continue
			}

			asciiBuffer = append(asciiBuffer, "\033["...)
			asciiBuffer = append(asciiBuffer, intToBytes(y+1)...)
			asciiBuffer = append(asciiBuffer, ';')
			asciiBuffer = append(asciiBuffer, intToBytes(x+1)...)
			asciiBuffer = append(asciiBuffer, 'H')

			if opts.UseColor {
				asciiBuffer = append(asciiBuffer, "\x1b[38;2;"...)
				asciiBuffer = append(asciiBuffer, intToBytes(int(color.R))...)
				asciiBuffer = append(asciiBuffer, ';')
				asciiBuffer = append(asciiBuffer, intToBytes(int(color.G))...)
				asciiBuffer = append(asciiBuffer, ';')
				asciiBuffer = append(asciiBuffer, intToBytes(int(color.B))...)
				asciiBuffer = append(asciiBuffer, 'm')
			}

			asciiBuffer = append(asciiBuffer, byte(char))

			if opts.UseColor {
				asciiBuffer = append(asciiBuffer, "\x1b[0m"...)
			}

			prevChars[y][x] = char
			if opts.UseColor {
				prevColors[y][x] = color
			}
		}
	}

	if len(asciiBuffer) > 0 {
		os.Stdout.Write(asciiBuffer)
		os.Stdout.Sync()
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

func AsciiToImage(ascii Ascii_t, opts Options,output string) *image.RGBA{
	if len(ascii.AsciiChars) == 0 {
		return nil
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
	if output != "" {
		 file, err := os.Create(output)
		
		 if err != nil {
			panic("Could not create image")
		 }
		 defer file.Close()
		 png.Encode(file,img)
	}
	return img
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
			width, height := GetTermBounds()	
			charAspect := 2.0
			opts.Height = int(float64(height) * charAspect) -2
			opts.Width = width
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
func GetTermBounds() (int,int){
		fd := int(os.Stdout.Fd())
		if term.IsTerminal(fd) {
			width, height, err := term.GetSize(fd)	
			if err != nil {
				panic("Could not get terminal size")
			}
			return width,height
		}
	return 0,0
}

func GetFileExtension(filename string) string{
	str := strings.Split(filename, ".")
	if len(str) != 0 {
		return str[len(str)-1]
	}
	return ""
}

func ResizeImg(img image.Image,opts Options) image.Image {
	if opts.FitTerminal {
			width, height:= GetTermBounds()
			charAspect := 2.0
			opts.Height = int(float64(height) * charAspect) -2
			opts.Width = width 
	}
	dst := image.NewRGBA(image.Rect(0,0,opts.Width,opts.Height))
	draw.BiLinear.Scale(dst,dst.Bounds(),img,img.Bounds(),draw.Over,nil)
	return dst
}
func ResizeRgba(img *image.RGBA,opts Options) *image.RGBA {
	if opts.FitTerminal {
			width, height:= GetTermBounds()
			charAspect := 2.0
			opts.Height = int(float64(height) * charAspect) -2
			opts.Width = width 
	}
	dst := image.NewRGBA(image.Rect(0,0,opts.Width,opts.Height))
	draw.BiLinear.Scale(dst,dst.Bounds(),img,img.Bounds(),draw.Over,nil)
	return dst
}

func AsciiToGif(imgs []Ascii_t,opts Options,delays []int,disposal []byte,palets []color.Palette,output string){
	anim := gif.GIF{
		LoopCount: 0,
	}
	file, err := os.Create(output)

	if err != nil {
		panic("Could not create gif")
	}
	defer file.Close()
	frames := make([]*image.Paletted,len(imgs))
	if opts.ShowProgress {
		PrintProgress(0,len(imgs))
	}
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
	

	width,_ := GetTermBounds()
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

func GifToAscii(g *gif.GIF, opts Options) ([]Ascii_t,[]color.Palette,error){
	 palets := make([]color.Palette, len(g.Image))
    gifImages := make([]Ascii_t, len(g.Image))

    bounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
    canvas := image.NewRGBA(bounds)
    prevCanvas := image.NewRGBA(bounds)


    for i, img := range g.Image {
		bgColor := image.NewUniform(g.Image[i].Palette[g.BackgroundIndex])
        if i > 0 {
            switch g.Disposal[i-1] {
            case gif.DisposalBackground:
                draw.Draw(canvas, g.Image[i-1].Bounds(), bgColor, image.Point{}, draw.Src)
            case gif.DisposalPrevious:
                draw.Draw(canvas, bounds, prevCanvas, image.Point{}, draw.Src)
            }
        }

        if g.Disposal[i] == gif.DisposalPrevious {
            draw.Draw(prevCanvas, bounds, canvas, image.Point{}, draw.Src)
        }

        draw.Draw(canvas, img.Bounds(), img, image.Point{}, draw.Over)

        resized := ResizeRgba(canvas, opts)
        ascii, err := ImageToAscii(resized, opts)
        if err != nil {
            return nil, nil, err
        }

        gifImages[i] = ascii
        palets[i] = img.Palette
    }
    return gifImages, palets, nil
}
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

func RgbToImage(frame []byte,width int, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0,0,width,height))	
	i := 0
	for y := 0; y< height; y++ {
		for x := 0; x< width; x++ {
			r := frame[i]
			g := frame[i+1]
			b := frame[i+2]
			img.Set(x,y,color.RGBA{r,g,b,255})
			i = i+3
		}
	}
	return img
}
var prevFrame []byte
func ImageToRgbBytes(frame image.Image) []byte {
    bounds := frame.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()
    frameSize := width * height * 3

    buf := make([]byte, frameSize)

    i := 0
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r, g, b, _ := frame.At(x, y).RGBA()
            r8 := uint8(r >> 8)
            g8 := uint8(g >> 8)
            b8 := uint8(b >> 8)

            if prevFrame != nil &&
                prevFrame[i] == r8 &&
                prevFrame[i+1] == g8 &&
                prevFrame[i+2] == b8 {
                buf[i] = prevFrame[i]
                buf[i+1] = prevFrame[i+1]
                buf[i+2] = prevFrame[i+2]
            } else {
                buf[i] = r8
                buf[i+1] = g8
                buf[i+2] = b8
            }
            i += 3
        }
    }

    prevFrame = buf

    return buf
}
func AsciiToRgbBytes(frame Ascii_t) []byte {
	var width,height = 0,0
	height = len(frame.RgbColors)
	if height != 0 {
		width = len(frame.RgbColors[0])
	}
	if width == 0 {
		return nil
	}
	frameSize := width * height * 3
	buf := make([]byte, frameSize)
	i := 0
	for y := 0;y < height ; y ++ {
		for x := 0;x < width ; x ++ {
			rgb := frame.RgbColors[y][x]
			r ,g,b,_:= rgb.GetValues()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			buf[i] = r8
			buf[i+1] = g8
			buf[i+2] = b8
			i = i+3
		}
	}
	return buf
}











