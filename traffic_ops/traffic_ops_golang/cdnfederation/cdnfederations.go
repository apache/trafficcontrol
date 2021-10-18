package cdnfederation

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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation"
)

// we need a type alias to define functions on
type TOCDNFederation struct {
	api.APIInfoImpl `json:"-"`
	tc.CDNFederation
	TenantID *int `json:"-" db:"tenant_id"`
}

func (v *TOCDNFederation) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "federation")
}

func (v *TOCDNFederation) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOCDNFederation) InsertQuery() string           { return insertQuery() }
func (v *TOCDNFederation) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(federation.last_updated) as t from federation 
		join federation_deliveryservice fds on fds.federation = federation.id 
		join deliveryservice ds on ds.id = fds.deliveryservice 
		join cdn c on c.id = ds.cdn_id ` + where + orderBy + pagination +
		` UNION ALL
		select max(last_updated) as t from last_deleted l where l.table_name='federation') as res`
}

func (v *TOCDNFederation) NewReadObj() interface{} { return &TOCDNFederation{} }
func (v *TOCDNFederation) SelectQuery() string {
	return selectByCDNName()
}
func (v *TOCDNFederation) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	cols := map[string]dbhelpers.WhereColumnInfo{
		"id":    dbhelpers.WhereColumnInfo{Column: "federation.id", Checker: api.IsInt},
		"name":  dbhelpers.WhereColumnInfo{Column: "c.name", Checker: nil},
		"cname": dbhelpers.WhereColumnInfo{Column: "federation.cname", Checker: nil},
	}
	return cols
}
func (v *TOCDNFederation) DeleteQuery() string { return deleteQuery() }
func (v *TOCDNFederation) UpdateQuery() string { return updateQuery() }

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetKeys() (map[string]interface{}, bool) {
	if fed.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *fed.ID}, true
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetAuditName() string {
	if fed.CName != nil {
		return *fed.CName
	}
	if fed.ID != nil {
		return strconv.Itoa(*fed.ID)
	}
	return "unknown"
}

// Fufills `Identifier' interface
func (fed TOCDNFederation) GetType() string {
	return "cdnfederation"
}

// Fufills `Create' interface
func (fed *TOCDNFederation) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) // non-panicking type assertion
	fed.ID = &i
}

// Fulfills `Validate' interface
func (fed *TOCDNFederation) Validate() error {

	isDNSName := validation.NewStringRule(govalidator.IsDNSName, "must be a valid hostname")
	endsWithDot := validation.NewStringRule(
		func(str string) bool {
			return strings.HasSuffix(str, ".")
		}, "must end with a period")

	// cname regex: (^\S*\.$), ttl regex: (^\d+$)
	validateErrs := validation.Errors{
		"cname": validation.Validate(fed.CName, validation.Required, endsWithDot, isDNSName),
		"ttl":   validation.Validate(fed.TTL, validation.Required, validation.Min(0)),
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

func (fed *TOCDNFederation) CheckIfCDNAndFederationMatch(cdnName string) (error, error, int) {
	var cdnFromDS string
	var err error
	if fed.DeliveryServiceIDs != nil {
		if fed.DsId != nil {
			cdnNames, err := dbhelpers.GetCDNNamesFromDSIds(fed.APIInfo().Tx.Tx, []int{*fed.DsId})
			if err != nil {
				return nil, fmt.Errorf("getting CDN names from DS IDs: %w", err), http.StatusInternalServerError
			}
			if len(cdnNames) != 1 {
				return fmt.Errorf("%d CDNs returned for one DS ID", len(cdnNames)), nil, http.StatusBadRequest
			}
			cdnFromDS = cdnNames[0]
		} else if fed.XmlId != nil {
			cdnFromDS, err = dbhelpers.GetCDNNameFromDSXMLID(fed.APIInfo().Tx.Tx, *fed.XmlId)
			if err != nil {
				return nil, fmt.Errorf("getting CDN name from DS XMLID: %w", err), http.StatusInternalServerError
			}
		}
	}
	if cdnFromDS != "" && cdnFromDS != cdnName {
		return errors.New("cdn names in request path and payload do not match"), nil, http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

// fedAPIInfo.Params["name"] is not used on creation, rather the cdn name
// is connected when the federations/:id/deliveryservice links a federation
// However, we use fedAPIInfo.Params["name"] to check whether or not another user has a hard lock on the CDN.
// Note: cdns and deliveryservies have a 1-1 relationship
func (fed *TOCDNFederation) Create() (error, error, int) {
	if cdn, ok := fed.APIInfo().Params["name"]; ok {
		if ok, err := dbhelpers.CDNExists(fed.APIInfo().Params["name"], fed.APIInfo().Tx.Tx); err != nil {
			return nil, errors.New("verifying CDN exists: " + err.Error()), http.StatusInternalServerError
		} else if !ok {
			return errors.New("cdn not found"), nil, http.StatusNotFound
		}
		userErr, sysErr, errCode := fed.CheckIfCDNAndFederationMatch(cdn)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(fed.APIInfo().Tx.Tx, cdn, fed.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	// Deliveryservice IDs should not be included on create.
	if fed.DeliveryServiceIDs != nil {
		fed.DsId = nil
		fed.XmlId = nil
		fed.DeliveryServiceIDs = nil
	}
	return api.GenericCreate(fed)
}

// returning true indicates the data related to the given tenantID should be visible
// `tenantIDs` is presumed to be unsorted, and a nil tenantID is viewable by everyone
func checkTenancy(tenantID *int, tenantIDs []int) bool {
	if tenantID == nil {
		return true
	}
	for _, id := range tenantIDs {
		if id == *tenantID {
			return true
		}
	}
	return false
}

func (fed *TOCDNFederation) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	if idstr, ok := fed.APIInfo().Params["id"]; ok {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return nil, errors.New("id must be an integer"), nil, http.StatusBadRequest, nil
		}
		fed.ID = util.IntPtr(id)
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(fed.APIInfo().Tx.Tx, fed.APIInfo().User.TenantID)
	if err != nil {
		return nil, nil, errors.New("getting tenant list for user: " + err.Error()), http.StatusInternalServerError, nil
	}

	api.DefaultSort(fed.APIInfo(), "cname")
	if ok, err := dbhelpers.CDNExists(fed.APIInfo().Params["name"], fed.APIInfo().Tx.Tx); err != nil {
		return nil, nil, errors.New("verifying CDN exists: " + err.Error()), http.StatusInternalServerError, nil
	} else if !ok {
		return nil, errors.New("cdn not found"), nil, http.StatusNotFound, nil
	}
	federations, userErr, sysErr, errCode, maxTime := api.GenericRead(h, fed, useIMS)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	if errCode == http.StatusNotModified {
		return []interface{}{}, nil, nil, http.StatusNotModified, maxTime
	}

	filteredFederations := []interface{}{}
	for _, ifederation := range federations {
		federation := ifederation.(*TOCDNFederation)
		if !checkTenancy(federation.TenantID, tenantIDs) {
			return nil, errors.New("user not authorized for requested federation"), nil, http.StatusForbidden, nil
		}
		filteredFederations = append(filteredFederations, federation.CDNFederation)
	}

	if len(filteredFederations) == 0 {
		if fed.ID != nil {
			return nil, errors.New("not found"), nil, http.StatusNotFound, nil
		}
	}
	return filteredFederations, nil, nil, errCode, maxTime
}

func (fed *TOCDNFederation) Update(h http.Header) (error, error, int) {
	userErr, sysErr, errCode := fed.isTenantAuthorized()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if cdn, ok := fed.APIInfo().Params["name"]; ok {
		if ok, err := dbhelpers.CDNExists(fed.APIInfo().Params["name"], fed.APIInfo().Tx.Tx); err != nil {
			return nil, errors.New("verifying CDN exists: " + err.Error()), http.StatusInternalServerError
		} else if !ok {
			return errors.New("cdn not found"), nil, http.StatusNotFound
		}
		userErr, sysErr, errCode := fed.CheckIfCDNAndFederationMatch(cdn)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(fed.APIInfo().Tx.Tx, cdn, fed.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	// Deliveryservice IDs should not be included on update.
	if fed.DeliveryServiceIDs != nil {
		fed.DsId = nil
		fed.XmlId = nil
		fed.DeliveryServiceIDs = nil
	}
	return api.GenericUpdate(h, fed)
}

// Delete implements the Deleter interface for TOCDNFederation.
func (fed *TOCDNFederation) Delete() (error, error, int) {
	userErr, sysErr, errCode := fed.isTenantAuthorized()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if cdn, ok := fed.APIInfo().Params["name"]; ok {
		if ok, err := dbhelpers.CDNExists(fed.APIInfo().Params["name"], fed.APIInfo().Tx.Tx); err != nil {
			return nil, errors.New("verifying CDN exists: " + err.Error()), http.StatusInternalServerError
		} else if !ok {
			return errors.New("cdn not found"), nil, http.StatusNotFound
		}
		userErr, sysErr, errCode := fed.CheckIfCDNAndFederationMatch(cdn)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(fed.APIInfo().Tx.Tx, cdn, fed.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericDelete(fed)
}

func (fed TOCDNFederation) isTenantAuthorized() (error, error, int) {
	tenantID, err := getTenantIDFromFedID(*fed.ID, fed.APIInfo().Tx.Tx)
	if err != nil {
		// If nobody has claimed a tenant, that federation is publicly visible.
		// This logically follows /federations/:id/deliveryservices
		if err == sql.ErrNoRows {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("getting tenant id from federation: " + err.Error()), http.StatusInternalServerError
	}

	// TODO: use IsResourceAuthorizedToUserTx instead
	list, err := tenant.GetUserTenantIDListTx(fed.APIInfo().Tx.Tx, fed.APIInfo().User.TenantID)
	if err != nil {
		return nil, errors.New("getting federation tenant list: " + err.Error()), http.StatusInternalServerError
	}
	for _, id := range list {
		if id == tenantID {
			return nil, nil, http.StatusOK
		}
	}
	return errors.New("unauthorized for tenant"), nil, http.StatusForbidden
}

func getTenantIDFromFedID(id int, tx *sql.Tx) (int, error) {
	tenantID := 0
	query := `
	SELECT ds.tenant_id FROM federation AS f
	JOIN federation_deliveryservice AS fd ON f.id = fd.federation
	JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice
	WHERE f.id = $1`
	err := tx.QueryRow(query, id).Scan(&tenantID)
	return tenantID, err
}

func selectByID() string {
	return `
	SELECT
	ds.tenant_id,
	federation.id AS id,
	federation.cname,
	federation.ttl,
	federation.description,
	federation.last_updated,
	ds.id AS ds_id,
	ds.xml_id
	FROM federation
	LEFT JOIN federation_deliveryservice AS fd ON federation.id = fd.federation
	LEFT JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice`
	// WHERE federation.id = :id (determined by dbhelper)
}

func selectByCDNName() string {
	return `
	SELECT
	ds.tenant_id,
	federation.id AS id,
	federation.cname,
	federation.ttl,
	federation.description,
	federation.last_updated,
	ds.id AS ds_id,
	ds.xml_id
	FROM federation
	JOIN federation_deliveryservice AS fd ON federation.id = fd.federation
	JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice
	JOIN cdn c ON c.id = ds.cdn_id`
	// WHERE cdn.name = :cdn_name (determined by dbhelper)
}

func updateQuery() string {
	return `
UPDATE federation SET
	cname = :cname,
	ttl = :ttl,
	description = :description
WHERE
  id=:id
RETURNING last_updated`
}

func insertQuery() string {
	return `
	INSERT INTO federation (
	cname,
 	ttl,
 	description
  ) VALUES (
 	:cname,
	:ttl,
	:description
	) RETURNING id, last_updated`
}

func deleteQuery() string {
	return `DELETE FROM federation WHERE id = :id`
}
