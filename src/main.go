package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Function to format and color file sizes
func formatSize(size int64) string {
	const (
		KB int64 = 1024
		MB       = 1024 * KB
	)
	if size >= MB {
		return color.New(color.FgHiYellow).Sprintf("%.2f MB", float64(size)/float64(MB))
	} else if size >= KB {
		return color.New(color.FgHiCyan).Sprintf("%.2f KB", float64(size)/float64(KB))
	}
	return color.New(color.FgHiRed).Sprintf("%d bytes", size)
}

// Function to get and color Git status of a file
func getGitStatus(path string) string {
	cmd := exec.Command("git", "status", "--short", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	status := string(output)
	if len(status) < 2 {
		return ""
	}

	statusCode := strings.TrimSpace(status[:2])
	switch statusCode {
	case "M ":
		return color.New(color.FgYellow).Sprint("Modified")
	case "??":
		return color.New(color.FgRed).Sprint("Untracked")
	case "A ":
		return color.New(color.FgGreen).Sprint("Added")
	case "D ":
		return color.New(color.FgMagenta).Sprint("Deleted")
	default:
		return color.New(color.FgBlue).Sprint(statusCode)
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

	// Set up the tab writer for formatted output
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

	fmt.Fprintln(w, "File\tSize\t Status\t|")

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(cwd, file.Name())
			size := file.Size()
			formattedSize := formatSize(size)
			gitStatus := getGitStatus(filePath) // Now using the absolute file path

			fmt.Fprintf(w, "%s\t| %s\t\t\t| %s\t|\n", file.Name(), formattedSize, gitStatus)
		}
	}

	// Ensure output is flushed and displayed correctly
	w.Flush()
}
