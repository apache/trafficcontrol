package deliveryservice

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/lego"
	"net/http"
)

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
	"github.com/go-acme/lego/registration"
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

func GenerateLetsEncryptCertificates(w http.ResponseWriter, r *http.Request) {
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

	domainName := *req.HostName
	deliveryService := *req.DeliveryService

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
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
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domainName},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, nil)
		return
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)

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
	dsSSLKeys.Certificate = tc.DeliveryServiceSSLKeysCertificate{Crt: string(certificates.Certificate), Key: string(certificates.PrivateKey), CSR: string(certificates.CSR)}
	if err := riaksvc.PutDeliveryServiceSSLKeysObj(dsSSLKeys, inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort); err != nil {
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
