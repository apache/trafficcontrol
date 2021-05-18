package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"

	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3c-preprocess/util"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
)

var returnRegex = regexp.MustCompile(`\s*__RETURN__\s*`)

func PreprocessConfigFile(server *atscfg.Server, cfgFile string) string {
	if server.TCPPort != nil && *server.TCPPort != 80 && *server.TCPPort != 0 {
		cfgFile = strings.Replace(cfgFile, `__SERVER_TCP_PORT__`, strconv.Itoa(*server.TCPPort), -1)
	} else {
		cfgFile = strings.Replace(cfgFile, `:__SERVER_TCP_PORT__`, ``, -1)
	}

	ipAddr := ""
	for _, iFace := range server.Interfaces {
		for _, addr := range iFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}
			addrStr := addr.Address
			ip := net.ParseIP(addrStr)
			if ip == nil {
				err := error(nil)
				ip, _, err = net.ParseCIDR(addrStr)
				if err != nil {
					ip = nil // don't bother with the error, just skip
				}
			}
			if ip == nil || ip.To4() == nil {
				continue
			}
			ipAddr = addrStr
			break
		}
	}
	if ipAddr != "" {
		cfgFile = strings.Replace(cfgFile, `__CACHE_IPV4__`, ipAddr, -1)
	} else {
		log.Errorln("Preprocessing: this server had a missing or malformed IPv4 Service Interface, cannot replace __CACHE_IPV4__ directives!")
	}

	if server.HostName == nil || *server.HostName == "" {
		log.Errorln("Preprocessing: this server missing HostName, cannot replace __HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__HOSTNAME__`, *server.HostName, -1)
	}
	if server.HostName == nil || *server.HostName == "" || server.DomainName == nil || *server.DomainName == "" {
		log.Errorln("Preprocessing: this server missing HostName or DomainName, cannot replace __FULL_HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__FULL_HOSTNAME__`, *server.HostName+`.`+*server.DomainName, -1)
	}
	cfgFile = returnRegex.ReplaceAllString(cfgFile, "\n")
	return cfgFile
}

func main() {
	dataFiles := &DataAndFiles{}
	if err := json.NewDecoder(os.Stdin).Decode(dataFiles); err != nil {
		log.Errorln("Error reading json input")
	}

	for fileI, file := range dataFiles.Files {
		txt := PreprocessConfigFile(dataFiles.Data.Server, file.Text)
		dataFiles.Files[fileI].Text = txt
	}
	sort.Sort(t3cutil.ATSConfigFiles(dataFiles.Files))
	if err := util.WriteConfigs(dataFiles.Files, os.Stdout); err != nil {
		hostName := ""
		if dataFiles.Data.Server.HostName != nil {
			hostName = *dataFiles.Data.Server.HostName
		}
		log.Errorln("Writing configs for server '" + hostName + "': " + err.Error())
		os.Exit(ExitCodeErrGeneric)
	}
}

const ExitCodeErrGeneric = 1

type DataAndFiles struct {
	Data  t3cutil.ConfigData      `json:"data"`
	Files []t3cutil.ATSConfigFile `json:"files"`
}
