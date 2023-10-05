package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asalih/go-ext"
)

func main() {
	// u14, err := os.Open("C:\\tmp\\raw_images\\ubuntu_20230920121211-1_image.001")
	u14, err := os.Open("C:\\tmp\\raw_images\\ubuntu1404_20230925221736-1_image.001")
	if err != nil {
		log.Fatalf("file err: %v", err)
	}

	fs, err := ext.NewFS(u14)
	if err != nil {
		log.Fatalf("stat err: %v", err)
	}

	ents, err := fs.ReadDir("/home")
	fmt.Println(err, ents)

}
