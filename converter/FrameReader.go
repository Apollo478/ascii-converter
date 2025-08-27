package converter

import (
	// "bytes"
	"fmt"
	// "image"
	// "time"
	"io"

	//"image"
	"os/exec"
)
type FrameReader struct {
	stdout io.Reader
	buf []byte
	cmd *exec.Cmd
	frames chan []byte
	active bool
}
func NewFrameReader(opts Options) (*FrameReader,error) {
	width,height := opts.Width,opts.Height
	cmd := exec.Command("ffmpeg",
		"-f", "v4l2",           
		"-i", "/dev/video0",    
		"-framerate", "30",
		"-vf", fmt.Sprintf("scale=%d:%d",width,height), 
		"-pix_fmt", "rgb24",    
		"-f", "rawvideo",      
		"pipe:1")
		stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil,err
	}
	frames := make(chan []byte)	
	frameSize := width * height * 3	
	fmt.Println(width,height)
	frameReader := FrameReader {
		stdout: stdout,
		buf: make([]byte,frameSize),
		cmd : cmd,
		frames: frames,
		active: true,
	}
	return &frameReader,nil
}

func (fr *FrameReader) Frames(skip int) (<-chan []byte, error) {
	frameCount := 0 
	var error error = nil
	go func() {
		for {
			_,err := io.ReadFull(fr.stdout,fr.buf)
			if err != nil {
				error = err
				return 
			}
			frameCount++
			if frameCount % (skip+1) == 0 {
				continue
			}
			tmp := make([]byte, len(fr.buf))
			copy(tmp,fr.buf)
			if fr.frames != nil && fr.active  {
				fr.frames <- tmp
			}
		}
	}()
	return fr.frames,error
}

func (fr *FrameReader) Stop() error {
	fr.active = false
	close(fr.frames)
	return fr.cmd.Wait()
}











