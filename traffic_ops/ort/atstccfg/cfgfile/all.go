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
	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

// GetAllConfigs gets all config files for cfg.CacheHostName.
func GetAllConfigs(cfg config.TCCfg) ([]ATSConfigFile, error) {
	toData, err := GetTOData(cfg)
	if err != nil {
		return nil, errors.New("getting data from traffic ops: " + err.Error())
	}

	meta, err := GetMeta(toData)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	hasSSLMultiCertConfig := false
	configs := []ATSConfigFile{}
	for _, fi := range meta.ConfigFiles {
		if cfg.RevalOnly && fi.FileNameOnDisk != atscfg.RegexRevalidateFileName {
			continue
		}
		txt, contentType, err := GetConfigFile(toData, fi)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.APIURI + "': " + err.Error())
		}
		if fi.FileNameOnDisk == atscfg.SSLMultiCertConfigFileName {
			hasSSLMultiCertConfig = true
		}
		txt = PreprocessConfigFile(toData.Server, txt)
		configs = append(configs, ATSConfigFile{ATSConfigMetaDataConfigFile: fi, Text: txt, ContentType: contentType})
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

// WriteConfigs writes the given configs as a RFC2046§5.1 MIME multipart/mixed message.
func WriteConfigs(configs []ATSConfigFile, output io.Writer) error {
	w := multipart.NewWriter(output)

	// Create a unique boundary. Because we're using a text encoding, we need to make sure the boundary text doesn't occur in any body.
	boundary := w.Boundary()
	randSet := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	for _, cfg := range configs {
		for strings.Contains(cfg.Text, boundary) {
			boundary += string(randSet[rand.Intn(len(randSet))])
		}
	}
	if err := w.SetBoundary(boundary); err != nil {
		return errors.New("setting multipart writer boundary '" + boundary + "': " + err.Error())
	}

	io.WriteString(output, `MIME-Version: 1.0`+"\r\n"+`Content-Type: multipart/mixed; boundary="`+boundary+`"`+"\r\n\r\n")

	for _, cfg := range configs {
		hdr := map[string][]string{
			rfc.ContentType:   {cfg.ContentType},
			HdrConfigFilePath: {filepath.Join(cfg.Location, cfg.FileNameOnDisk)},
		}
		partW, err := w.CreatePart(hdr)
		if err != nil {
			return errors.New("creating multipart part for config file '" + cfg.FileNameOnDisk + "': " + err.Error())
		}
		if _, err := io.WriteString(partW, cfg.Text); err != nil {
			return errors.New("writing to multipart part for config file '" + cfg.FileNameOnDisk + "': " + err.Error())
		}
	}

	if err := w.Close(); err != nil {
		return errors.New("closing multipart writer and writing final boundary: " + err.Error())
	}
	return nil
}

var returnRegex = regexp.MustCompile(`\s*__RETURN__\s*`)

// PreprocessConfigFile does global preprocessing on the given config file cfgFile.
// This is mostly string replacements of __X__ directives. See the code for the full list of replacements.
// These things were formerly done by ORT, but need to be processed by atstccfg now, because ORT no longer has the metadata necessary.
func PreprocessConfigFile(server tc.Server, cfgFile string) string {
	if server.TCPPort != 80 && server.TCPPort != 0 {
		cfgFile = strings.Replace(cfgFile, `__SERVER_TCP_PORT__`, strconv.Itoa(server.TCPPort), -1)
	} else {
		cfgFile = strings.Replace(cfgFile, `:__SERVER_TCP_PORT__`, ``, -1)
	}
	cfgFile = strings.Replace(cfgFile, `__CACHE_IPV4__`, server.IPAddress, -1)
	cfgFile = strings.Replace(cfgFile, `__HOSTNAME__`, server.HostName, -1)
	cfgFile = strings.Replace(cfgFile, `__FULL_HOSTNAME__`, server.HostName+`.`+server.DomainName, -1)
	cfgFile = returnRegex.ReplaceAllString(cfgFile, "\n")
	return cfgFile
}
