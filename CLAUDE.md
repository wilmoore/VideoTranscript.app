# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## VideoTranscript.app - YouTube Video Transcription API

> A Go-based API for transcribing YouTube videos using yt-dlp, FFmpeg, and whisper.cpp with async job processing.

## Development Commands

All development operations use the comprehensive Makefile:

```bash
# Essential commands
make help          # Show all available commands with descriptions
make dev           # Run server in development mode
make test-short    # Run quick tests (skip load tests)
make build         # Build application
make clean         # Clean all build artifacts

# Testing & quality
make test          # Run all tests including load tests
make test-coverage # Generate coverage report in coverage/coverage.html
make benchmark     # Run performance benchmarks
make perf-short    # Quick performance tests
make fmt           # Format code (go fmt + goimports if available)
make lint          # Run golangci-lint (requires setup)
make check         # Run fmt + lint + vet + test-short

# Single test execution
go test -run TestSpecificFunction ./package
go test -v -run TestPostTranscribe_ValidationErrors ./handlers
```

Environment configuration via `.env` file:
```bash
PORT=3000
API_KEY=dev-api-key-12345
WORK_DIR=/tmp/videotranscript
MAX_VIDEO_LENGTH=1800
FREE_JOB_LIMIT=5
```

## Architecture Overview

### Core Processing Pipeline
The application implements a three-stage video transcription pipeline:

1. **Download Stage** (`lib/transcription.go`): Uses go-ytdlp wrapper to extract audio from YouTube URLs
2. **Normalize Stage**: Uses ffmpeg-go to convert audio to 16kHz mono WAV format suitable for whisper.cpp
3. **Transcription Stage**: Executes whisper.cpp binary to generate timestamped transcripts

### Job Processing System
**Sync vs Async Decision**: Videos ≤2 minutes process synchronously with real-time response polling. Longer videos use async job queue with job ID for status checking.

**Job States** (`jobs/job.go`): `pending` → `running` → `complete`/`error`

**Job Queue** (`jobs/queue.go`): Thread-safe in-memory map with read/write mutex. Single global instance initialized in main.go.

**Processing Flow** (`handlers/transcribe.go`):
- `PostTranscribe`: Creates job → determines sync/async based on duration → launches goroutine
- `GetTranscribeJob`: Returns job status and results for async jobs
- Background processing calls `lib.ProcessTranscription` which orchestrates the pipeline

### Package Organization
- `main.go`: Fiber server setup, middleware, routing
- `config/`: Environment variable management with defaults
- `handlers/`: HTTP request handlers and response logic
- `jobs/`: Job data structures and thread-safe queue implementation
- `lib/`: Core business logic (transcription pipeline, auth middleware)
- `models/`: Request/response structs and URL validation
- `scripts/`: Performance testing automation

### External Dependencies
**Runtime Requirements**: The application shells out to external binaries:
- `yt-dlp`: Must be in PATH for video downloads
- `ffmpeg`: Must be in PATH for audio processing
- `whisper.cpp`: Must be in PATH as `whisper.cpp` binary with `ggml-base.en.bin` model

**Go Libraries**:
- `github.com/gofiber/fiber/v2`: HTTP framework (Express-like for Go)
- `github.com/lrstanley/go-ytdlp`: Go wrapper for yt-dlp
- `github.com/u2takey/ffmpeg-go`: Go wrapper for FFmpeg
- `github.com/google/uuid`: Job ID generation

## Key Implementation Details

### Authentication
Bearer token authentication via `lib.AuthMiddleware()` checks `Authorization: Bearer <token>` against `API_KEY` environment variable. Health endpoint bypasses authentication.

### File Management
Temporary files created in `WORK_DIR` with job ID naming. Automatic cleanup via defer statements in `ProcessTranscription`.

### Error Handling
Errors propagate up through the pipeline with context using `fmt.Errorf` wrapping. Job errors stored in Job.Error field and returned via API.

### Testing Architecture
- Unit tests in `handlers/transcribe_test.go` with Fiber test framework
- Performance benchmarks in `lib/transcription_bench_test.go`
- Load testing with configurable concurrency and performance assertions
- Performance test runner script at `scripts/run_perf_tests.sh`

### Configuration Patterns
Environment variables loaded once at startup via `config.Load()`. No runtime config changes supported.

## API Endpoints

- `GET /health`: Unauthenticated health check
- `POST /transcribe`: Submit video URL, returns transcript (sync) or job ID (async)
- `GET /transcribe/{job_id}`: Get job status and results

Request/response handled via `models.TranscribeRequest` and `models.TranscribeResponse` structs.

## Development Notes

- Use `make dev` for development server with automatic restarts
- Tests require external dependencies (yt-dlp, ffmpeg, whisper.cpp) to pass fully
- Load tests have strict performance requirements that may fail without external deps
- The `parseDuration` function currently returns hardcoded 120 seconds (TODO: implement actual parsing)
- Job queue is in-memory only - jobs lost on restart (suitable for MVP)