
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
	configfiles "github.com/apache/trafficcontrol/traffic_ops/experimental/ats_config/config_files"
	"github.com/apache/trafficcontrol/traffic_ops/experimental/ats_config/traffic_ops"
)

// Args encapsulates the command line arguments
type Args struct {
	ConfigFile       string
	ConfigFileServer string
	TrafficOpsUri    string
	TrafficOpsUser   string
	TrafficOpsPass   string
}

// getFlags parses and returns the command line arguments. The returned error
// will be non-nil if any expected arg is missing.
func getFlags() (Args, error) {
	var args Args
	flag.StringVar(&args.ConfigFile, "configfile", "", "the config file to get")
	flag.StringVar(&args.ConfigFile, "c", "", "the config file to get (shorthand)")
	flag.StringVar(&args.ConfigFileServer, "configfilehost", "", "the host to get the config file for")
	flag.StringVar(&args.ConfigFileServer, "h", "", "the host to get the config file for (shorthand)")
	flag.StringVar(&args.TrafficOpsUri, "uri", "", "the Traffic Ops URI")
	flag.StringVar(&args.TrafficOpsUri, "u", "", "the Traffic Ops URI (shorthand)")
	flag.StringVar(&args.TrafficOpsUser, "user", "", "the Traffic Ops username")
	flag.StringVar(&args.TrafficOpsUser, "U", "", "the Traffic Ops username (shorthand)")
	flag.StringVar(&args.TrafficOpsPass, "Pass", "", "the Traffic Ops password")
	flag.StringVar(&args.TrafficOpsPass, "P", "", "the Traffic Ops password (shorthand)")
	flag.Parse()
	if args.ConfigFile == "" {
		return args, errors.New("Missing config file")
	}
	if args.ConfigFileServer == "" {
		return args, errors.New("Missing config file host")
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
	fmt.Println("Example: ats-config -uri http://my-traffic-ops.mycdn -user bill -pass thelizard -configfile storage.config -configfilehost c2-atsec-01")
}

func main() {
	args, err := getFlags()
	if err != nil {
		fmt.Println(err)
		printUsage()
		return
	}

	trafficOpsCookie, err := traffic_ops.GetCookie(args.TrafficOpsUri, args.TrafficOpsUser, args.TrafficOpsPass)
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := configfiles.Get(args.TrafficOpsUri, trafficOpsCookie, args.ConfigFileServer, args.ConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", config)
}
