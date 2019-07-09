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
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}

func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

type DNSProviderTrafficRouter struct {
	db *sqlx.DB
}

func NewDNSProviderTrafficRouter() (*DNSProviderTrafficRouter, error) {
	return &DNSProviderTrafficRouter{}, nil
}

func (d *DNSProviderTrafficRouter) Timeout() (timeout, interval time.Duration) {
	return time.Minute * 10, time.Second * 30
}

func (d *DNSProviderTrafficRouter) Present(domain, token, keyAuth string) error {
	tx, err := d.db.Begin()
	fqdn, value := dns01.GetRecord(domain, keyAuth)

	q := `INSERT INTO dnschallenges (fqdn, record) VALUES ($1, $2)`
	response, err := tx.Exec(q, fqdn, value)
	tx.Commit()
	if err != nil {
		log.Errorf("Inserting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		return errors.New("Inserting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
	} else {
		rows, err := response.RowsAffected()
		if err != nil {
			log.Errorf("Determining rows affected dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
			return errors.New("Determining rows affected dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		}
		if rows == 0 {
			log.Errorf("Zero rows affected when inserting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
			return errors.New("Zero rows affected when inserting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		}
	}

	return nil
}

func (d *DNSProviderTrafficRouter) CleanUp(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	tx, err := d.db.Begin()

	q := `DELETE FROM dnschallenges WHERE fqdn = $1 and record = $2`
	response, err := tx.Exec(q, fqdn, value)
	tx.Commit()
	if err != nil {
		log.Errorf("Deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		return errors.New("Deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
	} else {
		rows, err := response.RowsAffected()
		if err != nil {
			log.Errorf("Determining rows affected when deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
			return errors.New("Determining rows affected when deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		}
		if rows == 0 {
			log.Errorf("Zero rows affected when deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
			return errors.New("Zero rows affected when deleting dns txt record for fqdn '" + fqdn + "' record '" + value + "': " + err.Error())
		}
	}

	return nil
}

func GenerateLetsEncryptCertificates(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ctx, _ := context.WithTimeout(r.Context(), time.Minute*10)
	db, err := api.GetDB(ctx)
	tx, err := db.BeginTxx(ctx, nil)

	req := tc.DeliveryServiceLetsEncryptSSLKeysReq{}
	if err := api.Parse(r.Body, nil, &req); err != nil {
		log.Errorf("Error parsing request: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}

	domainName := *req.HostName
	deliveryService := *req.DeliveryService

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorf("Error generating private key: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	myUser := MyUser{
		key: privateKey,
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = lego.LEDirectoryStaging // TODO take this out after testing
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Errorf("Error creating lets encrypt client: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	client.Challenge.Remove(challenge.HTTP01)
	client.Challenge.Remove(challenge.TLSALPN01)
	trafficRouterDns, err := NewDNSProviderTrafficRouter()
	trafficRouterDns.db = db
	if err != nil {
		log.Errorf("Error creating Traffic Router DNS provider: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
	}
	client.Challenge.SetDNS01Provider(trafficRouterDns)

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Errorf("Error registering lets encrypt client: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domainName},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Errorf("Error obtaining lets encrypt certificate: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	fmt.Printf("%#v\n", certificates)

	block, _ := pem.Decode([]byte(certificates.Certificate))
	if block == nil {
		log.Errorf("Error parsing cert - %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}
	x509cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Errorf("Error parsing cert to get expiry - %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		return
	}

	expiration := x509cert.NotAfter
	log.Errorf("MATT JACKSON - expiration = %s", expiration)

	// Save certs into Riak
	dsSSLKeys := tc.DeliveryServiceSSLKeys{
		CDN:             *req.CDN,
		DeliveryService: *req.DeliveryService,
		Hostname:        *req.HostName,
		Version:         *req.Version,
	}

	dsSSLKeys.Certificate = tc.DeliveryServiceSSLKeysCertificate{Crt: string(EncodePEMToLegacyPerlRiakFormat(certificates.Certificate)), Key: string(EncodePEMToLegacyPerlRiakFormat(certificates.PrivateKey)), CSR: "Not Applicable"}
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(dsSSLKeys, tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort); err != nil {
		log.Errorf("Error posting lets encrypt certificate to riak: %s", err.Error())
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, errors.New("putting riak keys: "+err.Error()), nil)
		return
	}
	api.WriteResp(w, r, "Successfully created ssl keys for "+deliveryService)

}

func RenewLetsEncryptCert(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	req := tc.DeliveryServiceLetsEncryptSSLKeysReq{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &req); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}
	deliveryService := *req.DeliveryService

	version := inf.Params["version"]
	xmlID := inf.Params["xmlId"]
	keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(xmlID, version, inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ssl keys: "+err.Error()))
		return
	}
	if !ok {
		api.WriteRespAlertObj(w, r, tc.InfoLevel, "no object found for the specified key", struct{}{}) // empty response object because Perl
		return
	}

	oldCert := certificate.Resource{
		PrivateKey:  []byte(keyObj.Certificate.Key),
		Certificate: []byte(keyObj.Certificate.Crt),
		CSR:         []byte(keyObj.Certificate.CSR),
	}

	myUser := MyUser{
		key: keyObj.Key,
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = lego.LEDirectoryStaging // TODO take this out after testing
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}

	client.Challenge.Remove(challenge.HTTP01)
	client.Challenge.Remove(challenge.TLSALPN01)
	trafficRouterDns, err := NewDNSProviderTrafficRouter()
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
	}
	client.Challenge.SetDNS01Provider(trafficRouterDns)

	renewedCert, err := client.Certificate.Renew(oldCert, true, false)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, errors.New("renewing certificate: "+err.Error()), nil)
		return
	}

	// Save certs into Riak
	dsSSLKeys := tc.DeliveryServiceSSLKeys{
		CDN:             *req.CDN,
		DeliveryService: *req.DeliveryService,
		BusinessUnit:    *req.BusinessUnit,
		City:            *req.City,
		Organization:    *req.Organization,
		Hostname:        *req.HostName,
		Country:         *req.Country,
		State:           *req.State,
		Key:             *req.Key,
		Version:         *req.Version,
	}

	dsSSLKeys.Certificate = tc.DeliveryServiceSSLKeysCertificate{Crt: string(renewedCert.Certificate), Key: string(renewedCert.PrivateKey), CSR: string(renewedCert.CSR)}
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(dsSSLKeys, inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, errors.New("putting riak keys: "+err.Error()), nil)
		return
	}

	api.WriteResp(w, r, "Successfully renewed ssl keys for "+deliveryService)

}
