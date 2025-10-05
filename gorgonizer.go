package main

import (
	"flag"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/pterm/pterm"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

var exact bool

func isClassFolder(name string) bool {
	switch name {
	case "Images", "Videos", "Audio", "Archives", "Documents", "Other":
		return true
	default:
		return false
	}
}

func isExactFolder(name string) bool {
	if name == "NoExt" {
		return true
	}
	if name == "" || isClassFolder(name) {
		return false
	}
	if strings.IndexByte(name, '.') >= 0 || strings.IndexByte(name, ' ') >= 0 {
		return false
	}
	return name == strings.ToUpper(name)
}

func boolMark(b bool) string {
	if b {
		return pterm.FgLightGreen.Sprint("✔")
	}
	return pterm.FgLightRed.Sprint("✖")
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

var verbose bool
var deferOutput bool
var outputMessages []string
var noColor bool
var log bool
var logMessages []string
var stats bool

type Stats struct {
	totalCount  int64
	totalBytes  int64
	countByType map[string]int64
	bytesByType map[string]int64
}

var statsAccumulator Stats

func initStats() {
	statsAccumulator = Stats{
		totalCount:  0,
		totalBytes:  0,
		countByType: make(map[string]int64),
		bytesByType: make(map[string]int64),
	}
}

func recordStats(typeKey string, sizeBytes int64) {
	if statsAccumulator.countByType == nil || statsAccumulator.bytesByType == nil {
		initStats()
	}
	if typeKey == "" || sizeBytes < 0 {
		return
	}
	statsAccumulator.totalCount += 1
	statsAccumulator.totalBytes += sizeBytes
	statsAccumulator.countByType[typeKey] = statsAccumulator.countByType[typeKey] + 1
	statsAccumulator.bytesByType[typeKey] = statsAccumulator.bytesByType[typeKey] + sizeBytes
}

func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return strconv.FormatInt(b, 10) + " B"
	}
	d := float64(b)
	units := []string{"KB", "MB", "GB", "TB"}
	for i := 0; i < len(units); i++ {
		d /= unit
		if d < unit {
			return fmt.Sprintf("%.2f %s", d, units[i])
		}
	}
	return fmt.Sprintf("%.2f %s", d, "PB")
}

func printStatsSummary() {
	styledTitle := pterm.FgLightCyan.Sprint("Stats Summary")
	lines := []string{}
	lines = append(lines, fmt.Sprintf("Total files: %d", statsAccumulator.totalCount))
	lines = append(lines, fmt.Sprintf("Total size: %s", humanizeBytes(statsAccumulator.totalBytes)))
	lines = append(lines, "")
	lines = append(lines, pterm.FgLightWhite.Sprint("Breakdown by type:"))
	for typeKey, count := range statsAccumulator.countByType {
		bytes := statsAccumulator.bytesByType[typeKey]
		lines = append(lines, fmt.Sprintf("- %s: %d files, %s", typeKey, count, humanizeBytes(bytes)))
	}
	pterm.DefaultBox.
		WithTitle(styledTitle).
		WithTitleTopLeft().
		WithLeftPadding(1).
		WithRightPadding(1).
		WithTopPadding(0).
		WithBottomPadding(0).
		Println(strings.Join(lines, "\n"))

	if log {
		logMessages = append(logMessages, timestamp()+" Stats => Total files: "+strconv.FormatInt(statsAccumulator.totalCount, 10)+
			" | Total size: "+humanizeBytes(statsAccumulator.totalBytes))
		for typeKey, count := range statsAccumulator.countByType {
			bytes := statsAccumulator.bytesByType[typeKey]
			logMessages = append(logMessages, timestamp()+" Stats Type => "+typeKey+": "+strconv.FormatInt(count, 10)+" files, "+humanizeBytes(bytes))
		}
	}
}

func printDeferLog(message string) {
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
	if log {
		messageWithTimestamp := timestamp() + " " + message
		logMessages = append(logMessages, messageWithTimestamp)
	}
}

func printOutputMessages() {
	for _, message := range outputMessages {
		pterm.Println(message)
	}
}

func copyDirectory(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove destination %s: %w", dst, err)
		}
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create destination %s: %w", dst, err)
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			if relPath == "." {
				return nil
			}
			return os.MkdirAll(targetPath, info.Mode())
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		return nil
	})
}

func organizeFile(directory string, file os.FileInfo, includeSubfolders bool) {
	if file.IsDir() {
		if includeSubfolders {
			if (!exact && isClassFolder(file.Name())) || (exact && isExactFolder(file.Name())) {
				return
			}
			organizeDirectory(filepath.Join(directory, file.Name()), includeSubfolders)
		}
		return
	}

	filePath := filepath.Join(directory, file.Name())
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		printDeferLog("Error reading file " + file.Name() + ": " + err.Error())
		return
	}

	if len(buf) == 0 {
		printDeferLog("Skipping empty file: " + file.Name())
		return
	}

	var destFolder string
	if exact {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == "" {
			destFolder = "NoExt"
		} else {
			destFolder = strings.ToUpper(strings.TrimPrefix(ext, "."))
		}
	} else {
		fileClass := getFileClass(buf)
		if filepath.Base(directory) == fileClass {
			return
		}
		destFolder = fileClass
	}

	if stats {
		recordStats(destFolder, file.Size())
	}

	if filepath.Base(directory) == destFolder {
		return
	}

	subfolder := filepath.Join(directory, destFolder)
	if _, err := os.Stat(subfolder); os.IsNotExist(err) {
		os.Mkdir(subfolder, 0755)
		printDeferLog("Created subfolder " + subfolder)
	}

	if err := os.Rename(filePath, filepath.Join(subfolder, file.Name())); err != nil {
		printDeferLog("Error moving file " + file.Name() + ": " + err.Error())
	} else {
		printDeferLog("Moved file " + file.Name() + " to " + subfolder)
	}
}

func organizeDirectory(directory string, includeSubfolders bool) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	printDeferLog("Organizing directory " + directory)
	printDeferLog("Found " + strconv.Itoa(len(files)) + " entries")
	for _, file := range files {
		organizeFile(directory, file, includeSubfolders)
	}
	return nil
}

func main() {
	directory := flag.String("dir", "", "The directory to organize.")
	includeSubfolders := flag.Bool("include-subfolders", false, "Include subfolders in the organization.")
	flag.BoolVar(&verbose, "verbose", false, "Print verbose output.")
	flag.BoolVar(&deferOutput, "defer-output", false, "Defer output until the end.")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored terminal output.")
	flag.BoolVar(&log, "log", false, "Save a log of operations.")
	flag.BoolVar(&exact, "exact", false, "Organize by exact type ONLY.")
	flag.BoolVar(&stats, "stats", false, "Print details about data organized.")
	detailed := flag.Bool("detailed", false, "Print detailed output (deferred output, log, and stats).")
	demo := flag.Bool("demo", false, "Run on copy of base-dummy-files.")

	flag.Parse()
	if *demo {
		if *directory == "" {
			*directory = "sample-folder"
		}
		if err := copyDirectory("base-dummy-files", *directory); err != nil {
			printDeferLog("Error preparing demo folder: " + err.Error())
		} else {
			printDeferLog("Prepared demo folder '" + *directory + "' from base-dummy-files")
		}
	}
	if *detailed {
		deferOutput = true
		log = true
		stats = true
	}
	if stats {
		initStats()
	}

	if noColor {
		pterm.DisableColor()
	}

	if *directory == "" {
		pterm.Error.Println("The directory to organize is required.")
		if log {
			logMessages = append(logMessages, timestamp()+" Error: The directory to organize is required.")
		}
		return
	}

	pterm.DefaultHeader.WithFullWidth().Println("Gorgonizer")
	options := fmt.Sprintf(
		"Directory: %s\nInclude subfolders %s\nVerbose %s\nDefer output %s\nNo color %s\nLog %s\nExact %s\nStats %s",
		*directory,
		boolMark(*includeSubfolders),
		boolMark(verbose),
		boolMark(deferOutput),
		boolMark(noColor),
		boolMark(log),
		boolMark(exact),
		boolMark(stats),
	)
	styledTitle := pterm.FgLightCyan.Sprint("Options")
	styledOptions := pterm.FgLightWhite.Sprint(options)
	pterm.DefaultBox.
		WithTitle(styledTitle).
		WithTitleTopLeft().
		WithLeftPadding(1).
		WithRightPadding(1).
		WithTopPadding(0).
		WithBottomPadding(0).
		Println(styledOptions)
	if log {
		logMessages = append(logMessages, timestamp()+" Options => Directory: "+*directory+
			" | Include subfolders: "+strconv.FormatBool(*includeSubfolders)+
			" | Verbose: "+strconv.FormatBool(verbose)+
			" | Defer output: "+strconv.FormatBool(deferOutput)+
			" | No color: "+strconv.FormatBool(noColor)+
			" | Log: "+strconv.FormatBool(log)+
			" | Exact: "+strconv.FormatBool(exact)+
			" | Stats: "+strconv.FormatBool(stats))
	}

	err := organizeDirectory(*directory, *includeSubfolders)
	if err != nil {
		printDeferLog("Error: Failed to read directory " + *directory + ": " + err.Error())
		return
	}

	pterm.Success.Println("Organization complete.")
	if log {
		logMessages = append(logMessages, timestamp()+" Organization complete.")
	}
	if deferOutput {
		pterm.FgDarkGray.Println("────────────────────────────────")
		printOutputMessages()
	}

	if stats {
		pterm.FgDarkGray.Println("────────────────────────────────")
		printStatsSummary()
	}

	if log {
		ioutil.WriteFile("log.txt", []byte(strings.Join(logMessages, "\n")), 0644)
		pterm.Success.Println("Log saved to " + "log.txt")
	}
}
