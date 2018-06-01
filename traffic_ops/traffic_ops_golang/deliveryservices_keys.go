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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
)

// Delivery Services: SSL Keys.

// returns the cdn_id found by domainname.
func getCDNIDByDomainname(domainName string, tx *sql.Tx) (int64, bool, error) {
	cdnID := int64(0)
	if err := tx.QueryRow(`SELECT id from cdn WHERE domain_name = $1`, domainName).Scan(&cdnID); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	return cdnID, true, nil
}

// returns a delivery service xmlId for a cdn by host regex.
func getXMLID(cdnID int64, hostRegex string, tx *sql.Tx) (string, bool, error) {
	q := `
SELECT ds.xml_id from deliveryservice ds
JOIN deliveryservice_regex dr on ds.id = dr.deliveryservice AND ds.cdn_id = $1
JOIN regex r on r.id = dr.regex
WHERE r.pattern = $2
`
	xmlID := ""
	if err := tx.QueryRow(q, cdnID, hostRegex).Scan(&xmlID); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying xml id: " + err.Error())
	}
	return xmlID, true, nil
}

func getDeliveryServiceSSLKeysByXMLID(xmlID string, version string, tx *sql.Tx, cfg config.Config) ([]byte, error) {
	if cfg.RiakEnabled == false {
		err := errors.New("Riak is not configured!")
		log.Errorln("getting delivery services SSL keys: " + err.Error())
		return nil, err
	}
	key, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(xmlID, version, tx, cfg.RiakAuthOptions)
	if err != nil {
		log.Errorln("getting delivery service keys: " + err.Error())
		return nil, err
	}
	if !ok {
		alert := tc.CreateAlerts(tc.InfoLevel, "no object found for the specified key")
		respBytes, err := json.Marshal(alert)
		if err != nil {
			log.Errorf("failed to marshal an alert response: %s\n", err)
			return nil, err
		}
		return respBytes, nil
	}

	respBytes := []byte{}
	resp := tc.DeliveryServiceSSLKeysResponse{Response: key}
	respBytes, err = json.Marshal(resp)
	if err != nil {
		log.Errorf("failed to marshal a sslkeys response: %s\n", err)
		return nil, err
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
		return "", fmt.Errorf("could not base64 decode the certificate %v", err)
	}

	// decode, verify, and order certs for storgae
	var bundle string
	certs := strings.SplitAfter(string(pemCerts), "-----END CERTIFICATE-----")
	if len(certs) > 1 {
		// decode and verify the server certificate
		block, _ := pem.Decode([]byte(certs[0]))
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("could not parse the server certificate %v", err)
		}
		if !(cert.KeyUsage&x509.KeyUsageKeyEncipherment > 0) {
			return "", fmt.Errorf("no key encipherment usage for the server certificate")
		}
		for i := 0; i < len(certs)-1; i++ {
			bundle += certs[i]
		}

		var opts x509.VerifyOptions

		rootPool := x509.NewCertPool()
		if rootCA != "" {
			if !rootPool.AppendCertsFromPEM([]byte(rootCA)) {
				return "", fmt.Errorf("root  CA certificate is empty, %v", err)
			}
		}

		intermediatePool := x509.NewCertPool()
		if !intermediatePool.AppendCertsFromPEM([]byte(bundle)) {
			return "", fmt.Errorf("certificate CA bundle is empty, %v", err)
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
			return "", fmt.Errorf("could verify the certificate chain %v", err)
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

func addDeliveryServiceSSLKeysHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)
		if !cfg.RiakEnabled {
			err := errors.New("Riak is not configured!")
			log.Errorln("adding Riak SSL keys for delivery service: " + err.Error())
			handleErr(http.StatusInternalServerError, err)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		keysObj := tc.DeliveryServiceSSLKeys{}
		if err := json.Unmarshal(data, &keysObj); err != nil {
			log.Errorf("ERROR: could not unmarshal the request, %v\n", err)
			handleErr(http.StatusBadRequest, err)
			return
		}

		tx, err := db.DB.Begin()
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("beginning transaction: "+err.Error()))
			return
		}
		commitTx := false
		defer dbhelpers.FinishTx(tx, &commitTx)

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(user, keysObj.DeliveryService, tx)
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

		var certChain string
		if certChain, err = verifyAndEncodeCertificate(keysObj.Certificate.Crt, ""); err != nil {
			log.Errorf("ERROR: could not unmarshal the request, %v\n", err)
			handleErr(http.StatusBadRequest, err)
			return
		}
		keysObj.Certificate.Crt = certChain

		// marshal the keysObj
		keysJSON, err := json.Marshal(&keysObj)
		if err != nil {
			log.Errorf("ERROR: could not marshal the keys object, %v\n", err)
			handleErr(http.StatusBadRequest, err)
			return
		}

		if err := riaksvc.PutDeliveryServiceSSLKeysObj(keysObj, tx, cfg.RiakAuthOptions); err != nil {
			log.Errorln("putting Riak SSL keys for delivery service '" + keysObj.DeliveryService + "': " + err.Error())
			handleErr(http.StatusInternalServerError, err)
			return
		}

		commitTx = true
		w.Header().Set("Content-Type", "application/json")
		w.Write(keysJSON)
	}
}

// fetch the ssl keys for a deliveryservice specified by the fully qualified hostname
func getDeliveryServiceSSLKeysByHostNameHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)
		var respBytes []byte
		var domainName string
		var hostName string
		var hostRegex string

		if cfg.RiakEnabled == false {
			handleErr(http.StatusServiceUnavailable, fmt.Errorf("The RIAK service is unavailable"))
			return
		}

		version := r.URL.Query().Get("version")

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
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

		tx, err := db.DB.Begin()
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("beginning transaction: "+err.Error()))
			return
		}
		commitTx := false
		defer dbhelpers.FinishTx(tx, &commitTx)

		// lookup the cdnID
		cdnID, ok, err := getCDNIDByDomainname(domainName, tx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		if !ok {
			alert := tc.CreateAlerts(tc.InfoLevel, fmt.Sprintf(" - a cdn does not exist for the domain: %s parsed from hostname: %s",
				domainName, hostName))
			respBytes, err = json.Marshal(alert)
			if err != nil {
				log.Errorf("failed to marshal an alert response: %s\n", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(respBytes)
		}
		// now lookup the deliveryservice xmlID
		xmlID, ok, err := getXMLID(cdnID, hostRegex, tx)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting xml id: "+err.Error()))
			return
		}
		if !ok {
			alert := tc.CreateAlerts(tc.InfoLevel, fmt.Sprintf("  - a delivery service does not exist for a host with hostname of %s",
				hostName))
			respBytes, err = json.Marshal(alert)
			if err != nil {
				log.Errorf("failed to marshal an alert response: %s\n", err)
				handleErr(http.StatusInternalServerError, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(respBytes)
		}

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(user, xmlID, tx)
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
		respBytes, err = getDeliveryServiceSSLKeysByXMLID(xmlID, version, tx, cfg)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		commitTx = true
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBytes)
	}
}

// fetch the deliveryservice ssl keys by the specified xmlID.
func getDeliveryServiceSSLKeysByXMLIDHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := tc.GetHandleErrorsFunc(w, r)
		var respBytes []byte

		if cfg.RiakEnabled == false {
			handleErr(http.StatusServiceUnavailable, fmt.Errorf("The RIAK service is unavailable"))
			return
		}

		version := r.URL.Query().Get("version")

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		xmlID := pathParams["xmlID"]

		tx, err := db.DB.Begin()
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("beginning transaction: "+err.Error()))
			return
		}
		commitTx := false
		defer dbhelpers.FinishTx(tx, &commitTx)

		// check user tenancy access to this resource.
		hasAccess, err, apiStatus := tenant.HasTenant(user, xmlID, tx)
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

		respBytes, err = getDeliveryServiceSSLKeysByXMLID(xmlID, version, tx, cfg)
		if err != nil {
			handleErr(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		commitTx = true
		w.Write(respBytes)
	}
}
