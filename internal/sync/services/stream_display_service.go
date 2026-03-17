package services

import (
	"context"
	"log"
	"time"

	"github.com/kbinani/screenshot"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
)

const FPS = 20

type StreamDisplayServiceImpl struct {
	width int
	ctx   context.Context
}

func NewStreamDisplayService(width int, ctx context.Context) domain.StreamService {
	return &StreamDisplayServiceImpl{
		width: width,
		ctx:   ctx,
	}
}

func (s *StreamDisplayServiceImpl) DisplayStream() (<-chan domain.Frame, error) {
	ch := make(chan domain.Frame, 1)

	go func() {
		defer close(ch)
		delay := time.Second / time.Duration(FPS)

		for {
			frame, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(1))
			if err != nil {
				log.Println("capture error:", err)
				continue
			}

			ch <- domain.Frame{
				Image:     frame,
				Timestamp: time.Now(),
			}

			time.Sleep(delay)
		}
	}()

	return ch, nil
}
