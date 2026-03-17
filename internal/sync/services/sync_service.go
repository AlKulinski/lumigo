package services

import (
	"fmt"
	"log"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
	"github.com/AlKulinski/lumigo/internal/sync/utils"
)

const THRESHOLD = 20

type SyncService interface {
	Sync(color domain.Color) error
}

type SyncServiceImpl struct {
	openHabRepository domain.OpenHabRepository
	lastValue         domain.Color
	topic             string
}

func NewSyncService(openHabRepository domain.OpenHabRepository, topic string) SyncService {
	return &SyncServiceImpl{
		openHabRepository: openHabRepository,
		lastValue: domain.Color{
			R:    0,
			G:    0,
			B:    0,
			Luma: 0,
		},
		topic: topic,
	}
}

func (s *SyncServiceImpl) Sync(color domain.Color) error {
	h, sat, v := utils.RGBToHSB(color.R, color.G, color.B)
	if v < 20 {
		v *= 5
	}
	if v > 60 {
		v = 100
	}

	log.Printf("h: %d, sat: %d, v: %d", int(h), int(sat), int(v))
	err := s.openHabRepository.SendEvent(
		domain.OpenHabMessage{
			Type:  domain.OpenHabTypeItemCommand,
			Topic: s.topic,
			Payload: domain.OpenHabPayload{
				Type:  domain.OpenHabPayloadTypeItemCommandHSB,
				Value: fmt.Sprintf("%d,%d,%d", int(h), int(sat), int(v)),
			},
		},
	)
	if err != nil {
		return err
	}
	s.lastValue = color
	return nil
}
