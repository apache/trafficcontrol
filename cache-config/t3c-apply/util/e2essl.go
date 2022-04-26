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
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
)

const E2ESSLCADestPathCert = "e2e-ssl-ca.cert"

const E2ESSLPathBase = "e2e-ssl"

const E2ESSLPathClientBase = E2ESSLPathBase + "-client"

const E2ESSLPathServerBase = E2ESSLPathBase + "-server"

// e2eSSLCertLifetime is the lifetime to generate for certificates. TODO make configurable.
const e2eSSLCertLifetime = time.Hour * 24 * 7

// e2eSSLCertRefreshAge is the age to refresh certificates after. TODO make configurable.
const e2eSSLCertRefreshAge = time.Hour * 24

// E2ESSLClientCertExists returns nil if the E2E SSL client cert and key exists, or an error if they don't exist or we were unable to determine.
func E2ESSLClientCertExists(certDir string) error {
	// TODO add verifying the files are actually keys?
	// That would probably be easiest to do by externally calling openssl (rather than Go parsing).
	files := []string{
		filepath.Join(certDir, E2ESSLPathBase+".key"),
		// filepath.Join(certDir, E2ESSLPathClientBase+".csr"),
		filepath.Join(certDir, E2ESSLPathClientBase+".cert"),
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

// E2ESSLGenerateClientCert generates the client cert and key.
// Note this key is also used for E2ESSL server certs.
func E2ESSLGenerateClientCert(certDir string, caKeyPath string, caCertPath string) error {
	clientKey := filepath.Join(certDir, E2ESSLPathBase+".key")
	clientCSR := filepath.Join(certDir, E2ESSLPathClientBase+".csr")
	clientCert := filepath.Join(certDir, E2ESSLPathClientBase+".cert")

	// client cert private key
	if stdOut, stdErr, code := t3cutil.Do(`sh`, `-c`, `(umask 060; openssl ecparam -name secp256r1 -genkey -noout -out `+clientKey+`)`); code != 0 {
		return fmt.Errorf("generating client private key returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	if err := E2ESetFileOwnerAndMode(clientKey); err != nil {
		return fmt.Errorf("setting client private key '%v' owner and mode: %v", clientKey, err)
	}

	// TODO only get once for the app. Is it needed anywhere else?

	hostnameFQDN, err := GetHostnameFQDN()
	if err != nil {
		return errors.New("getting hostname: " + err.Error())
	}

	certCN := hostnameFQDN
	certAltNames := []string{}

	if err := E2ESSLGenerateCert(caCertPath, caKeyPath, clientKey, clientCSR, clientCert, certCN, certAltNames); err != nil {
		return errors.New("generating client cert: " + err.Error())
	}
	return nil
}

func E2ESSLGenerateCert(caCertPath string, caKeyPath string, clientKeyPath string, csrPath string, certPath string, certCN string, altNames []string) error {
	certLifetimeDaysStr := strconv.Itoa(int(e2eSSLCertLifetime / time.Hour / 24))

	log.Infof("E2ESSLGenerateCert calling -days %v\n", certLifetimeDaysStr)

	csrConfPath := csrPath + ".conf"
	if err := E2ESSLWriteCSRConf(csrConfPath, certCN, altNames); err != nil {
		return fmt.Errorf("writing csr conf '%v': %v", csrConfPath, err)
	}

	if stdOut, stdErr, code := t3cutil.Do("openssl", "req", "-new", "-sha256", "-key", clientKeyPath, "-config", csrConfPath, "-out", csrPath); code != 0 {
		return fmt.Errorf("generating csr returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}
	if stdOut, stdErr, code := t3cutil.Do("openssl", "x509", "-req", "-in", csrPath, "-CA", caCertPath, "-CAkey", caKeyPath, "-CAcreateserial", "-out", certPath, "-days", certLifetimeDaysStr, "-sha256", "-extensions", "v3_req", "-extfile", csrConfPath); code != 0 {
		return fmt.Errorf("generating certificate returned code '%v' stdout '%v' stderr '%v'", code, string(stdOut), string(stdErr))
	}

	return E2ESetFileOwnerAndMode(certPath)
}

func E2ESSLWriteCSRConf(csrConfPath string, certCN string, altNames []string) error {
	if len(altNames) == 0 {
		altNames = append(altNames, certCN)
	}
	altNamesEntries := []string{}
	for i, name := range altNames {
		altNamesEntries = append(altNamesEntries, "DNS."+strconv.Itoa(i)+" = "+name)
	}
	altNamesTxt := strings.Join(altNamesEntries, "\n")
	conf := `
[req]
distinguished_name = dn
req_extensions = v3_req
prompt = no

[ dn ]
C=US
ST=CO
O=ApacheTrafficControl
CN = ` + certCN + `

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = critical,nonRepudiation,digitalSignature,keyEncipherment
extendedKeyUsage=critical,serverAuth,clientAuth
subjectAltName = @alt_names

[ alt_names ]
` + altNamesTxt + `
`
	return e2eWriteAtomic(csrConfPath, []byte(conf))
}

// TODO roll client key
// TODO check CSR for expiration and warn
// TODO log generation/refreshes
// TODO add sha512sum, openssl to RPM dependencies
// TODO refactor to not log. This should act like a library

// E2ESSLRefreshClientCert generates new certificates if they're about to expire.
// The mustRefresh parameter forces new generation, regardless of expiration. This is most often used when the Certificate Authority has changed.
func E2ESSLRefreshClientCert(certDir string, caKeyPath string, caCertPath string, mustRefresh bool) error {
	clientKey := filepath.Join(certDir, E2ESSLPathBase+".key")
	clientCSR := filepath.Join(certDir, E2ESSLPathClientBase+".csr")
	clientCert := filepath.Join(certDir, E2ESSLPathClientBase+".cert")

	hostnameFQDN, err := GetHostnameFQDN()
	if err != nil {
		return errors.New("getting hostname: " + err.Error())
	}

	certCN := hostnameFQDN
	certAltNames := []string{}

	if err := E2ESSLRefreshCert(caKeyPath, caCertPath, mustRefresh, clientKey, clientCert, clientCSR, certCN, certAltNames); err != nil {
		return errors.New("refreshing client cert: " + err.Error())
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

func E2ESSLRefreshCert(caKeyPath string, caCertPath string, mustRefresh bool, keyPath string, certPath string, csrPath string, certCN string, certAltNames []string) error {
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
		if err := E2ESSLGenerateCert(caCertPath, caKeyPath, keyPath, csrPath, certPath, certCN, certAltNames); err != nil {
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
		return "", fmt.Errorf("Certificate '" + certPath + "' returned empty CN!")
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

// E2ESSLGenerateOrRefreshClientCert creates the client cert and key if they don't exist,
// or refreshes them if they do.
// Note this key is also used for E2E server certs.
func E2ESSLGenerateOrRefreshClientCert(certDir string, caKeyPath string, caCertPath string, e2eCACertDestPath string) error {
	caIsNew, err := E2ECheckAndHandleNewCA(caCertPath, caKeyPath, e2eCACertDestPath)
	if err != nil {
		return errors.New("checking and refreshing new CA: " + err.Error())
	}
	mustRefresh := caIsNew // if the CA is new, we must regenerate certificates

	if err := E2ESSLClientCertExists(certDir); err == nil {
		return E2ESSLRefreshClientCert(certDir, caKeyPath, caCertPath, mustRefresh)
	} else if !os.IsNotExist(err) {
		return errors.New("checking if keys exist: " + err.Error())
	}
	if err := E2ESSLGenerateClientCert(certDir, caKeyPath, caCertPath); err != nil {
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

	if !changed {
		// even if the ca in etc/trafficcontrol-cache-config is unchanged,
		// the ca in etc/trafficserver may have been modified or deleted somehow.
		// Verify it exists, and if not, place it

		_, err := os.Stat(caDestPath)
		if err != nil && !os.IsNotExist(err) {
			return false, fmt.Errorf("checking if e2e ca cert '%v' exists: %v", caDestPath, err)
		}
		if err == nil {
			log.Infoln("E2ECheckAndHandleNewCA: ca and key are unchanged")
			return false, nil
		}

		log.Infoln("E2ECheckAndHandleNewCA: ca and key are unchanged in t3c lib directory, but changed in ats etc, writing CA to ats etc and will recreate certs")

		oldPlusNew := append(newCA, oldCA...)
		if err := e2eWriteAtomic(caDestPath, oldPlusNew); err != nil {
			return false, errors.New("writing concatenated new and old CA '" + caDestPath + "': " + err.Error())
		}

		return true, nil
	}

	log.Infoln("E2ECheckAndHandleNewCA: ca or key changed, rewriting ca")

	// the CA changed, so we need to
	// 1. copy the ca.new to ca.old
	// 2. copy the ca to ca.new
	// 3. copy the cakey.new to cakey.old
	// 4. copy the cakey to cakey.new
	// 5. cat ca.new ca.old > etc/trafficserver/ca.cert

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
// After writing, the file's owner is changed to ats.
//
func e2eWriteAtomic(path string, bts []byte) error {
	if err := ioutil.WriteFile(path+".temp", bts, 0600); err != nil {
		return errors.New("writing temp file '" + path + ".temp" + "': " + err.Error())
	}
	if err := os.Rename(path+".temp", path); err != nil {
		return errors.New("moving temp file to '" + path + "': " + err.Error())
	}

	return E2ESetFileOwnerAndMode(path)
}

// E2ESSLGenerateServerCerts generates the End-to-End SSL Server certs, used for
// internal remap targets to parent caches.
//
// If the certs exist, they are refreshed if necessary.
//
// Note the key previously generated for client certs is also used for server certs.
// Note internal sources use a single shared client cert, generated before config and passed to t3c-generate/atscfg.
//
// This func called after config generation. t3c-generate/atscfg return metadata about remaps and required server cert paths, which this func will now generate.
//
// remapData is the metadata received from ssl_multicert config gen, about remap sources and targets and their required cert paths.
//
func E2ESSLGenerateServerCerts(remapData []atscfg.E2ECertMetaData, certDir string, caKeyPath string, caCertPath string, e2eCACertDestPath string) error {
	remapData = FilterE2EMetaData(remapData)
	log.Infof("E2ESSL generating %v server certs\n", len(remapData))

	// all certs use a shared key, no reason to make one for each.
	keyPath := filepath.Join(certDir, E2ESSLPathBase+".key")

	for _, rd := range remapData {
		certPath := filepath.Join(certDir, rd.CertPath)
		csrPath := strings.TrimSuffix(certPath, filepath.Ext(certPath)) + ".csr"
		_, err := os.Stat(certPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("checking if e2e server cert '%v' exists: %v", certPath, err)
		}
		certExists := err == nil

		cnDSName := string(rd.DSName)
		cnSuffix := ".internal.cdn.comcast.invalid"
		if len(cnDSName) > (rfc.MaxCertificateCNLen - len(cnSuffix)) {
			cnDSName = cnDSName[:(rfc.MaxCertificateCNLen-len(cnSuffix)-3)] + "ETC"
		}
		certCN := cnDSName + cnSuffix
		altNames := []string{rd.URI.Hostname()}

		// if the cert exists, check if it needs refreshed
		if certExists {
			mustRefresh := false // don't force refresh if it isn't necessary
			if err := E2ESSLRefreshCert(caKeyPath, caCertPath, mustRefresh, keyPath, certPath, csrPath, certCN, altNames); err != nil {
				return fmt.Errorf("creating e2e ds '%v' %v server cert '%v': %v", rd.DSName, rd.Type, certPath, err)
			} else {
				log.Infof("E2ESSL refreshed ds '" + string(rd.DSName) + "' " + string(rd.Type) + " server cert " + certPath)
			}
			continue
		}

		if err := E2ESSLGenerateCert(caCertPath, caKeyPath, keyPath, csrPath, certPath, certCN, altNames); err != nil {
			return fmt.Errorf("creating e2e ds '%v' %v server cert '%v': %v", rd.DSName, rd.Type, certPath, err)
		} else {
			log.Infof("E2ESSL generated ds '" + string(rd.DSName) + "' " + string(rd.Type) + " server cert " + certPath)
		}
	}
	return nil
}

// FilterE2EMetaData filters the metadata from atscfg and returns only
// the DS metadata which need E2E server certs generated.
func FilterE2EMetaData(mds []atscfg.E2ECertMetaData) []atscfg.E2ECertMetaData {
	filtered := []atscfg.E2ECertMetaData{}
	for _, md := range mds {
		if !md.Internal {
			continue // non-internal sources (clients) and targets (origins) don't need E2E certs.
		}
		if md.Type != atscfg.RemapMapTypeSource {
			continue // only sources need E2E server certs; targets use a single shared client cert
		}
		if md.URI.Scheme != rfc.SchemeHTTPS {
			continue // http remaps don't need certs
		}
		filtered = append(filtered, md)
	}
	return filtered
}

func E2ESetFileOwnerAndMode(filePath string) error {
	atsUser, err := user.Lookup(config.TrafficServerOwner)
	if err != nil {
		// fatal: ATS can't load the file if it doesn't own it
		return fmt.Errorf("could not lookup the trafficserver, '%s', owner uid: %v", config.TrafficServerOwner, err)
	}
	atsUid, err := strconv.Atoi(atsUser.Uid)
	if err != nil {
		// fatal: ATS can't load the file if it doesn't own it
		return fmt.Errorf("got non-integer uid '%v': %v", atsUser.Uid, err)
	}
	atsGid, err := strconv.Atoi(atsUser.Gid)
	if err != nil {
		// not fatal: ATS will still work if the user is set but not the group
		log.Errorln("getting ATS user gid: got non-integer '%v', not setting: %v", atsUser.Gid, err)
		atsGid = -1
	}

	if err := os.Chown(filePath, atsUid, atsGid); err != nil {
		// fatal: ATS can't load the file if it doesn't own it
		return fmt.Errorf("chown ats user for '%v': %v", filePath, err)
	}

	mode := os.FileMode(0600)
	if err := os.Chmod(filePath, mode); err != nil {
		// not fatal: chmod failure isn't good, but ATS will still work
		log.Errorf("chmod 0600 for '%v': %v", filePath, err)
	}
	return nil
}
