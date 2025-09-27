package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PipelineResult struct {
	VideoFile       string
	AudioFile       string
	TranscriptFiles []string
}

func downloadVideo(url, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

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

	cmd := exec.Command("ffmpeg",
		"-i", videoFile,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		"-y",
		audioFile,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("audio extraction failed: %v", err)
	}

	return audioFile, nil
}

func createAudioSample(inputFile, outputFile string, durationSeconds int) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-ss", "00:00:00",
		"-t", fmt.Sprintf("00:00:%02d", durationSeconds),
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		"-y",
		outputFile,
	)

	return cmd.Run()
}

func transcribeAudio(audioFile, outputDir string) ([]string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}

	cmd := exec.Command("whisper",
		audioFile,
		"--model", "base.en",
		"--output_dir", outputDir,
		"--output_format", "all",
		"--language", "en",
		"--word_timestamps", "True",
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("transcription failed: %v", err)
	}

	// Find generated transcript files
	baseName := strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile))
	pattern := filepath.Join(outputDir, baseName+".*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find transcript files: %v", err)
	}

	return files, nil
}

func runFullPipeline(url, outputDir string) (*PipelineResult, error) {
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

	// Step 4: Transcribe
	transcriptFiles, err := transcribeAudio(sampleAudio, outputDir)
	if err != nil {
		return nil, fmt.Errorf("transcription step failed: %v", err)
	}

	return &PipelineResult{
		VideoFile:       videoFile,
		AudioFile:       audioFile,
		TranscriptFiles: transcriptFiles,
	}, nil
}