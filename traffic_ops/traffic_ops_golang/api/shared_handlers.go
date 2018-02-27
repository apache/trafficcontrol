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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

const PathParamsKey = "pathParams"

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
func decodeAndValidateRequestBody(r *http.Request, v Validator, db *sqlx.DB) (interface{}, []error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	payload := reflect.New(typ).Interface()
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
		//create local instance of the shared typeRef pointer
		//no operations should be made on the typeRef
		//decode the body and validate the request struct
		decoded, errs := decodeAndValidateRequestBody(r, typeRef, db)
		if len(errs) > 0 {
			handleErrs(http.StatusBadRequest, errs...)
			return
		}
		u := decoded.(Updater)
		//now we have a validated local object to update

		//collect path parameters and user from context
		ctx := r.Context()
		pathParams, err := GetPathParams(ctx)
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
		id, err := strconv.Atoi(pathParams["id"])
		if err != nil {
			log.Errorf("received error trying to convert id path parameter: %s", err)
			handleErrs(http.StatusBadRequest, errors.New("id from path not parseable as int"))
			return
		}

		iid, ok := u.GetID()
		if !ok || iid != id {
			handleErrs(http.StatusBadRequest, errors.New("id in body does not match id in path"))
			return
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
		pathParams, err := GetPathParams(ctx)
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

		id, err := strconv.Atoi(pathParams["id"])
		if err != nil {
			handleErrs(http.StatusBadRequest, errors.New("id from path not parseable as int"))
			return
		}
		d.SetID(id)

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
			tc.HandleErrorsWithType([]error{err}, errType, handleErrs)
			return
		}
		//audit here
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

		//decode the body and validate the request struct
		decoded, errs := decodeAndValidateRequestBody(r, typeRef, db)
		if len(errs) > 0 {
			handleErrs(http.StatusBadRequest, errs...)
			return
		}
		i := decoded.(Creator)
		log.Debugf("%++v", i)
		//now we have a validated local object to insert

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}

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
