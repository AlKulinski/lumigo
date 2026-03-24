package services

import (
	"image"
	"log"
	"math"

	"github.com/AlKulinski/lumigo/internal/framing/domain"
	"gocv.io/x/gocv"
)

type GocvFramingServiceImpl struct {
}

func NewFramingService() domain.FramingService {
	return &GocvFramingServiceImpl{}
}

func (s *GocvFramingServiceImpl) Frame(frame gocv.Mat) (domain.Framing, bool) {
	gray, blur, edges := gocv.NewMat(), gocv.NewMat(), gocv.NewMat()
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	defer func() {
		frame.Close()
		gray.Close()
		blur.Close()
		edges.Close()
		kernel.Close()
	}()

	gocv.CvtColor(frame, &gray, gocv.ColorBGRToGray)
	gocv.GaussianBlur(gray, &blur, image.Pt(5, 5), 0, 0, gocv.BorderDefault)
	gocv.Canny(blur, &edges, 20, 100)

	contours := gocv.FindContours(edges, gocv.RetrievalList, gocv.ChainApproxSimple)
	defer contours.Close()

	frameArea := float64(frame.Cols() * frame.Rows())

	bestScore := -1.0
	var bestRect image.Rectangle

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		perimeter := gocv.ArcLength(contour, true)
		approx := gocv.ApproxPolyDP(contour, 0.02*perimeter, true)

		rect := gocv.BoundingRect(approx)
		if rect.Dx() <= 0 || rect.Dy() <= 0 {
			continue
		}

		area := float64(rect.Dx() * rect.Dy())
		if area < frameArea*0.15 {
			continue
		}

		aspect := float64(rect.Dx()) / float64(rect.Dy())
		aspectScore := 1.0 - math.Abs(aspect-16.0/9.0)/(16.0/9.0)
		if aspectScore < 0 {
			aspectScore = 0
		}

		areaScore := area / frameArea

		score := areaScore*0.5 + aspectScore*0.5

		log.Printf("contour: %v, rect: %v, score: %f\n", contour, rect, score)

		if score > bestScore {
			bestScore = score
			bestRect = rect
		}
	}

	if bestScore < 0 {
		return image.Rectangle{}, false
	}

	rect := clampRect(bestRect, frame.Cols(), frame.Rows())

	return rect, true
}

func clampRect(r domain.Framing, maxW, maxH int) domain.Framing {
	if r.Min.X < 0 {
		r.Min.X = 0
	}
	if r.Min.Y < 0 {
		r.Min.Y = 0
	}
	if r.Max.X > maxW {
		r.Max.X = maxW
	}
	if r.Max.Y > maxH {
		r.Max.Y = maxH
	}
	return r
}
