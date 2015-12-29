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

package todb

type actionhandler func(method string, id int, payload []byte) (interface{}, error)
type actionmap map[string]actionhandler

var funcMap = actionmap{
	"asn":                  handleAsn,
	"cachegroup":           handleCachegroup,
	"cachegroup_parameter": handleCachegroupParameter,
	"cdn":                            handleCdn,
	"deliveryservice":                handleDeliveryservice,
	"deliveryservice_regex":          handleDeliveryserviceRegex,
	"deliveryservice_server":         handleDeliveryserviceServer,
	"deliveryservice_tmuser":         handleDeliveryserviceTmuser,
	"division":                       handleDivision,
	"federation":                     handleFederation,
	"federation_deliveryservice":     handleFederationDeliveryservice,
	"federation_federation_resolver": handleFederationFederationResolver,
	"federation_resolver":            handleFederationResolver,
	"federation_tmuser":              handleFederationTmuser,
	"goose_db_version":               handleGooseDbVersion,
	"hwinfo":                         handleHwinfo,
	"job":                            handleJob,
	"job_agent":                      handleJobAgent,
	"job_result":                     handleJobResult,
	"job_status":                     handleJobStatus,
	"log":                            handleLog,
	"parameter":                      handleParameter,
	"phys_location":                  handlePhysLocation,
	"profile":                        handleProfile,
	"profile_parameter":              handleProfileParameter,
	"regex":                          handleRegex,
	"region":                         handleRegion,
	"role":                           handleRole,
	"server":                         handleServer,
	"servercheck":                    handleServercheck,
	"staticdnsentry":                 handleStaticdnsentry,
	"stats_summary":                  handleStatsSummary,
	"status":                         handleStatus,
	"tm_user":                        handleTmUser,
	"to_extension":                   handleToExtension,
	"type":                           handleType,
}

func Action(tableName, method string, id int, payload []byte) (interface{}, error) {
	return funcMap[tableName](method, id, payload)
}
