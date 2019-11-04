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
	"errors"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/lib/pq"
)

func GetServersEligible(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsTenantID, ok, err := getDSTenantIDByID(inf.Tx.Tx, inf.IntParams["id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("delivery service "+inf.Params["id"]+" not found"), nil)
		return
	}
	if authorized, err := tenant.IsResourceAuthorizedToUserTx(*dsTenantID, inf.User, inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
		return
	} else if !authorized {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	servers, err := getEligibleServers(inf.Tx.Tx, inf.IntParams["id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting eligible servers: "+err.Error()))
		return
	}
	api.WriteResp(w, r, servers)
}

const JumboFrameBPS = 9000

func getEligibleServers(tx *sql.Tx, dsID int) ([]tc.DSServer, error) {
	q := `
WITH ds_id as (SELECT $1::bigint as v)
SELECT
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
s.ip6_gateway,
s.ip_address,
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
s.router_host_name,
s.router_port_name,
st.name as status,
s.status as status_id,
s.tcp_port,
t.name as server_type,
s.type as server_type_id,
s.upd_pending as upd_pending,
ARRAY(select ssc.server_capability from server_server_capability ssc where ssc.server = s.id order by ssc.server_capability) as server_capabilities,
ARRAY(select drc.required_capability from deliveryservices_required_capability drc where drc.deliveryservice_id = (select v from ds_id) order by drc.required_capability) as deliveryservice_capabilities
FROM server s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id
WHERE s.cdn_id = (SELECT cdn_id from deliveryservice where id = (select v from ds_id))
AND (t.name LIKE 'EDGE%' OR t.name LIKE 'ORG%')
`
	rows, err := tx.Query(q, dsID)
	if err != nil {
		return nil, errors.New("querying delivery service eligible servers: " + err.Error())
	}
	defer rows.Close()

	servers := []tc.DSServer{}
	for rows.Next() {
		s := tc.DSServer{}
		err := rows.Scan(
			&s.Cachegroup,
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
			&s.InterfaceMtu,
			&s.InterfaceName,
			&s.IP6Address,
			&s.IP6Gateway,
			&s.IPAddress,
			&s.IPGateway,
			&s.IPNetmask,
			&s.LastUpdated,
			&s.MgmtIPAddress,
			&s.MgmtIPGateway,
			&s.MgmtIPNetmask,
			&s.OfflineReason,
			&s.PhysLocation,
			&s.PhysLocationID,
			&s.Profile,
			&s.ProfileDesc,
			&s.ProfileID,
			&s.Rack,
			&s.RouterHostName,
			&s.RouterPortName,
			&s.Status,
			&s.StatusID,
			&s.TCPPort,
			&s.Type,
			&s.TypeID,
			&s.UpdPending,
			pq.Array(&s.ServerCapabilities),
			pq.Array(&s.DeliveryServiceCapabilities),
		)
		if err != nil {
			return nil, errors.New("scanning delivery service eligible servers: " + err.Error())
		}
		eligible := true

		for _, dsc := range s.DeliveryServiceCapabilities {
			if !util.ContainsStr(s.ServerCapabilities, dsc) {
				eligible = false
			}
		}
		if eligible {
			servers = append(servers, s)
		}
	}
	return servers, nil
}
