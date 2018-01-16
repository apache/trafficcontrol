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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	tcapi "github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TODeliveryService tcapi.DeliveryService

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TODeliveryService(tcapi.DeliveryService{})

func GetRefType() *TODeliveryService {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (ds *TODeliveryService) GetID() int {
	return ds.ID
}

func (ds *TODeliveryService) GetAuditName() string {
	xmlId := *ds.XMLID
	return xmlId
}

func (ds *TODeliveryService) GetType() string {
	return "ds"
}

func (ds *TODeliveryService) SetID(i int) {
	ds.ID = i
}

func (ds *TODeliveryService) Validate(db *sqlx.DB) []error {

	noSpaces := validation.Match(regexp.MustCompile("^\\S*$"))
	noSpaces.Error("cannot contain spaces")

	noPeriods := validation.Match(regexp.MustCompile("^[^\\.]*$"))
	noPeriods.Error("cannot contain periods")

	var typeName string
	var err error
	if db != nil && ds.TypeID != nil {
		typeID := *ds.TypeID
		typeName, err = getTypeName(db, typeID)
		if err != nil {
			return []error{err}
		}
	}

	DNSRegexType := "^DNS.*$"
	HTTPRegexType := "^HTTP.*$"
	SteeringRegexType := "^STEERING.*$"
	// Custom Examples:
	// Just add isCIDR as a parameter to Validate()
	// isCIDR := validation.NewStringRule(govalidator.IsCIDR, "must be a valid CIDR address")
	isHost := validation.NewStringRule(govalidator.IsHost, "must be a valid hostname")
	errs := validation.Errors{
		"active": validation.Validate(ds.Active, validation.NotNil),
		"cdnId":  validation.Validate(ds.CDNID, validation.NotNil),
		"displayName": validation.Validate(ds.DisplayName,
			validation.Required),
		"dnsBypassIp":  validation.Validate(ds.DNSBypassIP, is.IP),
		"dnsBypassIp6": validation.Validate(ds.DNSBypassIP6, is.IPv6),
		"dscp":         validation.Validate(ds.DSCP, validation.NotNil),
		"geoLimit":     validation.Validate(ds.GeoLimit, validation.NotNil),
		"geoProvider":  validation.Validate(ds.GeoProvider, validation.NotNil),
		"infoUrl":      validation.Validate(ds.InfoURL, is.URL),
		"initialDispersion": validation.Validate(ds.InitialDispersion,
			validation.By(tovalidate.GreaterThanZero),
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"ipv6RoutingEnabled": validation.Validate(ds.IPV6RoutingEnabled,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"logsEnabled": validation.Validate(ds.LogsEnabled, validation.NotNil),
		"missLat": validation.Validate(ds.MissLat,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"missLong": validation.Validate(ds.MissLong,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"multiSiteOrigin": validation.Validate(ds.MultiSiteOrigin,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"orgServerFqdn": validation.Validate(ds.OrgServerFQDN,
			is.URL,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"protocol": validation.Validate(ds.Protocol,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"qstringIgnore": validation.Validate(ds.QStringIgnore,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"rangeRequestHandling": validation.Validate(ds.RangeRequestHandling,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"regionalGeoBlocking": validation.Validate(ds.RegionalGeoBlocking,
			validation.NotNil),
		"routingName": validation.Validate(ds.RoutingName,
			isHost,
			noPeriods,
			validation.Length(1, 48)),
		"typeId": validation.Validate(ds.TypeID,
			validation.NotNil,
			validation.By(tovalidate.GreaterThanZero)),
		"xmlId": validation.Validate(ds.XMLID,
			validation.Required,
			noSpaces,
			validation.Length(1, 48)),
	}
	return tovalidate.ToErrors(errs)
}

func requiredIfMatchesTypeName(patterns []string, typeName string) func(interface{}) error {
	return func(value interface{}) error {

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

// TODO: drichardson - refactor to the types.go once implemented.
func getTypeName(db *sqlx.DB, typeID int) (string, error) {

	query := `SELECT name from type where id=$1`

	var rows *sqlx.Rows
	var err error

	rows, err = db.Queryx(query, typeID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	typeResults := []tc.Type{}
	for rows.Next() {
		var s tc.Type
		if err = rows.StructScan(&s); err != nil {
			return "", fmt.Errorf("getting Type: %v", err)
		}
		typeResults = append(typeResults, s)
	}

	typeName := typeResults[0].Name
	return typeName, err
}

//The TODeliveryService implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a delivery service with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (ds *TODeliveryService) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with ds: %++v", updateDSQuery(), ds)
	resultRows, err := tx.NamedQuery(updateDSQuery(), ds)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(err)
			if eType == tc.DataConflictError {
				return errors.New("a delivery service with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	var lastUpdated tc.Time
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	ds.LastUpdated = lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no delivery service found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//The TODeliveryService implementation of the Inserter interface
//all implementations of Inserter should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a ds with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted ds and have
//to be added to the struct
func (ds *TODeliveryService) Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	fmt.Printf("ds ---> %v\n", ds)
	resultRows, err := tx.NamedQuery(insertDSQuery(), ds)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqerr)
			return errors.New("a delivery service with " + err.Error()), eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	var id int
	var lastUpdated tc.Time
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no delivery service was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from delivery service insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	ds.SetID(id)
	ds.LastUpdated = lastUpdated
	return nil, tc.NoError
}

//The DeliveryService implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (ds *TODeliveryService) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with Delivery Service: %++v", deleteDSQuery(), ds)
	result, err := tx.NamedExec(deleteDSQuery(), ds)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no delivery service with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//TODO: drichardson - plumb these out!
func updateDSQuery() string {
	query := `UPDATE
cdn SET
dnssec_enabled=:dnssec_enabled,
domain_name=:domain_name,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

//TODO: drichardson - plumb these out!
func insertDSQuery() string {
	query := `INSERT INTO deliveryservice (
dnssec_enabled,
domain_name,
name) VALUES (
:dnssec_enabled,
:domain_name,
:name) RETURNING id,last_updated`
	return query
}

func deleteDSQuery() string {
	query := `DELETE FROM deliveryservice
WHERE id=:id`
	return query
}
