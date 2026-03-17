package infra

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
)

type OpenHabRepository struct {
	websocket *websocket.Conn
}

type OpenHabMessage struct {
	Type    string `json:"type"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

func NewOpenHabRepository(websocket *websocket.Conn) *OpenHabRepository {
	return &OpenHabRepository{
		websocket,
	}
}

func (r *OpenHabRepository) SendEvent(openHabMessage domain.OpenHabMessage) error {
	log.Println(time.Now().UnixMilli())
	jsonPayloadBytes, err := json.Marshal(openHabMessage.Payload)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(OpenHabMessage{
		Type:    string(openHabMessage.Type),
		Topic:   string(openHabMessage.Topic),
		Payload: string(jsonPayloadBytes),
	})
	if err != nil {
		return err
	}
	return r.websocket.WriteMessage(websocket.TextMessage, jsonBytes)
}
