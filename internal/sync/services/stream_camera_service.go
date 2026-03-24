package services

import (
	"context"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/AlKulinski/lumigo/internal/commons/timer"
	framing "github.com/AlKulinski/lumigo/internal/framing/domain"
	"github.com/AlKulinski/lumigo/internal/sync/domain"
	"gocv.io/x/gocv"
)

const defaultCameraStreamFPS = 10

type StreamCameraServiceImpl struct {
	ctx            context.Context
	fps            int
	webcam         *gocv.VideoCapture
	framingService framing.FramingService
}

func NewStreamCameraServiceImpl(ctx context.Context, deviceID int, fps int, framingService framing.FramingService) domain.StreamService {
	if fps <= 0 {
		fps = defaultCameraStreamFPS
	}

	webcam := openCamera(deviceID, fps)

	return &StreamCameraServiceImpl{
		ctx:            ctx,
		fps:            fps,
		framingService: framingService,
		webcam:         webcam,
	}
}

func (s *StreamCameraServiceImpl) DisplayStream() (<-chan domain.Frame, error) {
	ch := make(chan domain.Frame, 1)

	img := gocv.NewMat()

	if ok := s.webcam.Read(&img); !ok {
		return nil, fmt.Errorf("cannot read device")
	}

	framing := image.Rect(0, 0, img.Cols(), img.Rows())
	ticker := time.NewTicker(time.Second / time.Duration(s.fps))
	framingTicker := timer.NewImmediateTicker(2 * time.Second)
	go func() {
		defer func() {
			close(ch)
			img.Close()
			s.webcam.Close()
			ticker.Stop()
			framingTicker.Stop()
		}()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-framingTicker.C:
				tvRect, found := s.framingService.Frame(img.Clone())
				if found {
					log.Printf("framing: %v", tvRect)
					framing = tvRect
				}
			case <-ticker.C:
				if ok := s.webcam.Read(&img); !ok {
					log.Println("cannot read device")
					continue
				}

				if img.Empty() {
					log.Println("empty image")
					continue
				}

				img := img.Region(framing)

				parsedImg, err := img.ToImage()
				if err != nil {
					log.Printf("cannot convert image: %v", err)
					continue
				}

				select {
				case ch <- domain.Frame{
					Image:     parsedImg,
					Timestamp: time.Now(),
				}:
				case <-s.ctx.Done():
					return
				}
			}
		}
	}()

	return ch, nil
}

func openCamera(deviceId int, fps int) *gocv.VideoCapture {
	webcam, err := gocv.OpenVideoCapture(deviceId)
	if err != nil {
		panic(err)
	}

	webcam.Set(gocv.VideoCaptureFPS, float64(fps))

	return webcam
}
