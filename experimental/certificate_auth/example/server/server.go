package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", HelloHandler)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}

	server := http.Server{
		Addr:      "server.local:8443",
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	if err := server.ListenAndServeTLS("../certs/server.crt.pem", "../certs/server.key.pem"); err != nil {
		log.Fatalf("error listening to port: %v", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	if r.TLS.PeerCertificates != nil {
		clientCert := r.TLS.PeerCertificates[0]
		fmt.Println("Client cert subject: ", clientCert.Subject)
	}

	fmt.Println("Hello")
}
