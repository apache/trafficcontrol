package cfgfile

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
	"io"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// GetAllConfigs gets all config files for cfg.CacheHostName.
func GetAllConfigs(
	toData *t3cutil.ConfigData,
	cfg config.Cfg,
) ([]t3cutil.ATSConfigFile, error) {
	if toData.Server.HostName == nil {
		return nil, errors.New("server hostname is nil")
	}

	configFiles, warnings, err := MakeConfigFilesList(toData, cfg.Dir)
	logWarnings("generating config files list: ", warnings)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	genTime := time.Now()
	hdrCommentTxt := makeHeaderComment(*toData.Server.HostName, cfg.AppVersion(), toData.TrafficOpsURL, toData.TrafficOpsAddresses, genTime)

	configs := []t3cutil.ATSConfigFile{}
	for _, fi := range configFiles {
		if cfg.RevalOnly && fi.Name != atscfg.RegexRevalidateFileName {
			continue
		}
		txt, contentType, secure, lineComment, metaData, warnings, err := GetConfigFile(toData, fi, hdrCommentTxt, cfg)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.Name + "': " + err.Error())
		}

		metaDataBts, err := json.Marshal(metaData)
		if err != nil {
			return nil, errors.New("marshalling config file '" + fi.Name + "' metadata: " + err.Error())
		}

		configs = append(configs, t3cutil.ATSConfigFile{
			Name:        fi.Name,
			Path:        fi.Path,
			Text:        txt,
			Secure:      secure,
			ContentType: contentType,
			LineComment: lineComment,
			MetaData:    metaDataBts,
			Warnings:    warnings,
		})
	}

	// TODO currently "EDGE" type servers accept client requests, and "MID" don't.
	//      But in the future, Topologies shouldn't care about types, and this
	//      should have the logic for which DSes need which Certs for this Server
	//      (or probably, move that logic to lib/go-atscfg/meta.go)
	needsCertsForClients := tc.CacheType(toData.Server.Type) == tc.CacheTypeEdge

	if needsCertsForClients {
		sslConfigs, err := GetSSLCertsAndKeyFiles(toData)
		if err != nil {
			return nil, errors.New("getting ssl key and cert config files: " + err.Error())
		}
		configs = append(configs, sslConfigs...)
	}

	return configs, nil
}

const HdrConfigFilePath = "Path"
const HdrLineComment = "Line-Comment"

// WriteConfigs writes the given configs as a RFC2046§5.1 MIME multipart/mixed message.
func WriteConfigs(configs []t3cutil.ATSConfigFile, output io.Writer) error {
	if err := json.NewEncoder(output).Encode(configs); err != nil {
		return errors.New("encoding and writing configs: " + err.Error())
	}
	return nil
}

func makeHeaderComment(serverHostName string, appVersion string, toURL string, toIPs []string, genTime time.Time) string {
	return fmt.Sprintf(
		`DO NOT EDIT - Generated for %v by %v from %v ips %v on %v`,
		serverHostName,
		appVersion,
		toURL,
		makeIPStr(toIPs),
		genTime.UTC().Format(time.RFC3339Nano),
	)
}

func makeIPStr(ips []string) string {
	return `(` + strings.Join(ips, `,`) + `)`
}

// logWarnings writes all strings in warnings to the warning log, with the context prefix.
// If warnings is empty, no log is written.
func logWarnings(context string, warnings []string) {
	for _, warn := range warnings {
		log.Warnln(context + warn)
	}
}
