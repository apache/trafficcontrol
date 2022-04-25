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

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-generate"

const ExitCodeSuccess = 0
const ExitCodeErrGeneric = 1
const ExitCodeNotFound = 104
const ExitCodeBadRequest = 100

var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")

type Cfg struct {
	ListPlugins        bool
	LogLocationErr     string
	LogLocationInfo    string
	LogLocationDebug   string
	LogLocationWarn    string
	RevalOnly          bool
	Dir                string
	UseStrategies      t3cutil.UseStrategiesFlag
	ViaRelease         bool
	SetDNSLocalBind    bool
	NoOutgoingIP       bool
	ParentComments     bool
	DefaultEnableH2    bool
	DefaultTLSVersions []atscfg.TLSVersion
	Version            string
	GitRevision        string
	ExtraCertificates  []atscfg.SSLMultiCertDotConfigCertInf
	ClientCAPath       string
	ServerCAPath       string
	ClientCertPath     string
	ClientCertKeyPath  string
	InternalHTTPS      atscfg.InternalHTTPS
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // app doesn't use the event logger.

func (cfg Cfg) AppVersion() string { return t3cutil.VersionStr(AppName, cfg.Version, cfg.GitRevision) }

// GetCfg gets the application configuration, from arguments and environment variables.
func GetCfg(appVersion string, gitRevision string) (Cfg, error) {

	version := getopt.BoolLong("version", 'V', "Print version information and exit.")
	listPlugins := getopt.BoolLong("list-plugins", 'l', "Print the list of plugins.")
	help := getopt.BoolLong("help", 'h', "Print usage information and exit")
	revalOnly := getopt.BoolLong("revalidate-only", 'y', "Whether to exclude files not named 'regex_revalidate.config'")
	dir := getopt.StringLong("dir", 'D', "", "ATS config directory, used for config files without location parameters or with relative paths. May be blank. If blank and any required config file location parameter is missing or relative, will error.")
	viaRelease := getopt.BoolLong("via-string-release", 'r', "Whether to use the Release value from the RPM package as a replacement for the ATS version specified in the build that is returned in the Via and Server headers from ATS.")
	dnsLocalBind := getopt.BoolLong("dns-local-bind", 'b', "Whether to use the server's Service Addresses to set the ATS DNS local bind address.")
	disableParentConfigComments := getopt.BoolLong("disable-parent-config-comments", 'c', "Disable adding a comments to parent.config individual lines")
	defaultEnableH2 := getopt.BoolLong("default-client-enable-h2", '2', "Whether to enable HTTP/2 on Delivery Services by default, if they have no explicit Parameter. This is irrelevant if ATS records.config is not serving H2. If omitted, H2 is disabled.")
	defaultTLSVersionsStr := getopt.StringLong("default-client-tls-versions", 'T', "", "Comma-delimited list of default TLS versions for Delivery Services with no Parameter, e.g. '--default-tls-versions=1.1,1.2,1.3'. If omitted, all versions are enabled.")
	noOutgoingIP := getopt.BoolLong("no-outgoing-ip", 'i', "Whether to not set the records.config outgoing IP to the server's addresses in Traffic Ops. Default is false.")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	silentPtr := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	const useStrategiesFlagName = "use-strategies"
	const defaultUseStrategies = t3cutil.UseStrategiesFlagFalse
	useStrategiesPtr := getopt.EnumLong(useStrategiesFlagName, 0, []string{string(t3cutil.UseStrategiesFlagTrue), string(t3cutil.UseStrategiesFlagCore), string(t3cutil.UseStrategiesFlagFalse), string(t3cutil.UseStrategiesFlagCore), ""}, "", "[true | core| false] whether to generate config using strategies.yaml instead of parent.config. If true use the parent_select plugin, if 'core' use ATS core strategies, if false use parent.config.")

	extraCertsStr := getopt.StringLong("extra-certificates", 0, "", `List of extra certificates to add to the ATS ssl_multicert.config file. Paths will be inserted into the file verbatim, and must be relative to the ATS records.config proxy.config.ssl.server.cert.path (typically etc/trafficserver/ssl). Each cert, key, and optional CA path must be comma-separated, and each set of certs must be semicolon-delimited. For example, '--extra-certificates=e2e-client.cert,e2e-client.key,e2essl-ca.cert;e2e-server.cert,e2e-server.key,e2e-ca.cert;some-other.cert,some-other.key'`)

	serverCAPath := getopt.StringLong("server-ca-path", 0, "", `the path to the Certificate Authority used to sign certificates which will be presented by parent caches to child caches. This must be relative to the ATS proxy.config.ssl.server.cert.path (e.g. etc/trafficserver/ssl/)`)
	clientCAPath := getopt.StringLong("client-ca-path", 0, "", `the path to the Certificate Authority used to sign client certificates which will be presented by child caches to parent caches. This must be relative to the ATS config directory (e.g. etc/trafficserver/), _not_ the proxy.config.ssl.server.cert.path (e.g. etc/trafficserver/ssl/)`)
	clientCertPath := getopt.StringLong("client-cert-path", 0, "", `ClientCertPath is the path to the client certificate presented by child caches to parent caches. Relative to records.config proxy.config.ssl.client.cert.path`)
	clientCertKeyPath := getopt.StringLong("client-cert-key-path", 0, "", `ClientCertKeyPath is the path to the key for the client certificate presented by child caches to parent caches. Relative to records.config proxy.config.ssl.client.private_key.path`)

	const internalHTTPSFlagName = "internal-https"
	const defaultInternalHTTPS = atscfg.InternalHTTPSNo
	internalHTTPSPtr := getopt.EnumLong(internalHTTPSFlagName, 0, []string{string(atscfg.InternalHTTPSNo), string(atscfg.InternalHTTPSYes), string(atscfg.InternalHTTPSNoChild), ""}, "", "[no | yes | no-child] Whether to use HTTPS for internal cache communication. The no-child flag will generate config for both http and https on parents, and http on children; this makes it possible to upgrade without downtime: generate 'no-child' on all caches first, then 'yes'. Default is 'no'.")

	getopt.Parse()

	if *version {
		cfg := &Cfg{Version: appVersion, GitRevision: gitRevision}
		fmt.Println(cfg.AppVersion())
		os.Exit(0)
	} else if *help {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	} else if *listPlugins {
		return Cfg{ListPlugins: true}, nil
	}

	logLocationError := log.LogLocationStderr
	logLocationWarn := log.LogLocationNull
	logLocationInfo := log.LogLocationNull
	logLocationDebug := log.LogLocationNull
	if *silentPtr {
		logLocationError = log.LogLocationNull
	} else {
		if *verbosePtr >= 1 {
			logLocationWarn = log.LogLocationStderr
		}
		if *verbosePtr >= 2 {
			logLocationInfo = log.LogLocationStderr
			logLocationDebug = log.LogLocationStderr // t3c only has 3 verbosity options: none (-s), error (default or --verbose=0), warning (-v), and info (-vv). Any code calling log.Debug is treated as Info.
		}
	}

	if *verbosePtr > 2 {
		return Cfg{}, errors.New("Too many verbose options. The maximum log verbosity level is 2 (-vv or --verbose=2) for errors (0), warnings (1), and info (2)")
	}

	if *internalHTTPSPtr == "" {
		*internalHTTPSPtr = defaultInternalHTTPS.String()
	}

	defaultTLSVersions := atscfg.DefaultDefaultTLSVersions

	*defaultTLSVersionsStr = strings.TrimSpace(*defaultTLSVersionsStr)
	if len(*defaultTLSVersionsStr) > 0 {
		defaultTLSVersionsStrs := strings.Split(*defaultTLSVersionsStr, ",")

		defaultTLSVersions = []atscfg.TLSVersion{}
		for _, tlsVersionStr := range defaultTLSVersionsStrs {
			tlsVersion := atscfg.StringToTLSVersion(tlsVersionStr)
			if tlsVersion == atscfg.TLSVersionInvalid {
				return Cfg{}, errors.New("unknown TLS Version '" + tlsVersionStr + "' in '" + *defaultTLSVersionsStr + "'")
			}
			defaultTLSVersions = append(defaultTLSVersions, tlsVersion)
		}
	}

	if !getopt.IsSet(useStrategiesFlagName) {
		*useStrategiesPtr = defaultUseStrategies.String()
	}

	warnMsgs := []string{}

	extraCertsArr := strings.Split(strings.TrimSpace(*extraCertsStr), `;`)
	extraCerts := []atscfg.SSLMultiCertDotConfigCertInf{}
	for _, certEntry := range extraCertsArr {
		certEntry = strings.TrimSpace(certEntry)
		if certEntry == "" {
			continue // skip empty fields
		}
		certEntryArr := strings.Split(certEntry, `,`)
		if len(certEntryArr) < 2 {
			warnMsgs = append(warnMsgs, "--extra-certificates was malformed, entry '"+certEntry+"' was not a cert-key pair, skipping!")
			continue
		}
		cert := atscfg.SSLMultiCertDotConfigCertInf{CertPath: certEntryArr[0], KeyPath: certEntryArr[1]}
		if len(certEntryArr) > 2 {
			cert.CAPath = certEntryArr[2]
		}
		if len(certEntryArr) > 3 {
			warnMsgs = append(warnMsgs, "--extra-certificates was malformed, entry '"+certEntry+"' had more than 3 entries (cert, key, ca), ignoring entries after the 3rd!")
		}
		extraCerts = append(extraCerts, cert)
	}

	cfg := Cfg{
		LogLocationErr:     logLocationError,
		LogLocationWarn:    logLocationWarn,
		LogLocationInfo:    logLocationInfo,
		LogLocationDebug:   logLocationDebug,
		ListPlugins:        *listPlugins,
		RevalOnly:          *revalOnly,
		Dir:                *dir,
		ViaRelease:         *viaRelease,
		SetDNSLocalBind:    *dnsLocalBind,
		NoOutgoingIP:       *noOutgoingIP,
		ParentComments:     !(*disableParentConfigComments),
		DefaultEnableH2:    *defaultEnableH2,
		DefaultTLSVersions: defaultTLSVersions,
		Version:            appVersion,
		GitRevision:        gitRevision,
		UseStrategies:      t3cutil.UseStrategiesFlag(*useStrategiesPtr),
		ExtraCertificates:  extraCerts,
		ClientCAPath:       *clientCAPath,
		ServerCAPath:       *serverCAPath,
		ClientCertPath:     *clientCertPath,
		ClientCertKeyPath:  *clientCertKeyPath,
		InternalHTTPS:      atscfg.InternalHTTPS(*internalHTTPSPtr),
	}
	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}
	for _, warnMsg := range warnMsgs {
		log.Warnln(warnMsg)
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

// TOData is the Traffic Ops data needed to generate configs.
// See each field for details on the data required.
// - If a field says 'must', the creation of TOData is guaranteed to do so, and users of the struct may rely on that.
// - If it says 'may', the creation may or may not do so, and therefore users of the struct must filter if they
//   require the potential fields to be omitted to generate correctly.
type TOData struct {
	// Servers must be all the servers from Traffic Ops. May include servers not on the current cdn.
	Servers []atscfg.Server

	// CacheGroups must be all cachegroups in Traffic Ops with Servers on the current server's cdn. May also include CacheGroups without servers on the current cdn.
	CacheGroups []tc.CacheGroupNullable

	// GlobalParams must be all Parameters in Traffic Ops on the tc.GlobalProfileName Profile. Must not include other parameters.
	GlobalParams []tc.Parameter

	// ServerParams must be all Parameters on the Profile of the current server. Must not include other Parameters.
	ServerParams []tc.Parameter

	// RemapConfigParams must be all Parameters with the ConfigFile "remap.config". Also includes cachekey.config parameters
	RemapConfigParams []tc.Parameter

	// ParentConfigParams must be all Parameters with the ConfigFile "parent.config.
	ParentConfigParams []tc.Parameter

	// DeliveryServices must include all Delivery Services on the current server's cdn, including those not assigned to the server. Must not include delivery services on other cdns.
	DeliveryServices []atscfg.DeliveryService

	// DeliveryServiceServers must include all delivery service servers in Traffic Ops for all delivery services on the current cdn, including those not assigned to the current server.
	DeliveryServiceServers []tc.DeliveryServiceServer

	// Server must be the server we're fetching configs from
	Server *atscfg.Server

	// Jobs must be all Jobs on the server's CDN. May include jobs on other CDNs.
	Jobs []tc.Job

	// CDN must be the CDN of the server.
	CDN *tc.CDN

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

	// Topologies must be all the topologies for the server's cdn.
	// May incude topologies of other cdns.
	Topologies []tc.Topology
}
