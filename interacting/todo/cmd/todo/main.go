package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"pragprog.com/rggo/interacting/todo"
)

// Hardcoding the file name
var todoFileName = ".todo.json"

// Command line flags:
// -list: Boolean flag, when specified tool will list all to-do items
// -task: String flag, when used tool will include string argument as new to do item in the list
// -complete: Integer flag, when used tool will mark the item number as completed
func main() {
	// Display custom usage message for tool.
	// PrintDefaults will print usage information for each specified flag
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for the Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	// Assigned variables are pointers, so will need to be dereferenced with * when used later
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Parse()

	// Check if the user defined the Env Var for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Create pointer to type todo.List by using address operator & to get the address
	// of an empty instance of todo.List
	l := &todo.List{}

	// Read existing items from file
	// Good practice to use STDERR for error messages rather than STDOUT so user can
	// easily filter them out if they desire.
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on provided flags (need dereferencing with *)
	switch {
	// Check if -list flag set
	case *list:
		// List current to do items
		fmt.Print(l) // uses the fmt.Stringer String() interface implementation
		// for _, item := range *l {
		// print only items not marked as completed
		// if !item.Done {
		// 	fmt.Println(item.Task)
		// }
		// }

	// Check if -complete flag set with value greater than 0 (default)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	// Add a new task if -add flag set
	case *add:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		// ... suffix operator expands the slice into a list of values
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Add the task
		l.Add(t)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	// Print error message to STDERR if an unhandled flag provided
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the description for a new task from:
// arguments or STDIN
// ...string means 0 or more arguments of type string (makes it a variadic function)
func getTask(r io.Reader, args ...string) (string, error) {
	// If arguments provided, return them concatenated with a space
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// Otherwise scan for a single input line from the passed reader interface
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
