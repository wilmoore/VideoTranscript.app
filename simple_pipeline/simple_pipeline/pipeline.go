package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PipelineResult struct {
	VideoFile       string
	AudioFile       string
	TranscriptFiles []string
	VideoID         string
}

func extractVideoID(url string) string {
	// Extract YouTube video ID from URL
	re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return "unknown"
}

func downloadVideo(url, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Use external yt-dlp (working approach)
	cmd := exec.Command("yt-dlp",
		"--format", "best[ext=mp4]",
		"--output", filepath.Join(outputDir, "%(title)s.%(ext)s"),
		"--no-warnings",
		url,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("download failed: %v", err)
	}

	// Find downloaded video file
	files, err := filepath.Glob(filepath.Join(outputDir, "*.mp4"))
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("no video file found after download")
	}

	return files[0], nil
}

func extractAudio(videoFile string) (string, error) {
	ext := filepath.Ext(videoFile)
	audioFile := strings.TrimSuffix(videoFile, ext) + ".wav"

	// Use ffmpeg-go library
	err := ffmpeg.Input(videoFile).
		Output(audioFile, ffmpeg.KwArgs{
			"ar": 16000,           // Sample rate
			"ac": 1,               // Mono channel
			"c:a": "pcm_s16le",    // PCM format
			"y": "",               // Overwrite output
		}).
		OverWriteOutput().
		Run()

	if err != nil {
		return "", fmt.Errorf("audio extraction failed: %v", err)
	}

	return audioFile, nil
}

func createAudioSample(inputFile, outputFile string, durationSeconds int) error {
	// Use ffmpeg-go library
	return ffmpeg.Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{
			"ss": "00:00:00",                              // Start time
			"t":  fmt.Sprintf("00:00:%02d", durationSeconds), // Duration
			"ar": 16000,                                   // Sample rate
			"ac": 1,                                       // Mono channel
			"c:a": "pcm_s16le",                           // PCM format
		}).
		OverWriteOutput().
		Run()
}

func transcribeAudio(audioFile, outputDir, videoID string) ([]string, error) {
	// Create video-specific directory
	videoDir := filepath.Join(outputDir, videoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create video directory: %v", err)
	}

	cmd := exec.Command("whisper",
		audioFile,
		"--model", "base.en",
		"--output_dir", videoDir,
		"--output_format", "all",
		"--language", "en",
		"--word_timestamps", "True",
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("transcription failed: %v", err)
	}

	// Find generated transcript files in video-specific directory
	baseName := strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile))
	pattern := filepath.Join(videoDir, baseName+".*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find transcript files: %v", err)
	}

	return files, nil
}

func runFullPipeline(url, outputDir string) (*PipelineResult, error) {
	// Extract video ID for organization
	videoID := extractVideoID(url)

	// Step 1: Download video
	videoFile, err := downloadVideo(url, outputDir)
	if err != nil {
		return nil, fmt.Errorf("download step failed: %v", err)
	}

	// Step 2: Extract audio
	audioFile, err := extractAudio(videoFile)
	if err != nil {
		return nil, fmt.Errorf("audio extraction step failed: %v", err)
	}

	// Step 3: Create sample for testing (first 30 seconds)
	sampleAudio := strings.TrimSuffix(audioFile, ".wav") + "_sample.wav"
	if err := createAudioSample(audioFile, sampleAudio, 30); err != nil {
		return nil, fmt.Errorf("sample creation failed: %v", err)
	}

	// Step 4: Transcribe with video ID organization
	transcriptDir := filepath.Join(outputDir, "transcripts")
	transcriptFiles, err := transcribeAudio(sampleAudio, transcriptDir, videoID)
	if err != nil {
		return nil, fmt.Errorf("transcription step failed: %v", err)
	}

	return &PipelineResult{
		VideoFile:       videoFile,
		AudioFile:       audioFile,
		TranscriptFiles: transcriptFiles,
		VideoID:         videoID,
	}, nil
}