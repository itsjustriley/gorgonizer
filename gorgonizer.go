package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"github.com/h2non/filetype"
)

func main() {
	// TODO: implement flag behaviour
	directory := flag.String("dir", "sample-folder", "The directory to organize.")
	// includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
	// log := flag.Bool("log", false, "Save a log of operations.")
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
		filePath := *directory + "/" + file.Name()
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}
		
		kind, err := filetype.Match(buf)
		if err != nil {
			fmt.Printf("Error matching file type for %s: %v\n", file.Name(), err)
			continue
		}
		
		if kind == filetype.Unknown {
			fmt.Printf("%s: Unknown file type\n", file.Name())
		} else {
			fmt.Printf("%s: %s (%s)\n", file.Name(), kind.MIME.Value, kind.Extension)
		}
	}

	fmt.Println("Found", len(files), "files in directory.")

	fmt.Println("Organization complete.")
}