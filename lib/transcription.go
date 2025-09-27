package lib

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lrstanley/go-ytdlp"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"videotranscript-app/config"
	"videotranscript-app/models"
)

func ProcessTranscription(url, jobID string) (string, []models.Segment, error) {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	audioFile := filepath.Join(cfg.WorkDir, fmt.Sprintf("%s.wav", jobID))
	normalizedAudio := filepath.Join(cfg.WorkDir, fmt.Sprintf("%s_norm.wav", jobID))
	transcriptFile := filepath.Join(cfg.WorkDir, fmt.Sprintf("%s_transcript.txt", jobID))

	defer func() {
		os.Remove(audioFile)
		os.Remove(normalizedAudio)
		os.Remove(transcriptFile)
	}()

	if err := downloadAudio(url, audioFile); err != nil {
		return "", nil, fmt.Errorf("failed to download audio: %w", err)
	}

	if err := normalizeAudio(audioFile, normalizedAudio); err != nil {
		return "", nil, fmt.Errorf("failed to normalize audio: %w", err)
	}

	if err := transcribeAudio(normalizedAudio, transcriptFile); err != nil {
		return "", nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	transcript, segments, err := models.LoadTranscript(transcriptFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load transcript: %w", err)
	}

	return transcript, segments, nil
}

func downloadAudio(url, outputPath string) error {
	dl := ytdlp.New().
		ExtractAudio().
		AudioFormat("wav").
		AudioQuality("0").
		Output(outputPath)

	result, err := dl.Run(context.Background(), url)
	if err != nil {
		return fmt.Errorf("yt-dlp failed: %w", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("yt-dlp failed with code %d: %s", result.ExitCode, result.Stderr)
	}

	return nil
}

func normalizeAudio(inputPath, outputPath string) error {
	err := ffmpeg_go.Input(inputPath).
		Audio().
		Output(outputPath, ffmpeg_go.KwArgs{
			"ar":  16000,
			"ac":  1,
			"c:a": "pcm_s16le",
			"y":   nil,
		}).
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg normalization failed: %w", err)
	}

	return nil
}

func transcribeAudio(audioPath, outputPath string) error {
	// Use OpenAI Whisper with multiple output formats
	outputDir := filepath.Dir(outputPath)

	cmd := exec.Command("whisper", audioPath,
		"--model", "base.en",
		"--output_dir", outputDir,
		"--output_format", "all", // Generate txt, srt, vtt, json, tsv
		"--language", "en",
		"--word_timestamps", "True",
		"--max_line_width", "80",
		"--max_line_count", "2",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("whisper failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Whisper output: %s\n", string(output))
	return nil
}

func GetVideoDuration(url string) (int, error) {
	dl := ytdlp.New()

	result, err := dl.Run(context.Background(), url, "--get-duration", "--no-warnings")
	if err != nil {
		return 0, fmt.Errorf("failed to get video info: %w", err)
	}

	if result.ExitCode != 0 {
		return 0, fmt.Errorf("yt-dlp failed with code %d: %s", result.ExitCode, result.Stderr)
	}

	return parseDuration(result.Stdout), nil
}

func parseDuration(duration string) int {
	return 120
}
