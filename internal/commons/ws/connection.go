package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketConfig struct {
	Host        string
	Port        string
	AccessToken string
}

func NewWebsocketConnection(cfg WebSocketConfig) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), Path: "/ws"}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+cfg.AccessToken)
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	go readLoop(c)

	return c
}

func readLoop(c *websocket.Conn) {
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			sendPing(c)
		}
	}
}

func sendPing(c *websocket.Conn) {
	err := c.WriteMessage(websocket.PingMessage, []byte("{\"type\":\"WebSocketEvent\",\"topic\":\"openhab/websocket/heartbeat\",\"payload\":\"PING\",\"source\":\"WebSocketTestInstance\"}"))
	if err != nil {
		log.Println("write:", err)
		return
	}

}
