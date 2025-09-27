package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testURL = "https://www.youtube.com/watch?v=GOHdTwKdT14"

func TestDownloadVideo(t *testing.T) {
	// Create test directory
	testDir := "./test_downloads"
	defer os.RemoveAll(testDir)

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test download function
	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %s", videoFile)
	}

	// Verify file is not empty
	stat, err := os.Stat(videoFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatal("Downloaded file is empty")
	}

	// Verify file extension
	if !strings.HasSuffix(videoFile, ".mp4") {
		t.Fatalf("Expected .mp4 file, got: %s", videoFile)
	}

	t.Logf("Successfully downloaded: %s (%.2f MB)", videoFile, float64(stat.Size())/(1024*1024))
}

func TestExtractAudio(t *testing.T) {
	// First download a video for testing
	testDir := "./test_audio"
	defer os.RemoveAll(testDir)

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Test audio extraction
	audioFile, err := extractAudio(videoFile)
	if err != nil {
		t.Fatalf("Audio extraction failed: %v", err)
	}

	// Verify audio file exists
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		t.Fatalf("Audio file does not exist: %s", audioFile)
	}

	// Verify file is not empty
	stat, err := os.Stat(audioFile)
	if err != nil {
		t.Fatalf("Failed to stat audio file: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatal("Audio file is empty")
	}

	// Verify file extension
	if !strings.HasSuffix(audioFile, ".wav") {
		t.Fatalf("Expected .wav file, got: %s", audioFile)
	}

	t.Logf("Successfully extracted audio: %s (%.2f MB)", audioFile, float64(stat.Size())/(1024*1024))
}

func TestTranscribeAudio(t *testing.T) {
	// This test is slower, so we'll use a smaller sample
	testDir := "./test_transcribe"
	defer os.RemoveAll(testDir)

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Download and extract audio
	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	audioFile, err := extractAudio(videoFile)
	if err != nil {
		t.Fatalf("Audio extraction failed: %v", err)
	}

	// Create a short sample for testing (first 10 seconds)
	sampleAudio := strings.TrimSuffix(audioFile, ".wav") + "_sample.wav"
	if err := createAudioSample(audioFile, sampleAudio, 10); err != nil {
		t.Fatalf("Failed to create audio sample: %v", err)
	}

	// Test transcription
	transcriptFiles, err := transcribeAudio(sampleAudio, testDir)
	if err != nil {
		t.Fatalf("Transcription failed: %v", err)
	}

	// Verify transcript files exist
	expectedFormats := []string{".txt", ".srt", ".vtt", ".json", ".tsv"}
	for _, format := range expectedFormats {
		found := false
		for _, file := range transcriptFiles {
			if strings.HasSuffix(file, format) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing transcript format: %s", format)
		}
	}

	// Verify transcript content is not empty
	for _, file := range transcriptFiles {
		if strings.HasSuffix(file, ".txt") {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read transcript: %v", err)
			}
			if len(strings.TrimSpace(string(content))) == 0 {
				t.Fatal("Transcript content is empty")
			}
			t.Logf("Transcript content preview: %s", string(content)[:min(100, len(content))])
		}
	}

	t.Logf("Successfully generated %d transcript files", len(transcriptFiles))
}

func TestFullPipeline(t *testing.T) {
	// Set longer timeout for full pipeline test
	if testing.Short() {
		t.Skip("Skipping full pipeline test in short mode")
	}

	testDir := "./test_pipeline"
	defer os.RemoveAll(testDir)

	start := time.Now()

	// Run full pipeline
	result, err := runFullPipeline(testURL, testDir)
	if err != nil {
		t.Fatalf("Full pipeline failed: %v", err)
	}

	duration := time.Since(start)

	// Verify all components
	if result.VideoFile == "" {
		t.Fatal("No video file in result")
	}
	if result.AudioFile == "" {
		t.Fatal("No audio file in result")
	}
	if len(result.TranscriptFiles) == 0 {
		t.Fatal("No transcript files in result")
	}

	t.Logf("Full pipeline completed in %v", duration)
	t.Logf("Video: %s", result.VideoFile)
	t.Logf("Audio: %s", result.AudioFile)
	t.Logf("Transcripts: %v", result.TranscriptFiles)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}