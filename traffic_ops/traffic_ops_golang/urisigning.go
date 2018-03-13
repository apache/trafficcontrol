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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat/go-jwx/jwk"
)

// CDNURIKeysBucket is the namespace or bucket used for CDN URI signing keys.
const CDNURIKeysBucket = "cdn_uri_sig_keys"

// URISignerKeyset is the container for the CDN URI signing keys
type URISignerKeyset struct {
	RenewalKid *string               `json:"renewal_kid"`
	Keys       []jwk.EssentialHeader `json:"keys"`
}

// endpoint handler for fetching uri signing keys from riak
func getURIsignkeysHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)

		if cfg.RiakEnabled == false {
			handleErr(http.StatusServiceUnavailable, fmt.Errorf("The RIAK service is unavailable"))
			return
		}

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		xmlID := pathParams["xmlID"]

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(*user, xmlID, db)
		if !hasAccess {
			switch apiStatus {
			case tc.SystemError:
				handleErr(http.StatusInternalServerError, err)
				return
			case tc.DataMissingError:
				handleErr(http.StatusBadRequest, err)
				return
			case tc.ForbiddenError:
				handleErr(http.StatusForbidden, err)
				return
			}
		}

		// create and start a cluster
		cluster, err := riaksvc.GetRiakCluster(db, cfg.RiakAuthOptions)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		ro, err := riaksvc.FetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
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
func removeDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)

		if cfg.RiakEnabled == false {
			handleErr(http.StatusServiceUnavailable, fmt.Errorf("The RIAK service is unavailable"))
			return
		}

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		xmlID := pathParams["xmlID"]

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(*user, xmlID, db)
		if !hasAccess {
			switch apiStatus {
			case tc.SystemError:
				handleErr(http.StatusInternalServerError, err)
				return
			case tc.DataMissingError:
				handleErr(http.StatusBadRequest, err)
				return
			case tc.ForbiddenError:
				handleErr(http.StatusForbidden, err)
				return
			}
		}

		// create and start a cluster
		cluster, err := riaksvc.GetRiakCluster(db, cfg.RiakAuthOptions)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		defer func() {
			if err := cluster.Stop(); err != nil {
				log.Errorf("%v\n", err)
			}
		}()

		ro, err := riaksvc.FetchObjectValues(xmlID, CDNURIKeysBucket, cluster)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		// fetch the object and delete it if it exists.
		var alert tc.Alerts

		if ro == nil || ro[0].Value == nil {
			alert = tc.CreateAlerts(tc.InfoLevel, "not deleted, no object found to delete")
		} else if err := riaksvc.DeleteObject(xmlID, CDNURIKeysBucket, cluster); err != nil {
			handleErr(http.StatusInternalServerError, err)
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

// Http POST or PUT handler used to store urisigning keys to a delivery service.
func saveDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)

		defer r.Body.Close()

		if cfg.RiakEnabled == false {
			handleErr(http.StatusServiceUnavailable, fmt.Errorf("The RIAK service is unavailable"))
			return
		}

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		xmlID := pathParams["xmlID"]

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(*user, xmlID, db)
		if !hasAccess {
			switch apiStatus {
			case tc.SystemError:
				handleErr(http.StatusInternalServerError, err)
				return
			case tc.DataMissingError:
				handleErr(http.StatusBadRequest, err)
				return
			case tc.ForbiddenError:
				handleErr(http.StatusForbidden, err)
				return
			}
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		// validate that the received data is a valid jwk keyset
		var keySet map[string]URISignerKeyset
		if err := json.Unmarshal(data, &keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(http.StatusBadRequest, err)
			return
		}
		if err := validateURIKeyset(keySet); err != nil {
			log.Errorf("%v\n", err)
			handleErr(http.StatusBadRequest, err)
			return
		}

		// create and start a cluster
		cluster, err := riaksvc.GetRiakCluster(db, cfg.RiakAuthOptions)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		if err = cluster.Start(); err != nil {
			handleErr(http.StatusInternalServerError, err)
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

		err = riaksvc.SaveObject(obj, CDNURIKeysBucket, cluster)
		if err != nil {
			log.Errorf("%v\n", err)
			handleErr(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", data)
	}
}

// validates URISigingKeyset json.
func validateURIKeyset(msg map[string]URISignerKeyset) error {
	var renewalKidFound int
	var renewalKidMatched = false

	for key, value := range msg {
		issuer := key
		renewalKid := value.RenewalKid
		if issuer == "" {
			return errors.New("JSON Keyset has no issuer")
		}

		if renewalKid != nil {
			renewalKidFound++
		}

		for _, skey := range value.Keys {
			if skey.Algorithm == "" {
				return errors.New("A Key has no algorithm, alg, specified")
			}
			if skey.KeyID == "" {
				return errors.New("A Key has no key id, kid, specified")
			}
			if renewalKid != nil && strings.Compare(*renewalKid, skey.KeyID) == 0 {
				renewalKidMatched = true
			}
		}
	}

	// should only have one renewal_kid
	switch renewalKidFound {
	case 0:
		return errors.New("No renewal_kid was found in any keyset")
	case 1: // okay, this is what we want
		break
	default:
		return errors.New("More than one renewal_kid was found in the keysets")
	}

	// the renewal_kid should match the kid of one key
	if !renewalKidMatched {
		return errors.New("No key was found with a kid that matches the renewal kid")
	}

	return nil
}
