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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

const NewCertValidDuration = time.Hour * 24 * 365

// GenerateCert generates a key and certificate for serving HTTPS. The generated key is 2048-bit RSA, to match the old Perl code.
// The certificate will be valid for NewCertValidDuration time after now.
// Returns PEM-encoded certificate signing request (csr), certificate (crt), and key; or any error.
func GenerateCert(host, country, city, state, org, unit string) ([]byte, []byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, errors.New("generating key: " + err.Error())
	}
	now := time.Now()
	expires := now.Add(NewCertValidDuration)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, nil, errors.New("getting random int for serial number: " + err.Error())
	}

	subj := pkix.Name{
		CommonName:         host,
		Country:            []string{country},
		Province:           []string{state},
		Locality:           []string{city},
		Organization:       []string{org},
		OrganizationalUnit: []string{unit},
	}

	crt := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subj,
		NotBefore:             now,
		NotAfter:              expires,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host},
		Version:               1,
	}

	crtReq := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
		Version:            1,
	}

	crtDer, err := x509.CreateCertificate(rand.Reader, &crt, &crt, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, nil, errors.New("creating certificate: " + err.Error())
	}
	crtBuf := bytes.Buffer{}
	if err := pem.Encode(&crtBuf, &pem.Block{Type: "CERTIFICATE", Bytes: crtDer}); err != nil {
		return nil, nil, nil, errors.New("pem-encoding certificate: " + err.Error())
	}
	crtPem := crtBuf.Bytes()

	csrDer, err := x509.CreateCertificateRequest(rand.Reader, &crtReq, priv)
	if err != nil {
		return nil, nil, nil, errors.New("creating certificate request: " + err.Error())
	}
	csrBuf := bytes.Buffer{}
	if err := pem.Encode(&csrBuf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDer}); err != nil {
		return nil, nil, nil, errors.New("pem-encoding certificate request: " + err.Error())
	}
	csrPem := csrBuf.Bytes()

	keyDer := x509.MarshalPKCS1PrivateKey(priv)
	if keyDer == nil {
		return nil, nil, nil, errors.New("marshalling private key: nil der")
	}
	keyBuf := bytes.Buffer{}
	if err := pem.Encode(&keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDer}); err != nil {
		return nil, nil, nil, errors.New("pem-encoding private key: " + err.Error())
	}
	keyPem := keyBuf.Bytes()

	return EncodePEMToLegacyPerlRiakFormat(csrPem), EncodePEMToLegacyPerlRiakFormat(crtPem), EncodePEMToLegacyPerlRiakFormat(keyPem), nil
}

// EncodePEMToLegacyPerlRiakFormat takes a PEM-encoded byte (typically a certificate, csr, or key) and returns the format Perl Traffic Ops used to send to Riak.
func EncodePEMToLegacyPerlRiakFormat(pem []byte) []byte {
	b64Pem := []byte(base64.StdEncoding.EncodeToString(pem)) // Why are we base64-encoding a base64-encoded format? Because Perl
	b64Lines := util.BytesLenSplit(b64Pem, 76)               // Why 76? Because Perl
	joined := bytes.Join(b64Lines, []byte{'\n'})             // Why are we joining arbitrary base64-encoded characters with newlines? Because Perl
	return joined
}
