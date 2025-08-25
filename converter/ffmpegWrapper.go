package converter

import (
	// "bytes"
	"fmt"
	// "time"
	"io"

	//"image"
	"os/exec"
)
type FrameReader struct {
	stdout io.Reader
	buf []byte
}
func NewFrameReader(stdout io.Reader,width int, height int) *FrameReader {
	frameSize := width * height * 3	
	return  &FrameReader{
		stdout: stdout,
		buf: make([]byte,frameSize),
	}
}

func (fr *FrameReader) ReadNextFrame() <-chan []byte {
	frames := make(chan []byte)	
	go func(){
		defer close(frames)
		for {
			_,err := io.ReadFull(fr.stdout,fr.buf)
			if err != nil {
				panic("Could not read frame "+err.Error())
			}
			tmp := make([]byte, len(fr.buf))
			copy(tmp,fr.buf)
			frames <- tmp
		}
	}()
	return frames
}

func ReadFrames(opts Options) io.ReadCloser{
	width,height := 0,0
	if opts.FitTerminal {
		width,height = GetTermBounds()
		height = height * 2 -1
	}
	// frameSize := width * height *3
	cmd := exec.Command("ffmpeg",
		"-f", "v4l2",           
		"-i", "/dev/video0",    
		"-framerate", "30",
		"-vf", fmt.Sprintf("scale=%d.%d",width,height), 
		"-pix_fmt", "rgb24",    
		"-f", "rawvideo",      
		"pipe:1")
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		panic("Could not start ffmpeg " + err.Error())
	}
	return stdout
}
