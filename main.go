package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var config struct {
	addr    string
	dir     string
	verbose bool
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&config.addr, "addr", ":8888", "listen on the provided address and port")
	flag.StringVar(&config.dir, "dir", filepath.Join(homeDir, "public_html"), "directory of files")
	flag.BoolVar(&config.verbose, "verbose", false, "print request bodies")

	flag.Parse()

	log.Fatal(http.ListenAndServe(config.addr, logRequest(http.FileServer(http.Dir(config.dir)))))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v", r.RemoteAddr, r.Method, r.URL)
		if config.verbose {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			defer r.Body.Close()

			log.Println(string(body))
		}
		handler.ServeHTTP(w, r)
	})
}
