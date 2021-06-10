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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/pborman/getopt/v2"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

var (
	fromType  string
	toType    string
	fromCfg   string
	toCfg     string
	dry       bool
	compare   bool
	noConfirm bool
	dump      bool

	riakBE RiakBackend = RiakBackend{}
	pgBE   PGBackend   = PGBackend{}
)

func init() {
	fromTypePtr := getopt.StringLong("fromType", 't', riakBE.Name(), fmt.Sprintf("From server types (%v)", strings.Join(supportedTypes(), "|")))
	if fromTypePtr == nil {
		log.Fatal("unable to load fromType")
	}
	fromType = *fromTypePtr

	toTypePtr := getopt.StringLong("toType", 'o', pgBE.Name(), fmt.Sprintf("From server types (%v)", strings.Join(supportedTypes(), "|")))
	if toTypePtr == nil {
		log.Fatal("unable to load toType")
	}
	toType = *toTypePtr

	toCfgPtr := getopt.StringLong("toCfg", 'g', "pg.json", "To server config file")
	if toCfgPtr == nil {
		log.Fatal("unable to load toCfg")
	}
	toCfg = *toCfgPtr

	fromCfgPtr := getopt.StringLong("fromCfg", 'f', "riak.json", "From server config file")
	if fromCfgPtr == nil {
		log.Fatal("unable to load fromCfg")
	}
	fromCfg = *fromCfgPtr

	dryPtr := getopt.BoolLong("dry", 'r', "Do not perform writes")
	if dryPtr == nil {
		log.Fatal("unable to load dry")
	}
	dry = *dryPtr

	comparePtr := getopt.BoolLong("compare", 'c', "Compare to and from server records")
	if comparePtr == nil {
		log.Fatal("unable to load compare")
	}
	compare = *comparePtr

	noConfirmPtr := getopt.BoolLong("noConfirm", 'm', "Requires confirmation before inserting records")
	if noConfirmPtr == nil {
		log.Fatal("unable to load noConfirm")
	}
	noConfirm = *noConfirmPtr

	dumpPtr := getopt.BoolLong("dump", 'd', "Write keys (from 'from' server) to disk")
	if dumpPtr == nil {
		log.Fatal("unable to load dump")
	}
	dump = *dumpPtr
}

// supportBackends returns the backends available in this tool
func supportedBackends() []TVBackend {
	return []TVBackend{
		&riakBE, &pgBE,
	}
}

func main() {
	getopt.ParseV2()
	var fromSrv TVBackend
	var toSrv TVBackend

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

	log.Println("Reading configs...")
	if err := fromSrv.ReadConfig(fromCfg); err != nil {
		log.Fatalf("Unable to read fromSrv cfg: %v", err)
	}

	if toSrvUsed {
		if err := toSrv.ReadConfig(toCfg); err != nil {
			log.Fatalf("Unable to read toSrv cfg: %v", err)
		}
	}

	log.Println("Starting server connections...")
	if err := fromSrv.Start(); err != nil {
		log.Fatalf("issue starting fromSrv: %v", err)
	}
	defer func() {
		fromSrv.Close()
	}()
	if toSrvUsed {
		if err := toSrv.Start(); err != nil {
			log.Fatalf("issue starting toSrv: %v", err)
		}
		defer func() {
			toSrv.Close()
		}()
	}

	log.Println("Pinging servers...")
	if err := fromSrv.Ping(); err != nil {
		log.Fatalf("Unable to ping fromSrv: %v", err)
	}
	if toSrvUsed {
		if err := toSrv.Ping(); err != nil {
			log.Fatalf("Unable to ping toSrv: %v", err)
		}
	}

	log.Printf("Fetching data from %v...\n", fromSrv.Name())
	if err := fromSrv.Fetch(); err != nil {
		log.Fatalf("Unable to fetch fromSrv data: %v", err)
	}

	fromSecret, err := GetKeys(fromSrv)
	if err != nil {
		log.Fatal(err)
	}

	if err := Validate(fromSrv); err != nil {
		log.Fatal(err)
	}

	if dump {
		log.Printf("Dumping data from %v...\n", fromSrv.Name())
		fromSecret.dump("dump")
		return
	}

	if compare {
		log.Printf("Fetching data from %v...\n", toSrv.Name())
		if err := toSrv.Fetch(); err != nil {
			log.Fatalf("Unable to fetch toSrv data: %v\n", err)
		}

		toSecret, err := GetKeys(toSrv)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Validating " + toSrv.Name())
		if err := toSrv.ValidateKey(); err != nil && len(err) > 0 {
			log.Fatal(strings.Join(err, "\n"))
		}

		fromSecret.sort()
		toSecret.sort()

		log.Println(fromSrv.String())
		log.Println(toSrv.String())

		if !reflect.DeepEqual(fromSecret.sslkeys, toSecret.sslkeys) {
			log.Fatal("from sslkeys and to sslkeys don't match")
		}
		if !reflect.DeepEqual(fromSecret.dnssecKeys, toSecret.dnssecKeys) {
			log.Fatal("from dnssec and to dnssec don't match")
		}
		if !reflect.DeepEqual(fromSecret.uriKeys, toSecret.uriKeys) {
			log.Fatal("from uri and to uri don't match")
		}
		if !reflect.DeepEqual(fromSecret.urlKeys, toSecret.urlKeys) {
			log.Fatal("from url and to url don't match")
		}
		log.Println("Both data sources have the same keys")
		return
	}

	log.Printf("Setting %v keys...\n", toSrv.Name())
	if err := SetKeys(toSrv, fromSecret); err != nil {
		log.Fatal(err)
	}

	if err := Validate(toSrv); err != nil {
		log.Fatal(err)
	}

	log.Println(fromSrv.String())

	if dry {
		return
	}

	if !noConfirm {
		ans := "q"
		for {
			fmt.Print("Confirm data insertion (y/n): ")
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
	if err := toSrv.Insert(); err != nil {
		log.Fatal(err)
	}
}

// Validate runs the ValidateKey method on the backend
func Validate(be TVBackend) error {
	if errs := be.ValidateKey(); errs != nil && len(errs) > 0 {
		return errors.New(fmt.Sprintf("Validation Errors (%v): \n%v", be.Name(), strings.Join(errs, "\n")))
	}
	return nil
}

// SetKeys will set all of the keys for a backend
func SetKeys(be TVBackend, s Secrets) error {
	if err := be.SetSSLKeys(s.sslkeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v ssl keys: %v", be.Name(), err))
	}
	if err := be.SetDNSSecKeys(s.dnssecKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v dnssec keys: %v", be.Name(), err))
	}
	if err := be.SetURLSigKeys(s.urlKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v url keys: %v", be.Name(), err))
	}
	if err := be.SetURISignKeys(s.uriKeys); err != nil {
		return errors.New(fmt.Sprintf("Unable to set %v uri keys: %v", be.Name(), err))
	}
	return nil
}

// GetKeys will get all of the keys for a backend
func GetKeys(be TVBackend) (Secrets, error) {
	var secret Secrets
	var err error
	if secret.sslkeys, err = be.GetSSLKeys(); err != nil {
		return Secrets{}, errors.New(fmt.Sprintf("Unable to get %v sslkeys: %v", be.Name(), err))
	}
	if secret.dnssecKeys, err = be.GetDNSSecKeys(); err != nil {
		return Secrets{}, errors.New(fmt.Sprintf("Unable to get %v dnssec keys: %v", be.Name(), err))
	}
	if secret.uriKeys, err = be.GetURISignKeys(); err != nil {
		return Secrets{}, errors.New(fmt.Sprintf("Unable to get %v uri keys: %v", be.Name(), err))
	}
	if secret.urlKeys, err = be.GetURLSigKeys(); err != nil {
		return Secrets{}, errors.New(fmt.Sprintf("Unable to %v url keys: %v", be.Name(), err))
	}
	return secret, nil
}

// UnmarshalConfig takes in a config file and a type and will read the config file into the reflected type
func UnmarshalConfig(configFile string, config interface{}) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

// TVBackend represents a TV backend that can be have data migrated to/from
type TVBackend interface {
	// Start initiates the connection to the backend DB
	Start() error
	// Close terminates the connection to the backend DB
	Close() error
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

// Secrets contains every key to be migrated
type Secrets struct {
	sslkeys    []SSLKey
	dnssecKeys []DNSSecKey
	uriKeys    []URISignKey
	urlKeys    []URLSigKey
}

func (s *Secrets) sort() {
	sort.Slice(s.sslkeys[:], func(a, b int) bool {
		return s.sslkeys[a].CDN < s.sslkeys[b].CDN ||
			s.sslkeys[a].CDN == s.sslkeys[b].CDN && s.sslkeys[a].DeliveryService < s.sslkeys[b].DeliveryService
	})
	sort.Slice(s.dnssecKeys[:], func(a, b int) bool {
		return s.dnssecKeys[a].CDN < s.dnssecKeys[b].CDN
	})
	sort.Slice(s.uriKeys[:], func(a, b int) bool {
		return s.uriKeys[a].DeliveryService < s.uriKeys[b].DeliveryService
	})
	sort.Slice(s.urlKeys[:], func(a, b int) bool {
		return s.urlKeys[a].DeliveryService < s.urlKeys[b].DeliveryService
	})
}
func (s *Secrets) dump(directory string) {
	if err := os.Mkdir(directory, 0750); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	if err := writeKeys(directory+"/sslkeys.json", &s.sslkeys); err != nil {
		log.Fatal(err)
	}
	if err := writeKeys(directory+"/dnsseckeys.json", &s.dnssecKeys); err != nil {
		log.Fatal(err)
	}
	if err := writeKeys(directory+"/urlkeys.json", &s.urlKeys); err != nil {
		log.Fatal(err)
	}
	if err := writeKeys(directory+"/urikeys.json", &s.uriKeys); err != nil {
		log.Fatal(err)
	}
}

// SSLKey is the common representation of a SSL Key
type SSLKey struct {
	tc.DeliveryServiceSSLKeys
}

// DNSSecKey is the common representation of a DNSSec Key
type DNSSecKey struct {
	CDN string
	tc.DNSSECKeysTrafficVault
}

// URISignKey is the common representation of an URI Sign Key
type URISignKey struct {
	DeliveryService string
	Keys            map[string]tc.URISignerKeyset
}

// URLSigKey is the common representation of an URL Sig Key
type URLSigKey struct {
	DeliveryService string
	tc.URLSigKeys
}

func writeKeys(filename string, data interface{}) error {
	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, bytes, 0640); err != nil {
		return err
	}

	return nil
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
