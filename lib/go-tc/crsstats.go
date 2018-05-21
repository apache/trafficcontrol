package tc

type CRSStats struct {
	App   CRSStatsApp   `json:"app"`
	Stats CRSStatsStats `json:"stats"`
}

type CRSStatsApp struct {
	BuildTimestamp string `json:"buildTimestamp"`
	Name           string `json:"name"`
	DeployDir      string `json:"deploy-dir"`
	GitRevision    string `json:"git-revision"`
	Version        string `json:"version"`
}

type CRSStatsStats struct {
	DNSMap           map[string]CRSStatsStat
	HTTPMap          map[string]CRSStatsStat
	TotalDNSCount    uint64                `json:"totalDnsCount"`
	TotalHTTPCount   uint64                `json:"totalHttpCount"`
	TotalDSMissCount uint64                `json:"totalDsMissCount"`
	AppStartTime     uint64                `json:"appStartTime"`
	AverageDnsTime   uint64                `json:"averageDnsTime"`
	AverageHttpTime  uint64                `json:"averageHttpTime"`
	UpdateTracker    CRSStatsUpdateTracker `json:"updateTracker"`
}

type CRSStatsStat struct {
	CZCount                uint64 `json:"czCount"`
	GeoCount               uint64 `json:"geoCount"`
	DeepCZCount            uint64 `json:"deepCzCount"`
	MissCount              uint64 `json:"missCount"`
	DSRCount               uint64 `json:"dsrCount"`
	ErrCount               uint64 `json:"errCount"`
	StaticRouteCount       uint64 `json:"staticRouteCount"`
	FedCount               uint64 `json:"fedCount"`
	RegionalDeniedCount    uint64 `json:"regionalDeniedCount"`
	RegionalAlternateCount uint64 `json:"regionalAlternateCount"`
}

type CRSStatsUpdateTracker struct {
	LastHttpsCertificatesCheck           uint64 `json:"lastHttpsCertificatesCheck"`
	LastGeolocationDatabaseUpdaterUpdate uint64 `json:"lastGeolocationDatabaseUpdaterUpdate"`
	LastCacheStateCheck                  uint64 `json:"lastCacheStateCheck"`
	LastCacheStateChange                 uint64 `json:"lastCacheStateChange"`
	LastNetworkUpdaterUpdate             uint64 `json:"lastNetworkUpdaterUpdate"`
	LastHTTPSCertificatesUpdate          uint64 `json:"lastHttpsCertificatesUpdate"`
	LastConfigCheck                      uint64 `json:"lastConfigCheck"`
	LastConfigChange                     uint64 `json:"lastConfigChange"`
	LastHTTPSCertificatesFetchFail       uint64 `json:"lastHttpsCertificatesFetchFail"`
	LastNetworkUpdaterCheck              uint64 `json:"lastNetworkUpdaterCheck"`
	NewDNSSECKeysFound                   uint64 `json:"newDnsSecKeysFound"`
	LastGeolocationDatabaseUpdaterCheck  uint64 `json:"lastGeolocationDatabaseUpdaterCheck"`
	LastHTTPSCertificatesFetchSuccess    uint64 `json:"lastHttpsCertificatesFetchSuccess"`
	LastSteeringWatcherCheck             uint64 `json:"lastSteeringWatcherCheck"`
	LastDNSSECKeysCheck                  uint64 `json:"lastDnsSecKeysCheck"`
	LastFederationsWatcherCheck          uint64 `json:"lastFederationsWatcherCheck"`
	LastHTTPSCertificatesFetchAttempt    uint64 `json:"lastHttpsCertificatesFetchAttempt"`
}
