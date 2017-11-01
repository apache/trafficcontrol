package main

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
	"fmt"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

// Http POST handler used to store urisigning keys to a delivery service.
func assignDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		defer r.Body.Close()

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		xmlID := pathParams["xmlID"]
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// validate that the received data is a valid jwk keyset
		var keySet map[string]URISignerKeyset
		if err := json.Unmarshal(data, &keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}
		if err := validateURIKeyset(keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}

		// create and start a cluster
		cluster, err := getRiakCluster(db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		ro, err := fetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// object exists.
		if ro != nil && ro[0].Value != nil {
			handleErr(fmt.Errorf("a keyset already exists for this delivery service"), http.StatusBadRequest)
			return
		}

		// create a storage object and store the data
		obj := &riak.Object{
			ContentType:     "text/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             xmlID,
			Value:           []byte(data),
		}

		err = saveObject(obj, CDNURIKeysBucket, cluster)
		if err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", data)
	}
}

// endpoint handler for fetching uri signing keys from riak
func getURIsignkeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		xmlID := pathParams["xmlID"]

		// create and start a cluster
		cluster, err := getRiakCluster(db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		ro, err := fetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		var respBytes []byte

		if ro == nil {
			var empty URISignerKeyset
			respBytes, err = json.Marshal(empty)
			if err != nil {
				log.Errorf("failed to marshal an empty response: %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
				return
			}
		} else {
			respBytes = ro[0].Value
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBytes)
	}
}

// Http DELETE handler used to remove urisigning keys assigned to a delivery service.
func removeDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		xmlID := pathParams["xmlID"]

		// create and start a cluster
		cluster, err := getRiakCluster(db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		ro, err := fetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// fetch the object and delete it if it exists.
		var alert tc.Alerts

		if ro == nil || ro[0].Value == nil {
			alert = tc.CreateAlerts(tc.InfoLevel, "not deleted, no object found to delete.")
		} else if err := deleteObject(xmlID, CDNURIKeysBucket, cluster); err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		} else { // object successfully deleted
			alert = tc.CreateAlerts(tc.SuccessLevel, "object deleted")
		}

		// send response
		respBytes, err := json.Marshal(alert)
		if err != nil {
			log.Errorf("failed to marshal an alert response: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBytes)
	}
}

// Http POST handler used to store urisigning keys to a delivery service.
func updateDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		defer r.Body.Close()

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		xmlID := pathParams["xmlID"]
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// validate that the received data is a valid jwk keyset
		var keySet map[string]URISignerKeyset
		if err := json.Unmarshal(data, &keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}
		if err := validateURIKeyset(keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}

		// create and start a cluster
		cluster, err := getRiakCluster(db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		// create a storage object and store the data
		obj := &riak.Object{
			ContentType:     "text/json",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Key:             xmlID,
			Value:           []byte(data),
		}

		err = saveObject(obj, CDNURIKeysBucket, cluster)
		if err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", data)
	}
}
