package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"path/filepath"
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

func isClassFolder(name string) bool {
    switch name {
    case "Images", "Videos", "Audio", "Archives", "Documents", "Other":
        return true
    default:
        return false
    }
}

var verbose bool
var deferOutput bool
var outputMessages []string

func printAndDeferOutput(message string) {
    if verbose {
        fmt.Println(message)
    }
    if deferOutput {
        outputMessages = append(outputMessages, message)
    }
}

func printOutputMessages() {
    for _, message := range outputMessages {
        fmt.Println(message)
    }
}

func organizeFile(directory string, file os.FileInfo, includeSubfolders bool) {
    if file.IsDir() {
        if includeSubfolders && !isClassFolder(file.Name()) {
            organizeDirectory(directory+"/"+file.Name(), includeSubfolders)
        }
        return
    }

    filePath := directory + "/" + file.Name()
    buf, err := ioutil.ReadFile(filePath)
    if err != nil {
        printAndDeferOutput("Error reading file " + file.Name() + ": " + err.Error())
        return
    }

    if len(buf) == 0 {
        printAndDeferOutput("Skipping empty file: " + file.Name())
        return
    }

    fileClass := getFileClass(buf)

    if filepath.Base(directory) == fileClass {
        return
    }

    subfolder := directory + "/" + fileClass
    if _, err := os.Stat(subfolder); os.IsNotExist(err) {
        os.Mkdir(subfolder, 0755)
    }

    if err := os.Rename(filePath, subfolder+"/"+file.Name()); err != nil {
        printAndDeferOutput("Error moving file " + file.Name() + ": " + err.Error())
    }
}

func organizeDirectory(directory string, includeSubfolders bool) (int, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
			return 0, err
	}
	printAndDeferOutput("Organizing directory: " + directory)
  printAndDeferOutput("Found " + strconv.Itoa(len(files)) + " entries in: " + directory)
	for _, file := range files {
			organizeFile(directory, file, includeSubfolders)
	}
	return len(files), nil
}

func main() {
	// TODO: implement flag behaviour
	directory := flag.String("dir", "sample-folder", "The directory to organize.")
	includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
  flag.BoolVar(&verbose, "verbose", false, "Print verbose output.")
  flag.BoolVar(&deferOutput, "defer-output", false, "Defer output until the end.")
	// log := flag.Bool("log", false, "Save a log of operations.")
	// exactMatch := flag.Bool("exact", false, "Organize by exact type ONLY.")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Error: The directory to organize is required.")
		return
	}

	fmt.Println("Starting organization...")

	count, err := organizeDirectory(*directory, *includeSubfolders)
	if err != nil {
			printAndDeferOutput("Error: Failed to read directory " + *directory + ": " + err.Error())
			return
	}

	printAndDeferOutput("Found " + strconv.Itoa(count) + " files in directory " + *directory)

	fmt.Println("Organization complete.")
	if deferOutput {
		fmt.Println("--------------------------------")
		printOutputMessages()
	}
}