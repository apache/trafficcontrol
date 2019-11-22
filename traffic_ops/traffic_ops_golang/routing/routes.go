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
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercheck"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steering"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steeringtargets"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/systeminfo"
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

const perlBypass = true
const noPerlBypass = false

func handlerToFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func getRouteIDMap(IDs []int) map[int]struct{} {
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
		// 1.3 routes exist only in a Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		//ASN: CRUD
		{1.2, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 0, noPerlBypass},
		{1.1, http.MethodGet, `asns/?(\.json)?$`, asn.V11ReadAll, auth.PrivLevelReadOnly, Authenticated, nil, 1, noPerlBypass},
		{1.1, http.MethodGet, `asns/{id}$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 2, noPerlBypass},
		{1.1, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 3, noPerlBypass},
		{1.1, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 4, noPerlBypass},
		{1.1, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 5, noPerlBypass},

		// Traffic Stats access
		{1.2, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil, 6, perlBypass},
		{1.2, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil, 7, perlBypass},
		{1.2, http.MethodGet, `current_stats/?(\.json)?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil, 1785442893, perlBypass},

		{1.1, http.MethodGet, `caches/stats/?(\.json)?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil, 8, noPerlBypass},

		//CacheGroup: CRUD
		{1.1, http.MethodGet, `cachegroups/trimmed/?(\.json)?$`, cachegroup.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil, 9, noPerlBypass},
		{1.1, http.MethodGet, `cachegroups/?(\.json)?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, 10, noPerlBypass},
		{1.1, http.MethodGet, `cachegroups/{id}$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, 11, noPerlBypass},
		{1.1, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 12, noPerlBypass},
		{1.1, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 13, noPerlBypass},
		{1.1, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, 14, noPerlBypass},

		{1.1, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil, 15, noPerlBypass},
		{1.1, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandler, auth.PrivLevelOperations, Authenticated, nil, 16, noPerlBypass},

		//CacheGroup Parameters: CRUD
		{1.1, http.MethodGet, `cachegroups/{id}/parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 17, perlBypass},
		{1.1, http.MethodGet, `cachegroups/{id}/unassigned_parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupUnassignedParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 18, perlBypass},

		//CDN
		{1.1, http.MethodGet, `cdns/name/{name}/sslkeys/?(\.json)?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, 19, noPerlBypass},
		{1.1, http.MethodGet, `cdns/metric_types`, notImplementedHandler, 0, NoAuth, nil, 20, noPerlBypass}, // MUST NOT end in $, because the 1.x route is longer

		{1.1, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, 21, perlBypass},
		{1.1, http.MethodGet, `cdns/configs/?(\.json)?$`, cdn.GetConfigs, auth.PrivLevelReadOnly, Authenticated, nil, 22, noPerlBypass},

		{1.1, http.MethodGet, `cdns/domains/?(\.json)?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil, 23, noPerlBypass},
		{1.1, http.MethodGet, `cdns/health$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}, 24, noPerlBypass},
		{1.1, http.MethodGet, `cdns/routing$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}, 25, noPerlBypass},

		//CDN: CRUD
		{1.1, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 26, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 27, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 28, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 29, noPerlBypass},
		{1.1, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 30, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 31, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil, 32, noPerlBypass},

		//CDN: queue updates
		{1.1, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil, 33, noPerlBypass},
		{1.1, http.MethodPost, `cdns/dnsseckeys/generate(\.json)?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 34, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/delete/?(\.json)?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 35, noPerlBypass},
		{1.4, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, 36, noPerlBypass},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeysV11, auth.PrivLevelAdmin, Authenticated, nil, 37, noPerlBypass},

		{1.4, http.MethodGet, `cdns/dnsseckeys/refresh/?(\.json)?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil, 38, noPerlBypass},

		//CDN: Monitoring: Traffic Monitor
		{1.1, http.MethodGet, `cdns/{cdn}/configs/monitoring(\.json)?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil, 39, noPerlBypass},

		//Database dumps
		{1.1, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil, 40, perlBypass},

		//Division: CRUD
		{1.1, http.MethodGet, `divisions/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 41, noPerlBypass},
		{1.1, http.MethodGet, `divisions/{id}$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 42, noPerlBypass},
		{1.1, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 43, noPerlBypass},
		{1.1, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 44, noPerlBypass},
		{1.1, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, 45, noPerlBypass},
		{1.1, http.MethodGet, `divisions/name/{name}/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, 46, noPerlBypass},

		{1.1, http.MethodGet, `logs/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, 47, perlBypass},
		{1.1, http.MethodGet, `logs/{days}/days/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, 48, perlBypass},
		{1.1, http.MethodGet, `logs/newcount/?(\.json)?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil, 49, perlBypass},

		//HWInfo
		{1.1, http.MethodGet, `hwinfo/?(\.json)?$`, hwinfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, 50, noPerlBypass},

		//Content invalidation jobs
		{1.1, http.MethodGet, `jobs(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, 51, perlBypass},
		{1.4, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil, 52, noPerlBypass},
		{1.4, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil, 53, noPerlBypass},
		{1.4, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil, 54, noPerlBypass},
		{1.1, http.MethodGet, `jobs/{id}(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, 55, perlBypass},
		{1.1, http.MethodPost, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.CreateUserJob, auth.PrivLevelPortal, Authenticated, nil, 56, perlBypass},
		{1.1, http.MethodGet, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.GetUserJobs, auth.PrivLevelReadOnly, Authenticated, nil, 57, perlBypass},

		//Login
		{1.1, http.MethodGet, `users/{id}/deliveryservices/?(\.json)?$`, user.GetDSes, auth.PrivLevelReadOnly, Authenticated, nil, 58, noPerlBypass},
		{1.1, http.MethodGet, `user/{id}/deliveryservices/available/?(\.json)?$`, user.GetAvailableDSes, auth.PrivLevelReadOnly, Authenticated, nil, 59, noPerlBypass},
		{1.1, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil, 60, noPerlBypass},
		{1.1, http.MethodPost, `user/logout(/|\.json)?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil, 61, perlBypass},
		{1.4, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil, 62, noPerlBypass},
		{1.1, http.MethodPost, `user/login/token(/|\.json)?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil, 63, perlBypass},
		{1.1, http.MethodPost, `user/reset_password(/|\.json)?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil, 64, perlBypass},

		//ISO
		{1.1, http.MethodGet, `osversions(/|\.json)?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil, 65, perlBypass},

		//User: CRUD
		{1.1, http.MethodGet, `users/?(\.json)?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, 66, noPerlBypass},
		{1.1, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, 67, noPerlBypass},
		{1.1, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, 68, noPerlBypass},
		{1.1, http.MethodPost, `users/?(\.json)?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, 69, noPerlBypass},

		{1.1, http.MethodGet, `user/current/?(\.json)?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil, 70, noPerlBypass},

		//Parameter: CRUD
		{1.1, http.MethodGet, `parameters/?(\.json)?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 71, noPerlBypass},
		{1.1, http.MethodGet, `parameters/{id}$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 72, noPerlBypass},
		{1.1, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 73, noPerlBypass},
		{1.1, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 74, noPerlBypass},
		{1.1, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, 75, noPerlBypass},

		//Phys_Location: CRUD
		{1.1, http.MethodGet, `phys_locations/?(\.json)?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, 76, noPerlBypass},
		{1.1, http.MethodGet, `phys_locations/trimmed/?(\.json)?$`, physlocation.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil, 77, noPerlBypass},
		{1.1, http.MethodGet, `phys_locations/{id}$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, 78, noPerlBypass},
		{1.1, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 79, noPerlBypass},
		{1.1, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 80, noPerlBypass},
		{1.1, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, 81, noPerlBypass},

		//Ping
		{1.1, http.MethodGet, `ping$`, ping.PingHandler(), 0, NoAuth, nil, 82, noPerlBypass},
		{1.1, http.MethodGet, `riak/ping/?(\.json)?$`, ping.Riak, auth.PrivLevelReadOnly, Authenticated, nil, 83, noPerlBypass},
		{1.1, http.MethodGet, `keys/ping/?(\.json)?$`, ping.Keys, auth.PrivLevelReadOnly, Authenticated, nil, 84, noPerlBypass},

		//Profile: CRUD
		{1.1, http.MethodGet, `profiles/?(\.json)?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, 85, noPerlBypass},
		{1.1, http.MethodGet, `profiles/trimmed/?(\.json)?$`, profile.Trimmed, auth.PrivLevelReadOnly, Authenticated, nil, 86, noPerlBypass},

		{1.1, http.MethodGet, `profiles/{id}$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, 87, noPerlBypass},
		{1.1, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 88, noPerlBypass},
		{1.1, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 89, noPerlBypass},
		{1.1, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, 90, noPerlBypass},

		{1.1, http.MethodGet, `profiles/{id}/export/?(\.json)?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil, 91, perlBypass},
		{1.1, http.MethodPost, `profiles/import/?(\.json)?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil, 92, perlBypass},

		//Region: CRUDs
		{1.1, http.MethodGet, `regions/?(\.json)?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, 93, noPerlBypass},
		{1.1, http.MethodGet, `regions/{id}$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, 94, noPerlBypass},
		{1.1, http.MethodGet, `regions/name/{name}/?(\.json)?$`, region.GetName, auth.PrivLevelReadOnly, Authenticated, nil, 95, noPerlBypass},
		{1.1, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 96, noPerlBypass},
		{1.1, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 97, noPerlBypass},
		{1.1, http.MethodDelete, `regions/{id}$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, 98, noPerlBypass},

		{1.1, http.MethodDelete, `deliveryservice_server/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil, 99, noPerlBypass},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{1.4, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandlerV14, auth.PrivLevelReadOnly, Authenticated, nil, 100, noPerlBypass},
		{1.1, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil, 101, noPerlBypass},
		{1.1, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil, 102, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil, 103, noPerlBypass},
		{1.1, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 104, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil, 105, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, dsserver.GetReadUnassigned, auth.PrivLevelReadOnly, Authenticated, nil, 106, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil, 107, perlBypass},
		{1.1, http.MethodGet, `deliveryservice_matches/?(\.json)?$`, deliveryservice.GetMatches, auth.PrivLevelReadOnly, Authenticated, nil, 108, noPerlBypass},

		//Server
		{1.1, http.MethodGet, `servers/status$`, server.GetServersStatusCountsHandler, auth.PrivLevelReadOnly, Authenticated, nil, 109, perlBypass},
		{1.1, http.MethodGet, `servers/totals$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}, 110, noPerlBypass},

		//Serverchecks
		{1.1, http.MethodGet, `servers/checks$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}, 111, noPerlBypass},
		{1.1, http.MethodPost, `servercheck/?(\.json)?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil, 112, perlBypass},

		//Server Details
		{1.1, http.MethodGet, `servers/details/?(\.json)?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil, 113, noPerlBypass},
		{1.1, http.MethodGet, `servers/hostname/{hostName}/details/?(\.json)?$`, server.GetDetailHandler, auth.PrivLevelReadOnly, Authenticated, nil, 114, noPerlBypass},

		//Server status
		{1.1, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil, 115, perlBypass},

		//Server: CRUD
		{1.1, http.MethodGet, `servers/?(\.json)?$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil, 116, noPerlBypass},
		{1.1, http.MethodGet, `servers/{id}$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil, 117, noPerlBypass},
		{1.1, http.MethodPut, `servers/{id}$`, api.UpdateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 118, noPerlBypass},
		{1.1, http.MethodPost, `servers/?$`, api.CreateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 119, noPerlBypass},
		{1.1, http.MethodDelete, `servers/{id}$`, api.DeleteHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil, 120, noPerlBypass},

		//Server Capability
		{1.4, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 121, noPerlBypass},
		{1.4, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 122, noPerlBypass},
		{1.4, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 123, noPerlBypass},

		//Server Server Capabilities: CRUD
		{1.4, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 124, noPerlBypass},
		{1.4, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 125, noPerlBypass},
		{1.4, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, 126, noPerlBypass},

		//Status: CRUD
		{1.1, http.MethodGet, `statuses/?(\.json)?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, 127, noPerlBypass},
		{1.1, http.MethodGet, `statuses/{id}$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, 128, noPerlBypass},
		{1.1, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 129, noPerlBypass},
		{1.1, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 130, noPerlBypass},
		{1.1, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, 131, noPerlBypass},

		//System
		{1.1, http.MethodGet, `system/info/?(\.json)?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, 132, noPerlBypass},

		//Type: CRUD
		{1.1, http.MethodGet, `types/?(\.json)?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, 133, noPerlBypass},
		{1.1, http.MethodGet, `types/{id}$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, 134, noPerlBypass},
		{1.1, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 135, noPerlBypass},
		{1.1, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 136, noPerlBypass},
		{1.1, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, 137, noPerlBypass},

		//About
		{1.3, http.MethodGet, `about/?(\.json)?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil, 138, noPerlBypass},

		//Coordinates
		{1.3, http.MethodGet, `coordinates/?(\.json)?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, 139, noPerlBypass},
		{1.3, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, 140, noPerlBypass},
		{1.3, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 141, noPerlBypass},
		{1.3, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 142, noPerlBypass},
		{1.3, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, 143, noPerlBypass},

		//ASNs
		{1.3, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 144, noPerlBypass},
		{1.3, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 145, noPerlBypass},
		{1.3, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 146, noPerlBypass},
		{1.3, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, 147, noPerlBypass},

		//CDN generic handlers:
		{1.3, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 148, noPerlBypass},
		{1.3, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, 149, noPerlBypass},
		{1.3, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 150, noPerlBypass},
		{1.3, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 151, noPerlBypass},
		{1.3, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, 152, noPerlBypass},

		//Delivery service requests
		{1.3, http.MethodGet, `deliveryservice_requests/?(\.json)?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil, 153, noPerlBypass},
		{1.3, http.MethodGet, `deliveryservice_requests/?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil, 154, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_requests/?$`, api.UpdateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 155, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservice_requests/?$`, api.CreateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 156, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservice_requests/?$`, api.DeleteHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil, 157, noPerlBypass},

		//Delivery service request: Actions
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/assign$`, api.UpdateHandler(dsrequest.GetAssignmentSingleton()), auth.PrivLevelOperations, Authenticated, nil, 158, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/status$`, api.UpdateHandler(dsrequest.GetStatusSingleton()), auth.PrivLevelPortal, Authenticated, nil, 159, noPerlBypass},

		//Delivery service request comment: CRUD
		{1.3, http.MethodGet, `deliveryservice_request_comments/?(\.json)?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil, 160, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 161, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 162, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, 163, noPerlBypass},

		//Delivery service uri signing keys: CRUD
		{1.3, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 164, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 165, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 166, noPerlBypass},
		{1.3, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, 167, noPerlBypass},

		//Delivery Service Required Capabilities: CRUD
		{1.4, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, 168, noPerlBypass},
		{1.4, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, 169, noPerlBypass},
		{1.4, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, 170, noPerlBypass},

		// Federations by CDN (the actual table for federation)
		{1.1, http.MethodGet, `cdns/{name}/federations/?(\.json)?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, 171, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{name}/federations/{id}$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, 172, noPerlBypass},
		{1.1, http.MethodPost, `cdns/{name}/federations/?(\.json)?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 173, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 174, noPerlBypass},
		{1.1, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, 175, noPerlBypass},

		{1.4, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil, 176, noPerlBypass},

		//Origins
		{1.3, http.MethodGet, `origins/?(\.json)?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, 177, noPerlBypass},
		{1.3, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, 178, noPerlBypass},
		{1.3, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 179, noPerlBypass},
		{1.3, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 180, noPerlBypass},
		{1.3, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, 181, noPerlBypass},

		//Roles
		{1.1, http.MethodGet, `roles/?(\.json)?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil, 182, perlBypass},
		{1.3, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 183, noPerlBypass},
		{1.3, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 184, noPerlBypass},
		{1.3, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, 185, noPerlBypass},

		//Delivery Services Regexes
		{1.1, http.MethodGet, `deliveryservices_regexes/?(\.json)?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil, 186, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil, 187, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.DSGetID, auth.PrivLevelReadOnly, Authenticated, nil, 188, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil, 189, noPerlBypass},
		{1.1, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil, 190, noPerlBypass},
		{1.1, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil, 191, noPerlBypass},

		//Servers
		{1.3, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil, 192, noPerlBypass},
		{1.3, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil, 193, noPerlBypass},

		//StaticDNSEntries
		{1.1, http.MethodGet, `staticdnsentries/?(\.json)?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, 194, noPerlBypass},
		{1.3, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, 195, noPerlBypass},
		{1.3, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 196, noPerlBypass},
		{1.3, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 197, noPerlBypass},
		{1.3, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, 198, noPerlBypass},

		//ProfileParameters
		{1.1, http.MethodGet, `profiles/{id}/parameters/?(\.json)?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil, 199, noPerlBypass},
		{1.1, http.MethodGet, `profiles/{id}/unassigned_parameters/?(\.json)?$`, profileparameter.GetUnassigned, auth.PrivLevelReadOnly, Authenticated, nil, 200, noPerlBypass},
		{1.1, http.MethodGet, `profiles/name/{name}/parameters/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, 201, noPerlBypass},
		{1.1, http.MethodGet, `parameters/profile/{name}/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, 202, noPerlBypass},
		{1.1, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil, 203, noPerlBypass},
		{1.1, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil, 204, noPerlBypass},
		{1.1, http.MethodGet, `profileparameters/?(\.json)?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, 205, noPerlBypass},
		{1.1, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, 206, noPerlBypass},
		{1.1, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil, 207, noPerlBypass},
		{1.1, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil, 208, noPerlBypass},
		{1.1, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, 209, noPerlBypass},

		//Tenants
		{1.1, http.MethodGet, `tenants/?(\.json)?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, 210, noPerlBypass},
		{1.1, http.MethodGet, `tenants/{id}$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, 211, noPerlBypass},
		{1.1, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 212, noPerlBypass},
		{1.1, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 213, noPerlBypass},
		{1.1, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, 214, noPerlBypass},

		//CRConfig
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil, 215, noPerlBypass},
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil, 216, noPerlBypass},
		{1.1, http.MethodPut, `cdns/{id}/snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, 217, noPerlBypass},
		{1.1, http.MethodPut, `snapshot/{cdn}/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, 218, noPerlBypass},

		// ATS config files
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/?(\.json)?$`, atsserver.GetConfigMetaData, auth.PrivLevelOperations, Authenticated, nil, 219, perlBypass},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_revalidate.config/?(\.json)?$`, atscdn.GetRegexRevalidateDotConfig, auth.PrivLevelOperations, Authenticated, nil, 220, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_mid_{xml-id}.config/?(\.json)?$`, atscdn.GetMidHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 221, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_{xml-id}.config/?(\.json)?$`, atscdn.GetEdgeHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 222, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_revalidate.config/?(\.json)?$`, atscdn.GetRegexRevalidateDotConfig, auth.PrivLevelOperations, Authenticated, nil, 223, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_mid_{xml-id}.config/?(\.json)?$`, atscdn.GetMidHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 224, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_{xml-id}.config/?(\.json)?$`, atscdn.GetEdgeHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil, 225, perlBypass},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/bg_fetch.config/?(\.json)?$`, atscdn.GetBGFetchDotConfig, auth.PrivLevelOperations, Authenticated, nil, 226, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/cacheurl{filename}.config/?(\.json)?$`, atscdn.GetCacheURLDotConfig, auth.PrivLevelOperations, Authenticated, nil, 227, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_remap_{ds-name}.config/?(\.json)?$`, atscdn.GetRegexRemapDotConfig, auth.PrivLevelOperations, Authenticated, nil, 228, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/set_dscp_{dscp}.config/?(\.json)?$`, atscdn.GetSetDSCPDotConfig, auth.PrivLevelOperations, Authenticated, nil, 229, perlBypass},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/ssl_multicert.config/?(\.json)?$`, atscdn.GetSSLMultiCertDotConfig, auth.PrivLevelOperations, Authenticated, nil, 230, perlBypass},

		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/12M_facts/?$`, atsprofile.GetFacts, auth.PrivLevelOperations, Authenticated, nil, 231, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/50-ats.rules/?$`, atsprofile.GetATSDotRules, auth.PrivLevelOperations, Authenticated, nil, 232, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/astats.config/?$`, atsprofile.GetAstats, auth.PrivLevelOperations, Authenticated, nil, 233, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/cache.config/?$`, atsprofile.GetCache, auth.PrivLevelOperations, Authenticated, nil, 234, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/drop_qstring.config/?$`, atsprofile.GetDropQString, auth.PrivLevelOperations, Authenticated, nil, 235, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging.config/?$`, atsprofile.GetLogging, auth.PrivLevelOperations, Authenticated, nil, 236, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging.yaml/?$`, atsprofile.GetLoggingYAML, auth.PrivLevelOperations, Authenticated, nil, 237, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logs_xml.config/?$`, atsprofile.GetLogsXML, auth.PrivLevelOperations, Authenticated, nil, 238, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/plugin.config/?$`, atsprofile.GetPlugin, auth.PrivLevelOperations, Authenticated, nil, 239, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/records.config/?$`, atsprofile.GetRecords, auth.PrivLevelOperations, Authenticated, nil, 240, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/storage.config/?$`, atsprofile.GetStorage, auth.PrivLevelOperations, Authenticated, nil, 241, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/sysctl.conf/?$`, atsprofile.GetSysctl, auth.PrivLevelOperations, Authenticated, nil, 242, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/url_sig_{file}.config/?$`, atsprofile.GetURLSig, auth.PrivLevelOperations, Authenticated, nil, 243, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/uri_signing_{file}.config/?$`, atsprofile.GetURISigning, auth.PrivLevelOperations, Authenticated, nil, 244, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/volume.config/?$`, atsprofile.GetVolume, auth.PrivLevelOperations, Authenticated, nil, 245, perlBypass},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/{file}/?$`, atsprofile.GetUnknown, auth.PrivLevelOperations, Authenticated, nil, 246, perlBypass},

		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/parent.config/?(\.json)?$`, atsserver.GetParentDotConfig, auth.PrivLevelOperations, Authenticated, nil, 247, perlBypass},
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/remap.config/?(\.json)?$`, atsserver.GetServerConfigRemap, auth.PrivLevelOperations, Authenticated, nil, 248, perlBypass},

		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/cache.config/?(\.json)?$`, atsserver.GetCacheDotConfig, auth.PrivLevelOperations, Authenticated, nil, 249, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/ip_allow.config/?(\.json)?$`, atsserver.GetIPAllowDotConfig, auth.PrivLevelOperations, Authenticated, nil, 250, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/hosting.config/?(\.json)?$`, atsserver.GetHostingDotConfig, auth.PrivLevelOperations, Authenticated, nil, 251, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/packages/?(\.json)?$`, atsserver.GetPackages, auth.PrivLevelOperations, Authenticated, nil, 252, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/chkconfig/?(\.json)?$`, atsserver.GetChkconfig, auth.PrivLevelOperations, Authenticated, nil, 253, perlBypass},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/{file}/?(\.json)?$`, atsserver.GetUnknown, auth.PrivLevelOperations, Authenticated, nil, 254, perlBypass},

		// Federations
		{1.4, http.MethodGet, `federations/all/?(\.json)?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil, 255, noPerlBypass},
		{1.1, http.MethodGet, `federations/?(\.json)?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil, 256, noPerlBypass},
		{1.1, http.MethodPost, `federations(/|\.json)?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 257, perlBypass},
		{1.1, http.MethodDelete, `federations(/|\.json)?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 258, perlBypass},
		{1.1, http.MethodPut, `federations(/|\.json)?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, 259, perlBypass},
		{1.1, http.MethodPost, `federations/{id}/deliveryservices?(\.json)?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil, 260, noPerlBypass},
		{1.1, http.MethodGet, `federations/{id}/deliveryservices?(\.json)?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil, 261, perlBypass},
		{1.1, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?(\.json)?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil, 262, perlBypass},

		// Federation Resolvers
		{1.1, http.MethodPost, `federation_resolvers(/|\.json)?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil, 263, perlBypass},
		{1.1, http.MethodGet, `federation_resolvers(/|\.json)?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil, 264, perlBypass},

		// Federations Users
		{1.1, http.MethodPost, `federations/{id}/users?(\.json)?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil, 265, perlBypass},
		{1.1, http.MethodGet, `federations/{id}/users?(\.json)?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil, 266, perlBypass},
		{1.1, http.MethodDelete, `federations/{id}/users/{userID}/?(\.json)?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil, 267, perlBypass},

		////DeliveryServices
		{1.1, http.MethodGet, `deliveryservices/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 268, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, 269, noPerlBypass},

		{1.4, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV14, auth.PrivLevelOperations, Authenticated, nil, 270, noPerlBypass},
		{1.3, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV13, auth.PrivLevelOperations, Authenticated, nil, 271, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV12, auth.PrivLevelOperations, Authenticated, nil, 272, noPerlBypass},

		{1.4, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV14, auth.PrivLevelOperations, Authenticated, nil, 273, noPerlBypass},
		{1.3, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV13, auth.PrivLevelOperations, Authenticated, nil, 274, noPerlBypass},
		{1.1, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV12, auth.PrivLevelOperations, Authenticated, nil, 275, noPerlBypass},

		{1.1, http.MethodDelete, `deliveryservices/{id}/?(\.json)?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil, 276, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/{id}/servers/eligible/?(\.json)?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil, 277, noPerlBypass},

		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLID, auth.PrivLevelAdmin, Authenticated, nil, 278, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/hostname/{hostname}/sslkeys$`, deliveryservice.GetSSLKeysByHostName, auth.PrivLevelAdmin, Authenticated, nil, 279, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, 280, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys/delete$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil, 281, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/generate/?(\.json)?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil, 282, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?(\.json)?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil, 283, noPerlBypass},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?(\.json)?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil, 284, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil, 285, noPerlBypass},
		{1.1, http.MethodGet, `deliveryservices/{id}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil, 286, noPerlBypass},
		{1.1, http.MethodGet, `riak/bucket/{bucket}/key/{key}/values/?(\.json)?$`, apiriak.GetBucketKey, auth.PrivLevelAdmin, Authenticated, nil, 287, noPerlBypass},

		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/?(\.json)?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 288, noPerlBypass},
		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/{target}$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, 289, noPerlBypass},
		{1.1, http.MethodPost, `steering/{deliveryservice}/targets/?(\.json)?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 290, noPerlBypass},
		{1.1, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 291, noPerlBypass},
		{1.1, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, 292, noPerlBypass},

		//Pattern based consistent hashing endpoint
		{1.4, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil, 293, noPerlBypass},

		{1.4, http.MethodGet, `steering/?(\.json)?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil, 294, noPerlBypass},
	}

	// sanity check to make sure all Route IDs are unique
	routeIDs := make(map[int]struct{}, len(routes))
	for _, r := range routes {
		if _, found := routeIDs[r.ID]; !found {
			routeIDs[r.ID] = struct{}{}
		} else {
			return nil, nil, nil, fmt.Errorf("route ID %d is already taken. Please give it a unique Route ID", r.ID)
		}
	}

	// verify configured perl_routes are actually able to pass through to Perl
	perlRoutes := getRouteIDMap(d.PerlRoutes)
	for _, r := range routes {
		if _, isPerlRoute := perlRoutes[r.ID]; isPerlRoute && !r.CanBypassToPerl {
			return nil, nil, nil, fmt.Errorf("route '%s' is configured as a perl_route but cannot be passed through to Perl", r.String())
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
	loggingProxyHandler := wrapAccessLog(d.Secrets[0], rp)

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
