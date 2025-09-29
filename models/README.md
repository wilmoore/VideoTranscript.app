# Whisper Models

This directory contains pre-trained Whisper models for offline speech transcription.

## Available Models

| Model | Size | Language | Speed | Accuracy | Use Case |
|-------|------|----------|-------|----------|----------|
| ggml-tiny.en.bin | ~37MB | English | 32x realtime | Basic | Quick transcription, real-time |
| ggml-base.en.bin | ~148MB | English | 16x realtime | Good | **Recommended for most use cases** |
| ggml-small.en.bin | ~483MB | English | 6x realtime | Better | High-quality English transcription |
| ggml-medium.en.bin | ~1.5GB | English | 2x realtime | High | Professional English transcription |
| ggml-large-v3.bin | ~3GB | Multilingual | 1x realtime | Best | Best quality, supports 99+ languages |

## Model Selection Guide

- **For most users**: Use `ggml-base.en.bin` (default)
- **For real-time applications**: Use `ggml-tiny.en.bin`
- **For highest quality**: Use `ggml-large-v3.bin`
- **For non-English content**: Use `ggml-large-v3.bin`

## Configuration

Set the model path in your `.env` file:

```bash
WHISPER_MODEL_PATH=models/ggml-base.en.bin
```

## Usage

Models are automatically used by the transcription system when available.
The application will fall back to cloud services if models are not present.

## Download

Use the provided script to download models:

```bash
# Download default model
./scripts/download-models.sh

# Download specific model
./scripts/download-models.sh ggml-tiny.en.bin

# Download all models
./scripts/download-models.sh all
```

## License

Models are distributed under the MIT license by OpenAI.
See: https://github.com/openai/whisper
