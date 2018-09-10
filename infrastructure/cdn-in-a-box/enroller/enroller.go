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
	"net/url"
	"os"
	"strings"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/client"
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
	types, _, err := s.GetTypeByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(types) == 0 {
		return -1, errors.New("no type with name " + n)
	}
	return types[0].ID, err
}

func (s session) getCDNIDByName(n string) (int, error) {
	cdns, _, err := s.GetCDNByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(cdns) == 0 {
		return -1, errors.New("no CDN with name " + n)
	}
	return cdns[0].ID, err
}

func (s session) getRegionIDByName(n string) (int, error) {
	divisions, _, err := s.GetRegionByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(divisions) == 0 {
		return -1, errors.New("no division with name " + n)
	}
	return divisions[0].ID, err
}

func (s session) getDivisionIDByName(n string) (int, error) {
	divisions, _, err := s.GetDivisionByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(divisions) == 0 {
		return -1, errors.New("no division with name " + n)
	}
	return divisions[0].ID, err
}

func (s session) getPhysLocationIDByName(n string) (int, error) {
	physLocs, _, err := s.GetPhysLocationByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(physLocs) == 0 {
		return -1, errors.New("no physLocation with name " + n)
	}
	return physLocs[0].ID, err
}

func (s session) getCachegroupIDByName(n string) (int, error) {
	cgs, _, err := s.GetCacheGroupByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(cgs) == 0 {
		return -1, errors.New("no cachegroups with name" + n)
	}
	return cgs[0].ID, err
}

func (s session) getProfileIDByName(n string) (int, error) {
	profiles, _, err := s.GetProfileByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(profiles) == 0 {
		return -1, errors.New("no profile with name " + n)
	}
	return profiles[0].ID, err
}

func (s session) getStatusIDByName(n string) (int, error) {
	statuses, _, err := s.GetStatusByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(statuses) == 0 {
		return -1, errors.New("no status with name " + n)
	}
	return statuses[0].ID, err
}

func (s session) getTenantIDByName(n string) (int, error) {
	tenant, _, err := s.TenantByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if tenant == nil {
		return -1, errors.New("no tenant with name " + n)
	}
	return tenant.ID, err
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateType(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.CDN
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateCDN(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateASN(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.CacheGroupNullable
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.Type != nil {
		id, err := toSession.getTypeIDByName(*s.Type)
		if err != nil {
			return err
		}
		s.TypeID = &id
	}

	if s.ParentName != nil {
		id, err := toSession.getCachegroupIDByName(*s.ParentName)
		if err != nil {
			return err
		}
		s.ParentCachegroupID = &id
	}

	if s.SecondaryParentName != nil {
		id, err := toSession.getCachegroupIDByName(*s.SecondaryParentName)
		if err != nil {
			return err
		}
		s.SecondaryParentCachegroupID = &id
	}

	alerts, _, err := toSession.CreateCacheGroupNullable(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.DeliveryService
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.Type != "" {
		id, err := toSession.getTypeIDByName(s.Type.String())
		if err != nil {
			return err
		}
		s.TypeID = id
	}

	alerts, err := toSession.CreateDeliveryService(&s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateDivision(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.Origin
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateOrigin(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateParameter(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.RegionName != "" {
		id, err := toSession.getRegionIDByName(s.RegionName)
		if err != nil {
			return err
		}
		s.RegionID = id
	}
	alerts, _, err := toSession.CreatePhysLocation(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.DivisionName != "" {
		id, err := toSession.getDivisionIDByName(s.DivisionName)
		if err != nil {
			return err
		}
		s.Division = id
	}

	alerts, _, err := toSession.CreateRegion(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	alerts, _, err := toSession.CreateStatus(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.ParentName != "" {
		id, err := toSession.getTenantIDByName(s.ParentName)
		if err != nil {
			return err
		}
		s.ParentID = id
	}

	alerts, err := toSession.CreateTenant(&s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.Profile
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.CDNName != "" {
		id, err := toSession.getCDNIDByName(s.CDNName)
		if err != nil {
			return err
		}
		s.CDNID = id
	}

	alerts, _, err := toSession.CreateProfile(s)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
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
	var s tc.Server
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
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
		log.Printf("error creating from %s: %s\n", fn, err)
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
		const (
			processed = ".processed"
			rejected  = ".rejected"
		)

		for {
			select {
			case event, ok := <-dw.Events:
				if !ok {
					log.Printf("event not ok: %+v", event)
					return
				}

				//log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					if strings.HasSuffix(event.Name, processed) || strings.HasSuffix(event.Name, rejected) {
						continue
					}
					if i, err := os.Stat(event.Name); err != nil || i.IsDir() {
						log.Println("skipping " + event.Name)
						continue
					}
					log.Println("new file :", event.Name)
					p := strings.IndexRune(event.Name, '/')
					if p == -1 {
						continue
					}
					dir := event.Name[:p]
					log.Printf("dir is %s\n", dir)
					suffix := rejected
					if f, ok := dw.watched[dir]; ok {
						log.Printf("creating from %s\n", event.Name)
						time.Sleep(100 * time.Millisecond)

						err := f(toSession, event.Name)
						if err != nil {
							log.Printf("error creating %s from %s: %s\n", dir, event.Name, err.Error())
						} else {
							suffix = processed
						}
					} else {
						log.Printf("no method for creating %s\n", dir)
					}
					// rename the file indicating if processed or rejected
					err = os.Rename(event.Name, event.Name+suffix)
					if err != nil {
						log.Printf("error renaming %s to %s: %s\n", event.Name, event.Name+suffix, err.Error())
					}
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
		log.Fatalln("error starting TrafficOps session: " + err.Error())
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
	dw.watch("phys_locations", enrollPhysLocation)
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
