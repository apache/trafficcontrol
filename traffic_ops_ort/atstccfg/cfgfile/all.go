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
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

// GetAllConfigs gets all config files for cfg.CacheHostName.
func GetAllConfigs(
	toData *config.TOData,
	appVersion string,
	toIPs []net.Addr,
	cfg config.TCCfg,
) ([]config.ATSConfigFile, error) {
	if toData.Server.HostName == nil {
		return nil, errors.New("server hostname is nil")
	}

	configFiles, warnings, err := MakeConfigFilesList(toData, cfg.Dir)
	logWarnings("generating config files list: ", warnings)
	if err != nil {
		return nil, errors.New("creating meta: " + err.Error())
	}

	genTime := time.Now()
	hdrCommentTxt := makeHeaderComment(*toData.Server.HostName, appVersion, cfg.TOClient.C.URL, toIPs, genTime)

	hasSSLMultiCertConfig := false
	configs := []config.ATSConfigFile{}
	for _, fi := range configFiles {
		if cfg.RevalOnly && fi.Name != atscfg.RegexRevalidateFileName {
			continue
		}
		txt, contentType, lineComment, err := GetConfigFile(toData, fi, hdrCommentTxt, cfg)
		if err != nil {
			return nil, errors.New("getting config file '" + fi.Name + "': " + err.Error())
		}
		if fi.Name == atscfg.SSLMultiCertConfigFileName {
			hasSSLMultiCertConfig = true
		}
		txt = PreprocessConfigFile(toData.Server, txt)
		configs = append(configs, config.ATSConfigFile{Name: fi.Name, Path: fi.Path, Text: txt, ContentType: contentType, LineComment: lineComment})
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
			HdrConfigFilePath: {filepath.Join(cfg.Path, cfg.Name)},
		}
		partW, err := w.CreatePart(hdr)
		if err != nil {
			return errors.New("creating multipart part for config file '" + cfg.Name + "': " + err.Error())
		}
		if _, err := io.WriteString(partW, cfg.Text); err != nil {
			return errors.New("writing to multipart part for config file '" + cfg.Name + "': " + err.Error())
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
	if p[i].Path != p[j].Path {
		return p[i].Path < p[j].Path
	}
	return p[i].Name < p[j].Name
}

var returnRegex = regexp.MustCompile(`\s*__RETURN__\s*`)

// PreprocessConfigFile does global preprocessing on the given config file cfgFile.
// This is mostly string replacements of __X__ directives. See the code for the full list of replacements.
// These things were formerly done by ORT, but need to be processed by atstccfg now, because ORT no longer has the metadata necessary.
func PreprocessConfigFile(server *atscfg.Server, cfgFile string) string {
	if server.TCPPort != nil && *server.TCPPort != 80 && *server.TCPPort != 0 {
		cfgFile = strings.Replace(cfgFile, `__SERVER_TCP_PORT__`, strconv.Itoa(*server.TCPPort), -1)
	} else {
		cfgFile = strings.Replace(cfgFile, `:__SERVER_TCP_PORT__`, ``, -1)
	}

	ipAddr := ""
	for _, iFace := range server.Interfaces {
		for _, addr := range iFace.IPAddresses {
			if !addr.ServiceAddress {
				continue
			}
			addrStr := addr.Address
			ip := net.ParseIP(addrStr)
			if ip == nil {
				err := error(nil)
				ip, _, err = net.ParseCIDR(addrStr)
				if err != nil {
					ip = nil // don't bother with the error, just skip
				}
			}
			if ip == nil || ip.To4() == nil {
				continue
			}
			ipAddr = addrStr
			break
		}
	}
	if ipAddr != "" {
		cfgFile = strings.Replace(cfgFile, `__CACHE_IPV4__`, ipAddr, -1)
	} else {
		log.Errorln("Preprocessing: this server had a missing or malformed IPv4 Service Interface, cannot replace __CACHE_IPV4__ directives!")
	}

	if server.HostName == nil || *server.HostName == "" {
		log.Errorln("Preprocessing: this server missing HostName, cannot replace __HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__HOSTNAME__`, *server.HostName, -1)
	}
	if server.HostName == nil || *server.HostName == "" || server.DomainName == nil || *server.DomainName == "" {
		log.Errorln("Preprocessing: this server missing HostName or DomainName, cannot replace __FULL_HOSTNAME__ directives!")
	} else {
		cfgFile = strings.Replace(cfgFile, `__FULL_HOSTNAME__`, *server.HostName+`.`+*server.DomainName, -1)
	}
	cfgFile = returnRegex.ReplaceAllString(cfgFile, "\n")
	return cfgFile
}

func makeHeaderComment(serverHostName string, appVersion string, toURL string, toIPs []net.Addr, genTime time.Time) string {
	return fmt.Sprintf(
		`DO NOT EDIT - Generated for %v by %v from %v ips %v on %v`,
		serverHostName,
		appVersion,
		toURL,
		makeIPStr(toIPs),
		genTime.UTC().Format(time.RFC3339Nano),
	)
}

func makeIPStr(ips []net.Addr) string {
	ipStrM := map[string]struct{}{} // use a map to de-duplicate
	for _, ip := range ips {
		if ip == nil {
			continue // shouldn't happen, but not really an error if it does
		}
		ipStrM[ip.String()] = struct{}{}
	}
	ipStrArr := []string{}
	for ipStr, _ := range ipStrM {
		ipStrArr = append(ipStrArr, ipStr)
	}
	return `(` + strings.Join(ipStrArr, `,`) + `)`
}
