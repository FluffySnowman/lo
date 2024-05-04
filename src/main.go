package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

func formatSize(size int64) string {
	const (
		KB int64 = 1024
		MB       = KB * KB
		GB       = MB * KB
	)
	switch {
	case size >= GB:
		return color.New(color.FgHiRed).Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return color.New(color.FgHiYellow).Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return color.New(color.FgHiGreen).Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return color.New(color.FgHiCyan).Sprintf("%d bytes", size)
	}
}

func timeSince(modTime time.Time) string {
	hours := time.Since(modTime).Hours()
	switch {
	case hours < 1:
		return color.New(color.FgHiGreen).Sprintf("%d min ago", int(time.Since(modTime).Minutes()))
	case hours < 24:
		return color.New(color.FgHiYellow).Sprintf("%dh ago", int(hours))
	default:
		return color.New(color.FgHiRed).Sprintf("%dd ago", int(hours/24))
	}
}

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
		return color.New(color.FgYellow).Sprintf(" %s", filename)
	case "??":
		return color.New(color.FgRed).Sprintf(" %s", filename)
	case "A ":
		return color.New(color.FgGreen).Sprintf(" %s", filename)
	case "D ":
		return color.New(color.FgMagenta).Sprintf(" %s", filename)
	default:
		return color.New(color.FgBlue).Sprintf(" %s", filename)
	}
}

func printColoredName(file os.FileInfo) string {
	if file.IsDir() {
		return color.New(color.FgHiBlue, color.Bold).Sprint(file.Name())
	}
	return color.New(color.FgWhite).Sprint(file.Name())
}

type byModTime []os.FileInfo

func (a byModTime) Len() int      { return len(a) }
func (a byModTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byModTime) Less(i, j int) bool {
	return a[i].ModTime().After(a[j].ModTime())
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

	// Sort files by modification time, directories first
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		} else if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].ModTime().After(files[j].ModTime())
	})

	fmt.Printf("%-4s\t %-10s\t %-15s %s\n", "ID", "Size", "Modified", "File")
	for i, file := range files {
		filenameWithStatus := prependGitStatus(printColoredName(file), filepath.Join(cwd, file.Name()))
		fmt.Printf("%-4d\t %-10s\t %-15s\t %s\n",
			i, formatSize(file.Size()), timeSince(file.ModTime()), filenameWithStatus)
	}
}
