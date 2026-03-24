package services

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/AlKulinski/lumigo/internal/framing/domain"
)

type ImageWindowService struct {
	app    fyne.App
	window fyne.Window
	img    *canvas.Image
}

func NewImageWindowService() domain.DebugWindowService {
	a := app.New()
	w := a.NewWindow("Image Window")

	img := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(800, 600))

	w.SetContent(img)
	w.Resize(fyne.NewSize(800, 600))

	return &ImageWindowService{
		app:    a,
		window: w,
		img:    img,
	}
}

func (s *ImageWindowService) Run() {
	s.window.Show()
	s.app.Run()
}

func (s *ImageWindowService) Update(img image.Image) {
	fyne.Do(func() {
		s.img.Image = img
		s.img.Refresh()
	})
}
