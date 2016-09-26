package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func MustLoadCertificates(cafile string) *x509.CertPool {
	pem, err := ioutil.ReadFile(cafile)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return certPool
}

func MustGetTlsConfiguration(servername string, cafile string) *tls.Config {
	config := &tls.Config{}
	certPool := MustLoadCertificates(cafile)

	config.RootCAs = certPool
	config.ClientCAs = certPool

	config.ClientAuth = tls.RequireAndVerifyClientCert

	config.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}

	//Use only TLS v1.2
	config.MinVersion = tls.VersionTLS12

	//Don't allow session resumption
	config.SessionTicketsDisabled = true
	return config
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwzyz0123456789/"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func DoRequest(tlsConfig *tls.Config, host string) {
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}

	url := fmt.Sprintf("https://ccr.%v/%v/stuff?fakeClientIpAddress=68.87.25.123", host, RandomString(23))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error", err)
	}

	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 302 {
			fmt.Println("Received unexpected http status", resp.StatusCode)
		}
	}
}

func ExerciseDeliveryService(tlsConfig *tls.Config, host string, numRequests int, maxWorkers int) {
	fmt.Println("Sending", numRequests, "requests to", host)
	var waitGroup sync.WaitGroup

	workers := make(chan int, maxWorkers)
	start := time.Now()
	for worker := 0; worker < maxWorkers; worker++ {
		waitGroup.Add(1)

		go func(ch <-chan int) {
			defer waitGroup.Done()

			for _ = range ch {
				DoRequest(tlsConfig, host)
			}
		}(workers)
	}

	for i := 0; i < numRequests; i++ {
		if i%maxWorkers == 0 {
			fmt.Println(host, i)
		}
		time.Sleep(100 * time.Microsecond)
		workers <- i
	}
	close(workers)

	fmt.Println(host, "waiting")
	waitGroup.Wait()
	elapsed := time.Since(start)
	fmt.Println(host, "Took", elapsed)
}

func main() {
	txCountFlag := flag.Int("txcount", 100, "Number of transactions to execute per delivery service")
	workersFlag := flag.Int("workers", 10, "Number of concurrent d")
	cafileFlag := flag.String("cafile", "ca.crt", "root Certificate Authority file")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "\tloadtest")
		fmt.Fprintln(os.Stderr, "\t\t[-cafile <file name>]")
		fmt.Fprintln(os.Stderr, "\t\t[-txcount <number of transactions per delivery service>]")
		fmt.Fprintln(os.Stderr, "\t\t[-workers <number of concurrent requests per delivery service>]")
		fmt.Fprintln(os.Stderr, "\t\tmy-cdn.example.com delivery-service-1 delivery-service-2 ...")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tSome example locations of trusted certificate files to use for -cafile are:")
		fmt.Fprintln(os.Stderr, "/etc/ssl/certs");
		fmt.Fprintln(os.Stderr, "/etc/pki/tls/certs/ca-bundle.crt");
		fmt.Fprintln(os.Stderr, "/etc/ssl/certs/ca-bundle.crt");
		fmt.Fprintln(os.Stderr, "/etc/pki/tls/certs/ca-bundle.trust.crt");
		fmt.Fprintln(os.Stderr, "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem");
		fmt.Fprintln(os.Stderr, "/System/Library/OpenSSL");
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
	}
	flag.Parse()

	if _, err := os.Stat(*cafileFlag); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Loadtest requires a root certificate authority file, by default it looks for ca.crt in the current working directory")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	cdnDomain := flag.Args()[0]
	deliveryServices := flag.Args()[1:]

	fmt.Println(deliveryServices)

	var waitGroup sync.WaitGroup

	start := time.Now()
	for _, deliveryService := range deliveryServices {
		waitGroup.Add(1)
		go func(ds string) {
			defer waitGroup.Done()
			host := strings.Join([]string{ds, cdnDomain}, ".")
			fmt.Println(host)
			tlsConfig := MustGetTlsConfiguration(host, *cafileFlag)
			ExerciseDeliveryService(tlsConfig, host, *txCountFlag, *workersFlag)
		}(deliveryService)
	}

	waitGroup.Wait()
	elapsed := time.Since(start)
	totalRequests := *txCountFlag * len(deliveryServices)
	fmt.Println(totalRequests)
	tps := float64(totalRequests) / elapsed.Seconds()
	fmt.Println("Took", elapsed)
	fmt.Println("TPS", tps)
}
