package monitoring

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
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/lib/pq"
)

const CacheMonitorConfigFile = "rascal.properties"

const MonitorType = "RASCAL"
const RouterType = "CCR"
const MonitorProfilePrefix = "RASCAL"
const MonitorConfigFile = "rascal-config.txt"
const KilobitsPerMegabit = 1000
const DeliveryServiceStatus = "REPORTED"

type BasicServer struct {
	Profile    string `json:"profile"`
	Status     string `json:"status"`
	IP         string `json:"ip"`
	IP6        string `json:"ip6"`
	Port       int    `json:"port"`
	Cachegroup string `json:"cachegroup"`
	HostName   string `json:"hostname"`
	FQDN       string `json:"fqdn"`
}

type Monitor struct {
	BasicServer
}

type Cache struct {
	BasicServer
	InterfaceName string `json:"interfacename"`
	Type          string `json:"type"`
	HashID        string `json:"hashid"`
}

type Cachegroup struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Profile struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Monitoring struct {
	TrafficServers   []Cache                `json:"trafficServers"`
	TrafficMonitors  []Monitor              `json:"trafficMonitors"`
	Cachegroups      []Cachegroup           `json:"cacheGroups"`
	Profiles         []Profile              `json:"profiles"`
	DeliveryServices []DeliveryService      `json:"deliveryServices"`
	Config           map[string]interface{} `json:"config"`
}

type MonitoringResponse struct {
	Response Monitoring `json:"response"`
}

type Router struct {
	Type    string
	Profile string
}

type DeliveryService struct {
	XMLID              string  `json:"xmlId"`
	TotalTPSThreshold  float64 `json:"totalTpsThreshold"`
	Status             string  `json:"status"`
	TotalKBPSThreshold float64 `json:"totalKbpsThreshold"`
}

func GetMonitoringJSON(tx *sql.Tx, cdnName string, live bool) (*Monitoring, error) {
	mn := &Monitoring{}
	err := error(nil)
	routers := []Router{}
	if mn.TrafficMonitors, mn.TrafficServers, routers, err = getMonitoringServers(tx, cdnName, live); err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}
	if mn.Cachegroups, err = getCachegroups(tx, cdnName, live); err != nil {
		return nil, fmt.Errorf("error getting cachegroups: %v", err)
	}
	if mn.Profiles, err = getProfiles(tx, cdnName, mn.TrafficServers, routers, live); err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}
	if mn.DeliveryServices, err = getDeliveryServices(tx, routers, live); err != nil {
		return nil, fmt.Errorf("error getting deliveryservices: %v", err)
	}
	if mn.Config, err = getConfig(tx, cdnName, live); err != nil {
		return nil, fmt.Errorf("error getting config: %v", err)
	}
	return mn, nil
}

func getMonitoringServers(tx *sql.Tx, cdn string, live bool) ([]Monitor, []Cache, []Router, error) {
	with := ""
	selectedColumns := `host_name, fqdn, status, cachegroup, port, ip, ip6, profile, interface_name, server_type, hash_id`
	primaryKeys := `s.host_name`
	selectBody := `
  s.host_name,
  CONCAT(s.host_name, '.', s.domain_name) fqdn,
  st.name status,
  cg.name cachegroup,
  s.tcp_port port,
  s.ip_address ip,
  s.ip6_address ip6,
  pr.name profile,
  s.interface_name,
  tp.name server_type,
  s.xmpp_id hash_id,
  s.deleted
FROM
  server_snapshot s
  JOIN type_snapshot tp ON tp.id = s.type
  JOIN status_snapshot st ON st.id = s.status
  JOIN cachegroup_snapshot cg ON cg.id = s.cachegroup
  JOIN profile_snapshot pr ON pr.id = s.profile
  JOIN cdn_snapshot c ON c.id = s.cdn_id
`
	where := `c.name = (select v from cdn_name)`
	tableAliases := []string{"s", "tp", "st", "cg", "pr", "c"}
	qry := buildSnapshotQuery(live, with, selectedColumns, primaryKeys, selectBody, where, tableAliases)

	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}
	caches := []Cache{}
	routers := []Router{}

	for rows.Next() {
		var hostName sql.NullString
		var fqdn sql.NullString
		var status sql.NullString
		var cachegroup sql.NullString
		var port sql.NullInt64
		var ip sql.NullString
		var ip6 sql.NullString
		var profile sql.NullString
		var interfaceName sql.NullString
		var ttype sql.NullString
		var hashID sql.NullString

		if err := rows.Scan(&hostName, &fqdn, &status, &cachegroup, &port, &ip, &ip6, &profile, &interfaceName, &ttype, &hashID); err != nil {
			return nil, nil, nil, err
		}

		if ttype.String == tc.MonitorTypeName {
			monitors = append(monitors, Monitor{
				BasicServer: BasicServer{
					Profile:    profile.String,
					Status:     status.String,
					IP:         ip.String,
					IP6:        ip6.String,
					Port:       int(port.Int64),
					Cachegroup: cachegroup.String,
					HostName:   hostName.String,
					FQDN:       fqdn.String,
				},
			})
		} else if strings.HasPrefix(ttype.String, "EDGE") || strings.HasPrefix(ttype.String, "MID") {
			caches = append(caches, Cache{
				BasicServer: BasicServer{
					Profile:    profile.String,
					Status:     status.String,
					IP:         ip.String,
					IP6:        ip6.String,
					Port:       int(port.Int64),
					Cachegroup: cachegroup.String,
					HostName:   hostName.String,
					FQDN:       fqdn.String,
				},
				InterfaceName: interfaceName.String,
				Type:          ttype.String,
				HashID:        hashID.String,
			})
		} else if ttype.String == tc.RouterTypeName {
			routers = append(routers, Router{
				Type:    ttype.String,
				Profile: profile.String,
			})
		}
	}
	return monitors, caches, routers, nil
}

func getCachegroups(tx *sql.Tx, cdn string, live bool) ([]Cachegroup, error) {
	with := ""
	selectedColumns := "name, latitude, longitude"
	primaryKeys := "cg.name"
	selectBody := `
  cg.name,
  co.latitude,
  co.longitude,
  cg.deleted
FROM
  cachegroup_snapshot cg
  LEFT JOIN coordinate_snapshot co ON co.id = cg.coordinate
  JOIN server_snapshot s ON s.cachegroup = cg.id
  JOIN cdn_snapshot c ON c.id = s.cdn_id
`
	where := `c.name = (select v from cdn_name)`
	tableAliases := []string{"cg", "co", "s", "c"}
	qry := buildSnapshotQuery(live, with, selectedColumns, primaryKeys, selectBody, where, tableAliases)

	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cachegroups := []Cachegroup{}
	for rows.Next() {
		var name sql.NullString
		var lat sql.NullFloat64
		var lon sql.NullFloat64
		if err := rows.Scan(&name, &lat, &lon); err != nil {
			return nil, err
		}
		cachegroups = append(cachegroups, Cachegroup{
			Name: name.String,
			Coordinates: Coordinates{
				Latitude:  lat.Float64,
				Longitude: lon.Float64,
			},
		})
	}
	return cachegroups, nil
}

func getProfiles(tx *sql.Tx, cdn string, caches []Cache, routers []Router, live bool) ([]Profile, error) {
	with := ""
	selectedColumns := `profile, name, value`
	primaryKeys := "pr.name, pa.name, pa.value"
	selectBody := `
  pr.name as profile,
  pa.name,
  pa.value,
  pa.deleted
FROM
  parameter_snapshot pa
  JOIN profile_snapshot pr ON pr.name = ANY($2)
  JOIN profile_parameter_snapshot pp ON pp.profile = pr.id and pp.parameter = pa.id
`
	where := `pa.config_file = $3`
	tableAliases := []string{`pa`, `pr`, `pp`}
	qry := buildSnapshotQuery(live, with, selectedColumns, primaryKeys, selectBody, where, tableAliases)

	cacheProfileTypes := map[string]string{}
	profiles := map[string]Profile{}
	profileNames := []string{}
	for _, router := range routers {
		profiles[router.Profile] = Profile{
			Name: router.Profile,
			Type: router.Type,
		}
	}

	for _, cache := range caches {
		if _, ok := cacheProfileTypes[cache.Profile]; !ok {
			cacheProfileTypes[cache.Profile] = cache.Type
			profiles[cache.Profile] = Profile{
				Name: cache.Profile,
				Type: cache.Type,
			}
			profileNames = append(profileNames, cache.Profile)
		}
	}

	rows, err := tx.Query(qry, cdn, pq.Array(profileNames), CacheMonitorConfigFile)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profileName sql.NullString
		var name sql.NullString
		var value sql.NullString
		if err := rows.Scan(&profileName, &name, &value); err != nil {
			return nil, err
		}
		if name.String == "" {
			return nil, fmt.Errorf("null name") // TODO continue and warn?
		}
		profile := profiles[profileName.String]
		if profile.Parameters == nil {
			profile.Parameters = map[string]interface{}{}
		}

		if valNum, err := strconv.Atoi(value.String); err == nil {
			profile.Parameters[name.String] = valNum
		} else {
			profile.Parameters[name.String] = value.String
		}
		profiles[profileName.String] = profile

	}

	profilesArr := []Profile{} // TODO make for efficiency?
	for _, profile := range profiles {
		profilesArr = append(profilesArr, profile)
	}
	return profilesArr, nil
}

func getDeliveryServices(tx *sql.Tx, routers []Router, live bool) ([]DeliveryService, error) {
	profileNames := []string{}
	for _, router := range routers {
		profileNames = append(profileNames, router.Profile)
	}

	qry := `
SELECT
  xml_id,
  global_max_tps,
  global_max_mbps,
  ds_deleted
FROM (
SELECT DISTINCT ON (ds.xml_id)
  ds.xml_id,
  ds.global_max_tps,
  ds.global_max_mbps,
  ds.deleted as ds_deleted
FROM
  deliveryservice_snapshot ds
  JOIN profile_snapshot pr ON pr.id = ds.profile
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = ds.xml_id
WHERE
  pr.name = ANY($1)
  AND ds.active = true
`
	if !live {
		qry += `
  AND ds.last_updated <= dsn.time
  AND pr.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  ds.xml_id DESC,
  ds.last_updated DESC,
  pr.last_updated DESC
) s WHERE ds_deleted = false
`
	rows, err := tx.Query(qry, pq.Array(profileNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dses := []DeliveryService{}

	for rows.Next() {
		var xmlid sql.NullString
		var tps sql.NullFloat64
		var mbps sql.NullFloat64
		if err := rows.Scan(&xmlid, &tps, &mbps); err != nil {
			return nil, err
		}
		dses = append(dses, DeliveryService{
			XMLID:              xmlid.String,
			TotalTPSThreshold:  tps.Float64,
			Status:             DeliveryServiceStatus,
			TotalKBPSThreshold: mbps.Float64 * KilobitsPerMegabit,
		})
	}
	return dses, nil
}

func getConfig(tx *sql.Tx, cdn string, live bool) (map[string]interface{}, error) {
	with := ""
	selectedColumns := `name, value`
	primaryKeys := "pa.name, pa.value"
	// TODO remove 'like' in query? Slow?
	selectBody := `
  pa.name,
  pa.value
FROM
  parameter_snapshot pa
  JOIN profile_snapshot pr ON pr.name LIKE '` + MonitorProfilePrefix + `%%'
  JOIN profile_parameter_snapshot pp ON pp.profile = pr.id and pp.parameter = pa.id
`
	where := `pa.config_file = '` + MonitorConfigFile + `'`
	tableAliases := []string{`pa`, `pr`, `pp`}
	qry := buildSnapshotQuery(live, with, selectedColumns, primaryKeys, selectBody, where, tableAliases)

	rows, err := tx.Query(qry, cdn)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cfg := map[string]interface{}{}

	for rows.Next() {
		var name sql.NullString
		var val sql.NullString
		if err := rows.Scan(&name, &val); err != nil {
			return nil, err
		}
		if valNum, err := strconv.Atoi(val.String); err == nil {
			cfg[name.String] = valNum
		} else {
			cfg[name.String] = val.String
		}
	}
	return cfg, nil
}

// buildSnapshotQuery builds a query to select the latest timestamp, from the query parts of an ordinary query.
//
// The live arg is whether to query the latest timestamp. If false, the latest up to the snapshot is queried.
//
// This is just a helper, there are things it can't do (e.g. order by), and it's fine not to use it, if it doesn't fit.
//
// Note this requires the cdn name (in order to get the snapshot time). This is queried as $1. Hence:
// 1. The cdn name must be the first query parameter
// 2. The cdn name is available to the selectBody via `(select v from cdn_name)`
//
// The with parameter is any with-statement query parts. This should be a complete "WITH" query part. This may be blank.
// Examples:
//    with := `WITH cdn_id AS (select id from cdn where name = 'foo')`
//    with := `
//    WITH one AS (select 1),
//         two AS (select 2)
//    `
//
// The selectColumns must be the names of columns selected by the parameter selectBody, separated by commas. This should not include the deleted column, which will always be false.
// Example:
//    selectColumns := `name, value`
//
// The primaryKeys must be the primary keys of the selected statement, including the table aliases used in the selectBody, separated by commas. Note this is not necessarily the primary key(s) of a single table, but rather the unique values of the select statement itself (in technical terms, the "candidate key"). For example, this may be the delivery service xml_id; but it may also be "profile.id, parameter.id, parameter.name".
//  Examples:
//    primaryKeys := `ds.xml_id`
//    primaryKeys := `pr.name, pa.name, pa.value`
//
// The selectBody must be the select statement, including from and joins, including selecting the deleted column of the primary table, excluding the initial "select" keyword. A deleted column MUST be selected.
// Example:
//      pa.name,
//      pa.value,
//      pa.deleted
//    FROM
//      parameter pa
//    JOIN
//      profile pr ON pr.name LIKE '` + MonitorProfilePrefix + `%%'
//      JOIN profile_parameter pp ON pp.profile = pr.id and pp.parameter = pa.id
//
// The where is the where clause, including the "where" keyword. This may be blank.
// Examples:
//    where := `WHERE pa.config_file = '` + MonitorConfigFile + `'`
//    where := `
//    WHERE pa.config_file = 'CRConfig.json'
//          AND pa.name = 'something'
//    `
//
// The tableAliases is all table aliases (or names) used in the selectBody. For example, if the select body contains "FROM foo JOIN bar b on b.id = foo.bar", then tableAliases must be []string{"foo", "b"}.
// Examples:
//   tableAliases := []string{"pa", "pr"}
//   tableAliases := []string{"s", "tp", "st", "cg", "pr", "cdn"}
//
//
// All the tables selected, in the parameters selectBody and tableAliases, should be snapshot tables like deliveryservice_snapshot, not raw tables like "deliveryservice". The purpose of this function is to build a query to select the latest values, abstracting away the common boilerplate, and allowing callers to pass the query parts mostly unmodified from an ordinary "non-snapshot" query.
//
// To better explain its purpose with an example, it allows you to pass the query parts of a query such as:
//
//    SELECT
//      pa.name,
//      pa.value
//    FROM
//      parameter_snapshot pa
//    JOIN
//      profile_snapshot pr ON pr.name LIKE '` + MonitorProfilePrefix + `%%'
//      JOIN profile_parameter_snapshot pp ON pp.profile = pr.id and pp.parameter = pa.id
//    WHERE
//      pa.config_file = '` + MonitorConfigFile + `'
//
// and construct the "select latest <= snapshot" query such as:
//
//    WITH cdn_name AS (
//      SELECT $1::text as v
//    ),
//    snapshot_time AS (
//      SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
//    )
//    SELECT
//      name,
//      value
//    FROM (
//    SELECT DISTINCT ON (pa.name, pa.value)
//      pa.deleted,
//      pp.deleted,
//      pa.name,
//      pa.value,
//    FROM
//      parameter_snapshot pa
//    JOIN
//      profile_snapshot pr ON pr.name LIKE '` + MonitorProfilePrefix + `%%'
//      JOIN profile_parameter_snapshot pp ON pp.profile = pr.id and pp.parameter = pa.id
//    WHERE
//      pa.config_file = '` + MonitorConfigFile + `'
//
//      AND ds.last_updated <= (select v from snapshot_time)
//      AND pr.last_updated <= (select v from snapshot_time)
//
//    ORDER BY
//      pa.name,
//      pa.value,
//      pr.last_updated DESC,
//      pp.last_updated DESC
//    ) s WHERE
//      pa.deleted = false
//      AND pp.deleted = false
//
// While requiring only minimal changes to the original:
// 1. Query parts must be separated out.
// 2. Tables must be suffixed with '_snapshot', to select from the snapshot tables.
// 3. The primary key of the select statement must be identified.
//
// The usage for the above example would be:
//
//    live := true
//    with := ""
//    selectedColumns := "name, value"
//    primaryKeys := "pa.name, pa.value"
//    selectBody := `
//      pa.name,
//      pa.value
//    FROM
//      parameter_snapshot pa
//    JOIN
//      profile_snapshot pr ON pr.name LIKE '` + MonitorProfilePrefix + `%%'
//      JOIN profile_parameter_snapshot pp ON pp.profile = pr.id and pp.parameter = pa.id
//    `
//    where := `pa.config_file = '` + MonitorConfigFile + `'`
//    tableAliases := []string{"pa", "pp"}
//    qry := buildSnapshotQuery(live, with, selectedColumns, primaryKeys, selectBody, where, tableAliases)
//
func buildSnapshotQuery(
	live bool,
	with string,
	selectedColumns string,
	primaryKeys string,
	selectBody string,
	where string,
	tableAliases []string,
) string {
	qry := with

	if !live {
		if with == `` {
			qry += `WITH `
		} else {
			qry += `, `
		}
		qry += `
cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
`
	}

	if len(tableAliases) < 1 {
		// this function is never useful with no tables, so this should never happen; but I loathe panics.
		return `` // TODO log?
	}
	firstAlias := tableAliases[0]
	restAliases := tableAliases[1:]

	qry += ` SELECT ` + selectedColumns + `
FROM (SELECT DISTINCT ON (` + primaryKeys + `)
`

	for _, alias := range tableAliases {
		qry += alias + `.deleted as ` + alias + `_deleted,
`
	}

	qry += selectBody
	if where != "" {
		qry += "WHERE " + where
	}

	if !live {
		if where == `` {
			qry += `
WHERE `
		} else {
			qry += `
 AND `
		}
		qry += firstAlias + `.last_updated <= (select v from snapshot_time)
`
		for _, alias := range restAliases {
			qry += ` AND ` + alias + `.last_updated <= (select v from snapshot_time)
`
		}
	}

	qry += ` ORDER BY
` + primaryKeys
	for _, alias := range tableAliases {
		qry += `, ` + alias + `.last_updated DESC
`
	}
	qry += ` ) s WHERE ` + firstAlias + `_deleted = false
`
	for _, alias := range tableAliases {
		qry += ` AND ` + alias + `_deleted = false
`
	}

	return qry
}
