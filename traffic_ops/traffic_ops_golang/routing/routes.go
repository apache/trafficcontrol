package routing

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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apicapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apiriak"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apitenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/asn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats/atscdn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats/atsprofile"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats/atsserver"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroupparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachesstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdnfederation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/coordinate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/crconfig"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbdump"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/consistenthash"
	dsrequest "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request/comment"
	dsserver "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/servers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservicerequests"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservicesregexes"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/division"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/federation_resolvers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/federations"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/hwinfo"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/invalidationjobs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/iso"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/login"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/logs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/origin"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/physlocation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ping"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profile"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profileparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/region"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/role"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercheck"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steering"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steeringtargets"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/systeminfo"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/toextension"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/types"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/urisigning"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/user"

	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
)

// Authenticated ...
const Authenticated = true

// NoAuth ...
const NoAuth = false

// perlBypass means the route is ABLE to bypass to TO-Perl for handling. Generally, this should only be used for Routes
// that have been rewritten from Perl to Go but have not been "vetted" in production environment yet. Once a Route has
// been "vetted", it should be changed to use noPerlBypass instead.
const perlBypass = true

// noPerlBypass means the route is UNABLE to bypass to TO-Perl for handling. Generally, this should only be used for
// NEW Routes that have no Perl-equivalent, or Go-rewritten routes that have been "vetted" in production.
const noPerlBypass = false

func handlerToFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// GetRouteIDMap takes a []int Route IDs and converts it into a map for fast lookup.
func GetRouteIDMap(IDs []int) map[int]struct{} {
	m := make(map[int]struct{}, len(IDs))
	for _, id := range IDs {
		m[id] = struct{}{}
	}
	return m
}

// Routes returns the API routes, raw non-API root level routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, []RawRoute, http.Handler, error) {
	proxyHandler := rootHandler(d)

	routes := []Route{
		// 1.1 and 1.2 routes are simply a Go replacement for the equivalent Perl route. They may or may not conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).
		// 1.3 routes exist only in Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		// NOTE: Route IDs are immutable and unique. DO NOT change the ID of an existing Route; otherwise, existing
		// configurations may break. New Route IDs can be any integer between 0 and 2147483647 (inclusive), as long as
		// it's unique.

		// API Capability
		{1.1, http.MethodGet, `api_capabilities/?(\.json)?$`, apicapability.GetAPICapabilitiesHandler, auth.PrivLevelReadOnly, Authenticated, nil, 1813206589, perlBypass},

		//ASN: CRUD
		{1.2, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 473877722, noPerlBypass},
		{1.1, http.MethodGet, `asns/?(\.json)?$`, asn.V11ReadAll, auth.PrivLevelReadOnly, Authenticated, nil, 570341929, noPerlBypass},
		{1.1, http.MethodGet, `asns/{id}$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 123008984, noPerlBypass},
		{1.1, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 1951198629, noPerlBypass},
		{1.1, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 1999492188, noPerlBypass},
		{1.1, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 1672524769, noPerlBypass},

		// Traffic Stats access
		{1.2, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil, 1319569028, perlBypass},
		{1.2, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil, 1497997906, perlBypass},
		{1.2, http.MethodGet, `current_stats/?(\.json)?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil, 1785442893, perlBypass},

		{1.1, http.MethodGet, `caches/stats/?(\.json)?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil, 1813206588, noPerlBypass},

		//CacheGroup: CRUD
		{1.1, http.MethodGet, `cachegroups/trimmed/?(\.json)?$`, cachegroup.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil, 329527916, noPerlBypass},
		{1.1, http.MethodGet, `cachegroups/?(\.json)?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, 123079110, noPerlBypass},
		{1.1, http.MethodGet, `cachegroups/{id}$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, 691886338, noPerlBypass},
		{1.1, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 112954546, noPerlBypass},
		{1.1, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 32982665, noPerlBypass},
		{1.1, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 257869365, noPerlBypass},

		{1.1, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil, 1071644110, noPerlBypass},
		{1.1, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandler, auth.PrivLevelOperations, Authenticated, nil, 1520240431, noPerlBypass},

		//CacheGroup Parameters: CRUD
		{1.1, http.MethodGet, `cachegroups/{id}/parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 912449723, perlBypass},
		{1.1, http.MethodGet, `cachegroups/{id}/unassigned_parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupUnassignedParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 1457339250, perlBypass},
		{1.1, http.MethodDelete, `cachegroupparameters/{cachegroupID}/{parameterId}$`, api.DeleteHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelOperations, Authenticated, nil, 912449733, perlBypass},

		//CDN
		{1.1, http.MethodGet, `cdns/name/{name}/sslkeys/?(\.json)?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, 1278581772, noPerlBypass},
		{1.1, http.MethodGet, `cdns/metric_types`, notImplementedHandler, 0, NoAuth, nil, 683165463, noPerlBypass}, // MUST NOT end in $, because the 1.x route is longer

		{1.1, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, 697185281, perlBypass},
		{1.1, http.MethodGet, `cdns/configs/?(\.json)?$`, cdn.GetConfigs, auth.PrivLevelReadOnly, Authenticated, nil, 1768437852, noPerlBypass},

		{1.1, http.MethodGet, `cdns/{name}/health/?(\.json)?$`, cdn.GetNameHealth, auth.PrivLevelReadOnly, Authenticated, nil, 1135348194, perlBypass},
		{1.1, http.MethodGet, `cdns/health/?(\.json)?$`, cdn.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, 1085381134, perlBypass},

		{1.1, http.MethodGet, `cdns/domains/?(\.json)?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil, 296902560, noPerlBypass},
		{1.1, http.MethodGet, `cdns/routing$`, handlerToFunc(proxyHandler), 0, NoAuth, []middleware.Middleware{}, 66722982, noPerlBypass},

		//CDN: CRUD
		{1.1, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 1345914650, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 2122954075, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 2135233288, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 549326357, noPerlBypass},
		{1.1, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 24013912, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 1595587002, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil, 408804959, noPerlBypass},

		//CDN: queue updates
		{1.1, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil, 271515980, noPerlBypass},
		{1.1, http.MethodPost, `cdns/dnsseckeys/generate(\.json)?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 675336, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/delete/?(\.json)?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 571104207, noPerlBypass},
		{1.4, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 479010609, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeysV11, auth.PrivLevelAdmin, Authenticated, nil, 1427173311, noPerlBypass},

		{1.4, http.MethodGet, `cdns/dnsseckeys/refresh/?(\.json)?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil, 1771997116, noPerlBypass},

		//CDN: Monitoring: Traffic Monitor
		{1.1, http.MethodGet, `cdns/{cdn}/configs/monitoring(\.json)?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil, 2140847892, noPerlBypass},

		//Database dumps
		{1.1, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil, 274016647, perlBypass},

		//Division: CRUD
		{1.1, http.MethodGet, `divisions/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 1085181534, noPerlBypass},
		{1.1, http.MethodGet, `divisions/{id}$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 1241497902, noPerlBypass},
		{1.1, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 306369140, noPerlBypass},
		{1.1, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 553713800, noPerlBypass},
		{1.1, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 1325382237, noPerlBypass},
		{1.1, http.MethodGet, `divisions/name/{name}/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 1211408769, noPerlBypass},

		{1.1, http.MethodGet, `logs/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, 848340550, perlBypass},
		{1.1, http.MethodGet, `logs/{days}/days/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, 1192414145, perlBypass},
		{1.1, http.MethodGet, `logs/newcount/?(\.json)?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil, 1405833012, perlBypass},

		//HWInfo
		{1.1, http.MethodGet, `hwinfo/?(\.json)?$`, hwinfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, 621685998, noPerlBypass},

		//Content invalidation jobs
		{1.1, http.MethodGet, `jobs(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, 1966782041, perlBypass},
		{1.4, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil, 616780776, noPerlBypass},
		{1.4, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil, 186134226, noPerlBypass},
		{1.4, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil, 80450955, noPerlBypass},
		{1.1, http.MethodGet, `jobs/{id}(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, 2085189426, perlBypass},
		{1.1, http.MethodPost, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.CreateUserJob, auth.PrivLevelPortal, Authenticated, nil, 611328688, perlBypass},
		{1.1, http.MethodGet, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.GetUserJobs, auth.PrivLevelReadOnly, Authenticated, nil, 349163540, perlBypass},

		//Login
		{1.1, http.MethodGet, `users/{id}/deliveryservices/?(\.json)?$`, user.GetDSes, auth.PrivLevelReadOnly, Authenticated, nil, 988787789, noPerlBypass},
		{1.1, http.MethodGet, `user/{id}/deliveryservices/available/?(\.json)?$`, user.GetAvailableDSes, auth.PrivLevelReadOnly, Authenticated, nil, 757082995, noPerlBypass},
		{1.1, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil, 1392670821, noPerlBypass},
		{1.1, http.MethodPost, `user/logout(/|\.json)?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil, 443434825, perlBypass},
		{1.4, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil, 1415886009, noPerlBypass},
		{1.1, http.MethodPost, `user/login/token(/|\.json)?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil, 402408841, perlBypass},
		{1.1, http.MethodPost, `user/reset_password(/|\.json)?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil, 2092914630, perlBypass},
		{1.1, http.MethodPost, `users/register(/|\.json)?$`, login.RegisterUser, auth.PrivLevelOperations, Authenticated, nil, 1337, perlBypass},

		//ISO
		{1.1, http.MethodGet, `osversions(/|\.json)?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil, 576088657, perlBypass},

		//User: CRUD
		{1.1, http.MethodGet, `users/?(\.json)?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, 1491929900, noPerlBypass},
		{1.1, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, 713809980, noPerlBypass},
		{1.1, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, 135433404, noPerlBypass},
		{1.1, http.MethodPost, `users/?(\.json)?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, 876244816, noPerlBypass},

		{1.1, http.MethodGet, `user/current/?(\.json)?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil, 1610701614, noPerlBypass},
		{1.1, http.MethodPut, `user/current(/|\.json)?$`, user.ReplaceCurrent, auth.PrivLevelReadOnly, Authenticated, nil, 420, perlBypass},

		//Parameter: CRUD
		{1.1, http.MethodGet, `parameters/?(\.json)?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 2012554292, noPerlBypass},
		{1.1, http.MethodGet, `parameters/{id}$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 1221666841, noPerlBypass},
		{1.1, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 1873936115, noPerlBypass},
		{1.1, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 1669510859, noPerlBypass},
		{1.1, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 276277118, noPerlBypass},

		//Phys_Location: CRUD
		{1.1, http.MethodGet, `phys_locations/?(\.json)?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, 120405182, noPerlBypass},
		{1.1, http.MethodGet, `phys_locations/trimmed/?(\.json)?$`, physlocation.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil, 1097221000, noPerlBypass},
		{1.1, http.MethodGet, `phys_locations/{id}$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, 1554216025, noPerlBypass},
		{1.1, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 226795021, noPerlBypass},
		{1.1, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 2146456648, noPerlBypass},
		{1.1, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 15614221, noPerlBypass},

		//Ping
		{1.1, http.MethodGet, `ping$`, ping.PingHandler(), 0, NoAuth, nil, 1555661597, noPerlBypass},
		{1.1, http.MethodGet, `riak/ping/?(\.json)?$`, ping.Riak, auth.PrivLevelReadOnly, Authenticated, nil, 1884012114, noPerlBypass},
		{1.1, http.MethodGet, `keys/ping/?(\.json)?$`, ping.Keys, auth.PrivLevelReadOnly, Authenticated, nil, 318416022, noPerlBypass},

		//Profile: CRUD
		{1.1, http.MethodGet, `profiles/?(\.json)?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, 668758589, noPerlBypass},
		{1.1, http.MethodGet, `profiles/trimmed/?(\.json)?$`, profile.Trimmed, auth.PrivLevelReadOnly, Authenticated, nil, 644942941, noPerlBypass},

		{1.1, http.MethodGet, `profiles/{id}$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, 1570260672, noPerlBypass},
		{1.1, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 98439172, noPerlBypass},
		{1.1, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 1540211556, noPerlBypass},
		{1.1, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 2005594465, noPerlBypass},

		{1.1, http.MethodGet, `profiles/{id}/export/?(\.json)?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil, 30133517, perlBypass},
		{1.1, http.MethodPost, `profiles/import/?(\.json)?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil, 806143208, perlBypass},

		// Copy Profile
		{1.1, http.MethodPost, `profiles/name/{new_profile}/copy/{existing_profile}`, profile.CopyProfileHandler, auth.PrivLevelOperations, Authenticated, nil, 806143209, perlBypass},

		//Region: CRUDs
		{1.1, http.MethodGet, `regions/?(\.json)?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, 410037085, noPerlBypass},
		{1.1, http.MethodGet, `regions/{id}$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, 2024440051, noPerlBypass},
		{1.1, http.MethodGet, `regions/name/{name}/?(\.json)?$`, region.GetName, auth.PrivLevelReadOnly, Authenticated, nil, 503583197, noPerlBypass},
		{1.1, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 226308224, noPerlBypass},
		{1.1, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 1288334488, noPerlBypass},
		{1.1, http.MethodDelete, `regions/{id}$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 1181575271, noPerlBypass},

		{1.1, http.MethodDelete, `deliveryservice_server/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil, 1532184523, noPerlBypass},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{1.4, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandlerV14, auth.PrivLevelReadOnly, Authenticated, nil, 1946145033, noPerlBypass},
		{1.1, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil, 1928775049, noPerlBypass},
		{1.1, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil, 429799788, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil, 1428181206, noPerlBypass},
		{1.1, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 133115411, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil, 1345121223, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, dsserver.GetReadUnassigned, auth.PrivLevelReadOnly, Authenticated, nil, 2023944221, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil, 740875299, perlBypass},
		{1.1, http.MethodGet, `deliveryservice_matches/?(\.json)?$`, deliveryservice.GetMatches, auth.PrivLevelReadOnly, Authenticated, nil, 1191301170, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/{id}/capacity/?(\.json)?$`, deliveryservice.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, 1231409110, perlBypass},

		//Server
		{1.1, http.MethodGet, `servers/status$`, server.GetServersStatusCountsHandler, auth.PrivLevelReadOnly, Authenticated, nil, 2052786293, perlBypass},
		{1.1, http.MethodGet, `servers/totals$`, handlerToFunc(proxyHandler), 0, NoAuth, []middleware.Middleware{}, 2037840835, noPerlBypass},

		//Serverchecks
		{1.1, http.MethodGet, `servers/checks$`, handlerToFunc(proxyHandler), 0, NoAuth, []middleware.Middleware{}, 1796112922, noPerlBypass},
		{1.1, http.MethodPost, `servercheck/?(\.json)?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil, 1764281568, perlBypass},

		//Server Details
		{1.1, http.MethodGet, `servers/details/?(\.json)?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil, 1261264714, noPerlBypass},
		{1.1, http.MethodGet, `servers/hostname/{hostName}/details/?(\.json)?$`, server.GetDetailHandler, auth.PrivLevelReadOnly, Authenticated, nil, 372366128, noPerlBypass},

		//Server status
		{1.1, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil, 776663851, perlBypass},
		{1.1, http.MethodPost, `servers/{id}/queue_update$`, server.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, 9189471, perlBypass},

		//Server: CRUD
		{1.1, http.MethodGet, `servers/?(\.json)?$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil, 1720959285, noPerlBypass},
		{1.1, http.MethodGet, `servers/{id}$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil, 1543122028, noPerlBypass},
		{1.1, http.MethodPut, `servers/{id}$`, api.UpdateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 958634103, noPerlBypass},
		{1.1, http.MethodPost, `servers/?$`, api.CreateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 2025558061, noPerlBypass},
		{1.1, http.MethodDelete, `servers/{id}$`, api.DeleteHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 192322233, noPerlBypass},

		//Server Capability
		{1.4, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 610407391, noPerlBypass},
		{1.4, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 1074470708, noPerlBypass},
		{1.4, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 736415038, noPerlBypass},

		//Server Server Capabilities: CRUD
		{1.4, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 1800231889, noPerlBypass},
		{1.4, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 2093166834, noPerlBypass},
		{1.4, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 1058714058, noPerlBypass},

		//Status: CRUD
		{1.1, http.MethodGet, `statuses/?(\.json)?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, 2044905656, noPerlBypass},
		{1.1, http.MethodGet, `statuses/{id}$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, 1899095947, noPerlBypass},
		{1.1, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 1207966504, noPerlBypass},
		{1.1, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 1369123612, noPerlBypass},
		{1.1, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 755111360, noPerlBypass},

		//System
		{1.1, http.MethodGet, `system/info/?(\.json)?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, 211047475, noPerlBypass},

		//Type: CRUD
		{1.1, http.MethodGet, `types/?(\.json)?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, 2026701823, noPerlBypass},
		{1.1, http.MethodGet, `types/{id}$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, 86037256, noPerlBypass},
		{1.1, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 68860115, noPerlBypass},
		{1.1, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 1513308195, noPerlBypass},
		{1.1, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 93175773, noPerlBypass},

		//About
		{1.3, http.MethodGet, `about/?(\.json)?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil, 1317501166, noPerlBypass},

		//Coordinates
		{1.3, http.MethodGet, `coordinates/?(\.json)?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, 696700745, noPerlBypass},
		{1.3, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, 244546706, noPerlBypass},
		{1.3, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 368926174, noPerlBypass},
		{1.3, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 1428112157, noPerlBypass},
		{1.3, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 1303849889, noPerlBypass},

		//ASNs
		{1.3, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 2017162392, noPerlBypass},
		{1.3, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 2064172317, noPerlBypass},
		{1.3, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 859114392, noPerlBypass},
		{1.3, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 680204898, noPerlBypass},

		//CDN generic handlers:
		{1.3, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 1230318621, noPerlBypass},
		{1.3, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 632290057, noPerlBypass},
		{1.3, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 1311178934, noPerlBypass},
		{1.3, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 1160505289, noPerlBypass},
		{1.3, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 2007694657, noPerlBypass},

		//Delivery service requests
		{1.3, http.MethodGet, `deliveryservice_requests/?(\.json)?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil, 1681163935, noPerlBypass},
		{1.3, http.MethodGet, `deliveryservice_requests/?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil, 286812311, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_requests/?$`, api.UpdateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 2049907918, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservice_requests/?$`, api.CreateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 59385039, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservice_requests/?$`, api.DeleteHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 1296985025, noPerlBypass},

		//Delivery service request: Actions
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/assign$`, api.UpdateHandler(dsrequest.GetAssignmentSingleton()), auth.PrivLevelOperations, Authenticated, nil, 1703160290, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/status$`, api.UpdateHandler(dsrequest.GetStatusSingleton()), auth.PrivLevelPortal, Authenticated, nil, 668415099, noPerlBypass},

		//Delivery service request comment: CRUD
		{1.3, http.MethodGet, `deliveryservice_request_comments/?(\.json)?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil, 1032650737, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 860487847, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 727227672, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 199504668, noPerlBypass},

		//Delivery service uri signing keys: CRUD
		{1.3, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 1293078558, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 508466335, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 47648969, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 429925417, noPerlBypass},

		//Delivery Service Required Capabilities: CRUD
		{1.4, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 1158522227, noPerlBypass},
		{1.4, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, 1096873992, noPerlBypass},
		{1.4, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, 1496289304, noPerlBypass},

		// Federations by CDN (the actual table for federation)
		{1.1, http.MethodGet, `cdns/{name}/federations/?(\.json)?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, 989225032, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{name}/federations/{id}$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, 21850599, noPerlBypass},
		{1.1, http.MethodPost, `cdns/{name}/federations/?(\.json)?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 1954894219, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 2106065466, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 1442852902, noPerlBypass},

		{1.4, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil, 872924281, noPerlBypass},

		//Origins
		{1.3, http.MethodGet, `origins/?(\.json)?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, 844649256, noPerlBypass},
		{1.3, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, 1945936793, noPerlBypass},
		{1.3, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 141567746, noPerlBypass},
		{1.3, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 1099561643, noPerlBypass},
		{1.3, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 460273263, noPerlBypass},

		//Roles
		{1.1, http.MethodGet, `roles/?(\.json)?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil, 187088583, perlBypass},
		{1.3, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 1612897489, noPerlBypass},
		{1.3, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 430652406, noPerlBypass},
		{1.3, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 1356705982, noPerlBypass},

		//Delivery Services Regexes
		{1.1, http.MethodGet, `deliveryservices_regexes/?(\.json)?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil, 605501453, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil, 577432763, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.DSGetID, auth.PrivLevelReadOnly, Authenticated, nil, 1044974567, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil, 412737800, noPerlBypass},
		{1.1, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil, 1248339691, noPerlBypass},
		{1.1, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil, 2046731663, noPerlBypass},

		//Servers
		{1.3, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil, 880128253, noPerlBypass},
		{1.3, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil, 438451599, noPerlBypass},

		//StaticDNSEntries
		{1.1, http.MethodGet, `staticdnsentries/?(\.json)?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, 258939477, noPerlBypass},
		{1.3, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, 1116932668, noPerlBypass},
		{1.3, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 142457111, noPerlBypass},
		{1.3, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 1629148238, noPerlBypass},
		{1.3, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 1846031132, noPerlBypass},

		//ProfileParameters
		{1.1, http.MethodGet, `profiles/{id}/parameters/?(\.json)?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil, 876464975, noPerlBypass},
		{1.1, http.MethodGet, `profiles/{id}/unassigned_parameters/?(\.json)?$`, profileparameter.GetUnassigned, auth.PrivLevelReadOnly, Authenticated, nil, 574429262, noPerlBypass},
		{1.1, http.MethodGet, `profiles/name/{name}/parameters/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, 2067737832, noPerlBypass},
		{1.1, http.MethodGet, `parameters/profile/{name}/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, 1802599194, noPerlBypass},
		{1.1, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil, 1355945582, noPerlBypass},
		{1.1, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil, 316818708, noPerlBypass},
		{1.1, http.MethodGet, `profileparameters/?(\.json)?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 850609805, noPerlBypass},
		{1.1, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, 218809693, noPerlBypass},
		{1.1, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil, 234275, noPerlBypass},
		{1.1, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil, 1080610861, noPerlBypass},
		{1.1, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, 254839529, noPerlBypass},

		//Tenants
		{1.1, http.MethodGet, `tenants/?(\.json)?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, 1677967814, noPerlBypass},
		{1.1, http.MethodGet, `tenants/{id}$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, 171544338, noPerlBypass},
		{1.1, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 1094131478, noPerlBypass},
		{1.1, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 917248013, noPerlBypass},
		{1.1, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 516365558, noPerlBypass},

		//CRConfig
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil, 1957273695, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil, 676716889, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{id}/snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, 854424150, noPerlBypass},
		{1.1, http.MethodPut, `snapshot/{cdn}/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, 1969911829, noPerlBypass},

		// ATS config files
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/?(\.json)?$`, atsserver.GetConfigMetaData, auth.PrivLevelOperations, Authenticated, nil, 1755842214, perlBypass},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_revalidate\.config/?(\.json)?$`, atscdn.GetRegexRevalidateDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1810067775, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_mid_{xml-id}\.config/?(\.json)?$`, atscdn.GetMidHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 658322121, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_{xml-id}\.config/?(\.json)?$`, atscdn.GetEdgeHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1894063777, perlBypass},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/bg_fetch\.config/?(\.json)?$`, atscdn.GetBGFetchDotConfig, auth.PrivLevelOperations, Authenticated, nil, 160404036, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/cacheurl{filename}\.config/?(\.json)?$`, atscdn.GetCacheURLDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1373111113, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_remap_{ds-name}\.config/?(\.json)?$`, atscdn.GetRegexRemapDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1283602930, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/set_dscp_{dscp}\.config/?(\.json)?$`, atscdn.GetSetDSCPDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1889993740, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/ssl_multicert\.config/?(\.json)?$`, atscdn.GetSSLMultiCertDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1113687166, perlBypass},

		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/12M_facts/?$`, atsprofile.GetFacts, auth.PrivLevelOperations, Authenticated, nil, 2146608231, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/50-ats\.rules/?$`, atsprofile.GetATSDotRules, auth.PrivLevelOperations, Authenticated, nil, 1101032000, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/astats\.config/?$`, atsprofile.GetAstats, auth.PrivLevelOperations, Authenticated, nil, 1362661662, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/cache\.config/?$`, atsprofile.GetCache, auth.PrivLevelOperations, Authenticated, nil, 292387870, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/drop_qstring\.config/?$`, atsprofile.GetDropQString, auth.PrivLevelOperations, Authenticated, nil, 1097869291, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging\.config/?$`, atsprofile.GetLogging, auth.PrivLevelOperations, Authenticated, nil, 172702063, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging\.yaml/?$`, atsprofile.GetLoggingYAML, auth.PrivLevelOperations, Authenticated, nil, 453568059, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logs_xml\.config/?$`, atsprofile.GetLogsXML, auth.PrivLevelOperations, Authenticated, nil, 1309053227, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/plugin\.config/?$`, atsprofile.GetPlugin, auth.PrivLevelOperations, Authenticated, nil, 274047559, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/records\.config/?$`, atsprofile.GetRecords, auth.PrivLevelOperations, Authenticated, nil, 469014057, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/storage\.config/?$`, atsprofile.GetStorage, auth.PrivLevelOperations, Authenticated, nil, 121977329, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/sysctl\.conf/?$`, atsprofile.GetSysctl, auth.PrivLevelOperations, Authenticated, nil, 1202950646, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/url_sig_{file}\.config/?$`, atsprofile.GetURLSig, auth.PrivLevelOperations, Authenticated, nil, 448450070, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/uri_signing_{file}\.config/?$`, atsprofile.GetURISigning, auth.PrivLevelOperations, Authenticated, nil, 125995582, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/volume\.config/?$`, atsprofile.GetVolume, auth.PrivLevelOperations, Authenticated, nil, 792704719, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/{file}/?$`, atsprofile.GetUnknown, auth.PrivLevelOperations, Authenticated, nil, 1651257268, perlBypass},

		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/parent\.config/?(\.json)?$`, atsserver.GetParentDotConfig, auth.PrivLevelOperations, Authenticated, nil, 645056066, perlBypass},
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/remap\.config/?(\.json)?$`, atsserver.GetServerConfigRemap, auth.PrivLevelOperations, Authenticated, nil, 2038454899, perlBypass},

		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/cache\.config/?(\.json)?$`, atsserver.GetCacheDotConfig, auth.PrivLevelOperations, Authenticated, nil, 34686861, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/ip_allow\.config/?(\.json)?$`, atsserver.GetIPAllowDotConfig, auth.PrivLevelOperations, Authenticated, nil, 724651786, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/hosting\.config/?(\.json)?$`, atsserver.GetHostingDotConfig, auth.PrivLevelOperations, Authenticated, nil, 1387459113, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/packages/?(\.json)?$`, atsserver.GetPackages, auth.PrivLevelOperations, Authenticated, nil, 245024839, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/chkconfig/?(\.json)?$`, atsserver.GetChkconfig, auth.PrivLevelOperations, Authenticated, nil, 1012457987, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/{file}/?(\.json)?$`, atsserver.GetUnknown, auth.PrivLevelOperations, Authenticated, nil, 322079218, perlBypass},

		// Federations
		{1.4, http.MethodGet, `federations/all/?(\.json)?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil, 61059986, noPerlBypass},
		{1.1, http.MethodGet, `federations/?(\.json)?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil, 154954994, noPerlBypass},
		{1.1, http.MethodPost, `federations(/|\.json)?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 1894064742, perlBypass},
		{1.1, http.MethodDelete, `federations(/|\.json)?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 592098323, perlBypass},
		{1.1, http.MethodPut, `federations(/|\.json)?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 1283182516, perlBypass},
		{1.1, http.MethodPost, `federations/{id}/deliveryservices/?(\.json)?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil, 1682863513, noPerlBypass},
		{1.1, http.MethodGet, `federations/{id}/deliveryservices/?(\.json)?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil, 353773034, perlBypass},
		{1.1, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?(\.json)?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil, 1417402570, perlBypass},

		// Federation Resolvers
		{1.1, http.MethodPost, `federation_resolvers(/|\.json)?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil, 1134373661, perlBypass},
		{1.1, http.MethodGet, `federation_resolvers(/|\.json)?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil, 556608759, perlBypass},
		{1.1, http.MethodPost, `federations/{id}/federation_resolvers(/|\.json)?$`, federations.AssignFederationResolversToFederationHandler, auth.PrivLevelAdmin, Authenticated, nil, 556608760, perlBypass},
		{1.1, http.MethodGet, `federations/{id}/federation_resolvers(/|\.json)?$`, federations.GetFederationFederationResolversHandler, auth.PrivLevelReadOnly, Authenticated, nil, 556608761, perlBypass},
		{1.5, http.MethodDelete, `federation_resolvers(/|\.json)?$`, federation_resolvers.Delete, auth.PrivLevelAdmin, Authenticated, nil, 9001, noPerlBypass},
		{1.1, http.MethodDelete, `federation_resolvers/{id}(/|\.json)?$`, federation_resolvers.DeleteByID, auth.PrivLevelAdmin, Authenticated, nil, 42, perlBypass},

		// Federations Users
		{1.1, http.MethodPost, `federations/{id}/users/?(\.json)?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil, 1779334930, perlBypass},
		{1.1, http.MethodGet, `federations/{id}/users/?(\.json)?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil, 394075015, perlBypass},
		{1.1, http.MethodDelete, `federations/{id}/users/{userID}/?(\.json)?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil, 1949102882, perlBypass},

		////DeliveryServices
		{1.1, http.MethodGet, `deliveryservices/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 1238317294, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 444348195, noPerlBypass},

		{1.5, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV15, auth.PrivLevelOperations, Authenticated, nil, 506431432, noPerlBypass},
		{1.4, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV14, auth.PrivLevelOperations, Authenticated, nil, 506431431, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV13, auth.PrivLevelOperations, Authenticated, nil, 1705681904, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV12, auth.PrivLevelOperations, Authenticated, nil, 652813412, noPerlBypass},

		{1.5, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV15, auth.PrivLevelOperations, Authenticated, nil, 1766567527, noPerlBypass},
		{1.4, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV14, auth.PrivLevelOperations, Authenticated, nil, 1766567526, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV13, auth.PrivLevelOperations, Authenticated, nil, 1559124565, noPerlBypass},
		{1.1, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV12, auth.PrivLevelOperations, Authenticated, nil, 597160536, noPerlBypass},

		{1.1, http.MethodDelete, `deliveryservices/{id}/?(\.json)?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil, 242642074, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/{id}/servers/eligible/?(\.json)?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil, 474761584, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLID, auth.PrivLevelAdmin, Authenticated, nil, 1135772907, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/hostname/{hostname}/sslkeys$`, deliveryservice.GetSSLKeysByHostName, auth.PrivLevelAdmin, Authenticated, nil, 2105792225, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, 1872878583, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys/delete$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil, 1926734, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/generate/?(\.json)?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil, 753439051, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?(\.json)?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil, 1262501076, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?(\.json)?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil, 1530482824, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil, 2102719211, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil, 393197114, noPerlBypass},
		{1.1, http.MethodGet, `riak/bucket/{bucket}/key/{key}/values/?(\.json)?$`, apiriak.GetBucketKey, auth.PrivLevelAdmin, Authenticated, nil, 2020510801, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/{id}/health/?(\.json)?$`, deliveryservice.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, 2034590101, perlBypass},

		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/?(\.json)?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 1569607824, noPerlBypass},
		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/{target}$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 105995848, noPerlBypass},
		{1.1, http.MethodPost, `steering/{deliveryservice}/targets/?(\.json)?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 1338216397, noPerlBypass},
		{1.1, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 1438608295, noPerlBypass},
		{1.1, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 2088021515, noPerlBypass},

		// Stats Summary
		{1.1, http.MethodGet, `stats_summary/?(\.json)?$`, trafficstats.GetStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, 380498598, perlBypass},

		// TO Extensions
		{1.5, http.MethodPost, `to_extensions$`, toextension.CreateTOExtension, auth.PrivLevelReadOnly, Authenticated, nil, 380498599, noPerlBypass},

		//Pattern based consistent hashing endpoint
		{1.4, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil, 1960755076, noPerlBypass},

		{1.4, http.MethodGet, `steering/?(\.json)?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil, 1174852457, noPerlBypass},
	}

	// sanity check to make sure all Route IDs are unique
	knownRouteIDs := make(map[int]struct{}, len(routes))
	for _, r := range routes {
		if _, found := knownRouteIDs[r.ID]; !found {
			knownRouteIDs[r.ID] = struct{}{}
		} else {
			return nil, nil, nil, fmt.Errorf("route ID %d is already taken. Please give it a unique Route ID", r.ID)
		}
	}

	// verify configured perl_routes are actually able to pass through to Perl
	perlRoutes := GetRouteIDMap(d.PerlRoutes)
	for _, r := range routes {
		if _, isPerlRoute := perlRoutes[r.ID]; isPerlRoute && !r.CanBypassToPerl {
			return nil, nil, nil, fmt.Errorf("route '%s' is configured as a perl_route but cannot be passed through to Perl", r.String())
		}
	}

	// check for unknown route IDs in cdn.conf
	disabledRoutes := GetRouteIDMap(d.DisabledRoutes)
	unknownRouteIDs := []string{}
	for _, routeMap := range []map[int]struct{}{perlRoutes, disabledRoutes} {
		for routeID := range routeMap {
			if _, known := knownRouteIDs[routeID]; !known {
				unknownRouteIDs = append(unknownRouteIDs, fmt.Sprintf("%d", routeID))
			}
		}
	}
	if len(unknownRouteIDs) > 0 {
		msg := "unknown route IDs in routing_blacklist: " + strings.Join(unknownRouteIDs, ", ")
		if d.IgnoreUnknownRoutes {
			log.Warnln(msg)
		} else {
			return nil, nil, nil, errors.New(msg)
		}
	}

	// rawRoutes are served at the root path. These should be almost exclusively old Perl pre-API routes, which have yet to be converted in all clients. New routes should be in the versioned API path.
	rawRoutes := []RawRoute{
		// DEPRECATED - use PUT /api/1.2/snapshot/{cdn}
		{http.MethodGet, `tools/write_crconfig/{cdn}/?$`, crconfig.SnapshotOldGUIHandler, auth.PrivLevelOperations, Authenticated, nil},
		// DEPRECATED - use GET /api/1.2/cdns/{cdn}/snapshot
		{http.MethodGet, `CRConfig-Snapshots/{cdn}/CRConfig.json?$`, crconfig.SnapshotOldGetHandler, auth.PrivLevelReadOnly, Authenticated, nil},
	}

	return routes, rawRoutes, proxyHandler, nil
}

func MemoryStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		stats := runtime.MemStats{}
		runtime.ReadMemStats(&stats)

		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Errorln("unable to marshal stats: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("marshalling error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

func DBStatsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		stats := db.DB.Stats()

		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Errorln("unable to marshal stats: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("marshalling error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

// RootHandler returns the / handler for the service, which reverse-proxies the old Perl Traffic Ops
func rootHandler(d ServerData) http.Handler {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(d.Config.ProxyTimeout) * time.Second,
			KeepAlive: time.Duration(d.Config.ProxyKeepAlive) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(d.Config.ProxyTLSTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(d.Config.ProxyReadHeaderTimeout) * time.Second,
		//IdleConnTimeout: time.Duration(d.Config.ProxyIdleConnTimeout) * time.Second,
		//Other knobs we can turn: ExpectContinueTimeout,IdleConnTimeout
	}
	rp := httputil.NewSingleHostReverseProxy(d.URL)
	rp.Transport = tr

	rp.ErrorLog = log.StandardLogger(log.Error, "proxy error: ")
	riak.SetErrorLogger(log.StandardLogger(log.Error, "riak error: "))
	riak.SetLogger(log.StandardLogger(log.Info, "riak info: "))

	log.Debugf("our reverseProxy: %++v\n", rp)
	log.Debugf("our reverseProxy's transport: %++v\n", tr)
	loggingProxyHandler := middleware.WrapAccessLog(d.Secrets[0], rp)

	managerHandler := CreateThrottledHandler(loggingProxyHandler, d.BackendMaxConnections["mojolicious"])
	return managerHandler
}

// notImplementedHandler returns a 501 Not Implemented to the client. This should be used very rarely, and primarily for old API Perl routes which were broken long ago, which we don't have the resources to rewrite in Go for the time being.
func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotImplemented
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

//CreateThrottledHandler takes a handler, and a max and uses a channel to insure the handler is used concurrently by only max number of routines
func CreateThrottledHandler(handler http.Handler, maxConcurrentCalls int) ThrottledHandler {
	return ThrottledHandler{handler, make(chan struct{}, maxConcurrentCalls)}
}

// ThrottledHandler ...
type ThrottledHandler struct {
	Handler http.Handler
	ReqChan chan struct{}
}

func (m ThrottledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ReqChan <- struct{}{}
	defer func() { <-m.ReqChan }()
	m.Handler.ServeHTTP(w, r)
}
