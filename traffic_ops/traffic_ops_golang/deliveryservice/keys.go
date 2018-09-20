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
	"bytes"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("adding Riak SSL keys for delivery service:: riak is not configured"))
		return
	}
	keysObj := tc.DeliveryServiceSSLKeys{}
	if err := json.NewDecoder(r.Body).Decode(&keysObj); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON"), nil)
		return
	}
	if userErr, sysErr, errCode := tenant.Check(inf.User, keysObj.DeliveryService, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	certChain, err := verifyAndEncodeCertificate(keysObj.Certificate.Crt, "")
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("verifying certificate: "+err.Error()), nil)
		return
	}
	keysObj.Certificate.Crt = certChain
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(keysObj, inf.Tx.Tx, inf.Config.RiakAuthOptions); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, nil, errors.New("putting Riak SSL keys for delivery service '"+keysObj.DeliveryService+"': "+err.Error()))
		return
	}
	api.WriteRespRaw(w, r, keysObj)
}

// GetSSLKeysByHostName fetches the ssl keys for a deliveryservice specified by the fully qualified hostname
func GetSSLKeysByHostName(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"hostName"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("The RIAK service is unavailable"), errors.New("getting Riak SSL keys by host name: riak is not configured"))
		return
	}

	version := inf.Params["version"]
	hostName := inf.Params["hostName"]
	domainName := ""
	hostRegex := ""
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

	if userErr, sysErr, errCode := tenant.Check(inf.User, xmlID, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(xmlID, version, inf.Tx.Tx, inf.Config.RiakAuthOptions)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ssl keys: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlert(w, r, tc.InfoLevel, "no object found for the specified key")
		return
	}
	api.WriteResp(w, r, keyObj)
}

// GetSSLKeysByXMLID fetches the deliveryservice ssl keys by the specified xmlID.
func GetSSLKeysByXMLID(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"xmlID"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("The RIAK service is unavailable"), errors.New("getting Riak SSL keys by xml id: riak is not configured"))
		return
	}
	version := inf.Params["version"]
	xmlID := inf.Params["xmlID"]
	if userErr, sysErr, errCode := tenant.Check(inf.User, xmlID, inf.Tx.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(xmlID, version, inf.Tx.Tx, inf.Config.RiakAuthOptions)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ssl keys: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlert(w, r, tc.InfoLevel, "no object found for the specified key")
		return
	}
	api.WriteResp(w, r, keyObj)
}

func DeleteSSLKeys(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, userErr, errors.New("deliveryservice.DeleteSSLKeys: Riak is not configured!"))
		return
	}
	ds := tc.DeliveryServiceName(inf.Params["name"])
	if err := riaksvc.DeleteDSSSLKeys(inf.Tx.Tx, inf.Config.RiakAuthOptions, ds, inf.Params["version"]); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, userErr, errors.New("deliveryservice.DeleteSSLKeys: deleting SSL keys: "+err.Error()))
		return
	}
	api.WriteResp(w, r, "Successfully deleted ssl keys for "+string(ds))
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
