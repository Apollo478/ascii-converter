package cli

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"os"

	"github.com/Apollo478/ascii-converter/converter"
)
func Execute() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	opts := converter.Options{
		BlendPrev: true,
		UseAlpha: true,
		Compression: 1,
		ShowProgress: false,
	}

	converter.RevRamp = converter.SimpleRamp
	switch os.Args[1] {
	case "convert":
		convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
		input := convertCmd.String("input", "", "Input file")
		convertCmd.StringVar(input, "i", "", "Alias for --input")

		output := convertCmd.String("output", "", "output file")
		convertCmd.StringVar(output, "o", "", "Alias for --output")

		width := convertCmd.Int("width", converter.DefaultWidth, "ASCII width")
		convertCmd.IntVar(width, "w", converter.DefaultWidth, "Alias for --width")

		height := convertCmd.Int("height", converter.DefaultHeight, "ASCII height")
		convertCmd.IntVar(height, "h", converter.DefaultHeight, "Alias for --height")


		fitTerminal := convertCmd.Bool("fit-terminal", false, "Fit ASCII to terminal size")
		convertCmd.BoolVar(fitTerminal, "f", false, "Alias for --fit-terminal")

		color := convertCmd.Bool("color", false, "Enable colored ASCII")
		convertCmd.BoolVar(color, "C", false, "Alias for --color")

		parallel := convertCmd.Bool("parallel", false, "Process frames in parallel")
		convertCmd.BoolVar(parallel, "p", false, "Alias for --parallel")

		preview := convertCmd.Bool("preview", false, "Preview ascii while saving")
		convertCmd.BoolVar(preview, "P", false, "Alias for --preview")

		clearScreen := convertCmd.Bool("clear-screen", true, "Clear screen before printing frames")
		convertCmd.BoolVar(clearScreen, "s", true, "Alias for --clear-screen")

		inverse := convertCmd.Bool("invert", false, "Invert the ASCII scale")
		convertCmd.BoolVar(inverse, "I", false, "Alias for --invert")

		aspectRatio := convertCmd.Float64("aspect-ratio", 0.5, "Set aspect ratio of ASCII’s Y axis")
		convertCmd.Float64Var(aspectRatio, "a", 0.5, "Alias for --aspect-ratio")
		convertCmd.Parse(os.Args[2:])
		opts.Width = *width
		opts.Height = *height
		opts.FitTerminal = *fitTerminal
		opts.AspectRatio = *aspectRatio
		opts.UseColor = *color
		opts.AspectRatio = *aspectRatio
		opts.ClearScreen = *clearScreen
		opts.Parallel = *parallel
		opts.Invert = *inverse
		opts.Preview = *preview
		if *inverse {
			converter.RevRamp = converter.ReverseRamp(converter.RevRamp)
		}
		if opts.FitTerminal {
			opts.Width,opts.Height = converter.GetTermBounds()
		}
		opts.Height = opts.Height / opts.Compression
		opts.Width = opts.Width / opts.Compression
		if *input == "" {
			fmt.Println("Error: -input is required")
			convertCmd.Usage()
			os.Exit(1)
		}
		if *output == "" {
			fmt.Println("Error: -output is required")
			convertCmd.Usage()
			os.Exit(1)
		}
		extension := converter.GetFileExtension(*input)
		if extension == "mp4" || extension == "mov" || extension == "avi" {
			asciis, err := converter.VideoToAscii(opts,*input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			go converter.SaveAsciiToVideo(asciis,opts,*output)
			if *preview {
				converter.PrintAsciiVideo(asciis,opts)
			}
		}
		if extension =="png" || extension =="jpg" || extension =="jpeg" {
			file,err := os.Open(*input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			img,_,err := image.Decode(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			img = converter.ResizeImg(img,opts)
			ascii,err := converter.ImageToAscii(img,opts)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if *preview {
				converter.PrintAsciiImage(ascii,opts)	
			}
			converter.AsciiToImage(ascii,opts,*output)
		}
		if extension == "gif" {
			file,err := os.Open(*input)
			if err != nil {
				fmt.Println("Could not open file")
				os.Exit(1)
			}
			defer file.Close()

			g,err := gif.DecodeAll(file)
			if err != nil {
				fmt.Println("Could not decode gif " + err.Error())
				os.Exit(1)
			}
			asciis,palets,err := converter.GifToAscii(g,opts)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			
			 converter.AsciiToGif(asciis,opts,g.Delay,g.Disposal,palets,*output)
			if *preview {
				 converter.PrintAsciiGif(asciis,opts,g.Delay)
			}
		}

	case "preview":
		previewCmd := flag.NewFlagSet("preview", flag.ExitOnError)
		input := previewCmd.String("input", "", "Input file")
		previewCmd.StringVar(input, "i", "", "Alias for --input")

		width := previewCmd.Int("width", converter.DefaultWidth, "ASCII width")
		previewCmd.IntVar(width, "w", converter.DefaultWidth, "Alias for --width")

		height := previewCmd.Int("height", converter.DefaultHeight, "ASCII height")
		previewCmd.IntVar(height, "h", converter.DefaultHeight, "Alias for --height")


		fitTerminal := previewCmd.Bool("fit-terminal", false, "Fit ASCII to terminal size")
		previewCmd.BoolVar(fitTerminal, "f", false, "Alias for --fit-terminal")

		color := previewCmd.Bool("color", false, "Enable colored ASCII")
		previewCmd.BoolVar(color, "C", false, "Alias for --color")

		parallel := previewCmd.Bool("parallel", false, "Process frames in parallel")
		previewCmd.BoolVar(parallel, "p", false, "Alias for --parallel")

		clearScreen := previewCmd.Bool("clear-screen", true, "Clear screen before printing frames")
		previewCmd.BoolVar(clearScreen, "s", true, "Alias for --clear-screen")

		inverse := previewCmd.Bool("invert", false, "Invert the ASCII scale")
		previewCmd.BoolVar(inverse, "I", false, "Alias for --invert")

		aspectRatio := previewCmd.Float64("aspect-ratio", 0.5, "Set aspect ratio of ASCII’s Y axis")
		previewCmd.Float64Var(aspectRatio, "a", 0.5, "Alias for --aspect-ratio")
		previewCmd.Parse(os.Args[2:])
		opts.Width = *width
		opts.Height = *height
		opts.FitTerminal = *fitTerminal
		opts.AspectRatio = *aspectRatio
		opts.UseColor = *color
		opts.AspectRatio = *aspectRatio
		opts.ClearScreen = *clearScreen
		opts.Parallel = *parallel
		opts.Invert = *inverse
		opts.Preview = true
		if *inverse {
			converter.RevRamp = converter.ReverseRamp(converter.RevRamp)
		}
		if opts.FitTerminal {
			opts.Width,opts.Height = converter.GetTermBounds()
		}
		opts.Height = opts.Height / opts.Compression
		opts.Width = opts.Width / opts.Compression
		if *input == "" {
			fmt.Println("Error: -input is required")
			previewCmd.Usage()
			os.Exit(1)
		}
		extension := converter.GetFileExtension(*input)
		opts.PreviewInPreview = true
		if extension == "mp4" || extension == "mov" || extension == "avi" {
			opts.Preview = true
			_, err := converter.VideoToAscii(opts,*input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// converter.PrintAsciiVideo(asciis,opts)
		}
		if extension =="png" || extension =="jpg" || extension =="jpeg" {
			file,err := os.Open(*input)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()
			img,_,err := image.Decode(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			img = converter.ResizeImg(img,opts)
			ascii,err := converter.ImageToAscii(img,opts)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			converter.PrintAsciiImage(ascii,opts)	
		}
		if extension == "gif" {
			file,err := os.Open(*input)
			if err != nil {
				fmt.Println("Could not open file")
				os.Exit(1)
			}
			defer file.Close()

			g,err := gif.DecodeAll(file)
			if err != nil {
				fmt.Println("Could not decode gif " + err.Error())
				os.Exit(1)
			}
			asciis,_,err := converter.GifToAscii(g,opts)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			converter.PrintAsciiGif(asciis,opts,g.Delay)
		}

	case "camera": {
		cameraCmd := flag.NewFlagSet("camera", flag.ExitOnError)

		output := cameraCmd.String("output", "", "output file")
		cameraCmd.StringVar(output, "o", "", "Alias for --output")

		width := cameraCmd.Int("width", converter.DefaultWidth, "ASCII width")
		cameraCmd.IntVar(width, "w", converter.DefaultWidth, "Alias for --width")

		height := cameraCmd.Int("height", converter.DefaultHeight, "ASCII height")
		cameraCmd.IntVar(height, "h", converter.DefaultHeight, "Alias for --height")


		fitTerminal := cameraCmd.Bool("fit-terminal", false, "Fit ASCII to terminal size")
		cameraCmd.BoolVar(fitTerminal, "f", false, "Alias for --fit-terminal")

		color := cameraCmd.Bool("color", false, "Enable colored ASCII")
		cameraCmd.BoolVar(color, "C", false, "Alias for --color")

		parallel := cameraCmd.Bool("parallel", false, "Process frames in parallel")
		cameraCmd.BoolVar(parallel, "p", false, "Alias for --parallel")


		clearScreen := cameraCmd.Bool("clear-screen", true, "Clear screen before printing frames")
		cameraCmd.BoolVar(clearScreen, "s", true, "Alias for --clear-screen")

		inverse := cameraCmd.Bool("invert", false, "Invert the ASCII scale")
		cameraCmd.BoolVar(inverse, "I", false, "Alias for --invert")

		aspectRatio := cameraCmd.Float64("aspect-ratio", 0.5, "Set aspect ratio of ASCII’s Y axis")
		cameraCmd.Float64Var(aspectRatio, "a", 0.5, "Alias for --aspect-ratio")

		preview := cameraCmd.Bool("preview", false, "Preview ascii while saving")
		cameraCmd.BoolVar(preview, "P", false, "Alias for --preview")

		cameraCmd.Parse(os.Args[2:])
		opts.Width = *width
		opts.Height = *height
		opts.FitTerminal = *fitTerminal
		opts.AspectRatio = *aspectRatio
		opts.UseColor = *color
		opts.AspectRatio = *aspectRatio
		opts.ClearScreen = *clearScreen
		opts.Parallel = *parallel
		opts.Invert = *inverse
		opts.Preview =*preview 
		if *inverse {
			converter.RevRamp = converter.ReverseRamp(converter.RevRamp)
		}
		if opts.FitTerminal {
			opts.Width,opts.Height = converter.GetTermBounds()
		}
		opts.Height = opts.Height / opts.Compression
		opts.Width = opts.Width / opts.Compression
		if opts.Parallel && opts.Preview {
			fmt.Println("Cant preview paralelled frames")
			os.Exit(1)
		}
		err := converter.CameraToAscii(opts,0,*output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	default:
		printUsage()
		os.Exit(1)
	}
}
func printUsage() {
	fmt.Println("Usage: ascii-cli <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  convert   Convert image/gif/video to ASCII")
	fmt.Println("  preview   Preview ASCII frames in terminal")
	fmt.Println("  camera   Preview/convert camera ASCII frames ")
}
