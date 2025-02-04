package csv_extractor

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func TestCopyCSVColumns(t *testing.T) {

	f, err := os.Create("src.csv")
	if err != nil {
		t.Fatalf("failed to create csv file for testing")
	}
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)
	headers := []string{"Name", "Age", "Phone"}
	if err := writer.Write(headers); err != nil {
		t.Fatalf("failed to write headers to csv file for testing with err %v", err)
	}

	for i := 1; i <= 100; i++ {
		writer.Write([]string{fmt.Sprintf("name%d", i), fmt.Sprintf("age%d", i), fmt.Sprintf("123456789%d", i)})
	}
	writer.Flush()
	f.Close()

	result, err := CopyCSVColumns("src.csv", "dst1.csv",
		ExtractorOptions{SkipHeader: true, Columns: []int{0, 2}})
	if err != nil {
		t.Fatalf("failed to copy csv file columns with err %v", err)
	}
	defer os.Remove("dst1.csv")

	if result.TotalColumnsCopied != 2 {
		t.Fatalf("expected 2 columns to be copied, but copied %v", result.TotalColumnsCopied)
	}
	if result.TotalRowsCopied != 100 {
		t.Fatalf("expected 100 rows to be copied, but copied %v", result.TotalRowsCopied)
	}
	if result.DestinationFile != "dst1.csv" {
		t.Fatalf("expected destination file to be dst1.csv but got %v", result.DestinationFile)
	}
}

func TestSplitCSVFile(t *testing.T) {
	f, err := os.Create("src.csv")
	if err != nil {
		t.Fatalf("failed to create csv file for testing")
	}
	defer os.Remove(f.Name())

	writer := csv.NewWriter(f)
	headers := []string{"Name", "Age", "Phone"}
	if err := writer.Write(headers); err != nil {
		t.Fatalf("failed to write headers to csv file for testing with err %v", err)
	}

	for i := 1; i <= 100; i++ {
		writer.Write([]string{fmt.Sprintf("name%d", i), fmt.Sprintf("age%d", i), fmt.Sprintf("123456789%d", i)})
	}
	writer.Flush()
	f.Close()

	result, err := SplitCSVFile("src.csv", "dst", 5, true)
	if err != nil {
		t.Fatalf("failed to copy csv file columns with err %v", err)
	}
	for i, r := range result {
		if r.TotalColumnsCopied != 3 {
			t.Fatalf("expected 2 columns to be copied, but copied %v", r.TotalColumnsCopied)
		}
		if r.TotalRowsCopied != 5 {
			t.Fatalf("expected 100 rows to be copied, but copied %v", r.TotalRowsCopied)
		}
		if destP := fmt.Sprintf("dst_%d.csv", i+1); r.DestinationFile != destP {
			t.Fatalf("expected destination file to be %s but got %v", destP, r.DestinationFile)
		}
		defer os.Remove(r.DestinationFile)
	}
}
