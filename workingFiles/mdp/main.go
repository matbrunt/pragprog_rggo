package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// content type represents the HTML content to add into the template
type content struct {
	Title string
	Body  template.HTML
}

// Defines HTML template with dynamic content blocks to wrap markdown HTML generated content
const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
    {{ .Body }}
  </body>
</html>
`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	tFname := flag.String("t", "", "Alternate template name")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Coordinates execution of remaining functions
func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	// Read all data from input file and check for errors
	// ReadFile reads content of input markdown file into slice of bytes
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Converts markdown to HTML
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

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

	// Tidy up generated preview files once the function has terminated
	// (calling os.Exit would terminate immediately, and not run any deferred function calls)
	defer os.Remove(outName)

	return preview(outName)
}

// Receives a slice of bytes with markdown content, outputs slice of bytes with converted HTML content
func parseContent(input []byte, tFname string) ([]byte, error) {
	// Parse the markdown content through blackfriday and bluemonday to generate valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse the contents of the defaultTemplate const into a new template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content type, adding the title and body
	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	// Create a byte buffer to store the template executions result
	var buffer bytes.Buffer

	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	// Extract the contents of the buffer as a byte slice
	return buffer.Bytes(), nil
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
	err = exec.Command(cPath, cParams...).Run()

	// Simple hack to resolve race condition of browser not opening the file before the run function defer
	// call deletes the temporary file
	// Give the browser some time to open the file before deleting it
	// Better way is to clean up resources using a signal, will be implemented later in "Handling Signals"
	time.Sleep(2 * time.Second)

	return err
}
