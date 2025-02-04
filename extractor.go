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

type CopyResult struct {
	DestinationFile    string
	TotalRowsCopied    int
	TotalColumnsCopied int
}
type CsvFileDimensions struct {
	Rows    int
	Columns int
}

// Copies only given columns of a CSV file to a new CSV file.
// The columns are specified by their index.
// The first column has index 0.
// The columns are copied in the order they are specified.
// src is the path to the source CSV file.
// dst is the path to the destination CSV file. If the file does not exist, it will be created.
func CopyCSVColumns(src, dst string, options ExtractorOptions) (CopyResult, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return CopyResult{}, err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return CopyResult{}, err
	}
	defer dstFile.Close()

	l, err := csv.NewReader(srcFile).Read()
	if err != nil {
		return CopyResult{}, err
	}

	if err := validate(len(l), options.Columns); err != nil {
		return CopyResult{}, err
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
			return CopyResult{}, err
		}
		rowIndex++
	}
	dstwr.Flush()
	return CopyResult{TotalColumnsCopied: len(options.Columns), TotalRowsCopied: rowIndex, DestinationFile: dstFile.Name()}, nil
}

func GetCsvFileDimensions(csvFileReader *csv.Reader) (CsvFileDimensions, error) {
	first, err := csvFileReader.Read()
	if err != nil {
		return CsvFileDimensions{}, err
	}
	dim := CsvFileDimensions{Columns: len(first), Rows: 1}

	for {
		_, err := csvFileReader.Read()
		if err != nil {
			break
		}
		dim.Rows += 1
	}
	return dim, nil
}

func SplitCSVFile(inputFilePath, outputBasePath string, chunkSize int, hasHeader bool) error {
	// Open input file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	// Read header if present
	var header []string
	if hasHeader {
		header, err = reader.Read()
		if err != nil {
			return fmt.Errorf("error reading header: %w", err)
		}
	}

	chunkNumber := 1
	var currentFile *os.File
	var writer *csv.Writer
	recordsWritten := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record: %w", err)
		}

		// Create new chunk file if needed
		if recordsWritten == 0 {
			filename := fmt.Sprintf("%s_%d.csv", outputBasePath, chunkNumber)
			currentFile, err = os.Create(filename)
			if err != nil {
				return fmt.Errorf("error creating chunk file: %w", err)
			}

			writer = csv.NewWriter(currentFile)
			// Write header to new chunk
			if hasHeader {
				if err := writer.Write(header); err != nil {
					currentFile.Close()
					return fmt.Errorf("error writing header: %w", err)
				}
			}
		}

		// Write record to current chunk
		if err := writer.Write(record); err != nil {
			currentFile.Close()
			return fmt.Errorf("error writing record: %w", err)
		}
		recordsWritten++

		// Finalize chunk if we've reached the size limit
		if recordsWritten == chunkSize {
			writer.Flush()
			if err := writer.Error(); err != nil {
				currentFile.Close()
				return fmt.Errorf("error flushing writer: %w", err)
			}
			currentFile.Close()
			chunkNumber++
			recordsWritten = 0
		}
	}

	// Finalize last chunk if there are remaining records
	if recordsWritten > 0 {
		writer.Flush()
		if err := writer.Error(); err != nil {
			currentFile.Close()
			return fmt.Errorf("error flushing final chunk: %w", err)
		}
		return currentFile.Close()
	}
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
