package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8888", "listen on the provided address and port")

	var dir string
	flag.StringVar(&dir, "dir", filepath.Join(homeDir, "public_html"), "directory of files")

	flag.Parse()

	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}
