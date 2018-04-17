package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	hostName := flag.String("h", "localhost", "hostname to serve on")
	port := flag.String("p", "8000", "port to serve on")
	flag.Parse()

	swaggerFile := "/swaggerspec/swagger.json"

	http.HandleFunc(swaggerFile, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, swaggerFile)
	})
	log.Printf("Serving swagger spec file here: http://%s:%s%s\n", *hostName, *port, swaggerFile)
	log.Printf("Serving Swagger UI here: http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
