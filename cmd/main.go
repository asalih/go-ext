package main

import (
	"fmt"
	"io"
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

	f, err := fs.Open("/opt/binalyze/air/agent/config.yml")
	if err != nil {
		log.Fatalf("open err: %v", err)
	}

	all, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("open err: %v", err)
	}

	// os.WriteFile("C:\\tmp\\wtf.txt", all, os.ModePerm)
	fmt.Println(len(all))

	st, err := fs.Stat("/")
	if err != nil {
		log.Fatalf("open err: %v", err)
	}
	fmt.Println(st.Name(), st.IsDir(), st.ModTime())

}
