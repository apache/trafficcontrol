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

package api

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type ApiMethod int

const (
	GET ApiMethod = iota
	POST
	PUT
	DELETE
	OPTIONS
)

const API_PATH = "/api/2.0/"

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
	case OPTIONS:
		return "OPTIONS"
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

type ApiHandlerFunc func(pathParams map[string]string, payload []byte, dbb *sqlx.DB) (interface{}, error)
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
		"cdn":                                  ApiHandlerFuncMap{GET: emptyWrap(getCdnss), POST: bodyWrap(postCdns)},
		"cdn/{key}":                            ApiHandlerFuncMap{GET: stringWrap(getCdnsById), PUT: stringBodyWrap(putCdns), DELETE: stringWrap(delCdns)},
		"asn":                                  ApiHandlerFuncMap{GET: emptyWrap(getAsnss), POST: bodyWrap(postAsns)},
		"asn/{key}":                            ApiHandlerFuncMap{GET: int64Wrap(getAsnsById), PUT: int64BodyWrap(putAsns), DELETE: int64Wrap(delAsns)},
		"cachegroup":                           ApiHandlerFuncMap{GET: emptyWrap(getCachegroupss), POST: bodyWrap(postCachegroups)},
		"cachegroup/{key}":                     ApiHandlerFuncMap{GET: stringWrap(getCachegroupsById), PUT: stringBodyWrap(putCachegroups), DELETE: stringWrap(delCachegroups)},
		"cachegroup_parameter":                 ApiHandlerFuncMap{GET: emptyWrap(getCachegroupsParameterss), POST: bodyWrap(postCachegroupsParameters)},
		"cachegroup_parameter/{key}":           ApiHandlerFuncMap{GET: stringWrap(getCachegroupsParametersById), PUT: stringBodyWrap(putCachegroupsParameters), DELETE: stringWrap(delCachegroupsParameters)},
		"deliveryservice":                      ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicess), POST: bodyWrap(postDeliveryservices)},
		"deliveryservice/{key}":                ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesById), PUT: stringBodyWrap(putDeliveryservices), DELETE: stringWrap(delDeliveryservices)},
		"deliveryserviceRegex/{key}":           ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesRegexesById), PUT: stringBodyWrap(putDeliveryservicesRegexes), DELETE: stringWrap(delDeliveryservicesRegexes)},
		"deliveryservice_regex":                ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesRegexess), POST: bodyWrap(postDeliveryservicesRegexes)},
		"deliveryservice_regex/{key}":          ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesRegexesById), PUT: stringBodyWrap(putDeliveryservicesRegexes), DELETE: stringWrap(delDeliveryservicesRegexes)},
		"deliveryservice_server":               ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesServerss), POST: bodyWrap(postDeliveryservicesServers)},
		"deliveryservice_server/{key}":         ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesServersById), PUT: stringBodyWrap(putDeliveryservicesServers), DELETE: stringWrap(delDeliveryservicesServers)},
		"deliveryservice_tmuser":               ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesUserss), POST: bodyWrap(postDeliveryservicesUsers)},
		"deliveryservice_tmuser/{key}":         ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesUsersById), PUT: stringBodyWrap(putDeliveryservicesUsers), DELETE: stringWrap(delDeliveryservicesUsers)},
		"division":                             ApiHandlerFuncMap{GET: emptyWrap(getDivisionss), POST: bodyWrap(postDivisions)},
		"division/{key}":                       ApiHandlerFuncMap{GET: stringWrap(getDivisionsById), PUT: stringBodyWrap(putDivisions), DELETE: stringWrap(delDivisions)},
		"federation":                           ApiHandlerFuncMap{GET: emptyWrap(getFederationss), POST: bodyWrap(postFederations)},
		"federation/{key}":                     ApiHandlerFuncMap{GET: int64Wrap(getFederationsById), PUT: int64BodyWrap(putFederations), DELETE: int64Wrap(delFederations)},
		"federation_deliveryservice":           ApiHandlerFuncMap{GET: emptyWrap(getFederationsDeliveryservicess), POST: bodyWrap(postFederationsDeliveryservices)},
		"federation_deliveryservice/{key}":     ApiHandlerFuncMap{GET: int64Wrap(getFederationsDeliveryservicesById), PUT: int64BodyWrap(putFederationsDeliveryservices), DELETE: int64Wrap(delFederationsDeliveryservices)},
		"federation_federation_resolver":       ApiHandlerFuncMap{GET: emptyWrap(getFederationsFederationResolverss), POST: bodyWrap(postFederationsFederationResolvers)},
		"federation_federation_resolver/{key}": ApiHandlerFuncMap{GET: int64Wrap(getFederationsFederationResolversById), PUT: int64BodyWrap(putFederationsFederationResolvers), DELETE: int64Wrap(delFederationsFederationResolvers)},
		"federation_resolver":                  ApiHandlerFuncMap{GET: emptyWrap(getFederationResolverss), POST: bodyWrap(postFederationResolvers)},
		"federation_resolver/{key}":            ApiHandlerFuncMap{GET: int64Wrap(getFederationResolversById), PUT: int64BodyWrap(putFederationResolvers), DELETE: int64Wrap(delFederationResolvers)},
		"federation_tmuser":                    ApiHandlerFuncMap{GET: emptyWrap(getFederationUserss), POST: bodyWrap(postFederationUsers)},
		"federation_tmuser/{key}":              ApiHandlerFuncMap{GET: int64Wrap(getFederationUsersById), PUT: int64BodyWrap(putFederationUsers), DELETE: int64Wrap(delFederationUsers)},
		// "job":                    ApiHandlerFuncMap{GET: emptyWrap(getJobss), POST: bodyWrap(postJobs)},
		// "job/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getJobsById), PUT: int64BodyWrap(putJobs), DELETE: int64Wrap(delJobs)},
		// "job_result":             ApiHandlerFuncMap{GET: emptyWrap(getJobResultss), POST: bodyWrap(postJobResults)},
		// "job_result/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobResultsById), PUT: int64BodyWrap(putJobResults), DELETE: int64Wrap(delJobResults)},
		// "job_status":             ApiHandlerFuncMap{GET: emptyWrap(getJobStatusess), POST: bodyWrap(postJobStatuses)},
		// "job_status/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobStatusesById), PUT: int64BodyWrap(putJobStatuses), DELETE: int64Wrap(delJobStatuses)},
		// "log":                    ApiHandlerFuncMap{GET: emptyWrap(getLogss), POST: bodyWrap(postLogs)},
		// "log/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getLogsById), PUT: int64BodyWrap(putLogs), DELETE: int64Wrap(delLogs)},
		"parameter":               ApiHandlerFuncMap{GET: emptyWrap(getParameterss), POST: bodyWrap(postParameters)},
		"parameter/{key}":         ApiHandlerFuncMap{GET: int64Wrap(getParametersById), PUT: int64BodyWrap(putParameters), DELETE: int64Wrap(delParameters)},
		"phys_location":           ApiHandlerFuncMap{GET: emptyWrap(getPhysLocationss), POST: bodyWrap(postPhysLocations)},
		"phys_location/{key}":     ApiHandlerFuncMap{GET: stringWrap(getPhysLocationsById), PUT: stringBodyWrap(putPhysLocations), DELETE: stringWrap(delPhysLocations)},
		"profile":                 ApiHandlerFuncMap{GET: emptyWrap(getProfiless), POST: bodyWrap(postProfiles)},
		"profile/{key}":           ApiHandlerFuncMap{GET: stringWrap(getProfilesById), PUT: stringBodyWrap(putProfiles), DELETE: stringWrap(delProfiles)},
		"profile_parameter":       ApiHandlerFuncMap{GET: emptyWrap(getProfilesParameterss), POST: bodyWrap(postProfilesParameters)},
		"profile_parameter/{key}": ApiHandlerFuncMap{GET: stringWrap(getProfilesParametersById), PUT: stringBodyWrap(putProfilesParameters), DELETE: stringWrap(delProfilesParameters)},
		"regex":                   ApiHandlerFuncMap{GET: emptyWrap(getRegexess), POST: bodyWrap(postRegexes)},
		"regex/{key}":             ApiHandlerFuncMap{GET: int64Wrap(getRegexesById), PUT: int64BodyWrap(putRegexes), DELETE: int64Wrap(delRegexes)},
		"region":                  ApiHandlerFuncMap{GET: emptyWrap(getRegionss), POST: bodyWrap(postRegions)},
		"region/{key}":            ApiHandlerFuncMap{GET: stringWrap(getRegionsById), PUT: stringBodyWrap(putRegions), DELETE: stringWrap(delRegions)},
		"role":                    ApiHandlerFuncMap{GET: emptyWrap(getRoless), POST: bodyWrap(postRoles)},
		"role/{key}":              ApiHandlerFuncMap{GET: stringWrap(getRolesById), PUT: stringBodyWrap(putRoles), DELETE: stringWrap(delRoles)},
		"server":                  ApiHandlerFuncMap{GET: emptyWrap(getServerss), POST: bodyWrap(postServers)},
		"server/{key}":            ApiHandlerFuncMap{GET: stringWrap(GetServersById), PUT: stringBodyWrap(putServers), DELETE: stringWrap(delServers)},
		// "servercheck":            ApiHandlerFuncMap{GET: emptyWrap(getServercheckss), POST: bodyWrap(postServerchecks)},
		// "servercheck/{key}":       ApiHandlerFuncMap{GET: int64Wrap(getServerchecksById), PUT: int64BodyWrap(putServerchecks), DELETE: int64Wrap(delServerchecks)},
		"staticdnsentry":       ApiHandlerFuncMap{GET: emptyWrap(getStaticdnsentriess), POST: bodyWrap(postStaticdnsentries)},
		"staticdnsentry/{key}": ApiHandlerFuncMap{GET: int64Wrap(getStaticdnsentriesById), PUT: int64BodyWrap(putStaticdnsentries), DELETE: int64Wrap(delStaticdnsentries)},
		"stats_summary":        ApiHandlerFuncMap{GET: emptyWrap(getStatsSummarys), POST: bodyWrap(postStatsSummary)},
		"stats_summary/{key}":  ApiHandlerFuncMap{GET: stringWrap(getStatsSummaryById), PUT: stringBodyWrap(putStatsSummary), DELETE: stringWrap(delStatsSummary)},
		"status":               ApiHandlerFuncMap{GET: emptyWrap(getStatusess), POST: bodyWrap(postStatuses)},
		"status/{key}":         ApiHandlerFuncMap{GET: stringWrap(getStatusesById), PUT: stringBodyWrap(putStatuses), DELETE: stringWrap(delStatuses)},
		"user":                 ApiHandlerFuncMap{GET: emptyWrap(getUserss), POST: bodyWrap(postUsers)},
		"user/{key}":           ApiHandlerFuncMap{GET: stringWrap(GetUsersById), PUT: stringBodyWrap(putUsers), DELETE: stringWrap(delUsers)},
		"extension":            ApiHandlerFuncMap{GET: emptyWrap(getExtensionss), POST: bodyWrap(postExtensions)},
		"extension/{key}":      ApiHandlerFuncMap{GET: stringWrap(getExtensionsById), PUT: stringBodyWrap(putExtensions), DELETE: stringWrap(delExtensions)},
		// "type":                   ApiHandlerFuncMap{GET: emptyWrap(getTypess), POST: bodyWrap(postTypes)},
		// "type/{key}":              ApiHandlerFuncMap{GET: int64Wrap(getTypesById), PUT: int64BodyWrap(putTypes), DELETE: int64Wrap(delTypes)},
		//		"snapshot/crconfig/{key}": ApiHandlerFuncMap{GET: stringWrap(snapshotCrconfigs)},
	}
}

type EmptyHandlerFunc func(db *sqlx.DB) (interface{}, error)
type Int64HandlerFunc func(key int64, db *sqlx.DB) (interface{}, error)
type BodyHandlerFunc func(payload []byte, db *sqlx.DB) (interface{}, error)
type Int64BodyHandlerFunc func(key int64, payload []byte, db *sqlx.DB) (interface{}, error)
type StringHandlerFunc func(key string, db *sqlx.DB) (interface{}, error)
type StringBodyHandlerFunc func(key string, payload []byte, db *sqlx.DB) (interface{}, error)

func int64BodyWrap(f Int64BodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		if strKey, ok := pathParams["key"]; !ok {
			return nil, errors.New("int64 key missing")
		} else if key, err := strconv.Atoi(strKey); err != nil {
			return nil, errors.New("key is not an int64: " + strKey)
		} else {
			return f(int64(key), payload, db)
		}
	}
}

func int64Wrap(f Int64HandlerFunc) ApiHandlerFunc {
	return int64BodyWrap(func(key int64, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key, db)
	})
}

func bodyWrap(f BodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(payload, db)
	}
}

func stringBodyWrap(f StringBodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		if key, ok := pathParams["key"]; !ok {
			return nil, errors.New("string key missing")
		} else {
			return f(key, payload, db)
		}
	}
}

func stringWrap(f StringHandlerFunc) ApiHandlerFunc {
	return stringBodyWrap(func(key string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key, db)
	})
}

func emptyWrap(f EmptyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(db)
	}
}
