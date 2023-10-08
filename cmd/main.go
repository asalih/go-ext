package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asalih/go-ext"
)

func main() {
	extFile, err := os.Open(os.Getenv("EXT_PATH"))
	if err != nil {
		log.Fatalf("file err: %v", err)
	}

	fs, err := ext.NewFS(extFile)
	if err != nil {
		log.Fatalf("stat err: %v", err)
	}

	ents, err := fs.ReadDir("/dev")
	fmt.Println(err, ents)

}
