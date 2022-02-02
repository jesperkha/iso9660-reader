package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/jesperkha/iso-reader/reader"
)

//go:embed help.txt
var helpMessage string
var scanner = bufio.NewScanner(os.Stdin)

// Runs a simple cli program to navigate around the disk and open text files.
// Also allows to extract a file from the disk with the "get" command.
func RunTerminalMode(fs *reader.FileSystem) {
	currentDir := []string{""}
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, '	', 0)

	for {
		path := strings.Join(currentDir, "/")
		ct.Foreground(ct.Blue, true)
		fmt.Print("\niso-reader ")
		ct.Foreground(ct.Green, false)
		fmt.Printf("%s ", fs.Descriptor.VolumeIdentifier)
		ct.Foreground(ct.Yellow, false)
		fmt.Printf("~%s\n", path)
		ct.ResetColor()
		fmt.Printf("$ ")

		scanner.Scan()
		args := strings.Split(scanner.Text(), " ")
		command := args[0]
		if command == "exit" {
			break
		}

		if command == "help" {
			fmt.Println(helpMessage)
			continue
		}

		if command == "ls" {
			dirs, err := fs.ReadDirectory(path)
			if err != nil {
				printError(err.Error())
				continue
			}

			totalSize := 0
			for _, d := range dirs {
				slash := "/"
				if d.IsFile {
					slash = ""
				}

				totalSize += d.ExtentSize
				extentSize := formatFileSize(d.ExtentSize)
				date, time := d.Date.FormatDate(), d.Date.FormatTime()
				fmt.Fprintf(writer, "%s %s %s %s%s\n", extentSize, date, time, d.Name, slash)
			}

			fmt.Printf("total %s\n", formatFileSize(totalSize))
			writer.Flush()
			continue
		}

		if command == "cd" {
			if len(args) == 1 {
				currentDir = []string{""}
				continue
			}

			path := strings.Split(args[1], "/")
			for _, p := range path {
				switch p {
				case ".":
					continue
				case "..":
					if len(currentDir) > 1 {
						currentDir = currentDir[:len(currentDir)-1]
					}
					continue
				}

				newPath := append(currentDir, p)
				if _, err := fs.ReadDirectory(strings.Join(newPath, "/")); err != nil {
					printError(err.Error())
					continue
				}

				currentDir = newPath
			}

			continue
		}

		if command == "open" {
			if len(args) == 1 {
				printError("expected filename after 'open'")
				continue
			}

			filepath := fmt.Sprintf("%s/%s", strings.Join(currentDir, "/"), args[1])
			file, err := fs.ReadFile(filepath)
			if err != nil {
				printError(err.Error())
				continue
			}

			fmt.Println(file.String())
			continue
		}

		if command == "get" {
			printError("not implemented yet")
			continue
		}

		printError("error: uknown command")
	}
}

// Formats file size to nearest thousand
func formatFileSize(size int) string {
	unit := ""
	units := []string{"K", "M", "G"}
	for i := 0; size > 1000; i++ {
		size /= 1000
		unit = units[i]
	}

	return fmt.Sprintf("%d%s", size, unit)
}

func printError(msg string) {
	ct.Foreground(ct.Red, false)
	fmt.Println(msg)
	ct.ResetColor()
}
