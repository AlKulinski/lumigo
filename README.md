# Lumigo

Lumigo is a real-time ambient lighting synchronization tool that captures colors from your screen and syncs them to smart lights via MQTT. It analyzes the dominant colors in your display and sends HSV color commands to OpenHAB-compatible smart lighting systems.

## Features

- **Real-time Screen Capture**: Captures your screen at 20 FPS for smooth color transitions
- **Intelligent Color Analysis**: Analyzes pixel luminance to extract the most representative colors
- **Smart Light Integration**: Sends color commands via MQTT to OpenHAB smart home systems
- **WebSocket Support**: Includes WebSocket connectivity for additional integrations
- **Configurable**: Environment-based configuration for easy deployment

## How It Works

1. **Screen Capture**: Takes screenshots of your primary display at 20 FPS
2. **Color Analysis**: 
   - Calculates luminance for each pixel using standard RGB-to-luma conversion
   - Sorts pixels by luminance and samples the darkest 10% for dominant color extraction
   - Averages the sampled colors to get a representative color
3. **Color Conversion**: Converts RGB values to HSV (Hue, Saturation, Value) format
4. **Smart Light Control**: Sends HSV commands via MQTT to your smart lighting system

## Prerequisites

- Go 1.24.0 or later
- MQTT broker (tested with OpenHAB)
- Smart lights compatible with OpenHAB
- Access token for WebSocket connections (if using WebSocket features)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/AlKulinski/lumigo.git
cd lumigo
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o lumigo
```

## Configuration

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Configure the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MQTT_BROKER` | MQTT broker URL | `tcp://192.168.0.124:1883` |
| `MQTT_CLIENT_ID` | MQTT client identifier | `go-mqtt-lumigo` |
| `MQTT_TOPIC` | MQTT topic for color commands | `lumigo:hsv` |
| `WS_HOST` | WebSocket host | `192.168.0.228` |
| `WS_PORT` | WebSocket port | `8080` |
| `WS_ACCESS_TOKEN` | WebSocket access token | **Required** |
| `SYNC_SERVICE` | Default sync input service: `camera` or `file` | `camera` |
| `SYNC_CAMERA_ID` | Default camera device index for `sync --service camera` | `0` |
| `SYNC_FILE_PATH` | Default media file path for `sync --service file` | Empty |
| `FFMPEG_PATH` | Path to the `ffmpeg` binary used for file streaming | `ffmpeg` |

## Usage

Run the sync command to start capturing and syncing colors.

### Default camera input

```bash
./lumigo sync
```

This uses the `camera` service by default.

### Choose the stream service from the CLI

Use the camera input service explicitly:

```bash
./lumigo sync --service camera
```

Select a specific camera input device:

```bash
./lumigo sync --service camera --camera-input 1
```

Set a custom camera FPS:

```bash
./lumigo sync --service camera --camera-input 1 --camera-fps 10
```

Use a media file as the sync source:

```bash
./lumigo sync --service file --file-path ./test_mov/sample.mov
```

### Configure defaults with environment variables

You can set default sync behavior in your `.env` file:

```bash
SYNC_SERVICE=camera
SYNC_CAMERA_ID=1
SYNC_FILE_PATH=./test_mov/sample.mov
FFMPEG_PATH=/opt/homebrew/bin/ffmpeg
```

CLI flags override environment-based defaults.

The application will:
- Start capturing from the selected input source
- Analyze colors in real-time
- Send HSV color commands to your MQTT broker
- Continue running until stopped with Ctrl+C

## Architecture

The project follows clean architecture principles:

```
├── cmd/                    # CLI commands
├── internal/
│   ├── commons/           # Shared utilities
│   │   ├── mqtt/         # MQTT client
│   │   └── ws/           # WebSocket connection
│   ├── config/           # Configuration management
│   └── sync/             # Core synchronization logic
│       ├── domain/       # Business entities
│       ├── infra/        # Infrastructure layer
│       ├── services/     # Application services
│       ├── usecases/     # Business logic
│       └── utils/        # Utility functions
```

## Color Processing

The application uses sophisticated color analysis:

- **Luminance Calculation**: Uses the standard formula `0.299*R + 0.587*G + 0.114*B`
- **Sampling Strategy**: Focuses on darker pixels (bottom 10% by luminance) for better ambient lighting
- **HSV Conversion**: Converts RGB to HSV for better color representation in lighting systems
- **Brightness Adjustment**: Automatically adjusts brightness values for optimal lighting

## MQTT Integration

Lumigo sends OpenHAB-compatible MQTT messages:

```json
{
  "type": "ItemCommand",
  "topic": "lumigo:hsv",
  "payload": {
    "type": "HSB",
    "value": "240,75,85"
  }
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Troubleshooting

### Common Issues

- **Screen capture fails**: Ensure the application has screen recording permissions on macOS
- **Camera input fails**: Try a different `--camera-input` value or set `SYNC_CAMERA_ID` in your environment
- **File input fails**: Verify `--file-path` points to a valid media file and `FFMPEG_PATH` points to an installed `ffmpeg` binary
- **MQTT connection issues**: Verify your broker URL and network connectivity
- **WebSocket errors**: Check your access token and WebSocket server availability

### Performance Tips

- The application captures at 20 FPS by default - adjust the `FPS` constant in `stream_display_service.go` if needed
- Color sampling uses 10% of pixels - modify the `pickSample` function for different sampling rates