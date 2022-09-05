package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func sum(data []float64) float64 {
	sum := 0.0

	for _, v := range data {
		sum += v
	}

	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// auxiliary type statsFunc uses same signature as sum & avg to represent a class of
// functions with this signature, so any function matching this signature qualifies
// as this type. Use it as an input parameter for any new calculation functions.
// statsFunc defines a generic statistical function
type statsFunc func(data []float64) float64

// csv reads all data as string. This function will parse the contents of the csv file
// into a slice of floating point numbers.
func csv2float(r io.Reader, column int) ([]float64, error) {
	// Create the csv reader used to read in data from csv files
	cr := csv.NewReader(r)

	// Adjust column number for 0 based index
	column--

	// Read in the entire file to memory
	// ReadAll reads in all lines (records) as a slice of fields (columns), with each field
	// being a slice of strings (values)
	// This means the data structure will be [][]string.
	allData, err := cr.ReadAll()
	if err != nil {
		// %w wraps the original error allowing you to decorate the error with additional
		// information whilst keeping the original error available for inspection
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	var data []float64

	// Loop through all records
	for i, row := range allData {
		if i == 0 {
			continue
		}

		// Checking number of columns in csv file
		if len(row) <= column {
			// File does not have that many columns
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		// Try to convert read data into a float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	return data, nil
}
