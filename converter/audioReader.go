package converter

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
)

type AudioReader struct {
	stdout     io.ReadCloser
	buf        []byte
	cmd        *exec.Cmd
	chunkSize  int
	sampleRate int
	channels   int
	active     bool
}
func NewAudioReader(input string, sampleRate, channels, chunkSize int) (*AudioReader, error) {
	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-f", "s16le",                  
		"-ac", fmt.Sprintf("%d", channels), 
		"-ar", fmt.Sprintf("%d", sampleRate), 
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sampleSize := 2 * channels
	bufSize := chunkSize * sampleSize

	audioReader := AudioReader{
		stdout:     stdout,
		buf:        make([]byte, bufSize),
		cmd:        cmd,
		chunkSize:  chunkSize,
		sampleRate: sampleRate,
		channels:   channels,
		active:     true,
	}

	return &audioReader, nil
}

func (a *AudioReader) ReadChunk() ([]int16, error) {
	n, err := io.ReadFull(a.stdout, a.buf)
	if err != nil {
		return nil, err
	}
	samples := make([]int16, n/2) 
	for i := 0; i < len(samples); i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(a.buf[i*2:]))
	}

	return samples, nil
}
