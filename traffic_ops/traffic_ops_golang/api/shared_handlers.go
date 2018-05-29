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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

const PathParamsKey = "pathParams"

type KeyFieldInfo struct {
	Field string
	Func  func(string) (interface{}, error)
}

func GetIntKey(s string) (interface{}, error) {
	if strings.HasSuffix(s, ".json") {
		s = s[:len(s) - len(".json")]
	}
	return strconv.Atoi(s)
}

func GetStringKey(s string) (interface{}, error) {
	return s, nil
}

func GetPathParams(ctx context.Context) (map[string]string, error) {
	val := ctx.Value(PathParamsKey)
	if val != nil {
		switch v := val.(type) {
		case map[string]string:
			return v, nil
		default:
			return nil, fmt.Errorf("path parameters found with bad type: %T", v)
		}
	}
	return nil, errors.New("no PathParams found in Context")
}

func IsInt(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		err = errors.New("cannot parse to integer")
	}
	return err
}

func IsBool(s string) error {
	_, err := strconv.ParseBool(s)
	if err != nil {
		err = errors.New("cannot parse to boolean")
	}
	return err
}

func GetCombinedParams(r *http.Request) (map[string]string, error) {
	combinedParams := make(map[string]string)
	q := r.URL.Query()
	for k, v := range q {
		combinedParams[k] = v[0] //we take the first value and do not support multiple keys in query parameters
	}

	ctx := r.Context()
	pathParams, err := GetPathParams(ctx)
	if err != nil {
		return combinedParams, fmt.Errorf("no path parameters: %s", err)
	}
	//path parameters will overwrite query parameters
	for k, v := range pathParams {
		combinedParams[k] = v
	}

	return combinedParams, nil
}

//decodes and validates a pointer to a struct implementing the Validator interface
//      we lose the ability to unmarshal the struct if a struct implementing the interface is passed in,
//      because when when it is de-referenced it is a pointer to an interface. A new copy is created so that
//      there are no issues with concurrent goroutines
func decodeAndValidateRequestBody(r *http.Request, v Validator, db *sqlx.DB, user auth.CurrentUser) (interface{}, []error) {
	payload := reflect.Indirect(reflect.ValueOf(v)).Addr().Interface() // does a shallow copy v's internal struct members
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return nil, []error{err}
	}
	return payload, payload.(Validator).Validate(db)
}

//this creates a handler function from the pointer to a struct implementing the Reader interface
//      this handler retrieves the user from the context
//      combines the path and query parameters
//      produces the proper status code based on the error code returned
//      marshals the structs returned into the proper response json
func ReadHandler(typeRef Reader, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create error function with ResponseWriter and Request
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()

		// Load the PathParams into the query parameters for pass through
		params, err := GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		results, errs, errType := typeRef.Read(db, params, *user)
		if len(errs) > 0 {
			tc.HandleErrorsWithType(errs, errType, handleErrs)
			return
		}
		resp := struct {
			Response []interface{} `json:"response"`
		}{results}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

//this creates a handler function from the pointer to a struct implementing the Updater interface
//it must be immediately assigned to a local variable
//   this generic handler encapsulates the logic for handling:
//   *fetching the id from the path parameter
//   *current user
//   *decoding and validating the struct
//   *change log entry
//   *forming and writing the body over the wire
func UpdateHandler(typeRef Updater, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create error function with ResponseWriter and Request
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		//collect path parameters and user from context
		ctx := r.Context()
		params, err := GetCombinedParams(r)
		if err != nil {
			log.Errorf("received error trying to get path parameters: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		//create local instance of the shared typeRef pointer
		//no operations should be made on the typeRef
		//decode the body and validate the request struct
		decoded, errs := decodeAndValidateRequestBody(r, typeRef, db, *user)
		if len(errs) > 0 {
			handleErrs(http.StatusBadRequest, errs...)
			return
		}
		u := decoded.(Updater)
		//now we have a validated local object to update

		keyFields := u.GetKeyFieldsInfo() //expecting a slice of the key fields info which is a struct with the field name and a function to convert a string into a {}interface of the right type. in most that will be [{Field:"id",Func: func(s string)({}interface,error){return strconv.Atoi(s)}}]
		keys, ok := u.GetKeys()           // a map of keyField to keyValue where keyValue is an {}interface
		if !ok {
			log.Errorf("unable to parse keys from request: %++v", u)
			handleErrs(http.StatusBadRequest, errors.New("unable to parse required keys from request body"))
			return // TODO verify?
		}
		for _, keyFieldInfo := range keyFields {
			paramKey := params[keyFieldInfo.Field]
			if paramKey == "" {
				log.Errorf("missing key: %s", keyFieldInfo.Field)
				handleErrs(http.StatusBadRequest, errors.New("missing key: "+keyFieldInfo.Field))
				return
			}

			paramValue, err := keyFieldInfo.Func(paramKey)
			if err != nil {
				log.Errorf("failed to parse key %s: %s", keyFieldInfo.Field, err)
				handleErrs(http.StatusBadRequest, errors.New("failed to parse key: "+keyFieldInfo.Field))
				return
			}

			if paramValue != keys[keyFieldInfo.Field] {
				handleErrs(http.StatusBadRequest, errors.New("key in body does not match key in params"))
				return
			}
		}

		// if the object has tenancy enabled, check that user is able to access the tenant
		if t, ok := u.(Tenantable); ok {
			authorized, err := t.IsTenantAuthorized(*user, db)
			if err != nil {
				handleErrs(http.StatusBadRequest, err)
				return
			}
			if !authorized {
				handleErrs(http.StatusForbidden, errors.New("not authorized on this tenant"))
				return
			}
		}

		//run the update and handle any error
		err, errType := u.Update(db, *user)
		if err != nil {
			tc.HandleErrorsWithType([]error{err}, errType, handleErrs)
			return
		}
		//auditing here
		CreateChangeLog(ApiChange, Updated, u, *user, db)
		//form response to send across the wire
		resp := struct {
			Response interface{} `json:"response"`
			tc.Alerts
		}{u, tc.CreateAlerts(tc.SuccessLevel, u.GetType()+" was updated.")}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

//this creates a handler function from the pointer to a struct implementing the Deleter interface
//it must be immediately assigned to a local variable
//   this generic handler encapsulates the logic for handling:
//   *fetching the id from the path parameter
//   *current user
//   *change log entry
//   *forming and writing the body over the wire
func DeleteHandler(typeRef Deleter, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		d := typeRef

		ctx := r.Context()
		params, err := GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		keyFields := d.GetKeyFieldsInfo() // expecting a slice of the key fields info which is a struct with the field name and a function to convert a string into a interface{} of the right type. in most that will be [{Field:"id",Func: func(s string)(interface{},error){return strconv.Atoi(s)}}]
		keys := make(map[string]interface{})
		for _, keyFieldInfo := range keyFields {
			paramKey := params[keyFieldInfo.Field]
			if paramKey == "" {
				log.Errorf("missing key: %s", keyFieldInfo.Field)
				handleErrs(http.StatusBadRequest, errors.New("missing key: "+keyFieldInfo.Field))
				return
			}

			paramValue, err := keyFieldInfo.Func(paramKey)
			if err != nil {
				log.Errorf("failed to parse key %s: %s", keyFieldInfo.Field, err)
				handleErrs(http.StatusBadRequest, errors.New("failed to parse key: "+keyFieldInfo.Field))
			}
			keys[keyFieldInfo.Field] = paramValue
		}
		d.SetKeys(keys) // if the type assertion of a key fails it will be should be set to the zero value of the type and the delete should fail (this means the code is not written properly no changes of user input should cause this.)

		// if the object has tenancy enabled, check that user is able to access the tenant
		if t, ok := d.(Tenantable); ok {
			authorized, err := t.IsTenantAuthorized(*user, db)
			if err != nil {
				handleErrs(http.StatusBadRequest, err)
				return
			}
			if !authorized {
				handleErrs(http.StatusForbidden, errors.New("not authorized on this tenant"))
				return
			}
		}

		log.Debugf("calling delete on object: %++v", d) //should have id set now
		err, errType := d.Delete(db, *user)
		if err != nil {
			log.Errorf("error deleting: %++v", err)
			tc.HandleErrorsWithType([]error{err}, errType, handleErrs)
			return
		}
		//audit here
		log.Debugf("changelog for delete on object")
		CreateChangeLog(ApiChange, Deleted, d, *user, db)
		//
		resp := struct {
			tc.Alerts
		}{tc.CreateAlerts(tc.SuccessLevel, d.GetType()+" was deleted.")}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

//this creates a handler function from the pointer to a struct implementing the Creator interface
//it must be immediately assigned to a local variable
//   this generic handler encapsulates the logic for handling:
//   *fetching the id from the path parameter
//   *current user
//   *decoding and validating the struct
//   *change log entry
//   *forming and writing the body over the wire
func CreateHandler(typeRef Creator, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		//decode the body and validate the request struct
		decoded, errs := decodeAndValidateRequestBody(r, typeRef, db, *user)
		if len(errs) > 0 {
			handleErrs(http.StatusBadRequest, errs...)
			return
		}
		i := decoded.(Creator)
		log.Debugf("%++v", i)
		//now we have a validated local object to insert

		// if the object has tenancy enabled, check that user is able to access the tenant
		if t, ok := i.(Tenantable); ok {
			authorized, err := t.IsTenantAuthorized(*user, db)
			if err != nil {
				handleErrs(http.StatusBadRequest, err)
				return
			}
			if !authorized {
				handleErrs(http.StatusForbidden, errors.New("not authorized on this tenant"))
				return
			}
		}

		err, errType := i.Create(db, *user)
		if err != nil {
			tc.HandleErrorsWithType([]error{err}, errType, handleErrs)
			return
		}

		CreateChangeLog(ApiChange, Created, i, *user, db)

		resp := struct {
			Response interface{} `json:"response"`
			tc.Alerts
		}{i, tc.CreateAlerts(tc.SuccessLevel, i.GetType()+" was created.")}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

// WriteResp takes any object, serializes it as JSON, and writes that to w. Any errors are logged and written to w via tc.GetHandleErrorsFunc.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func WriteResp(w http.ResponseWriter, r *http.Request, v interface{}) {
	resp := struct {
		Response interface{} `json:"response"`
	}{v}
	respBts, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("marshalling JSON for %T: %v", v, err)
		tc.GetHandleErrorsFunc(w, r)(http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBts)
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

// IntParams parses integer parameters, and returns map of the given params, or an error. This guarantees if error is nil, all requested parameters successfully parsed and exist in the returned map, hence if error is nil there's no need to check for existence. The intParams may be nil if no integer parameters are required.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func IntParams(params map[string]string, intParamNames []string) (map[string]int, error) {
	intParams := map[string]int{}
	for _, intParam := range intParamNames {
		realParam, ok := params[intParam]
		if !ok {
			return nil, errors.New("missing required integer parameter '" + intParam + "'")
		}
		intVal, err := strconv.Atoi(realParam)
		if err != nil {
			return nil, errors.New("required parameter '" + intParam + "'" + " not an integer")
		}
		intParams[intParam] = intVal
	}
	return intParams, nil
}

// AllParams takes the request (in which the router has inserted context for path parameters), and an array of parameters required to be integers, and returns the map of combined parameters, and the map of int parameters; or a user or system error and the HTTP error code. The intParams may be nil if no integer parameters are required.
// This is a helper for the common case; not using this in unusual cases is perfectly acceptable.
func AllParams(req *http.Request, intParamNames []string) (map[string]string, map[string]int, error, error, int) {
	params, err := GetCombinedParams(req)
	if err != nil {
		return nil, nil, errors.New("getting combined URI parameters: " + err.Error()), nil, http.StatusBadRequest
	}
	intParams, err := IntParams(params, intParamNames)
	if err != nil {
		return nil, nil, nil, errors.New("getting combined URI parameters: " + err.Error()), http.StatusInternalServerError
	}
	return params, intParams, nil, nil, 0
}
