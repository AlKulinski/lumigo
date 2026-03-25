package usecases

import (
	"image"

	"github.com/AlKulinski/lumigo/internal/commons/quickselect"
	"github.com/AlKulinski/lumigo/internal/sync/domain"
	"github.com/AlKulinski/lumigo/internal/sync/services"
	"github.com/AlKulinski/lumigo/internal/sync/utils"
)

type SyncUsecase struct {
	syncService   services.SyncService
	streamService domain.StreamService
}

func NewSyncUsecase(syncService services.SyncService, streamService domain.StreamService) *SyncUsecase {
	return &SyncUsecase{syncService: syncService, streamService: streamService}
}

func (u *SyncUsecase) Execute() error {
	frames, err := u.streamService.DisplayStream()
	if err != nil {
		return err
	}

	for frame := range frames {
		lumas := calculateLumas(frame.Image)

		topLumas := pickSample(lumas)
		r, g, b, luma := utils.AverageColor(topLumas)

		u.syncService.Sync(domain.Color{
			R:    r,
			G:    g,
			B:    b,
			Luma: float64(luma),
		})
	}
	return nil
}

func calculateLumas(image image.Image) []domain.Pixel {
	rect := image.Bounds()
	yMax, xMax := rect.Max.Y, rect.Max.X
	lumas := make([]domain.Pixel, 0, xMax*yMax)

	for y := rect.Min.Y; y < yMax; y++ {
		for x := rect.Min.X; x < xMax; x++ {
			r, g, b, _ := image.At(x, y).RGBA()
			r8 := uint32(r >> 8)
			g8 := uint32(g >> 8)
			b8 := uint32(b >> 8)
			luma := utils.CalculateLuma(r8, g8, b8)

			lumas = append(lumas, domain.Pixel{
				X: x,
				Y: y,
				Color: domain.Color{
					R:    r8,
					G:    g8,
					B:    b8,
					Luma: float64(luma),
				},
			})
		}
	}

	return lumas
}

func pickSample(lumas []domain.Pixel) []domain.Pixel {
	return quickselect.TopN(lumas, len(lumas)/10, func(a domain.Pixel) int {
		return int(a.Color.Luma)
	})
}
