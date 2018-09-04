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
	"log"
	"os"
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
	if err != nil || len(types) == 0 {
		fmt.Printf("unknown type %s\n", n)
		return -1, err
	}
	fmt.Printf("type %s: %++v\n", n, types)
	return types[0].ID, err
}

func (s session) getCDNIDByName(n string) (int, error) {
	cdns, _, err := s.GetCDNByName(n)
	if err != nil {
		fmt.Println("cannot get CDNS")
		return -1, err
	}
	if len(cdns) < 1 {
		panic(fmt.Sprintf("CDNS: %v;  err: %v", cdns, err))
	}
	return cdns[0].ID, err
}

func (s session) getCachegroupIDByName(n string) (int, error) {
	cgs, _, err := s.GetCacheGroupByName(n)
	if err != nil {
		fmt.Println("cannot get Cachegroup")
		return -1, err
	}
	if len(cgs) == 0 {
		return -1, errors.New("No cachegroups found")
	}
	return cgs[0].ID, err
}

func (s session) getPhysLocationIDByName(n string) (int, error) {
	physLocs, _, err := s.GetPhysLocationByName(n)
	if err != nil {
		fmt.Println("cannot get physlocations")
		return -1, err
	}
	return physLocs[0].ID, err
}

func (s session) getProfileIDByName(n string) (int, error) {
	profiles, _, err := s.GetProfileByName(n)
	if err != nil {
		fmt.Println("cannot get profiles")
		return -1, err
	}
	return profiles[0].ID, err
}

func (s session) getStatusIDByName(n string) (int, error) {
	statuses, _, err := s.GetStatusByName(n)
	if err != nil {
		fmt.Printf("unknown Status %s\n", n)
		return -1, err
	}
	return statuses[0].ID, err
}

var to struct {
	URL      string `envconfig:"TO_URL"`
	User     string `envconfig:"TO_USER"`
	Password string `envconfig:"TO_PASSWORD"`
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

// watch starts f when a new file is created in dir
func watch(watcher *fsnotify.Watcher, toSession *session, dir string, f func(*session, string) error) {
	go func() {
		log.Println("started watching " + dir)
		defer func() { log.Println("stopped watching " + dir) }()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Printf("event not ok: %+v", event)
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("new "+dir+" file:", event.Name)
					err := f(toSession, event.Name)
					if err != nil {
						log.Print(err)
						continue
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Printf("error not ok: %+v", err)
				}
				log.Println("error:", err)
			}
		}
	}()
}

const startedFile = "enroller-started"

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: enroller <dir> [<dir> ...]")
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	watch(watcher, &toSession, "cdns", enrollCDN)
	watch(watcher, &toSession, "profiles", enrollProfile)
	watch(watcher, &toSession, "parameters", enrollParameter)
	watch(watcher, &toSession, "servers", enrollServer)

	var waitforever chan struct{}
	<-waitforever
}
