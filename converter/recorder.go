package converter

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type Recorder struct {
	cmd *exec.Cmd
	stdin io.WriteCloser 
	active bool
}
var frameSize int
func NewRecorder(opts Options,output string) (*Recorder, error){
	 width,height := opts.Width,opts.Height
	if width %2 != 0 || height %2 != 0 {
		return nil,errors.New("Height or width are not divisible by 2")
	}
	cmd := exec.Command("ffmpeg",
        "-y",                
        "-f", "rawvideo",
        "-pix_fmt", "rgb24",
        "-s", fmt.Sprintf("%dx%d", width, height),
        "-r", "10",          
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-qp", "0",
        output,
    )
	frameSize = width * height*3
	stdin,err := cmd.StdinPipe()
	if err != nil {
		return nil,err
	}
	recorder := Recorder{
		cmd: cmd,
		stdin: stdin,
		active: true,
	}
	if err := cmd.Start(); err  != nil {
		return nil,err
	}
	return &recorder,nil
}
func PadVideo(input string) error{
	cmd := exec.Command("ffmpeg",
        "-i", input,
		"-vf","scale=iw*4:ih*4:flags=neighbor,pad=800:600:(ow-iw*4)/2:(oh-ih*4)/2",
        input,
    )
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (r *Recorder) WriteFrame(buf []byte) error {
	if !r.active {
		return nil
	}
	_, err := r.stdin.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
func (r *Recorder) Stop() error {
	r.active = false
	r.stdin.Close()
	return  r.cmd.Wait()
}




