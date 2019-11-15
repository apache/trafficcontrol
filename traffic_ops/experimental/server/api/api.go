
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
		"cdn":                                                          ApiHandlerFuncMap{GET: emptyWrap(getCdns), POST: bodyWrap(postCdn)},
		"cdn/{key}":                                                    ApiHandlerFuncMap{GET: stringWrap(getCdn), PUT: stringBodyWrap(putCdn), DELETE: stringWrap(delCdn)},
		"asn":                                                          ApiHandlerFuncMap{GET: emptyWrap(getAsns), POST: bodyWrap(postAsn)},
		"asn/{key}":                                                    ApiHandlerFuncMap{GET: int64Wrap(getAsn), PUT: int64BodyWrap(putAsn), DELETE: int64Wrap(delAsn)},
		"cachegroup":                                                   ApiHandlerFuncMap{GET: emptyWrap(getCachegroups), POST: bodyWrap(postCachegroup)},
		"cachegroup/{key}":                                             ApiHandlerFuncMap{GET: stringWrap(getCachegroup), PUT: stringBodyWrap(putCachegroup), DELETE: stringWrap(delCachegroup)},
		"cachegroup_parameter":                                         ApiHandlerFuncMap{GET: emptyWrap(getCachegroupsParameters), POST: bodyWrap(postCachegroupsParameter)},
		"cachegroup_parameter/cachegroup/{key0}/parameter_id/{key1}":   ApiHandlerFuncMap{GET: stringInt64Wrap(getCachegroupsParameter), PUT: stringInt64BodyWrap(putCachegroupsParameter), DELETE: stringInt64Wrap(delCachegroupsParameter)},
		"deliveryservice":                                              ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservices), POST: bodyWrap(postDeliveryservice)},
		"deliveryservice/{key}":                                        ApiHandlerFuncMap{GET: stringWrap(getDeliveryservice), PUT: stringBodyWrap(putDeliveryservice), DELETE: stringWrap(delDeliveryservice)},
		"deliveryservice_regex":                                        ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesRegexes), POST: bodyWrap(postDeliveryservicesRegex)},
		"deliveryservice_regex/deliveryservice/{key0}/regex_id/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(getDeliveryservicesRegex), PUT: stringInt64BodyWrap(putDeliveryservicesRegex), DELETE: stringInt64Wrap(delDeliveryservicesRegex)},
		"deliveryservice_server":                                       ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesServers), POST: bodyWrap(postDeliveryservicesServer)},
		"deliveryservice_server/deliveryservice/{key0}/server/{key1}":  ApiHandlerFuncMap{GET: stringStringWrap(getDeliveryservicesServer), PUT: stringStringBodyWrap(putDeliveryservicesServer), DELETE: stringStringWrap(delDeliveryservicesServer)},
		"deliveryservice_tmuser":                                       ApiHandlerFuncMap{GET: emptyWrap(getDeliveryservicesUsers), POST: bodyWrap(postDeliveryservicesUser)},
		"deliveryservice_user/deliveryservice/{key0}/user/{key1}":      ApiHandlerFuncMap{GET: stringStringWrap(getDeliveryservicesUser), PUT: stringStringBodyWrap(putDeliveryservicesUser), DELETE: stringStringWrap(delDeliveryservicesUser)},
		"division":                                                                       ApiHandlerFuncMap{GET: emptyWrap(getDivisions), POST: bodyWrap(postDivision)},
		"division/{key}":                                                                 ApiHandlerFuncMap{GET: stringWrap(getDivision), PUT: stringBodyWrap(putDivision), DELETE: stringWrap(delDivision)},
		"federation":                                                                     ApiHandlerFuncMap{GET: emptyWrap(getFederations), POST: bodyWrap(postFederation)},
		"federation/{key}":                                                               ApiHandlerFuncMap{GET: int64Wrap(getFederation), PUT: int64BodyWrap(putFederation), DELETE: int64Wrap(delFederation)},
		"federation_deliveryservice":                                                     ApiHandlerFuncMap{GET: emptyWrap(getFederationsDeliveryservices), POST: bodyWrap(postFederationsDeliveryservice)},
		"federation_deliveryservice/federation_id/{key0}/deliveryservice/{key1}":         ApiHandlerFuncMap{GET: int64StringWrap(getFederationsDeliveryservice), PUT: int64StringBodyWrap(putFederationsDeliveryservice), DELETE: int64StringWrap(delFederationsDeliveryservice)},
		"federation_federation_resolver":                                                 ApiHandlerFuncMap{GET: emptyWrap(getFederationsFederationResolvers), POST: bodyWrap(postFederationsFederationResolver)},
		"federation_federation_resolver/federation_id/{key0}/federation_resolver/{key1}": ApiHandlerFuncMap{GET: int64Int64Wrap(getFederationsFederationResolver), PUT: int64Int64BodyWrap(putFederationsFederationResolver), DELETE: int64Int64Wrap(delFederationsFederationResolver)},
		"federation_resolver":                                                            ApiHandlerFuncMap{GET: emptyWrap(getFederationResolvers), POST: bodyWrap(postFederationResolver)},
		"federation_resolver/{key}":                                                      ApiHandlerFuncMap{GET: int64Wrap(getFederationResolver), PUT: int64BodyWrap(putFederationResolver), DELETE: int64Wrap(delFederationResolver)},
		"federation_user":                                                                ApiHandlerFuncMap{GET: emptyWrap(getFederationUsers), POST: bodyWrap(postFederationUser)},
		"federation_user/federation_id/{key0}/user/{key1}":                               ApiHandlerFuncMap{GET: int64StringWrap(getFederationUser), PUT: int64StringBodyWrap(putFederationUser), DELETE: int64StringWrap(delFederationUser)},
		// "job":                    ApiHandlerFuncMap{GET: emptyWrap(getJobs), POST: bodyWrap(postJob)},
		// "job/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getJob), PUT: int64BodyWrap(putJob), DELETE: int64Wrap(delJob)},
		// "job_result":             ApiHandlerFuncMap{GET: emptyWrap(getJobResults), POST: bodyWrap(postJobResult)},
		// "job_result/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobResult), PUT: int64BodyWrap(putJobResult), DELETE: int64Wrap(delJobResult)},
		// "job_status":             ApiHandlerFuncMap{GET: emptyWrap(getJobStatuses), POST: bodyWrap(postJobStatuse)},
		// "job_status/{key}":        ApiHandlerFuncMap{GET: int64Wrap(getJobStatuse), PUT: int64BodyWrap(putJobStatuse), DELETE: int64Wrap(delJobStatuse)},
		// "log":                    ApiHandlerFuncMap{GET: emptyWrap(getLogs), POST: bodyWrap(postLog)},
		// "log/{key}":               ApiHandlerFuncMap{GET: int64Wrap(getLog), PUT: int64BodyWrap(putLog), DELETE: int64Wrap(delLog)},
		"parameter":                                         ApiHandlerFuncMap{GET: emptyWrap(getParameters), POST: bodyWrap(postParameter)},
		"parameter/{key}":                                   ApiHandlerFuncMap{GET: int64Wrap(getParameter), PUT: int64BodyWrap(putParameter), DELETE: int64Wrap(delParameter)},
		"phys_location":                                     ApiHandlerFuncMap{GET: emptyWrap(getPhysLocations), POST: bodyWrap(postPhysLocation)},
		"phys_location/{key}":                               ApiHandlerFuncMap{GET: stringWrap(getPhysLocation), PUT: stringBodyWrap(putPhysLocation), DELETE: stringWrap(delPhysLocation)},
		"profile":                                           ApiHandlerFuncMap{GET: emptyWrap(getProfiles), POST: bodyWrap(postProfile)},
		"profile/{key}":                                     ApiHandlerFuncMap{GET: stringWrap(getProfile), PUT: stringBodyWrap(putProfile), DELETE: stringWrap(delProfile)},
		"profile_parameter":                                 ApiHandlerFuncMap{GET: emptyWrap(getProfilesParameters), POST: bodyWrap(postProfilesParameter)},
		"profile_parameter/profile/{key0}/parameter/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(getProfilesParameter), PUT: stringInt64BodyWrap(putProfilesParameter), DELETE: stringInt64Wrap(delProfilesParameter)},
		"regex":        ApiHandlerFuncMap{GET: emptyWrap(getRegexes), POST: bodyWrap(postRegex)},
		"regex/{key}":  ApiHandlerFuncMap{GET: int64Wrap(getRegex), PUT: int64BodyWrap(putRegex), DELETE: int64Wrap(delRegex)},
		"region":       ApiHandlerFuncMap{GET: emptyWrap(getRegions), POST: bodyWrap(postRegion)},
		"region/{key}": ApiHandlerFuncMap{GET: stringWrap(getRegion), PUT: stringBodyWrap(putRegion), DELETE: stringWrap(delRegion)},
		"role":         ApiHandlerFuncMap{GET: emptyWrap(getRoles), POST: bodyWrap(postRole)},
		"role/{key}":   ApiHandlerFuncMap{GET: stringWrap(getRole), PUT: stringBodyWrap(putRole), DELETE: stringWrap(delRole)},
		"server":       ApiHandlerFuncMap{GET: emptyWrap(getServers), POST: bodyWrap(postServer)},
		"server/host_name/{key0}/tcp_port/{key1}": ApiHandlerFuncMap{GET: stringInt64Wrap(GetServer), PUT: stringInt64BodyWrap(putServer), DELETE: stringInt64Wrap(delServer)},
		// "servercheck":            ApiHandlerFuncMap{GET: emptyWrap(getServerchecks), POST: bodyWrap(postServercheck)},
		// "servercheck/{key}":       ApiHandlerFuncMap{GET: int64Wrap(getServercheck), PUT: int64BodyWrap(putServercheck), DELETE: int64Wrap(delServercheck)},
		"staticdnsentry":       ApiHandlerFuncMap{GET: emptyWrap(getStaticdnsentries), POST: bodyWrap(postStaticdnsentry)},
		"staticdnsentry/{key}": ApiHandlerFuncMap{GET: int64Wrap(getStaticdnsentry), PUT: int64BodyWrap(putStaticdnsentry), DELETE: int64Wrap(delStaticdnsentry)},
		"stats_summary":        ApiHandlerFuncMap{GET: emptyWrap(getStatsSummaries), POST: bodyWrap(postStatsSummary)},
		"stats_summary/cdn_name/{key0}/deliveryservice/{key1}/stat_name/{key2}/stat_date/{key3}": ApiHandlerFuncMap{GET: stringStringStringTimeWrap(getStatsSummary), PUT: stringStringStringTimeBodyWrap(putStatsSummary), DELETE: stringStringStringTimeWrap(delStatsSummary)},
		"status":          ApiHandlerFuncMap{GET: emptyWrap(getStatuses), POST: bodyWrap(postStatus)},
		"status/{key}":    ApiHandlerFuncMap{GET: stringWrap(getStatus), PUT: stringBodyWrap(putStatus), DELETE: stringWrap(delStatus)},
		"user":            ApiHandlerFuncMap{GET: emptyWrap(getUsers), POST: bodyWrap(postUser)},
		"user/{key}":      ApiHandlerFuncMap{GET: stringWrap(GetUser), PUT: stringBodyWrap(putUser), DELETE: stringWrap(delUser)},
		"extension":       ApiHandlerFuncMap{GET: emptyWrap(getExtensions), POST: bodyWrap(postExtension)},
		"extension/{key}": ApiHandlerFuncMap{GET: stringWrap(getExtension), PUT: stringBodyWrap(putExtension), DELETE: stringWrap(delExtension)},
		// "type":                   ApiHandlerFuncMap{GET: emptyWrap(getTypes), POST: bodyWrap(postType)},
		// "type/{key}":              ApiHandlerFuncMap{GET: int64Wrap(getType), PUT: int64BodyWrap(putType), DELETE: int64Wrap(delType)},
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
