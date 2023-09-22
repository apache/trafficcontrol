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

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// GetAllConfigs gets all config files for cfg.CacheHostName.
func GetAllConfigs(
	toData *t3cutil.ConfigData,
	cfg config.Cfg,
) ([]t3cutil.ATSConfigFile, error) {
	if toData.Server.HostName == "" {
		return nil, errors.New("server hostname is nil")
	}
	// if 0 get dataconfig.go was unable to get DS capabilities using APIv5
	// because the end point has been removed and was added to the DS struct
	// so we will get the data from DeliveryServices
	if len(toData.DSRequiredCapabilities) == 0 {
		toData.DSRequiredCapabilities = makeNewDsCaps(toData.DeliveryServices)
	}

	configFiles, warnings, err := MakeConfigFilesList(toData, cfg.Dir, cfg.ATSMajorVersion)
	logWarnings("generating config files list: ", warnings)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	genTime := time.Now()
	hdrCommentTxt := makeHeaderComment(toData.Server.HostName, cfg.AppVersion(), toData.TrafficOpsURL, toData.TrafficOpsAddresses, genTime)

	hasSSLMultiCertConfig := false
	configs := []t3cutil.ATSConfigFile{}
	for _, fi := range configFiles {
		if cfg.RevalOnly && fi.Name != atscfg.RegexRevalidateFileName {
			continue
		}
		txt, contentType, secure, lineComment, warnings, err := GetConfigFile(toData, fi, hdrCommentTxt, cfg)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.Name + "': " + err.Error())
		}
		if fi.Name == atscfg.SSLMultiCertConfigFileName {
			hasSSLMultiCertConfig = true
		}
		configs = append(configs, t3cutil.ATSConfigFile{
			Name:        fi.Name,
			Path:        fi.Path,
			Text:        txt,
			Secure:      secure,
			ContentType: contentType,
			LineComment: lineComment,
			Warnings:    warnings,
		})
	}

	if hasSSLMultiCertConfig {
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

func makeNewDsCaps(deliveryServices []atscfg.DeliveryService) map[int]map[atscfg.ServerCapability]struct{} {
	svcReqCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	for _, service := range deliveryServices {
		for _, dsCap := range service.RequiredCapabilities {
			if _, ok := svcReqCaps[*service.ID]; !ok {
				svcReqCaps[*service.ID] = map[atscfg.ServerCapability]struct{}{}
			}
			svcReqCaps[*service.ID][atscfg.ServerCapability(dsCap)] = struct{}{}
		}
	}
	return svcReqCaps
}
