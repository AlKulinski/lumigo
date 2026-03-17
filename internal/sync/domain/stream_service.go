package domain

type StreamService interface {
	DisplayStream() (<-chan Frame, error)
}
