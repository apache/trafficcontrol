package comment

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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

// we need a type alias to define functions on
type TODeliveryServiceRequestComment struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceRequestCommentNullable
}

func (v *TODeliveryServiceRequestComment) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "deliveryservice_request_comment")
}

func (v *TODeliveryServiceRequestComment) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TODeliveryServiceRequestComment) InsertQuery() string           { return insertQuery() }
func (v *TODeliveryServiceRequestComment) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(dsrc.last_updated) as t from deliveryservice_request_comment dsrc
JOIN tm_user a ON dsrc.author_id = a.id
JOIN deliveryservice_request dsr ON dsrc.deliveryservice_request_id = dsr.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='deliveryservice_request_comment') as res`
}

func selectMaxLastUpdatedQuery(where, orderBy, pagination string) string {
	return `SELECT max(t) from (
		SELECT max(dsrc.last_updated) as t from deliveryservice_request_comment dsrc
JOIN tm_user a ON dsrc.author_id = a.id
JOIN deliveryservice_request dsr ON dsrc.deliveryservice_request_id = dsr.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='deliveryservice_request_comment') as res`
}

func (v *TODeliveryServiceRequestComment) NewReadObj() interface{} {
	return &tc.DeliveryServiceRequestCommentNullable{}
}
func (v *TODeliveryServiceRequestComment) SelectQuery() string { return selectQuery() }
func (v *TODeliveryServiceRequestComment) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"authorId":                 dbhelpers.WhereColumnInfo{Column: "dsrc.author_id"},
		"author":                   dbhelpers.WhereColumnInfo{Column: "a.username"},
		"deliveryServiceRequestId": dbhelpers.WhereColumnInfo{Column: "dsrc.deliveryservice_request_id"},
		"id":                       dbhelpers.WhereColumnInfo{Column: "dsrc.id", Checker: api.IsInt},
	}
}
func (v *TODeliveryServiceRequestComment) UpdateQuery() string { return updateQuery() }
func (v *TODeliveryServiceRequestComment) DeleteQuery() string { return deleteQuery() }

func (comment TODeliveryServiceRequestComment) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (comment TODeliveryServiceRequestComment) GetKeys() (map[string]interface{}, bool) {
	if comment.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *comment.ID}, true
}

func (comment *TODeliveryServiceRequestComment) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	comment.ID = &i
}

func (comment TODeliveryServiceRequestComment) GetAuditName() string {
	if comment.ID != nil {
		return strconv.Itoa(*comment.ID)
	}
	return "unknown"
}

func (comment TODeliveryServiceRequestComment) GetType() string {
	return "deliveryservice_request_comment"
}

func (comment TODeliveryServiceRequestComment) Validate() (error, error) {
	errs := validation.Errors{
		"deliveryServiceRequestId": validation.Validate(comment.DeliveryServiceRequestID, validation.NotNil),
		"value":                    validation.Validate(comment.Value, validation.NotNil),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (comment *TODeliveryServiceRequestComment) Create() (error, error, int) {
	au := tc.IDNoMod(comment.ReqInfo.User.ID)
	comment.AuthorID = &au
	return api.GenericCreate(comment)
}

func (comment *TODeliveryServiceRequestComment) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(comment.APIInfo(), "xmlId")
	return api.GenericRead(h, comment, useIMS)
}

func (comment *TODeliveryServiceRequestComment) Update(h http.Header) (error, error, int) {
	current := TODeliveryServiceRequestComment{}
	err := comment.ReqInfo.Tx.QueryRowx(selectQuery() + `WHERE dsrc.id=` + strconv.Itoa(*comment.ID)).StructScan(&current)
	if err != nil {
		return api.ParseDBError(err)
	}

	userID := tc.IDNoMod(comment.ReqInfo.User.ID)
	if *current.AuthorID != userID {
		return errors.New("Comments can only be updated by the author"), nil, http.StatusBadRequest
	}

	return api.GenericUpdate(h, comment)
}

func (comment *TODeliveryServiceRequestComment) Delete() (error, error, int) {
	var current TODeliveryServiceRequestComment
	err := comment.ReqInfo.Tx.QueryRowx(selectQuery() + `WHERE dsrc.id=` + strconv.Itoa(*comment.ID)).StructScan(&current)
	if err != nil {
		return nil, errors.New("querying DeliveryServiceRequestComments: " + err.Error()), http.StatusInternalServerError
	}

	if userID := tc.IDNoMod(comment.ReqInfo.User.ID); *current.AuthorID != userID {
		// TODO determine if users should be able to delete sub-tenant users' comments? Else, a deleted user's comments can never be removed.
		return errors.New("Comments can only be deleted by the author"), nil, http.StatusBadRequest
	}

	return api.GenericDelete(comment)
}

func insertQuery() string {
	query := `INSERT INTO deliveryservice_request_comment (
author_id,
deliveryservice_request_id,
value) VALUES (
:author_id,
:deliveryservice_request_id,
:value) RETURNING id,last_updated`
	return query
}

func selectQuery() string {
	query := `SELECT
a.username AS author,
dsrc.author_id,
dsrc.deliveryservice_request_id,
dsr.deliveryservice->>'xmlId' as xml_id,
dsrc.id,
dsrc.last_updated,
dsrc.value
FROM deliveryservice_request_comment dsrc
JOIN tm_user a ON dsrc.author_id = a.id
JOIN deliveryservice_request dsr ON dsrc.deliveryservice_request_id = dsr.id
`
	return query
}

func updateQuery() string {
	query := `UPDATE
deliveryservice_request_comment SET
deliveryservice_request_id=:deliveryservice_request_id,
value=:value
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	return `DELETE FROM deliveryservice_request_comment WHERE id = :id`
}

func Get(w http.ResponseWriter, r *http.Request) {
	var maxTime time.Time
	var runSecond bool
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cols := map[string]dbhelpers.WhereColumnInfo{
		"authorId":                 dbhelpers.WhereColumnInfo{Column: "dsrc.author_id"},
		"author":                   dbhelpers.WhereColumnInfo{Column: "a.username"},
		"deliveryServiceRequestId": dbhelpers.WhereColumnInfo{Column: "dsrc.deliveryservice_request_id"},
		"id":                       dbhelpers.WhereColumnInfo{Column: "dsrc.id", Checker: api.IsInt},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, cols)
	if len(errs) > 0 {
		errCode = http.StatusBadRequest
		userErr = util.JoinErrs(errs)
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	if inf.Config.UseIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where, orderBy, pagination))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.AddLastModifiedHdr(w, maxTime)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	deliveryServiceRequestComments := []tc.DeliveryServiceRequestCommentV5{}
	query := selectQuery() + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying cdn locks: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var deliveryServiceRequestComment tc.DeliveryServiceRequestCommentNullable
		if err = rows.StructScan(&deliveryServiceRequestComment); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning delivery service request comments: "+err.Error()))
			return
		}
		dsrcV5, err := deliveryServiceRequestComment.UpgradeFromNullable()
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("converting to deliveryServiceRequestCommentV5: "+err.Error()))
			return
		}
		deliveryServiceRequestComments = append(deliveryServiceRequestComments, dsrcV5)
	}
	api.WriteResp(w, r, deliveryServiceRequestComments)
}

func Validate(dsrc tc.DeliveryServiceRequestCommentV5) error {
	errs := validation.Errors{
		"deliveryServiceRequestId": validation.Validate(dsrc.DeliveryServiceRequestID, validation.NotNil),
		"value":                    validation.Validate(dsrc.Value, validation.NotNil),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func Update(w http.ResponseWriter, r *http.Request) {
	var deliveryServiceRequestComment tc.DeliveryServiceRequestCommentV5

	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	if err := json.NewDecoder(r.Body).Decode(&deliveryServiceRequestComment); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err := Validate(deliveryServiceRequestComment); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	var current TODeliveryServiceRequestComment
	err := inf.Tx.QueryRowx(selectQuery() + `WHERE dsrc.id=` + inf.Params["id"]).StructScan(&current)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning deliveryservice_request_comment: "+err.Error()))
		return
	}

	userID := tc.IDNoMod(inf.User.ID)
	if *current.AuthorID != userID {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("comments can only be updated by the author"), nil)
		return
	}

	idParam := inf.Params["id"]
	id, parseErr := strconv.Atoi(idParam)
	if parseErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("id must be an integer"), nil)
		return
	}
	deliveryServiceRequestComment.ID = id
	userErr, sysErr, sc := api.CheckIfUnModified(r.Header, inf.Tx, id, "deliveryservice_request_comment")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, sc, userErr, sysErr)
		return
	}

	rows, err := inf.Tx.NamedQuery(updateQuery(), deliveryServiceRequestComment)
	if err != nil {
		userErr, sysErr, sc := api.ParseDBError(err)
		api.HandleErr(w, r, tx, sc, userErr, sysErr)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no deliveryservice_request_comment found with this id"), nil)
		return
	}
	lastUpdated := tc.TimeNoMod{}
	if err := rows.Scan(&lastUpdated); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning lastUpdated from deliveryservice_request_comment insert: "+err.Error()))
		return
	}
	deliveryServiceRequestComment.LastUpdated = time.Now()

	t := lastUpdated.Time
	updatedTime, err := util.ConvertTimeFormat(t, time.RFC3339)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("converting last_updated to RFC3339 format: "+err.Error()))
		return
	}
	deliveryServiceRequestComment.LastUpdated = *updatedTime

	if rows.Next() {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("deliveryservice_request_comment update affected too many rows: >1"))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "deliveryservice_request_comment was updated.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, deliveryServiceRequestComment)
	changeLogMsg := fmt.Sprintf("DELIVERYSERVICE_REQUEST_COMMENT: %d, ID: %d, ACTION: Updated deliveryservice_request_comment", deliveryServiceRequestComment.ID, deliveryServiceRequestComment.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func Create(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	idParam := inf.Params["id"]
	id, parseErr := strconv.Atoi(idParam)
	if parseErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("id must be an integer"), nil)
		return
	}
	rows, err := tx.Query(deleteQuery(), id)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deleting deliveryservice_request_comment: %w", err))
		return
	}
	rows.Close()
	alerts := tc.CreateAlerts(tc.SuccessLevel, "deliveryservice_request_comment was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("DELIVERYSERVICE_REQUEST_COMMENT: %d, ACTION: Deleted Deliveryservice_request_comment", id)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
