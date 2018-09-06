package main

// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	v13 "github.com/apache/trafficcontrol/lib/go-tc/v13"
	client "github.com/apache/trafficcontrol/traffic_ops/client/v13"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/fsnotify.v1"
)

type session struct {
	*client.Session
}

func newSession(reqTimeout time.Duration, toURL string, toUser string, toPass string) (session, error) {
	s, _, err := client.LoginWithAgent(toURL, toUser, toPass, true, "cdn-in-a-box-enroller", true, reqTimeout)
	return session{s}, err
}

func printJSON(label string, b interface{}) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent(``, `  `)
	enc.Encode(b)
	fmt.Println(label, buf.String())
}

func (s session) getTypeIDByName(n string) (int, error) {
	types, _, err := s.GetTypeByName(n)
	if err != nil {
		return -1, err
	}
	if len(types) == 0 {
		return -1, errors.New("no type with name " + n)
	}
	return types[0].ID, err
}

func (s session) getCDNIDByName(n string) (int, error) {
	cdns, _, err := s.GetCDNByName(n)
	if err != nil {
		return -1, err
	}
	if len(cdns) == 0 {
		return -1, errors.New("no CDN with name " + n)
	}
	return cdns[0].ID, err
}

func (s session) getCachegroupIDByName(n string) (int, error) {
	cgs, _, err := s.GetCacheGroupByName(n)
	if err != nil {
		return -1, err
	}
	if len(cgs) == 0 {
		return -1, errors.New("no cachegroups with name" + n)
	}
	return cgs[0].ID, err
}

func (s session) getPhysLocationIDByName(n string) (int, error) {
	physLocs, _, err := s.GetPhysLocationByName(n)
	if err != nil {
		return -1, err
	}
	if len(physLocs) == 0 {
		return -1, errors.New("no physLocation with name " + n)
	}
	return physLocs[0].ID, err
}

func (s session) getProfileIDByName(n string) (int, error) {
	profiles, _, err := s.GetProfileByName(n)
	if err != nil {
		return -1, err
	}
	if len(profiles) == 0 {
		return -1, errors.New("no profile with name " + n)
	}
	return profiles[0].ID, err
}

func (s session) getStatusIDByName(n string) (int, error) {
	statuses, _, err := s.GetStatusByName(n)
	if err != nil {
		return -1, err
	}
	if len(statuses) == 0 {
		return -1, errors.New("no status with name " + n)
	}
	return statuses[0].ID, err
}

var to struct {
	URL      string `envconfig:"TO_URL"`
	User     string `envconfig:"TO_USER"`
	Password string `envconfig:"TO_PASSWORD"`
}

// enrollType takes a json file and creates a Type object using the TO API
func enrollType(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Type
	err = dec.Decode(&s)
	if err != io.EOF {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateType(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollCDN takes a json file and creates a CDN object using the TO API
func enrollCDN(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s v13.CDN
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateCDN(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollASN(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.ASN
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateASN(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollCachegroup takes a json file and creates a Cachegroup object using the TO API
func enrollCachegroup(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s v13.CacheGroup
	err = dec.Decode(&s)
	if err != io.EOF {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateCacheGroup(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}
func enrollDeliveryService(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.DeliveryServiceV13
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, err := toSession.CreateDeliveryService(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollDivision(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Division
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateDivision(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollOrigin(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s v13.Origin
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateOrigin(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollParameter(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Parameter
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateParameter(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollPhysLocation(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.PhysLocation
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreatePhysLocation(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollRegion(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Region
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateRegion(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollStatus(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Status
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateStatus(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollTenant(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Tenant
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, err := toSession.CreateTenant(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollProfile takes a json file and creates a Profile object using the TO API
func enrollProfile(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s v13.Profile
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	alerts, _, err := toSession.CreateProfile(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollServer takes a json file and creates a Server object using the TO API
func enrollServer(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s v13.Server
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return err
	}

	if s.Type != "" {
		id, err := toSession.getTypeIDByName(s.Type)
		if err != nil {
			return err
		}
		s.TypeID = id
	}

	if s.Profile != "" {
		id, err := toSession.getProfileIDByName(s.Profile)
		if err != nil {
			return err
		}
		s.ProfileID = id
	}
	if s.Status != "" {
		id, err := toSession.getStatusIDByName(s.Status)
		if err != nil {
			return err
		}
		s.StatusID = id
	}
	if s.CDNName != "" {
		id, err := toSession.getCDNIDByName(s.CDNName)
		if err != nil {
			return err
		}
		s.CDNID = id
	}
	if s.Cachegroup != "" {
		id, err := toSession.getCachegroupIDByName(s.Cachegroup)
		if err != nil {
			return err
		}
		s.CachegroupID = id
	}
	if s.PhysLocation != "" {
		id, err := toSession.getPhysLocationIDByName(s.Cachegroup)
		if err != nil {
			return err
		}
		s.PhysLocationID = id
	}

	alerts, _, err := toSession.CreateServer(s)
	if err != nil {
		log.Println(err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

type dirWatcher struct {
	*fsnotify.Watcher
	TOSession *session
	watched   map[string]func(toSession *session, fn string) error
}

func newDirWatcher(toSession *session) (*dirWatcher, error) {
	var err error
	var dw dirWatcher
	dw.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dw.watched = make(map[string]func(toSession *session, fn string) error)
	go func() {
		for {
			select {
			case event, ok := <-dw.Events:
				if !ok {
					log.Printf("event not ok: %+v", event)
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("new file :", event.Name)
					p := strings.IndexRune(event.Name, '/')
					if p == -1 {
						continue
					}
					dir := event.Name[:p]
					log.Printf("dir is %s\n", dir)
					if f, ok := dw.watched[dir]; ok {
						err := f(toSession, event.Name)
						if err != nil {
							log.Printf("error creating %s from %s: %+v\n", dir, event.Name, err)
						}
					} else {
						log.Printf("no method for creating %s\n", dir)
					}

					//err := f(toSession, event.Name)
					//if err != nil {
					//		log.Print(err)
					//			continue
					//			}
				}
			case err, ok := <-dw.Errors:
				if !ok {
					log.Printf("error not ok: %+v", err)
				}
				log.Println("error:", err)
			}
		}
	}()
	return &dw, err
}

// watch starts f when a new file is created in dir
func (dw *dirWatcher) watch(dir string, f func(*session, string) error) {
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		// attempt to create dir
		if err = os.Mkdir(dir, os.ModeDir); err != nil {
			log.Println("cannot watch " + dir + ": not a directory")
			return
		}
	}
	log.Println("watching " + dir)
	dw.Add(dir)
	dw.watched[dir] = f
}

const startedFile = "enroller-started"

func main() {
	watchDir := "."
	if len(os.Args) > 1 {
		watchDir = os.Args[1]
	}
	if stat, err := os.Stat(watchDir); err != nil || !stat.IsDir() {
		log.Fatalln("expected " + watchDir + " to be a directory")
	}
	if err := os.Chdir(watchDir); err != nil {
		log.Fatalf("cannot chdir to %s: %s", watchDir, err)
	}
	envconfig.Process("", &to)

	reqTimeout := time.Second * time.Duration(60)

	log.Println("Starting TrafficOps session")
	toSession, err := newSession(reqTimeout, to.URL, to.User, to.Password)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TrafficOps session established")

	// watch for file creation in directories
	dw, err := newDirWatcher(&toSession)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer dw.Close()

	dw.watch("types", enrollType)
	dw.watch("cdns", enrollCDN)
	dw.watch("cachegroups", enrollCachegroup)
	dw.watch("profiles", enrollProfile)
	dw.watch("parameters", enrollParameter)
	dw.watch("servers", enrollServer)
	dw.watch("asns", enrollASN)
	dw.watch("deliveryservices", enrollDeliveryService)
	dw.watch("divisions", enrollDivision)
	dw.watch("origins", enrollOrigin)
	dw.watch("physlocations", enrollPhysLocation)
	dw.watch("regions", enrollRegion)
	dw.watch("statuses", enrollStatus)
	dw.watch("tenants", enrollTenant)

	// create this file to indicate the enroller is ready
	f, err := os.Create(startedFile)
	if err != nil {
		panic(err)
	}
	f.Close()

	var waitforever chan struct{}
	<-waitforever
}
