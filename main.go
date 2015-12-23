package main

import (
	db "./todb"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	db.InitializeDatabase(os.Args[1], os.Args[2], os.Args[3])
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", index).Methods("GET")
	router.HandleFunc("/api/2.0/raw/{table}.json", handleTable)
	router.HandleFunc("/api/2.0/{cdn}/CRConfig.json", handleCRConfig)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleCRConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cdn := vars["cdn"]
	resp, _ := db.GetCRConfig(cdn)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func handleTable(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /api request")
	log.Println(r.UserAgent())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type")

	vars := mux.Vars(r)
	table := vars["table"]

	rows, _ := db.GetTable(table)
	// fmt.Print(rows)
	enc := json.NewEncoder(w)
	enc.Encode(rows)
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, "table:", table)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}
