package services

import (
	"container/ring"
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
	ctx                context.Context
	fps                int
	webcam             *gocv.VideoCapture
	framingService     framing.FramingService
	debugWindowService framing.DebugWindowService
}

func NewStreamCameraServiceImpl(ctx context.Context, deviceID int, fps int, framingService framing.FramingService, debugWindowService framing.DebugWindowService) domain.StreamService {
	if fps <= 0 {
		fps = defaultCameraStreamFPS
	}

	webcam := openCamera(deviceID, fps)

	return &StreamCameraServiceImpl{
		ctx:                ctx,
		fps:                fps,
		framingService:     framingService,
		debugWindowService: debugWindowService,
		webcam:             webcam,
	}
}

func (s *StreamCameraServiceImpl) averageFraming(f framing.Framing, r *ring.Ring) framing.Framing {
	r.Value = f
	r = r.Next()
	sum := image.Rect(0, 0, 0, 0)
	r.Do(func(value interface{}) {
		if rect, ok := value.(framing.Framing); ok {
			sum = sum.Union(rect)
		}
	})
	return sum
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
	framingRing := ring.New(10)
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
					framing = s.averageFraming(tvRect, framingRing)
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

				if s.debugWindowService != nil {
					go s.debugWindowService.Update(parsedImg)
				}

				ch <- domain.Frame{
					Image:     parsedImg,
					Timestamp: time.Now(),
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
