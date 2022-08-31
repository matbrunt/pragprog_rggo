package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// Defines HTML header and footer constants to wrap markdown HTML generated content
const (
	header = `<!DOCTYPE html>
	<html>
		<head>
		<meta http-equiv="content-type" content="text/html"; charset=utf-8">
		<title>Markdown Preview Tool</title>
		</head>
		<body>
	`

	footer = `
	</body>
	</html>
	`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Coordinates execution of remaining functions
func run(filename string) error {
	// Read all data from input file and check for errors
	// ReadFile reads content of input markdown file into slice of bytes
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Converts markdown to HTML
	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	// Save HTML content to a file
	return saveHTML(outName, htmlData)
}

// Receives a slice of bytes with markdown content, outputs slice of bytes with converted HTML content
func parseContent(input []byte) []byte {
	// Parse the markdown content through blackfriday and bluemonday to generate valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Use a buffer of bytes (bytes.Buffer) to join the HTML header, body and footer components together
	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Write HTML to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	// Extract the contents of the buffer as a byte slice
	return buffer.Bytes()
}

// Save HTML content to an html file
func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file (0644 is read/write by owner, but only readable by anyone else)
	return os.WriteFile(outFname, data, 0644)
}
