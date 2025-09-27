package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run simple_download.go <youtube_url>")
	}

	url := os.Args[1]
	outputDir := "./downloads"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Printf("Downloading from: %s\n", url)

	// Simple download using the working yt-dlp approach from troubleshooting
	cmd := exec.Command("yt-dlp",
		"--format", "best[ext=mp4]",
		"--output", filepath.Join(outputDir, "%(title)s.%(ext)s"),
		"--no-warnings",
		url,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	fmt.Println("\nDownload completed successfully!")

	// List what was downloaded
	files, err := filepath.Glob(filepath.Join(outputDir, "*"))
	if err == nil && len(files) > 0 {
		fmt.Println("\nDownloaded files:")
		for _, file := range files {
			fmt.Printf("  %s\n", file)
		}
	}
}