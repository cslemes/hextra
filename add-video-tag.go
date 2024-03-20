package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func convertMarkdownToHTMLVideo(directoryPath string) {
	// Regular expression to match the Markdown image syntax for videos
	pattern := regexp.MustCompile(`!\[(.*?)\]\(videos/.*?\.mp4\)`)

	// Iterate over all files in the directory
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a Markdown file
		if strings.HasSuffix(info.Name(), ".md") {
			// Open the Markdown file and read its content
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Replace the Markdown image syntax with the HTML video tag
			modifiedContent := pattern.ReplaceAllStringFunc(string(content), func(match string) string {
				// Extract the video source from the Markdown syntax
				videoSource := strings.TrimPrefix(strings.TrimSuffix(match, ")"), "![")
				return fmt.Sprintf(`<video controls width="320" height="240"><source src="%s" type="video/mp4">Your browser does not support the video tag.</video>`, videoSource)
			})

			// Write the modified content back to the file
			err = ioutil.WriteFile(path, []byte(modifiedContent), 0644)
			if err != nil {
				return err
			}

			fmt.Printf("Processed: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", directoryPath, err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run markdown_to_html_video.go <path_to_markdown_files>")
		os.Exit(1)
	}

	directoryPath := os.Args[1]
	convertMarkdownToHTMLVideo(directoryPath)
}
