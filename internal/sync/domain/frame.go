package domain

import (
	"image"
	"time"
)

type Frame struct {
	Image     image.Image
	Timestamp time.Time
}
