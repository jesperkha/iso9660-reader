package cmd

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/jesperkha/iso-reader/reader"
)

var (
	ErrExpectedFilename = errors.New("expected filename")
	ErrUnknownCommand   = errors.New("unknown command. use 'help' to view a list of commands")

	//go:embed help.txt
	helpMessage string
	scanner     = bufio.NewScanner(os.Stdin)
	writer      = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
)

// Runs a simple cli program to navigate around the disk and open files.
// Also allows to extract a file from the disk with the "get" command.
func RunTerminalMode(fs *reader.FileSystem) {
	currentDir := []string{""}
	fmt.Println("type 'exit' to terminate session")

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

		newDir, err := runCommand(fs, currentDir, command, args[1:])
		if err != nil {
			printError(err.Error())
		}

		currentDir = newDir
	}
}

func runCommand(fs *reader.FileSystem, dir []string, command string, args []string) (newDir []string, err error) {
	path := strings.Join(dir, "/")

	if command == "ls" {
		dirs, err := fs.ReadDirectory(path)
		if err != nil {
			return dir, err
		}

		totalSize := 0
		for _, d := range dirs {
			// Both the . and .. directory should be hidden
			if d.Name == "." || d.Name == ".." {
				continue
			}

			slash := "/"
			if d.IsFile {
				slash = ""
			}

			totalSize += d.ExtentSize
			extentSize := formatFileSize(d.ExtentSize)
			date, time := d.Date.FormatDate(), d.Date.FormatTime()
			fmt.Fprintf(writer, "%s \t %s \t %s \t %s%s \n", extentSize, date, time, d.Name, slash)
		}

		fmt.Printf("total %s\n", formatFileSize(totalSize))
		writer.Flush()
		return dir, err
	}

	if command == "cd" {
		if len(args) == 0 {
			return []string{""}, err
		}

		path := strings.Split(args[0], "/")
		for _, p := range path {
			switch p {
			case ".":
				continue
			case "..":
				if len(dir) > 1 {
					dir = (dir)[:len(dir)-1]
				}
				continue
			}

			newPath := append(dir, p)
			if _, err := fs.ReadDirectory(strings.Join(newPath, "/")); err != nil {
				return dir, err
			}

			dir = newPath
		}

		return dir, err
	}

	if command == "open" {
		if len(args) == 0 {
			return dir, ErrExpectedFilename
		}

		filepath := fmt.Sprintf("%s/%s", strings.Join(dir, "/"), args[0])
		file, err := fs.ReadFile(filepath)
		if err != nil {
			return dir, err
		}

		fmt.Println(file.String())
		return dir, err
	}

	// Extract a file from the disk. Keeps the files name
	if command == "get" {
		if len(args) == 0 {
			return dir, ErrExpectedFilename
		}

		filepath := fmt.Sprintf("%s/%s", strings.Join(dir, "/"), args[0])
		file, err := fs.ReadFile(filepath)
		if err != nil {
			return dir, err
		}

		output, err := os.Create(file.Name)
		if err != nil {
			return dir, err
		}

		defer output.Close()
		_, err = output.Write(file.Bytes)
		return dir, err
	}

	return dir, ErrUnknownCommand
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
