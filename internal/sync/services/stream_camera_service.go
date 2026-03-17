package services

import (
	"context"
	"log"
	"time"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
	"gocv.io/x/gocv"
)

const defaultCameraStreamFPS = 10

type StreamCameraServiceImpl struct {
	ctx      context.Context
	deviceID int
	fps      int
}

func NewStreamCameraServiceImpl(ctx context.Context, deviceID int, fps int) domain.StreamService {
	if fps <= 0 {
		fps = defaultCameraStreamFPS
	}

	return &StreamCameraServiceImpl{
		ctx:      ctx,
		deviceID: deviceID,
		fps:      fps,
	}
}

func (s *StreamCameraServiceImpl) DisplayStream() (<-chan domain.Frame, error) {
	ch := make(chan domain.Frame, 1)

	webcam := s.openCamera()
	img := gocv.NewMat()

	ticker := time.NewTicker(time.Second / time.Duration(s.fps))

	go func() {
		defer func() {
			close(ch)
			img.Close()
			webcam.Close()
			ticker.Stop()
		}()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				if ok := webcam.Read(&img); !ok {
					log.Println("cannot read device")
					continue
				}

				if img.Empty() {
					log.Println("empty image")
					continue
				}

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

func (s *StreamCameraServiceImpl) openCamera() *gocv.VideoCapture {
	webcam, err := gocv.OpenVideoCapture(s.deviceID)
	if err != nil {
		panic(err)
	}

	webcam.Set(gocv.VideoCaptureFPS, float64(s.fps))

	fps := webcam.Get(gocv.VideoCaptureFPS)
	log.Println("Camera FPS:", fps)

	return webcam
}
