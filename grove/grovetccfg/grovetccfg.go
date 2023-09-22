package main

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v3-client"

	"github.com/apache/trafficcontrol/v8/grove/config"
	"github.com/apache/trafficcontrol/v8/grove/remap"
	"github.com/apache/trafficcontrol/v8/grove/remapdata"
	"github.com/apache/trafficcontrol/v8/grove/web"
)

const Version = "0.2"
const UserAgent = "grove-tc-cfg/" + Version
const TrafficOpsTimeout = time.Second * 90
const DefaultCertificateDir = "/etc/grove/ssl"
const GroveConfigFile = "grove.cfg"
const GroveConfigPath = "/etc/grove/" + GroveConfigFile
const ConfigHistory = "cfg_history/"
const RemapHistory = "remap_history/"

// Exit codes are defined in the documentation, DO NOT change to iota, to avoid ambiguity.
const (
	ExitSuccess                 = 0
	ExitError                   = 1
	ExitErrorReloadingService   = 2
	ExitErrorClearingUpdateFlag = 3
)

func AvailableStatuses() map[string]struct{} {
	return map[string]struct{}{
		"reported": {},
		"online":   {},
	}
}

func GetRemapPath() (string, error) {
	cfg, err := config.LoadConfig(GroveConfigPath)
	if err != nil {
		return "", errors.New("loading Grove config file: " + err.Error())
	}
	return cfg.RemapRulesFile, nil
}

// CopyAndGzipFile reads the src file, gzips the contents, and writes the result to dst.
func CopyAndGzipFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return errors.New("opening source file: " + err.Error())
	}
	defer srcF.Close()
	dstF, err := os.Create(dst)
	if err != nil {
		return errors.New("creating destination file: " + err.Error())
	}
	defer dstF.Close()

	dstFGzip := gzip.NewWriter(dstF)

	if _, err = io.Copy(dstFGzip, srcF); err != nil {
		return errors.New("copying source to destination: " + err.Error())
	}
	if err := dstFGzip.Close(); err != nil {
		return errors.New("closing destination gzip writer: " + err.Error())
	}
	if err := dstF.Sync(); err != nil {
		return errors.New("flushing copy to destination: " + err.Error())
	}
	return nil
}

// BackupFile copies the given file to a new file in the subdirectory "/remap_history", with the name suffixed by the current timestamp. Returns nil if the given file doesn't exist (nothing to back up).
func BackupFile(path string, backupDirectory string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	fileTimeFormat := "2006-01-02T15_04_05_999999999Z07_00" // this is time.RFC3339Nano with : replaced by _
	backupDir := filepath.Join(filepath.Dir(path), backupDirectory)
	backupFile := filepath.Base(path) + "." + time.Now().Format(fileTimeFormat) + ".gz"
	backupPath := filepath.Join(backupDir, backupFile)

	os.MkdirAll(backupDir, os.ModePerm)

	if err := CopyAndGzipFile(path, backupPath); err != nil {
		return errors.New("backuping up remap file: " + err.Error())
	}
	return nil
}

func NewFilename(path string) string {
	return path + ".new"
}

func WriteNewFile(path string, bts []byte) error {
	newPath := NewFilename(path)
	f, err := os.Create(newPath)
	if err != nil {
		return errors.New("creating file: " + err.Error())
	}
	defer f.Close()
	if _, err := f.Write(bts); err != nil {
		return errors.New("writing file: " + err.Error())
	}
	if err := f.Sync(); err != nil {
		return errors.New("flushing file: " + err.Error())
	}
	return nil
}

// WriteAndBackup creates a backup of the existing file at the given path, then writes the given bytes to the path. The write is fail-safe on operating systems with atomic file rename (Linux is).
func WriteAndBackup(path string, backupDirectory string, bts []byte) error {
	if err := BackupFile(path, backupDirectory); err != nil {
		return errors.New("backing up file: " + err.Error())
	}
	if err := WriteNewFile(path, bts); err != nil {
		return errors.New("writing new file: " + err.Error())
	}
	if err := os.Rename(NewFilename(path), path); err != nil {
		return errors.New("copying new file to real location: " + err.Error())
	}
	return nil
}

// hasUpdatePending returns whether an update is pending, the revalPending status (which will be needed later in the clear update POST), and any error.
func hasUpdatePending(toc *to.Session, hostname string) (bool, bool, error) {
	response, _, err := toc.GetServersWithHdr(&url.Values{"hostName": {hostname}}, nil)
	if err != nil {
		return false, false, errors.New("getting update from Traffic Ops: " + err.Error())
	} else if len(response.Response) != 1 {
		return false, false, fmt.Errorf("Want exactly one server with hostname '%s', got %d", hostname, len(response.Response))
	}
	upd := response.Response
	return *upd[0].UpdPending, *upd[0].RevalPending, nil
}

// clearUpdatePending clears the given host's update pending flag in Traffic Ops. It takes the host to clear, and the old revalPending flag to send.
func clearUpdatePending(toc *to.Session, hostname string, revalPending bool) error {
	response, _, err := toc.GetServersWithHdr(&url.Values{"hostName": {hostname}}, nil)
	if err != nil {
		return fmt.Errorf("Failed to update reval_pending: %v", err)
	} else if len(response.Response) != 1 {
		return fmt.Errorf("Want exactly one server with hostname '%s', got '%d'", hostname, len(response.Response))
	}

	srv := response.Response
	srv[0].RevalPending = &revalPending
	srv[0].UpdPending = new(bool)
	alerts, _, err := toc.UpdateServerByIDWithHdr(*srv[0].ID, srv[0], nil)
	if err != nil {
		return fmt.Errorf("setting update pending on Traffic Ops: %v (Alerts: %+v)", err, alerts.Alerts)
	}
	return nil
}

func main() {
	toURL := flag.String("tourl", "", "The Traffic Ops URL")
	toUser := flag.String("touser", "", "The Traffic Ops username")
	toPass := flag.String("topass", "", "The Traffic Ops password")
	pretty := flag.Bool("pretty", false, "Whether to pretty-print output")
	ignoreUpdateFlag := flag.Bool("ignore-update-flag", false, "Whether to fetch and apply the config, without checking or updating the Traffic Ops Update Pending flag")
	host := flag.String("host", "", "The hostname of the server whose config to generate")
	toInsecure := flag.Bool("insecure", false, "Whether to allow invalid certificates with Traffic Ops")
	certDir := flag.String("certdir", DefaultCertificateDir, "Directory to save certificates to")
	noServiceReload := flag.Bool("no-service-reload", false, "Whether to avoid trying to reload the Grove service")
	flag.Parse()

	if host == nil || *host == "" {
		// if the host was not specified on the command line, get it from the kernel
		h, err := os.Hostname()
		if err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + err.Error() + " and 'host' is required, use '-host'")
			os.Exit(ExitError)
		}
		sl := strings.Split(h, ".")
		*host = sl[0]
	}

	useCache := false
	toc, _, err := to.LoginWithAgent(*toURL, *toUser, *toPass, *toInsecure, UserAgent, useCache, TrafficOpsTimeout)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error connecting to Traffic Ops: " + err.Error())
		os.Exit(ExitError)
	}

	revalPendingStatus := false
	if !*ignoreUpdateFlag {
		needsUpdate := false
		needsUpdate, revalPendingStatus, err = hasUpdatePending(toc, *host)
		if err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error checking Traffic Ops update pending: " + err.Error())
			os.Exit(ExitError)
		}
		if !needsUpdate {
			os.Exit(ExitSuccess) // if no error and no update necessary, return success and print nothing
		}
	}

	// TODO remove this once converted to 1.3 API
	// using the 1.2 API
	var hostServer tc.ServerV30
	var hostProfile tc.Profile
	var ok bool
	var profiles map[string]tc.Profile
	var servers map[string]tc.ServerV30

	serversArr, _, err := toc.GetServersWithHdr(nil, nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Servers: " + err.Error())
		os.Exit(ExitError)
	}
	servers = makeServersHostnameMap(serversArr.Response)

	hostServer, ok = servers[*host]
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error: host '" + *host + "' not in Servers\n")
		os.Exit(ExitError)
	}

	profilesArr, _, err := toc.GetProfilesWithHdr(nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Profiles: " + err.Error())
		os.Exit(ExitError)
	}
	profiles = makeProfileNameMap(profilesArr)

	hostProfile, ok = profiles[*hostServer.Profile]
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error: profile '" + *hostServer.Profile + "' not in Profiles\n")
		os.Exit(ExitError)
	}
	// end of API 1.2 stuff

	if hostProfile.Type == tc.GroveProfileType {
		updateRequired, cfg, err := createGroveCfg(toc, hostServer)
		if err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting config rules for '" + GroveConfigPath + "' :" + err.Error())
			os.Exit(ExitError)
		}
		if updateRequired {
			cfgBytes := []byte{}
			if *pretty {
				cfgBytes, err = json.MarshalIndent(cfg, "", "  ")
			} else {
				cfgBytes, err = json.Marshal(cfg)
			}
			if err != nil {
				fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error creating JSON Remap Rules: " + err.Error())
				os.Exit(ExitError)
			}
			if err := WriteAndBackup(GroveConfigPath, ConfigHistory, cfgBytes); err != nil {
				fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error writing new config file: " + err.Error())
				os.Exit(ExitError)
			}
		}
	} else {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Warning: the profile '" + *hostServer.Profile + "' is not a '" + tc.GroveProfileType + "', will not build a config from it.")
	}

	rules := remap.RemapRules{}
	// if *api == "1.3" {
	// 	rules, err = createRulesNewAPI(toc, *host, *certDir)
	// } else {
	rules, err = createRulesOldAPI(toc, *host, *certDir, servers) // TODO remove once 1.3 / traffic_ops_golang is deployed to production.
	// }
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error creating rules: " + err.Error())
		os.Exit(ExitError)
	}

	jsonRules, err := remap.RemapRulesToJSON(rules)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error creating JSON Remap Rules: " + err.Error())
		os.Exit(ExitError)
	}

	bts := []byte{}
	if *pretty {
		bts, err = json.MarshalIndent(jsonRules, "", "  ")
	} else {
		bts, err = json.Marshal(jsonRules)
	}

	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error marshalling rules JSON: " + err.Error())
		os.Exit(ExitError)
	}

	// TODO add app/option to print config to stdout

	remapPath, err := GetRemapPath()
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting remap config path: " + err.Error())
		os.Exit(ExitError)
	}

	if err := WriteAndBackup(remapPath, RemapHistory, bts); err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error writing new config file: " + err.Error())
		os.Exit(ExitError)
	}

	if !*noServiceReload {
		if err := exec.Command("service", "grove", "reload").Run(); err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error restarting grove service (but successfully updated config file): " + err.Error())
			os.Exit(ExitErrorReloadingService)
		}
	}

	if !*ignoreUpdateFlag {
		if err := clearUpdatePending(toc, *host, revalPendingStatus); err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error clearing update pending flag in Traffic Ops (but successfully updated config): " + err.Error())
			os.Exit(ExitErrorClearingUpdateFlag)
		}
	}

	os.Exit(ExitSuccess)
}

func createGroveCfg(toc *to.Session, server tc.ServerV30) (bool, config.Config, error) {
	var newCfg config.Config
	var currCfg config.Config
	var pluginParams = []string{}

	// load the servers current config parameters.
	if _, err := os.Stat(GroveConfigPath); err == nil {
		currCfg, err = config.LoadConfig(GroveConfigPath)
		if err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error loading current config from '" + GroveConfigPath + "' " + err.Error())
			return false, currCfg, err
		} else {
			// make sure this array is sorted for later comparison
			sort.Strings(currCfg.Plugins)
		}
	}

	serverParameters, _, err := toc.GetParametersByProfileNameWithHdr(*server.Profile, nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Parameters for host '" + *server.HostName + "' profile '" + *server.Profile + "': " + err.Error())
		return false, currCfg, err
	} else {
		// load config parameters from the servers profile
		for _, p := range serverParameters {
			if p.ConfigFile == GroveConfigFile {
				if p.Name == "plugins" {
					pluginParams = append(pluginParams, p.Value)
				} else {
					err := setConfigParameter(&newCfg, p.Name, p.Value)
					if err != nil {
						fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error setting config parameter '" + p.Name + "' :" + err.Error())
						return false, currCfg, err
					}
				}
			}
		}
		sort.Strings(pluginParams)
		newCfg.Plugins = pluginParams
	}
	// no update is required if the configs are the same
	areEqual := reflect.DeepEqual(newCfg, currCfg)
	if areEqual == true {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " There are no changes needed to '" + GroveConfigPath + "'")
		return false, currCfg, nil
	} else { // updates are required, send the new config
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Config updates are required to '" + GroveConfigPath + "'")
		return true, newCfg, nil
	}
}

func setConfigParameter(cfg *config.Config, name string, value string) error {
	var err error

	switch name {
	case "rfc_compliant":
		cfg.RFCCompliant, err = strconv.ParseBool(value)
	case "port":
		cfg.Port, err = strconv.Atoi(value)
	case "https_port":
		cfg.HTTPSPort, err = strconv.Atoi(value)
	case "cache_size_bytes":
		cfg.CacheSizeBytes, err = strconv.Atoi(value)
	case "concurrent_rule_requests":
		cfg.ConcurrentRuleRequests, err = strconv.Atoi(value)
	case "remap_rules_file":
		cfg.RemapRulesFile = value
	case "cert_file":
		cfg.CertFile = value
	case "key_file":
		cfg.KeyFile = value
	case "interface_name":
		cfg.InterfaceName = value
	case "connection_close":
		cfg.ConnectionClose, err = strconv.ParseBool(value)
	case "log_location_error":
		cfg.LogLocationError = value
	case "log_location_warning":
		cfg.LogLocationWarning = value
	case "log_location_info":
		cfg.LogLocationInfo = value
	case "log_location_debug":
		cfg.LogLocationDebug = value
	case "log_location_event":
		cfg.LogLocationEvent = value
	case "parent_request_timeout_ms":
		cfg.ReqTimeoutMS, err = strconv.Atoi(value)
	case "parent_request_keep_alive_ms":
		cfg.ReqKeepAliveMS, err = strconv.Atoi(value)
	case "parent_request_max_idle_connections":
		cfg.ReqMaxIdleConns, err = strconv.Atoi(value)
	case "parent_request_idle_connection_timeout_ms":
		cfg.ReqIdleConnTimeoutMS, err = strconv.Atoi(value)
	case "server_idle_timeout_ms":
		cfg.ServerIdleTimeoutMS, err = strconv.Atoi(value)
	case "server_write_timeout_ms":
		cfg.ServerWriteTimeoutMS, err = strconv.Atoi(value)
	case "server_read_timeout_ms":
		cfg.ServerReadTimeoutMS, err = strconv.Atoi(value)
	case "file_mem_bytes":
		cfg.FileMemBytes, err = strconv.Atoi(value)
	default:
		err = fmt.Errorf(time.Now().Format(time.RFC3339Nano) + "No such config parameter '" + name + "', parameter ignored")
	}
	return err
}

func createRulesOldAPI(toc *to.Session, host string, certDir string, servers map[string]tc.ServerV30) (remap.RemapRules, error) {
	cachegroupsArr, _, err := toc.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Cachegroups: " + err.Error())
		os.Exit(1)
	}
	cachegroups := makeCachegroupsNameMap(cachegroupsArr)

	hostServer, ok := servers[host]
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error: host '" + host + "' not in Servers\n")
		os.Exit(1)
	}

	deliveryservices, _, err := toc.GetDeliveryServicesByServerWithHdr(*hostServer.ID, nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Deliveryservices: " + err.Error())
		os.Exit(1)
	}

	deliveryserviceRegexArr, _, err := toc.GetDeliveryServiceRegexesWithHdr(nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Deliveryservice Regexes: " + err.Error())
		os.Exit(1)
	}
	deliveryserviceRegexes := makeDeliveryserviceRegexMap(deliveryserviceRegexArr)

	cdnsArr, _, err := toc.GetCDNsWithHdr(nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops CDNs: " + err.Error())
		os.Exit(1)
	}
	cdns := makeCDNMap(cdnsArr)

	serverParameters, _, err := toc.GetParametersByProfileNameWithHdr(*hostServer.Profile, nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting Traffic Ops Parameters for host '" + host + "' profile '" + *hostServer.Profile + "': " + err.Error())
		os.Exit(1)
	}

	parents, err := getParents(host, servers, cachegroups)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting '" + host + "' parents: " + err.Error())
		os.Exit(1)
	}

	sameCDN := func(s tc.ServerV30) bool {
		return s.CDNName == hostServer.CDNName
	}

	serverAvailable := func(s tc.ServerV30) bool {
		status := strings.ToLower(*s.Status)
		statuses := AvailableStatuses()
		_, ok := statuses[status]
		return ok
	}

	parents = filterParents(parents, sameCDN)
	parents = filterParents(parents, serverAvailable)

	cdnSSLKeys, _, err := toc.GetCDNSSLKeysWithHdr(*hostServer.CDNName, nil)
	if err != nil {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error getting '" + *hostServer.CDNName + "' SSL keys: " + err.Error())
		os.Exit(1)
	}
	dsCerts := makeDSCertMap(cdnSSLKeys)

	return createRulesOld(host, deliveryservices, parents, deliveryserviceRegexes, cdns, serverParameters, dsCerts, certDir)
}

// func createRulesNewAPI(toc *to.Session, host string, certDir string) (remap.RemapRules, error) {
// 	cacheCfg, err := toc.CacheConfig(host)
// 	if err != nil {
// 		fmt.Printf("Error getting Traffic Ops Cache Config: %v\n", err)
// 		os.Exit(1)
// 	}

// 	rules := []remapdata.RemapRule{}

// 	allowedIPs, err := makeAllowIP(cacheCfg.AllowIP)
// 	if err != nil {
// 		return remap.RemapRules{}, fmt.Errorf("creating allowed IPs: %v", err)
// 	}

// 	cdnSSLKeys, err := toc.CDNSSLKeys(cacheCfg.CDN)
// 	if err != nil {
// 		fmt.Printf("Error getting %v SSL keys: %v\n", cacheCfg.CDN, err)
// 		os.Exit(1)
// 	}
// 	dsCerts := makeDSCertMap(cdnSSLKeys)

// 	weight := DefaultRuleWeight
// 	retryNum := DefaultRetryNum
// 	timeout := DefaultTimeout
// 	parentSelection := DefaultRuleParentSelection

// 	for _, ds := range cacheCfg.DeliveryServices {
// 		protocol := ds.Protocol
// 		queryStringRule, err := getQueryStringRule(ds.QueryStringIgnore)
// 		if err != nil {
// 			return remap.RemapRules{}, fmt.Errorf("getting deliveryservice %v Query String Rule: %v", ds.XMLID, err)
// 		}

// 		protocolStrs := []ProtocolStr{}
// 		switch protocol {
// 		case ProtocolHTTP:
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
// 		case ProtocolHTTPS:
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
// 		case ProtocolHTTPAndHTTPS:
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
// 		case ProtocolHTTPToHTTPS:
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "https"})
// 			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
// 		}

// 		cert, hasCert := dsCerts[ds.XMLID]
// 		// DEBUG
// 		// if protocol != ProtocolHTTP {
// 		// 	if !hasCert {
// 		// 		fmt.Fprint(os.Stderr, "HTTPS delivery service: "+ds.XMLID+" has no certificate!\n")
// 		// 	} else if err := createCertificateFiles(cert, certDir); err != nil {
// 		// 		fmt.Fprint(os.Stderr, "HTTPS delivery service "+ds.XMLID+" failed to create certificate: "+err.Error()+"\n")
// 		// 	}
// 		// }

// 		dsType := strings.ToLower(ds.Type)
// 		if !strings.HasPrefix(dsType, "http") && !strings.HasPrefix(dsType, "dns") {
// 			fmt.Printf("createRules skipping deliveryservice %v - unknown type %v", ds.XMLID, ds.Type)
// 			continue
// 		}

// 		for _, protocolStr := range protocolStrs {
// 			for _, dsRegex := range ds.Regexes {
// 				rule := remapdata.RemapRule{}
// 				pattern, patternLiteralRegex := trimLiteralRegex(dsRegex)
// 				rule.Name = fmt.Sprintf("%s.%s.%s.%s", ds.XMLID, protocolStr.From, protocolStr.To, pattern)
// 				rule.From = buildFrom(protocolStr.From, pattern, patternLiteralRegex, host, dsType, cacheCfg.Domain)

// 				if protocolStr.From == "https" && hasCert {
// 					rule.CertificateFile = getCertFileName(cert, certDir)
// 					rule.CertificateKeyFile = getCertKeyFileName(cert, certDir)
// 					// fmt.Fprintf(os.Stderr, "HTTPS delivery service: "+ds.XMLID+" certificate %+v\n", cert)
// 				}

// 				for _, parent := range cacheCfg.Parents {
// 					to, proxyURLStr := buildToNew(parent, protocolStr.To, ds.OriginFQDN, dsType)
// 					proxyURL, err := url.Parse(proxyURLStr)
// 					if err != nil {
// 						return remap.RemapRules{}, fmt.Errorf("error parsing deliveryservice %v parent %v proxy_url: %v", ds.XMLID, parent.Host, proxyURLStr)
// 					}

// 					ruleTo := remapdata.RemapRuleTo{
// 						RemapRuleToBase: remapdata.RemapRuleToBase{
// 							URL:      to,
// 							Weight:   &weight,
// 							RetryNum: &retryNum,
// 						},
// 						ProxyURL:   proxyURL,
// 						RetryCodes: DefaultRetryCodes(),
// 						Timeout:    &timeout,
// 					}
// 					rule.To = append(rule.To, ruleTo)
// 					// TODO get from TO?
// 					rule.RetryNum = &retryNum
// 					rule.Timeout = &timeout
// 					rule.RetryCodes = DefaultRetryCodes()
// 					rule.QueryString = queryStringRule
// 					rule.DSCP = ds.DSCP
// 					if err != nil {
// 						return remap.RemapRules{}, err
// 					}
// 					rule.ConnectionClose = DefaultRuleConnectionClose
// 					rule.ParentSelection = &parentSelection
// 				}
// 				rules = append(rules, rule)
// 			}
// 		}
// 	}

// 	remapRules := remap.RemapRules{
// 		Rules:           rules,
// 		RetryCodes:      DefaultRetryCodes(),
// 		Timeout:         &timeout,
// 		ParentSelection: &parentSelection,
// 		Stats:           remap.RemapRulesStats{Allow: allowedIPs},
// 	}

// 	return remapRules, nil
// }

func makeProfileNameMap(profiles []tc.Profile) map[string]tc.Profile {
	m := map[string]tc.Profile{}
	for _, profile := range profiles {
		m[profile.Name] = profile
	}
	return m
}

func makeServersHostnameMap(servers []tc.ServerV30) map[string]tc.ServerV30 {
	m := map[string]tc.ServerV30{}
	for _, server := range servers {
		m[*server.HostName] = server
	}
	return m
}

func makeCachegroupsNameMap(cgs []tc.CacheGroupNullable) map[string]tc.CacheGroupNullable {
	m := map[string]tc.CacheGroupNullable{}
	for _, cg := range cgs {
		if cg.Name != nil {
			m[*cg.Name] = cg
		}
	}
	return m
}

func makeDeliveryserviceRegexMap(dsrs []tc.DeliveryServiceRegexes) map[string][]tc.DeliveryServiceRegex {
	m := map[string][]tc.DeliveryServiceRegex{}
	for _, dsr := range dsrs {
		m[dsr.DSName] = dsr.Regexes
	}
	return m
}

func makeCDNMap(cdns []tc.CDN) map[string]tc.CDN {
	m := map[string]tc.CDN{}
	for _, cdn := range cdns {
		m[cdn.Name] = cdn
	}
	return m
}

func makeDSCertMap(sslKeys []tc.CDNSSLKeys) map[string]tc.CDNSSLKeys {
	m := map[string]tc.CDNSSLKeys{}
	for _, sslkey := range sslKeys {
		m[sslkey.DeliveryService] = sslkey
	}
	return m
}

func getParents(hostname string, servers map[string]tc.ServerV30, cachegroups map[string]tc.CacheGroupNullable) ([]tc.ServerV30, error) {
	server, ok := servers[hostname]
	if !ok {
		return nil, fmt.Errorf("hostname not found in Servers")
	}

	cachegroup, ok := cachegroups[*server.Cachegroup]
	if !ok {
		return nil, fmt.Errorf("server cachegroup '%v' not found in Cachegroups", server.Cachegroup)
	}

	parents := []tc.ServerV30{}
	for _, server := range servers {
		if cachegroup.ParentName != nil && *server.Cachegroup == *cachegroup.ParentName {
			parents = append(parents, server)
		}
	}
	return parents, nil
}

func filterParents(parents []tc.ServerV30, include func(tc.ServerV30) bool) []tc.ServerV30 {
	newParents := []tc.ServerV30{}
	for _, parent := range parents {
		if include(parent) {
			newParents = append(newParents, parent)
		}
	}
	return newParents
}

const ProtocolHTTP = 0
const ProtocolHTTPS = 1
const ProtocolHTTPAndHTTPS = 2
const ProtocolHTTPToHTTPS = 3

type ProtocolStr struct {
	From string
	To   string
}

// trimLiteralRegex removes the prefix and suffix in .*\.foo\.* delivery service regexes. Traffic Ops Delivery Services have regexes of this form, which aren't really regexes, and the .*\ and \.* need stripped to construct the "to" FQDN. Returns the trimmed string, and whether it was of the form `.*\.foo\.*`
func trimLiteralRegex(s string) (string, bool) {
	prefix := `.*\.`
	suffix := `\..*`
	if strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix) {
		return s[len(prefix) : len(s)-len(suffix)], true
	}
	return s, false
}

// buildFrom builds the remap "from" URI prefix. It assumes ttype is a delivery service type HTTP or DNS, behavior is undefined for any other ttype.
func buildFrom(protocol string, pattern string, patternLiteralRegex bool, host string, dsType string, cdnDomain string) string {
	if !patternLiteralRegex {
		return protocol + "://" + pattern
	}

	if isHTTP := strings.HasPrefix(dsType, "http"); isHTTP {
		return protocol + "://" + host + "." + pattern + "." + cdnDomain
	}

	return protocol + "://" + "edge." + pattern + "." + cdnDomain
}

func dsTypeSkipsMid(ttype string) bool {
	ttype = strings.ToLower(ttype)
	if ttype == "http_no_cache" || ttype == "http_live" || ttype == "dns_live" {
		return true
	}
	if strings.Contains(ttype, "live") && !strings.Contains(ttype, "natnl") {
		return true
	}
	return false
}

// buildTo returns the to URL, and the Proxy URL (if any)
func buildTo(parentServer tc.ServerV30, protocol string, originURI string, dsType string) (string, string) {
	// TODO add port?
	to := originURI
	proxy := ""
	if !dsTypeSkipsMid(dsType) {
		proxy = "http://" + *parentServer.HostName + "." + *parentServer.DomainName + ":" + strconv.Itoa(*parentServer.TCPPort)
	}
	return to, proxy
}

// // buildToNew returns the to URL, and the Proxy URL (if any)
// func buildToNew(parent tc.CacheConfigParent, protocol string, originURI string, dsType string) (string, string) {
// 	// TODO add port?
// 	to := originURI
// 	proxy := ""
// 	if !dsTypeSkipsMid(dsType) {
// 		proxy = "http://" + parent.Host + "." + parent.Domain + ":" + strconv.FormatUint(uint64(parent.Port), 10)
// 	}
// 	return to, proxy
// }

const DeliveryServiceQueryStringCacheAndRemap = 0
const DeliveryServiceQueryStringNoCacheRemap = 1
const DeliveryServiceQueryStringNoCacheNoRemap = 2

func getQueryStringRule(dsQstringIgnore *int) (remapdata.QueryStringRule, error) {
	val := 0
	if dsQstringIgnore != nil {
		val = *dsQstringIgnore
	}
	switch val {
	case DeliveryServiceQueryStringCacheAndRemap:
		return remapdata.QueryStringRule{Remap: true, Cache: true}, nil
	case DeliveryServiceQueryStringNoCacheRemap:
		return remapdata.QueryStringRule{Remap: true, Cache: false}, nil
	case DeliveryServiceQueryStringNoCacheNoRemap:
		return remapdata.QueryStringRule{Remap: false, Cache: false}, nil
	default:
		return remapdata.QueryStringRule{}, fmt.Errorf("unknown delivery service qstringIgnore value '%v'", dsQstringIgnore)
	}
}

func DefaultRetryCodes() map[int]struct{} {
	return map[int]struct{}{}
}

const DefaultRuleWeight = 1.0
const DefaultRetryNum = 5
const DefaultTimeout = time.Millisecond * 5000
const DefaultRuleConnectionClose = false
const DefaultRuleParentSelection = remapdata.ParentSelectionTypeConsistentHash

func getAllowIP(params []tc.Parameter) ([]*net.IPNet, error) {
	ips := []string{}
	for _, param := range params {
		if (param.Name == "allow_ip" || param.Name == "allow_ip6") && param.ConfigFile == "astats.config" {
			ips = append(ips, strings.Split(param.Value, ",")...)
		}
	}
	return makeAllowIP(ips)
}

func makeAllowIP(ips []string) ([]*net.IPNet, error) {
	cidrs := make([]*net.IPNet, len(ips))
	for i, ip := range ips {
		ip = strings.TrimSpace(ip)
		if !strings.Contains(ip, "/") {
			if strings.Contains(ip, ":") {
				ip += "/128"
			} else {
				ip += "/32"
			}
		}
		_, cidrnet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, fmt.Errorf("error parsing CIDR '%s': %v", ip, err)
		}
		cidrs[i] = cidrnet
	}
	return cidrs, nil
}

func createRulesOld(
	hostname string,
	dses []tc.DeliveryServiceNullable,
	parents []tc.ServerV30,
	dsRegexes map[string][]tc.DeliveryServiceRegex,
	cdns map[string]tc.CDN,
	hostParams []tc.Parameter,
	dsCerts map[string]tc.CDNSSLKeys,
	certDir string,
) (remap.RemapRules, error) {
	rules := []remapdata.RemapRule{}
	allowedIPs, err := getAllowIP(hostParams)
	if err != nil {
		return remap.RemapRules{}, fmt.Errorf("getting allowed IPs: %v", err)
	}

	weight := DefaultRuleWeight
	retryNum := DefaultRetryNum
	timeout := DefaultTimeout
	parentSelection := DefaultRuleParentSelection

	for _, ds := range dses {
		protocol := *ds.Protocol
		queryStringRule, err := getQueryStringRule(ds.QStringIgnore)
		if err != nil {
			return remap.RemapRules{}, fmt.Errorf("getting deliveryservice %v Query String Rule: %v", ds.XMLID, err)
		}

		cdn, ok := cdns[*ds.CDNName]
		if !ok {
			return remap.RemapRules{}, fmt.Errorf("deliveryservice '%v' CDN '%v' not found", *ds.XMLID, *ds.CDNName)
		}

		protocolStrs := []ProtocolStr{}
		switch protocol {
		case ProtocolHTTP:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
		case ProtocolHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPAndHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "http"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		case ProtocolHTTPToHTTPS:
			protocolStrs = append(protocolStrs, ProtocolStr{From: "http", To: "https"})
			protocolStrs = append(protocolStrs, ProtocolStr{From: "https", To: "https"})
		}

		cert, hasCert := dsCerts[*ds.XMLID]
		if protocol != ProtocolHTTP {
			if !hasCert {
				fmt.Fprint(os.Stderr, time.Now().Format(time.RFC3339Nano)+" HTTPS delivery service: "+*ds.XMLID+" has no certificate!\n")
			} else if err := createCertificateFiles(cert, certDir); err != nil {
				fmt.Fprint(os.Stderr, time.Now().Format(time.RFC3339Nano)+" HTTPS delivery service "+*ds.XMLID+" failed to create certificate: "+err.Error()+"\n")
			}
		}

		dsType := strings.ToLower(string(*ds.Type))
		if !strings.HasPrefix(dsType, "http") && !strings.HasPrefix(dsType, "dns") {
			fmt.Printf(time.Now().Format(time.RFC3339Nano)+" createRules skipping deliveryservice %v - unknown type %v", *ds.XMLID, *ds.Type)
			continue
		}

		toClientHeaders, toOriginHeaders, err := makeModHdrs(ds.EdgeHeaderRewrite, ds.RemapText)
		if err != nil {
			return remap.RemapRules{}, errors.New("Making headers for delivery service '" + *ds.XMLID + "':" + err.Error())
		}
		dsRemap := ""
		if ds.RemapText != nil {
			dsRemap = *ds.RemapText
		}
		acl, err := makeACL(dsRemap)
		if err != nil {
			fmt.Println(time.Now().Format(time.RFC3339Nano) + " createRules skipping deliveryservice '" + *ds.XMLID + "' - unsupported ACL " + dsRemap)
			continue
		}

		orgServerFQDN := ""
		if ds.OrgServerFQDN != nil {
			orgServerFQDN = *ds.OrgServerFQDN
		}

		for _, protocolStr := range protocolStrs {
			regexes, ok := dsRegexes[*ds.XMLID]
			if !ok {
				return remap.RemapRules{}, fmt.Errorf("deliveryservice '%v' has no regexes", *ds.XMLID)
			}

			for _, dsRegex := range regexes {
				rule := remapdata.RemapRule{}
				pattern, patternLiteralRegex := trimLiteralRegex(dsRegex.Pattern)
				rule.Name = fmt.Sprintf("%s.%s.%s.%s", *ds.XMLID, protocolStr.From, protocolStr.To, pattern)
				rule.From = buildFrom(protocolStr.From, pattern, patternLiteralRegex, hostname, dsType, cdn.DomainName)

				if protocolStr.From == "https" && hasCert {
					rule.CertificateFile = getCertFileName(cert, certDir)
					rule.CertificateKeyFile = getCertKeyFileName(cert, certDir)
					// fmt.Fprintf(os.Stderr, "HTTPS delivery service: "+ds.XMLID+" certificate %+v\n", cert)
				}

				rule.PluginsShared = map[string]json.RawMessage{}
				// if the delivery service skips the mid's ie, http_no_cache, http_live, and dns_live
				// only add the url rule to the origin.
				if dsTypeSkipsMid(dsType) {
					var proxyURLStr = ""
					proxyURL, err := url.Parse(proxyURLStr)
					if err != nil {
						return remap.RemapRules{}, fmt.Errorf("error parsing deliveryservice %v proxy_url: %v", *ds.XMLID, proxyURLStr)
					}
					ruleTo := remapdata.RemapRuleTo{
						RemapRuleToBase: remapdata.RemapRuleToBase{
							URL:      orgServerFQDN,
							Weight:   &weight,
							RetryNum: &retryNum,
						},
						RetryCodes: DefaultRetryCodes(),
						ProxyURL:   proxyURL,
						Timeout:    &timeout,
					}
					rule.To = append(rule.To, ruleTo)
					rule.RetryNum = &retryNum
					rule.Timeout = &timeout
					rule.RetryCodes = DefaultRetryCodes()
					rule.QueryString = queryStringRule
					rule.DSCP = *ds.DSCP
					rule.ConnectionClose = DefaultRuleConnectionClose
					rule.Allow = acl
					rule.Plugins = map[string]interface{}{}
					rule.Plugins["modify_headers"] = toClientHeaders
					rule.Plugins["modify_parent_request_headers"] = toOriginHeaders
					remapTextJSON, err := json.Marshal(dsRemap)
					if err != nil {
						return remap.RemapRules{}, fmt.Errorf("parsing deliveryservice '%v' remap text '%v' marshalling JSON: %v", *ds.XMLID, dsRemap, err)
					}
					rule.PluginsShared[web.RemapTextKey] = remapTextJSON
				} else {
					for _, parent := range parents {
						to, proxyURLStr := buildTo(parent, protocolStr.To, orgServerFQDN, dsType)
						proxyURL, err := url.Parse(proxyURLStr)
						if err != nil {
							return remap.RemapRules{}, fmt.Errorf("error parsing deliveryservice %v parent %v proxy_url: %v", *ds.XMLID, parent.HostName, proxyURLStr)
						}

						ruleTo := remapdata.RemapRuleTo{
							RemapRuleToBase: remapdata.RemapRuleToBase{
								URL:      to,
								Weight:   &weight,
								RetryNum: &retryNum,
							},
							ProxyURL:   proxyURL,
							RetryCodes: DefaultRetryCodes(),
							Timeout:    &timeout,
						}
						rule.To = append(rule.To, ruleTo)
						// TODO get from TO?
						rule.RetryNum = &retryNum
						rule.Timeout = &timeout
						rule.RetryCodes = DefaultRetryCodes()
						rule.QueryString = queryStringRule
						rule.DSCP = *ds.DSCP
						rule.ConnectionClose = DefaultRuleConnectionClose
						rule.ParentSelection = &parentSelection
						rule.Allow = acl
						rule.Plugins = map[string]interface{}{}
						rule.Plugins["modify_headers"] = toClientHeaders
						rule.Plugins["modify_parent_request_headers"] = toOriginHeaders
						remapTextJSON, err := json.Marshal(dsRemap)
						if err != nil {
							return remap.RemapRules{}, fmt.Errorf("parsing deliveryservice '%v' remap text '%v' marshalling JSON: %v", *ds.XMLID, dsRemap, err)
						}
						rule.PluginsShared[web.RemapTextKey] = remapTextJSON
					}
				}
				rules = append(rules, rule)
			}
		}
	}

	globalPlugins := map[string]interface{}{}
	serverHeader := web.Hdr{Name: "Server", Value: "Grove/0.33"}
	setHeaders := []web.Hdr{}
	setHeaders = append(setHeaders, serverHeader)
	globalHeaders := web.ModHdrs{Set: setHeaders}
	globalPlugins["modify_response_headers_global"] = globalHeaders
	remapRules := remap.RemapRules{
		Rules:           rules,
		RetryCodes:      DefaultRetryCodes(),
		Timeout:         &timeout,
		ParentSelection: &parentSelection,
		Stats:           remapdata.RemapRulesStats{Allow: allowedIPs},
		Plugins:         globalPlugins,
	}

	return remapRules, nil
}

func getCertFileName(cert tc.CDNSSLKeys, dir string) string {
	return dir + string(os.PathSeparator) + strings.Replace(cert.Hostname, "*.", "", -1) + ".crt"
}

func getCertKeyFileName(cert tc.CDNSSLKeys, dir string) string {
	return dir + string(os.PathSeparator) + strings.Replace(cert.Hostname, "*.", "", -1) + ".key"
}

func createCertificateFiles(cert tc.CDNSSLKeys, dir string) error {
	certFileName := getCertFileName(cert, dir)
	crt, err := base64.StdEncoding.DecodeString(cert.Certificate.Crt)
	if err != nil {
		return errors.New("base64decoding certificate file " + certFileName + ": " + err.Error())
	}
	if err := ioutil.WriteFile(certFileName, crt, 0644); err != nil {
		return errors.New("writing certificate file " + certFileName + ": " + err.Error())
	}

	keyFileName := getCertKeyFileName(cert, dir)
	key, err := base64.StdEncoding.DecodeString(cert.Certificate.Key)
	if err != nil {
		return errors.New("base64decoding certificate key " + keyFileName + ": " + err.Error())
	}
	if err := ioutil.WriteFile(keyFileName, key, 0644); err != nil {
		return errors.New("writing certificate key file " + keyFileName + ": " + err.Error())
	}
	return nil
}

// makeACL is a hack to take the very ATS/TrafficControl remap_text field ACLs, and turn them into grove ACLs
// note that the astats ACL input already has CIDR notation, but the DS ACL input is IP ranges.
func makeACL(remapTxt string) ([]*net.IPNet, error) {
	allow := []*net.IPNet{}
	remapTxt = strings.Join(strings.Fields(remapTxt), " ")
	// We only have @action=allow not @action=deny, and only @scr_ip= not @in_op=
	// so only worrying about that here.
	if strings.HasPrefix(remapTxt, "@action=allow") {
		remaps := strings.Split(remapTxt, " ")
		if len(remaps) < 2 {
			return nil, errors.New("malformed remapTxt '" + remapTxt + "'")
		}
		for _, allowStr := range remaps[1:] {
			allBits := 32
			if strings.Contains(allowStr, ":") {
				allBits = 128 // not sure if v6 works, only tested with our v4 ACLs.
			}
			if strings.Contains(allowStr, "-") {
				a := strings.Split(strings.TrimPrefix(allowStr, "@src_ip="), "-")
				if len(a) < 2 {
					return nil, errors.New("malformed remapTxt '" + remapTxt + "'")
				}
				startAddrStr := a[0]
				endAddrStr := a[1]
				start := net.ParseIP(startAddrStr)
				end := net.ParseIP(endAddrStr)
				if start == nil || end == nil {
					// This error catches all unexpected (but valid) options.
					return nil, fmt.Errorf("error parsing allow string: %v, %v", allowStr, remapTxt)
				}
				maskBits := allBits
				mask := net.CIDRMask(maskBits, allBits)
				for !start.Mask(mask).Equal(end.Mask(mask)) {
					maskBits--
					mask = net.CIDRMask(maskBits, allBits)
				}
				fmt.Println(time.Now().Format(time.RFC3339Nano)+" DEBUG base: ", startAddrStr, " end:", endAddrStr, " maskBits:", maskBits) // TODO remove?
				allow = append(allow, &net.IPNet{IP: start.Mask(mask), Mask: mask})
			} else {
				addrStr := strings.Trim(allowStr, "@src_ip=")
				addr := net.ParseIP(addrStr)
				if addr == nil {
					// This error catches all unexpected (but valid) options.
					return nil, fmt.Errorf("error parsing allow string: %v, %v", allowStr, remapTxt)
				}
				mask := net.CIDRMask(allBits, allBits)
				allow = append(allow, &net.IPNet{IP: addr.Mask(mask), Mask: mask})
			}
		}
	}
	return allow, nil
}

// makeModHdrs is a pretty nasty hack to take the very ATS/TrafficControl specific config stuff from Traffic Ops and turn it into header manipulation rules for grove.
// Returns the client header modifications, the origin header modifications, and any error.
func makeModHdrs(edgeHRWin *string, remapTXTin *string) (web.ModHdrs, web.ModHdrs, error) {

	edgeHRW := ""
	remapTXT := ""
	if edgeHRWin != nil {
		edgeHRW = *edgeHRWin
	}
	if remapTXTin != nil {
		remapTXT = *remapTXTin
	}
	if edgeHRW == "" && remapTXT == "" {
		return web.ModHdrs{}, web.ModHdrs{}, nil
	}

	// Normalize the whitespaces to just a single " "
	edgeHRW = strings.Join(strings.Fields(edgeHRW), " ")
	remapTXT = strings.Join(strings.Fields(remapTXT), " ")

	toClientList := web.ModHdrs{}
	toOriginList := web.ModHdrs{}
	if strings.Contains(edgeHRW, "-header") {
		toOrigin := true
		for _, line := range strings.Split(edgeHRW, "__RETURN__") {
			line = strings.TrimSuffix(line, "[L]")
			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}
			if len(parts) < 2 {
				return web.ModHdrs{}, web.ModHdrs{}, errors.New("edge header rewrite: malformed line '" + line + "'")
			}
			switch {
			case parts[0] == "cond":
				if parts[1] == "%{SEND_RESPONSE_HDR_HOOK}" {
					toOrigin = false
				}
				if parts[1] == "%{SEND_REQUEST_HDR_HOOK}" {
					toOrigin = true
				}
			case parts[0] == "set-header" || parts[0] == "add-header": // Technically these are different
				if len(parts) < 3 {
					return web.ModHdrs{}, web.ModHdrs{}, errors.New("edge header rewrite: malformed line '" + line + "'")
				}
				hdr := web.Hdr{Name: parts[1], Value: strings.Join(parts[2:], " ")}
				if toOrigin {
					toOriginList.Set = append(toOriginList.Set, hdr)
				} else {
					toClientList.Set = append(toClientList.Set, hdr)
				}
			case parts[0] == "rm-header":
				if toOrigin {
					toOriginList.Drop = append(toOriginList.Drop, parts[1])
				} else {
					toClientList.Drop = append(toClientList.Drop, parts[1])
				}
			}
		}
	}

	return toClientList, toOriginList, nil
}
