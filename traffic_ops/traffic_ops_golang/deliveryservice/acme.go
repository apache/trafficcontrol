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
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/jmoiron/sqlx"
)

const validAccountStatus = "valid"

// RenewAcmeCertificate renews the SSL certificate for a delivery service if possible through ACME protocol.
func RenewAcmeCertificate(w http.ResponseWriter, r *http.Request) {
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

	ctx, _ := context.WithTimeout(r.Context(), LetsEncryptTimeout)

	userErr, sysErr, statusCode := renewAcmeCerts(inf.Config, xmlID, ctx, inf.User)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
	}

	api.WriteRespAlert(w, r, tc.SuccessLevel, "Certificate for "+xmlID+" successfully renewed.")

}

func renewAcmeCerts(cfg *config.Config, dsName string, ctx context.Context, currentUser *auth.CurrentUser) (error, error, int) {
	db, err := api.GetDB(ctx)
	if err != nil {
		log.Errorf(dsName+": Error getting db: %s", err.Error())
		return nil, err, http.StatusInternalServerError
	}

	tx, err := db.Begin()
	if err != nil {
		log.Errorf(dsName+": Error getting tx: %s", err.Error())
		return nil, err, http.StatusInternalServerError
	}

	userTx, err := db.Begin()
	if err != nil {
		log.Errorf(dsName+": Error getting userTx: %s", err.Error())
		return nil, err, http.StatusInternalServerError
	}
	defer userTx.Commit()

	logTx, err := db.Begin()
	if err != nil {
		log.Errorf(dsName+": Error getting logTx: %s", err.Error())
		return nil, err, http.StatusInternalServerError
	}
	defer logTx.Commit()

	dsID, certVersion, err := getDSIdAndVersionFromName(db, dsName)
	if err != nil {
		return nil, errors.New("querying DS info: " + err.Error()), http.StatusInternalServerError
	}
	if dsID == nil || *dsID == 0 {
		return errors.New("DS id for " + dsName + " was nil or 0"), nil, http.StatusBadRequest
	}
	if certVersion == nil || *certVersion == 0 {
		return errors.New("certificate for " + dsName + " could not be renewed because version was nil or 0"), nil, http.StatusBadRequest
	}

	if cfg == nil {
		return nil, errors.New("acme: config was nil"), http.StatusInternalServerError
	}
	keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObjV15(dsName, strconv.Itoa(int(*certVersion)), tx, cfg.RiakAuthOptions, cfg.RiakPort)
	if err != nil {
		return nil, errors.New("getting ssl keys for xmlId: " + dsName + " and version: " + strconv.Itoa(int(*certVersion)) + " : " + err.Error()), http.StatusInternalServerError
	}
	if !ok {
		return nil, errors.New("no object found for the specified key with xmlId: " + dsName + " and version: " + strconv.Itoa(int(*certVersion))), http.StatusInternalServerError
	}

	err = base64DecodeCertificate(&keyObj.Certificate)
	if err != nil {
		return nil, errors.New("decoding cert for XMLID " + dsName + " : " + err.Error()), http.StatusInternalServerError
	}

	acmeAccount := GetAcmeAccountConfig(cfg, keyObj.AuthType)
	if acmeAccount == nil {
		return nil, errors.New("No acme account information in cdn.conf for " + keyObj.AuthType), http.StatusInternalServerError
	}

	client, err := GetAcmeClient(acmeAccount, userTx, db)
	if err != nil {
		api.CreateChangeLogRawTx(api.ApiChange, "DS: "+dsName+", ID: "+strconv.Itoa(*dsID)+", ACTION: FAILED to add SSL keys with "+acmeAccount.AcmeProvider, currentUser, logTx)
		return nil, errors.New("getting acme client: " + err.Error()), http.StatusInternalServerError
	}

	renewRequest := certificate.Resource{
		Certificate: []byte(keyObj.Certificate.Crt),
	}

	cert, err := client.Certificate.Renew(renewRequest, true, false)
	if err != nil || cert == nil {
		log.Errorf("Error obtaining acme certificate: %s", err.Error())
		return nil, err, http.StatusInternalServerError
	}

	newCertObj := tc.DeliveryServiceSSLKeys{
		AuthType:        keyObj.AuthType,
		CDN:             keyObj.CDN,
		DeliveryService: keyObj.DeliveryService,
		Key:             keyObj.DeliveryService,
		Hostname:        keyObj.Hostname,
		Version:         keyObj.Version + 1,
	}

	newCertObj.Certificate = tc.DeliveryServiceSSLKeysCertificate{
		Crt: string(EncodePEMToLegacyPerlRiakFormat(cert.Certificate)),
		Key: string(EncodePEMToLegacyPerlRiakFormat(cert.PrivateKey)),
		CSR: string(EncodePEMToLegacyPerlRiakFormat([]byte("ACME Generated"))),
	}

	if err := riaksvc.PutDeliveryServiceSSLKeysObj(newCertObj, tx, cfg.RiakAuthOptions, cfg.RiakPort); err != nil {
		log.Errorf("Error posting acme certificate to riak: %s", err.Error())
		api.CreateChangeLogRawTx(api.ApiChange, "DS: "+dsName+", ID: "+strconv.Itoa(*dsID)+", ACTION: FAILED to add SSL keys with "+acmeAccount.AcmeProvider, currentUser, logTx)
		return nil, errors.New(dsName + ": putting riak keys: " + err.Error()), http.StatusInternalServerError
	}

	tx2, err := db.Begin()
	if err != nil {
		log.Errorf("starting sql transaction for delivery service " + dsName + ": " + err.Error())
		return nil, errors.New("starting sql transaction for delivery service " + dsName + ": " + err.Error()), http.StatusInternalServerError
	}

	if err := updateSSLKeyVersion(dsName, *certVersion+1, tx2); err != nil {
		log.Errorf("updating SSL key version for delivery service '" + dsName + "': " + err.Error())
		return nil, errors.New("updating SSL key version for delivery service '" + dsName + "': " + err.Error()), http.StatusInternalServerError
	}
	tx2.Commit()

	api.CreateChangeLogRawTx(api.ApiChange, "DS: "+dsName+", ID: "+strconv.Itoa(*dsID)+", ACTION: Added SSL keys with "+acmeAccount.AcmeProvider, currentUser, logTx)

	return nil, nil, http.StatusOK
}

// GetAcmeAccountConfig returns the ACME account information from cdn.conf for a given provider.
func GetAcmeAccountConfig(cfg *config.Config, acmeProvider string) *config.ConfigAcmeAccount {
	for _, acmeCfg := range cfg.AcmeAccounts {
		if acmeCfg.AcmeProvider == acmeProvider {
			return &acmeCfg
		}
	}
	return nil
}

func getDSIdAndVersionFromName(db *sqlx.DB, xmlId string) (*int, *int64, error) {
	var dsID int
	var certVersion int64

	if err := db.QueryRow(`SELECT id, ssl_key_version FROM deliveryservice WHERE xml_id = $1`, xmlId).Scan(&dsID, &certVersion); err != nil {
		return nil, nil, err
	}

	return &dsID, &certVersion, nil
}

// GetAcmeClient uses the ACME account information in either cdn.conf or the database to create and register an ACME client.
func GetAcmeClient(acmeAccount *config.ConfigAcmeAccount, userTx *sql.Tx, db *sqlx.DB) (*lego.Client, error) {
	if acmeAccount.UserEmail == "" {
		log.Errorf("An email address must be provided to use ACME with %v", acmeAccount.AcmeProvider)
		return nil, errors.New("An email address must be provided to use ACME with " + acmeAccount.AcmeProvider)
	}
	storedAcmeInfo, err := getStoredAcmeAccountInfo(userTx, acmeAccount.UserEmail, acmeAccount.AcmeProvider)
	if err != nil {
		log.Errorf("Error finding stored ACME information: %s", err.Error())
		return nil, err
	}

	myUser := MyUser{}
	foundPreviousAccount := false
	userPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorf("Error generating private key: %s", err.Error())
		return nil, err
	}

	if storedAcmeInfo == nil || acmeAccount.UserEmail == "" {
		myUser = MyUser{
			key:   userPrivateKey,
			Email: acmeAccount.UserEmail,
		}
	} else {
		foundPreviousAccount = true
		myUser = MyUser{
			key:   &storedAcmeInfo.PrivateKey,
			Email: storedAcmeInfo.Email,
			Registration: &registration.Resource{
				URI: storedAcmeInfo.URI,
			},
		}
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = acmeAccount.AcmeUrl
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Errorf("Error creating acme client: %s", err.Error())
		return nil, err
	}

	if acmeAccount.AcmeProvider == tc.LetsEncryptAuthType {
		client.Challenge.Remove(challenge.HTTP01)
		client.Challenge.Remove(challenge.TLSALPN01)
		trafficRouterDns := NewDNSProviderTrafficRouter()
		trafficRouterDns.db = db
		if err != nil {
			log.Errorf("Error creating Traffic Router DNS provider: %s", err.Error())
			return nil, err
		}
		client.Challenge.SetDNS01Provider(trafficRouterDns)
	}

	if foundPreviousAccount {
		log.Debugf("Found existing account with %s", acmeAccount.AcmeProvider)
		reg, err := client.Registration.QueryRegistration()
		if err != nil {
			log.Errorf("Error querying %s for existing account: %s", acmeAccount.AcmeProvider, err.Error())
			return nil, err
		}
		myUser.Registration = reg
		if reg.Body.Status != validAccountStatus {
			log.Debugf("Account found with %s is not valid.", acmeAccount.AcmeProvider)
			foundPreviousAccount = false
		}
	}
	if !foundPreviousAccount {
		if acmeAccount.Kid != "" && acmeAccount.HmacEncoded != "" {
			reg, err := client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
				TermsOfServiceAgreed: true,
				Kid:                  acmeAccount.Kid,
				HmacEncoded:          acmeAccount.HmacEncoded,
			})
			if err != nil {
				log.Errorf("Error registering acme client with external account binding: %s", err.Error())
				return nil, err
			}
			myUser.Registration = reg
			log.Debugf("Creating a new account with %s", acmeAccount.AcmeProvider)
		} else {
			reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
			if err != nil {
				log.Errorf("Error registering acme client: %s", err.Error())
				return nil, err
			}
			myUser.Registration = reg
			log.Debugf("Creating a new account with %s", acmeAccount.AcmeProvider)
		}

		// save account info
		userKeyPem, err := ConvertPrivateKeyToKeyPem(userPrivateKey)
		if err != nil {
			return nil, err
		}
		err = storeAcmeAccountInfo(userTx, myUser.Email, string(userKeyPem), myUser.Registration.URI, acmeAccount.AcmeProvider)
		if err != nil {
			log.Errorf("storing user account info: " + err.Error())
			return nil, errors.New("storing user account info: " + err.Error())
		}
	}

	return client, nil
}

// ConvertPrivateKeyToKeyPem converts an rsa.PrivateKey to be PEM encoded.
func ConvertPrivateKeyToKeyPem(userPrivateKey *rsa.PrivateKey) ([]byte, error) {
	userKeyDer := x509.MarshalPKCS1PrivateKey(userPrivateKey)
	if userKeyDer == nil {
		log.Errorf("marshalling private key: nil der")
		return nil, errors.New("marshalling private key: nil der")
	}
	userKeyBuf := bytes.Buffer{}
	if err := pem.Encode(&userKeyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: userKeyDer}); err != nil {
		log.Errorf("pem-encoding private key: " + err.Error())
		return nil, errors.New("pem-encoding private key: " + err.Error())
	}
	return userKeyBuf.Bytes(), nil
}

func getStoredAcmeAccountInfo(tx *sql.Tx, email string, provider string) (*AcmeInfo, error) {
	acmeInfo := AcmeInfo{}
	selectQuery := `SELECT email, private_key, uri FROM acme_account WHERE email = $1 AND provider = $2 LIMIT 1`
	if err := tx.QueryRow(selectQuery, email, provider).Scan(&acmeInfo.Email, &acmeInfo.Key, &acmeInfo.URI); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("getting lets encrypt account record: " + err.Error())
	}

	decodedKeyBlock, _ := pem.Decode([]byte(acmeInfo.Key))
	decodedKey, err := x509.ParsePKCS1PrivateKey(decodedKeyBlock.Bytes)
	if err != nil {
		return nil, errors.New("decoding private key for user account")
	}
	acmeInfo.PrivateKey = *decodedKey

	return &acmeInfo, nil
}

func storeAcmeAccountInfo(tx *sql.Tx, email string, privateKey string, uri string, provider string) error {
	q := `INSERT INTO acme_account (email, private_key, uri, provider) VALUES ($1, $2, $3, $4)`
	response, err := tx.Exec(q, email, privateKey, uri, provider)
	if err != nil {
		return err
	}

	rows, err := response.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("zero rows affected when inserting Let's Encrypt account information")
	}

	return nil
}

type AcmeInfo struct {
	Email      string `db:"email"`
	Key        string `db:"private_key"`
	URI        string `db:"uri"`
	PrivateKey rsa.PrivateKey
}
