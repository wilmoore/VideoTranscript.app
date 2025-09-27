package main

import (
	"os"
	"strings"
	"testing"
)

const testURL = "https://www.youtube.com/watch?v=GOHdTwKdT14"

func TestDownloadVideo(t *testing.T) {
	testDir := "./test_downloads"
	defer os.RemoveAll(testDir)

	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %s", videoFile)
	}

	if !strings.HasSuffix(videoFile, ".mp4") {
		t.Fatalf("Expected .mp4 file, got: %s", videoFile)
	}

	stat, _ := os.Stat(videoFile)
	t.Logf("Successfully downloaded: %s (%.2f MB)", videoFile, float64(stat.Size())/(1024*1024))
}

func TestExtractAudio(t *testing.T) {
	testDir := "./test_audio"
	defer os.RemoveAll(testDir)

	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	audioFile, err := extractAudio(videoFile)
	if err != nil {
		t.Fatalf("Audio extraction failed: %v", err)
	}

	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		t.Fatalf("Audio file does not exist: %s", audioFile)
	}

	if !strings.HasSuffix(audioFile, ".wav") {
		t.Fatalf("Expected .wav file, got: %s", audioFile)
	}

	stat, _ := os.Stat(audioFile)
	t.Logf("Successfully extracted audio: %s (%.2f MB)", audioFile, float64(stat.Size())/(1024*1024))
}

func TestTranscription(t *testing.T) {
	testDir := "./test_transcription"
	defer os.RemoveAll(testDir)

	// Download and extract audio
	videoFile, err := downloadVideo(testURL, testDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	audioFile, err := extractAudio(videoFile)
	if err != nil {
		t.Fatalf("Audio extraction failed: %v", err)
	}

	// Create 30-second sample for faster testing
	sampleFile := strings.TrimSuffix(audioFile, ".wav") + "_sample.wav"
	if err := createAudioSample(audioFile, sampleFile, 30); err != nil {
		t.Fatalf("Sample creation failed: %v", err)
	}

	// Extract video ID and transcribe with organization
	videoID := extractVideoID(testURL)
	t.Logf("Video ID: %s", videoID)

	transcriptFiles, err := transcribeAudio(sampleFile, testDir, videoID)
	if err != nil {
		t.Fatalf("Transcription failed: %v", err)
	}

	t.Logf("Generated %d transcript files in directory: %s/%s", len(transcriptFiles), testDir, videoID)

	// Show the transcript content
	for _, file := range transcriptFiles {
		if strings.HasSuffix(file, ".txt") {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read transcript: %v", err)
			}
			t.Logf("Transcript content:\n%s", string(content))
			break
		}
	}
}