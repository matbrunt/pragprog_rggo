package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// We don't need to repeat the API unit tests as they are handled by `todo_test.go`.
// Instead we can execute integration tests to the CLI wrapper implementation around the API,
// so we are testing the user interface of the tool rather than the business logic.

// os/exec allows us to execute external commands

var (
	binName  = "todo" // binary name of file built during tests
	fileName = ".todo.json"
)

// Use TestMain to execute extra setup before your tests
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	os.Remove(binName)
	os.Remove(fileName)

	// Call the Go build tool to build the executable binary
	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	// Execute the tests using m.Run()
	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

// Create test cases in TestTodoCLI, using subtests feature to execute tests that depend
// on each other by using the t.Run method from testing package.
func TestTodoCLI(t *testing.T) {
	// task name
	task := "test task number 1"

	// current working directory
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// path to command executable compiled in TestMain
	cmdPath := filepath.Join(dir, binName)

	// Create first test to ensure tool can add a new task by using t.Run
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Ensure the tool can list the tasks
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n", task)

		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
