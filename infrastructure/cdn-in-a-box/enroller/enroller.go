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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	clientv13 "github.com/apache/trafficcontrol/traffic_ops/client/v13"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
)

type session struct {
	*clientv13.Session
	*dockerclient.Client
	addr net.Addr
}

func newSession(reqTimeout time.Duration, toURL string, toUser string, toPass string) (*session, error) {
	s, addr, err := clientv13.LoginWithAgent(toURL, toUser, toPass, true, "cdn-in-a-box-enroller", true, reqTimeout)
	if err != nil {
		return nil, err
	}

	dockerCli, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &session{Session: s, addr: addr, Client: dockerCli}, err
}

func printJSON(label string, b interface{}) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent(``, `  `)
	enc.Encode(b)
	fmt.Println(label, buf.String())
}

func (s *session) getExposedPorts(containerName string) ([]int, error) {
	c, err := s.Client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return nil, err
	}

	var ports []int
	for port := range c.Config.ExposedPorts {
		ports = append(ports, port.Int())
	}
	return ports, nil
}

func (s *session) getNetwork(containerName string) (*network.EndpointSettings, error) {
	c, err := s.Client.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return nil, err
	}

	networkName := c.HostConfig.NetworkMode.UserDefined()
	net := c.NetworkSettings.Networks[networkName]
	return net, err
}

// Matches service name (container) with type in traffic ops db
var serviceTypes = map[string]string{
	"db":               "TRAFFIC_OPS_DB",
	"edge":             "EDGE",
	"influxdb":         "INFLUXDB",
	"mid":              "MID",
	"origin":           "ORG",
	"trafficanalytics": "TRAFFIC_ANALYTICS",
	"trafficmonitor":   "RASCAL",
	"trafficops":       "TRAFFIC_OPS",
	"trafficportal":    "TRAFFIC_PORTAL",
	"trafficrouter":    "CCR",
	"trafficstats":     "TRAFFIC_STATS",
	"trafficvault":     "RIAK",
}

func serverType(serviceName string) string {
	for s, t := range serviceTypes {
		if s == serviceName {
			return t
		}
	}
	// unknown -- let caller deal with it
	return serviceName
}

func (s *session) getTypeIDByName(typeName string) (int, error) {
	types, _, err := s.GetTypeByName(typeName)
	if err != nil {
		fmt.Printf("unknown type %s\n", typeName)
		return -1, err
	}
	return types[0].ID, err
}

func (s *session) getCDNIDByName(name string) (int, error) {
	cdns, _, err := s.GetCDNByName(name)
	if err != nil {
		fmt.Println("cannot get CDNS")
		return -1, err
	}
	if len(cdns) < 1 {
		panic(fmt.Sprintf("CDNS: %v;  err: %v", cdns, err))
	}
	return cdns[0].ID, err
}

func (s *session) getCachegroupID() (int, error) {
	cgs, _, err := s.GetCacheGroups()
	if err != nil {
		fmt.Println("cannot get Cachegroup")
		return -1, err
	}
	return cgs[0].ID, err
}

func (s *session) getPhysLocationID() (int, error) {
	physLocs, _, err := s.GetPhysLocations()
	if err != nil {
		fmt.Println("cannot get physlocations")
		return -1, err
	}
	return physLocs[0].ID, err
}

func (s *session) getProfileID() (int, error) {
	profiles, _, err := s.GetProfiles()
	if err != nil {
		fmt.Println("cannot get profiles")
		return -1, err
	}
	return profiles[0].ID, err
}

func (s *session) getStatusIDByName(cdnName string) (int, error) {
	statuses, _, err := s.GetStatusByName(cdnName)
	if err != nil {
		fmt.Printf("unknown Status %s\n", cdnName)
		return -1, err
	}
	return statuses[0].ID, err
}

func getMask(m []byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func (s *session) enrollService(containerName string) (*v13.Server, error) {

	split := strings.Split(containerName, `_`)
	serviceName := containerName
	if len(split) > 2 {
		serviceName = split[1]
	}
	server := v13.Server{
		HostName:   serviceName,
		DomainName: os.Getenv("DOMAINNAME"),
	}
	var err error

	fmt.Println("type is ", serverType(serviceName))
	server.TypeID, err = s.getTypeIDByName(serverType(serviceName))
	if err != nil {
		fmt.Printf("cannot get type for %s", serviceName)
	}

	server.StatusID, err = s.getStatusIDByName("PRE_PROD")
	if err != nil {
		fmt.Printf("cannot get status for %s", serviceName)
	}

	server.CDNID, err = s.getCDNIDByName("ALL")
	if err != nil {
		fmt.Printf("cannot get CDN for %s", serviceName)
	}

	server.ProfileID, err = s.getProfileID()
	if err != nil {
		fmt.Printf("cannot get profile for %s", serviceName)
	}

	server.CachegroupID, err = s.getCachegroupID()
	if err != nil {
		fmt.Printf("cannot get Cachegroup for %s", containerName)
	}

	server.PhysLocationID, err = s.getPhysLocationID()
	if err != nil {
		fmt.Printf("cannot get PhysLocation for %s", containerName)
	}

	dnet, err := s.getNetwork(containerName)
	if err != nil {
		fmt.Printf("cannot get network: %v\n", err)
		return nil, err
	}

	server.IPAddress = dnet.IPAddress
	server.IPNetmask = getMask(net.CIDRMask(dnet.IPPrefixLen, net.IPv4len*8))
	server.IPGateway = dnet.Gateway

	ports, err := s.getExposedPorts(containerName)
	if err != nil {
		fmt.Printf("cannot get exposed ports: %v\n", err)
		return nil, err
	}

	if len(ports) > 0 {
		// TODO: for now, assuming there's only 1
		server.TCPPort = ports[0]
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent(``, `  `)
	enc.Encode(server)
	fmt.Println("Server: ", buf.String())

	resp, _, err := s.CreateServer(server)
	fmt.Printf("Response: %s\n", resp)
	return &server, err
}

var to struct {
	URL      string `envconfig:"TO_URL"`
	User     string `envconfig:"TO_USER"`
	Password string `envconfig:"TO_PASSWORD"`
}

func (s *session) enrollerHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		containerName := r.FormValue("containerName")

		switch r.Method {
		case "GET":
			containers, err := s.ContainerList(context.Background(), dockertypes.ContainerListOptions{})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var names []string
			for _, container := range containers {
				name := container.ID
				if len(container.Names) > 0 {
					name = container.Names[0]
				}
				names = append(names, name)
			}
			enc := json.NewEncoder(w)
			enc.Encode(map[string][]string{containerName: names})
			return

		case "POST":
			fmt.Println("enrolling ", containerName)
			server, err := s.enrollService(containerName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(server); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "unhandled method "+r.Method, http.StatusBadRequest)
			return
		}
	}
}

func main() {
	envconfig.Process("", &to)
	reqTimeout := time.Second * time.Duration(60)

	toSession, err := newSession(reqTimeout, to.URL, to.User, to.Password)
	if err != nil {
		panic(err)
	}
	fmt.Println("TO session established")
	http.HandleFunc("/", toSession.enrollerHandler())

	log.Fatal(http.ListenAndServeTLS(":8080", "./server.crt", "./server.key", nil))
}
