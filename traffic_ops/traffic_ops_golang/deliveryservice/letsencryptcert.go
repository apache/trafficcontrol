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

	req := tc.DeliveryServiceLetsEncryptSSLKeysReq{}
	if err := api.Parse(r.Body, nil, &req); err != nil {
		log.Errorf("Error parsing request: %s", err.Error())
		api.HandleErr(w, r, nil, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}

	error := GetLetsEncryptCertificates(inf, req, ctx)

	if error != nil {
		api.HandleErr(w, r, nil, http.StatusInternalServerError, error, nil)
	}

	api.WriteResp(w, r, "Successfully created ssl keys for "+*req.DeliveryService)

}

func GetLetsEncryptCertificates(inf *api.APIInfo, req tc.DeliveryServiceLetsEncryptSSLKeysReq, ctx context.Context) error {

	db, err := api.GetDB(ctx)
	tx, err := db.BeginTxx(ctx, nil)

	domainName := *req.HostName
	deliveryService := *req.DeliveryService

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorf(deliveryService+": Error generating private key: %s", err.Error())
		return err
	}

	myUser := MyUser{
		key:   privateKey,
		Email: inf.Config.ConfigLetsEncrypt.Email,
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = lego.LEDirectoryStaging // TODO take this out after testing
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Errorf(deliveryService+": Error creating lets encrypt client: %s", err.Error())
		return err
	}

	client.Challenge.Remove(challenge.HTTP01)
	client.Challenge.Remove(challenge.TLSALPN01)
	trafficRouterDns, err := NewDNSProviderTrafficRouter()
	trafficRouterDns.db = db
	if err != nil {
		log.Errorf(deliveryService+": Error creating Traffic Router DNS provider: %s", err.Error())
		return err
	}
	client.Challenge.SetDNS01Provider(trafficRouterDns)

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Errorf(deliveryService+": Error registering lets encrypt client: %s", err.Error())
		return err
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domainName},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Errorf(deliveryService+": Error obtaining lets encrypt certificate: %s", err.Error())
		return err
	}

	fmt.Printf("%#v\n", certificates)

	block, _ := pem.Decode([]byte(certificates.Certificate))
	if block == nil {
		log.Errorf(deliveryService + ": Error parsing cert")
		return errors.New(deliveryService + ": parsing cert")
	}
	x509cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Errorf(deliveryService+": Error parsing cert to get expiry - %s", err.Error())
		return err
	}

	expiration := x509cert.NotAfter

	// Save certs into Riak
	dsSSLKeys := tc.DeliveryServiceSSLKeys{
		AuthType:        tc.LetsEncryptAuthType,
		CDN:             *req.CDN,
		DeliveryService: *req.DeliveryService,
		Hostname:        *req.HostName,
		Version:         *req.Version,
		Expiration:      expiration,
	}

	crtBuf := bytes.Buffer{}
	if err := pem.Encode(&crtBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certificates.Certificate}); err != nil {
		log.Errorf(deliveryService + ": pem-encoding certificate: " + err.Error())
		return errors.New(deliveryService + ": pem-encoding certificate: " + err.Error())
	}
	crtPem := crtBuf.Bytes()

	keyBuf := bytes.Buffer{}
	if err := pem.Encode(&keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: certificates.PrivateKey}); err != nil {
		log.Errorf(deliveryService + ": pem-encoding key: " + err.Error())
		return errors.New(deliveryService + ": pem-encoding key: " + err.Error())
	}
	keyPem := keyBuf.Bytes()

	dsSSLKeys.Certificate = tc.DeliveryServiceSSLKeysCertificate{Crt: string(EncodePEMToLegacyPerlRiakFormat(crtPem)), Key: string(EncodePEMToLegacyPerlRiakFormat(keyPem)), CSR: "Not Applicable"}
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(dsSSLKeys, tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort); err != nil {
		log.Errorf("Error posting lets encrypt certificate to riak: %s", err.Error())
		return errors.New(deliveryService + ": putting riak keys: " + err.Error())
	}

	return nil
}
