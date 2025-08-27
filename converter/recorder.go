package converter

import (
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
	cmd := exec.Command("ffmpeg",
        "-y",                
        "-f", "rawvideo",
        "-pix_fmt", "rgb24",
        "-s", fmt.Sprintf("%d:%d", width, height),
        "-r", "10",          
        "-i", "pipe:0",
        "-c:v", "libx264",
        "-pix_fmt", "rgb24",
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

func (r *Recorder) WriteFrame(buf []byte) error {
	if !r.active {
		return nil
	}
	total := 0
    for total < len(buf) {
        n, err := r.stdin.Write(buf[total:])
        if err != nil {
            return err
        }
        total += n
    }
	return nil
}
func (r *Recorder) Stop() error {
fmt.Println("here")
	r.active = false
	r.stdin.Close()
	return  r.cmd.Wait()
}




