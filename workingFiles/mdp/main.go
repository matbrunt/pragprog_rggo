package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Coordinates execution of remaining functions
func run(filename string, out io.Writer, skipPreview bool) error {
	// Read all data from input file and check for errors
	// ReadFile reads content of input markdown file into slice of bytes
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Converts markdown to HTML
	htmlData := parseContent(input)

	// TempFile replaces the * character with a random number
	// Create the temporary file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	// We close the temp file after creating it as we don't want to write to it just yet
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()

	// Write the temp filename to the writer
	// This allows us to pass Stdout when running via CLI, and bytes.Buffer to capture output in a buffer
	// when running via a test
	fmt.Fprintln(out, outName)

	// Save HTML content to a file
	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outName)
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

// Automatically preview generated html file
func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}

	// Open the file using default program
	return exec.Command(cPath, cParams...).Run()
}
