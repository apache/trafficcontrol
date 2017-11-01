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
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat/go-jwx/jwk"
	"io/ioutil"
	"net/http"
	"strings"
)

// RiakPort is the port RIAK is listening on.
const RiakPort = 8087

// CDNURIKeysBucket is the namespace or bucket used for CDN URI signing keys.
const CDNURIKeysBucket = "cdn_uri_sig_keys"

// URISignerKeyset is the container for the CDN URI signing keys
type URISignerKeyset struct {
	RenewalKid *string               `json:"renewal_kid"`
	Keys       []jwk.EssentialHeader `json:"keys"`
}

// deletes an object from riak storage
func deleteObject(key string, bucket string, cluster *riak.Cluster) error {
	// build store command and execute.
	cmd, err := riak.NewDeleteValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return err
	}
	if err := cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}

// fetch an object from riak storage
func fetchObjectValues(key string, bucket string, cluster *riak.Cluster) ([]*riak.Object, error) {
	// build the fetch command
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return nil, err
	}

	if err = cluster.Execute(cmd); err != nil {
		return nil, err
	}
	fvc := cmd.(*riak.FetchValueCommand)

	// no object found with given key
	if fvc.Response == nil || fvc.Response.IsNotFound {
		return nil, nil
	}
	return fvc.Response.Values, nil
}

// saves an object to riak storage
func saveObject(obj *riak.Object, bucket string, cluster *riak.Cluster) error {
	// build store command and execute.
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(bucket).
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}
	if err := cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}

// returns a riak cluster of online riak nodes.
func getRiakCluster(db *sqlx.DB, cfg Config) (*riak.Cluster, error) {
	riakServerQuery := `
		SELECT s.host_name, s.domain_name FROM server s 
		INNER JOIN type t on s.type = t.id 
		INNER JOIN status st on s.status = st.id 
		WHERE t.name = 'RIAK' AND st.name = 'ONLINE'
		`

	if cfg.RiakAuthOptions == nil {
		return nil, errors.New("ERROR: no riak auth information from riak.conf, cannot authenticate to any riak servers")
	}

	var nodes []*riak.Node
	rows, err := db.Query(riakServerQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s tc.Server
		var n *riak.Node
		if err := rows.Scan(&s.HostName, &s.DomainName); err != nil {
			return nil, err
		}
		addr := fmt.Sprintf("%s.%s:%d", s.HostName, s.DomainName, RiakPort)
		nodeOpts := &riak.NodeOptions{
			RemoteAddress: addr,
			AuthOptions:   cfg.RiakAuthOptions,
		}
		nodeOpts.AuthOptions.TlsConfig.ServerName = fmt.Sprintf("%s.%s", s.HostName, s.DomainName)
		n, err := riak.NewNode(nodeOpts)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	if len(nodes) == 0 {
		return nil, errors.New("ERROR: no available riak servers")
	}

	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}
	cluster, err := riak.NewCluster(opts)

	return cluster, err
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
