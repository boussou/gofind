# gofind  (WIP)

Go find with  **directory Traversal using Goroutines**.

Search files which name contains a certain string (case-insensitive).

It will output each matching file with its full path (and size in bytes if requested)

Sample call:

```
gofind  -search propa  /tmp/sandbox
gofind -search contains -size /tmp/sandbox
gofind /tmp/sandbox -search contains
gofind /tmp/sandbox -search contains -size
```

parameters:
- search Search string to match in file names (case-insensitive)
- size  displays the file size right after the filename
- printdir:  Optionally prints directory names if -printdir is set
             It allow to totally mimic the output of the linux find command.


## Execute

1. **Using `go run`:**  
   ```bash
   go run main.go -search=pattern /path/to/directory
   ```
   This command runs the program directly. Replace `/path/to/directory` with the directory you want to search and `pattern` with the string to match in file names.

2. **Building a binary and then running it:**  

   Build the binary:
   ```bash
   go build 
   
   or specifying the executable name :
   
   go build -o gofind main.go
   ```

   Then execute the binary:
   ```bash
   ./gofind /path/to/directory -search=pattern
   ```
  Here specify the root directory to begin the search from and the `-search` pattern to filter file names, 



### Features

- **Directory Traversal with Goroutines:**  
  The `walkDir` function recursively processes each directory. When it finds a subdirectory, it launches a new goroutine for concurrent traversal.

- **Symlink Check:**  
  For each file, before calling os.Stat, the code checks if the entry is a symlink by using the bitwise check `entry.Type() & os.ModeSymlink`. If the file is a symlink, it avoids calling os.Stat and outputs the file path with "symlink" string instead.  
This helps avoid usual issues with symlinks (such as cycles or broken links).
  

- **File Size for Regular Files:**  
  If the file is not a symlink and its name matches the search string, os.Stat is called to retrieve the file size, which is then printed along with the file path.

- **Channel and WaitGroup:**  
  The matching file information is sent through a channel and printed in the main goroutine. A sync.WaitGroup ensures all goroutines complete before closing the channel.
