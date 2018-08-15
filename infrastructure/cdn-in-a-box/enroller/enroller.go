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
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	clientv13 "github.com/apache/trafficcontrol/traffic_ops/client/v13"
	dockertypes "github.com/docker/docker/api/types"
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

//docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' trafficopsdb_db_1

func printJSON(label string, b interface{}) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent(``, `  `)
	enc.Encode(b)
	fmt.Println(label, buf.String())
}

func (s *session) inspectIPAddress(service string) (string, error) {

	networks, err := s.Client.NetworkList(context.Background(), dockertypes.NetworkListOptions{})
	if err != nil {
		return "", err
	}
	printJSON("Networks: ", networks)

	const inspectIPAddressFormat = "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}"
	cmdArgs := []string{"inspect", "--format='" + inspectIPAddressFormat + "'", service}
	ipAddressBytes, err := runDockerCommand(cmdArgs)
	ipAddress := string(ipAddressBytes)
	ipAddress = strings.TrimSuffix(ipAddress, "\n")
	ipAddress = trimQuotes(ipAddress)
	return ipAddress, err
}

func (s *session) inspectPort(service string) (int, error) {
	const inspectPortFormat = "{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}}{{end}}"
	cmdArgs := []string{"inspect", "--format='" + inspectPortFormat + "'", service}
	portBytes, err := runDockerCommand(cmdArgs)
	if err != nil {
		fmt.Printf("cannot runDockerCommand: %s", cmdArgs)
		return 0, err
	}

	portStr := string(portBytes)
	portStr = strings.TrimSuffix(portStr, "\n")
	portStr = trimQuotes(portStr)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("cannot convert portBytes to integer: %s\n", string(portBytes))
		return 0, err
	}
	return port, err
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

func serverType(service string) string {
	for s, t := range serviceTypes {
		if strings.Contains(service, s) {
			return t
		}
	}
	// unknown -- let caller deal with it
	return service
}

func (s *session) getTypeIDByName(typeName string) (int, error) {
	types, _, err := s.GetTypeByName(typeName)
	if err != nil {
		fmt.Printf("unknown type %s\n", typeName)
		return -1, err
	}
	return types[0].ID, err
}

func (s *session) getCDNID() (int, error) {
	cdns, _, err := s.GetCDNs()
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

func (s *session) enrollService(service string) (*v13.Server, error) {
	server := v13.Server{
		HostName: service,
	}
	var err error

	server.TypeID, err = s.getTypeIDByName(serverType(service))
	if err != nil {
		fmt.Printf("cannot get type for %s", service)
	}

	server.StatusID, err = s.getStatusIDByName("PRE_PROD")
	if err != nil {
		fmt.Printf("cannot get status for %s", service)
	}

	server.CDNID, err = s.getCDNID()
	if err != nil {
		fmt.Printf("cannot get CDN for %s", service)
	}

	server.ProfileID, err = s.getProfileID()
	if err != nil {
		fmt.Printf("cannot get profile for %s", service)
	}

	server.CachegroupID, err = s.getCachegroupID()
	if err != nil {
		fmt.Printf("cannot get Cachegroup for %s", service)
	}

	server.PhysLocationID, err = s.getPhysLocationID()
	if err != nil {
		fmt.Printf("cannot get PhysLocation for %s", service)
	}

	server.IPAddress, err = s.inspectIPAddress(service)
	if err != nil {
		fmt.Printf("cannot lookup ipaddress: %v\n", err)
	}

	server.TCPPort, err = s.inspectPort(service)
	if err != nil {
		fmt.Printf("cannot lookup port: %v", err)
		return nil, err
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

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func runDockerCommand(cmdArgs []string) ([]byte, error) {
	dockerCmd, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println("cannot find the docker executeable")
		return nil, err
	}
	fmt.Printf("Executing: %s %v\n", dockerCmd, strings.Join(cmdArgs, " "))
	cmdOut, err := exec.Command(dockerCmd, cmdArgs...).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error running %s: %v\n", dockerCmd, err)
	}
	return cmdOut, err
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
		service := r.FormValue("service")

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
			enc.Encode(map[string][]string{service: names})
			return

		case "POST":
			fmt.Println("enrolling ", service)
			server, err := s.enrollService(service)
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
