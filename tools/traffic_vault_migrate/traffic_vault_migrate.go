package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
)

var (
	fromType string
	toType   string
	fromCfg  string
	toCfg    string
	dry      bool
	compare  bool
	confirm  bool
	dump     bool

	riakBE RiakBackend = RiakBackend{}
	pgBE   PGBackend   = PGBackend{}
)

func init() {
	flag.StringVar(&fromType, "from_type", riakBE.Name(), fmt.Sprintf("From server types (%v)", strings.Join(supportedTypes(), "|")))
	flag.StringVar(&toType, "to_type", pgBE.Name(), fmt.Sprintf("From server types (%v)", strings.Join(supportedTypes(), "|")))
	flag.StringVar(&toCfg, "to_cfg", "pg.json", "To server config file")
	flag.StringVar(&fromCfg, "from_cfg", "riak.json", "From server config file")
	flag.BoolVar(&dry, "dry", false, "Do not perform writes")
	flag.BoolVar(&compare, "compare", false, "Compare to and from server records")
	flag.BoolVar(&confirm, "confirm", true, "Requires confirmation before inserting records")
	flag.BoolVar(&dump, "dump", false, "Write keys (from 'from' server) to disk")
}

func main() {
	flag.Parse()
	var fromSrv TVBackend
	var toSrv TVBackend

	//fromType = "PG"
	//toType = "Riak"
	//fromCfg = "pg.json"
	//toCfg = "riak.json"

	toSrvUsed := !dump && !dry

	if !validateType(fromType) {
		log.Fatal("Unknown fromType " + fromType)
	}
	if toSrvUsed && !validateType(toType) {
		log.Fatal("Unknown toType " + toType)
	}

	fromSrv = getBackendFromType(fromType)
	if toSrvUsed {
		toSrv = getBackendFromType(toType)
	}

	var toTimer time.Time
	var toTime float64
	var fromTimer time.Time
	var fromTime float64

	log.Println("Reading configs...")
	fromTimer = time.Now()
	if err := fromSrv.ReadConfig(fromCfg); err != nil {
		log.Fatalf("Unable to read fromSrv cfg: %v", err)
	}
	fromTime = time.Now().Sub(fromTimer).Seconds()

	if toSrvUsed {
		toTimer = time.Now()
		if err := toSrv.ReadConfig(toCfg); err != nil {
			log.Fatalf("Unable to read toSrv cfg: %v", err)
		}
		toTime := time.Now().Sub(toTimer).Seconds()
		log.Printf("Done [%v seconds]\n\tto: [%v seconds]\n\tfrom: [%v seconds]\n", toTime+fromTime, toTime, fromTime)
	} else {
		log.Printf("Done [%v seconds]\n", fromTime)
	}

	log.Println("Starting servers...")
	fromTimer = time.Now()
	if err := fromSrv.Start(); err != nil {
		log.Fatalf("issue starting fromSrv: %v", err)
	}
	fromTime = time.Now().Sub(fromTimer).Seconds()
	defer func() {
		fromSrv.Stop()
	}()
	if toSrvUsed {
		toTimer = time.Now()
		if err := toSrv.Start(); err != nil {
			log.Fatalf("issue starting toSrv: %v", err)
		}
		toTime = time.Now().Sub(toTimer).Seconds()
		defer func() {
			toSrv.Stop()
		}()
		log.Printf("Done [%v seconds]\n\tto: [%v seconds]\n\tfrom: [%v seconds]\n", toTime+fromTime, toTime, fromTime)
	} else {
		log.Printf("Done [%v seconds]\n", fromTime)
	}

	log.Println("Pinging servers...")
	fromTimer = time.Now()
	if err := fromSrv.Ping(); err != nil {
		log.Fatalf("Unable to ping fromSrv: %v", err)
	}
	fromTime = time.Now().Sub(fromTimer).Seconds()
	if toSrvUsed {
		toTimer = time.Now()
		if err := toSrv.Ping(); err != nil {
			log.Fatalf("Unable to ping toSrv: %v", err)
		}
		toTime = time.Now().Sub(toTimer).Seconds()
		log.Printf("Done [%v seconds]\n\tto: [%v seconds]\n\tfrom: [%v seconds]\n", toTime+fromTime, toTime, fromTime)
	} else {
		log.Printf("Done [%v seconds]\n", fromTime)
	}

	log.Printf("Fetching data from %v...\n", fromSrv.Name())
	fromTimer = time.Now()
	if err := fromSrv.Fetch(); err != nil {
		log.Fatalf("Unable to fetch fromSrv data: %v", err)
	}

	fromSSLKeys, fromDNSSecKeys, fromURIKeys, fromURLKeys, err := GetKeys(fromSrv)
	if err != nil {
		log.Fatal(err)
	}

	if err := Validate(fromSrv); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done [%v seconds]\n", time.Now().Sub(fromTimer).Seconds())

	if compare {
		log.Printf("Fetching data from %v...\n", toSrv.Name())
		toTimer = time.Now()
		if err := toSrv.Fetch(); err != nil {
			log.Fatalf("Unable to fetch toSrv data: %v\n", err)
		}

		toSSLKeys, toDNSSecKeys, toURIKeys, toURLKeys, err := GetKeys(toSrv)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Validating " + toSrv.Name())
		if err := toSrv.ValidateKey(); err != nil && len(err) > 0 {
			log.Fatal(strings.Join(err, "\n"))
		}
		log.Printf("Done [%v seconds]\n", time.Now().Sub(toTimer).Seconds())
		log.Println(fromSrv.String())
		log.Println(toSrv.String())

		if !reflect.DeepEqual(fromSSLKeys, toSSLKeys) {
			log.Fatal("from sslkeys and to sslkeys don't match")
		}
		if !reflect.DeepEqual(fromDNSSecKeys, toDNSSecKeys) {
			log.Fatal("from dnssec and to dnssec don't match")
		}
		if !reflect.DeepEqual(fromURIKeys, toURIKeys) {
			log.Fatal("from uri and to uri don't match")
		}
		if !reflect.DeepEqual(fromURLKeys, toURLKeys) {
			log.Fatal("from url and to url don't match")
		}
		log.Println("Both datasources have the same keys!")
		return
	}

	log.Printf("Setting %v keys...\n", toSrv.Name())
	toTimer = time.Now()
	if err := SetKeys(toSrv, fromSSLKeys, fromDNSSecKeys, fromURIKeys, fromURLKeys); err != nil {
		log.Fatal(err)
	}

	if err := Validate(toSrv); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done [%v seconds]\n", time.Now().Sub(toTimer).Seconds())

	log.Println(fromSrv.String())

	if dry {
		return
	}

	if confirm {
		ans := "q"
		for {
			fmt.Print("Confirm data insertion (y/n):")
			if _, err := fmt.Scanln(&ans); err != nil {
				log.Fatal("unable to get user input")
			}

			if ans == "y" {
				break
			} else if ans == "n" {
				return
			}
		}
	}
	log.Printf("Inserting data into %v...\n", toSrv.Name())
	toTimer = time.Now()
	if err := toSrv.Insert(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done [%v seconds]\n", time.Now().Sub(toTimer).Seconds())

}

// Validate runs the ValidateKey method on the backend
func Validate(be TVBackend) error {
	if errs := be.ValidateKey(); errs != nil && len(errs) > 0 {
		return errors.New(fmt.Sprintf("Validation Errors (%v): \n%v", be.Name(), strings.Join(errs, "\n")))
	} else {
		log.Println("Validated " + be.Name())
	}
	return nil
}

// SetKeys will set all of the keys for a backend
func SetKeys(be TVBackend, sslkeys []SSLKey, dnssecKeys []DNSSecKey, uriKeys []URISignKey, urlKeys []URLSigKey) error {
	if err := be.SetSSLKeys(sslkeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v ssl keys: %v", be.Name(), err))
	}
	if err := be.SetDNSSecKeys(dnssecKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v dnssec keys: %v", be.Name(), err))
	}
	if err := be.SetURLSigKeys(urlKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v url keys: %v", be.Name(), err))
	}
	if err := be.SetURISignKeys(uriKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v uri keys: %v", be.Name(), err))
	}
	return nil
}

// GetKeys will get all of the keys for a backend
func GetKeys(be TVBackend) ([]SSLKey, []DNSSecKey, []URISignKey, []URLSigKey, error) {
	var sslkeys []SSLKey
	var dnssec []DNSSecKey
	var uri []URISignKey
	var url []URLSigKey
	var err error
	if sslkeys, err = be.GetSSLKeys(); err != nil {
		return nil, nil, nil, nil, errors.New(fmt.Sprintf("Unable to get %v sslkeys: %v", be.Name(), err))
	}
	if dnssec, err = be.GetDNSSecKeys(); err != nil {
		return nil, nil, nil, nil, errors.New(fmt.Sprintf("Unable to get %v dnssec keys: %v", be.Name(), err))
	}
	if uri, err = be.GetURISignKeys(); err != nil {
		return nil, nil, nil, nil, errors.New(fmt.Sprintf("Unable to get %v uri keys: %v", be.Name(), err))
	}
	if url, err = be.GetURLSigKeys(); err != nil {
		return nil, nil, nil, nil, errors.New(fmt.Sprintf("Unable to %v url keys: %v", be.Name(), err))
	}
	sort.Slice(sslkeys[:], func(a, b int) bool {
		return sslkeys[a].CDN < sslkeys[b].CDN && sslkeys[a].DeliveryService < sslkeys[b].DeliveryService && sslkeys[a].Version < sslkeys[b].Version
	})
	sort.Slice(dnssec[:], func(a, b int) bool {
		return dnssec[a].CDN < dnssec[b].CDN
	})
	sort.Slice(uri[:], func(a, b int) bool {
		return uri[a].DeliveryService < uri[b].DeliveryService
	})
	sort.Slice(url[:], func(a, b int) bool {
		return url[a].DeliveryService < url[b].DeliveryService
	})
	return sslkeys, dnssec, uri, url, nil
}

// UnmarshalConfig takes in a config file and a type and will read the config file into the reflected type
func UnmarshalConfig(configFile string, t reflect.Type) (reflect.Value, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return reflect.Value{}, err
	}
	val := reflect.New(t)
	err = json.Unmarshal(data, val.Interface())
	if err != nil {
		return reflect.Value{}, err
	}

	return val, nil
}

// TVBackend represents a TV backend that can be have data migrated to/from
type TVBackend interface {
	// Start initiates the connection to the backend DB
	Start() error
	// Stop terminates the connection to the backend DB
	Stop() error
	// Ping checks the connection to the backend DB
	Ping() error
	// ValidateKey validates that the keys are valid (in most cases, certain fields are not null)
	ValidateKey() []string
	// Name returns the name for this backend
	Name() string
	// ReadConfig takes in a filename and will read it into the backends config
	ReadConfig(string) error
	// String returns a high level overview of the backend and its keys
	String() string

	// Fetch gets all of the keys from the backend DB
	Fetch() error
	// Insert takes the current keys and inserts them into the backend DB
	Insert() error

	// GetSSLKeys converts the backends internal key representation into the common representation (SSLKey)
	GetSSLKeys() ([]SSLKey, error)
	// SetSSLKeys takes in keys and converts & encrypts the data into the backends internal format
	SetSSLKeys([]SSLKey) error

	// GetDNSSecKeys converts the backends internal key representation into the common representation (DNSSecKey)
	GetDNSSecKeys() ([]DNSSecKey, error)
	// SetDNSSecKeys takes in keys and converts & encrypts the data into the backends internal format
	SetDNSSecKeys([]DNSSecKey) error

	// GetURISignKeys converts the pg internal key representation into the common representation (URISignKey)
	GetURISignKeys() ([]URISignKey, error)
	// SetURISignKeys takes in keys and converts & encrypts the data into the backends internal format
	SetURISignKeys([]URISignKey) error

	// GetURLSigKeys converts the backends internal key representation into the common representation (URLSigKey)
	GetURLSigKeys() ([]URLSigKey, error)
	// SetURLSigKeys takes in keys and converts & encrypts the data into the backends internal format
	SetURLSigKeys([]URLSigKey) error
}

type CommonRecord struct{}

// SSLKey is the common representation of a SSL Key
type SSLKey struct {
	tc.DeliveryServiceSSLKeys
	CommonRecord
}

// DNSSecKey is the common representation of a DNSSec Key
type DNSSecKey struct {
	CDN string
	tc.DNSSECKeysTrafficVault
	CommonRecord
}

// URISignKey is the common representation of an URI Sign Key
type URISignKey struct {
	DeliveryService string
	Keys            map[string]tc.URISignerKeyset
	CommonRecord
}

// URLSigKey is the common representation of an URL Sig Key
type URLSigKey struct {
	DeliveryService string
	tc.URLSigKeys
	CommonRecord
}

func supportedBackends() []TVBackend {
	return []TVBackend{
		&riakBE, &pgBE,
	}
}
func supportedTypes() []string {
	var typs []string
	for _, be := range supportedBackends() {
		typs = append(typs, be.Name())
	}
	return typs
}
func validateType(typ string) bool {
	for _, be := range supportedBackends() {
		if typ == be.Name() {
			return true
		}
	}
	return false
}
func getBackendFromType(typ string) TVBackend {
	for _, be := range supportedBackends() {
		if be.Name() == typ {
			return be
		}
	}
	return nil
}
