package csv_extractor

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

type ExtractorOptions struct {
	// Skip copying first row into destination file
	SkipHeader bool
	// Columns to copy
	// The first column has index 0.
	Columns []int
	// Copy only the first N rows
	// If 0, all rows will be copied
	MaxRows int
}

// Copies only given columns of a CSV file to a new CSV file.
// The columns are specified by their index.
// The first column has index 0.
// The columns are copied in the order they are specified.
// src is the path to the source CSV file.
// dst is the path to the destination CSV file. If the file does not exist, it will be created.
func CopyCSVColumns(src, dst string, options ExtractorOptions) error {

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	l, err := csv.NewReader(srcFile).Read()
	if err != nil {
		return err
	}

	if err := validate(len(l), options.Columns); err != nil {
		return err
	}

	srcFile.Seek(0, io.SeekStart)

	dstwr := csv.NewWriter(dstFile)
	srcReader := csv.NewReader(srcFile)

	// skips header
	if options.SkipHeader {
		srcReader.Read()
	}

	rowIndex := 0
	for {

		if options.MaxRows > 0 && rowIndex >= options.MaxRows {
			break
		}

		record, err := srcReader.Read()

		if err != nil {
			break
		}

		var newRecord []string
		for _, column := range options.Columns {
			newRecord = append(newRecord, record[column])
		}

		if err := dstwr.Write(newRecord); err != nil {
			return err
		}

		rowIndex++
	}

	dstwr.Flush()

	return nil
}

func validate(totalColumns int, columns []int) error {
	if len(columns) == 0 {
		return errors.New("no columns specified")
	}

	for _, column := range columns {
		if column < 0 {
			return fmt.Errorf("column index must be positive, given %d", column)
		}

		if column >= totalColumns {
			return fmt.Errorf("column index %d out of range", column)
		}
	}

	return nil
}
