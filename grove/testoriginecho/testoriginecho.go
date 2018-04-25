package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v %v %v\n", time.Now(), r.RemoteAddr, r.RequestURI)
		// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%v %v %v\n", time.Now(), r.RemoteAddr, r.RequestURI)))
}

func main() {
	port := flag.Int("port", -1, "The port to serve on")
	flag.Parse()
	if *port < 0 {
		fmt.Printf("usage: testorigin -port 8080\n")
		os.Exit(1)
	}

	fmt.Printf("Serving on %d\n", *port)

	handle := func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r)
	}

	http.HandleFunc("/", handle)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		os.Exit(1)
	}
}
