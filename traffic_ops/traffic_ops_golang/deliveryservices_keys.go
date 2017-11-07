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
	"bytes"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/basho/riak-go-client"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"strings"
)

// Delivery Services: SSL Keys.

// returns the cdn_id found by domainname.
func getCDNIDByDomainname(domainName string, db *sqlx.DB) (sql.NullInt64, error) {
	cdnQuery := `SELECT id from cdn WHERE domain_name = $1`
	var cdnID sql.NullInt64

	noCdnID := sql.NullInt64{
		Int64: 0,
		Valid: false,
	}

	rows, err := db.Query(cdnQuery, domainName)
	if err != nil {
		return noCdnID, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&cdnID); err != nil {
			return noCdnID, err
		}
	}

	return cdnID, nil
}

func getDeliveryServiceCountByXmlID(xmlID string, db *sqlx.DB) (int64, error) {
	dsQuery := `SELECT count(*)  from deliveryservice WHERE xml_id = $1`
	var count sql.NullInt64

	rows, err := db.Query(dsQuery, xmlID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count.Int64, nil
}

// returns a delivery service xmlId for a cdn by host regex.
func getXmlIDByCDNAndRegex(cdnID sql.NullInt64, hostRegex string, db *sqlx.DB) (sql.NullString, error) {
	dsQuery := `
			SELECT ds.xml_id from deliveryservice ds
			INNER JOIN deliveryservice_regex dr 
			on ds.id = dr.deliveryservice AND ds.cdn_id = $1
			INNER JOIN regex r on r.id = dr.regex
			WHERE r.pattern = $2
		`
	var xmlID sql.NullString

	rows, err := db.Query(dsQuery, cdnID.Int64, hostRegex)
	if err != nil {
		xmlID.Valid = false
		return xmlID, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&xmlID); err != nil {
			xmlID.Valid = false
			return xmlID, err
		}
	}

	return xmlID, nil
}

func getDeliveryServiceSSLKeysByXmlID(xmlID string, version string, db *sqlx.DB, cfg Config) ([]byte, error) {
	var respBytes []byte
	// create and start a cluster
	cluster, err := getRiakCluster(db, cfg)
	if err != nil {
		return nil, err
	}
	if err = cluster.Start(); err != nil {
		return nil, err
	}
	defer func() {
		if err := cluster.Stop(); err != nil {
			log.Errorf("%v\n", err)
		}
	}()

	if version == "" {
		xmlID = xmlID + "-latest"
	} else {
		xmlID = xmlID + "-" + version
	}

	// get the deliveryservice ssl keys by xmlID and version
	ro, err := fetchObjectValues(xmlID, SSLKeysBucket, cluster)
	if err != nil {
		return nil, err
	}

	// no keys we're found
	if ro == nil {
		alert := tc.CreateAlerts(tc.InfoLevel, "no object found for the specified key")
		respBytes, err = json.Marshal(alert)
		if err != nil {
			log.Errorf("failed to marshal an alert response: %s\n", err)
			return nil, err
		}
	} else { // keys were found
		var key tc.DeliveryServiceSSLKeys

		// unmarshal into a response tc.DeliveryServiceSSLKeysResponse object.
		if err := json.Unmarshal(ro[0].Value, &key); err != nil {
			log.Errorf("failed at unmarshaling sslkey response: %s\n", err)
			return nil, err
		}
		resp := tc.DeliveryServiceSSLKeysResponse{
			Response: key,
		}
		respBytes, err = json.Marshal(resp)
		if err != nil {
			log.Errorf("failed to marshal a sslkeys response: %s\n", err)
			return nil, err
		}
	}

	return respBytes, nil
}

// verify the server certificate chain and return the
// certificate and its chain in the proper order. Returns a  verified,
// ordered, and base64 encoded certificate and CA chain.
func verifyAndEncodeCertificate(certificate string, rootCA string) (string, error) {
	var pemEncodedChain string
	var b64crt string

	// strip newlines from encoded crt and decode it from base64.
	crtArr := strings.Split(certificate, "\\n")
	for i := 0; i < len(crtArr); i++ {
		b64crt += crtArr[i]
	}
	pemCerts := make([]byte, base64.StdEncoding.EncodedLen(len(b64crt)))
	_, err := base64.StdEncoding.Decode(pemCerts, []byte(b64crt))
	if err != nil {
		return "", fmt.Errorf("ERROR: could not base64 decode the certificate, %v\n", err)
	}

	// decode, verify, and order certs for storgae
	var bundle string
	certs := strings.SplitAfter(string(pemCerts), "-----END CERTIFICATE-----")
	if len(certs) > 1 {
		// decode and verify the server certificate
		block, _ := pem.Decode([]byte(certs[0]))
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("ERROR: could not parse the server certificate, %v\n", err)
		}
		if !(cert.KeyUsage&x509.KeyUsageKeyEncipherment > 0) {
			return "", fmt.Errorf("ERROR: no key encipherment usage for the server certificate\n")
		}
		for i := 0; i < len(certs)-1; i++ {
			bundle += certs[i]
		}

		var opts x509.VerifyOptions

		rootPool := x509.NewCertPool()
		if rootCA != "" {
			if !rootPool.AppendCertsFromPEM([]byte(rootCA)) {
				return "", fmt.Errorf("ERROR: root  CA certificate is empty, %v\n", err)
			}
		}

		intermediatePool := x509.NewCertPool()
		if !intermediatePool.AppendCertsFromPEM([]byte(bundle)) {
			return "", fmt.Errorf("ERROR: certificate CA bundle is empty, %v\n", err)
		}

		if rootCA != "" {
			// verify the certificate chain.
			opts = x509.VerifyOptions{
				Intermediates: intermediatePool,
				Roots:         rootPool,
			}
		} else {
			opts = x509.VerifyOptions{
				Intermediates: intermediatePool,
			}
		}

		chain, err := cert.Verify(opts)
		if err != nil {
			return "", fmt.Errorf("ERROR: could verify the certificate chain, %v\n", err)
		}
		if len(chain) > 0 {
			for _, link := range chain[0] {
				// Only print non-self signed elements of the chain
				if link.AuthorityKeyId != nil && !bytes.Equal(link.AuthorityKeyId, link.SubjectKeyId) {
					block := &pem.Block{Type: "CERTIFICATE", Bytes: link.Raw}
					pemEncodedChain += string(pem.EncodeToMemory(block))
				}
			}
		} else {
			return "", fmt.Errorf("Can't find valid chain for cert in file in request")
		}
	} else {
		return "", fmt.Errorf("ERROR: no certificate chain to verify")
	}

	base64EncodedStr := base64.StdEncoding.EncodeToString([]byte(pemEncodedChain))

	return base64EncodedStr, nil
}

func addDeliveryServiceSSLKeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)
		var keysObj tc.DeliveryServiceSSLKeys
		var respBytes []byte

		defer r.Body.Close()

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// unmarshal the request
		if err := json.Unmarshal(data, &keysObj); err != nil {
			log.Errorf("ERROR: could not unmarshal the request, %v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}

		dsCount, err := getDeliveryServiceCountByXmlID(keysObj.DeliveryService, db)
		if err != nil {
			log.Errorf("ERROR: querying deliveryservice, %v\n", err)
			handleErr(err, http.StatusInternalServerError)
			return
		}
		if dsCount != 1 {
			alert := tc.CreateAlerts(tc.InfoLevel, fmt.Sprintf(" - a delivery service does not exist named: %s",
				keysObj.DeliveryService))
			respBytes, err = json.Marshal(alert)
			if err != nil {
				log.Errorf("failed to marshal an alert response: %s\n", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", respBytes)
			return
		}

		var certChain string
		if certChain, err = verifyAndEncodeCertificate(keysObj.Certificate.Crt, ""); err != nil {
			log.Errorf("ERROR: could not unmarshal the request, %v\n", err)
			handleErr(err, http.StatusBadRequest)
			return
		}
		keysObj.Certificate.Crt = certChain

		// marshal the keysObj
		keysJson, err := json.Marshal(&keysObj)
		if err != nil {
			log.Errorf("ERROR: could not marshal the keys object, %v\n", err)
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
			Key:             keysObj.DeliveryService,
			Value:           []byte(keysJson),
		}

		err = saveObject(obj, SSLKeysBucket, cluster)
		if err != nil {
			log.Errorf("%v\n", err)
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", keysJson)
	}
}

// fetch the ssl keys for a deliveryservice specified by the fully qualified hostname
func getDeliveryServiceSSLKeysByHostNameHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)
		var respBytes []byte
		var domainName string
		var hostName string
		var hostRegex string

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

		version := r.URL.Query().Get("version")

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		hostName = pathParams["hostName"]

		strArr := strings.Split(hostName, ".")
		ln := len(strArr)

		if ln > 1 {
			for i := 2; i < ln-1; i++ {
				domainName += strArr[i] + "."
			}
			domainName += strArr[ln-1]
			hostRegex = ".*\\." + strArr[1] + "\\..*"
		}

		// lookup the cdnID
		cdnID, err := getCDNIDByDomainname(domainName, db)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// verify that a valid cdnID was returned.
		if !cdnID.Valid {
			alert := tc.CreateAlerts(tc.InfoLevel, fmt.Sprintf(" - a cdn does not exist for the domain: %s parsed from hostname: %s",
				domainName, hostName))
			respBytes, err = json.Marshal(alert)
			if err != nil {
				log.Errorf("failed to marshal an alert response: %s\n", err)
				return
			}
		} else {
			// now lookup the deliveryservice xmlID
			xmlIDStr, err := getXmlIDByCDNAndRegex(cdnID, hostRegex, db)
			if err != nil {
				handleErr(err, http.StatusInternalServerError)
				return
			}

			// verify that the xmlIDStr returned is valid, ie not nil
			if !xmlIDStr.Valid {
				alert := tc.CreateAlerts(tc.InfoLevel, fmt.Sprintf("  - a delivery service does not exist for a host with hostname of %s",
					hostName))
				respBytes, err = json.Marshal(alert)
				if err != nil {
					log.Errorf("failed to marshal an alert response: %s\n", err)
					handleErr(err, http.StatusInternalServerError)
					return
				}
			} else {
				xmlID := xmlIDStr.String
				respBytes, err = getDeliveryServiceSSLKeysByXmlID(xmlID, version, db, cfg)
				if err != nil {
					handleErr(err, http.StatusInternalServerError)
					return
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBytes)
	}
}

// fetch the deliveryservice ssl keys by the specified xmlID.
func getDeliveryServiceSSLKeysByXmlIDHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)
		var respBytes []byte

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

		version := r.URL.Query().Get("version")

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		xmlID := pathParams["xmlID"]

		respBytes, err = getDeliveryServiceSSLKeysByXmlID(xmlID, version, db, cfg)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBytes)
	}
}

// Delivery Services: URI Sign Keys.

// Http POST handler used to store urisigning keys to a delivery service.
func assignDeliveryServiceURIKeysHandler(db *sqlx.DB, cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		defer r.Body.Close()

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

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

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

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

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

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
			alert = tc.CreateAlerts(tc.InfoLevel, "not deleted, no object found to delete")
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

		if cfg.RiakEnabled == false {
			handleErr(fmt.Errorf("The RIAK service is unavailable"), http.StatusServiceUnavailable)
			return
		}

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
