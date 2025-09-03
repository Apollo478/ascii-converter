package converter

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func rosePineGradient(scale float32) Rgb {
    if scale < 0.5 {
        t := scale / 0.5
        return lerpColor(Rgb{156,207,216,255}, Rgb{163,190,140,255}, t)
    } else {
        t := (scale-0.5)/0.5
        return lerpColor(Rgb{163,190,140,255}, Rgb{235,188,186,255}, t)
    }
}

func lerpColor(a, b Rgb, t float32) Rgb {
    return Rgb{
        R: uint32(float32(a.R)*(1-t) + float32(b.R)*t),
        G: uint32(float32(a.G)*(1-t) + float32(b.G)*t),
        B: uint32(float32(a.B)*(1-t) + float32(b.B)*t),
    }
}
func samplesToAscii2D(samples [][]int16, width, height int) Ascii_t {
	var ascii Ascii_t
    ascii.AsciiChars = make([][]rune, height)
    ascii.RgbColors = make([][]Rgb, height)

    for i := range ascii.AsciiChars {
        ascii.AsciiChars[i] = make([]rune, width)
        ascii.RgbColors[i] = make([]Rgb, width)
        for j := 0; j < width; j++ {
            ascii.AsciiChars[i][j] = ' '
        }
    }

    if len(samples) < 2 {
        return ascii 
    }

    left := samples[0]
    right := samples[1]
    n := min(len(left), len(right))

    for i := 0; i < n; i++ {
		x := int((float64(left[i]) / 32768.0) * float64(width/2)) + width/2
		y := height/2 - int((float64(right[i]) / 32768.0) * float64(height/2))
        if x >= 0 && x < width && y >= 0 && y < height {
            scale := float32(i) / float32(n)
            charIndex := int(scale * float32(len(RevRamp)-1))

            ascii.AsciiChars[y][x] = rune(RevRamp[charIndex])
            ascii.RgbColors[y][x] = rosePineGradient(scale)
        }
    }

    return ascii
}
func samplesToAscii(samples []int16, width, height int) Ascii_t {
	var ascii Ascii_t
	ascii.AsciiChars = make([][]rune, height)
	ascii.RgbColors = make([][]Rgb, height)

	for i := range ascii.AsciiChars {
		ascii.AsciiChars[i] = make([]rune, width)
		ascii.RgbColors[i] = make([]Rgb, width)
		for j := 0; j < width; j++ {
			ascii.AsciiChars[i][j] = ' '
		}
	}

	mid := height / 2
	for x := 0; x < width && x < len(samples); x++ {
		amp := int(float64(samples[x]) / 32768.0 * float64(height/2))
		y := mid - amp
		if y >= 0 && y < height {
			if y < mid {
				barHeight := mid - y
				for i, row := 0, y; row <= mid; row, i = row+1, i+1 {
					scale := float32(i) / float32(barHeight) 
					charIndex := int(scale * float32(len(RevRamp)-1))
					ascii.AsciiChars[row][x] =rune(RevRamp[charIndex] )
					ascii.RgbColors[row][x] = rosePineGradient(scale)
				}
			} else {
				barHeight := y - mid
				if barHeight == 0 {
					barHeight = 1
				}
				for i, row := 0, mid; row <= y; row, i = row+1, i+1 {
					scale := float32(i) / float32(barHeight) 
					charIndex := int((1-scale) * float32(len(RevRamp)-1))
					ascii.AsciiChars[row][x] =rune(RevRamp[charIndex] )
					ascii.RgbColors[row][x] = rosePineGradient(scale)
				}

			}
		}
	}

	return ascii
}
func samplesToBars(samples []int16, width, height int) Ascii_t {
    step := int(math.Max(1, float64(len(samples))/float64(width)))
    bars := make([]float64, width)

    for i := 0; i < width; i++ {
        sum := 0.0
        count := 0
        for j := i * step; j < (i+1)*step && j < len(samples); j++ {
            sum += math.Abs(float64(samples[j]))
            count++
        }
        if count > 0 {
            bars[i] = sum / float64(count)
        }
    }
    maxAmp := 0.0
    for _, v := range bars {
        if v > maxAmp {
            maxAmp = v
        }
    }
    if maxAmp == 0 {
        maxAmp = 1
    }
    chars :=RevRamp 

	var ascii Ascii_t
	ascii.AsciiChars = make([][]rune, height)
	ascii.RgbColors = make([][]Rgb, height)
    for y := 0; y < height; y++ {
        ascii.AsciiChars[y] = []rune(strings.Repeat(" ", width))
        ascii.RgbColors[y] =make([]Rgb,width) 
    }

    for x, v := range bars {
        barHeight := int((v / maxAmp) * float64(height))
        for y := 0; y < barHeight && y < height; y++ {

            level := (y * len(chars)) / height
            ascii.AsciiChars[height-1-y][x] = rune(chars[level])
            ascii.RgbColors[height-1-y][x] =rosePineAmpColor(barHeight,height) 
        }
    }
	return ascii
}

func samplesToSpectrum(samples []int16,width int,height int) Ascii_t  {
	sampleLenght := len(samples)	
	buf := make([]float64,sampleLenght)
	for i,sample := range samples {
		buf[i] = float64(sample)
	}

	step := int(math.Max(1,float64(sampleLenght/2)/float64(width)))
	spectrum := make([]float64,width)

	for x := 0; x < width; x ++{
		re,im := 0.0,0.0

		for n:=0; n < sampleLenght; n ++ {
			angle := -2.0 * math.Pi * float64(x * step) * float64(n) /float64(sampleLenght)

			re += buf[n] * math.Cos(angle)
			im += buf[n] * math.Sin(angle)
		}
		spectrum[x] = math.Sqrt(re*re + im*im)
	}
	maxMag := 0.0
	for _,m:= range spectrum {
		if m > maxMag {
			maxMag = m
		}
	}
    if maxMag == 0 {
        maxMag = 1
    }
	var ascii Ascii_t
	ascii.AsciiChars = make([][]rune, height)
	ascii.RgbColors = make([][]Rgb, height)
    for y := 0; y < height; y++ {
        ascii.AsciiChars[y] = []rune(strings.Repeat(" ", width))
        ascii.RgbColors[y] =make([]Rgb,width) 
    }

	for x,mag := range spectrum {
		barHeight  := int((mag/maxMag) * float64(height))
		for y := 0; y < barHeight; y++ {
			level := (y*len(RevRamp)) / height
			ascii.AsciiChars[height-1-y][x] = rune(RevRamp[level])
            ascii.RgbColors[height-1-y][x] =rosePineAmpColor(barHeight,height) 
		}
	}

	return ascii
}
func AudioToAscii(input string,opts Options) {
	reader, err := NewAudioReader(input, 44100, 2, 2084)
    if err != nil {
        panic(err)
    }

    frameDuration := time.Duration(reader.chunkSize) * time.Second / time.Duration(reader.sampleRate)

    for {
        samples, err := reader.ReadChunk2d()
        if err != nil {
            break
        }

        ascii := samplesToAscii2D(samples, opts.Width, opts.Height)
        PrintAsciiImage(ascii, opts)

        time.Sleep(frameDuration)
    }
}
func PrintAudio(s []string) {
	for _, line := range s {
			fmt.Println(line)
		}
		fmt.Print("\033[H")
}
func rosePineAmpColor(amp int, height int) Rgb {
    strength := float64(math.Abs(float64(amp)) / float64(height/2))
    switch {
    case strength < 0.20:
        return Rgb{224, 222, 244,255}
    case strength < 0.40:
        return Rgb{156, 207, 216,255}
    case strength < 0.60:
        return Rgb{163, 190, 140,255}
    case strength < 0.80:
        return Rgb{197, 199, 198,255}

    default:
        return Rgb{235, 188, 186,255}
    }
}
