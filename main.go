package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// search file which name contains a string (case-insensitive)
// sample calls:
//   gofind -search contains -size /tmp/sandbox
//   gofind /tmp/sandbox -search contains
//   gofind /tmp/sandbox -search contains -size

func main() {
	// Default root directory.
	root := "."

	// If the first argument is not a flag (doesn't start with "-"), treat it as the root.
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		root = os.Args[1]
		// Remove the root argument so that flag.Parse doesn't treat it as a flag.
		newArgs := []string{os.Args[0]}
		if len(os.Args) > 2 {
			newArgs = append(newArgs, os.Args[2:]...)
		}
		os.Args = newArgs
	}

	// Define flags for search string and size.
	search := flag.String("search", "", "Search string to match in file names")
	sizeFlag := flag.Bool("size", false, "Display file size if set")
	flag.Parse()

	// If a positional argument remains after flag.Parse(), use it as the root.
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Tilde Expansion: Expand "~" to the user's home directory.
	// Expand tilde if present (e.g., "~/D" -> "/home/username/D")	
	if strings.HasPrefix(root, "~") {
		home, err := os.UserHomeDir() // get the current user's home directory
		if err != nil {
			log.Fatalf("failed to get user home directory: %v", err)
		}
		if root == "~" {
			root = home
		} else if strings.HasPrefix(root, "~/") {
			root = filepath.Join(home, root[2:])
		}
	}

	// Convert search string to lower-case for case-insensitive matching.
	lowerSearch := strings.ToLower(*search)

	// Display the search parameter only if it is not empty.
	if lowerSearch != "" {
		fmt.Printf("Search parameter: %s\n", lowerSearch)
	}

    // Collect results in a slice
    var results []string

	// Walk through the directory tree.
	walkDir(root, lowerSearch, *sizeFlag, &results)

	// Print out the results.
	for _, fileInfo := range results {
		fmt.Println(fileInfo)
	}
}

// walkDir recursively processes the given directory.
func walkDir(dir string, search string, sizeFlag bool, results *[]string) {
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
			walkDir(fullPath, search, sizeFlag, results)
		} else {
            // Check if the file name contains the search string.
	    	// When search is empty, strings.Contains always returns true.
			// Convert file name to lower-case before comparing.			
			if strings.Contains(strings.ToLower(entry.Name()), search) {

				if sizeFlag {
				
					// Get file details (e.g., size) using os.Stat (Slow !)
					info, err := os.Stat(fullPath)
					if err != nil {
						log.Printf("failed to stat file %s: %v\n", fullPath, err)
						continue
					}

                    // Append formatted output to results slice
					*results = append(*results, fmt.Sprintf("%s\t%d", fullPath, info.Size()))
				} else {
					*results = append(*results, fullPath)
				}

			}
		}
	}
}
