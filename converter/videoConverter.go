package converter

import (
	"errors"
	"runtime"
	"time"
)

func SaveAsciiToVideo(frames []Ascii_t, opts Options,output string) error  {
	recorder, err := NewRecorder(opts,output)
	if err != nil {
		return  err
	}
	for _,ascii := range  frames {
		img := AsciiToImage(ascii,opts,"")
		img = ResizeRgba(img,opts)
		if img != nil {
			bytes := ImageToRgbBytes(img)
			err := recorder.WriteFrame(bytes)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func VideoToAscii(opts Options, filename string) ([]Ascii_t,error){
	reader ,err:= NewVideoReader(opts,filename)
	var asciis []Ascii_t
	if err != nil {
		return nil,err
	}
	frames ,err := reader.Frames(1)
	if err != nil {
		return nil,err
	}
	for frame := range frames {
		ascii := RgbBufferToAscii(frame,opts)
		if opts.PreviewInPreview {
			PrintAsciiImage(ascii,opts)
			time.Sleep(100 * time.Millisecond)
		}
		asciis = append(asciis,ascii)
	}
	return asciis,nil
}
func PrintAsciiVideo(asciis []Ascii_t,opts Options) {
	for _,ascii := range asciis {
		PrintAsciiImage(ascii,opts)
		time.Sleep(33 * time.Millisecond)
	}
}
func CameraToAscii(opts Options,camera int,output string) error {
	
	camReader ,err:= NewCamReader(opts,camera)
	if err != nil {
		return errors.New("Error creating cam frame reader "+err.Error())
	}
	frame ,err:= camReader.Frames(1)
	var processed chan []byte
	processed = make(chan []byte,opts.Height * opts.Width * 3)
	if err != nil {
		return errors.New("Error reading frames "+err.Error())
	}
	if opts.Parallel {
		numWorkers := runtime.NumCPU()
		for i := 0; i < numWorkers; i++ {
			go func(){
				for frame := range  frame {
					ascii := RgbBufferToAscii(frame,opts)
					if opts.Preview {
						PrintAsciiImage(ascii,opts)
					}
					img := AsciiToImage(ascii,opts,"")
					img = ResizeRgba(img,opts)

					if img != nil {
						processed <- ImageToRgbBytes(img)
					}

				}
			}()
		}
	} else {
			go func(){
				for frame := range  frame {
					ascii := RgbBufferToAscii(frame,opts)
					if opts.Preview {
						PrintAsciiImage(ascii,opts)
					}
					img := AsciiToImage(ascii,opts,"")
					img = ResizeRgba(img,opts)

					if img != nil {
						processed <- ImageToRgbBytes(img)
					}

				}
			}()
	}
	if output != ""{
		recorder, err := NewRecorder(opts,output)
		if err != nil {
			return errors.New("Error creating recorder "+err.Error())
		}
		for buf := range processed {
			err := recorder.WriteFrame(buf)
			if err != nil {
				return errors.New("Could not record frame " + err.Error())
			}
			time.Sleep(100 * time.Millisecond)
		}
	} else {
		for range processed  {
			
		}
	}

	return nil
}
