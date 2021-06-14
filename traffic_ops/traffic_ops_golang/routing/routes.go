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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/acme"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apicapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apitenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/asn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroupparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachesstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/capabilities"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdn_lock"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdnfederation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdnnotification"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/coordinate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/crconfig"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/crstats"
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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/invalidationjobs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/iso"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/login"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/logs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/origin"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/physlocation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ping"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/plugins"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profile"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profileparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/region"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/role"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercheck"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercheck/extensions"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servicecategory"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steering"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steeringtargets"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/systeminfo"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/topology"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/types"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/urisigning"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/user"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/vault"

	"github.com/jmoiron/sqlx"
)

// Authenticated indicates that a route requires authentication for use.
const Authenticated = true

// NoAuth indicates that a route does not require authentication for use.
const NoAuth = false

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

		// 2.x routes exist only in Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		// NOTE: Route IDs are immutable and unique. DO NOT change the ID of an existing Route; otherwise, existing
		// configurations may break. New Route IDs can be any integer between 0 and 2147483647 (inclusive), as long as
		// it's unique.

		/**
		 * 4.x API
		 */

		// CDN lock
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdn_locks/?$`, cdn_lock.Read, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4134390561},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdn_locks/?$`, cdn_lock.Create, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4134390562},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdn_locks/?$`, cdn_lock.Delete, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4134390564},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `acme_accounts/providers?$`, acme.ReadProviders, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4034390565},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/acme/?$`, deliveryservice.GenerateAcmeCertificates, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2534390576},

		// ACME account information
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `acme_accounts/?$`, acme.Read, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4034390561},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `acme_accounts/?$`, acme.Create, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4034390562},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `acme_accounts/?$`, acme.Update, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4034390563},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `acme_accounts/{provider}/{email}?$`, acme.Delete, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4034390564},

		//Delivery service ACME
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{xmlid}/sslkeys/renew$`, deliveryservice.RenewAcmeCertificate, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2534390573},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `acme_autorenew/?$`, deliveryservice.RenewCertificates, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2534390574},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `async_status/{id}$`, api.GetAsyncStatus, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2534390575},

		// API Capability
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `api_capabilities/?$`, apicapability.GetAPICapabilitiesHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 48132065893},

		//ASNs
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42641723173},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 402048983},

		//ASN: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `asns/?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4738777223},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 49511986293},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 49994921883},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 46725247693},

		// Traffic Stats access
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 43195690283},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 44979979063},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `current_stats/?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 47854428933},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `caches/stats/?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 48132065883},

		//CacheGroup: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cachegroups/?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4230791103},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4129545463},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 429826653},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4278693653},

		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil, DoCache, 40716441103},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandlerV40, auth.PrivLevelOperations, Authenticated, nil, DoCache, 45202404313},

		//Capabilities
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `capabilities/?$`, capabilities.Read, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 40081353},

		//CDN
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/name/{name}/sslkeys/?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 42785817723},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4971852813},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/{name}/health/?$`, cdn.GetNameHealth, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 41353481943},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/health/?$`, cdn.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 40853811343},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/domains/?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4269025603},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/routing$`, crstats.GetCDNRouting, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 467229823},

		//CDN: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4088049593},

		//CDN: queue updates
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4215159803},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdns/dnsseckeys/generate?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4753363},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdns/name/{name}/dnsseckeys?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4711042073},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/name/{name}/dnsseckeys/?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4790106093},

		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `cdns/dnsseckeys/refresh/?$`, cdn.RefreshDNSSECKeysV4, auth.PrivLevelOperations, Authenticated, nil, DoCache, 47719971163},

		//CDN: Monitoring: Traffic Monitor
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/{cdn}/configs/monitoring?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42408478923},

		//Database dumps
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4240166473},

		//Division: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `divisions/?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 40851815343},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4063691403},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4537138003},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 43253822373},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `logs/?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4483405503},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `logs/newcount/?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 44058330123},

		//Content invalidation jobs
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `jobs/?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 49667820413},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil, DoCache, 4167807763},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil, DoCache, 4861342263},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil, DoCache, 404509553},

		//Login
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 43926708213},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `user/logout/?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil, NoCache, 4434348253},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 44158860093},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `user/login/token/?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 4024088413},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `user/reset_password/?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil, NoCache, 42929146303},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `users/register/?$`, login.RegisterUser, auth.PrivLevelOperations, Authenticated, nil, NoCache, 43373},

		//ISO
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `osversions/?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4760886573},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `isos/?$`, iso.ISOs, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4760336573},

		//User: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `users/?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 44919299003},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4138099803},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4354334043},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `users/?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4762448163},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `user/current/?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 46107016143},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `user/current/?$`, user.ReplaceCurrent, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 4203},

		//Parameter: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `parameters/?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42125542923},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 48739361153},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 46695108593},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4262771183},

		//Phys_Location: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `phys_locations/?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4204051823},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4227950213},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42464566483},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 456142213},

		//Ping
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `ping$`, ping.Handler, 0, NoAuth, nil, DoCache, 45556615973},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `vault/ping/?$`, ping.Vault, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 48840121143},

		//Profile: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `profiles/?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4687585893},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 484391723},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 45402115563},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42055944653},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `profiles/{id}/export/?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 401335173},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profiles/import/?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4061432083},

		// Copy Profile
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profiles/name/{new_profile}/copy/{existing_profile}`, profile.CopyProfileHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4061432093},

		//Region: CRUDs
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `regions/?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4100370853},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4223082243},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42883344883},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `regions/?$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42326267583},

		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `topologies/?$`, api.CreateHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4871452221},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `topologies/?$`, api.ReadHandler(&topology.TOTopology{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4871452222},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `topologies/?$`, api.UpdateHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4871452223},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `topologies/?$`, api.DeleteHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4871452224},

		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `topologies/{name}/queue_update$`, topology.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4205351748},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryserviceserver/?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 49461450333},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4297997883},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryserviceserver/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil, DoCache, 45321845233},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 44281812063},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4331154113},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4801282533},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 43451212233},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/capacity/?$`, deliveryservice.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42314091103},
		//Serverchecks
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servercheck/?$`, servercheck.ReadServerCheck, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 47961129223},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servercheck/?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil, DoCache, 47642815683},

		// Servercheck Extensions
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servercheck/extensions$`, extensions.Create, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4804985993},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servercheck/extensions$`, extensions.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4834985993},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `servercheck/extensions/{id}$`, extensions.Delete, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4804982993},

		//Server Details
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servers/details/?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42612647143},

		//Server status
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4766638513},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servers/{id}/queue_update$`, server.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 41894713},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4384515993},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servers/{id-or-name}/update$`, server.UpdateHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 443813233},

		//Server: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `servers/?$`, server.Read, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 47209592853},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `servers/{id}$`, server.Update, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4586341033},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `servers/?$`, server.Create, auth.PrivLevelOperations, Authenticated, nil, DoCache, 42255580613},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `servers/{id}$`, server.Delete, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4923222333},

		//Server Capability
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4104073913},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 40744707083},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `server_capabilities$`, api.UpdateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42543770109},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4364150383},

		//Server Server Capabilities: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 48002318893},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42931668343},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 40587140583},

		//Status: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `statuses/?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42449056563},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 42079665043},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 43691236123},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4551113603},

		//System
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `system/info/?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4210474753},

		//Type: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `types/?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42267018233},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 488601153},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 45133081953},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 431757733},

		//About
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `about/?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 43175011663},

		//Coordinates
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4967007453},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4689261743},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 44281121573},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 43038498893},

		//CDN notification
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdn_notifications/?$`, cdnnotification.Read, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 2221224514},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdn_notifications/?$`, cdnnotification.Create, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2765223513},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdn_notifications/?$`, cdnnotification.Delete, auth.PrivLevelOperations, Authenticated, nil, DoCache, 2722411851},

		//CDN generic handlers:
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42303186213},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 43111789343},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 41605052893},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4276946573},

		//Delivery service requests
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservice_requests/?$`, dsrequest.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 46811639353},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservice_requests/?$`, dsrequest.Put, auth.PrivLevelPortal, Authenticated, nil, DoCache, 42499079183},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservice_requests/?$`, dsrequest.Post, auth.PrivLevelPortal, Authenticated, nil, DoCache, 493850393},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservice_requests/?$`, dsrequest.Delete, auth.PrivLevelPortal, Authenticated, nil, DoCache, 42969850253},

		//Delivery service request: Actions
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservice_requests/{id}/assign$`, dsrequest.GetAssignment, auth.PrivLevelOperations, Authenticated, nil, DoCache, 47031602904},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/assign$`, dsrequest.PutAssignment, auth.PrivLevelOperations, Authenticated, nil, DoCache, 47031602903},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservice_requests/{id}/status$`, dsrequest.GetStatus, auth.PrivLevelPortal, Authenticated, nil, DoCache, 4684150994},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/status$`, dsrequest.PutStatus, auth.PrivLevelPortal, Authenticated, nil, DoCache, 4684150993},

		//Delivery service request comment: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservice_request_comments/?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 40326507373},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, DoCache, 4604878473},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, DoCache, 4272276723},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, DoCache, 4995046683},

		//Delivery service uri signing keys: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 42930785583},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4084663353},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 476489693},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4299254173},

		//Delivery Service Required Capabilities: CRUD
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 41585222273},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 40968739923},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 44962893043},

		// Federations by CDN (the actual table for federation)
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/{name}/federations/?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4892250323},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdns/{name}/federations/?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 49548942193},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4260654663},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 44428529023},

		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4729242813},

		//Origins
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4446492563},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 415677463},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 40995616433},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4602732633},

		//Roles
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `roles/?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4870885833},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 46128974893},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4306524063},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 43567059823},

		//Delivery Services Regexes
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices_regexes/?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4055014533},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4774327633},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4127378003},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil, DoCache, 42483396913},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil, DoCache, 42467316633},

		//ServiceCategories
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `service_categories/?$`, api.ReadHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4085181543},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `service_categories/{name}/?$`, servicecategory.Update, auth.PrivLevelOperations, Authenticated, nil, DoCache, 406369141},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `service_categories/?$`, api.CreateHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 453713801},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `service_categories/{name}$`, api.DeleteHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4325382238},

		//StaticDNSEntries
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4289394773},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4424571113},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 46291482383},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 48460311323},

		//ProfileParameters
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `profiles/{id}/parameters/?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4764649753},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `profiles/name/{name}/parameters/?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42677378323},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil, DoCache, 43559455823},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4168187083},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `profileparameters/?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4506098053},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4288096933},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4242753},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil, DoCache, 40806108613},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4248395293},

		//Tenants
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `tenants/?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 46779678143},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 40941314783},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4172480133},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4163655583},

		//CRConfig
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 49572736953},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4767168893},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, DoCache, 49699118293},

		// Federations
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federations/all/?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 410599863},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federations/?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil, DoCache, 4549549943},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `federations/?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, DoCache, 48940647423},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `federations/?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, DoCache, 420983233},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `federations/?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, DoCache, 42831825163},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `federations/{id}/deliveryservices/?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 46828635133},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federations/{id}/deliveryservices/?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4537730343},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 44174025703},

		// Federation Resolvers
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `federation_resolvers/?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 41343736613},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federation_resolvers/?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4566087593},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `federations/{id}/federation_resolvers/?$`, federations.AssignFederationResolversToFederationHandler, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 4566087603},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federations/{id}/federation_resolvers/?$`, federations.GetFederationFederationResolversHandler, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4566087613},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `federation_resolvers/?$`, federation_resolvers.Delete, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 40013},

		// Federations Users
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `federations/{id}/users/?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 47793349303},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `federations/{id}/users/?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4940750153},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `federations/{id}/users/{userID}/?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil, DoCache, 49491028823},

		////DeliveryServices
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42383172943},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/?$`, deliveryservice.CreateV40, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4064315323},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservices/{id}/?$`, deliveryservice.UpdateV40, auth.PrivLevelOperations, Authenticated, nil, DoCache, 47665675673},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `deliveryservices/{id}/safe/?$`, deliveryservice.UpdateSafe, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4472109313},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/{id}/?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil, DoCache, 4226420743},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers/eligible/?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4747615843},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLIDV15, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 41357729073},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, DoCache, 48728785833},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil, DoCache, 49267343},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4534390513},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil, DoCache, 42625010763},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil, DoCache, 45304828243},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42027192113},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/xmlId/{name}/urlkeys/?$`, deliveryservice.DeleteURLKeysByName, auth.PrivLevelOperations, Authenticated, nil, DoCache, 42027192114},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/urlkeys/?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4931971143},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `deliveryservices/{id}/urlkeys/?$`, deliveryservice.DeleteURLKeysByID, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4931971144},

		//Delivery service LetsEncrypt
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/letsencrypt/?$`, deliveryservice.GenerateLetsEncryptCertificates, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4534390523},
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `letsencrypt/dnsrecords/?$`, deliveryservice.GetDnsChallengeRecords, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4534390553},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `letsencrypt/autorenew/?$`, deliveryservice.RenewCertificatesDeprecated, auth.PrivLevelOperations, Authenticated, nil, DoCache, 4534390563},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/health/?$`, deliveryservice.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 42345901013},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `deliveryservices/{id}/routing$`, crstats.GetDSRouting, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 467339833},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `steering/{deliveryservice}/targets/?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 45696078243},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `steering/{deliveryservice}/targets/?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, DoCache, 43382163973},
		{api.Version{Major: 4, Minor: 0}, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, DoCache, 44386082953},
		{api.Version{Major: 4, Minor: 0}, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, DoCache, 42880215153},

		// Stats Summary
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `stats_summary/?$`, trafficstats.GetStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4804985983},
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `stats_summary/?$`, trafficstats.CreateStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4804915983},

		//Pattern based consistent hashing endpoint
		{api.Version{Major: 4, Minor: 0}, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4607550763},

		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `steering/?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil, DoCache, 41748524573},

		// Plugins
		{api.Version{Major: 4, Minor: 0}, http.MethodGet, `plugins/?$`, plugins.Get(d.Plugins), auth.PrivLevelReadOnly, Authenticated, nil, DoCache, 4834985393},

		/**
		 * 3.x API
		 */
		////DeliveryServices
		{api.Version{Major: 3, Minor: 1}, http.MethodPost, `deliveryservices/?$`, deliveryservice.CreateV31, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2064315323},
		{api.Version{Major: 3, Minor: 1}, http.MethodPut, `deliveryservices/{id}/?$`, deliveryservice.UpdateV31, auth.PrivLevelOperations, Authenticated, nil, NoCache, 27665675673},

		// Acme account information
		{api.Version{Major: 3, Minor: 1}, http.MethodGet, `acme_accounts/?$`, acme.Read, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2034390561},
		{api.Version{Major: 3, Minor: 1}, http.MethodPost, `acme_accounts/?$`, acme.Create, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2034390562},
		{api.Version{Major: 3, Minor: 1}, http.MethodPut, `acme_accounts/?$`, acme.Update, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2034390563},
		{api.Version{Major: 3, Minor: 1}, http.MethodDelete, `acme_accounts/{provider}/{email}?$`, acme.Delete, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2034390564},

		// API Capability
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `api_capabilities/?$`, apicapability.GetAPICapabilitiesHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 28132065893},

		//ASNs
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22641723173},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 202048983},

		//ASN: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `asns/?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2738777223},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 29511986293},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 29994921883},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 26725247693},

		// Traffic Stats access
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 23195690283},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 24979979063},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `current_stats/?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 27854428933},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `caches/stats/?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 28132065883},

		//CacheGroup: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cachegroups/?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2230791103},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2129545463},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 229826653},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2278693653},

		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 20716441103},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandlerV31, auth.PrivLevelOperations, Authenticated, nil, NoCache, 25202404313},

		//CacheGroup Parameters: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cachegroupparameters/?$`, cachegroupparameter.ReadAllCacheGroupParameters, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2124497243},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cachegroupparameters/?$`, cachegroupparameter.AddCacheGroupParameters, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2124497253},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cachegroups/{id}/parameters/?$`, api.DeprecatedReadHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2124497233},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cachegroupparameters/{cachegroupID}/{parameterId}$`, api.DeprecatedDeleteHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2124497333},

		//Capabilities
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `capabilities/?$`, capabilities.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 20081353},

		//CDN
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/name/{name}/sslkeys/?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 22785817723},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2971852813},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/{name}/health/?$`, cdn.GetNameHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 21353481943},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/health/?$`, cdn.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 20853811343},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/domains/?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2269025603},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/routing$`, crstats.GetCDNRouting, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 267229823},

		//CDN: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2088049593},

		//CDN: queue updates
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2215159803},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cdns/dnsseckeys/generate?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2753363},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cdns/name/{name}/dnsseckeys?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2711042073},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/name/{name}/dnsseckeys/?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2790106093},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/dnsseckeys/refresh/?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 27719971163},

		//CDN: Monitoring: Traffic Monitor
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/{cdn}/configs/monitoring?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22408478923},

		//Database dumps
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2240166473},

		//Division: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `divisions/?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 20851815343},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2063691403},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2537138003},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 23253822373},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `logs/?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2483405503},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `logs/newcount/?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 24058330123},

		//Content invalidation jobs
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `jobs/?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 29667820413},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2167807763},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2861342263},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil, NoCache, 204509553},

		//Login
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 23926708213},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `user/logout/?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil, NoCache, 2434348253},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 24158860093},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `user/login/token/?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 2024088413},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `user/reset_password/?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil, NoCache, 22929146303},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `users/register/?$`, login.RegisterUser, auth.PrivLevelOperations, Authenticated, nil, NoCache, 23373},

		//ISO
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `osversions/?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2760886573},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `isos/?$`, iso.ISOs, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2760336573},

		//User: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `users/?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 24919299003},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2138099803},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2354334043},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `users/?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2762448163},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `user/current/?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 26107016143},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `user/current/?$`, user.ReplaceCurrent, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2203},

		//Parameter: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `parameters/?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22125542923},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 28739361153},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 26695108593},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2262771183},

		//Phys_Location: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `phys_locations/?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2204051823},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2227950213},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22464566483},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 256142213},

		//Ping
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `ping$`, ping.Handler, 0, NoAuth, nil, NoCache, 25556615973},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `vault/ping/?$`, ping.Vault, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 28840121143},

		//Profile: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `profiles/?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2687585893},

		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 284391723},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 25402115563},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22055944653},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `profiles/{id}/export/?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 201335173},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profiles/import/?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2061432083},

		// Copy Profile
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profiles/name/{new_profile}/copy/{existing_profile}`, profile.CopyProfileHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2061432093},

		//Region: CRUDs
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `regions/?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2100370853},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2223082243},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22883344883},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `regions/?$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22326267583},

		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `topologies/?$`, api.CreateHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 3871452221},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `topologies/?$`, api.ReadHandler(&topology.TOTopology{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 3871452222},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `topologies/?$`, api.UpdateHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 3871452223},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `topologies/?$`, api.DeleteHandler(&topology.TOTopology{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 3871452224},

		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `topologies/{name}/queue_update$`, topology.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 3205351748},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryserviceserver/?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 29461450333},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2297997883},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryserviceserver/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 25321845233},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 24281812063},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2331154113},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2801282533},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 23451212233},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2408752993},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/capacity/?$`, deliveryservice.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22314091103},
		//Serverchecks
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servercheck/?$`, servercheck.ReadServerCheck, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 27961129223},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servercheck/?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil, NoCache, 27642815683},

		// Servercheck Extensions
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servercheck/extensions$`, extensions.Create, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2804985993},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servercheck/extensions$`, extensions.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2834985993},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `servercheck/extensions/{id}$`, extensions.Delete, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2804982993},

		//Server Details
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servers/details/?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22612647143},

		//Server status
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2766638513},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servers/{id}/queue_update$`, server.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 21894713},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2384515993},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servers/{id-or-name}/update$`, server.UpdateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 143813233},

		//Server: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `servers/?$`, server.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 27209592853},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `servers/{id}$`, server.Update, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2586341033},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `servers/?$`, server.Create, auth.PrivLevelOperations, Authenticated, nil, NoCache, 22255580613},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `servers/{id}$`, server.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2923222333},

		//Server Capability
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2104073913},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20744707083},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2364150383},

		//Server Server Capabilities: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 28002318893},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22931668343},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20587140583},

		//Status: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `statuses/?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22449056563},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22079665043},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 23691236123},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2551113603},

		//System
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `system/info/?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2210474753},

		//Type: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `types/?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22267018233},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 288601153},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 25133081953},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 231757733},

		//About
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `about/?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 23175011663},

		//Coordinates
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2967007453},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2689261743},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 24281121573},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 23038498893},

		//CDN generic handlers:
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22303186213},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 23111789343},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 21605052893},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2276946573},

		//Delivery service requests
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservice_requests/?$`, dsrequest.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 26811639353},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservice_requests/?$`, dsrequest.Put, auth.PrivLevelPortal, Authenticated, nil, NoCache, 22499079183},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservice_requests/?$`, dsrequest.Post, auth.PrivLevelPortal, Authenticated, nil, NoCache, 293850393},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservice_requests/?$`, dsrequest.Delete, auth.PrivLevelPortal, Authenticated, nil, NoCache, 22969850253},

		//Delivery service request: Actions
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/assign$`, dsrequest.PutAssignment, auth.PrivLevelOperations, Authenticated, nil, NoCache, 27031602903},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/status$`, dsrequest.PutStatus, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2684150993},

		//Delivery service request comment: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservice_request_comments/?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 20326507373},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 2604878473},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 2272276723},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 2995046683},

		//Delivery service uri signing keys: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 22930785583},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2084663353},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 276489693},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2299254173},

		//Delivery Service Required Capabilities: CRUD
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 21585222273},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20968739923},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 24962893043},

		// Federations by CDN (the actual table for federation)
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/{name}/federations/?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2892250323},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cdns/{name}/federations/?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 29548942193},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2260654663},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 24428529023},

		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2729242813},

		//Origins
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2446492563},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 215677463},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20995616433},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2602732633},

		//Roles
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `roles/?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2870885833},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 26128974893},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2306524063},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 23567059823},

		//Delivery Services Regexes
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices_regexes/?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2055014533},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2774327633},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2127378003},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil, NoCache, 22483396913},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 22467316633},

		//ServiceCategories
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `service_categories/?$`, api.ReadHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 1085181543},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `service_categories/{name}/?$`, servicecategory.Update, auth.PrivLevelOperations, Authenticated, nil, NoCache, 306369141},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `service_categories/?$`, api.CreateHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 553713801},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `service_categories/{name}$`, api.DeleteHandler(&servicecategory.TOServiceCategory{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 1325382238},

		//StaticDNSEntries
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2289394773},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2424571113},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 26291482383},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 28460311323},

		//ProfileParameters
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `profiles/{id}/parameters/?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2764649753},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `profiles/name/{name}/parameters/?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22677378323},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil, NoCache, 23559455823},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2168187083},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `profileparameters/?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2506098053},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2288096933},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2242753},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil, NoCache, 20806108613},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2248395293},

		//Tenants
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `tenants/?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 26779678143},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20941314783},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2172480133},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2163655583},

		//CRConfig
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 29572736953},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2767168893},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 29699118293},

		// Federations
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federations/all/?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 210599863},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federations/?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil, NoCache, 2549549943},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `federations/?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 28940647423},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `federations/?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 220983233},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `federations/?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 22831825163},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `federations/{id}/deliveryservices/?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 26828635133},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federations/{id}/deliveryservices/?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2537730343},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 24174025703},

		// Federation Resolvers
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `federation_resolvers/?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 21343736613},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federation_resolvers/?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2566087593},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `federations/{id}/federation_resolvers/?$`, federations.AssignFederationResolversToFederationHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2566087603},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federations/{id}/federation_resolvers/?$`, federations.GetFederationFederationResolversHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2566087613},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `federation_resolvers/?$`, federation_resolvers.Delete, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 20013},

		// Federations Users
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `federations/{id}/users/?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 27793349303},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `federations/{id}/users/?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2940750153},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `federations/{id}/users/{userID}/?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 29491028823},

		////DeliveryServices
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22383172943},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/?$`, deliveryservice.CreateV30, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2064314323},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservices/{id}/?$`, deliveryservice.UpdateV30, auth.PrivLevelOperations, Authenticated, nil, NoCache, 27665675273},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `deliveryservices/{id}/safe/?$`, deliveryservice.UpdateSafe, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2472109313},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservices/{id}/?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2226420743},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers/eligible/?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2747615843},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLIDV15, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 21357729073},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 28728785833},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 29267343},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2534390513},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 22625010763},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 25304828243},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22027192113},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/urlkeys/?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2931971143},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `vault/bucket/{bucket}/key/{key}/values/?$`, vault.GetBucketKeyDeprecated, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 22205108013},

		//Delivery service LetsEncrypt
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/letsencrypt/?$`, deliveryservice.GenerateLetsEncryptCertificates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2534390523},
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `letsencrypt/dnsrecords/?$`, deliveryservice.GetDnsChallengeRecords, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2534390553},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `letsencrypt/autorenew/?$`, deliveryservice.RenewCertificates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2534390563},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/health/?$`, deliveryservice.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 22345901013},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `deliveryservices/{id}/routing$`, crstats.GetDSRouting, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 667339833},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `steering/{deliveryservice}/targets/?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 25696078243},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `steering/{deliveryservice}/targets/?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 23382163973},
		{api.Version{Major: 3, Minor: 0}, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 24386082953},
		{api.Version{Major: 3, Minor: 0}, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 22880215153},

		// Stats Summary
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `stats_summary/?$`, trafficstats.GetStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2804985983},
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `stats_summary/?$`, trafficstats.CreateStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2804915983},

		//Pattern based consistent hashing endpoint
		{api.Version{Major: 3, Minor: 0}, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2607550763},

		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `steering/?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil, NoCache, 21748524573},

		// Plugins
		{api.Version{Major: 3, Minor: 0}, http.MethodGet, `plugins/?$`, plugins.Get(d.Plugins), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2834985393},

		/**
		 * 2.x API
		 */
		// API Capability
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `api_capabilities/?$`, apicapability.GetAPICapabilitiesHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2813206589},

		//ASNs
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2264172317},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 20204898},

		//ASN: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `asns/?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 273877722},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2951198629},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2999492188},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2672524769},

		// Traffic Stats access
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2319569028},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2497997906},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `current_stats/?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2785442893},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `caches/stats/?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2813206588},

		//CacheGroup: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cachegroups/?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 223079110},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 212954546},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 22982665},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 227869365},

		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2071644110},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandlerV31, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2520240431},

		//CacheGroup Parameters: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cachegroupparameters/?$`, cachegroupparameter.ReadAllCacheGroupParameters, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 212449724},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cachegroupparameters/?$`, cachegroupparameter.AddCacheGroupParameters, auth.PrivLevelOperations, Authenticated, nil, NoCache, 212449725},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cachegroups/{id}/parameters/?$`, api.DeprecatedReadHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 212449723},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cachegroupparameters/{cachegroupID}/{parameterId}$`, api.DeprecatedDeleteHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), auth.PrivLevelOperations, Authenticated, nil, NoCache, 212449733},

		//Capabilities
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `capabilities/?$`, capabilities.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2008135},

		//CDN
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/name/{name}/sslkeys/?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2278581772},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 297185281},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/{name}/health/?$`, cdn.GetNameHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2135348194},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/health/?$`, cdn.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2085381134},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/domains/?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 226902560},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/routing$`, crstats.GetCDNRouting, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 26722982},

		//CDN: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil, NoCache, 208804959},

		//CDN: queue updates
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil, NoCache, 221515980},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cdns/dnsseckeys/generate?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 275336},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cdns/name/{name}/dnsseckeys?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 271104207},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/name/{name}/dnsseckeys/?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 279010609},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/dnsseckeys/refresh/?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2771997116},

		//CDN: Monitoring: Traffic Monitor
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/{cdn}/configs/monitoring?$`, crconfig.SnapshotGetMonitoringLegacyHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2240847892},

		//Database dumps
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 224016647},

		//Division: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `divisions/?$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2085181534},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 206369140},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 253713800},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2325382237},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `logs/?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 248340550},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `logs/newcount/?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2405833012},

		//Content invalidation jobs
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `jobs/?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2966782041},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil, NoCache, 216780776},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil, NoCache, 286134226},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil, NoCache, 20450955},

		//Login
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 2392670821},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `user/logout/?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil, NoCache, 243434825},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 2415886009},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `user/login/token/?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil, NoCache, 202408841},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `user/reset_password/?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil, NoCache, 2292914630},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `users/register/?$`, login.RegisterUser, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2337},

		//ISO
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `osversions/?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 276088657},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `isos/?$`, iso.ISOs, auth.PrivLevelOperations, Authenticated, nil, NoCache, 276033657},

		//User: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `users/?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2491929900},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 213809980},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 235433404},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `users/?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 276244816},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `user/current/?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2610701614},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `user/current/?$`, user.ReplaceCurrent, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 220},

		//Parameter: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `parameters/?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2212554292},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2873936115},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2669510859},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 226277118},

		//Phys_Location: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `phys_locations/?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 220405182},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 222795021},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2246456648},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 25614221},

		//Ping
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `ping$`, ping.Handler, 0, NoAuth, nil, NoCache, 2555661597},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `vault/ping/?$`, ping.Vault, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2884012114},

		//Profile: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `profiles/?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 268758589},

		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 28439172},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2540211556},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2205594465},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `profiles/{id}/export/?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 20133517},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profiles/import/?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 206143208},

		// Copy Profile
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profiles/name/{new_profile}/copy/{existing_profile}`, profile.CopyProfileHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 206143209},

		//Region: CRUDs
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `regions/?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 210037085},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 222308224},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2288334488},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `regions/?$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2232626758},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryserviceserver/?$`, dsserver.ReadDSSHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2946145033},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 229799788},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryserviceserver/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2532184523},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2428181206},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 233115411},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 280128253},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2345121223},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil, NoCache, 240875299},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/capacity/?$`, deliveryservice.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2231409110},
		//Serverchecks
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servercheck/?$`, servercheck.ReadServerCheck, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2796112922},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servercheck/?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil, NoCache, 2764281568},

		// Servercheck Extensions
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servercheck/extensions$`, extensions.Create, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 280498599},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servercheck/extensions$`, extensions.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 283498599},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `servercheck/extensions/{id}$`, extensions.Delete, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 280498299},

		//Server Details
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servers/details/?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2261264714},

		//Server status
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 276663851},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servers/{id}/queue_update$`, server.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2189471},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 238451599},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servers/{id-or-name}/update$`, server.UpdateHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 14381323},

		//Server: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `servers/?$`, server.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2720959285},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `servers/{id}$`, server.Update, auth.PrivLevelOperations, Authenticated, nil, NoCache, 258634103},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `servers/?$`, server.Create, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2225558061},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `servers/{id}$`, server.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 292322233},

		//Server Capability
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 210407391},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2074470708},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 236415038},

		//Server Server Capabilities: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2800231889},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2293166834},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2058714058},

		//Status: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `statuses/?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2244905656},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2207966504},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2369123612},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 255111360},

		//System
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `system/info/?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 221047475},

		//Type: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `types/?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2226701823},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 28860115},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2513308195},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 23175773},

		//About
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `about/?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2317501166},

		//Coordinates
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 296700745},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 268926174},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2428112157},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2303849889},

		//CDN generic handlers:
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2230318621},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2311178934},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2160505289},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 227694657},

		//Delivery service requests
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservice_requests/?$`, dsrequest.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2681163935},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservice_requests/?$`, dsrequest.Put, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2249907918},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservice_requests/?$`, dsrequest.Post, auth.PrivLevelPortal, Authenticated, nil, NoCache, 29385039},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservice_requests/?$`, dsrequest.Delete, auth.PrivLevelPortal, Authenticated, nil, NoCache, 2296985025},

		//Delivery service request: Actions
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/assign$`, dsrequest.PutAssignment, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2703160290},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservice_requests/{id}/status$`, dsrequest.PutStatus, auth.PrivLevelPortal, Authenticated, nil, NoCache, 268415099},

		//Delivery service request comment: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservice_request_comments/?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2032650737},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 260487847},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 227227672},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil, NoCache, 299504668},

		//Delivery service uri signing keys: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2293078558},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 208466335},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 27648969},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 229925417},

		//Delivery Service Required Capabilities: CRUD
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2158522227},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2096873992},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2496289304},

		// Federations by CDN (the actual table for federation)
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/{name}/federations/?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 289225032},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cdns/{name}/federations/?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2954894219},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 226065466},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2442852902},

		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 272924281},

		//Origins
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 244649256},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 21567746},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2099561643},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 260273263},

		//Roles
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `roles/?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 287088583},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2612897489},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 230652406},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2356705982},

		//Delivery Services Regexes
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices_regexes/?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 205501453},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 277432763},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil, NoCache, 212737800},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2248339691},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2246731663},

		//StaticDNSEntries
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 228939477},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 242457111},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2629148238},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2846031132},

		//ProfileParameters
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `profiles/{id}/parameters/?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 276464975},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `profiles/name/{name}/parameters/?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2267737832},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2355945582},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil, NoCache, 216818708},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `profileparameters/?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 250609805},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 228809693},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil, NoCache, 224275},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2080610861},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 224839529},

		//Tenants
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `tenants/?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2677967814},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 2094131478},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 217248013},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 216365558},

		//CRConfig
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2957273695},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 276716889},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2969911829},

		// Federations
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federations/all/?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 21059986},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federations/?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil, NoCache, 254954994},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `federations/?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 2894064742},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `federations/?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 22098323},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `federations/?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil, NoCache, 2283182516},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `federations/{id}/deliveryservices/?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2682863513},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federations/{id}/deliveryservices/?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 253773034},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2417402570},

		// Federation Resolvers
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `federation_resolvers/?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2134373661},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federation_resolvers/?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 256608759},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `federations/{id}/federation_resolvers/?$`, federations.AssignFederationResolversToFederationHandler, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 256608760},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federations/{id}/federation_resolvers/?$`, federations.GetFederationFederationResolversHandler, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 256608761},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `federation_resolvers/?$`, federation_resolvers.Delete, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2001},

		// Federations Users
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `federations/{id}/users/?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2779334930},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `federations/{id}/users/?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 294075015},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `federations/{id}/users/{userID}/?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2949102882},

		////DeliveryServices
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2238317294},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/?$`, deliveryservice.CreateV15, auth.PrivLevelOperations, Authenticated, nil, NoCache, 206431432},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservices/{id}/?$`, deliveryservice.UpdateV15, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2766567527},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `deliveryservices/{id}/safe/?$`, deliveryservice.UpdateSafe, auth.PrivLevelOperations, Authenticated, nil, NoCache, 247210931},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservices/{id}/?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil, NoCache, 222642074},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/servers/eligible/?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 274761584},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLIDV15, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2135772907},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2872878583},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2926734},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 253439051},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2262501076},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil, NoCache, 2530482824},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2202719211},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/urlkeys/?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 293197114},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `vault/bucket/{bucket}/key/{key}/values/?$`, vault.GetBucketKeyDeprecated, auth.PrivLevelAdmin, Authenticated, nil, NoCache, 2220510801},

		//Delivery service LetsEncrypt
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `deliveryservices/sslkeys/generate/letsencrypt/?$`, deliveryservice.GenerateLetsEncryptCertificates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 253439052},
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `letsencrypt/dnsrecords/?$`, deliveryservice.GetDnsChallengeRecords, auth.PrivLevelOperations, Authenticated, nil, NoCache, 253439055},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `letsencrypt/autorenew/?$`, deliveryservice.RenewCertificates, auth.PrivLevelOperations, Authenticated, nil, NoCache, 253439056},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/health/?$`, deliveryservice.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2234590101},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `deliveryservices/{id}/routing$`, crstats.GetDSRouting, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 66733983},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `steering/{deliveryservice}/targets/?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 2569607824},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `steering/{deliveryservice}/targets/?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 2338216397},
		{api.Version{Major: 2, Minor: 0}, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 2438608295},
		{api.Version{Major: 2, Minor: 0}, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil, NoCache, 2288021515},

		// Stats Summary
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `stats_summary/?$`, trafficstats.GetStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 280498598},
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `stats_summary/?$`, trafficstats.CreateStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 280491598},

		//Pattern based consistent hashing endpoint
		{api.Version{Major: 2, Minor: 0}, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 260755076},

		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `steering/?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil, NoCache, 2174852457},

		// Plugins
		{api.Version{Major: 2, Minor: 0}, http.MethodGet, `plugins/?$`, plugins.Get(d.Plugins), auth.PrivLevelReadOnly, Authenticated, nil, NoCache, 283498539},
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

	// check for unknown route IDs in cdn.conf
	disabledRoutes := GetRouteIDMap(d.DisabledRoutes)
	unknownRouteIDs := []string{}
	for _, routeMap := range []map[int]struct{}{disabledRoutes} {
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
	rawRoutes := []RawRoute{}

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
		api.WriteAndLogErr(w, r, bytes)
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
		api.WriteAndLogErr(w, r, bytes)
	}
}

type root struct {
	Handler http.Handler
}

func (root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	alerts := tc.CreateAlerts(tc.ErrorLevel, fmt.Sprintf("The requested path '%s' does not exist.", r.URL.Path))
	api.WriteAlerts(w, r, http.StatusNotFound, alerts)
}

// rootHandler returns the / handler for the service, which simply returns a "not found" response.
func rootHandler(d ServerData) http.Handler {
	return root{}
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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) >= 3 {
		version, err := strconv.ParseFloat(pathParts[2], 64)
		if err == nil && version >= 2 { // do not default to Perl for versions over 2.x
			api.WriteRespAlertNotFound(w, r)
			return
		}
	}

	m.ReqChan <- struct{}{}
	defer func() { <-m.ReqChan }()
	m.Handler.ServeHTTP(w, r)
}
