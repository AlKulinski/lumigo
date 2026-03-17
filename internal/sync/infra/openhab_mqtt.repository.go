package infra

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/AlKulinski/lumigo/internal/sync/domain"
)

type OpenHabMqttRepository struct {
	client mqtt.Client
}

func NewOpenHabMqttRepository(client mqtt.Client) *OpenHabMqttRepository {
	return &OpenHabMqttRepository{
		client,
	}
}

func (r *OpenHabMqttRepository) SendEvent(message domain.OpenHabMessage) error {
	log.Println(message.Payload.Value)
	token := r.client.Publish(message.Topic, 0, false, message.Payload.Value)
	token.Wait()
	log.Println("published to topic:", message.Topic)

	if token.Error() != nil {
		log.Println("error publishing:", token.Error())
		log.Println("error publishing to topic:", message.Topic)
		return token.Error()
	}
	return nil
}
