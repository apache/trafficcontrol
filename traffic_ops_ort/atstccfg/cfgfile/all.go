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
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

// GetAllConfigs gets all config files for cfg.CacheHostName.
func GetAllConfigs(toData *config.TOData, revalOnly bool, dir string) ([]config.ATSConfigFile, error) {
	meta, err := GetMeta(toData, dir)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	hasSSLMultiCertConfig := false
	configs := []config.ATSConfigFile{}
	for _, fi := range meta.ConfigFiles {
		if revalOnly && fi.FileNameOnDisk != atscfg.RegexRevalidateFileName {
			continue
		}
		txt, contentType, lineComment, err := GetConfigFile(toData, fi)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.FileNameOnDisk + "': " + err.Error())
		}
		if fi.FileNameOnDisk == atscfg.SSLMultiCertConfigFileName {
			hasSSLMultiCertConfig = true
		}
		txt = PreprocessConfigFile(toData.Server, txt)
		configs = append(configs, config.ATSConfigFile{ATSConfigMetaDataConfigFile: fi, Text: txt, ContentType: contentType, LineComment: lineComment})
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
func WriteConfigs(configs []config.ATSConfigFile, output io.Writer) error {
	sort.Sort(ATSConfigFiles(configs))
	w := multipart.NewWriter(output)

	// Create a unique boundary. Because we're using a text encoding, we need to make sure the boundary text doesn't occur in any body.
	// Always start with the same random UUID, so generating twice diffs the same (except in the unlikely chance this string is in a config somewhere).
	boundary := `dc5p7zOLNkyTzdcZSme6tg` // random UUID
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
			HdrLineComment:    {cfg.LineComment},
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

// ATSConfigFiles implements sort.Interface to sort by path.
type ATSConfigFiles []config.ATSConfigFile

func (p ATSConfigFiles) Len() int      { return len(p) }
func (p ATSConfigFiles) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ATSConfigFiles) Less(i, j int) bool {
	if p[i].Location != p[j].Location {
		return p[i].Location < p[j].Location
	}
	return p[i].FileNameOnDisk < p[j].FileNameOnDisk
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
