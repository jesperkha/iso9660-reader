package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesperkha/iso-reader/cmd"
	"github.com/jesperkha/iso-reader/reader"
)

// https://wiki.osdev.org/ISO_9660

func main() {
	if len(os.Args) == 1 {
		fmt.Println("error: expected filename")
		return
	}

	fs, err := reader.ReadDisk(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Close()
	cmd.RunTerminalMode(fs)
}
