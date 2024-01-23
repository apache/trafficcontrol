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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-check-refs/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

var (
	cfg          config.Cfg
	atsPlugins   = make(map[string]int)
	pluginChecks = make(map[string]bool)
	pluginParams = make(map[string]bool)
)

// This function accepts config line data from either ATS
// a 'plugin.config' or a 'remap.config' format.
//
// It checks the configuration file line by line and verifies
// that any specified plugin exists in the file system at the
// complete file path or relative to the ATS plugins installation
// directory. Also, any plugin arguments or plugin parameters that
// end in '.config', '.cfg', '.txt', '.yml', or '.yaml'  are assumed
// to be plugin configuration files and they will be verified
// that the exist at the absolute path in the file name or
// relative to the ATS configuration files directory.
//
// Returns '0' if all plugins on the config line successfully verify
// otherwise, returns the the count of plugins that failed to verify.
func checkConfigLine(line string, lineNumber int, filesAdding map[string]struct{}) int {

	pluginErrorCount := 0
	exists := false
	verified := false

	log.Debugf("line: %s\n", line)

	// create an array of whitespace delimited fields
	l := regexp.MustCompile(`\s+`)
	fields := l.Split(line, -1)
	length := len(fields)

	log.Debugf("length: %d, fields: %v", length, fields)

	// processing a line from remap.config
	if length > 3 && (fields[0] == "map" ||
		fields[0] == "map_with_recv_port" ||
		fields[0] == "map_with_referer" ||
		fields[0] == "reverse_map" ||
		fields[0] == "redirect" ||
		fields[0] == "redirect_temporary") {

		for ii := 3; ii < len(fields); ii++ {
			if strings.HasPrefix(fields[ii], "@plugin=") {
				sa := strings.Split(fields[ii], "=")
				if len(sa) != 2 {
					log.Errorf("malformed @plugin definition on line '%d'\n", lineNumber)
				} else {
					key := strings.TrimSpace(sa[1])
					verified, exists = pluginChecks[key]
					log.Debugf("Verified plugin '%s', exists: %v\n", key, verified)
					if !exists {
						verified = verifyPlugin(key)
						pluginChecks[key] = verified
					}
					if !verified {
						log.Errorf("the plugin '%s' in remap.config on line '%d' is not available to the installed trafficserver\n",
							key, lineNumber)
						pluginErrorCount++
					} else {
						log.Infof("then plugin DSO '%s' in remap.config on line '%d' has been verified\n", key, lineNumber)
					}
				}
			} else if strings.HasPrefix(fields[ii], "@pparam") {
				// any plugin parameters that end in '.config | .cfg | .txt | yml | .yaml | .lua'
				// are assumed to be configuration files and are checked that they
				// exist in the filesystem at the absolute location in the name
				// or relative to the ATS configuration files directory.
				m := regexp.MustCompile(`^.+(\.config|\.cfg|\.txt|\.yml|\.yaml|\.lua)$`)
				sa := strings.Split(fields[ii], "=")
				if len(sa) != 2 && len(sa) != 3 {
					log.Errorf("malformed @pparam definition in remap.config on line '%d': %v\n", lineNumber, fields)
					pluginErrorCount++
				} else {
					param := strings.TrimSpace(sa[1])
					if m.MatchString(param) {
						verified, exists = pluginParams[param]
						if !exists {
							verified = verifyPluginConfigfile(param, filesAdding)
							pluginParams[param] = verified
						}
						if !verified {
							log.Errorf("the plugin config file '%s' on line '%d' of remap.config does not exist or is empty\n",
								param, lineNumber)
							pluginErrorCount++
						} else {
							log.Infof("the plugin config file '%s' on line '%d' of remap.config has been verified\n",
								param, lineNumber)
						}
					}
				}
			}
		}
	} else { // process a line from plugin.config

		// process a line from plugin.config
		if length > 0 && strings.HasSuffix(fields[0], ".so") {
			key := strings.TrimSpace(fields[0])
			verified, exists = pluginChecks[key]
			if !exists {
				verified = verifyPlugin(key)
				pluginChecks[key] = verified
			}
			if !verified {
				log.Errorf("the plugin '%s' on line '%d' of plugin.config is not available to the the installed trafficserver\n",
					key, lineNumber)
				pluginErrorCount++
			} else {
				log.Infof("the plugin '%s' on line '%d' of plugin.config has been verified\n", key, lineNumber)
			}
		}
		// Check the arguments in a plugin.config file for possible plugin config files.
		// Any plugin argument that ends in '.config | .cfg | .txt | .yml | .yaml | .lua' are
		// assumed to be configuration files and are checked that they
		// exist in the filesystem at the absolute location in the name
		// or relative to the ATS configuration files directory.
		m := regexp.MustCompile(`([^=]+\.config$|[^=]\.cfg$|[^=]+\.txt$|[^=]+\.yml$|[^=]+\.yaml$|[^=]+\.lua$)`)
		for ii := 1; ii < length; ii++ {
			param := strings.TrimSpace(fields[ii])
			cfg := m.FindStringSubmatch(param)
			if len(cfg) == 2 {
				verified, exists = pluginParams[cfg[0]]
				if !exists {
					verified = verifyPluginConfigfile(cfg[0], filesAdding)
					pluginParams[cfg[0]] = verified
				}
				if !verified {
					log.Errorf("the plugin config file '%s' on line '%d' of plugin.config does not exist or is empty\n",
						cfg[0], lineNumber)
					pluginErrorCount++
				} else {
					log.Infof("the plugin config file '%s' on line '%d' of plugin.config has been verified\n", cfg[0], lineNumber)
				}
			}
		}
	}
	return pluginErrorCount
}

// returns 'filename' exists 'true' or 'false'
func fileExists(filename string) bool {
	log.Debugf("verifying plugin file at %s\n", filename)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// read the names of all available plugins in the
// installed trafficservers plugin directory.
func loadAvailablePlugins() {
	files, err := ioutil.ReadDir(cfg.TrafficServerPluginDir)
	if err != nil {
		log.Errorf("%v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			log.Debugf("loaded plugin %s\n", file.Name())
			atsPlugins[file.Name()] = 1
		}
	}
}

func verifyPluginConfigfile(filename string, filesAdding map[string]struct{}) bool {
	// filename isn't necessarily just a filename, it's whatever was in the plugin param, and can include a relative or absolute path.
	// So, get just the file name for the filesAdding check, because filesAdding is just the name.
	// TODO smarter path checking. This would wrongly succeed if a file was being created but in a different path.
	filenameForAdding := filepath.Base(filename)
	if _, ok := filesAdding[filenameForAdding]; ok {
		return true
	}
	if filepath.IsAbs(filename) {
		return fileExists(filename)
	} else {
		return fileExists(filepath.Join(cfg.TrafficServerConfigDir, filename))
	}
}

// returns plugin is verified (filename exists), 'true' or 'false'
func verifyPlugin(filename string) bool {

	if !strings.HasSuffix(filename, ".so") {
		return false
	}

	if filepath.IsAbs(filename) {
		return fileExists(filename)
	} else {
		return fileExists(filepath.Join(cfg.TrafficServerPluginDir, filename))
	}
}

func main() {
	// The count of plugins that could not be verified is returned
	// to the calling program.
	//
	// A count of '0' is successful meaning all ATS plugins named
	// in the config file have been verified to exist where
	// named or in the ATS plugins directory.
	pluginErrorCount := 0

	var err error
	cfg, err = config.InitConfig(Version, GitRevision)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	args := cfg.CommandArgs

	// load up the names of available plugins (at cfg.TrafficServerPluginDir).
	loadAvailablePlugins()

	var reader io.Reader
	var scanner *bufio.Scanner

	// open the indicated 'filename' argument or os.Stdin.
	length := len(args)
	switch length {
	case 0:
		reader = os.Stdin
	case 1:
		reader, err = os.Open(args[0])
		if err != nil {
			log.Errorf("%v\n", err)
			os.Exit(-1)
		}
	default:
		config.Usage()
		os.Exit(-1)
	}

	// If --files-adding=input, then the input (stdin or file) is a JSON object containing the file contents and the list of files to add, and we need to set cfg.FilesAdding from it.
	// Otherwise, the input is the raw file, and config loading has already set cfg.FilesAdding.
	if cfg.FilesAddingInput {
		inputObj := &t3cutil.CheckRefsInputFileAndAdding{}
		if err := json.NewDecoder(reader).Decode(inputObj); err != nil {
			log.Errorln("decoding input: " + err.Error())
			os.Exit(-1)
		}
		cfg.FilesAdding = util.SliceToSet(inputObj.Adding)
		scanner = bufio.NewScanner(bytes.NewBuffer(inputObj.File))
	} else {
		scanner = bufio.NewScanner(reader)
	}

	// process the config file contents verifying plugins.

	lineNumber := 1
	line := ""
	textArray := make([]string, 0)

	// scan the stream line by line
	for scanner.Scan() {
		text := scanner.Text()
		log.Debugf("parsing: %s\n", text)

		// skip lines beginning with a comment.
		if strings.HasPrefix(text, "#") {
			continue
		}

		textArray = append(textArray, scanner.Text())

		// check for and concatenate lines that have the '\' continuation marker
		if strings.HasSuffix(scanner.Text(), "\\") {
			lineNumber++
			continue
		}

		line = strings.Join(textArray, " ")
		line = strings.ReplaceAll(line, "\\", " ")

		pluginErrorCount += checkConfigLine(line, lineNumber, cfg.FilesAdding)
		lineNumber++
		textArray = make([]string, 0)
	}

	if pluginErrorCount > 0 {
		log.Errorf("there are '%d' plugins that could not be verified\n", pluginErrorCount)
		os.Exit(pluginErrorCount)
	} else {
		log.Infoln("All configured plugins have successfully been verified")
	}
	os.Exit(0)
}
