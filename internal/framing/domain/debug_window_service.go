package domain

import "image"

type DebugWindowService interface {
	Run()
	Update(image image.Image)
}
