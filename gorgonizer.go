package main

import (
	"fmt"
	"flag"
	"io/ioutil"
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

	// TODO: implement organization logic
	files, err := ioutil.ReadDir(*directory)
	if err != nil {
		fmt.Println("Error: Failed to read directory.")
		return
	}

	fmt.Println("Found", len(files), "files in directory.")

	fmt.Println("Organization complete.")
}