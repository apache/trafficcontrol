package auth

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
	"encoding/pem"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc/ldap"
)

// VerifyClientCertificate takes a http.Request, pulls the (optionally) provided client TLS
// certificates and attempts to verify them against the directory of provided Root CA
// certificates. The Root CA certificates can be different than those utilized by the
// http.Server. Returns an error if the verification process fails
func VerifyClientCertificate(r *http.Request, rootCertsDirPath string, insecureSkipVerify bool) error {
	// TODO: Parse client headers as alternative to TLS in the request

	if err := loadRootCerts(rootCertsDirPath); err != nil {
		return fmt.Errorf("failed to load root certificates")
	}

	if err := verifyClientRootChain(r.TLS.PeerCertificates, insecureSkipVerify); err != nil {
		return fmt.Errorf("failed to verify client to root certificate chain")
	}

	return nil
}

func verifyClientRootChain(clientChain []*x509.Certificate, insecureSkipVerify bool) error {
	if len(clientChain) == 0 {
		return fmt.Errorf("empty client chain")
	}

	if rootPool == nil {
		return fmt.Errorf("uninitialized root cert pool")
	}

	intermediateCertPool := x509.NewCertPool()
	for _, intermediate := range clientChain[1:] {
		intermediateCertPool.AddCert(intermediate)
	}

	opts := x509.VerifyOptions{
		Intermediates: intermediateCertPool,
		Roots:         rootPool,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	_, err := clientChain[0].Verify(opts)
	if err != nil {
		if insecureSkipVerify {
			return nil
		}
		return fmt.Errorf("failed to verify client cert chain. err: %w", err)
	}
	return nil
}

// Lazy initialized
var rootPool *x509.CertPool

func loadRootCerts(dirPath string) error {
	// Root cert pool already populated
	// TODO: This will prevent rolling cert renewals at runtime and will require a TO restart
	// to pick up additional certificates.
	if rootPool != nil {
		return nil
	}

	if dirPath == "" {
		return fmt.Errorf("empty path supplied for root cert directory")
	}

	err := filepath.WalkDir(dirPath,
		// walk function to perform on each file in the supplied
		// directory path for root certificiates.
		//
		// For each file in the directory, first check if it, too, is a dir. If so,
		// return the filepath.SkipDir error to allow for it to be skipped without
		// stopping the subsequent executions.
		//
		// If of type File, then load the PEM encoded string from the file and
		// attempt to decode the PEM block into an x509 certificate. If successful,
		// add that certificate to the Root Cert Pool to be used for verification.
		//
		// Must be a closure for access to the `dirPath` value
		func(path string, file fs.DirEntry, e error) error {
			if e != nil {
				return e
			}

			// Skip logic if root directory
			if path == dirPath {
				return nil
			}

			// Don't traverse nested directories
			if file.IsDir() {
				return filepath.SkipDir
			}

			if info, err := file.Info(); err != nil {
				return fmt.Errorf("getting info for file %s: %s", file.Name(), err)
			} else if groupWritable := fs.FileMode(020); info.Mode().Perm()&groupWritable == groupWritable {
				log.Errorf("refusing to use group-writable file err: %s", err)
				return nil
			} else if worldWritable := fs.FileMode(002); info.Mode().Perm()&worldWritable == worldWritable {
				log.Errorf("refusing to use world-writable file err: %s", err)
				return nil
			}

			// Read file
			pemBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to open cert at %s. err: %w", path, err)
			}
			pemBlock, _ := pem.Decode(pemBytes)
			// Failed to decode PEM, skip file
			if pemBlock == nil {
				return nil
			}
			certificate, err := x509.ParseCertificate(pemBlock.Bytes)
			if err != nil {
				log.Errorf("failed to parse PEM into x509. err: %s", err)
				return nil
			}

			if rootPool == nil {
				rootPool = x509.NewCertPool()
			}
			rootPool.AddCert(certificate)

			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to load root certs from path %s. err: %s", dirPath, err)
	}

	return nil
}

// ParseClientCertificateUID takes an x509 Certificate and loops through the Names in the
// Subject. If it finds an asn.ObjectIdentifier that matches UID, it returns the
// corresponding value. Otherwise returns empty string. If more than one UID is present,
// the first result found to match is returned (order not guaranteed).
func ParseClientCertificateUID(cert *x509.Certificate) (string, error) {
	foundUID := false
	uid := ""
	err := error(nil)
	for _, name := range cert.Subject.Names {
		if name.Type.Equal(ldap.OIDType) {
			if foundUID {
				err = fmt.Errorf("found more than 1 UID in certificate subject")
				break
			}
			uid = name.Value.(string)
			foundUID = true
		}
	}
	if !foundUID {
		err = fmt.Errorf("no UID found")
	}
	return uid, err
}
