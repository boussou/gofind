package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const maxConcurrent = 1000

// gofind searches for files whose names contain a search string (case-insensitive).
// Optionally, it displays file sizes when -size is set, and prints directory names if -printdir is set.
// printing directory names allows to totally mimic the linux find command output.
// Sample calls:
//   gofind -search contains -size /tmp/sandbox
//   gofind /tmp/sandbox -search contains  -printdir
//   gofind /tmp/sandbox -search contains -size  -printdir
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

	// Define flags for search string, size display, and printing directory names.
	search := flag.String("search", "", "Search string to match in file names (case-insensitive)")
	sizeFlag := flag.Bool("size", false, "Display file size if set")
	printDirFlag := flag.Bool("printdir", false, "Also print directory names that match the search string")
	flag.Parse()

	// If a positional argument remains after flag.Parse(), use it as the root.
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Tilde Expansion: Expand "~" to the user's home directory.
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

	// Channel to send matching file (or directory) info.
	fileCh := make(chan string)
	var wg sync.WaitGroup

	// Create a semaphore channel to limit concurrent walkDir invocations.
	sem := make(chan struct{}, maxConcurrent)

	// walkDir recursively processes the given directory and its subdirectories concurrently.
	// It uses the semaphore (sem) to limit concurrent calls.
	var walkDir func(dir string, search string, sizeFlag bool, printDir bool, sem chan struct{})
	walkDir = func(dir string, search string, sizeFlag bool, printDir bool, sem chan struct{}) {
		// Acquire a slot in the semaphore.
		sem <- struct{}{}
		// Release the slot when done.
		defer func() { <-sem }()
		defer wg.Done()

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
				// Optionally print the directory name if -printDir is set and its name matches the search.
				if printDir {
					lowerDirName := strings.ToLower(entry.Name())
		                	// When search is empty, strings.Contains always returns true.					
					if strings.Contains(lowerDirName, search) {
						// Mark the output as a directory.
						if sizeFlag {
							fileCh <- fmt.Sprintf("%s\tDIR", fullPath)
						} else {
							fileCh <- fullPath
						}
					}
				}

				// Spawn a new goroutine for subdirectories.
				wg.Add(1)
				go walkDir(fullPath, search, sizeFlag,  printDir, sem)
			} else {
				// Convert the file name to lower-case before comparing.
				lowerName := strings.ToLower(entry.Name())
		                // Check if the file name contains the search string.
		                // When search is empty, strings.Contains always returns true.
				if strings.Contains(lowerName, search) {
					if sizeFlag {
						// Check if the entry is a symlink: avoid calling os.Stat on symlinks.
						if entry.Type()&os.ModeSymlink != 0 {
							fileCh <- fmt.Sprintf("%s\tSYMLINK", fullPath)
						} else {
							// For regular files, call os.Stat to get the file size.

		                            		// Get file details (e.g., size) using os.Stat (Slow !)
							info, err := os.Stat(fullPath)
							if err != nil {
								log.Printf("failed to stat file %s: %v\n", fullPath, err)
								continue
							}

                            // Send the formatted output to the channel.
							fileCh <- fmt.Sprintf("%s\t%d", fullPath, info.Size())
						}
					} else {
						fileCh <- fullPath
					}
				}
			}
		}
	}

	// Start traversing from the root directory.
	wg.Add(1)
	go walkDir(root, lowerSearch, *sizeFlag, *printDirFlag, sem)

	// Close the channel once all goroutines have finished.
    // (so after all directories have been processed)
	go func() {
		wg.Wait()
		close(fileCh)
	}()

    // Print the results: matching files
	for fileInfo := range fileCh {
		fmt.Println(fileInfo)
	}
}
