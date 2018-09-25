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

// TODO: Some GetxxxByxxx() methods escape the string passed in; others don't
//  Here we escape the name if not escaped in the Getxxx method being called
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

func (s session) getCoordinateIDByName(n string) (int, error) {
	coordinates, _, err := s.GetCoordinateByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(coordinates) == 0 {
		return -1, errors.New("no coordinate with name " + n)
	}
	return coordinates[0].ID, err
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
	profiles, _, err := s.GetProfileByName(n)
	if err != nil {
		return -1, err
	}
	if len(profiles) == 0 {
		return -1, errors.New("no profile with name " + n)
	}
	return profiles[0].ID, err
}

func (s session) getParameterIDMatching(m tc.Parameter) (int, error) {
	// TODO: s.GetParameterByxxx() does not seem to work with values with spaces --
	// doing this the hard way for now
	parameters, _, err := s.GetParameters()
	if err != nil {
		return -1, err
	}
	for _, p := range parameters {
		if p.Name == m.Name && p.Value == m.Value && p.ConfigFile == m.ConfigFile {
			return p.ID, nil
		}
	}
	return -1, fmt.Errorf("no parameter matching name %s, configFile %s, value %s", m.Name, m.ConfigFile, m.Value)
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

func (s session) getRoleIDByName(n string) (int, error) {
	roles, _, _, err := s.GetRoleByName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(roles) == 0 || roles[0].ID == nil {
		return -1, errors.New("no role with name " + n)
	}
	return *roles[0].ID, err
}

func (s session) getServerIDByHostName(n string) (int, error) {
	servers, _, err := s.GetServerByHostName(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(servers) == 0 {
		return -1, errors.New("no server with hostName " + n)
	}
	return servers[0].ID, err
}

func (s session) getDeliveryServiceIDByXMLID(n string) (int, error) {
	dses, _, err := s.GetDeliveryServiceByXMLID(url.QueryEscape(n))
	if err != nil {
		return -1, err
	}
	if len(dses) == 0 {
		return -1, errors.New("no deliveryservice with name " + n)
	}
	return dses[0].ID, err
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("type %s already exists\n", s.Name)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("cdn %s already exists\n", s.Name)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("asn %d already exists\n", s.ASN)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("cachegroup %s already exists\n", *s.Name)
			return nil
		}
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
	var s tc.DeliveryServiceNullable
	err = dec.Decode(&s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.Type != nil && *s.Type != "" {
		id, err := toSession.getTypeIDByName(s.Type.String())
		if err != nil {
			return err
		}
		s.TypeID = &id
	}

	if s.CDNName != nil && *s.CDNName != "" {
		id, err := toSession.getCDNIDByName(*s.CDNName)
		if err != nil {
			return err
		}
		s.CDNID = &id
	}

	if s.ProfileName != nil && *s.ProfileName != "" {
		id, err := toSession.getProfileIDByName(*s.ProfileName)
		if err != nil {
			return err
		}
		s.ProfileID = &id
	}

	if s.Tenant != nil && *s.Tenant != "" {
		id, err := toSession.getTenantIDByName(*s.Tenant)
		if err != nil {
			return err
		}
		s.TenantID = &id
	}

	alerts, err := toSession.CreateDeliveryServiceNullable(&s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("deliveryservice %s already exists\n", s.XMLID)
			return nil
		}
		log.Printf("error creating from %s: %s\n", fn, err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollDeliveryServiceServer(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)

	// DeliveryServiceServers lists ds xmlid and array of server names.  Use that to create multiple DeliveryServiceServer objects
	var dss tc.DeliveryServiceServers
	err = dec.Decode(&dss)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	dsID, err := toSession.getDeliveryServiceIDByXMLID(dss.XmlId)
	if err != nil {
		return err
	}

	var serverIDs []int
	for _, sn := range dss.ServerNames {
		id, err := toSession.getServerIDByHostName(sn)
		if err != nil {
			log.Println("error finding " + sn + ": " + err.Error())
			continue
		}
		serverIDs = append(serverIDs, id)
	}
	_, err = toSession.CreateDeliveryServiceServers(dsID, serverIDs, true)
	if err != nil {
		log.Printf("error creating from %s: %s\n", fn, err)
	}

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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("division %s already exists\n", s.Name)
			return nil
		}
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

	if s.Cachegroup != nil && *s.Cachegroup != "" {
		id, err := toSession.getCachegroupIDByName(*s.Cachegroup)
		if err != nil {
			return err
		}
		s.CachegroupID = &id
	}

	if s.DeliveryService != nil && *s.DeliveryService != "" {
		id, err := toSession.getDeliveryServiceIDByXMLID(*s.DeliveryService)
		if err != nil {
			return err
		}
		s.DeliveryServiceID = &id
	}

	if s.Profile != nil && *s.Profile != "" {
		id, err := toSession.getProfileIDByName(*s.Profile)
		if err != nil {
			return err
		}
		s.ProfileID = &id
	}

	if s.Coordinate != nil && *s.Coordinate != "" {
		id, err := toSession.getCoordinateIDByName(*s.Coordinate)
		if err != nil {
			return err
		}
		s.CoordinateID = &id
	}

	if s.Tenant != nil && *s.Tenant != "" {
		id, err := toSession.getTenantIDByName(*s.Tenant)
		if err != nil {
			return err
		}
		s.TenantID = &id
	}

	alerts, _, err := toSession.CreateOrigin(s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("origin %s already exists\n", *s.Name)
			return nil
		}
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
	paramID, err := toSession.getParameterIDMatching(s)
	var alerts tc.Alerts
	if err == nil {
		// existing param -- update
		alerts, _, err = toSession.UpdateParameterByID(paramID, s)
		if err != nil {
			log.Printf("error updating parameter %d: %s with %+v ", paramID, err.Error(), s)
		}
	} else {
		alerts, _, err = toSession.CreateParameter(s)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("parameter %s already exists\n", s.Name)
				return nil
			}
			log.Printf("error creating from %s: %s\n", fn, err)
			return err
		}
	}
	// link parameter with profiles
	if len(s.Profiles) > 0 {
		paramID, err := toSession.getParameterIDMatching(s)
		if err != nil {
			return err
		}

		var profiles []string
		err = json.Unmarshal(s.Profiles, &profiles)
		if err != nil {
			log.Printf("%v", err)
		}

		for _, n := range profiles {
			pid, err := toSession.getProfileIDByName(n)
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			pp := tc.ProfileParameter{ParameterID: paramID, ProfileID: pid}
			_, _, err = toSession.CreateProfileParameter(pp)
			if err != nil {
				if strings.Contains(err.Error(), "already exists") {
					continue
				}
				log.Printf("%v", err)
				continue
			}
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("physLocation %s already exists\n", s.Name)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("region %s already exists\n", s.Name)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("status %s already exists\n", s.Name)
			return nil
		}
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
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("tenant %s already exists\n", s.Name)
			return nil
		}
		log.Printf("error creating from %s: %s\n", fn, err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

func enrollUser(toSession *session, fn string) error {
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() {
		fh.Close()
	}()

	dec := json.NewDecoder(fh)
	var s tc.User
	err = dec.Decode(&s)
	log.Printf("User is %++v\n", s)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	if s.Tenant != nil && *s.Tenant != "" {
		id, err := toSession.getTenantIDByName(*s.Tenant)
		if err != nil {
			return err
		}
		s.TenantID = &id
	}

	if s.RoleName != nil && *s.RoleName != "" {
		id, err := toSession.getRoleIDByName(*s.RoleName)
		if err != nil {
			return err
		}
		s.Role = &id
	}

	alerts, _, err := toSession.CreateUser(&s)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("user %s already exists\n", *s.Username)
			return nil
		}
		log.Printf("error creating from %s: %s\n", fn, err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&alerts)

	return err
}

// using documented import form for profiles
type profileImport struct {
	Profile    tc.Profile             `json:"profile"`
	Parameters []tc.ParameterNullable `json:"parameters"`
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
	var profile tc.Profile

	err = dec.Decode(&profile)
	if err != nil && err != io.EOF {
		log.Printf("error decoding %s: %s\n", fn, err)
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("  ", "")
	enc.Encode(profile)

	if len(profile.Name) == 0 {
		log.Println("missing name on profile from " + fn)
		return errors.New("missing name on profile from " + fn)
	}

	profiles, _, err := toSession.GetProfileByName(profile.Name)

	createProfile := false
	if err != nil || len(profiles) == 0 {
		// no profile by that name -- need to create it
		createProfile = true
	} else {
		// updating - ID needs to match
		profile.ID = profiles[0].ID
	}

	// these need to be done whether creating or updating
	if profile.CDNName != "" {
		id, err := toSession.getCDNIDByName(profile.CDNName)
		if err != nil {
			return err
		}
		profile.CDNID = id
	}

	var alerts tc.Alerts
	var action string
	if createProfile {
		alerts, _, err = toSession.CreateProfile(profile)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("profile %s already exists\n", profile.Name)
			} else {
				log.Printf("error creating profile from %+v: %s\n", profile, err.Error())
			}
		}
		profile.ID, err = toSession.getProfileIDByName(profile.Name)
		if err != nil {
			log.Printf("error getting profile ID from %+v: %s\n", profile, err.Error())
		}
		action = "creating"
	} else {
		alerts, _, err = toSession.UpdateProfileByID(profile.ID, profile)
		action = "updating"
	}

	if err != nil {
		log.Printf("error "+action+" from %s: %s\n", fn, err)
		return err
	}

	//log.Printf("total profile is  %+v\n", profile)
	for _, p := range profile.Parameters {
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
		log.Printf("creating param %+v\n", param)
		id, err := toSession.getParameterIDMatching(param)
		if err != nil {
			// create it
			_, _, err = toSession.CreateParameter(param)
			if err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					log.Printf("can't create parameter %+v: %s\n", param, err.Error())
				}
				continue
			}
			param.ID, err = toSession.getParameterIDMatching(param)
			if err != nil {
				log.Printf("error getting new parameter %+v\n", param)
				param.ID, err = toSession.getParameterIDMatching(param)
				log.Printf(err.Error())

			}
		} else {
			param.ID = id
			toSession.UpdateParameterByID(param.ID, param)
		}

		if param.ID < 1 {
			panic(fmt.Sprintf("param ID not found for %v", param))

		}
		pp := tc.ProfileParameter{ProfileID: profile.ID, ParameterID: param.ID}
		_, _, err = toSession.CreateProfileParameter(pp)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.Printf("error creating profileparameter %+v: %s\n", pp, err.Error())
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
		id, err := toSession.getPhysLocationIDByName(s.PhysLocation)
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
					i, err := os.Stat(event.Name)
					if err != nil || i.IsDir() {
						log.Println("skipping " + event.Name)
						continue
					}
					log.Println("new file :", event.Name)
					p := strings.IndexRune(event.Name, '/')
					if p == -1 {
						continue
					}
					dir := event.Name[:p]
					suffix := rejected
					if f, ok := dw.watched[dir]; ok {
						log.Printf("creating from %s\n", event.Name)
						// TODO: ensure file content is there before attempting to read
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
		if err = os.Mkdir(dir, os.ModeDir|0700); err != nil {
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
	dw.watch("deliveryservice_servers", enrollDeliveryServiceServer)
	dw.watch("divisions", enrollDivision)
	dw.watch("origins", enrollOrigin)
	dw.watch("phys_locations", enrollPhysLocation)
	dw.watch("regions", enrollRegion)
	dw.watch("statuses", enrollStatus)
	dw.watch("tenants", enrollTenant)
	dw.watch("users", enrollUser)

	// create this file to indicate the enroller is ready
	f, err := os.Create(startedFile)
	if err != nil {
		panic(err)
	}
	f.Close()

	var waitforever chan struct{}
	<-waitforever
}
