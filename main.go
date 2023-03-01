package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var config struct {
	addr    string
	dir     string
	verbose bool
	lag     time.Duration
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&config.addr, "addr", ":8888", "listen on the provided address and port")
	flag.StringVar(&config.dir, "dir", filepath.Join(homeDir, "public_html"), "directory of files")
	flag.BoolVar(&config.verbose, "verbose", false, "print request bodies")
	flag.DurationVar(&config.lag, "lag", 0, "sleep for `duration` before each HTTP response")

	flag.Parse()

	log.Fatal(http.ListenAndServe(config.addr, lagRequest(logRequest(http.FileServer(http.Dir(config.dir))))))
}

func lagRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("sleeping for %v", config.lag)
		time.Sleep(config.lag)
		handler.ServeHTTP(w, r)
	})
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v", r.RemoteAddr, r.Method, r.URL)
		if config.verbose {
			log.Println(r.Header)
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
