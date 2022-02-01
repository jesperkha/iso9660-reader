package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jesperkha/iso-reader/reader"
)

// https://wiki.osdev.org/ISO_9660

func main() {
	f, err := os.Open("output.iso")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fs, err := reader.ReadDisk(f)
	if err != nil {
		log.Fatal(err)
	}

	file, err := fs.ReadFile("folder/myfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file.String())
}
