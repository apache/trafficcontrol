package crconfig

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// UpdateSnapshotTables inserts the latest "real" rows into the snapshot tables, for all tables used in the CRConfig EXCEPT the "deliveryservices" key and "contentServers/{server}/deliveryservices" keys.
// It should be possible to simply remove this, if and when the "real" tables are changed to be time-series and always-insert-never-update-never-delete.
//
// Note Delivery services must be snapshotted separately, their data is NOT inserted by this function.
//
func UpdateSnapshotTables(tx *sql.Tx) error {
	// TODO add cdn to queries, to not unnecessarily copy all cdns, since snapshots are always per-cdn.
	// TODO benchmark, see how long it takes
	for _, t := range SnapshotTables {
		if err := t.UpdateSnapshot(tx); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSnapshotTablesForDS is an optiization of UpdateSnapshotTables for when only a single delivery service is being snapshotted.
// This inserts updated data into snapshot tables for rows affecting the given delivery service.
//
// It does NOT necessarily avoid inserting snapshot data for unrelated delivery services. This is only an optimization, in cases where performance is irrelevant, other data may be included for code simplicity. This shouldn't matter, as all functions which query snapshotted data should only select up to their own snapshot time.
//

// UpdateSnapshotTables inserts the latest "real" rows into the snapshot tables, for all tables used in the CRConfig EXCEPT the "deliveryservices" key and "contentServers/{server}/deliveryservices" keys,
// PLUS the rows in tables used for the given delivery service.
//
// It should be possible to simply remove this, if and when the "real" tables are changed to be time-series and always-insert-never-update-never-delete.
//
// Note Delivery services must be snapshotted separately, this function ONLY inserts the snapshot rows for the given delivery service.
//
func UpdateSnapshotTablesForDS(tx *sql.Tx, ds tc.DeliveryServiceName) error {
	// TODO benchmark, see how long it takes
	for _, t := range DSSnapshotTables {
		if err := t.UpdateSnapshot(tx); err != nil {
			return err
		}
	}

	if err := UpdateDSDSSSnapshot(tx, ds); err != nil {
		return err
	}
	return nil
}

// SnapshotTables includes all tables used by the CRConfig NOT INCLUDING the "deliveryservices" key OR "contentServers/{server}/deliveryService" keys
var SnapshotTables = []SnapshotTable{
	CacheGroupTable,
	CacheGroupFallBacksTable,
	CacheGroupLocalizationMethodTable,
	CDNTable,
	CoordinateTable,
	ParameterTable,
	ProfileTable,
	ProfileParameterTable,
	ServerTable,
	StatusTable,
	TypeTable,
}

// SnapshotTables includes all tables used by the CRConfig NOT INCLUDING DeliveryServiceServer, which should be populated with special logic to only include a given delivery service being snapshotted, because it's so large.
var DSSnapshotTables = append([]SnapshotTable{
	DeliveryServiceTable,
	DeliveryServiceRegexTable,
	StaticDNSEntryTable,
	RegexTable,
	// DeliveryServiceServerTable, // not included, because it needs special logic to only do the given DS
}, SnapshotTables...)

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
	//debug
	// return errors.New("inserting " + t.Table + " into snapshot (query: QQ" + t.InsertSnapshotQuery() + "QQ): DEBUG")

	if _, err := tx.Exec(t.InsertSnapshotQuery()); err != nil {
		return errors.New("inserting " + t.Table + " into snapshot (query: QQ" + t.InsertSnapshotQuery() + "QQ): " + err.Error())
	}

	if _, err := tx.Exec(t.InsertDeletedSnapshotQuery()); err != nil {
		return errors.New("inserting " + t.Table + " deleted rows into snapshot (query QQ" + t.InsertDeletedSnapshotQuery() + "QQ): " + err.Error())
	}
	return nil
}

const SnapshotTableSuffix = `_snapshot`
const SnapshotTableLatestSuffix = `_latest`

func (t *SnapshotTable) SnapshotTable() string { return t.Table + SnapshotTableSuffix }
func (t *SnapshotTable) SnapshotLatestTable() string {
	return t.Table + SnapshotTableSuffix + SnapshotTableLatestSuffix
}

func (t *SnapshotTable) InsertSnapshotQuery() string {
	// TODO test performance, vs inserting only selected things newer than the values in the snapshot table.
	return `INSERT INTO "` + t.SnapshotTable() + `" SELECT * FROM "` + t.Table + `" ON CONFLICT DO NOTHING`
}

func (t *SnapshotTable) InsertDeletedSnapshotQuery() string {
	return `WITH
 ` + WithCDNSnapshotTimeLive() + `,
 ` + t.WithLatest() + `
INSERT INTO "` + t.SnapshotTable() + `" (
` + t.Columns + `, last_updated, deleted
) SELECT
` + t.Columns + `, now() as last_updated, true as deleted
FROM "` + t.SnapshotLatestTable() + `"
WHERE (` + t.PK + `) NOT IN (SELECT ` + t.PK + ` FROM ` + t.Table + `)
`
}

// WithLatest returns the WITH query part to get the latest snapshot times of this table,
// according to the snapshot table (which is per-CDN).
// This SHOULD NOT be used for tables belonging to the delivery service snapshot. Use WithLatestDS instead.
// The with table name is returned by SnapshotTable.SnapshotLatestTable.
// Note the query part requires CDNSnapshotTimeTable exist. See WithCDNSNapshotTime.
func (t *SnapshotTable) WithLatest() string {
	return t.SnapshotLatestTable() + ` AS (
  SELECT DISTINCT ON (` + t.PK + `) *
  FROM "` + t.SnapshotTable() + `" dt
  WHERE dt.last_updated <= (SELECT time FROM ` + CDNSnapshotTimeTable + `)
  ORDER BY ` + t.PK + `, last_updated DESC
)
`
}

// WithCDNSnapshotTime returns a query part, adding the CDN snapshot time as a WITH statement, thereafter available to select from as a table named CDNSNapshotTimeTable.
// The cdn argument is the name of the cdn to add.
// The live argument is whether to use the 'live' cdn table; if false, the cdn_snapshot table is used.
// The queryArgs argument is the existing query arguments, which should be all query arguments prior to this query part.
// The query part, and new query arguments are returned. Depending whether the live table is used, the query args may be returned unchanged, or with the cdn added.
//
func WithCDNSnapshotTime(cdn string, live bool, qryArgs []interface{}) (string, []interface{}) {
	if live {
		return WithCDNSnapshotTimeLive(), qryArgs
	}
	return WithCDNSnapshotTimeReal(cdn, qryArgs)
}

const CDNSnapshotTimeTable = "cdn_snapshot_time"

// WithCDNSnapshotTime returns CDN snapshot times as a "with" query part.
// Note the word "with" and a trailing comma are not included, so the withs can be chained.
func WithCDNSnapshotTimeReal(cdn string, qryArgs []interface{}) (string, []interface{}) {
	return `
` + CDNSnapshotTimeTable + ` AS (
  SELECT time from snapshot where cdn = $` + strconv.Itoa(len(qryArgs)+1) + `
)
`, append(qryArgs, cdn)
}

// WithCDNSnapshotTimeLive returns CDN snapshot times as a "with" query part, except the time is actually "now", to get a "live" snapshot.
// Note the word "with" and a trailing comma are not included, so the withs can be chained.
func WithCDNSnapshotTimeLive() string {
	return `
` + CDNSnapshotTimeTable + ` AS (
  SELECT now() as time
)
`
}

const DSSnapshotTimesTable = "ds_snapshot_times"

// WithDSSnapshotTime returns delivery service snapshot times as a "with" query part.
// This exists primarily to allow easy switching between snapshotted and live data.
// Note the word "with" and a trailing comma are not included, so the withs can be chained.
func WithDSSnapshotTimes(live bool) string {
	if live {
		return WithDSSnapshotTimesLive()
	}
	return WithDSSnapshotTimesReal()
}

func WithDSSnapshotTimesReal() string {
	return `
` + DSSnapshotTimesTable + ` AS (
  SELECT deliveryservice, time from deliveryservice_snapshots
)
`
}

// // WithDSSnapshotTimesLive returns DS snapshot times as a "with" query part, except it actually returns "now" for all DSes, in order to get a "live" snapshot.
func WithDSSnapshotTimesLive() string {
	return `
WITH ds_snapshot_time AS (
  SELECT xml_id as deliveryservice, now() from deliveryservice
)
`
}

// WithCDNSnapshotTimeLive returns CDN snapshot times as a "with" query part, except the time is actually "now", to get a "live" snapshot.
// Note the word "with" and a trailing comma are not included, so the withs can be chained.
func WithCDNSnapshotTimesLive() string {
	return `
` + DSSnapshotTimesTable + ` AS (
  SELECT deliveryservice, now() as time from deliveryservice_snapshots
)
`
}

func UpdateDSDSSSnapshot(tx *sql.Tx, ds tc.DeliveryServiceName) error {
	qry, qryArgs := InsertDSDSSSnapshotQuery(ds)
	if _, err := tx.Exec(qry, qryArgs...); err != nil {
		return errors.New("inserting  ds '" + string(ds) + "' dss snapshot (query: QQ" + qry + "QQ): " + err.Error())
	}

	qry, qryArgs = InsertDSDSSDeletedSnapshotQuery(ds)
	if _, err := tx.Exec(qry, qryArgs); err != nil {
		return errors.New("inserting  ds '" + string(ds) + "' dss deleted rows into snapshot (query: QQ" + qry + "QQ): " + err.Error())
	}

	return nil
}

// InsertDSDSSSnapshotQuery builds the insert query, to insert into the delivery service snapshot table.
// It returns the query, and the query arguments.
func InsertDSDSSSnapshotQuery(ds tc.DeliveryServiceName) (string, []interface{}) {
	// TODO put ds id in with statement, to only query once?
	// TODO test performance, vs inserting only selected things newer than the values in the snapshot table.
	t := DeliveryServiceServerTable
	return `INSERT INTO "` + t.SnapshotTable() + `"
SELECT * FROM "` + t.Table + `" dss
WHERE dss.deliveryservice = (SELECT id FROM deliveryservice WHERE xml_id = $1)
ON CONFLICT DO NOTHING`, []interface{}{ds}
}

// DSDSSWithLatest builds the 'with' query part, to select the latest delivery service as a table named DeliveryServiceServerTable.SnapshotLatestTable().
// It takes the delivery service to selet, and any prior query arguments, and returns the query part and new query arguments.
// It returns the query, and the query arguments.
func DSDSSWithLatest(ds tc.DeliveryServiceName, qryArgs []interface{}) (string, []interface{}) {
	// TODO put ds id in with statement, to only query once?
	t := DeliveryServiceServerTable
	return t.SnapshotLatestTable() + ` AS (
  SELECT DISTINCT ON (` + t.PK + `) *
  FROM "` + t.SnapshotTable() + `" dsssnap
  WHERE dsssnap.deliveryservice = (select id from deliveryservice where xml_id = $` + strconv.Itoa(len(qryArgs)+1) + `
  ORDER BY ` + t.PK + `, last_updated DESC
)
`, append(qryArgs, ds)
}

// InsertDSDSSDeletedSnapshotQuery builds the query to insert any deleted delivery services into the snapshot table.
// It returns the query, and the query arguments.
func InsertDSDSSDeletedSnapshotQuery(ds tc.DeliveryServiceName) (string, []interface{}) {
	// TODO put ds id in with statement, to only query once?
	t := DeliveryServiceServerTable
	qryArgs := []interface{}{}
	withDSDSSLatestQryPart, qryArgs := DSDSSWithLatest(ds, qryArgs)

	return `WITH ` + withDSDSSLatestQryPart + `
INSERT INTO "` + t.SnapshotTable() + `" (
` + t.Columns + `, last_updated, deleted
) SELECT
` + t.Columns + `, now() as last_updated, true as deleted
FROM "` + t.SnapshotLatestTable() + `" dssl
WHERE dssl.deliveryservice = (SELECT id FROM deliveryservice WHERE xml_id = $` + strconv.Itoa(len(qryArgs)+1) + `
  AND (` + t.PK + `) NOT IN (
    SELECT ` + t.PK + ` FROM ` + t.Table + `
  )
`, append(qryArgs, ds)
}
