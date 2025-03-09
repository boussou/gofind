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

// search file which name contains a string (case-insensitive)
// sample calls:
//   gofind -search contains -excludeFile tmp -size /tmp/sandbox
//   gofind /tmp/sandbox -search contains -excludeFile tmp
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

	// Define flags for search string, size display and exclusion.
	search := flag.String("search", "", "Search string to match in file names (case-insensitive)")
	sizeFlag := flag.Bool("size", false, "Display file size if set")
	excludeFile := flag.String("excludeFile", "", "Exclude files whose names contain this string (case-insensitive)")
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
	// Convert search and exclude strings to lower-case for case-insensitive matching.
	lowerSearch := strings.ToLower(*search)
	
	// Display the search parameter only if it is not empty.
	if lowerSearch != "" {
		fmt.Printf("Search parameter: %s\n", lowerSearch)
	}
	if excludeFile != "" {
		fmt.Printf("Exclude parameter: %s\n", excludeFile)
	}

    //...............//...............//...............

    // Load excluded directory names from the exclude file (if provided).
    excludedDirs := make(map[string]bool)
    if excludeFile != "" {
        data, err := os.ReadFile(excludeFile)
        if err != nil {
            log.Fatalf("failed to read exclude file %s: %v", excludeFile, err)
        }
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            dirName := strings.TrimSpace(line)
            if dirName != "" {
                excludedDirs[dirName] = true
            }
        }
    }
    //...............//...............//...............



	// Channel to send matching file info.
	fileCh := make(chan string)
	var wg sync.WaitGroup

	// walkDir recursively processes the given directory and its subdirectories concurrently.
	// It now takes an additional parameter, exclude, that skips files whose names contain that substring.
	var walkDir func(dir, search, exclude string, sizeFlag bool)
	walkDir = func(dir, search, exclude string, sizeFlag bool) {
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
				// Exclude directories starting with ".cache"
		                //if strings.HasPrefix(entry.Name(), ".cache") {
		                //    continue
		                //}  

                // Check if the directory should be excluded.
                if excludedDirs[entry.Name()] {
                    continue
                }

				// Spawn a new goroutine for subdirectories.
				wg.Add(1)
				go walkDir(fullPath, search, exclude, sizeFlag)

			} else {
				// Convert the file name to lower-case before comparing.
				lowerName := strings.ToLower(entry.Name())

				// Check if file matches the search criteria.
				if strings.Contains(lowerName, search) {
				
                    // If exclude is provided and the file name contains it, skip this file.
					//if exclude != "" && strings.Contains(lowerName, exclude) {
					//	continue
					//}

					if sizeFlag {
						// Check if the entry is a symlink.
						if entry.Type()&os.ModeSymlink != 0 {
							// For symlinks, avoid calling os.Stat and simply note it's a symlink.
							fileCh <- fmt.Sprintf("%s\tsymlink", fullPath)
						} else {
							// For non-symlink files, call os.Stat to get the file size.

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

	// Start traversing from the root directory the directory tree    
	wg.Add(1)
	go walkDir(root, lowerSearch, excludeFile, *sizeFlag)

	// Close the channel once all goroutines have finished
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
