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
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

const DefaultKSKTTLSeconds = 60
const DefaultKSKEffectiveMultiplier = 2

func GenerateKSK(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName := tc.CDNName(inf.Params["name"])
	req := tc.CDNGenerateKSKReq{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &req); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}

	ttl, multiplier, err := getKSKParams(inf.Tx.Tx, cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting CDN KSK parameters: "+err.Error()))
		return
	}
	if ttl == nil {
		kskTTL := uint64(DefaultKSKTTLSeconds)
		ttl = &kskTTL
	}
	if multiplier == nil {
		mult := uint64(DefaultKSKEffectiveMultiplier)
		multiplier = &mult
	}

	dnssecKeys, ok, err := riaksvc.GetDNSSECKeys(string(cdnName), inf.Tx.Tx, inf.Config.RiakAuthOptions)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting CDN DNSSEC keys: "+err.Error()))
		return
	}
	if !ok {
		log.Warnln("Generating CDN '" + string(cdnName) + "' KSK: no keys found in Riak, generating and inserting new key anyway")
	}

	isKSK := true
	newKey, err := regenExpiredKeys(isKSK, dnssecKeys[string(cdnName)], *req.EffectiveDate, true, true)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("regenerating CDN DNSSEC keys: "+err.Error()))
		return
	}
	dnssecKeys[string(cdnName)] = newKey

	if err := riaksvc.PutDNSSECKeys(dnssecKeys, string(cdnName), inf.Tx.Tx, inf.Config.RiakAuthOptions); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("putting CDN DNSSEC keys: "+err.Error()))
		return
	}
	api.WriteResp(w, r, "Successfully generated ksk dnssec keys for "+string(cdnName))
}

// getKSKParams returns the CDN's profile's tld.ttls.DNSKEY and DNSKEY.effective.multiplier parameters. If either parameter doesn't exist, nil is returned.
func getKSKParams(tx *sql.Tx, cdn tc.CDNName) (*uint64, *uint64, error) {
	qry := `
WITH cdn_profile_id AS (
  SELECT
    p.id
  FROM
    profile p
    JOIN cdn c ON c.id = p.cdn
  WHERE
    c.name = $1
    AND (p.name like 'CCR%' OR p.name like 'TR%')
  FETCH FIRST 1 ROWS ONLY
)
SELECT
  pa.name,
  pa.value
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN cdn_profile_id pi on pi.id = pr.id
WHERE
  (pa.name = 'tld.ttls.DNSKEY' OR pa.name = 'DNSKEY.effective.multiplier')
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, nil, errors.New("getting cdn ksk parameters: " + err.Error())
	}
	defer rows.Close()
	ttl := (*uint64)(nil)
	mult := (*uint64)(nil)

	for rows.Next() {
		name := ""
		valStr := ""
		if err := rows.Scan(&name, &valStr); err != nil {
			return nil, nil, errors.New("scanning cdn ksk parameters: " + err.Error())
		}
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			log.Warnln("getting CDN KSK parameters: parameter '" + name + "' value '" + valStr + "' is not a number, skipping")
			continue
		}
		if name == "tld.ttls.DNSKEY" {
			ttl = &val
		} else {
			mult = &val
		}
	}
	return ttl, mult, nil
}

// regenExpiredKeys regenerates expired keys. The key is the map key into the keys object, which may be a CDN name or a delivery service name.
func regenExpiredKeys(typeKSK bool, existingKeys tc.DNSSECKeySet, effectiveDate time.Time, tld bool, resetExp bool) (tc.DNSSECKeySet, error) {

	existingKey := ([]tc.DNSSECKey)(nil)
	if typeKSK {
		existingKey = existingKeys.KSK
	} else {
		existingKey = existingKeys.ZSK
	}

	oldKey := tc.DNSSECKey{}
	oldKeyFound := false
	for _, key := range existingKey {
		if key.Status == tc.DNSSECKeyStatusNew {
			oldKey = key
			oldKeyFound = true
			break
		}
	}

	if !oldKeyFound {
		return existingKeys, errors.New("no old key found") // TODO verify this is correct (Perl doesn't check)
	}

	name := oldKey.Name
	ttl := time.Duration(oldKey.TTLSeconds) * time.Second
	expiration := oldKey.ExpirationDateUnix
	inception := oldKey.InceptionDateUnix
	const secPerDay = 86400
	expirationDays := (expiration - inception) / secPerDay

	newInception := time.Now()
	newExpiration := time.Now().Add(time.Duration(expirationDays) * time.Hour * 24)

	keyType := tc.DNSSECKSKType
	if !typeKSK {
		keyType = tc.DNSSECZSKType
	}
	newKey, err := deliveryservice.GetDNSSECKeys(keyType, name, ttl, newInception, newExpiration, tc.DNSSECKeyStatusNew, effectiveDate, tld)
	if err != nil {
		return tc.DNSSECKeySet{}, errors.New("getting and generating DNSSEC keys: " + err.Error())
	}

	oldKey.Status = tc.DNSSECKeyStatusExpired
	if resetExp {
		oldKey.ExpirationDateUnix = effectiveDate.Unix()
	}

	regenKeys := tc.DNSSECKeySet{}
	if typeKSK {
		regenKeys = tc.DNSSECKeySet{ZSK: existingKeys.ZSK, KSK: []tc.DNSSECKey{newKey, oldKey}}
	} else {
		regenKeys = tc.DNSSECKeySet{ZSK: []tc.DNSSECKey{newKey, oldKey}, KSK: existingKeys.KSK}
	}
	return regenKeys, nil
}
