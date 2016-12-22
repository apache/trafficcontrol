// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// To add a config file:
//   add your function which creates the text of the config file
//   to the case statement in `GetConfig`
//

package config_files

import (
	"fmt"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"time"
)

var ErrServerNotFound = fmt.Errorf("Server not found")

// Get returns the requested config for the requested server. This is a convenience function equivalent to GetConfig(cfgFile, toClient.URL(), server, toClient.Parameters)
func Get(toClient towrap.ITrafficOpsSession, serverHostname string, configFileName string) (string, error) {
	profile, err := GetServerProfileName(toClient, serverHostname)
	if err != nil {
		return "", fmt.Errorf("Error getting server profile name: %v", err)
	}

	toURL, err := toClient.URL()
	if err != nil {
		return "", fmt.Errorf("Error getting Traffic Ops URL: %v", err)
	}

	params, err := toClient.Parameters(profile)
	if err != nil {
		return "", fmt.Errorf("Error getting Traffic Ops parameters: %v", err)
	}

	return GetConfig(configFileName, toURL, serverHostname, params)
}

// GetServerProfileName returns the name of the given server's profile in Traffic Ops.
// TODO move to a utiliy package/file?
// TODO add to.Client.Server(name) for efficiency
func GetServerProfileName(toClient towrap.ITrafficOpsSession, serverHostname string) (string, error) {
	servers, err := toClient.Servers()
	if err != nil {
		return "", err
	}
	for _, server := range servers {
		if server.HostName == serverHostname {
			return server.Profile, nil
		}
	}
	return "", ErrServerNotFound
}

// GetConfig takes the name of the config file, and the Traffic Ops parameters for a server,
// and returns the text of that config file for that server.
func GetConfig(configFileName string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	switch configFileName {
	case "storage.config":
		return createStorageDotConfig(trafficOpsHost, trafficServerHost, params)
	case "volume.config":
		return createVolumeDotConfig(trafficOpsHost, trafficServerHost, params)
	default:
		return "", fmt.Errorf("Config file '%s' not valid", configFileName)
	}
}

// createParamsMap returns a map[ConfigFile]map[ParameterName]ParameterValue.
// Helper function for createStorageDotConfig.
func createParamsMap(params []to.Parameter) map[string]map[string]string {
	m := make(map[string]map[string]string)
	for _, param := range params {
		if m[param.ConfigFile] == nil {
			m[param.ConfigFile] = make(map[string]string)
		}
		m[param.ConfigFile][param.Name] = param.Value
	}
	return m
}

func createStorageDotConfig(trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	// # DO NOT EDIT - Generated for my-edge-0 by Traffic Ops (https://localhost) on Fri Feb 19 22:16:34 UTC 2016
	// /dev/ram0 volume=1
	// /dev/ram1 volume=2
	// /dev/ram2 volume=3

	s := "# DO NOT EDIT - Generated for " + trafficServerHost + " by Traffic Ops (" + trafficOpsHost + ") on " + time.Now().String() + "\n"

	paramMap := createParamsMap(params)

	if _, ok := paramMap["storage.config"]; !ok {
		return "", fmt.Errorf("No storage config parameters")
	}

	storageConfigParams := paramMap["storage.config"]

	volumePrefixes := []string{"", "RAM_", "SSD_"}

	numVolumes := 0
	for _, prefix := range volumePrefixes {
		if _, ok := storageConfigParams[prefix+"Drive_Prefix"]; ok {
			numVolumes++
		}
	}

	hasMultipleVolumes := numVolumes > 1

	volumeText := func(volume, prefix, letters string, hasMultipleVolumes bool) string {
		s := ""
		lettersSlice := strings.Split(letters, ",")
		for _, letter := range lettersSlice {
			s += prefix + letter
			if hasMultipleVolumes {
				s += " volume=" + volume
			}
			s += "\n"
		}
		return s
	}

	for _, prefix := range volumePrefixes {
		volumeParamName := "Volume"
		if prefix != "" {
			volumeParamName = prefix + volumeParamName
		} else {
			volumeParamName = "Disk_" + volumeParamName
		}

		drivePrefix, hasDrivePrefix := storageConfigParams[prefix+"Drive_Prefix"]
		driveLetters, hasDriveLetters := storageConfigParams[prefix+"Drive_Letters"]
		driveVolume, hasDriveVolume := storageConfigParams[volumeParamName]
		if hasDrivePrefix && hasDriveLetters && hasDriveVolume {
			s += volumeText(driveVolume, drivePrefix, driveLetters, hasMultipleVolumes)
		}
	}
	return s, nil
}
