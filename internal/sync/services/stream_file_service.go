package services

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
)

const (
	defaultFileStreamFPS = 20
	defaultFFmpegPath    = "ffmpeg"
)

type StreamFileServiceImpl struct {
	path       string
	ctx        context.Context
	ffmpegPath string
	fps        int
}

func NewStreamFileServiceImpl(path string, ctx context.Context) domain.StreamService {
	return &StreamFileServiceImpl{
		path:       path,
		ctx:        ctx,
		ffmpegPath: getEnv("FFMPEG_PATH", defaultFFmpegPath),
		fps:        getEnvInt("FILE_STREAM_FPS", defaultFileStreamFPS),
	}
}

func (s *StreamFileServiceImpl) DisplayStream() (<-chan domain.Frame, error) {
	cmd := exec.CommandContext(
		s.ctx,
		s.ffmpegPath,
		"-hide_banner",
		"-loglevel", "error",
		"-i", s.path,
		"-vf", fmt.Sprintf("fps=%d", s.fps),
		"-f", "image2pipe",
		"-vcodec", "png",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("failed to start ffmpeg process: %v", err)
		return nil, err
	}

	ch := make(chan domain.Frame, 1)
	log.Println("started ffmpeg process for file stream")
	go func() {
		defer close(ch)
		defer func() {
			// Ensure process is reaped. If ctx was cancelled, this returns quickly.
			_ = cmd.Wait()
		}()

		// Drain stderr in background so ffmpeg can't block on a full pipe.
		go io.Copy(io.Discard, stderr)

		r := bufio.NewReader(stdout)
		for {
			// Each iteration: read one PNG frame from the stream.
			img, derr := decodeNextPNG(r)
			if derr == io.EOF {
				return
			}
			if derr != nil {
				// If decode fails mid-stream, stop; could also choose to continue.
				return
			}

			frame := domain.Frame{
				Image:     img,
				Timestamp: time.Now(),
			}

			select {
			case ch <- frame:
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// decodeNextPNG reads exactly one PNG image from a stream that contains
// concatenated PNG images back-to-back (ffmpeg image2pipe).
func decodeNextPNG(r *bufio.Reader) (image.Image, error) {
	// PNG files end with the 12-byte IEND chunk:
	// 00 00 00 00 49 45 4E 44 AE 42 60 82
	const iend = "\x00\x00\x00\x00IEND\xaeB`\x82"

	var buf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			if buf.Len() == 0 {
				return nil, err // likely EOF
			}
			return nil, err
		}
		_ = buf.WriteByte(b)

		if buf.Len() >= len(iend) {
			tail := buf.Bytes()[buf.Len()-len(iend):]
			if bytes.Equal(tail, []byte(iend)) {
				break
			}
		}
	}

	return png.Decode(bytes.NewReader(buf.Bytes()))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return defaultValue
	}

	return parsed
}
