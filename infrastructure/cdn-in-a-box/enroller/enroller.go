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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/v4-client"
	"github.com/fsnotify/fsnotify"
	"github.com/kelseyhightower/envconfig"
)

var startedFile = "enroller-started"

type session struct {
	*client.Session
}

func newSession(reqTimeout time.Duration, toURL string, toUser string, toPass string) (session, error) {
	s, _, err := client.LoginWithAgent(toURL, toUser, toPass, true, "cdn-in-a-box-enroller", true, reqTimeout)
	return session{s}, err
}

func (s session) getParameter(m tc.Parameter, header http.Header) (tc.Parameter, error) {
	// TODO: s.GetParameterByxxx() does not seem to work with values with spaces --
	// doing this the hard way for now
	parameters, _, err := s.GetParametersWithHdr(header)
	if err != nil {
		return m, err
	}
	for _, p := range parameters {
		if p.Name == m.Name && p.Value == m.Value && p.ConfigFile == m.ConfigFile {
			return p, nil
		}
	}
	return m, fmt.Errorf("no parameter matching name %s, configFile %s, value %s", m.Name, m.ConfigFile, m.Value)
}

// enrollType takes a json file and creates a Type object using the TO API
func enrollType(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Type
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Type: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateType(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("type %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Type: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollCDN takes a json file and creates a CDN object using the TO API
func enrollCDN(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.CDN
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding CDN: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateCDN(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("cdn %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating CDN: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollASN(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.ASN
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding ASN: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateASN(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("asn %d already exists\n", s.ASN)
			return nil
		}
		log.Infof("error creating ASN: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollCachegroup takes a json file and creates a Cachegroup object using the TO API
func enrollCachegroup(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.CacheGroupNullable
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Cachegroup: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateCacheGroupNullable(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("cachegroup %s already exists\n", *s.Name)
			return nil
		}
		log.Infof("error creating Cachegroup: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollTopology(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Topology
	err := dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Infof("error decoding Topology: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateTopology(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("topology %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Topology: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollDeliveryService(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.DeliveryServiceNullableV30
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding DeliveryService: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateDeliveryServiceV30(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("deliveryservice %s already exists\n", *s.XMLID)
			return nil
		}
		log.Infof("error creating DeliveryService: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollDeliveryServicesRequiredCapability takes a json file and creates a DeliveryServicesRequiredCapability object using the TO API
func enrollDeliveryServicesRequiredCapability(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var dsrc tc.DeliveryServicesRequiredCapability
	err := dec.Decode(&dsrc)
	if err != nil {
		log.Infof("error decoding Delivery Services Required Capability: %s\n", err)
		return err
	}

	dses, _, err := toSession.GetDeliveryServiceByXMLIDNullableWithHdr(*dsrc.XMLID, nil)
	if err != nil {
		log.Infof("getting Delivery Service by XMLID %s: %s", *dsrc.XMLID, err.Error())
		return err
	}
	if len(dses) < 1 {
		err = errors.New("could not find a Delivey Service with XMLID %s")
		log.Infoln(err)
		return err
	}
	dsrc.DeliveryServiceID = dses[0].ID

	alerts, _, err := toSession.CreateDeliveryServicesRequiredCapability(dsrc)
	if err != nil {
		log.Infof("error creating Delivery Services Required Capability: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollDeliveryServiceServer(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)

	// DeliveryServiceServers lists ds xmlid and array of server names.  Use that to create multiple DeliveryServiceServer objects
	var dss tc.DeliveryServiceServers
	err := dec.Decode(&dss)
	if err != nil {
		log.Infof("error decoding DeliveryServiceServer: %s\n", err)
		return err
	}

	params := url.Values{"xmlId": []string{dss.XmlId}}
	dses, _, err := toSession.GetDeliveryServicesV4(nil, params)
	if err != nil {
		return err
	}
	if len(dses) == 0 {
		return errors.New("no deliveryservice with name " + dss.XmlId)
	}
	if dses[0].ID == nil {
		return errors.New("Deliveryservice with name " + dss.XmlId + " has a nil ID")
	}
	dsID := *dses[0].ID

	params = url.Values{}
	var serverIDs []int
	for _, sn := range dss.ServerNames {
		params.Set("hostName", sn)
		servers, _, err := toSession.GetServersWithHdr(&params, nil)
		if err != nil {
			return err
		}
		if len(servers.Response) == 0 {
			return errors.New("no server with hostName " + sn)
		}
		serverIDs = append(serverIDs, *servers.Response[0].ID)
	}
	_, _, err = toSession.CreateDeliveryServiceServers(dsID, serverIDs, true)
	if err != nil {
		log.Infof("error creating DeliveryServiceServer: %s\n", err)
	}

	return err
}

func enrollDivision(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Division
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Division: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateDivision(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("division %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Division: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollOrigin(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Origin
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Origin: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateOrigin(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("origin %s already exists\n", *s.Name)
			return nil
		}
		log.Infof("error creating Origin: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollParameter(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var params []tc.Parameter
	err := dec.Decode(&params)
	if err != nil {
		log.Infof("error decoding Parameter: %s\n", err)
		return err
	}

	for _, p := range params {
		eparam, err := toSession.getParameter(p, nil)
		var alerts tc.Alerts
		if err == nil {
			// existing param -- update
			alerts, _, err = toSession.UpdateParameterByID(eparam.ID, p)
			if err != nil {
				log.Infof("error updating parameter %d: %s with %+v ", eparam.ID, err.Error(), p)
				break
			}
		} else {
			alerts, _, err = toSession.CreateParameter(p)
			if err != nil {
				log.Infof("error creating parameter: %s from %+v\n", err.Error(), p)
				return err
			}
			eparam, err = toSession.getParameter(p, nil)
			if err != nil {
				return err
			}
		}

		// link parameter with profiles
		if len(p.Profiles) > 0 {
			var profiles []string
			err = json.Unmarshal(p.Profiles, &profiles)
			if err != nil {
				log.Infof("%v", err)
				return err
			}

			for _, n := range profiles {
				profiles, _, err := toSession.GetProfileByNameWithHdr(n, nil)
				if err != nil {
					return err
				}
				if len(profiles) == 0 {
					return errors.New("no profile with name " + n)
				}

				pp := tc.ProfileParameter{ParameterID: eparam.ID, ProfileID: profiles[0].ID}
				_, _, err = toSession.CreateProfileParameter(pp)
				if err != nil {
					if strings.Contains(err.Error(), "already exists") {
						continue
					}
				}
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(&alerts)
	}
	return err
}

func enrollPhysLocation(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.PhysLocation
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding PhysLocation: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreatePhysLocation(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("physLocation %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating PhysLocation: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollRegion(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Region
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Region: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateRegion(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("region %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Region: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollStatus(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.StatusNullable
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Status: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateStatusNullable(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("status %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Status: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollTenant(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.Tenant
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Tenant: %s\n", err)
		return err
	}

	alerts, err := toSession.CreateTenant(&s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("tenant %s already exists\n", s.Name)
			return nil
		}
		log.Infof("error creating Tenant: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollUser(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.User
	err := dec.Decode(&s)
	log.Infof("User is %++v\n", s)
	if err != nil {
		log.Infof("error decoding User: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateUser(&s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Infof("user %s already exists\n", *s.Username)
			return nil
		}
		log.Infof("error creating User: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollProfile takes a json file and creates a Profile object using the TO API
func enrollProfile(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var profile tc.Profile

	err := dec.Decode(&profile)
	if err != nil {
		log.Infof("error decoding Profile: %s\n", err)
		return err
	}
	// get a copy of the parameters
	parameters := profile.Parameters

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("  ", "")
	enc.Encode(profile)

	if len(profile.Name) == 0 {
		log.Infoln("missing name on profile")
		return errors.New("missing name on profile")
	}

	profiles, _, err := toSession.GetProfileByNameWithHdr(profile.Name, nil)

	createProfile := false
	if err != nil || len(profiles) == 0 {
		// no profile by that name -- need to create it
		createProfile = true
	} else {
		// updating - ID needs to match
		profile = profiles[0]
	}

	var alerts tc.Alerts
	var action string
	if createProfile {
		alerts, _, err = toSession.CreateProfile(profile)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Infof("profile %s already exists\n", profile.Name)
			} else {
				log.Infof("error creating profile from %+v: %s\n", profile, err.Error())
			}
		}
		profiles, _, err = toSession.GetProfileByNameWithHdr(profile.Name, nil)
		if err != nil {
			log.Infof("error getting profile ID from %+v: %s\n", profile, err.Error())
		}
		if len(profiles) == 0 {
			log.Infof("no results returned for getting profile ID from %+v", profile)
		}
		profile = profiles[0]
		action = "creating"
	} else {
		alerts, _, err = toSession.UpdateProfileByID(profile.ID, profile)
		action = "updating"
	}

	if err != nil {
		log.Infof("error "+action+" from %s: %s\n", err)
		return err
	}

	for _, p := range parameters {
		var name, configFile, value string
		var secure bool
		if p.ConfigFile != nil {
			configFile = *p.ConfigFile
		}
		if p.Name != nil {
			name = *p.Name
		}
		if p.Value != nil {
			value = *p.Value
		}
		param := tc.Parameter{ConfigFile: configFile, Name: name, Value: value, Secure: secure}
		eparam, err := toSession.getParameter(param, nil)
		if err != nil {
			// create it
			log.Infof("creating param %+v\n", param)
			_, _, err = toSession.CreateParameter(param)
			if err != nil {
				log.Infof("can't create parameter %+v: %s\n", param, err.Error())
				continue
			}
			eparam, err = toSession.getParameter(param, nil)
			if err != nil {
				log.Infof("error getting new parameter %+v: \n", param)
				log.Infof(err.Error())
				continue
			}
		} else {
			log.Infof("found param %+v\n", eparam)
		}

		if eparam.ID < 1 {
			log.Infof("param ID not found for %v", eparam)
			continue
		}
		pp := tc.ProfileParameter{ProfileID: profile.ID, ParameterID: eparam.ID}
		_, _, err = toSession.CreateProfileParameter(pp)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.Infof("error creating profileparameter %+v: %s\n", pp, err.Error())
			}
			continue
		}
	}

	//enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollServer takes a json file and creates a Server object using the TO API
func enrollServer(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.ServerV40
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Server: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateServer(s, nil)
	if err != nil {
		log.Infof("error creating Server: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollServerCapability takes a json file and creates a ServerCapability object using the TO API
func enrollServerCapability(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.ServerCapability
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Server Capability: %s\n", err)
		return err
	}

	alerts, _, err := toSession.CreateServerCapability(s)
	if err != nil {
		log.Infof("error creating Server Capability: %s\n", err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// enrollServerServerCapability takes a json file and creates a ServerServerCapability object using the TO API
func enrollServerServerCapability(toSession *session, r io.Reader) error {
	dec := json.NewDecoder(r)
	var s tc.ServerServerCapability
	err := dec.Decode(&s)
	if err != nil {
		log.Infof("error decoding Server: %s\n", err)
		return err
	}

	resp, _, err := toSession.GetServersWithHdr(&url.Values{"hostName": []string{*s.Server}}, nil)
	if err != nil {
		log.Infof("getting server %s: %s\n", *s.Server, err.Error())
		return err
	}
	if len(resp.Response) < 1 {
		err = fmt.Errorf("could not find Server %s", *s.Server)
		log.Infoln(err.Error())
		return err
	}
	if len(resp.Response) > 1 {
		err = fmt.Errorf("found more than 1 Server with hostname %s", *s.Server)
		log.Infoln(err.Error())
		return err
	}
	s.ServerID = resp.Response[0].ID

	alerts, _, err := toSession.CreateServerServerCapability(s)
	if err != nil {
		log.Infof("error creating Server Server Capability: %s\n", err)
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
			retry     = ".retry"
		)
		originalNameRegex := regexp.MustCompile(`(\.retry)*$`)

		emptyCount := map[string]int{}
		const maxEmptyTries = 10

		for {
			select {
			case event, ok := <-dw.Events:
				if !ok {
					log.Infoln("event not ok")
					continue
				}

				// ignore all but Create events
				if event.Op&fsnotify.Create != fsnotify.Create {
					continue
				}

				// skip already processed files
				if strings.HasSuffix(event.Name, processed) || strings.HasSuffix(event.Name, rejected) {
					continue
				}

				i, err := os.Stat(event.Name)
				if err != nil || i.IsDir() {
					log.Infoln("skipping " + event.Name)
					continue
				}
				log.Infoln("new file :", event.Name)

				// what directory is the file in?  Invoke the matching func
				dir := filepath.Base(filepath.Dir(event.Name))
				suffix := rejected
				if f, ok := dw.watched[dir]; ok {
					t := filepath.Base(dir)
					log.Infoln("creating " + t + " from " + event.Name)
					// Sleep for 100 milliseconds so that the file content is probably there when the directory watcher
					// sees the file
					time.Sleep(100 * time.Millisecond)

					err := f(toSession, event.Name)
					// If a file is empty, try reading from it 10 times before giving up on that file
					if err == io.EOF {
						originalName := originalNameRegex.ReplaceAllString(event.Name, "")
						if _, exists := emptyCount[originalName]; !exists {
							emptyCount[originalName] = 0
						}
						emptyCount[originalName]++
						log.Infof("empty json object %s: %s\ntried file %d out of %d times", originalName, err.Error(), emptyCount[originalName], maxEmptyTries)
						if emptyCount[originalName] < maxEmptyTries {
							newName := event.Name + retry
							if err := os.Rename(event.Name, newName); err != nil {
								log.Infof("error renaming %s to %s: %s", event.Name, newName, err)
							}
							continue
						}
					}
					if err != nil {
						log.Infof("error creating %s from %s: %s\n", dir, event.Name, err.Error())
					} else {
						suffix = processed
					}
				} else {
					log.Infof("no method for creating %s\n", dir)
				}
				// rename the file indicating if processed or rejected
				err = os.Rename(event.Name, event.Name+suffix)
				if err != nil {
					log.Infof("error renaming %s to %s: %s\n", event.Name, event.Name+suffix, err.Error())
				}
			case err, ok := <-dw.Errors:
				log.Infof("error from fsnotify: ok? %v;  error: %v\n", ok, err)
				continue
			}
		}
	}()
	return &dw, err
}

// watch starts f when a new file is created in dir
func (dw *dirWatcher) watch(watchdir, t string, f func(*session, io.Reader) error) {
	dir := watchdir + "/" + t
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		// attempt to create dir
		if err = os.Mkdir(dir, os.ModeDir|0700); err != nil {
			log.Infoln("cannot watch " + dir + ": not a directory")
			return
		}
	}

	log.Infoln("watching " + dir)
	dw.Add(dir)
	dw.watched[t] = func(toSession *session, fn string) error {
		fh, err := os.Open(fn)
		if err != nil {
			return err
		}
		defer log.Close(fh, "could not close file")
		return f(toSession, fh)
	}
}

func startWatching(watchDir string, toSession *session, dispatcher map[string]func(*session, io.Reader) error) (*dirWatcher, error) {
	// watch for file creation in directories
	dw, err := newDirWatcher(toSession)
	if err == nil {
		for d, f := range dispatcher {
			dw.watch(watchDir, d, f)
		}
	}
	return dw, err
}

func startServer(httpPort string, toSession *session, dispatcher map[string]func(*session, io.Reader) error) error {
	baseEP := "/api/2.0/"
	for d, f := range dispatcher {
		http.HandleFunc(baseEP+d, func(w http.ResponseWriter, r *http.Request) {
			defer log.Close(r.Body, "could not close reader")
			f(toSession, r.Body)
		})
	}

	go func() {
		server := &http.Server{
			Addr:      httpPort,
			TLSConfig: nil,
			ErrorLog:  log.Error,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Errorf("stopping server: %v\n", err)
			panic(err)
		}
	}()

	log.Infoln("http service started on " + httpPort)
	return nil
}

// Set up the log config -- all messages go to stdout
type logConfig struct{}

func (cfg logConfig) ErrorLog() log.LogLocation {
	return log.LogLocationStdout
}
func (cfg logConfig) WarningLog() log.LogLocation {
	return log.LogLocationStdout
}
func (cfg logConfig) InfoLog() log.LogLocation {
	return log.LogLocationStdout
}
func (cfg logConfig) DebugLog() log.LogLocation {
	return log.LogLocationStdout
}
func (cfg logConfig) EventLog() log.LogLocation {
	return log.LogLocationStdout
}

func main() {
	var watchDir, httpPort string

	flag.StringVar(&startedFile, "started", startedFile, "file indicating service was started")
	flag.StringVar(&watchDir, "dir", "", "base directory to watch")
	flag.StringVar(&httpPort, "http", "", "act as http server for POST on this port (e.g. :7070)")
	flag.Parse()

	err := log.InitCfg(logConfig{})
	if err != nil {
		panic(err.Error())
	}
	if watchDir == "" && httpPort == "" {
		// if neither -dir nor -http provided, default to watching the current dir
		watchDir = "."
	}

	var toCreds struct {
		URL      string `envconfig:"TO_URL"`
		User     string `envconfig:"TO_USER"`
		Password string `envconfig:"TO_PASSWORD"`
	}

	envconfig.Process("", &toCreds)

	reqTimeout := time.Second * time.Duration(60)

	log.Infoln("Starting TrafficOps session")
	toSession, err := newSession(reqTimeout, toCreds.URL, toCreds.User, toCreds.Password)
	if err != nil {
		log.Errorln("error starting TrafficOps session: " + err.Error())
		os.Exit(1)
	}
	log.Infoln("TrafficOps session established")

	// dispatcher maps an API endpoint name to a function to act on the JSON input Reader
	dispatcher := map[string]func(*session, io.Reader) error{
		"types":                                  enrollType,
		"cdns":                                   enrollCDN,
		"cachegroups":                            enrollCachegroup,
		"topologies":                             enrollTopology,
		"profiles":                               enrollProfile,
		"parameters":                             enrollParameter,
		"servers":                                enrollServer,
		"server_capabilities":                    enrollServerCapability,
		"server_server_capabilities":             enrollServerServerCapability,
		"asns":                                   enrollASN,
		"deliveryservices":                       enrollDeliveryService,
		"deliveryservices_required_capabilities": enrollDeliveryServicesRequiredCapability,
		"deliveryservice_servers":                enrollDeliveryServiceServer,
		"divisions":                              enrollDivision,
		"origins":                                enrollOrigin,
		"phys_locations":                         enrollPhysLocation,
		"regions":                                enrollRegion,
		"statuses":                               enrollStatus,
		"tenants":                                enrollTenant,
		"users":                                  enrollUser,
	}

	if len(httpPort) != 0 {
		log.Infoln("Starting http server on " + httpPort)
		err := startServer(httpPort, &toSession, dispatcher)
		if err != nil {
			log.Errorln("http server on " + httpPort + " failed: " + err.Error())
		}
	}

	if len(watchDir) != 0 {
		log.Infoln("Watching directory " + watchDir)
		dw, err := startWatching(watchDir, &toSession, dispatcher)
		defer log.Close(dw, "could not close dirwatcher")
		if err != nil {
			log.Errorf("dirwatcher on %s failed: %s", watchDir, err.Error())
		}
	}

	// create this file to indicate the enroller is ready
	f, err := os.Create(startedFile)
	if err != nil {
		panic(err)
	}
	log.Infoln("Created " + startedFile)
	log.Close(f, "could not close file")

	var waitforever chan struct{}
	<-waitforever
}
