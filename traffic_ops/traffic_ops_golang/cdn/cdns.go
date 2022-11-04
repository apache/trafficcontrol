package cdn

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

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
)

// we need a type alias to define functions on
type TOCDN struct {
	api.APIInfoImpl `json:"-"`
	tc.CDNNullable
}

func (v *TOCDN) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "cdn")
}

func (v *TOCDN) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` c ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (v *TOCDN) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOCDN) InsertQuery() string           { return insertQuery(v.APIInfo().Version) }
func (v *TOCDN) NewReadObj() interface{}       { return &tc.CDNNullable{} }
func (v *TOCDN) SelectQuery() string           { return selectQuery(v.APIInfo().Version) }
func (v *TOCDN) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	columnInfo := map[string]dbhelpers.WhereColumnInfo{
		"domainName":    dbhelpers.WhereColumnInfo{Column: "domain_name"},
		"dnssecEnabled": dbhelpers.WhereColumnInfo{Column: "dnssec_enabled"},
		"id":            dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"name":          dbhelpers.WhereColumnInfo{Column: "name"},
	}
	if v.APIInfo().Version.GreaterThanOrEqualTo(&api.Version{Major: 4, Minor: 1}) {
		columnInfo["ttlOverride"] = dbhelpers.WhereColumnInfo{Column: "ttl_override", Checker: api.IsInt}
	}
	return columnInfo
}
func (v *TOCDN) UpdateQuery() string { return updateQuery(v.APIInfo().Version) }
func (v *TOCDN) DeleteQuery() string { return deleteQuery() }

func (cdn TOCDN) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (cdn TOCDN) GetKeys() (map[string]interface{}, bool) {
	if cdn.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cdn.ID}, true
}

func (cdn TOCDN) GetAuditName() string {
	if cdn.Name != nil {
		return *cdn.Name
	}
	if cdn.ID != nil {
		return strconv.Itoa(*cdn.ID)
	}
	return "0"
}

func (cdn TOCDN) GetType() string {
	return "cdn"
}

func (cdn *TOCDN) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cdn.ID = &i
}

func isValidCDNchar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	if r == '.' || r == '-' {
		return true
	}
	return false
}

// IsValidCDNName returns true if the name contains only characters valid for a CDN name
func IsValidCDNName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCDNchar(r) })
	return i == -1
}

// Validate fulfills the api.Validator interface.
func (cdn TOCDN) Validate() (error, error) {
	validName := validation.NewStringRule(IsValidCDNName, "invalid characters found - Use alphanumeric . or - .")
	validDomainName := validation.NewStringRule(govalidator.IsDNSName, "not a valid domain name")
	errs := validation.Errors{
		"name":       validation.Validate(cdn.Name, validation.Required, validName),
		"domainName": validation.Validate(cdn.DomainName, validation.Required, validDomainName),
	}
	if cdn.APIInfo().Version.GreaterThanOrEqualTo(&api.Version{Major: 4, Minor: 1}) {
		errs["ttlOverride"] = validation.Validate(cdn.TTLOverride, validation.By(tovalidate.IsGreaterThanZero))
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (cdn *TOCDN) Create() (error, error, int) {
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	if cdn.APIInfo().Version.LessThan(&api.Version{Major: 4, Minor: 1}) {
		cdn.TTLOverride = nil
	}
	return api.GenericCreate(cdn)
}

func (cdn *TOCDN) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(cdn.APIInfo(), "name")
	return api.GenericRead(h, cdn, useIMS)
}

func (cdn *TOCDN) Update(h http.Header) (error, error, int) {
	if cdn.ID != nil {
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(cdn.APIInfo().Tx.Tx, int64(*cdn.ID), cdn.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	if cdn.APIInfo().Version.LessThan(&api.Version{Major: 4, Minor: 1}) {
		cdn.TTLOverride = nil
	}
	return api.GenericUpdate(h, cdn)
}

func (cdn *TOCDN) Delete() (error, error, int) {
	if cdn.ID != nil {
		userErr, sysErr, errCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(cdn.APIInfo().Tx.Tx, int64(*cdn.ID), cdn.APIInfo().User.UserName)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return api.GenericDelete(cdn)
}

func formatQueryByAPIVersion(apiVersion *api.Version, minimumAPIVersion *api.Version, queryFormatString string, columnStrs []string, lowAPIVersionColumnStrs []string) string {
	if apiVersion.LessThan(&api.Version{Major: 4, Minor: 1}) {
		for index, _ := range columnStrs {
			columnStrs[index] = lowAPIVersionColumnStrs[index]
		}
	}
	columnStrArgs := make([]interface{}, len(columnStrs))
	for index, _ := range columnStrs {
		columnStrArgs[index] = columnStrs[index]
	}
	query := fmt.Sprintf(queryFormatString, columnStrArgs...)
	return query
}

func selectQuery(apiVersion *api.Version) string {
	query := `SELECT
dnssec_enabled,
domain_name,
id,
last_updated,
%s
name

FROM cdn c`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`
			ttl_override,
`}, []string{``})
}

func updateQuery(apiVersion *api.Version) string {
	query := `UPDATE
cdn SET
dnssec_enabled=:dnssec_enabled,
domain_name=:domain_name,
name=:name
%s
WHERE id=:id RETURNING last_updated`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`,
ttl_override=:ttl_override
`}, []string{``})
}

func insertQuery(apiVersion *api.Version) string {
	query := `INSERT INTO cdn (
dnssec_enabled,
domain_name,
name
%s
) VALUES (
:dnssec_enabled,
:domain_name,
:name
%s
) RETURNING id,last_updated`
	return formatQueryByAPIVersion(apiVersion, &api.Version{Major: 4, Minor: 1}, query, []string{`,
ttl_override
`, `,
:ttl_override
`}, []string{``, ``})
}

func deleteQuery() string {
	query := `DELETE FROM cdn
WHERE id=:id`
	return query
}
