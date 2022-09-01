package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// filter rules
// 1. path points to a directory
// 2. file size (bytes) is less than minimum file size provided by user
// 3. file extension does not match extensions provided by user
func filterOut(path string, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
