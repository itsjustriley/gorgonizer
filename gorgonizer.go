package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"github.com/h2non/filetype"
)

func getFileClass(buf []byte) string {
	if filetype.IsImage(buf) {
		return "Images"
	} else if filetype.IsVideo(buf) {
		return "Videos"
	} else if filetype.IsAudio(buf) {
		return "Audio"
	} else if filetype.IsArchive(buf) {
		return "Archives"
	} else if filetype.IsDocument(buf) {
		return "Documents"
	} else {
		return "Other"
	}
}

func organizeFile(directory string, file os.FileInfo, includeSubfolders bool) {
    if file.IsDir() {
        if includeSubfolders {
            organizeDirectory(directory+"/"+file.Name(), includeSubfolders)
        }
        return
    }

    filePath := directory + "/" + file.Name()
    buf, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
        return
    }

    if len(buf) == 0 {
        fmt.Printf("Skipping empty file: %s\n", file.Name())
        return
    }

    fileClass := getFileClass(buf)
    subfolder := directory + "/" + fileClass
    if _, err := os.Stat(subfolder); os.IsNotExist(err) {
        os.Mkdir(subfolder, 0755)
    }

    if err := os.Rename(filePath, subfolder+"/"+file.Name()); err != nil {
        fmt.Printf("Error moving file %s: %v\n", file.Name(), err)
    }
}

func organizeDirectory(directory string, includeSubfolders bool) (int, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
			return 0, err
	}
	fmt.Println("Organizing directory:", directory)
  fmt.Println("Found", len(files), "entries in:", directory)
	for _, file := range files {
			organizeFile(directory, file, includeSubfolders)
	}
	return len(files), nil
}

func main() {
	// TODO: implement flag behaviour
	directory := flag.String("dir", "sample-folder", "The directory to organize.")
	includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
	// log := flag.Bool("log", false, "Save a log of operations.")
	// exactMatch := flag.Bool("exact", false, "Organize by exact type ONLY.")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Error: The directory to organize is required.")
		return
	}

    count, err := organizeDirectory(*directory, *includeSubfolders)
    if err != nil {
        fmt.Println("Error: Failed to read directory.")
        return
    }

    fmt.Println("Found", count, "files in directory.")

	fmt.Println("Organization complete.")
}