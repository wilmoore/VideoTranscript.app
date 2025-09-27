# VideoTranscript.app API

A fast, simple API for transcribing YouTube videos using Go, yt-dlp, ffmpeg, and whisper.cpp.

## Features

- **YouTube Video Transcription**: Convert YouTube videos to text with timestamps
- **Fast Processing**: Optimized pipeline with Go libraries (yt-dlp, ffmpeg, whisper)
- **Async & Sync**: Short videos (<2min) return transcripts immediately, longer videos use job queue
- **Simple Pricing**: $1 for videos <10min, $2 for 10-30min, first job free

## API Endpoints

### `POST /transcribe`

Transcribe a YouTube video.

**Request:**
```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}
```

**Response (Short Videos):**
```json
{
  "transcript": "Never gonna give you up, never gonna let you down...",
  "segments": [
    {
      "start": 0.0,
      "end": 3.5,
      "text": "Never gonna give you up"
    }
  ]
}
```

**Response (Long Videos):**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### `GET /transcribe/{job_id}`

Get the status and result of a transcription job.

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "complete",
  "transcript": "Never gonna give you up...",
  "segments": [...],
  "created_at": "2024-01-01T12:00:00Z",
  "completed_at": "2024-01-01T12:02:30Z"
}
```

Possible statuses: `pending`, `running`, `complete`, `error`

### `GET /health`

Health check endpoint (no authentication required).

## Usage Examples

### cURL

```bash
# Transcribe a video
curl -X POST http://localhost:3000/transcribe \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'

# Check job status
curl -X GET http://localhost:3000/transcribe/YOUR_JOB_ID \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### JavaScript

```javascript
const response = await fetch('http://localhost:3000/transcribe', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
  })
});

const result = await response.json();
```

## Development Setup

### Prerequisites

1. **Go 1.23+**
   ```bash
   # Install via official installer or package manager
   # https://golang.org/doc/install
   ```

2. **yt-dlp**
   ```bash
   # macOS
   brew install yt-dlp

   # Linux (Ubuntu/Debian)
   sudo apt update && sudo apt install yt-dlp

   # Or install via pip
   pip install yt-dlp
   ```

3. **FFmpeg**
   ```bash
   # macOS
   brew install ffmpeg

   # Linux (Ubuntu/Debian)
   sudo apt update && sudo apt install ffmpeg

   # Windows
   # Download from https://ffmpeg.org/download.html
   ```

4. **whisper.cpp**
   ```bash
   # Clone and build whisper.cpp
   git clone https://github.com/ggerganov/whisper.cpp.git
   cd whisper.cpp
   make

   # Download a model
   bash ./models/download-ggml-model.sh base.en

   # Copy binary to PATH or update PATH to include whisper.cpp directory
   sudo cp main /usr/local/bin/whisper.cpp
   # OR add to PATH: export PATH="$PATH:/path/to/whisper.cpp"

   # Make sure model is accessible
   sudo mkdir -p /usr/local/share/whisper
   sudo cp models/ggml-base.en.bin /usr/local/share/whisper/
   ```

### Local Development

1. **Clone and install dependencies:**
   ```bash
   git clone <your-repo>
   cd videotranscript-app
   go mod tidy
   ```

2. **Configure environment:**
   ```bash
   # Copy and edit .env
   cp .env .env.local
   vim .env.local
   ```

   ```env
   PORT=3000
   API_KEY=dev-api-key-12345
   WORK_DIR=/tmp/videotranscript
   MAX_VIDEO_LENGTH=1800
   FREE_JOB_LIMIT=5
   ```

3. **Run the server:**
   ```bash
   # Using Makefile (recommended)
   make dev           # Development mode with hot reload
   make run           # Build and run

   # Or manually
   go run main.go
   # or build and run
   go build -o videotranscript-app && ./videotranscript-app
   ```

4. **Test the API:**
   ```bash
   # Using Makefile
   make test          # Run all tests
   make test-short    # Run quick tests

   # Manual testing
   # Health check
   curl http://localhost:3000/health

   # Test transcription (replace with your API key)
   curl -X POST http://localhost:3000/transcribe \
     -H "Authorization: Bearer dev-api-key-12345" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
   ```

## Make Commands

The project includes a comprehensive Makefile with organized tasks:

```bash
# See all available commands
make help

# Development
make setup         # Install development dependencies
make dev           # Run with hot reload
make run           # Build and run

# Building
make build         # Build for current platform
make build-all     # Build for all platforms
make install       # Install to $GOPATH/bin

# Testing
make test          # Run all tests
make test-short    # Run quick tests
make test-coverage # Generate coverage report
make benchmark     # Run benchmark tests
make perf          # Run comprehensive performance tests

# Code Quality
make fmt           # Format code
make lint          # Run linter
make vet           # Run go vet
make check         # Run all quality checks

# Docker
make docker-build  # Build Docker image
make docker-run    # Run Docker container

# Utilities
make clean         # Clean build artifacts
make deps          # Download dependencies
make version       # Show version info
```

## Production Deployment

### Docker Deployment (Recommended)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o videotranscript-app

FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    && pip3 install yt-dlp

# Install whisper.cpp
RUN apk add --no-cache git make g++ \
    && git clone https://github.com/ggerganov/whisper.cpp.git /tmp/whisper \
    && cd /tmp/whisper \
    && make \
    && cp main /usr/local/bin/whisper.cpp \
    && bash ./models/download-ggml-model.sh base.en \
    && mkdir -p /usr/local/share/whisper \
    && cp models/ggml-base.en.bin /usr/local/share/whisper/ \
    && rm -rf /tmp/whisper

WORKDIR /app
COPY --from=builder /app/videotranscript-app .
COPY .env .

EXPOSE 3000

CMD ["./videotranscript-app"]
```

Build and run:
```bash
docker build -t videotranscript-app .
docker run -p 3000:3000 --env-file .env videotranscript-app
```

### Cloud Deployment

**Environment Variables to Set:**
- `PORT`: Server port (default: 3000)
- `API_KEY`: API key for authentication
- `WORK_DIR`: Working directory for temp files
- `MAX_VIDEO_LENGTH`: Max video length in seconds
- `FREE_JOB_LIMIT`: Free jobs per user

**Deployment Checklist:**
1. Install Go, yt-dlp, ffmpeg, whisper.cpp on the server
2. Set environment variables
3. Configure firewall to allow traffic on your port
4. Set up process manager (systemd, PM2, or Docker)
5. Configure reverse proxy (nginx/traefik) for HTTPS
6. Monitor logs and set up health checks

## Architecture

- **Framework**: Fiber (Go) - Fast HTTP framework
- **Job Queue**: In-memory map (simple MVP, can be upgraded to Redis)
- **Video Download**: go-ytdlp (Go wrapper for yt-dlp)
- **Audio Processing**: ffmpeg-go (Go wrapper for FFmpeg)
- **Transcription**: whisper.cpp via command execution
- **Auth**: API key Bearer token
- **Storage**: Temporary files in configurable work directory

## Performance & Scaling

**Current Implementation (MVP):**
- In-memory job queue (single instance)
- Synchronous processing for videos <2min
- Asynchronous processing for longer videos

**Production Scaling:**
- Replace in-memory queue with Redis
- Add horizontal scaling with load balancer
- Implement caching for repeated videos
- Add database for user management and billing
- Use dedicated storage (S3) for temporary files

## API Documentation

Full OpenAPI/Swagger documentation available at [`docs/swagger.yaml`](docs/swagger.yaml).

## Troubleshooting

**Common Issues:**

1. **"whisper.cpp: command not found"**
   - Make sure whisper.cpp is built and in your PATH
   - Or update the command path in `utils/transcription.go`

2. **"ggml-base.en.bin: no such file"**
   - Download the whisper model: `bash ./models/download-ggml-model.sh base.en`
   - Make sure it's accessible at the path specified in the code

3. **"ffmpeg: command not found"**
   - Install FFmpeg: `brew install ffmpeg` (macOS) or `apt install ffmpeg` (Linux)

4. **"yt-dlp failed"**
   - Update yt-dlp: `pip install --upgrade yt-dlp`
   - Some videos may be restricted or unavailable

**Debug Mode:**
Set environment variable `DEBUG=1` for detailed logging.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.