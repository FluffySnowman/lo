package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	git "github.com/go-git/go-git/v5"
)

// Function to format file sizes
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

// Function to get Git status of a file using go-git
func getGitStatus(path string) string {
	// Open the existing repository
	r, err := git.PlainOpen(".")
	if err != nil {
		return ""
	}

	// Get the worktree
	w, err := r.Worktree()
	if err != nil {
		return ""
	}

	// Get the status of a single file
	status, err := w.Status()
	if err != nil {
		return ""
	}

	fileStatus, exists := status[path]
	if !exists {
		return ""
	}

	// Construct the status string based on the status of the file
	statusStrings := []string{}
	if fileStatus.Staging != git.Unmodified {
		statusStrings = append(statusStrings, "Staged")
	}
	if fileStatus.Worktree != git.Unmodified {
		statusStrings = append(statusStrings, "Modified")
	}
	// if fileStatus.Extra != git.Unmodified {
	// 	statusStrings = append(statusStrings, "Untracked")
	// }

	return strings.Join(statusStrings, " ")
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

	// Initialize a tab writer for columnar output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header row
	fmt.Fprintln(w, "File\tSize\tStatus\t")

	// Loop through each file and print details in columns
	for _, file := range files {
		if !file.IsDir() {
			size := file.Size()
			formattedSize := formatSize(size)
			gitStatus := getGitStatus(file.Name())

			fmt.Fprintf(w, "%s\t%s\t[%s]\t\n", file.Name(), formattedSize, gitStatus)
		}
	}

	// Ensure output is flushed and displayed
	w.Flush()
}
