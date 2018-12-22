package crconfig

import (
	"database/sql"
	"errors"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// UpdateSnapshotTables inserts the latest "real" rows into the snapshot tables.
// It should be possible to simply remove this, if and when the "real" tables are changed to be time-series and always-insert-never-update-never-delete.
func UpdateSnapshotTables(tx *sql.Tx) error {
	// TODO benchmark, see how long it takes
	for _, t := range SnapshotTables {
		if err := t.UpdateSnapshot(tx); err != nil {
			return err
		}
	}
	return nil
}

var SnapshotTables = []SnapshotTable{
	CacheGroupTable,
	CacheGroupFallBacksTable,
	CacheGroupLocalizationMethodTable,
	CDNTable,
	CoordinateTable,
	DeliveryServiceTable,
	DeliveryServiceRegexTable,
	ParameterTable,
	ProfileTable,
	ProfileParameterTable,
	RegexTable,
	ServerTable,
	StaticDNSEntryTable,
	StatusTable,
	TypeTable,
	DeliveryServiceServerTable,
}

var CacheGroupTable = SnapshotTable{
	Table: `cachegroup`,
	PK:    `id`,
	Columns: `
id,
name,
short_name,
parent_cachegroup_id,
secondary_parent_cachegroup_id,
type,
fallback_to_closest,
coordinate
`,
}

var CacheGroupFallBacksTable = SnapshotTable{
	Table: `cachegroup_fallbacks`,
	PK:    `primary_cg, backup_cg`,
	Columns: `
primary_cg,
backup_cg,
set_order
`,
}

// var CacheGroupFallBacksTable = SnapshotTable{
// 	Table: `cachegroup_fallbacks`,
// 	PK:    `primary_cg, backup_cg`,
// 	Columns: `
// primary_cg,
// backup_cg,
// set_order
// `,
// }

var CacheGroupLocalizationMethodTable = SnapshotTable{
	Table: `cachegroup_localization_method`,
	PK:    `cachegroup, method`,
	Columns: `
cachegroup,
method
`,
}

var CDNTable = SnapshotTable{
	Table: `cdn`,
	PK:    `name`,
	Columns: `
id,
name,
dnssec_enabled,
domain_name
`,
}

var CoordinateTable = SnapshotTable{
	Table: `coordinate`,
	PK:    `id`,
	Columns: `
id,
name,
latitude,
longitude
`,
}

var DeliveryServiceTable = SnapshotTable{
	Table: `deliveryservice`,
	PK:    `xml_id`,
	Columns: `
id,
xml_id,
active,
dscp,
signing_algorithm,
qstring_ignore,
geo_limit,
http_bypass_fqdn,
dns_bypass_ip,
dns_bypass_ip6,
dns_bypass_ttl,
type,
profile,
cdn_id,
ccr_dns_ttl,
global_max_mbps,
global_max_tps,
long_desc,
long_desc_1,
long_desc_2,
max_dns_answers,
info_url,
miss_lat,
miss_long,
check_path,
protocol,
ssl_key_version,
ipv6_routing_enabled,
range_request_handling,
edge_header_rewrite,
origin_shield,
mid_header_rewrite,
regex_remap,
cacheurl,
remap_text,
multi_site_origin,
display_name,
tr_response_headers,
initial_dispersion,
dns_bypass_cname,
tr_request_headers,
regional_geo_blocking,
geo_provider,
geo_limit_countries,
logs_enabled,
multi_site_origin_algorithm,
geolimit_redirect_url,
tenant_id,
routing_name,
deep_caching_type,
fq_pacing_rate,
anonymous_blocking_enabled
`,
}

var DeliveryServiceRegexTable = SnapshotTable{
	Table: `deliveryservice_regex`,
	PK:    `deliveryservice, regex`,
	Columns: `
deliveryservice,
regex,
set_number
`,
}

var DeliveryServiceServerTable = SnapshotTable{
	Table: `deliveryservice_server`,
	PK:    `deliveryservice, server`,
	Columns: `
deliveryservice,
server
`,
}

var ParameterTable = SnapshotTable{
	Table: `parameter`,
	PK:    `name, config_file, value`,
	Columns: `
id,
name,
config_file,
value,
secure
`,
}

var ProfileTable = SnapshotTable{
	Table: `profile`,
	PK:    `name`,
	Columns: `
id,
name,
description,
type,
cdn,
routing_disabled
`,
}

var ProfileParameterTable = SnapshotTable{
	Table: `profile_parameter`,
	PK:    `profile, parameter`,
	Columns: `
profile,
parameter
`,
}

var RegexTable = SnapshotTable{
	Table: `regex`,
	PK:    `id`,
	Columns: `
id,
pattern,
type
`,
}

var ServerTable = SnapshotTable{
	Table: `server`,
	PK:    `ip_address, profile`,
	Columns: `
id,
host_name,
domain_name,
tcp_port,
xmpp_id,
xmpp_passwd,
interface_name,
ip_address,
ip_netmask,
ip_gateway,
ip6_address,
ip6_gateway,
interface_mtu,
phys_location,
rack,
cachegroup,
type,
status,
offline_reason,
upd_pending,
profile,
cdn_id,
mgmt_ip_address,
mgmt_ip_netmask,
mgmt_ip_gateway,
ilo_ip_address,
ilo_ip_netmask,
ilo_ip_gateway,
ilo_username,
ilo_password,
router_host_name,
router_port_name,
guid,
https_port,
reval_pending
`,
}

var StaticDNSEntryTable = SnapshotTable{
	Table: `staticdnsentry`,
	PK:    `host, address, deliveryservice, cachegroup`,
	Columns: `
id,
host,
address,
type,
ttl,
deliveryservice,
cachegroup
`,
}

var StatusTable = SnapshotTable{
	Table: `status`,
	PK:    `name`,
	Columns: `
id,
name,
description
`,
}

var TypeTable = SnapshotTable{
	Table: `type`,
	PK:    `name`,
	Columns: `
id,
name,
description,
use_in_table
`,
}

// SnapshotTable contains the data necessary for building common queries for snapshot tables.
type SnapshotTable struct {
	// Table is the name of the base snapshot table
	Table string
	// PK is the primary key for both the base and snapshot table, not including last_updated.
	PK string
	// Columns is a list of the table columns, minus last_updated and deleted (which are on all tables)
	Columns string
}

func (t *SnapshotTable) UpdateSnapshot(tx *sql.Tx) error {
	log.Errorln("DEBUG UpdateSnapshot " + t.Table + " inserting new")

	//debug
	// return errors.New("inserting " + t.Table + " into snapshot (query: QQ" + t.InsertSnapshotQuery() + "QQ): DEBUG")

	if _, err := tx.Exec(t.InsertSnapshotQuery()); err != nil {
		return errors.New("inserting " + t.Table + " into snapshot (query: QQ" + t.InsertSnapshotQuery() + "QQ): " + err.Error())
	}
	log.Errorln("DEBUG UpdateSnapshot " + t.Table + " inserting deleted")
	if _, err := tx.Exec(t.InsertDeletedSnapshotQuery()); err != nil {
		return errors.New("inserting " + t.Table + " deleted rows into snapshot (query QQ" + t.InsertDeletedSnapshotQuery() + "QQ): " + err.Error())
	}
	return nil
}

const SnapshotTableSuffix = `_snapshot`

func (t *SnapshotTable) SnapshotTable() string { return t.Table + SnapshotTableSuffix }

func (t *SnapshotTable) InsertSnapshotQuery() string {
	// TODO test performance, vs inserting only selected things newer than the values in the snapshot table.
	return `INSERT INTO "` + t.SnapshotTable() + `" SELECT * FROM "` + t.Table + `" ON CONFLICT DO NOTHING`
}

func (t *SnapshotTable) InsertDeletedSnapshotQuery() string {
	return t.WithLatest() + `
INSERT INTO "` + t.SnapshotTable() + `" (
` + t.Columns + `, last_updated, deleted
) SELECT
` + t.Columns + `, now() as last_updated, true as deleted
FROM "` + t.SnapshotTable() + `_latest"
WHERE (` + t.PK + `) NOT IN (SELECT ` + t.PK + ` FROM ` + t.Table + `)
`
}

func (t *SnapshotTable) WithLatest() string {
	return `
WITH ` + t.SnapshotTable() + `_latest AS (
  SELECT DISTINCT ON (` + t.PK + `) *
  FROM "` + t.SnapshotTable() + `"
  ORDER BY ` + t.PK + `, last_updated DESC
)
`
}

// UpdateSnapshotTablesForDS is an optiization of UpdateSnapshotTables for when only a single delivery service is being snapshotted.
// This inserts updated data into snapshot tables for rows affecting the given delivery service.
//
// It does NOT necessarily avoid inserting snapshot data for unrelated delivery services. This is only an optimization, in cases where performance is irrelevant, other data may be included for code simplicity. This shouldn't matter, as all functions which query snapshotted data should only select up to their own snapshot time.
//
func UpdateSnapshotTablesForDS(tx *sql.Tx, ds tc.DeliveryServiceName) error {
	// TODO benchmark, see how long it takes
	for _, t := range SnapshotTables {
		if t.Table == DeliveryServiceServerTable.Table {
			continue // skip DSS and only do this DS's rows, because it's slow.
		}
		if err := t.UpdateSnapshot(tx); err != nil {
			return err
		}
	}

	if err := UpdateDSDSSSnapshot(tx, ds); err != nil {
		return err
	}
	return nil
}

func UpdateDSDSSSnapshot(tx *sql.Tx, ds tc.DeliveryServiceName) error {
	log.Errorln("DEBUG UpdateDSDSSSnapshot inserting")
	if _, err := tx.Exec(InsertDSDSSSnapshotQuery(ds)); err != nil {
		return errors.New("inserting  ds '" + string(ds) + "' dss snapshot (query: QQ" + InsertDSDSSSnapshotQuery(ds) + "QQ): " + err.Error())
	}

	log.Errorln("DEBUG UpdateDSDSSSnapshot inserting deleted")
	if _, err := tx.Exec(InsertDSDSSDeletedSnapshotQuery(ds)); err != nil {
		return errors.New("inserting  ds '" + string(ds) + "' dss deleted rows into snapshot (query: QQ" + InsertDSDSSDeletedSnapshotQuery(ds) + "QQ): " + err.Error())
	}

	log.Errorln("DEBUG UpdateDSDSSSnapshot returning")
	return nil
}

func InsertDSDSSSnapshotQuery(ds tc.DeliveryServiceName) string {
	// TODO put ds id in with statement, to only query once?
	// TODO test performance, vs inserting only selected things newer than the values in the snapshot table.
	return `INSERT INTO "` + DeliveryServiceServerTable.SnapshotTable() + `"
SELECT * FROM "` + DeliveryServiceServerTable.Table + `" dss
WHERE dss.deliveryservice = (select id from deliveryservice where xml_id = '` + string(ds) + `')
ON CONFLICT DO NOTHING`
}

func DSDSSWithLatest(ds tc.DeliveryServiceName) string {
	// TODO put ds id in with statement, to only query once?
	t := DeliveryServiceServerTable
	return `
WITH ` + t.SnapshotTable() + `_latest AS (
  SELECT DISTINCT ON (` + t.PK + `) *
  FROM "` + t.SnapshotTable() + `" dsssnap
  WHERE dsssnap.deliveryservice = (select id from deliveryservice where xml_id = '` + string(ds) + `')
  ORDER BY ` + t.PK + `, last_updated DESC
)
`
}

func InsertDSDSSDeletedSnapshotQuery(ds tc.DeliveryServiceName) string {
	// TODO put ds id in with statement, to only query once?
	t := DeliveryServiceServerTable
	return DSDSSWithLatest(ds) + `
INSERT INTO "` + t.SnapshotTable() + `" (
` + t.Columns + `, last_updated, deleted
) SELECT
` + t.Columns + `, now() as last_updated, true as deleted
FROM "` + t.SnapshotTable() + `_latest" dssl
WHERE dssl.deliveryservice = (select id from deliveryservice where xml_id = '` + string(ds) + `')
  AND (` + t.PK + `) NOT IN (
    SELECT ` + t.PK + ` FROM ` + t.Table + `
  )
`
}
