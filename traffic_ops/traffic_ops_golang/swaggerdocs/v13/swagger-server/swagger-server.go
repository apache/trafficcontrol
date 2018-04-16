package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8000", "port to serve on")
	flag.Parse()

	swaggerFile := "swagger.json"

	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, swaggerFile)
	})
	log.Printf("Serving %s on HTTP port: %s\n", swaggerFile, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
