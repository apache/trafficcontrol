package tc

import (
	"strings"
)

type ATSConfigMetaData struct {
	Info        ATSConfigMetaDataInfo         `json:"info"`
	ConfigFiles []ATSConfigMetaDataConfigFile `json:"configFiles"`
}

type ATSConfigMetaDataInfo struct {
	CDNID             int    `json:"cdnId"`
	CDNName           string `json:"cdnName"`
	ServerID          int    `json:"serverId"`
	ServerIPv4        string `json:"serverIpv4"`
	ServerName        string `json:"serverName"`
	ServerPort        int    `json:"serverTcpPort"`
	ProfileID         int    `json:"profileId"`
	ProfileName       string `json:"profileName"`
	TOReverseProxyURL string `json:"toRevProxyUrl"`
	TOURL             string `json:"toUrl"`
}

type ATSConfigMetaDataConfigFileScope string

const ATSConfigMetaDataConfigFileScopeProfiles = ATSConfigMetaDataConfigFileScope("profiles")
const ATSConfigMetaDataConfigFileScopeServers = ATSConfigMetaDataConfigFileScope("servers")
const ATSConfigMetaDataConfigFileScopeCDNs = ATSConfigMetaDataConfigFileScope("cdns")
const ATSConfigMetaDataConfigFileScopeInvalid = ATSConfigMetaDataConfigFileScope("")

type ATSConfigMetaDataConfigFile struct {
	FileNameOnDisk string `json:"fnameOnDisk"`
	Location       string `json:"location"`
	APIURI         string `json:"apiUri, omitempty"`
	URL            string `json:"url, omitempty"`
	Scope          string `json:"scope"`
}

func (t ATSConfigMetaDataConfigFileScope) String() string {
	switch t {
	case ATSConfigMetaDataConfigFileScopeProfiles:
		fallthrough
	case ATSConfigMetaDataConfigFileScopeServers:
		fallthrough
	case ATSConfigMetaDataConfigFileScopeCDNs:
		return string(t)
	default:
		return "invalid"
	}
}

func ATSConfigMetaDataConfigFileScopeFromString(s string) ATSConfigMetaDataConfigFileScope {
	s = strings.ToLower(s)
	switch s {
	case "profiles":
		return ATSConfigMetaDataConfigFileScopeProfiles
	case "servers":
		return ATSConfigMetaDataConfigFileScopeServers
	case "cdns":
		return ATSConfigMetaDataConfigFileScopeCDNs
	default:
		return ATSConfigMetaDataConfigFileScopeInvalid
	}
}
