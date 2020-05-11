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

// TOServer combines data about a server with metadata from an API request and
// provides methods that implement several interfaces from the api package.
type TOServer struct {
	api.APIInfoImpl `json:"-"`
	tc.ServerNullableV2
}

const unfilteredServersQuery = `
SELECT COUNT(server.id)
FROM server
`

const selectQuery = `
SELECT
	cg.name AS cachegroup,
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
	s.offline_reason,
	pl.name as phys_location,
	p.name as profile,
	p.description as profile_desc,
	s.rack,
	s.router_host_name,
	s.router_port_name,
	st.name as status,
	s.tcp_port,
	t.name as server_type,
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

const selectInterfacesQuery = `
SELECT (
	ARRAY ( SELECT (
		json_build_object (
			'ipAddresses',
			ARRAY (
				SELECT (
					json_build_object (
						'address', ip_address.address,
						'gateway', ip_address.gateway,
						'service_address', ip_address.service_address
					)
				)
				FROM ip_address
				WHERE ip_address.interface = interface.name
				AND ip_address.server = server.id
			),
			'max_bandwidth', interface.max_bandwidth,
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
) RETURNING id,last_updated
`

const insertInterfacesQuery = `
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

const insertIPsQuery = `
INSERT INTO ip_address (
	address,
	gateway,
	interface,
	server,
	service_address
) VALUES UNNEST (
	$1
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
RETURNING last_updated
`

func (s *TOServer) SetLastUpdated(t tc.TimeNoMod) { s.LastUpdated = &t }
func (*TOServer) InsertQuery() string             { return insertQuery }
func (*TOServer) UpdateQuery() string             { return updateQuery }
func (*TOServer) DeleteQuery() string             { return deleteQuery() }

func (TOServer) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

func (s TOServer) GetKeys() (map[string]interface{}, bool) {
	if s.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *s.ID}, true
}

func (s *TOServer) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	s.ID = &i
}

func (s *TOServer) GetAuditName() string {
	if s.DomainName != nil {
		return *s.DomainName
	}
	if s.ID != nil {
		return strconv.Itoa(*s.ID)
	}
	return "unknown"
}

func (s *TOServer) GetType() string {
	return "server"
}

func (s *TOServer) Sanitize() {
	if s.IP6Address != nil && *s.IP6Address == "" {
		s.IP6Address = nil
	}
}

func validateCommon(s tc.CommonServerProperties, tx *sql.Tx) []error {
	if s.XMPPID == nil || *s.XMPPID == "" {
		hostName := *s.HostName
		s.XMPPID = &hostName
	}

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

	if _, err := tc.ValidateTypeID(tx, s.TypeID, "server"); err != nil {
		errs = append(errs, err)
	}

	var cdnID int
	if err := tx.QueryRow("SELECT cdn from profile WHERE id=$1", s.ProfileID).Scan(&cdnID); err != nil {
		log.Error.Printf("could not execute select cdnID from profile: %s\n", err)
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

func validateV2(s tc.ServerNullableV2, tx *sql.Tx) error {
	var errs []error

	if err := validateV1(s.ServerNullableV11, tx); err != nil {
		return err
	}

	if (s.IPIsService == nil || !*s.IPIsService) && (s.IP6IsService == nil || !*s.IP6IsService) {
		errs = append(errs, tc.NeedsAtLeastOneServiceAddressError)
	}

	if s.IPIsService != nil && *s.IPIsService && (s.IPAddress == nil) {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}

	if s.IP6IsService != nil && *s.IP6IsService && (s.IP6Address == nil) {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}
	return util.JoinErrs(errs)
}

func validateV3(tc.ServerNullableV2, *sql.Tx) error {
	return nil
}

// ChangeLogMessage implements the api.ChangeLogger interface for a custom log message
func (s TOServer) ChangeLogMessage(action string) (string, error) {

	var status string
	if s.Status != nil {
		status = *s.Status
	}

	var hostName string
	if s.HostName != nil {
		hostName = *s.HostName
	}

	var domainName string
	if s.DomainName != nil {
		domainName = *s.DomainName
	}

	var serverID string
	if s.ID != nil {
		serverID = strconv.Itoa(*s.ID)
	}

	message := action + ` ` + status + ` server: { "hostName":"` + hostName + `", "domainName":"` + domainName + `", id:` + serverID + ` }`

	return message, nil
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

	var servers []tc.ServerNullable
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
		return
	}

	legacyServers := make([]tc.ServerNullableV11, 0, len(servers))
	log.Debugf("servers len=%d", len(servers))
	for _, server := range servers {
		legacyServer, err := server.ToServerV2()
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Failed to convert servers to legacy format: %v", err))
			return
		}
		legacyServers = append(legacyServers, legacyServer.ServerNullableV11)
	}
	log.Debugf("legacyServers len=%d", len(legacyServers))
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

	var servers []tc.ServerNullable
	servers, _, userErr, sysErr, errCode = getServers(inf.Params, inf.Tx, inf.User)

	legacyServers := make([]tc.ServerNullableV11, len(servers))
	for _, server := range servers {
		legacyServer, err := server.ToServerV2()
		if err != nil {
			api.HandleDeprecatedErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("Failed to convert servers to legacy format: %v", err), &alternative)
			return
		}
		legacyServers = append(legacyServers, legacyServer.ServerNullableV11)
	}
	deprecationAlerts := api.CreateDeprecationAlerts(&alternative)
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

	interfaceRows, err := tx.Tx.Query(selectInterfacesQuery, pq.Array(ids))
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
			log.Warnf("interfaces query returned interfaces for server #%d that was not in original query")
		} else {
			s.Interfaces = ifaces
			servers[id] = s
		}
	}

	var returnable []tc.ServerNullable
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

func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

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

	// see if cdn or type changed
	var cdnID int
	var typeID int

	if err := inf.Tx.QueryRow("SELECT type, cdn_id FROM server WHERE id = $1", server.ID).Scan(&typeID, &cdnID); err != nil {
		if err == sql.ErrNoRows {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no server found with this ID"), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("getting current server type: %v", err))
		return
	}

	var dsIDs []int64
	if err := inf.Tx.QueryRowx("SELECT ARRAY(SELECT deliveryservice FROM deliveryservice_server WHERE server = $1)", server.ID).Scan(pq.Array(&dsIDs)); err != nil && err != sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("getting server assigned delivery services: %v", err))
		return
	}
	// Check to see if the user is trying to change the CDN of a server, which is already linked with a DS
	if cdnID != *server.CDNID && len(dsIDs) != 0 {
		api.HandleErr(w, r, tx, http.StatusConflict, errors.New("server cdn can not be updated when it is currently assigned to delivery services"), nil)
		return
	}
	// If type is changing ensure it isn't assigned to any DSes.
	if typeID != *server.TypeID {
		if len(dsIDs) != 0 {
			api.HandleErr(w, r, tx, http.StatusConflict, errors.New("server type can not be updated when it is currently assigned to Delivery Services"), nil)
			return
		}
	}

	// current := TOServer{}
	// err := inf.Tx.QueryRowx(selectV20UpdatesQuery()+` WHERE sv.id=$1`, strconv.Itoa(*s.ID)).StructScan(&current)
	// if err != nil {
	// 	return api.ParseDBError(err)
	// }
	// defaultIsService := true
	// if s.IPIsService == nil {
	// 	if current.IPIsService != nil {
	// 		s.IPIsService = current.IPIsService
	// 	} else {
	// 		s.IPIsService = &defaultIsService
	// 	}
	// }
	// if s.IP6IsService == nil {
	// 	if current.IP6IsService != nil {
	// 		s.IP6IsService = current.IP6IsService
	// 	} else {
	// 		s.IP6IsService = &defaultIsService
	// 	}
	// }

	rows, err := inf.Tx.NamedQuery(updateQuery, server)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no server found with this id"), nil)
	}
	var lastUpdated tc.TimeNoMod
	if err := rows.Scan(&lastUpdated); err != nil {
		api.HandleErr(w, r, tx, http.StatusNotFound, nil, fmt.Errorf("scanning lastUpdated from server insert: %v", err))
		return
	}
	server.LastUpdated = &lastUpdated

	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Server updated", server)
}

func createInterfaces(s tc.ServerNullableV11, tx *sql.Tx) error {
	if err := tx.QueryRow(insertInterfacesQuery, nil, true, s.InterfaceMtu, s.InterfaceName, s.ID).Scan(); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Inserting interface: %v", err)
	}

	var ips []tc.ServerIpAddress
	if s.IPAddress != nil && *s.IPAddress != "" {
		ips = append(ips, tc.ServerIpAddress{*s.IPAddress, s.IPGateway, *s.InterfaceName, uint64(*s.ID), true})
	}

	if s.IP6Address != nil && *s.IP6Address != "" {
		ips = append(ips, tc.ServerIpAddress{*s.IP6Address, s.IP6Gateway, *s.InterfaceName, uint64(*s.ID), true})
	}

	if err := tx.QueryRow(insertIPsQuery, pq.Array(&ips)).Scan(); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Inserting IPs: %v", err)
	}
	return nil
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

	var id int
	var lastUpdated tc.TimeNoMod

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
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
	server.ID = &id
	server.LastUpdated = &lastUpdated

	if err := createInterfaces(server, tx); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)
}

func createV2(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerNullableV2

	tx := inf.Tx.Tx

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err := validateV2(server, tx); err != nil {
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


	var id int
	var lastUpdated tc.TimeNoMod

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
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
	server.ID = &id
	server.LastUpdated = &lastUpdated

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)
}

func createV3(inf *api.APIInfo, w http.ResponseWriter, r *http.Request) {
	var server tc.ServerNullableV2

	tx := inf.Tx.Tx

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err := validateV3(server, tx); err != nil {
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


	var id int
	var lastUpdated tc.TimeNoMod

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
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
	server.ID = &id
	server.LastUpdated = &lastUpdated

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Server created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, server)
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

func (s *TOServer) Delete() (error, error, int) { return api.GenericDelete(s) }

func selectV20UpdatesQuery() string {
	return `SELECT
sv.ip_address_is_service,
sv.ip6_address_is_service
FROM
	server sv`
}



func deleteQuery() string {
	return `DELETE FROM server WHERE id = :id`
}
