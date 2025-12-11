package infra

import (
	"fmt"

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
	println(message.Payload.Value)
	token := r.client.Publish(message.Topic, 0, false, message.Payload.Value)
	token.Wait()
	fmt.Println("Published to topic: " + message.Topic)

	if token.Error() != nil {
		fmt.Println("Error publishing:", token.Error())
		fmt.Println("Error publishing to topic: " + message.Topic)
		return token.Error()
	}
	return nil
}
