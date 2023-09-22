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
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-preprocess/util"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-preprocess"

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

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

	if server.HostName == "" {
		log.Errorln("Preprocessing: this server missing HostName, cannot replace __HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__HOSTNAME__`, server.HostName, -1)
	}
	if server.HostName == "" || server.DomainName == "" {
		log.Errorln("Preprocessing: this server missing HostName or DomainName, cannot replace __FULL_HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__FULL_HOSTNAME__`, server.HostName+`.`+server.DomainName, -1)
	}
	if server.CacheGroup != "" {
		cfgFile = strings.Replace(cfgFile, `__CACHEGROUP__`, server.CacheGroup, -1)
	} else {
		log.Errorln("Preprocessing: this server missing Cachegroup, cannot replace __CACHEGROUP__ directives!")
	}

	cfgFile = returnRegex.ReplaceAllString(cfgFile, "\n")
	return cfgFile
}

func main() {
	flagHelp := getopt.BoolLong("help", 'h', "Print usage information and exit")
	flagVersion := getopt.BoolLong("version", 'V', "Print version information and exit.")
	flagVerbose := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	flagSilent := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	getopt.Parse()
	if *flagHelp {
		fmt.Println(usageStr())
		os.Exit(0)
	} else if *flagVersion {
		fmt.Println(t3cutil.VersionStr(AppName, Version, GitRevision))
		os.Exit(0)
	}

	logErr := io.WriteCloser(os.Stderr)
	logWarn := io.WriteCloser(nil)
	logInf := io.WriteCloser(nil)
	logDebug := io.WriteCloser(nil)
	if *flagSilent {
		logErr = io.WriteCloser(nil)
	} else {
		if *flagVerbose >= 1 {
			logWarn = os.Stderr
		}
		if *flagVerbose >= 2 {
			logInf = os.Stderr
			logDebug = os.Stderr
		}
	}
	log.Init(nil, logErr, logWarn, logInf, logDebug)

	// TODO read log location arguments
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
		if dataFiles.Data.Server.HostName == "" {
			hostName = dataFiles.Data.Server.HostName
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

func usageStr() string {
	return `usage: t3c-preprocess [--help] [--version]
       <command> [<args>]

The 't3c-preprocess' app preprocesses generated config files, replacing directives with relevant data.

The stdin must be the JSON '{"data": \<data\>, "files": \<files\>}' where \<data\> is the output of 't3c-request --get-data=config' and \<files\> is the output of 't3c-generate'.
`
}
