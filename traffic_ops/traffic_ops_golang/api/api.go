package api

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
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/jmoiron/sqlx"
)

const DBContextKey = "db"
const ConfigContextKey = "context"
const ReqIDContextKey = "reqid"

type CRUDFactory func(reqInfo *APIInfo) CRUDer

// WriteResp takes any object, serializes it as JSON, and writes that to w. Any errors are logged and written to w via tc.GetHandleErrorsFunc.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func WriteResp(w http.ResponseWriter, r *http.Request, v interface{}) {
	resp := struct {
		Response interface{} `json:"response"`
	}{v}
	WriteRespRaw(w, r, resp)
}

// WriteRespRaw acts like WriteResp, but doesn't wrap the object in a `{"response":` object. This should be used to respond with endpoints which don't wrap their response in a "response" object.
func WriteRespRaw(w http.ResponseWriter, r *http.Request, v interface{}) {
	bts, err := json.Marshal(v)
	if err != nil {
		log.Errorf("marshalling JSON (raw) for %T: %v", v, err)
		tc.GetHandleErrorsFunc(w, r)(http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bts)
}

// WriteRespVals is like WriteResp, but also takes a map of root-level values to write. The API most commonly needs these for meta-parameters, like size, limit, and orderby.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func WriteRespVals(w http.ResponseWriter, r *http.Request, v interface{}, vals map[string]interface{}) {
	vals["response"] = v
	respBts, err := json.Marshal(vals)
	if err != nil {
		log.Errorf("marshalling JSON for %T: %v", v, err)
		tc.GetHandleErrorsFunc(w, r)(http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBts)
}

// HandleErr handles an API error, writing the given statusCode and userErr to the user, and logging the sysErr. If userErr is nil, the text of the HTTP statusCode is written.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func HandleErr(w http.ResponseWriter, r *http.Request, statusCode int, userErr error, sysErr error) {
	if sysErr != nil {
		log.Errorln(r.RemoteAddr + " " + sysErr.Error())
	}
	if userErr == nil {
		userErr = errors.New(http.StatusText(statusCode))
	}
	respBts, err := json.Marshal(tc.CreateErrorAlerts(userErr))
	if err != nil {
		log.Errorln("marshalling error: " + err.Error())
		*r = *r.WithContext(context.WithValue(r.Context(), tc.StatusKey, http.StatusInternalServerError))
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	*r = *r.WithContext(context.WithValue(r.Context(), tc.StatusKey, statusCode))
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write(respBts)
}

// RespWriter is a helper to allow a one-line response, for endpoints with a function that returns the object that needs to be written and an error.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func RespWriter(w http.ResponseWriter, r *http.Request) func(v interface{}, err error) {
	return func(v interface{}, err error) {
		if err != nil {
			HandleErr(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		WriteResp(w, r, v)
	}
}

// RespWriterVals is like RespWriter, but also takes a map of root-level values to write. The API most commonly needs these for meta-parameters, like size, limit, and orderby.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func RespWriterVals(w http.ResponseWriter, r *http.Request, vals map[string]interface{}) func(v interface{}, err error) {
	return func(v interface{}, err error) {
		if err != nil {
			HandleErr(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		WriteRespVals(w, r, v, vals)
	}
}

// WriteRespAlert creates an alert, serializes it as JSON, and writes that to w. Any errors are logged and written to w via tc.GetHandleErrorsFunc.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func WriteRespAlert(w http.ResponseWriter, r *http.Request, level tc.AlertLevel, msg string) {
	resp := struct{ tc.Alerts }{tc.CreateAlerts(level, msg)}
	respBts, err := json.Marshal(resp)
	if err != nil {
		HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("marshalling JSON: "+err.Error()))
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write(respBts)
}

// WriteRespAlertObj Writes the given alert, and the given response object.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func WriteRespAlertObj(w http.ResponseWriter, r *http.Request, level tc.AlertLevel, msg string, obj interface{}) {
	resp := struct {
		tc.Alerts
		Response interface{} `json:"response"`
	}{
		Alerts:   tc.CreateAlerts(level, msg),
		Response: obj,
	}
	respBts, err := json.Marshal(resp)
	if err != nil {
		HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("marshalling JSON: "+err.Error()))
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write(respBts)
}

// IntParams parses integer parameters, and returns map of the given params, or an error if any integer param is not an integer. The intParams may be nil if no integer parameters are required. Note this does not check existence; if an integer paramter is required, it should be included in the requiredParams given to NewInfo.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func IntParams(params map[string]string, intParamNames []string) (map[string]int, error) {
	intParams := map[string]int{}
	for _, intParam := range intParamNames {
		valStr, ok := params[intParam]
		if !ok {
			continue
		}
		valInt, err := strconv.Atoi(valStr)
		if err != nil {
			return nil, errors.New("parameter '" + intParam + "'" + " not an integer")
		}
		intParams[intParam] = valInt
	}
	return intParams, nil
}

// ParamsHaveRequired checks that params have all the required parameters, and returns nil on success, or an error providing information on which params are missing.
func ParamsHaveRequired(params map[string]string, required []string) error {
	missing := []string{}
	for _, requiredParam := range required {
		if _, ok := params[requiredParam]; !ok {
			missing = append(missing, requiredParam)
		}
	}
	if len(missing) > 0 {
		return errors.New("missing required parameters: " + strings.Join(missing, ", "))
	}
	return nil
}

// StripParamJSON removes ".json" trailing any parameter value, and returns the modified params.
// This allows the API handlers to transparently accept /id.json routes, as allowed by the 1.x API.
func StripParamJSON(params map[string]string) map[string]string {
	for name, val := range params {
		if strings.HasSuffix(val, ".json") {
			params[name] = val[:len(val)-len(".json")]
		}
	}
	return params
}

// AllParams takes the request (in which the router has inserted context for path parameters), and an array of parameters required to be integers, and returns the map of combined parameters, and the map of int parameters; or a user or system error and the HTTP error code. The intParams may be nil if no integer parameters are required.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func AllParams(req *http.Request, required []string, ints []string) (map[string]string, map[string]int, error, error, int) {
	params, err := GetCombinedParams(req)
	if err != nil {
		return nil, nil, errors.New("getting combined URI parameters: " + err.Error()), nil, http.StatusBadRequest
	}
	params = StripParamJSON(params)
	if err := ParamsHaveRequired(params, required); err != nil {
		return nil, nil, errors.New("required parameters missing: " + err.Error()), nil, http.StatusBadRequest
	}
	intParams, err := IntParams(params, ints)
	if err != nil {
		return nil, nil, nil, errors.New("getting integer parameters: " + err.Error()), http.StatusInternalServerError
	}
	return params, intParams, nil, nil, 0
}

type ParseValidator interface {
	Validate(tx *sql.Tx) []error
}

// Decode decodes a JSON object from r into the given v, validating and sanitizing the input. This helper should be used in API endpoints, rather than the json package, to safely decode and validate PUT and POST requests.
// TODO change to take data loaded from db, to remove sql from tc package.
func Parse(r io.Reader, tx *sql.Tx, v ParseValidator) error {
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return errors.New("decoding: " + err.Error())
	}
	if errs := v.Validate(tx); len(errs) > 0 {
		return errors.New("validating: " + util.JoinErrs(errs).Error())
	}
	return nil
}

type APIInfo struct {
	Params    map[string]string
	IntParams map[string]int
	User      *auth.CurrentUser
	ReqID     uint64
	Tx        *sqlx.Tx
	CommitTx  *bool
	Config    *config.Config
}

// NewInfo get and returns the context info needed by handlers. It also returns any user error, any system error, and the status code which should be returned to the client if an error occurred.
// Close() must be called to free resources, and should be called in a defer immediately after NewInfo(), to commit or rollback the transaction.
//
// Example:
//  func handler(w http.ResponseWriter, r *http.Request) {
//    inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
//    if userErr != nil || sysErr != nil {
//      api.HandleErr(w, r, errCode, userErr, sysErr)
//      return
//    }
//    defer inf.Close()
//
//    ...
//
//    err := finalDatabaseOperation(inf.Tx)
//    if err == nil {
//      *inf.CommitTx = true
//    }
//
func NewInfo(r *http.Request, requiredParams []string, intParamNames []string) (*APIInfo, error, error, int) {
	db, err := getDB(r.Context())
	if err != nil {
		return nil, errors.New("getting db: " + err.Error()), nil, http.StatusInternalServerError
	}
	cfg, err := getConfig(r.Context())
	if err != nil {
		return nil, errors.New("getting config: " + err.Error()), nil, http.StatusInternalServerError
	}
	reqID, err := getReqID(r.Context())
	if err != nil {
		return nil, errors.New("getting reqID: " + err.Error()), nil, http.StatusInternalServerError
	}

	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		return nil, errors.New("getting user: " + err.Error()), nil, http.StatusInternalServerError
	}
	params, intParams, userErr, sysErr, errCode := AllParams(r, requiredParams, intParamNames)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode
	}
	tx, err := db.Beginx() // must be last, MUST not return an error if this suceeds, without closing the tx
	if err != nil {
		return nil, userErr, errors.New("could not begin transaction: " + err.Error()), http.StatusInternalServerError
	}
	return &APIInfo{
		Config:    cfg,
		ReqID:     reqID,
		Params:    params,
		IntParams: intParams,
		User:      user,
		Tx:        tx,
		CommitTx:  util.BoolPtr(false),
	}, nil, nil, http.StatusOK
}

// Close implements the io.Closer interface. It should be called in a defer immediately after NewInfo().
//
// Close will commit or rollback the transaction, depending whether *info.CommitTx is true.
func (inf *APIInfo) Close() {
	dbhelpers.FinishTxX(inf.Tx, inf.CommitTx)
}

func getDB(ctx context.Context) (*sqlx.DB, error) {
	val := ctx.Value(DBContextKey)
	if val != nil {
		switch v := val.(type) {
		case *sqlx.DB:
			return v, nil
		default:
			return nil, fmt.Errorf("Tx found with bad type: %T", v)
		}
	}
	return nil, errors.New("No db found in Context")
}

func getConfig(ctx context.Context) (*config.Config, error) {
	val := ctx.Value(ConfigContextKey)
	if val != nil {
		switch v := val.(type) {
		case *config.Config:
			return v, nil
		default:
			return nil, fmt.Errorf("Config found with bad type: %T", v)
		}
	}
	return nil, errors.New("No config found in Context")
}

func getReqID(ctx context.Context) (uint64, error) {
	val := ctx.Value(ReqIDContextKey)
	if val != nil {
		switch v := val.(type) {
		case uint64:
			return v, nil
		default:
			return 0, fmt.Errorf("ReqID found with bad type: %T", v)
		}
	}
	return 0, errors.New("No ReqID found in Context")
}
