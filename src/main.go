package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

func removeANSICodes(input string) string {
	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return re.ReplaceAllString(input, "")
}

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
		return color.New(color.FgHiMagenta).Sprintf("%dd ago", int(hours/24))
	}
}

// func prependGitStatus(filename, path string) string {
// 	cmd := exec.Command("git", "status", "--short", path)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return filename
// 	}
// 	status := string(output)
// 	if len(status) < 2 {
// 		return filename
// 	}

// 	statusCode := strings.TrimSpace(status[:2])
// 	switch statusCode {
// 	case "M ":
// 		return color.New(color.FgYellow).Sprintf(" %s", filename)
// 	case "??":
// 		return color.New(color.FgRed).Sprintf(" %s", filename)
// 	case "A ":
// 		return color.New(color.FgGreen).Sprintf(" %s", filename)
// 	case "D ":
// 		return color.New(color.FgMagenta).Sprintf(" %s", filename)
// 	default:
// 		return color.New(color.FgCyan).Sprintf(" %s", filename)
// 	}
// }


func prependGitStatus(path string) string {
	cmd := exec.Command("git", "status", "--short", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "  " // Return two spaces if there's an error or no status.
	}
	status := string(output)
	if len(status) < 2 {
		return "  " // Return two spaces if no status is detected.
	}

	statusCode := strings.TrimSpace(status[:2])
	switch statusCode {
	case "M ":
		return color.New(color.FgYellow).Sprintf(" ")
	case "??":
		return color.New(color.FgRed).Sprintf(" ")
	case "A ":
		return color.New(color.FgGreen).Sprintf(" ")
	case "D ":
		return color.New(color.FgMagenta).Sprintf(" ")
	default:
		return color.New(color.FgCyan).Sprintf(" ")
	}
}

func printColoredName(file os.FileInfo) string {
	if file.IsDir() {
		return color.New(color.FgHiBlue, color.Bold).Sprint(file.Name())
	} else if file.Mode().Perm()&0111 != 0 { // Check for executable bit
		return color.New(color.FgHiGreen).Sprint(file.Name())
	}
	return color.New(color.FgWhite).Sprint(file.Name())
}

func gitDiffStat(path string) string {
	cmd := exec.Command("git", "diff", "--numstat", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	stats := strings.Split(strings.TrimSpace(string(output)), "\t")
	if len(stats) < 3 {
		return ""
	}
	additions := color.New(color.FgGreen).Sprintf("+%s", stats[0])
	deletions := color.New(color.FgRed).Sprintf("-%s", stats[1])
	return additions + " " + deletions
}

var (
	detailMode bool
	dirPath    string
)

func init() {
	flag.BoolVar(&detailMode, "d", false, "Show detailed file change stats")
	flag.BoolVar(&detailMode, "o", false, "Same as -d")
	// flag.StringVar(&dirPath, "path", ".", "Specify the path to list")
	// flag.StringVar(&dirPath, "p", ".", "Specify the path to list")
}

func main() {
	flag.Parse()

	var dirPath string
	if args := flag.Args(); len(args) > 0 {
		dirPath = args[0] 
	} else {
		dirPath = "." 
	}

	if dirPath == "." {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		dirPath = cwd
	} else {
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			fmt.Println("Error resolving absolute path:", err)
			return
		}
		dirPath = absPath
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory files:", err)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		} else if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].ModTime().After(files[j].ModTime())
	})

	var totalSize int64
	for _, file := range files {
		totalSize += file.Size()
	}

	fmt.Printf("\nCWD: %s\tTotal Files %s\n", color.New(color.FgHiMagenta).Sprint(dirPath), color.New(color.FgHiCyan, color.Italic, color.Bold).Sprintf("%d", len(files)))
	fmt.Printf("Total Size: %s\n\n", formatSize(totalSize))
	fmt.Printf("%-4s\t %-10s\t %-25s %s\n", "ID", "Size", "Modified", "File")
	for i, file := range files {
		filename := printColoredName(file)
		gitStatus := prependGitStatus(filepath.Join(dirPath, file.Name()))
		detail := ""
		if detailMode {
			detail = gitDiffStat(filepath.Join(dirPath, file.Name()))
		}
		fmt.Printf("%-4d\t %-10s\t %-25s\t %s%s %s\n", i, formatSize(file.Size()), timeSince(file.ModTime()), gitStatus, filename, detail)
	}
	println()
}
