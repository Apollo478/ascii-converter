package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"time"

	"github.com/Apollo478/ascii-converter/converter"
)

func main(){
	var chosenFile string = ""
	var chosenOption int = 0
	options := converter.Options{
		UseColor: true,
		UseAlpha: true,
		AspectRatio: 0.5,
		ClearScreen: true,
		BlendPrev: true,
		Parallel: true,
		Compression: 1,
		// Width: 300,
		// Height:200,
		FitTerminal: true,
	}
	fmt.Print("Convert image, gif, camera? (1 for image, 2 for gif,3 for camera): ");
	fmt.Scanf("%d",&chosenOption)
	converter.RevRamp = converter.SimpleRamp
	if chosenOption == 3 {
		converter.CameraToAscii(options,0)
		os.Exit(0)
	}
	fmt.Print("choose the file you wish to convert : ");
	fmt.Scanf("%s",&chosenFile)
	file,err := os.Open(chosenFile)
	if err != nil {
		fmt.Println("Could not open file")
		os.Exit(1)
	}
	defer file.Close()
	fmt.Println(file.Name())
	if options.Invert {
		converter.RevRamp = converter.ReverseRamp(converter.RevRamp)
	}
	 if chosenOption == 1 {
	
	 	img,_,err := image.Decode(file)
	 	if err != nil {
	 		panic("Could not decode image")
	 	}
	 	if options.Height == 0 {
	 		options.Height = img.Bounds().Dy()
	 	}
	 	if options.Width == 0 {
	 		options.Width = img.Bounds().Dx()
	 	}
	 	fmt.Print("choose your compression factor :")
	 	fmt.Scanf("%d",&options.Compression)
	 	options.Height = options.Height / options.Compression
	 	options.Width = options.Width / options.Compression
	 	img = converter.ResizeImg(img,options)
	
	 	fmt.Println("conversion")
	 	ascii := converter.ImageToAscii(img,options,nil)
	 	converter.PrintAsciiImage(ascii,options)
	 	fmt.Println("done conversion")
	 	converter.AsciiToImage(ascii,options,"")
	 	fmt.Println("imaged")
	 } else if chosenOption == 2 {
	
	 	g,err := gif.DecodeAll(file)
	 	if err != nil {
	 		panic(" Could not decode gif \n" + err.Error() )
	 	}
	
	 	palets  := make([]color.Palette,len(g.Image))
	
	 	fmt.Print("choose your compression factor :")
	 	fmt.Scanf("%d",&options.Compression)
	 	gifImages := make([]converter.Ascii_t, len(g.Image))
	 	fmt.Println(len(gifImages))
	 	fmt.Println(time.Now())
	 	if options.Height == 0 {
	 		options.Height = g.Config.Height
	 	}
	 	if options.Width == 0 {
	 		options.Width = g.Config.Width
	 	}
	 	options.Height = options.Height / options.Compression
	 	options.Width = options.Width / options.Compression
	 	for i, img := range g.Image {
	 		var frameToPass image.Image
	 		if i == 0 {
	 			frameToPass = nil 
	 		} else if g.Disposal[i-1] == gif.DisposalNone || g.Disposal[i-1] == gif.DisposalPrevious {
	 			frameToPass = g.Image[i-1]
	 		}
	 		 img = converter.ResizePaletted(img,options)
	 		 palets[i] = img.Palette
	
	 		var ascii converter.Ascii_t
	 			ascii = converter.ImageToAscii(img, options,frameToPass)
	
	 		if len(ascii.AsciiChars)!= 0{
	 			gifImages[i] = ascii
	 		}
	 	}
	 	fmt.Println(time.Now())
			// go func(){
			// 	converter.AsciiToGifSlow(gifImages,options,g.Delay,g.Disposal,palets)
			// }()
				converter.PrintAsciiGif(gifImages,options,g.Delay)
	 	fmt.Println(time.Now())
	} else if chosenOption == 4 {
		asciis,err := converter.Mp4ToAscii(options,chosenFile)
		if err != nil {
			fmt.Println("Could not convert mp4 file"+err.Error())
		}
		converter.SaveAsciiToMp4(asciis,options)
	}


	 

}
