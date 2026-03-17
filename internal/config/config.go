package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MQTT struct {
		Broker   string
		ClientID string
		Topic    string
	}
	WebSocket struct {
		Host        string
		Port        string
		AccessToken string
	}
	Sync struct {
		Service    string
		CameraID   string
		FilePath   string
		FFmpegPath string
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.MQTT.Broker = getEnv("MQTT_BROKER", "tcp://192.168.0.124:1883")
	cfg.MQTT.ClientID = getEnv("MQTT_CLIENT_ID", "go-mqtt-lumigo")
	cfg.MQTT.Topic = getEnv("MQTT_TOPIC", "lumigo:hsv")

	cfg.WebSocket.Host = getEnv("WS_HOST", "192.168.0.228")
	cfg.WebSocket.Port = getEnv("WS_PORT", "8080")
	cfg.WebSocket.AccessToken = getEnv("WS_ACCESS_TOKEN", "")

	cfg.Sync.Service = getEnv("SYNC_SERVICE", "camera")
	cfg.Sync.CameraID = getEnv("SYNC_CAMERA_ID", "0")
	cfg.Sync.FilePath = getEnv("SYNC_FILE_PATH", "")
	cfg.Sync.FFmpegPath = getEnv("FFMPEG_PATH", "ffmpeg")

	if cfg.WebSocket.AccessToken == "" {
		return nil, fmt.Errorf("WS_ACCESS_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
