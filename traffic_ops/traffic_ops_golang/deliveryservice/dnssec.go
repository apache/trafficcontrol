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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"

	"github.com/miekg/dns"
)

func PutDNSSecKeys(tx *sql.Tx, cfg *config.Config, xmlID string, cdnName string, exampleURLs []string) error {
	keys, ok, err := riaksvc.GetDNSSECKeys(cdnName, tx, cfg.RiakAuthOptions)
	if err != nil {
		return errors.New("getting DNSSec keys from Riak: " + err.Error())
	} else if !ok {
		return errors.New("getting DNSSec keys from Riak: no DNSSec keys found")
	}
	cdnKeys, ok := keys[cdnName]
	// TODO warn and continue?
	if !ok {
		return errors.New("getting DNSSec keys from Riak: no DNSSec keys for CDN")
	}
	kExp := getKeyExpiration(cdnKeys.KSK, dnssecDefaultKSKExpiration)
	zExp := getKeyExpiration(cdnKeys.ZSK, dnssecDefaultZSKExpiration)
	overrideTTL := false
	dsKeys, err := CreateDNSSECKeys(tx, cfg, xmlID, exampleURLs, cdnKeys, kExp, zExp, dnssecDefaultTTL, overrideTTL)
	if err != nil {
		return errors.New("creating DNSSEC keys for delivery service '" + xmlID + "': " + err.Error())
	}
	keys[xmlID] = dsKeys
	if err := riaksvc.PutDNSSECKeys(keys, cdnName, tx, cfg.RiakAuthOptions); err != nil {
		return errors.New("putting Riak DNSSEC keys: " + err.Error())
	}
	return nil
}

// CreateDNSSECKeys creates DNSSEC keys for the given delivery service, updating existing keys if they exist. The overrideTTL parameter determines whether to reuse existing key TTLs if they exist, or to override existing TTLs with the ttl parameter's value.
func CreateDNSSECKeys(tx *sql.Tx, cfg *config.Config, xmlID string, exampleURLs []string, cdnKeys tc.DNSSECKeySet, kskExpiration time.Duration, zskExpiration time.Duration, ttl time.Duration, overrideTTL bool) (tc.DNSSECKeySet, error) {
	if len(cdnKeys.ZSK) == 0 {
		return tc.DNSSECKeySet{}, errors.New("getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN")
	}
	if len(cdnKeys.KSK) == 0 {
		return tc.DNSSECKeySet{}, errors.New("getting DNSSec keys from Riak: no DNSSec ZSK keys for CDN")
	}
	if !overrideTTL {
		ttl = getKeyTTL(cdnKeys.KSK, ttl)
	}
	dsName, err := GetDSDomainName(exampleURLs)
	if err != nil {
		return tc.DNSSECKeySet{}, errors.New("creating DS domain name: " + err.Error())
	}
	inception := time.Now()
	zExpiration := inception.Add(zskExpiration)
	kExpiration := inception.Add(kskExpiration)

	tld := false
	effectiveDate := inception
	zsk, err := GetDNSSECKeys(tc.DNSSECZSKType, dsName, ttl, inception, zExpiration, tc.DNSSECKeyStatusNew, effectiveDate, tld)
	if err != nil {
		return tc.DNSSECKeySet{}, errors.New("getting DNSSEC keys for ZSK: " + err.Error())
	}
	ksk, err := GetDNSSECKeys(tc.DNSSECKSKType, dsName, ttl, inception, kExpiration, tc.DNSSECKeyStatusNew, effectiveDate, tld)
	if err != nil {
		return tc.DNSSECKeySet{}, errors.New("getting DNSSEC keys for KSK: " + err.Error())
	}
	return tc.DNSSECKeySet{ZSK: []tc.DNSSECKey{zsk}, KSK: []tc.DNSSECKey{ksk}}, nil
}

func GetDNSSECKeys(keyType string, dsName string, ttl time.Duration, inception time.Time, expiration time.Time, status string, effectiveDate time.Time, tld bool) (tc.DNSSECKey, error) {
	key := tc.DNSSECKey{
		InceptionDateUnix:  inception.Unix(),
		ExpirationDateUnix: expiration.Unix(),
		Name:               dsName,
		TTLSeconds:         uint64(ttl / time.Second),
		Status:             status,
		EffectiveDateUnix:  effectiveDate.Unix(),
	}
	isKSK := keyType != tc.DNSSECZSKType
	err := error(nil)
	key.Public, key.Private, key.DSRecord, err = genKeys(dsName, isKSK, ttl, tld)
	return key, err
}

// genKeys generates keys for DNSSEC for a delivery service. Returns the public key, private key, and DS record (which will be nil if ksk or tld is false).
// This emulates the old Perl Traffic Ops behavior: the public key is of the RFC1035 single-line zone file format, base64 encoded; the private key is of the BIND private-key-file format, base64 encoded; the DSRecord contains the algorithm, digest type, and digest.
func genKeys(dsName string, ksk bool, ttl time.Duration, tld bool) (string, string, *tc.DNSSECKeyDSRecord, error) {
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
			Ttl:    uint32(ttl / time.Second),
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

func GetDSDomainName(dsExampleURLs []string) (string, error) {
	// TODO move somewhere generic
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

const dnssecDefaultKSKExpiration = time.Duration(365) * time.Hour * 24
const dnssecDefaultZSKExpiration = time.Duration(30) * time.Hour * 24
const dnssecDefaultTTL = 60

func getKeyExpiration(keys []tc.DNSSECKey, defaultExpiration time.Duration) time.Duration {
	for _, key := range keys {
		if key.Status != tc.DNSSECKeyStatusNew {
			continue
		}
		return time.Duration(key.ExpirationDateUnix-key.InceptionDateUnix) * time.Second
	}
	return defaultExpiration
}

func getKeyTTL(keys []tc.DNSSECKey, defaultTTL time.Duration) time.Duration {
	for _, key := range keys {
		if key.Status != tc.DNSSECKeyStatusNew {
			continue
		}
		return time.Duration(key.TTLSeconds) * time.Second
	}
	return defaultTTL
}
