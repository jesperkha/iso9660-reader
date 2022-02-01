package main

import (
	"fmt"
	"log"

	"github.com/jesperkha/iso-reader/reader"
)

// https://wiki.osdev.org/ISO_9660

func main() {
	fs, err := reader.ReadDisk("output.iso")
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Close()

	file, err := fs.ReadFile("folder/myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file.String())
}
