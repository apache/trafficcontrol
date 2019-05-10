package cacheassignmentgroup
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
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
)


/*******
Temporary Docs
----------------

POST,GET /api/1.4/cacheassignmentgroups/
PUT /api/1.4/cacheassignmentgroups/
{"id": 1,
"name": "name1",
"description": "description1",
"cdnId": 1,
"servers": [1,2,...n],
"lastUpdated": "",
}


Read includes DB contents of cacheassignmentgrouptable AND server assignments
Read can be filtered by CAG ID, Server ID, CDN Id or Name

PUT, DELETE methods REQUIRE and id parameter

**************/

type TOCacheAssignmentGroup struct {
	api.APIInfoImpl `json:"-"`
	tc.CacheAssignmentGroupNullable
}

func (cag *TOCacheAssignmentGroup) SetLastUpdated(t tc.TimeNoMod) { cag.LastUpdated = &t }

func (cag *TOCacheAssignmentGroup) Validate() error {
	errs := validation.Errors{
		"name":   		validation.Validate(cag.Name, validation.Required),
		"description":	validation.Validate(cag.Description, validation.Required),
		"cdnId": 	    validation.Validate(cag.CDNID, validation.Required, validation.Min(0)),
		"servers":      validation.Validate(cag.Servers, validation.Required),
	}
	if errs != nil {
		return util.JoinErrs(tovalidate.ToErrors(errs))
	}
	return nil
}

func (cag *TOCacheAssignmentGroup) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":     dbhelpers.WhereColumnInfo{"cag.id", api.IsInt},
		"name":   dbhelpers.WhereColumnInfo{"cag.name", nil},
		"cdnId":  dbhelpers.WhereColumnInfo{"cag.cdn_id", api.IsInt},
		// filtering by server is handled below, this also means we cannot orderby server
	}
}

// Lots of copy and paste here from generic_crud.GenericRead(). Would like to refactor to use GenericRead() but there
//  doesn't seem to be a good way to add custom WHERE and GROUP BY clauses without affecting all other GenericReaders
func (cag *TOCacheAssignmentGroup) Read() ([]interface{}, error, error, int) {
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(cag.APIInfo().Params, cag.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}

	if serverId, ok := cag.APIInfo().Params["server"]; ok {
		if err := api.IsInt(serverId); err != nil {
			return nil, errors.New("server "+err.Error()), nil, http.StatusBadRequest
		}
		queryValues["server"] = serverId

		serverFilter := dbhelpers.BaseWhere
		if len(where) > 0 {
			serverFilter = " AND"
		}

		serverFilter += " cag.id in (SELECT cacheassignmentgroup FROM cacheassignmentgroup_server WHERE server = :server) "
		where += serverFilter
		log.Debugln("Updated Where clause")
		log.Debugln(where)
	}

	selectQuery, groupBy := cag.SelectQuery()

	query :=  selectQuery + where + groupBy + orderBy
	log.Debugln(query)
	rows, err := cag.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("querying " + cag.GetType() + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	vals := []interface{}{}
	for rows.Next() {
		v := &tc.CacheAssignmentGroupNullable{}
		if err = rows.StructScan(v); err != nil {
			return nil, nil, errors.New("scanning " + cag.GetType() + ": " + err.Error()), http.StatusInternalServerError
		}
		vals = append(vals, v)
	}
	return vals, nil, nil, http.StatusOK
}

func (cag *TOCacheAssignmentGroup) SelectQuery() (string, string) {
	selectQuery := `SELECT 
       cag.id AS id,
       name,
       description,
       cag.cdn_id AS cdn_id,
       cag.last_updated AS last_updated,
       ARRAY_REMOVE(ARRAY_AGG(server.id), NULL) AS servers
FROM   cacheassignmentgroup cag
       LEFT JOIN cacheassignmentgroup_server AS cag_s
               ON cag_s.cacheassignmentgroup = cag.id
       LEFT OUTER JOIN server
               ON cag_s.server = server.id
`

	groupBy := " GROUP BY cag.id"
	return selectQuery, groupBy
 }


func (cag *TOCacheAssignmentGroup) Update() (error, error, int) {
	usrErr, sysErr, errCode := api.GenericUpdate(cag)
	if usrErr != nil || sysErr != nil {
		return usrErr, sysErr, errCode
	}

	err := cag.DeleteServerAssignments()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if len(cag.Servers) > 0 {
		usrErr, sysErr, errCode = cag.UpdateServerAssignments()
		if usrErr != nil || sysErr != nil {
			return usrErr, sysErr, errCode
		}
	}

	return nil, nil, http.StatusOK
}

func (cag *TOCacheAssignmentGroup) UpdateQuery() string {
	return `UPDATE
cacheassignmentgroup SET
name=:name,
description=:description,
cdn_id=:cdn_id
WHERE id=:id RETURNING last_updated`
}

func (cag *TOCacheAssignmentGroup) DeleteServerAssignments() error {
	deleteQuery := "DELETE FROM cacheassignmentgroup_server where cacheassignmentgroup = $1"
	_, err := cag.APIInfo().Tx.Exec(deleteQuery, *cag.ID)
	if (err != nil) {
		return err
	}

	return nil
}

func (cag *TOCacheAssignmentGroup) UpdateServerAssignments() (error, error, int) {
	serverInsertQuery, params := cag.BuildServerInsertQuery()
	log.Debugln(serverInsertQuery)
	log.Debugln(params)

	// Because Go typing is draconian...
	paramInterface := make([]interface{}, len(params))
	for i, v := range params {
		paramInterface[i] = v
	}

	result, err := cag.APIInfo().Tx.Exec(serverInsertQuery, paramInterface...)
	if err != nil {
		return api.ParseDBError(err)
	}

	rowsAffected, err := result.RowsAffected();
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if rowsAffected != int64(len(cag.Servers)) {
		return nil, errors.New("wrong number of cag_server rows updated in " + cag.GetType()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func (cag *TOCacheAssignmentGroup) BuildServerInsertQuery() (string, []int) {
	insertQuery := "INSERT INTO cacheassignmentgroup_server (cacheassignmentgroup, server) VALUES "
	var params []int

	paramIdx := 1
	for idx, server := range cag.Servers {
		if idx > 0 	{
			insertQuery += ", \n"
		}

		insertQuery +=  "($"+ strconv.Itoa(paramIdx) +", $"+ strconv.Itoa(paramIdx+1) +")"
		paramIdx +=2

		params = append(params, *cag.ID)
		params = append(params, int(server))
	}

	return insertQuery, params
}


func (cag *TOCacheAssignmentGroup) Create() (error, error, int) {
	usrErr, sysErr, errCode := api.GenericCreate(cag)
	if usrErr != nil || sysErr != nil {
		return usrErr, sysErr, errCode
	}

	if len(cag.Servers) > 0 {
		usrErr, sysErr, errCode = cag.UpdateServerAssignments()
		if usrErr != nil || sysErr != nil {
			return usrErr, sysErr, errCode
		}
	}

	return nil, nil, http.StatusOK
}

func (cag *TOCacheAssignmentGroup) InsertQuery() string {
	return `INSERT INTO cacheassignmentgroup (
name,
description,
cdn_id) VALUES (
:name,
:description,
:cdn_id) RETURNING id,last_updated`
}


func (cag *TOCacheAssignmentGroup) Delete() (error, error, int)              { return api.GenericDelete(cag) }

func (cag *TOCacheAssignmentGroup)  DeleteQuery() string {
	return `DELETE FROM cacheassignmentgroup
WHERE id=:id`
}


// Boilerplate for Identifier interface
func (cag TOCacheAssignmentGroup) GetKeys() (map[string]interface{}, bool) {
	if cag.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cag.ID}, true
}

func (cag *TOCacheAssignmentGroup)  SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cag.ID = &i
}
func (cag *TOCacheAssignmentGroup) GetType() string { return "cacheassignmentgroup" }
func (cag *TOCacheAssignmentGroup) GetAuditName() string {
	if cag.Name != nil {
		return *cag.Name
	}
	if cag.ID != nil {
		return strconv.Itoa(*cag.ID)
	}
	return "unknown"
}

func (cag TOCacheAssignmentGroup) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}