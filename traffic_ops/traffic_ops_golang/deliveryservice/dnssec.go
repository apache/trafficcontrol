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
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"

	"github.com/miekg/dns"
)

func createDNSSecKeys(tx *sql.Tx, cfg config.Config, dsID int, xmlID string, cdnName string, cdnDomain string, dnssecEnabled bool, exampleURLs []string) error {
	if !dnssecEnabled {
		return nil
	}

	keys, ok, err := riaksvc.GetDNSSECKeys(cdnName, tx, cfg.RiakAuthOptions)
	if err != nil {
		log.Errorln("Getting DNSSec keys from Riak: " + err.Error())
		return errors.New("getting DNSSec keys from Riak: " + err.Error())
	}
	if !ok {
		log.Errorln("Getting DNSSec keys from Riak: no DNSSec keys found")
		return errors.New("getting DNSSec keys from Riak: no DNSSec keys found")
	}

	cdnKeys, ok := keys[cdnName]
	// TODO warn and continue?
	if !ok {
		log.Errorln("Getting DNSSec keys from Riak: no DNSSec keys for CDN '" + cdnName + "'")
		return errors.New("getting DNSSec keys from Riak: no DNSSec keys for CDN")
	}
	if len(cdnKeys.ZSK) == 0 {
		log.Errorln("Getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN '" + cdnName + "'")
		return errors.New("getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN")
	}
	if len(cdnKeys.KSK) == 0 {
		log.Errorln("Getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN '" + cdnName + "'")
		return errors.New("getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN")
	}

	kExpDays := getKeyExpirationDays(cdnKeys.KSK, dnssecDefaultKSKExpirationDays)
	zExpDays := getKeyExpirationDays(cdnKeys.ZSK, dnssecDefaultZSKExpirationDays)
	ttl := getKeyTTL(cdnKeys.KSK, dnssecDefaultTTL)
	dsName, err := getDSDomainName(exampleURLs)
	if err != nil {
		log.Errorln("creating DS domain name: " + err.Error())
		return errors.New("creating DS domain name: " + err.Error())
	}
	inception := time.Now()
	zExpiration := inception.Add(time.Duration(zExpDays) * time.Hour * 24)
	kExpiration := inception.Add(time.Duration(kExpDays) * time.Hour * 24)

	tld := false
	effectiveDate := inception
	zsk, err := getDNSSECKeys(dnssecZSKType, dsName, ttl, inception, zExpiration, dnssecKeyStatusNew, effectiveDate, tld)
	if err != nil {
		log.Errorln("getting DNSSEC keys for ZSK: " + err.Error())
		return errors.New("getting DNSSEC keys for ZSK: " + err.Error())
	}
	ksk, err := getDNSSECKeys(dnssecKSKType, dsName, ttl, inception, kExpiration, dnssecKeyStatusNew, effectiveDate, tld)
	if err != nil {
		log.Errorln("getting DNSSEC keys for KSK: " + err.Error())
		return errors.New("getting DNSSEC keys for KSK: " + err.Error())
	}
	keys[xmlID] = tc.DNSSECKeySet{ZSK: []tc.DNSSECKey{zsk}, KSK: []tc.DNSSECKey{ksk}}

	if err := riaksvc.PutDNSSECKeys(keys, cdnName, tx, cfg.RiakAuthOptions); err != nil {
		log.Errorln("putting Riak DNSSEC keys: " + err.Error())
		return errors.New("putting Riak DNSSEC keys: " + err.Error())
	}
	return nil
}

func getDNSSECKeys(keyType string, dsName string, ttl uint64, inception time.Time, expiration time.Time, status string, effectiveDate time.Time, tld bool) (tc.DNSSECKey, error) {
	key := tc.DNSSECKey{
		InceptionDateUnix:  inception.Unix(),
		ExpirationDateUnix: expiration.Unix(),
		Name:               dsName + ".",
		TTLSeconds:         ttl,
		Status:             status,
		EffectiveDateUnix:  effectiveDate.Unix(),
	}
	isKSK := keyType != dnssecZSKType
	err := error(nil)
	key.Public, key.Private, key.DSRecord, err = genKeys(dsName, isKSK, ttl, tld)
	return key, err
}

// genKeys generates keys for DNSSEC for a delivery service. Returns the public key, private key, and DS record (which will be nil if ksk or tld is false).
// This emulates the old Perl Traffic Ops behavior: the public key is of the RFC1035 single-line zone file format, base64 encoded; the private key is of the BIND private-key-file format, base64 encoded; the DSRecord contains the algorithm, digest type, and digest.
func genKeys(dsName string, ksk bool, ttl uint64, tld bool) (string, string, *tc.DNSSECKeyDSRecord, error) {
	bits := 1024
	flags := 256
	algorithm := dns.RSASHA1 // 5 - http://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
	protocol := 3

	if ksk {
		flags |= 1
		bits *= 2
	}

	dnskey := dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name:   dsName,
			Rrtype: dns.TypeDNSKEY,
			Class:  dns.ClassINET,
			Ttl:    uint32(ttl),
		},
		Flags:     uint16(flags),
		Protocol:  uint8(protocol),
		Algorithm: algorithm,
	}

	priKey, err := dnskey.Generate(bits)
	if err != nil {
		return "", "", nil, errors.New("error generating DNS key: " + err.Error())
	}

	priKeyStr := dnskey.PrivateKeyString(priKey) // BIND9 private-key-file format; cooresponds to Perl Net::DNS::SEC::Private->generate_rsa.dump_rsa_priv
	priKeyStrBase64 := base64.StdEncoding.EncodeToString([]byte(priKeyStr))

	pubKeyStr := dnskey.String() // RFC1035 single-line zone file format; cooresponds to Perl Net::DNS::RR.plain
	pubKeyStrBase64 := base64.StdEncoding.EncodeToString([]byte(pubKeyStr))

	keyDS := (*tc.DNSSECKeyDSRecord)(nil)
	if ksk && tld {
		dsRecord := dnskey.ToDS(dns.SHA1) // TODO update to SHA512
		keyDS = &tc.DNSSECKeyDSRecord{Algorithm: int64(dsRecord.Algorithm), DigestType: int64(dsRecord.DigestType), Digest: dsRecord.Digest}
	}

	return pubKeyStrBase64, priKeyStrBase64, keyDS, nil
}

// TODO change ttl to time.Duration

const dnssecKSKType = "ksk"
const dnssecZSKType = "zsk"

func getDSDomainName(dsExampleURLs []string) (string, error) {
	if len(dsExampleURLs) == 0 {
		return "", errors.New("no example URLs")
	}

	dsName := dsExampleURLs[0] + "."
	firstDot := strings.Index(dsName, ".")
	if firstDot == -1 {
		return "", errors.New("malformed example URL, no dots")
	}
	if len(dsName) < firstDot+2 {
		return "", errors.New("malformed example URL, nothing after first dot")
	}
	dsName = dsName[firstDot+1:]
	return dsName, nil
}

const dnssecKeyStatusNew = "new"
const secondsPerDay = 86400
const dnssecDefaultKSKExpirationDays = 365
const dnssecDefaultZSKExpirationDays = 30
const dnssecDefaultTTL = 60

func getKeyExpirationDays(keys []tc.DNSSECKey, defaultExpirationDays uint64) uint64 {
	for _, key := range keys {
		if key.Status != dnssecKeyStatusNew {
			continue
		}
		return uint64((key.ExpirationDateUnix - key.InceptionDateUnix) / secondsPerDay)
	}
	return defaultExpirationDays
}

func getKeyTTL(keys []tc.DNSSECKey, defaultTTL uint64) uint64 {
	for _, key := range keys {
		if key.Status != dnssecKeyStatusNew {
			continue
		}
		return key.TTLSeconds
	}
	return defaultTTL
}
