package deliveryservice

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
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
)

const (
	PemCertEndMarker = "-----END CERTIFICATE-----"
)

// AddSSLKeys adds the given ssl keys to the given delivery service.
func AddSSLKeys(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if !inf.Config.RiakEnabled {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("adding SSL keys to Riak for delivery service: Riak is not configured"))
		return
	}
	req := tc.DeliveryServiceAddSSLKeysReq{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &req); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}
	if userErr, sysErr, errCode := tenant.Check(inf.User, *req.DeliveryService, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	certChain, isUnknownAuth, err := verifyCertificate(req.Certificate.Crt, "")
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("verifying certificate: "+err.Error()), nil)
		return
	}
	req.Certificate.Crt = certChain
	base64EncodeCertificate(req.Certificate)
	dsSSLKeys := tc.DeliveryServiceSSLKeys{
		CDN:             *req.CDN,
		DeliveryService: *req.DeliveryService,
		Hostname:        *req.HostName,
		Key:             *req.Key,
		Version:         *req.Version,
		Certificate:     *req.Certificate,
	}
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(dsSSLKeys, inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("putting SSL keys in Riak for delivery service '"+*req.DeliveryService+"': "+err.Error()))
		return
	}
	if err := updateSSLKeyVersion(*req.DeliveryService, req.Version.ToInt64(), inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("adding SSL keys to delivery service '"+*req.DeliveryService+"': "+err.Error()))
		return
	}
	if isUnknownAuth {
		api.WriteRespAlert(w, r, tc.WarnLevel, "WARNING: SSL keys were successfully added for '"+*req.DeliveryService+"', but the certificate is signed by an unknown authority and may be invalid")
		return
	}
	api.WriteResp(w, r, "Successfully added ssl keys for "+*req.DeliveryService)
}

// GetSSLKeysByHostName fetches the ssl keys for a deliveryservice specified by the fully qualified hostname
func GetSSLKeysByHostName(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"hostname"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("the Riak service is unavailable"), errors.New("getting SSL keys from Riak by host name: Riak is not configured"))
		return
	}

	hostName := inf.Params["hostname"]
	domainName := ""
	hostRegex := ""
	strArr := strings.Split(hostName, ".")
	ln := len(strArr)
	if ln > 1 {
		for i := 2; i < ln-1; i++ {
			domainName += strArr[i] + "."
		}
		domainName += strArr[ln-1]
		hostRegex = `.*\.` + strArr[1] + `\..*`
	}

	// lookup the cdnID
	cdnID, ok, err := getCDNIDByDomainname(domainName, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn id by domain name: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlert(w, r, tc.InfoLevel, " - a cdn does not exist for the domain: "+domainName+" parsed from hostname: "+hostName)
		return
	}
	// now lookup the deliveryservice xmlID
	xmlID, ok, err := getXMLID(cdnID, hostRegex, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting xml id: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlert(w, r, tc.InfoLevel, "  - a delivery service does not exist for a host with hostname of "+hostName)
		return
	}

	getSSLKeysByXMLIDHelper(xmlID, inf, w, r)
}

// GetSSLKeysByXMLID fetches the deliveryservice ssl keys by the specified xmlID.
func GetSSLKeysByXMLID(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"xmlid"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("the Riak service is unavailable"), errors.New("getting SSL keys from Riak by xml id: Riak is not configured"))
		return
	}
	xmlID := inf.Params["xmlid"]
	getSSLKeysByXMLIDHelper(xmlID, inf, w, r)
}

func getSSLKeysByXMLIDHelper(xmlID string, inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	version := inf.Params["version"]
	decode := inf.Params["decode"]
	if userErr, sysErr, errCode := tenant.Check(inf.User, xmlID, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(xmlID, version, inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ssl keys: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlertObj(w, r, tc.InfoLevel, "no object found for the specified key", struct{}{}) // empty response object because Perl
		return
	}
	if decode != "" && decode != "0" { // the Perl version checked the decode string as: if ( $decode )
		err = base64DecodeCertificate(&keyObj.Certificate)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting SSL keys for XMLID '"+xmlID+"': "+err.Error()))
			return
		}
	}
	api.WriteResp(w, r, keyObj)
}

func base64DecodeCertificate(cert *tc.DeliveryServiceSSLKeysCertificate) error {
	csrDec, err := base64.StdEncoding.DecodeString(cert.CSR)
	if err != nil {
		return errors.New("base64 decoding csr: " + err.Error())
	}
	cert.CSR = string(csrDec)
	crtDec, err := base64.StdEncoding.DecodeString(cert.Crt)
	if err != nil {
		return errors.New("base64 decoding crt: " + err.Error())
	}
	cert.Crt = string(crtDec)
	keyDec, err := base64.StdEncoding.DecodeString(cert.Key)
	if err != nil {
		return errors.New("base64 decoding key: " + err.Error())
	}
	cert.Key = string(keyDec)
	return nil
}

func base64EncodeCertificate(cert *tc.DeliveryServiceSSLKeysCertificate) {
	cert.CSR = base64.StdEncoding.EncodeToString([]byte(cert.CSR))
	cert.Crt = base64.StdEncoding.EncodeToString([]byte(cert.Crt))
	cert.Key = base64.StdEncoding.EncodeToString([]byte(cert.Key))
}

func DeleteSSLKeys(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"xmlid"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, userErr, errors.New("deliveryservice.DeleteSSLKeys: Riak is not configured"))
		return
	}
	xmlID := inf.Params["xmlid"]
	if userErr, sysErr, errCode := tenant.Check(inf.User, xmlID, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	if err := riaksvc.DeleteDSSSLKeys(inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort, xmlID, inf.Params["version"]); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, userErr, errors.New("deliveryservice.DeleteSSLKeys: deleting SSL keys: "+err.Error()))
		return
	}
	api.WriteResp(w, r, "Successfully deleted ssl keys for "+xmlID)
}

func updateSSLKeyVersion(xmlID string, version int64, tx *sql.Tx) error {
	q := `UPDATE deliveryservice SET ssl_key_version = $1 WHERE xml_id = $2`
	if _, err := tx.Exec(q, version, xmlID); err != nil {
		return errors.New("updating delivery service ssl_key_version: " + err.Error())
	}
	return nil
}

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

// verify the server certificate chain and return the
// certificate and its chain in the proper order. Returns a verified
// and ordered certificate and CA chain.
// If the cert verification returns UnknownAuthorityError, return true to
// indicate that the certs are signed by an unknown authority (e.g. self-signed).
func verifyCertificate(certificate string, rootCA string) (string, bool, error) {
	// decode, verify, and order certs for storage
	certs := strings.SplitAfter(certificate, PemCertEndMarker)
	if len(certs) <= 1 {
		return "", false, errors.New("no certificate chain to verify")
	}

	// decode and verify the server certificate
	block, _ := pem.Decode([]byte(certs[0]))
	if block == nil {
		return "", false, errors.New("could not decode pem-encoded server certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", false, errors.New("could not parse the server certificate: " + err.Error())
	}
	if !(cert.KeyUsage&x509.KeyUsageKeyEncipherment > 0) {
		return "", false, errors.New("no key encipherment usage for the server certificate")
	}
	bundle := ""
	for i := 0; i < len(certs)-1; i++ {
		bundle += certs[i]
	}

	intermediatePool := x509.NewCertPool()
	if !intermediatePool.AppendCertsFromPEM([]byte(bundle)) {
		return "", false, errors.New("certificate CA bundle is empty")
	}

	opts := x509.VerifyOptions{
		Intermediates: intermediatePool,
	}
	if rootCA != "" {
		// verify the certificate chain.
		rootPool := x509.NewCertPool()
		if !rootPool.AppendCertsFromPEM([]byte(rootCA)) {
			return "", false, errors.New("unable to parse root CA certificate")
		}
		opts.Roots = rootPool
	}

	chain, err := cert.Verify(opts)
	if err != nil {
		if _, ok := err.(x509.UnknownAuthorityError); ok {
			return certificate, true, nil
		}
		return "", false, errors.New("could not verify the certificate chain: " + err.Error())
	}
	if len(chain) < 1 {
		return "", false, errors.New("can't find valid chain for cert in file in request")
	}
	pemEncodedChain := ""
	for _, link := range chain[0] {
		// Include all certificates in the chain, since verification was successful.
		block := &pem.Block{Type: "CERTIFICATE", Bytes: link.Raw}
		pemEncodedChain += string(pem.EncodeToMemory(block))
	}
   
  	if len(pemEncodedChain) < 1 {
		return "", false, errors.New("Invalid empty certicate chain in request")
  	}

	return pemEncodedChain, false, nil
}
