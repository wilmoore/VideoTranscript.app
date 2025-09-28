# VideoTranscript.app

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/doc/go1.23)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/wilmoore/VideoTranscript.app?style=flat&logo=github)](https://github.com/wilmoore/VideoTranscript.app/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/wilmoore/VideoTranscript.app)](https://github.com/wilmoore/VideoTranscript.app/issues)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker&logoColor=white)](Dockerfile)
[![Encore.dev](https://img.shields.io/badge/Encore.dev-Ready-6366F1?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyIDJMMjIgOFYxNkwxMiAyMkwyIDI2VjhMMTIgMloiIGZpbGw9IiM2MzY2RjEiLz4KPC9zdmc+)](https://encore.dev)

> A fast, production-ready API for transcribing YouTube videos using Go with native libraries and OpenAI Whisper, featuring comprehensive async job processing.

![logo](./docs/logo.png)

## Features

- **YouTube Video Transcription**: Convert YouTube videos to text with timestamps
- **Fast Processing**: Optimized pipeline with native Go libraries and FFmpeg
- **Async & Sync**: Short videos (<2min) return transcripts immediately, longer videos use job queue
- **Production Ready**: Database support, metrics, monitoring, and webhooks
- **Multiple Output Formats**: SRT, VTT, JSON, TSV, and plain text
- **Real-time Dashboard**: Live job monitoring and system metrics
- **Docker & Cloud Ready**: Easy deployment with comprehensive guides

## 📚 Documentation

<div align="center">

<table>
  <tr>
    <td>
      <h3>🚀 Quick Start</h3>
      <p>Get up and running in minutes.</p>
      <a href="#quick-start"><strong>Read guide →</strong></a>
    </td>
    <td>
      <h3>📖 API Docs</h3>
      <p>Complete reference with request & response examples.</p>
      <a href="docs/api.md"><strong>Explore →</strong></a>
    </td>
    <td>
      <h3>🏗️ Architecture</h3>
      <p>Deep dive into the system design and patterns.</p>
      <a href="docs/architecture.md"><strong>Understand →</strong></a>
    </td>
    <td>
      <h3>🚀 Deployment</h3>
      <p>Production-ready deployment playbooks.</p>
      <a href="docs/deployment.md"><strong>Deploy →</strong></a>
    </td>
  </tr>
  <tr>
    <td>
      <h3>💻 Development</h3>
      <p>Local setup, workflows, and contributor tooling.</p>
      <a href="docs/development.md"><strong>Build →</strong></a>
    </td>
    <td>
      <h3>🔧 Troubleshooting</h3>
      <p>Quick fixes for common pitfalls and errors.</p>
      <a href="docs/troubleshooting.md"><strong>Fix →</strong></a>
    </td>
    <td>
      <h3>🤝 Contributing</h3>
      <p>Guidelines for issues, pull requests, and reviews.</p>
      <a href="docs/contributing.md"><strong>Join →</strong></a>
    </td>
    <td>
      <h3>📋 Changelog</h3>
      <p>Track version history and notable updates.</p>
      <a href="docs/changelog.md"><strong>Review →</strong></a>
    </td>
  </tr>
</table>

</div>

## Quick Start

### 🐳 Docker (Recommended)
```bash
# Run with Docker (coming soon)
docker run -d \
  --name videotranscript \
  -p 3000:3000 \
  -e API_KEY=your-api-key-here \
  wilmoore/videotranscript:latest
```

### 🚀 Local Development
```bash
# 1. Install dependencies
# macOS
brew install go ffmpeg
pip install openai-whisper

# Ubuntu/Debian
sudo apt install golang-go ffmpeg python3-pip
pip install openai-whisper

# 2. Clone and run
git clone https://github.com/wilmoore/VideoTranscript.app.git
cd VideoTranscript.app
cp .env.example .env  # Edit with your settings
make dev
```

### ⚡ Encore.dev (Production)
```bash
# Deploy to production in one command
curl -L https://encore.dev/install.sh | bash
encore deploy --env production
```

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

**📖 Complete API documentation:** [docs/api.md](docs/api.md)

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

**💻 More usage examples and SDKs:** [docs/api.md](docs/api.md)

## Development Setup

### Prerequisites
- **Go 1.23+** - [Install Go](https://golang.org/doc/install)
- **FFmpeg** - `brew install ffmpeg` (macOS) or `sudo apt install ffmpeg` (Linux)
- **OpenAI Whisper** - `pip install openai-whisper`

### Quick Setup
```bash
# Clone and setup
git clone https://github.com/wilmoore/VideoTranscript.app.git
cd VideoTranscript.app
cp .env.example .env  # Edit with your settings
make dev
```

**💻 Detailed development guide:** [docs/development.md](docs/development.md)

## Make Commands

```bash
make help          # Show all available commands
make dev           # Development server with hot reload
make test          # Run all tests
make build         # Build for current platform
make check         # Run quality checks (fmt + lint + vet + test)
```

**🛠️ Complete command reference:** [docs/development.md#using-the-makefile](docs/development.md#using-the-makefile)

## Production Deployment

### 🐳 Docker (Recommended)
```bash
docker build -t videotranscript-app .
docker run -p 3000:3000 --env-file .env videotranscript-app
```

### ⚡ Encore.dev (Zero-Config)
```bash
encore deploy --env production
```

### ☁️ Cloud Platforms
- **AWS ECS/Fargate** - Container-based deployment
- **Google Cloud Run** - Serverless containers
- **Azure Container Instances** - Managed containers
- **Kubernetes** - Full orchestration

**🚀 Complete deployment guides:** [docs/deployment.md](docs/deployment.md)

## Architecture

**Three-Stage Pipeline:**
1. **Download** - Extract audio from YouTube URLs (native Go library)
2. **Normalize** - Convert to 16kHz mono WAV (FFmpeg)
3. **Transcribe** - Generate timestamped transcripts (OpenAI Whisper)

**Smart Processing:**
- Videos ≤2min: Synchronous (immediate results)
- Videos >2min: Asynchronous (job queue with status tracking)

**🏗️ Detailed architecture:** [docs/architecture.md](docs/architecture.md)

## API Documentation

Full OpenAPI/Swagger documentation available at [`docs/swagger.yaml`](docs/swagger.yaml).

## Troubleshooting

**Common issues and solutions:** [docs/troubleshooting.md](docs/troubleshooting.md)

## Contributing

We welcome contributions! Please see our [contributing guide](docs/contributing.md) for details on:
- Setting up your development environment
- Coding standards and best practices
- Submitting issues and pull requests
- Community guidelines

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**⭐ Star this repo if it's helpful!**

[🐛 Report Bug](https://github.com/wilmoore/VideoTranscript.app/issues) • [✨ Request Feature](https://github.com/wilmoore/VideoTranscript.app/issues) • [💬 Discussions](https://github.com/wilmoore/VideoTranscript.app/discussions)

**Built with ❤️ by [wilmoore](https://github.com/wilmoore)**

</div>
