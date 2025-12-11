/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/AlKulinski/lumigo/internal/commons/mqtt"
	"github.com/AlKulinski/lumigo/internal/commons/ws"
	"github.com/AlKulinski/lumigo/internal/config"
	"github.com/AlKulinski/lumigo/internal/sync/infra"
	"github.com/AlKulinski/lumigo/internal/sync/services"
	"github.com/AlKulinski/lumigo/internal/sync/usecases"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Used to colors from a video to a lamp",
	Long:  `Used to colors from a video to a lamp`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		c := mqtt.NewClient(mqtt.MQTTConfig{
			Broker:   cfg.MQTT.Broker,
			ClientID: cfg.MQTT.ClientID,
		})
		_ = ws.NewWebsocketConnection(ws.WebSocketConfig{
			Host:        cfg.WebSocket.Host,
			Port:        cfg.WebSocket.Port,
			AccessToken: cfg.WebSocket.AccessToken,
		})
		openHabRepository := infra.NewOpenHabMqttRepository(c)
		syncService := services.NewSyncService(openHabRepository, cfg.MQTT.Topic)
		streamService := services.NewStreamDisplayService(2560, context.Background())
		usecase := usecases.NewSyncUsecase(syncService, streamService)

		usecase.Execute()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
