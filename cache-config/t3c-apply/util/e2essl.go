package util

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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-log"
)

const E2ESSLCADestPathCert = "e2e-ssl-ca.cert"

const E2ESSLPathBase = "e2e-ssl"

const E2ESSLPathClientBase = E2ESSLPathBase + "-client"

const E2ESSLPathServerBase = E2ESSLPathBase + "-server"

// e2eSSLCertLifetime is the lifetime to generate for certificates. TODO make configurable.
const e2eSSLCertLifetime = time.Hour * 24 * 7

// e2eSSLCertRefreshAge is the age to refresh certificates after. TODO make configurable.
const e2eSSLCertRefreshAge = time.Hour * 24

// E2ESSLKeysExist returns nil if the E2E SSL keys exist, or an error if they don't exist or we were unable to determine.
func E2ESSLKeysExist(certDir string) error {
	// TODO add verifying the files are actually keys?
	// That would probably be easiest to do by externally calling openssl (rather than Go parsing).
	files := []string{
		filepath.Join(certDir, E2ESSLPathBase+".key"),
		// filepath.Join(certDir, E2ESSLPathClientBase+".csr"),
		filepath.Join(certDir, E2ESSLPathClientBase+".cert"),
		// filepath.Join(certDir, E2ESSLPathServerBase+".csr"),
		filepath.Join(certDir, E2ESSLPathServerBase+".cert"),
	}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Infof("E2ESSLKeysExist path %v doesn't exist, returning err IsNotExist\n", file)
			return err
		} else if err != nil {
			log.Infoln("E2ESSLKeysExist returning real error")
			return err
		}
		log.Infof("E2ESSLKeysExist path %v exists\n", file)
	}
	log.Infoln("E2ESSLKeysExist returning nil, all paths exist")
	return nil
}

func E2ESSLGenerateKeys(certDir string, caKeyPath string, caCertPath string) error {
	clientKey := filepath.Join(certDir, E2ESSLPathBase+".key")
	clientCSR := filepath.Join(certDir, E2ESSLPathClientBase+".csr")
	clientCert := filepath.Join(certDir, E2ESSLPathClientBase+".cert")
	serverCSR := filepath.Join(certDir, E2ESSLPathServerBase+".csr")
	serverCert := filepath.Join(certDir, E2ESSLPathServerBase+".cert")

	// client cert private key
	if stdOut, stdErr, code := t3cutil.Do("openssl", "ecparam", "-name", "secp256r1", "-genkey", "-noout", "-out", clientKey); code != 0 {
		return fmt.Errorf("generating client private key returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	// TODO only get once for the app. Is it needed anywhere else?

	hostnameFQDN, err := GetHostnameFQDN()
	if err != nil {
		return errors.New("getting hostname: " + err.Error())
	}

	if err := E2ESSLGenerateCert(caCertPath, caKeyPath, clientKey, clientCSR, clientCert, hostnameFQDN); err != nil {
		return errors.New("generating client cert: " + err.Error())
	}
	if err := E2ESSLGenerateCert(caCertPath, caKeyPath, clientKey, serverCSR, serverCert, "*"); err != nil {
		return errors.New("generating server cert: " + err.Error())
	}
	return nil
}

func E2ESSLGenerateCert(caCertPath string, caKeyPath string, clientKeyPath string, csrPath string, certPath string, certCN string) error {
	certLifetimeDaysStr := strconv.Itoa(int(e2eSSLCertLifetime / time.Hour / 24))

	log.Infof("E2ESSLGenerateCert calling -days %v\n", certLifetimeDaysStr)

	if stdOut, stdErr, code := t3cutil.Do("openssl", "req", "-new", "-sha256", "-key", clientKeyPath, "-subj", "/C=US/ST=CO/O=ApacheTrafficControl/CN="+certCN, "-out", csrPath); code != 0 {
		return fmt.Errorf("generating csr returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if stdOut, stdErr, code := t3cutil.Do("openssl", "x509", "-req", "-in", csrPath, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-out", certPath, "-days", certLifetimeDaysStr, "-sha256"); code != 0 {
		return fmt.Errorf("generating certificate returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	return nil
}

// TODO roll client key
// TODO check CSR for expiration and warn
// TODO log generation/refreshes
// TODO add sha512sum, openssl to RPM dependencies
// TODO refactor to not log. This should act like a library

// E2ESSLRefreshCerts generates new certificates if they're about to expire.
// The mustRefresh parameter forces new generation, regardless of expiration. This is most often used when the Certificate Authority has changed.
func E2ESSLRefreshCerts(certDir string, caKeyPath string, caCertPath string, mustRefresh bool) error {
	clientKey := filepath.Join(certDir, E2ESSLPathBase+".key")
	clientCSR := filepath.Join(certDir, E2ESSLPathClientBase+".csr")
	serverCSR := filepath.Join(certDir, E2ESSLPathServerBase+".csr")
	clientCert := filepath.Join(certDir, E2ESSLPathClientBase+".cert")
	serverCert := filepath.Join(certDir, E2ESSLPathServerBase+".cert")

	if err := E2ESSLRefreshCert(caKeyPath, caCertPath, mustRefresh, clientKey, clientCert, clientCSR); err != nil {
		return errors.New("refreshing client cert: " + err.Error())
	}
	if err := E2ESSLRefreshCert(caKeyPath, caCertPath, mustRefresh, clientKey, serverCert, serverCSR); err != nil {
		return errors.New("refreshing server cert: " + err.Error())
	}
	return nil
}

func e2eCAChanged(caKeyPath string, caCertPath string) (bool, error) {
	caKeyChanged, err := e2eFileChanged(caKeyPath)
	if err != nil {
		return false, errors.New("checking if '" + caKeyPath + "' changed: " + err.Error())
	}

	caCertChanged, err := e2eFileChanged(caCertPath)
	if err != nil {
		return false, errors.New("checking if '" + caCertPath + "' changed: " + err.Error())
	}

	log.Infof("CA cert %v change %v key %v change %v\n", caCertPath, caCertChanged, caKeyPath, caKeyChanged)

	return caKeyChanged || caCertChanged, nil
}

// e2eFileChanged determines whether the given file has changed.
// It does this by reading and writing a path.hash file.
//
// Note this creates a checksum of the path's file name in the app var/lib directory.
// Therefore, this function must not be used for multiple files of the same name in different directories, or they will attempt to share checksums.
//
func e2eFileChanged(path string) (bool, error) {
	caKeyHashFileName := filepath.Join(config.VarDir, filepath.Base(path)+".sha512sum")
	caKeyHash, err := ioutil.ReadFile(caKeyHashFileName)
	if err != nil {
		if os.IsNotExist(err) {
			caKeyHash = []byte("") // set it empty so it creates a diff
		} else {
			return true, errors.New("error checking if '" + caKeyHashFileName + "' exists: " + err.Error())
		}
	}
	oldHash := strings.TrimSpace(string(caKeyHash))

	changed := false

	stdOut, stdErr, code := t3cutil.Do("sha512sum", path)
	if code != 0 {
		return true, fmt.Errorf("creating hash of '"+path+"'  returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	newHash := strings.TrimSpace(string(stdOut))

	if newHash != oldHash {
		changed = true
		if err := ioutil.WriteFile(caKeyHashFileName, []byte(newHash), 0600); err != nil {
			return true, errors.New("writing hash file '" + caKeyHashFileName + "': " + err.Error())
		}
	}
	return changed, nil
}

func E2ESSLRefreshCert(caKeyPath string, caCertPath string, mustRefresh bool, keyPath string, certPath string, csrPath string) error {
	if !mustRefresh {
		expiration, err := GetCertExpiration(certPath)
		if err != nil {
			return errors.New("getting cert '" + certPath + "' expiration: " + err.Error())
		}
		if time.Now().Add(e2eSSLCertRefreshAge).After(expiration) {
			log.Infof("Cert %v age %v exceeds threshold %v, refreshing\n", certPath, expiration, e2eSSLCertRefreshAge)
			mustRefresh = true
		} else {
			log.Infof("Cert %v age %v under threshold %v, not refreshing\n", certPath, expiration, e2eSSLCertRefreshAge)
		}
	}
	if mustRefresh {
		certCN, err := GetCertCN(certPath)
		if err != nil {
			return errors.New("getting cert '" + certPath + "' CN: " + err.Error())
		}
		if err := E2ESSLGenerateCert(caCertPath, caKeyPath, keyPath, csrPath, certPath, certCN); err != nil {
			return errors.New("generating cert: " + err.Error())
		}
	}
	return nil
}

func GetCertExpiration(certPath string) (time.Time, error) {
	stdOut, stdErr, code := t3cutil.Do("openssl", "x509", "-in", certPath, "-noout", "-enddate")
	if code != 0 {
		return time.Time{}, fmt.Errorf("getting certificate '"+certPath+"'  returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	naDate := strings.Split(strings.TrimSpace(string(stdOut)), "=")
	if len(naDate) != 2 {
		return time.Time{}, fmt.Errorf("getting certificate '"+certPath+"' returned code 0 but unexpected format, expected notAfter=date actual stdout '%v' stderr '%v'", string(stdOut), string(stdErr))
	}
	dateStr := naDate[1]

	dateFormat := `Jan _2 15:04:05 2006 MST`
	tm, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("getting certificate '"+certPath+"' returned code 0 but failed to parse, stdout '%v' stderr '%v' parse error: %v", string(stdOut), string(stdErr), err)
	}
	return tm, nil
}

func GetCertCN(certPath string) (string, error) {
	stdOut, stdErr, code := t3cutil.Do("openssl", "x509", "-in", certPath, "-noout", "-subject")
	if code != 0 {
		return "", fmt.Errorf("getting certificate '"+certPath+"' returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	// replace all spaces. Some certs have the form "CN = foo.com".
	subjectStr := strings.Replace(strings.TrimSpace(string(stdOut)), " ", "", -1)
	subjectFields := strings.Split(subjectStr, ",")
	cn := ""
	for _, field := range subjectFields {
		if strings.HasPrefix(field, "CN=") {
			cnVal := strings.Split(field, "=")
			if len(cnVal) > 1 {
				cn = cnVal[1]
			}
			break
		}
	}
	if cn == "" {
		cn = "*" // if the cert had no CN, replace with a wildcard. TODO warn?
	}
	return cn, nil
}

func GetHostnameFQDN() (string, error) {
	stdOut, stdErr, code := t3cutil.Do("hostname", "--fqdn")
	if code != 0 {
		return "", fmt.Errorf("getting 'hostname --fqdn' returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	return strings.TrimSpace(string(stdOut)), nil
}

func E2ESSLGenerateKeysIfNotExist(certDir string, caKeyPath string, caCertPath string, e2eCACertDestPath string) error {
	caIsNew, err := E2ECheckAndHandleNewCA(caCertPath, caKeyPath, e2eCACertDestPath)
	if err != nil {
		return errors.New("checking and refreshing new CA: " + err.Error())
	}
	mustRefresh := caIsNew // if the CA is new, we must regenerate certificates

	if err := E2ESSLKeysExist(certDir); err == nil {
		return E2ESSLRefreshCerts(certDir, caKeyPath, caCertPath, mustRefresh)
	} else if !os.IsNotExist(err) {
		return errors.New("checking if keys exist: " + err.Error())
	}
	if err := E2ESSLGenerateKeys(certDir, caKeyPath, caCertPath); err != nil {
		return errors.New("generating keys: " + err.Error())
	}
	return nil
}

var E2ESSLOldCAPath = filepath.Join(config.VarDir, config.DefaultE2ESSLCACertFileName+".old")
var E2ESSLNewCAPath = filepath.Join(config.VarDir, config.DefaultE2ESSLCACertFileName+".new")
var E2ESSLOldCAKeyPath = filepath.Join(config.VarDir, config.DefaultE2ESSLCAKeyFileName+".old")
var E2ESSLNewCAKeyPath = filepath.Join(config.VarDir, config.DefaultE2ESSLCAKeyFileName+".new")

// E2ECheckAndHandleNewCA handles copying the CA from var/lib to etc/trafficserver
// This is not trivial, because ATS must continue to serve the previous CA
// This function:
// 1. For both the CA cert and key files, checks if they have changed, or if this is the first run. (by checking for cafilename.new and/or cafilename.new.sha512). If so:
// 1.1. moves any existing cafilename.new to cafilename.old
// 1.2. creates cafilename.new to cafilename.old
// 1.3. creates cafilename.new.sha512
// 1.5. concatenates cafilename.new and cafilename.old, and writes the concatenated file to etc/trafficserver/ssl/e2e-ssl-ca.{cert|key}
//
// The caPath is the path to the certificate authority, placed by something other than t3c (manual, Ansible, etc).
// The caDestPath is the path where the CA will be placed, to be read by ATS.
//
// Returns whether the CA cert or key was changed (and therefore whether e2e certs will need regenerated).
//
func E2ECheckAndHandleNewCA(caPath string, caKeyPath string, caDestPath string) (bool, error) {
	// we need the var dir for the old and new cert copies, and checksums
	// Try to make it, if it doesn't exist
	if err := os.MkdirAll(config.VarDir, 0644); err != nil {
		return false, errors.New("ensuring '" + config.VarDir + "': " + err.Error())
	}

	certChanged, err := e2eFileChanged(caPath)
	if err != nil {
		return false, errors.New("checking ca '" + caPath + "' changed: " + err.Error())
	}

	// Note we check if the key changed, even if the cert changed and it won't affect the changed variable
	// to ensure the checksum is created and updated
	keyChanged, err := e2eFileChanged(caKeyPath)
	if err != nil {
		return false, errors.New("checking ca key '" + caPath + "' changed: " + err.Error())
	}

	changed := certChanged || keyChanged

	if !changed {
		log.Infoln("E2ECheckAndHandleNewCA: ca and key are unchanged")
		return false, nil
	}

	log.Infoln("E2ECheckAndHandleNewCA: ca or key changed, rewriting ca")

	// the CA changed, so we need to
	// 1. copy the ca.new to ca.old
	// 2. copy the ca to ca.new
	// 3. copy the cakey.new to cakey.old
	// 4. copy the cakey to cakey.new
	// 5. cat ca.new ca.old > etc/trafficserver/ca.cert

	// TODO don't read twice (the fileChanged func probably also loads the file)?
	// it doesn't really matter, the performance here is negligible.
	newCA, err := ioutil.ReadFile(caPath)
	if err != nil {
		return false, errors.New("reading CA '" + caPath + "': " + err.Error())
	}

	// this is the "new old" ca, what was previously named ca.new but is now old.
	oldCA, err := ioutil.ReadFile(E2ESSLNewCAPath)
	if err != nil {
		if os.IsNotExist(err) {
			oldCA = []byte{} // if there was no previous CA, set empty bytes (which will make all the concatenating and moving "just work")
		} else {
			return false, errors.New("reading previous CA '" + E2ESSLNewCAPath + "': " + err.Error())
		}
	}

	newCAKey, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		return false, errors.New("reading CA key '" + caKeyPath + "': " + err.Error())
	}

	// this is the "new old" key, what was previously named cakey.new but is now old
	oldCAKey, err := ioutil.ReadFile(E2ESSLNewCAKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			oldCA = []byte{} // if there was no previous CA key, set empty bytes (which will make all the concatenating and moving "just work")
		} else {
			return false, errors.New("reading previous CA key '" + E2ESSLNewCAKeyPath + "': " + err.Error())
		}
	}

	if err := e2eWriteAtomic(E2ESSLNewCAPath, newCA); err != nil {
		return false, errors.New("writing new CA '" + E2ESSLNewCAPath + "': " + err.Error())
	}
	if err := e2eWriteAtomic(E2ESSLOldCAPath, oldCA); err != nil {
		return false, errors.New("writing old CA '" + E2ESSLOldCAPath + "': " + err.Error())
	}
	if err := e2eWriteAtomic(E2ESSLNewCAKeyPath, newCAKey); err != nil {
		return false, errors.New("writing new CA '" + E2ESSLNewCAKeyPath + "': " + err.Error())
	}
	if err := e2eWriteAtomic(E2ESSLOldCAKeyPath, oldCAKey); err != nil {
		return false, errors.New("writing old CA '" + E2ESSLOldCAKeyPath + "': " + err.Error())
	}

	oldPlusNew := append(newCA, oldCA...)
	if err := e2eWriteAtomic(caDestPath, oldPlusNew); err != nil {
		return false, errors.New("writing concatenated new and old CA '" + caDestPath + "': " + err.Error())
	}

	return true, nil
}

// e2eWriteAtomic writes a file atomically, by writing path.temp and then moving it.
//
// We must always write a temp file, then mv, for atomicity.
// Otherwise, a system crash would result in malformed files
//
// This is designed for E2E SSL work, and always writes 0600.
//
func e2eWriteAtomic(path string, bts []byte) error {
	if err := ioutil.WriteFile(path+".temp", bts, 0600); err != nil {
		return errors.New("writing temp file '" + path + ".temp" + "': " + err.Error())
	}
	if err := os.Rename(path+".temp", path); err != nil {
		return errors.New("moving temp file to '" + path + "': " + err.Error())
	}
	return nil
}
