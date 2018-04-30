package main

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
