package dbhelpers

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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/topology/topology_validation"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WhereColumnInfo struct {
	Column  string
	Checker func(string) error
}

const BaseWhere = "\nWHERE"
const BaseOrderBy = "\nORDER BY"
const BaseLimit = "\nLIMIT"
const BaseOffset = "\nOFFSET"

const getDSTenantIDFromXMLIDQuery = `
SELECT deliveryservice.tenant_id
FROM deliveryservice
WHERE deliveryservice.xml_id = $1
`

const getFederationIDForUserIDByXMLIDQuery = `
SELECT federation_deliveryservice.federation
FROM federation_deliveryservice
WHERE federation_deliveryservice.deliveryservice IN (
	SELECT deliveryservice.id
	FROM deliveryservice
	WHERE deliveryservice.xml_id = $1
) AND federation_deliveryservice.federation IN (
	SELECT federation_tmuser.federation
	FROM federation_tmuser
	WHERE federation_tmuser.tm_user = $2
)
`

const getUserBaseQuery = `
SELECT tm_user.address_line1,
       tm_user.address_line2,
       tm_user.city,
       tm_user.company,
       tm_user.country,
       tm_user.email,
       tm_user.full_name,
       tm_user.gid,
       tm_user.id,
       tm_user.last_updated,
       tm_user.new_user,
       tm_user.phone_number,
       tm_user.postal_code,
       tm_user.public_ssh_key,
       tm_user.registration_sent,
       tm_user.role,
       role.name AS role_name,
       tm_user.state_or_province,
       tenant.name AS tenant,
       tm_user.tenant_id,
       tm_user.token,
       tm_user.uid,
       tm_user.username
FROM tm_user
LEFT OUTER JOIN role ON role.id = tm_user.role
LEFT OUTER JOIN tenant ON tenant.id = tm_user.tenant_id
`
const getUserByIDQuery = getUserBaseQuery + `
WHERE tm_user.id = $1
`

const getUserByEmailQuery = getUserBaseQuery + `
WHERE tm_user.email = $1
`

// CheckIfCurrentUserHasCdnLock checks if the current user has the lock on the cdn that the requested operation is to be performed on.
// This will succeed if the either there is no lock by any user on the CDN, or if the current user has the lock on the CDN.
func CheckIfCurrentUserHasCdnLock(tx *sql.Tx, cdn, user string) (error, error, int) {
	query := `
SELECT c.username, ARRAY_REMOVE(ARRAY_AGG(u.username), NULL) AS shared_usernames
FROM cdn_lock c
    LEFT JOIN cdn_lock_user u ON c.username = u.owner AND c.cdn = u.cdn
WHERE c.cdn=$1
GROUP BY c.username`
	var userName string
	var sharedUserNames []string
	rows, err := tx.Query(query, cdn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("querying cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userName, pq.Array(&sharedUserNames))
		if err != nil {
			return nil, errors.New("scanning cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
		}
	}
	if userName != "" && user != userName {
		for _, u := range sharedUserNames {
			if u == user {
				return nil, nil, http.StatusOK
			}
		}
		return errors.New("user " + user + " currently does not have the lock on cdn " + cdn), nil, http.StatusForbidden
	}
	return nil, nil, http.StatusOK
}

func BuildWhereAndOrderByAndPagination(parameters map[string]string, queryParamsToSQLCols map[string]WhereColumnInfo) (string, string, string, map[string]interface{}, []error) {
	whereClause := BaseWhere
	orderBy := BaseOrderBy
	paginationClause := BaseLimit
	var criteria string
	var queryValues map[string]interface{}
	var errs []error
	criteria, queryValues, errs = parseCriteriaAndQueryValues(queryParamsToSQLCols, parameters)

	if len(queryValues) > 0 {
		whereClause += " " + criteria
	}
	if len(errs) > 0 {
		return "", "", "", queryValues, errs
	}

	if orderby, ok := parameters["orderby"]; ok {
		log.Debugln("orderby: ", orderby)
		if colInfo, ok := queryParamsToSQLCols[orderby]; ok {
			log.Debugln("orderby column ", colInfo)
			orderBy += " " + colInfo.Column

			// if orderby is specified and valid, also check for sortOrder
			if sortOrder, exists := parameters["sortOrder"]; exists {
				log.Debugln("sortOrder: ", sortOrder)
				if sortOrder == "desc" {
					orderBy += " DESC"
				} else if sortOrder != "asc" {
					log.Debugln("sortOrder value must be desc or asc. Invalid value provided: ", sortOrder)
				}
			}
		} else {
			log.Debugln("This column is not configured to support orderby: ", orderby)
		}
	}

	if limit, exists := parameters["limit"]; exists {
		// try to convert to int, if it fails the limit parameter is invalid, so return an error
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt < -1 {
			errs = append(errs, errors.New("limit parameter must be bigger than -1"))
			return "", "", "", queryValues, errs
		}
		log.Debugln("limit: ", limit)
		if limitInt == -1 {
			paginationClause = ""
		} else {
			paginationClause += " " + limit
		}
		if offset, exists := parameters["offset"]; exists {
			// check that offset is valid
			offsetInt, err := strconv.Atoi(offset)
			if err != nil || offsetInt < 1 {
				errs = append(errs, errors.New("offset parameter must be a positive integer"))
				return "", "", "", queryValues, errs
			}
			paginationClause += BaseOffset + " " + offset
		} else if page, exists := parameters["page"]; exists {
			// check that offset is valid
			page, err := strconv.Atoi(page)
			if err != nil || page < 1 {
				errs = append(errs, errors.New("page parameter must be a positive integer"))
				return "", "", "", queryValues, errs
			}
			paginationClause += BaseOffset + " " + strconv.Itoa((page-1)*limitInt)
		}
	}

	if whereClause == BaseWhere {
		whereClause = ""
	}
	if orderBy == BaseOrderBy {
		orderBy = ""
	}
	if paginationClause == BaseLimit {
		paginationClause = ""
	}
	log.Debugf("\n--\n Where: %s \n Order By: %s \n Limit+Offset: %s", whereClause, orderBy, paginationClause)
	return whereClause, orderBy, paginationClause, queryValues, errs
}

// CheckIfCurrentUserCanModifyCDNs checks if the current user has the lock on the list of cdns that the requested operation is to be performed on.
// This will succeed if the either there is no lock by any user on any of the CDNs, or if the current user has the lock on any of the CDNs.
func CheckIfCurrentUserCanModifyCDNs(tx *sql.Tx, cdns []string, user string) (error, error, int) {
	query := `SELECT c.username, c.soft, c.cdn, ARRAY_REMOVE(ARRAY_AGG(u.username), NULL) AS shared_usernames FROM cdn_lock c LEFT JOIN cdn_lock_user u ON c.username = u.owner AND c.cdn = u.cdn WHERE c.cdn=ANY($1) GROUP BY c.username, c.soft, c.cdn`
	var userName, cdn string
	var soft bool
	var sharedUserNames []string
	rows, err := tx.Query(query, pq.Array(cdns))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("querying cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userName, &soft, &cdn, pq.Array(&sharedUserNames))
		if err != nil {
			return nil, errors.New("scanning cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
		}
		if userName != "" && user != userName && !soft {
			for _, u := range sharedUserNames {
				if u == user {
					return nil, nil, http.StatusOK
				}
			}
			return errors.New("user " + userName + " currently has a hard lock on cdn " + cdn), nil, http.StatusForbidden
		}
	}
	return nil, nil, http.StatusOK
}

// CheckIfCurrentUserCanModifyCDNs checks if the current user has the lock on the list of cdns(identified by ID) that the requested operation is to be performed on.
// This will succeed if the either there is no lock by any user on any of the CDNs, or if the current user has the lock on any of the CDNs.
func CheckIfCurrentUserCanModifyCDNsByID(tx *sql.Tx, cdns []int, user string) (error, error, int) {
	query := `SELECT name FROM cdn WHERE id=ANY($1)`
	var name string
	var cdnNames []string
	rows, err := tx.Query(query, pq.Array(cdns))
	if err != nil {
		return nil, errors.New("no cdn names found for the given IDs: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, errors.New("scanning cdn name: " + err.Error()), http.StatusInternalServerError
		}
		cdnNames = append(cdnNames, name)
	}
	return CheckIfCurrentUserCanModifyCDNs(tx, cdnNames, user)
}

// CheckIfCurrentUserCanModifyCDN checks if the current user has the lock on the cdn that the requested operation is to be performed on.
// This will succeed if the either there is no lock by any user on the CDN, or if the current user has the lock on the CDN.
func CheckIfCurrentUserCanModifyCDN(tx *sql.Tx, cdn, user string) (error, error, int) {
	query := `SELECT c.username, c.soft, ARRAY_REMOVE(ARRAY_AGG(u.username), NULL) AS shared_usernames FROM cdn_lock c LEFT JOIN cdn_lock_user u ON c.username = u.owner AND c.cdn = u.cdn WHERE c.cdn=$1 GROUP BY c.username, c.soft`
	var userName string
	var soft bool
	var sharedUserNames []string
	rows, err := tx.Query(query, cdn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("querying cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userName, &soft, pq.Array(&sharedUserNames))
		if err != nil {
			return nil, errors.New("scanning cdn_lock for user " + user + " and cdn " + cdn + ": " + err.Error()), http.StatusInternalServerError
		}
		if userName != "" && user != userName && !soft {
			for _, u := range sharedUserNames {
				if u == user {
					return nil, nil, http.StatusOK
				}
			}
			return errors.New("user " + userName + " currently has a hard lock on cdn " + cdn), nil, http.StatusForbidden
		}
	}
	return nil, nil, http.StatusOK
}

// CheckIfCurrentUserCanModifyCDNWithID checks if the current user has the lock on the cdn (identified by ID) that the requested operation is to be performed on.
// This will succeed if the either there is no lock by any user on the CDN, or if the current user has the lock on the CDN.
func CheckIfCurrentUserCanModifyCDNWithID(tx *sql.Tx, cdnID int64, user string) (error, error, int) {
	cdnName, ok, err := GetCDNNameFromID(tx, cdnID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	} else if !ok {
		return errors.New("CDN not found"), nil, http.StatusNotFound
	}
	return CheckIfCurrentUserCanModifyCDN(tx, string(cdnName), user)
}

// CheckIfCurrentUserCanModifyCachegroup checks if the current user has the lock on the cdns that are associated with the provided cachegroup ID.
// This will succeed if no other user has a hard lock on any of the CDNs that relate to the cachegroup in question.
func CheckIfCurrentUserCanModifyCachegroup(tx *sql.Tx, cachegroupID int, user string) (error, error, int) {
	query := `
SELECT c.username, c.cdn, c.soft, ARRAY_REMOVE(ARRAY_AGG(u.username), NULL) AS shared_usernames
FROM cdn_lock c LEFT JOIN cdn_lock_user u
    ON c.username = u.owner
           AND c.cdn = u.cdn
WHERE c.cdn IN (
    SELECT name FROM cdn
    WHERE id IN (
        SELECT cdn_id FROM server
        WHERE cachegroup = ($1)))
GROUP BY c.username, c.cdn, c.soft`
	var userName string
	var cdn string
	var soft bool
	var sharedUserNames []string
	rows, err := tx.Query(query, cachegroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("querying cdn_lock for user " + user + " and cachegroup ID " + strconv.Itoa(cachegroupID) + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userName, &cdn, &soft, pq.Array(&sharedUserNames))
		if err != nil {
			return nil, errors.New("scanning cdn_lock for user " + user + " and cachegroup ID " + strconv.Itoa(cachegroupID) + ": " + err.Error()), http.StatusInternalServerError
		}
		if userName != "" && user != userName && !soft {
			for _, u := range sharedUserNames {
				if u == user {
					return nil, nil, http.StatusOK
				}
			}
			return errors.New("user " + userName + " currently has a hard lock on cdn " + cdn), nil, http.StatusForbidden
		}
	}
	return nil, nil, http.StatusOK
}

// CheckIfCurrentUserCanModifyCachegroups checks if the current user has the lock on the cdns that are associated with the provided cachegroup IDs.
// This will succeed if no other user has a hard lock on any of the CDNs that relate to the cachegroups in question.
func CheckIfCurrentUserCanModifyCachegroups(tx *sql.Tx, cachegroupIDs []int, user string) (error, error, int) {
	query := `SELECT c.username, c.cdn, c.soft, ARRAY_REMOVE(ARRAY_AGG(u.username), NULL) AS shared_usernames FROM cdn_lock c
    LEFT JOIN cdn_lock_user u
        ON c.username = u.owner
               AND c.cdn = u.cdn
WHERE c.cdn IN (
    SELECT name FROM cdn
    WHERE id IN (
        SELECT cdn_id FROM server
        WHERE cachegroup = ANY($1)))
        GROUP BY c.username, c.cdn, c.soft`
	var userName string
	var cdn string
	var soft bool
	var sharedUserNames []string
	rows, err := tx.Query(query, pq.Array(cachegroupIDs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, http.StatusOK
		}
		return nil, errors.New("querying cachegroups cdn_lock for user " + user + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userName, &cdn, &soft, pq.Array(&sharedUserNames))
		if err != nil {
			return nil, errors.New("scanning cachegroups cdn_lock for user " + user + ": " + err.Error()), http.StatusInternalServerError
		}
		if userName != "" && user != userName && !soft {
			for _, u := range sharedUserNames {
				if u == user {
					return nil, nil, http.StatusOK
				}
			}
			return errors.New("user " + userName + " currently has a hard lock on cdn " + cdn), nil, http.StatusForbidden
		}
	}
	return nil, nil, http.StatusOK
}

func parseCriteriaAndQueryValues(queryParamsToSQLCols map[string]WhereColumnInfo, parameters map[string]string) (string, map[string]interface{}, []error) {
	var criteria string

	var criteriaArgs []string
	errs := []error{}
	queryValues := make(map[string]interface{})
	for key, colInfo := range queryParamsToSQLCols {
		if urlValue, ok := parameters[key]; ok {
			var err error
			if colInfo.Checker != nil {
				err = colInfo.Checker(urlValue)
			}
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %w", key, err))
			} else {
				criteria = colInfo.Column + "=:" + key
				criteriaArgs = append(criteriaArgs, criteria)
				queryValues[key] = urlValue
			}
		}
	}
	criteria = strings.Join(criteriaArgs, " AND ")

	return criteria, queryValues, errs
}

// AddTenancyCheck takes a WHERE clause (can be ""), the associated queryValues (can be empty),
// a tenantColumnName that should provide a bigint corresponding to the tenantID of the object being checked (this may require a CAST),
// and an array of the tenantIDs the user has access to; it returns a where clause and associated queryValues including filtering based on tenancy.
func AddTenancyCheck(where string, queryValues map[string]interface{}, tenantColumnName string, tenantIDs []int) (string, map[string]interface{}) {
	if where == "" {
		where = BaseWhere + " " + tenantColumnName + " = ANY(CAST(:accessibleTenants AS bigint[]))"
	} else {
		where += " AND " + tenantColumnName + " = ANY(CAST(:accessibleTenants AS bigint[]))"
	}
	queryValues["accessibleTenants"] = pq.Array(tenantIDs)

	return where, queryValues
}

// CommitIf commits if doCommit is true at the time of execution.
// This is designed as a defer helper.
//
// Example:
//
//	tx, err := db.Begin()
//	txCommit := false
//	defer dbhelpers.CommitIf(tx, &txCommit)
//	if err := tx.Exec("select ..."); err != nil {
//	  return errors.New("executing: " + err.Error())
//	}
//	txCommit = true
//	return nil
func CommitIf(tx *sql.Tx, doCommit *bool) {
	if *doCommit {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}

// GetPrivLevelFromRoleID returns the priv_level associated with a role, whether it exists, and any error.
// This method exists on a temporary basis. After priv_level is fully deprecated and capabilities take over,
// this method will not only no longer be needed, but the corresponding new privilege check should be done
// via the primary database query for the users endpoint. The users json response will contain a list of
// capabilities in the future, whereas now the users json response currently does not contain privLevel.
// See the wiki page on the roles/capabilities as a system:
// https://cwiki.apache.org/confluence/pages/viewpage.action?pageId=68715910
func GetPrivLevelFromRoleID(tx *sql.Tx, id int) (int, bool, error) {
	var privLevel int
	err := tx.QueryRow(`SELECT priv_level FROM role WHERE role.id = $1`, id).Scan(&privLevel)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("getting priv_level from role ID: %w", err)
	}
	return privLevel, true, nil
}

// GetPrivLevelFromRole returns the priv_level associated with a role, whether it exists, and any error.
// This method exists on a temporary basis. After priv_level is fully deprecated and capabilities take over,
// this method will not only no longer be needed, but the corresponding new privilege check should be done
// via the primary database query for the users endpoint. The users json response will contain a list of
// capabilities in the future, whereas now the users json response currently does not contain privLevel.
// See the wiki page on the roles/capabilities as a system:
// https://cwiki.apache.org/confluence/pages/viewpage.action?pageId=68715910
func GetPrivLevelFromRole(tx *sql.Tx, role string) (int, bool, error) {
	var privLevel int
	err := tx.QueryRow(`SELECT priv_level FROM role WHERE role.name = $1`, role).Scan(&privLevel)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("getting priv_level from role: %w", err)
	}
	return privLevel, true, nil
}

// GetCapabilitiesFromRoleID returns the capabilities for the supplied role ID.
func GetCapabilitiesFromRoleID(tx *sql.Tx, roleID int) ([]string, error) {
	var caps []string
	var cap string

	rows, err := tx.Query(`SELECT cap_name FROM role_capability WHERE role_id = $1`, roleID)

	if err != nil {
		return caps, fmt.Errorf("getting capabilities from role ID: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&cap)
		if err != nil {
			return caps, fmt.Errorf("scanning capabilities: %w", err)
		}
		caps = append(caps, cap)
	}
	return caps, nil
}

// GetCapabilitiesFromRoleName returns the capabilities for the supplied role name.
func GetCapabilitiesFromRoleName(tx *sql.Tx, role string) ([]string, error) {
	var caps []string
	var cap string

	rows, err := tx.Query(`SELECT cap_name FROM role_capability rc JOIN role r ON r.id = rc.role_id WHERE r.name = $1`, role)

	if err != nil {
		return caps, fmt.Errorf("getting capabilities from role name: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&cap)
		if err != nil {
			return caps, fmt.Errorf("scanning capabilities: %w", err)
		}
		caps = append(caps, cap)
	}
	return caps, nil
}

// RoleExists returns whether or not the role with the given roleName exists, and any error that occurred.
func RoleExists(tx *sql.Tx, roleID int) (bool, error) {
	exists := false
	err := tx.QueryRow(`SELECT EXISTS(SELECT * FROM role WHERE role.id = $1)`, roleID).Scan(&exists)
	return exists, err
}

// GetDSNameFromID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func GetDSNameFromID(tx *sql.Tx, id int) (tc.DeliveryServiceName, bool, error) {
	name := tc.DeliveryServiceName("")
	if err := tx.QueryRow(`SELECT xml_id FROM deliveryservice WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return tc.DeliveryServiceName(""), false, nil
		}
		return tc.DeliveryServiceName(""), false, fmt.Errorf("querying xml_id for delivery service ID '%v': %v", id, err)
	}
	return name, true, nil
}

// GetDSNameFromID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func GetDSIDFromXMLID(tx *sql.Tx, xmlID string) (int, bool, error) {
	var id int
	if err := tx.QueryRow(`SELECT id FROM deliveryservice WHERE xml_id = $1`, xmlID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}
		return id, false, fmt.Errorf("querying ID for delivery service XMLID '%v': %v", xmlID, err)
	}
	return id, true, nil
}

// GetDSCDNIdFromID loads the DeliveryService's cdn ID from the database, from the delivery service ID. Returns whether the delivery service was found, and any error.
func GetDSCDNIdFromID(tx *sql.Tx, dsID int) (int, bool, error) {
	var cdnID int
	if err := tx.QueryRow(`SELECT cdn_id FROM deliveryservice WHERE id = $1`, dsID).Scan(&cdnID); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("querying cdn_id for delivery service ID '%v': %v", dsID, err)
	}
	return cdnID, true, nil
}

// GetDSTenantIDFromXMLID fetches the ID of the Tenant to whom the Delivery Service identified by the
// the provided XMLID belongs. It returns, in order, the requested ID (if one could be found), a
// boolean indicating whether or not a Delivery Service with the provided xmlid could be found, and
// an error for logging in case something unexpected goes wrong.
func GetDSTenantIDFromXMLID(tx *sql.Tx, xmlid string) (int, bool, error) {
	var id int
	if err := tx.QueryRow(getDSTenantIDFromXMLIDQuery, xmlid).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return -1, false, nil
		}
		return -1, false, fmt.Errorf("Fetching Tenant ID for DS %s: %v", xmlid, err)
	}
	return id, true, nil
}

// GetServerDSNamesByCDN returns a map of ONLINE/REPORTED/ADMIN_DOWN cache names to slice of
// strings which are the XML IDs of the active, non-ANYMAP delivery services to which the cache
// is assigned in the given CDN.
func GetServerDSNamesByCDN(tx *sql.Tx, cdn string) (map[tc.CacheName][]string, error) {
	q := `
SELECT s.host_name, ds.xml_id
FROM deliveryservice_server AS dss
INNER JOIN server AS s ON dss.server = s.id
INNER JOIN deliveryservice AS ds ON ds.id = dss.deliveryservice
INNER JOIN type AS dt ON dt.id = ds.type
INNER JOIN profile AS p ON p.id = s.profile
INNER JOIN status AS st ON st.id = s.status
WHERE ds.cdn_id = (SELECT id FROM cdn WHERE name = $1)
AND ds.active = $2
AND dt.name != $3
AND p.routing_disabled = false
AND (
	st.name = $4
	OR st.name = $5
	OR st.name = $6
)
`
	rows, err := tx.Query(q, cdn, tc.DSActiveStateActive, tc.DSTypeAnyMap, tc.CacheStatusOnline, tc.CacheStatusReported, tc.CacheStatusAdminDown)
	if err != nil {
		return nil, errors.New("querying server deliveryservice names by CDN: " + err.Error())
	}
	defer log.Close(rows, "closing rows after querying server deliveryservice names by CDN")

	serverDSes := map[tc.CacheName][]string{}
	for rows.Next() {
		ds := ""
		server := ""
		if err := rows.Scan(&server, &ds); err != nil {
			return nil, errors.New("scanning server deliveryservice names: " + err.Error())
		}
		serverDSes[tc.CacheName(server)] = append(serverDSes[tc.CacheName(server)], ds)
	}
	return serverDSes, nil
}

// returns returns the delivery service name and cdn, whether it existed, and any error.
func GetDSNameAndCDNFromID(tx *sql.Tx, id int) (tc.DeliveryServiceName, tc.CDNName, bool, error) {
	name := tc.DeliveryServiceName("")
	cdn := tc.CDNName("")
	if err := tx.QueryRow(`
SELECT ds.xml_id, cdn.name
FROM deliveryservice as ds
JOIN cdn on cdn.id = ds.cdn_id
WHERE ds.id = $1
`, id).Scan(&name, &cdn); err != nil {
		if err == sql.ErrNoRows {
			return tc.DeliveryServiceName(""), tc.CDNName(""), false, nil
		}
		return tc.DeliveryServiceName(""), tc.CDNName(""), false, errors.New("querying delivery service name and CDN name: " + err.Error())
	}
	return name, cdn, true, nil
}

// GetDSIDAndCDNFromName returns the delivery service ID and cdn name given from the delivery service name, whether a result existed, and any error.
func GetDSIDAndCDNFromName(tx *sql.Tx, xmlID string) (int, tc.CDNName, bool, error) {
	dsId := 0
	cdn := tc.CDNName("")
	if err := tx.QueryRow(`
SELECT ds.id, cdn.name
FROM deliveryservice as ds
JOIN cdn on cdn.id = ds.cdn_id
WHERE ds.xml_id = $1
`, xmlID).Scan(&dsId, &cdn); err != nil {
		if err == sql.ErrNoRows {
			return dsId, tc.CDNName(""), false, nil
		}
		return dsId, tc.CDNName(""), false, errors.New("querying delivery service ID and CDN name: " + err.Error())
	}
	return dsId, cdn, true, nil
}

// GetFederationResolversByFederationID fetches all of the federation resolvers currently assigned to a federation.
// In the event of an error, it will return an empty slice and the error.
func GetFederationResolversByFederationID(tx *sql.Tx, fedID int) ([]tc.FederationResolver, error) {
	qry := `
		SELECT
		  fr.ip_address,
		  frt.name as resolver_type,
		  ffr.federation_resolver
		FROM
		  federation_federation_resolver ffr
		  JOIN federation_resolver fr ON ffr.federation_resolver = fr.id
		  JOIN type frt on fr.type = frt.id
		WHERE
		  ffr.federation = $1
		ORDER BY fr.ip_address
	`
	rows, err := tx.Query(qry, fedID)
	if err != nil {
		return nil, fmt.Errorf(
			"error querying federation_resolvers by federation ID [%d]: %s", fedID, err.Error(),
		)
	}
	defer rows.Close()

	resolvers := []tc.FederationResolver{}
	for rows.Next() {
		fr := tc.FederationResolver{}
		err := rows.Scan(
			&fr.IPAddress,
			&fr.Type,
			&fr.ID,
		)
		if err != nil {
			return resolvers, fmt.Errorf(
				"error scanning federation_resolvers rows for federation ID [%d]: %s", fedID, err.Error(),
			)
		}
		resolvers = append(resolvers, fr)
	}
	return resolvers, nil
}

// GetFederationNameFromID returns the federation's name, whether a federation with ID exists, or any error.
func GetFederationNameFromID(id int, tx *sql.Tx) (string, bool, error) {
	var name string
	if err := tx.QueryRow(`SELECT cname from federation where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return name, false, fmt.Errorf(
			"error querying federation name from id [%d]: %s", id, err.Error(),
		)
	}
	return name, true, nil
}

// GetCDNIDFromFedID returns the ID of the CDN for the current federation.
func GetCDNIDFromFedID(id int, tx *sql.Tx) (int, bool, error) {
	var cdnID int
	if err := tx.QueryRow(`SELECT cdn_id FROM deliveryservice WHERE id = (SELECT deliveryservice FROM federation_deliveryservice WHERE federation = $1)`, id).Scan(&cdnID); err != nil {
		if err == sql.ErrNoRows {
			return cdnID, false, nil
		}
		return cdnID, false, err
	}
	return cdnID, true, nil
}

// GetCDNIDFromFedResolverID returns the IDs of the CDNs that the fed resolver is associated with.
func GetCDNIDsFromFedResolverID(id int, tx *sql.Tx) ([]int, bool, error) {
	var cdnIDs []int
	var cdnID int
	rows, err := tx.Query(`SELECT cdn_id FROM deliveryservice WHERE id = ANY(SELECT deliveryservice FROM federation_deliveryservice fds JOIN federation_federation_resolver ffr ON ffr.federation = fds.federation WHERE ffr.federation_resolver = $1)`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return cdnIDs, false, nil
		}
		return cdnIDs, false, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&cdnID); err != nil {
			return cdnIDs, false, errors.New("scanning cdn IDs: " + err.Error())
		}
		cdnIDs = append(cdnIDs, cdnID)
	}
	return cdnIDs, true, nil
}

// GetProfileNameFromID returns the profile's name, whether a profile with ID exists, or any error.
func GetProfileNameFromID(id int, tx *sql.Tx) (string, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name from profile where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying profile name from id: " + err.Error())
	}
	return name, true, nil
}

// GetProfileIDFromName returns the profile's ID, whether a profile with name exists, or any error.
func GetProfileIDFromName(name string, tx *sql.Tx) (int, bool, error) {
	id := 0
	if err := tx.QueryRow(`SELECT id from profile where name = $1`, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying profile id from name: " + err.Error())
	}
	return id, true, nil
}

// GetServerCapabilitiesFromName returns the server's capabilities.
func GetServerCapabilitiesFromName(name string, tx *sql.Tx) ([]string, error) {
	var caps []string
	q := `SELECT ARRAY(SELECT ssc.server_capability FROM server s JOIN server_server_capability ssc ON s.id = ssc.server WHERE s.host_name = $1 ORDER BY ssc.server_capability);`
	rows, err := tx.Query(q, name)
	if err != nil {
		return nil, errors.New("querying server capabilities from name: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(pq.Array(&caps)); err != nil {
			return nil, errors.New("scanning capability: " + err.Error())
		}
	}
	return caps, nil
}

// GetServerCapabilitiesOfServers gets all of the server capabilities of the given server hostnames.
func GetServerCapabilitiesOfServers(names []string, tx *sql.Tx) (map[string][]string, error) {
	serverCaps := make(map[string][]string, len(names))
	q := `
SELECT
  s.host_name,
  ARRAY_REMOVE(ARRAY_AGG(ssc.server_capability ORDER BY ssc.server_capability), NULL) AS capabilities
FROM server s
LEFT JOIN server_server_capability ssc ON s.id = ssc.server
WHERE
  s.host_name = ANY($1)
GROUP BY s.host_name
`
	rows, err := tx.Query(q, pq.Array(&names))
	if err != nil {
		return nil, errors.New("querying server capabilities by host names: " + err.Error())
	}
	defer log.Close(rows, "closing rows in GetServerCapabilitiesOfServers")

	for rows.Next() {
		hostname := ""
		caps := []string{}
		if err := rows.Scan(&hostname, pq.Array(&caps)); err != nil {
			return nil, errors.New("scanning server capabilities: " + err.Error())
		}
		serverCaps[hostname] = caps
	}
	return serverCaps, nil
}

// GetRequiredCapabilitiesOfDeliveryServices gets all of the required capabilities of the given delivery service IDs.
func GetRequiredCapabilitiesOfDeliveryServices(ids []int, tx *sql.Tx) (map[int][]string, error) {
	queryIDs := make([]int64, 0, len(ids))
	for _, id := range ids {
		queryIDs = append(queryIDs, int64(id))
	}
	dsCaps := make(map[int][]string, len(ids))
	q := `
SELECT
  ds.id,
  ARRAY_REMOVE((ds.required_capabilities), NULL) AS required_capabilities
FROM deliveryservice ds
WHERE ds.id = ANY($1)
GROUP BY ds.id, ds.required_capabilities
`
	rows, err := tx.Query(q, pq.Array(&queryIDs))
	if err != nil {
		return nil, errors.New("querying delivery service required capabilities by IDs: " + err.Error())
	}
	defer log.Close(rows, "closing rows in GetRequiredCapabilitiesOfDeliveryServices")

	for rows.Next() {
		id := 0
		caps := []string{}
		if err := rows.Scan(&id, pq.Array(&caps)); err != nil {
			return nil, errors.New("scanning required capabilities: " + err.Error())
		}
		dsCaps[id] = caps
	}
	return dsCaps, nil
}

const dsrExistsQuery = `
SELECT EXISTS(
	SELECT id
	FROM deliveryservice_request
	WHERE status <> 'complete' AND
		status <> 'rejected' AND
		deliveryservice ->> 'xmlId' = $1
)
`

// DSRExistsWithXMLID returns whether or not an **open** Delivery Service
// Request with the given xmlid exists, and any error that occurred.
func DSRExistsWithXMLID(xmlid string, tx *sql.Tx) (bool, error) {
	if tx == nil {
		return false, errors.New("checking for DSR with nil transaction")
	}

	var exists bool
	err := tx.QueryRow(dsrExistsQuery, xmlid).Scan(&exists)
	return exists, err
}

// ScanCachegroupsServerCapabilities, given rows of (server ID, CDN ID, cachegroup name, server capabilities),
// returns a map of cachegroup names to server IDs, a map of server IDs to a map of their capabilities,
// a map of server IDs to CDN IDs, and an error (if one occurs).
func ScanCachegroupsServerCapabilities(rows *sql.Rows) (map[string][]int, map[int]map[string]struct{}, map[int]int, error) {
	defer log.Close(rows, "closing rows in ScanCachegroupsServerCapabilities")

	cachegroupServers := make(map[string][]int)
	serverCapabilities := make(map[int]map[string]struct{})
	serverCDNs := make(map[int]int)
	for rows.Next() {
		serverID := 0
		cdnID := 0
		cachegroup := ""
		serverCap := []string{}
		if err := rows.Scan(&serverID, &cdnID, &cachegroup, pq.Array(&serverCap)); err != nil {
			return nil, nil, nil, fmt.Errorf("scanning rows in ScanCachegroupsServerCapabilities: %v", err)
		}
		cachegroupServers[cachegroup] = append(cachegroupServers[cachegroup], serverID)
		serverCapabilities[serverID] = make(map[string]struct{}, len(serverCap))
		serverCDNs[serverID] = cdnID
		for _, sc := range serverCap {
			serverCapabilities[serverID][sc] = struct{}{}
		}
	}
	return cachegroupServers, serverCapabilities, serverCDNs, nil
}

// GetDSRequiredCapabilitiesFromID returns the server's capabilities.
func GetDSRequiredCapabilitiesFromID(id int, tx *sql.Tx) ([]string, error) {
	q := `
	SELECT required_capabilities
	FROM deliveryservice
	WHERE id = $1
	ORDER BY required_capabilities`

	caps := []string{}
	if err := tx.QueryRow(q, id).Scan(pq.Array(&caps)); err != nil {
		return nil, errors.New("getting/ scanning capability: " + err.Error())
	}
	return caps, nil
}

// Returns true if the cdn exists
func CDNExists(cdnName string, tx *sql.Tx) (bool, error) {
	var id int
	if err := tx.QueryRow(`SELECT id FROM cdn WHERE name = $1`, cdnName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("Error querying CDN name: " + err.Error())
	}
	return true, nil
}

// GetCDNNameFromID gets the CDN name from the given CDN ID.
func GetCDNNameFromID(tx *sql.Tx, id int64) (tc.CDNName, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying CDN name from ID: " + err.Error())
	}
	return tc.CDNName(name), true, nil
}

// GetCDNNameFromServerID gets the CDN name for the server with the given ID.
func GetCDNNameFromServerID(tx *sql.Tx, serverId int64) (tc.CDNName, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = (SELECT cdn_id FROM server WHERE id=$1)`, serverId).Scan(&name); err != nil {
		return "", fmt.Errorf("querying CDN name from server ID: %w", err)
	}
	return tc.CDNName(name), nil
}

// GetCDNIDFromName returns the ID of the CDN if a CDN with the name exists
func GetCDNIDFromName(tx *sql.Tx, name tc.CDNName) (int, bool, error) {
	id := 0
	if err := tx.QueryRow(`SELECT id FROM cdn WHERE name = $1`, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}
		return id, false, errors.New("querying CDN ID from name: " + err.Error())
	}
	return id, true, nil
}

// GetCDNDomainFromName returns the domain, whether the cdn exists, and any error.
func GetCDNDomainFromName(tx *sql.Tx, cdnName tc.CDNName) (string, bool, error) {
	domain := ""
	if err := tx.QueryRow(`SELECT domain_name FROM cdn WHERE name = $1`, cdnName).Scan(&domain); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("Error querying CDN name: " + err.Error())
	}
	return domain, true, nil
}

// GetServerInterfaces, given the IDs of one or more servers, returns all of their network
// interfaces mapped by their ids, or an error if one occurs during retrieval.
func GetServersInterfaces(ids []int, tx *sql.Tx) (map[int]map[string]tc.ServerInterfaceInfoV40, error) {
	q := `
	SELECT max_bandwidth,
	       monitor,
	       mtu,
	       name,
	       server,
	       router_host_name,
	       router_port_name
	FROM interface
	WHERE interface.server = ANY ($1)
	`
	ifaceRows, err := tx.Query(q, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer ifaceRows.Close()

	infs := map[int]map[string]tc.ServerInterfaceInfoV40{}
	for ifaceRows.Next() {
		var inf tc.ServerInterfaceInfoV40
		var server int
		if err := ifaceRows.Scan(&inf.MaxBandwidth, &inf.Monitor, &inf.MTU, &inf.Name, &server, &inf.RouterHostName, &inf.RouterPortName); err != nil {
			return nil, err
		}
		if _, ok := infs[server]; !ok {
			infs[server] = make(map[string]tc.ServerInterfaceInfoV40)
		}

		infs[server][inf.Name] = inf
	}

	q = `
	SELECT address,
	       gateway,
	       service_address,
	       interface,
	       server
	FROM ip_address
	WHERE ip_address.server = ANY ($1)
	`
	ipRows, err := tx.Query(q, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer ipRows.Close()

	for ipRows.Next() {
		var ip tc.ServerIPAddress
		var inf string
		var server int
		if err = ipRows.Scan(&ip.Address, &ip.Gateway, &ip.ServiceAddress, &inf, &server); err != nil {
			return nil, err
		}

		ifaces, ok := infs[server]
		if !ok {
			return nil, fmt.Errorf("retrieved ip_address with server not previously found: %d", server)
		}

		iface, ok := ifaces[inf]
		if !ok {
			return nil, fmt.Errorf("retrieved ip_address with interface not previously found: %s", inf)
		}
		iface.IPAddresses = append(iface.IPAddresses, ip)
		infs[server][inf] = iface
	}

	return infs, nil
}

// GetStatusByID returns a Status struct, a bool for whether or not a status of the given ID exists, and an error (if one occurs).
func GetStatusByID(id int, tx *sql.Tx) (tc.StatusNullable, bool, error) {
	q := `
SELECT
  description,
  id,
  last_updated,
  name
FROM
  status s
WHERE
  id = $1
`
	row := tc.StatusNullable{}
	if err := tx.QueryRow(q, id).Scan(
		&row.Description,
		&row.ID,
		&row.LastUpdated,
		&row.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return row, false, nil
		}
		return row, false, fmt.Errorf("querying status id %d: %v", id, err.Error())
	}
	return row, true, nil
}

// GetStatusByName returns a Status struct, a bool for whether or not a status of the given name exists, and an error (if one occurs).
func GetStatusByName(name string, tx *sql.Tx) (tc.StatusNullable, bool, error) {
	q := `
SELECT
  description,
  id,
  last_updated,
  name
FROM
  status s
WHERE
  name = $1
`
	row := tc.StatusNullable{}
	if err := tx.QueryRow(q, name).Scan(
		&row.Description,
		&row.ID,
		&row.LastUpdated,
		&row.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return row, false, nil
		}
		return row, false, fmt.Errorf("querying status name %s: %v", name, err.Error())
	}
	return row, true, nil
}

// GetServerIDFromName gets server id from a given name
func GetServerIDFromName(serverName string, tx *sql.Tx) (int, bool, error) {
	id := 0
	if err := tx.QueryRow(`SELECT id FROM server WHERE host_name = $1`, serverName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}
		return id, false, errors.New("querying server name: " + err.Error())
	}
	return id, true, nil
}

func GetServerNameFromID(tx *sql.Tx, id int64) (string, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT host_name FROM server WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying server name: " + err.Error())
	}
	return name, true, nil
}

// language=sql
const getServerInfoBaseQuery = `
SELECT
  s.cachegroup,
  c.name,
  s.host_name,
  s.domain_name,
  s.cdn_id,
  t.name,
  s.id,
  status.name
FROM
  server s JOIN type t ON s.type = t.id
  JOIN cachegroup c on s.cachegroup = c.id
  JOIN status on status.id = s.status
`

// GetServerInfosFromIDs returns the ServerInfo structs of the given server IDs or an error if any occur.
func GetServerInfosFromIDs(tx *sql.Tx, ids []int) ([]tc.ServerInfo, error) {
	qry := getServerInfoBaseQuery + `
WHERE s.id = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(ids))
	if err != nil {
		return nil, errors.New("querying server info: " + err.Error())
	}
	return scanServerInfoRows(rows)
}

// GetServerInfosFromHostNames returns the ServerInfo structs of the given server host names or an error if any occur.
func GetServerInfosFromHostNames(tx *sql.Tx, hostNames []string) ([]tc.ServerInfo, error) {
	qry := getServerInfoBaseQuery + `
WHERE s.host_name = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(hostNames))
	if err != nil {
		return nil, errors.New("querying server info: " + err.Error())
	}
	return scanServerInfoRows(rows)
}

func scanServerInfoRows(rows *sql.Rows) ([]tc.ServerInfo, error) {
	defer log.Close(rows, "error closing rows")
	servers := []tc.ServerInfo{}
	for rows.Next() {
		s := tc.ServerInfo{}
		if err := rows.Scan(&s.CachegroupID, &s.Cachegroup, &s.HostName, &s.DomainName, &s.CDNID, &s.Type, &s.ID, &s.Status); err != nil {
			return nil, errors.New("scanning server info: " + err.Error())
		}
		servers = append(servers, s)
	}
	return servers, nil
}

// GetServerInfo returns a ServerInfo struct, whether the server exists, and an error (if one occurs).
func GetServerInfo(serverID int, tx *sql.Tx) (tc.ServerInfo, bool, error) {
	servers, err := GetServerInfosFromIDs(tx, []int{serverID})
	if err != nil {
		return tc.ServerInfo{}, false, fmt.Errorf("getting server info: %v", err)
	}
	if len(servers) == 0 {
		return tc.ServerInfo{}, false, nil
	}
	if len(servers) != 1 {
		return tc.ServerInfo{}, false, fmt.Errorf("getting server info - expected row count: 1, actual: %d", len(servers))
	}
	return servers[0], true, nil
}

func GetCDNDSes(tx *sql.Tx, cdn tc.CDNName) (map[string]struct{}, error) {
	dses := map[string]struct{}{}
	qry := `SELECT xml_id from deliveryservice where cdn_id = (select id from cdn where name = $1)`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		ds := ""
		if err := rows.Scan(&ds); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		dses[ds] = struct{}{}
	}
	return dses, nil
}

func GetCDNs(tx *sql.Tx) (map[tc.CDNName]struct{}, error) {
	cdns := map[tc.CDNName]struct{}{}
	qry := `SELECT name from cdn;`
	rows, err := tx.Query(qry)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		cdn := tc.CDNName("")
		if err := rows.Scan(&cdn); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		cdns[cdn] = struct{}{}
	}
	return cdns, nil
}

// GetGlobalParams returns the value of the global param, whether it existed, or any error
func GetGlobalParam(tx *sql.Tx, name string) (string, bool, error) {
	return GetParam(tx, name, "global")
}

// GetParam returns the value of the param, whether it existed, or any error.
func GetParam(tx *sql.Tx, name string, configFile string) (string, bool, error) {
	val := ""
	if err := tx.QueryRow(`select value from parameter where name = $1 and config_file = $2`, name, configFile).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("Error querying global paramter '" + name + "': " + err.Error())
	}
	return val, true, nil
}

// GetParamNameByID returns the name of the param, whether it existed, or any error.
func GetParamNameByID(tx *sql.Tx, id int) (string, bool, error) {
	name := ""
	if err := tx.QueryRow(`select name from parameter where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("Error querying global paramter %v: %v", id, err.Error())
	}
	return name, true, nil
}

// GetCacheGroupNameFromID Get Cache Group name from a given ID
func GetCacheGroupNameFromID(tx *sql.Tx, id int) (tc.CacheGroupName, bool, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cachegroup WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying cachegroup ID: " + err.Error())
	}
	return tc.CacheGroupName(name), true, nil
}

// TopologyExists checks if a Topology with the given name exists.
// Returns whether or not the Topology exists, along with any encountered error.
func TopologyExists(tx *sql.Tx, name string) (bool, error) {
	q := `
	SELECT COUNT("name")
	FROM topology
	WHERE name = $1
	`
	var count int
	var err error
	if err = tx.QueryRow(q, name).Scan(&count); err != nil {
		err = fmt.Errorf("querying topologies: %s", err)
	}
	return count > 0, err
}

// CheckTopology returns an error if the given Topology does not exist or if
// one of the Topology's Cache Groups is empty with respect to the Delivery
// Service's CDN. Note that this can panic if ds does not have a properly set
// CDNID.
func CheckTopology(tx *sqlx.Tx, ds tc.DeliveryServiceV5) (int, error, error) {
	if ds.Topology == nil {
		return http.StatusOK, nil, nil
	}

	cacheGroupIDs, _, err := GetTopologyCachegroups(tx.Tx, *ds.Topology)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("getting topology cachegroups: %w", err)
	}
	if len(cacheGroupIDs) == 0 {
		return http.StatusBadRequest, fmt.Errorf("no such Topology '%s'", *ds.Topology), nil
	}

	if err = topology_validation.CheckForEmptyCacheGroups(tx, cacheGroupIDs, []int{ds.CDNID}, true, []int{}); err != nil {
		return http.StatusBadRequest, fmt.Errorf("empty cachegroups in Topology %s found for CDN %d: %w", *ds.Topology, ds.CDNID, err), nil
	}

	return http.StatusOK, nil, nil
}

// GetTopologyCachegroups returns an array of cachegroup IDs and an array of cachegroup
// names for the given topology, or any error.
func GetTopologyCachegroups(tx *sql.Tx, name string) ([]int, []string, error) {
	q := `
	SELECT ARRAY_AGG(c.id), ARRAY_AGG(tc.cachegroup)
	FROM topology_cachegroup tc
	JOIN cachegroup c ON tc.cachegroup = c."name"
	WHERE tc.topology = $1
`
	int64Ids := []int64{}
	names := []string{}
	if err := tx.QueryRow(q, name).Scan(pq.Array(&int64Ids), pq.Array(&names)); err != nil {
		return nil, nil, fmt.Errorf("querying topology '%s' cachegroups: %s", name, err)
	}

	ids := make([]int, len(int64Ids))
	for index, int64Id := range int64Ids {
		ids[index] = int(int64Id)
	}

	return ids, names, nil
}

// GetDeliveryServicesWithTopologies returns a list containing the delivery services in the given dsIDs
// list that have a topology assigned. An error indicates unexpected errors that occurred when querying.
func GetDeliveryServicesWithTopologies(tx *sql.Tx, dsIDs []int) ([]int, error) {
	q := `
SELECT
  id
FROM
  deliveryservice
WHERE
  id = ANY($1::bigint[])
  AND topology IS NOT NULL
`
	rows, err := tx.Query(q, pq.Array(dsIDs))
	if err != nil {
		return nil, errors.New("querying deliveryservice topologies: " + err.Error())
	}
	defer log.Close(rows, "error closing rows")
	dses := make([]int, 0)
	for rows.Next() {
		id := 0
		if err := rows.Scan(&id); err != nil {
			return nil, errors.New("scanning deliveryservice id: " + err.Error())
		}
		dses = append(dses, id)
	}
	return dses, nil
}

// GetDeliveryServiceCDNsByTopology returns a slice of CDN IDs for all delivery services
// assigned to the given topology.
func GetDeliveryServiceCDNsByTopology(tx *sql.Tx, topology string) ([]int, error) {
	q := `
SELECT
  COALESCE(ARRAY_AGG(DISTINCT d.cdn_id), '{}'::BIGINT[])
FROM
  deliveryservice d
WHERE
  d.topology = $1
`
	cdnIDs := []int64{}
	if err := tx.QueryRow(q, topology).Scan(pq.Array(&cdnIDs)); err != nil {
		return nil, fmt.Errorf("in GetDeliveryServiceCDNsByTopology: querying deliveryservices by topology '%s': %v", topology, err)
	}
	res := make([]int, len(cdnIDs))
	for i, id := range cdnIDs {
		res[i] = int(id)
	}
	return res, nil
}

// CheckCachegroupHasTopologyBasedDeliveryServicesOnCDN returns true if the given cachegroup is assigned to
// any topologies with delivery services assigned on the given CDN.
func CachegroupHasTopologyBasedDeliveryServicesOnCDN(tx *sql.Tx, cachegroupID int, CDNID int) (bool, error) {
	q := `
SELECT EXISTS(
  SELECT
    1
  FROM cachegroup c
  JOIN topology_cachegroup tc on c.name = tc.cachegroup
  JOIN topology t ON tc.topology = t.name
  JOIN deliveryservice d on t.name = d.topology
  WHERE
    c.id = $1
    AND d.cdn_id = $2
)
`
	res := false
	if err := tx.QueryRow(q, cachegroupID, CDNID).Scan(&res); err != nil {
		return false, fmt.Errorf("in CachegroupHasTopologyBasedDeliveryServicesOnCDN: %v", err)
	}
	return res, nil
}

// GetFederationIDForUserIDByXMLID retrieves the ID of the Federation assigned to the user defined by
// userID on the Delivery Service identified by xmlid. If no such federation exists, the boolean
// returned will be 'false', while the error indicates unexpected errors that occurred when querying.
func GetFederationIDForUserIDByXMLID(tx *sql.Tx, userID int, xmlid string) (uint, bool, error) {
	var id uint
	if err := tx.QueryRow(getFederationIDForUserIDByXMLIDQuery, xmlid, userID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("Getting Federation ID for user #%d by DS XMLID '%s': %v", userID, xmlid, err)
	}
	return id, true, nil
}

// UsernameExists reports whether or not the the given username exists as a user in the database to
// which the passed transaction refers. If anything goes wrong when checking the existence of said
// user, the error is directly returned to the caller. Note that in that case, no real meaning
// should be assigned to the returned boolean value.
func UsernameExists(uname string, tx *sql.Tx) (bool, error) {
	row := tx.QueryRow(`SELECT EXISTS(SELECT * FROM tm_user WHERE tm_user.username=$1)`, uname)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

// GetTypeIDByName reports the id of the type and whether or not a type exists with the given name.
func GetTypeIDByName(t string, tx *sql.Tx) (int, bool, error) {
	id := 0
	if err := tx.QueryRow(`SELECT id FROM type WHERE name = $1`, t).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}
		return id, false, errors.New("querying type id: " + err.Error())
	}
	return id, true, nil
}

// GetUserByEmail retrieves the user with the given email. If no such user exists, the boolean
// returned will be 'false', while the error indicates unexpected errors that occurred when querying.
func GetUserByEmail(email string, tx *sql.Tx) (tc.User, bool, error) {
	row := tx.QueryRow(getUserByEmailQuery, email)
	return scanUserRow(row)
}

// GetUserByID returns the user with the requested ID if one exists. The second return value is a
// boolean indicating whether said user actually did exist, and the third contains any error
// encountered along the way.
func GetUserByID(id int, tx *sql.Tx) (tc.User, bool, error) {
	row := tx.QueryRow(getUserByIDQuery, id)
	return scanUserRow(row)
}

func scanUserRow(row *sql.Row) (tc.User, bool, error) {
	var u tc.User
	err := row.Scan(&u.AddressLine1,
		&u.AddressLine2,
		&u.City,
		&u.Company,
		&u.Country,
		&u.Email,
		&u.FullName,
		&u.GID,
		&u.ID,
		&u.LastUpdated,
		&u.NewUser,
		&u.PhoneNumber,
		&u.PostalCode,
		&u.PublicSSHKey,
		&u.RegistrationSent,
		&u.Role,
		&u.RoleName,
		&u.StateOrProvince,
		&u.Tenant,
		&u.TenantID,
		&u.Token,
		&u.UID,
		&u.Username)
	if err == sql.ErrNoRows {
		return u, false, nil
	}
	return u, true, err
}

// CachegroupParameterAssociationExists returns whether a cachegroup parameter association with the given parameter id exists, and any error.
func CachegroupParameterAssociationExists(id int, cachegroup int, tx *sql.Tx) (bool, error) {
	count := 0
	if err := tx.QueryRow(`SELECT count(*) from cachegroup_parameter where parameter = $1 and cachegroup = $2`, id, cachegroup).Scan(&count); err != nil {
		return false, errors.New("querying cachegroup parameter existence: " + err.Error())
	}
	return count > 0, nil
}

// GetDeliveryServiceTypeAndCDNName returns the type and the CDN name of the deliveryservice.
func GetDeliveryServiceTypeAndCDNName(dsID int, tx *sql.Tx) (tc.DSType, string, bool, error) {
	var dsType tc.DSType
	var cdnName string
	if err := tx.QueryRow(`SELECT t.name, c.name as cdn FROM deliveryservice as ds JOIN type t ON ds.type = t.id JOIN cdn c ON c.id = ds.cdn_id WHERE ds.id=$1`, dsID).Scan(&dsType, &cdnName); err != nil {
		if err == sql.ErrNoRows {
			return tc.DSTypeInvalid, cdnName, false, nil
		}
		return tc.DSTypeInvalid, cdnName, false, errors.New("querying type from delivery service: " + err.Error())
	}
	return dsType, cdnName, true, nil
}

// GetDeliveryServiceTypeRequiredCapabilitiesAndTopology returns the type of the deliveryservice and the name of its topology.
func GetDeliveryServiceTypeRequiredCapabilitiesAndTopology(dsID int, tx *sql.Tx) (tc.DSType, []string, *string, bool, error) {
	var dsType tc.DSType
	var reqCap []string
	var topology *string
	q := `
SELECT
  t.name,
  ARRAY_REMOVE(ds.required_capabilities, NULL) AS required_capabilities,
  ds.topology
FROM deliveryservice AS ds
JOIN type t ON ds.type = t.id
WHERE ds.id = $1
GROUP BY t.name, ds.topology, ds.required_capabilities
`
	if err := tx.QueryRow(q, dsID).Scan(&dsType, pq.Array(&reqCap), &topology); err != nil {
		if err == sql.ErrNoRows {
			return tc.DSTypeInvalid, nil, nil, false, nil
		}
		return tc.DSTypeInvalid, nil, nil, false, errors.New("querying type from delivery service: " + err.Error())
	}
	return dsType, reqCap, topology, true, nil
}

// CheckOriginServerInDSCG checks if a DS has ORG server and if it does, to make sure the cachegroup is part of DS
func CheckOriginServerInDSCG(tx *sql.Tx, dsID int, dsTopology string) (error, error, int) {
	// get servers and respective cachegroup name that have ORG type in a delivery service
	q := `
		SELECT s.host_name, c.name
		FROM server s
			INNER JOIN deliveryservice_server ds ON ds.server = s.id
			INNER JOIN type t ON t.id = s.type
			INNER JOIN cachegroup c ON c.id = s.cachegroup
		WHERE ds.deliveryservice=$1 AND t.name=$2
	`

	serverName := ""
	cacheGroupName := ""
	servers := make(map[string]string)
	var offendingSCG []string
	rows, err := tx.Query(q, dsID, tc.OriginTypeName)
	if err != nil {
		return nil, fmt.Errorf("querying deliveryservice origin server: %s", err), http.StatusInternalServerError
	}
	defer log.Close(rows, "error closing rows")
	for rows.Next() {
		if err := rows.Scan(&serverName, &cacheGroupName); err != nil {
			return nil, fmt.Errorf("querying deliveryservice origin server: %s", err), http.StatusInternalServerError
		}
		servers[cacheGroupName] = serverName
	}

	if len(servers) > 0 {
		//Validation for DS
		_, cachegroups, sysErr := GetTopologyCachegroups(tx, dsTopology)
		if sysErr != nil {
			return nil, fmt.Errorf("validating %s servers in topology: %v", tc.OriginTypeName, sysErr), http.StatusInternalServerError
		}
		// put slice values into map for DS's validation
		topoCachegroups := make(map[string]string)
		for _, cg := range cachegroups {
			topoCachegroups[cg] = ""
		}
		for cg, s := range servers {
			_, ok := topoCachegroups[cg]
			if !ok {
				offendingSCG = append(offendingSCG, fmt.Sprintf("%s (%s)", cg, s))
			}
		}
	}
	if len(offendingSCG) > 0 {
		return errors.New("the following ORG server cachegroups are not in the delivery service's topology (" + dsTopology + "): " + strings.Join(offendingSCG, ", ")), nil, http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

// CheckTopologyOrgServerCGInDSCG checks if ORG server are part of DS. IF they are then the user is not allowed to remove the ORG servers from the associated DS's topology
func CheckTopologyOrgServerCGInDSCG(tx *sql.Tx, cdnIds []int, dsTopology string, topologyCGNames []string) (error, error, int) {
	// get servers and respective cachegroup name that have ORG type for evert delivery service
	q := `
		SELECT ARRAY_AGG(d.xml_id), s.id, c.name
		FROM server s
			INNER JOIN deliveryservice_server ds ON ds.server = s.id
			INNER JOIN deliveryservice d ON d.id = ds.deliveryservice
			INNER JOIN type t ON t.id = s.type
			INNER JOIN cachegroup c ON c.id = s.cachegroup
		WHERE d.cdn_id =ANY($1) AND t.name=$2 AND d.topology=$3
		GROUP BY s.id, c.name
	`
	serverId := ""
	cacheGroupName := ""
	dsNames := []string{}
	serversCG := make(map[string]string)
	serversDS := make(map[string][]string)
	rows, err := tx.Query(q, pq.Array(cdnIds), tc.OriginTypeName, dsTopology)
	if err != nil {
		return nil, fmt.Errorf("querying deliveryservice origin server: %s", err), http.StatusInternalServerError
	}
	defer log.Close(rows, "error closing rows")
	for rows.Next() {
		if err := rows.Scan(pq.Array(&dsNames), &serverId, &cacheGroupName); err != nil {
			return nil, fmt.Errorf("querying deliveryservice origin server: %s", err), http.StatusInternalServerError
		}
		serversCG[cacheGroupName] = serverId
		serversDS[serverId] = dsNames
	}

	var offendingDSSerCG []string
	// put slice values into map for Topology's validation
	topoCacheGroupNames := make(map[string]string)
	for _, currentCG := range topologyCGNames {
		topoCacheGroupNames[currentCG] = ""
	}
	for cg, s := range serversCG {
		_, currentTopoCGOk := topoCacheGroupNames[cg]
		if !currentTopoCGOk {
			offendingDSSerCG = append(offendingDSSerCG, fmt.Sprintf("cachegroup=%s (serverID=%s, delivery_services=%s)", cg, s, serversDS[s]))
		}
	}
	if len(offendingDSSerCG) > 0 {
		return errors.New("ORG servers are assigned to delivery services that use this topology, and their cachegroups cannot be removed: " + strings.Join(offendingDSSerCG, ", ")), nil, http.StatusBadRequest
	}
	return nil, nil, http.StatusOK
}

// GetCDNNameFromProfileID returns the cdn name for the provided profile ID.
func GetCDNNameFromProfileID(tx *sql.Tx, id int) (tc.CDNName, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = (SELECT cdn FROM profile WHERE id = $1)`, id).Scan(&name); err != nil {
		return "", errors.New("querying CDN name from profile ID: " + err.Error())
	}
	return tc.CDNName(name), nil
}

// GetCDNNameFromProfileName returns the cdn name for the provided profile name.
func GetCDNNameFromProfileName(tx *sql.Tx, profileName string) (tc.CDNName, error) {
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = (SELECT cdn FROM profile WHERE name = $1)`, profileName).Scan(&name); err != nil {
		return "", errors.New("querying CDN name from profile name: " + err.Error())
	}
	return tc.CDNName(name), nil
}

// GetServerIDsFromCachegroupNames returns a list of servers IDs for a list of cachegroup IDs.
func GetServerIDsFromCachegroupNames(tx *sql.Tx, cgID []string) ([]int64, error) {
	var serverIDs []int64
	var serverID int64
	query := `SELECT server.id FROM server JOIN cachegroup cg ON cg.id = server.cachegroup where cg.name = ANY($1)`
	rows, err := tx.Query(query, pq.Array(cgID))
	if err != nil {
		return serverIDs, errors.New("getting server IDs from cachegroup names : " + err.Error())
	}
	defer log.Close(rows, "could not close rows in GetServerIDsFromCachegroupNames")
	for rows.Next() {
		err = rows.Scan(&serverID)
		if err != nil {
			return serverIDs, errors.New("scanning server ID : " + err.Error())
		}
		serverIDs = append(serverIDs, serverID)
	}
	return serverIDs, nil
}

// GetCDNNamesFromServerIds returns a list of cdn names for a list of server IDs.
func GetCDNNamesFromServerIds(tx *sql.Tx, serverIds []int64) ([]string, error) {
	var cdns []string
	cdn := ""
	query := `SELECT DISTINCT(name) FROM cdn JOIN server ON cdn.id = server.cdn_id WHERE server.id = ANY($1)`
	rows, err := tx.Query(query, pq.Array(serverIds))
	if err != nil {
		return cdns, errors.New("getting cdn name for server : " + err.Error())
	}
	defer log.Close(rows, "could not close rows in GetCDNNamesFromServerIds")
	for rows.Next() {
		err = rows.Scan(&cdn)
		if err != nil {
			return cdns, errors.New("scanning cdn name " + cdn + ": " + err.Error())
		}
		cdns = append(cdns, cdn)
	}
	return cdns, nil
}

// GetCDNNameFromDSXMLID returns the CDN name of the DS associated with the supplied XML ID
func GetCDNNameFromDSXMLID(tx *sql.Tx, dsXMLID string) (string, error) {
	var cdnName string
	query := `SELECT name FROM cdn JOIN deliveryservice ON cdn.id = deliveryservice.cdn_id WHERE deliveryservice.xml_id = $1`
	err := tx.QueryRow(query, dsXMLID).Scan(&cdnName)
	if err != nil {
		return "", err
	}
	return cdnName, nil
}

// GetCDNNamesFromDSIds returns a list of cdn names for a list of DS IDs.
func GetCDNNamesFromDSIds(tx *sql.Tx, dsIds []int) ([]string, error) {
	var cdns []string
	cdn := ""
	query := `SELECT DISTINCT(name) FROM cdn JOIN deliveryservice ON cdn.id = deliveryservice.cdn_id WHERE deliveryservice.id = ANY($1)`
	rows, err := tx.Query(query, pq.Array(dsIds))
	if err != nil {
		return cdns, errors.New("getting cdn name for DS : " + err.Error())
	}
	defer log.Close(rows, "could not close rows in GetCDNNamesFromDSIds")
	for rows.Next() {
		err = rows.Scan(&cdn)
		if err != nil {
			return cdns, errors.New("scanning cdn name " + cdn + ": " + err.Error())
		}
		cdns = append(cdns, cdn)
	}
	return cdns, nil
}

// GetCDNNamesFromProfileIDs returns a list of cdn names for a list of profile IDs.
func GetCDNNamesFromProfileIDs(tx *sql.Tx, profileIDs []int64) ([]string, error) {
	var cdns []string
	cdn := ""
	query := `SELECT DISTINCT(cdn.name) FROM cdn JOIN profile ON cdn.id = profile.cdn WHERE profile.id = ANY($1)`
	rows, err := tx.Query(query, pq.Array(profileIDs))
	if err != nil {
		return cdns, errors.New("getting cdn name for profiles : " + err.Error())
	}
	defer log.Close(rows, "could not close rows in GetCDNNamesFromProfileIDs")
	for rows.Next() {
		err = rows.Scan(&cdn)
		if err != nil {
			return cdns, errors.New("scanning cdn name " + cdn + ": " + err.Error())
		}
		cdns = append(cdns, cdn)
	}
	return cdns, nil
}

// GetDSIDFromStaticDNSEntry returns the delivery service ID associated with the static DNS entry
func GetDSIDFromStaticDNSEntry(tx *sql.Tx, staticDNSEntryID int) (int, error) {
	var dsID int
	query := `SELECT deliveryservice FROM staticdnsentry WHERE id = $1`
	if err := tx.QueryRow(query, staticDNSEntryID).Scan(&dsID); err != nil {
		return -1, errors.New("querying DS ID from static dns entry: " + err.Error())
	}
	return dsID, nil
}

// AppendWhere appends 'extra' safely to the WHERE clause 'where'. What is
// returned is guaranteed to be a valid WHERE clause (including a blank string),
// provided the supplied 'where' and 'extra' clauses are valid.
func AppendWhere(where, extra string) string {
	if where == "" && extra == "" {
		return ""
	}
	if where == "" {
		where = BaseWhere + " "
	} else {
		where += " AND "
	}
	return where + extra
}

// GetRoleIDFromName returns the ID of the role associated with the supplied name.
func GetRoleIDFromName(tx *sql.Tx, roleName string) (int, bool, error) {
	var id int
	if err := tx.QueryRow(`SELECT id FROM role WHERE name = $1`, roleName).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, false, nil
		}
		return id, false, fmt.Errorf("querying role ID from name: %w", err)
	}
	return id, true, nil
}

// GetCDNNameDomain returns the name and domain for a given CDN ID.
func GetCDNNameDomain(cdnID int, tx *sql.Tx) (string, string, error) {
	q := `SELECT cdn.name, cdn.domain_name from cdn where cdn.id = $1`
	cdnName := ""
	cdnDomain := ""
	if err := tx.QueryRow(q, cdnID).Scan(&cdnName, &cdnDomain); err != nil {
		return "", "", fmt.Errorf("getting cdn name and domain for cdn '%v': "+err.Error(), cdnID)
	}
	return cdnName, cdnDomain, nil
}

// GetRegionNameFromID returns the name of the region associated with the supplied ID.
func GetRegionNameFromID(tx *sql.Tx, regionID int) (string, bool, error) {
	var regionName string
	if err := tx.QueryRow(`SELECT name FROM region WHERE id = $1`, regionID).Scan(&regionName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return regionName, false, nil
		}
		return regionName, false, fmt.Errorf("querying region name from ID: %w", err)
	}
	return regionName, true, nil
}

// QueueUpdateForServer sets the config update time for the server to now.
func QueueUpdateForServer(tx *sql.Tx, serverID int64) error {
	query := `
UPDATE public.server
SET config_update_time = now()
WHERE server.id = $1;`

	if _, err := tx.Exec(query, serverID); err != nil {
		return fmt.Errorf("queueing config update for ServerID %d: %w", serverID, err)
	}

	return nil
}

// QueueUpdateForServerWithCachegroupCDN sets the config update time
// for all servers belonging to the specified cachegroup (id) and cdn (id).
func QueueUpdateForServerWithCachegroupCDN(tx *sql.Tx, cgID int, cdnID int64) ([]tc.CacheName, error) {
	q := `
UPDATE public.server
SET config_update_time = now()
WHERE server.cachegroup = $1 AND server.cdn_id = $2
RETURNING server.host_name;`
	rows, err := tx.Query(q, cgID, cdnID)
	if err != nil {
		return nil, fmt.Errorf("updating server config_update_time: %w", err)
	}
	defer log.Close(rows, "error closing rows for QueueUpdateForServerWithCachegroupCDN")
	names := []tc.CacheName{}
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning queue updates: %w", err)
		}
		names = append(names, tc.CacheName(name))
	}
	return names, nil
}

// QueueUpdateForServerWithTopologyCDN sets the config update time
// for all servers within the specific topology (name) and cdn (id).
func QueueUpdateForServerWithTopologyCDN(tx *sql.Tx, topologyName tc.TopologyName, cdnId int64) error {
	query := `
UPDATE public.server
SET config_update_time = now()
FROM public.cachegroup AS cg
INNER JOIN public.topology_cachegroup AS tc ON tc.cachegroup = cg."name"
WHERE cg.id = server.cachegroup
AND tc.topology = $1
AND server.cdn_id = $2;`
	var err error
	if _, err = tx.Exec(query, topologyName, cdnId); err != nil {
		err = fmt.Errorf("queueing updates: %w", err)
	}
	return err
}

// DequeueUpdateForServer sets the config update time equal to the
// config apply time, thereby effectively dequeueing any pending
// updates for the server specified.
func DequeueUpdateForServer(tx *sql.Tx, serverID int64) error {
	query := `
UPDATE public.server
SET config_update_time = config_apply_time
WHERE server.id = $1;`

	if _, err := tx.Exec(query, serverID); err != nil {
		return fmt.Errorf("applying config update for ServerID %d: %w", serverID, err)
	}

	return nil
}

// DequeueUpdateForServerWithCachegroupCDN sets the config update time equal to
// the config apply time, thereby effectively dequeueing any pending
// updates for the servers specified by the cachegroup (id) and the cdn (id).
func DequeueUpdateForServerWithCachegroupCDN(tx *sql.Tx, cgID int, cdnID int64) ([]tc.CacheName, error) {
	q := `
UPDATE public.server
SET config_update_time = config_apply_time
WHERE server.cachegroup = $1
AND server.cdn_id = $2
RETURNING server.host_name;`
	rows, err := tx.Query(q, cgID, cdnID)
	if err != nil {
		return nil, fmt.Errorf("querying queue updates: %w", err)
	}
	defer log.Close(rows, "error closing rows for DequeueUpdateForServerWithCachegroupCDN")
	names := []tc.CacheName{}
	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning queue updates: %w", err)
		}
		names = append(names, tc.CacheName(name))
	}
	return names, nil
}

// DequeueUpdateForServerWithTopologyCDN sets the config update time equal to
// the config apply time, thereby effectively dequeueing any pending
// updates for the servers specified by the topology (name) and the cdn (id).
func DequeueUpdateForServerWithTopologyCDN(tx *sql.Tx, topologyName tc.TopologyName, cdnId int64) error {
	query := `
UPDATE public.server
SET config_update_time = config_apply_time
FROM cachegroup cg
INNER JOIN topology_cachegroup tc ON tc.cachegroup = cg."name"
WHERE cg.id = server.cachegroup
AND tc.topology = $1
AND server.cdn_id = $2;`
	var err error
	if _, err = tx.Exec(query, topologyName, cdnId); err != nil {
		err = fmt.Errorf("queueing updates: %w", err)
	}
	return err
}

// SetApplyUpdateForServer sets the config apply time for the server to now.
func SetApplyUpdateForServer(tx *sql.Tx, serverID int64) error {
	query := `
UPDATE public.server
SET config_apply_time = now()
WHERE server.id = $1;`

	if _, err := tx.Exec(query, serverID); err != nil {
		return fmt.Errorf("applying config update for ServerID %d: %w", serverID, err)
	}

	return nil
}

// SetApplyUpdateForServerWithTime sets the config apply time for the
// server to the provided time.
func SetApplyUpdateForServerWithTime(tx *sql.Tx, serverID int64, applyUpdateTime time.Time) error {
	query := `
UPDATE public.server
SET config_apply_time = $1
WHERE server.id = $2;`

	if _, err := tx.Exec(query, applyUpdateTime, serverID); err != nil {
		return fmt.Errorf("applying config update for ServerID %d with time %v: %w", serverID, applyUpdateTime, err)
	}

	return nil
}

// SetUpdateFailedForServer sets the update failed flag for the server.
func SetUpdateFailedForServer(tx *sql.Tx, serverID int64, failed bool) error {
	query := `
UPDATE public.server
SET config_update_failed = $1
WHERE server.id = $2`
	if _, err := tx.Exec(query, failed, serverID); err != nil {
		return fmt.Errorf("setting config update failed for ServerID %d with value %v: %w", serverID, failed, err)
	}
	return nil
}

// SetRevalFailedForServer sets the reval failed flag for the server.
func SetRevalFailedForServer(tx *sql.Tx, serverID int64, failed bool) error {
	query := `
UPDATE public.server
SET revalidate_update_failed = $1
WHERE server.id = $2`
	if _, err := tx.Exec(query, failed, serverID); err != nil {
		return fmt.Errorf("setting reval update failed for ServerID %d with value %v: %w", serverID, failed, err)
	}
	return nil
}

// QueueRevalForServer sets the revalidate update time for the server to now.
func QueueRevalForServer(tx *sql.Tx, serverID int64) error {
	query := `
UPDATE public.server
SET revalidate_update_time = now()
WHERE server.id = $1;`

	if _, err := tx.Exec(query, serverID); err != nil {
		return fmt.Errorf("queueing reval update for ServerID %d: %w", serverID, err)
	}

	return nil
}

// SetApplyRevalForServer sets the revalidate apply time for the server to now.
func SetApplyRevalForServer(tx *sql.Tx, serverID int64) error {
	query := `
UPDATE public.server
SET revalidate_apply_time = now()
WHERE server.id = $1;`

	if _, err := tx.Exec(query, serverID); err != nil {
		return fmt.Errorf("queueing reval update for ServerID %d: %w", serverID, err)
	}

	return nil
}

// SetApplyRevalForServerWithTime sets the revalidate apply time for the
// server to the provided time.
func SetApplyRevalForServerWithTime(tx *sql.Tx, serverID int64, applyRevalTime time.Time) error {
	query := `
UPDATE public.server
SET revalidate_apply_time = $1
WHERE server.id = $2;`

	if _, err := tx.Exec(query, applyRevalTime, serverID); err != nil {
		return fmt.Errorf("applying config update for ServerID %d with time %v: %w", serverID, applyRevalTime, err)
	}

	return nil
}

// GetCommonServerPropertiesFromV4 converts ServerV40 to CommonServerProperties struct.
func GetCommonServerPropertiesFromV4(s tc.ServerV40, tx *sql.Tx) (tc.CommonServerProperties, error) {
	var id int
	var desc string
	if len(s.ProfileNames) == 0 {
		return tc.CommonServerProperties{}, fmt.Errorf("profileName doesnot exist in server: %v", *s.ID)
	}
	rows, err := tx.Query("SELECT id, description from profile WHERE name=$1", (s.ProfileNames)[0])
	if err != nil {
		return tc.CommonServerProperties{}, fmt.Errorf("querying profile id and description by profile_name: %w", err)
	}
	defer log.Close(rows, "closing rows in GetCommonServerPropertiesFromV4")

	for rows.Next() {
		if err := rows.Scan(&id, &desc); err != nil {
			return tc.CommonServerProperties{}, fmt.Errorf("scanning profile: %w", err)
		}
	}

	return tc.CommonServerProperties{
		Cachegroup:       s.Cachegroup,
		CachegroupID:     s.CachegroupID,
		CDNID:            s.CDNID,
		CDNName:          s.CDNName,
		DeliveryServices: s.DeliveryServices,
		DomainName:       s.DomainName,
		FQDN:             s.FQDN,
		FqdnTime:         s.FqdnTime,
		GUID:             s.GUID,
		HostName:         s.HostName,
		HTTPSPort:        s.HTTPSPort,
		ID:               s.ID,
		ILOIPAddress:     s.ILOIPAddress,
		ILOIPGateway:     s.ILOIPGateway,
		ILOIPNetmask:     s.ILOIPNetmask,
		ILOPassword:      s.ILOPassword,
		ILOUsername:      s.ILOUsername,
		LastUpdated:      s.LastUpdated,
		MgmtIPAddress:    s.MgmtIPAddress,
		MgmtIPGateway:    s.MgmtIPGateway,
		MgmtIPNetmask:    s.MgmtIPNetmask,
		OfflineReason:    s.OfflineReason,
		Profile:          &(s.ProfileNames)[0],
		ProfileDesc:      &desc,
		ProfileID:        &id,
		PhysLocation:     s.PhysLocation,
		PhysLocationID:   s.PhysLocationID,
		Rack:             s.Rack,
		RevalPending:     s.RevalPending,
		Status:           s.Status,
		StatusID:         s.StatusID,
		TCPPort:          s.TCPPort,
		Type:             s.Type,
		TypeID:           s.TypeID,
		UpdPending:       s.UpdPending,
		XMPPID:           s.XMPPID,
		XMPPPasswd:       s.XMPPPasswd,
	}, nil
}

// UpdateServerProfilesForV4 updates server_profile table via update function for APIv4.
func UpdateServerProfilesForV4(id int, profile []string, tx *sql.Tx) error {
	profileNames := make([]string, 0, len(profile))
	priority := make([]int, 0, len(profile))
	for i, _ := range profile {
		priority = append(priority, i)
	}

	//Delete existing rows from server_profile to get the priority correct for profile_name changes
	_, err := tx.Exec("DELETE FROM server_profile WHERE server=$1", id)
	if err != nil {
		return fmt.Errorf("updating server_profile by server id: %d, error: %w", id, err)
	}

	query := `WITH inserted AS (
		INSERT INTO server_profile
		SELECT $1, "profile_name", "priority"
		FROM UNNEST($2::text[], $3::int[]) AS tmp("profile_name", "priority")
		RETURNING profile_name, priority
	)
	SELECT ARRAY_AGG(profile_name)
	FROM (
		SELECT profile_name
		FROM inserted
		ORDER BY priority ASC
	) AS returned(profile_name)
`
	err = tx.QueryRow(query, id, pq.Array(profile), pq.Array(priority)).Scan(pq.Array(&profileNames))
	if err != nil {
		return fmt.Errorf("failed to insert/read into/from server_profile table, %w", err)
	}
	return nil
}

// UpdateServerProfileTableForV3 updates CommonServerPropertiesV40 struct and server_profile table via Update (server) function for API v3.
func UpdateServerProfileTableForV3(id *int, newProfileId *int, origProfile string, tx *sql.Tx) error {
	newProfile, _, err := GetProfileNameFromID(*newProfileId, tx)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("selecting profile by name: %w", err)
	}

	query := `UPDATE server_profile SET profile_name=$1 WHERE server=$2 AND profile_name=$3`
	_, err = tx.Exec(query, newProfile, *id, origProfile)
	if err != nil {
		return fmt.Errorf("updating server_profile by profile_name: %w", err)
	}
	return nil
}

// GetServerDetailFromV4 function converts server details from V4 to V3
func GetServerDetailFromV4(sd tc.ServerDetailV40, tx *sql.Tx) (tc.ServerDetail, error) {
	var profileDesc *string
	if err := tx.QueryRow(`SELECT p.description FROM profile p WHERE p.name=$1`, sd.ProfileNames[0]).Scan(&profileDesc); err != nil {
		return tc.ServerDetail{}, fmt.Errorf("querying profile description by profile name: %w", err)
	}
	return tc.ServerDetail{
		CacheGroup:         sd.CacheGroup,
		CDNName:            sd.CDNName,
		DeliveryServiceIDs: sd.DeliveryServiceIDs,
		DomainName:         sd.DomainName,
		GUID:               sd.GUID,
		HardwareInfo:       sd.HardwareInfo,
		HostName:           sd.HostName,
		HTTPSPort:          sd.HTTPSPort,
		ID:                 sd.ID,
		ILOIPAddress:       sd.ILOIPAddress,
		ILOIPGateway:       sd.ILOIPGateway,
		ILOIPNetmask:       sd.ILOIPNetmask,
		ILOPassword:        sd.ILOPassword,
		ILOUsername:        sd.ILOUsername,
		MgmtIPAddress:      sd.MgmtIPAddress,
		MgmtIPGateway:      sd.MgmtIPGateway,
		MgmtIPNetmask:      sd.MgmtIPNetmask,
		OfflineReason:      sd.OfflineReason,
		PhysLocation:       sd.PhysLocation,
		Profile:            &sd.ProfileNames[0],
		ProfileDesc:        profileDesc,
		Rack:               sd.Rack,
		Status:             sd.Status,
		TCPPort:            sd.TCPPort,
		Type:               sd.Type,
		XMPPID:             sd.XMPPID,
		XMPPPasswd:         sd.XMPPPasswd,
	}, nil
}

// GetProfileIDDesc gets profile ID and desc for V3 servers
func GetProfileIDDesc(tx *sql.Tx, name string) (id int, desc string) {
	err := tx.QueryRow(`SELECT id, description from "profile" p WHERE p.name=$1`, name).Scan(&id, &desc)
	if err != nil {
		log.Errorf("scanning id and description in GetProfileIDDesc: " + err.Error())
	}
	return
}

// GetSCInfo confirms whether the server capability exists, and an error (if one occurs).
func GetSCInfo(tx *sql.Tx, name string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(name) FROM server_capability AS sc WHERE sc.name=$1", name).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting server capability info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting server capability info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// ServiceCategoryExists confirms whether the service category exists, and an error (if one occurs).
func ServiceCategoryExists(tx *sql.Tx, name string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(name) FROM service_category AS sc WHERE sc.name=$1", name).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting service category info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting service category info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// ASNExists confirms whether the asn exists, and an error (if one occurs).
func ASNExists(tx *sql.Tx, id string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(asn) FROM asn WHERE id=$1", id).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting asn info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting asn info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// CacheGroupExists confirms whether the cache group exists, and an error (if one occurs).
func CacheGroupExists(tx *sql.Tx, name string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(name) FROM cachegroup AS cg WHERE cg.name=$1", name).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting cache group info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting cache group info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// DivisionExists confirms whether the division exists, and an error (if one occurs).
func DivisionExists(tx *sql.Tx, id string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(id) FROM division AS div WHERE div.id=$1", id).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting divisions info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting division info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// PhysLocationExists confirms whether the PhysLocation exists, and an error (if one occurs).
func PhysLocationExists(tx *sql.Tx, id string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(name) FROM phys_location WHERE id=$1", id).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting PhysLocation info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting PhysLocation info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// ParameterExists confirms whether the Parameter exists, and an error (if one occurs).
func ParameterExists(tx *sql.Tx, id string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(name) FROM parameter WHERE id=$1", id).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting Parameter info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting Parameter info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// ProfileExists confirms whether the profile exists, and an error (if one occurs).
func ProfileExists(tx *sql.Tx, id string) (bool, error) {
	var count int
	if err := tx.QueryRow(`SELECT count(name) FROM profile WHERE id=$1`, id).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting profile info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting profile info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}

// GetCoordinateID obtains coordinateID, and an error (if one occurs)
func GetCoordinateID(tx *sql.Tx, id int) (*int, error) {
	q := `SELECT coordinate FROM cachegroup WHERE id = $1`

	var coordinateID *int
	if err := tx.QueryRow(q, id).Scan(&coordinateID); err != nil {
		return nil, err
	}

	return coordinateID, nil
}

// DeleteCoordinate deletes coordinate by id, and an error (if one occurs)
func DeleteCoordinate(tx *sql.Tx, cacheGroupID int, coordinateID int) error {
	q := `UPDATE cachegroup SET coordinate = NULL WHERE id = $1`
	result, err := tx.Exec(q, cacheGroupID)
	if err != nil {
		return fmt.Errorf("updating cachegroup %d coordinate to null: %s", cacheGroupID, err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updating cachegroup %d coordinate to null, getting rows affected: %s", coordinateID, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updating cachegroup %d coordinate to null, zero rows affected", coordinateID)
	}

	q = `DELETE FROM coordinate WHERE id = $1`
	result, err = tx.Exec(q, coordinateID)
	if err != nil {
		return fmt.Errorf("delete coordinate %d for cachegroup %d: %s", coordinateID, coordinateID, err.Error())
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete coordinate %d for cachegroup %d, getting rows affected: %s", coordinateID, coordinateID, err.Error())
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete coordinate %d for cachegroup %d, zero rows affected", coordinateID, coordinateID)
	}
	return nil
}

// ProfileParameterExists confirms whether the ProfileParameter exists, and an error (if one occurs).
func ProfileParameterExists(tx *sql.Tx, profileID string, parameterID string) (bool, error) {
	var count int
	if err := tx.QueryRow("SELECT count(*) FROM profile_parameter WHERE profile=$1 and parameter=$2", profileID, parameterID).Scan(&count); err != nil {
		return false, fmt.Errorf("error getting profile_parameter info: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	if count != 1 {
		return false, fmt.Errorf("getting profile_parameter info - expected row count: 1, actual: %d", count)
	}
	return true, nil
}
