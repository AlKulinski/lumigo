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
			Luma: luma,
		})
	}
	return nil
}

func calculateLumas(img image.Image) []domain.Pixel {
	rect := img.Bounds()
	yMin, xMin := rect.Min.Y, rect.Min.X
	yMax, xMax := rect.Max.Y, rect.Max.X
	lumas := make([]domain.Pixel, 0, (xMax-xMin)*(yMax-yMin))

	appendPixel := func(x, y int, r, g, b uint32) {
		lumas = append(lumas, domain.Pixel{
			X: x,
			Y: y,
			Color: domain.Color{
				R:    r,
				G:    g,
				B:    b,
				Luma: utils.CalculateLuma(r, g, b),
			},
		})
	}

	switch m := img.(type) {
	case *image.RGBA:
		for y := yMin; y < yMax; y++ {
			i := m.PixOffset(xMin, y)
			for x := xMin; x < xMax; x++ {
				appendPixel(x, y, uint32(m.Pix[i]), uint32(m.Pix[i+1]), uint32(m.Pix[i+2]))
				i += 4
			}
		}
	case *image.NRGBA:
		for y := yMin; y < yMax; y++ {
			i := m.PixOffset(xMin, y)
			for x := xMin; x < xMax; x++ {
				appendPixel(x, y, uint32(m.Pix[i]), uint32(m.Pix[i+1]), uint32(m.Pix[i+2]))
				i += 4
			}
		}
	default:
		for y := yMin; y < yMax; y++ {
			for x := xMin; x < xMax; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				appendPixel(x, y, r>>8, g>>8, b>>8)
			}
		}
	}

	return lumas
}

func pickSample(lumas []domain.Pixel) []domain.Pixel {
	return quickselect.TopN(lumas, len(lumas)/10, func(a domain.Pixel) int {
		return int(a.Color.Luma)
	})
}
