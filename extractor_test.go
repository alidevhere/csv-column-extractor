package csv_extractor

import (
	"fmt"
	"os"
	"testing"
)

func TestCopyCSVColumns(t *testing.T) {

	f, _ := os.Create("src.csv")
	f.WriteString("Name,age,phone")
	for i := 0; i < 100; i++ {
		f.WriteString(fmt.Sprintf("name%d,age%d,1234567890%d", i, i, i))
	}
	f.Close()

	err := CopyCSVColumns("src.csv", "dst1.csv",
		ExtractorOptions{SkipHeader: true, Columns: []int{0, 2}})
	if err != nil {
		fmt.Println(err)
	}

	os.Remove(f.Name())
	os.Remove("dst1.csv")

}
