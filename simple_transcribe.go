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
		log.Fatal("Usage: go run simple_transcribe.go <audio_file>")
	}

	audioFile := os.Args[1]

	// Check if audio file exists
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		log.Fatalf("Audio file does not exist: %s", audioFile)
	}

	// Create output directory for transcription files
	outputDir := "./transcripts"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Printf("Transcribing audio: %s\n", audioFile)
	fmt.Printf("Output directory: %s\n", outputDir)

	// Get base name for output files
	baseName := strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile))

	// Transcribe using Whisper
	cmd := exec.Command("whisper",
		audioFile,
		"--model", "base.en",
		"--output_dir", outputDir,
		"--output_format", "all",
		"--language", "en",
		"--word_timestamps", "True",
		"--verbose", "True",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Transcription failed: %v", err)
	}

	fmt.Println("\nTranscription completed successfully!")

	// List generated files
	pattern := filepath.Join(outputDir, baseName+".*")
	files, err := filepath.Glob(pattern)
	if err == nil && len(files) > 0 {
		fmt.Println("\nGenerated files:")
		for _, file := range files {
			if stat, err := os.Stat(file); err == nil {
				fmt.Printf("  %s (%.2f KB)\n", file, float64(stat.Size())/1024)
			}
		}
	}
}