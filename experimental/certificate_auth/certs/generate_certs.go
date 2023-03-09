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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net"
	"time"
)

// CertificateKeyPair contains the parsed representation of a certificate
// and private key.
type CertificateKeyPair struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

// CertificatePEMPair contains the PEM encoded certificate and private key.
type CertificatePEMPair struct {
	CertificatePEM, PrivateKeyPEM string
}

var (
	rootCN   = "root.local"
	interCN  = "intermediate.local"
	clientCN = "client.local"
	serverCN = "server.local"

	uid = "userid"
	// useEcdsa = false //TODO: Enable and refactor
)

func main() {
	flag.StringVar(&uid, "uid", uid, "[Optional] The User ID value to be added to the client certificate")

	flag.Parse()

	rootCAPEMPair, err := GenerateRootCACertificate()
	if err != nil {
		log.Fatalf("Failed to generate and sign Root CA certificate\nErr: %s\n", err)
	}
	ioutil.WriteFile("rootca.crt.pem", []byte(rootCAPEMPair.CertificatePEM), 0644)
	ioutil.WriteFile("rootca.key.pem", []byte(rootCAPEMPair.PrivateKeyPEM), 0644)

	intermediatePEMPair, err := GenerateIntermediateCertificate(rootCAPEMPair)
	if err != nil {
		log.Fatalf("Failed to generate and sign Intermediate certificate\nErr: %s\n", err)
	}
	ioutil.WriteFile("intermediate.crt.pem", []byte(intermediatePEMPair.CertificatePEM), 0644)
	ioutil.WriteFile("intermediate.key.pem", []byte(intermediatePEMPair.PrivateKeyPEM), 0644)

	serverPEMPair, err := GenerateServerCertificate(intermediatePEMPair)
	if err != nil {
		log.Fatalf("Failed to generate and sign Server certificate\nErr: %s\n", err)
	}

	ioutil.WriteFile("server.crt.pem", []byte(serverPEMPair.CertificatePEM), 0644)
	ioutil.WriteFile("server.key.pem", []byte(serverPEMPair.PrivateKeyPEM), 0644)

	clientPEMPair, err := GenerateClientCertificate(intermediatePEMPair)
	if err != nil {
		log.Fatalf("Failed to generate and sign Client certificate\nErr: %s\n", err)
	}

	ioutil.WriteFile("client.crt.pem", []byte(clientPEMPair.CertificatePEM), 0644)
	ioutil.WriteFile("client.key.pem", []byte(clientPEMPair.PrivateKeyPEM), 0644)

	clientIntermediateChain := clientPEMPair.CertificatePEM + intermediatePEMPair.CertificatePEM
	ioutil.WriteFile("client-intermediate-chain.crt.pem", []byte(clientIntermediateChain), 0644)

	if err := VerifyCertificates(rootCAPEMPair, intermediatePEMPair, clientPEMPair, serverPEMPair); err != nil {
		log.Fatalf("failed to verify certificate: %s", err)
	}
}

// ParseCertificateKeyPair decodes the provided PEM pair (key, cert) and returns a
// parsed private key and x509 certificate.
func ParseCertificateKeyPair(pemPair *CertificatePEMPair) (*CertificateKeyPair, error) {

	keyPair := new(CertificateKeyPair)

	privPemBlock, _ := pem.Decode([]byte(pemPair.PrivateKeyPEM))

	privateKey, err := x509.ParsePKCS8PrivateKey(privPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	keyPair.PrivateKey = privateKey.(*rsa.PrivateKey)

	certPemBlock, _ := pem.Decode([]byte(pemPair.CertificatePEM))

	certificate, err := x509.ParseCertificate(certPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	keyPair.Certificate = certificate

	return keyPair, nil
}

// GenereateRootCACertificate creates a Root CA certificate that can be used
// for signing intermediate, client, and server x509 certificates.
func GenerateRootCACertificate() (*CertificatePEMPair, error) {

	now := time.Now()

	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			OrganizationalUnit: []string{"ATC"},
			Organization:       []string{"Apache"},
			Country:            []string{"US"},
			Province:           []string{"Colorado"},
			Locality:           []string{"Denver"},
			CommonName:         rootCN,
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEMPair := new(CertificatePEMPair)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDERBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	certPrivKeyByes, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key to PKCS8: %w", err)
	}

	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: certPrivKeyByes,
	})

	certPEMPair.CertificatePEM = certPEM.String()
	certPEMPair.PrivateKeyPEM = certPrivKeyPEM.String()

	return certPEMPair, nil
}

// GenerateIntermediateCeertificate creates an intermediate based on the provided Root certificate.
// This certificate can be used for signing client and server certificates to establish
// a chain to the Root certificate.
func GenerateIntermediateCertificate(root *CertificatePEMPair) (*CertificatePEMPair, error) {

	rootKeyPair, err := ParseCertificateKeyPair(root)
	if err != nil {
		log.Fatalln("Failed to parse root cert and key")
	}

	now := time.Now()

	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			OrganizationalUnit: []string{"ATC"},
			Organization:       []string{"Apache"},
			Country:            []string{"US"},
			Province:           []string{"Colorado"},
			Locality:           []string{"Denver"},
			CommonName:         interCN,
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(1, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		MaxPathLenZero:        true,
		IsCA:                  true,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, cert, rootKeyPair.Certificate, &certPrivKey.PublicKey, rootKeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEMPair := new(CertificatePEMPair)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDERBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	certPrivKeyByes, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key to PKCS8: %w", err)
	}

	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: certPrivKeyByes,
	})

	certPEMPair.CertificatePEM = certPEM.String()
	certPEMPair.PrivateKeyPEM = certPrivKeyPEM.String()

	return certPEMPair, nil
}

// GenerateClientCertificate creates and signs a certificate based on the provided RootCA. This differs
// from the Server certificate in that it includes the OID for LDAP UID as well as Client Auth key usage.
//
// Currently the key is an RSA key, which also entails adding KeyEncipherment key usage.
func GenerateClientCertificate(intermediate *CertificatePEMPair) (*CertificatePEMPair, error) {

	intermediateKeyPair, err := ParseCertificateKeyPair(intermediate)
	if err != nil {
		log.Fatalln("Failed to parse root cert and key")
	}

	now := time.Now()

	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	// LDAP OID reference: https://ldap.com/ldap-oid-reference-guide/
	// 0.9.2342.19200300.100.1.1 	uid Attribute Type
	uidPkix := pkix.AttributeTypeAndValue{
		Type:  asn1.ObjectIdentifier([]int{0, 9, 2342, 19200300, 100, 1, 1}),
		Value: uid,
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			OrganizationalUnit: []string{"ATC"},
			Organization:       []string{"Apache"},
			Country:            []string{"US"},
			Province:           []string{"Colorado"},
			Locality:           []string{"Denver"},
			CommonName:         clientCN,
			ExtraNames:         []pkix.AttributeTypeAndValue{uidPkix},
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, cert, intermediateKeyPair.Certificate, &certPrivKey.PublicKey, intermediateKeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEMPair := new(CertificatePEMPair)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDERBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	certPrivKeyByes, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key to PKCS8: %w", err)
	}

	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: certPrivKeyByes,
	})

	certPEMPair.CertificatePEM = certPEM.String()
	certPEMPair.PrivateKeyPEM = certPrivKeyPEM.String()

	return certPEMPair, nil
}

// GenerateServerCertificate creates and signs a certificate based on the provided RootCA. This differs
// from the Client certificate in that it ServerAuth key usage. It also does NOT include the OID for LDAP UID.
//
// Currently the key is an RSA key, which also entails adding KeyEncipherment key usage.
func GenerateServerCertificate(intermediate *CertificatePEMPair) (*CertificatePEMPair, error) {

	intermediateKeyPair, err := ParseCertificateKeyPair(intermediate)
	if err != nil {
		log.Fatalln("Failed to parse root cert and key")
	}

	now := time.Now()

	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			OrganizationalUnit: []string{"ATC"},
			Organization:       []string{"Apache"},
			Country:            []string{"US"},
			Province:           []string{"Colorado"},
			Locality:           []string{"Denver"},
			CommonName:         serverCN,
		},
		DNSNames:    []string{serverCN},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   now,
		NotAfter:    now.AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, cert, intermediateKeyPair.Certificate, &certPrivKey.PublicKey, intermediateKeyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEMPair := new(CertificatePEMPair)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDERBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	certPrivKeyByes, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key to PKCS8: %w", err)
	}

	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: certPrivKeyByes,
	})

	certPEMPair.CertificatePEM = certPEM.String()
	certPEMPair.PrivateKeyPEM = certPrivKeyPEM.String()

	return certPEMPair, nil
}

// VerifyCertificates checks that the client and server certificates match the
// Root and Intermediate chains.
func VerifyCertificates(root, intermediate, client, server *CertificatePEMPair) error {

	rootKeyPair, err := ParseCertificateKeyPair(root)
	if err != nil {
		log.Fatalln("Failed to parse root cert and key")
	}
	intermediateKeyPair, err := ParseCertificateKeyPair(intermediate)
	if err != nil {
		log.Fatalln("Failed to parse intermediate cert and key")
	}

	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootKeyPair.Certificate)
	intermediatePool := x509.NewCertPool()
	intermediatePool.AddCert(intermediateKeyPair.Certificate)

	opts := x509.VerifyOptions{
		Intermediates: intermediatePool,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Roots:         rootPool,
	}

	clientCert, err := ParseCertificateKeyPair(client)
	if err != nil {
		return fmt.Errorf("failed to parse client cert and key: %w", err)
	}

	if _, err := clientCert.Certificate.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify client cert and key: %w", err)
	}

	serverCert, err := ParseCertificateKeyPair(server)
	if err != nil {
		return fmt.Errorf("failed to parse server cert and key: %w", err)
	}

	if _, err := serverCert.Certificate.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify client cert and key: %w", err)
	}

	return nil
}
