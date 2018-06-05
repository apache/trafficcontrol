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
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-util"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/utils"

	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
)

type TODeliveryServiceV12 struct {
	tc.DeliveryServiceNullableV12
	Cfg config.Config
	DB  *sqlx.DB
}

func GetRefTypeV12(cfg config.Config, db *sqlx.DB) *TODeliveryServiceV12 {
	return &TODeliveryServiceV12{Cfg: cfg, DB: db}
}

func (ds TODeliveryServiceV12) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

func (ds TODeliveryServiceV12) GetKeys() (map[string]interface{}, bool) {
	if ds.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *ds.ID}, true
}

func (ds *TODeliveryServiceV12) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ds.ID = &i
}

func (ds *TODeliveryServiceV12) GetAuditName() string {
	if ds.XMLID != nil {
		return *ds.XMLID
	}
	return ""
}

func (ds *TODeliveryServiceV12) GetType() string {
	return "ds"
}

func ValidateV12(db *sqlx.DB, ds *tc.DeliveryServiceNullableV12) []error {
	if ds == nil {
		return []error{}
	}
	tods := TODeliveryServiceV12{DB: db} // TODO pass config?
	return tods.Validate(db)
}

func (ds *TODeliveryServiceV12) Sanitize(db *sqlx.DB) {
	sanitizeV12(&ds.DeliveryServiceNullableV12)
}

func sanitizeV12(ds *tc.DeliveryServiceNullableV12) {
	if ds.GeoLimitCountries != nil {
		*ds.GeoLimitCountries = strings.ToUpper(strings.Replace(*ds.GeoLimitCountries, " ", "", -1))
	}
	if ds.ProfileID != nil && *ds.ProfileID == -1 {
		ds.ProfileID = nil
	}
	if ds.EdgeHeaderRewrite != nil && strings.TrimSpace(*ds.EdgeHeaderRewrite) == "" {
		ds.EdgeHeaderRewrite = nil
	}
	if ds.MidHeaderRewrite != nil && strings.TrimSpace(*ds.MidHeaderRewrite) == "" {
		ds.MidHeaderRewrite = nil
	}
}

// getDSTenantIDByID returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
func getDSTenantIDByID(db *sql.DB, id int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := db.QueryRow(`SELECT tenant_id FROM deliveryservice where id = $1`, id).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service ID '%v': %v", id, err)
	}
	return tenantID, true, nil
}

// getDSTenantIDByName returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
func getDSTenantIDByName(db *sql.DB, name string) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := db.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, name).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service name '%v': %v", name, err)
	}
	return tenantID, true, nil
}

// GetXMLID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func (ds *TODeliveryServiceV12) GetXMLID(db *sqlx.DB) (string, bool, error) {
	if ds.ID == nil {
		return "", false, errors.New("missing ID")
	}
	xmlID := ""
	if err := db.QueryRow(`SELECT xml_id FROM deliveryservice where id = $1`, ds.ID).Scan(&xmlID); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("querying xml_id for delivery service ID '%v': %v", *ds.ID, err)
	}
	return xmlID, true, nil
}

// IsTenantAuthorized checks that the user is authorized for both the delivery service's existing tenant, and the new tenant they're changing it to (if different).
func (ds *TODeliveryServiceV12) IsTenantAuthorized(user auth.CurrentUser, db *sqlx.DB) (bool, error) {
	return isTenantAuthorized(user, db, &ds.DeliveryServiceNullableV12)
}

func isTenantAuthorized(user auth.CurrentUser, db *sqlx.DB, ds *tc.DeliveryServiceNullableV12) (bool, error) {
	if ds.ID == nil && ds.XMLID == nil {
		return false, errors.New("delivery service has no ID or XMLID")
	}

	existingID, err := (*int)(nil), error(nil)
	if ds.ID != nil {
		existingID, _, err = getDSTenantIDByID(db.DB, *ds.ID) // ignore exists return - if the DS is new, we only need to check the user input tenant
	} else {
		existingID, _, err = getDSTenantIDByName(db.DB, *ds.XMLID) // ignore exists return - if the DS is new, we only need to check the user input tenant
	}
	if err != nil {
		return false, errors.New("getting tenant ID: " + err.Error())
	}

	if ds.TenantID == nil {
		ds.TenantID = existingID
	}
	if existingID != nil && existingID != ds.TenantID {
		userAuthorizedForExistingDSTenant, err := tenant.IsResourceAuthorizedToUser(*existingID, user, db)
		if err != nil {
			return false, errors.New("checking authorization for existing DS ID: " + err.Error())
		}
		if !userAuthorizedForExistingDSTenant {
			return false, nil
		}
	}
	if ds.TenantID != nil {
		userAuthorizedForNewDSTenant, err := tenant.IsResourceAuthorizedToUser(*ds.TenantID, user, db)
		if err != nil {
			return false, errors.New("checking authorization for new DS ID: " + err.Error())
		}
		if !userAuthorizedForNewDSTenant {
			return false, nil
		}
	}
	return true, nil
}

func (ds *TODeliveryServiceV12) Validate(db *sqlx.DB) []error {
	return validateV12(db, &ds.DeliveryServiceNullableV12)
}

func validateV12(db *sqlx.DB, ds *tc.DeliveryServiceNullableV12) []error {
	sanitizeV12(ds)
	// Custom Examples:
	// Just add isCIDR as a parameter to Validate()
	// isCIDR := validation.NewStringRule(govalidator.IsCIDR, "must be a valid CIDR address")
	isDNSName := validation.NewStringRule(govalidator.IsDNSName, "must be a valid hostname")
	noPeriods := validation.NewStringRule(tovalidate.NoPeriods, "cannot contain periods")
	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	// Validate that the required fields are sent first to prevent panics below
	errs := validation.Errors{
		"active":              validation.Validate(ds.Active, validation.NotNil),
		"cdnId":               validation.Validate(ds.CDNID, validation.Required),
		"displayName":         validation.Validate(ds.DisplayName, validation.Required, validation.Length(1, 48)),
		"dscp":                validation.Validate(ds.DSCP, validation.NotNil, validation.Min(0)),
		"geoLimit":            validation.Validate(ds.GeoLimit, validation.NotNil),
		"geoProvider":         validation.Validate(ds.GeoProvider, validation.NotNil),
		"logsEnabled":         validation.Validate(ds.LogsEnabled, validation.NotNil),
		"regionalGeoBlocking": validation.Validate(ds.RegionalGeoBlocking, validation.NotNil),
		"routingName":         validation.Validate(ds.RoutingName, isDNSName, noPeriods, validation.Length(1, 48)),
		"typeId":              validation.Validate(ds.TypeID, validation.Required, validation.Min(1)),
		"xmlId":               validation.Validate(ds.XMLID, noSpaces, noPeriods, validation.Length(1, 48)),
	}
	toErrs := tovalidate.ToErrors(errs)
	if fieldErrs := validateTypeFields(db, ds); len(fieldErrs) > 0 {
		toErrs = append(toErrs, fieldErrs...)
	}
	if len(toErrs) > 0 {
		return toErrs
	}
	return nil
}

func validateTypeFields(db *sqlx.DB, ds *tc.DeliveryServiceNullableV12) []error {
	// Validate the TypeName related fields below
	var typeName string
	var err error
	DNSRegexType := "^DNS.*$"
	HTTPRegexType := "^HTTP.*$"
	SteeringRegexType := "^STEERING.*$"

	if ds.TypeID == nil {
		return []error{errors.New("missing typeID")}
	}

	typeName, ok, err := getTypeName(db, *ds.TypeID)
	if err != nil {
		return []error{err}
	}
	if !ok {
		return []error{errors.New("type not found")}
	}

	errs := validation.Errors{
		"initialDispersion": validation.Validate(ds.InitialDispersion,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"ipv6RoutingEnabled": validation.Validate(ds.IPV6RoutingEnabled,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"missLat": validation.Validate(ds.MissLat,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"missLong": validation.Validate(ds.MissLong,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"multiSiteOrigin": validation.Validate(ds.MultiSiteOrigin,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"orgServerFqdn": validation.Validate(ds.OrgServerFQDN,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"protocol": validation.Validate(ds.Protocol,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"qstringIgnore": validation.Validate(ds.QStringIgnore,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"rangeRequestHandling": validation.Validate(ds.RangeRequestHandling,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
	}
	toErrs := tovalidate.ToErrors(errs)
	if len(toErrs) > 0 {
		return toErrs
	}
	return nil
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
		case *string:
			if v != nil {
				return nil
			}
		case *float64:
			if v != nil {
				return nil
			}
		default:
			return fmt.Errorf("validation failure: unknown type %T", value)
		}
		pattern := strings.Join(patterns, "|")
		var err error
		var match bool
		if typeName != "" {
			match, err = regexp.MatchString(pattern, typeName)
			if match {
				return fmt.Errorf("is required if type is '%s'", typeName)
			}
		}
		return err
	}
}

func getTypeName(db *sqlx.DB, typeID int) (string, bool, error) {
	name := ""
	if err := db.QueryRow(`SELECT name from type where id=$1`, typeID).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying type name: " + err.Error())
	}
	return name, true, nil
}

func CreateV12(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}

		ds := tc.DeliveryServiceNullableV12{}
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}

		if ds.RoutingName == nil || *ds.RoutingName == "" {
			ds.RoutingName = utils.StrPtr("cdn")
		}

		if errs := validateV12(db, &ds); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}

		dsv13 := tc.DeliveryServiceNullableV13{DeliveryServiceNullableV12: tc.DeliveryServiceNullableV12(ds)}

		if authorized, err := isTenantAuthorized(*user, db, &ds); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
			return
		} else if !authorized {
			api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
			return
		}

		dsv13, errCode, userErr, sysErr := create(db.DB, cfg, user, dsv13)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		api.WriteResp(w, r, []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
	}
}

func (ds *TODeliveryServiceV12) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	returnable := []interface{}{}
	dses, errs, errType := readGetDeliveryServices(params, db, user)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` {
				return nil, []error{errors.New("Resource not found.")}, tc.DataMissingError //matches perl response
			}
		}
		return nil, errs, errType
	}

	for _, ds := range dses {
		returnable = append(returnable, ds.DeliveryServiceNullableV12)
	}
	return returnable, nil, tc.NoError
}

func (ds *TODeliveryServiceV12) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	v13 := &TODeliveryServiceV13{
		Cfg: ds.Cfg,
		DB:  ds.DB,
		DeliveryServiceNullableV13: tc.DeliveryServiceNullableV13{
			DeliveryServiceNullableV12: ds.DeliveryServiceNullableV12,
		},
	}
	err, errType := v13.Delete(db, user)
	ds.DeliveryServiceNullableV12 = v13.DeliveryServiceNullableV12 // TODO avoid copy
	return err, errType
}

func UpdateV12(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		user, err := auth.GetCurrentUser(r.Context())
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting user: "+err.Error()))
			return
		}

		ds := tc.DeliveryServiceNullableV12{}
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
			return
		}

		if errs := validateV12(db, &ds); len(errs) > 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("invalid request: "+util.JoinErrs(errs).Error()), nil)
			return
		}

		dsv13 := tc.DeliveryServiceNullableV13{DeliveryServiceNullableV12: tc.DeliveryServiceNullableV12(ds)}

		if authorized, err := isTenantAuthorized(*user, db, &ds); err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
			return
		} else if !authorized {
			api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
			return
		}

		dsv13, errCode, userErr, sysErr := update(db.DB, cfg, *user, &dsv13)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		api.WriteResp(w, r, []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
	}
}
