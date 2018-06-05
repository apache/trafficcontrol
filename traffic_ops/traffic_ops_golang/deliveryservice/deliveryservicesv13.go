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
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-util"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/utils"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TODeliveryServiceV13 struct {
	tc.DeliveryServiceNullableV13
	Cfg config.Config
	DB  *sqlx.DB
}

func (ds *TODeliveryServiceV13) V12() *TODeliveryServiceV12 {
	return &TODeliveryServiceV12{DeliveryServiceNullableV12: ds.DeliveryServiceNullableV12, DB: ds.DB, Cfg: ds.Cfg}
}

func (ds TODeliveryServiceV13) MarshalJSON() ([]byte, error) {
	return json.Marshal(ds.DeliveryServiceNullableV13)
}
func (ds *TODeliveryServiceV13) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, ds.DeliveryServiceNullableV13)
}

func GetRefTypeV13(cfg config.Config, db *sqlx.DB) *TODeliveryServiceV13 {
	return &TODeliveryServiceV13{Cfg: cfg, DB: db}
}

func (ds TODeliveryServiceV13) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return ds.V12().GetKeyFieldsInfo()
}

//Implementation of the Identifier, Validator interface functions
func (ds TODeliveryServiceV13) GetKeys() (map[string]interface{}, bool) {
	return ds.V12().GetKeys()
}

func (ds *TODeliveryServiceV13) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ds.ID = &i
}

func (ds *TODeliveryServiceV13) GetAuditName() string {
	return ds.V12().GetAuditName()
}

func (ds *TODeliveryServiceV13) GetType() string {
	return ds.V12().GetType()
}

func ValidateV13(db *sqlx.DB, ds *tc.DeliveryServiceNullableV13) []error {
	if ds == nil {
		return []error{}
	}
	tods := TODeliveryServiceV13{DeliveryServiceNullableV13: *ds, DB: db} // TODO set Cfg?
	return tods.Validate(db)
}

func (ds *TODeliveryServiceV13) Sanitize(db *sqlx.DB) { sanitizeV13(&ds.DeliveryServiceNullableV13) }

func sanitizeV13(ds *tc.DeliveryServiceNullableV13) {
	sanitizeV12(&ds.DeliveryServiceNullableV12)
	signedAlgorithm := "url_sig"
	if ds.Signed && (ds.SigningAlgorithm == nil || *ds.SigningAlgorithm == "") {
		ds.SigningAlgorithm = &signedAlgorithm
	}
	if !ds.Signed && ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == signedAlgorithm {
		ds.Signed = true
	}
	if ds.DeepCachingType == nil {
		s := tc.DeepCachingType("")
		ds.DeepCachingType = &s
	}
	*ds.DeepCachingType = tc.DeepCachingTypeFromString(string(*ds.DeepCachingType))
}

func (ds *TODeliveryServiceV13) Validate(db *sqlx.DB) []error {
	return validateV13(db, &ds.DeliveryServiceNullableV13)
}

func validateV13(db *sqlx.DB, ds *tc.DeliveryServiceNullableV13) []error {
	sanitizeV13(ds)
	neverOrAlways := validation.NewStringRule(tovalidate.IsOneOfStringICase("NEVER", "ALWAYS"),
		"must be one of 'NEVER' or 'ALWAYS'")
	errs := tovalidate.ToErrors(validation.Errors{
		"deepCachingType": validation.Validate(ds.DeepCachingType, neverOrAlways),
	})
	oldErrs := validateV12(db, &ds.DeliveryServiceNullableV12)
	return append(errs, oldErrs...)
}

// Create implements the Creator interface.
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a ds with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted ds and have
//to be added to the struct
// func (ds *TODeliveryServiceV13) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) { //
//
// 	TODO allow users to post names (type, cdn, etc) and get the IDs from the names. This isn't trivial to do in a single query, without dynamically building the entire insert query, and ideally inserting would be one query. But it'd be much more convenient for users. Alternatively, remove IDs from the database entirely and use real candidate keys.
func CreateV13(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}

		ds := tc.DeliveryServiceNullableV13{}
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}

		if ds.RoutingName == nil || *ds.RoutingName == "" {
			ds.RoutingName = utils.StrPtr("cdn")
		}

		if errs := validateV13(db, &ds); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}

		if authorized, err := isTenantAuthorized(*user, db, &ds.DeliveryServiceNullableV12); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
			return
		} else if !authorized {
			api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
			return
		}

		ds, errCode, userErr, sysErr := create(db.DB, cfg, user, ds)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		api.WriteResp(w, r, []tc.DeliveryServiceNullableV13{ds})
	}
}

// create creates the given ds in the database, and returns the DS with its id and other fields created on insert set. On error, the HTTP status cdoe, user error, and system error are returned. The status code SHOULD NOT be used, if both errors are nil.
func create(db *sql.DB, cfg config.Config, user *auth.CurrentUser, ds tc.DeliveryServiceNullableV13) (tc.DeliveryServiceNullableV13, int, error, error) {
	// TODO change DeepCachingType to implement sql.Valuer and sql.Scanner, so sqlx struct scan can be used.
	deepCachingType := tc.DeepCachingType("").String()
	if ds.DeepCachingType != nil {
		deepCachingType = ds.DeepCachingType.String() // necessary, because DeepCachingType's default needs to insert the string, not "", and Query doesn't call .String().
	}

	tx, err := db.Begin()
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("could not begin transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	resultRows, err := tx.Query(insertQuery(), &ds.Active, &ds.CacheURL, &ds.CCRDNSTTL, &ds.CDNID, &ds.CheckPath, &deepCachingType, &ds.DisplayName, &ds.DNSBypassCNAME, &ds.DNSBypassIP, &ds.DNSBypassIP6, &ds.DNSBypassTTL, &ds.DSCP, &ds.EdgeHeaderRewrite, &ds.GeoLimitRedirectURL, &ds.GeoLimit, &ds.GeoLimitCountries, &ds.GeoProvider, &ds.GlobalMaxMBPS, &ds.GlobalMaxTPS, &ds.FQPacingRate, &ds.HTTPBypassFQDN, &ds.InfoURL, &ds.InitialDispersion, &ds.IPV6RoutingEnabled, &ds.LogsEnabled, &ds.LongDesc, &ds.LongDesc1, &ds.LongDesc2, &ds.MaxDNSAnswers, &ds.MidHeaderRewrite, &ds.MissLat, &ds.MissLong, &ds.MultiSiteOrigin, &ds.OrgServerFQDN, &ds.OriginShield, &ds.ProfileID, &ds.Protocol, &ds.QStringIgnore, &ds.RangeRequestHandling, &ds.RegexRemap, &ds.RegionalGeoBlocking, &ds.RemapText, &ds.RoutingName, &ds.SigningAlgorithm, &ds.SSLKeyVersion, &ds.TenantID, &ds.TRRequestHeaders, &ds.TRResponseHeaders, &ds.TypeID, &ds.XMLID)

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			err, _ := dbhelpers.ParsePQUniqueConstraintError(pqerr)
			return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("a delivery service with " + err.Error())
		}
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("inserting ds: " + err.Error())
	}

	id := 0
	lastUpdated := tc.TimeNoMod{}
	if !resultRows.Next() {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("no deliveryservice request inserted, no id was returned")
	}
	if err := resultRows.Scan(&id, &lastUpdated); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("could not scan id from insert: " + err.Error())
	}
	if resultRows.Next() {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("too many ids returned from deliveryservice request insert")
	}
	ds.ID = &id

	if ds.ID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing id after insert")
	}
	if ds.XMLID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing xml_id after insert")
	}
	if ds.TypeID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing type after insert")
	}
	dsType, err := getTypeNameFromID(*ds.TypeID, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting delivery service type: " + err.Error())
	}
	ds.Type = &dsType

	if err := createDefaultRegex(tx, *ds.ID, *ds.XMLID); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating default regex: " + err.Error())
	}

	matchlists, err := readGetDeliveryServicesMatchLists([]string{*ds.XMLID}, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating DS: reading matchlists: " + err.Error())
	}
	if matchlist, ok := matchlists[*ds.XMLID]; !ok {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating DS: reading matchlists: not found")
	} else {
		ds.MatchList = &matchlist
	}

	cdnName, cdnDomain, dnssecEnabled, err := getCDNNameDomainDNSSecEnabled(*ds.ID, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating DS: getting CDN info: " + err.Error())
	}

	ds.ExampleURLs = makeExampleURLs(ds.Protocol, *ds.Type, *ds.RoutingName, *ds.MatchList, cdnDomain)

	if err := ensureHeaderRewriteParams(tx, *ds.ID, *ds.XMLID, ds.EdgeHeaderRewrite, edgeTier, *ds.Type); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating edge header rewrite parameters: " + err.Error())
	}
	if err := ensureHeaderRewriteParams(tx, *ds.ID, *ds.XMLID, ds.MidHeaderRewrite, midTier, *ds.Type); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating mid header rewrite parameters: " + err.Error())
	}
	if err := ensureRegexRemapParams(tx, *ds.ID, *ds.XMLID, ds.RegexRemap); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating regex remap parameters: " + err.Error())
	}
	if err := ensureCacheURLParams(tx, *ds.ID, *ds.XMLID, ds.CacheURL); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating cache url parameters: " + err.Error())
	}
	if err := createDNSSecKeys(tx, cfg, *ds.ID, *ds.XMLID, cdnName, cdnDomain, dnssecEnabled, ds.ExampleURLs); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating DNSSEC keys: " + err.Error())
	}
	ds.LastUpdated = &lastUpdated
	commitTx = true
	api.CreateChangeLogRaw(api.ApiChange, "Created ds: "+*ds.XMLID+" id: "+strconv.Itoa(*ds.ID), *user, db)
	return ds, http.StatusOK, nil, nil
}

func (ds *TODeliveryServiceV13) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	returnable := []interface{}{}
	dses, errs, errType := readGetDeliveryServices(params, db, user)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` { // TODO create const for string
				return nil, []error{errors.New("Resource not found.")}, tc.DataMissingError //matches perl response
			}
		}
		return nil, errs, errType
	}

	for _, ds := range dses {
		returnable = append(returnable, ds)
	}
	return returnable, nil, tc.NoError
}

func createDefaultRegex(tx *sql.Tx, dsID int, xmlID string) error {
	regexStr := `.*\.` + xmlID + `\..*`
	regexID := 0
	if err := tx.QueryRow(`INSERT INTO regex (type, pattern) VALUES ((select id from type where name = 'HOST_REGEXP'), $1::text) RETURNING id`, regexStr).Scan(&regexID); err != nil {
		return errors.New("insert regex: " + err.Error())
	}
	if _, err := tx.Exec(`INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES ($1::bigint, $2::bigint, 0)`, dsID, regexID); err != nil {
		return errors.New("executing parameter query to insert location: " + err.Error())
	}
	return nil
}

func getOldHostName(id int, tx *sql.Tx) (string, error) {
	q := `
SELECT ds.xml_id, ds.protocol, type.name, ds.routing_name, cdn.domain_name
FROM  deliveryservice as ds
JOIN type ON ds.type = type.id
JOIN cdn ON ds.cdn_id = cdn.id
WHERE ds.id=$1
`
	xmlID := ""
	protocol := (*int)(nil)
	dsType := ""
	routingName := ""
	cdnDomain := ""
	if err := tx.QueryRow(q, id).Scan(&xmlID, &protocol, &dsType, &routingName, &cdnDomain); err != nil {
		return "", fmt.Errorf("querying delivery service %v host name: "+err.Error()+"\n", id)
	}
	matchLists, err := readGetDeliveryServicesMatchLists([]string{xmlID}, tx)
	if err != nil {
		return "", errors.New("getting delivery services matchlist: " + err.Error())
	}
	matchList, ok := matchLists[xmlID]
	if !ok {
		return "", errors.New("delivery service has no match lists (is your delivery service missing regexes?)")
	}
	host, err := getHostName(protocol, dsType, routingName, matchList, cdnDomain) // protocol defaults to 0: doesn't need to check Valid()
	if err != nil {
		return "", errors.New("getting hostname: " + err.Error())
	}
	return host, nil
}

func getTypeNameFromID(id int, tx *sql.Tx) (string, error) {
	// TODO combine with getOldHostName, to only make one query?
	name := ""
	if err := tx.QueryRow(`SELECT name FROM type WHERE id = $1`, id).Scan(&name); err != nil {
		return "", fmt.Errorf("querying type ID %v: "+err.Error()+"\n", id)
	}
	return name, nil
}

func UpdateV13(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}

		params, _, userErr, sysErr, errCode := api.AllParams(r, nil)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		if strings.HasSuffix(params["id"], ".json") {
			params["id"] = params["id"][:len(params["id"])-len(".json")]
		}
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("id must be an integer"), sysErr)
		}

		ds := tc.DeliveryServiceNullableV13{}
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}
		ds.ID = &id

		if errs := validateV13(db, &ds); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}

		if authorized, err := isTenantAuthorized(*user, db, &ds.DeliveryServiceNullableV12); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
			return
		} else if !authorized {
			api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
			return
		}

		ds, errCode, userErr, sysErr = update(db.DB, cfg, *user, &ds)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		api.WriteResp(w, r, []tc.DeliveryServiceNullableV13{ds})
	}
}

func update(db *sql.DB, cfg config.Config, user auth.CurrentUser, ds *tc.DeliveryServiceNullableV13) (tc.DeliveryServiceNullableV13, int, error, error) {
	tx, err := db.Begin()
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("could not begin transaction: " + err.Error())
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	if ds.XMLID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusBadRequest, errors.New("missing xml_id"), nil
	}
	if ds.ID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusBadRequest, errors.New("missing id"), nil
	}

	oldHostName, err := getOldHostName(*ds.ID, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting existing delivery service hostname: " + err.Error())
	}

	// TODO change DeepCachingType to implement sql.Valuer and sql.Scanner, so sqlx struct scan can be used.
	deepCachingType := tc.DeepCachingType("").String()
	if ds.DeepCachingType != nil {
		deepCachingType = ds.DeepCachingType.String() // necessary, because DeepCachingType's default needs to insert the string, not "", and Query doesn't call .String().
	}

	resultRows, err := tx.Query(updateDSQuery(), &ds.Active, &ds.CacheURL, &ds.CCRDNSTTL, &ds.CDNID, &ds.CheckPath, &deepCachingType, &ds.DisplayName, &ds.DNSBypassCNAME, &ds.DNSBypassIP, &ds.DNSBypassIP6, &ds.DNSBypassTTL, &ds.DSCP, &ds.EdgeHeaderRewrite, &ds.GeoLimitRedirectURL, &ds.GeoLimit, &ds.GeoLimitCountries, &ds.GeoProvider, &ds.GlobalMaxMBPS, &ds.GlobalMaxTPS, &ds.FQPacingRate, &ds.HTTPBypassFQDN, &ds.InfoURL, &ds.InitialDispersion, &ds.IPV6RoutingEnabled, &ds.LogsEnabled, &ds.LongDesc, &ds.LongDesc1, &ds.LongDesc2, &ds.MaxDNSAnswers, &ds.MidHeaderRewrite, &ds.MissLat, &ds.MissLong, &ds.MultiSiteOrigin, &ds.OrgServerFQDN, &ds.OriginShield, &ds.ProfileID, &ds.Protocol, &ds.QStringIgnore, &ds.RangeRequestHandling, &ds.RegexRemap, &ds.RegionalGeoBlocking, &ds.RemapText, &ds.RoutingName, &ds.SigningAlgorithm, &ds.SSLKeyVersion, &ds.TenantID, &ds.TRRequestHeaders, &ds.TRResponseHeaders, &ds.TypeID, &ds.XMLID, &ds.ID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			if eType == tc.DataConflictError {
				return tc.DeliveryServiceNullableV13{}, http.StatusBadRequest, errors.New("a delivery service with " + err.Error()), nil
			}
			return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("query updating delivery service: pq: " + err.Error())
		}
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("query updating delivery service: " + err.Error())
	}
	if !resultRows.Next() {
		return tc.DeliveryServiceNullableV13{}, http.StatusNotFound, errors.New("no delivery service found with this id"), nil
	}
	lastUpdated := tc.TimeNoMod{}
	if err := resultRows.Scan(&lastUpdated); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("scan updating delivery service: " + err.Error())
	}
	if resultRows.Next() {
		xmlID := ""
		if ds.XMLID != nil {
			xmlID = *ds.XMLID
		}
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("updating delivery service " + xmlID + ": " + "this update affected too many rows: > 1")
	}

	if ds.ID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing id after update")
	}
	if ds.XMLID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing xml_id after update")
	}
	if ds.TypeID == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing type after update")
	}
	if ds.RoutingName == nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("missing routing name after update")
	}
	dsType, err := getTypeNameFromID(*ds.TypeID, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting delivery service type after update: " + err.Error())
	}
	ds.Type = &dsType

	cdnDomain, err := getCDNDomain(*ds.ID, db) // need to get the domain again, in case it changed.
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting CDN domain after update: " + err.Error())
	}

	newHostName, err := getHostName(ds.Protocol, *ds.Type, *ds.RoutingName, *ds.MatchList, cdnDomain)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting hostname after update: " + err.Error())
	}

	matchLists, err := readGetDeliveryServicesMatchLists([]string{*ds.XMLID}, tx)
	if err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("getting matchlists after update: " + err.Error())
	}
	if ml, ok := matchLists[*ds.XMLID]; !ok {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("no matchlists after update")
	} else {
		ds.MatchList = &ml
	}

	if oldHostName != newHostName {
		if err := updateSSLKeys(ds, newHostName, db, cfg); err != nil {
			return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("updating delivery service " + *ds.XMLID + ": updating SSL keys: " + err.Error())
		}
	}

	if err := ensureHeaderRewriteParams(tx, *ds.ID, *ds.XMLID, ds.EdgeHeaderRewrite, edgeTier, *ds.Type); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating edge header rewrite parameters: " + err.Error())
	}
	if err := ensureHeaderRewriteParams(tx, *ds.ID, *ds.XMLID, ds.MidHeaderRewrite, midTier, *ds.Type); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating mid header rewrite parameters: " + err.Error())
	}
	if err := ensureRegexRemapParams(tx, *ds.ID, *ds.XMLID, ds.RegexRemap); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating mid regex remap parameters: " + err.Error())
	}
	if err := ensureCacheURLParams(tx, *ds.ID, *ds.XMLID, ds.CacheURL); err != nil {
		return tc.DeliveryServiceNullableV13{}, http.StatusInternalServerError, nil, errors.New("creating mid cacheurl parameters: " + err.Error())
	}
	ds.LastUpdated = &lastUpdated
	commitTx = true
	api.CreateChangeLogRaw(api.ApiChange, "Updated ds: "+*ds.XMLID+" id: "+strconv.Itoa(*ds.ID), user, db)
	return *ds, http.StatusOK, nil, nil
}

//The DeliveryService implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (ds *TODeliveryServiceV13) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	log.Debugln("TODeliveryServiceV13.Delete calling id '%v' xmlid '%v'\n", ds.ID, ds.XMLID)
	// return nil, tc.NoError // debug

	tx, err := db.Begin()
	if err != nil {
		log.Errorln("could not begin transaction: " + err.Error())
		return tc.DBError, tc.SystemError
	}
	commitTx := false
	defer dbhelpers.FinishTx(tx, &commitTx)

	if err != nil {
		log.Errorln("could not begin transaction: " + err.Error())
		return tc.DBError, tc.SystemError
	}

	if ds.ID == nil {
		log.Errorln("TODeliveryServiceV13.Delete called with nil ID")
		return tc.DBError, tc.DataMissingError
	}
	xmlID, ok, err := ds.V12().GetXMLID(db)
	if err != nil {
		log.Errorln("TODeliveryServiceV13.Delete ID '" + string(*ds.ID) + "' loading XML ID: " + err.Error())
		return tc.DBError, tc.SystemError
	}
	if !ok {
		log.Errorln("TODeliveryServiceV13.Delete ID '" + string(*ds.ID) + "' had no delivery service!")
		return tc.DBError, tc.DataMissingError
	}
	ds.XMLID = &xmlID

	// Note ds regexes MUST be deleted before the ds, because there's a ON DELETE CASCADE on deliveryservice_regex (but not on regex).
	// Likewise, it MUST happen in a transaction with the later DS delete, so they aren't deleted if the DS delete fails.
	if _, err := tx.Exec(`DELETE FROM regex WHERE id IN (SELECT regex FROM deliveryservice_regex WHERE deliveryservice=$1)`, *ds.ID); err != nil {
		log.Errorln("TODeliveryServiceV13.Delete deleting regexes for delivery service: " + err.Error())
		return tc.DBError, tc.SystemError
	}

	if _, err := tx.Exec(`DELETE FROM deliveryservice_regex WHERE deliveryservice=$1`, *ds.ID); err != nil {
		log.Errorln("TODeliveryServiceV13.Delete deleting delivery service regexes: " + err.Error())
		return tc.DBError, tc.SystemError
	}

	result, err := tx.Exec(`DELETE FROM deliveryservice WHERE id=$1`, *ds.ID)
	if err != nil {
		log.Errorln("TODeliveryServiceV13.Delete deleting delivery service: " + err.Error())
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no delivery service with that id found"), tc.DataMissingError
		}
		return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
	}

	paramConfigFilePrefixes := []string{"hdr_rw_", "hdr_rw_mid_", "regex_remap_", "cacheurl_"}
	configFiles := []string{}
	for _, prefix := range paramConfigFilePrefixes {
		configFiles = append(configFiles, prefix+*ds.XMLID+".config")
	}

	if _, err := tx.Exec(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`, pq.Array(configFiles)); err != nil {
		log.Errorln("TODeliveryServiceV13.Delete deleting delivery service parameters: " + err.Error())
		return tc.DBError, tc.SystemError
	}

	commitTx = true
	return nil, tc.NoError
}

// IsTenantAuthorized implements the Tenantable interface to ensure the user is authorized on the deliveryservice tenant
func (ds *TODeliveryServiceV13) IsTenantAuthorized(user auth.CurrentUser, db *sqlx.DB) (bool, error) {
	return ds.V12().IsTenantAuthorized(user, db)
}

func filterAuthorized(dses []tc.DeliveryServiceNullableV13, user auth.CurrentUser, db *sqlx.DB) ([]tc.DeliveryServiceNullableV13, error) {
	newDSes := []tc.DeliveryServiceNullableV13{}
	for _, ds := range dses {
		// TODO add/use a helper func to make a single SQL call, for performance
		ok, err := tenant.IsResourceAuthorizedToUser(*ds.TenantID, user, db)
		if err != nil {
			if ds.XMLID == nil {
				return nil, errors.New("isResourceAuthorized for delivery service with nil XML ID: " + err.Error())
			} else {
				return nil, errors.New("isResourceAuthorized for '" + *ds.XMLID + "': " + err.Error())
			}
		}
		if !ok {
			continue
		}
		newDSes = append(newDSes, ds)
	}
	return newDSes, nil
}

func readGetDeliveryServices(params map[string]string, db *sqlx.DB, user auth.CurrentUser) ([]tc.DeliveryServiceNullableV13, []error, tc.ApiErrorType) {
	if strings.HasSuffix(params["id"], ".json") {
		params["id"] = params["id"][:len(params["id"])-len(".json")]
	}
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":       dbhelpers.WhereColumnInfo{"ds.id", api.IsInt},
		"hostName": dbhelpers.WhereColumnInfo{"s.host_name", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	if tenant.IsTenancyEnabled(db) {
		log.Debugln("Tenancy is enabled")
		tenantIDs, err := tenant.GetUserTenantIDList(user, db)
		if err != nil {
			log.Errorln("received error querying for user's tenants: " + err.Error())
			return nil, []error{tc.DBError}, tc.SystemError
		}
		where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "ds.tenant_id", tenantIDs)
	}
	query := selectQuery() + where + orderBy

	log.Debugln("generated deliveryServices query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	tx, err := db.Beginx()
	if err != nil {
		log.Errorln("could not begin transaction: " + err.Error())
		return nil, []error{tc.DBError}, tc.SystemError
	}
	commitTx := false
	defer dbhelpers.FinishTxX(tx, &commitTx)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{fmt.Errorf("querying: %v", err)}, tc.SystemError
	}
	defer rows.Close()

	dses := []tc.DeliveryServiceNullableV13{}
	dsCDNDomains := map[string]string{}
	for rows.Next() {
		ds := tc.DeliveryServiceNullableV13{}
		cdnDomain := ""
		err := rows.Scan(&ds.Active, &ds.CacheURL, &ds.CCRDNSTTL, &ds.CDNID, &ds.CDNName, &ds.CheckPath, &ds.DeepCachingType, &ds.DisplayName, &ds.DNSBypassCNAME, &ds.DNSBypassIP, &ds.DNSBypassIP6, &ds.DNSBypassTTL, &ds.DSCP, &ds.EdgeHeaderRewrite, &ds.GeoLimitRedirectURL, &ds.GeoLimit, &ds.GeoLimitCountries, &ds.GeoProvider, &ds.GlobalMaxMBPS, &ds.GlobalMaxTPS, &ds.FQPacingRate, &ds.HTTPBypassFQDN, &ds.ID, &ds.InfoURL, &ds.InitialDispersion, &ds.IPV6RoutingEnabled, &ds.LastUpdated, &ds.LogsEnabled, &ds.LongDesc, &ds.LongDesc1, &ds.LongDesc2, &ds.MaxDNSAnswers, &ds.MidHeaderRewrite, &ds.MissLat, &ds.MissLong, &ds.MultiSiteOrigin, &ds.OrgServerFQDN, &ds.OriginShield, &ds.ProfileID, &ds.ProfileName, &ds.ProfileDesc, &ds.Protocol, &ds.QStringIgnore, &ds.RangeRequestHandling, &ds.RegexRemap, &ds.RegionalGeoBlocking, &ds.RemapText, &ds.RoutingName, &ds.SigningAlgorithm, &ds.SSLKeyVersion, &ds.TenantID, &ds.Tenant, &ds.TRRequestHeaders, &ds.TRResponseHeaders, &ds.Type, &ds.TypeID, &ds.XMLID, &cdnDomain)
		if err != nil {
			return nil, []error{fmt.Errorf("getting delivery services: %v", err)}, tc.SystemError
		}
		dsCDNDomains[*ds.XMLID] = cdnDomain
		if ds.DeepCachingType != nil {
			*ds.DeepCachingType = tc.DeepCachingTypeFromString(string(*ds.DeepCachingType))
		}
		ds.Signed = ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == "url_sig"
		dses = append(dses, ds)
	}

	dsNames := make([]string, len(dses), len(dses))
	for i, ds := range dses {
		dsNames[i] = *ds.XMLID
	}

	matchLists, err := readGetDeliveryServicesMatchLists(dsNames, tx.Tx)
	if err != nil {
		return nil, []error{errors.New("getting delivery service matchlists: " + err.Error())}, tc.SystemError
	}
	for i, ds := range dses {
		matchList, ok := matchLists[*ds.XMLID]
		if !ok {
			continue
		}
		ds.MatchList = &matchList
		ds.ExampleURLs = makeExampleURLs(ds.Protocol, *ds.Type, *ds.RoutingName, *ds.MatchList, dsCDNDomains[*ds.XMLID])
		dses[i] = ds
	}

	commitTx = true
	return dses, nil, tc.NoError
}

func updateSSLKeys(ds *tc.DeliveryServiceNullableV13, hostName string, db *sql.DB, cfg config.Config) error {
	if ds.XMLID == nil {
		return errors.New("delivery services has no XMLID!")
	}
	key, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(*ds.XMLID, "latest", db, cfg.RiakAuthOptions)
	if err != nil {
		return errors.New("getting SSL key: " + err.Error())
	}
	if !ok {
		return nil // no keys to update
	}
	key.DeliveryService = *ds.XMLID
	key.Hostname = hostName
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(key, db, cfg.RiakAuthOptions); err != nil {
		return errors.New("putting updated SSL key: " + err.Error())
	}
	return nil
}

// getHostName gets the host name used for delivery service requests. The dsProtocol may be nil, if the delivery service type doesn't have a protocol (e.g. ANY_MAP).
func getHostName(dsProtocol *int, dsType string, dsRoutingName string, dsMatchList []tc.DeliveryServiceMatch, cdnDomain string) (string, error) {
	exampleURLs := makeExampleURLs(dsProtocol, dsType, dsRoutingName, dsMatchList, cdnDomain)

	exampleURL := ""
	if dsProtocol != nil && *dsProtocol == 2 {
		if len(exampleURLs) < 2 {
			return "", errors.New("missing example URLs (does your delivery service have matchsets?)")
		}
		exampleURL = exampleURLs[1]
	} else {
		if len(exampleURLs) < 1 {
			return "", errors.New("missing example URLs (does your delivery service have matchsets?)")
		}
		exampleURL = exampleURLs[0]
	}

	host := strings.NewReplacer(`http://`, ``, `https://`, ``).Replace(exampleURL)
	if strings.HasPrefix(dsType, "HTTP") {
		if firstDot := strings.Index(host, "."); firstDot == -1 {
			host = "*" // TODO warn? error?
		} else {
			host = "*" + host[firstDot:]
		}
	}
	return host, nil
}

func getCDNDomain(dsID int, db *sql.DB) (string, error) {
	q := `SELECT cdn.domain_name from cdn where cdn.id = (SELECT ds.cdn_id from deliveryservice as ds where ds.id = $1)`
	cdnDomain := ""
	if err := db.QueryRow(q, dsID).Scan(&cdnDomain); err != nil {
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

// makeExampleURLs creates the example URLs for a delivery service. The dsProtocol may be nil, if the delivery service type doesn't have a protocol (e.g. ANY_MAP).
func makeExampleURLs(protocol *int, dsType string, routingName string, matchList []tc.DeliveryServiceMatch, cdnDomain string) []string {
	examples := []string{}
	scheme := ""
	scheme2 := ""
	if protocol != nil {
		switch *protocol {
		case 0:
			scheme = "http"
		case 1:
			scheme = "https"
		case 2:
			fallthrough
		case 3:
			scheme = "http"
			scheme2 = "https"
		default:
			scheme = "http"
		}
	} else {
		scheme = "http"
	}
	dsIsDNS := strings.HasPrefix(strings.ToLower(dsType), "DNS")
	regexReplacer := strings.NewReplacer(`\`, ``, `.*`, ``, `.`, ``)
	for _, match := range matchList {
		switch {
		case dsIsDNS:
			fallthrough
		case match.Type == `HOST_REGEXP`:
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
		case match.Type == `PATH_REGEXP`:
			examples = append(examples, match.Pattern)
		}
	}
	return examples
}

func readGetDeliveryServicesMatchLists(dses []string, tx *sql.Tx) (map[string][]tc.DeliveryServiceMatch, error) {
	q := `
SELECT ds.xml_id as ds_name, t.name as type, r.pattern, COALESCE(dsr.set_number, 0)
FROM regex as r
JOIN deliveryservice_regex as dsr ON dsr.regex = r.id
JOIN deliveryservice as ds on ds.id = dsr.deliveryservice
JOIN type as t ON r.type = t.id
WHERE ds.xml_id = ANY($1)
`
	rows, err := tx.Query(q, pq.Array(dses))
	if err != nil {
		return nil, errors.New("getting delivery service regexes: " + err.Error())
	}

	matches := map[string][]tc.DeliveryServiceMatch{}
	for rows.Next() {
		m := tc.DeliveryServiceMatch{}
		dsName := ""
		if err := rows.Scan(&dsName, &m.Type, &m.Pattern, &m.SetNumber); err != nil {
			return nil, errors.New("scanning delivery service regexes: " + err.Error())
		}
		matches[dsName] = append(matches[dsName], m)
	}
	return matches, nil
}

type tierType int

const (
	midTier tierType = iota
	edgeTier
)

func ensureHeaderRewriteParams(tx *sql.Tx, dsID int, xmlID string, hdrRW *string, tier tierType, dsType string) error {
	if tier == midTier && strings.Contains(dsType, "LIVE") && !strings.Contains(dsType, "NATNL") {
		return nil // live local DSes don't get remap rules
	}
	configFile := "hdr_rw_" + xmlID + ".config"
	if tier == midTier {
		configFile = "hdr_rw_mid_" + xmlID + ".config"
	}
	if hdrRW == nil || *hdrRW == "" {
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

func ensureCacheURLParams(tx *sql.Tx, dsID int, xmlID string, cacheURL *string) error {
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

// createDSLocationProfileParams adds the given parameter to all profiles assigned to servers which are assigned to the given delivery service.
func createDSLocationProfileParams(tx *sql.Tx, locationParamID int, deliveryServiceID int) error {
	profileParameterQuery := `
INSERT INTO profile_parameter (profile, parameter)
SELECT DISTINCT(profile), $1::bigint FROM server
WHERE server.id IN (SELECT server from deliveryservice_server where deliveryservice = $2)
ON CONFLICT DO NOTHING
`
	if _, err := tx.Exec(profileParameterQuery, locationParamID, deliveryServiceID); err != nil {
		return errors.New("inserting profile_parameters: " + err.Error())
	}
	return nil
}

// ensureLocation ensures a location parameter exists for the given config file. If not, it creates one, with the same value as the 'remap.config' file parameter. Returns the ID of the location parameter.
func ensureLocation(tx *sql.Tx, configFile string) (int, error) {
	atsConfigLocation := ""
	if err := tx.QueryRow(`SELECT value FROM parameter WHERE name = 'location' AND config_file = 'remap.config'`).Scan(&atsConfigLocation); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("executing parameter query for ATS config location: parameter missing (do you have a name=location config_file=remap.config parameter?")
		}
		return 0, errors.New("executing parameter query for ATS config location: " + err.Error())
	}
	atsConfigLocation = strings.TrimRight(atsConfigLocation, `/`)

	locationParamID := 0
	existingLocationErr := tx.QueryRow(`SELECT id FROM parameter WHERE name = 'location' AND config_file = $1`, configFile).Scan(&locationParamID)
	if existingLocationErr != nil && existingLocationErr != sql.ErrNoRows {
		return 0, errors.New("executing parameter query for existing location: " + existingLocationErr.Error())
	}

	if existingLocationErr == sql.ErrNoRows {
		resultRows, err := tx.Query(`INSERT INTO parameter (config_file, name, value) VALUES ($1, 'location', $2) RETURNING id`, configFile, atsConfigLocation)
		if err != nil {
			return 0, errors.New("executing parameter query to insert location: " + err.Error())
		}
		defer resultRows.Close()
		if !resultRows.Next() {
			return 0, errors.New("parameter query to insert location didn't return id")
		}
		if err := resultRows.Scan(&locationParamID); err != nil {
			return 0, errors.New("parameter query to insert location returned id scan: " + err.Error())
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
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Errorln("deleting name=location config_file=" + configFile + " parameter: " + err.Error())
		return errors.New("executing parameter delete: " + err.Error())
	}
	if _, err := tx.Exec(`DELETE FROM profile_parameter WHERE parameter = $1`, id); err != nil {
		log.Errorf("deleting parameter name=location config_file=%v id=%v profile_parameter: %v", configFile, id, err)
		return errors.New("executing parameter profile_parameter delete: " + err.Error())
	}
	return nil
}

func selectQuery() string {
	return `
SELECT
ds.active,
ds.cacheurl,
ds.ccr_dns_ttl,
ds.cdn_id,
cdn.name as cdnName,
ds.check_path,
CAST(ds.deep_caching_type AS text) as deep_caching_type,
ds.display_name,
ds.dns_bypass_cname,
ds.dns_bypass_ip,
ds.dns_bypass_ip6,
ds.dns_bypass_ttl,
ds.dscp,
ds.edge_header_rewrite,
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
ds.ipv6_routing_enabled,
ds.last_updated,
ds.logs_enabled,
ds.long_desc,
ds.long_desc_1,
ds.long_desc_2,
ds.max_dns_answers,
ds.mid_header_rewrite,
COALESCE(ds.miss_lat, 0.0),
COALESCE(ds.miss_long, 0.0),
ds.multi_site_origin,
ds.org_server_fqdn,
ds.origin_shield,
ds.profile as profileID,
profile.name as profile_name,
profile.description  as profile_description,
ds.protocol,
ds.qstring_ignore,
ds.range_request_handling,
ds.regex_remap,
ds.regional_geo_blocking,
ds.remap_text,
ds.routing_name,
ds.signing_algorithm,
ds.ssl_key_version,
ds.tenant_id,
tenant.name,
ds.tr_request_headers,
ds.tr_response_headers,
type.name,
ds.type as type_id,
ds.xml_id,
cdn.domain_name as cdn_domain
from deliveryservice as ds
JOIN type ON ds.type = type.id
JOIN cdn ON ds.cdn_id = cdn.id
LEFT JOIN profile ON ds.profile = profile.id
LEFT JOIN tenant ON ds.tenant_id = tenant.id
`
}

func updateDSQuery() string {
	return `
UPDATE
deliveryservice SET
active=$1,
cacheurl=$2,
ccr_dns_ttl=$3,
cdn_id=$4,
check_path=$5,
deep_caching_type=$6,
display_name=$7,
dns_bypass_cname=$8,
dns_bypass_ip=$9,
dns_bypass_ip6=$10,
dns_bypass_ttl=$11,
dscp=$12,
edge_header_rewrite=$13,
geolimit_redirect_url=$14,
geo_limit=$15,
geo_limit_countries=$16,
geo_provider=$17,
global_max_mbps=$18,
global_max_tps=$19,
fq_pacing_rate=$20,
http_bypass_fqdn=$21,
info_url=$22,
initial_dispersion=$23,
ipv6_routing_enabled=$24,
logs_enabled=$25,
long_desc=$26,
long_desc_1=$27,
long_desc_2=$28,
max_dns_answers=$29,
mid_header_rewrite=$30,
miss_lat=$31,
miss_long=$32,
multi_site_origin=$33,
org_server_fqdn=$34,
origin_shield=$35,
profile=$36,
protocol=$37,
qstring_ignore=$38,
range_request_handling=$39,
regex_remap=$40,
regional_geo_blocking=$41,
remap_text=$42,
routing_name=$43,
signing_algorithm=$44,
ssl_key_version=$45,
tenant_id=$46,
tr_request_headers=$47,
tr_response_headers=$48,
type=$49,
xml_id=$50
WHERE id=$51
RETURNING last_updated
`
}

func insertQuery() string {
	return `
INSERT INTO deliveryservice (
active,
cacheurl,
ccr_dns_ttl,
cdn_id,
check_path,
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
mid_header_rewrite,
miss_lat,
miss_long,
multi_site_origin,
org_server_fqdn,
origin_shield,
profile,
protocol,
qstring_ignore,
range_request_handling,
regex_remap,
regional_geo_blocking,
remap_text,
routing_name,
signing_algorithm,
ssl_key_version,
tenant_id,
tr_request_headers,
tr_response_headers,
type,
xml_id
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50)
RETURNING id, last_updated
`
}
