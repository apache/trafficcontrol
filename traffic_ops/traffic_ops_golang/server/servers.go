// Package server provides tools for manipulating the server database table and
// corresponding http handlers.
package server

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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/topology/topology_validation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const serversFromAndJoin = `
FROM server AS s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id
`

const joinProfileV4 = `JOIN server_profile sp ON p.name = sp.profile_name AND s.id = sp.server
`

/* language=SQL */
const dssTopologiesJoinSubquery = `
(SELECT
	ARRAY_AGG(CAST(ROW(td.id, s.id, NULL) AS deliveryservice_server))
FROM "server" s
JOIN cachegroup c on s.cachegroup = c.id
JOIN topology_cachegroup tc ON c.name = tc.cachegroup
JOIN deliveryservice td ON td.topology = tc.topology
JOIN type t ON s.type = t.id
LEFT JOIN deliveryservice_server dss
	ON s.id = dss."server"
	AND dss.deliveryservice = td.id
WHERE td.id = :dsId
AND (
	t.name != '` + tc.OriginTypeName + `'
	OR dss.deliveryservice IS NOT NULL
)),
`

/* language=SQL */
const deliveryServiceServersJoin = `
FULL OUTER JOIN (
SELECT (dss.dss_record).deliveryservice, (dss.dss_record).server FROM (
	SELECT UNNEST(COALESCE(
		%s
		(SELECT
			ARRAY_AGG(CAST(ROW(dss.deliveryservice, dss."server", NULL) AS deliveryservice_server))
		FROM deliveryservice_server dss)
	)) AS dss_record) AS dss
) dss ON dss.server = s.id
JOIN deliveryservice d ON cdn.id = d.cdn_id AND dss.deliveryservice = d.id
`

/* language=SQL */
const requiredCapabilitiesCondition = `
AND (
	SELECT ARRAY_AGG(ssc.server_capability)
	FROM server_server_capability ssc
	WHERE ssc."server" = s.id
) @> (
	SELECT d.required_capabilities
	FROM deliveryservice d
	WHERE d.id = :dsId
)
`

const serverCountQuery = `
SELECT COUNT(s.id)
` + serversFromAndJoin

const selectQuery = `
SELECT
	cg.name AS cachegroup,
	s.cachegroup AS cachegroup_id,
	s.cdn_id,
	cdn.name AS cdn_name,
	s.domain_name,
	s.guid,
	s.host_name,
	s.https_port,
	s.id,
	s.ilo_ip_address,
	s.ilo_ip_gateway,
	s.ilo_ip_netmask,
	s.ilo_password,
	s.ilo_username,
	s.last_updated,
	s.mgmt_ip_address,
	s.mgmt_ip_gateway,
	s.mgmt_ip_netmask,
	s.offline_reason,
	pl.name AS phys_location,
	s.phys_location AS phys_location_id,
	(SELECT ARRAY_AGG(sp.profile_name ORDER BY sp.priority ASC) FROM server_profile AS sp where sp.server=s.id) AS profile_name,
	s.rack,
	s.revalidate_update_time > s.revalidate_apply_time AS reval_pending,
	s.revalidate_update_time,
	s.revalidate_apply_time,
	st.name AS status,
	s.status AS status_id,
	s.tcp_port,
	t.name AS server_type,
	s.type AS server_type_id,
	s.config_update_time > s.config_apply_time AS upd_pending,
	s.config_update_time,
	s.config_apply_time,
	s.xmpp_id,
	s.xmpp_passwd,
	s.status_last_updated
` + serversFromAndJoin

const selectIDQuery = `
SELECT
	s.id
` + serversFromAndJoin

const midWhereClause = `
WHERE t.name = :cache_type_mid AND s.cachegroup IN (
	SELECT cg.parent_cachegroup_id FROM cachegroup AS cg
	WHERE cg.id IN (
	SELECT s.cachegroup FROM server AS s
	WHERE s.id = ANY(:edge_ids)))
	AND (SELECT d.topology
		FROM deliveryservice d
		WHERE d.id = :ds_id) IS NULL
`

const insertQueryV3 = `
INSERT INTO server (
	cachegroup,
	cdn_id,
	domain_name,
	host_name,
	https_port,
	ilo_ip_address,
	ilo_ip_netmask,
	ilo_ip_gateway,
	ilo_username,
	ilo_password,
	mgmt_ip_address,
	mgmt_ip_netmask,
	mgmt_ip_gateway,
	offline_reason,
	phys_location,
	profile,
	rack,
	status,
	tcp_port,
	type,
	xmpp_id,
	xmpp_passwd,
	status_last_updated
) VALUES (
	:cachegroup_id,
	:cdn_id,
	:domain_name,
	:host_name,
	:https_port,
	:ilo_ip_address,
	:ilo_ip_netmask,
	:ilo_ip_gateway,
	:ilo_username,
	:ilo_password,
	:mgmt_ip_address,
	:mgmt_ip_netmask,
	:mgmt_ip_gateway,
	:offline_reason,
	:phys_location_id,
	:profile_id,
	:rack,
	:status_id,
	:tcp_port,
	:server_type_id,
	:xmpp_id,
	:xmpp_passwd,
	:status_last_updated
) RETURNING
	id
`

const insertQueryV4 = `
INSERT INTO server (
	cachegroup,
	cdn_id,
	domain_name,
	host_name,
	https_port,
	ilo_ip_address,
	ilo_ip_netmask,
	ilo_ip_gateway,
	ilo_username,
	ilo_password,
	mgmt_ip_address,
	mgmt_ip_netmask,
	mgmt_ip_gateway,
	offline_reason,
	phys_location,
	profile,
	rack,
	status,
	tcp_port,
	type,
	xmpp_id,
	xmpp_passwd
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
) RETURNING
	(SELECT name FROM cachegroup WHERE cachegroup.id=server.cachegroup) AS cachegroup,
	cachegroup AS cachegroup_id,
	cdn_id,
	(SELECT name FROM cdn WHERE cdn.id=server.cdn_id) AS cdn_name,
	domain_name,
	guid,
	host_name,
	https_port,
	id,
	ilo_ip_address,
	ilo_ip_gateway,
	ilo_ip_netmask,
	ilo_password,
	ilo_username,
	last_updated,
	mgmt_ip_address,
	mgmt_ip_gateway,
	mgmt_ip_netmask,
	offline_reason,
	(SELECT name FROM phys_location WHERE phys_location.id=server.phys_location) AS phys_location,
	phys_location AS phys_location_id,
	(SELECT ARRAY[name] FROM profile WHERE profile.id=server.profile) AS profile_name,
	rack,
	(SELECT name FROM status WHERE status.id=server.status) AS status,
	status AS status_id,
	tcp_port,
	(SELECT name FROM type WHERE type.id=server.type) AS server_type,
	type AS server_type_id
`

const updateQuery = `
UPDATE server SET
	cachegroup=:cachegroup_id,
	cdn_id=:cdn_id,
	domain_name=:domain_name,
	host_name=:host_name,
	https_port=:https_port,
	ilo_ip_address=:ilo_ip_address,
	ilo_ip_netmask=:ilo_ip_netmask,
	ilo_ip_gateway=:ilo_ip_gateway,
	ilo_username=:ilo_username,
	ilo_password=:ilo_password,
	mgmt_ip_address=:mgmt_ip_address,
	mgmt_ip_netmask=:mgmt_ip_netmask,
	mgmt_ip_gateway=:mgmt_ip_gateway,
	offline_reason=:offline_reason,
	phys_location=:phys_location_id,
	profile=(SELECT id from profile where name=(SELECT profile_name from server_profile sp WHERE sp.server=:id and priority=0)),
	rack=:rack,
	status=:status_id,
	tcp_port=:tcp_port,
	type=:server_type_id,
	xmpp_passwd=:xmpp_passwd,
	status_last_updated=:status_last_updated
WHERE id=:id
RETURNING
	(SELECT name FROM cachegroup WHERE cachegroup.id=server.cachegroup) AS cachegroup,
	cachegroup AS cachegroup_id,
	cdn_id,
	(SELECT name FROM cdn WHERE cdn.id=server.cdn_id) AS cdn_name,
	domain_name,
	guid,
	host_name,
	https_port,
	id,
	ilo_ip_address,
	ilo_ip_gateway,
	ilo_ip_netmask,
	ilo_password,
	ilo_username,
	last_updated,
	mgmt_ip_address,
	mgmt_ip_gateway,
	mgmt_ip_netmask,
	offline_reason,
	(SELECT name FROM phys_location WHERE phys_location.id=server.phys_location) AS phys_location,
	phys_location AS phys_location_id,
	(SELECT ARRAY_AGG(profile_name ORDER BY priority ASC) FROM server_profile WHERE server_profile.server=server.id) AS profile_name,
	rack,
	(SELECT name FROM status WHERE status.id=server.status) AS status,
	status AS status_id,
	tcp_port,
	(SELECT name FROM type WHERE type.id=server.type) AS server_type,
	type AS server_type_id,
	status_last_updated
`

const originServerQuery = `
JOIN deliveryservice_server dsorg
ON dsorg.server = s.id
WHERE t.name = '` + tc.OriginTypeName + `'
AND dsorg.deliveryservice=:dsId
`
const deleteServerQuery = `DELETE FROM server WHERE id=$1`
const deleteInterfacesQuery = `DELETE FROM interface WHERE server=$1`
const deleteIPsQuery = `DELETE FROM ip_address WHERE server = $1`

func newUUID() *string {
	uuidReference := uuid.New().String()
	return &uuidReference
}

func validateCommon(s *tc.CommonServerProperties, tx *sql.Tx) ([]error, error) {

	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	errs := tovalidate.ToErrors(validation.Errors{
		"cachegroupId":   validation.Validate(s.CachegroupID, validation.NotNil),
		"cdnId":          validation.Validate(s.CDNID, validation.NotNil),
		"domainName":     validation.Validate(s.DomainName, validation.Required, noSpaces),
		"hostName":       validation.Validate(s.HostName, validation.Required, noSpaces),
		"physLocationId": validation.Validate(s.PhysLocationID, validation.NotNil),
		"profileId":      validation.Validate(s.ProfileID, validation.NotNil),
		"statusId":       validation.Validate(s.StatusID, validation.NotNil),
		"typeId":         validation.Validate(s.TypeID, validation.NotNil),
		"httpsPort":      validation.Validate(s.HTTPSPort, validation.By(tovalidate.IsValidPortNumber)),
		"tcpPort":        validation.Validate(s.TCPPort, validation.By(tovalidate.IsValidPortNumber)),
	})

	if len(errs) > 0 {
		return errs, nil
	}

	if _, err := tc.ValidateTypeID(tx, s.TypeID, "server"); err != nil {
		errs = append(errs, err)
	}

	var cdnID int
	if err := tx.QueryRow("SELECT cdn from profile WHERE id=$1", s.ProfileID).Scan(&cdnID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errs = append(errs, fmt.Errorf("no such Profile: #%d", *s.ProfileID))
			return errs, nil
		}
		return nil, fmt.Errorf("could not execute select cdnID from profile: %w", err)
	}

	log.Infof("got cdn id: %d from profile and cdn id: %d from server", cdnID, *s.CDNID)
	if cdnID != *s.CDNID {
		errs = append(errs, fmt.Errorf("CDN id '%d' for profile '%d' does not match Server CDN '%d'", cdnID, *s.ProfileID, *s.CDNID))
	}
	return errs, nil
}

func validateCommonV40(s *tc.ServerV40, tx *sql.Tx) ([]error, error) {

	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	errs := tovalidate.ToErrors(validation.Errors{
		"cachegroupId":   validation.Validate(s.CachegroupID, validation.NotNil),
		"cdnId":          validation.Validate(s.CDNID, validation.NotNil),
		"domainName":     validation.Validate(s.DomainName, validation.Required, noSpaces),
		"hostName":       validation.Validate(s.HostName, validation.Required, noSpaces),
		"physLocationId": validation.Validate(s.PhysLocationID, validation.NotNil),
		"profileNames":   validation.Validate(s.ProfileNames, validation.NotNil),
		"statusId":       validation.Validate(s.StatusID, validation.NotNil),
		"typeId":         validation.Validate(s.TypeID, validation.NotNil),
		"httpsPort":      validation.Validate(s.HTTPSPort, validation.By(tovalidate.IsValidPortNumber)),
		"tcpPort":        validation.Validate(s.TCPPort, validation.By(tovalidate.IsValidPortNumber)),
	})

	if len(errs) > 0 {
		return errs, nil
	}

	if _, err := tc.ValidateTypeID(tx, s.TypeID, "server"); err != nil {
		errs = append(errs, err)
	}

	if len(s.ProfileNames) == 0 {
		errs = append(errs, fmt.Errorf("no profiles exists"))
	}

	var cdnID int
	for _, profile := range s.ProfileNames {
		if err := tx.QueryRow("SELECT cdn from profile WHERE name=$1", profile).Scan(&cdnID); err != nil {
			log.Errorf("could not execute select cdnID from profile: %s\n", err)
			if errors.Is(err, sql.ErrNoRows) {
				errs = append(errs, fmt.Errorf("no such profileName: '%s'", profile))
			} else {
				return nil, fmt.Errorf("unable to get CDN ID for profile name '%s': %w", profile, err)
			}
			return errs, nil
		}

		log.Infof("got cdn id: %d from profile and cdn id: %d from server", cdnID, *s.CDNID)
		if cdnID != *s.CDNID {
			errs = append(errs, fmt.Errorf("CDN id '%d' for profile '%v' does not match Server CDN '%d'", cdnID, profile, *s.CDNID))
		}
	}

	return errs, nil
}

func validateMTU(mtu interface{}) error {
	m, ok := mtu.(*uint64)
	if !ok {
		return errors.New("must be an unsigned integer with 64-bit precision")
	}
	if m == nil {
		return nil
	}

	if *m < 1280 {
		return errors.New("must be at least 1280")
	}
	return nil
}

func validateV4(s *tc.ServerV40, tx *sql.Tx) (string, error, error) {
	if len(s.Interfaces) == 0 {
		return "", errors.New("a server must have at least one interface"), nil
	}
	var errs []error
	var serviceAddrV4Found bool
	var ipv4 string
	var serviceAddrV6Found bool
	var ipv6 string
	var serviceInterface string
	for _, iface := range s.Interfaces {
		ruleName := fmt.Sprintf("interface '%s' ", iface.Name)
		errs = append(errs, tovalidate.ToErrors(validation.Errors{
			ruleName + "name":        validation.Validate(iface.Name, validation.Required),
			ruleName + "mtu":         validation.Validate(iface.MTU, validation.By(validateMTU)),
			ruleName + "ipAddresses": validation.Validate(iface.IPAddresses, validation.Required),
		})...)

		for _, addr := range iface.IPAddresses {
			ruleName += fmt.Sprintf("address '%s'", addr.Address)

			var parsedIP net.IP
			var err error
			if parsedIP, _, err = net.ParseCIDR(addr.Address); err != nil {
				if parsedIP = net.ParseIP(addr.Address); parsedIP == nil {
					errs = append(errs, fmt.Errorf("%s: address: %v", ruleName, err))
					continue
				}
			}

			if addr.Gateway != nil {
				if gateway := net.ParseIP(*addr.Gateway); gateway == nil {
					errs = append(errs, fmt.Errorf("%s: gateway: could not parse '%s' as a network gateway", ruleName, *addr.Gateway))
				} else if (gateway.To4() == nil && parsedIP.To4() != nil) || (gateway.To4() != nil && parsedIP.To4() == nil) {
					errs = append(errs, errors.New(ruleName+": address family mismatch between address and gateway"))
				}
			}

			if addr.ServiceAddress {
				if serviceInterface != "" && serviceInterface != iface.Name {
					errs = append(errs, fmt.Errorf("interfaces: both %s and %s interfaces contain service addresses - only one service-address-containing-interface is allowed", serviceInterface, iface.Name))
				}
				serviceInterface = iface.Name
				if parsedIP.To4() != nil {
					if serviceAddrV4Found {
						errs = append(errs, fmt.Errorf("interfaces: address '%s' of interface '%s' is marked as a service address, but an IPv4 service address appears earlier in the list", addr.Address, iface.Name))
					}
					serviceAddrV4Found = true
					ipv4 = addr.Address
				} else {
					if serviceAddrV6Found {
						errs = append(errs, fmt.Errorf("interfaces: address '%s' of interface '%s' is marked as a service address, but an IPv6 service address appears earlier in the list", addr.Address, iface.Name))
					}
					serviceAddrV6Found = true
					ipv6 = addr.Address
				}
			}
		}
	}

	if !serviceAddrV6Found && !serviceAddrV4Found {
		errs = append(errs, errors.New("a server must have at least one service address"))
	}
	usrErr, sysErr := validateCommonV40(s, tx)
	errs = append(errs, usrErr...)
	if sysErr != nil || len(errs) > 0 {
		return serviceInterface, util.JoinErrs(errs), sysErr
	}
	query := `
SELECT tmp.server, ip.address
FROM (
  SELECT server, ARRAY_AGG(profile_name order by priority) AS profiles
	FROM server_profile
	GROUP BY server
) AS tmp
JOIN ip_address ip on ip.server = tmp.server
WHERE (profiles = $1::text[])
`
	var rows *sql.Rows
	var err error
	//ProfileID already validated
	if s.ID != nil {
		rows, err = tx.Query(query+" and tmp.server != $2", pq.Array(s.ProfileNames), *s.ID)
	} else {
		rows, err = tx.Query(query, pq.Array(s.ProfileNames))
	}
	if err != nil {
		errs = append(errs, errors.New("unable to determine service address uniqueness"))
	} else if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var ipaddress string
			err = rows.Scan(&id, &ipaddress)
			if err != nil {
				errs = append(errs, errors.New("unable to determine service address uniqueness"))
			} else if (ipaddress == ipv4 || ipaddress == ipv6) && (s.ID == nil || *s.ID != id) {
				errs = append(errs, fmt.Errorf("there exists a server with id %v on the same profile that has the same service address %s", id, ipaddress))
			}
		}
	}

	return serviceInterface, util.JoinErrs(errs), nil
}

func validateV3(s *tc.ServerV30, tx *sql.Tx) (string, error, error) {

	if len(s.Interfaces) == 0 {
		return "", errors.New("a server must have at least one interface"), nil
	}
	var errs []error
	var serviceAddrV4Found bool
	var ipv4 string
	var serviceAddrV6Found bool
	var ipv6 string
	var serviceInterface string
	for _, iface := range s.Interfaces {

		ruleName := fmt.Sprintf("interface '%s' ", iface.Name)
		errs = append(errs, tovalidate.ToErrors(validation.Errors{
			ruleName + "name":        validation.Validate(iface.Name, validation.Required),
			ruleName + "mtu":         validation.Validate(iface.MTU, validation.By(validateMTU)),
			ruleName + "ipAddresses": validation.Validate(iface.IPAddresses, validation.Required),
		})...)

		for _, addr := range iface.IPAddresses {
			ruleName += fmt.Sprintf("address '%s'", addr.Address)

			var parsedIP net.IP
			var err error
			if parsedIP, _, err = net.ParseCIDR(addr.Address); err != nil {
				if parsedIP = net.ParseIP(addr.Address); parsedIP == nil {
					errs = append(errs, fmt.Errorf("%s: address: %v", ruleName, err))
					continue
				}
			}

			if addr.Gateway != nil {
				if gateway := net.ParseIP(*addr.Gateway); gateway == nil {
					errs = append(errs, fmt.Errorf("%s: gateway: could not parse '%s' as a network gateway", ruleName, *addr.Gateway))
				} else if (gateway.To4() == nil && parsedIP.To4() != nil) || (gateway.To4() != nil && parsedIP.To4() == nil) {
					errs = append(errs, errors.New(ruleName+": address family mismatch between address and gateway"))
				}
			}

			if addr.ServiceAddress {
				if serviceInterface != "" && serviceInterface != iface.Name {
					errs = append(errs, fmt.Errorf("interfaces: both %s and %s interfaces contain service addresses - only one service-address-containing-interface is allowed", serviceInterface, iface.Name))
				}
				serviceInterface = iface.Name
				if parsedIP.To4() != nil {
					if serviceAddrV4Found {
						errs = append(errs, fmt.Errorf("interfaces: address '%s' of interface '%s' is marked as a service address, but an IPv4 service address appears earlier in the list", addr.Address, iface.Name))
					}
					serviceAddrV4Found = true
					ipv4 = addr.Address
				} else {
					if serviceAddrV6Found {
						errs = append(errs, fmt.Errorf("interfaces: address '%s' of interface '%s' is marked as a service address, but an IPv6 service address appears earlier in the list", addr.Address, iface.Name))
					}
					serviceAddrV6Found = true
					ipv6 = addr.Address
				}
			}
		}
	}

	if !serviceAddrV6Found && !serviceAddrV4Found {
		errs = append(errs, errors.New("a server must have at least one service address"))
	}

	commonErrs, sysErr := validateCommon(&s.CommonServerProperties, tx)
	errs = append(errs, commonErrs...)
	if len(errs) > 0 || sysErr != nil {
		return serviceInterface, util.JoinErrs(errs), sysErr
	}
	query := `
SELECT s.ID, ip.address FROM server s
JOIN profile p on p.Id = s.Profile
JOIN interface i on i.server = s.ID
JOIN ip_address ip on ip.Server = s.ID and ip.interface = i.name
WHERE ip.service_address = true
and p.id = $1
`
	var rows *sql.Rows
	var err error
	//ProfileID already validated
	if s.ID != nil {
		rows, err = tx.Query(query+" and s.id != $2", *s.ProfileID, *s.ID)
	} else {
		rows, err = tx.Query(query, *s.ProfileID)
	}
	if err != nil {
		return serviceInterface, util.JoinErrs(errs), fmt.Errorf("unable to determine service address uniqueness: querying: %w", err)
	} else if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var ipaddress string
			err = rows.Scan(&id, &ipaddress)
			if err != nil {
				return serviceInterface, util.JoinErrs(errs), fmt.Errorf("unable to determine service address uniqueness: scanning: %w", err)
			} else if (ipaddress == ipv4 || ipaddress == ipv6) && (s.ID == nil || *s.ID != id) {
				errs = append(errs, fmt.Errorf("there exists a server with id %v on the same profile that has the same service address %s", id, ipaddress))
			}
		}
	}

	return serviceInterface, util.JoinErrs(errs), nil
}

// Read is the handler for GET requests to /servers.
func Read(w http.ResponseWriter, r *http.Request) {
	var maxTime *time.Time
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	servers := []tc.ServerV40{}
	var serverCount uint64
	cfg, e := api.GetConfig(r.Context())
	useIMS := false
	if e == nil && cfg != nil {
		useIMS = cfg.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	servers, serverCount, userErr, sysErr, errCode, maxTime = getServers(r.Header, inf.Params, inf.Tx, inf.User, useIMS, *version)
	if maxTime != nil && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, *maxTime)
	}
	if errCode == http.StatusNotModified {
		w.WriteHeader(errCode)
		api.WriteResp(w, r, []tc.ServerV40{})
		return
	}
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if version.Major >= 4 {
		if version.Minor >= 1 {
			api.WriteRespWithSummary(w, r, servers, serverCount)
			return
		}
		v40Servers := make([]tc.ServerV40, 0)
		for _, server := range servers {
			v40Servers = append(v40Servers, server)
		}
		api.WriteRespWithSummary(w, r, v40Servers, serverCount)
		return
	}
	v3Servers := make([]tc.ServerV30, 0)
	for _, server := range servers {
		csp, err := dbhelpers.GetCommonServerPropertiesFromV4(server, tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("failed to get common server properties from V4 server struct: %w", err))
			return
		}
		v3Server, err := server.ToServerV3FromV4(csp)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("failed to convert servers to V3 format: %w", err))
			return
		}
		v3Servers = append(v3Servers, v3Server)
	}
	api.WriteRespWithSummary(w, r, v3Servers, serverCount)
	return
}

func selectMaxLastUpdatedQuery(queryAddition string, where string) string {
	return `SELECT max(t) from (
		SELECT max(s.last_updated) as t from server s JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id ` +
		queryAddition + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='server') as res`
}

func getServerCount(tx *sqlx.Tx, query string, queryValues map[string]interface{}) (uint64, error) {
	var serverCount uint64
	ns, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("couldn't prepare the query to get server count : %v", err)
	}
	err = tx.NamedStmt(ns).QueryRow(queryValues).Scan(&serverCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get servers count: %v", err)
	}
	return serverCount, nil
}

func getServers(h http.Header, params map[string]string, tx *sqlx.Tx, user *auth.CurrentUser, useIMS bool, version api.Version) ([]tc.ServerV40, uint64, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"cachegroup":       {Column: "s.cachegroup", Checker: api.IsInt},
		"parentCachegroup": {Column: "cg.parent_cachegroup_id", Checker: api.IsInt},
		"cachegroupName":   {Column: "cg.name", Checker: nil},
		"cdn":              {Column: "s.cdn_id", Checker: api.IsInt},
		"id":               {Column: "s.id", Checker: api.IsInt},
		"hostName":         {Column: "s.host_name", Checker: nil},
		"physLocation":     {Column: "s.phys_location", Checker: api.IsInt},
		"status":           {Column: "st.name", Checker: nil},
		"topology":         {Column: "tc.topology", Checker: nil},
		"type":             {Column: "t.name", Checker: nil},
		"dsId":             {Column: "dss.deliveryservice", Checker: nil},
	}

	if version.Major == 3 {
		queryParamsToSQLCols["profileId"] = dbhelpers.WhereColumnInfo{
			Column:  "s.profile",
			Checker: api.IsInt,
		}
	}

	if version.Major >= 4 {
		queryParamsToSQLCols["profileName"] = dbhelpers.WhereColumnInfo{
			Column:  "sp.profile_name",
			Checker: nil,
		}
	}

	usesMids := false
	queryAddition := ""
	dsHasRequiredCapabilities := false
	var requiredCapabilities []string
	var dsID int
	var cdnID int
	var serverCount uint64
	var err error

	if dsIDStr, ok := params[`dsId`]; ok {
		// don't allow query on ds outside user's tenant
		dsID, err = strconv.Atoi(dsIDStr)
		if err != nil {
			return nil, 0, errors.New("dsId must be an integer"), nil, http.StatusNotFound, nil
		}
		cdnID, _, err = dbhelpers.GetDSCDNIdFromID(tx.Tx, dsID)
		if err != nil {
			return nil, 0, nil, err, http.StatusInternalServerError, nil
		}

		userErr, sysErr, _ := tenant.CheckID(tx.Tx, user, dsID)
		if userErr != nil || sysErr != nil {
			return nil, 0, errors.New("Forbidden"), sysErr, http.StatusForbidden, nil
		}

		var joinSubQuery string
		if err = tx.QueryRow(deliveryservice.GetRequiredCapabilitiesQuery, dsID).Scan(pq.Array(&requiredCapabilities)); err != nil {
			err = fmt.Errorf("unable to get required capabilities for deliveryservice %d: %s", dsID, err)
			return nil, 0, nil, err, http.StatusInternalServerError, nil
		}
		if requiredCapabilities != nil && len(requiredCapabilities) > 0 {
			dsHasRequiredCapabilities = true
		}
		joinSubQuery = dssTopologiesJoinSubquery
		// only if dsId is part of params: add join on deliveryservice_server table
		queryAddition = fmt.Sprintf(deliveryServiceServersJoin, joinSubQuery)

		// depending on ds type, also need to add mids
		dsType, _, _, err := dbhelpers.GetDeliveryServiceTypeAndCDNName(dsID, tx.Tx)
		if err != nil {
			return nil, 0, nil, err, http.StatusInternalServerError, nil
		}
		usesMids = dsType.UsesMidCache()
		log.Debugf("Servers for ds %d; uses mids? %v\n", dsID, usesMids)
	}

	if _, ok := params[`topology`]; ok {
		/* language=SQL */
		queryAddition += `
			JOIN topology_cachegroup tc ON cg."name" = tc.cachegroup
`
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if dsHasRequiredCapabilities {
		where += requiredCapabilitiesCondition
	}
	if len(errs) > 0 {
		return nil, 0, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	var queryString, countQueryString string
	queryString = selectQuery
	countQueryString = serverCountQuery
	if version.Major >= 4 {
		if _, ok := params["profileName"]; ok {
			queryString = queryString + `
JOIN server_profile sp ON s.id = sp.server`
			countQueryString = countQueryString + `
JOIN server_profile sp ON s.id = sp.server`
		} else {
			queryString = queryString + ` ` + joinProfileV4
			countQueryString = countQueryString + ` ` + joinProfileV4
		}
	}
	countQuery := countQueryString + queryAddition + where
	// If we are querying for a DS that has reqd capabilities, we need to make sure that we also include all the ORG servers directly assigned to this DS
	if _, ok := params["dsId"]; ok && dsHasRequiredCapabilities {
		countQuery = `SELECT (` + countQuery + `) + (` + countQueryString + originServerQuery + `) AS total`
	}
	serverCount, err = getServerCount(tx, countQuery, queryValues)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to get servers count: %v", err), http.StatusInternalServerError, nil
	}

	serversList := []tc.ServerV40{}
	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, h, queryValues, selectMaxLastUpdatedQuery(queryAddition, where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return serversList, 0, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	query := queryString + queryAddition + where + orderBy + pagination
	// If you're looking to get the servers for a particular delivery service, make sure you're also querying the ORG servers from the deliveryservice_server table
	if _, ok := params[`dsId`]; ok {
		query = `(` + queryString + queryAddition + where + orderBy + pagination + `) UNION ` + queryString + originServerQuery
	}

	log.Debugln("Query is ", query)
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, serverCount, nil, errors.New("querying: " + err.Error()), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	HiddenField := "********"

	servers := make(map[int]tc.ServerV40)
	ids := []int{}
	for rows.Next() {
		s := tc.ServerV40{}
		err := rows.Scan(&s.Cachegroup,
			&s.CachegroupID,
			&s.CDNID,
			&s.CDNName,
			&s.DomainName,
			&s.GUID,
			&s.HostName,
			&s.HTTPSPort,
			&s.ID,
			&s.ILOIPAddress,
			&s.ILOIPGateway,
			&s.ILOIPNetmask,
			&s.ILOPassword,
			&s.ILOUsername,
			&s.LastUpdated,
			&s.MgmtIPAddress,
			&s.MgmtIPGateway,
			&s.MgmtIPNetmask,
			&s.OfflineReason,
			&s.PhysLocation,
			&s.PhysLocationID,
			pq.Array(&s.ProfileNames),
			&s.Rack,
			&s.RevalPending,
			&s.RevalUpdateTime,
			&s.RevalApplyTime,
			&s.Status,
			&s.StatusID,
			&s.TCPPort,
			&s.Type,
			&s.TypeID,
			&s.UpdPending,
			&s.ConfigUpdateTime,
			&s.ConfigApplyTime,
			&s.XMPPID,
			&s.XMPPPasswd,
			&s.StatusLastUpdated)
		if err != nil {
			return nil, serverCount, nil, errors.New("getting servers: " + err.Error()), http.StatusInternalServerError, nil
		}
		if user.PrivLevel < auth.PrivLevelOperations {
			s.ILOPassword = &HiddenField
			s.XMPPPasswd = &HiddenField
		}

		if s.ID == nil {
			return nil, serverCount, nil, errors.New("found server with nil ID"), http.StatusInternalServerError, nil
		}
		if _, ok := servers[*s.ID]; ok {
			return nil, serverCount, nil, fmt.Errorf("found more than one server with ID #%d", *s.ID), http.StatusInternalServerError, nil
		}
		servers[*s.ID] = s
		ids = append(ids, *s.ID)
	}

	// if ds requested uses mid-tier caches, add those to the list as well
	if usesMids {
		midIDs, userErr, sysErr, errCode := getMidServers(ids, servers, dsID, cdnID, tx, dsHasRequiredCapabilities)

		log.Debugf("getting mids: %v, %v, %s\n", userErr, sysErr, http.StatusText(errCode))

		serverCount = serverCount + uint64(len(midIDs))
		if userErr != nil || sysErr != nil {
			return nil, serverCount, userErr, sysErr, errCode, nil
		}
		ids = append(ids, midIDs...)
	}

	if len(ids) < 1 {
		return []tc.ServerV40{}, serverCount, nil, nil, http.StatusOK, nil
	}

	query, args, err := sqlx.In(`SELECT max_bandwidth, monitor, mtu, name, server, router_host_name, router_port_name FROM interface WHERE server IN (?)`, ids)
	if err != nil {
		return nil, serverCount, nil, fmt.Errorf("building interfaces query: %v", err), http.StatusInternalServerError, nil
	}
	query = tx.Rebind(query)
	interfaces := map[int]map[string]tc.ServerInterfaceInfoV40{}
	interfaceRows, err := tx.Queryx(query, args...)
	if err != nil {
		return nil, serverCount, nil, fmt.Errorf("querying for interfaces: %v", err), http.StatusInternalServerError, nil
	}
	defer interfaceRows.Close()

	for interfaceRows.Next() {
		iface := tc.ServerInterfaceInfoV40{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{},
			},
		}
		var server int
		var routerHostName string
		var routerPort string
		if err = interfaceRows.Scan(&iface.MaxBandwidth, &iface.Monitor, &iface.MTU, &iface.Name, &server, &routerHostName, &routerPort); err != nil {
			return nil, serverCount, nil, fmt.Errorf("getting server interfaces: %v", err), http.StatusInternalServerError, nil
		}

		if _, ok := servers[server]; !ok {
			continue
		}

		if _, ok := interfaces[server]; !ok {
			interfaces[server] = map[string]tc.ServerInterfaceInfoV40{}
		}
		iface.RouterHostName = routerHostName
		iface.RouterPortName = routerPort
		interfaces[server][iface.Name] = iface
	}

	query, args, err = sqlx.In(`SELECT address, gateway, service_address, server, interface FROM ip_address WHERE server IN (?)`, ids)
	if err != nil {
		return nil, serverCount, nil, fmt.Errorf("building IP addresses query: %v", err), http.StatusInternalServerError, nil
	}
	query = tx.Rebind(query)
	ipRows, err := tx.Tx.Query(query, args...)
	if err != nil {
		return nil, serverCount, nil, fmt.Errorf("querying for IP addresses: %v", err), http.StatusInternalServerError, nil
	}
	defer ipRows.Close()

	for ipRows.Next() {
		var ip tc.ServerIPAddress
		var server int
		var iface string

		if err = ipRows.Scan(&ip.Address, &ip.Gateway, &ip.ServiceAddress, &server, &iface); err != nil {
			return nil, serverCount, nil, fmt.Errorf("getting server IP addresses: %v", err), http.StatusInternalServerError, nil
		}

		if _, ok := interfaces[server]; !ok {
			continue
		}
		if i, ok := interfaces[server][iface]; !ok {
			log.Warnf("IP addresses query returned addresses for an interface that was not found in interfaces query: %s", iface)
		} else {
			i.IPAddresses = append(i.IPAddresses, ip)
			interfaces[server][iface] = i
		}
	}

	returnable := make([]tc.ServerV40, 0, len(ids))

	for _, id := range ids {
		server := servers[id]
		for _, iface := range interfaces[id] {
			server.Interfaces = append(server.Interfaces, iface)
		}
		returnable = append(returnable, server)
	}

	return returnable, serverCount, nil, nil, http.StatusOK, &maxTime
}

// getMidServers gets the mids used by the edges provided with an option to filter for a given cdn
func getMidServers(edgeIDs []int, servers map[int]tc.ServerV40, dsID int, cdnID int, tx *sqlx.Tx, includeCapabilities bool) ([]int, error, error, int) {
	if len(edgeIDs) == 0 {
		return nil, nil, nil, http.StatusOK
	}

	filters := map[string]interface{}{
		"cache_type_mid": tc.CacheTypeMid,
		"edge_ids":       pq.Array(edgeIDs),
		"ds_id":          dsID,
	}

	midIDs := []int{}
	query := ""
	if includeCapabilities {
		// Query to select the associated mids for this DS
		q := selectIDQuery + midWhereClause
		rows, err := tx.NamedQuery(q, filters)
		if err != nil {
			return nil, err, nil, http.StatusBadRequest
		}
		defer rows.Close()

		for rows.Next() {
			var midID int
			if err := rows.Scan(&midID); err != nil {
				log.Errorf("could not scan mid server id: %s\n", err)
				return nil, nil, err, http.StatusInternalServerError
			}
			midIDs = append(midIDs, midID)
		}
		filters["mid_ids"] = pq.Array(midIDs)

		// Query to select only those mids that match the required capabilities of the DS
		query = selectQuery + midWhereClause + `
		AND s.id IN (
		WITH capabilities AS (
		SELECT ARRAY_AGG(ssc.server_capability), server
		FROM server_server_capability ssc
		WHERE ssc.server = ANY(:mid_ids)
		GROUP BY server)
		SELECT server
		FROM capabilities WHERE
		capabilities.array_agg
		@>
		(
		SELECT ds.required_capabilities
		FROM deliveryservice ds
		WHERE ds.id=:ds_id)
		)`
	} else {
		// TODO: include secondary parent?
		query = selectQuery + midWhereClause
	}

	if cdnID > 0 {
		query += ` AND s.cdn_id = :cdn_id`
		filters["cdn_id"] = cdnID
	}

	rows, err := tx.NamedQuery(query, filters)
	if err != nil {
		return nil, err, nil, http.StatusBadRequest
	}
	defer rows.Close()

	ids := []int{}

	for rows.Next() {
		var s tc.ServerV40
		if err := rows.Scan(&s.Cachegroup,
			&s.CachegroupID,
			&s.CDNID,
			&s.CDNName,
			&s.DomainName,
			&s.GUID,
			&s.HostName,
			&s.HTTPSPort,
			&s.ID,
			&s.ILOIPAddress,
			&s.ILOIPGateway,
			&s.ILOIPNetmask,
			&s.ILOPassword,
			&s.ILOUsername,
			&s.LastUpdated,
			&s.MgmtIPAddress,
			&s.MgmtIPGateway,
			&s.MgmtIPNetmask,
			&s.OfflineReason,
			&s.PhysLocation,
			&s.PhysLocationID,
			pq.Array(&s.ProfileNames),
			&s.Rack,
			&s.RevalPending,
			&s.RevalUpdateTime,
			&s.RevalApplyTime,
			&s.Status,
			&s.StatusID,
			&s.TCPPort,
			&s.Type,
			&s.TypeID,
			&s.UpdPending,
			&s.ConfigUpdateTime,
			&s.ConfigApplyTime,
			&s.XMPPID,
			&s.XMPPPasswd,
			&s.StatusLastUpdated); err != nil {
			log.Errorf("could not scan mid servers: %s\n", err)
			return nil, nil, err, http.StatusInternalServerError
		}
		if s.ID == nil {
			return nil, nil, errors.New("found server with nil ID"), http.StatusInternalServerError
		}

		// This may mean that the server was caught by other query parameters,
		// so not technically an error, unlike earlier in 'getServers'.
		if _, ok := servers[*s.ID]; ok {
			continue
		}

		servers[*s.ID] = s
		ids = append(ids, *s.ID)
	}

	return ids, nil, nil, http.StatusOK
}

func checkTypeChangeSafety(server tc.ServerV40, tx *sqlx.Tx) (error, error, int) {
	// see if cdn or type changed
	var cdnID int
	var typeID int
	if err := tx.QueryRow("SELECT type, cdn_id FROM server WHERE id = $1", *server.ID).Scan(&typeID, &cdnID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no server found with this ID"), nil, http.StatusNotFound
		}
		return nil, fmt.Errorf("getting current server type: %v", err), http.StatusInternalServerError
	}

	var dsIDs []int64
	if err := tx.QueryRowx("SELECT ARRAY(SELECT deliveryservice FROM deliveryservice_server WHERE server = $1)", server.ID).Scan(pq.Array(&dsIDs)); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("getting server assigned delivery services: %v", err), http.StatusInternalServerError
	}
	// If type is changing ensure it isn't assigned to any DSes.
	if typeID != *server.TypeID {
		if len(dsIDs) != 0 {
			return errors.New("server type can not be updated when it is currently assigned to Delivery Services"), nil, http.StatusConflict
		}
	}
	// Check to see if the user is trying to change the CDN of a server, which is already linked with a DS
	if cdnID != *server.CDNID && len(dsIDs) != 0 {
		return errors.New("server cdn can not be updated when it is currently assigned to delivery services"), nil, http.StatusConflict
	}

	return nil, nil, http.StatusOK
}

func updateStatusLastUpdatedTime(id int, statusLastUpdatedTime *time.Time, tx *sql.Tx) (error, error, int) {
	query := `UPDATE server SET
	status_last_updated=$1
WHERE id=$2 `
	if _, err := tx.Exec(query, statusLastUpdatedTime, id); err != nil {
		return errors.New("updating status last updated: " + err.Error()), nil, http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}

func createInterfaces(id int, interfaces []tc.ServerInterfaceInfoV40, tx *sql.Tx) (error, error, int) {
	ifaceQry := `
	INSERT INTO interface (
		max_bandwidth,
		monitor,
		mtu,
		name,
		server,
		router_host_name,
		router_port_name
	) VALUES
	`
	ipQry := `
	INSERT INTO ip_address (
		address,
		gateway,
		interface,
		server,
		service_address
	) VALUES
	`

	ifaceQueryParts := make([]string, 0, len(interfaces))
	ipQueryParts := make([]string, 0, len(interfaces))
	ifaceArgs := make([]interface{}, 0, len(interfaces))
	ipArgs := make([]interface{}, 0, len(interfaces))
	for i, iface := range interfaces {
		argStart := i * 7
		ifaceQueryParts = append(ifaceQueryParts, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", argStart+1, argStart+2, argStart+3, argStart+4, argStart+5, argStart+6, argStart+7))
		ifaceArgs = append(ifaceArgs, iface.MaxBandwidth, iface.Monitor, iface.MTU, iface.Name, id, iface.RouterHostName, iface.RouterPortName)
		for _, ip := range iface.IPAddresses {
			argStart = len(ipArgs)
			ipQueryParts = append(ipQueryParts, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", argStart+1, argStart+2, argStart+3, argStart+4, argStart+5))
			ipArgs = append(ipArgs, ip.Address, ip.Gateway, iface.Name, id, ip.ServiceAddress)
		}
	}

	ifaceQry += strings.Join(ifaceQueryParts, ",")
	log.Debugf("Inserting interfaces for new server, query is: %s", ifaceQry)

	_, err := tx.Exec(ifaceQry, ifaceArgs...)
	if err != nil {
		return api.ParseDBError(err)
	}

	ipQry += strings.Join(ipQueryParts, ",")
	log.Debugf("Inserting IP addresses for new server, query is: %s", ipQry)

	_, err = tx.Exec(ipQry, ipArgs...)
	if err != nil {
		return api.ParseDBError(err)
	}
	return nil, nil, http.StatusOK
}

func deleteInterfaces(id int, tx *sql.Tx) (error, error, int) {
	if _, err := tx.Exec(deleteIPsQuery, id); err != nil && err != sql.ErrNoRows {
		return api.ParseDBError(err)
	}

	if _, err := tx.Exec(deleteInterfacesQuery, id); err != nil && err != sql.ErrNoRows {
		return api.ParseDBError(err)
	}

	return nil, nil, http.StatusOK
}

// Update is the handler for PUT requests to /servers.
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	id := inf.IntParams["id"]

	// Get original server
	originals, _, userErr, sysErr, errCode, _ := getServers(r.Header, inf.Params, inf.Tx, inf.User, false, *version)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if len(originals) < 1 {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("the server doesn't exist, cannot update"), nil)
		return
	}
	if len(originals) > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("too many servers by ID %d: %d", id, len(originals)))
		return
	}

	original := originals[0]
	if original.XMPPID == nil || *original.XMPPID == "" {
		log.Warnf("original server %s had no XMPPID\n", *original.HostName)
	}
	if original.StatusID == nil {
		sysErr = errors.New("original server had no status ID")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if original.Status == nil {
		sysErr = errors.New("original server had no status name")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if original.CachegroupID == nil {
		sysErr = errors.New("original server had no Cache Group ID")
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if original.StatusLastUpdated == nil {
		log.Warnln("original server had no Status Last Updated time")
		if original.LastUpdated == nil {
			sysErr = errors.New("original server had no Last Updated time")
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		original.StatusLastUpdated = &original.LastUpdated.Time
	}

	var originalXMPPID string
	if original.XMPPID != nil {
		originalXMPPID = *original.XMPPID
	}
	originalStatusID := *original.StatusID

	var server tc.ServerV40
	var serverV3 tc.ServerV30
	var statusLastUpdatedTime time.Time

	if inf.Version.Major >= 4 {
		server.ID = new(int)
		*server.ID = inf.IntParams["id"]
		if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
		if server.StatusID != nil && *server.StatusID != originalStatusID {
			currentTime := time.Now()
			server.StatusLastUpdated = &currentTime
			statusLastUpdatedTime = currentTime
		} else {
			server.StatusLastUpdated = original.StatusLastUpdated
			statusLastUpdatedTime = *original.StatusLastUpdated
		}
		_, userErr, sysErr := validateV4(&server, tx)
		if userErr != nil || sysErr != nil {
			errCode := http.StatusBadRequest
			if sysErr != nil {
				errCode = http.StatusInternalServerError
			}
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
	} else {
		serverV3.ID = new(int)
		*serverV3.ID = inf.IntParams["id"]
		if err := json.NewDecoder(r.Body).Decode(&serverV3); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
		if serverV3.StatusID != nil && *serverV3.StatusID != originalStatusID {
			currentTime := time.Now()
			serverV3.StatusLastUpdated = &currentTime
			statusLastUpdatedTime = currentTime
		} else {
			serverV3.StatusLastUpdated = original.StatusLastUpdated
			statusLastUpdatedTime = *original.StatusLastUpdated
		}
		_, userErr, sysErr := validateV3(&serverV3, tx)
		if userErr != nil || sysErr != nil {
			code := http.StatusBadRequest
			if sysErr != nil {
				code = http.StatusInternalServerError
			}
			api.HandleErr(w, r, tx, code, userErr, sysErr)
			return
		}

		profileName, exists, err := dbhelpers.GetProfileNameFromID(*serverV3.ProfileID, tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		if !exists {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("profile does not exist"), nil)
			return
		}
		profileNames := []string{profileName}

		server, err = serverV3.UpgradeToV40(profileNames)
		if err != nil {
			sysErr = fmt.Errorf("error upgrading valid V3 server to V4 structure: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
	}

	if *original.CachegroupID != *server.CachegroupID || *original.CDNID != *server.CDNID {
		hasDSOnCDN, err := dbhelpers.CachegroupHasTopologyBasedDeliveryServicesOnCDN(tx, *original.CachegroupID, *original.CDNID)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		CDNIDs := []int{}
		if hasDSOnCDN {
			CDNIDs = append(CDNIDs, *original.CDNID)
		}
		cacheGroupIds := []int{*original.CachegroupID}
		serverIds := []int{*original.ID}
		if err = topology_validation.CheckForEmptyCacheGroups(inf.Tx, cacheGroupIds, CDNIDs, true, serverIds); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("server is the last one in its cachegroup, which is used by a topology, so it cannot be moved to another cachegroup: "+err.Error()), nil)
			return
		}
	}

	status, ok, err := dbhelpers.GetStatusByID(*server.StatusID, tx)
	if err != nil {
		sysErr = fmt.Errorf("getting server #%d status (#%d): %v", id, *server.StatusID, err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if !ok {
		log.Warnf("previously existent status #%d not found when fetching later", *server.StatusID)
		userErr = fmt.Errorf("no such Status: #%d", *server.StatusID)
		api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
		return
	}
	if status.Name == nil {
		sysErr = fmt.Errorf("status #%d had no name", *server.StatusID)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	}
	if *status.Name != string(tc.CacheStatusOnline) && *status.Name != string(tc.CacheStatusReported) {
		dsIDs, err := getActiveDeliveryServicesThatOnlyHaveThisServerAssigned(id, original.Type, tx)
		if err != nil {
			sysErr = fmt.Errorf("getting Delivery Services to which server #%d is assigned that have no other servers: %v", id, err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}
		if len(dsIDs) > 0 {
			prefix := fmt.Sprintf("setting server status to '%s' would leave Active Delivery Service", *status.Name)
			alertText := InvalidStatusForDeliveryServicesAlertText(prefix, original.Type, dsIDs)
			api.WriteAlerts(w, r, http.StatusConflict, tc.CreateAlerts(tc.ErrorLevel, alertText))
			return
		}
	}

	if userErr, sysErr, errCode = checkTypeChangeSafety(server, inf.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if server.XMPPID != nil && *server.XMPPID != "" && originalXMPPID != "" && *server.XMPPID != originalXMPPID {
		api.WriteAlerts(w, r, http.StatusBadRequest, tc.CreateAlerts(tc.ErrorLevel, fmt.Sprintf("server cannot be updated due to requested XMPPID change. XMPIDD is immutable")))
		return
	}

	userErr, sysErr, statusCode := api.CheckIfUnModified(r.Header, inf.Tx, *server.ID, "server")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	if server.CDNName != nil {
		userErr, sysErr, statusCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, *server.CDNName, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	} else if server.CDNID != nil {
		userErr, sysErr, statusCode = dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(inf.Tx.Tx, int64(*server.CDNID), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	if inf.Version.Major >= 4 {
		if err = dbhelpers.UpdateServerProfilesForV4(*server.ID, server.ProfileNames, tx); err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	} else {
		if err = dbhelpers.UpdateServerProfileTableForV3(serverV3.ID, serverV3.ProfileID, (original.ProfileNames)[0], tx); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("failed to update server_profile: %w", err))
			return
		}
	}

	serverID, errCode, userErr, sysErr := updateServer(inf.Tx, server)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	if userErr, sysErr, errCode = deleteInterfaces(id, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if userErr, sysErr, errCode = createInterfaces(id, server.Interfaces, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	where := `WHERE s.id = $1`
	var selquery string
	if version.Major <= 4 {
		selquery = selectQuery + joinProfileV4 + where
	} else {
		selquery = selectQuery + where
	}
	var srvr tc.ServerV40
	err = inf.Tx.QueryRow(selquery, serverID).Scan(&srvr.Cachegroup,
		&srvr.CachegroupID,
		&srvr.CDNID,
		&srvr.CDNName,
		&srvr.DomainName,
		&srvr.GUID,
		&srvr.HostName,
		&srvr.HTTPSPort,
		&srvr.ID,
		&srvr.ILOIPAddress,
		&srvr.ILOIPGateway,
		&srvr.ILOIPNetmask,
		&srvr.ILOPassword,
		&srvr.ILOUsername,
		&srvr.LastUpdated,
		&srvr.MgmtIPAddress,
		&srvr.MgmtIPGateway,
		&srvr.MgmtIPNetmask,
		&srvr.OfflineReason,
		&srvr.PhysLocation,
		&srvr.PhysLocationID,
		pq.Array(&srvr.ProfileNames),
		&srvr.Rack,
		&srvr.RevalPending,
		&srvr.RevalUpdateTime,
		&srvr.RevalApplyTime,
		&srvr.Status,
		&srvr.StatusID,
		&srvr.TCPPort,
		&srvr.Type,
		&srvr.TypeID,
		&srvr.UpdPending,
		&srvr.ConfigUpdateTime,
		&srvr.ConfigApplyTime,
		&srvr.XMPPID,
		&srvr.XMPPPasswd,
		&srvr.StatusLastUpdated)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	serversInterfaces, err := dbhelpers.GetServersInterfaces([]int{*srvr.ID}, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	if interfacesMap, ok := serversInterfaces[*srvr.ID]; ok {
		for _, intfc := range interfacesMap {
			srvr.Interfaces = append(srvr.Interfaces, intfc)
		}
	}

	if userErr, sysErr, errCode = updateStatusLastUpdatedTime(id, &statusLastUpdatedTime, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if inf.Version.Major >= 5 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", srvr)
	} else if inf.Version.Major >= 4 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", srvr)
	} else if inf.Version.Major >= 3 {
		csp, err := dbhelpers.GetCommonServerPropertiesFromV4(srvr, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
		srvrV30, err := srvr.ToServerV3FromV4(csp)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", srvrV30)
	}

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: updated", *srvr.HostName, *srvr.DomainName, *srvr.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func updateServer(tx *sqlx.Tx, server tc.ServerV40) (int64, int, error, error) {

	rows, err := tx.NamedQuery(updateQuery, server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return 0, errCode, userErr, sysErr
	}
	defer rows.Close()

	var serverId int64
	rowsAffected := 0
	for rows.Next() {
		if err := rows.Scan(&server.Cachegroup,
			&server.CachegroupID,
			&server.CDNID,
			&server.CDNName,
			&server.DomainName,
			&server.GUID,
			&server.HostName,
			&server.HTTPSPort,
			&serverId,
			&server.ILOIPAddress,
			&server.ILOIPGateway,
			&server.ILOIPNetmask,
			&server.ILOPassword,
			&server.ILOUsername,
			&server.LastUpdated,
			&server.MgmtIPAddress,
			&server.MgmtIPGateway,
			&server.MgmtIPNetmask,
			&server.OfflineReason,
			&server.PhysLocation,
			&server.PhysLocationID,
			pq.Array(&server.ProfileNames),
			&server.Rack,
			&server.Status,
			&server.StatusID,
			&server.TCPPort,
			&server.Type,
			&server.TypeID,
			&server.StatusLastUpdated); err != nil {
			return 0, http.StatusNotFound, nil, fmt.Errorf("scanning lastUpdated from server insert: %v", err)
		}
		rowsAffected++
	}

	if rowsAffected < 1 {
		return 0, http.StatusNotFound, fmt.Errorf("no server found with id %d", *server.ID), nil
	}
	if rowsAffected > 1 {
		return 0, http.StatusInternalServerError, nil, fmt.Errorf("update for server #%d affected too many rows (%d)", *server.ID, rowsAffected)
	}

	return serverId, http.StatusOK, nil, nil
}

func insertServerProfile(id int, pName []string, tx *sql.Tx) (error, error, int) {
	priority := make([]int, 0, len(pName))
	for i, _ := range pName {
		priority = append(priority, i)
	}
	insertQuery := `
	INSERT INTO server_profile (
		server,
		profile_name,
		priority
	)SELECT $1, profile_name, priority
	FROM UNNEST($2::text[], $3::int[]) WITH ORDINALITY AS tmp(profile_name, priority)
	`

	if _, err := tx.Exec(insertQuery, id, pq.Array(pName), pq.Array(priority)); err != nil {
		return api.ParseDBError(err)
	}
	return nil, nil, http.StatusOK
}

func createV3(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerV30

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if server.ID != nil {
		var prevID int
		err := inf.Tx.Tx.QueryRow("SELECT id from server where id = $1", server.ID).Scan(&prevID)
		if err != nil && err != sql.ErrNoRows {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("checking if server with id %d exists", *server.ID))
			return
		}
		if prevID != 0 {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("server with id %d already exists. Please do not provide an id", *server.ID), nil)
			return
		}
	}

	server.XMPPID = newUUID()

	_, userErr, sysErr := validateV3(&server, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, code, userErr, sysErr)
		return
	}

	currentTime := time.Now()
	server.StatusLastUpdated = &currentTime

	if server.CDNName != nil {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, *server.CDNName, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	} else if server.CDNID != nil {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(inf.Tx.Tx, int64(*server.CDNID), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	serverID, err := createServerV3(inf.Tx, server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	interfaces, err := tc.ToInterfacesV4(server.Interfaces, server.RouterHostName, server.RouterPortName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	userErr, sysErr, errCode := createInterfaces(int(serverID), interfaces, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	var origProfile string
	err = inf.Tx.Tx.QueryRow("SELECT name from profile where id = $1", server.ProfileID).Scan(&origProfile)
	if err != nil && err != sql.ErrNoRows {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("retreiving profile with id %d", *server.ProfileID))
		return
	}
	var origProfiles = []string{origProfile}

	userErr, sysErr, statusCode := insertServerProfile(int(serverID), origProfiles, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	where := `WHERE s.id = $1`
	selquery := selectQuery + where
	var s4 tc.ServerV40
	err = inf.Tx.QueryRow(selquery, serverID).Scan(&s4.Cachegroup,
		&s4.CachegroupID,
		&s4.CDNID,
		&s4.CDNName,
		&s4.DomainName,
		&s4.GUID,
		&s4.HostName,
		&s4.HTTPSPort,
		&s4.ID,
		&s4.ILOIPAddress,
		&s4.ILOIPGateway,
		&s4.ILOIPNetmask,
		&s4.ILOPassword,
		&s4.ILOUsername,
		&s4.LastUpdated,
		&s4.MgmtIPAddress,
		&s4.MgmtIPGateway,
		&s4.MgmtIPNetmask,
		&s4.OfflineReason,
		&s4.PhysLocation,
		&s4.PhysLocationID,
		pq.Array(&s4.ProfileNames),
		&s4.Rack,
		&s4.RevalPending,
		&s4.RevalUpdateTime,
		&s4.RevalApplyTime,
		&s4.Status,
		&s4.StatusID,
		&s4.TCPPort,
		&s4.Type,
		&s4.TypeID,
		&s4.UpdPending,
		&s4.ConfigUpdateTime,
		&s4.ConfigApplyTime,
		&s4.XMPPID,
		&s4.XMPPPasswd,
		&s4.StatusLastUpdated)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	s4.Interfaces = interfaces

	csp, err := dbhelpers.GetCommonServerPropertiesFromV4(s4, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
	}
	srvr, err := s4.ToServerV3FromV4(csp)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, srvr)

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: created", *srvr.HostName, *srvr.DomainName, *srvr.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, inf.Tx.Tx)
}

func createV4(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerV40

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, err, nil)
		return
	}

	if server.ID != nil {
		var prevID int
		err := inf.Tx.Tx.QueryRow("SELECT id from server where id = $1", server.ID).Scan(&prevID)
		if err != nil && err != sql.ErrNoRows {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("checking if server with id %d exists", *server.ID))
			return
		}
		if prevID != 0 {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("server with id %d already exists. Please do not provide an id", *server.ID), nil)
			return
		}
	}

	server.XMPPID = newUUID()

	_, userErr, sysErr := validateV4(&server, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		errCode := http.StatusBadRequest
		if sysErr != nil {
			errCode = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	currentTime := time.Now()
	server.StatusLastUpdated = &currentTime

	if server.CDNName != nil {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, *server.CDNName, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	} else if server.CDNID != nil {
		userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(inf.Tx.Tx, int64(*server.CDNID), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	origProfiles := server.ProfileNames
	serverID, err := createServerV4(inf.Tx, server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	userErr, sysErr, errCode := createInterfaces(int(serverID), server.Interfaces, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	userErr, sysErr, statusCode := insertServerProfile(int(serverID), origProfiles, inf.Tx.Tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	where := `WHERE s.id = $1`
	selquery := selectQuery + joinProfileV4 + where
	var srvr tc.ServerV40
	err = inf.Tx.QueryRow(selquery, serverID).Scan(&srvr.Cachegroup,
		&srvr.CachegroupID,
		&srvr.CDNID,
		&srvr.CDNName,
		&srvr.DomainName,
		&srvr.GUID,
		&srvr.HostName,
		&srvr.HTTPSPort,
		&srvr.ID,
		&srvr.ILOIPAddress,
		&srvr.ILOIPGateway,
		&srvr.ILOIPNetmask,
		&srvr.ILOPassword,
		&srvr.ILOUsername,
		&srvr.LastUpdated,
		&srvr.MgmtIPAddress,
		&srvr.MgmtIPGateway,
		&srvr.MgmtIPNetmask,
		&srvr.OfflineReason,
		&srvr.PhysLocation,
		&srvr.PhysLocationID,
		pq.Array(&srvr.ProfileNames),
		&srvr.Rack,
		&srvr.RevalPending,
		&srvr.RevalUpdateTime,
		&srvr.RevalApplyTime,
		&srvr.Status,
		&srvr.StatusID,
		&srvr.TCPPort,
		&srvr.Type,
		&srvr.TypeID,
		&srvr.UpdPending,
		&srvr.ConfigUpdateTime,
		&srvr.ConfigApplyTime,
		&srvr.XMPPID,
		&srvr.XMPPPasswd,
		&srvr.StatusLastUpdated)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}

	// TODO: Use returned values from SQL insert to ensure inserted values match
	srvr.Interfaces = server.Interfaces

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, srvr)

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: created", *srvr.HostName, *srvr.DomainName, *srvr.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, inf.Tx.Tx)
}

func createServerV4(tx *sqlx.Tx, server tc.ServerV40) (int64, error) {
	//rows, err := tx.NamedQuery(insertQueryV4, server)
	var profileID int
	err := tx.QueryRow("SELECT id FROM profile p WHERE name=$1", (server.ProfileNames)[0]).Scan(&profileID)
	if err != nil {
		return 0, fmt.Errorf("unable to get profileID for a profile name: %w", err)
	}

	rows, err := tx.Query(insertQueryV4, server.CachegroupID, server.CDNID, server.DomainName,
		server.HostName, server.HTTPSPort, server.ILOIPAddress, server.ILOIPNetmask, server.ILOIPGateway,
		server.ILOUsername, server.ILOPassword, server.MgmtIPAddress, server.MgmtIPNetmask, server.MgmtIPGateway,
		server.OfflineReason, server.PhysLocationID, profileID, server.Rack, server.StatusID,
		server.TCPPort, server.TypeID, server.XMPPID, server.XMPPPasswd)
	if err != nil {
		return 0, err
	}
	defer log.Close(rows, "failed to close rows for createServerV4")

	rowsAffected := 0
	var serverID int64
	for rows.Next() {
		rowsAffected++
		err := rows.Scan(&server.Cachegroup,
			&server.CachegroupID,
			&server.CDNID,
			&server.CDNName,
			&server.DomainName,
			&server.GUID,
			&server.HostName,
			&server.HTTPSPort,
			&serverID,
			&server.ILOIPAddress,
			&server.ILOIPGateway,
			&server.ILOIPNetmask,
			&server.ILOPassword,
			&server.ILOUsername,
			&server.LastUpdated,
			&server.MgmtIPAddress,
			&server.MgmtIPGateway,
			&server.MgmtIPNetmask,
			&server.OfflineReason,
			&server.PhysLocation,
			&server.PhysLocationID,
			pq.Array(&server.ProfileNames),
			&server.Rack,
			&server.Status,
			&server.StatusID,
			&server.TCPPort,
			&server.Type,
			&server.TypeID)
		if err != nil {
			return 0, err
		}
	}

	if rowsAffected == 0 {
		return 0, errors.New("server create: no server was inserted, no id was returned")
	} else if rowsAffected > 1 {
		return 0, fmt.Errorf("too many ids returned from server insert: %d", rowsAffected)
	}

	return serverID, nil
}

func createServerV3(tx *sqlx.Tx, server tc.ServerV30) (int64, error) {
	rows, err := tx.NamedQuery(insertQueryV3, server)
	if err != nil {
		return 0, err
	}
	defer log.Close(rows, "failed to close rows for createServerV3")

	rowsAffected := 0
	var serverID int64
	for rows.Next() {
		rowsAffected++
		if err := rows.Scan(&serverID); err != nil {
			return 0, err
		}
	}

	if rowsAffected == 0 {
		return 0, errors.New("server create: no server was inserted, no id was returned")
	} else if rowsAffected > 1 {
		return 0, fmt.Errorf("too many ids returned from server insert: %d", rowsAffected)
	}

	return serverID, nil
}

// Create is the handler for POST requests to /servers.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	switch {
	case inf.Version.Major == 3:
		createV3(inf, w, r)
	default:
		createV4(inf, w, r)
	}
}

const lastServerTypeOfDSesQuery = `
SELECT ds.id, ds.multi_site_origin, ds.topology
FROM deliveryservice_server dss
JOIN server s ON dss.server = s.id
JOIN type t ON s.type = t.id
JOIN deliveryservice ds ON dss.deliveryservice = ds.id
WHERE t.name LIKE $1 AND ds.active
GROUP BY ds.id, ds.multi_site_origin, ds.topology
HAVING COUNT(dss.server) = 1 AND $2 = ANY(ARRAY_AGG(dss.server));
`

// getActiveDeliveryServicesThatOnlyHaveThisServerAssigned returns the IDs of active delivery services for which the given
// server ID is either the last EDGE-type server or last ORG-type server (if MSO is enabled) assigned to them.
func getActiveDeliveryServicesThatOnlyHaveThisServerAssigned(id int, serverType string, tx *sql.Tx) ([]int, error) {
	var ids []int
	var like string
	isEdge := strings.HasPrefix(serverType, tc.CacheTypeEdge.String())
	isOrigin := strings.HasPrefix(serverType, tc.OriginTypeName)
	if isEdge {
		like = tc.CacheTypeEdge.String() + "%"
	} else if isOrigin {
		like = tc.OriginTypeName + "%"
	} else {
		// by definition, only EDGE-type or ORG-type servers can be assigned
		return ids, nil
	}
	if tx == nil {
		return ids, errors.New("nil transaction")
	}

	rows, err := tx.Query(lastServerTypeOfDSesQuery, like, id)
	if err != nil {
		return ids, fmt.Errorf("querying: %v", err)
	}
	defer log.Close(rows, "closing rows from getActiveDeliveryServicesThatOnlyHaveThisServerAssigned")

	for rows.Next() {
		var dsID int
		var mso bool
		var topology *string
		err = rows.Scan(&dsID, &mso, &topology)
		if err != nil {
			return ids, fmt.Errorf("scanning: %v", err)
		}
		if (isEdge && topology == nil) || (isOrigin && mso) {
			ids = append(ids, dsID)
		}
	}

	return ids, nil
}

// Delete is the handler for DELETE requests to the /servers API endpoint.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	// Middleware should've already handled this, so idk why this is a pointer at all tbh
	version := inf.Version
	if version == nil {
		middleware.NotImplementedHandler().ServeHTTP(w, r)
		return
	}

	id := inf.IntParams["id"]
	serverInfo, exists, err := dbhelpers.GetServerInfo(id, tx)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	if !exists {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no server exists by id #%d", id), nil)
		return
	}

	if dsIDs, err := getActiveDeliveryServicesThatOnlyHaveThisServerAssigned(id, serverInfo.Type, tx); err != nil {
		sysErr = fmt.Errorf("checking if server #%d is the last server assigned to any Delivery Services: %v", id, err)
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
		return
	} else if len(dsIDs) > 0 {
		alertText := fmt.Sprintf("deleting server #%d would leave Active Delivery Service", id)
		alertText = InvalidStatusForDeliveryServicesAlertText(alertText, serverInfo.Type, dsIDs)

		api.WriteAlerts(w, r, http.StatusConflict, tc.CreateAlerts(tc.ErrorLevel, alertText))
		return
	}

	var servers []tc.ServerV40
	servers, _, userErr, sysErr, errCode, _ = getServers(r.Header, map[string]string{"id": inf.Params["id"]}, inf.Tx, inf.User, false, *version)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if len(servers) < 1 {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no server exists by id #%d", id), nil)
		return
	}
	if len(servers) > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("there are somehow two servers with id %d - cannot delete", id))
		return
	}
	server := servers[0]
	if server.CDNName != nil {
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDN(inf.Tx.Tx, *server.CDNName, inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	} else if server.CDNID != nil {
		userErr, sysErr, errCode = dbhelpers.CheckIfCurrentUserCanModifyCDNWithID(inf.Tx.Tx, int64(*server.CDNID), inf.User.UserName)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}
	}
	cacheGroupIds := []int{*server.CachegroupID}
	serverIds := []int{*server.ID}
	hasDSOnCDN, err := dbhelpers.CachegroupHasTopologyBasedDeliveryServicesOnCDN(inf.Tx.Tx, *server.CachegroupID, *server.CDNID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	CDNIDs := []int{}
	if hasDSOnCDN {
		CDNIDs = append(CDNIDs, *server.CDNID)
	}
	if err := topology_validation.CheckForEmptyCacheGroups(inf.Tx, cacheGroupIds, CDNIDs, true, serverIds); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("server is the last one in its cachegroup, which is used by a topology: "+err.Error()), nil)
		return
	}

	if result, err := tx.Exec(deleteServerQuery, id); err != nil {
		log.Errorf("Raw error: %v", err)
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("getting rows affected by server delete: %v", err))
		return
	} else if rowsAffected != 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("incorrect number of rows affected: %d", rowsAffected))
		return
	}

	if inf.Version.Major >= 4 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server deleted", server)
	} else {
		csp, err := dbhelpers.GetCommonServerPropertiesFromV4(server, tx)
		if err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		serverv3, err := server.ToServerV3FromV4(csp)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server deleted", serverv3)
	}

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: deleted", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
