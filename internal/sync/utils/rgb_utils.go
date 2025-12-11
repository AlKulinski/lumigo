package utils

import (
	"math"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
)

func AverageColor(img []domain.Pixel) (uint32, uint32, uint32, uint32) {
	var r, g, b, luma uint32
	count := uint32(len(img))

	for i := 0; i < int(count); i++ {
		r += uint32(img[i].Color.R)
		g += uint32(img[i].Color.G)
		b += uint32(img[i].Color.B)
		luma += uint32(img[i].Color.Luma)
	}

	if count == 0 {
		return 0, 0, 0, 0
	}

	return (r / count), (g / count), (b / count), (luma / count)
}

func CalculateLuma(r, g, b uint32) uint32 {
	return uint32((0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)) / 255.0)
}

func RGBToHSB(r, g, b uint32) (h, s, v float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min

	if delta == 0 {
		h = 0
	} else {
		switch {
		case max == rf:
			h = 60 * ((gf - bf) / delta)
			if h < 0 {
				h += 360
			}
		case max == gf:
			h = 60 * ((bf-rf)/delta + 2)
		case max == bf:
			h = 60 * ((rf-gf)/delta + 4)
		}
	}

	if max == 0 {
		s = 0
	} else {
		s = (delta / max) * 100.0
	}

	v = max * 100.0

	h = math.Round(h*100) / 100
	s = math.Round(s*100) / 100
	v = math.Round(v*100) / 100

	return h, s, v
}
