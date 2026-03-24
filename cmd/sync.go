package cmd

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/AlKulinski/lumigo/internal/commons/mqtt"
	"github.com/AlKulinski/lumigo/internal/commons/ws"
	"github.com/AlKulinski/lumigo/internal/config"
	framing "github.com/AlKulinski/lumigo/internal/framing/services"
	"github.com/AlKulinski/lumigo/internal/sync/domain"
	"github.com/AlKulinski/lumigo/internal/sync/infra"
	"github.com/AlKulinski/lumigo/internal/sync/services"
	"github.com/AlKulinski/lumigo/internal/sync/usecases"
)

var (
	syncStreamService   string
	syncCameraInput     int
	syncCameraFPS       int
	syncFilePath        string
	syncDebugWindowMode bool
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
		streamService, runUI := buildStreamService(context.Background())
		usecase := usecases.NewSyncUsecase(syncService, streamService)

		if runUI != nil {
			go usecase.Execute()
			runUI()
		} else {
			usecase.Execute()
		}
	},
}

func buildStreamService(ctx context.Context) (domain.StreamService, func()) {
	switch strings.ToLower(syncStreamService) {
	case "camera":
		framingService := framing.NewFramingService()

		if syncDebugWindowMode {
			debugWindowService := framing.NewImageWindowService()
			return services.NewStreamCameraServiceImpl(ctx, syncCameraInput, syncCameraFPS, framingService, debugWindowService), debugWindowService.Run
		}

		return services.NewStreamCameraServiceImpl(ctx, syncCameraInput, syncCameraFPS, framingService, nil), nil
	case "file":
		if syncFilePath == "" {
			log.Fatal("sync file stream requires --file-path")
		}
		return services.NewStreamFileServiceImpl(syncFilePath, ctx), nil
	default:
		log.Fatalf("unsupported stream service %q, expected one of: camera, file", syncStreamService)
		return nil, nil
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	cameraInputDefault, err := strconv.Atoi(cfg.Sync.CameraID)
	if err != nil {
		cameraInputDefault = 0
	}

	syncCmd.Flags().StringVar(
		&syncStreamService,
		"service",
		cfg.Sync.Service,
		"Stream service to use: camera or file",
	)
	syncCmd.Flags().IntVar(
		&syncCameraInput,
		"camera-input",
		cameraInputDefault,
		"Camera input device index to use when service=camera",
	)
	syncCmd.Flags().IntVar(
		&syncCameraFPS,
		"camera-fps",
		10,
		"Camera FPS to use when service=camera",
	)
	syncCmd.Flags().StringVar(
		&syncFilePath,
		"file-path",
		cfg.Sync.FilePath,
		"Path to input media file when service=file",
	)
	syncCmd.Flags().BoolVar(
		&syncDebugWindowMode,
		"debug",
		false,
		"Debug window enabled debug=true",
	)
}
