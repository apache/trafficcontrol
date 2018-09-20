# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## [Unreleased]
### Added
- Backup Edge Cache group: If the matched group in the CZF is not available, this list of backup edge cache group configured via Traffic Ops API can be used as backup. In the event of all backup edge cache groups not available, GEO location can be optionally used as further backup. APIs detailed [here](http://traffic-control-cdn.readthedocs.io/en/latest/development/traffic_ops_api/v12/cachegroup_fallbacks.html)
- Traffic Ops Golang Proxy Endpoints
  - /api/1.3/origins `(GET,POST,PUT,DELETE)`
  - /api/1.3/coordinates `(GET,POST,PUT,DELETE)`
  - /api/1.3/staticdnsentries `(GET,POST,PUT,DELETE)`
  - /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys `GET`
  - /api/1.1/deliveryservices/hostname/:hostname/sslkeys `GET`
  - /api/1.1/deliveryservices/sslkeys/add `POST`
  - /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete `GET`
- Delivery Service Origins Refactor: The Delivery Service API now creates/updates an Origin entity on Delivery Service creates/updates, and the `org_server_fqdn` column in the `deliveryservice` table has been removed. The `org_server_fqdn` data is now computed from the Delivery Service's primary origin (note: the name of the primary origin is the `xml_id` of its delivery service).
- Cachegroup-Coordinate Refactor: The Cachegroup API now creates/updates a Coordinate entity on Cachegroup creates/updates, and the `latitude` and `longitude` columns in the `cachegroup` table have been replaced with `coordinate` (a foreign key to Coordinate). Coordinates created from Cachegroups are given the name `from_cachegroup_\<cachegroup name\>`.
- Geolocation-based Client Steering: two new steering target types are available to use for `CLIENT_STEERING` delivery services: `STEERING_GEO_ORDER` and `STEERING_GEO_WEIGHT`. When targets of these types have an Origin with a Coordinate, Traffic Router will order and prioritize them based upon the shortest total distance from client -> edge -> origin. Co-located targets are grouped together and can be weighted or ordered within the same location using `STEERING_GEO_WEIGHT` or `STEERING_GEO_ORDER`, respectively.
- Tenancy is now the default behavior in Traffic Ops.  All database entries that reference a tenant now have a default of the root tenant.  This eliminates the need for the `use_tenancy` global parameter and will allow for code to be simplified as a result. If all user and delivery services reference the root tenant, then there will be no difference from having `use_tenancy` set to 0.
- Traffic Monitor Client Update: Traffic Monitor is updated to use the Traffic Ops v13 client.
- Removed previously deprecated `traffic_monitor_java`
- Added `infrastructure/cdn-in-a-box` for Apachecon 2018 demonstration

### Changed
- Issue 2821: Fixed "Traffic Router may choose wrong certificate when SNI names overlap"

## [2.2.0] - 2018-06-07
### Added
- Per-DeliveryService Routing Names: you can now choose a Delivery Service's Routing Name (rather than a hardcoded "tr" or "edge" name). This might require a few pre-upgrade steps detailed [here](http://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops/migration_from_20_to_22.html#per-deliveryservice-routing-names)
- [Delivery Service Requests](http://traffic-control-cdn.readthedocs.io/en/latest/admin/quick_howto/ds_requests.html#ds-requests): When enabled, delivery service requests are created when ALL users attempt to create, update or delete a delivery service. This allows users with higher level permissions to review delivery service changes for completeness and accuracy before deploying the changes.
- Traffic Ops Golang Proxy Endpoints
  - /api/1.3/about `(GET)`
  - /api/1.3/asns `(GET,POST,PUT,DELETE)`
  - /api/1.3/cachegroups `(GET,POST,PUT,DELETE)`
  - /api/1.3/cdns `(GET,POST,PUT,DELETE)`
  - /api/1.3/cdns/capacity `(GET)`
  - /api/1.3/cdns/configs `(GET)`
  - /api/1.3/cdns/dnsseckeys `(GET)`
  - /api/1.3/cdns/domain `(GET)`
  - /api/1.3/cdns/monitoring `(GET)`
  - /api/1.3/cdns/health `(GET)`
  - /api/1.3/cdns/routing `(GET)`
  - /api/1.3/deliveryservice_requests `(GET,POST,PUT,DELETE)`
  - /api/1.3/divisions `(GET,POST,PUT,DELETE)`
  - /api/1.3/hwinfos `(GET)`
  - /api/1.3/login `(POST)`
  - /api/1.3/parameters `(GET,POST,PUT,DELETE)`
  - /api/1.3/profileparameters `(GET,POST,PUT,DELETE)`
  - /api/1.3/phys_locations `(GET,POST,PUT,DELETE)`
  - /api/1.3/ping `(GET)`
  - /api/1.3/profiles `(GET,POST,PUT,DELETE)`
  - /api/1.3/regions `(GET,POST,PUT,DELETE)`
  - /api/1.3/servers `(GET,POST,PUT,DELETE)`
  - /api/1.3/servers/checks `(GET)`
  - /api/1.3/servers/details `(GET)`
  - /api/1.3/servers/status `(GET)`
  - /api/1.3/servers/totals `(GET)`
  - /api/1.3/statuses `(GET,POST,PUT,DELETE)`
  - /api/1.3/system/info `(GET)`
  - /api/1.3/types `(GET,POST,PUT,DELETE)`
- Fair Queuing Pacing: Using the FQ Pacing Rate parameter in Delivery Services allows operators to limit the rate of individual sessions to the edge cache. This feature requires a Trafficserver RPM containing the fq_pacing experimental plugin AND setting 'fq' as the default Linux qdisc in sysctl.
- Traffic Ops rpm changed to remove world-read permission from configuration files.

### Changed
- Reformatted this CHANGELOG file to the keep-a-changelog format

[Unreleased]: https://github.com/apache/trafficcontrol/compare/RELEASE-2.2.0...HEAD
[2.2.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-2.1.0...RELEASE-2.2.0
