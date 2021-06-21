package cdn

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
	// "context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault"

	"github.com/lib/pq"
)

func RefreshDNSSECKeys(w http.ResponseWriter, r *http.Request) {
	if setInDNSSECKeyRefresh() {
		db, err := api.GetDB(r.Context())
		noTx := (*sql.Tx)(nil) // make a variable instead of passing nil directly, to reduce copy-paste errors
		if err != nil {
			api.HandleErr(w, r, noTx, http.StatusInternalServerError, nil, errors.New("RefresHDNSSECKeys getting db from context: "+err.Error()))
			unsetInDNSSECKeyRefresh()
			return
		}
		cfg, err := api.GetConfig(r.Context())
		if err != nil {
			api.HandleErr(w, r, noTx, http.StatusInternalServerError, nil, errors.New("RefresHDNSSECKeys getting config from context: "+err.Error()))
			unsetInDNSSECKeyRefresh()
			return
		}
		if !cfg.TrafficVaultEnabled {
			api.HandleErr(w, r, noTx, http.StatusInternalServerError, nil, errors.New("refreshing DNSSEC keys: Traffic Vault not enabled"))
			unsetInDNSSECKeyRefresh()
			return
		}
		tv, err := api.GetTrafficVault(r.Context())
		if err != nil {
			api.HandleErr(w, r, noTx, http.StatusInternalServerError, nil, errors.New("RefresHDNSSECKeys getting Traffic Vault from context: "+err.Error()))
			unsetInDNSSECKeyRefresh()
			return
		}

		tx, err := db.Begin()
		if err != nil {
			api.HandleErr(w, r, noTx, http.StatusInternalServerError, nil, errors.New("RefresHDNSSECKeys beginning tx: "+err.Error()))
			unsetInDNSSECKeyRefresh()
			return
		}
		go doDNSSECKeyRefresh(tx, tv) // doDNSSECKeyRefresh takes ownership of tx and MUST close it.
	} else {
		log.Infoln("RefreshDNSSECKeys called, while server was concurrently executing a refresh, doing nothing")
	}

	api.WriteResp(w, r, "Checking DNSSEC keys for refresh in the background")
}

const DNSSECKeyRefreshDefaultTTL = time.Duration(60) * time.Second
const DNSSECKeyRefreshDefaultGenerationMultiplier = uint64(10)
const DNSSECKeyRefreshDefaultEffectiveMultiplier = uint64(10)
const DNSSECKeyRefreshDefaultKSKExpiration = time.Duration(365) * time.Hour * 24
const DNSSECKeyRefreshDefaultZSKExpiration = time.Duration(30) * time.Hour * 24

// doDNSSECKeyRefresh refreshes the CDN's DNSSEC keys, as necessary.
// This takes ownership of tx, and MUST call `tx.Close()`.
// This SHOULD only be called if setInDNSSECKeyRefresh() returned true, in which case this MUST call unsetInDNSSECKeyRefresh() before returning.
func doDNSSECKeyRefresh(tx *sql.Tx, tv trafficvault.TrafficVault) {
	doCommit := true
	defer func() {
		if doCommit {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	defer unsetInDNSSECKeyRefresh()

	updatedAny := false

	cdnDNSSECKeyParams, err := getDNSSECKeyRefreshParams(tx)
	if err != nil {
		log.Errorln("refreshing DNSSEC Keys: getting cdn parameters: " + err.Error())
		doCommit = false
		return
	}
	cdns := []string{}
	for _, inf := range cdnDNSSECKeyParams {
		if inf.DNSSECEnabled {
			cdns = append(cdns, string(inf.CDNName))
		}
	}
	// TODO change to return a slice, map is slow and unnecessary
	dsInfo, err := getDNSSECKeyRefreshDSInfo(tx, cdns)
	if err != nil {
		log.Errorln("refreshing DNSSEC Keys: getting ds info: " + err.Error())
		doCommit = false
		return
	}
	dses := []string{}
	for ds, _ := range dsInfo {
		dses = append(dses, string(ds))
	}

	dsMatchlists, err := deliveryservice.GetDeliveryServicesMatchLists(dses, tx)
	if err != nil {
		log.Errorln("refreshing DNSSEC Keys: getting ds matchlists: " + err.Error())
		doCommit = false
		return
	}
	exampleURLs := map[tc.DeliveryServiceName][]string{}
	for ds, inf := range dsInfo {
		exampleURLs[ds] = deliveryservice.MakeExampleURLs(inf.Protocol, inf.Type, inf.RoutingName, dsMatchlists[string(ds)], inf.CDNDomain)
	}

	for _, cdnInf := range cdnDNSSECKeyParams {
		keys, ok, err := tv.GetDNSSECKeys(string(cdnInf.CDNName), tx, context.Background()) // TODO get all in a map beforehand
		if err != nil {
			log.Warnln("refreshing DNSSEC Keys: getting cdn '" + string(cdnInf.CDNName) + "' keys from Traffic Vault, skipping: " + err.Error())
			continue
		}
		if !ok {
			log.Warnln("refreshing DNSSEC Keys: cdn '" + string(cdnInf.CDNName) + "' has no keys in Traffic Vault, skipping")
			continue
		}

		ttl := DNSSECKeyRefreshDefaultTTL
		if cdnInf.TLDTTLsDNSKEY != nil {
			ttl = time.Duration(*cdnInf.TLDTTLsDNSKEY) * time.Second
		}

		genMultiplier := DNSSECKeyRefreshDefaultGenerationMultiplier
		if cdnInf.DNSKEYGenerationMultiplier != nil {
			genMultiplier = *cdnInf.DNSKEYGenerationMultiplier
		}

		effectiveMultiplier := DNSSECKeyRefreshDefaultEffectiveMultiplier
		if cdnInf.DNSKEYEffectiveMultiplier != nil {
			effectiveMultiplier = *cdnInf.DNSKEYEffectiveMultiplier
		}

		nowPlusTTL := time.Now().Add(ttl * time.Duration(genMultiplier)) // "key_expiration" in the Perl this was transliterated from

		defaultKSKExpiration := DNSSECKeyRefreshDefaultKSKExpiration
		for _, key := range keys[string(cdnInf.CDNName)].KSK {
			if key.Status != tc.DNSSECKeyStatusNew {
				continue
			}
			defaultKSKExpiration = time.Unix(key.ExpirationDateUnix, 0).Sub(time.Unix(key.InceptionDateUnix, 0))
			break
		}

		defaultZSKExpiration := DNSSECKeyRefreshDefaultZSKExpiration
		for _, key := range keys[string(cdnInf.CDNName)].ZSK {
			if key.Status != tc.DNSSECKeyStatusNew {
				continue
			}
			expiration := time.Unix(key.ExpirationDateUnix, 0)
			inception := time.Unix(key.InceptionDateUnix, 0)
			defaultZSKExpiration = expiration.Sub(inception)

			if expiration.After(nowPlusTTL) {
				continue
			}
			log.Infoln("The ZSK keys for '" + string(cdnInf.CDNName) + "' are expired!")
			effectiveDate := expiration.Add(ttl * time.Duration(effectiveMultiplier) * -1) // -1 to subtract
			isKSK := false
			cdnDNSDomain := cdnInf.CDNDomain + "."
			newKeys, err := regenExpiredKeys(isKSK, cdnDNSDomain, keys[string(cdnInf.CDNName)], effectiveDate, false, false)
			if err != nil {
				log.Errorln("refreshing DNSSEC Keys: regenerating expired ZSK keys: " + err.Error())
			} else {
				keys[string(cdnInf.CDNName)] = newKeys
				updatedAny = true
			}
		}

		for _, ds := range dsInfo {
			if ds.CDNName != cdnInf.CDNName {
				continue
			}
			if t := ds.Type; !t.UsesDNSSECKeys() {
				continue
			}

			dsKeys, dsKeysExist := keys[string(ds.DSName)]
			if !dsKeysExist {
				log.Infoln("Keys do not exist for ds '" + string(ds.DSName) + "'")

				cdnKeys, ok := keys[string(ds.CDNName)]
				if !ok {
					log.Errorln("refreshing DNSSEC Keys: cdn has no keys, cannot create ds keys")
					continue
				}

				overrideTTL := false
				dsKeys, err := deliveryservice.CreateDNSSECKeys(exampleURLs[ds.DSName], cdnKeys, defaultKSKExpiration, defaultZSKExpiration, ttl, overrideTTL)
				if err != nil {
					log.Errorln("refreshing DNSSEC Keys: creating missing ds keys: " + err.Error())
				}
				keys[string(ds.DSName)] = dsKeys
				updatedAny = true
				continue
			}

			for _, key := range dsKeys.KSK {
				if key.Status != tc.DNSSECKeyStatusNew {
					continue
				}
				expiration := time.Unix(key.ExpirationDateUnix, 0)
				if expiration.After(nowPlusTTL) {
					continue
				}
				log.Infoln("The KSK keys for '" + ds.DSName + "' are expired!")
				effectiveDate := expiration.Add(ttl * time.Duration(effectiveMultiplier) * -1) // -1 to subtract
				isKSK := true
				newKeys, err := regenExpiredKeys(isKSK, string(ds.DSName), dsKeys, effectiveDate, false, false)
				if err != nil {
					log.Errorln("refreshing DNSSEC Keys: regenerating expired KSK keys for ds '" + string(ds.DSName) + "': " + err.Error())
				} else {
					keys[string(ds.DSName)] = newKeys
					updatedAny = true
				}
			}

			for _, key := range dsKeys.ZSK {
				if key.Status != tc.DNSSECKeyStatusNew {
					continue
				}
				expiration := time.Unix(key.ExpirationDateUnix, 0)
				if expiration.After(nowPlusTTL) {
					continue
				}
				log.Infoln("The ZSK keys for '" + ds.DSName + "' are expired!")
				effectiveDate := expiration.Add(ttl * time.Duration(effectiveMultiplier) * -1) // -1 to subtract
				isKSK := false
				newKeys, err := regenExpiredKeys(isKSK, string(ds.DSName), dsKeys, effectiveDate, false, false)
				if err != nil {
					log.Errorln("refreshing DNSSEC Keys: regenerating expired ZSK keys for ds '" + string(ds.DSName) + "': " + err.Error())
				} else {
					if existingNewKeys, ok := keys[string(ds.DSName)]; ok {
						existingNewKeys.ZSK = newKeys.ZSK
						newKeys = existingNewKeys
					}
					keys[string(ds.DSName)] = newKeys
					updatedAny = true
				}
			}
		}
		if updatedAny {
			if err := tv.PutDNSSECKeys(string(cdnInf.CDNName), keys, tx, context.Background()); err != nil {
				log.Errorln("refreshing DNSSEC Keys: putting keys into Traffic Vault for cdn '" + string(cdnInf.CDNName) + "': " + err.Error())
			}
		}
	}
	log.Infoln("Done refreshing DNSSEC keys")
}

type DNSSECKeyRefreshCDNInfo struct {
	CDNName                    tc.CDNName
	CDNDomain                  string
	DNSSECEnabled              bool
	TLDTTLsDNSKEY              *uint64
	DNSKEYEffectiveMultiplier  *uint64
	DNSKEYGenerationMultiplier *uint64
}

// getDNSSECKeyRefreshParams returns returns the CDN's profile's tld.ttls.DNSKEY, DNSKEY.effective.multiplier, and DNSKEY.generation.multiplier parameters. If either parameter doesn't exist, nil is returned.
// If a CDN exists, but has no parameters, it is returned as a key in the map with a nil value.
func getDNSSECKeyRefreshParams(tx *sql.Tx) (map[tc.CDNName]DNSSECKeyRefreshCDNInfo, error) {
	qry := `
WITH cdn_profile_ids AS (
  SELECT
    DISTINCT(c.name) as cdn_name,
    c.domain_name as cdn_domain,
    c.dnssec_enabled as cdn_dnssec_enabled,
    MAX(p.id) as profile_id -- We only want 1 profile, so get the probably-newest if there's more than one.
  FROM
    cdn c
    LEFT JOIN profile p ON c.id = p.cdn AND (p.name like 'CCR%' OR p.name like 'TR%')
    GROUP BY c.name, c.dnssec_enabled, c.domain_name
)
SELECT
  DISTINCT(pi.cdn_name),
  pi.cdn_domain,
  pi.cdn_dnssec_enabled,
  MAX(pa.name) as parameter_name,
  MAX(pa.value) as parameter_value
FROM
  cdn_profile_ids pi
  LEFT JOIN profile pr ON pi.profile_id = pr.id
  LEFT JOIN profile_parameter pp ON pr.id = pp.profile
  LEFT JOIN parameter pa ON pp.parameter = pa.id AND (
    pa.name = 'tld.ttls.DNSKEY'
    OR pa.name = 'DNSKEY.effective.multiplier'
    OR pa.name = 'DNSKEY.generation.multiplier'
  )
GROUP BY pi.cdn_name, pi.cdn_domain, pi.cdn_dnssec_enabled
`
	rows, err := tx.Query(qry)
	if err != nil {
		return nil, errors.New("getting cdn dnssec key refresh parameters: " + err.Error())
	}
	defer rows.Close()

	params := map[tc.CDNName]DNSSECKeyRefreshCDNInfo{}
	for rows.Next() {
		cdnName := tc.CDNName("")
		cdnDomain := ""
		dnssecEnabled := false
		name := util.StrPtr("")
		valStr := util.StrPtr("")
		if err := rows.Scan(&cdnName, &cdnDomain, &dnssecEnabled, &name, &valStr); err != nil {
			return nil, errors.New("scanning cdn dnssec key refresh parameters: " + err.Error())
		}

		inf := params[cdnName]
		inf.CDNName = cdnName
		inf.CDNDomain = cdnDomain
		inf.DNSSECEnabled = dnssecEnabled

		if name == nil || valStr == nil {
			// no DNSKEY parameters, but the CDN still exists.
			params[cdnName] = inf
			continue
		}

		val, err := strconv.ParseUint(*valStr, 10, 64)
		if err != nil {
			log.Warnln("getting CDN dnssec refresh parameters: parameter '" + *name + "' value '" + *valStr + "' is not a number, skipping")
			params[cdnName] = inf
			continue
		}

		switch *name {
		case "tld.ttls.DNSKEY":
			inf.TLDTTLsDNSKEY = &val
		case "DNSKEY.effective.multiplier":
			inf.DNSKEYEffectiveMultiplier = &val
		case "DNSKEY.generation.multiplier":
			inf.DNSKEYGenerationMultiplier = &val
		default:
			log.Warnln("getDNSSECKeyRefreshParams got unknown parameter '" + *name + "', skipping")
			continue
		}
		params[cdnName] = inf
	}
	return params, nil
}

type DNSSECKeyRefreshDSInfo struct {
	DSName      tc.DeliveryServiceName
	Type        tc.DSType
	Protocol    *int
	CDNName     tc.CDNName
	CDNDomain   string
	RoutingName string
}

func getDNSSECKeyRefreshDSInfo(tx *sql.Tx, cdns []string) (map[tc.DeliveryServiceName]DNSSECKeyRefreshDSInfo, error) {
	qry := `
SELECT
  ds.xml_id,
  tp.name as type,
  ds.protocol,
  c.name as cdn_name,
  c.domain_name as cdn_domain,
  ds.routing_name
FROM
  deliveryservice ds
  JOIN type tp ON tp.id = ds.type
  JOIN cdn c ON c.id = ds.cdn_id
WHERE
  c.name = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(cdns))
	if err != nil {
		return nil, errors.New("getting cdn dnssec key refresh ds info: " + err.Error())
	}
	defer rows.Close()

	dsInf := map[tc.DeliveryServiceName]DNSSECKeyRefreshDSInfo{}
	for rows.Next() {
		i := DNSSECKeyRefreshDSInfo{}
		if err := rows.Scan(&i.DSName, &i.Type, &i.Protocol, &i.CDNName, &i.CDNDomain, &i.RoutingName); err != nil {
			return nil, errors.New("scanning cdn dnssec key refresh ds info: " + err.Error())
		}
		dsInf[i.DSName] = i
	}
	return dsInf, nil
}

// inDNSSECKeyRefresh is whether the server is currently processing a refresh in the background.
// This is used to only perform 1 refresh at a time.
// This MUST NOT be changed outside of atomic operations.
// This MUST NOT be changed to a boolean, or set without atomics. Atomic semantics involve more than just setting a memory location.
var inDNSSECKeyRefresh = uint64(0)

// setInDNSSECKeyRefresh attempts to set whether the server is currently executing a DNSSEC key refresh operation.
// Returns false if a refresh operation is already executing.
// If this returns true, the caller MUST call unsetInDNSSECKeyRefresh().
func setInDNSSECKeyRefresh() bool { return atomic.CompareAndSwapUint64(&inDNSSECKeyRefresh, 0, 1) }

// unsetInDNSSECKeyRefresh sets the flag indicating that the server is currently executing a DNSSEC key refresh operation to false.
// This MUST NOT be called, unless setInDNSSECKeyRefresh() was previously called and returned true.
func unsetInDNSSECKeyRefresh() { atomic.StoreUint64(&inDNSSECKeyRefresh, 0) }
