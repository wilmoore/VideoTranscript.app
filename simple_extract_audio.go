package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run simple_extract_audio.go <video_file>")
	}

	videoFile := os.Args[1]

	// Check if video file exists
	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		log.Fatalf("Video file does not exist: %s", videoFile)
	}

	// Create output filename
	ext := filepath.Ext(videoFile)
	audioFile := strings.TrimSuffix(videoFile, ext) + ".wav"

	fmt.Printf("Extracting audio from: %s\n", videoFile)
	fmt.Printf("Output audio file: %s\n", audioFile)

	// Extract audio using ffmpeg
	cmd := exec.Command("ffmpeg",
		"-i", videoFile,
		"-ar", "16000",  // Sample rate for Whisper
		"-ac", "1",      // Mono channel
		"-c:a", "pcm_s16le", // PCM format
		"-y",            // Overwrite output
		audioFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Audio extraction failed: %v", err)
	}

	fmt.Println("\nAudio extraction completed successfully!")

	// Check output file
	if stat, err := os.Stat(audioFile); err == nil {
		fmt.Printf("Audio file size: %.2f MB\n", float64(stat.Size())/(1024*1024))
	}
}