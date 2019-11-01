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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/federations"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/hwinfo"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/invalidationjobs"
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
var Authenticated = true

// NoAuth ...
var NoAuth = false

func handlerToFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// Routes returns the API routes, raw non-API root level routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, []RawRoute, http.Handler, error) {
	proxyHandler := rootHandler(d)

	routes := []Route{
		// 1.1 and 1.2 routes are simply a Go replacement for the equivalent Perl route. They may or may not conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).
		// 1.3 routes exist only in a Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		//ASN: CRUD
		{1.2, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `asns/?(\.json)?$`, asn.V11ReadAll, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `asns/{id}$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},

		// Traffic Stats access
		{1.2, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.2, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `caches/stats/?(\.json)?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil},

		//CacheGroup: CRUD
		{1.1, http.MethodGet, `cachegroups/trimmed/?(\.json)?$`, cachegroup.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cachegroups/?(\.json)?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cachegroups/{id}$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandler, auth.PrivLevelOperations, Authenticated, nil},

		//CacheGroup Parameters: CRUD
		{1.1, http.MethodGet, `cachegroups/{id}/parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cachegroups/{id}/unassigned_parameters/?(\.json)?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupUnassignedParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},

		//CDN
		{1.1, http.MethodGet, `cdns/name/{name}/sslkeys/?(\.json)?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/metric_types`, notImplementedHandler, 0, NoAuth, nil}, // MUST NOT end in $, because the 1.x route is longer

		{1.1, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/configs/?(\.json)?$`, cdn.GetConfigs, auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `cdns/domains/?(\.json)?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/health$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `cdns/routing$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},

		//CDN: CRUD
		{1.1, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/name/{name}/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil},

		//CDN: queue updates
		{1.1, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cdns/dnsseckeys/generate(\.json)?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/delete/?(\.json)?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{1.4, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/name/{name}/dnsseckeys/?(\.json)?$`, cdn.GetDNSSECKeysV11, auth.PrivLevelAdmin, Authenticated, nil},

		{1.4, http.MethodGet, `cdns/dnsseckeys/refresh/?(\.json)?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil},

		//CDN: Monitoring: Traffic Monitor
		{1.1, http.MethodGet, `cdns/{cdn}/configs/monitoring(\.json)?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		//Database dumps
		{1.1, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil},

		//Division: CRUD
		{1.1, http.MethodGet, `divisions/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `divisions/{id}$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `divisions/name/{name}/?(\.json)?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `logs/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `logs/{days}/days/?(\.json)?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `logs/newcount/?(\.json)?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil},

		//HWInfo
		{1.1, http.MethodGet, `hwinfo-wip/?(\.json)?$`, hwinfo.Get, auth.PrivLevelReadOnly, Authenticated, nil},

		//Content invalidation jobs
		{1.1, http.MethodGet, `jobs(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.4, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil},
		{1.4, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil},
		{1.4, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil},
		{1.1, http.MethodGet, `jobs/{id}(/|\.json/?)?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.CreateUserJob, auth.PrivLevelPortal, Authenticated, nil},
		{1.1, http.MethodGet, `user/current/jobs(/|\.json/?)?$`, invalidationjobs.GetUserJobs, auth.PrivLevelReadOnly, Authenticated, nil},

		//Login
		{1.1, http.MethodGet, `users/{id}/deliveryservices/?(\.json)?$`, user.GetDSes, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `user/{id}/deliveryservices/available/?(\.json)?$`, user.GetAvailableDSes, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{1.1, http.MethodPost, `user/logout(/|\.json)?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil},
		{1.4, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{1.1, http.MethodPost, `user/login/token(/|\.json)?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{1.1, http.MethodPost, `user/reset_password(/|\.json)?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil},

		//User: CRUD
		{1.1, http.MethodGet, `users/?(\.json)?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `users/?(\.json)?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `user/current/?(\.json)?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil},

		//Parameter: CRUD
		{1.1, http.MethodGet, `parameters/?(\.json)?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `parameters/{id}$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},

		//Phys_Location: CRUD
		{1.1, http.MethodGet, `phys_locations/?(\.json)?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `phys_locations/trimmed/?(\.json)?$`, physlocation.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `phys_locations/{id}$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},

		//Ping
		{1.1, http.MethodGet, `ping$`, ping.PingHandler(), 0, NoAuth, nil},
		{1.1, http.MethodGet, `riak/ping/?(\.json)?$`, ping.Riak, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `keys/ping/?(\.json)?$`, ping.Keys, auth.PrivLevelReadOnly, Authenticated, nil},

		//Profile: CRUD
		{1.1, http.MethodGet, `profiles/?(\.json)?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/trimmed/?(\.json)?$`, profile.Trimmed, auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `profiles/{id}$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `profiles/{id}/export/?(\.json)?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/import/?(\.json)?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil},

		//Region: CRUDs
		{1.1, http.MethodGet, `regions/?(\.json)?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `regions/{id}$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `regions/name/{name}/?(\.json)?$`, region.GetName, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `regions/{id}$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodDelete, `deliveryservice_server/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{1.4, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandlerV14, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryserviceserver/?(\.json)?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, dsserver.GetReadUnassigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservice_matches/?(\.json)?$`, deliveryservice.GetMatches, auth.PrivLevelReadOnly, Authenticated, nil},

		//Server
		{1.1, http.MethodGet, `servers/status$`, server.GetServersStatusCountsHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `servers/totals$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},

		//Serverchecks
		{1.1, http.MethodGet, `servers/checks$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodPost, `servercheck/?(\.json)?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil},

		//Server Details
		{1.1, http.MethodGet, `servers/details/?(\.json)?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `servers/hostname/{hostName}/details/?(\.json)?$`, server.GetDetailHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		//Server status
		{1.1, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil},

		//Server: CRUD
		{1.1, http.MethodGet, `servers/?(\.json)?$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id}$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `servers/{id}$`, api.UpdateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `servers/?$`, api.CreateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `servers/{id}$`, api.DeleteHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},

		//Server Capability
		{1.4, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.4, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.4, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		//Server Server Capabilities: CRUD
		{1.4, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.4, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.4, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		//Status: CRUD
		{1.1, http.MethodGet, `statuses/?(\.json)?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `statuses/{id}$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},

		//System
		{1.1, http.MethodGet, `system/info/?(\.json)?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil},

		//Type: CRUD
		{1.1, http.MethodGet, `types/?(\.json)?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `types/{id}$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},

		//About
		{1.3, http.MethodGet, `about/?(\.json)?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil},

		//Coordinates
		{1.3, http.MethodGet, `coordinates/?(\.json)?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},

		//ASNs
		{1.3, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},

		//CDN generic handlers:
		{1.3, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},

		//Delivery service requests
		{1.3, http.MethodGet, `deliveryservice_requests/?(\.json)?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `deliveryservice_requests/?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_requests/?$`, api.UpdateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservice_requests/?$`, api.CreateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservice_requests/?$`, api.DeleteHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service request: Actions
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/assign$`, api.UpdateHandler(dsrequest.GetAssignmentSingleton()), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/status$`, api.UpdateHandler(dsrequest.GetStatusSingleton()), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service request comment: CRUD
		{1.3, http.MethodGet, `deliveryservice_request_comments/?(\.json)?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service uri signing keys: CRUD
		{1.3, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},

		//Delivery Service Required Capabilities: CRUD
		{1.4, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.4, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.4, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		// Federations by CDN (the actual table for federation)
		{1.1, http.MethodGet, `cdns/{name}/federations/?(\.json)?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{name}/federations/{id}$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `cdns/{name}/federations/?(\.json)?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},

		{1.4, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil},

		//Origins
		{1.3, http.MethodGet, `origins/?(\.json)?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},

		//Roles
		{1.1, http.MethodGet, `roles/?(\.json)?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},

		//Delivery Services Regexes
		{1.1, http.MethodGet, `deliveryservices_regexes/?(\.json)?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.DSGetID, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil},

		//Servers
		{1.3, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		//StaticDNSEntries
		{1.1, http.MethodGet, `staticdnsentries/?(\.json)?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},

		//ProfileParameters
		{1.1, http.MethodGet, `profiles/{id}/parameters/?(\.json)?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{id}/unassigned_parameters/?(\.json)?$`, profileparameter.GetUnassigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/name/{name}/parameters/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `parameters/profile/{name}/?(\.json)?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profileparameters/?(\.json)?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil},

		//Tenants
		{1.1, http.MethodGet, `tenants/?(\.json)?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `tenants/{id}$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},

		//CRConfig
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `cdns/{id}/snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPut, `snapshot/{cdn}/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil},

		// ATS config files
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/?(\.json)?$`, atsserver.GetConfigMetaData, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_revalidate.config/?(\.json)?$`, atscdn.GetRegexRevalidateDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_mid_{xml-id}.config/?(\.json)?$`, atscdn.GetMidHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_{xml-id}.config/?(\.json)?$`, atscdn.GetEdgeHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_revalidate.config/?(\.json)?$`, atscdn.GetRegexRevalidateDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_mid_{xml-id}.config/?(\.json)?$`, atscdn.GetMidHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/hdr_rw_{xml-id}.config/?(\.json)?$`, atscdn.GetEdgeHeaderRewriteDotConfig, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/bg_fetch.config/?(\.json)?$`, atscdn.GetBGFetchDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/cacheurl{filename}.config/?(\.json)?$`, atscdn.GetCacheURLDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/regex_remap_{ds-name}.config/?(\.json)?$`, atscdn.GetRegexRemapDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/set_dscp_{dscp}.config/?(\.json)?$`, atscdn.GetSetDSCPDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn-name-or-id}/configfiles/ats/ssl_multicert.config/?(\.json)?$`, atscdn.GetSSLMultiCertDotConfig, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/12M_facts/?$`, atsprofile.GetFacts, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/50-ats.rules/?$`, atsprofile.GetATSDotRules, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/astats.config/?$`, atsprofile.GetAstats, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/cache.config/?$`, atsprofile.GetCache, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/drop_qstring.config/?$`, atsprofile.GetDropQString, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging.config/?$`, atsprofile.GetLogging, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logging.yaml/?$`, atsprofile.GetLoggingYAML, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/logs_xml.config/?$`, atsprofile.GetLogsXML, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/plugin.config/?$`, atsprofile.GetPlugin, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/records.config/?$`, atsprofile.GetRecords, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/storage.config/?$`, atsprofile.GetStorage, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/sysctl.conf/?$`, atsprofile.GetSysctl, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/url_sig_{file}.config/?$`, atsprofile.GetURLSig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/uri_signing_{file}.config/?$`, atsprofile.GetURISigning, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/volume.config/?$`, atsprofile.GetVolume, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{profile-name-or-id}/configfiles/ats/{file}/?$`, atsprofile.GetUnknown, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/parent.config/?(\.json)?$`, atsserver.GetParentDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{server-name-or-id}/configfiles/ats/remap.config/?(\.json)?$`, atsserver.GetServerConfigRemap, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/cache.config/?(\.json)?$`, atsserver.GetCacheDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/ip_allow.config/?(\.json)?$`, atsserver.GetIPAllowDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/hosting.config/?(\.json)?$`, atsserver.GetHostingDotConfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/packages/?(\.json)?$`, atsserver.GetPackages, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/chkconfig/?(\.json)?$`, atsserver.GetChkconfig, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id-or-host}/configfiles/ats/{file}/?(\.json)?$`, atsserver.GetUnknown, auth.PrivLevelOperations, Authenticated, nil},

		// Federations
		{1.4, http.MethodGet, `federations/all/?(\.json)?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `federations/?(\.json)?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil},
		{1.1, http.MethodPost, `federations(/|\.json)?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{1.1, http.MethodDelete, `federations(/|\.json)?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{1.1, http.MethodPut, `federations(/|\.json)?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{1.1, http.MethodPost, `federations/{id}/deliveryservices?(\.json)?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `federations/{id}/deliveryservices?(\.json)?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?(\.json)?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil},

		// Federations Users
		{1.1, http.MethodPost, `federations/{id}/users?(\.json)?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `federations/{id}/users?(\.json)?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodDelete, `federations/{id}/users/{userID}/?(\.json)?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil},

		////DeliveryServices
		{1.1, http.MethodGet, `deliveryservices/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/?(\.json)?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},

		{1.4, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV14, auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV13, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/?(\.json)?$`, deliveryservice.CreateV12, auth.PrivLevelOperations, Authenticated, nil},

		{1.4, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV14, auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV13, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPut, `deliveryservices/{id}/?(\.json)?$`, deliveryservice.UpdateV12, auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodDelete, `deliveryservices/{id}/?(\.json)?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodGet, `deliveryservices/{id}/servers/eligible/?(\.json)?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLID, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/hostname/{hostname}/sslkeys$`, deliveryservice.GetSSLKeysByHostName, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys/delete$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/sslkeys/generate/?(\.json)?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?(\.json)?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?(\.json)?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/urlkeys/?(\.json)?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `riak/bucket/{bucket}/key/{key}/values/?(\.json)?$`, apiriak.GetBucketKey, auth.PrivLevelAdmin, Authenticated, nil},

		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/?(\.json)?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `steering/{deliveryservice}/targets/{target}$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `steering/{deliveryservice}/targets/?(\.json)?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},
		{1.1, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},
		{1.1, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?(\.json)?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},

		//Pattern based consistent hashing endpoint
		{1.4, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil},

		{1.4, http.MethodGet, `steering/?(\.json)?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil},
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
