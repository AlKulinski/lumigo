package domain

import (
	"math"
)

type Color struct {
	R    uint32
	G    uint32
	B    uint32
	Luma uint32
}

func (c *Color) Distance(color *Color) float64 {
	dr := float64(c.R) - float64(color.R)
	dg := float64(c.G) - float64(color.G)
	db := float64(c.B) - float64(color.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func (c *Color) IsClose(color *Color, threshold float64) bool {
	return c.Distance(color) <= threshold
}

type Pixel struct {
	X     int
	Y     int
	Color Color
}
