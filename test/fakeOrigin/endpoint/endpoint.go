package endpoint

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
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// DefaultConfigFile is the default configuration file path
const DefaultConfigFile = "config.json"

// DefaultOutputDirectory is the default output directory for generated content
const DefaultOutputDirectory = "./out"

// DefaultHTTPSKeyFile is the default path to the SSL key
const DefaultHTTPSKeyFile = "server.key"

// DefaultHTTPSCertFile is the default path to the SSL certificate
const DefaultHTTPSCertFile = "server.cert"

// ServerInfo contains relevant info for serving content
type ServerInfo struct {
	HTTPListeningPort  int           `json:"http_port"`
	HTTPSListeningPort int           `json:"https_port"`
	SSLCert            string        `json:"ssl_cert"`
	SSLKey             string        `json:"ssl_key"`
	BindingAddress     string        `json:"binding_address"`
	CrossdomainFile    string        `json:"crossdomain_xml_file"`
	ReadTimeout        time.Duration `json:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout"`
}

// Endpoint defines all kinds of endpoints to be served
type Endpoint struct {
	ID              string              `json:"id"`
	DiskID          string              `json:"override_disk_id,omitempty"`
	Source          string              `json:"source,omitempty"`
	OutputDirectory string              `json:"outputdir,omitempty"`
	EndpointType    Type                `json:"type"`
	ManualCommand   []string            `json:"manual_command,omitempty"`
	DefaultHeaders  map[string][]string `json:"default_headers,omitempty"`
	NoCache         bool                `json:"no_cache,omitempty"`
	ABRManifests    []string            `json:"abr_manifests,omitempty"`

	// Testing endpoint specific config
	LogReqHeaders  bool          `json:"log_request_headers,omitempty"`
	LogRespHeaders bool          `json:"log_response_headers,omitempty"`
	StallDuration  time.Duration `json:"stall_duration,omitempty"`
	EnablePprof    bool          `json:"enable_pprof,omitempty"`
	EnableDebug    bool          `json:"enable_debug,omitempty"`
}

// Config defines the application configuration
type Config struct {
	ServerConf ServerInfo `json:"server"`
	Endpoints  []Endpoint `json:"endpoints"`
}

func contains(args []string, param string) bool {
	for _, str := range args {
		if strings.Contains(str, param) {
			return true
		}
	}
	return false
}

func replace(tokenDict map[string]string, input string) string {
	for k, v := range tokenDict {
		input = strings.Replace(input, k, v, -1)
	}
	return input
}

// DefaultConfig generates a basic default configuration
func DefaultConfig() Config {
	return Config{
		ServerConf: ServerInfo{
			HTTPListeningPort:  8080,
			HTTPSListeningPort: 8443,
			SSLCert:            DefaultHTTPSCertFile,
			SSLKey:             DefaultHTTPSKeyFile,
			BindingAddress:     "",
			CrossdomainFile:    "./example/crossdomain.xml",
		},
		Endpoints: []Endpoint{
			{
				ID:              "SampleFile",
				Source:          "./example/video/kelloggs.mp4",
				OutputDirectory: DefaultOutputDirectory,
				EndpointType:    Static,
			},
			{
				ID:              "SampleDirectory",
				Source:          "./example/video/",
				OutputDirectory: DefaultOutputDirectory,
				EndpointType:    Dir,
			},
		},
	}
}

// WriteConfig saves a fakeOrigin config as a pretty-printed json file
func WriteConfig(cfg Config, path string) error {
	bts, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return errors.New("marshalling JSON: " + err.Error())
	}
	if err = ioutil.WriteFile(path, append(bts, '\n'), 0644); err != nil {
		return errors.New("writing file '" + path + "': " + err.Error())
	}
	return nil
}

// ProcessConfig processes the config loaded from disk, or generated the first time. This must be called before the config can be used to transcode or serve.
func ProcessConfig(out Config) (Config, error) {
	for i := range out.Endpoints {
		var err error
		// Resolve relative paths to absolute paths
		if out.Endpoints[i].EndpointType != Testing {
			out.Endpoints[i].Source, err = filepath.Abs(out.Endpoints[i].Source)
			if err != nil {
				return Config{}, errors.New("resolving relative path: " + err.Error())
			}
			out.Endpoints[i].OutputDirectory, _ = filepath.Abs(out.Endpoints[i].OutputDirectory)

			if out.Endpoints[i].OutputDirectory == "" && out.Endpoints[i].EndpointType != Testing {
				out.Endpoints[i].OutputDirectory = DefaultOutputDirectory
			}
		}
		if out.Endpoints[i].DiskID == "" {
			out.Endpoints[i].DiskID = out.Endpoints[i].ID
		}
	}
	return out, nil
}

// LoadAndGenerateDefaultConfig loads the config from a given json file and puts a default value in place if you havn't stored anything
func LoadAndGenerateDefaultConfig(path string) (Config, error) {
	out := DefaultConfig()
	defaultEndpoints := make([]Endpoint, len(out.Endpoints))
	copy(defaultEndpoints, out.Endpoints)
	raw, err := ioutil.ReadFile(path)
	if err != nil || len(raw) == 0 {
		raw = []byte("{}")
	}
	err = json.Unmarshal(raw, &out)
	if err != nil {
		return out, err
	}
	if err := WriteConfig(out, path); err != nil {
		return out, errors.New("writing config to file: " + err.Error())
	}
	if fmt.Sprintf("%v", out.Endpoints) == fmt.Sprintf("%v", defaultEndpoints) {
		return out, errors.New("default endpoints generated, please provide real input")
	}
	for _, ep := range out.Endpoints {
		if len(ep.ManualCommand) > 0 {
			//if !contains(ep.ManualCommand, "%MASTERMANIFEST%") {
			//	return out, errors.New("Manual commands must include the %MASTERMANIFEST% token")
			//}
			if !contains(ep.ManualCommand, `%OUTPUTDIRECTORY%`) {
				return out, errors.New(`manual commands must include the %OUTPUTDIRECTORY% token`)
			}
			if !contains(ep.ManualCommand, `%SOURCE%`) {
				return out, errors.New(`manual commands must include the %SOURCE% token`)
			}
		}
		if len(ep.ABRManifests) > 0 {
			return out, errors.New("paths of ABR Layer manifests must not be set via configuration")
		}
	}
	out, _ = ProcessConfig(out)
	if err = WriteConfig(out, path); err != nil {
		return out, errors.New("processing config file: " + err.Error())
	}
	return out, nil
}

// GetTranscoderCommand produces an instruction for the transcode phase to execute
func GetTranscoderCommand(ep Endpoint) (string, []string, error) {
	out := ""
	args := []string{}
	if ep.EndpointType.String() == Static.String() {
		out = "static"
	} else if ep.EndpointType.String() == Dir.String() {
		out = "dir"
	} else if len(ep.ManualCommand) > 0 {
		tokenmap := map[string]string{
			`%DISKID%`:          ep.DiskID,
			`%ENDPOINTTYPE%`:    ep.EndpointType.String(),
			`%ID%`:              ep.ID,
			`%MASTERMANIFEST%`:  ep.OutputDirectory + "/" + ep.DiskID + ".m3u8",
			`%OUTPUTDIRECTORY%`: ep.OutputDirectory,
			`%SOURCE%`:          ep.Source,
		}
		for _, cmdPart := range ep.ManualCommand {
			args = append(args, replace(tokenmap, cmdPart))
		}
		out = args[0]
		args = args[1:]
	}

	return out, args, nil
}
