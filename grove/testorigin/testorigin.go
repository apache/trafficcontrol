package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func GetRandPage(pageLenBytes int) []byte {
	page := make([]byte, pageLenBytes, pageLenBytes)
	rand.Seed(time.Now().Unix())
	rand.Read(page)
	return page
}

func Handle(w http.ResponseWriter, r *http.Request, page []byte) {
	fmt.Printf("%v %v %v\n", time.Now(), r.RemoteAddr, r.RequestURI)
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func main() {
	port := flag.Int("port", -1, "The port to serve on")
	pageBytes := flag.Int("pagebytes", -1, "The number of random bytes to serve as a page")
	flag.Parse()
	if *port < 0 || *pageBytes < 0 {
		fmt.Printf("usage: testorigin -port 8080 -pagebytes 350000\n")
		os.Exit(1)
	}

	page := GetRandPage(*pageBytes)

	fmt.Printf("Serving on %d\n", *port)

	handle := func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r, page)
	}

	http.HandleFunc("/", handle)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		os.Exit(1)
	}
}
