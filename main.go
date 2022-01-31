package main

// https://wiki.osdev.org/ISO_9660

import (
	"fmt"
	"log"
	"os"

	"github.com/jesperkha/iso-reader/reader"
)

func main() {
	f, err := os.Open("output.iso")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fs, err := reader.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fs.Descriptor)
}
