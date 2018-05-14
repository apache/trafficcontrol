package main

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
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	tclog "github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/asn"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/cdn"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/coordinate"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/crconfig"
	dsrequest "github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request/comment"
	dsserver "github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/servers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservicesregexes"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/division"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/hwinfo"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/physlocation"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/ping"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/profile"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/profileparameter"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/region"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/role"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/systeminfo"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/types"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/origin"
	"github.com/basho/riak-go-client"
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
		{1.2, http.MethodGet, `asns/?(\.json)?$`, api.ReadHandler(asn.GetRefTypeV12(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `asns/?(\.json)?$`, asn.V11ReadAll(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `asns/{id}$`, api.ReadHandler(asn.GetRefTypeV11(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `asns/{id}$`, api.UpdateHandler(asn.GetRefTypeV11(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `asns/?$`, api.CreateHandler(asn.GetRefTypeV11(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(asn.GetRefTypeV11(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//CacheGroup: CRUD
		{1.1, http.MethodGet, `cachegroups/trimmed/?(\.json)?$`, cachegroup.TrimmedHandler(d.DB.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cachegroups/?(\.json)?$`, api.ReadHandler(cachegroup.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cachegroups/{id}$`, api.ReadHandler(cachegroup.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(cachegroup.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cachegroups/?$`, api.CreateHandler(cachegroup.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(cachegroup.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//CDN
		{1.1, http.MethodGet, `cdns/metric_types`, notImplementedHandler, 0, NoAuth, nil}, // MUST NOT end in $, because the 1.x route is longer
		{1.1, http.MethodGet, `cdns/capacity$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `cdns/configs$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `cdns/domains$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `cdns/health$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `cdns/routing$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},

		//CDN: CRUD
		{1.1, http.MethodGet, `cdns/?(\.json)?$`, api.ReadHandler(cdn.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{id}$`, api.ReadHandler(cdn.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(cdn.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `cdns/?$`, api.CreateHandler(cdn.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(cdn.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		{1.1, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName(d.DB.DB), auth.PrivLevelOperations, Authenticated, nil},

		//CDN: Monitoring: Traffic Monitor
		{1.1, http.MethodGet, `cdns/{name}/configs/monitoring(\.json)?$`, monitoringHandler(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},

		//Division: CRUD
		{1.1, http.MethodGet, `divisions/?(\.json)?$`, api.ReadHandler(division.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `divisions/{id}$`, api.ReadHandler(division.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(division.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `divisions/?$`, api.CreateHandler(division.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(division.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//HWInfo
		{1.1, http.MethodGet, `hwinfo-wip/?(\.json)?$`, hwinfo.HWInfoHandler(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},

		//Login
		{1.2, http.MethodPost, `user/login/?$`, auth.LoginHandler(d.DB, d.Config), 0, NoAuth, nil}, {1.3, http.MethodPost, `user/login/?$`, auth.LoginHandler(d.DB, d.Config), 0, NoAuth, nil},

		//Parameter: CRUD
		{1.1, http.MethodGet, `parameters/?(\.json)?$`, api.ReadHandler(parameter.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `parameters/{id}$`, api.ReadHandler(parameter.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(parameter.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `parameters/?$`, api.CreateHandler(parameter.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(parameter.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Phys_Location: CRUD
		{1.1, http.MethodGet, `phys_locations/?(\.json)?$`, api.ReadHandler(physlocation.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `phys_locations/{id}$`, api.ReadHandler(physlocation.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(physlocation.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `phys_locations/?$`, api.CreateHandler(physlocation.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(physlocation.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Ping
		{1.1, http.MethodGet, `ping$`, ping.PingHandler(), 0, NoAuth, nil},

		//Profile: CRUD
		{1.1, http.MethodGet, `profiles/?(\.json)?$`, api.ReadHandler(profile.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{id}$`, api.ReadHandler(profile.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(profile.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/?$`, api.CreateHandler(profile.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(profile.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Region: CRUD
		{1.1, http.MethodGet, `regions/?(\.json)?$`, api.ReadHandler(region.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `regions/{id}$`, api.ReadHandler(region.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `regions/{id}$`, api.UpdateHandler(region.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `regions/?$`, api.CreateHandler(region.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `regions/{id}$`, api.DeleteHandler(region.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		// get all edge servers associated with a delivery service (from deliveryservice_server table)
		{1.1, http.MethodGet, `deliveryserviceserver$`, dsserver.ReadDSSHandler(d.DB),auth.PrivLevelReadOnly, Authenticated, nil},
		//{1.1, http.MethodPost,`deliveryserviceserver$`, dsserver.GetReplaceHandler(d.DB),auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost,`deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler( d.DB ) ,auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(dsserver.GetDServiceRef(), d.DB),auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadHandler(d.DB, "assigned"),auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, dsserver.GetReadHandler(d.DB, "unassigned"),auth.PrivLevelReadOnly, Authenticated, nil},
		//{1.1, http.MethodGet, `deliveryservices/{id}/servers/eligible$`, dsserver.GetReadHandler(d.DB, "eligible"),auth.PrivLevelReadOnly, Authenticated, nil},

		//Server
		{1.1, http.MethodGet, `servers/checks$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `servers/details$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `servers/status$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},
		{1.1, http.MethodGet, `servers/totals$`, handlerToFunc(proxyHandler), 0, NoAuth, []Middleware{}},

		//Server: CRUD
		{1.1, http.MethodGet, `servers/?(\.json)?$`, api.ReadHandler(server.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `servers/{id}$`, api.ReadHandler(server.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `servers/{id}$`, api.UpdateHandler(server.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `servers/?$`, api.CreateHandler(server.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `servers/{id}$`, api.DeleteHandler(server.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Status: CRUD
		{1.1, http.MethodGet, `statuses/?(\.json)?$`, api.ReadHandler(status.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `statuses/{id}$`, api.ReadHandler(status.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(status.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `statuses/?$`, api.CreateHandler(status.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(status.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//System
		{1.1, http.MethodGet, `system/info/?(\.json)?$`, systeminfo.Handler(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},

		//Type: CRUD
		{1.1, http.MethodGet, `types/?(\.json)?$`, api.ReadHandler(types.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `types/{id}$`, api.ReadHandler(types.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPut, `types/{id}$`, api.UpdateHandler(types.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `types/?$`, api.CreateHandler(types.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `types/{id}$`, api.DeleteHandler(types.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//About
		{1.3, http.MethodGet, `about/?(\.json)?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil},

		//Coordinates
		{1.3, http.MethodGet, `coordinates/?(\.json)?$`, api.ReadHandler(coordinate.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `coordinates/?$`, api.ReadHandler(coordinate.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `coordinates/?$`, api.UpdateHandler(coordinate.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `coordinates/?$`, api.CreateHandler(coordinate.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(coordinate.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Delivery service request: CRUD
		{1.3, http.MethodGet, `deliveryservice_requests/?(\.json)?$`, api.ReadHandler(dsrequest.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_requests/?$`, api.UpdateHandler(dsrequest.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservice_requests/?$`, api.CreateHandler(dsrequest.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservice_requests/?$`, api.DeleteHandler(dsrequest.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service request: Actions
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/assign$`, api.UpdateHandler(dsrequest.GetAssignRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_requests/{id}/status$`, api.UpdateHandler(dsrequest.GetStatusRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service request comment: CRUD
		{1.3, http.MethodGet, `deliveryservice_request_comments/?(\.json)?$`, api.ReadHandler(comment.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(comment.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(comment.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(comment.GetRefType(), d.DB), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery service uri signing keys: CRUD
		{1.3, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, getURIsignkeysHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, saveDeliveryServiceURIKeysHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, saveDeliveryServiceURIKeysHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, removeDeliveryServiceURIKeysHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},

		//Origins
		{1.3, http.MethodGet, `origins/?(\.json)?$`, api.ReadHandler(origin.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodGet, `origins/?$`, api.ReadHandler(origin.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `origins/?$`, api.UpdateHandler(origin.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodPost, `origins/?$`, api.CreateHandler(origin.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodDelete, `origins/?$`, api.DeleteHandler(origin.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Roles
		{1.3, http.MethodGet, `roles/?(\.json)?$`, api.ReadHandler(role.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.3, http.MethodPut, `roles/?$`, api.UpdateHandler(role.GetRefType(), d.DB), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, `roles/?$`, api.CreateHandler(role.GetRefType(), d.DB), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodDelete, `roles/?$`, api.DeleteHandler(role.GetRefType(), d.DB), auth.PrivLevelAdmin, Authenticated, nil},

		//Delivery Services Regexes
		{1.1, http.MethodGet, `deliveryservices_regexes/?(\.json)?$`, deliveryservicesregexes.Get(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.DSGet(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.DSGetID(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `deliveryservices/{dsid}/regexes/?(\.json)?$`, deliveryservicesregexes.Post(d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Put(d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?(\.json)?$`, deliveryservicesregexes.Delete(d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//Servers
		{1.3, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler(d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.3, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler(d.DB), auth.PrivLevelReadOnly, Authenticated, nil},

		{1.1, http.MethodGet, `staticdnsentries/?(\.json)?$`, api.ReadHandler(staticdnsentry.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},

		//ProfileParameters
		{1.1, http.MethodGet, `profiles/{id}/parameters/?(\.json)?$`, profileparameter.GetProfileID(d.DB.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/{id}/unassigned_parameters/?(\.json)?$`, profileparameter.GetUnassigned(d.DB.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `profiles/name/{name}/parameters/?(\.json)?$`, profileparameter.GetProfileName(d.DB.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodGet, `parameters/profile/{name}/?(\.json)?$`, profileparameter.GetProfileName(d.DB.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName(d.DB.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID(d.DB.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodGet, `profileparameters/?(\.json)?$`, api.ReadHandler(profileparameter.GetRefType(), d.DB), auth.PrivLevelReadOnly, Authenticated, nil},
		{1.1, http.MethodPost, `profileparameters/?$`, api.CreateHandler(profileparameter.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam(d.DB.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile(d.DB.DB), auth.PrivLevelOperations, Authenticated, nil},
		{1.1, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(profileparameter.GetRefType(), d.DB), auth.PrivLevelOperations, Authenticated, nil},

		//SSLKeys deliveryservice endpoints here that are marked  marked as '-wip' need to have tenancy checks added
		{1.3, http.MethodGet, `deliveryservices-wip/xmlId/{xmlID}/sslkeys$`, getDeliveryServiceSSLKeysByXMLIDHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodGet, `deliveryservices-wip/hostname/{hostName}/sslkeys$`, getDeliveryServiceSSLKeysByHostNameHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},
		{1.3, http.MethodPost, `deliveryservices-wip/hostname/{hostName}/sslkeys/add$`, addDeliveryServiceSSLKeysHandler(d.DB, d.Config), auth.PrivLevelAdmin, Authenticated, nil},

		//CRConfig
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
		{1.1, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
		{1.1, http.MethodPut, `cdns/{id}/snapshot/?$`, crconfig.SnapshotHandler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
		{1.1, http.MethodPut, `snapshot/{cdn}/?$`, crconfig.SnapshotHandler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
	}

	// rawRoutes are served at the root path. These should be almost exclusively old Perl pre-API routes, which have yet to be converted in all clients. New routes should be in the versioned API path.
	rawRoutes := []RawRoute{
		// DEPRECATED - use PUT /api/1.2/snapshot/{cdn}
		{http.MethodGet, `tools/write_crconfig/{cdn}/?$`, crconfig.SnapshotOldGUIHandler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
		// DEPRECATED - use GET /api/1.2/cdns/{cdn}/snapshot
		{http.MethodGet, `CRConfig-Snapshots/{cdn}/CRConfig.json?$`, crconfig.SnapshotOldGetHandler(d.DB, d.Config), crconfig.PrivLevel, Authenticated, nil},
	}

	return routes, rawRoutes, proxyHandler, nil
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
		//Other knobs we can turn: ExpectContinueTimeout,IdleConnTimeout
	}
	rp := httputil.NewSingleHostReverseProxy(d.URL)
	rp.Transport = tr

	var errorLogger interface{}
	errorLogger, err := tclog.GetLogWriter(d.Config.ErrorLog())
	if err != nil {
		tclog.Errorln("could not create error log writer for proxy: ", err)
	}
	if errorLogger != nil {
		rp.ErrorLog = log.New(errorLogger.(io.Writer), "proxy error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC) //if we don't provide a logger to the reverse proxy it logs to stdout/err and is lost when ran by a script.
		riak.SetErrorLogger(log.New(errorLogger.(io.Writer), "riak error: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC))
	}
	var infoLogger interface{}
	infoLogger, err = tclog.GetLogWriter(d.Config.InfoLog())
	if err != nil {
		tclog.Errorln("could not create info log writer for proxy: ", err)
	}
	if infoLogger != nil {
		riak.SetLogger(log.New(infoLogger.(io.Writer), "riak info: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC))
	}
	tclog.Debugf("our reverseProxy: %++v\n", rp)
	tclog.Debugf("our reverseProxy's transport: %++v\n", tr)
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
