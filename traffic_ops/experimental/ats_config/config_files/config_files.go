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
	"regexp"
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

	return GetConfig(toClient, configFileName, toURL, serverHostname, params)
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

type ConfigFileCreatorFunc func(toClient towrap.ITrafficOpsSession, filename string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error)

// ConfigFileFuncMap returns the dispatch map, of regular expressions to config file creator functions
// TODO change apps to cache this, namely for the long-running service to only compile the regexes once. Or, put in init()?
func ConfigFileFuncMap() map[*regexp.Regexp]ConfigFileCreatorFunc {
	spaceCreateConfig := createGenericDotConfigFunc(" ")
	spacedEqualsCreateConfig := createGenericDotConfigFunc(" = ")
	equalsCreateConfig := createGenericDotConfigFunc("=")
	return map[*regexp.Regexp]ConfigFileCreatorFunc{
		regexp.MustCompile(`^storage\.config$`):          createStorageDotConfig,
		regexp.MustCompile(`^volume\.config$`):           createVolumeDotConfig,
		regexp.MustCompile(`^logs_xml\.config$`):         createLogsXmlDotConfig,
		regexp.MustCompile(`^cacheurl\.config$`):         createCacheurlDotConfig,
		regexp.MustCompile(`^cacheurl_qstring\.config$`): createCacheurlQstringDotConfig,
		regexp.MustCompile(`^cacheurl_(.*)\.config$`):    createCacheurlStarDotConfig,
		regexp.MustCompile(`^records\.config$`):          spaceCreateConfig,
		regexp.MustCompile(`^plugin\.config$`):           spaceCreateConfig,
		regexp.MustCompile(`^astats\.config$`):           equalsCreateConfig,
		regexp.MustCompile(`^sysctl\.config$`):           spacedEqualsCreateConfig,
		regexp.MustCompile(`^hosting\.config$`):          createHostingDotConfig,
		regexp.MustCompile(`^50-ats\.rules$`):            createAtsDotRules,
		regexp.MustCompile(`^cache\.config$`):            createCacheDotConfig,
	}
}

// GetConfig takes the name of the config file, and the Traffic Ops parameters for a server,
// and returns the text of that config file for that server.
func GetConfig(toClient towrap.ITrafficOpsSession, configFileName string, trafficOpsHost string, trafficServerHost string, params []to.Parameter) (string, error) {
	fileFuncs := ConfigFileFuncMap()
	for r, f := range fileFuncs {
		if r.MatchString(configFileName) {
			return f(toClient, configFileName, trafficOpsHost, trafficServerHost, params)
		}
	}
	return "", fmt.Errorf("Config file '%s' not valid", configFileName)
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
