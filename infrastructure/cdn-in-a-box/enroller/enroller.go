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
	"errors"
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

	dockerCli, err := dockerclient.NewClientWithOpts(dockerclient.WithVersion("1.38"), dockerclient.FromEnv)
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

func (s *session) getExposedPorts(c dockertypes.ContainerJSON) []int {
	var ports []int
	for port := range c.Config.ExposedPorts {
		ports = append(ports, port.Int())
	}
	return ports
}

func (s *session) getNetwork(c dockertypes.ContainerJSON) (*network.EndpointSettings, error) {
	if c.NetworkSettings == nil {
		return nil, errors.New("cannot get network from container")
	}
	mode := string(c.HostConfig.NetworkMode)
	net, ok := c.NetworkSettings.Networks[mode]
	if !ok {
		return nil, errors.New("no network for " + mode)
	}
	return net, nil
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
	"trafficops-perl":  "TRAFFIC_OPS",
	"trafficportal":    "TRAFFIC_PORTAL",
	"trafficrouter":    "CCR",
	"trafficstats":     "TRAFFIC_STATS",
	"trafficvault":     "RIAK",
}

func containerName(c dockertypes.ContainerJSON) string {
	return strings.Trim(c.Name, "/")
}

func serviceName(c dockertypes.ContainerJSON) string {
	if s, ok := c.Config.Labels["com.docker.compose.service"]; ok {
		return s
	}
	return containerName(c)
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
	if err != nil || len(types) == 0 {
		fmt.Printf("unknown type %s\n", typeName)
		return -1, err
	}
	fmt.Printf("type %s: %++v\n", typeName, types)
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
	if len(cgs) == 0 {
		return -1, errors.New("No cachegroups found")
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

func (s *session) enrollContainer(c dockertypes.ContainerJSON) (*v13.Server, error) {
	hostName := serviceName(c)
	cName := containerName(c)
	fmt.Printf("enrolling %s(%s)\n", cName, hostName)
	server := v13.Server{
		HostName:   hostName,
		DomainName: os.Getenv("DOMAINNAME"),
		HTTPSPort:  443,
	}

	fmt.Println("type is ", serverType(hostName))
	fmt.Println("hostName is ", hostName)
	var err error
	server.TypeID, err = s.getTypeIDByName(serverType(hostName))
	if err != nil {
		fmt.Printf("cannot get type for %s", hostName)
	}

	server.StatusID, err = s.getStatusIDByName("PRE_PROD")
	if err != nil {
		fmt.Printf("cannot get status for %s", hostName)
	}

	server.CDNID, err = s.getCDNIDByName("ALL")
	if err != nil {
		fmt.Printf("cannot get CDN for %s", hostName)
	}

	server.ProfileID, err = s.getProfileID()
	if err != nil {
		fmt.Printf("cannot get profile for %s", hostName)
	}

	server.CachegroupID, err = s.getCachegroupID()
	if err != nil {
		fmt.Printf("cannot get Cachegroup for %s", cName)
	}

	server.PhysLocationID, err = s.getPhysLocationID()
	if err != nil {
		fmt.Printf("cannot get PhysLocation for %s", cName)
	}
	dnet, err := s.getNetwork(c)
	if err != nil {
		fmt.Printf("cannot get network: %v\n", err)
		return nil, err
	}

	server.IPAddress = dnet.IPAddress
	server.IPNetmask = getMask(net.CIDRMask(dnet.IPPrefixLen, net.IPv4len*8))
	server.IPGateway = dnet.Gateway
	server.IP6Address = dnet.GlobalIPv6Address
	server.IP6Gateway = dnet.IPv6Gateway

	ports := s.getExposedPorts(c)
	if err != nil {
		fmt.Printf("cannot get exposed ports: %v\n", err)
		return nil, err
	}

	if len(ports) > 0 {
		// TODO: for now, assuming there's only 1
		server.TCPPort = ports[0]
		server.HTTPSPort = ports[0]
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
		hostName := r.FormValue("host")
		cName := r.FormValue("name")

		match := func(dockertypes.ContainerJSON) bool { return true }
		switch {
		case len(hostName) > 0:
			match = func(c dockertypes.ContainerJSON) bool {
				return hostName == c.Name
			}
		case len(cName) > 0:
			match = func(c dockertypes.ContainerJSON) bool {
				return cName == c.Config.Hostname
			}
		}

		net, err := s.NetworkInspect(context.Background(), "cdn-in-a-box_tcnet", dockertypes.NetworkInspectOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var names []string
		var containers []dockertypes.ContainerJSON
		for _, epr := range net.Containers {
			c, err := s.ContainerInspect(context.Background(), epr.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// check against any query params provided
			if !match(c) {
				continue
			}

			fmt.Printf("including %s\n", c.Name)
			names = append(names, c.Name)
			containers = append(containers, c)
		}

		switch r.Method {
		case "GET":
			// just list the container names
			enc := json.NewEncoder(w)
			enc.Encode(names)
			return

		case "POST":
			// enroll each container
			var servers []*v13.Server
			for _, c := range containers {
				server, err := s.enrollContainer(c)
				if err != nil {
					fmt.Printf("failed to enroll %s\n", containerName(c))
					continue
				}
				servers = append(servers, server)
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(servers); err != nil {
				fmt.Println("failed to encode servers")
			}
			return

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

	log.Fatal(http.ListenAndServeTLS(":443", "./server.crt", "./server.key", nil))
}
