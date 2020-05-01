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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"
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
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TOServer combines data about a server with metadata from an API request and
// provides methods that implement several interfaces from the api package.
type TOServer struct {
	api.APIInfoImpl `json:"-"`
	tc.ServerNullable
}

func (s *TOServer) SetLastUpdated(t tc.TimeNoMod)     { s.LastUpdated = &t }
func (*TOServer) InsertQuery() string                 { return insertQuery() }
func (*TOServer) UpdateQuery() string                 { return updateQuery() }
func (*TOServer) DeleteQuery() string                 { return deleteQuery() }

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

func (s *TOServer) Validate() error {
	s.Sanitize()
	version := s.APIInfo().Version
	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")

	errs := []error{}
	if (s.IPAddress == nil || *s.IPAddress == "") && (s.IP6Address == nil || *s.IP6Address == "") {
		errs = append(errs, tc.NeedsAtLeastOneIPError)
	}

	if s.IPIsService != nil && *s.IPIsService && (s.IPAddress == nil || *s.IPAddress == "") {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}

	if s.IP6IsService != nil && *s.IP6IsService && (s.IP6Address == nil || *s.IP6Address == "") {
		errs = append(errs, tc.EmptyAddressCannotBeAServiceAddressError)
	}

	if version.Major >= 2 {
		if (s.IPIsService == nil || !*s.IPIsService) && (s.IP6IsService == nil || !*s.IP6IsService) {
			errs = append(errs, tc.NeedsAtLeastOneServiceAddressError)
		}
	}

	validateErrs := validation.Errors{
		"cachegroupId":   validation.Validate(s.CachegroupID, validation.NotNil),
		"cdnId":          validation.Validate(s.CDNID, validation.NotNil),
		"domainName":     validation.Validate(s.DomainName, validation.NotNil, noSpaces),
		"hostName":       validation.Validate(s.HostName, validation.NotNil, noSpaces),
		"interfaceMtu":   validation.Validate(s.InterfaceMtu, validation.NotNil),
		"interfaceName":  validation.Validate(s.InterfaceName, validation.NotNil),
		"physLocationId": validation.Validate(s.PhysLocationID, validation.NotNil),
		"profileId":      validation.Validate(s.ProfileID, validation.NotNil),
		"statusId":       validation.Validate(s.StatusID, validation.NotNil),
		"typeId":         validation.Validate(s.TypeID, validation.NotNil),
		"updPending":     validation.Validate(s.UpdPending, validation.NotNil),
		"httpsPort":      validation.Validate(s.HTTPSPort, validation.By(tovalidate.IsValidPortNumber)),
		"tcpPort":        validation.Validate(s.TCPPort, validation.By(tovalidate.IsValidPortNumber)),
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
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}

	if _, err := tc.ValidateTypeID(s.ReqInfo.Tx.Tx, s.TypeID, "server"); err != nil {
		return err
	}

	rows, err := s.ReqInfo.Tx.Tx.Query("select cdn from profile where id=$1", s.ProfileID)
	if err != nil {
		log.Error.Printf("could not execute select cdnID from profile: %s\n", err)
		errs = append(errs, tc.DBError)
		return util.JoinErrs(errs)
	}
	defer rows.Close()
	var cdnID int
	for rows.Next() {
		if err := rows.Scan(&cdnID); err != nil {
			log.Error.Printf("could not scan cdnID from profile: %s\n", err)
			errs = append(errs, errors.New("associated profile must have a cdn associated"))
			return util.JoinErrs(errs)
		}
	}
	log.Infof("got cdn id: %d from profile and cdn id: %d from server", cdnID, *s.CDNID)
	if cdnID != *s.CDNID {
		errs = append(errs, errors.New(fmt.Sprintf("CDN id '%d' for profile '%d' does not match Server CDN '%d'", cdnID, *s.ProfileID, *s.CDNID)))
	}
	return util.JoinErrs(errs)
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

func (s *TOServer) Read(h http.Header) ([]interface{}, error, error, int) {
	version := s.APIInfo().Version
	if version == nil {
		return nil, nil, errors.New("TOServer.Read called with nil API version"), http.StatusInternalServerError
	}

	returnable := []interface{}{}

	servers, userErr, sysErr, errCode := getServers(h, s.ReqInfo.Params, s.ReqInfo.Tx, s.ReqInfo.User)

	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode
	}

	for _, server := range servers {
		switch {
		// NOTE: it's required to handle minor version cases in a descending >= manner
		case version.Major >= 2:
			returnable = append(returnable, server)
		case version.Major == 1 && version.Minor >= 1:
			returnable = append(returnable, server.ServerNullableV11)
		default:
			return nil, nil, fmt.Errorf("TOServer.Read called with invalid API version: %d.%d", version.Major, version.Minor), http.StatusInternalServerError
		}
	}

	return returnable, nil, nil, errCode
}

func selectMaxLastUpdatedQuery(where, orderBy, pagination string) string {
	return `SELECT max(t) from (
		SELECT max(s.last_updated) as t from server s JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.tab_name='server') as res`
}

func getServers(h http.Header, params map[string]string, tx *sqlx.Tx, user *auth.CurrentUser) ([]tc.ServerNullable, error, error, int) {
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
			return nil, errors.New("dsId must be an integer"), nil, http.StatusNotFound
		}
		userErr, sysErr, _ := tenant.CheckID(tx.Tx, user, dsID)
		if userErr != nil || sysErr != nil {
			return nil, errors.New("Forbidden"), sysErr, http.StatusForbidden
		}
		// only if dsId is part of params: add join on deliveryservice_server table
		queryAddition = `
FULL OUTER JOIN deliveryservice_server dss ON dss.server = s.id
`
		// depending on ds type, also need to add mids
		dsType, exists, err := dbhelpers.GetDeliveryServiceType(dsID, tx.Tx)
		if err != nil {
			return nil, nil, err, http.StatusInternalServerError
		}
		if !exists {
			return nil, fmt.Errorf("a deliveryservice with id %v was not found", dsID), nil, http.StatusBadRequest
		}
		usesMids = dsType.UsesMidCache()
		log.Debugf("Servers for ds %d; uses mids? %v\n", dsID, usesMids)
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}
	servers := []tc.ServerNullable{}
	runSecond := ims.MakeFirstQuery(tx, h, queryValues, selectMaxLastUpdatedQuery(where, orderBy, pagination))
	if !runSecond {
		return servers, nil, nil, http.StatusNotModified
	}
	// Case where we need to run the second query
	query := selectQuery() + queryAddition + where + orderBy + pagination
	log.Debugln("Query is ", query)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("querying: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	HiddenField := "********"

	for rows.Next() {
		var s tc.ServerNullable
		if err = rows.StructScan(&s); err != nil {
			return nil, nil, errors.New("getting servers: " + err.Error()), http.StatusInternalServerError
		}
		if user.PrivLevel < auth.PrivLevelOperations {
			s.ILOPassword = &HiddenField
			s.XMPPPasswd = &HiddenField
		}
		servers = append(servers, s)
	}

	// if ds requested uses mid-tier caches, add those to the list as well
	if usesMids {
		mids, userErr, sysErr, errCode := getMidServers(servers, tx)

		log.Debugf("getting mids: %v, %v, %s\n", userErr, sysErr, http.StatusText(errCode))

		if userErr != nil || sysErr != nil {
			return nil, userErr, sysErr, errCode
		}
		servers = append(servers, mids...)
	}

	return servers, nil, nil, http.StatusOK
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
	q := selectQuery() + `
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

func (s *TOServer) Update() (error, error, int) {
	if s.IP6Address != nil && len(strings.TrimSpace(*s.IP6Address)) == 0 {
		s.IP6Address = nil
	}

	// see if type changed
	typeID := -1
	if err := s.APIInfo().Tx.QueryRow("SELECT type FROM server WHERE id = $1", s.ID).Scan(&typeID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no server found with this id"), nil, http.StatusNotFound
		}
		return nil, fmt.Errorf("getting current server type: %v", err), http.StatusInternalServerError
	}

	// If type is changing ensure it isn't assigned to any DSes.
	if typeID != *s.TypeID {
		dsIDs := []int64{}
		if err := s.APIInfo().Tx.QueryRowx("SELECT ARRAY(SELECT deliveryservice FROM deliveryservice_server WHERE server = $1)", s.ID).Scan(pq.Array(&dsIDs)); err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("getting server assigned delivery services: %v", err), http.StatusInternalServerError
		}
		if len(dsIDs) != 0 {
			return errors.New("server type can not be updated when it is currently assigned to delivery services"), nil, http.StatusConflict
		}
	}

	current := TOServer{}
	err := s.ReqInfo.Tx.QueryRowx(selectV20UpdatesQuery()+` WHERE sv.id=$1`, strconv.Itoa(*s.ID)).StructScan(&current)
	if err != nil {
		return api.ParseDBError(err)
	}
	defaultIsService := true
	if s.IPIsService == nil {
		if current.IPIsService != nil {
			s.IPIsService = current.IPIsService
		} else {
			s.IPIsService = &defaultIsService
		}
	}
	if s.IP6IsService == nil {
		if current.IP6IsService != nil {
			s.IP6IsService = current.IP6IsService
		} else {
			s.IP6IsService = &defaultIsService
		}
	}

	return api.GenericUpdate(s)
}

func (s *TOServer) Create() (error, error, int) {
	// TODO put in Validate()
	if s.IP6Address != nil && len(strings.TrimSpace(*s.IP6Address)) == 0 {
		s.IP6Address = nil
	}
	if s.XMPPID == nil || *s.XMPPID == "" {
		hostName := *s.HostName
		s.XMPPID = &hostName
	}

	// default the is service field to true if omitted and to upgrade version < 1.4
	defaultIsService := true
	if s.IPIsService == nil {
		s.IPIsService = &defaultIsService
	}
	if s.IP6IsService == nil {
		s.IP6IsService = &defaultIsService
	}

	return api.GenericCreate(s)
}

func (s *TOServer) Delete() (error, error, int) { return api.GenericDelete(s) }

func selectV20UpdatesQuery() string {
	return `SELECT 
sv.ip_address_is_service, 
sv.ip6_address_is_service 
FROM 
	server sv`
}

func selectQuery() string {
	const JumboFrameBPS = 9000
	return `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroup_id,
s.cdn_id,
cdn.name as cdn_name,
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
COALESCE(s.interface_mtu, ` + strconv.Itoa(JumboFrameBPS) + `) as interface_mtu,
s.interface_name,
s.ip6_address,
s.ip6_address_is_service,
s.ip6_gateway,
s.ip_address,
s.ip_address_is_service,
s.ip_gateway,
s.ip_netmask,
s.last_updated,
s.mgmt_ip_address,
s.mgmt_ip_gateway,
s.mgmt_ip_netmask,
s.offline_reason,
pl.name as phys_location,
s.phys_location as phys_location_id,
p.name as profile,
p.description as profile_desc,
s.profile as profile_id,
s.rack,
s.reval_pending,
s.router_host_name,
s.router_port_name,
st.name as status,
s.status as status_id,
s.tcp_port,
t.name as server_type,
s.type as server_type_id,
s.upd_pending as upd_pending,
s.xmpp_id,
s.xmpp_passwd
FROM
  server s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`
}

func insertQuery() string {
	query := `INSERT INTO server (
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
interface_mtu,
interface_name,
ip6_address,
ip6_address_is_service,
ip6_gateway,
ip_address,
ip_address_is_service,
ip_netmask,
ip_gateway,
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
:interface_mtu,
:interface_name,
:ip6_address,
:ip6_address_is_service,
:ip6_gateway,
:ip_address,
:ip_address_is_service,
:ip_netmask,
:ip_gateway,
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
) RETURNING id,last_updated`
	return query
}

func updateQuery() string {
	query := `UPDATE
server SET
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
interface_mtu=:interface_mtu,
interface_name=:interface_name,
ip6_address=:ip6_address,
ip6_address_is_service=:ip6_address_is_service,
ip6_gateway=:ip6_gateway,
ip_address=:ip_address,
ip_address_is_service=:ip_address_is_service,
ip_netmask=:ip_netmask,
ip_gateway=:ip_gateway,
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
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	return `DELETE FROM server WHERE id = :id`
}
