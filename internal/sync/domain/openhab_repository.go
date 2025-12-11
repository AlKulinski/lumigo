package domain

type OpenHabRepository interface {
	SendEvent(message OpenHabMessage) error
}
