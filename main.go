package main

import (
	"fmt"
	"log"
	"os"
)

const path = "C:/Users/jesper/Downloads/linuxmint-20.3-cinnamon-64bit.iso"

func prettyPrint(data map[string]interface{}) {
	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func main() {
	f, err := os.Open("output.iso")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 96
	// 682
	// 50

	// data, err := reader.ReadPrimaryDescriptor(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// prettyPrint(data)
	
	buf := make([]byte, 64)
	_, err = f.ReadAt(buf, 18 * 2048 + 2 * 2048)
	if err != nil {
		log.Fatal(err)
	}

	// .
	// ..
	// folder

	// .
	// ..
	// myfile.txt

	// pos 16
	// len 2048
	// fmt.Println(binary.LittleEndian.Uint32(buf[10:14]))
	fmt.Println(string(buf))
	// fmt.Println(int(buf[0]))
	// fmt.Println(string(buf[33:33+12]))
}
