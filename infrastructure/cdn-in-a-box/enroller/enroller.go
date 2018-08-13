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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
)

type session struct {
	*clientv13.Session
	addr net.Addr
}

func newSession(reqTimeout time.Duration, toURL string, toUser string, toPass string) (*session, error) {
	s, addr, err := clientv13.LoginWithAgent(toURL, toUser, toPass, true, "cdn-in-a-box-enroller", true, reqTimeout)

	return &session{Session: s, addr: addr}, err
}

//docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' trafficopsdb_db_1

func (s *session) inspectIPAddress(service string) (string, error) {

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
	fmt.Printf("portBytes ---> %v\n", string(portBytes))

	portStr := string(portBytes)
	portStr = strings.TrimSuffix(portStr, "\n")
	portStr = trimQuotes(portStr)

	port, err := strconv.Atoi(portStr)
	fmt.Printf("err ---> %v\n", err)
	fmt.Printf("port ---> %v\n", port)
	if err != nil {
		fmt.Printf("cannot convert portBytes to integer: %s\n", string(portBytes))
		return 0, err
	}
	return port, err
}

func (s *session) enrollService(service string) (*v13.Server, error) {
	IPAddress, err := s.inspectIPAddress(service)
	if err != nil {
		fmt.Printf("cannot lookup ipaddress: %v\n", err)
	}
	fmt.Printf("IPAddress ---> %v\n", IPAddress)

	port, err := s.inspectPort(service)
	if err != nil {
		fmt.Printf("cannot lookup port: %v", err)
	}
	fmt.Printf("port ---> %v\n", port)

	server := v13.Server{
		CDNID:          1,
		CachegroupID:   1,
		PhysLocationID: 1,
		ProfileID:      1,
		StatusID:       2,
		TypeID:         1,
		HostName:       service,
		IPAddress:      IPAddress,
		TCPPort:        port,
	}

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
	var (
		cmdOut []byte
		err    error
	)
	dockerCmd, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println("cannot find the docker executeable")
	}
	fmt.Printf("Executing: %s %v\n", dockerCmd, strings.Join(cmdArgs, " "))
	if cmdOut, err = exec.Command(dockerCmd, cmdArgs...).Output(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error running %s: %v\n", dockerCmd, err)
		os.Exit(1)
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

		dockerCli, err := client.NewEnvClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			containers, err := dockerCli.ContainerList(context.Background(), types.ContainerListOptions{})
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
