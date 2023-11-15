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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/acme"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/apicapability"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/apitenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/asn"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cachegroupparameter"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cachesstats"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/capabilities"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cdn"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cdn_lock"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cdnfederation"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cdni"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cdnnotification"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/coordinate"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/crconfig"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/crstats"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbdump"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice/consistenthash"
	dsrequest "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice/request"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice/request/comment"
	dsserver "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice/servers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservicerequests"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservicesregexes"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/division"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/federation_resolvers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/federations"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/invalidationjobs"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/iso"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/login"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/logs"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/origin"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/physlocation"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/ping"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/plugins"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/profile"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/profileparameter"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/region"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/role"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/servercapability"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/servercheck"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/servercheck/extensions"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/servicecategory"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/steering"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/steeringtargets"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/systeminfo"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/topology"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficstats"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/types"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/urisigning"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/user"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/vault"

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
func Routes(d ServerData) ([]Route, http.Handler, error) {
	proxyHandler := rootHandler(d)

	routes := []Route{
		// 1.1 and 1.2 routes are simply a Go replacement for the equivalent Perl route. They may or may not conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).
		// 1.3 routes exist only in Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		// 2.x routes exist only in Go. There is NO equivalent Perl route. They should conform with the API guidelines (https://cwiki.apache.org/confluence/display/TC/API+Guidelines).

		// NOTE: Route IDs are immutable and unique. DO NOT change the ID of an existing Route; otherwise, existing
		// configurations may break. New Route IDs can be any integer between 0 and 2147483647 (inclusive), as long as
		// it's unique.
		/**
		 * 5.x API
		 */

		// CDNI integration
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `OC/FCI/advertisement/?$`, Handler: cdni.GetCapabilities, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 5413577290771},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/?$`, Handler: cdni.PutConfiguration, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 5413577290781},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/{host}$`, Handler: cdni.PutHostConfiguration, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 5413577290791},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/request/{id}/{approved}$`, Handler: cdni.PutConfigurationResponse, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"CDNI-ADMIN:READ", "CDNI-ADMIN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 5413577290801},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `OC/CI/configuration/requests/?$`, Handler: cdni.GetRequests, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"CDNI-ADMIN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 5413577290811},

		// SSL Keys
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `sslkey_expirations/?$`, Handler: deliveryservice.GetSSlKeyExpirationInformation, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"SSL-KEY-EXPIRATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 413577290751},

		// CDN lock
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdn_locks/?$`, Handler: cdn_lock.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41343905611},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdn_locks/?$`, Handler: cdn_lock.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-LOCK:CREATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41343905621},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdn_locks/?$`, Handler: cdn_lock.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-LOCK:DELETE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41343905641},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `acme_accounts/providers?$`, Handler: acme.ReadProviders, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40343905651},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/acme/?$`, Handler: deliveryservice.GenerateAcmeCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:UPDATE", "ACME:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 25343905761},

		// ACME account information
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `acme_accounts/?$`, Handler: acme.Read, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40343905611},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `acme_accounts/?$`, Handler: acme.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:CREATE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40343905621},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `acme_accounts/?$`, Handler: acme.Update, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:UPDATE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40343905631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `acme_accounts/{provider}/{email}?$`, Handler: acme.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:DELETE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40343905641},

		//Delivery service ACME
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{xmlid}/sslkeys/renew$`, Handler: deliveryservice.RenewAcmeCertificate, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ", "DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:CREATE", "DS-SECURITY-KEY:UPDATE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 25343905731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `acme_autorenew/?$`, Handler: deliveryservice.RenewCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ", "DS-SECURITY-KEY:UPDATE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 25343905741},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `async_status/{id}$`, Handler: api.GetAsyncStatus, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASYNC-STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 25343905751},

		//ASNs
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `asns/?$`, Handler: asn.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ASN:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47387772231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `asns/?$`, Handler: asn.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:UPDATE", "ASN:READ", "CACHE-GROUP:UPDATE", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 426417231731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `asns/{id}$`, Handler: asn.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:UPDATE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 495119862931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `asns/?$`, Handler: asn.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:CREATE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 499949218831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `asns/?$`, Handler: asn.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:DELETE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4020489831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `asns/{id}$`, Handler: asn.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:DELETE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 467252476931},

		// Traffic Stats access
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_stats`, Handler: trafficstats.GetDSStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 431956902831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cache_stats`, Handler: trafficstats.GetCacheStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 449799790631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `current_stats/?$`, Handler: trafficstats.GetCurrentStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 478544289331},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `caches/stats/?$`, Handler: cachesstats.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 481320658831},

		//CacheGroup: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cachegroups/?$`, Handler: cachegroup.GetCacheGroup, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42307911031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `cachegroups/{id}$`, Handler: cachegroup.UpdateCacheGroup, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:UPDATE", "CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41295454631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/?$`, Handler: cachegroup.CreateCacheGroup, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:CREATE", "CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4298266531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cachegroups/{id}$`, Handler: cachegroup.DeleteCacheGroup, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:DELETE", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42786936531},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/queue_update$`, Handler: cachegroup.QueueUpdates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:READ", "CDN:READ", "SERVER:READ", "SERVER:QUEUE"}, Authenticated: Authenticated, Middlewares: nil, ID: 407164411031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/deliveryservices/?$`, Handler: cachegroup.DSPostHandlerV40, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:UPDATE", "DELIVERY-SERVICE:UPDATE", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 452024043131},

		//CDN
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/sslkeys/?$`, Handler: cdn.GetSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 427858177231},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/capacity$`, Handler: cdn.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49718528131},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/health/?$`, Handler: cdn.GetNameHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 413534819431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/health/?$`, Handler: cdn.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 408538113431},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/domains/?$`, Handler: cdn.DomainsHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42690256031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/routing$`, Handler: crstats.GetCDNRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4672298231},

		//CDN: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}$`, Handler: cdn.DeleteName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40880495931},

		//CDN: queue updates
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdns/{id}/queue_update$`, Handler: cdn.Queue, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42151598031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdns/dnsseckeys/generate?$`, Handler: cdn.CreateDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:CREATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47533631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}/dnsseckeys?$`, Handler: cdn.DeleteDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:DELETE", "CDN:UPDATE", "DELIVERY-SERVICE:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47110420731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/dnsseckeys/?$`, Handler: cdn.GetDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47901060931},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `cdns/dnsseckeys/refresh/?$`, Handler: cdn.RefreshDNSSECKeysV4, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DNS-SEC:UPDATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 477199711631},

		//CDN: Monitoring: Traffic Monitor
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/configs/monitoring?$`, Handler: crconfig.SnapshotGetMonitoringHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"MONITOR-CONFIG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424084789231},

		//Database dumps
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `dbdump/?`, Handler: dbdump.DBDump, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DBDUMP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42401664731},

		//Division: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `divisions/?$`, Handler: division.GetDivisions, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 408518153431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `divisions/{id}$`, Handler: division.UpdateDivision, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:UPDATE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40636914031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `divisions/?$`, Handler: division.CreateDivision, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:CREATE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45371380031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `divisions/{id}$`, Handler: division.DeleteDivision, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:DELETE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 432538223731},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `logs/?$`, Handler: logs.Getv40, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"LOG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44834055031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `logs/newcount/?$`, Handler: logs.GetNewCount, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"LOG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 440583301231},

		//Content invalidation jobs
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `jobs/?$`, Handler: api.ReadHandler(&invalidationjobs.InvalidationJobV4{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 496678204131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `jobs/?$`, Handler: invalidationjobs.DeleteV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:DELETE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41678077631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `jobs/?$`, Handler: invalidationjobs.UpdateV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:UPDATE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48613422631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `jobs/?`, Handler: invalidationjobs.CreateV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:CREATE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4045095531},

		//Login
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `user/login/?$`, Handler: login.LoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 439267082131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `user/logout/?$`, Handler: login.LogoutHandler(d.Config.Secrets[0]), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 44343482531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `user/login/oauth/?$`, Handler: login.OauthLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 441588600931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `user/login/token/?$`, Handler: login.TokenLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 40240884131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `user/reset_password/?$`, Handler: login.ResetPassword(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 429291463031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `users/register/?$`, Handler: login.RegisterUser, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:CREATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 433731},

		//ISO
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `osversions/?$`, Handler: iso.GetOSVersions, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ISO:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47608865731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `isos/?$`, Handler: iso.ISOs, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ISO:GENERATE", "ISO:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47603365731},

		//User: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `users/?$`, Handler: user.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 449192990031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `users/{id}$`, Handler: user.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41380998031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `users/{id}$`, Handler: user.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:UPDATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43543340431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `users/?$`, Handler: user.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:CREATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47624481631},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `user/current/?$`, Handler: user.Current, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 461070161431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `user/current/?$`, Handler: user.ReplaceCurrentV4, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 42031},

		//Parameter: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `parameters/?$`, Handler: parameter.GetParameters, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 421255429231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `parameters/{id}$`, Handler: parameter.UpdateParameter, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:UPDATE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 487393611531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `parameters/?$`, Handler: parameter.CreateParameter, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:CREATE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 466951085931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `parameters/{id}$`, Handler: parameter.DeleteParameter, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:DELETE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42627711831},

		//Phys_Location: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `phys_locations/?$`, Handler: physlocation.GetPhysLocation, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42040518231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `phys_locations/{id}$`, Handler: physlocation.UpdatePhysLocation, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:UPDATE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42279502131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `phys_locations/?$`, Handler: physlocation.CreatePhysLocation, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:CREATE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424645664831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `phys_locations/{id}$`, Handler: physlocation.DeletePhysLocation, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:DELETE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4561422131},

		//Ping
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `ping$`, Handler: ping.Handler, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 455566159731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `vault/ping/?$`, Handler: ping.Vault, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TRAFFIC-VAULT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 488401211431},

		//Profile: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `profiles/?$`, Handler: profile.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46875858931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `profiles/{id}$`, Handler: profile.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4843917231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profiles/?$`, Handler: profile.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 454021155631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `profiles/{id}$`, Handler: profile.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:DELETE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 420559446531},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/export/?$`, Handler: profile.ExportProfileHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4013351731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profiles/import/?$`, Handler: profile.ImportProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PARAMETER:CREATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40614320831},

		// Copy Profile
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{new_profile}/copy/{existing_profile}`, Handler: profile.CopyProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40614320931},

		//Region: CRUDs
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `regions/?$`, Handler: region.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41003708531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `regions/{id}$`, Handler: region.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:UPDATE", "REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42230822431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `regions/?$`, Handler: region.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:CREATE", "REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 428833448831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `regions/?$`, Handler: region.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 423262675831},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `topologies/?$`, Handler: topology.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:CREATE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48714522211},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `topologies/?$`, Handler: topology.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48714522221},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `topologies/?$`, Handler: topology.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:UPDATE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48714522231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `topologies/?$`, Handler: topology.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:DELETE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48714522241},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `topologies/{name}/queue_update$`, Handler: topology.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "TOPOLOGY:READ", "SERVER:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42053517481},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryserviceserver/?$`, Handler: dsserver.ReadDSSHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 494614503331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryserviceserver$`, Handler: dsserver.GetReplaceHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "SERVER:UPDATE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42979978831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryserviceserver/{dsid}/{serverid}`, Handler: dsserver.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE", "SERVER:READ", "SERVER:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 453218452331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xml_id}/servers$`, Handler: dsserver.GetCreateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "SERVER:UPDATE", "DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 442818120631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `servers/{id}/deliveryservices$`, Handler: api.ReadHandler(&dsserver.TODSSDeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43311541131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/deliveryservices$`, Handler: server.AssignDeliveryServicesToServerHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "DELIVERY-SERVICE:UPDATE", "SERVER:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 48012825331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers$`, Handler: dsserver.GetReadAssigned, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "CDN:READ", "TYPE:READ", "PROFILE:READ", "DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 434512122331},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/capacity/?$`, Handler: deliveryservice.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 423140911031},
		//Serverchecks
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `servercheck/?$`, Handler: servercheck.ReadServerCheck, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 479611292231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servercheck/?$`, Handler: servercheck.CreateUpdateServercheck, RequiredPrivLevel: auth.PrivLevelInvalid, RequiredPermissions: []string{"SERVER-CHECK:CREATE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 476428156831},

		// Servercheck Extensions
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servercheck/extensions$`, Handler: extensions.Create, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:CREATE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48049859931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `servercheck/extensions$`, Handler: extensions.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48349859931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `servercheck/extensions/{id}$`, Handler: extensions.Delete, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:DELETE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48049829931},

		//Server status
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}/status$`, Handler: server.UpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47666385131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/queue_update$`, Handler: server.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 418947131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `servers/{host_name}/update_status$`, Handler: server.GetServerUpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43845159931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servers/{id-or-name}/update$`, Handler: server.UpdateHandlerV4, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4438132331},

		//Server: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `servers/?$`, Handler: api.Wrap(server.Read, nil, nil), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 472095928531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}$`, Handler: api.Wrap(server.Update, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45863410331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `servers/?$`, Handler: api.Wrap(server.Create, nil, nil), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:CREATE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 422555806131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `servers/{id}$`, Handler: api.Wrap(server.Delete, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:DELETE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49232223331},

		//Server Capability
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `server_capabilities$`, Handler: servercapability.GetServerCapabilityV5, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41040739131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `server_capabilities$`, Handler: servercapability.CreateServerCapabilityV5, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:CREATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 407447070831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `server_capabilities$`, Handler: servercapability.UpdateServerCapabilityV5, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:UPDATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 425437701091},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `server_capabilities$`, Handler: servercapability.DeleteServerCapability, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:DELETE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43641503831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `multiple_servers_capabilities/?$`, Handler: server.AssignMultipleServersCapabilities, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:READ", "SERVER:CREATE", "SERVER-CAPABILITY:READ", "SERVER-CAPABILITY:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 407924192581},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `multiple_servers_capabilities/?$`, Handler: server.DeleteMultipleServersCapabilities, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:READ", "SERVER:DELETE", "SERVER-CAPABILITY:READ", "SERVER-CAPABILITY:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 407924192781},

		//Server Server Capabilities: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `server_server_capabilities/?$`, Handler: api.ReadHandler(&server.TOServerServerCapabilityV5{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 480023188931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `server_server_capabilities/?$`, Handler: api.CreateHandler(&server.TOServerServerCapabilityV5{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 429316683431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `server_server_capabilities/?$`, Handler: api.DeleteHandler(&server.TOServerServerCapabilityV5{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 405871405831},

		//Status: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `statuses/?$`, Handler: api.ReadHandler(&status.TOStatusV5{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424490565631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `statuses/{id}$`, Handler: api.UpdateHandler(&status.TOStatusV5{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:UPDATE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 420796650431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `statuses/?$`, Handler: api.CreateHandler(&status.TOStatusV5{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:CREATE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 436912361231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `statuses/{id}$`, Handler: api.DeleteHandler(&status.TOStatusV5{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:DELETE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45511136031},

		//System
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `system/info/?$`, Handler: systeminfo.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 42104747531},

		//Type: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `types/?$`, Handler: types.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 422670182331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `types/{id}$`, Handler: types.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:UPDATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4886011531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `types/?$`, Handler: types.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:CREATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 451330819531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `types/{id}$`, Handler: types.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:DELETE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4317577331},

		//About
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `about/?$`, Handler: about.Handler(), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 431750116631},

		//Coordinates
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `coordinates/?$`, Handler: coordinate.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49670074531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `coordinates/?$`, Handler: coordinate.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:UPDATE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46892617431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `coordinates/?$`, Handler: coordinate.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:CREATE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 442811215731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `coordinates/?$`, Handler: coordinate.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:DELETE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 430384988931},

		//CDN notification
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdn_notifications/?$`, Handler: cdnnotification.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 22212245141},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdn_notifications/?$`, Handler: cdnnotification.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 27652235131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdn_notifications/?$`, Handler: cdnnotification.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 27224118511},

		//CDN generic handlers:
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/?$`, Handler: cdn.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 423031862131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `cdns/{id}$`, Handler: cdn.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 431117893431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdns/?$`, Handler: cdn.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:READ", "CDN:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 416050528931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{id}$`, Handler: cdn.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:DELETE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42769465731},

		//Delivery service requests
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-REQUEST:READ", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 468116393531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Put, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424990791831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Post, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:CREATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4938503931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Delete, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:DELETE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 429698502531},

		//Delivery service request: Actions
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/{id}/assign$`, Handler: dsrequest.GetAssignment, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-REQUEST:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 470316029041},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/assign$`, Handler: dsrequest.PutAssignment, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DS-REQUEST:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 470316029031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/{id}/status$`, Handler: dsrequest.GetStatus, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46841509941},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/status$`, Handler: dsrequest.PutStatus, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DS-REQUEST:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46841509931},

		//Delivery service request comment: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_request_comments/?$`, Handler: comment.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-REQUEST:READ", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 403265073731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_request_comments/?$`, Handler: comment.Update, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46048784731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_request_comments/?$`, Handler: comment.Create, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42722767231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_request_comments/?$`, Handler: comment.Delete, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49950466831},

		//Delivery service uri signing keys: CRUD
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.GetURIsignkeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 429307855831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40846633531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4764896931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.RemoveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42992541731},

		// Federations by CDN (the actual table for federation)
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/federations/?$`, Handler: api.Wrap(cdnfederation.Read, []string{"name"}, nil), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48922503231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/federations/?$`, Handler: api.Wrap(cdnfederation.Create, nil, nil), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:CREATE", "FEDERATION:READ, CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 495489421931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/federations/{id}/?$`, Handler: api.Wrap(cdnfederation.ReadID, []string{"name"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48922503232},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `cdns/{name}/federations/{id}$`, Handler: api.Wrap(cdnfederation.Update, nil, []string{"id"}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42606546631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{name}/federations/{id}$`, Handler: api.Wrap(cdnfederation.Delete, nil, []string{"id"}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:DELETE", "FEDERATION:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 444285290231},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/dnsseckeys/ksk/generate$`, Handler: cdn.GenerateKSK, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:CREATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47292428131},

		//Origins
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `origins/?$`, Handler: origin.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ORIGIN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44464925631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `origins/?$`, Handler: origin.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:UPDATE", "ORIGIN:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4156774631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `origins/?$`, Handler: origin.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:CREATE", "ORIGIN:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 409956164331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `origins/?$`, Handler: origin.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:DELETE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 46027326331},

		//Roles
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `roles/?$`, Handler: role.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48708858331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `roles/?$`, Handler: role.Update, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:UPDATE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 461289748931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `roles/?$`, Handler: role.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:CREATE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43065240631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `roles/?$`, Handler: role.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:DELETE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 435670598231},

		//Delivery Services Regexes
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices_regexes/?$`, Handler: deliveryservicesregexes.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40550145331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.DSGet, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47743276331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.Post, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41273780031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Put, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424833969131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 424673166331},

		//ServiceCategories
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `service_categories/?$`, Handler: servicecategory.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40851815431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `service_categories/{name}/?$`, Handler: servicecategory.UpdateServiceCategory, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:UPDATE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4063691411},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `service_categories/?$`, Handler: servicecategory.CreateServiceCategory, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:CREATE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4537138011},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `service_categories/{name}$`, Handler: servicecategory.DeleteServiceCategory, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:DELETE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43253822381},

		//StaticDNSEntries
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `staticdnsentries/?$`, Handler: staticdnsentry.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42893947731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `staticdnsentries/?$`, Handler: staticdnsentry.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:UPDATE", "STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44245711131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `staticdnsentries/?$`, Handler: staticdnsentry.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:CREATE", "STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 462914823831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `staticdnsentries/?$`, Handler: staticdnsentry.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:DELETE", "STATIC-DN:READ", "DELIVERY-SERVICE:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 484603113231},

		//ProfileParameters
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.GetProfileID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47646497531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.GetProfileName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 426773783231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.PostProfileParamsByName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 435594558231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.PostProfileParamsByID, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41681870831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `parameterprofile/?$`, Handler: profileparameter.PostParamProfile, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 408061086131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profileparameter/?$`, Handler: profileparameter.PostProfileParam, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ", "PROFILE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42427531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `profileparameters/?$`, Handler: profileparameter.GetProfileParameter, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45060980531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `profileparameters/?$`, Handler: profileparameter.CreateProfileParameter, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ", "PROFILE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42880969331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `profileparameters/{profileId}/{parameterId}$`, Handler: profileparameter.DeleteProfileParameter, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42483952931},

		//Tenants
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `tenants/?$`, Handler: apitenant.GetTenant, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 467796781431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `tenants/{id}$`, Handler: apitenant.UpdateTenant, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:UPDATE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 409413147831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `tenants/?$`, Handler: apitenant.CreateTenant, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:CREATE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41724801331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `tenants/{id}$`, Handler: apitenant.DeleteTenant, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:DELETE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41636555831},

		//CRConfig
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/?$`, Handler: crconfig.SnapshotGetHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 495727369531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/new/?$`, Handler: crconfig.Handler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47671688931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `snapshot/?$`, Handler: crconfig.SnapshotHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-SNAPSHOT:CREATE", "CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 496991182931},

		// Federations
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federations/all/?$`, Handler: federations.GetAll, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4105998631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federations/?$`, Handler: federations.Get, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45495499431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `federations/?$`, Handler: federations.AddFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:CREATE", "FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 489406474231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `federations/?$`, Handler: federations.RemoveFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4209832331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `federations/?$`, Handler: federations.ReplaceFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE", "FEDERATION-RESOLVER:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 428318251631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/deliveryservices/?$`, Handler: federations.PostDSes, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "DELIVERY-SERVICE:UPDATE", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 468286351331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/deliveryservices/?$`, Handler: api.ReadHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45377303431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/deliveryservices/{dsID}/?$`, Handler: api.DeleteHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "DELIVERY-SERVICE:UPDATE", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 441740257031},

		// Federation Resolvers
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:CREATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 413437366131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45660875931},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.AssignFederationResolversToFederationHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "FEDERATION-RESOLVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45660876031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.GetFederationFederationResolversHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "FEDERATION-RESOLVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45660876131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 400131},

		// Federations Users
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/users/?$`, Handler: federations.PostUsers, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "USER:READ", "FEDERATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 477933493031},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/users/?$`, Handler: api.ReadHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49407501531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/users/{userID}/?$`, Handler: api.DeleteHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 494910288231},

		////DeliveryServices
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/?$`, Handler: api.ReadHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 423831729431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/?$`, Handler: deliveryservice.CreateV50, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:CREATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40643153231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/?$`, Handler: deliveryservice.UpdateV50, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 476656756731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/safe/?$`, Handler: deliveryservice.UpdateSafe, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: []string{"DELIVERY-SERVICE-SAFE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44721093131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{id}/?$`, Handler: api.DeleteHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:DELETE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42264207431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers/eligible/?$`, Handler: deliveryservice.GetServersEligible, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "CACHE-GROUP:READ", "TYPE:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47476158431},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.GetSSLKeysByXMLID, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 413577290731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/add$`, Handler: deliveryservice.AddSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 487287858331},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.DeleteSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DELIVERY-SERVICE:READ", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 492673431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/?$`, Handler: deliveryservice.GenerateSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 45343905131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, Handler: deliveryservice.CopyURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 426250107631},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, Handler: deliveryservice.GenerateURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 453048282431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{name}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 420271921131},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/xmlId/{name}/urlkeys/?$`, Handler: deliveryservice.DeleteURLKeysByName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 420271921141},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49319711431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{id}/urlkeys/?$`, Handler: deliveryservice.DeleteURLKeysByID, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 49319711441},

		//Delivery service LetsEncrypt
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/letsencrypt/?$`, Handler: deliveryservice.GenerateLetsEncryptCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45343905231},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `letsencrypt/dnsrecords/?$`, Handler: deliveryservice.GetDnsChallengeRecords, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45343905531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `letsencrypt/autorenew/?$`, Handler: deliveryservice.RenewCertificatesDeprecated, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 45343905631},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/health/?$`, Handler: deliveryservice.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 423459010131},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/routing$`, Handler: crstats.GetDSRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4673398331},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 456960782431},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:CREATE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 433821639731},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPut, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:UPDATE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 443860829531},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodDelete, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:DELETE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 428802151531},

		// Stats Summary
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `stats_summary/?$`, Handler: trafficstats.GetStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48049859831},
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `stats_summary/?$`, Handler: trafficstats.CreateStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:CREATE", "STAT:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48049159831},

		//Pattern based consistent hashing endpoint
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodPost, Path: `consistenthash/?$`, Handler: consistenthash.Post, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46075507631},

		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `steering/?$`, Handler: steering.Get, RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 417485245731},

		// Plugins
		{Version: api.Version{Major: 5, Minor: 0}, Method: http.MethodGet, Path: `plugins/?$`, Handler: plugins.Get(d.Plugins), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PLUGIN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48349853931},

		/**
		 * 4.x API
		 */
		// Assign Multiple Servers to a capability and assign multiple capabilities to a server
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodPost, Path: `multiple_servers_capabilities/?$`, Handler: server.AssignMultipleServersCapabilities, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:READ", "SERVER:CREATE", "SERVER-CAPABILITY:READ", "SERVER-CAPABILITY:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40792419258},
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodDelete, Path: `multiple_servers_capabilities/?$`, Handler: server.DeleteMultipleServersCapabilities, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:READ", "SERVER:DELETE", "SERVER-CAPABILITY:READ", "SERVER-CAPABILITY:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40792419278},

		//Server Capability
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodGet, Path: `server_capabilities$`, Handler: servercapability.GetServerCapability, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4104073912},
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodPost, Path: `server_capabilities$`, Handler: servercapability.CreateServerCapability, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:CREATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40744707025},
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodPut, Path: `server_capabilities$`, Handler: servercapability.UpdateServerCapability, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:UPDATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 425437701023},

		// CDNI integration
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `OC/FCI/advertisement/?$`, Handler: cdni.GetCapabilities, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 541357729077},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/?$`, Handler: cdni.PutConfiguration, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 541357729078},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/{host}$`, Handler: cdni.PutHostConfiguration, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDNI-CAPACITY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 541357729079},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `OC/CI/configuration/request/{id}/{approved}$`, Handler: cdni.PutConfigurationResponse, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"CDNI-ADMIN:READ", "CDNI-ADMIN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 541357729080},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `OC/CI/configuration/requests/?$`, Handler: cdni.GetRequests, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"CDNI-ADMIN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 541357729081},

		// SSL Keys
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `sslkey_expirations/?$`, Handler: deliveryservice.GetSSlKeyExpirationInformation, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"SSL-KEY-EXPIRATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41357729075},

		// CDN lock
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdn_locks/?$`, Handler: cdn_lock.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4134390561},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdn_locks/?$`, Handler: cdn_lock.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-LOCK:CREATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4134390562},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdn_locks/?$`, Handler: cdn_lock.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-LOCK:DELETE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4134390564},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `acme_accounts/providers?$`, Handler: acme.ReadProviders, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4034390565},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/acme/?$`, Handler: deliveryservice.GenerateAcmeCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:UPDATE", "ACME:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 2534390576},

		// ACME account information
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `acme_accounts/?$`, Handler: acme.Read, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4034390561},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `acme_accounts/?$`, Handler: acme.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:CREATE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4034390562},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `acme_accounts/?$`, Handler: acme.Update, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:UPDATE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4034390563},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `acme_accounts/{provider}/{email}?$`, Handler: acme.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ACME:DELETE", "ACME:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4034390564},

		//Delivery service ACME
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{xmlid}/sslkeys/renew$`, Handler: deliveryservice.RenewAcmeCertificate, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ", "DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:CREATE", "DS-SECURITY-KEY:UPDATE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 2534390573},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `acme_autorenew/?$`, Handler: deliveryservice.RenewCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ACME:READ", "DS-SECURITY-KEY:UPDATE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 2534390574},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `async_status/{id}$`, Handler: api.GetAsyncStatus, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASYNC-STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 2534390575},

		//ASNs
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `asns/?$`, Handler: api.UpdateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:UPDATE", "ASN:READ", "CACHE-GROUP:UPDATE", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42641723173},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `asns/?$`, Handler: api.DeleteHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:DELETE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 402048983},

		//ASN: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `asns/?$`, Handler: api.ReadHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ASN:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4738777223},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `asns/{id}$`, Handler: api.UpdateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:UPDATE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 49511986293},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `asns/?$`, Handler: api.CreateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:CREATE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 49994921883},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `asns/{id}$`, Handler: api.DeleteHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ASN:DELETE", "ASN:READ", "CACHE-GROUP:READ", "CACHE-GROUP:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 46725247693},

		// Traffic Stats access
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_stats`, Handler: trafficstats.GetDSStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43195690283},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cache_stats`, Handler: trafficstats.GetCacheStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44979979063},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `current_stats/?$`, Handler: trafficstats.GetCurrentStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47854428933},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `caches/stats/?$`, Handler: cachesstats.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48132065883},

		//CacheGroup: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cachegroups/?$`, Handler: api.ReadHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4230791103},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `cachegroups/{id}$`, Handler: api.UpdateHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:UPDATE", "CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4129545463},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/?$`, Handler: api.CreateHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:CREATE", "CACHE-GROUP:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 429826653},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cachegroups/{id}$`, Handler: api.DeleteHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:DELETE", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4278693653},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/queue_update$`, Handler: cachegroup.QueueUpdates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:READ", "CDN:READ", "SERVER:READ", "SERVER:QUEUE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40716441103},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/deliveryservices/?$`, Handler: cachegroup.DSPostHandlerV40, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CACHE-GROUP:UPDATE", "DELIVERY-SERVICE:UPDATE", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45202404313},

		//CDN
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/sslkeys/?$`, Handler: cdn.GetSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42785817723},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/capacity$`, Handler: cdn.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4971852813},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/health/?$`, Handler: cdn.GetNameHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41353481943},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/health/?$`, Handler: cdn.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40853811343},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/domains/?$`, Handler: cdn.DomainsHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4269025603},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/routing$`, Handler: crstats.GetCDNRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 467229823},

		//CDN: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}$`, Handler: cdn.DeleteName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4088049593},

		//CDN: queue updates
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdns/{id}/queue_update$`, Handler: cdn.Queue, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4215159803},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdns/dnsseckeys/generate?$`, Handler: cdn.CreateDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:CREATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4753363},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}/dnsseckeys?$`, Handler: cdn.DeleteDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:DELETE", "CDN:UPDATE", "DELIVERY-SERVICE:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4711042073},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/dnsseckeys/?$`, Handler: cdn.GetDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4790106093},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `cdns/dnsseckeys/refresh/?$`, Handler: cdn.RefreshDNSSECKeysV4, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DNS-SEC:UPDATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47719971163},

		//CDN: Monitoring: Traffic Monitor
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/configs/monitoring?$`, Handler: crconfig.SnapshotGetMonitoringHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"MONITOR-CONFIG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42408478923},

		//Database dumps
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `dbdump/?`, Handler: dbdump.DBDump, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DBDUMP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4240166473},

		//Division: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `divisions/?$`, Handler: api.ReadHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40851815343},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `divisions/{id}$`, Handler: api.UpdateHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:UPDATE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4063691403},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `divisions/?$`, Handler: api.CreateHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:CREATE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4537138003},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `divisions/{id}$`, Handler: api.DeleteHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DIVISION:DELETE", "DIVISION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43253822373},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `logs/?$`, Handler: logs.Getv40, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"LOG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4483405503},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `logs/newcount/?$`, Handler: logs.GetNewCount, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"LOG:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44058330123},

		//Content invalidation jobs
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `jobs/?$`, Handler: api.ReadHandler(&invalidationjobs.InvalidationJobV4{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49667820413},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `jobs/?$`, Handler: invalidationjobs.DeleteV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:DELETE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4167807763},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `jobs/?$`, Handler: invalidationjobs.UpdateV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:UPDATE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4861342263},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `jobs/?`, Handler: invalidationjobs.CreateV40, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"JOB:CREATE", "JOB:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 404509553},

		//Login
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `user/login/?$`, Handler: login.LoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 43926708213},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `user/logout/?$`, Handler: login.LogoutHandler(d.Config.Secrets[0]), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 4434348253},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `user/login/oauth/?$`, Handler: login.OauthLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 44158860093},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `user/login/token/?$`, Handler: login.TokenLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 4024088413},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `user/reset_password/?$`, Handler: login.ResetPassword(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 42929146303},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `users/register/?$`, Handler: login.RegisterUser, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:CREATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43373},

		//ISO
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `osversions/?$`, Handler: iso.GetOSVersions, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ISO:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4760886573},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `isos/?$`, Handler: iso.ISOs, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ISO:GENERATE", "ISO:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4760336573},

		//User: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `users/?$`, Handler: user.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44919299003},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `users/{id}$`, Handler: user.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4138099803},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `users/{id}$`, Handler: user.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:UPDATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4354334043},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `users/?$`, Handler: user.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"USER:CREATE", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4762448163},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `user/current/?$`, Handler: user.Current, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 46107016143},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `user/current/?$`, Handler: user.ReplaceCurrentV4, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 4203},

		//Parameter: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `parameters/?$`, Handler: api.ReadHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42125542923},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `parameters/{id}$`, Handler: api.UpdateHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:UPDATE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48739361153},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `parameters/?$`, Handler: api.CreateHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:CREATE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46695108593},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `parameters/{id}$`, Handler: api.DeleteHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PARAMETER:DELETE", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4262771183},

		//Phys_Location: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `phys_locations/?$`, Handler: api.ReadHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4204051823},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `phys_locations/{id}$`, Handler: api.UpdateHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:UPDATE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4227950213},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `phys_locations/?$`, Handler: api.CreateHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:CREATE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42464566483},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `phys_locations/{id}$`, Handler: api.DeleteHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PHYSICAL-LOCATION:DELETE", "PHYSICAL-LOCATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 456142213},

		//Ping
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `ping$`, Handler: ping.Handler, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 45556615973},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `vault/ping/?$`, Handler: ping.Vault, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TRAFFIC-VAULT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48840121143},

		//Profile: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `profiles/?$`, Handler: api.ReadHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4687585893},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `profiles/{id}$`, Handler: api.UpdateHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 484391723},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profiles/?$`, Handler: api.CreateHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45402115563},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `profiles/{id}$`, Handler: api.DeleteHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:DELETE", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42055944653},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/export/?$`, Handler: profile.ExportProfileHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 401335173},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profiles/import/?$`, Handler: profile.ImportProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PARAMETER:CREATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4061432083},

		// Copy Profile
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{new_profile}/copy/{existing_profile}`, Handler: profile.CopyProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:CREATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4061432093},

		//Region: CRUDs
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `regions/?$`, Handler: api.ReadHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4100370853},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `regions/{id}$`, Handler: api.UpdateHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:UPDATE", "REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4223082243},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `regions/?$`, Handler: api.CreateHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:CREATE", "REGION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42883344883},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `regions/?$`, Handler: api.DeleteHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"REGION:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42326267583},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `topologies/?$`, Handler: api.CreateHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:CREATE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4871452221},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `topologies/?$`, Handler: api.ReadHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4871452222},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `topologies/?$`, Handler: api.UpdateHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:UPDATE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4871452223},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `topologies/?$`, Handler: api.DeleteHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TOPOLOGY:DELETE", "TOPOLOGY:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4871452224},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `topologies/{name}/queue_update$`, Handler: topology.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "TOPOLOGY:READ", "SERVER:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4205351748},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryserviceserver/?$`, Handler: dsserver.ReadDSSHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49461450333},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryserviceserver$`, Handler: dsserver.GetReplaceHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "SERVER:UPDATE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4297997883},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryserviceserver/{dsid}/{serverid}`, Handler: dsserver.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE", "SERVER:READ", "SERVER:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 45321845233},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xml_id}/servers$`, Handler: dsserver.GetCreateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "SERVER:UPDATE", "DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44281812063},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `servers/{id}/deliveryservices$`, Handler: api.ReadHandler(&dsserver.TODSSDeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4331154113},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/deliveryservices$`, Handler: server.AssignDeliveryServicesToServerHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "DELIVERY-SERVICE:UPDATE", "SERVER:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4801282533},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers$`, Handler: dsserver.GetReadAssigned, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CACHE-GROUP:READ", "CDN:READ", "TYPE:READ", "PROFILE:READ", "DELIVERY-SERVICE:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43451212233},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/capacity/?$`, Handler: deliveryservice.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42314091103},
		//Serverchecks
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `servercheck/?$`, Handler: servercheck.ReadServerCheck, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47961129223},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servercheck/?$`, Handler: servercheck.CreateUpdateServercheck, RequiredPrivLevel: auth.PrivLevelInvalid, RequiredPermissions: []string{"SERVER-CHECK:CREATE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47642815683},

		// Servercheck Extensions
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servercheck/extensions$`, Handler: extensions.Create, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:CREATE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4804985993},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `servercheck/extensions$`, Handler: extensions.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4834985993},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `servercheck/extensions/{id}$`, Handler: extensions.Delete, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CHECK:DELETE", "SERVER-CHECK:READ", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4804982993},

		//Server status
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}/status$`, Handler: server.UpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4766638513},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/queue_update$`, Handler: server.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:QUEUE", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41894713},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `servers/{host_name}/update_status$`, Handler: server.GetServerUpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4384515993},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servers/{id-or-name}/update$`, Handler: server.UpdateHandlerV4, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 443813233},

		//Server: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `servers/?$`, Handler: api.Wrap(server.Read, nil, nil), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47209592853},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}$`, Handler: api.Wrap(server.Update, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4586341033},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `servers/?$`, Handler: api.Wrap(server.Create, nil, nil), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:CREATE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42255580613},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `servers/{id}$`, Handler: api.Wrap(server.Delete, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:DELETE", "SERVER:READ", "DELIVERY-SERVICE:READ", "CDN:READ", "PHYSICAL-LOCATION:READ", "CACHE-GROUP:READ", "TYPE:READ", "PROFILE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4923222333},

		//Server Capability
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `server_capabilities$`, Handler: api.ReadHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4104073913},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `server_capabilities$`, Handler: api.CreateHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:CREATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40744707083},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `server_capabilities$`, Handler: api.UpdateHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:UPDATE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42543770109},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `server_capabilities$`, Handler: api.DeleteHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER-CAPABILITY:DELETE", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4364150383},

		//Server Server Capabilities: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `server_server_capabilities/?$`, Handler: api.ReadHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48002318893},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `server_server_capabilities/?$`, Handler: api.CreateHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42931668343},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `server_server_capabilities/?$`, Handler: api.DeleteHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVER:UPDATE", "SERVER:READ", "SERVER-CAPABILITY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40587140583},

		//Status: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `statuses/?$`, Handler: api.ReadHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42449056563},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `statuses/{id}$`, Handler: api.UpdateHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:UPDATE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42079665043},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `statuses/?$`, Handler: api.CreateHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:CREATE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43691236123},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `statuses/{id}$`, Handler: api.DeleteHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATUS:DELETE", "STATUS:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4551113603},

		//System
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `system/info/?$`, Handler: systeminfo.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 4210474753},

		//Type: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `types/?$`, Handler: api.ReadHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42267018233},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `types/{id}$`, Handler: api.UpdateHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:UPDATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 488601153},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `types/?$`, Handler: api.CreateHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:CREATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45133081953},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `types/{id}$`, Handler: api.DeleteHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TYPE:DELETE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 431757733},

		//About
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `about/?$`, Handler: about.Handler(), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 43175011663},

		//Coordinates
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `coordinates/?$`, Handler: api.ReadHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4967007453},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `coordinates/?$`, Handler: api.UpdateHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:UPDATE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4689261743},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `coordinates/?$`, Handler: api.CreateHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:CREATE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44281121573},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `coordinates/?$`, Handler: api.DeleteHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"COORDINATE:DELETE", "COORDINATE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43038498893},

		//CDN notification
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdn_notifications/?$`, Handler: cdnnotification.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 2221224514},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdn_notifications/?$`, Handler: cdnnotification.Create, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 2765223513},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdn_notifications/?$`, Handler: cdnnotification.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 2722411851},

		//CDN generic handlers:
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/?$`, Handler: api.ReadHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42303186213},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `cdns/{id}$`, Handler: api.UpdateHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43111789343},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdns/?$`, Handler: api.CreateHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:READ", "CDN:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 41605052893},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{id}$`, Handler: api.DeleteHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN:DELETE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4276946573},

		//Delivery service requests
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-REQUEST:READ", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46811639353},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Put, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42499079183},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Post, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:CREATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 493850393},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Delete, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:DELETE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42969850253},

		//Delivery service request: Actions
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/{id}/assign$`, Handler: dsrequest.GetAssignment, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-REQUEST:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47031602904},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/assign$`, Handler: dsrequest.PutAssignment, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DS-REQUEST:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47031602903},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/{id}/status$`, Handler: dsrequest.GetStatus, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4684150994},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/status$`, Handler: dsrequest.PutStatus, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DS-REQUEST:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4684150993},

		//Delivery service request comment: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_request_comments/?$`, Handler: api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-REQUEST:READ", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40326507373},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_request_comments/?$`, Handler: api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4604878473},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_request_comments/?$`, Handler: api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4272276723},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_request_comments/?$`, Handler: api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: []string{"DS-REQUEST:UPDATE", "DELIVERY-SERVICE:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4995046683},

		//Delivery service uri signing keys: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.GetURIsignkeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42930785583},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4084663353},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 476489693},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.RemoveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4299254173},

		//Delivery Service Required Capabilities: CRUD
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices_required_capabilities/?$`, Handler: api.DeprecatedReadHandler(&deliveryservice.RequiredCapability{}, util.StrPtr("the deliveryservices endpoint")), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41585222273},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices_required_capabilities/?$`, Handler: api.DeprecatedCreateHandler(&deliveryservice.RequiredCapability{}, util.StrPtr("the deliveryservices endpoint")), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40968739923},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices_required_capabilities/?$`, Handler: api.DeprecatedDeleteHandler(&deliveryservice.RequiredCapability{}, util.StrPtr("the deliveryservices endpoint")), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 44962893043},

		// Federations by CDN (the actual table for federation)
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/federations/?$`, Handler: api.ReadHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4892250323},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/federations/?$`, Handler: api.CreateHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:CREATE", "FEDERATION:READ, CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49548942193},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `cdns/{name}/federations/{id}$`, Handler: api.UpdateHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4260654663},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{name}/federations/{id}$`, Handler: api.DeleteHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:DELETE", "FEDERATION:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44428529023},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/dnsseckeys/ksk/generate$`, Handler: cdn.GenerateKSK, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DNS-SEC:CREATE", "CDN:UPDATE", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4729242813},

		//Origins
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `origins/?$`, Handler: api.ReadHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ORIGIN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4446492563},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `origins/?$`, Handler: api.UpdateHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:UPDATE", "ORIGIN:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 415677463},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `origins/?$`, Handler: api.CreateHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:CREATE", "ORIGIN:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 40995616433},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `origins/?$`, Handler: api.DeleteHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"ORIGIN:DELETE", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4602732633},

		//Roles
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `roles/?$`, Handler: role.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4870885833},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `roles/?$`, Handler: role.Update, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:UPDATE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46128974893},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `roles/?$`, Handler: role.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:CREATE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4306524063},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `roles/?$`, Handler: role.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"ROLE:DELETE", "ROLE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43567059823},

		//Delivery Services Regexes
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices_regexes/?$`, Handler: deliveryservicesregexes.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4055014533},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.DSGet, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4774327633},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.Post, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4127378003},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Put, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42483396913},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42467316633},

		//ServiceCategories
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `service_categories/?$`, Handler: api.ReadHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4085181543},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `service_categories/{name}/?$`, Handler: servicecategory.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:UPDATE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 406369141},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `service_categories/?$`, Handler: api.CreateHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:CREATE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 453713801},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `service_categories/{name}$`, Handler: api.DeleteHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"SERVICE-CATEGORY:DELETE", "SERVICE-CATEGORY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4325382238},

		//StaticDNSEntries
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `staticdnsentries/?$`, Handler: api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4289394773},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `staticdnsentries/?$`, Handler: api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:UPDATE", "STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4424571113},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `staticdnsentries/?$`, Handler: api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:CREATE", "STATIC-DN:READ", "CACHE-GROUP:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46291482383},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `staticdnsentries/?$`, Handler: api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"STATIC-DN:DELETE", "STATIC-DN:READ", "DELIVERY-SERVICE:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48460311323},

		//ProfileParameters
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.GetProfileID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4764649753},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.GetProfileName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42677378323},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.PostProfileParamsByName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43559455823},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.PostProfileParamsByID, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4168187083},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `profileparameters/?$`, Handler: api.ReadHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4506098053},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profileparameters/?$`, Handler: api.CreateHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ", "PROFILE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4288096933},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `profileparameter/?$`, Handler: profileparameter.PostProfileParam, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:READ", "PARAMETER:READ", "PROFILE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4242753},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `parameterprofile/?$`, Handler: profileparameter.PostParamProfile, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40806108613},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `profileparameters/{profileId}/{parameterId}$`, Handler: api.DeleteHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"PROFILE:UPDATE", "PROFILE:READ", "PARAMETER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4248395293},

		//Tenants
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `tenants/?$`, Handler: api.ReadHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46779678143},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `tenants/{id}$`, Handler: api.UpdateHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:UPDATE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40941314783},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `tenants/?$`, Handler: api.CreateHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:CREATE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4172480133},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `tenants/{id}$`, Handler: api.DeleteHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"TENANT:DELETE", "TENANT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4163655583},

		//CRConfig
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/?$`, Handler: crconfig.SnapshotGetHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49572736953},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/new/?$`, Handler: crconfig.Handler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4767168893},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `snapshot/?$`, Handler: crconfig.SnapshotHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"CDN-SNAPSHOT:CREATE", "CDN-SNAPSHOT:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49699118293},

		// Federations
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federations/all/?$`, Handler: federations.GetAll, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 410599863},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federations/?$`, Handler: federations.Get, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4549549943},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `federations/?$`, Handler: federations.AddFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:CREATE", "FEDERATION-RESOLVER:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48940647423},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `federations/?$`, Handler: federations.RemoveFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE"}, Authenticated: Authenticated, Middlewares: nil, ID: 420983233},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `federations/?$`, Handler: federations.ReplaceFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE", "FEDERATION-RESOLVER:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42831825163},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/deliveryservices/?$`, Handler: federations.PostDSes, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "DELIVERY-SERVICE:UPDATE", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 46828635133},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/deliveryservices/?$`, Handler: api.ReadHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4537730343},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/deliveryservices/{dsID}/?$`, Handler: api.DeleteHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "DELIVERY-SERVICE:UPDATE", "FEDERATION:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44174025703},

		// Federation Resolvers
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:CREATE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41343736613},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION-RESOLVER:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4566087593},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.AssignFederationResolversToFederationHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "FEDERATION-RESOLVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4566087603},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.GetFederationFederationResolversHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "FEDERATION-RESOLVER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4566087613},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION-RESOLVER:DELETE", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 40013},

		// Federations Users
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/users/?$`, Handler: federations.PostUsers, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "USER:READ", "FEDERATION:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47793349303},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/users/?$`, Handler: api.ReadHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"FEDERATION:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4940750153},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/users/{userID}/?$`, Handler: api.DeleteHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"FEDERATION:UPDATE", "FEDERATION:READ", "USER:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 49491028823},

		////DeliveryServices
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/?$`, Handler: api.ReadHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42383172943},
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodPost, Path: `deliveryservices/?$`, Handler: deliveryservice.CreateV41, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:CREATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4064315324},
		{Version: api.Version{Major: 4, Minor: 1}, Method: http.MethodPut, Path: `deliveryservices/{id}/?$`, Handler: deliveryservice.UpdateV41, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47665675674},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/?$`, Handler: deliveryservice.CreateV40, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:CREATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4064315323},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/?$`, Handler: deliveryservice.UpdateV40, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:UPDATE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 47665675673},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/safe/?$`, Handler: deliveryservice.UpdateSafe, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: []string{"DELIVERY-SERVICE-SAFE:UPDATE", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4472109313},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{id}/?$`, Handler: api.DeleteHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DELIVERY-SERVICE:DELETE", "DELIVERY-SERVICE:READ", "CDN:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4226420743},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers/eligible/?$`, Handler: deliveryservice.GetServersEligible, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "SERVER:READ", "CACHE-GROUP:READ", "TYPE:READ", "CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4747615843},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.GetSSLKeysByXMLID, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41357729073},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/add$`, Handler: deliveryservice.AddSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 48728785833},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.DeleteSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DELIVERY-SERVICE:READ", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 49267343},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/?$`, Handler: deliveryservice.GenerateSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4534390513},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, Handler: deliveryservice.CopyURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42625010763},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, Handler: deliveryservice.GenerateURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 45304828243},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{name}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42027192113},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/xmlId/{name}/urlkeys/?$`, Handler: deliveryservice.DeleteURLKeysByName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 42027192114},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DS-SECURITY-KEY:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4931971143},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{id}/urlkeys/?$`, Handler: deliveryservice.DeleteURLKeysByID, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:DELETE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4931971144},

		//Delivery service LetsEncrypt
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/letsencrypt/?$`, Handler: deliveryservice.GenerateLetsEncryptCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4534390523},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `letsencrypt/dnsrecords/?$`, Handler: deliveryservice.GetDnsChallengeRecords, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4534390553},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `letsencrypt/autorenew/?$`, Handler: deliveryservice.RenewCertificatesDeprecated, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: []string{"DS-SECURITY-KEY:CREATE", "DELIVERY-SERVICE:READ", "DELIVERY-SERVICE:UPDATE"}, Authenticated: Authenticated, Middlewares: nil, ID: 4534390563},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/health/?$`, Handler: deliveryservice.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ", "CACHE-GROUP:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42345901013},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/routing$`, Handler: crstats.GetDSRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 467339833},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 45696078243},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:CREATE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 43382163973},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPut, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:UPDATE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 44386082953},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodDelete, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:DELETE", "STEERING:READ", "DELIVERY-SERVICE:READ", "TYPE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 42880215153},

		// Stats Summary
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `stats_summary/?$`, Handler: trafficstats.GetStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4804985983},
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `stats_summary/?$`, Handler: trafficstats.CreateStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"STAT:CREATE", "STAT:READ", "CDN:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4804915983},

		//Pattern based consistent hashing endpoint
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodPost, Path: `consistenthash/?$`, Handler: consistenthash.Post, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"CDN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4607550763},

		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `steering/?$`, Handler: steering.Get, RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: []string{"STEERING:READ", "DELIVERY-SERVICE:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 41748524573},

		// Plugins
		{Version: api.Version{Major: 4, Minor: 0}, Method: http.MethodGet, Path: `plugins/?$`, Handler: plugins.Get(d.Plugins), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: []string{"PLUGIN:READ"}, Authenticated: Authenticated, Middlewares: nil, ID: 4834985393},

		/**
		 * 3.x API
		 */
		////DeliveryServices
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodPost, Path: `deliveryservices/?$`, Handler: deliveryservice.CreateV31, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2064315323},
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodPut, Path: `deliveryservices/{id}/?$`, Handler: deliveryservice.UpdateV31, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27665675673},

		// Acme account information
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodGet, Path: `acme_accounts/?$`, Handler: acme.Read, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2034390561},
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodPost, Path: `acme_accounts/?$`, Handler: acme.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2034390562},
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodPut, Path: `acme_accounts/?$`, Handler: acme.Update, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2034390563},
		{Version: api.Version{Major: 3, Minor: 1}, Method: http.MethodDelete, Path: `acme_accounts/{provider}/{email}?$`, Handler: acme.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2034390564},

		// API Capability
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `api_capabilities/?$`, Handler: apicapability.GetAPICapabilitiesHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28132065893},

		//ASNs
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `asns/?$`, Handler: api.UpdateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22641723173},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `asns/?$`, Handler: api.DeleteHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 202048983},

		//ASN: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `asns/?$`, Handler: api.ReadHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2738777223},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `asns/{id}$`, Handler: api.UpdateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29511986293},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `asns/?$`, Handler: api.CreateHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29994921883},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `asns/{id}$`, Handler: api.DeleteHandler(&asn.TOASNV11{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26725247693},

		// Traffic Stats access
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_stats`, Handler: trafficstats.GetDSStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23195690283},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cache_stats`, Handler: trafficstats.GetCacheStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24979979063},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `current_stats/?$`, Handler: trafficstats.GetCurrentStats, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27854428933},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `caches/stats/?$`, Handler: cachesstats.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28132065883},

		//CacheGroup: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cachegroups/?$`, Handler: api.ReadHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2230791103},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `cachegroups/{id}$`, Handler: api.UpdateHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2129545463},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/?$`, Handler: api.CreateHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 229826653},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cachegroups/{id}$`, Handler: api.DeleteHandler(&cachegroup.TOCacheGroup{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2278693653},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/queue_update$`, Handler: cachegroup.QueueUpdates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20716441103},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cachegroups/{id}/deliveryservices/?$`, Handler: cachegroup.DSPostHandlerV31, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25202404313},

		//CacheGroup Parameters: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cachegroupparameters/?$`, Handler: cachegroupparameter.ReadAllCacheGroupParameters, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2124497243},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cachegroupparameters/?$`, Handler: cachegroupparameter.AddCacheGroupParameters, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2124497253},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cachegroups/{id}/parameters/?$`, Handler: api.DeprecatedReadHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2124497233},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cachegroupparameters/{cachegroupID}/{parameterId}$`, Handler: api.DeprecatedDeleteHandler(&cachegroupparameter.TOCacheGroupParameter{}, nil), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2124497333},

		//Capabilities
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `capabilities/?$`, Handler: capabilities.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20081353},

		//CDN
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/sslkeys/?$`, Handler: cdn.GetSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22785817723},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/capacity$`, Handler: cdn.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2971852813},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/health/?$`, Handler: cdn.GetNameHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21353481943},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/health/?$`, Handler: cdn.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20853811343},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/domains/?$`, Handler: cdn.DomainsHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2269025603},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/routing$`, Handler: crstats.GetCDNRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 267229823},

		//CDN: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}$`, Handler: cdn.DeleteName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2088049593},

		//CDN: queue updates
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cdns/{id}/queue_update$`, Handler: cdn.Queue, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2215159803},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cdns/dnsseckeys/generate?$`, Handler: cdn.CreateDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2753363},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cdns/name/{name}/dnsseckeys?$`, Handler: cdn.DeleteDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2711042073},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/name/{name}/dnsseckeys/?$`, Handler: cdn.GetDNSSECKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2790106093},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/dnsseckeys/refresh/?$`, Handler: cdn.RefreshDNSSECKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27719971163},

		//CDN: Monitoring: Traffic Monitor
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/configs/monitoring?$`, Handler: crconfig.SnapshotGetMonitoringHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22408478923},

		//Database dumps
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `dbdump/?`, Handler: dbdump.DBDump, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2240166473},

		//Division: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `divisions/?$`, Handler: api.ReadHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20851815343},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `divisions/{id}$`, Handler: api.UpdateHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2063691403},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `divisions/?$`, Handler: api.CreateHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2537138003},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `divisions/{id}$`, Handler: api.DeleteHandler(&division.TODivision{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23253822373},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `logs/?$`, Handler: logs.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2483405503},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `logs/newcount/?$`, Handler: logs.GetNewCount, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24058330123},

		//Content invalidation jobs
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `jobs/?$`, Handler: api.ReadHandler(&invalidationjobs.InvalidationJob{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29667820413},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `jobs/?$`, Handler: invalidationjobs.Delete, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2167807763},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `jobs/?$`, Handler: invalidationjobs.Update, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2861342263},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `jobs/?`, Handler: invalidationjobs.Create, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 204509553},

		//Login
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `user/login/?$`, Handler: login.LoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 23926708213},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `user/logout/?$`, Handler: login.LogoutHandler(d.Config.Secrets[0]), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2434348253},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `user/login/oauth/?$`, Handler: login.OauthLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 24158860093},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `user/login/token/?$`, Handler: login.TokenLoginHandler(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 2024088413},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `user/reset_password/?$`, Handler: login.ResetPassword(d.DB, d.Config), RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 22929146303},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `users/register/?$`, Handler: login.RegisterUser, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23373},

		//ISO
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `osversions/?$`, Handler: iso.GetOSVersions, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2760886573},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `isos/?$`, Handler: iso.ISOs, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2760336573},

		//User: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `users/?$`, Handler: api.ReadHandler(&user.TOUser{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24919299003},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `users/{id}$`, Handler: api.ReadHandler(&user.TOUser{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2138099803},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `users/{id}$`, Handler: api.UpdateHandler(&user.TOUser{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2354334043},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `users/?$`, Handler: api.CreateHandler(&user.TOUser{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2762448163},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `user/current/?$`, Handler: user.Current, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26107016143},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `user/current/?$`, Handler: user.ReplaceCurrent, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2203},

		//Parameter: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `parameters/?$`, Handler: api.ReadHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22125542923},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `parameters/{id}$`, Handler: api.UpdateHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28739361153},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `parameters/?$`, Handler: api.CreateHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26695108593},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `parameters/{id}$`, Handler: api.DeleteHandler(&parameter.TOParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2262771183},

		//Phys_Location: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `phys_locations/?$`, Handler: api.ReadHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2204051823},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `phys_locations/{id}$`, Handler: api.UpdateHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2227950213},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `phys_locations/?$`, Handler: api.CreateHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22464566483},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `phys_locations/{id}$`, Handler: api.DeleteHandler(&physlocation.TOPhysLocation{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 256142213},

		//Ping
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `ping$`, Handler: ping.Handler, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: NoAuth, Middlewares: nil, ID: 25556615973},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `vault/ping/?$`, Handler: ping.Vault, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28840121143},

		//Profile: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `profiles/?$`, Handler: api.ReadHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2687585893},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `profiles/{id}$`, Handler: api.UpdateHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 284391723},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profiles/?$`, Handler: api.CreateHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25402115563},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `profiles/{id}$`, Handler: api.DeleteHandler(&profile.TOProfile{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22055944653},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/export/?$`, Handler: profile.ExportProfileHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 201335173},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profiles/import/?$`, Handler: profile.ImportProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2061432083},

		// Copy Profile
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{new_profile}/copy/{existing_profile}`, Handler: profile.CopyProfileHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2061432093},

		//Region: CRUDs
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `regions/?$`, Handler: api.ReadHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2100370853},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `regions/{id}$`, Handler: api.UpdateHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2223082243},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `regions/?$`, Handler: api.CreateHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22883344883},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `regions/?$`, Handler: api.DeleteHandler(&region.TORegion{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22326267583},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `topologies/?$`, Handler: api.CreateHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 3871452221},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `topologies/?$`, Handler: api.ReadHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 3871452222},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `topologies/?$`, Handler: api.UpdateHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 3871452223},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `topologies/?$`, Handler: api.DeleteHandler(&topology.TOTopology{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 3871452224},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `topologies/{name}/queue_update$`, Handler: topology.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 3205351748},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryserviceserver/?$`, Handler: dsserver.ReadDSSHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29461450333},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryserviceserver$`, Handler: dsserver.GetReplaceHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2297997883},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryserviceserver/{dsid}/{serverid}`, Handler: dsserver.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25321845233},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xml_id}/servers$`, Handler: dsserver.GetCreateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24281812063},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servers/{id}/deliveryservices$`, Handler: api.ReadHandler(&dsserver.TODSSDeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2331154113},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/deliveryservices$`, Handler: server.AssignDeliveryServicesToServerHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2801282533},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers$`, Handler: dsserver.GetReadAssigned, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23451212233},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/request`, Handler: deliveryservicerequests.Request, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2408752993},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/capacity/?$`, Handler: deliveryservice.GetCapacity, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22314091103},
		//Serverchecks
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servercheck/?$`, Handler: servercheck.ReadServerCheck, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27961129223},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servercheck/?$`, Handler: servercheck.CreateUpdateServercheck, RequiredPrivLevel: auth.PrivLevelInvalid, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27642815683},

		// Servercheck Extensions
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servercheck/extensions$`, Handler: extensions.Create, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2804985993},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servercheck/extensions$`, Handler: extensions.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2834985993},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `servercheck/extensions/{id}$`, Handler: extensions.Delete, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2804982993},

		//Server Details
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servers/details/?$`, Handler: server.GetDetailParamHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22612647143},

		//Server status
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}/status$`, Handler: server.UpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2766638513},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servers/{id}/queue_update$`, Handler: server.QueueUpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21894713},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servers/{host_name}/update_status$`, Handler: server.GetServerUpdateStatusHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2384515993},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servers/{id-or-name}/update$`, Handler: server.UpdateHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 143813233},

		//Server: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `servers/?$`, Handler: api.Wrap(server.Read, nil, nil), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27209592853},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `servers/{id}$`, Handler: api.Wrap(server.Update, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2586341033},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `servers/?$`, Handler: api.Wrap(server.Create, nil, nil), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22255580613},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `servers/{id}$`, Handler: api.Wrap(server.Delete, []string{"id"}, []string{"id"}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2923222333},

		//Server Capability
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `server_capabilities$`, Handler: api.ReadHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2104073913},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `server_capabilities$`, Handler: api.CreateHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20744707083},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `server_capabilities$`, Handler: api.DeleteHandler(&servercapability.TOServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2364150383},

		//Server Server Capabilities: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `server_server_capabilities/?$`, Handler: api.ReadHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28002318893},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `server_server_capabilities/?$`, Handler: api.CreateHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22931668343},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `server_server_capabilities/?$`, Handler: api.DeleteHandler(&server.TOServerServerCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20587140583},

		//Status: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `statuses/?$`, Handler: api.ReadHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22449056563},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `statuses/{id}$`, Handler: api.UpdateHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22079665043},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `statuses/?$`, Handler: api.CreateHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23691236123},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `statuses/{id}$`, Handler: api.DeleteHandler(&status.TOStatus{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2551113603},

		//System
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `system/info/?$`, Handler: systeminfo.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2210474753},

		//Type: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `types/?$`, Handler: api.ReadHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22267018233},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `types/{id}$`, Handler: api.UpdateHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 288601153},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `types/?$`, Handler: api.CreateHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25133081953},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `types/{id}$`, Handler: api.DeleteHandler(&types.TOType{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 231757733},

		//About
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `about/?$`, Handler: about.Handler(), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23175011663},

		//Coordinates
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `coordinates/?$`, Handler: api.ReadHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2967007453},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `coordinates/?$`, Handler: api.UpdateHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2689261743},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `coordinates/?$`, Handler: api.CreateHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24281121573},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `coordinates/?$`, Handler: api.DeleteHandler(&coordinate.TOCoordinate{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23038498893},

		//CDN generic handlers:
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/?$`, Handler: api.ReadHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22303186213},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `cdns/{id}$`, Handler: api.UpdateHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23111789343},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cdns/?$`, Handler: api.CreateHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21605052893},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{id}$`, Handler: api.DeleteHandler(&cdn.TOCDN{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2276946573},

		//Delivery service requests
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26811639353},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Put, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22499079183},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Post, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 293850393},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_requests/?$`, Handler: dsrequest.Delete, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22969850253},

		//Delivery service request: Actions
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/assign$`, Handler: dsrequest.PutAssignment, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27031602903},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_requests/{id}/status$`, Handler: dsrequest.PutStatus, RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2684150993},

		//Delivery service request comment: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservice_request_comments/?$`, Handler: api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20326507373},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservice_request_comments/?$`, Handler: api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2604878473},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservice_request_comments/?$`, Handler: api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2272276723},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservice_request_comments/?$`, Handler: api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), RequiredPrivLevel: auth.PrivLevelPortal, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2995046683},

		//Delivery service uri signing keys: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.GetURIsignkeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22930785583},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2084663353},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.SaveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 276489693},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{xmlID}/urisignkeys$`, Handler: urisigning.RemoveDeliveryServiceURIKeysHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2299254173},

		//Delivery Service Required Capabilities: CRUD
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices_required_capabilities/?$`, Handler: api.ReadHandler(&deliveryservice.RequiredCapability{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21585222273},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices_required_capabilities/?$`, Handler: api.CreateHandler(&deliveryservice.RequiredCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20968739923},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices_required_capabilities/?$`, Handler: api.DeleteHandler(&deliveryservice.RequiredCapability{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24962893043},

		// Federations by CDN (the actual table for federation)
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/{name}/federations/?$`, Handler: api.ReadHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2892250323},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/federations/?$`, Handler: api.CreateHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29548942193},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `cdns/{name}/federations/{id}$`, Handler: api.UpdateHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2260654663},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `cdns/{name}/federations/{id}$`, Handler: api.DeleteHandler(&cdnfederation.TOCDNFederation{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24428529023},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `cdns/{name}/dnsseckeys/ksk/generate$`, Handler: cdn.GenerateKSK, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2729242813},

		//Origins
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `origins/?$`, Handler: api.ReadHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2446492563},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `origins/?$`, Handler: api.UpdateHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 215677463},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `origins/?$`, Handler: api.CreateHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20995616433},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `origins/?$`, Handler: api.DeleteHandler(&origin.TOOrigin{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2602732633},

		//Roles
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `roles/?$`, Handler: api.ReadHandler(&role.TORole{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2870885833},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `roles/?$`, Handler: api.UpdateHandler(&role.TORole{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26128974893},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `roles/?$`, Handler: api.CreateHandler(&role.TORole{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2306524063},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `roles/?$`, Handler: api.DeleteHandler(&role.TORole{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23567059823},

		//Delivery Services Regexes
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices_regexes/?$`, Handler: deliveryservicesregexes.Get, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2055014533},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.DSGet, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2774327633},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/{dsid}/regexes/?$`, Handler: deliveryservicesregexes.Post, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2127378003},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Put, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22483396913},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{dsid}/regexes/{regexid}?$`, Handler: deliveryservicesregexes.Delete, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22467316633},

		//ServiceCategories
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `service_categories/?$`, Handler: api.ReadHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 1085181543},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `service_categories/{name}/?$`, Handler: servicecategory.Update, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 306369141},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `service_categories/?$`, Handler: api.CreateHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 553713801},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `service_categories/{name}$`, Handler: api.DeleteHandler(&servicecategory.TOServiceCategory{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 1325382238},

		//StaticDNSEntries
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `staticdnsentries/?$`, Handler: api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2289394773},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `staticdnsentries/?$`, Handler: api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2424571113},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `staticdnsentries/?$`, Handler: api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26291482383},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `staticdnsentries/?$`, Handler: api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28460311323},

		//ProfileParameters
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.GetProfileID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2764649753},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.GetProfileName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22677378323},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profiles/name/{name}/parameters/?$`, Handler: profileparameter.PostProfileParamsByName, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23559455823},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profiles/{id}/parameters/?$`, Handler: profileparameter.PostProfileParamsByID, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2168187083},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `profileparameters/?$`, Handler: api.ReadHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2506098053},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profileparameters/?$`, Handler: api.CreateHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2288096933},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `profileparameter/?$`, Handler: profileparameter.PostProfileParam, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2242753},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `parameterprofile/?$`, Handler: profileparameter.PostParamProfile, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20806108613},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `profileparameters/{profileId}/{parameterId}$`, Handler: api.DeleteHandler(&profileparameter.TOProfileParameter{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2248395293},

		//Tenants
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `tenants/?$`, Handler: api.ReadHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26779678143},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `tenants/{id}$`, Handler: api.UpdateHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20941314783},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `tenants/?$`, Handler: api.CreateHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2172480133},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `tenants/{id}$`, Handler: api.DeleteHandler(&apitenant.TOTenant{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2163655583},

		//CRConfig
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/?$`, Handler: crconfig.SnapshotGetHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29572736953},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `cdns/{cdn}/snapshot/new/?$`, Handler: crconfig.Handler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2767168893},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `snapshot/?$`, Handler: crconfig.SnapshotHandler, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29699118293},

		// Federations
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federations/all/?$`, Handler: federations.GetAll, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 210599863},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federations/?$`, Handler: federations.Get, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2549549943},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `federations/?$`, Handler: federations.AddFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28940647423},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `federations/?$`, Handler: federations.RemoveFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 220983233},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `federations/?$`, Handler: federations.ReplaceFederationResolverMappingsForCurrentUser, RequiredPrivLevel: auth.PrivLevelFederation, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22831825163},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/deliveryservices/?$`, Handler: federations.PostDSes, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 26828635133},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/deliveryservices/?$`, Handler: api.ReadHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2537730343},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/deliveryservices/{dsID}/?$`, Handler: api.DeleteHandler(&federations.TOFedDSes{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24174025703},

		// Federation Resolvers
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Create, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21343736613},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Read, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2566087593},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.AssignFederationResolversToFederationHandler, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2566087603},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/federation_resolvers/?$`, Handler: federations.GetFederationFederationResolversHandler, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2566087613},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `federation_resolvers/?$`, Handler: federation_resolvers.Delete, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 20013},

		// Federations Users
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `federations/{id}/users/?$`, Handler: federations.PostUsers, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27793349303},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `federations/{id}/users/?$`, Handler: api.ReadHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2940750153},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `federations/{id}/users/{userID}/?$`, Handler: api.DeleteHandler(&federations.TOUsers{}), RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29491028823},

		////DeliveryServices
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/?$`, Handler: api.ReadHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22383172943},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/?$`, Handler: deliveryservice.CreateV30, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2064314323},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/?$`, Handler: deliveryservice.UpdateV30, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 27665675273},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `deliveryservices/{id}/safe/?$`, Handler: deliveryservice.UpdateSafe, RequiredPrivLevel: auth.PrivLevelUnauthenticated, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2472109313},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/{id}/?$`, Handler: api.DeleteHandler(&deliveryservice.TODeliveryService{}), RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2226420743},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/servers/eligible/?$`, Handler: deliveryservice.GetServersEligible, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2747615843},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.GetSSLKeysByXMLID, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21357729073},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/add$`, Handler: deliveryservice.AddSSLKeys, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 28728785833},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `deliveryservices/xmlId/{xmlid}/sslkeys$`, Handler: deliveryservice.DeleteSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 29267343},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/?$`, Handler: deliveryservice.GenerateSSLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2534390513},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, Handler: deliveryservice.CopyURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22625010763},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, Handler: deliveryservice.GenerateURLKeys, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25304828243},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/xmlId/{name}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByName, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22027192113},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/urlkeys/?$`, Handler: deliveryservice.GetURLKeysByID, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2931971143},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `vault/bucket/{bucket}/key/{key}/values/?$`, Handler: vault.GetBucketKeyDeprecated, RequiredPrivLevel: auth.PrivLevelAdmin, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22205108013},

		//Delivery service LetsEncrypt
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `deliveryservices/sslkeys/generate/letsencrypt/?$`, Handler: deliveryservice.GenerateLetsEncryptCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2534390523},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `letsencrypt/dnsrecords/?$`, Handler: deliveryservice.GetDnsChallengeRecords, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2534390553},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `letsencrypt/autorenew/?$`, Handler: deliveryservice.RenewCertificates, RequiredPrivLevel: auth.PrivLevelOperations, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2534390563},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/health/?$`, Handler: deliveryservice.GetHealth, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22345901013},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `deliveryservices/{id}/routing$`, Handler: crstats.GetDSRouting, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 667339833},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 25696078243},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `steering/{deliveryservice}/targets/?$`, Handler: api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 23382163973},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPut, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 24386082953},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodDelete, Path: `steering/{deliveryservice}/targets/{target}/?$`, Handler: api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 22880215153},

		// Stats Summary
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `stats_summary/?$`, Handler: trafficstats.GetStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2804985983},
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `stats_summary/?$`, Handler: trafficstats.CreateStatsSummary, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2804915983},

		//Pattern based consistent hashing endpoint
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodPost, Path: `consistenthash/?$`, Handler: consistenthash.Post, RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2607550763},

		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `steering/?$`, Handler: steering.Get, RequiredPrivLevel: auth.PrivLevelSteering, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 21748524573},

		// Plugins
		{Version: api.Version{Major: 3, Minor: 0}, Method: http.MethodGet, Path: `plugins/?$`, Handler: plugins.Get(d.Plugins), RequiredPrivLevel: auth.PrivLevelReadOnly, RequiredPermissions: nil, Authenticated: Authenticated, Middlewares: nil, ID: 2834985393},
	}

	// sanity check to make sure all Route IDs are unique
	knownRouteIDs := make(map[int]struct{}, len(routes))
	for _, r := range routes {
		if _, found := knownRouteIDs[r.ID]; !found {
			knownRouteIDs[r.ID] = struct{}{}
		} else {
			return nil, nil, fmt.Errorf("route ID %d is already taken. Please give it a unique Route ID", r.ID)
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
			return nil, nil, errors.New(msg)
		}
	}

	return routes, proxyHandler, nil
}

func MemoryStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := runtime.MemStats{}
		runtime.ReadMemStats(&stats)

		bytes, err := json.Marshal(stats)
		if err != nil {
			api.HandleErr(w, r, nil, http.StatusInternalServerError, nil, fmt.Errorf("unable to marshal stats: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		api.WriteAndLogErr(w, r, bytes)
	}
}

func DBStatsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := db.DB.Stats()

		bytes, err := json.Marshal(stats)
		if err != nil {
			api.HandleErr(w, r, nil, http.StatusInternalServerError, nil, fmt.Errorf("unable to marshal stats: %w", err))
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

// CreateThrottledHandler takes a handler, and a max and uses a channel to insure the handler is used concurrently by only max number of routines
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
