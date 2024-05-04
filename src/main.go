package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Helper to format file sizes.
func formatSize(size int64) string {
	const (
		KB int64 = 1024
		MB       = 1024 * KB
	)
	if size >= MB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else if size >= KB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}
	return fmt.Sprintf("%d bytes", size)
}

// Helper to format modification time.
func timeSince(modTime time.Time) string {
	duration := time.Since(modTime)
	if minutes := duration.Minutes(); minutes < 60 {
		return fmt.Sprintf("%d min ago", int(minutes))
	}
	return fmt.Sprintf("%dh ago", int(duration.Hours()))
}

// Helper to prepend Git status.
func prependGitStatus(filename, path string) string {
	cmd := exec.Command("git", "status", "--short", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return filename
	}
	status := string(output)
	if len(status) < 2 {
		return filename
	}

	statusCode := strings.TrimSpace(status[:2])
	switch statusCode {
	case "M ":
		return color.New(color.FgYellow).Sprintf("⎇ %s", filename)
	case "??":
		return color.New(color.FgRed).Sprintf("⎇ %s", filename)
	case "A ":
		return color.New(color.FgGreen).Sprintf("⎇ %s", filename)
	case "D ":
		return color.New(color.FgMagenta).Sprintf("⎇ %s", filename)
	default:
		return color.New(color.FgBlue).Sprintf("⎇ %s", filename)
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		fmt.Println("Error reading directory files:", err)
		return
	}

	// Calculate maximum lengths for dynamic spacing
	maxNameLength := 0
	for _, file := range files {
		filename := prependGitStatus(file.Name(), filepath.Join(cwd, file.Name()))
		if len(filename) > maxNameLength {
			maxNameLength = len(filename)
		}
	}

	// Print headers
	fmt.Printf("%-4s %-"+fmt.Sprint(maxNameLength)+"s %-10s %-15s\n", "ID", "File", "Size", "Modified")

	// Print file details
	for i, file := range files {
		filenameWithStatus := prependGitStatus(file.Name(), filepath.Join(cwd, file.Name()))
		fmt.Printf("%-4d %-"+fmt.Sprint(maxNameLength)+"s %-10s %-15s\n",
			i, filenameWithStatus, formatSize(file.Size()), timeSince(file.ModTime()))
	}
}
