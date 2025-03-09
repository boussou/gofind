package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
)

// sample call:
// gofind  -search propa  /tmp/sandbox

func main() {
	// Command-line flag for search string.
	var search = flag.String("search", "", "Search string to match in file names")
	flag.Parse()

    // Use the first positional argument as the root directory; default to "." if not provided.
    root := "."
    if flag.NArg() > 0 {
        root = flag.Arg(0)
    }


	// Display the search parameter only if it is not empty.
	if *search != "" {
		fmt.Printf("Search parameter: %s\n", *search)
	}

    // Collect results in a slice
    var results []string

    // Walk through the directory tree without goroutines
    walkDir(root, *search, &results)

    // Print out the results
    for _, fileInfo := range results {
        fmt.Println(fileInfo)
    }
}

// walkDir recursively processes the given directory without goroutines
func walkDir(dir string, search string, results *[]string) {
    // List the directory entries.
    entries, err := os.ReadDir(dir)
    if err != nil {
        log.Printf("failed to read directory %s: %v\n", dir, err)
        return
    }

    // Iterate over each entry.
    for _, entry := range entries {
        fullPath := filepath.Join(dir, entry.Name())
        if entry.IsDir() {
            // Recursively process subdirectories
            walkDir(fullPath, search, results)
        } else {
            // Check if the file name contains the search string.
	    // When search is empty, strings.Contains always returns true.
            if strings.Contains(entry.Name(), search) {

                // Append formatted output to results slice
                *results = append(*results, fmt.Sprintf("%s", fullPath))
            }
        }
    }
}
