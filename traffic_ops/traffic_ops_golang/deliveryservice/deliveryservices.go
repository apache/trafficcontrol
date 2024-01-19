// Package deliveryservice contains logic and handlers for /deliveryservices and
// related sub-routes (e.g. /deliveryservices/{{ID}}/servers).
package deliveryservice

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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type tierType int

const (
	midTier tierType = iota
	edgeTier
)

// TODeliveryService is a CRUDder implementation for Delivery Services; it's
// only used for GET and DELETE requests.
type TODeliveryService struct {
	api.APIInfoImpl
	tc.DeliveryServiceV5
}

// TODeliveryServiceOldDetails is used to store the "old" details while
// updating a DS.
type TODeliveryServiceOldDetails struct {
	OldOrgServerFQDN *string
	OldCDNName       string
	OldCDNID         int
	OldRoutingName   string
	OldSSLKeyVersion *int
}

// MarshalJSON implements encoding/json.Marshaler.
func (ds TODeliveryService) MarshalJSON() ([]byte, error) {
	return json.Marshal(ds.DeliveryServiceV5)
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (ds *TODeliveryService) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &ds.DeliveryServiceV5)
}

// APIInfo implements
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.APIInfoer.
func (ds *TODeliveryService) APIInfo() *api.Info { return ds.ReqInfo }

// SetKeys implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interface.
func (ds *TODeliveryService) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int)
	ds.ID = &i
}

// GetKeys implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interface.
func (ds TODeliveryService) GetKeys() (map[string]interface{}, bool) {
	if ds.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *ds.ID}, true
}

// GetKeyFieldsInfo implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interface.
func (ds TODeliveryService) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// GetAuditName implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interface.
func (ds *TODeliveryService) GetAuditName() string {
	return ds.XMLID
}

// GetType implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interface.
func (ds *TODeliveryService) GetType() string {
	return "ds"
}

// IsTenantAuthorized implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Tenantable
// interface.
func (ds *TODeliveryService) IsTenantAuthorized(user *auth.CurrentUser) (bool, error) {
	return isTenantAuthorized(ds.ReqInfo, &ds.DeliveryServiceV5)
}

const baseTLSVersionsQuery = `SELECT ARRAY_AGG(tls_version ORDER BY tls_version) FROM deliveryservice_tls_version`

const getTLSVersionsQuery = baseTLSVersionsQuery + `
WHERE deliveryservice = $1
`

// GetDSTLSVersions retrieves the TLS versions explicitly supported by a
// Delivery Service identified by dsID. This will panic if handed a nil
// transaction.
func GetDSTLSVersions(dsID int, tx *sql.Tx) ([]string, error) {
	var vers []string
	err := tx.QueryRow(getTLSVersionsQuery, dsID).Scan(pq.Array(&vers))
	if err != nil {
		err = fmt.Errorf("querying: %w", err)
	}
	return vers, err
}

// CreateV30 is used to handle POST requests to create a Delivery Service at
// version 3.0 of the Traffic Ops API.
func CreateV30(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceV30{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("decoding: %w", err), nil)
		return
	}

	res, status, userErr, sysErr := createV30(w, r, inf, ds)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service creation was successful", []tc.DeliveryServiceV30{*res})
}

// CreateV31 is used to handle POST requests to create a Delivery Service at
// version 3.1 of the Traffic Ops API.
func CreateV31(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceV31{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("decoding: %w", err), nil)
		return
	}

	res, status, userErr, sysErr := createV31(w, r, inf, ds)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service creation was successful", []tc.DeliveryServiceV31{*res})
}

// CreateV50 is used to handle POST requests to create a Delivery Service at
// version 5.0 of the Traffic Ops API.
func CreateV50(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var ds tc.DeliveryServiceV5
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("decoding: %w", err), nil)
		return
	}
	res, status, userErr, sysErr := createV50(w, r, inf, ds, true, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service creation was successful")

	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/deliveryservices?id=%d", inf.Version, *res.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, *res)
}

// CreateV40 is used to handle POST requests to create a Delivery Service at
// version 4.0 of the Traffic Ops API (and isomorphic API versions thereof, with
// respect to Delivery Service representations).
func CreateV40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceV40{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("decoding: %w", err), nil)
		return
	}
	res, status, userErr, sysErr := createV40(w, r, inf, ds, true)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service creation was successful")

	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/deliveryservices?id=%d", inf.Version, *res.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, []tc.DeliveryServiceV40{*res})
}

// CreateV41 is a handler for POST requests to create Delivery Services in
// version 4.1 of the Traffic Ops API.
func CreateV41(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var ds tc.DeliveryServiceV41
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("decoding: %w", err), nil)
		return
	}
	res, status, userErr, sysErr := createV41(w, r, inf, ds, true)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service creation was successful")

	w.Header().Set("Location", fmt.Sprintf("/api/4.0/deliveryservices?id=%d", *res.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, []tc.DeliveryServiceV41{*res})
}

func createV30(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV30 tc.DeliveryServiceV30) (*tc.DeliveryServiceV30, int, error, error) {
	ds := tc.DeliveryServiceV31{DeliveryServiceV30: dsV30}
	res, status, userErr, sysErr := createV31(w, r, inf, ds)
	if res != nil {
		return &res.DeliveryServiceV30, status, userErr, sysErr
	}
	return nil, status, userErr, sysErr
}

func createV31(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV31 tc.DeliveryServiceV31) (*tc.DeliveryServiceV31, int, error, error) {
	tx := inf.Tx.Tx
	dsNullable := tc.DeliveryServiceNullableV30(dsV31)
	ds := dsNullable.UpgradeToV4()
	res, status, userErr, sysErr := createV40(w, r, inf, ds.DeliveryServiceV40, false)
	if res == nil {
		return nil, status, userErr, sysErr
	}

	ds = tc.DeliveryServiceV41{DeliveryServiceV40: *res}
	if dsV31.CacheURL != nil {
		_, err := tx.Exec("UPDATE deliveryservice SET cacheurl = $1 WHERE ID = $2",
			dsV31.CacheURL,
			ds.ID,
		)
		if err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			return nil, code, usrErr, sysErr
		}
	}

	if err := EnsureCacheURLParams(tx, *ds.ID, *ds.XMLID, dsV31.CacheURL); err != nil {
		return nil, http.StatusInternalServerError, nil, err
	}

	oldRes := tc.DeliveryServiceV31(ds.DowngradeToV31())
	return &oldRes, status, userErr, sysErr
}

// 'ON CONFLICT DO NOTHING' should be unnecessary because all data should be
// dumped from the table before re-insertion, but it's also harmless because
// the only conflict that could occur is a fully duplicate row, which is fine
// since we're intending to create that data anyway. Although it is weird.
const insertTLSVersionsQuery = `
INSERT INTO public.deliveryservice_tls_version (deliveryservice, tls_version)
	SELECT
		$1 AS deliveryservice,
		UNNEST($2::text[]) AS tls_version
ON CONFLICT DO NOTHING
`

func recreateTLSVersions(versions []string, dsid int, tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM public.deliveryservice_tls_version WHERE deliveryservice = $1`, dsid)
	if err != nil {
		return fmt.Errorf("cleaning up existing TLS version for DS #%d: %w", dsid, err)
	}

	if len(versions) < 1 {
		return nil
	}

	_, err = tx.Exec(insertTLSVersionsQuery, dsid, pq.Array(versions))
	if err != nil {
		return fmt.Errorf("inserting new TLS versions: %w", err)
	}
	return nil
}

// createV40 creates the given ds in the database, and returns the DS with its id and other fields created on insert set. On error, the HTTP status code, user error, and system error are returned. The status code SHOULD NOT be used, if both errors are nil.
func createV40(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV4 tc.DeliveryServiceV40, omitExtraLongDescFields bool) (*tc.DeliveryServiceV40, int, error, error) {
	ds, code, userErr, sysErr := createV41(w, r, inf, tc.DeliveryServiceV41{DeliveryServiceV40: dsV4}, omitExtraLongDescFields)
	if userErr != nil || sysErr != nil || ds == nil {
		return nil, code, userErr, sysErr
	}
	d := ds.DeliveryServiceV40
	return &d, code, nil, nil
}

// createV41 creates the given Delivery Service in the database, and returns a
// reference to the Delivery Service with its ID and other fields which are
// created on insert set. On error, an HTTP status code, user error, and system
// error are returned. The status code SHOULD NOT be used, if both errors are
// nil.
func createV41(w http.ResponseWriter, r *http.Request, inf *api.Info, ds tc.DeliveryServiceV41, omitExtraLongDescFields bool) (*tc.DeliveryServiceV41, int, error, error) {
	res, code, userErr, sysErr := createV50(w, r, inf, ds.Upgrade(), omitExtraLongDescFields, ds.LongDesc1, ds.LongDesc2)
	if res != nil {
		ds := res.Downgrade()
		return &ds, code, userErr, sysErr
	}
	return nil, code, userErr, sysErr
}

// createV50 creates the given Delivery Service in the database, and returns a
// reference to the Delivery Service with its ID and other fields which are
// created on insert set. On error, an HTTP status code, user error, and system
// error are returned. The status code SHOULD NOT be used, if both errors are
// nil.
func createV50(w http.ResponseWriter, r *http.Request, inf *api.Info, ds tc.DeliveryServiceV5, omitExtraLongDescFields bool, longDesc1, longDesc2 *string) (*tc.DeliveryServiceV5, int, error, error) {
	var err error
	tx := inf.Tx.Tx
	userErr, sysErr := Validate(tx, &ds)
	if userErr != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request: %w", userErr), nil
	}
	if sysErr != nil {
		return nil, http.StatusInternalServerError, nil, sysErr
	}
	if authorized, err := isTenantAuthorized(inf, &ds); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("checking tenant: %w", err)
	} else if !authorized {
		return nil, http.StatusForbidden, errors.New("not authorized on this tenant"), nil
	}

	// TODO change DeepCachingType to implement sql.Valuer and sql.Scanner, so sqlx struct scan can be used.
	deepCachingType := tc.DeepCachingType("").String()
	if ds.DeepCachingType != "" {
		deepCachingType = ds.DeepCachingType.String() // necessary, because DeepCachingType's default needs to insert the string, not "", and Query doesn't call .String().
	}

	if errCode, userErr, sysErr := dbhelpers.CheckTopology(inf.Tx, ds); userErr != nil || sysErr != nil {
		return nil, errCode, userErr, sysErr
	}

	userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(inf.Tx.Tx, int64(ds.CDNID), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		return nil, errCode, userErr, sysErr
	}

	geoLimitCountries := strings.Join(ds.GeoLimitCountries, ",")
	var resultRows *sql.Rows
	if omitExtraLongDescFields {
		if longDesc1 != nil || longDesc2 != nil {
			return nil, http.StatusBadRequest, errors.New("the longDesc1 and longDesc2 fields are no longer supported in API 4.0 onwards"), nil
		}
		resultRows, err = tx.Query(insertQueryWithoutLD1AndLD2(),
			ds.Active,
			ds.AnonymousBlockingEnabled,
			ds.CCRDNSTTL,
			ds.CDNID,
			ds.CheckPath,
			ds.ConsistentHashRegex,
			deepCachingType,
			ds.DisplayName,
			ds.DNSBypassCNAME,
			ds.DNSBypassIP,
			ds.DNSBypassIP6,
			ds.DNSBypassTTL,
			ds.DSCP,
			ds.EdgeHeaderRewrite,
			ds.GeoLimitRedirectURL,
			ds.GeoLimit,
			geoLimitCountries,
			ds.GeoProvider,
			ds.GlobalMaxMBPS,
			ds.GlobalMaxTPS,
			ds.FQPacingRate,
			ds.HTTPBypassFQDN,
			ds.InfoURL,
			ds.InitialDispersion,
			ds.IPV6RoutingEnabled,
			ds.LogsEnabled,
			ds.LongDesc,
			ds.MaxDNSAnswers,
			ds.MaxOriginConnections,
			ds.MidHeaderRewrite,
			ds.MissLat,
			ds.MissLong,
			ds.MultiSiteOrigin,
			ds.OriginShield,
			ds.ProfileID,
			ds.Protocol,
			ds.QStringIgnore,
			ds.RangeRequestHandling,
			ds.RegexRemap,
			ds.Regional,
			ds.RegionalGeoBlocking,
			ds.RemapText,
			ds.RoutingName,
			ds.SigningAlgorithm,
			ds.SSLKeyVersion,
			ds.TenantID,
			ds.Topology,
			ds.TRRequestHeaders,
			ds.TRResponseHeaders,
			ds.TypeID,
			ds.XMLID,
			ds.EcsEnabled,
			ds.RangeSliceBlockSize,
			ds.FirstHeaderRewrite,
			ds.InnerHeaderRewrite,
			ds.LastHeaderRewrite,
			ds.ServiceCategory,
			ds.MaxRequestHeaderBytes,
			pq.Array(ds.RequiredCapabilities),
		)
	} else {
		resultRows, err = tx.Query(insertQuery(),
			ds.Active,
			ds.AnonymousBlockingEnabled,
			ds.CCRDNSTTL,
			ds.CDNID,
			ds.CheckPath,
			ds.ConsistentHashRegex,
			deepCachingType,
			ds.DisplayName,
			ds.DNSBypassCNAME,
			ds.DNSBypassIP,
			ds.DNSBypassIP6,
			ds.DNSBypassTTL,
			ds.DSCP,
			ds.EdgeHeaderRewrite,
			ds.GeoLimitRedirectURL,
			ds.GeoLimit,
			geoLimitCountries,
			ds.GeoProvider,
			ds.GlobalMaxMBPS,
			ds.GlobalMaxTPS,
			ds.FQPacingRate,
			ds.HTTPBypassFQDN,
			ds.InfoURL,
			ds.InitialDispersion,
			ds.IPV6RoutingEnabled,
			ds.LogsEnabled,
			ds.LongDesc,
			longDesc1,
			longDesc2,
			ds.MaxDNSAnswers,
			ds.MaxOriginConnections,
			ds.MidHeaderRewrite,
			ds.MissLat,
			ds.MissLong,
			ds.MultiSiteOrigin,
			ds.OriginShield,
			ds.ProfileID,
			ds.Protocol,
			ds.QStringIgnore,
			ds.RangeRequestHandling,
			ds.RegexRemap,
			ds.Regional,
			ds.RegionalGeoBlocking,
			ds.RemapText,
			ds.RoutingName,
			ds.SigningAlgorithm,
			ds.SSLKeyVersion,
			ds.TenantID,
			ds.Topology,
			ds.TRRequestHeaders,
			ds.TRResponseHeaders,
			ds.TypeID,
			ds.XMLID,
			ds.EcsEnabled,
			ds.RangeSliceBlockSize,
			ds.FirstHeaderRewrite,
			ds.InnerHeaderRewrite,
			ds.LastHeaderRewrite,
			ds.ServiceCategory,
			ds.MaxRequestHeaderBytes,
			pq.Array(ds.RequiredCapabilities),
		)
	}

	if err != nil {
		usrErr, sysErr, code := api.ParseDBError(err)
		return nil, code, usrErr, sysErr
	}
	defer log.Close(resultRows, "inserting Delivery Service")

	id := 0
	if !resultRows.Next() {
		return nil, http.StatusInternalServerError, nil, errors.New("no deliveryservice request inserted, no id was returned")
	}
	var lastUpdated time.Time
	if err := resultRows.Scan(&id, &lastUpdated); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("could not scan id from insert: %w", err)
	}
	if resultRows.Next() {
		return nil, http.StatusInternalServerError, nil, errors.New("too many ids returned from deliveryservice request insert")
	}
	ds.ID = &id

	if ds.ID == nil {
		return nil, http.StatusInternalServerError, nil, errors.New("missing id after insert")
	}

	dsType, err := getTypeFromID(ds.TypeID, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting delivery service type: %w", err)
	}
	ds.Type = (*string)(&dsType)

	if len(ds.TLSVersions) < 1 {
		ds.TLSVersions = nil
	} else if err = recreateTLSVersions(ds.TLSVersions, *ds.ID, tx); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("creating TLS versions for new Delivery Service: %w", err)
	}

	if err := createDefaultRegex(tx, *ds.ID, ds.XMLID); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("creating default regex: %w", err)
	}

	if _, err := createConsistentHashQueryParams(tx, *ds.ID, ds.ConsistentHashQueryParams); err != nil {
		usrErr, sysErr, code := api.ParseDBError(err)
		return nil, code, usrErr, sysErr
	}

	matchlists, err := GetDeliveryServicesMatchLists([]string{ds.XMLID}, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("creating DS: reading matchlists: %w", err)
	}
	matchlist, ok := matchlists[ds.XMLID]
	if !ok {
		return nil, http.StatusInternalServerError, nil, errors.New("creating DS: reading matchlists: not found")
	}
	ds.MatchList = matchlist

	cdnName, cdnDomain, dnssecEnabled, err := getCDNNameDomainDNSSecEnabled(*ds.ID, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("creating DS: getting CDN info: %w", err)
	}

	ds.ExampleURLs = MakeExampleURLs(ds.Protocol, tc.DSType(*ds.Type), ds.RoutingName, ds.MatchList, cdnDomain)

	if err := EnsureParams(tx, *ds.ID, ds.XMLID, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.SigningAlgorithm, dsType, ds.MaxOriginConnections); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("ensuring ds parameters: %w", err)
	}

	if dnssecEnabled && tc.DSType(*ds.Type).UsesDNSSECKeys() {
		if !inf.Config.TrafficVaultEnabled {
			return nil, http.StatusInternalServerError, nil, errors.New("cannot create DNSSEC keys for delivery service: Traffic Vault is not configured")
		}
		if userErr, sysErr, statusCode := PutDNSSecKeys(tx, ds.XMLID, cdnName, ds.ExampleURLs, inf.Vault, r.Context()); userErr != nil || sysErr != nil {
			return nil, statusCode, userErr, sysErr
		}
	}

	user := inf.User
	if err := createPrimaryOrigin(tx, user, ds); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("creating delivery service: %w", err)
	}

	ds.LastUpdated = lastUpdated
	if err := api.CreateChangeLogRawErr(api.ApiChange, "DS: "+ds.XMLID+", ID: "+strconv.Itoa(*ds.ID)+", ACTION: Created delivery service", user, tx); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("error writing to audit log: %w", err)
	}

	if inf.Config.TrafficVaultEnabled && ds.Protocol != nil && (*ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS) {
		err, errCode := GeneratePlaceholderSelfSignedCert(ds, inf, r.Context())
		if err != nil || errCode != http.StatusOK {
			return nil, errCode, nil, fmt.Errorf("creating self signed default cert: %w", err)
		}
	}

	return &ds, http.StatusOK, nil, nil
}

func createDefaultRegex(tx *sql.Tx, dsID int, xmlID string) error {
	regexStr := `.*\.` + xmlID + `\..*`
	regexID := 0
	if err := tx.QueryRow(`INSERT INTO regex (type, pattern) VALUES ((select id from type where name = 'HOST_REGEXP'), $1::text) RETURNING id`, regexStr).Scan(&regexID); err != nil {
		return fmt.Errorf("insert regex: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES ($1::bigint, $2::bigint, 0)`, dsID, regexID); err != nil {
		return fmt.Errorf("executing parameter query to insert location: %w", err)
	}
	return nil
}

func createConsistentHashQueryParams(tx *sql.Tx, dsID int, consistentHashQueryParams []string) (int, error) {
	if len(consistentHashQueryParams) == 0 {
		return 0, nil
	}
	c := 0
	q := `INSERT INTO deliveryservice_consistent_hash_query_param (name, deliveryservice_id) VALUES ($1, $2)`
	for _, k := range consistentHashQueryParams {
		if _, err := tx.Exec(q, k, dsID); err != nil {
			return c, err
		}
		c++
	}

	return c, nil
}

// Read implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Reader
// interface (given that
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.APIInfoer
// is already implemented).
func (ds *TODeliveryService) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	version := ds.APIInfo().Version
	if version == nil {
		return nil, nil, errors.New("TODeliveryService.Read called with nil API version"), http.StatusInternalServerError, nil
	}

	dses, userErr, sysErr, errCode, maxTime := readGetDeliveryServices(h, ds.APIInfo().Params, ds.APIInfo().Tx, ds.APIInfo().User, useIMS, *version)
	if sysErr != nil {
		sysErr = fmt.Errorf("reading dses: %w", sysErr)
		errCode = http.StatusInternalServerError
	}
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	returnable := make([]interface{}, 0, len(dses))
	for _, d := range dses {
		ds := d.DS
		switch {
		case version.Major > 4:
			returnable = append(returnable, ds)
		case version.Major == 4:
			d := ds.Downgrade()
			if version.Minor >= 1 {
				returnable = append(returnable, d.RemoveLD1AndLD2())
			} else {
				returnable = append(returnable, d.DeliveryServiceV40.RemoveLD1AndLD2())
			}
		case version.Major >= 3 && version.Minor >= 1:
			legacyDS := ds.Downgrade()
			legacyDS.LongDesc1 = d.LongDesc1
			legacyDS.LongDesc2 = d.LongDesc2
			returnable = append(returnable, legacyDS.DowngradeToV31())
		case version.Major >= 3:
			legacyDS := ds.Downgrade()
			legacyDS.LongDesc1 = d.LongDesc1
			legacyDS.LongDesc2 = d.LongDesc2
			returnable = append(returnable, legacyDS.DowngradeToV31().DeliveryServiceV30)
		default:
			return nil, nil, fmt.Errorf("TODeliveryService.Read called with invalid API version: %d.%d", version.Major, version.Minor), http.StatusInternalServerError, nil
		}
	}
	return returnable, nil, nil, errCode, maxTime
}

// UpdateV30 is used to handle PUT requests to update a Delivery Service in
// version 3.0 of the Traffic Ops API.
func UpdateV30(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	ds := tc.DeliveryServiceV30{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %w", err), nil)
		return
	}
	ds.ID = &id

	res, status, userErr, sysErr := updateV30(w, r, inf, &ds)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service update was successful", []tc.DeliveryServiceV30{*res})
}

// UpdateV31 is used to handle PUT requests to update a Delivery Service in
// version 3.1 of the Traffic Ops API.
func UpdateV31(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	ds := tc.DeliveryServiceV31{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %w", err), nil)
		return
	}
	ds.ID = &id
	res, status, userErr, sysErr := updateV31(w, r, inf, &ds)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Delivery Service update was successful", []tc.DeliveryServiceV31{*res})
}

// UpdateV40 is used to handle PUT requests to update a Delivery Service in
// version 4.0 of the Traffic Ops API (and isomorphic API versions thereof, with
// respect to Delivery Service representations).
func UpdateV40(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	id := inf.IntParams["id"]

	ds := tc.DeliveryServiceV40{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %w", err), nil)
		return
	}
	ds.ID = &id
	_, cdn, _, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, id)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("deliveryservice update: getting CDN from DS ID %w", err))
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, string(cdn), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	res, status, userErr, sysErr := updateV40(w, r, inf, &ds, true)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service update was successful")

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, []tc.DeliveryServiceV40{*res})
}

// UpdateV41 is used to handle PUT requests to update a Delivery Service in
// version 4.1 of the Traffic Ops API (and isomorphic API versions thereof, with
// respect to Delivery Service representations).
func UpdateV41(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	id := inf.IntParams["id"]

	var ds tc.DeliveryServiceV41
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %w", err), nil)
		return
	}
	ds.ID = &id
	_, cdn, exists, err := dbhelpers.GetDSNameAndCDNFromID(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deliveryservice update: getting CDN from DS ID %w", err))
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no such Delivery Service: #%d", id), nil)
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdn), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}
	res, status, userErr, sysErr := updateV41(w, r, inf, &ds, true)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service update was successful")

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, []tc.DeliveryServiceV41{*res})
}

// UpdateV50 is used to handle PUT requests to update a Delivery Service in
// version 5.0 of the Traffic Ops API (and isomorphic API versions thereof, with
// respect to Delivery Service representations).
func UpdateV50(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	id := inf.IntParams["id"]

	var ds tc.DeliveryServiceV5
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("malformed JSON: %w", err), nil)
		return
	}
	ds.ID = &id
	_, cdn, exists, err := dbhelpers.GetDSNameAndCDNFromID(tx, id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deliveryservice update: getting CDN from DS ID %w", err))
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no such Delivery Service: #%d", id), nil)
		return
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(tx, string(cdn), inf.User.UserName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}
	res, status, userErr, sysErr := updateV50(w, r, inf, &ds, true, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, status, userErr, sysErr)
		return
	}
	alerts := res.TLSVersionsAlerts()
	alerts.AddNewAlert(tc.SuccessLevel, "Delivery Service update was successful")

	api.WriteAlertsObj(w, r, http.StatusOK, alerts, *res)
}

func updateV30(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV30 *tc.DeliveryServiceV30) (*tc.DeliveryServiceV30, int, error, error) {
	dsV31 := tc.DeliveryServiceV31{DeliveryServiceV30: *dsV30}
	// query the DB for existing 3.1 fields in order to "upgrade" this 3.0 request into a 3.1 request
	query := `
SELECT
  ds.max_request_header_bytes
FROM
  deliveryservice ds
WHERE
  ds.id = $1`
	if err := inf.Tx.Tx.QueryRow(query, *dsV30.ID).Scan(
		&dsV31.MaxRequestHeaderBytes,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound, fmt.Errorf("delivery service ID %d not found", *dsV31.ID), nil
		}
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("querying delivery service ID %d: %w", *dsV31.ID, err)
	}
	res, status, userErr, sysErr := updateV31(w, r, inf, &dsV31)
	if res != nil {
		return &res.DeliveryServiceV30, status, userErr, sysErr
	}
	return nil, status, userErr, sysErr
}

func updateV31(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV31 *tc.DeliveryServiceV31) (*tc.DeliveryServiceV31, int, error, error) {
	dsNull := tc.DeliveryServiceNullableV30(*dsV31)
	ds := dsNull.UpgradeToV4()
	dsV41 := ds
	dsMap, err := dbhelpers.GetRequiredCapabilitiesOfDeliveryServices([]int{*dsV31.ID}, inf.Tx.Tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, err
	}
	if caps, ok := dsMap[*dsV31.ID]; ok {
		dsV41.RequiredCapabilities = caps
	}
	if dsV41.ID == nil {
		return nil, http.StatusInternalServerError, nil, errors.New("cannot update a Delivery Service with nil ID")
	}

	tx := inf.Tx.Tx
	var sysErr error
	if dsV41.TLSVersions, sysErr = GetDSTLSVersions(*dsV41.ID, tx); sysErr != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting TLS versions for DS #%d in API version < 4.0: %w", *dsV41.ID, sysErr)
	}

	res, status, usrErr, sysErr := updateV41(w, r, inf, &dsV41, false)
	if res == nil || usrErr != nil || sysErr != nil {
		return nil, status, usrErr, sysErr
	}
	ds.DeliveryServiceV40 = res.DeliveryServiceV40
	if dsV31.CacheURL != nil {
		_, err := tx.Exec("UPDATE deliveryservice SET cacheurl = $1 WHERE id = $2",
			*dsV31.CacheURL,
			*ds.ID,
		)
		if err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			return nil, code, usrErr, sysErr
		}
	}
	if dsV31.LongDesc1 != nil {
		_, err := tx.Exec("UPDATE deliveryservice SET long_desc_1 = $1 WHERE id = $2",
			*dsV31.LongDesc1,
			*ds.ID,
		)
		if err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			return nil, code, usrErr, sysErr
		}
	}
	if dsV31.LongDesc2 != nil {
		_, err := tx.Exec("UPDATE deliveryservice SET long_desc_2 = $1 WHERE id = $2",
			*dsV31.LongDesc2,
			*ds.ID,
		)
		if err != nil {
			usrErr, sysErr, code := api.ParseDBError(err)
			return nil, code, usrErr, sysErr
		}
	}

	if err := EnsureCacheURLParams(tx, *ds.ID, *ds.XMLID, dsV31.CacheURL); err != nil {
		return nil, http.StatusInternalServerError, nil, err
	}

	oldRes := tc.DeliveryServiceV31(ds.DowngradeToV31())
	return &oldRes, http.StatusOK, nil, nil
}

func updateV40(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV40 *tc.DeliveryServiceV40, omitExtraLongDescFields bool) (*tc.DeliveryServiceV40, int, error, error) {
	ds, code, userErr, sysErr := updateV41(w, r, inf, &tc.DeliveryServiceV41{DeliveryServiceV40: *dsV40}, omitExtraLongDescFields)
	if userErr != nil || sysErr != nil || ds == nil {
		return nil, code, userErr, sysErr
	}
	d := ds.DeliveryServiceV40
	return &d, code, nil, nil
}
func updateV41(w http.ResponseWriter, r *http.Request, inf *api.Info, dsV4 *tc.DeliveryServiceV41, omitExtraLongDescFields bool) (*tc.DeliveryServiceV41, int, error, error) {
	upgraded := dsV4.Upgrade()
	res, code, userErr, sysErr := updateV50(w, r, inf, &upgraded, omitExtraLongDescFields, dsV4.LongDesc1, dsV4.LongDesc2)
	if userErr != nil || sysErr != nil || res == nil {
		return nil, code, userErr, sysErr
	}
	ds := res.Downgrade()
	return &ds, code, userErr, sysErr
}

func updateV50(w http.ResponseWriter, r *http.Request, inf *api.Info, ds *tc.DeliveryServiceV5, omitExtraLongDescFields bool, longDesc1, longDesc2 *string) (*tc.DeliveryServiceV5, int, error, error) {
	tx := inf.Tx.Tx
	user := inf.User
	userErr, sysErr := Validate(tx, ds)
	if userErr != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid request: %w", userErr), nil
	}
	if sysErr != nil {
		return nil, http.StatusInternalServerError, nil, sysErr
	}

	if authorized, err := isTenantAuthorized(inf, ds); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("checking tenant: %w", err)
	} else if !authorized {
		return nil, http.StatusForbidden, errors.New("not authorized on this tenant"), nil
	}

	if ds.ID == nil {
		return nil, http.StatusBadRequest, errors.New("missing id"), nil
	}

	dsType, ok, err := getDSType(tx, ds.XMLID)
	if !ok {
		return nil, http.StatusNotFound, errors.New("delivery service '" + ds.XMLID + "' not found"), nil
	}
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting delivery service type during update: %w", err)
	}

	var errCode int
	var oldDetails TODeliveryServiceOldDetails
	if dsType.HasSSLKeys() {
		oldDetails, userErr, sysErr, errCode = getOldDetails(*ds.ID, tx)
		if userErr != nil || sysErr != nil {
			return nil, errCode, userErr, sysErr
		}
		sslKeysExist, err := getSSLVersion(ds.XMLID, tx)
		if err != nil {
			return nil, http.StatusInternalServerError, nil, fmt.Errorf("querying delivery service with sslKeyVersion failed: %w", err)
		}
		if sslKeysExist {
			errStr := "delivery service has ssl keys that cannot be automatically changed, therefore %s is immutable"
			if oldDetails.OldCDNID != ds.CDNID {
				return nil, http.StatusBadRequest, fmt.Errorf(errStr, "CDN ID"), nil
			}
			if ds.CDNName != nil && oldDetails.OldCDNName != *ds.CDNName {
				return nil, http.StatusBadRequest, fmt.Errorf(errStr, "CDN Name"), nil
			}
			if oldDetails.OldRoutingName != ds.RoutingName {
				return nil, http.StatusBadRequest, fmt.Errorf(errStr, "Routing Name"), nil
			}
			ds.SSLKeyVersion = oldDetails.OldSSLKeyVersion
		}
	}

	// TODO change DeepCachingType to implement sql.Valuer and sql.Scanner, so sqlx struct scan can be used.
	deepCachingType := tc.DeepCachingType("").String()
	if ds.DeepCachingType != "" {
		deepCachingType = ds.DeepCachingType.String() // necessary, because DeepCachingType's default needs to insert the string, not "", and Query doesn't call .String().
	}

	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, *ds.ID, "deliveryservice")
	if userErr != nil || sysErr != nil {
		return nil, errCode, userErr, sysErr
	}

	if errCode, userErr, sysErr = dbhelpers.CheckTopology(inf.Tx, *ds); userErr != nil || sysErr != nil {
		return nil, errCode, userErr, sysErr
	}

	if ds.Topology != nil {
		if len(ds.RequiredCapabilities) > 0 {
			if userErr, sysErr, status := EnsureTopologyBasedRequiredCapabilities(tx, *ds.ID, *ds.Topology, ds.RequiredCapabilities); userErr != nil || sysErr != nil {
				return nil, status, userErr, sysErr
			}
		}

		userErr, sysErr, status := dbhelpers.CheckOriginServerInDSCG(tx, *ds.ID, *ds.Topology)
		if userErr != nil || sysErr != nil {
			return nil, status, userErr, sysErr
		}
	}

	var geoLimitCountries string
	if ds.GeoLimitCountries != nil {
		geoLimitCountries = strings.Join(ds.GeoLimitCountries, ",")
	}
	var resultRows *sql.Rows
	if omitExtraLongDescFields {
		if longDesc1 != nil || longDesc2 != nil {
			return nil, http.StatusBadRequest, errors.New("the longDesc1 and longDesc2 fields are no longer supported in API 4.0 onwards"), nil
		}
		resultRows, err = tx.Query(updateDSQueryWithoutLD1AndLD2(),
			ds.Active,
			ds.CCRDNSTTL,
			ds.CDNID,
			ds.CheckPath,
			deepCachingType,
			ds.DisplayName,
			ds.DNSBypassCNAME,
			ds.DNSBypassIP,
			ds.DNSBypassIP6,
			ds.DNSBypassTTL,
			ds.DSCP,
			ds.EdgeHeaderRewrite,
			ds.GeoLimitRedirectURL,
			ds.GeoLimit,
			geoLimitCountries,
			ds.GeoProvider,
			ds.GlobalMaxMBPS,
			ds.GlobalMaxTPS,
			ds.FQPacingRate,
			ds.HTTPBypassFQDN,
			ds.InfoURL,
			ds.InitialDispersion,
			ds.IPV6RoutingEnabled,
			ds.LogsEnabled,
			ds.LongDesc,
			ds.MaxDNSAnswers,
			ds.MidHeaderRewrite,
			ds.MissLat,
			ds.MissLong,
			ds.MultiSiteOrigin,
			ds.OriginShield,
			ds.ProfileID,
			ds.Protocol,
			ds.QStringIgnore,
			ds.RangeRequestHandling,
			ds.RegexRemap,
			ds.Regional,
			ds.RegionalGeoBlocking,
			ds.RemapText,
			ds.RoutingName,
			ds.SigningAlgorithm,
			ds.SSLKeyVersion,
			ds.TenantID,
			ds.TRRequestHeaders,
			ds.TRResponseHeaders,
			ds.TypeID,
			ds.XMLID,
			ds.AnonymousBlockingEnabled,
			ds.ConsistentHashRegex,
			ds.MaxOriginConnections,
			ds.EcsEnabled,
			ds.RangeSliceBlockSize,
			ds.Topology,
			ds.FirstHeaderRewrite,
			ds.InnerHeaderRewrite,
			ds.LastHeaderRewrite,
			ds.ServiceCategory,
			ds.MaxRequestHeaderBytes,
			pq.Array(ds.RequiredCapabilities),
			ds.ID,
		)
	} else {
		resultRows, err = tx.Query(updateDSQuery(),
			ds.Active,
			ds.CCRDNSTTL,
			ds.CDNID,
			ds.CheckPath,
			deepCachingType,
			ds.DisplayName,
			ds.DNSBypassCNAME,
			ds.DNSBypassIP,
			ds.DNSBypassIP6,
			ds.DNSBypassTTL,
			ds.DSCP,
			ds.EdgeHeaderRewrite,
			ds.GeoLimitRedirectURL,
			ds.GeoLimit,
			geoLimitCountries,
			ds.GeoProvider,
			ds.GlobalMaxMBPS,
			ds.GlobalMaxTPS,
			ds.FQPacingRate,
			ds.HTTPBypassFQDN,
			ds.InfoURL,
			ds.InitialDispersion,
			ds.IPV6RoutingEnabled,
			ds.LogsEnabled,
			ds.LongDesc,
			longDesc1,
			longDesc2,
			ds.MaxDNSAnswers,
			ds.MidHeaderRewrite,
			ds.MissLat,
			ds.MissLong,
			ds.MultiSiteOrigin,
			ds.OriginShield,
			ds.ProfileID,
			ds.Protocol,
			ds.QStringIgnore,
			ds.RangeRequestHandling,
			ds.RegexRemap,
			ds.Regional,
			ds.RegionalGeoBlocking,
			ds.RemapText,
			ds.RoutingName,
			ds.SigningAlgorithm,
			ds.SSLKeyVersion,
			ds.TenantID,
			ds.TRRequestHeaders,
			ds.TRResponseHeaders,
			ds.TypeID,
			ds.XMLID,
			ds.AnonymousBlockingEnabled,
			ds.ConsistentHashRegex,
			ds.MaxOriginConnections,
			ds.EcsEnabled,
			ds.RangeSliceBlockSize,
			ds.Topology,
			ds.FirstHeaderRewrite,
			ds.InnerHeaderRewrite,
			ds.LastHeaderRewrite,
			ds.ServiceCategory,
			ds.MaxRequestHeaderBytes,
			pq.Array(ds.RequiredCapabilities),
			ds.ID,
		)
	}

	if err != nil {
		usrErr, sysErr, code := api.ParseDBError(err)
		return nil, code, usrErr, sysErr
	}
	defer log.Close(resultRows, "updating Delivery Service")
	if !resultRows.Next() {
		return nil, http.StatusNotFound, errors.New("no delivery service found with this id"), nil
	}
	var lastUpdated time.Time
	if err := resultRows.Scan(&lastUpdated); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("scan updating delivery service: %w", err)
	}
	if resultRows.Next() {
		return nil, http.StatusInternalServerError, nil, errors.New("updating delivery service " + ds.XMLID + ": this update affected too many rows: > 1")
	}
	if ds.ID == nil {
		return nil, http.StatusInternalServerError, nil, errors.New("missing id after update")
	}

	if len(ds.TLSVersions) < 1 {
		ds.TLSVersions = nil
	}
	err = recreateTLSVersions(ds.TLSVersions, *ds.ID, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("updating TLS versions for DS #%d: %w", *ds.ID, err)
	}

	newDSType, err := getTypeFromID(ds.TypeID, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting delivery service type after update: %w", err)
	}
	ds.Type = (*string)(&newDSType)

	cdnDomain, err := getCDNDomain(*ds.ID, tx) // need to get the domain again, in case it changed.
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting CDN domain: (%s) after update: %w", cdnDomain, err)
	}

	matchLists, err := GetDeliveryServicesMatchLists([]string{ds.XMLID}, tx)
	if err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("getting matchlists after update: %w", err)
	}
	ml, ok := matchLists[ds.XMLID]
	if !ok {
		return nil, http.StatusInternalServerError, nil, errors.New("no matchlists after update")
	}
	ds.MatchList = ml

	if err := EnsureParams(tx, *ds.ID, ds.XMLID, ds.EdgeHeaderRewrite, ds.MidHeaderRewrite, ds.RegexRemap, ds.SigningAlgorithm, newDSType, ds.MaxOriginConnections); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("ensuring ds parameters: %w", err)
	}

	if oldDetails.OldOrgServerFQDN != nil && ds.OrgServerFQDN != nil && *oldDetails.OldOrgServerFQDN != *ds.OrgServerFQDN {
		if err := updatePrimaryOrigin(tx, user, *ds); err != nil {
			return nil, http.StatusInternalServerError, nil, fmt.Errorf("updating delivery service: %w", err)
		}
	}

	ds.LastUpdated = lastUpdated

	// the update may change or delete the query params -- delete existing and re-add if any provided
	q := `DELETE FROM deliveryservice_consistent_hash_query_param WHERE deliveryservice_id = $1`
	if res, err := tx.Exec(q, *ds.ID); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("deleting consistent hash query params for ds %s: %w", ds.XMLID, err)
	} else if c, _ := res.RowsAffected(); c > 0 {
		api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.XMLID+", ID: "+strconv.Itoa(*ds.ID)+", ACTION: Deleted "+strconv.FormatInt(c, 10)+" consistent hash query params", user, tx)
	}

	if _, err = createConsistentHashQueryParams(tx, *ds.ID, ds.ConsistentHashQueryParams); err != nil {
		usrErr, sysErr, code := api.ParseDBError(err)
		return nil, code, usrErr, sysErr
	}

	if err := api.CreateChangeLogRawErr(api.ApiChange, "Updated ds: "+ds.XMLID+" id: "+strconv.Itoa(*ds.ID), user, tx); err != nil {
		return nil, http.StatusInternalServerError, nil, fmt.Errorf("writing change log entry: %w", err)
	}

	if inf.Config.TrafficVaultEnabled && ds.Protocol != nil && (*ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPAndHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS) {
		err, errCode := GeneratePlaceholderSelfSignedCert(*ds, inf, r.Context())
		if err != nil || errCode != http.StatusOK {
			return nil, errCode, nil, fmt.Errorf("creating self signed default cert: %w", err)
		}
	}

	return ds, http.StatusOK, nil, nil
}

// Delete implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Deleter
// interface (given that the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.APIInfoer
// and
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.Identifier
// interfaces are already implemented).
func (ds *TODeliveryService) Delete() (error, error, int) {
	if ds.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}

	xmlID, ok, err := GetXMLID(ds.ReqInfo.Tx.Tx, *ds.ID)
	if err != nil {
		return nil, fmt.Errorf("ds delete: getting xmlid: %w", err), http.StatusInternalServerError
	} else if !ok {
		return errors.New("delivery service not found"), nil, http.StatusNotFound
	}
	ds.XMLID = xmlID

	var userErr error
	var sysErr error
	var errCode int
	if ds.CDNID != 0 {
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(ds.APIInfo().Tx.Tx, int64(ds.CDNID), ds.APIInfo().User.UserName)
	} else if ds.CDNName != nil {
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(ds.APIInfo().Tx.Tx, *ds.CDNName, ds.APIInfo().User.UserName)
	} else {
		_, cdnName, _, err := dbhelpers.GetDSNameAndCDNFromID(ds.ReqInfo.Tx.Tx, *ds.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get cdn name for DS: %w", err), http.StatusBadRequest
		}
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(ds.APIInfo().Tx.Tx, string(cdnName), ds.APIInfo().User.UserName)
	}
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	// Note ds regexes MUST be deleted before the ds, because there's a ON DELETE CASCADE on deliveryservice_regex (but not on regex).
	// Likewise, it MUST happen in a transaction with the later DS delete, so they aren't deleted if the DS delete fails.
	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM regex WHERE id IN (SELECT regex FROM deliveryservice_regex WHERE deliveryservice=$1)`, *ds.ID); err != nil {
		return nil, fmt.Errorf("TODeliveryService.Delete deleting regexes for delivery service: %w", err), http.StatusInternalServerError
	}

	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM deliveryservice_regex WHERE deliveryservice=$1`, *ds.ID); err != nil {
		return nil, fmt.Errorf("TODeliveryService.Delete deleting delivery service regexes: %w", err), http.StatusInternalServerError
	}

	userErr, sysErr, errCode = api.GenericDelete(ds)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	paramConfigFilePrefixes := []string{"hdr_rw_", "hdr_rw_mid_", "regex_remap_", "cacheurl_"}
	configFiles := []string{}
	for _, prefix := range paramConfigFilePrefixes {
		configFiles = append(configFiles, prefix+ds.XMLID+".config")
	}

	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`, pq.Array(configFiles)); err != nil {
		return nil, fmt.Errorf("TODeliveryService.Delete deleting delivery service parameters: %w", err), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

// DeleteQuery implements part of the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.GenericDeleter
// interface.
func (*TODeliveryService) DeleteQuery() string {
	return `DELETE FROM deliveryservice WHERE id = :id`
}

// addActive adds the query parameter for a DS's "active" state for the
// appropriate given API version to the provided 'WHERE' clause, and updates
// queryValues with the client-provided value from params as necessary.
func addActive(where string, params map[string]string, queryValues map[string]interface{}, version api.Version) (string, error) {
	active, ok := params["active"]
	if !ok {
		return where, nil
	}
	if version.Major < 5 {
		b, err := strconv.ParseBool(active)
		if err != nil {
			return where, fmt.Errorf("active cannot parse to boolean: %w", err)
		}
		if b {
			return dbhelpers.AppendWhere(where, "ds.active = '"+string(tc.DSActiveStateActive)+"'"), nil
		}
		return dbhelpers.AppendWhere(where, "ds.active <> '"+string(tc.DSActiveStateActive)+"'"), nil
	}
	switch tc.DeliveryServiceActiveState(active) {
	case tc.DSActiveStateActive:
		fallthrough
	case tc.DSActiveStateInactive:
		fallthrough
	case tc.DSActiveStatePrimed:
		queryValues["active"] = active
		return dbhelpers.AppendWhere(where, "ds.active=:active"), nil
	}
	return where, fmt.Errorf("active illegal value '%s' must be one of '%s', '%s', or '%s'", active, tc.DSActiveStateActive, tc.DSActiveStateInactive, tc.DSActiveStatePrimed)
}

func readGetDeliveryServices(h http.Header, params map[string]string, tx *sqlx.Tx, user *auth.CurrentUser, useIMS bool, version api.Version) ([]DSWithLegacyFields, error, error, int, *time.Time) {
	if tx == nil {
		return nil, nil, errors.New("nil transaction passed to readGetDeliveryServices"), http.StatusInternalServerError, nil
	}
	if user == nil {
		return nil, nil, errors.New("nil user passed to readGetDeliveryServices"), http.StatusInternalServerError, nil
	}
	if params == nil {
		params = make(map[string]string)
	}

	var maxTime time.Time
	var runSecond bool
	// TODO: is this necessary?
	if idParam, ok := params["id"]; ok {
		params["id"] = strings.TrimSuffix(idParam, ".json")
	}
	if _, ok := params["orderby"]; !ok {
		params["orderby"] = "xml_id"
	}

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":               {Column: "ds.id", Checker: api.IsInt},
		"cdn":              {Column: "ds.cdn_id", Checker: api.IsInt},
		"xml_id":           {Column: "ds.xml_id"},
		"xmlId":            {Column: "ds.xml_id"},
		"profile":          {Column: "ds.profile", Checker: api.IsInt},
		"type":             {Column: "ds.type", Checker: api.IsInt},
		"logsEnabled":      {Column: "ds.logs_enabled", Checker: api.IsBool},
		"signingAlgorithm": {Column: "ds.signing_algorithm"},
		"serviceCategory":  {Column: "ds.service_category"},
		"tenant":           {Column: "ds.tenant_id", Checker: api.IsInt},
		"topology":         {Column: "ds.topology"},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	where, err := addActive(where, params, queryValues, version)
	if err != nil {
		return nil, err, nil, http.StatusBadRequest, nil
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return nil, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	tenantIDs, err := tenant.GetUserTenantIDListTx(tx.Tx, user.TenantID)

	if err != nil {
		err = fmt.Errorf("received error querying for user's tenants: %w", err)
		return nil, nil, err, http.StatusInternalServerError, &maxTime
	}

	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "ds.tenant_id", tenantIDs)

	if accessibleTo, ok := params["accessibleTo"]; ok {
		if err := api.IsInt(accessibleTo); err != nil {
			log.Errorln("unknown parameter value: " + err.Error())
			return nil, errors.New("accessibleTo must be an integer"), nil, http.StatusBadRequest, &maxTime
		}
		accessibleTo, _ := strconv.Atoi(accessibleTo)
		accessibleTenants, err := tenant.GetUserTenantIDListTx(tx.Tx, accessibleTo)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to get tenants: %w", err), http.StatusInternalServerError, &maxTime
		}
		where += " AND ds.tenant_id = ANY(CAST(:accessibleTo AS bigint[])) "
		queryValues["accessibleTo"] = pq.Array(accessibleTenants)
	}

	if reqCap, ok := params["requiredCapability"]; ok {
		where += " AND '" + reqCap + "'=ANY(ds.required_capabilities)"
	}

	query := SelectDeliveryServicesQuery + where + orderBy + pagination
	log.Debugln("generated deliveryServices query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	r, e1, e2, code := GetDeliveryServices(query, queryValues, tx)
	return r, e1, e2, code, &maxTime
}

func requiredIfMatchesTypeName(patterns []string, typeName string) func(interface{}) error {
	return func(value interface{}) error {
		switch v := value.(type) {
		case *int:
			if v != nil {
				return nil
			}
		case *bool:
			if v != nil {
				return nil
			}
		case bool:
			return nil
		case *string:
			if v != nil {
				if *v != "" {
					return nil
				}
			}
		case *float64:
			if v != nil {
				return nil
			}
		default:
			return fmt.Errorf("validation failure: unknown type %T", value)
		}
		pattern := strings.Join(patterns, "|")
		err := error(nil)
		match := false
		if typeName != "" {
			match, err = regexp.MatchString(pattern, typeName)
			if match {
				return fmt.Errorf("is required if type is '%s'", typeName)
			}
		}
		return err
	}
}

var validTLSVersionPattern = regexp.MustCompile(`^\d+\.\d+$`)

// Validate checks that the given Delivery Service is valid according to various
// criteria. *It also modifies the given DS under certain circumstances,
// providing default values to some properties when they are zero-valued or nil
// references*. This will panic if either argument is nil. The error returned is
// safe for clients to see.
func Validate(tx *sql.Tx, ds *tc.DeliveryServiceV5) (error, error) {
	sanitize(ds)
	neverOrAlways := validation.NewStringRule(tovalidate.IsOneOfStringICase("NEVER", "ALWAYS"),
		"must be one of 'NEVER' or 'ALWAYS'")
	isDNSName := validation.NewStringRule(govalidator.IsDNSName, "must be a valid hostname")
	noPeriods := validation.NewStringRule(tovalidate.NoPeriods, "cannot contain periods")
	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")
	noLineBreaks := validation.NewStringRule(tovalidate.NoLineBreaks, "cannot contain line breaks")
	errs := tovalidate.ToErrors(validation.Errors{
		"cdnId":               validation.Validate(ds.CDNID, validation.Required),
		"deepCachingType":     validation.Validate(ds.DeepCachingType, neverOrAlways),
		"displayName":         validation.Validate(ds.DisplayName, validation.Required, validation.Length(1, 48)),
		"dscp":                validation.Validate(ds.DSCP, validation.NotNil, validation.Min(0)),
		"geoLimit":            validation.Validate(ds.GeoLimit, validation.NotNil),
		"geoProvider":         validation.Validate(ds.GeoProvider, validation.NotNil),
		"httpByPassFqdn":      validation.Validate(ds.HTTPBypassFQDN, isDNSName),
		"logsEnabled":         validation.Validate(ds.LogsEnabled, validation.NotNil),
		"regional":            validation.Validate(ds.Regional, validation.NotNil),
		"regionalGeoBlocking": validation.Validate(ds.RegionalGeoBlocking, validation.NotNil),
		"remapText":           validation.Validate(ds.RemapText, noLineBreaks),
		"routingName":         validation.Validate(ds.RoutingName, isDNSName, noPeriods, validation.Length(1, 48)),
		"typeId":              validation.Validate(ds.TypeID, validation.Required, validation.Min(1)),
		"xmlId":               validation.Validate(ds.XMLID, validation.Required, noSpaces, noPeriods, validation.Length(1, 48)),
		"tlsVersions": validation.Validate(ds.TLSVersions, validation.By(
			func(value interface{}) error {
				vers, ok := value.([]string)
				if !ok {
					return fmt.Errorf("must be an array of string, got: %T", value)
				}
				seen := make(map[string]struct{}, len(vers))
				for _, tlsVersion := range vers {
					if _, ok := seen[tlsVersion]; ok {
						return fmt.Errorf("duplicate version '%s'", tlsVersion)
					}
					seen[tlsVersion] = struct{}{}
					if !validTLSVersionPattern.Match([]byte(tlsVersion)) {
						return fmt.Errorf("invalid TLS version '%s'", tlsVersion)
					}
				}
				return nil
			},
		)),
	})
	if err := validateGeoLimitCountries(ds); err != nil {
		errs = append(errs, err)
	}
	if err := validateTopologyFields(ds); err != nil {
		errs = append(errs, err)
	}
	if err := validateTypeFields(tx, ds); err != nil {
		errs = append(errs, fmt.Errorf("type fields: %w", err))
	}
	userErr, sysErr := validateRequiredCapabilities(tx, ds)
	if sysErr != nil {
		return nil, fmt.Errorf("reading/ scanning required capabilities: %w", sysErr)
	}
	if userErr != nil {
		errs = append(errs, errors.New("required capabilities: "+userErr.Error()))
	}
	if len(errs) == 0 {
		return nil, nil
	}
	return util.JoinErrs(errs), nil
}

func validateRequiredCapabilities(tx *sql.Tx, ds *tc.DeliveryServiceV5) (error, error) {
	missing := make([]string, 0)
	var missingCap string
	query := `SELECT missing
FROM (
    SELECT UNNEST($1::TEXT[])
    EXCEPT
    SELECT UNNEST(ARRAY_AGG(name)) FROM server_capability) t(missing)
`
	if len(ds.RequiredCapabilities) > 0 {
		rows, err := tx.Query(query, pq.Array(ds.RequiredCapabilities))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&missingCap)
			if err != nil {
				return nil, err
			}
			missing = append(missing, missingCap)
		}
		if len(missing) > 0 {
			msg := strings.Join(missing, ",")
			userErr := fmt.Errorf("the following capabilities do not exist: %s", msg)
			return userErr, nil
		}
	}
	return nil, nil
}

func validateGeoLimitCountries(ds *tc.DeliveryServiceV5) error {
	var IsLetter = regexp.MustCompile(`^[A-Z]+$`).MatchString
	if len(ds.GeoLimitCountries) == 0 {
		return nil
	}
	for _, cc := range ds.GeoLimitCountries {
		if cc != "" && !IsLetter(cc) {
			return errors.New("country codes can only contain alphabetical characters")
		}
	}
	return nil
}

func validateTopologyFields(ds *tc.DeliveryServiceV5) error {
	if ds.Topology != nil && (ds.EdgeHeaderRewrite != nil || ds.MidHeaderRewrite != nil) {
		return errors.New("cannot set edgeHeaderRewrite or midHeaderRewrite while a Topology is assigned. Use firstHeaderRewrite, innerHeaderRewrite, and/or lastHeaderRewrite instead")
	}
	if ds.Topology == nil && (ds.FirstHeaderRewrite != nil || ds.InnerHeaderRewrite != nil || ds.LastHeaderRewrite != nil) {
		return errors.New("cannot set firstHeaderRewrite, innerHeaderRewrite, or lastHeaderRewrite unless this delivery service is assigned to a Topology. Use edgeHeaderRewrite and/or midHeaderRewrite instead")
	}
	return nil
}

func parseOrgServerFQDN(orgServerFQDN string) (*string, *string, *string, error) {
	originRegex := regexp.MustCompile(`^(https?)://([^:]+)(:(\d+))?$`)
	matches := originRegex.FindStringSubmatch(orgServerFQDN)
	if len(matches) == 0 {
		return nil, nil, nil, fmt.Errorf("unable to parse invalid orgServerFqdn: '%s'", orgServerFQDN)
	}

	protocol := strings.ToLower(matches[1])
	FQDN := matches[2]

	if len(protocol) == 0 || len(FQDN) == 0 {
		return nil, nil, nil, fmt.Errorf("empty Origin protocol or FQDN parsed from '%s'", orgServerFQDN)
	}

	var port *string
	if len(matches[4]) != 0 {
		port = &matches[4]
	}
	return &protocol, &FQDN, port, nil
}

func validateOrgServerFQDN(orgServerFQDN string) bool {
	_, fqdn, port, err := parseOrgServerFQDN(orgServerFQDN)
	if err != nil || !govalidator.IsHost(*fqdn) || (port != nil && !govalidator.IsPort(*port)) {
		return false
	}
	return true
}

const (
	dnsTypeRegexp      = "^DNS.*$"
	httpTypeRegexp     = "^HTTP.*$"
	steeringTypeRegexp = "^STEERING.*$"
)

const (
	latitudeErr  = "Must be a floating point number within the range +-90"
	longitudeErr = "Must be a floating point number within the range +-180"
)

// validateTypeFields validates the TypeName-related field.
func validateTypeFields(tx *sql.Tx, ds *tc.DeliveryServiceV5) error {
	typeName, err := tc.ValidateTypeID(tx, &ds.TypeID, "deliveryservice")
	if err != nil {
		return err
	}

	errs := validation.Errors{
		"consistentHashQueryParams": validation.Validate(ds,
			validation.By(func(dsi interface{}) error {
				ds := dsi.(*tc.DeliveryServiceV5)
				if len(ds.ConsistentHashQueryParams) == 0 {
					return nil
				}
				if !tc.DSType(typeName).IsHTTP() {
					return fmt.Errorf("consistentHashQueryParams not allowed for '%s' deliveryservice type", typeName)
				}

				for _, param := range ds.ConsistentHashQueryParams {
					if param == tc.ReservedConsistentHashingQueryParameterFormat || param == tc.ReservedConsistentHashingQueryParameterTRRED || param == tc.ReservedConsistentHashingQueryParameterFakeClientIPAddress {
						return fmt.Errorf("'%s' cannot be used in consistent hashing, because it is reserved for use by Traffic Router", param)
					}
				}

				return nil
			})),
		"initialDispersion": validation.Validate(ds.InitialDispersion,
			validation.By(requiredIfMatchesTypeName([]string{httpTypeRegexp}, typeName)),
			validation.By(tovalidate.IsGreaterThanZero)),
		"ipv6RoutingEnabled": validation.Validate(ds.IPV6RoutingEnabled,
			validation.By(requiredIfMatchesTypeName([]string{steeringTypeRegexp, dnsTypeRegexp, httpTypeRegexp}, typeName))),
		"missLat": validation.Validate(ds.MissLat,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName)),
			validation.Min(-90.0).Error(latitudeErr),
			validation.Max(90.0).Error(latitudeErr)),
		"missLong": validation.Validate(ds.MissLong,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName)),
			validation.Min(-180.0).Error(longitudeErr),
			validation.Max(180.0).Error(longitudeErr)),
		"multiSiteOrigin": validation.Validate(ds.MultiSiteOrigin,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName))),
		"orgServerFqdn": validation.Validate(ds.OrgServerFQDN,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName)),
			validation.NewStringRule(validateOrgServerFQDN, "must start with http:// or https:// and be followed by a valid hostname with an optional port (no trailing slash)")),
		"rangeSliceBlockSize": validation.Validate(ds,
			validation.By(func(dsi interface{}) error {
				ds := dsi.(*tc.DeliveryServiceV5)
				if ds.RangeRequestHandling != nil {
					if *ds.RangeRequestHandling == 3 {
						return validation.Validate(ds.RangeSliceBlockSize, validation.Required,
							// Per Slice Plugin implementation
							validation.Min(tc.MinRangeSliceBlockSize), // 256KiB
							validation.Max(tc.MaxRangeSliceBlockSize), // 32MiB
						)
					}
					if ds.RangeSliceBlockSize != nil {
						return errors.New("rangeSliceBlockSize can only be set if the rangeRequestHandling is set to 3 (Use the Slice Plugin)")
					}
				}
				return nil
			})),
		"protocol": validation.Validate(ds.Protocol,
			validation.By(requiredIfMatchesTypeName([]string{steeringTypeRegexp, dnsTypeRegexp, httpTypeRegexp}, typeName))),
		"qstringIgnore": validation.Validate(ds.QStringIgnore,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName))),
		"rangeRequestHandling": validation.Validate(ds.RangeRequestHandling,
			validation.By(requiredIfMatchesTypeName([]string{dnsTypeRegexp, httpTypeRegexp}, typeName))),
		"tlsVersions": validation.Validate(
			&ds.TLSVersions,
			validation.By(
				func(value interface{}) error {
					vers := value.(*[]string)
					if vers != nil && len(*vers) > 0 && tc.DSType(typeName).IsSteering() {
						return errors.New("must be 'null' for STEERING-Type and CLIENT_STEERING-Type Delivery Services")
					}
					return nil
				},
			),
		),
		"topology": validation.Validate(ds,
			validation.By(func(dsi interface{}) error {
				ds := dsi.(*tc.DeliveryServiceV5)
				if ds.Topology != nil && tc.DSType(typeName).IsSteering() {
					return fmt.Errorf("steering deliveryservice types cannot be assigned to a topology")
				}
				return nil
			})),
		"maxRequestHeaderBytes": validation.Validate(ds,
			validation.By(func(dsi interface{}) error {
				ds := dsi.(*tc.DeliveryServiceV5)
				if ds.MaxRequestHeaderBytes == nil {
					return errors.New("maxRequestHeaderBytes empty, must be a valid positive value")
				}
				if *ds.MaxRequestHeaderBytes < 0 || *ds.MaxRequestHeaderBytes > 2147483647 {
					return errors.New("maxRequestHeaderBytes must be a valid non negative value between 0 and 2147483647")
				}
				return nil
			})),
	}
	toErrs := tovalidate.ToErrors(errs)
	if len(toErrs) > 0 {
		return errors.New(util.JoinErrsStr(toErrs))
	}
	return nil
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(ds.last_updated) as t from deliveryservice as ds
	JOIN type ON ds.type = type.id
	JOIN cdn ON ds.cdn_id = cdn.id
	LEFT JOIN profile ON ds.profile = profile.id
	LEFT JOIN tenant ON ds.tenant_id = tenant.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='deliveryservice') as res`
}

func getOldDetails(id int, tx *sql.Tx) (TODeliveryServiceOldDetails, error, error, int) {
	q := `
SELECT ds.routing_name, ds.ssl_key_version, cdn.name, cdn.id,
(SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
FROM origin o
WHERE o.deliveryservice = ds.id
AND o.is_primary) as org_server_fqdn
FROM  deliveryservice as ds
JOIN cdn ON ds.cdn_id = cdn.id
WHERE ds.id=$1
`
	var oldDetails TODeliveryServiceOldDetails
	if err := tx.QueryRow(q, id).Scan(&oldDetails.OldRoutingName, &oldDetails.OldSSLKeyVersion, &oldDetails.OldCDNName, &oldDetails.OldCDNID, &oldDetails.OldOrgServerFQDN); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return oldDetails, fmt.Errorf("querying delivery service %v host name: no such delivery service exists", id), nil, http.StatusNotFound
		}
		return oldDetails, nil, fmt.Errorf("querying delivery service %v host name: "+err.Error(), id), http.StatusInternalServerError
	}
	return oldDetails, nil, nil, http.StatusOK
}

func getTypeFromID(id int, tx *sql.Tx) (tc.DSType, error) {
	// TODO combine with getOldDetails, to only make one query?
	name := ""
	if err := tx.QueryRow(`SELECT name FROM type WHERE id = $1`, id).Scan(&name); err != nil {
		return "", fmt.Errorf("querying type ID %v: "+err.Error()+"\n", id)
	}
	return tc.DSTypeFromString(name), nil
}

func updatePrimaryOrigin(tx *sql.Tx, user *auth.CurrentUser, ds tc.DeliveryServiceV5) error {
	count := 0
	q := `SELECT count(*) FROM origin WHERE deliveryservice = $1 AND is_primary`
	if err := tx.QueryRow(q, *ds.ID).Scan(&count); err != nil {
		return fmt.Errorf("querying existing primary origin for ds %s: %w", ds.XMLID, err)
	}

	if ds.OrgServerFQDN == nil || *ds.OrgServerFQDN == "" {
		if count == 1 {
			// the update is removing the existing orgServerFQDN, so the existing row needs to be deleted
			q = `DELETE FROM origin WHERE deliveryservice = $1 AND is_primary`
			if _, err := tx.Exec(q, *ds.ID); err != nil {
				return fmt.Errorf("deleting primary origin for ds %s: %w", ds.XMLID, err)
			}
			api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.XMLID+", ID: "+strconv.Itoa(*ds.ID)+", ACTION: Deleted primary origin", user, tx)
		}
		return nil
	}

	if count == 0 {
		// orgServerFQDN is going from null to not null, so the primary origin needs to be created
		return createPrimaryOrigin(tx, user, ds)
	}

	protocol, fqdn, port, err := parseOrgServerFQDN(*ds.OrgServerFQDN)
	if err != nil {
		return fmt.Errorf("updating primary origin: %w", err)
	}

	name := ""
	q = `UPDATE origin SET protocol = $1, fqdn = $2, port = $3 WHERE is_primary AND deliveryservice = $4 RETURNING name`
	if err := tx.QueryRow(q, protocol, fqdn, port, *ds.ID).Scan(&name); err != nil {
		return fmt.Errorf("update primary origin for ds %s from '%s': %w", ds.XMLID, *ds.OrgServerFQDN, err)
	}

	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.XMLID+", ID: "+strconv.Itoa(*ds.ID)+", ACTION: Updated primary origin: "+name, user, tx)

	return nil
}

func createPrimaryOrigin(tx *sql.Tx, user *auth.CurrentUser, ds tc.DeliveryServiceV5) error {
	if ds.OrgServerFQDN == nil {
		return nil
	}

	protocol, fqdn, port, err := parseOrgServerFQDN(*ds.OrgServerFQDN)
	if err != nil {
		return fmt.Errorf("creating primary origin: %w", err)
	}

	originID := 0
	q := `INSERT INTO origin (name, fqdn, protocol, is_primary, port, deliveryservice, tenant) VALUES ($1, $2, $3, TRUE, $4, $5, $6) RETURNING id`
	if err := tx.QueryRow(q, ds.XMLID, fqdn, protocol, port, ds.ID, ds.TenantID).Scan(&originID); err != nil {
		return fmt.Errorf("insert origin from '%s': %w", *ds.OrgServerFQDN, err)
	}

	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+ds.XMLID+", ID: "+strconv.Itoa(*ds.ID)+", ACTION: Created primary origin id: "+strconv.Itoa(originID), user, tx)

	return nil
}

func getDSType(tx *sql.Tx, xmlid string) (tc.DSType, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM type WHERE id = (select type from deliveryservice where xml_id = $1)`, xmlid).Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("querying deliveryservice type name: " + err.Error())
	}
	return tc.DSTypeFromString(name), true, nil
}

// DSWithLegacyFields contains a Delivery Service along with associated fields
// that have been removed from its representation(s).
type DSWithLegacyFields struct {
	DS        tc.DeliveryServiceV5
	LongDesc1 *string
	LongDesc2 *string
}

// GetDeliveryServices can be used with a valid query and map of values to
// retrieve Delivery Services from the database. Note that Tenancy must be built
// into the query already. The returned Delivery Services come with accompanying
// information about fields that have been removed from the model, in case a
// legacy handler needs that information.
func GetDeliveryServices(query string, queryValues map[string]interface{}, tx *sqlx.Tx) ([]DSWithLegacyFields, error, error, int) {
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, fmt.Errorf("querying: %w", err), http.StatusInternalServerError
	}
	defer log.Close(rows, "getting Delivery Services")

	dses := []DSWithLegacyFields{}
	dsCDNDomains := map[string]string{}

	// ensure json generated from this slice won't come out as `null` if empty
	dsQueryParams := []string{}

	geoLimitCountries := new(string)
	for rows.Next() {
		var ds tc.DeliveryServiceV5
		var longDesc1 *string
		var longDesc2 *string
		var cdnDomain string
		err := rows.Scan(
			&ds.Active,
			&ds.AnonymousBlockingEnabled,
			&ds.CCRDNSTTL,
			&ds.CDNID,
			&ds.CDNName,
			&ds.CheckPath,
			&ds.ConsistentHashRegex,
			&ds.DeepCachingType,
			&ds.DisplayName,
			&ds.DNSBypassCNAME,
			&ds.DNSBypassIP,
			&ds.DNSBypassIP6,
			&ds.DNSBypassTTL,
			&ds.DSCP,
			&ds.EcsEnabled,
			&ds.EdgeHeaderRewrite,
			&ds.FirstHeaderRewrite,
			&ds.GeoLimitRedirectURL,
			&ds.GeoLimit,
			&geoLimitCountries,
			&ds.GeoProvider,
			&ds.GlobalMaxMBPS,
			&ds.GlobalMaxTPS,
			&ds.FQPacingRate,
			&ds.HTTPBypassFQDN,
			&ds.ID,
			&ds.InfoURL,
			&ds.InitialDispersion,
			&ds.InnerHeaderRewrite,
			&ds.IPV6RoutingEnabled,
			&ds.LastHeaderRewrite,
			&ds.LastUpdated,
			&ds.LogsEnabled,
			&ds.LongDesc,
			&longDesc1,
			&longDesc2,
			&ds.MaxDNSAnswers,
			&ds.MaxOriginConnections,
			&ds.MaxRequestHeaderBytes,
			&ds.MidHeaderRewrite,
			&ds.MissLat,
			&ds.MissLong,
			&ds.MultiSiteOrigin,
			&ds.OrgServerFQDN,
			&ds.OriginShield,
			&ds.ProfileID,
			&ds.ProfileName,
			&ds.ProfileDesc,
			&ds.Protocol,
			&ds.QStringIgnore,
			pq.Array(&dsQueryParams),
			&ds.RangeRequestHandling,
			&ds.RegexRemap,
			&ds.Regional,
			&ds.RegionalGeoBlocking,
			&ds.RemapText,
			pq.Array(&ds.RequiredCapabilities),
			&ds.RoutingName,
			&ds.ServiceCategory,
			&ds.SigningAlgorithm,
			&ds.RangeSliceBlockSize,
			&ds.SSLKeyVersion,
			&ds.TenantID,
			&ds.Tenant,
			pq.Array(&ds.TLSVersions),
			&ds.Topology,
			&ds.TRRequestHeaders,
			&ds.TRResponseHeaders,
			&ds.Type,
			&ds.TypeID,
			&ds.XMLID,
			&cdnDomain)

		if err != nil {
			return nil, nil, fmt.Errorf("getting delivery services: %w", err), http.StatusInternalServerError
		}

		if geoLimitCountries != nil && *geoLimitCountries != "" {
			geo := strings.Split(*geoLimitCountries, ",")
			ds.GeoLimitCountries = geo
		}
		ds.ConsistentHashQueryParams = make([]string, 0, len(dsQueryParams))
		// ensure unique and in consistent order
		m := make(map[string]struct{}, len(dsQueryParams))
		for _, k := range dsQueryParams {
			if _, exists := m[k]; exists {
				continue
			}
			m[k] = struct{}{}
			ds.ConsistentHashQueryParams = append(ds.ConsistentHashQueryParams, k)
		}

		dsCDNDomains[ds.XMLID] = cdnDomain
		ds.DeepCachingType = tc.DeepCachingTypeFromString(string(ds.DeepCachingType))

		ds.Signed = ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == tc.SigningAlgorithmURLSig

		if len(ds.TLSVersions) < 1 {
			ds.TLSVersions = nil
		}

		dses = append(dses, DSWithLegacyFields{
			DS:        ds,
			LongDesc1: longDesc1,
			LongDesc2: longDesc2,
		})
	}

	dsNames := make([]string, len(dses))
	for i, ds := range dses {
		dsNames[i] = ds.DS.XMLID
	}

	matchLists, err := GetDeliveryServicesMatchLists(dsNames, tx.Tx)
	if err != nil {
		return nil, nil, fmt.Errorf("getting delivery service matchlists: %w", err), http.StatusInternalServerError
	}
	for i, d := range dses {
		ds := d.DS
		matchList, ok := matchLists[ds.XMLID]
		if !ok {
			continue
		}
		ds.MatchList = matchList
		ds.ExampleURLs = MakeExampleURLs(ds.Protocol, tc.DSType(*ds.Type), ds.RoutingName, ds.MatchList, dsCDNDomains[ds.XMLID])
		d.DS = ds
		dses[i] = d
	}

	return dses, nil, nil, http.StatusOK
}

func getCDNDomain(dsID int, tx *sql.Tx) (string, error) {
	q := `SELECT cdn.domain_name from cdn where cdn.id = (SELECT ds.cdn_id from deliveryservice as ds where ds.id = $1)`
	cdnDomain := ""
	if err := tx.QueryRow(q, dsID).Scan(&cdnDomain); err != nil {
		return "", fmt.Errorf("getting CDN domain for delivery service '%v': "+err.Error(), dsID)
	}
	return cdnDomain, nil
}

func getCDNNameDomainDNSSecEnabled(dsID int, tx *sql.Tx) (string, string, bool, error) {
	q := `SELECT cdn.name, cdn.domain_name, cdn.dnssec_enabled from cdn where cdn.id = (SELECT ds.cdn_id from deliveryservice as ds where ds.id = $1)`
	cdnName := ""
	cdnDomain := ""
	dnssecEnabled := false
	if err := tx.QueryRow(q, dsID).Scan(&cdnName, &cdnDomain, &dnssecEnabled); err != nil {
		return "", "", false, fmt.Errorf("getting dnssec_enabled for delivery service '%v': "+err.Error(), dsID)
	}
	return cdnName, cdnDomain, dnssecEnabled, nil
}

// MakeExampleURLs creates the example URLs for a delivery service. The dsProtocol may be nil, if the delivery service type doesn't have a protocol (e.g. ANY_MAP).
func MakeExampleURLs(protocol *int, dsType tc.DSType, routingName string, matchList []tc.DeliveryServiceMatch, cdnDomain string) []string {
	examples := []string{}
	scheme := ""
	scheme2 := ""
	if protocol != nil {
		switch *protocol {
		case 0:
			scheme = string(tc.ProtocolHTTP)
		case 1:
			scheme = string(tc.ProtocolHTTPS)
		case 2:
			fallthrough
		case 3:
			scheme = string(tc.ProtocolHTTP)
			scheme2 = string(tc.ProtocolHTTPS)
		default:
			scheme = string(tc.ProtocolHTTP)
		}
	} else {
		scheme = string(tc.ProtocolHTTP)
	}
	dsIsDNS := dsType.IsDNS()
	regexReplacer := strings.NewReplacer(`\`, ``, `.*`, ``, `.`, ``)
	for _, match := range matchList {
		if dsIsDNS || match.Type == tc.DSMatchTypeHostRegex {
			host := regexReplacer.Replace(match.Pattern)
			if match.SetNumber == 0 {
				examples = append(examples, scheme+`://`+routingName+`.`+host+`.`+cdnDomain)
				if scheme2 != "" {
					examples = append(examples, scheme2+`://`+routingName+`.`+host+`.`+cdnDomain)
				}
				continue
			}
			examples = append(examples, scheme+`://`+match.Pattern)
			if scheme2 != "" {
				examples = append(examples, scheme2+`://`+match.Pattern)
			}
		} else if match.Type == tc.DSMatchTypePathRegex {
			examples = append(examples, match.Pattern)
		}
	}
	return examples
}

// GetDeliveryServicesMatchLists retrieves "Match Lists" for each of the
// Delivery Services having the provided XMLID(s). The error return value is not
// safe to return to clients.
func GetDeliveryServicesMatchLists(dses []string, tx *sql.Tx) (map[string][]tc.DeliveryServiceMatch, error) {
	// TODO move somewhere generic
	q := `
SELECT ds.xml_id as ds_name, t.name as type, r.pattern, COALESCE(dsr.set_number, 0)
FROM regex as r
JOIN deliveryservice_regex as dsr ON dsr.regex = r.id
JOIN deliveryservice as ds on ds.id = dsr.deliveryservice
JOIN type as t ON r.type = t.id
WHERE ds.xml_id = ANY($1)
ORDER BY dsr.set_number
`
	rows, err := tx.Query(q, pq.Array(dses))
	if err != nil {
		return nil, fmt.Errorf("getting delivery service regexes: %w", err)
	}
	defer log.Close(rows, "getting Delivery Service Match Lists")

	matches := make(map[string][]tc.DeliveryServiceMatch, len(dses))
	for rows.Next() {
		m := tc.DeliveryServiceMatch{}
		dsName := ""
		matchTypeStr := ""
		if err := rows.Scan(&dsName, &matchTypeStr, &m.Pattern, &m.SetNumber); err != nil {
			return nil, fmt.Errorf("scanning delivery service regexes: %w", err)
		}
		matchType := tc.DSMatchTypeFromString(matchTypeStr)
		if matchType == tc.DSMatchTypeInvalid {
			return nil, errors.New("getting delivery service regexes: got invalid delivery service match type '" + matchTypeStr + "'")
		}
		m.Type = matchType
		matches[dsName] = append(matches[dsName], m)
	}
	return matches, nil
}

// EnsureParams ensures the given delivery service's necessary parameters exist on profiles of servers assigned to the delivery service.
// Note the edgeHeaderRewrite, midHeaderRewrite, regexRemap may be nil, if the delivery service does not have those values.
func EnsureParams(tx *sql.Tx, dsID int, xmlID string, edgeHeaderRewrite *string, midHeaderRewrite *string, regexRemap *string, signingAlgorithm *string, dsType tc.DSType, maxOriginConns *int) error {
	if err := ensureHeaderRewriteParams(tx, dsID, xmlID, edgeHeaderRewrite, edgeTier, dsType, maxOriginConns); err != nil {
		return fmt.Errorf("creating edge header rewrite parameters: %w", err)
	}
	if err := ensureHeaderRewriteParams(tx, dsID, xmlID, midHeaderRewrite, midTier, dsType, maxOriginConns); err != nil {
		return fmt.Errorf("creating mid header rewrite parameters: %w", err)
	}
	if err := ensureRegexRemapParams(tx, dsID, xmlID, regexRemap); err != nil {
		return fmt.Errorf("creating mid regex remap parameters: %w", err)
	}
	if err := ensureURLSigParams(tx, dsID, xmlID, signingAlgorithm); err != nil {
		return fmt.Errorf("creating urlsig parameters: %w", err)
	}
	return nil
}

// EnsureCacheURLParams ensures the given delivery service's cachrurl parameters exist on profiles of servers assigned to the delivery service.
func EnsureCacheURLParams(tx *sql.Tx, dsID int, xmlID string, cacheURL *string) error {
	configFile := "cacheurl_" + xmlID + ".config"
	if cacheURL == nil || *cacheURL == "" {
		return deleteLocationParam(tx, configFile)
	}
	locationParamID, err := ensureLocation(tx, configFile)
	if err != nil {
		return err
	}
	return createDSLocationProfileParams(tx, locationParamID, dsID)
}

func ensureHeaderRewriteParams(tx *sql.Tx, dsID int, xmlID string, hdrRW *string, tier tierType, dsType tc.DSType, maxOriginConns *int) error {
	configFile := "hdr_rw_" + xmlID + ".config"
	if tier == midTier {
		configFile = "hdr_rw_mid_" + xmlID + ".config"
	}

	if tier == midTier && dsType.IsLive() && !dsType.IsNational() {
		// live local DSes don't get header rewrite rules on the mid so cleanup any location params related to mids
		return deleteLocationParam(tx, configFile)
	}

	hasMaxOriginConns := *maxOriginConns > 0 && ((tier == midTier) == dsType.UsesMidCache())
	if (hdrRW == nil || *hdrRW == "") && !hasMaxOriginConns {
		return deleteLocationParam(tx, configFile)
	}
	locationParamID, err := ensureLocation(tx, configFile)
	if err != nil {
		return err
	}
	if tier != midTier {
		return createDSLocationProfileParams(tx, locationParamID, dsID)
	}
	profileParameterQuery := `
INSERT INTO profile_parameter (profile, parameter)
SELECT DISTINCT(profile), $1::bigint FROM server
WHERE server.type IN (SELECT id from type where type.name like 'MID%' and type.use_in_table = 'server')
AND server.cdn_id = (select cdn_id from deliveryservice where id = $2)
ON CONFLICT DO NOTHING
`
	if _, err := tx.Exec(profileParameterQuery, locationParamID, dsID); err != nil {
		return fmt.Errorf("parameter query to insert profile_parameters query '"+profileParameterQuery+"' location parameter ID '%v' delivery service ID '%v': %v", locationParamID, dsID, err)
	}
	return nil
}

func ensureURLSigParams(tx *sql.Tx, dsID int, xmlID string, signingAlgorithm *string) error {
	configFile := "url_sig_" + xmlID + ".config"
	if signingAlgorithm == nil || *signingAlgorithm != tc.SigningAlgorithmURLSig {
		return deleteLocationParam(tx, configFile)
	}
	locationParamID, err := ensureLocation(tx, configFile)
	if err != nil {
		return err
	}
	return createDSLocationProfileParams(tx, locationParamID, dsID)
}

func ensureRegexRemapParams(tx *sql.Tx, dsID int, xmlID string, regexRemap *string) error {
	configFile := "regex_remap_" + xmlID + ".config"
	if regexRemap == nil || *regexRemap == "" {
		return deleteLocationParam(tx, configFile)
	}
	locationParamID, err := ensureLocation(tx, configFile)
	if err != nil {
		return err
	}
	return createDSLocationProfileParams(tx, locationParamID, dsID)
}

// createDSLocationProfileParams adds the given parameter to all profiles assigned to servers which are assigned to the given delivery service.
func createDSLocationProfileParams(tx *sql.Tx, locationParamID int, deliveryServiceID int) error {
	profileParameterQuery := `
INSERT INTO profile_parameter (profile, parameter)
SELECT DISTINCT(profile), $1::bigint FROM server
WHERE server.id IN (SELECT server from deliveryservice_server where deliveryservice = $2)
ON CONFLICT DO NOTHING
`
	if _, err := tx.Exec(profileParameterQuery, locationParamID, deliveryServiceID); err != nil {
		return fmt.Errorf("inserting profile_parameters: %w", err)
	}
	return nil
}

// ensureLocation ensures a location parameter exists for the given config file.
// If not, it creates one, with the same value as the 'remap.config' file
// parameter. Returns the ID of the location parameter.
func ensureLocation(tx *sql.Tx, configFile string) (int, error) {
	atsConfigLocation := ""
	if err := tx.QueryRow(`SELECT value FROM parameter WHERE name = 'location' AND config_file = 'remap.config'`).Scan(&atsConfigLocation); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("executing parameter query for ATS config location: parameter missing (do you have a name=location config_file=remap.config parameter?)")
		}
		return 0, fmt.Errorf("executing parameter query for ATS config location: %w", err)
	}
	atsConfigLocation = strings.TrimRight(atsConfigLocation, `/`)

	locationParamID := 0
	existingLocationErr := tx.QueryRow(`SELECT id FROM parameter WHERE name = 'location' AND config_file = $1`, configFile).Scan(&locationParamID)
	if existingLocationErr != nil && existingLocationErr != sql.ErrNoRows {
		return 0, fmt.Errorf("executing parameter query for existing location: %w", existingLocationErr)
	}

	if errors.Is(existingLocationErr, sql.ErrNoRows) {
		resultRows, err := tx.Query(`INSERT INTO parameter (config_file, name, value) VALUES ($1, 'location', $2) RETURNING id`, configFile, atsConfigLocation)
		if err != nil {
			return 0, fmt.Errorf("executing parameter query to insert location: %w", err)
		}
		defer log.Close(resultRows, "inserting location parameter for '"+configFile+"'")
		if !resultRows.Next() {
			return 0, errors.New("parameter query to insert location didn't return id")
		}
		if err := resultRows.Scan(&locationParamID); err != nil {
			return 0, fmt.Errorf("parameter query to insert location returned id scan: %w", err)
		}
		if resultRows.Next() {
			return 0, errors.New("parameter query to insert location returned too many rows (>1)")
		}
	}
	return locationParamID, nil
}

func deleteLocationParam(tx *sql.Tx, configFile string) error {
	id := 0
	err := tx.QueryRow(`DELETE FROM parameter WHERE name = 'location' AND config_file = $1 RETURNING id`, configFile).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		log.Errorln("deleting name=location config_file=" + configFile + " parameter: " + err.Error())
		return fmt.Errorf("executing parameter delete: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM profile_parameter WHERE parameter = $1`, id); err != nil {
		log.Errorf("deleting parameter name=location config_file=%v id=%v profile_parameter: %v", configFile, id, err)
		return fmt.Errorf("executing parameter profile_parameter delete: %w", err)
	}
	return nil
}

// getTenantID returns a pointer to the tenant ID of the given Delivery Service,
// or any error encountered while determining said ID.
// This will panic if the transaction is nil.
func getTenantID(tx *sql.Tx, ds tc.DeliveryServiceV5) (*int, error) {
	if ds.ID != nil {
		id, _, err := getDSTenantIDByID(tx, *ds.ID)
		return id, err
	}
	id, _, err := getDSTenantIDByName(tx, tc.DeliveryServiceName(ds.XMLID))
	return id, err
}

func isTenantAuthorized(inf *api.Info, ds *tc.DeliveryServiceV5) (bool, error) {
	tx := inf.Tx.Tx
	user := inf.User

	existingID, err := getTenantID(inf.Tx.Tx, *ds)
	if err != nil {
		return false, fmt.Errorf("getting tenant ID: %w", err)
	}
	if ds.TenantID <= 0 {
		if existingID == nil {
			return false, errors.New("checking Tenant authorization: Delivery Service doesn't have a specified Tenant ID and doesn't already exist")
		}
		ds.TenantID = *existingID
	}
	if existingID != nil && *existingID != ds.TenantID {
		userAuthorizedForExistingDSTenant, err := tenant.IsResourceAuthorizedToUserTx(*existingID, user, tx)
		if err != nil {
			return false, fmt.Errorf("checking authorization for existing DS ID: %w", err)
		}
		if !userAuthorizedForExistingDSTenant {
			return false, nil
		}
	}
	userAuthorizedForNewDSTenant, err := tenant.IsResourceAuthorizedToUserTx(ds.TenantID, user, tx)
	if err != nil {
		return false, fmt.Errorf("checking authorization for new DS ID: %w", err)
	}
	return userAuthorizedForNewDSTenant, nil
}

// getDSTenantIDByID returns the tenant ID, whether the delivery service exists, and any error.
func getDSTenantIDByID(tx *sql.Tx, id int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where id = $1`, id).Scan(&tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service ID '%v': %w", id, err)
	}
	return tenantID, true, nil
}

// getDSTenantIDByName returns the tenant ID, whether the delivery service exists, and any error.
func getDSTenantIDByName(tx *sql.Tx, ds tc.DeliveryServiceName) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, ds).Scan(&tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service name '%v': %w", ds, err)
	}
	return tenantID, true, nil
}

// GetXMLID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func GetXMLID(tx *sql.Tx, id int) (string, bool, error) {
	xmlID := ""
	if err := tx.QueryRow(`SELECT xml_id FROM deliveryservice where id = $1`, id).Scan(&xmlID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("querying xml_id for delivery service ID '%v': %w", id, err)
	}
	return xmlID, true, nil
}

// getSSLVersion reports a boolean value, confirming whether DS has a SSL
// version or not.
func getSSLVersion(xmlId string, tx *sql.Tx) (bool, error) {
	var exists bool
	row := tx.QueryRow(`SELECT EXISTS(SELECT * FROM deliveryservice WHERE xml_id = $1 AND ssl_key_version>=1)`, xmlId)
	err := row.Scan(&exists)
	return exists, err
}

func setNilIfEmpty(ptrs ...**string) {
	for _, s := range ptrs {
		if *s != nil && strings.TrimSpace(**s) == "" {
			*s = nil
		}
	}
}

func sanitizeGeoLimitCountries(ds *tc.DeliveryServiceV5) {
	if len(ds.GeoLimitCountries) == 0 {
		return
	}
	for i, cc := range ds.GeoLimitCountries {
		ds.GeoLimitCountries[i] = strings.ToUpper(strings.ReplaceAll(cc, " ", ""))
	}
}

func sanitize(ds *tc.DeliveryServiceV5) {
	sanitizeGeoLimitCountries(ds)
	if ds.ProfileID != nil && *ds.ProfileID == -1 {
		ds.ProfileID = nil
	}
	setNilIfEmpty(
		&ds.EdgeHeaderRewrite,
		&ds.MidHeaderRewrite,
		&ds.FirstHeaderRewrite,
		&ds.InnerHeaderRewrite,
		&ds.LastHeaderRewrite,
	)
	if ds.RoutingName == "" {
		ds.RoutingName = tc.DefaultRoutingName
	}
	signedAlgorithm := tc.SigningAlgorithmURLSig
	if ds.Signed && (ds.SigningAlgorithm == nil || *ds.SigningAlgorithm == "") {
		ds.SigningAlgorithm = &signedAlgorithm
	}
	if !ds.Signed && ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == signedAlgorithm {
		ds.Signed = true
	}
	if ds.MaxOriginConnections == nil || *ds.MaxOriginConnections < 0 {
		ds.MaxOriginConnections = util.IntPtr(0)
	}
	ds.DeepCachingType = tc.DeepCachingTypeFromString(string(ds.DeepCachingType))
	if ds.MaxRequestHeaderBytes == nil {
		ds.MaxRequestHeaderBytes = util.IntPtr(tc.DefaultMaxRequestHeaderBytes)
	}
}

// SelectDeliveryServicesQuery is a PostgreSQL query used to fetch Delivery
// Services from the Traffic Ops Database.
const SelectDeliveryServicesQuery = `
SELECT
ds.active,
	ds.anonymous_blocking_enabled,
	ds.ccr_dns_ttl,
	ds.cdn_id,
	cdn.name AS cdnName,
	ds.check_path,
	ds.consistent_hash_regex,
	CAST(ds.deep_caching_type AS text) AS deep_caching_type,
	ds.display_name,
	ds.dns_bypass_cname,
	ds.dns_bypass_ip,
	ds.dns_bypass_ip6,
	ds.dns_bypass_ttl,
	ds.dscp,
	ds.ecs_enabled,
	ds.edge_header_rewrite,
	ds.first_header_rewrite,
	ds.geolimit_redirect_url,
	ds.geo_limit,
	ds.geo_limit_countries,
	ds.geo_provider,
	ds.global_max_mbps,
	ds.global_max_tps,
	ds.fq_pacing_rate,
	ds.http_bypass_fqdn,
	ds.id,
	ds.info_url,
	ds.initial_dispersion,
	ds.inner_header_rewrite,
	ds.ipv6_routing_enabled,
	ds.last_header_rewrite,
	ds.last_updated,
	ds.logs_enabled,
	COALESCE(ds.long_desc, '') AS long_desc,
	ds.long_desc_1,
	ds.long_desc_2,
	ds.max_dns_answers,
	ds.max_origin_connections,
	ds.max_request_header_bytes,
	ds.mid_header_rewrite,
	COALESCE(ds.miss_lat, 0.0),
	COALESCE(ds.miss_long, 0.0),
	ds.multi_site_origin,
	(SELECT o.protocol::::text || ':://' || o.fqdn || rtrim(concat('::', o.port::::text), '::')
		FROM origin o
		WHERE o.deliveryservice = ds.id
		AND o.is_primary) AS org_server_fqdn,
	ds.origin_shield,
	ds.profile AS profileID,
	profile.name AS profile_name,
	profile.description  AS profile_description,
	ds.protocol,
	ds.qstring_ignore,
	(SELECT ARRAY_AGG(name ORDER BY name)
		FROM deliveryservice_consistent_hash_query_param
				WHERE deliveryservice_id = ds.id) AS query_keys,
	ds.range_request_handling,
	ds.regex_remap,
	ds.regional,
	ds.regional_geo_blocking,
	ds.remap_text,
	COALESCE(ds.required_capabilities, '{}'),
	ds.routing_name,
	ds.service_category,
	ds.signing_algorithm,
	ds.range_slice_block_size,
	ds.ssl_key_version,
	ds.tenant_id,
	tenant.name,
	(` + baseTLSVersionsQuery + ` WHERE deliveryservice = ds.id) AS tls_versions,
	ds.topology,
	ds.tr_request_headers,
	ds.tr_response_headers,
	type.name,
	ds.type AS type_id,
	ds.xml_id,
	cdn.domain_name AS cdn_domain
FROM deliveryservice AS ds
JOIN type ON ds.type = type.id
JOIN cdn ON ds.cdn_id = cdn.id
LEFT JOIN profile ON ds.profile = profile.id
LEFT JOIN tenant ON ds.tenant_id = tenant.id
`

func updateDSQuery() string {
	return `
UPDATE
deliveryservice SET
active=$1,
ccr_dns_ttl=$2,
cdn_id=$3,
check_path=$4,
deep_caching_type=$5,
display_name=$6,
dns_bypass_cname=$7,
dns_bypass_ip=$8,
dns_bypass_ip6=$9,
dns_bypass_ttl=$10,
dscp=$11,
edge_header_rewrite=$12,
geolimit_redirect_url=$13,
geo_limit=$14,
geo_limit_countries=$15,
geo_provider=$16,
global_max_mbps=$17,
global_max_tps=$18,
fq_pacing_rate=$19,
http_bypass_fqdn=$20,
info_url=$21,
initial_dispersion=$22,
ipv6_routing_enabled=$23,
logs_enabled=$24,
long_desc=$25,
long_desc_1=$26,
long_desc_2=$27,
max_dns_answers=$28,
mid_header_rewrite=$29,
miss_lat=$30,
miss_long=$31,
multi_site_origin=$32,
origin_shield=$33,
profile=$34,
protocol=$35,
qstring_ignore=$36,
range_request_handling=$37,
regex_remap=$38,
regional=$39,
regional_geo_blocking=$40,
remap_text=$41,
routing_name=$42,
signing_algorithm=$43,
ssl_key_version=$44,
tenant_id=$45,
tr_request_headers=$46,
tr_response_headers=$47,
type=$48,
xml_id=$49,
anonymous_blocking_enabled=$50,
consistent_hash_regex=$51,
max_origin_connections=$52,
ecs_enabled=$53,
range_slice_block_size=$54,
topology=$55,
first_header_rewrite=$56,
inner_header_rewrite=$57,
last_header_rewrite=$58,
service_category=$59,
max_request_header_bytes=$60,
required_capabilities=$61
WHERE id=$62
RETURNING last_updated
`
}

func updateDSQueryWithoutLD1AndLD2() string {
	return `
UPDATE
deliveryservice SET
active=$1,
ccr_dns_ttl=$2,
cdn_id=$3,
check_path=$4,
deep_caching_type=$5,
display_name=$6,
dns_bypass_cname=$7,
dns_bypass_ip=$8,
dns_bypass_ip6=$9,
dns_bypass_ttl=$10,
dscp=$11,
edge_header_rewrite=$12,
geolimit_redirect_url=$13,
geo_limit=$14,
geo_limit_countries=$15,
geo_provider=$16,
global_max_mbps=$17,
global_max_tps=$18,
fq_pacing_rate=$19,
http_bypass_fqdn=$20,
info_url=$21,
initial_dispersion=$22,
ipv6_routing_enabled=$23,
logs_enabled=$24,
long_desc=$25,
max_dns_answers=$26,
mid_header_rewrite=$27,
miss_lat=$28,
miss_long=$29,
multi_site_origin=$30,
origin_shield=$31,
profile=$32,
protocol=$33,
qstring_ignore=$34,
range_request_handling=$35,
regex_remap=$36,
regional=$37,
regional_geo_blocking=$38,
remap_text=$39,
routing_name=$40,
signing_algorithm=$41,
ssl_key_version=$42,
tenant_id=$43,
tr_request_headers=$44,
tr_response_headers=$45,
type=$46,
xml_id=$47,
anonymous_blocking_enabled=$48,
consistent_hash_regex=$49,
max_origin_connections=$50,
ecs_enabled=$51,
range_slice_block_size=$52,
topology=$53,
first_header_rewrite=$54,
inner_header_rewrite=$55,
last_header_rewrite=$56,
service_category=$57,
max_request_header_bytes=$58,
required_capabilities=$59
WHERE id=$60
RETURNING last_updated
`
}

func insertQuery() string {
	return `
INSERT INTO deliveryservice (
active,
anonymous_blocking_enabled,
ccr_dns_ttl,
cdn_id,
check_path,
consistent_hash_regex,
deep_caching_type,
display_name,
dns_bypass_cname,
dns_bypass_ip,
dns_bypass_ip6,
dns_bypass_ttl,
dscp,
edge_header_rewrite,
geolimit_redirect_url,
geo_limit,
geo_limit_countries,
geo_provider,
global_max_mbps,
global_max_tps,
fq_pacing_rate,
http_bypass_fqdn,
info_url,
initial_dispersion,
ipv6_routing_enabled,
logs_enabled,
long_desc,
long_desc_1,
long_desc_2,
max_dns_answers,
max_origin_connections,
mid_header_rewrite,
miss_lat,
miss_long,
multi_site_origin,
origin_shield,
profile,
protocol,
qstring_ignore,
range_request_handling,
regex_remap,
regional,
regional_geo_blocking,
remap_text,
routing_name,
signing_algorithm,
ssl_key_version,
tenant_id,
topology,
tr_request_headers,
tr_response_headers,
type,
xml_id,
ecs_enabled,
range_slice_block_size,
first_header_rewrite,
inner_header_rewrite,
last_header_rewrite,
service_category,
max_request_header_bytes,
required_capabilities
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,$61)
RETURNING id, last_updated
`
}

func insertQueryWithoutLD1AndLD2() string {
	return `
INSERT INTO deliveryservice (
active,
anonymous_blocking_enabled,
ccr_dns_ttl,
cdn_id,
check_path,
consistent_hash_regex,
deep_caching_type,
display_name,
dns_bypass_cname,
dns_bypass_ip,
dns_bypass_ip6,
dns_bypass_ttl,
dscp,
edge_header_rewrite,
geolimit_redirect_url,
geo_limit,
geo_limit_countries,
geo_provider,
global_max_mbps,
global_max_tps,
fq_pacing_rate,
http_bypass_fqdn,
info_url,
initial_dispersion,
ipv6_routing_enabled,
logs_enabled,
long_desc,
max_dns_answers,
max_origin_connections,
mid_header_rewrite,
miss_lat,
miss_long,
multi_site_origin,
origin_shield,
profile,
protocol,
qstring_ignore,
range_request_handling,
regex_remap,
regional,
regional_geo_blocking,
remap_text,
routing_name,
signing_algorithm,
ssl_key_version,
tenant_id,
topology,
tr_request_headers,
tr_response_headers,
type,
xml_id,
ecs_enabled,
range_slice_block_size,
first_header_rewrite,
inner_header_rewrite,
last_header_rewrite,
service_category,
max_request_header_bytes,
required_capabilities
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59)
RETURNING id, last_updated
`
}
