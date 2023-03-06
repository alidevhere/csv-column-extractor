# CSV Column Extractor
Extracts columns from one file and copies into other file in defined order.


## Download package

```
go get github.com/alidevhere/csv-column-extractor

```

# Example:
In this example column number 0 and 2 are copied into a new file, columns are copied into new file in defined order i.e 0 and 2 in this case.

```

err := CopyCSVColumns("src.csv", "dst.csv",ExtractorOptions{SkipHeader: true, Columns: []int{0, 2}})

if err != nil {
	fmt.Println(err)
}


```