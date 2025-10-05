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

// TODO: extract organize directory into a function so it can be recursive if subfolders included
// TODO: extract file organization into a function for readability and testing?

func main() {
	// TODO: implement flag behaviour
	directory := flag.String("dir", "sample-folder", "The directory to organize.")
	// includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
	// log := flag.Bool("log", false, "Save a log of operations.")
	// exactMatch := flag.Bool("exact", false, "Organize by exact type ONLY.")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Error: The directory to organize is required.")
		return
	}

	fmt.Println("Organizing directory:", *directory)

	files, err := ioutil.ReadDir(*directory)
	if err != nil {
		fmt.Println("Error: Failed to read directory.")
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		filePath := *directory + "/" + file.Name()
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}
		
		if len(buf) == 0 {
			fmt.Printf("Skipping empty file: %s\n", file.Name())
			continue
		}
		
		fileClass := getFileClass(buf)
			subfolder := *directory + "/" + fileClass
			if _, err := os.Stat(subfolder); os.IsNotExist(err) {
				os.Mkdir(subfolder, 0755)
			}
			os.Rename(filePath, subfolder+"/"+file.Name())
		}

	fmt.Println("Found", len(files), "files in directory.")

	fmt.Println("Organization complete.")
}