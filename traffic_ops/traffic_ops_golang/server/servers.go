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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing/middleware"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const unfilteredServersQuery = `
SELECT COUNT(server.id)
FROM server
`

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
	p.name AS profile,
	p.description AS profile_desc,
	s.profile AS profile_id,
	s.rack,
	s.reval_pending,
	s.router_host_name,
	s.router_port_name,
	st.name AS status,
	s.status AS status_id,
	s.tcp_port,
	t.name AS server_type,
	s.type AS server_type_id,
	s.upd_pending AS upd_pending,
	s.xmpp_id,
	s.xmpp_passwd
FROM server AS s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id
`

const SelectInterfacesQuery = `
SELECT (
	ARRAY ( SELECT (
		json_build_object (
			'ipAddresses',
			ARRAY (
				SELECT (
					json_build_object (
						'address', ip_address.address,
						'gateway', ip_address.gateway,
						'serviceAddress', ip_address.service_address
					)
				)
				FROM ip_address
				WHERE ip_address.interface = interface.name
				AND ip_address.server = server.id
			),
			'maxBandwidth', interface.max_bandwidth,
			'monitor', interface.monitor,
			'mtu', interface.mtu,
			'name', interface.name
		)
	)
	FROM interface
	WHERE interface.server = server.id
)) AS interfaces,
server.id
FROM server
WHERE server.id = ANY ($1)
`

const insertQuery = `
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
	interface_name,
	mgmt_ip_address,
	mgmt_ip_netmask,
	mgmt_ip_gateway,
	offline_reason,
	phys_location,
	profile,
	rack,
	router_host_name,
	router_port_name,
	status,
	tcp_port,
	type,
	upd_pending,
	xmpp_id,
	xmpp_passwd
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
	:interface_name,
	:mgmt_ip_address,
	:mgmt_ip_netmask,
	:mgmt_ip_gateway,
	:offline_reason,
	:phys_location_id,
	:profile_id,
	:rack,
	:router_host_name,
	:router_port_name,
	:status_id,
	:tcp_port,
	:server_type_id,
	:upd_pending,
	:xmpp_id,
	:xmpp_passwd
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
	profile AS profile_id,
	(SELECT description FROM profile WHERE profile.id=server.profile) AS profile_desc,
	(SELECT name FROM profile WHERE profile.id=server.profile) AS profile,
	rack,
	reval_pending,
	router_host_name,
	router_port_name,
	(SELECT name FROM status WHERE status.id=server.status) AS status,
	status AS status_id,
	tcp_port,
	(SELECT name FROM type WHERE type.id=server.type) AS server_type,
	type AS server_type_id,
	upd_pending
`

const insertInterfaceQuery = `
INSERT INTO interface (
	max_bandwidth,
	monitor,
	mtu,
	name,
	server
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
`

const insertIPQuery = `
INSERT INTO ip_address (
	address,
	gateway,
	interface,
	server,
	service_address
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
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
	profile=:profile_id,
	rack=:rack,
	router_host_name=:router_host_name,
	router_port_name=:router_port_name,
	status=:status_id,
	tcp_port=:tcp_port,
	type=:server_type_id,
	upd_pending=:upd_pending,
	xmpp_id=:xmpp_id,
	xmpp_passwd=:xmpp_passwd
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
	profile AS profile_id,
	(SELECT description FROM profile WHERE profile.id=server.profile) AS profile_desc,
	(SELECT name FROM profile WHERE profile.id=server.profile) AS profile,
	rack,
	reval_pending,
	router_host_name,
	router_port_name,
	(SELECT name FROM status WHERE status.id=server.status) AS status,
	status AS status_id,
	tcp_port,
	(SELECT name FROM type WHERE type.id=server.type) AS server_type,
	type AS server_type_id,
	upd_pending
`

const deleteServerQuery = `DELETE FROM server WHERE id=$1`
const deleteInterfacesQuery = `DELETE FROM interface WHERE server=$1`
const deleteIPsQuery = `DELETE FROM ip_address WHERE server = $1`

func validateCommon(s tc.CommonServerProperties, tx *sql.Tx) []error {

	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	errs := tovalidate.ToErrors(validation.Errors{
		"cachegroupId":   validation.Validate(s.CachegroupID, validation.NotNil),
		"cdnId":          validation.Validate(s.CDNID, validation.NotNil),
		"domainName":     validation.Validate(s.DomainName, validation.NotNil, noSpaces),
		"hostName":       validation.Validate(s.HostName, validation.NotNil, noSpaces),
		"physLocationId": validation.Validate(s.PhysLocationID, validation.NotNil),
		"profileId":      validation.Validate(s.ProfileID, validation.NotNil),
		"statusId":       validation.Validate(s.StatusID, validation.NotNil),
		"typeId":         validation.Validate(s.TypeID, validation.NotNil),
		"updPending":     validation.Validate(s.UpdPending, validation.NotNil),
		"httpsPort":      validation.Validate(s.HTTPSPort, validation.By(tovalidate.IsValidPortNumber)),
		"tcpPort":        validation.Validate(s.TCPPort, validation.By(tovalidate.IsValidPortNumber)),
	})

	if len(errs) > 0 {
		return errs
	}

	if s.XMPPID == nil || *s.XMPPID == "" {
		hostName := *s.HostName
		s.XMPPID = &hostName
	}

	if _, err := tc.ValidateTypeID(tx, s.TypeID, "server"); err != nil {
		errs = append(errs, err)
	}

	var cdnID int
	if err := tx.QueryRow("SELECT cdn from profile WHERE id=$1", s.ProfileID).Scan(&cdnID); err != nil {
		log.Errorf("could not execute select cdnID from profile: %s\n", err)
		if err == sql.ErrNoRows {
			errs = append(errs, errors.New("associated profile must have a cdn associated"))
		} else {
			errs = append(errs, tc.DBError)
		}
		return errs
	}

	log.Infof("got cdn id: %d from profile and cdn id: %d from server", cdnID, *s.CDNID)
	if cdnID != *s.CDNID {
		errs = append(errs, fmt.Errorf("CDN id '%d' for profile '%d' does not match Server CDN '%d'", cdnID, *s.ProfileID, *s.CDNID))
	}

	return errs
}

func validateV1(s tc.ServerNullableV11, tx *sql.Tx) error {
	if s.IP6Address != nil && len(strings.TrimSpace(*s.IP6Address)) == 0 {
		s.IP6Address = nil
	}


	errs := []error{}
	if (s.IPAddress == nil || *s.IPAddress == "") && s.IP6Address == nil {
		errs = append(errs, tc.NeedsAtLeastOneIPError)
	}

	validateErrs := validation.Errors{
		"interfaceMtu":   validation.Validate(s.InterfaceMtu, validation.NotNil),
		"interfaceName":  validation.Validate(s.InterfaceName, validation.NotNil),
	}

	if s.IPAddress != nil && *s.IPAddress != "" {
		validateErrs["ipAddress"] = validation.Validate(s.IPAddress, is.IPv4)
		validateErrs["ipNetmask"] = validation.Validate(s.IPNetmask, validation.NotNil)
		validateErrs["ipGateway"] = validation.Validate(s.IPGateway, validation.NotNil)
	}
	if s.IP6Address != nil && *s.IP6Address != "" {
		validateErrs["ip6Address"] = validation.Validate(s.IP6Address, validation.By(tovalidate.IsValidIPv6CIDROrAddress))
	}
	errs = append(errs, tovalidate.ToErrors(validateErrs)...)
	errs = append(errs, validateCommon(s.CommonServerProperties, tx)...)

	return util.JoinErrs(errs)
}

func validateV2(s *tc.ServerNullableV2, tx *sql.Tx) error {
	var errs []error

	if err := validateV1(s.ServerNullableV11, tx); err != nil {
		return err
	}

	// default boolean value is false
	if s.IPIsService == nil {
		s.IPIsService = new(bool)
	}
	if s.IP6IsService == nil {
		s.IP6IsService = new(bool)
	}

	if !*s.IPIsService && !*s.IP6IsService {
		errs = append(errs, tc.NeedsAtLeastOneServiceAddressError)
	}

	if *s.IPIsService && s.IPAddress == nil {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}

	if *s.IP6IsService && s.IP6Address == nil {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}
	return util.JoinErrs(errs)
}

func validateMTU(mtu interface{}) error {
	m := mtu.(*uint64)
	if m == nil {
		return nil
	}

	if *m < 1280 {
		return errors.New("must be at least 1280")
	}
	return nil
}

func validateGateway(g interface{}) error {
	if g == nil {
		return nil
	}

	if gtwy := net.ParseIP(*g.(*string)); gtwy == nil {
		return errors.New("gateway not a valid IP address")
	}
	return nil
}

func validateV3(s tc.ServerNullable, tx *sql.Tx) (string, error) {

	if len(s.Interfaces) == 0 {
		return "", errors.New("a server must have at least one interface")
	}
	var errs []error
	var serviceAddrV4Found bool
	var serviceAddrV6Found bool
	var serviceInterface string
	for _, iface := range s.Interfaces {

		ruleName := fmt.Sprintf("interface '%s' ", iface.Name)
		errs = append(errs, tovalidate.ToErrors(validation.Errors{
			ruleName + "name": validation.Validate(iface.Name, validation.Required),
			ruleName + "mtu": validation.Validate(iface.MaxBandwidth, validation.By(validateMTU)),
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
					errs = append(errs, errors.New(ruleName + ": address family mismatch between address and gateway"))
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
				} else {
					if serviceAddrV6Found {
						errs = append(errs, fmt.Errorf("interfaces: address '%s' of interface '%s' is marked as a service address, but an IPv6 service address appears earlier in the list", addr.Address, iface.Name))
					}
					serviceAddrV6Found = true
				}
			}
		}
	}

	errs = append(errs, validateCommon(s.CommonServerProperties, tx)...)
	return serviceInterface, util.JoinErrs(errs)
}

func Read(w http.ResponseWriter, r *http.Request) {
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

	servers := []tc.ServerNullable{}
	var unfiltered uint64
	servers, unfiltered, userErr, sysErr, errCode = getServers(inf.Params, inf.Tx, inf.User)

	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if version.Major >= 3 {
		api.WriteRespWithSummary(w, r, servers, unfiltered)
		return
	}

	if version.Major <= 1 {
		legacyServers := make([]tc.ServerNullableV11, 0, len(servers))
		for _, server := range servers {
			legacyServer, err := server.ToServerV2()
			if err != nil {
				api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Failed to convert servers to legacy format: %v", err))
				return
			}
			legacyServers = append(legacyServers, legacyServer.ServerNullableV11)
		}
		api.WriteResp(w, r, legacyServers)
		return
	}

	legacyServers := make([]tc.ServerNullableV2, 0, len(servers))
	for _, server := range servers {
		legacyServer, err := server.ToServerV2()
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Failed to convert servers to legacy format: %v", err))
			return
		}
		legacyServers = append(legacyServers, legacyServer)
	}
	api.WriteResp(w, r, legacyServers)
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	alternative := "GET /servers with query parameter id"
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleDeprecatedErr(w, r, tx, errCode, userErr, sysErr, &alternative)
		return
	}
	defer inf.Close()

	servers := []tc.ServerNullable{}
	servers, _, userErr, sysErr, errCode = getServers(inf.Params, inf.Tx, inf.User)

	if len(servers) > 1 {
		api.HandleDeprecatedErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("ID '%d' matched more than one server (%d total)", inf.IntParams["id"], len(servers)), &alternative)
		return
	}

	deprecationAlerts := api.CreateDeprecationAlerts(&alternative)

	// No need to bother converting if there's no data
	if len(servers) < 1 {
		api.WriteAlertsObj(w, r, http.StatusOK, deprecationAlerts, servers)
	}

	legacyServers := make([]tc.ServerNullableV11, 0, len(servers))
	for _, server := range servers {
		legacyServer, err := server.ToServerV2()
		if err != nil {
			api.HandleDeprecatedErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Failed to convert servers to legacy format: %v", err), &alternative)
			return
		}
		legacyServers = append(legacyServers, legacyServer.ServerNullableV11)
	}
	api.WriteAlertsObj(w, r, http.StatusOK, deprecationAlerts, legacyServers)
}

func getServers(params map[string]string, tx *sqlx.Tx, user *auth.CurrentUser) ([]tc.ServerNullable, uint64, error, error, int) {
	var unfiltered uint64
	if err := tx.QueryRow(unfilteredServersQuery).Scan(&unfiltered); err != nil {
		return nil, 0, nil, fmt.Errorf("Failed to get servers count: %v", err), http.StatusInternalServerError
	}

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"cachegroup":       dbhelpers.WhereColumnInfo{"s.cachegroup", api.IsInt},
		"parentCachegroup": dbhelpers.WhereColumnInfo{"cg.parent_cachegroup_id", api.IsInt},
		"cdn":              dbhelpers.WhereColumnInfo{"s.cdn_id", api.IsInt},
		"id":               dbhelpers.WhereColumnInfo{"s.id", api.IsInt},
		"hostName":         dbhelpers.WhereColumnInfo{"s.host_name", nil},
		"physLocation":     dbhelpers.WhereColumnInfo{"s.phys_location", api.IsInt},
		"profileId":        dbhelpers.WhereColumnInfo{"s.profile", api.IsInt},
		"status":           dbhelpers.WhereColumnInfo{"st.name", nil},
		"type":             dbhelpers.WhereColumnInfo{"t.name", nil},
		"dsId":             dbhelpers.WhereColumnInfo{"dss.deliveryservice", nil},
	}

	usesMids := false
	queryAddition := ""
	if dsIDStr, ok := params[`dsId`]; ok {
		// don't allow query on ds outside user's tenant
		dsID, err := strconv.Atoi(dsIDStr)
		if err != nil {
			return nil, unfiltered, errors.New("dsId must be an integer"), nil, http.StatusNotFound
		}
		userErr, sysErr, _ := tenant.CheckID(tx.Tx, user, dsID)
		if userErr != nil || sysErr != nil {
			return nil, unfiltered, errors.New("Forbidden"), sysErr, http.StatusForbidden
		}
		// only if dsId is part of params: add join on deliveryservice_server table
		queryAddition = "\nFULL OUTER JOIN deliveryservice_server dss ON dss.server = s.id\n"
		// depending on ds type, also need to add mids
		dsType, exists, err := dbhelpers.GetDeliveryServiceType(dsID, tx.Tx)
		if err != nil {
			return nil, unfiltered, nil, err, http.StatusInternalServerError
		}
		if !exists {
			return nil, unfiltered, fmt.Errorf("a deliveryservice with id %v was not found", dsID), nil, http.StatusBadRequest
		}
		usesMids = dsType.UsesMidCache()
		log.Debugf("Servers for ds %d; uses mids? %v\n", dsID, usesMids)
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, unfiltered, util.JoinErrs(errs), nil, http.StatusBadRequest
	}

	query := selectQuery + queryAddition + where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, unfiltered, nil, errors.New("querying: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()


	HiddenField := "********"


	servers := make(map[int]tc.ServerNullable)
	ids := []int{}
	for rows.Next() {
		var s tc.ServerNullable
		if err = rows.StructScan(&s); err != nil {
			return nil, unfiltered, nil, errors.New("getting servers: " + err.Error()), http.StatusInternalServerError
		}
		if user.PrivLevel < auth.PrivLevelOperations {
			s.ILOPassword = &HiddenField
			s.XMPPPasswd = &HiddenField
		}

		if s.ID == nil {
			return nil, unfiltered, nil, errors.New("found server with nil ID"), http.StatusInternalServerError
		}
		if _, ok := servers[*s.ID]; ok {
			return nil, unfiltered, nil, fmt.Errorf("found more than one server with ID #%d", *s.ID), http.StatusInternalServerError
		}
		servers[*s.ID] = s
		ids = append(ids, *s.ID)
	}

	interfaceRows, err := tx.Tx.Query(SelectInterfacesQuery, pq.Array(ids))
	if err != nil {
		return nil, unfiltered, nil, fmt.Errorf("querying for interfaces: %v", err), http.StatusInternalServerError
	}
	defer interfaceRows.Close()

	for interfaceRows.Next() {
		ifaces := []tc.ServerInterfaceInfo{}
		var id int
		if err = interfaceRows.Scan(pq.Array(&ifaces), &id); err != nil {
			return nil, unfiltered, nil, fmt.Errorf("getting server interfaces: %v", err), http.StatusInternalServerError
		}

		if s, ok := servers[id]; !ok {
			log.Warnf("interfaces query returned interfaces for server #%d that was not in original query", id)
		} else {
			s.Interfaces = ifaces
			servers[id] = s
		}
	}

	returnable := make([]tc.ServerNullable, 0, len(servers))
	for _, server := range servers {
		returnable = append(returnable, server)
	}

	// if ds requested uses mid-tier caches, add those to the list as well
	if usesMids {
		mids, userErr, sysErr, errCode := getMidServers(returnable, tx)

		log.Debugf("getting mids: %v, %v, %s\n", userErr, sysErr, http.StatusText(errCode))

		if userErr != nil || sysErr != nil {
			return nil, unfiltered, userErr, sysErr, errCode
		}
		returnable = append(returnable, mids...)
	}

	return returnable, unfiltered, nil, nil, http.StatusOK
}

// getMidServers gets the mids used by the servers in this DS.
//
// Original comment from the Perl code:
//
// If the delivery service employs mids, we're gonna pull mid servers too by
// pulling the cachegroups of the edges and finding those cachegroups parent
// cachegroup... then we see which servers have cachegroup in parent cachegroup
// list...that's how we find mids for the ds :)
func getMidServers(servers []tc.ServerNullable, tx *sqlx.Tx) ([]tc.ServerNullable, error, error, int) {
	if len(servers) == 0 {
		return nil, nil, nil, http.StatusOK
	}
	var ids []string
	for _, s := range servers {
		ids = append(ids, strconv.Itoa(*s.ID))
	}

	edgeIDs := strings.Join(ids, ",")
	// TODO: include secondary parent?
	q := selectQuery + `
	WHERE t.name = 'MID' AND s.cachegroup IN (
	SELECT cg.parent_cachegroup_id FROM cachegroup AS cg
	WHERE cg.id IN (
	SELECT s.cachegroup FROM server AS s
	WHERE s.id IN (` + edgeIDs + `)))
	`
	rows, err := tx.Queryx(q)
	if err != nil {
		return nil, err, nil, http.StatusBadRequest
	}
	defer rows.Close()

	var mids []tc.ServerNullable
	for rows.Next() {
		var s tc.ServerNullable
		if err := rows.StructScan(&s); err != nil {
			log.Error.Printf("could not scan mid servers: %s\n", err)
			return nil, nil, err, http.StatusInternalServerError
		}
		mids = append(mids, s)
	}
	return mids, nil, nil, http.StatusOK
}

func checkTypeChangeSafety(server tc.CommonServerProperties, tx *sqlx.Tx) (error, error, int) {
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

func createInterfaces(id int, interfaces []tc.ServerInterfaceInfo, tx *sql.Tx) (error, error, int) {
	ifaceQry := `
	INSERT INTO interface (
		max_bandwidth,
		monitor,
		mtu,
		name,
		server
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
		argStart := i * 5
		ifaceQueryParts = append(ifaceQueryParts, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", argStart+1, argStart+2, argStart+3, argStart+4, argStart+5))
		ifaceArgs = append(ifaceArgs, iface.MaxBandwidth, iface.Monitor, iface.MTU, iface.Name, id)
		for _, ip := range iface.IPAddresses {
			argStart = len(ipArgs)
			ipQueryParts = append(ipQueryParts, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",  argStart+1, argStart+2, argStart+3, argStart+4, argStart+5))
			ipArgs = append(ipArgs, ip.Address, ip.Gateway, iface.Name, id, ip.ServiceAddress)
		}
	}

	ifaceQry += strings.Join(ifaceQueryParts, ",")
	log.Debugf("Inserting interfaces for new server, query is: %s", ifaceQry)

	ifaceRows, err := tx.Query(ifaceQry, ifaceArgs...)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer ifaceRows.Close()
	insertedIfaces := 0
	for ifaceRows.Next() {
		insertedIfaces++
	}
	log.Debugf("Inserted %d interfaces", insertedIfaces)

	ipQry += strings.Join(ipQueryParts, ",")
	log.Debugf("Inserting IP addresses for new server, query is: %s", ipQry)

	ipRows, err := tx.Query(ipQry, ipArgs...)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer ipRows.Close()

	return nil, nil, http.StatusOK
}

func deleteInterfaces(id int, tx *sql.Tx) (error, error, int) {
	if err := tx.QueryRow(deleteIPsQuery, id).Scan(); err != nil && err != sql.ErrNoRows {
		return api.ParseDBError(err)
	}

	if err := tx.QueryRow(deleteInterfacesQuery, id).Scan(); err != nil && err != sql.ErrNoRows {
		return api.ParseDBError(err)
	}

	return nil, nil, http.StatusOK
}

func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()


	var server tc.ServerNullableV2
	var interfaces []tc.ServerInterfaceInfo
	if inf.Version.Major >= 3 {
		var newServer tc.ServerNullable
		if err := json.NewDecoder(r.Body).Decode(&newServer); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
		serviceInterface, err := validateV3(newServer, tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}

		server, err = newServer.ToServerV2()
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Converting v3 server to v2 for update: %v", err))
			return
		}
		server.InterfaceName = util.StrPtr(serviceInterface)
		interfaces = newServer.Interfaces
	} else if inf.Version.Major == 2 {
		if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}

		err := validateV2(&server, tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}

		interfaces, err = server.LegacyInterfaceDetails.ToInterfaces(*server.IPIsService, *server.IP6IsService)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Converting server legacy interfaces to interface array: %v", err))
			return
		}
	} else {
		var legacyServer tc.ServerNullableV11
		if err := json.NewDecoder(r.Body).Decode(&legacyServer); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}

		err := validateV1(legacyServer, tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}

		interfaces, err = legacyServer.LegacyInterfaceDetails.ToInterfaces(true, true)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Converting server legacy interfaces to interface array: %v", err))
			return
		}
		server = tc.ServerNullableV2{
			ServerNullableV11: legacyServer,
			IPIsService: util.BoolPtr(true),
			IP6IsService: util.BoolPtr(true),
		}
	}

	server.ID = new(int)
	*server.ID = inf.IntParams["id"]

	if userErr, sysErr, errCode = checkTypeChangeSafety(server.CommonServerProperties, inf.Tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	rows, err := inf.Tx.NamedQuery(updateQuery, server)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer rows.Close()

	rowsAffected := 0
	for rows.Next() {
		if err := rows.StructScan(&server); err != nil {
			api.HandleErr(w, r, tx, http.StatusNotFound, nil, fmt.Errorf("scanning lastUpdated from server insert: %v", err))
			return
		}
		rowsAffected++
	}

	if rowsAffected < 1 {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no server found with this id"), nil)
		return
	}
	if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("update for server #%d affected too many rows (%d)", *server.ID, rowsAffected))
		return
	}

	if userErr, sysErr, errCode = deleteInterfaces(inf.IntParams["id"], tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if userErr, sysErr, errCode = createInterfaces(inf.IntParams["id"], interfaces, tx); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if inf.Version.Major >= 3 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", tc.ServerNullable{CommonServerProperties: server.CommonServerProperties, Interfaces: interfaces})
	} else if inf.Version.Minor <= 1 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", server.ServerNullableV11)
	} else {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", server)
	}

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: updated", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func createV1(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerNullableV11

	tx := inf.Tx.Tx

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err := validateV1(server, tx); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	resultRows, err := inf.Tx.NamedQuery(insertQuery, server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&server.CommonServerProperties); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("server create scanning: %v", err))
			return
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("server create: no server was inserted, no id was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many ids returned from server insert"))
	}

	ifaces, err := server.LegacyInterfaceDetails.ToInterfaces(true, true)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
	}

	if userErr, sysErr, errCode := createInterfaces(*server.ID, ifaces, tx); err != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: created", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func createV2(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerNullableV2

	tx := inf.Tx.Tx

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}


	if err := validateV2(&server, tx); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	resultRows, err := inf.Tx.NamedQuery(insertQuery, server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&server.CommonServerProperties); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("server create scanning: %v", err))
			return
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("server create: no server was inserted, no id was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many ids returned from server insert"))
	}

	ifaces, err := server.LegacyInterfaceDetails.ToInterfaces(*server.IPIsService, *server.IP6IsService)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
	}

	if userErr, sysErr, errCode := createInterfaces(*server.ID, ifaces, tx); err != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: created", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func createV3(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerNullable

	tx := inf.Tx.Tx

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	serviceInterface, err := validateV3(server, tx)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	v2Server, err := server.ToServerV2()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}

	v2Server.InterfaceName = &serviceInterface

	resultRows, err := inf.Tx.NamedQuery(insertQuery, v2Server)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&server.CommonServerProperties); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("server create scanning: %v", err))
			return
		}
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("server create: no server was inserted, no id was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many ids returned from server insert"))
		return
	}

	userErr, sysErr, errCode := createInterfaces(*server.ID, server.Interfaces, tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)

	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: created", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	switch {
	case inf.Version.Major <= 1:
		createV1(inf, w, r)
	case inf.Version.Major == 2:
		createV2(inf, w, r)
	default:
		createV3(inf, w, r)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	var servers []tc.ServerNullable
	servers, _, userErr, sysErr, errCode = getServers(map[string]string{"id": inf.Params["id"]}, inf.Tx, inf.User)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if len(servers) < 1 {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("No server exists by id #%d", id), nil)
		return
	}
	if len(servers) > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("There are somehow two servers with id %d - cannot delete", id))
		return
	}

	userErr, sysErr, errCode = deleteInterfaces(id, tx)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	if err := tx.QueryRow(deleteServerQuery, id).Scan(); err != nil && err != sql.ErrNoRows {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	server := servers[0]

	if inf.Version.Major >= 3 {
		api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server deleted", server)
	} else {

		serverV2, err := server.ToServerV2()
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}

		if inf.Version.Major <= 1 {
			api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server deleted", serverV2.ServerNullableV11)
		} else {
			api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server deleted", serverV2)
		}
	}
	changeLogMsg := fmt.Sprintf("SERVER: %s.%s, ID: %d, ACTION: deleted", *server.HostName, *server.DomainName, *server.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
