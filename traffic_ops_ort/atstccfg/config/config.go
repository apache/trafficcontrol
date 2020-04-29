package config

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
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/toreq"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/toreqnew"

	flag "github.com/ogier/pflag"
)

const AppName = "atstccfg"
const Version = "0.2"
const UserAgent = AppName + "/" + Version

const ExitCodeSuccess = 0
const ExitCodeErrGeneric = 1
const ExitCodeNotFound = 104
const ExitCodeBadRequest = 100

var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")

type Cfg struct {
	CacheHostName   string
	DisableProxy    bool
	GetData         string
	ListPlugins     bool
	LogLocationErr  string
	LogLocationInfo string
	LogLocationWarn string
	NumRetries      int
	RevalOnly       bool
	SetQueueStatus  string
	SetRevalStatus  string
	TOInsecure      bool
	TOPass          string
	TOTimeout       time.Duration
	TOURL           *url.URL
	TOUser          string
}

type TCCfg struct {
	Cfg
	TOClient    *toreq.TOClient
	TOClientNew *toreqnew.TOClient
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // atstccfg doesn't use the debug logger, use Info instead.
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // atstccfg doesn't use the event logger.

// GetCfg gets the application configuration, from arguments and environment variables.
func GetCfg() (Cfg, error) {
	toURLPtr := flag.StringP("traffic-ops-url", "u", "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL.")
	toUserPtr := flag.StringP("traffic-ops-user", "U", "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER.")
	toPassPtr := flag.StringP("traffic-ops-password", "P", "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS.")
	numRetriesPtr := flag.IntP("num-retries", "r", 5, "The number of times to retry getting a file if it fails.")
	logLocationErrPtr := flag.StringP("log-location-error", "e", "stderr", "Where to log errors. May be a file path, stdout, stderr, or null.")
	logLocationWarnPtr := flag.StringP("log-location-warning", "w", "stderr", "Where to log warnings. May be a file path, stdout, stderr, or null.")
	logLocationInfoPtr := flag.StringP("log-location-info", "i", "stderr", "Where to log information messages. May be a file path, stdout, stderr, or null.")
	toInsecurePtr := flag.BoolP("traffic-ops-insecure", "s", false, "Whether to ignore HTTPS certificate errors from Traffic Ops. It is HIGHLY RECOMMENDED to never use this in a production environment, but only for debugging.")
	toTimeoutMSPtr := flag.IntP("traffic-ops-timeout-milliseconds", "t", 30000, "Timeout in seconds for Traffic Ops requests.")
	versionPtr := flag.BoolP("version", "v", false, "Print version information and exit.")
	listPluginsPtr := flag.BoolP("list-plugins", "l", false, "Print the list of plugins.")
	helpPtr := flag.BoolP("help", "h", false, "Print usage information and exit")
	cacheHostNamePtr := flag.StringP("cache-host-name", "n", "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	getDataPtr := flag.StringP("get-data", "d", "", "non-config-file Traffic Ops Data to get. Valid values are update-status, packages, chkconfig, system-info, and statuses")
	setQueueStatusPtr := flag.StringP("set-queue-status", "q", "", "POSTs to Traffic Ops setting the queue status of the server. Must be 'true' or 'false'. Requires --set-reval-status also be set")
	setRevalStatusPtr := flag.StringP("set-reval-status", "a", "", "POSTs to Traffic Ops setting the revalidate status of the server. Must be 'true' or 'false'. Requires --set-queue-status also be set")
	revalOnlyPtr := flag.BoolP("revalidate-only", "y", false, "Whether to exclude files not named 'regex_revalidate.config'")
	disableProxyPtr := flag.BoolP("traffic-ops-disable-proxy", "p", false, "Whether to not use the Traffic Ops proxy specified in the GLOBAL Parameter tm.rev_proxy.url")

	flag.Parse()

	if *versionPtr {
		fmt.Println(AppName + " v" + Version)
		os.Exit(0)
	} else if *helpPtr {
		flag.PrintDefaults()
		os.Exit(0)
	} else if *listPluginsPtr {
		return Cfg{ListPlugins: true}, nil
	}

	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	numRetries := *numRetriesPtr
	logLocationErr := *logLocationErrPtr
	logLocationWarn := *logLocationWarnPtr
	logLocationInfo := *logLocationInfoPtr
	toInsecure := *toInsecurePtr
	toTimeout := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	listPlugins := *listPluginsPtr
	cacheHostName := *cacheHostNamePtr
	getData := *getDataPtr
	setQueueStatus := *setQueueStatusPtr
	setRevalStatus := *setRevalStatusPtr
	revalOnly := *revalOnlyPtr
	disableProxy := *disableProxyPtr

	urlSourceStr := "argument" // for error messages
	if toURL == "" {
		urlSourceStr = "environment variable"
		toURL = os.Getenv("TO_URL")
	}
	if toUser == "" {
		toUser = os.Getenv("TO_USER")
	}
	if toPass == "" {
		toPass = os.Getenv("TO_PASS")
	}

	usageStr := "Usage: ./" + AppName + " --traffic-ops-url=myurl --traffic-ops-user=myuser --traffic-ops-password=mypass --cache-host-name=my-cache"
	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-url or TO_URL environment variable. " + usageStr)
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-user or TO_USER environment variable. " + usageStr)
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-password or TO_PASS environment variable. " + usageStr)
	}
	if strings.TrimSpace(cacheHostName) == "" {
		return Cfg{}, errors.New("Missing required argument --cache-host-name. " + usageStr)
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err := ValidateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	cfg := Cfg{
		LogLocationErr:  logLocationErr,
		LogLocationWarn: logLocationWarn,
		LogLocationInfo: logLocationInfo,
		NumRetries:      numRetries,
		TOInsecure:      toInsecure,
		TOPass:          toPass,
		TOTimeout:       toTimeout,
		TOURL:           toURLParsed,
		TOUser:          toUser,
		ListPlugins:     listPlugins,
		CacheHostName:   cacheHostName,
		GetData:         getData,
		SetRevalStatus:  setRevalStatus,
		SetQueueStatus:  setQueueStatus,
		RevalOnly:       revalOnly,
		DisableProxy:    disableProxy,
	}
	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}
	return cfg, nil
}

func ValidateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}

type ATSConfigFile struct {
	tc.ATSConfigMetaDataConfigFile
	Text        string
	ContentType string
	LineComment string
}

// ATSConfigFiles implements sort.Interface and sorts by the Location and then FileNameOnDisk, i.e. the full file path.
type ATSConfigFiles []ATSConfigFile

func (fs ATSConfigFiles) Len() int { return len(fs) }
func (fs ATSConfigFiles) Less(i, j int) bool {
	if fs[i].Location != fs[j].Location {
		return fs[i].Location < fs[j].Location
	}
	return fs[i].FileNameOnDisk < fs[j].FileNameOnDisk
}
func (fs ATSConfigFiles) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

// TOData is the Traffic Ops data needed to generate configs.
// See each field for details on the data required.
// - If a field says 'must', the creation of TOData is guaranteed to do so, and users of the struct may rely on that.
// - If it says 'may', the creation may or may not do so, and therefore users of the struct must filter if they
//   require the potential fields to be omitted to generate correctly.
type TOData struct {
	// Servers must be all the servers from Traffic Ops. May include servers not on the current cdn.
	Servers []tc.Server

	// CacheGroups must be all cachegroups in Traffic Ops with Servers on the current server's cdn. May also include CacheGroups without servers on the current cdn.
	CacheGroups []tc.CacheGroupNullable

	// GlobalParams must be all Parameters in Traffic Ops on the tc.GlobalProfileName Profile. Must not include other parameters.
	GlobalParams []tc.Parameter

	// ScopeParams must be all Parameters in Traffic Ops with the name "scope". Must not include other Parameters.
	ScopeParams []tc.Parameter

	// ServerParams must be all Parameters on the Profile of the current server. Must not include other Parameters.
	ServerParams []tc.Parameter

	// CacheKeyParams must be all Parameters with the ConfigFile atscfg.CacheKeyParameterConfigFile.
	CacheKeyParams []tc.Parameter

	// ParentConfigParams must be all Parameters with the ConfigFile "parent.config.
	ParentConfigParams []tc.Parameter

	// DeliveryServices must include all Delivery Services on the current server's cdn, including those not assigned to the server. Must not include delivery services on other cdns.
	DeliveryServices []tc.DeliveryServiceNullable

	// DeliveryServiceServers must include all delivery service servers in Traffic Ops for all delivery services on the current cdn, including those not assigned to the current server.
	DeliveryServiceServers []tc.DeliveryServiceServer

	// Server must be the server we're fetching configs from
	Server tc.Server

	// TOToolName must be the Parameter named 'tm.toolname' on the tc.GlobalConfigFileName Profile.
	TOToolName string

	// TOToolName must be the Parameter named 'tm.url' on the tc.GlobalConfigFileName Profile.
	TOURL string

	// Jobs must be all Jobs on the server's CDN. May include jobs on other CDNs.
	Jobs []tc.Job

	// CDN must be the CDN of the server.
	CDN tc.CDN

	// DeliveryServiceRegexes must be all regexes on all delivery services on this server's cdn.
	DeliveryServiceRegexes []tc.DeliveryServiceRegexes

	// Profile must be the Profile of the server being requested.
	Profile tc.Profile

	// URISigningKeys must be a map of every delivery service which is URI Signed, to its keys.
	URISigningKeys map[tc.DeliveryServiceName][]byte

	// URLSigKeys must be a map of every delivery service which uses URL Sig, to its keys.
	URLSigKeys map[tc.DeliveryServiceName]tc.URLSigKeys

	// ServerCapabilities must be a map of all server IDs on this server's CDN, to a set of their capabilities. May also include servers from other cdns.
	ServerCapabilities map[int]map[atscfg.ServerCapability]struct{}

	// DSRequiredCapabilities must be a map of all delivery service IDs on this server's CDN, to a set of their required capabilities. Delivery Services with no required capabilities may not have an entry in the map.
	DSRequiredCapabilities map[int]map[atscfg.ServerCapability]struct{}

	// SSLKeys must be all the ssl keys for the server's cdn.
	SSLKeys []tc.CDNSSLKeys
}
