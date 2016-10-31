
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// Args encapsulates the command line arguments
type Args struct {
	Port           int
	TrafficOpsUri  string
	TrafficOpsUser string
	TrafficOpsPass string
}

// getFlags parses and returns the command line arguments. The returned error
// will be non-nil if any expected arg is missing.
func getFlags() (Args, error) {
	var args Args
	flag.IntVar(&args.Port, "port", -1, "the port to serve on")
	flag.IntVar(&args.Port, "p", -1, "the port to serve on (shorthand)")
	flag.StringVar(&args.TrafficOpsUri, "uri", "", "the Traffic Ops URI")
	flag.StringVar(&args.TrafficOpsUri, "u", "", "the Traffic Ops URI (shorthand)")
	flag.StringVar(&args.TrafficOpsUser, "user", "", "the Traffic Ops username")
	flag.StringVar(&args.TrafficOpsUser, "U", "", "the Traffic Ops username (shorthand)")
	flag.StringVar(&args.TrafficOpsPass, "Pass", "", "the Traffic Ops password")
	flag.StringVar(&args.TrafficOpsPass, "P", "", "the Traffic Ops password (shorthand)")
	flag.Parse()
	if args.Port == -1 {
		return args, errors.New("Missing port")
	}
	if args.Port < 0 || args.Port > 65535 {
		return args, errors.New("Invalid port")
	}
	if args.TrafficOpsUri == "" {
		return args, errors.New("Missing CDN URI")
	}
	if args.TrafficOpsUser == "" {
		return args, errors.New("Missing CDN user")
	}
	if args.TrafficOpsPass == "" {
		return args, errors.New("Missing CDN password")
	}
	return args, nil
}

func printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	fmt.Println("Example: ats-config -port 3001 -uri http://my-traffic-ops.mycdn -user bill -pass thelizard")
}

// route routes HTTP requests to /traffic-server-host/config-file.config
// This should be registered with http.HandleFunc at "/"
// This could be changed to serve at an arbitrary endpoint, by removing the ^ in the regex
func route(trafficOpsUri string, trafficOpsCookie string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	hostnameRegex := `((?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*(?:[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))`
	configfileRegex := `([a-zA-z]*\.config)`
	// \todo precompile
	re := regexp.MustCompile(`^` + `/` + hostnameRegex + `/` + configfileRegex + `$`)

	url := r.URL.String()
	match := re.FindStringSubmatch(url)
	if len(match) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error this address is not a config file: '%s'", url)
		return
	}

	server := match[1]
	configFile := match[2]

	profile, err := GetServerProfileName(trafficOpsUri, trafficOpsCookie, server)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error getting profile for server '%s': '%v'", server, err)
		return
	}

	params, err := GetParameters(trafficOpsUri, trafficOpsCookie, profile) // \todo fix magic profile
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error getting parameters for server '%s' profile '%s': '%v'", server, profile, err)
		return
	}

	config, err := GetConfig(configFile, trafficOpsUri, server, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error getting config for server '%s' profile '%s': '%v'", server, profile, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, config)
}

func main() {
	args, err := getFlags()
	if err != nil {
		fmt.Println(err)
		printUsage()
		return
	}

	trafficOpsCookie, err := GetTrafficOpsCookie(args.TrafficOpsUri, args.TrafficOpsUser, args.TrafficOpsPass)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route(args.TrafficOpsUri, trafficOpsCookie, w, r)
	})

	err = http.ListenAndServe(":"+strconv.Itoa(args.Port), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
