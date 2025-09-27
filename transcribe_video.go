package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func extractVideoID(url string) string {
	// Extract YouTube video ID from URL
	re := regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return "unknown"
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run transcribe_video.go <youtube_url>")
	}

	url := os.Args[1]

	// Create directories
	downloadDir := "./downloads"
	transcriptDir := "./transcripts"

	for _, dir := range []string{downloadDir, transcriptDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	fmt.Printf("🎬 Downloading video from: %s\n", url)

	// Step 1: Download video
	cmd := exec.Command("yt-dlp",
		"--format", "best[ext=mp4]",
		"--output", filepath.Join(downloadDir, "%(title)s.%(ext)s"),
		"--no-warnings",
		url,
	)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	// Find downloaded video file
	files, err := filepath.Glob(filepath.Join(downloadDir, "*.mp4"))
	if err != nil || len(files) == 0 {
		log.Fatal("No video file found after download")
	}
	videoFile := files[0]
	fmt.Printf("✅ Downloaded: %s\n", filepath.Base(videoFile))

	// Step 2: Extract audio
	fmt.Printf("🎵 Extracting audio...\n")
	ext := filepath.Ext(videoFile)
	audioFile := strings.TrimSuffix(videoFile, ext) + ".wav"

	cmd = exec.Command("ffmpeg",
		"-i", videoFile,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		"-y",
		audioFile,
	)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Audio extraction failed: %v", err)
	}
	fmt.Printf("✅ Audio extracted: %s\n", filepath.Base(audioFile))

	// Step 3: Transcribe
	// Extract video ID for organization
	videoID := extractVideoID(url)
	fmt.Printf("🎙️  Transcribing with Whisper (Video ID: %s)...\n", videoID)

	// Create video-specific directory
	videoDir := filepath.Join(transcriptDir, videoID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		log.Fatalf("Failed to create video directory: %v", err)
	}

	// Create audio file with video ID name for cleaner output
	videoIDAudio := filepath.Join(filepath.Dir(audioFile), videoID+".wav")
	if err := os.Rename(audioFile, videoIDAudio); err == nil {
		audioFile = videoIDAudio // Use renamed file
	}

	// Extract title from video file name for metadata
	videoFileName := filepath.Base(videoFile)
	title := strings.TrimSuffix(videoFileName, filepath.Ext(videoFileName))

	// Create metadata file
	metadata := map[string]interface{}{
		"video_id":    videoID,
		"url":         url,
		"title":       title,
		"timestamp":   time.Now().Format(time.RFC3339),
		"audio_file":  audioFile,
		"video_file":  videoFile,
	}

	metadataJSON, _ := json.MarshalIndent(metadata, "", "  ")
	metadataPath := filepath.Join(videoDir, "metadata.json")
	os.WriteFile(metadataPath, metadataJSON, 0644)

	cmd = exec.Command("whisper",
		audioFile,
		"--model", "base.en",
		"--output_dir", videoDir,
		"--output_format", "all",
		"--language", "en",
		"--word_timestamps", "True",
	)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Transcription failed: %v", err)
	}

	// Show results
	baseName := strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile))
	pattern := filepath.Join(videoDir, baseName+".*")
	files, _ = filepath.Glob(pattern)

	fmt.Printf("✅ Transcription complete! Generated files:\n")
	for _, file := range files {
		fmt.Printf("   📄 %s\n", file)
	}

	// Show sample of transcript
	txtFile := filepath.Join(transcriptDir, baseName+".txt")
	if content, err := os.ReadFile(txtFile); err == nil {
		fmt.Printf("\n📝 Sample transcript:\n%s\n", string(content)[:min(500, len(content))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}