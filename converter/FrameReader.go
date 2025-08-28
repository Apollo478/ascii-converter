package converter

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)
type FrameReader struct {
	stdout io.Reader
	buf []byte
	cmd *exec.Cmd
	frames chan []byte
	active bool
}
func NewCamReader(opts Options,input int) (*FrameReader,error) {
	width,height := opts.Width,opts.Height
	if width %2 != 0 || height %2 != 0 {
		return nil,errors.New("Height or width are not divisible by 2")
	}
	cmd := exec.Command("ffmpeg",
		"-f", "v4l2",           
		"-i", fmt.Sprintf("/dev/video%d",input),    
		"-r", "10",
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

func NewVideoReader(opts Options,input string) (*FrameReader,error) {
	width,height := opts.Width,opts.Height
	if width %2 != 0 || height %2 != 0 {
		return nil,errors.New("Height or width are not divisible by 2")
	}
	cmd := exec.Command("ffmpeg",
		"-i", input,    
		// "-framerate", "30",
		"-s", fmt.Sprintf("%dx%d",width,height), 
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
	go func() {
		defer close(fr.frames)
		for {

			_,err := io.ReadFull(fr.stdout,fr.buf)
			if err != nil {
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
	return fr.frames,nil
}

func (fr *FrameReader) Stop() error {
	fr.active = false
	close(fr.frames)
	return fr.cmd.Wait()
}











