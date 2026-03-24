package domain

import "gocv.io/x/gocv"

type FramingService interface {
	Frame(frame gocv.Mat) (Framing, bool)
}
