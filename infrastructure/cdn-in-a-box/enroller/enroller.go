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

	tc "github.com/apache/trafficcontrol/lib/go-tc/v13"
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

// enrollServer takes a json file and creates a Server object using the TO API
func enrollServer(toSession session, fn string) (*tc.Server, error) {
	fh, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.Server
	err = dec.Decode(&s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if s.Type != "" {
		id, err := toSession.getTypeIDByName(s.Type)
		if err != nil {
			return &s, err
		}
		s.TypeID = id
	}

	if s.Profile != "" {
		id, err := toSession.getProfileIDByName(s.Profile)
		if err != nil {
			return &s, err
		}
		s.ProfileID = id
	}
	if s.Status != "" {
		id, err := toSession.getStatusIDByName(s.Status)
		if err != nil {
			return &s, err
		}
		s.StatusID = id
	}
	if s.CDNName != "" {
		id, err := toSession.getCDNIDByName(s.CDNName)
		if err != nil {
			return &s, err
		}
		s.CDNID = id
	}
	if s.Cachegroup != "" {
		id, err := toSession.getCachegroupIDByName(s.Cachegroup)
		if err != nil {
			return &s, err
		}
		s.CachegroupID = id
	}
	if s.PhysLocation != "" {
		id, err := toSession.getPhysLocationIDByName(s.Cachegroup)
		if err != nil {
			return &s, err
		}
		s.PhysLocationID = id
	}

	alerts, _, err := toSession.CreateServer(s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return &s, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: enroller <dir> [<dir> ...]")
	}

	// watch for file creation in directories
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	for _, dir := range os.Args[1:] {
		if err := watcher.Add(dir); err != nil {
			log.Fatalf("error watching directory %s: %v", dir, err)
		}
		log.Println("watching ", dir)
	}

	envconfig.Process("", &to)
	reqTimeout := time.Second * time.Duration(60)

	log.Println("Starting TrafficOps session")
	toSession, err := newSession(reqTimeout, to.URL, to.User, to.Password)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TrafficOps session established")

	// watch for file creation forever
	done := make(chan struct{})
	const startedFile = "enroller-started"
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					if event.Name == startedFile {
						continue
					}
					log.Println("created file:", event.Name)
					s, err := enrollServer(toSession, event.Name)
					if err != nil {
						log.Print(err)
						continue
					}
					log.Printf("Server %s.%s created\n", s.HostName, s.DomainName)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	<-done
}
