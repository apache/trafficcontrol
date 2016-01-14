// Copyright 2015 Comcast Cable Communications Management, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @APIVersion 2.0 alpha
// @APITitle Traffic Operations
// @APIDescription Traffic Ops API
// @Contact https://traffic-control-cdn.net
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0

// @SubApi Version 2.0 API [/api/2.0]

package api

import (
	"errors"
	"strconv"
)

type ApiMethod int

const (
	GET ApiMethod = iota
	POST
	PUT
	DELETE
)

func (m ApiMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	}
	return "INVALID"
}

type ApiMethods []ApiMethod

// String returns a comma-separated list of the methods, as expected in headers such as Access-Control-Allow-Methods
func (methods ApiMethods) String() string {
	var s string
	for _, method := range methods {
		s += method.String() + ","
	}
	if s != "" {
		s = s[:len(s)-1] // strip trailing ,
	}
	return s
}

type ApiHandlerFunc func(pathParams map[string]string, payload []byte) (interface{}, error)
type ApiHandlerFuncMap map[ApiMethod]ApiHandlerFunc

func (handlerMap ApiHandlerFuncMap) Methods() ApiMethods {
	var methods []ApiMethod
	for method, _ := range handlerMap {
		methods = append(methods, method)
	}
	return methods
}

func ApiHandlers() map[string]ApiHandlerFuncMap {
	return map[string]ApiHandlerFuncMap{
		"cdn":                                 ApiHandlerFuncMap{GET: emptyWrap(getCdns), POST: bodyWrap(postCdn)},
		"cdn/{id}":                            ApiHandlerFuncMap{GET: idWrap(getCdnById), PUT: idBodyWrap(putCdn), DELETE: idWrap(delCdn)},
		"asn":                                 ApiHandlerFuncMap{GET: emptyWrap(getAsns), POST: bodyWrap(postAsn)},
		"asn/{id}":                            ApiHandlerFuncMap{GET: idWrap(getAsnById), PUT: idBodyWrap(putAsn), DELETE: idWrap(delAsn)},
		"cachegroup":                          ApiHandlerFuncMap{GET: emptyWrap(getCachegroups), POST: bodyWrap(postCachegroup)},
		"cachegroup/{id}":                     ApiHandlerFuncMap{GET: idWrap(getCachegroupById), PUT: idBodyWrap(putCachegroup), DELETE: idWrap(delCachegroup)},
		"cachegroup_parameter":                ApiHandlerFuncMap{GET: emptyWrap(getCachegroupParameters), POST: bodyWrap(postCachegroupParameter)},
		"cachegroup_parameter/{id}":           ApiHandlerFuncMap{GET: idWrap(getCachegroupParameterById), PUT: idBodyWrap(putCachegroupParameter), DELETE: idWrap(delCachegroupParameter)},
		"deliveryservice":                     ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservices), POST: bodyWrap(postDeliveryservice)},
		"deliveryserviceRegex/{id}":           ApiHandlerFuncMap{GET: idWrap(getDeliveryserviceRegexById), PUT: idBodyWrap(putDeliveryserviceRegex), DELETE: idWrap(delDeliveryserviceRegex)},
		"deliveryservice_regex":               ApiHandlerFuncMap{GET: emptyWrap(getDeliveryserviceRegexs), POST: bodyWrap(postDeliveryserviceRegex)},
		"deliveryservice_regex/{id}":          ApiHandlerFuncMap{GET: idWrap(getDeliveryserviceRegexById), PUT: idBodyWrap(putDeliveryserviceRegex), DELETE: idWrap(delDeliveryserviceRegex)},
		"deliveryservice_server":              ApiHandlerFuncMap{GET: emptyWrap(getDeliveryserviceServers), POST: bodyWrap(postDeliveryserviceServer)},
		"deliveryservice_server/{id}":         ApiHandlerFuncMap{GET: idWrap(getDeliveryserviceServerById), PUT: idBodyWrap(putDeliveryserviceServer), DELETE: idWrap(delDeliveryserviceServer)},
		"deliveryservice_tmuser":              ApiHandlerFuncMap{GET: emptyWrap(getDeliveryserviceTmusers), POST: bodyWrap(postDeliveryserviceTmuser)},
		"deliveryservice_tmuser/{id}":         ApiHandlerFuncMap{GET: idWrap(getDeliveryserviceTmuserById), PUT: idBodyWrap(putDeliveryserviceTmuser), DELETE: idWrap(delDeliveryserviceTmuser)},
		"division":                            ApiHandlerFuncMap{GET: emptyWrap(getDivisions), POST: bodyWrap(postDivision)},
		"division/{id}":                       ApiHandlerFuncMap{GET: idWrap(getDivisionById), PUT: idBodyWrap(putDivision), DELETE: idWrap(delDivision)},
		"federation":                          ApiHandlerFuncMap{GET: emptyWrap(getFederations), POST: bodyWrap(postFederation)},
		"federation/{id}":                     ApiHandlerFuncMap{GET: idWrap(getFederationById), PUT: idBodyWrap(putFederation), DELETE: idWrap(delFederation)},
		"federation_deliveryservice":          ApiHandlerFuncMap{GET: emptyWrap(getFederationDeliveryservices), POST: bodyWrap(postFederationDeliveryservice)},
		"federation_deliveryservice/{id}":     ApiHandlerFuncMap{GET: idWrap(getFederationDeliveryserviceById), PUT: idBodyWrap(putFederationDeliveryservice), DELETE: idWrap(delFederationDeliveryservice)},
		"federation_federation_resolver":      ApiHandlerFuncMap{GET: emptyWrap(getFederationFederationResolvers), POST: bodyWrap(postFederationFederationResolver)},
		"federation_federation_resolver/{id}": ApiHandlerFuncMap{GET: idWrap(getFederationFederationResolverById), PUT: idBodyWrap(putFederationFederationResolver), DELETE: idWrap(delFederationFederationResolver)},
		"federation_resolver":                 ApiHandlerFuncMap{GET: emptyWrap(getFederationResolvers), POST: bodyWrap(postFederationResolver)},
		"federation_resolver/{id}":            ApiHandlerFuncMap{GET: idWrap(getFederationResolverById), PUT: idBodyWrap(putFederationResolver), DELETE: idWrap(delFederationResolver)},
		"federation_tmuser":                   ApiHandlerFuncMap{GET: emptyWrap(getFederationTmusers), POST: bodyWrap(postFederationTmuser)},
		"federation_tmuser/{id}":              ApiHandlerFuncMap{GET: idWrap(getFederationTmuserById), PUT: idBodyWrap(putFederationTmuser), DELETE: idWrap(delFederationTmuser)},
		"job":                    ApiHandlerFuncMap{GET: emptyWrap(getJobs), POST: bodyWrap(postJob)},
		"job/{id}":               ApiHandlerFuncMap{GET: idWrap(getJobById), PUT: idBodyWrap(putJob), DELETE: idWrap(delJob)},
		"job_result":             ApiHandlerFuncMap{GET: emptyWrap(getJobResults), POST: bodyWrap(postJobResult)},
		"job_result/{id}":        ApiHandlerFuncMap{GET: idWrap(getJobResultById), PUT: idBodyWrap(putJobResult), DELETE: idWrap(delJobResult)},
		"job_status":             ApiHandlerFuncMap{GET: emptyWrap(getJobStatuss), POST: bodyWrap(postJobStatus)},
		"job_status/{id}":        ApiHandlerFuncMap{GET: idWrap(getJobStatusById), PUT: idBodyWrap(putJobStatus), DELETE: idWrap(delJobStatus)},
		"log":                    ApiHandlerFuncMap{GET: emptyWrap(getLogs), POST: bodyWrap(postLog)},
		"log/{id}":               ApiHandlerFuncMap{GET: idWrap(getLogById), PUT: idBodyWrap(putLog), DELETE: idWrap(delLog)},
		"parameter":              ApiHandlerFuncMap{GET: emptyWrap(getParameters), POST: bodyWrap(postParameter)},
		"parameter/{id}":         ApiHandlerFuncMap{GET: idWrap(getParameterById), PUT: idBodyWrap(putParameter), DELETE: idWrap(delParameter)},
		"phys_location":          ApiHandlerFuncMap{GET: emptyWrap(getPhysLocations), POST: bodyWrap(postPhysLocation)},
		"phys_location/{id}":     ApiHandlerFuncMap{GET: idWrap(getPhysLocationById), PUT: idBodyWrap(putPhysLocation), DELETE: idWrap(delPhysLocation)},
		"profile":                ApiHandlerFuncMap{GET: emptyWrap(getProfiles), POST: bodyWrap(postProfile)},
		"profile/{id}":           ApiHandlerFuncMap{GET: idWrap(getProfileById), PUT: idBodyWrap(putProfile), DELETE: idWrap(delProfile)},
		"profile_parameter":      ApiHandlerFuncMap{GET: emptyWrap(getProfileParameters), POST: bodyWrap(postProfileParameter)},
		"profile_parameter/{id}": ApiHandlerFuncMap{GET: idWrap(getProfileParameterById), PUT: idBodyWrap(putProfileParameter), DELETE: idWrap(delProfileParameter)},
		"regex":                  ApiHandlerFuncMap{GET: emptyWrap(getRegexs), POST: bodyWrap(postRegex)},
		"regex/{id}":             ApiHandlerFuncMap{GET: idWrap(getRegexById), PUT: idBodyWrap(putRegex), DELETE: idWrap(delRegex)},
		"region":                 ApiHandlerFuncMap{GET: emptyWrap(getRegions), POST: bodyWrap(postRegion)},
		"region/{id}":            ApiHandlerFuncMap{GET: idWrap(getRegionById), PUT: idBodyWrap(putRegion), DELETE: idWrap(delRegion)},
		"role":                   ApiHandlerFuncMap{GET: emptyWrap(getRoles), POST: bodyWrap(postRole)},
		"role/{id}":              ApiHandlerFuncMap{GET: idWrap(getRoleById), PUT: idBodyWrap(putRole), DELETE: idWrap(delRole)},
		"server":                 ApiHandlerFuncMap{GET: emptyWrap(getServers), POST: bodyWrap(postServer)},
		"server/{id}":            ApiHandlerFuncMap{GET: idWrap(getServerById), PUT: idBodyWrap(putServer), DELETE: idWrap(delServer)},
		"servercheck":            ApiHandlerFuncMap{GET: emptyWrap(getServerchecks), POST: bodyWrap(postServercheck)},
		"servercheck/{id}":       ApiHandlerFuncMap{GET: idWrap(getServercheckById), PUT: idBodyWrap(putServercheck), DELETE: idWrap(delServercheck)},
		"staticdnsentry":         ApiHandlerFuncMap{GET: emptyWrap(getStaticdnsentrys), POST: bodyWrap(postStaticdnsentry)},
		"staticdnsentry/{id}":    ApiHandlerFuncMap{GET: idWrap(getStaticdnsentryById), PUT: idBodyWrap(putStaticdnsentry), DELETE: idWrap(delStaticdnsentry)},
		"stats_summary":          ApiHandlerFuncMap{GET: emptyWrap(getStatsSummarys), POST: bodyWrap(postStatsSummary)},
		"stats_summary/{id}":     ApiHandlerFuncMap{GET: idWrap(getStatsSummaryById), PUT: idBodyWrap(putStatsSummary), DELETE: idWrap(delStatsSummary)},
		"status":                 ApiHandlerFuncMap{GET: emptyWrap(getStatuss), POST: bodyWrap(postStatus)},
		"status/{id}":            ApiHandlerFuncMap{GET: idWrap(getStatusById), PUT: idBodyWrap(putStatus), DELETE: idWrap(delStatus)},
		"tm_user":                ApiHandlerFuncMap{GET: emptyWrap(getTmUsers), POST: bodyWrap(postTmUser)},
		"tm_user/{id}":           ApiHandlerFuncMap{GET: idWrap(getTmUserById), PUT: idBodyWrap(putTmUser), DELETE: idWrap(delTmUser)},
		"to_extension":           ApiHandlerFuncMap{GET: emptyWrap(getToExtensions), POST: bodyWrap(postToExtension)},
		"to_extension/{id}":      ApiHandlerFuncMap{GET: idWrap(getToExtensionById), PUT: idBodyWrap(putToExtension), DELETE: idWrap(delToExtension)},
		"type":                   ApiHandlerFuncMap{GET: emptyWrap(getTypes), POST: bodyWrap(postType)},
		"type/{id}":              ApiHandlerFuncMap{GET: idWrap(getTypeById), PUT: idBodyWrap(putType), DELETE: idWrap(delType)},
	}
}

type EmptyHandlerFunc func() (interface{}, error)
type IntHandlerFunc func(id int) (interface{}, error)
type BodyHandlerFunc func(payload []byte) (interface{}, error)
type IntBodyHandlerFunc func(id int, payload []byte) (interface{}, error)

func idBodyWrap(f IntBodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte) (interface{}, error) {
		if strid, ok := pathParams["id"]; !ok {
			return nil, errors.New("Id missing")
		} else if id, err := strconv.Atoi(strid); err != nil {
			return nil, errors.New("Id is not an integer: " + strid)
		} else {
			return f(id, payload)
		}
	}
}

func idWrap(f IntHandlerFunc) ApiHandlerFunc {
	return idBodyWrap(func(id int, payload []byte) (interface{}, error) {
		return f(id)
	})
}

func bodyWrap(f BodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte) (interface{}, error) {
		return f(payload)
	}
}

func emptyWrap(f EmptyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte) (interface{}, error) {
		return f()
	}
}
