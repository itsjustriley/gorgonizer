package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"path/filepath"
	"github.com/h2non/filetype"
	"github.com/pterm/pterm"
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
var noColor bool

func printAndDeferOutput(message string) {
    if verbose {
        switch {
        case len(message) >= 5 && message[:5] == "Error":
            pterm.Error.Println(message)
        case len(message) >= 8 && message[:8] == "Skipping":
            pterm.Warning.Println(message)
        default:
            pterm.Info.Println(message)
        }
    }
    if deferOutput {
        var styled string
        switch {
        case len(message) >= 5 && message[:5] == "Error":
            styled = pterm.Error.Sprint(message)
        case len(message) >= 8 && message[:8] == "Skipping":
            styled = pterm.Warning.Sprint(message)
        default:
            styled = pterm.Info.Sprint(message)
        }
        outputMessages = append(outputMessages, styled)
    }
}

func printOutputMessages() {
    for _, message := range outputMessages {
        pterm.Println(message)
    }
}

func organizeFile(directory string, file os.FileInfo, includeSubfolders bool) {
    if file.IsDir() {
        if includeSubfolders && !isClassFolder(file.Name()) {
            organizeDirectory(filepath.Join(directory, file.Name()), includeSubfolders)
        }
        return
    }

    filePath := filepath.Join(directory, file.Name())
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

    subfolder := filepath.Join(directory, fileClass)
    if _, err := os.Stat(subfolder); os.IsNotExist(err) {
        os.Mkdir(subfolder, 0755)
    }

    if err := os.Rename(filePath, filepath.Join(subfolder, file.Name())); err != nil {
        printAndDeferOutput("Error moving file " + file.Name() + ": " + err.Error())
    }
}

func organizeDirectory(directory string, includeSubfolders bool) (error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
			return err
	}
	printAndDeferOutput("Organizing directory" + directory)
  printAndDeferOutput("Found " + strconv.Itoa(len(files)) + " entries")
	for _, file := range files {
			organizeFile(directory, file, includeSubfolders)
	}
	return nil
}

func main() {
	// TODO: implement flag behaviour
	directory := flag.String("dir", "sample-folder", "The directory to organize.")
	includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
  flag.BoolVar(&verbose, "verbose", false, "Print verbose output.")
  flag.BoolVar(&deferOutput, "defer-output", false, "Defer output until the end.")
  flag.BoolVar(&noColor, "no-color", false, "Disable colored terminal output.")
	// log := flag.Bool("log", false, "Save a log of operations.")
	// exactMatch := flag.Bool("exact", false, "Organize by exact type ONLY.")
	flag.Parse()

    if noColor {
        pterm.DisableColor()
    }

    if *directory == "" {
        pterm.Error.Println("The directory to organize is required.")
		return
	}

    pterm.DefaultHeader.WithFullWidth().Println("Gorgonizer")
    pterm.FgLightWhite.Println("Starting organization…")
    pterm.FgDarkGray.Println("Directory:", *directory, "| Include subfolders:", strconv.FormatBool(*includeSubfolders))

	err := organizeDirectory(*directory, *includeSubfolders)
	if err != nil {
			printAndDeferOutput("Error: Failed to read directory " + *directory + ": " + err.Error())
			return
	}

    pterm.Success.Println("Organization complete.")
	if deferOutput {
        pterm.FgDarkGray.Println("────────────────────────────────")
		printOutputMessages()
	}
}