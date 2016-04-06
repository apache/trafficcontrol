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
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
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
		"cdn":                                                          ApiHandlerFuncMap{GET: emptyWrap(getCdnss), POST: bodyWrap(postCdns)},
		"cdn/{key}":                                                    ApiHandlerFuncMap{GET: stringWrap(getCdnsById), PUT: stringBodyWrap(putCdns), DELETE: stringWrap(delCdns)},
		"asn":                                                          ApiHandlerFuncMap{GET: emptyWrap(getAsnss), POST: bodyWrap(postAsns)},
		"asn/{key}":                                                    ApiHandlerFuncMap{GET: int64Wrap(getAsnsById), PUT: int64BodyWrap(putAsns), DELETE: int64Wrap(delAsns)},
		"cachegroup":                                                   ApiHandlerFuncMap{GET: emptyWrap(getCachegroupss), POST: bodyWrap(postCachegroups)},
		"cachegroup/{key}":                                             ApiHandlerFuncMap{GET: stringWrap(getCachegroupsById), PUT: stringBodyWrap(putCachegroups), DELETE: stringWrap(delCachegroups)},
		"cachegroup_parameter":                                         ApiHandlerFuncMap{GET: emptyWrap(getCachegroupsParameterss), POST: bodyWrap(postCachegroupsParameters)},
		"cachegroup_parameter/cachegroup/{key0}/parameter_id/{key1}":   ApiHandlerFuncMap{GET: stringInt64Wrap(getCachegroupsParametersById), PUT: stringInt64BodyWrap(putCachegroupsParameters), DELETE: stringInt64Wrap(delCachegroupsParameters)},
		"deliveryservice":                                              ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicess), POST: bodyWrap(postDeliveryservices)},
		"deliveryservice/{key}":                                        ApiHandlerFuncMap{GET: stringWrap(getDeliveryservicesById), PUT: stringBodyWrap(putDeliveryservices), DELETE: stringWrap(delDeliveryservices)},
		"deliveryservice_regex":                                        ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesRegexess), POST: bodyWrap(postDeliveryservicesRegexes)},
		"deliveryservice_regex/deliveryservice/{key0}/regex_id/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(getDeliveryservicesRegexesById), PUT: stringInt64BodyWrap(putDeliveryservicesRegexes), DELETE: stringInt64Wrap(delDeliveryservicesRegexes)},
		"deliveryservice_server":                                       ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesServerss), POST: bodyWrap(postDeliveryservicesServers)},
		"deliveryservice_server/deliveryservice/{key0}/server/{key1}":  ApiHandlerFuncMap{GET: stringStringWrap(getDeliveryservicesServersById), PUT: stringStringBodyWrap(putDeliveryservicesServers), DELETE: stringStringWrap(delDeliveryservicesServers)},
		"deliveryservice_tmuser":                                       ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesUserss), POST: bodyWrap(postDeliveryservicesUsers)},
		"deliveryservice_user/deliveryservice/{key0}/user/{key1}":      ApiHandlerFuncMap{GET: stringStringWrap(getDeliveryservicesUsersById), PUT: stringStringBodyWrap(putDeliveryservicesUsers), DELETE: stringStringWrap(delDeliveryservicesUsers)},
		"division":                                                                       ApiHandlerFuncMap{GET: emptyWrap(getDivisionss), POST: bodyWrap(postDivisions)},
		"division/{key}":                                                                 ApiHandlerFuncMap{GET: stringWrap(getDivisionsById), PUT: stringBodyWrap(putDivisions), DELETE: stringWrap(delDivisions)},
		"federation":                                                                     ApiHandlerFuncMap{GET: emptyWrap(getFederationss), POST: bodyWrap(postFederations)},
		"federation/{key}":                                                               ApiHandlerFuncMap{GET: int64Wrap(getFederationsById), PUT: int64BodyWrap(putFederations), DELETE: int64Wrap(delFederations)},
		"federation_deliveryservice":                                                     ApiHandlerFuncMap{GET: emptyWrap(getFederationsDeliveryservicess), POST: bodyWrap(postFederationsDeliveryservices)},
		"federation_deliveryservice/federation_id/{key0}/deliveryservice/{key1}":         ApiHandlerFuncMap{GET: int64StringWrap(getFederationsDeliveryservicesById), PUT: int64StringBodyWrap(putFederationsDeliveryservices), DELETE: int64StringWrap(delFederationsDeliveryservices)},
		"federation_federation_resolver":                                                 ApiHandlerFuncMap{GET: emptyWrap(getFederationsFederationResolverss), POST: bodyWrap(postFederationsFederationResolvers)},
		"federation_federation_resolver/federation_id/{key0}/federation_resolver/{key1}": ApiHandlerFuncMap{GET: int64Int64Wrap(getFederationsFederationResolversById), PUT: int64Int64BodyWrap(putFederationsFederationResolvers), DELETE: int64Int64Wrap(delFederationsFederationResolvers)},
		"federation_resolver":                                                            ApiHandlerFuncMap{GET: emptyWrap(getFederationResolverss), POST: bodyWrap(postFederationResolvers)},
		"federation_resolver/{key}":                                                      ApiHandlerFuncMap{GET: int64Wrap(getFederationResolversById), PUT: int64BodyWrap(putFederationResolvers), DELETE: int64Wrap(delFederationResolvers)},
		"federation_user":                                                                ApiHandlerFuncMap{GET: emptyWrap(getFederationUserss), POST: bodyWrap(postFederationUsers)},
		"federation_user/federation_id/{key0}/user/{key1}":                               ApiHandlerFuncMap{GET: int64StringWrap(getFederationUsersById), PUT: int64StringBodyWrap(putFederationUsers), DELETE: int64StringWrap(delFederationUsers)},
		// "job":                    ApiHandlerFuncMap{GET: emptyWrap(getJobss), POST: bodyWrap(postJobs)},
		// "job/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getJobsById), PUT: int64BodyWrap(putJobs), DELETE: int64Wrap(delJobs)},
		// "job_result":             ApiHandlerFuncMap{GET: emptyWrap(getJobResultss), POST: bodyWrap(postJobResults)},
		// "job_result/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobResultsById), PUT: int64BodyWrap(putJobResults), DELETE: int64Wrap(delJobResults)},
		// "job_status":             ApiHandlerFuncMap{GET: emptyWrap(getJobStatusess), POST: bodyWrap(postJobStatuses)},
		// "job_status/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobStatusesById), PUT: int64BodyWrap(putJobStatuses), DELETE: int64Wrap(delJobStatuses)},
		// "log":                    ApiHandlerFuncMap{GET: emptyWrap(getLogss), POST: bodyWrap(postLogs)},
		// "log/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getLogsById), PUT: int64BodyWrap(putLogs), DELETE: int64Wrap(delLogs)},
		"parameter":                                         ApiHandlerFuncMap{GET: emptyWrap(getParameterss), POST: bodyWrap(postParameters)},
		"parameter/{key}":                                   ApiHandlerFuncMap{GET: int64Wrap(getParametersById), PUT: int64BodyWrap(putParameters), DELETE: int64Wrap(delParameters)},
		"phys_location":                                     ApiHandlerFuncMap{GET: emptyWrap(getPhysLocationss), POST: bodyWrap(postPhysLocations)},
		"phys_location/{key}":                               ApiHandlerFuncMap{GET: stringWrap(getPhysLocationsById), PUT: stringBodyWrap(putPhysLocations), DELETE: stringWrap(delPhysLocations)},
		"profile":                                           ApiHandlerFuncMap{GET: emptyWrap(getProfiless), POST: bodyWrap(postProfiles)},
		"profile/{key}":                                     ApiHandlerFuncMap{GET: stringWrap(getProfilesById), PUT: stringBodyWrap(putProfiles), DELETE: stringWrap(delProfiles)},
		"profile_parameter":                                 ApiHandlerFuncMap{GET: emptyWrap(getProfilesParameterss), POST: bodyWrap(postProfilesParameters)},
		"profile_parameter/profile/{key0}/parameter/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(getProfilesParametersById), PUT: stringInt64BodyWrap(putProfilesParameters), DELETE: stringInt64Wrap(delProfilesParameters)},
		"regex":        ApiHandlerFuncMap{GET: emptyWrap(getRegexess), POST: bodyWrap(postRegexes)},
		"regex/{key}":  ApiHandlerFuncMap{GET: int64Wrap(getRegexesById), PUT: int64BodyWrap(putRegexes), DELETE: int64Wrap(delRegexes)},
		"region":       ApiHandlerFuncMap{GET: emptyWrap(getRegionss), POST: bodyWrap(postRegions)},
		"region/{key}": ApiHandlerFuncMap{GET: stringWrap(getRegionsById), PUT: stringBodyWrap(putRegions), DELETE: stringWrap(delRegions)},
		"role":         ApiHandlerFuncMap{GET: emptyWrap(getRoless), POST: bodyWrap(postRoles)},
		"role/{key}":   ApiHandlerFuncMap{GET: stringWrap(getRolesById), PUT: stringBodyWrap(putRoles), DELETE: stringWrap(delRoles)},
		"server":       ApiHandlerFuncMap{GET: emptyWrap(getServerss), POST: bodyWrap(postServers)},
		"server/host_name/{key0}/tcp_port/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(GetServersById), PUT: stringInt64BodyWrap(putServers), DELETE: stringInt64Wrap(delServers)},
		// "servercheck":            ApiHandlerFuncMap{GET: emptyWrap(getServercheckss), POST: bodyWrap(postServerchecks)},
		// "servercheck/{key}":       ApiHandlerFuncMap{GET: int64Wrap(getServerchecksById), PUT: int64BodyWrap(putServerchecks), DELETE: int64Wrap(delServerchecks)},
		"staticdnsentry":       ApiHandlerFuncMap{GET: emptyWrap(getStaticdnsentriess), POST: bodyWrap(postStaticdnsentries)},
		"staticdnsentry/{key}": ApiHandlerFuncMap{GET: int64Wrap(getStaticdnsentriesById), PUT: int64BodyWrap(putStaticdnsentries), DELETE: int64Wrap(delStaticdnsentries)},
		"stats_summary":        ApiHandlerFuncMap{GET: emptyWrap(getStatsSummarys), POST: bodyWrap(postStatsSummary)},
		"stats_summary/cdn_name/{key0}/deliveryservice/{key1}/stat_name/{key2}/stat_date/{key3}": ApiHandlerFuncMap{GET: stringStringStringTimeWrap(getStatsSummaryById), PUT: stringStringStringTimeBodyWrap(putStatsSummary), DELETE: stringStringStringTimeWrap(delStatsSummary)},
		"status":          ApiHandlerFuncMap{GET: emptyWrap(getStatusess), POST: bodyWrap(postStatuses)},
		"status/{key}":    ApiHandlerFuncMap{GET: stringWrap(getStatusesById), PUT: stringBodyWrap(putStatuses), DELETE: stringWrap(delStatuses)},
		"user":            ApiHandlerFuncMap{GET: emptyWrap(getUserss), POST: bodyWrap(postUsers)},
		"user/{key}":      ApiHandlerFuncMap{GET: stringWrap(GetUsersById), PUT: stringBodyWrap(putUsers), DELETE: stringWrap(delUsers)},
		"extension":       ApiHandlerFuncMap{GET: emptyWrap(getExtensionss), POST: bodyWrap(postExtensions)},
		"extension/{key}": ApiHandlerFuncMap{GET: stringWrap(getExtensionsById), PUT: stringBodyWrap(putExtensions), DELETE: stringWrap(delExtensions)},
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
type StringInt64HandlerFunc func(key0 string, key1 int64, db *sqlx.DB) (interface{}, error)
type StringInt64BodyHandlerFunc func(key0 string, key1 int64, payload []byte, db *sqlx.DB) (interface{}, error)
type Int64StringHandlerFunc func(key0 int64, key1 string, db *sqlx.DB) (interface{}, error)
type Int64StringBodyHandlerFunc func(key0 int64, key1 string, payload []byte, db *sqlx.DB) (interface{}, error)
type StringStringHandlerFunc func(key0 string, key1 string, db *sqlx.DB) (interface{}, error)
type StringStringBodyHandlerFunc func(key0 string, key1 string, payload []byte, db *sqlx.DB) (interface{}, error)
type Int64Int64HandlerFunc func(key0 int64, key1 int64, db *sqlx.DB) (interface{}, error)
type Int64Int64BodyHandlerFunc func(key0 int64, key1 int64, payload []byte, db *sqlx.DB) (interface{}, error)
type StringStringStringTimeHandlerFunc func(key0 string, key1 string, key2 string, key3 time.Time, db *sqlx.DB) (interface{}, error)
type StringStringStringTimeBodyHandlerFunc func(key0 string, key1 string, key2 string, key3 time.Time, payload []byte, db *sqlx.DB) (interface{}, error)

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

func stringInt64BodyWrap(f StringInt64BodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		var key0 string
		var key1 int
		var err error
		var ok bool
		if key0, ok = pathParams["key0"]; !ok {
			return nil, errors.New("string key missing")
		}
		if strKey1, ok := pathParams["key1"]; !ok {
			return nil, errors.New("second int64 key missing")
		} else if key1, err = strconv.Atoi(strKey1); err != nil {
			return nil, errors.New("second key is not an int64: " + strKey1)
		}
		return f(key0, int64(key1), payload, db)
	}
}

func stringInt64Wrap(f StringInt64HandlerFunc) ApiHandlerFunc {
	return stringInt64BodyWrap(func(key0 string, key1 int64, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key0, key1, db)
	})
}

func int64StringBodyWrap(f Int64StringBodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		var key0 int
		var key1 string
		var err error
		var ok bool
		if strKey0, ok := pathParams["key0"]; !ok {
			return nil, errors.New("int64 key missing")
		} else if key0, err = strconv.Atoi(strKey0); err != nil {
			return nil, errors.New("key is not an int64: " + strKey0)
		}
		if key1, ok = pathParams["key1"]; !ok {
			return nil, errors.New("second string key missing")
		}
		return f(int64(key0), key1, payload, db)
	}
}

func int64StringWrap(f Int64StringHandlerFunc) ApiHandlerFunc {
	return int64StringBodyWrap(func(key0 int64, key1 string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key0, key1, db)
	})
}

func stringStringBodyWrap(f StringStringBodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		var key0 string
		var key1 string
		var ok bool
		if key0, ok = pathParams["key0"]; !ok {
			return nil, errors.New("string key missing")
		}
		if key1, ok = pathParams["key1"]; !ok {
			return nil, errors.New("second string key missing")
		}
		return f(key0, key1, payload, db)
	}
}

func stringStringWrap(f StringStringHandlerFunc) ApiHandlerFunc {
	return stringStringBodyWrap(func(key0 string, key1 string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key0, key1, db)
	})
}

func int64Int64BodyWrap(f Int64Int64BodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		var key0 int
		var key1 int
		var err error
		if strKey, ok := pathParams["key0"]; !ok {
			return nil, errors.New("int64 key missing")
		} else if key0, err = strconv.Atoi(strKey); err != nil {
			return nil, errors.New("key is not an int64: " + strKey)
		}
		if strKey, ok := pathParams["key1"]; !ok {
			return nil, errors.New("second int64 key missing")
		} else if key1, err = strconv.Atoi(strKey); err != nil {
			return nil, errors.New("second key is not an int64: " + strKey)
		}
		return f(int64(key0), int64(key1), payload, db)
	}
}

func int64Int64Wrap(f Int64Int64HandlerFunc) ApiHandlerFunc {
	return int64Int64BodyWrap(func(key0 int64, key1 int64, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key0, key1, db)
	})
}

func stringStringStringTimeBodyWrap(f StringStringStringTimeBodyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		var key0 string
		var key1 string
		var key2 string
		var key3 time.Time
		var err error
		var ok bool
		if key0, ok = pathParams["key0"]; !ok {
			return nil, errors.New("string key missing")
		}
		if key1, ok = pathParams["key1"]; !ok {
			return nil, errors.New("second string key missing")
		}
		if key2, ok = pathParams["key2"]; !ok {
			return nil, errors.New("third string key missing")
		}
		if key3Str, ok := pathParams["key3"]; !ok {
			return nil, errors.New("fourth time key missing")
		} else if key3, err = time.Parse(time.RFC3339Nano, key3Str); err != nil {
			return nil, fmt.Errorf("Fourth time key is not a valid RFC 3339 time: %v", err)
		}
		return f(key0, key1, key2, key3, payload, db)
	}
}

func stringStringStringTimeWrap(f StringStringStringTimeHandlerFunc) ApiHandlerFunc {
	return stringStringStringTimeBodyWrap(func(key0 string, key1 string, key2 string, key3 time.Time, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(key0, key1, key2, key3, db)
	})
}

func emptyWrap(f EmptyHandlerFunc) ApiHandlerFunc {
	return func(pathParams map[string]string, payload []byte, db *sqlx.DB) (interface{}, error) {
		return f(db)
	}
}
