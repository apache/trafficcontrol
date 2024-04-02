# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## [unreleased]
### Added
- [#8014](https://github.com/apache/trafficcontrol/pull/8014) *Traffic Ops* Added logs to indicate which mechanism a client used to login to TO.
- [#7812](https://github.com/apache/trafficcontrol/pull/7812) *Traffic Portal*: Expose the `configUpdateFailed` and `revalUpdateFailed` fields on the server table.
- [#7870](https://github.com/apache/trafficcontrol/pull/7870) *Traffic Portal*: Adds a hyperlink to the DSR page to the DS itself for ease of navigation.
- [#7896](https://github.com/apache/trafficcontrol/pull/7896) *ATC Build system*: Count commits since the last release, not commits
- [#7927](https://github.com/apache/trafficcontrol/pull/7927) *Traffic Stats*: Migrate dynamic scripted Grafana Dashboards to Scenes

### Changed
- [#7614](https://github.com/apache/trafficcontrol/pull/7614) *Traffic Ops* The database upgrade process no longer overwrites changes users may have made to the initially seeded data.
- [#7832](https://github.com/apache/trafficcontrol/pull/7832) *t3c* Removed perl dependency
- Updated the CacheGroups Traffic Portal page to use a more performant AG-Grid-based table.
- Updated Go version to 1.22.0
- [#7958] Updated build ATS to 9.2.4
- [#7979](https://github.com/apache/trafficcontrol/pull/7979) *Traffic Router*, *Traffic Monitor*, *Traffic Stats*: Store logs in /var/log
- [#7999](https://github.com/apache/trafficcontrol/pull/7999) *Traffic Router*, *Traffic Monitor*, *Traffic Stats*: Symlink from /opt/<component>/var/log to /var/log/<component>. These symlinks are deprecated with the intent of removing them in ATC 9.0.0.
- [#7990](https://github.com/apache/trafficcontrol/pull/7990) *Traffic Router*: Updated Apache Tomcat from 9.0.43, 9.0.67, 9.0.83, and 9.0.86 to 9.0.87.
- [#7933](https://github.com/apache/trafficcontrol/pull/7933), [#8005](https://github.com/apache/trafficcontrol/pull/8005) *Traffic Portal v2*: Update NodeJS version to 18.
- [#8009](https://github.com/apache/trafficcontrol/pull/8009) *Traffic Portal v2*: Update NodeJS version to 20.
- [#8040](https://github.com/apache/trafficcontrol/pull/8040) *Traffic Router*: Get the Tomcat version from .env and update Tomcat to 9.0.90.
- [##8056](https://github.com/apache/trafficcontrol/pull/8056) Remove the `version` key from compose files and use `docker compose` instead of `docker-compose`.
- [7980](https://github.com/apache/trafficcontrol/pull/7980) *Traffic Server*: Store logs in /var/log

### Fixed
- [#8008](https://github.com/apache/trafficcontrol/pull/8008) *Traffic Router* Fix czf temp file deletion issue.
- [#7998](https://github.com/apache/trafficcontrol/pull/7998) *Traffic Portalv2* Fixed (create and update) page titles across every feature
- [#7984](https://github.com/apache/trafficcontrol/pull/7984) *Traffic Ops* Fixed TO Client cert authentication with respect to returning response cookie.
- [#7957](https://github.com/apache/trafficcontrol/pull/7957) *Traffic Ops* Fix the incorrect display of delivery services assigned to ORG servers.
- [#7917](https://github.com/apache/trafficcontrol/pull/7917) *Traffic Ops* Removed `Alerts` field from struct `ProfileExportResponse`.
- [#7918](https://github.com/apache/trafficcontrol/pull/7918) *Traffic Portal* Fixed topology link under DS-Servers tables page
- [#7846](https://github.com/apache/trafficcontrol/pull/7846) *Traffic Portal* Increase State character limit
- [#8010](https://github.com/apache/trafficcontrol/pull/8010) *Traffic Stats* Omit NPM dev dependencies from Traffic Stats RPM

### Removed
- [#7832](https://github.com/apache/trafficcontrol/pull/7832) *t3c* Removed Perl dependency
- [#7841](https://github.com/apache/trafficcontrol/pull/7841) *Postinstall* Removed Perl implementation and Python 2.x support

## [8.0.0] - 2024-01-30
### Added
- [#7672](https://github.com/apache/trafficcontrol/pull/7672) *Traffic Control Health Client*: Added peer monitor flag while using `strategies.yaml`.
- [#7609](https://github.com/apache/trafficcontrol/pull/7609) *Traffic Portal*: Added Scope Query Param to SSO login.
- [#7450](https://github.com/apache/trafficcontrol/pull/7450) *Traffic Ops*: Removed hypnotoad section and added listen field to traffic_ops_golang section in order to simplify cdn config.
- [#7291](https://github.com/apache/trafficcontrol/pull/7291) *Traffic Ops*: Extended Layered Profile feature to aggregate parameters for all server profiles.
- [#7314](https://github.com/apache/trafficcontrol/pull/7314) *Traffic Portal*: Added capability feature to Delivery Service Form (HTTP, DNS).
- [#7295](https://github.com/apache/trafficcontrol/pull/7295) *Traffic Portal*: Added description and priority order for Layered Profile on server form.
- [#6234](https://github.com/apache/trafficcontrol/issues/6234) *Traffic Ops, Traffic Portal*: Added description field to Server Capabilities.
- [#6033](https://github.com/apache/trafficcontrol/issues/6033) *Traffic Ops, Traffic Portal*: Added ability to assign multiple servers per capability.
- [#7081](https://github.com/apache/trafficcontrol/issues/7081) *Traffic Router*: Added better log messages for TR connection exceptions.
- [#7089](https://github.com/apache/trafficcontrol/issues/7089) *Traffic Router*: Added the ability to specify HTTPS certificate attributes.
- [#7109](https://github.com/apache/trafficcontrol/pull/7109) *Traffic Router*: Removed `dnssec.zone.diffing.enabled` and `dnssec.rrsig.cache.enabled` parameters.
- [#7075](https://github.com/apache/trafficcontrol/pull/7075) *Traffic Portal*: Added the `lastUpdated` field to all delivery service forms.
- [#7055](https://github.com/apache/trafficcontrol/issues/7055) *Traffic Portal*: Made `Clear Table Filters` option visible to the user.
- [#7024](https://github.com/apache/trafficcontrol/pull/7024) *Traffic Monitor*: Added logging for `ipv4Availability` and `ipv6Availability` in TM.
- [#7063](https://github.com/apache/trafficcontrol/pull/7063) *Traffic Ops*: Added API version 5.0 (IN DEVELOPMENT).
- [#7645](https://github.com/apache/trafficcontrol/pull/7645) *Traffic Ops*: Added a client method to be able to login with certificates.
- [#2101](https://github.com/apache/trafficcontrol/issues/2101) *Traffic Portal*: Added the ability to tell if a Delivery Service is the target of another steering DS.
- [#6021](https://github.com/apache/trafficcontrol/issues/6021) *Traffic Portal*: Added the ability to view a change logs message in it's entirety by clicking on it.
- [#7078](https://github.com/apache/trafficcontrol/issues/7078) *Traffic Ops, Traffic Portal*: Added ability to assign multiple server capabilities to a server.
- [#7096](https://github.com/apache/trafficcontrol/issues/7096) *Traffic Control Health Client*: Added health client parent health.
- [#7032](https://github.com/apache/trafficcontrol/issues/7032) *Traffic Control Cache Config (t3c)*: Add t3c-apply flag to use local ATS version for config generation rather than Server package Parameter, to allow managing the ATS OS package via external tools. See 'man t3c-apply' and 'man t3c-generate' for details.
- [#7097](https://github.com/apache/trafficcontrol/issues/7097) *Traffic Ops, Traffic Portal, Traffic Control Cache Config (t3c)*: Added the `regional` field to Delivery Services, which affects whether `maxOriginConnections` should be per Cache Group.
- [#2388](https://github.com/apache/trafficcontrol/issues/2388) *Trafic Ops, Traffic Portal*: Added the `TTLOverride` field to CDNs, which lets you override all TTLs in all Delivery Services of a CDN's snapshot with a single value.
- [#7176](https://github.com/apache/trafficcontrol/pull/7176) *ATC Build system*: Support building ATC for the `aarch64` CPU architecture.
- [#7113](https://github.com/apache/trafficcontrol/pull/7113) *Traffic Portal*: Minimize the Server Server Capability part of the *Traffic Servers* section of the Snapshot Diff.
- [#7273](https://github.com/apache/trafficcontrol/pull/7273) *Traffic Ops*: Adds `SSL-KEY-EXPIRATION:READ` permission to operations, portal, read-only, federation and steering roles.
- [#7343](https://github.com/apache/trafficcontrol/pull/7343) *Traffic Ops* Adds `ACME:READ`, `CDNI-ADMIN:READ` and `CDNI-CAPACITY:READ` permissions to operations, portal, read-only, federation and steering roles.
- [#7296](https://github.com/apache/trafficcontrol/pull/7296) *Traffic Portal*: New configuration option in `traffic_portal_properties.json` at `deliveryServices.exposeInactive` controls exposing APIv5 DS Active State options in the TP UI.
- [#7332](https://github.com/apache/trafficcontrol/pull/7332) *Traffic Ops*: Creates new role needed for TR to watch TO resources.
- [#7322](https://github.com/apache/trafficcontrol/issues/7322) *Traffic Control Cache Config (t3c)*: Adds support for anycast on http routed edges.
- [#7367](https://github.com/apache/trafficcontrol/pull/7367) *Traffic Ops*: Adds `ACME:CREATE`, `ACME:DELETE`, `ACME:DELETE`, and `ACME:READ` permissions to operations role.
- [#7380](https://github.com/apache/trafficcontrol/pull/7380) *Traffic Portal*: Adds strikethrough (expired), red (7 days until expiration) and yellow (30 days until expiration) visuals to delivery service cert expiration grid rows.
- [#7388](https://github.com/apache/trafficcontrol/pull/7388) *TC go Client*: Adds sslkey_expiration methodology in v4 and v5 clients.
- [#7543](https://github.com/apache/trafficcontrol/pull/7543) *Traffic Portal*: New Ansible Role to use Traffic Portal v2.
- [#7516](https://github.com/apache/trafficcontrol/pull/7516) *Traffic Control Cache Config (t3c)*: added command line arg to control go_direct in parent.config.
- [#7602](https://github.com/apache/trafficcontrol/pull/7602) *Traffic Control Cache Config (t3c)*: added installed package data to t3c-apply-metadata.json.
- [#7618](https://github.com/apache/trafficcontrol/pull/7618) *Traffic Portal*: Add the ability to inspect a user provider cert, or the cert chain on DS SSL keys.
- [#7619](https://github.com/apache/trafficcontrol/pull/7619) *Traffic Ops*: added optional field `oauth_user_attribute` for OAuth login credentials.
- [#7641](https://github.com/apache/trafficcontrol/pull/7641) *Traffic Router*: Added further optimization to TR's algorithm of figuring out the zone for an incoming request.
- [#7646](https://github.com/apache/trafficcontrol/pull/7646) *Traffic Portal*: Add the ability to delete a cert.
- [#7652](https://github.com/apache/trafficcontrol/pull/7652) *Traffic Control Cache Config (t3c)*: added rpmdb checks and use package data from t3c-apply-metadata.json if rpmdb is corrupt.
- [#7674](https://github.com/apache/trafficcontrol/issues/7674) *Traffic Ops*: Add the ability to indicate if a server failed its revalidate/config update.
- [#7784](https://github.com/apache/trafficcontrol/pull/7784) *Traffic Portal*: Added revert certificate functionality to the ssl-keys page.
- [#7719](https://github.com/apache/trafficcontrol/pull/7719) *t3c* self-healing will be added automatically when using the slice plugin.

### Changed
- [#7776](https://github.com/apache/trafficcontrol/pull/7776) *tc-health-client*: Added error message while issues interacting with Traffic Ops.
- [#7766](https://github.com/apache/trafficcontrol/pull/7766) *Traffic Portal*: now uses Traffic Ops APIv5.
- [#7765](https://github.com/apache/trafficcontrol/pull/7765) *Traffic Stats*: now uses Traffic Ops APIv5.
- [#7761](https://github.com/apache/trafficcontrol/pull/7761) *Traffic Monitor*: Use API v5 for the TM's Traffic Ops client, and use TO API v4 for TM's Traffic Ops legacy client.
- [#7757](https://github.com/apache/trafficcontrol/pull/7757) *Traffic Router*: Changed Traffic Router to point to API version 5.0 of Traffic Ops.
- [#7732](https://github.com/apache/trafficcontrol/pull/7732) *Traffic Router*: Increased negative TTL value to 900 seconds.
- [#7665](https://github.com/apache/trafficcontrol/pull/7665) *Automation*: Changes to Ansible role dataset_loader to add ATS 9 support.
- [#7761](https://github.com/apache/trafficcontrol/pull/7761) *Traffic Monitor*: Use API v5 for the TM's Traffic Ops client, and use TO API v4 for TM's Traffic Ops legacy client.
- [#7584](https://github.com/apache/trafficcontrol/pull/7584) *Docs*: Upgrade Traffic Control Sphinx documentation Makefile OS intelligent.
- [#7521](https://github.com/apache/trafficcontrol/pull/7521) *Traffic Ops*: Returns empty array instead of null when no permissions are given for roles endpoint using POST or PUT request.
- [#7369](https://github.com/apache/trafficcontrol/pull/7369) *Traffic Portal*: Adds better labels to routing methods widget on the TP dashboard.
- [#7371](https://github.com/apache/trafficcontrol/pull/7371) *Traffic Portal*: Simplifies DS button bar by moving DS changes / DSRs under More menu and renaming to 'View Change Requests'.
- [#7224](https://github.com/apache/trafficcontrol/pull/7224) *Traffic Ops*: Required Capabilities are now a part of the `DeliveryService` structure.
- [#7063](https://github.com/apache/trafficcontrol/pull/7063) *Traffic Ops*: Python client now uses Traffic Ops API 4.1 by default.
- [#6981](https://github.com/apache/trafficcontrol/pull/6981) *Traffic Portal*: Obscures sensitive text in Delivery Service "Raw Remap" fields, private SSL keys, "Header Rewrite" rules, and ILO interface passwords by default.
- [#7037](https://github.com/apache/trafficcontrol/pull/7037) *Traffic Router*: Uses Traffic Ops API 4.0 by default.
- [#7191](https://github.com/apache/trafficcontrol/issues/7191) *tc-health-client*: Uses Traffic Ops API 4.0. Also added reload option to systemd service file.
- [#4654](https://github.com/apache/trafficcontrol/pull/4654) *Traffic Ops, Traffic Portal*: Switched Delivery Service active state to a three-value system, adding a state that will be used to prevent cache servers from deploying DS configuration.
- [#7242](https://github.com/apache/trafficcontrol/pull/7276) *Traffic Portal*: Now depends on NodeJS version 16 or later.
- [#7120](https://github.com/apache/trafficcontrol/pull/7120) *Docs*: Update t3c documentation regarding parent.config parent_retry.
- [#7044](https://github.com/apache/trafficcontrol/issues/7044) *CDN in a Box*: [CDN in a Box, the t3c integration tests, and the tc health client integration tests now use Apache Traffic Server 9.1.
- [#7366](https://github.com/apache/trafficcontrol/pull/7366) *Traffic Control Cache Config (t3c)*: Removed timestamp from metadata file since it's changing every minute and causing excessive commits to git repo.
- [#7386](https://github.com/apache/trafficcontrol/pull/7386) *Traffic Portal*: Increased the number of events that are logged to the TP access log.
- [#7469](https://github.com/apache/trafficcontrol/pull/7469) *Traffic Ops*: Changed logic to not report empty or missing cookies into TO error.log.
- [#7586](https://github.com/apache/trafficcontrol/pull/7586) *Traffic Ops*: Add permission to Operations Role to read from dnsseckeys endpoint.
- [#7600](https://github.com/apache/trafficcontrol/pull/7600) *Traffic Control Cache Config (t3c)*: changed default go-direct command line arg to be old to avoid unexpected config changes upon upgrade.
- [#7621](https://github.com/apache/trafficcontrol/pull/7621) *Traffic Ops*: Use ID token for OAuth authentication, not Access Token.
- [#7694](https://github.com/apache/trafficcontrol/pull/7694) *Traffic Control Cache Config (t3c)*, *Traffic Control Health Client*: Upgrade to ATS 9.2.
- [#7966](https://github.com/apache/trafficcontrol/pull/7696) *Traffic Control Cache Config (t3c)*: will no longer clear update flag when config failure occurs and will also give a cache config error msg on exit.
- [#7716](https://github.com/apache/trafficcontrol/pull/7716) *Apache Traffic Server*: Use GCC 11 for building.
- [#7742](https://github.com/apache/trafficcontrol/pull/7742) *Traffic Ops*: Changed api tests to supply the absolute path of certs.
- [#7814](https://github.com/apache/trafficcontrol/issues/7814) All Go components: Updated the module path to [`github.com/apache/trafficcontrol/v8`](https://pkg.go.dev/github.com/apache/trafficcontrol/v8). Module https://pkg.go.dev/github.com/apache/trafficcontrol will not receive further updates.

### Fixed
- [#7891](https://github.com/apache/trafficcontrol/pull/7891) *Traffic Ops*: Updated job routes for v4 and v5.
- [#7890](https://github.com/apache/trafficcontrol/pull/7890) *Traffic Ops*: Fixed missing changelog entries to v5 routes.
- [#7887](https://github.com/apache/trafficcontrol/pull/7887) *Traffic Ops*: Limit Delivery Services returned for GET /servers/{id}/deliveryservices to ones in the same CDN
- [#7885](https://github.com/apache/trafficcontrol/pull/7885) *Traffic Portal*: Fixed the issue where Compare Profiles page was not being displayed.
- [#7879](https://github.com/apache/trafficcontrol/7879) *Traffic Ops, Traffic Portal*: Fixed broken capability links for delivery service and added required capability as a column in DS table.
- [#7878](https://github.com/apache/trafficcontrol/pull/7878) *Traffic Ops, Traffic Portal*: Fixed the case where TO was failing to assign delivery services to a server, due to a bug in the way the list of preexisting delivery services was being returned.
- [#7819](https://github.com/apache/trafficcontrol/pull/7819) *Traffic Ops*: API v5 routes should not use `privLevel` comparisons.
- [#7802](https://github.com/apache/trafficcontrol/pull/7802) *Traffic Control Health Client*: Fixed ReadMe.md typos and duplicates.
- [#7764](https://github.com/apache/trafficcontrol/pull/7764) *Traffic Ops*: Collapsed DB migrations.
- [#7767](https://github.com/apache/trafficcontrol/pull/7767) *Traffic Ops*: Fixed ASN update logic for APIv5.
- [RFC3339](https://github.com/apache/trafficcontrol/issues/5911)
    - [#7806](https://github.com/apache/trafficcontrol/pull/7806) *Traffic Ops*: Fixed `cdns/{{name}}/federations` and `cdns/{{name}}/federations/{{ID}}` v5 APIs to respond with `RFC3339` timestamps.
    - [#7759](https://github.com/apache/trafficcontrol/pull/7759) *Traffic Ops*: Fixed `/profiles/{{ID}}/parameters` and `profiles/name/{{name}}/parameters` v5 APIs to respond with `RFC3339` timestamps.
    - [#7734](https://github.com/apache/trafficcontrol/pull/7734) *Traffic Ops*: Fixed `/profiles` v5 APIs to respond with `RFC3339` timestamps.
    - [#7718](https://github.com/apache/trafficcontrol/pull/7718) *Traffic Ops*: `/servers` endpoint now responds with `RFC3339` timestamps for all timestamp fields. Cleaned up naming conventions and superfluous data.
    - [#7708](https://github.com/apache/trafficcontrol/pull/7708) *Traffic Ops*: Fixed `/parameters` v5 APIs to respond with `RFC3339` timestamps.
    - [#7749](https://github.com/apache/trafficcontrol/pull/7749) *Traffic Ops*: Fixed `/tenants` v5 APIs to respond with `RFC3339` timestamp.
    - [#7740](https://github.com/apache/trafficcontrol/pull/7740) *Traffic Ops*: Fixed `/staticDNSEntries` v5 APIs to respond with `RFC3339` timestamps.
    - [#7738](https://github.com/apache/trafficcontrol/pull/7738) *Traffic Ops*: Fixed `/profileparameters` v5 APIs to respond with `RFC3339` timestamps.
    - [#7690](https://github.com/apache/trafficcontrol/pull/7690) *Traffic Ops*: Fixed `/logs` v5 APIs to respond with RFC3339 timestamps.
    - [#7605](https://github.com/apache/trafficcontrol/pull/7605) *Traffic Ops*: Fixed `/cachegroups_request_comments` v5 APIs to respond with `RFC3339` timestamps.
    - [#7720](https://github.com/apache/trafficcontrol/pull/7720) *Traffic Ops*: Fixed Delivery Service Servers v5 APIs to respond with `RFC3339` timestamps.
    - [#7631](https://github.com/apache/trafficcontrol/pull/7631) *Traffic Ops*: Fixed `/phys_location` v5 APIs to respond with `RFC3339` timestamps.
    - [#7612](https://github.com/apache/trafficcontrol/pull/7612) *Traffic Ops*: Fixed `/divisions` v5 APIs to respond with `RFC3339` timestamps.
    - [#7561](https://github.com/apache/trafficcontrol/pull/7561) *Traffic Ops*: Fixed `/asns` v5 APIs to respond with `RFC3339` timestamps.
    - [#7575](https://github.com/apache/trafficcontrol/pull/7575) *Traffic Ops*: Fixed `/types` v5 APIs to respond with `RFC3339` timestamps.
    - [#7698](https://github.com/apache/trafficcontrol/pull/7698) *Traffic Ops*: Fixed `/region` v5 APIs to respond with `RFC3339` timestamps.
    - [#7660](https://github.com/apache/trafficcontrol/pull/7660) *Traffic Ops*: Fixed `/deliveryServices` v5 APIs to respond with `RFC3339` timestamps.
    - [#7570](https://github.com/apache/trafficcontrol/pull/7570) *Traffic Ops*: Fixed `/deliveryservice_request_comments` v5 APIs to respond with `RFC3339` timestamps.
    - [#7596](https://github.com/apache/trafficcontrol/pull/7596) *Traffic Ops*: Fixed `/federation_resolvers` v5 APIs to respond with `RFC3339` timestamps.
    - [#7572](https://github.com/apache/trafficcontrol/pull/7572) *Traffic Ops*: Fixed `/deliveryservice_requests` v5 APIs docs with `RFC3339` timestamps
    - [#7545](https://github.com/apache/trafficcontrol/pull/7545) *Traffic Ops*: Fixed `/stats_summary` v5 APIs to respond with `RFC3339` timestamps.
    - [#7542](https://github.com/apache/trafficcontrol/pull/7542) *Traffic Ops*: Fixed `cdn_locks` documentation to reflect the correct `RFC3339` timestamps.
    - [#7482](https://github.com/apache/trafficcontrol/pull/7482) *Traffic Ops*: Fixed `/server_capabilities` v5 APIs to respond with RFC3339 timestamps.
    - [#7691](https://github.com/apache/trafficcontrol/pull/7691) *Traffic Ops*: Fixed `/topologies` v5 APIs to respond with `RFC3339` timestamps.
    - [#7408](https://github.com/apache/trafficcontrol/pull/7408) *Traffic Ops*: Fixed `/service_category` v5 APIs to respond with `RFC3339` timestamps.
    - [#7707](https://github.com/apache/trafficcontrol/pull/7707) *Traffic Ops*: Fixed `/statuses` v5 APIs to respond with `RFC3339` timestamps.
    - [#7733](https://github.com/apache/trafficcontrol/pull/7733) *Traffic Ops*: Fixed `/origins` v5 APIs to respond with `RFC3339` timestamps.
    - [#7744](https://github.com/apache/trafficcontrol/pull/7744) *Traffic Ops*: Fixed `/server_server_capabilities` v5 APIs to respond with `RFC3339` timestamps.
- [#7762](https://github.com/apache/trafficcontrol/pull/7762) *Traffic Ops*: Fixed `/phys_locations` update API to remove error related to mismatching region name and ID.
- [#7730](https://github.com/apache/trafficcontrol/pull/7730) *Traffic Monitor*: Fixed the panic seen in TM when `plugin.system_stats.timestamp_ms` appears as float and not string.
- [#4393](https://github.com/apache/trafficcontrol/issues/4393) *Traffic Ops*: Fixed the error code and alert structure when TO is queried for a delivery service with no ssl keys.
- [#7623](https://github.com/apache/trafficcontrol/pull/7623) *Traffic Ops*: Removed TryIfModifiedSinceQuery from `servicecategories.go` and reused from `ims.go`.
- [#7608](https://github.com/apache/trafficcontrol/pull/7608) *Traffic Monitor*: Use stats_over_http(plugin.system_stats.timestamp_ms) timestamp field to calculate bandwidth for TM's caches.
- [#6318](https://github.com/apache/trafficcontrol/issues/6318) *Docs*: Included docs for POST, PUT, DELETE (v3,v4,v5) for statuses and statuses{id}.
- [#7598](https://github.com/apache/trafficcontrol/pull/7598) *Traffic Ops*: Fixed Server Capability V5 Type Name Minor version.
- [#7312](https://github.com/apache/trafficcontrol/issues/7312) *Docs*: Changing docs for `cdn_locks` for DELETE response structure v4 and v5.
- [#6340](https://github.com/apache/trafficcontrol/issues/6340) *Traffic Ops*: Fixed alert messages for POST and PUT invalidation job APIs.
- [#7519](https://github.com/apache/trafficcontrol/issues/7519) *Traffic Ops*: Fixed TO API `/servers/{id}/deliveryservices` endpoint to responding with all DS's on cache that are directly assigned and inherited through topology.
- [#7511](https://github.com/apache/trafficcontrol/pull/7511) *Traffic Ops*: Fixed the changelog registration message to include the username instead of duplicate email entry.
- [#7441](https://github.com/apache/trafficcontrol/pull/7441) *Traffic Ops*: Fixed the invalidation jobs endpoint to respect CDN locks.
- [#7414](https://github.com/apache/trafficcontrol/pull/7414) *Traffic Portal*: Fixed DSR difference for DS required capability.
- [#7130](https://github.com/apache/trafficcontrol/issues/7130) *Traffic Ops*: Fixed service_categories response to POST API.
- [#7340](https://github.com/apache/trafficcontrol/pull/7340) *Traffic Router*: Fixed TR logging for the `cqhv` field when absent.
- [#5557](https://github.com/apache/trafficcontrol/issues/5557) *Traffic Portal*: Moved `Fair Queueing Pacing Rate Bps` DS field to `Cache Configuration Settings` section.
- [#7252](https://github.com/apache/trafficcontrol/issues/7252) *Traffic Router*: Fixed integer overflow for `czCount`, by resetting the count to max value when it overflows.
- [#7221](https://github.com/apache/trafficcontrol/issues/7221) *Docs*: Fixed request structure spec in cdn locks description in APIv4.
- [#7225](https://github.com/apache/trafficcontrol/issues/7225) *Docs*: Fixed docs for `/roles` response description in APIv4.
- [#7246](https://github.com/apache/trafficcontrol/issues/7246) *Docs*: Fixed docs for `/jobs` response description in APIv4 and APIv5.
- [#6229](https://github.com/apache/trafficcontrol/issues/6229) *Traffic Ops*: Fixed error message for assignment of non-existent parameters to a profile.
- [#7231](https://github.com/apache/trafficcontrol/pull/7231) *Traffic Ops, Traffic Portal*: Fixed `sharedUserNames` display while retrieving CDN locks.
- [#7216](https://github.com/apache/trafficcontrol/pull/7216) *Traffic Portal*: Fixed sort for Server's Capabilities Table.
- [#4428](https://github.com/apache/trafficcontrol/issues/4428) *Traffic Ops*: Fixed Internal Server Error with POST to `profileparameters` when POST body is empty.
- [#7179](https://github.com/apache/trafficcontrol/issues/7179) *Traffic Portal*: Fixed search filter for Delivery Service Table.
- [#7174](https://github.com/apache/trafficcontrol/issues/7174) *Traffic Portal*: Fixed topologies sort (table and Delivery Service's form).
- [#5970](https://github.com/apache/trafficcontrol/issues/5970) *Traffic Portal*: Fixed numeric sort in Delivery Service's form for DSCP.
- [#5971](https://github.com/apache/trafficcontrol/issues/5971) *Traffic Portal*: Fixed Max DNS Tool Top link to open in a new page.
- [#7131](https://github.com/apache/trafficcontrol/issues/7131) *Docs*: Fixed Docs for staticdnsentries API endpoint missing lastUpdated response property description in APIv3, APIv4 and APIv5.
- [#6947](https://github.com/apache/trafficcontrol/issues/6947) *Docs*: Fixed docs for `cdns/{{name}}/federations` in APIv3, APIv4 and APIv5.
- [#6903](https://github.com/apache/trafficcontrol/issues/6903) *Docs*: Fixed docs for /cdns/dnsseckeys/refresh in APIv4 and APIv5.
- [#7049](https://github.com/apache/trafficcontrol/issues/7049), [#7052](https://github.com/apache/trafficcontrol/issues/7052) *Traffic Portal*: Fixed server table's quick search and filter option for multiple profiles.
- [#7080](https://github.com/apache/trafficcontrol/issues/7080), [#6335](https://github.com/apache/trafficcontrol/issues/6335) *Traffic Portal*: Fixed redirect links for server capability.
- [#7022](https://github.com/apache/trafficcontrol/pull/7022) *Traffic Stats*: Reuse InfluxDB client handle to prevent potential connection leaks.
- [#7021](https://github.com/apache/trafficcontrol/issues/7021) *Traffic Control Cache Config (t3c)*: Fixed cache config for Delivery Services with IP Origins.
- [#7043](https://github.com/apache/trafficcontrol/issues/7043) *Traffic Control Cache Config (t3c)*: Fixed cache config missing retry parameters for non-topology MSO Delivery Services going direct from edge to origin.
- [#7047](https://github.com/apache/trafficcontrol/issues/7047) *Traffic Ops*: allow `apply_time` query parameters on the `servers/{id-name}/update` when the CDN is locked.
- [#7163](https://github.com/apache/trafficcontrol/issues/7163) *Traffic Control Cache Config (t3c)*: Fix cache config for multiple profiles.
- [#7048](https://github.com/apache/trafficcontrol/issues/7048) *Traffic Stats*: Add configuration value to set the client request timeout for calls to Traffic Ops.
- [#7093](https://github.com/apache/trafficcontrol/issues/7093) *Traffic Router*: Updated Apache Tomcat from 9.0.43 to 9.0.67.
- [#7125](https://github.com/apache/trafficcontrol/issues/7125) *Docs*: Reflect implementation and deprecation notice for `letsencrypt/autorenew` endpoint.
- [#7046](https://github.com/apache/trafficcontrol/issues/7046) *Traffic Ops*: `deliveryservices/sslkeys/add` now checks that each cert in the chain is related.
- [#7158](https://github.com/apache/trafficcontrol/issues/7158) *Traffic Vault*: Fix the `reencrypt` utility to uniquely reencrypt each version of the SSL Certificates.
- [#7137](https://github.com/apache/trafficcontrol/pull/7137) *Traffic Control Cache Config (t3c)*: parent.config simulate topology for non topo delivery services.
- [#7153](https://github.com/apache/trafficcontrol/pull/7153) *Traffic Control Cache Config (t3c)*: Adds an extra T3C check for validity of an ssl cert (crash fix).
- [#3965](https://github.com/apache/trafficcontrol/issues/3965) *Traffic Router*: TR now always includes a `Content-Length` header in the response.
- [#6533](https://github.com/apache/trafficcontrol/issues/6533) *Traffic Router*: TR should not rename/recreate log files on rollover
- [#7182](https://github.com/apache/trafficcontrol/pull/7182) *Traffic Control Cache Config (t3c)*: Sort peers used in strategy.yaml to prevent false positive for reload.
- [#7204](https://github.com/apache/trafficcontrol/pull/7204) *Traffic Control Cache Config (t3c)*: strategies.yaml hash_key only for consistent_hash.
- [#7277](https://github.com/apache/trafficcontrol/pull/7277) *Traffic Control Cache Config (t3c)*: remapdotconfig: remove skip check at mids for nocache/live.
- [#7282](https://github.com/apache/trafficcontrol/pull/7282) *Traffic Ops*: Fixed issue with user getting correctly logged when using an access or bearer token authentication.
- [#7346](https://github.com/apache/trafficcontrol/pull/7346) *Traffic Control Cache Config (t3c)*: Fixed issue with stale lock file when using git to track changes.
- [#7352](https://github.com/apache/trafficcontrol/pull/7352) *Traffic Control Cache Config (t3c)*: Fixed issue with application locking which would allow multiple instances of `t3c apply` to run concurrently.
- [#6775](https://github.com/apache/trafficcontrol/issues/6775) *Traffic Ops*: Invalid "orgServerFqdn" in Delivery Service creation/update causes Internal Server Error.
- [#6695](https://github.com/apache/trafficcontrol/issues/6695) *Traffic Control Cache Config (t3c)*: Directory creation was erroneously reporting an error when actually succeeding.
- [#7411](https://github.com/apache/trafficcontrol/pull/7411) *Traffic Control Cache Config (t3c)*: Fixed issue with wrong parent ordering with MSO non-topology delivery services.
- [#7425](https://github.com/apache/trafficcontrol/pull/7425) *Traffic Control Cache Config (t3c)*: Fixed issue with layered profile iteration being done in the wrong order.
- [#6385](https://github.com/apache/trafficcontrol/issues/6385) *Traffic Ops*: Reserved consistentHashQueryParameters cause internal server error.
- [#7471](https://github.com/apache/trafficcontrol/pull/7471) *Traffic Control Cache Config (t3c)*: Fixed issue with MSO non topo origins from multiple cache groups.
- [#7590](https://github.com/apache/trafficcontrol/issues/7590) *Traffic Control Cache Config (t3c)*: Fixed issue with git detected dubious ownership in repository.
- [#7628](https://github.com/apache/trafficcontrol/pull/7628) *Traffic Ops*: Fixed an issue where certificate chain validation failed based on leading or trailing whitespace.
- [#7688](https://github.com/apache/trafficcontrol/pull/7688) *Traffic Ops*: Fixed secured parameters being visible when role has proper permissions.
- [#7697](https://github.com/apache/trafficcontrol/pull/7697) *Traffic Ops*: Fixed `iloPassword` and `xmppPassword` checking for priv-level instead of using permissions.
- [#7817](https://github.com/apache/trafficcontrol/pull/7817) *Traffic Control Cache Config (t3c)* fixed issue that would cause null ptr panic on client fallback.
- [#7866](https://github.com/apache/trafficcontrol/pull/7866) *Traffic Control Cache Config (t3c)* fixed rpm db check to work with rocky linux 9

### Removed
- [#7808](https://github.com/apache/trafficcontrol/pull/7808) *Traffic Router*: Set SOA `minimum` field to a custom value defined in the `tld.soa.minimum` param, and remove the previously added `dns.negative.caching.ttl` property.
- [#7804](https://github.com/apache/trafficcontrol/pull/7804) Removed unneeded V5 client methods for `deliveryServiceRequiredcapabilities`.
- [#7271](https://github.com/apache/trafficcontrol/pull/7271) Removed components in `infrastructre/docker/`, not in use as cdn-in-a-box performs the same functionality.
- [#7271](https://github.com/apache/trafficcontrol/pull/7271) Removed `misc/jira_github_issue_import.py`, the project does not use JIRA.
- [#7271](https://github.com/apache/trafficcontrol/pull/7271) Removed `traffic_ops/install/bin/convert_profile/`, this script is outdated and is for use on an EOL ATS version.
- [#7271](https://github.com/apache/trafficcontrol/pull/7271) Removed `traffic_ops/install/bin/install_go.sh`, `traffic_ops/install/bin/todb_bootstrap.sh` and `traffic_ops/install/bin/install_goose.sh` are no longer in use.
- [#7829](https://github.com/apache/trafficcontrol/pull/7829) Removed `cache-config/supermicro_udev_mapper.pl` and `traffic_ops_ort.pl` and any references

## [7.0.0] - 2022-07-19
### Added
- [Traffic Portal] Added Layered Profile feature to /servers/
- [#6448](https://github.com/apache/trafficcontrol/issues/6448) Added `status` and `lastPoll` fields to the `publish/CrStates` endpoint of Traffic Monitor (TM).
- Added back to the health-client the `status` field logging with the addition of the filed to `publish/CrStates`
- Added a new Traffic Ops endpoint to `GET` capacity and telemetry data for CDNi integration.
- Added SOA (Service Oriented Architecture) capability to CDN-In-A-Box.
- Added a Traffic Ops endpoints to `PUT` a requested configuration change for a full configuration or per host and an endpoint to approve or deny the request.
- Traffic Monitor config option `distributed_polling` which enables the ability for Traffic Monitor to poll a subset of the CDN and divide into "local peer groups" and "distributed peer groups". Traffic Monitors in the same group are local peers, while Traffic Monitors in other groups are distributed peers. Each TM group polls the same set of cachegroups and gets availability data for the other cachegroups from other TM groups. This allows each TM to be responsible for polling a subset of the CDN while still having a full view of CDN availability. In order to use this, `stat_polling` must be disabled.
- Added support for a new Traffic Ops GLOBAL profile parameter -- `tm_query_status_override` -- to override which status of Traffic Monitors to query (default: ONLINE).
- Traffic Ops: added new `cdn.conf` option -- `user_cache_refresh_interval_sec` -- which enables an in-memory users cache to improve performance. Default: 0 (disabled).
- Traffic Ops: added new `cdn.conf` option -- `server_update_status_cache_refresh_interval_sec` -- which enables an in-memory server update status cache to improve performance. Default: 0 (disabled).
- Traffic Router: Add support for `file`-protocol URLs for the `geolocation.polling.url` for the Geolocation database.
- Replaces all Traffic Portal Tenant select boxes with a novel tree select box [#6427](https://github.com/apache/trafficcontrol/issues/6427).
- Traffic Monitor: Add support for `access.log` to TM.
- Added functionality for login to provide a Bearer token and for that token to be later used for authorization.
- [Traffic Portal] Added the ability for users to view Delivery Service Requests corresponding to individual Delivery Services in TP.
- [Traffic Ops] Added support for backend configurations so that Traffic Ops can act as a reverse proxy for these services [#6754](https://github.com/apache/trafficcontrol/pull/6754).
- Added functionality for CDN locks, so that they can be shared amongst a list of specified usernames.
- [Traffic Ops | Traffic Go Clients | T3C] Add additional timestamp fields to server for queuing and dequeueing config and revalidate updates.
- Added layered profile feature to 4.0 for `GET` /servers/, `POST` /servers/, `PUT` /servers/{id} and `DELETE` /servers/{id}.
- Added a Traffic Ops endpoint and Traffic Portal page to view all CDNi configuration update requests and approve or deny.
- Added layered profile feature to 4.0 for `GET` /deliveryservices/{id}/servers/ and /deliveryservices/{id}/servers/eligible.
- Change to t3c regex_revalidate so that STALE is no longer explicitly added for default revalidate rule for ATS version backwards compatibility.
- Change to t3c diff to flag a config file for replacement if owner/group settings are not `ats` [#6879](https://github.com/apache/trafficcontrol/issues/6879).
- t3c now looks in the executable dir path for t3c- utilities
- Added support for parent.config markdown/retry DS parameters using first./inner./last. prefixes.  mso. and <null> prefixes should be deprecated.
- Add new __REGEX_REMAP_DIRECTIVE__ support to raw remap text to allow moving the regex_remap placement.
- t3c change `t3c diff` call to `t3c-diff` to fix a performance regression.
- Added a sub-app t3c-tail to tail diags.log and capture output when t3c reloads and restarts trafficserver

### Fixed
- Fixed TO to default route ID to 0, if it is not present in the request context.
- [#6291](https://github.com/apache/trafficcontrol/issues/6291) Prevent Traffic Ops from modifying and/or deleting reserved statuses.
- Update traffic\_portal dependencies to mitigate `npm audit` issues.
- Fixed a cdn-in-a-box build issue when using `RHEL_VERSION=7`
- `dequeueing` server updates should not require checking for cdn locks.
- Fixed Traffic Ops ignoring the configured database port value, which was prohibiting the use of anything other than port 5432 (the PostgreSQL default)
- [#6580](https://github.com/apache/trafficcontrol/issues/6580) Fixed cache config generation remap.config targets for MID-type servers in a Topology with other caches as parents and HTTPS origins.
- Traffic Router: fixed a null pointer exception that caused snapshots to be rejected if a topology cachegroup did not have any online/reported/admin\_down caches
- [#6271](https://github.com/apache/trafficcontrol/issues/6271) `api/{{version}/deliveryservices/{id}/health` returns no info if the delivery service uses a topology.
- [#6549](https://github.com/apache/trafficcontrol/issues/6549) Fixed internal server error while deleting a delivery service created from a DSR (Traafic Ops).
- [#6538](https://github.com/apache/trafficcontrol/pull/6538) Fixed the incorrect use of secure.port on TrafficRouter and corrected to the httpsPort value from the TR server configuration.
- [#6562](https://github.com/apache/trafficcontrol/pull/6562) Fixed incorrect template in Ansible dataset loader role when fallbackToClosest is defined.
- [#6590](https://github.com/apache/trafficcontrol/pull/6590) Python client: Corrected parameter name in decorator for get\_parameters\_by\_profile\_id
- [#6368](https://github.com/apache/trafficcontrol/pull/6368) Fixed validation response message from `/acme_accounts`
- [#6603](https://github.com/apache/trafficcontrol/issues/6603) Fixed users with "admin" "Priv Level" not having Permission to view or delete DNSSEC keys.
- Fixed Traffic Router to handle aggressive NSEC correctly.
- [#6907](https://github.com/apache/trafficcontrol/issues/6907) Fixed Traffic Ops to return the correct server structure (based on the API version) upon a server deletion.
- [#6626](https://github.com/apache/trafficcontrol/pull/6626) Fixed t3c Capabilities request failure issue which could result in malformed config.
- [#6370](https://github.com/apache/trafficcontrol/pull/6370) Fixed docs for `POST` and response code for `PUT` to `/acme_accounts` endpoint
- Only `operations` and `admin` roles should have the `DELIVERY-SERVICE:UPDATE` permission.
- [#6369](https://github.com/apache/trafficcontrol/pull/6369) Fixed `/acme_accounts` endpoint to validate email and URL fields
- Fixed searching of the ds parameter merge_parent_groups slice.
- [#6806](https://github.com/apache/trafficcontrol/issues/6806) t3c calculates max_origin_connections incorrectly for topology-based delivery services
- [#6944](https://github.com/apache/trafficcontrol/issues/6944) Fixed cache config generation for ATS 9 sni.yaml from disable_h2 to http2 directive. ATS 9 documents disable_h2, but it doesn't seem to work.
- Fixed TO API `PUT /servers/:id/status` to only queue updates on the same CDN as the updated server
- t3c-generate fix for combining remapconfig and cachekeyconfig parameters for MakeRemapDotConfig call.
- [#6780](https://github.com/apache/trafficcontrol/issues/6780) Fixed t3c to use secondary parents when there are no primary parents available.
- Correction where using the placeholder `__HOSTNAME__` in "unknown" files (others than the defaults ones), was being replaced by the full FQDN instead of the shot hostname.
- [#6800](https://github.com/apache/trafficcontrol/issues/6800) Fixed incorrect error message for `/server/details` associated with query parameters.
- [#6712](https://github.com/apache/trafficcontrol/issues/6712) - Fixed error when loading the Traffic Vault schema from `create_tables.sql` more than once.
- [#6883](https://github.com/apache/trafficcontrol/issues/6883) Fix t3c cache to invalidate on version change
- [#6834](https://github.com/apache/trafficcontrol/issues/6834) - In API 4.0, fixed `GET` for `/servers` to display all profiles irrespective of the index position. Also, replaced query param `profileId` with `profileName`.
- [#6299](https://github.com/apache/trafficcontrol/issues/6299) User representations don't match
- [#6896](https://github.com/apache/trafficcontrol/issues/6896) Fixed the `POST api/cachegroups/id/queue_updates` endpoint so that it doesn't give an internal server error anymore.
- [#6994](https://github.com/apache/trafficcontrol/issues/6994) Fixed the Health Client to not crash if parent.config has a blank line.
- [#6933](https://github.com/apache/trafficcontrol/issues/6933) Fixed tc-health-client to handle credentials files with special characters in variables
- [#6776](https://github.com/apache/trafficcontrol/issues/6776) User properties only required sometimes
- Fixed TO API `GET /deliveryservicesserver` causing error when an IMS request is made with the `cdn` and `maxRevalDurationDays` parameters set.
- [#6792](https://github.com/apache/trafficcontrol/issues/6792) Remove extraneous field from Topologies and Server Capability POST/PUT.
- [#6795](https://github.com/apache/trafficcontrol/issues/6795) Removed an unnecessary response wrapper object from being returned in a POST to the federation resolvers endpoint.

### Removed
- Remove `client.steering.forced.diversity` feature flag(profile parameter) from Traffic Router (TR). Client steering responses now have cache diversity by default.
- Remove traffic\_portal dependencies to mitigate `npm audit` issues, specifically `grunt-concurrent`, `grunt-contrib-concat`, `grunt-contrib-cssmin`, `grunt-contrib-jsmin`, `grunt-contrib-uglify`, `grunt-contrib-htmlmin`, `grunt-newer`, and `grunt-wiredep`
- Replace `forever` with `pm2` for process management of the traffic portal node server to remediate security issues.
- Removed the Traffic Monitor `peer_polling_protocol` option. Traffic Monitor now just uses hostnames to request peer states, which can be handled via IPv4 or IPv6 depending on the underlying IP version in use.
- Dropped CentOS 8 support
- The `/servers/details` endpoint of the Traffic Ops API has been dropped in version 4.0, and marked deprecated in earlier versions.
- Remove Traffic Ops API version 2

### Changed
- [#6694](https://github.com/apache/trafficcontrol/issues/6694) Traffic Stats now uses the TO API 3.0
- [#6654](https://github.com/apache/trafficcontrol/issues/6654) Traffic Monitor now uses the TO API 4.0 by default and falls back to 3.1
- Added Rocky Linux 8 support
- Traffic Monitors now peer with other Traffic Monitors of the same status (e.g. ONLINE with ONLINE, OFFLINE with OFFLINE), instead of all peering with ONLINE.
- Changed the Traffic Ops user last_authenticated update query to only update once per minute to avoid row-locking when the same user logs in frequently.
- Added new fields to the monitoring.json snapshot and made Traffic Monitor prefer data in monitoring.json to the CRConfig snapshot
- Changed the default Traffic Ops API version requsted by Traffic Router from 2.0 to 3.1
- Added permissions to the role form in traffic portal
- Updated the Cache Stats Traffic Portal page to use a more performant AG-Grid-based table.
- Updated the CDNs Traffic Portal page to use a more performant AG-Grid-based table.
- Updated the Profiles Traffic Portal page to use a more performant AG-Grid-based table.
- Updated Go version to 1.18
- Removed the unused `deliveryservice_tmuser` table from Traffic Ops database
- Adds updates to the trafficcontrol-health-client to, use new ATS Host status formats, detect and use proper
  traffic_ctl commands, and adds new markup-poll-threshold config.
- Traffic Monitor now defaults to 100 historical "CRConfig" Snapshots stored internally if not specified in configuration (previous default was 20,000)
- Updated Traffic Router dependencies:
  - commons-io: 2.0.1 -> 2.11.0
  - commons-codec: 1.6 -> 1.15
  - guava: 18.0 -> 31.1-jre
  - async-http-client: 2.12.1 -> 2.12.3
  - spring: 5.2.20.RELEASE -> 5.3.20
- `TRAFFIC_ROUTER`-type Profiles no longer need to have names that match any kind of pattern (e.g. `CCR_.*`)
- [#4351](https://github.com/apache/trafficcontrol/issues/4351) Updated message to an informative one when deleting a delivery service.
- Updated Grove to use the TO API v3 client library
- Updated Ansible Roles to use Traffic Ops API v3
- Update Go version to 1.19

## [6.1.0] - 2022-01-18
### Added
- Added permission based roles for better access control.
- [#5674](https://github.com/apache/trafficcontrol/issues/5674) Added new query parameters `cdn` and `maxRevalDurationDays` to the `GET /api/x/jobs` Traffic Ops API to filter by CDN name and within the start_time window defined by the `maxRevalDurationDays` GLOBAL profile parameter, respectively.
- Added a new Traffic Ops cdn.conf option -- `disable_auto_cert_deletion` -- in order to optionally prevent the automatic deletion of certificates for delivery services that no longer exist whenever a CDN snapshot is taken.
- [#6034](https://github.com/apache/trafficcontrol/issues/6034) Added new query parameter `cdn` to the `GET /api/x/deliveryserviceserver` Traffic Ops API to filter by CDN name
- Added a new Traffic Monitor configuration option -- `short_hostname_override` -- to traffic_monitor.cfg to allow overriding the system hostname that Traffic Monitor uses.
- Added a new Traffic Monitor configuration option -- `stat_polling` (default: true) -- to traffic_monitor.cfg to disable stat polling.
- A new Traffic Portal server command-line option `-c` to specify a configuration file, and the ability to set `log: null` to log to stdout (consult documentation for details).
- Multiple improvements to Ansible roles as discussed at ApacheCon 2021
- SANs information to the SSL key endpoint and Traffic Portal page.
- Added definition for `heartbeat.polling.interval` for CDN Traffic Monitor config in API documentation.
- New `pkg` script options, `-h`, `-s`, `-S`, and `-L`.
- Added `Invalidation Type` (REFRESH or REFETCH) for invalidating content to Traffic Portal.
- cache config t3c-apply retrying when another t3c-apply is running.
- IMS warnings to Content Invalidation requests in Traffic Portal and documentation.
- [#6032](https://github.com/apache/trafficcontrol/issues/6032) Add t3c setting mode 0600 for secure files
- [#6405](https://github.com/apache/trafficcontrol/issues/6405) Added cache config version to all t3c apps and config file headers
- Traffic Vault: Added additional flag to TV Riak (Deprecated) Util
- Added Traffic Vault Postgres columns, a Traffic Ops API endpoint, and Traffic Portal page to show SSL certificate expiration information.
- Added cache config `__CACHEGROUP__` preprocess directive, to allow injecting the local server's cachegroup name into any config file
- Added t3c experimental strategies generation.
- Added support for a DS profile parameter 'LastRawRemapPre' and 'LastRawRemapPost' which allows raw text lines to be pre or post pended to remap.config.
- Added DS parameter 'merge_parent_groups' to allow primary and secondary parents to be merged into the primary parent list in parent.config.

### Fixed
- [#6411](https://github.com/apache/trafficcontrol/pull/6411) Removes invalid 'ALL cdn' options from TP
- Fixed Traffic Router crs/stats to prevent overflow and to correctly record the time used in averages.
- [#6209](https://github.com/apache/trafficcontrol/pull/6209) Updated Traffic Router to use Java 11 to compile and run
- [#5893](https://github.com/apache/trafficcontrol/issues/5893) - A self signed certificate is created when an HTTPS delivery service is created or an HTTP delivery service is updated to HTTPS.
- [#6255](https://github.com/apache/trafficcontrol/issues/6255) - Unreadable Prod Mode CDN Notifications in Traffic Portal
- [#6378](https://github.com/apache/trafficcontrol/issues/6378) - Cannot update or delete Cache Groups with null latitude and longitude.
- Fixed broken `GET /cdns/routing` Traffic Ops API
- [#6259](https://github.com/apache/trafficcontrol/issues/6259) - Traffic Portal No Longer Allows Spaces in Server Object "Router Port Name"
- [#6392](https://github.com/apache/trafficcontrol/issues/6392) - Traffic Ops prevents assigning ORG servers to topology-based delivery services (as well as a number of other valid operations being prohibited by "last server assigned to DS" validations which don't apply to topology-based delivery services)
- [#6175](https://github.com/apache/trafficcontrol/issues/6175) - POST request to /api/4.0/phys_locations accepts mismatch values for regionName.
- Fixed Traffic Monitor parsing stats_over_http output so that multiple stats for the same underlying delivery service (when the delivery service has more than 1 regex) are properly summed together. This makes the resulting data more accurate in addition to fixing the "new stat is lower than last stat" warnings.
- [#6457](https://github.com/apache/trafficcontrol/issues/6457) - Fix broken user registration and password reset, due to the last_authenticated value being null.
- [#6367](https://github.com/apache/trafficcontrol/issues/6367) - Fix PUT `user/current` to work with v4 User Roles and Permissions
- [#6266](https://github.com/apache/trafficcontrol/issues/6266) - Removed postgresql13-devel requirement for traffic_ops
- [#6446](https://github.com/apache/trafficcontrol/issues/6446) - Revert Traffic Router rollover file pattern to the one previously used in `log4j.properties` with Log4j 1.2
- [#5118](https://github.com/apache/trafficcontrol/issues/5118) - Added support for Kafka to Traffic Stats

### Changed
- Updated `t3c` to request less unnecessary deliveryservice-server assignment and invalidation jobs data via new query params supported by Traffic Ops
- [#6179](https://github.com/apache/trafficcontrol/issues/6179) Updated the Traffic Ops rpm to include the `ToDnssecRefresh` binary and make the `trafops_dnssec_refresh` cron job use it
- [#6382](https://github.com/apache/trafficcontrol/issues/6382) Accept Geo Limit Countries as strings and arrays.
- Traffic Portal no longer uses `ruby compass` to compile sass and now uses `dart-sass`.
- Changed Invalidation Jobs throughout (TO, TP, T3C, etc.) to account for the ability to do both REFRESH and REFETCH requests for resources.
- Changed the `maxConnections` value on Traffic Router, to prevent the thundering herd problem (TR).
- The `admin` Role is now always guaranteed to exist, and can't be deleted or modified.
- [#6376](https://github.com/apache/trafficcontrol/issues/6376) Updated TO/TM so that TM doesn't overwrite monitoring snapshot data with CR config snapshot data.
- Updated `t3c-apply` to reduce mutable state in `TrafficOpsReq` struct.
- Updated Golang dependencies
- [#6506](https://github.com/apache/trafficcontrol/pull/6506) - Updated `jackson-databind` and `jackson-annotations` Traffic Router dependencies to version 2.13.1

### Deprecated
- Deprecated the endpoints and docs associated with `/api_capability` and `/capabilities`.
- The use of a seelog configuration file to configure Traffic Stats logging is deprecated, and logging configuration should instead be present in the `logs` property of the Traffic Stats configuration file (refer to documentation for details).

### Removed
- Removed the `user_role` table.
- The `traffic_ops.sh` shell profile no longer sets `GOPATH` or adds its `bin` folder to the `PATH`
- `/capabilities` removed from Traffic Ops API version 4.

## [6.0.2] - 2021-12-17
### Changed
- Updated `log4j` module in Traffic Router from version 1.2.17 to 2.17.0

## [6.0.1] - 2021-11-04
### Added
- [#2770](https://github.com/apache/trafficcontrol/issues/2770) Added validation for httpBypassFqdn as hostname in Traffic Ops

### Fixed
- [#6125](https://github.com/apache/trafficcontrol/issues/6125) - Fix `/cdns/{name}/federations?id=#` to search for CDN.
- [#6285](https://github.com/apache/trafficcontrol/issues/6285) - The Traffic Ops Postinstall script will work in CentOS 7, even if Python 3 is installed
- [#5373](https://github.com/apache/trafficcontrol/issues/5373) - Traffic Monitor logs not consistent
- [#6197](https://github.com/apache/trafficcontrol/issues/6197) - TO `/deliveryservices/:id/routing` makes requests to all TRs instead of by CDN.
- Traffic Ops: Sanitize username before executing LDAP query

### Changed
- [#5927](https://github.com/apache/trafficcontrol/issues/5927) Updated CDN-in-a-Box to not run a Riak container by default but instead only run it if the optional flag is provided.
- Changed the DNSSEC refresh Traffic Ops API to only create a new change log entry if any keys were actually refreshed or an error occurred (in order to reduce changelog noise)

## [6.0.0] - 2021-08-30
### Added
- [#4982](https://github.com/apache/trafficcontrol/issues/4982) Added the ability to support queueing updates by server type and profile
- [#5412](https://github.com/apache/trafficcontrol/issues/5412) Added last authenticated time to user API's (`GET /user/current, GET /users, GET /user?id=`) response payload
- [#5451](https://github.com/apache/trafficcontrol/issues/5451) Added change log count to user API's response payload and query param (username) to logs API
- Added support for CDN locks
- Added support for PostgreSQL as a Traffic Vault backend
- Added the tc-health-client to Trafficcontrol used to manage traffic server parents.
- [#5449](https://github.com/apache/trafficcontrol/issues/5449) The `todb-tests` GitHub action now runs the Traffic Ops DB tests
- Python client: [#5611](https://github.com/apache/trafficcontrol/pull/5611) Added server_detail endpoint
- Ported the Postinstall script to Python. The Perl version has been moved to `install/bin/_postinstall.pl` and has been deprecated, pending removal in a future release.
- CDN-in-a-Box: Generate config files using the Postinstall script
- CDN-in-a-Box: Add Federation with CNAME foo.kabletown.net.
- Apache Traffic Server: [#5627](https://github.com/apache/trafficcontrol/pull/5627) - Added the creation of Centos8 RPMs for Apache Traffic Server
- Traffic Ops/Traffic Portal: [#5479](https://github.com/apache/trafficcontrol/issues/5479) - Added the ability to change a server capability name
- Traffic Ops: [#3577](https://github.com/apache/trafficcontrol/issues/3577) - Added a query param (server host_name or ID) for servercheck API
- Traffic Portal: [#5318](https://github.com/apache/trafficcontrol/issues/5318) - Rename server columns for IPv4 address fields.
- Traffic Portal: [#5361](https://github.com/apache/trafficcontrol/issues/5361) - Added the ability to change the name of a topology.
- Traffic Portal: [#5340](https://github.com/apache/trafficcontrol/issues/5340) - Added the ability to resend a user registration from user screen.
- Traffic Portal: Adds the ability for operations/admin users to create a CDN-level notification.
- Traffic Portal: upgraded delivery service UI tables to use more powerful/performant ag-grid component
- Traffic Router: added new 'dnssec.rrsig.cache.enabled' profile parameter to enable new DNSSEC RRSIG caching functionality. Enabling this greatly reduces CPU usage during the DNSSEC signing process.
- Traffic Router: added new 'strip.special.query.params' profile parameter to enable stripping the 'trred' and 'fakeClientIpAddress' query parameters from responses: [#1065](https://github.com/apache/trafficcontrol/issues/1065)
- [#5316](https://github.com/apache/trafficcontrol/issues/5316) - Add router host names and ports on a per interface basis, rather than a per server basis.
- Traffic Ops: Adds API endpoints to fetch (GET), create (POST) or delete (DELETE) a cdn notification. Create and delete are limited to users with operations or admin role.
- Added ACME certificate renewals and ACME account registration using external account binding
- Added functionality to automatically renew ACME certificates.
- Traffic Ops: [#6069](https://github.com/apache/trafficcontrol/issues/6069) - prevent unassigning all ONLINE ORG servers from an MSO-enabled delivery service
- Added ORT flag to set local.dns bind address from server service addresses
- Added an endpoint for statuses on asynchronous jobs and applied it to the ACME renewal endpoint.
- Added two new cdn.conf options to make Traffic Vault configuration more backend-agnostic: `traffic_vault_backend` and `traffic_vault_config`
- Traffic Ops API version 4.0 - This version is **unstable** meaning that breaking changes can occur at any time - use at your own peril!
- `GET` request method for `/deliveryservices/{{ID}}/assign`
- `GET` request method for `/deliveryservices/{{ID}}/status`
- [#5644](https://github.com/apache/trafficcontrol/issues/5644) ORT config generation: Added ATS9 ip_allow.yaml support, and automatic generation if the server's package Parameter is 9.\*
- t3c: Added option to track config changes in git.
- ORT config generation: Added a rule to ip_allow such that PURGE requests are allowed over localhost
- Added integration to use ACME to generate new SSL certificates.
- Add a Federation to the Ansible Dataset Loader
- Added `GetServersByDeliveryService` method to the TO Go client
- Added asynchronous status to ACME certificate generation.
- Added per Delivery Service HTTP/2 and TLS Versions support, via ssl_server_name.yaml and sni.yaml. See overview/delivery_services and t3c docs.
- Added headers to Traffic Portal, Traffic Ops, and Traffic Monitor to opt out of tracking users via Google FLoC.
- Add logging scope for logging.yaml generation for ATS 9 support
- `DELETE` request method for `deliveryservices/xmlId/{name}/urlkeys` and `deliveryservices/{id}/urlkeys`.
- t3c now uses separate apps, full run syntax changed to `t3c apply ...`, moved to cache-config and RPM changed to trafficcontrol-cache-config. See cache-config README.md.
- t3c: bug fix to consider plugin config files for reloading remap.config
- t3c: add flag to wait for parents in syncds mode
- t3c: Change syncds so that it only warns on package version mismatch.
- atstccfg: add ##REFETCH## support to regex_revalidate.config processing.
- Traffic Router: Added `svc="..."` field to request logging output.
- Added t3c caching Traffic Ops data and using If-Modified-Since to avoid unnecessary requests.
- Added t3c --no-outgoing-ip flags.
- Added a Traffic Monitor integration test framework.
- Added `traffic_ops/app/db/traffic_vault_migrate` to help with migrating Traffic Ops Traffic Vault backends
- Added a tool at `/traffic_ops/app/db/reencrypt` to re-encrypt the data in the Postgres Traffic Vault with a new key.
- Enhanced ort integration test for reload states
- Added a new field to Delivery Services - `tlsVersions` - that explicitly lists the TLS versions that may be used to retrieve their content from Cache Servers.
- Added support for DS plugin parameters for cachekey, slice, cache_range_requests, background_fetch, url_sig  as remap.config parameters.
- Updated T3C changes in Ansible playbooks
- Updated all endpoints in infrastructure code to use API version 2.0

### Fixed
- [#5690](https://github.com/apache/trafficcontrol/issues/5690) - Fixed github action for added/modified db migration file.
- [#2471](https://github.com/apache/trafficcontrol/issues/2471) - A PR check to ensure added db migration file is the latest.
- [#6129](https://github.com/apache/trafficcontrol/issues/6129) - Traffic Monitor start doesn't recover when Traffic Ops is unavailable
- [#5609](https://github.com/apache/trafficcontrol/issues/5609) - Fixed GET /servercheck filter for an extra query param.
- [#5954](https://github.com/apache/trafficcontrol/issues/5954) - Traffic Ops HTTP response write errors are ignored
- [#6048](https://github.com/apache/trafficcontrol/issues/6048) - TM sometimes sets cachegroups to unavailable even though all caches are online
- [#6104](https://github.com/apache/trafficcontrol/issues/6104) - PUT /api/x/federations only respects first item in request payload
- [#5288](https://github.com/apache/trafficcontrol/issues/5288) - Fixed the ability to create and update a server with MTU value >= 1280.
- [#5284](https://github.com/apache/trafficcontrol/issues/5284) - Fixed error message when creating a server with non-existent profile
- Fixed a NullPointerException in TR when a client passes a null SNI hostname in a TLS request
- Fixed a logging bug in Traffic Monitor where it wouldn't log errors in certain cases where a backup file could be used instead. Also, Traffic Monitor now rejects monitoring snapshots that have no delivery services.
- [#5739](https://github.com/apache/trafficcontrol/issues/5739) - Prevent looping in case of a failed login attempt
- [#5407](https://github.com/apache/trafficcontrol/issues/5407) - Make sure that you cannot add two servers with identical content
- [#2881](https://github.com/apache/trafficcontrol/issues/2881) - Some API endpoints have incorrect Content-Types
- [#5863](https://github.com/apache/trafficcontrol/issues/5863) - Traffic Monitor logs warnings to `log_location_info` instead of `log_location_warning`
- [#5492](https://github.com/apache/trafficcontrol/issues/5942) - Traffic Stats does not failover to another Traffic Monitor when one stops responding
- [#5363](https://github.com/apache/trafficcontrol/issues/5363) - Postgresql version changeable by env variable
- [#5405](https://github.com/apache/trafficcontrol/issues/5405) - Prevent Tenant update from choosing child as new parent
- [#5384](https://github.com/apache/trafficcontrol/issues/5384) - New grids will now properly remember the current page number.
- [#5548](https://github.com/apache/trafficcontrol/issues/5548) - Don't return a `403 Forbidden` when the user tries to get servers of a non-existent DS using `GET /servers?dsId={{nonexistent DS ID}}`
- [#5732](https://github.com/apache/trafficcontrol/issues/5732) - TO API POST /cdns/dnsseckeys/generate times out with large numbers of delivery services
- [#5902](https://github.com/apache/trafficcontrol/issues/5902) - Fixed issue where the TO API wouldn't properly query all SSL certificates from Riak.
- Fixed server creation through legacy API versions to default `monitor` to `true`.
- Fixed t3c to generate topology parents correctly for parents with the Type MID+ or EDGE+ versus just the literal. Naming cache types to not be exactly 'EDGE' or 'MID' is still discouraged and not guaranteed to work, but it's unfortunately somewhat common, so this fixes it in one particular case.
- [#5965](https://github.com/apache/trafficcontrol/issues/5965) - Fixed Traffic Ops /deliveryserviceservers If-Modified-Since requests.
- Fixed t3c to create config files and directories as ats.ats
- Fixed t3c-apply service restart and ats config reload logic.
- Reduced TR dns.max-threads ansible default from 10000 to 100.
- Converted TP Delivery Service Servers Assignment table to ag-grid
- Converted TP Cache Checks table to ag-grid
- [#5981](https://github.com/apache/trafficcontrol/issues/5891) - `/deliveryservices/{{ID}}/safe` returns incorrect response for the requested API version
- [#5984](https://github.com/apache/trafficcontrol/issues/5894) - `/servers/{{ID}}/deliveryservices` returns incorrect response for the requested API version
- [#6027](https://github.com/apache/trafficcontrol/issues/6027) - Collapsed DB migrations
- [#6091](https://github.com/apache/trafficcontrol/issues/6091) - Fixed cache config of internal cache communication for https origins
- [#6066](https://github.com/apache/trafficcontrol/issues/6066) - Fixed missing/incorrect indices on some tables
- [#6169](https://github.com/apache/trafficcontrol/issues/6169) - Fixed t3c-update not updating server status when a fallback to a previous Traffic Ops API version occurred
- [#5576](https://github.com/apache/trafficcontrol/issues/5576) - Inconsistent Profile Name restrictions
- [#6327](https://github.com/apache/trafficcontrol/issues/6327) - Fixed cache config to invalidate its cache if the Server's Profile or CDN changes
- [#6174](https://github.com/apache/trafficcontrol/issues/6174) - Fixed t3c-apply with no hostname failing if the OS hostname returns a full FQDN
- Fixed Federations IMS so TR federations watcher will get updates.
- [#5129](https://github.com/apache/trafficcontrol/issues/5129) - Updated TM so that it returns a 404 if the endpoint is not supported.
- [#5992](https://github.com/apache/trafficcontrol/issues/5992) - Updated Traffic Router Integration tests to use a mock Traffic Monitor and Traffic Ops server
- [#6093](https://github.com/apache/trafficcontrol/issues/6093) - Fixed Let's Encrypt to work for delivery services where the domain does not contain the XMLID.

### Changed
- Migrated completely off of bower in favor of npm
- Updated the Traffic Ops Python client to 3.0
- Updated Flot libraries to supported versions
- [apache/trafficcontrol](https://github.com/apache/trafficcontrol) is now a Go module
- Updated Traffic Ops supported database version from PostgreSQL 9.6 to 13.2
- [#3342](https://github.com/apache/trafficcontrol/issues/3342) - Updated the [`db/admin`](https://traffic-control-cdn.readthedocs.io/en/latest/development/traffic_ops.html#database-management) tool to use [Migrate](https://github.com/golang-migrate/migrate) instead of Goose and converted the migrations to Migrate format (split up/down for each migration into separate files)
- Set Traffic Router to also accept TLSv1.3 protocols by default in server.xml
- Disabled TLSv1.1 for Traffic Router in Ansible role by default
- Updated the Traffic Monitor Ansible role to set `serve_write_timeout_ms` to `20000` by default because 10 seconds can be too short for relatively large CDNs.
- Refactored the Traffic Ops - Traffic Vault integration to more easily support the development of new Traffic Vault backends
- Changed the Traffic Router package structure from com.comcast.cdn.\* to org.apache.\*
- Updated Apache Tomcat from 8.5.63 to 9.0.43
- Improved the DNSSEC refresh Traffic Ops API (`/cdns/dnsseckeys/refresh`). As of TO API v4, its method is `PUT` instead of `GET`, its response format was changed to return an alert instead of a string-based response, it returns a 202 instead of a 200, and it now works with the `async_status` API in order for the client to check the status of the async job: [#3054](https://github.com/apache/trafficcontrol/issues/3054)
- Delivery Service Requests now keep a record of the changes they make.
- Changed the `goose` provider to the maintained fork [`github.com/kevinburke/goose`](https://github.com/kevinburke/goose)
- The format of the `/servers/{{host name}}/update_status` Traffic Ops API endpoint has been changed to use a top-level `response` property, in keeping with (most of) the rest of the API.
- The API v4 Traffic Ops Go client has been overhauled compared to its predecessors to have a consistent call signature that allows passing query string parameters and HTTP headers to any client method.
- Updated BouncyCastle libraries in Traffic Router to v1.68.
- lib/go-atscfg Make funcs to take Opts, to reduce future breaking changes.
- CDN in a Box now uses `t3c` for cache configuration.
- CDN in a Box now uses Apache Traffic Server 8.1.
- Customer names in payloads sent to the `/deliveryservices/request` Traffic Ops API endpoint can no longer contain characters besides alphanumerics, @, !, #, $, %, ^, &amp;, *, (, ), [, ], '.', ' ', and '-'. This fixes a vulnerability that allowed email content injection.
- Go version 1.17 is used to compile Traffic Ops, T3C, Traffic Monitor, Traffic Stats, and Grove.

### Deprecated
- The Riak Traffic Vault backend is now deprecated and its support may be removed in a future release. It is highly recommended to use the new PostgreSQL backend instead.
- The `riak.conf` config file and its corresponding `--riakcfg` option in `traffic_ops_golang` have been deprecated. Please use `"traffic_vault_backend": "riak"` and `"traffic_vault_config"` (with the existing contents of riak.conf) instead.
- The Traffic Ops API route `GET /api/{version}/vault/bucket/{bucket}/key/{key}/values` has been deprecated and will no longer be available as of Traffic Ops API v4
- The Traffic Ops API route `POST /api/{version}/deliveryservices/request` has been deprecated and will no longer be available as of Traffic Ops API v4
- The Traffic Ops API routes `GET /api/{version}/cachegroupparameters`, `POST /api/{version}/cachegroupparameters`, `GET /api/{version}/cachegroups/{id}/parameters`, and `DELETE /api/{version}/cachegroupparameters/{cachegroupID}/{parameterId}` have been deprecated and will no longer be available as of Traffic Ops API v4
- The `riak_port` option in cdn.conf is now deprecated. Please use the `"port"` field in `traffic_vault_config` instead.
- The `traffic_ops_ort.pl` tool has been deprecated in favor of `t3c`, and will be removed in the next major version.
- With the release of ATC v6.0, major API version 2 is now deprecated, subject to removal with the next ATC major version release, at the earliest.

### Removed
- Removed the unused `backend_max_connections` option from `cdn.conf`.
- Removed the unused `http_poll_no_sleep`, `max_stat_history`, `max_health_history`, `cache_health_polling_interval_ms`, `cache_stat_polling_interval_ms`, and `peer_polling_interval_ms` Traffic Monitor config options.
- Removed the `Long Description 2` and `Long Description 3` fields of `DeliveryService` from the UI, and changed the backend so that routes corresponding API 4.0 and above no longer accept or return these fields.
- The Perl implementation of Traffic Ops has been stripped out, along with the Go implementation's "fall-back to Perl" behavior.
- Traffic Ops no longer includes an `app/public` directory, as the static webserver has been removed along with the Perl Traffic Ops implementation. Traffic Ops also no longer attempts to download MaxMind GeoIP City databases when running the Traffic Ops Postinstall script.
- The `compare` tool stack has been removed, as it no longer serves a purpose.
- Removed the Perl-only `cdn.conf` option `geniso.iso_root_path`
- t3c dispersion flags. These flags existed in ort.pl and t3c, but the feature has been removed in t3c-apply. The t3c run is fast enough now, there's no value or need in internal logic, operators can easily use shell pipelines to randomly sleep before running if necessary.
- Traffic Ops API version 1


## [5.1.2] - 2021-05-17
### Fixed
- Fixed the return error for GET api `cdns/routing` to avoid incorrect success response.
- [#5712](https://github.com/apache/trafficcontrol/issues/5712) - Ensure that 5.x Traffic Stats is compatible with 5.x Traffic Monitor and 5.x Traffic Ops, and that it doesn't log all 0's for `cache_stats`
- Fixed ORT being unable to update URLSIG keys for Delivery Services
- Fixed ORT service category header rewrite for mids and topologies.
- Fixed an issue where Traffic Ops becoming unavailable caused Traffic Monitor to segfault and crash
- [#5754](https://github.com/apache/trafficcontrol/issues/5754) - Ensure Health Threshold Parameters use legacy format for legacy Monitoring Config handler
- [#5695](https://github.com/apache/trafficcontrol/issues/5695) - Ensure vitals are calculated only against monitored interfaces
- Fixed Traffic Monitor to report `ONLINE` caches as available.
- [#5744](https://github.com/apache/trafficcontrol/issues/5744) - Sort TM Delivery Service States page by DS name
- [#5724](https://github.com/apache/trafficcontrol/issues/5724) - Set XMPPID to hostname if the server had none, don't error on server update when XMPPID is empty

## [5.1.1] - 2021-03-19
### Added
- Atscfg: Added a rule to ip_allow such that PURGE requests are allowed over localhost

### Fixed
- [#5565](https://github.com/apache/trafficcontrol/issues/5565) - TO GET /caches/stats panic converting string to uint64
- [#5558](https://github.com/apache/trafficcontrol/issues/5558) - Fixed `TM UI` and `/api/cache-statuses` to report aggregate `bandwidth_kbps` correctly.
- [#5192](https://github.com/apache/trafficcontrol/issues/5192) - Fixed TO log warnings when generating snapshots for topology-based delivery services.
- Fixed Invalid TS logrotate configuration permissions causing TS logs to be ignored by logrotate.
- [#5604](https://github.com/apache/trafficcontrol/issues/5604) - traffic_monitor.log is no longer truncated when restarting Traffic Monitor

## [5.1.0] - 2021-03-11
### Added
- Traffic Ops: added a feature so that the user can specify `maxRequestHeaderBytes` on a per delivery service basis
- Traffic Router: log warnings when requests to Traffic Monitor return a 503 status code
- [#5344](https://github.com/apache/trafficcontrol/issues/5344) - Add a page that addresses migrating from Traffic Ops API v1 for each endpoint
- [#5296](https://github.com/apache/trafficcontrol/issues/5296) - Fixed a bug where users couldn't update any regex in Traffic Ops/ Traffic Portal
- Added API endpoints for ACME accounts
- Traffic Ops: Added validation to ensure that the cachegroups of a delivery services' assigned ORG servers are present in the topology
- Traffic Ops: Added validation to ensure that the `weight` parameter of `parent.config` is a float
- Traffic Ops Client: New Login function with more options, including falling back to previous minor versions. See traffic_ops/v3-client documentation for details.
- Added license files to the RPMs

### Fixed
- [#5288](https://github.com/apache/trafficcontrol/issues/5288) - Fixed the ability to create and update a server with MTU value >= 1280.
- [#1624](https://github.com/apache/trafficcontrol/issues/1624) - Fixed ORT to reload Traffic Server if LUA scripts are added or changed.
- [#5445](https://github.com/apache/trafficcontrol/issues/5445) - When updating a registered user, ignore updates on registration_sent field.
- [#5335](https://github.com/apache/trafficcontrol/issues/5335) - Don't create a change log entry if the delivery service primary origin hasn't changed
- [#5333](https://github.com/apache/trafficcontrol/issues/5333) - Don't create a change log entry for any delivery service consistent hash query params updates
- [#5341](https://github.com/apache/trafficcontrol/issues/5341) - For a DS with existing SSLKeys, fixed HTTP status code from 403 to 400 when updating CDN and Routing Name (in TO) and made CDN and Routing Name fields immutable (in TP).
- [#5192](https://github.com/apache/trafficcontrol/issues/5192) - Fixed TO log warnings when generating snapshots for topology-based delivery services.
- [#5284](https://github.com/apache/trafficcontrol/issues/5284) - Fixed error message when creating a server with non-existent profile
- [#5287](https://github.com/apache/trafficcontrol/issues/5287) - Fixed error message when creating a Cache Group with no typeId
- [#5382](https://github.com/apache/trafficcontrol/issues/5382) - Fixed API documentation and TP helptext for "Max DNS Answers" field with respect to DNS, HTTP, Steering Delivery Service
- [#5396](https://github.com/apache/trafficcontrol/issues/5396) - Return the correct error type if user tries to update the root tenant
- [#5378](https://github.com/apache/trafficcontrol/issues/5378) - Updating a non existent DS should return a 404, instead of a 500
- Fixed a potential Traffic Router race condition that could cause erroneous 503s for CLIENT_STEERING delivery services when loading new steering changes
- [#5195](https://github.com/apache/trafficcontrol/issues/5195) - Correctly show CDN ID in Changelog during Snap
- [#5438](https://github.com/apache/trafficcontrol/issues/5438) - Correctly specify nodejs version requirements in traffic_portal.spec
- Fixed Traffic Router logging unnecessary warnings for IPv6-only caches
- [#5294](https://github.com/apache/trafficcontrol/issues/5294) - TP ag grid tables now properly persist column filters on page refresh.
- [#5295](https://github.com/apache/trafficcontrol/issues/5295) - TP types/servers table now clears all filters instead of just column filters
- [#5407](https://github.com/apache/trafficcontrol/issues/5407) - Make sure that you cannot add two servers with identical content
- [#2881](https://github.com/apache/trafficcontrol/issues/2881) - Some API endpoints have incorrect Content-Types
- [#5311](https://github.com/apache/trafficcontrol/issues/5311) - Better TO log messages when failures calling TM CacheStats
- [#5364](https://github.com/apache/trafficcontrol/issues/5364) - Cascade server deletes to delete corresponding IP addresses and interfaces
- [#5390](https://github.com/apache/trafficcontrol/issues/5390) - Improve the way TO deals with delivery service server assignments
- [#5339](https://github.com/apache/trafficcontrol/issues/5339) - Ensure Changelog entries for SSL key changes
- [#5461](https://github.com/apache/trafficcontrol/issues/5461) - Fixed steering endpoint to be ordered consistently
- [#5395](https://github.com/apache/trafficcontrol/issues/5395) - Added validation to prevent changing the Type any Cache Group that is in use by a Topology
- Fixed an issue with 2020082700000000_server_id_primary_key.sql trying to create multiple primary keys when there are multiple schemas.
- Fix for public schema in 2020062923101648_add_deleted_tables.sql
- Fix for config gen missing max_origin_connections on mids in certain scenarios
- [#5642](https://github.com/apache/trafficcontrol/issues/5642) - Fixed ORT to fall back to previous minor Traffic Ops versions, allowing ORT to be upgraded before Traffic Ops when the minor has changed.
- Moved move_lets_encrypt_to_acme.sql, add_max_request_header_size_delivery_service.sql, and server_interface_ip_address_cascade.sql past last migration in 5.0.0
- [#5505](https://github.com/apache/trafficcontrol/issues/5505) - Make `parent_reval_pending` for servers in a Flexible Topology CDN-specific on `GET /servers/{name}/update_status`
- [#5317](https://github.com/apache/trafficcontrol/issues/5317) - Clicking IP addresses in the servers table no longer navigates to server details page.
- [#5554](https://github.com/apache/trafficcontrol/issues/5554) - TM UI overflows screen width and hides table data

### Changed
- [#5553](https://github.com/apache/trafficcontrol/pull/5553) - Removing Tomcat specific build requirement
- Refactored the Traffic Ops Go client internals so that all public methods have a consistent behavior/implementation
- Pinned external actions used by Documentation Build and TR Unit Tests workflows to commit SHA-1 and the Docker image used by the Weasel workflow to a SHA-256 digest
- Set Traffic Router to only accept TLSv1.1 and TLSv1.2 protocols in server.xml
- Updated Apache Tomcat from 8.5.57 to 8.5.63
- Updated Apache Tomcat Native from 1.2.16 to 1.2.23
- Traffic Portal: [#5394](https://github.com/apache/trafficcontrol/issues/5394) - Converts the tenant table to a tenant tree for usability
- Traffic Portal: upgraded delivery service UI tables to use more powerful/performant ag-grid component

## [5.0.0] - 2020-10-20
### Added
- Traffic Ops Ort: Disabled ntpd verification (ntpd is deprecated in CentOS)
- Traffic Ops Ort: Adds a transliteration of the traffic_ops_ort.pl perl script to the go language. See traffic_ops_ort/t3c/README.md.
- Traffic Ops API v3
- Added an optional readiness check service to cdn-in-a-box that exits successfully when it is able to get a `200 OK` from all delivery services
- Added health checks to Traffic Ops and Traffic Monitor in cdn-in-a-box
- [Flexible Topologies](https://github.com/apache/trafficcontrol/blob/master/blueprints/flexible-topologies.md)
    - Traffic Ops: Added an API 3.0 endpoint, `GET /api/3.0/topologies`, to create, read, update and delete flexible topologies.
    - Traffic Ops: Added an API 3.0 endpoint, `POST /api/3.0/topologies/{name}/queue_update`, to queue or dequeue updates for all servers assigned to the Cachegroups in a given Topology.
    - Traffic Ops: Added new `topology` field to the /api/3.0/deliveryservices APIs
    - Traffic Ops: Added support for `topology` query parameter to `GET /api/3.0/cachegroups` to return all cachegroups used in the given topology.
    - Traffic Ops: Added support for `topology` query parameter to `GET /api/3.0/deliveryservices` to return all delivery services that employ a given topology.
    - Traffic Ops: Added support for `dsId` query parameter for `GET /api/3.0/servers` for topology-based delivery services.
    - Traffic Ops: Excluded ORG-type servers from `GET /api/3.0/servers?dsId=#` for Topology-based Delivery Services unless the ORG server is assigned to that Delivery Service.
    - Traffic Ops: Added support for `topology` query parameter for `GET /api/3.0/servers` to return all servers whose cachegroups are in a given topology.
    - Traffic Ops: Added new topology-based delivery service fields for header rewrites: `firstHeaderRewrite`, `innerHeaderRewrite`, `lastHeaderRewrite`
    - Traffic Ops: Added validation to prohibit assigning caches to topology-based delivery services
    - Traffic Ops: Added validation to prohibit removing a capability from a server if no other server in the same cachegroup can satisfy the required capabilities of the delivery services assigned to it via topologies.
    - Traffic Ops: Added validation to ensure that updated topologies are still valid with respect to the required capabilities of their assigned delivery services.
    - Traffic Ops: Added validation to ensure that at least one server per cachegroup in a delivery service's topology has the delivery service's required capabilities.
    - Traffic Ops: Added validation to ensure that at least one server exists in each cachegroup that is used in a Topology on the `/api/3.0/topologies` endpoint and the `/api/3.0/servers/{{ID}}` endpoint.
    - Traffic Ops: Consider Topologies parentage when queueing or checking server updates
    - ORT: Added Topologies to Config Generation.
    - Traffic Portal: Added the ability to create, read, update and delete flexible topologies.
    - Traffic Portal: Added the ability to assign topologies to delivery services.
    - Traffic Portal: Added the ability to view all delivery services, cache groups and servers associated with a topology.
    - Traffic Portal: Added the ability to define first, inner and last header rewrite values for DNS* and HTTP* delivery services that employ a topology.
    - Traffic Portal: Adds the ability to view all servers utilized by a topology-based delivery service.
    - Traffic Portal: Added topology section to cdn snapshot diff.
    - Added to TP the ability to assign ORG servers to topology-based delivery services
    - Traffic Router: Added support for topology-based delivery services
    - Traffic Monitor: Added the ability to mark topology-based delivery services as available
    - CDN-in-a-Box: Add a second mid to CDN-in-a-Box, add topology `demo1-top`, and make the `demo1` delivery service topology-based
    - Traffic Ops: Added validation to ensure assigned ORG server cachegroups are in the topology when updating a delivery service
- Updated /servers/details to use multiple interfaces in API v3
- Added [Edge Traffic Routing](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_router.html#edge-traffic-routing) feature which allows Traffic Router to localize more DNS record types than just the routing name for DNS delivery services
- Added the ability to speedily build development RPMs from any OS without needing Docker
- Added the ability to perform a quick search, override default pagination size and clear column filters on the Traffic Portal servers table.
- Astats csv support - astats will now respond to `Accept: text/csv` and return a csv formatted stats list
- Updated /deliveryservices/{{ID}}/servers to use multiple interfaces in API v3
- Updated /deliveryservices/{{ID}}/servers/eligible to use multiple interfaces in API v3
- Added the ability to view Hash ID field (aka xmppID) on Traffic Portals' server summary page
- Added the ability to delete invalidation requests in Traffic Portal
- Added the ability to set TLS config provided here: https://golang.org/pkg/crypto/tls/#Config in Traffic Ops
- Added support for the `cachegroupName` query parameter for `GET /api/3.0/servers` in Traffic Ops
- Added an indiciator to the Traffic Monitor UI when using a disk backup of Traffic Ops.
- Added debugging functionality to CDN-in-a-Box for Traffic Stats.
- Added If-Match and If-Unmodified-Since Support in Server and Clients.
- Added debugging functionality to the Traffic Router unit tests runner at [`/traffic_router/tests`](https://github.com/apache/trafficcontrol/tree/master/traffic_router/tests)
- Made the Traffic Router unit tests runner at [`/traffic_router/tests`](https://github.com/apache/trafficcontrol/tree/master/traffic_router/tests) run in Alpine Linux
- Added GitHub Actions workflow for building RPMs and running the CDN-in-a-Box readiness check
- Added the `Status Last Updated` field to servers, and the UI, so that we can see when the last status change took place for a server.
- Added functionality in TR, so that it uses the default miss location of the DS, in case the location(for the  client IP) returned was the default location of the country.
- Added ability to set DNS Listening IPs in dns.properties
- Added Traffic Monitor: Support astats CSV output. Includes http_polling_format configuration option to specify the Accept header sent to stats endpoints. Adds CSV parsing ability (~100% faster than JSON) to the astats plugin
- Added Traffic Monitor: Support stats over http CSV output. Officially supported in ATS 9.0 unless backported by users. Users must also include `system_stats.so` when using stats over http in order to keep all the same functionality (and included stats) that astats_over_http provides.
- Added ability for Traffic Monitor to determine health of cache based on interface data and aggregate data. Using the new `stats_over_http` `health.polling.format` value that allows monitoring of multiple interfaces will first require that *all* Traffic Monitors monitoring the affected cache server be upgraded.
- Added ORT option to try all primaries before falling back to secondary parents, via Delivery Service Profile Parameter "try_all_primaries_before_secondary".
- Traffic Ops, Traffic Ops ORT, Traffic Monitor, Traffic Stats, and Grove are now compiled using Go version 1.15.
- Added `--traffic_ops_insecure=<0|1>` optional option to traffic_ops_ort.pl
- Added User-Agent string to Traffic Router log output.
- Added default sort logic to GET API calls using Read()
- Traffic Ops: added validation for assigning ORG servers to topology-based delivery services
- Added locationByDeepCoverageZone to the `crs/stats/ip/{ip}` endpoint in the Traffic Router API
- Traffic Ops: added validation for topology updates and server updates/deletions to ensure that topologies have at least one server per cachegroup in each CDN of any assigned delivery services
- Traffic Ops: added validation for delivery service updates to ensure that topologies have at least one server per cachegroup in each CDN of any assigned delivery services
- Traffic Ops: added a feature to get delivery services filtered by the `active` flag
- Traffic Portal: upgraded change log UI table to use more powerful/performant ag-grid component
- Traffic Portal: change log days are now configurable in traffic_portal_properties.json (default is 7 days) and can be overridden by the user in TP
- [#5319](https://github.com/apache/trafficcontrol/issues/5319) - Added support for building RPMs that target CentOS 8
- [#5360](https://github.com/apache/trafficcontrol/issues/5360) - Adds the ability to clone a topology

### Fixed
- Fixed #5188 - DSR (delivery service request) incorrectly marked as complete and error message not displaying when DSR fulfilled and DS update fails in Traffic Portal. [Related Github issue](https://github.com/apache/trafficcontrol/issues/5188)
- Fixed #3455 - Alphabetically sorting CDN Read API call [Related Github issue](https://github.com/apache/trafficcontrol/issues/3455)
- Fixed #5010 - Fixed Reference urls for Cache Config on Delivery service pages (HTTP, DNS) in Traffic Portal. [Related Github issue](https://github.com/apache/trafficcontrol/issues/5010)
- Fixed #5147 - GET /servers?dsId={id} should only return mid servers (in addition to edge servers) for the cdn of the delivery service if the mid tier is employed. [Related github issue](https://github.com/apache/trafficcontrol/issues/5147)
- Fixed #4981 - Cannot create routing regular expression with a blank pattern param in Delivery Service [Related github issue](https://github.com/apache/trafficcontrol/issues/4981)
- Fixed #4979 - Returns a Bad Request error during server creation with missing profileId [Related github issue](https://github.com/apache/trafficcontrol/issues/4979)
- Fixed #4237 - Do not return an internal server error when delivery service's capacity is zero. [Related github issue](https://github.com/apache/trafficcontrol/issues/4237)
- Fixed #2712 - Invalid TM logrotate configuration permissions causing TM logs to be ignored by logrotate. [Related github issue](https://github.com/apache/trafficcontrol/issues/2712)
- Fixed #3400 - Allow "0" as a TTL value for Static DNS entries [Related github issue](https://github.com/apache/trafficcontrol/issues/3400)
- Fixed #5050 - Allows the TP administrator to name a TP instance (production, staging, etc) and flag whether it is production or not in traffic_portal_properties.json [Related github issue](https://github.com/apache/trafficcontrol/issues/5050)
- Fixed #4743 - Validate absolute DNS name requirement on Static DNS entry for CNAME type [Related github issue](https://github.com/apache/trafficcontrol/issues/4743)
- Fixed #4848 - `GET /api/x/cdns/capacity` gives back 500, with the message `capacity was zero`
- Fixed #2156 - Renaming a host in TC, does not impact xmpp_id and thereby hashid [Related github issue](https://github.com/apache/trafficcontrol/issues/2156)
- Fixed #5038 - Adds UI warning when server interface IP CIDR is too large [Related github issue](https://github.com/apache/trafficcontrol/issues/5038)
- Fixed #3661 - Anonymous Proxy ipv4 whitelist does not work
- Fixed #1847 - Delivery Service with SSL keys are no longer allowed to be updated when the fields changed are relevant to the SSL Keys validity.
- Fixed #5153 - Right click context menu on new ag-grid tables appearing at the wrong place after scrolling. [Related github issue](https://github.com/apache/trafficcontrol/issues/5153)
- Fixed the `GET /api/x/jobs` and `GET /api/x/jobs/:id` Traffic Ops API routes to allow falling back to Perl via the routing blacklist
- Fixed ORT config generation not using the coalesce_number_v6 Parameter.
- Fixed POST deliveryservices/request (designed to simple send an email) regression which erroneously required deep caching type and routing name. [Related github issue](https://github.com/apache/trafficcontrol/issues/4735)
- Removed audit logging from the `POST /api/x/serverchecks` Traffic Ops API endpoint in order to reduce audit log spam
- Fixed an issue that causes Traffic Router to mistakenly route to caches that had recently been set from ADMIN_DOWN to OFFLINE
- Fixed an issue that caused Traffic Monitor to poll caches that did not have the status ONLINE/REPORTED/ADMIN_DOWN
- Fixed /deliveryservice_stats regression restricting metric type to a predefined set of values. [Related github issue](https://github.com/apache/trafficcontrol/issues/4740)
- Fixed audit logging from the `/jobs` APIs to bring them back to the same level of information provided by TO-Perl
- Fixed `maxRevalDurationDays` validation for `POST /api/1.x/user/current/jobs` and added that validation to the `/api/x/jobs` endpoints
- Fixed slice plugin error in delivery service request view. [Related github issue](https://github.com/apache/trafficcontrol/issues/4770)
- Fixed update procedure of servers, so that if a server is linked to one or more delivery services, you cannot change its "cdn". [Related github issue](https://github.com/apache/trafficcontrol/issues/4116)
- Fixed `POST /api/x/steering` and `PUT /api/x/steering` so that a steering target with an invalid `type` is no longer accepted. [Related github issue](https://github.com/apache/trafficcontrol/issues/3531)
- Fixed `cachegroups` READ endpoint, so that if a request is made with the `type` specified as a non integer value, you get back a `400` with error details, instead of a `500`. [Related github issue](https://github.com/apache/trafficcontrol/issues/4703)
- Fixed ORT bug miscalculating Mid Max Origin Connections as all servers, usually resulting in 1.
- Fixed ORT atstccfg helper log to append and not overwrite old logs. Also changed to log to /var/log/ort and added a logrotate to the RPM. See the ORT README.md for details.
- Added Delivery Service Raw Remap `__RANGE_DIRECTIVE__` directive to allow inserting the Range Directive after the Raw Remap text. This allows Raw Remaps which manipulate the Range.
- Added an option for `coordinateRange` in the RGB configuration file, so that in case a client doesn't have a postal code, we can still determine if it should be allowed or not, based on whether or not the latitude/ longitude of the client falls within the supplied ranges. [Related github issue](https://github.com/apache/trafficcontrol/issues/4372)
- Fixed TR build configuration (pom.xml) to invoke preinstall.sh. [Related github issue](https://github.com/apache/trafficcontrol/issues/4882)
- Fixed #3548 - Prevents DS regexes with non-consecutive order from generating invalid CRconfig/snapshot.
- Fixes #4984 - Lets `create_tables.sql` be run concurrently without issue
- Fixed #5020, #5021 - Creating an ASN with the same number and same cache group should not be allowed.
- Fixed #5006 - Traffic Ops now generates the Monitoring on-the-fly if the snapshot doesn't exist, and logs an error. This fixes upgrading to 4.x to not break the CDN until a Snapshot is done.
- Fixed #4680 - Change Content-Type to application/json for TR auth calls
- Fixed #4292 - Traffic Ops not looking for influxdb.conf in the right place
- Fixed #5102 - Python client scripts fail silently on authentication failures
- Fixed #5103 - Python client scripts crash on connection errors
- Fixed matching of wildcards in subjectAlternateNames when loading TLS certificates
- Fixed #5180 - Global Max Mbps and Tps is not send to TM
- Fixed #3528 - Fix Traffic Ops monitoring.json missing DeliveryServices
- Fixed an issue where the jobs and servers table in Traffic Portal would not clear a column's filter when it's hidden
- Fixed an issue with Traffic Router failing to authenticate if secrets are changed
- Fixed validation error message for Traffic Ops `POST /api/x/profileparameters` route
- Fixed #5216 - Removed duplicate button to link delivery service to server [Related Github issue](https://github.com/apache/trafficcontrol/issues/5216)
- Fixed an issue where Traffic Router would erroneously return 503s or NXDOMAINs if the caches in a cachegroup were all unavailable for a client's requested IP version, rather than selecting caches from the next closest available cachegroup.
- Fixed an issue where downgrading the database would fail while having server interfaces with null gateways, MTU, and/or netmasks.
- Fixed an issue where partial upgrades of the database would occasionally fail to apply 2020081108261100_add_server_ip_profile_trigger.
- Fixed #5197 - Allows users to assign topology-based DS to ORG servers [Related Github issue](https://github.com/apache/trafficcontrol/issues/5197)
- Fixed #5161 - Fixes topology name character validation [Related Github issue](https://github.com/apache/trafficcontrol/issues/5161)
- Fixed #5237 - /isos API endpoint rejecting valid IPv6 addresses with CIDR-notation network prefixes.
- Fixed an issue with Traffic Monitor to fix peer polling to work as expected
- Fixed #5274 - CDN in a Box's Traffic Vault image failed to build due to Basho's repo responding with 402 Payment Required. The repo has been removed from the image.
- Fixed #5069 - For LetsEncryptDnsChallengerWatcher in Traffic Router, the cr-config location is configurable instead of only looking at `/opt/traffic_router/db/cr-config.json`
- Fixed #5191 - Error from IMS requests to /federations/all
- Fixed Astats csv issue where it could crash if caches dont return proc data
- Fixed #5380 - Show the correct servers (including ORGs) when a topology based DS with required capabilities + ORG servers is queried for the assigned servers
- Fixed parent.config generation for topology-based delivery services (inline comments not supported)
- Fixed parent.config generation for MSO delivery services with required capabilities

### Changed
- Changed some Traffic Ops Go Client methods to use `DeliveryServiceNullable` inputs and outputs.
- When creating invalidation jobs through TO/TP, if an identical regex is detected that overlaps its time, then warnings
will be returned indicating that overlap exists.
- Changed Traffic Portal to disable browser caching on GETs until it utilizes the If-Modified-Since functionality that the TO API now provides.
- Changed Traffic Portal to use Traffic Ops API v3
- Changed Traffic Portal to use the more performant and powerful ag-grid for all server and invalidation request tables.
- Changed ORT Config Generation to be deterministic, which will prevent spurious diffs when nothing actually changed.
- Changed ORT to find the local ATS config directory and use it when location Parameters don't exist for many required configs, including all Delivery Service files (Header Rewrites, Regex Remap, URL Sig, URI Signing).
- Changed ORT to not update ip_allow.config but log an error if it needs updating in syncds mode, and only actually update in badass mode.
    - ATS has a known bug, where reloading when ip_allow.config has changed blocks arbitrary addresses. This will break things by not allowing any new necessary servers, but prevents breaking the Mid server. There is no solution that doesn't break something, until ATS fixes the bug, and breaking an Edge is better than breaking a Mid.
- Changed the access logs in Traffic Ops to now show the route ID with every API endpoint call. The Route ID is appended to the end of the access log line.
- Changed Traffic Monitor's `tmconfig.backup` to store the result of `GET /api/2.0/cdns/{{name}}/configs/monitoring` instead of a transformed map
- Changed OAuth workflow to use Basic Auth if client secret is provided per RFC6749 section 2.3.1.
- [Multiple Interface Servers](https://github.com/apache/trafficcontrol/blob/master/blueprints/multi-interface-servers.md)
    - Interface data is constructed from IP Address/Gateway/Netmask (and their IPv6 counterparts) and Interface Name and Interface MTU fields on services. These **MUST** have proper, valid data before attempting to upgrade or the upgrade **WILL** fail. In particular IP fields need to be valid IP addresses/netmasks, and MTU must only be positive integers of at least 1280.
    - The `/servers` and `/servers/{{ID}}}` TO API endpoints have been updated to use and reflect multi-interface servers.
    - Updated `/cdns/{{name}}/configs/monitoring` TO API endpoint to return multi-interface data.
    - CDN Snapshots now use a server's "service addresses" to provide its IP addresses.
    - Changed the `Cache States` tab of the Traffic Monitor UI to properly handle multiple interfaces.
    - Changed the `/publish/CacheStats` in Traffic Monitor to support multiple interfaces.
    - Changed the CDN-in-a-Box server enrollment template to support multiple interfaces.
- Changed Tomcat Java dependency to 8.5.57.
- Changed Spring Framework Java dependency to 4.2.5.
- Changed certificate loading code in Traffic Router to use Bouncy Castle instead of deprecated Sun libraries.
- Changed deprecated AsyncHttpClient Java dependency to use new active mirror and updated to version 2.12.1.
- Changed Traffic Portal to use the more performant and powerful ag-grid for the delivery service request (DSR) table.
- Traffic Ops: removed change log entry created during server update/revalidation unqueue
- Updated CDN in a Box to CentOS 8 and added `RHEL_VERSION` Docker build arg so CDN in a Box can be built for CentOS 7, if desired
- Added Delivery Service Raw Remap `__CACHEKEY_DIRECTIVE__` directive to allow inserting the cachekey directive into the Raw Remap text. This allows Raw Remaps which manipulate the cachekey.

### Deprecated
- Deprecated the non-nullable `DeliveryService` Go struct and other structs that use it. `DeliveryServiceNullable` structs should be used instead.
- Deprecated the `insecure` option in `traffic_ops_golang` in favor of `"tls_config": { "InsecureSkipVerify": <bool> }`
- Importing Traffic Ops Go clients via the un-versioned `github.com/apache/trafficcontrol/traffic_ops/client` is now deprecated in favor of versioned import paths e.g. `github.com/apache/trafficcontrol/traffic_ops/v3-client`.

### Removed
- Removed deprecated Traffic Ops Go Client methods.
- Configuration generation logic in the TO API (v1) for all files and the "meta" route - this means that versions of Traffic Ops ORT earlier than 4.0.0 **will not work any longer** with versions of Traffic Ops moving forward.
- Removed from Traffic Portal the ability to view cache server config files as the contents are no longer reliable through the TO API due to the introduction of atstccfg.


## [4.1.0] - 2020-04-23
### Added
- Added support for use of ATS Slice plugin as an additonal option to range request handling on HTTP/DNS DSes.
- Added a boolean to delivery service in Traffic Portal and Traffic Ops to enable EDNS0 client subnet at the delivery service level and include it in the cr-config.
- Updated Traffic Router to read new EDSN0 client subnet field and route accordingly only for enabled delivery services. When enabled and a subnet is present in the request, the subnet appears in the `chi` field and the resolver address is in the `rhi` field.
- Traffic Router DNSSEC zone diffing: if enabled via the new "dnssec.zone.diffing.enabled" TR profile parameter, TR will diff existing zones against newly generated zones in order to determine if a zone needs to be re-signed. Zones are typically generated on every snapshot and whenever new DNSSEC keys are found, and since signing a zone is a relatively CPU-intensive operation, this optimization can drastically reduce the CPU time taken to process new snapshots and new DNSSEC keys.
- Added an optimistic quorum feature to Traffic Monitor to prevent false negative states from propagating to downstream components in the event of network isolation.
- Added the ability to fetch users by role
- Added an API 1.5 endpoint to generate delivery service certificates using Let's Encrypt
- Added an API 1.5 endpoint to GET a single or all records for Let's Encrypt DNS challenge
- Added an API 1.5 endpoint to renew certificates
- Added ability to create multiple objects from generic API Create with a single POST.
- Added debugging functionality to CDN-in-a-Box.
- Added an SMTP server to CDN-in-a-Box.
- Cached builder Docker images on Docker Hub to speed up build time
- Added functionality in the GET endpoints to support the "If-Modified-Since" header in the incoming requests.
- Traffic Ops Golang Endpoints
  - /api/2.0 for all of the most recent route versions
  - /api/1.1/cachegroupparameters/{{cachegroupID}}/{{parameterID}} `(DELETE)`
  - /api/1.5/stats_summary `(POST)`
  - /api/1.1/cdns/routing
  - /api/1.1/cachegroupparameters/ `(GET, POST)`
  - /api/2.0/isos
  - /api/1.5/deliveryservice/:id/routing
  - /api/1.5/deliveryservices/sslkeys/generate/letsencrypt `POST`
  - /api/2.0/deliveryservices/xmlId/:XMLID/sslkeys `DELETE`
  - /deliveryserviceserver/:dsid/:serverid
  - /api/1.5/letsencrypt/autorenew `POST`
  - /api/1.5/letsencrypt/dnsrecords `GET`
  - /api/2.0/vault/ping `GET`
  - /api/2.0/vault/bucket/:bucket/key/:key/values `GET`
  - /api/2.0/servercheck `GET`
  - /api/2.0/servercheck/extensions/:id `(DELETE)`
  - /api/2.0/servercheck/extensions `(GET, POST)`
  - /api/2.0/servers/:name-or-id/update `POST`
  - /api/2.0/plugins `(GET)`
  - /api/2.0/snapshot `PUT`

### Changed
- Add null check in astats plugin before calling strtok to find ip mask values in the config file
- Fix to traffic_ops_ort.pl to strip specific comment lines before checking if a file has changed.  Also promoted a changed file message from DEBUG to ERROR for report mode.
- Fixed Traffic Portal regenerating CDN DNSSEC keys with the wrong effective date
- Fixed issue #4583: POST /users/register internal server error caused by failing DB query
- Type mutation through the api is now restricted to only those types that apply to the "server" table
- Updated The Traffic Ops Python, Go and Java clients to use API version 2.0 (when possible)
- Updated CDN-in-a-Box scripts and enroller to use TO API version 2.0
- Updated numerous, miscellaneous tools to use TO API version 2.0
- Updated TP to use TO API v2
- Updated TP application build dependencies
- Modified Traffic Monitor to poll over IPv6 as well as IPv4 and separate the availability statuses.
- Modified Traffic Router to separate availability statuses between IPv4 and IPv6.
- Modified Traffic Portal and Traffic Ops to accept IPv6 only servers.
- Updated Traffic Monitor to default to polling both IPv4 and IPv6.
- Traffic Ops, Traffic Monitor, Traffic Stats, and Grove are now compiled using Go version 1.14. This requires a Traffic Vault config update (see note below).
- Existing installations **must** enable TLSv1.1 for Traffic Vault in order for Traffic Ops to reach it. See [Enabling TLS 1.1](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_vault.html#tv-admin-enable-tlsv1-1) in the Traffic Vault administrator's guide for instructions.
- Changed the `totalBytes` property of responses to GET requests to `/deliveryservice_stats` to the more appropriate `totalKiloBytes` in API 2.x
- Fix to traffic_ops_ort to generate logging.yaml files correctly.
- Fixed issue #4650: add the "Vary: Accept-Encoding" header to all responses from Traffic Ops

### Deprecated/Removed
- The Traffic Ops `db/admin.pl` script has now been removed. Please use the `db/admin` binary instead.
- Traffic Ops Python client no longer supports Python 2.
- Traffic Ops API Endpoints
  - /api_capabilities/:id
  - /asns/:id
  - /cachegroups/:id (GET)
  - /cachegroup/:parameterID/parameter
  - /cachegroups/:parameterID/parameter/available
  - /cachegroups/:id/unassigned_parameters
  - /cachegroups/trimmed
  - /cdns/:name/configs/routing
  - /cdns/:name/federations/:id (GET)
  - /cdns/configs
  - /cdns/:id (GET)
  - /cdns/:id/snapshot
  - /cdns/name/:name (GET)
  - /cdns/usage/overview
  - /deliveryservice_matches
  - /deliveryservice_server/:dsid/:serverid
  - /deliveryservice_user
  - /deliveryservice_user/:dsId/:userId
  - /deliveryservices/hostname/:name/sslkeys
  - /deliveryservices/{dsid}/regexes/{regexid} (GET)
  - /deliveryservices/:id (GET)
  - /deliveryservices/:id/state
  - /deliveryservices/xmlId/:XMLID/sslkeys/delete
  - /divisions/:division_name/regions
  - /divisions/:id
  - /divisions/name/:name
  - /hwinfo/dtdata
  - /jobs/:id
  - /keys/ping
  - /logs/:days/days
  - /parameters/:id (GET)
  - /parameters/:id/profiles
  - /parameters/:id/unassigned_profiles
  - /parameters/profile/:name
  - /parameters/validate
  - /phys_locations/trimmed
  - /phys_locations/:id (GET)
  - /profile/:id (GET)
  - /profile/:id/unassigned_parameters
  - /profile/trimmed
  - /regions/:id (GET, DELETE)
  - /regions/:region_name/phys_locations
  - /regions/name/:region_name
  - /riak/bucket/:bucket/key/:key/vault
  - /riak/ping
  - /riak/stats
  - /servercheck/aadata
  - /servers/hostname/:hostName/details
  - /servers/status
  - /servers/:id (GET)
  - /servers/totals
  - /snapshot/:cdn
  - /stats_summary/create
  - /steering/:deliveryservice/targets/:target (GET)
  - /tenants/:id (GET)
  - /statuses/:id (GET)
  - /to_extensions/:id/delete
  - /to_extensions
  - /traffic_monitor/stats
  - /types/trimmed
  - /types/{{ID}} (GET)
  - /user/current/jobs
  - /users/:id/deliveryservices
  - /servers/checks
  - /user/{{user ID}}/deliveryservices/available

## [4.0.0] - 2019-12-16
### Added
- Traffic Router: TR now generates a self-signed certificate at startup and uses it as the default TLS cert.
  The default certificate is used whenever a client attempts an SSL handshake for an SNI host which does not match
  any of the other certificates.
- Client Steering Forced Diversity: force Traffic Router to return more unique edge caches in CLIENT_STEERING results instead of the default behavior which can sometimes return a result of multiple targets using the same edge cache. In the case of edge cache failures, this feature will give clients a chance to retry a different edge cache. This can be enabled with the new "client.steering.forced.diversity" Traffic Router profile parameter.
- Traffic Ops Golang Endpoints
  - /api/1.4/deliveryservices `(GET,POST,PUT)`
  - /api/1.4/users `(GET,POST,PUT)`
  - /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys `GET`
  - /api/1.1/deliveryservices/hostname/:hostname/sslkeys `GET`
  - /api/1.1/deliveryservices/sslkeys/add `POST`
  - /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete `GET`
  - /api/1.4/deliveryservices_required_capabilities `(GET,POST,DELETE)`
  - /api/1.1/servers/status `GET`
  - /api/1.4/cdns/dnsseckeys/refresh `GET`
  - /api/1.1/cdns/name/:name/dnsseckeys `GET`
  - /api/1.1/roles `GET`
  - /api/1.4/cdns/name/:name/dnsseckeys `GET`
  - /api/1.4/user/login/oauth `POST`
  - /api/1.1/servers/:name/configfiles/ats `GET`
  - /api/1.1/servers/:id/queue_update `POST`
  - /api/1.1/profiles/:name/configfiles/ats/* `GET`
  - /api/1.4/profiles/name/:name/copy/:copy
  - /api/1.1/servers/:name/configfiles/ats/* `GET`
  - /api/1.1/cdns/:name/configfiles/ats/* `GET`
  - /api/1.1/servers/:id/status `PUT`
  - /api/1.1/dbdump `GET`
  - /api/1.1/servers/:name/configfiles/ats/parent.config
  - /api/1.1/servers/:name/configfiles/ats/remap.config
  - /api/1.1/user/login/token `POST`
  - /api/1.4/deliveryservice_stats `GET`
  - /api/1.1/deliveryservices/request
  - /api/1.1/federations/:id/users
  - /api/1.1/federations/:id/users/:userID
  - /api/1.2/current_stats
  - /api/1.1/osversions
  - /api/1.1/stats_summary `GET`
  - /api/1.1/api_capabilities `GET`
  - /api/1.1/user/current `PUT`
  - /api/1.1/federations/:id/federation_resolvers `(GET, POST)`

- Traffic Router: Added a tunable bounded queue to support DNS request processing.
- Traffic Ops API Routing Blacklist: via the `routing_blacklist` field in `cdn.conf`, enable certain whitelisted Go routes to be handled by Perl instead (via the `perl_routes` list) in case a regression is found in the Go handler, and explicitly disable any routes via the `disabled_routes` list. Requests to disabled routes are immediately given a 503 response. Both fields are lists of Route IDs, and route information (ID, version, method, path, and whether or not it can bypass to Perl) can be found by running `./traffic_ops_golang --api-routes`. To disable a route or have it bypassed to Perl, find its Route ID using the previous command and put it in the `disabled_routes` or `perl_routes` list, respectively.
- To support reusing a single riak cluster connection, an optional parameter is added to riak.conf: "HealthCheckInterval". This options takes a 'Duration' value (ie: 10s, 5m) which affects how often the riak cluster is health checked.  Default is currently set to: "HealthCheckInterval": "5s".
- Added a new Go db/admin binary to replace the Perl db/admin.pl script which is now deprecated and will be removed in a future release. The new db/admin binary is essentially a drop-in replacement for db/admin.pl since it supports all of the same commands and options; therefore, it should be used in place of db/admin.pl for all the same tasks.
- Added an API 1.4 endpoint, /api/1.4/cdns/dnsseckeys/refresh, to perform necessary behavior previously served outside the API under `/internal`.
- Added the DS Record text to the cdn dnsseckeys endpoint in 1.4.
- Added monitoring.json snapshotting. This stores the monitoring json in the same table as the crconfig snapshot. Snapshotting is now required in order to push out monitoring changes.
- To traffic_ops_ort.pl added the ability to handle ##OVERRIDE## delivery service ANY_MAP raw remap text to replace and comment out a base delivery service remap rules. THIS IS A TEMPORARY HACK until versioned delivery services are implemented.
- Snapshotting the CRConfig now deletes HTTPS certificates in Riak for delivery services which have been deleted in Traffic Ops.
- Added a context menu in place of the "Actions" column from the following tables in Traffic Portal: cache group tables, CDN tables, delivery service tables, parameter tables, profile tables, server tables.
- Traffic Portal standalone Dockerfile
- In Traffic Portal, removes the need to specify line breaks using `__RETURN__` in delivery service edge/mid header rewrite rules, regex remap expressions, raw remap text and traffic router additional request/response headers.
- In Traffic Portal, provides the ability to clone delivery service assignments from one cache to another cache of the same type. Issue #2963.
- Added an API 1.4 endpoint, /api/1.4/server_capabilities, to create, read, and delete server capabilities.
- Traffic Ops now allows each delivery service to have a set of query parameter keys to be retained for consistent hash generation by Traffic Router.
- In Traffic Portal, delivery service table columns can now be rearranged and their visibility toggled on/off as desired by the user. Hidden table columns are excluded from the table search. These settings are persisted in the browser.
- Added an API 1.4 endpoint, /api/1.4/user/login/oauth to handle SSO login using OAuth.
- Added /#!/sso page to Traffic Portal to catch redirects back from OAuth provider and POST token into the API.
- In Traffic Portal, server table columns can now be rearranged and their visibility toggled on/off as desired by the user. Hidden table columns are excluded from the table search. These settings are persisted in the browser.
- Added pagination support to some Traffic Ops endpoints via three new query parameters, limit and offset/page
- Traffic Ops now supports a "sortOrder" query parameter on some endpoints to return API responses in descending order
- Traffic Ops now uses a consistent format for audit logs across all Go endpoints
- Added cache-side config generator, atstccfg, installed with ORT. Includes all configs. Includes a plugin system.
- Fixed ATS config generation to omit regex remap, header rewrite, URL Sig, and URI Signing files for delivery services not assigned to that server.
- In Traffic Portal, all tables now include a 'CSV' link to enable the export of table data in CSV format.
- Pylint configuration now enforced (present in [a file in the Python client directory](./traffic_control/clients/python/pylint.rc))
- Added an optional SMTP server configuration to the TO configuration file, api now has unused abilitiy to send emails
- Traffic Monitor now has "gbps" calculated stat, allowing operators to monitor bandwidth in Gbps.
- Added an API 1.4 endpoint, /api/1.4/deliveryservices_required_capabilities, to create, read, and delete associations between a delivery service and a required capability.
- Added ATS config generation omitting parents without Delivery Service Required Capabilities.
- In Traffic Portal, added the ability to create, view and delete server capabilities and associate those server capabilities with servers and delivery services. See [blueprint](./blueprints/server-capabilitites.md)
- Added validation to prevent assigning servers to delivery services without required capabilities.
- Added deep coverage zone routing percentage to the Traffic Portal dashboard.
- Added a `traffic_ops/app/bin/osversions-convert.pl` script to convert the `osversions.cfg` file from Perl to JSON as part of the `/osversions` endpoint rewrite.
- Added [Experimental] - Emulated Vault suppling a HTTP server mimicking RIAK behavior for usage as traffic-control vault.
- Added Traffic Ops Client function that returns a Delivery Service Nullable Response when requesting for a Delivery Service by XMLID

### Changed
- Traffic Router:  TR will now allow steering DSs and steering target DSs to have RGB enabled. (fixes #3910)
- Traffic Portal:  Traffic Portal now allows Regional Geo Blocking to be enabled for a Steering Delivery Service.
- Traffic Ops: fixed a regression where the `Expires` cookie header was not being set properly in responses. Also, added the `Max-Age` cookie header in responses.
- Traffic Router, added TLS certificate validation on certificates imported from Traffic Ops
  - validates modulus of private and public keys
  - validates current timestamp falls within the certificate date bracket
  - validates certificate subjects against the DS URL
- Traffic Ops Golang Endpoints
  - Updated /api/1.1/cachegroups: Cache Group Fallbacks are included
  - Updated /api/1.1/cachegroups: fixed so fallbackToClosest can be set through API
    - Warning:  a PUT of an old Cache Group JSON without the fallbackToClosest field will result in a `null` value for that field
- Traffic Router: fixed a bug which would cause `REFUSED` DNS answers if the zone priming execution did not complete within the configured `zonemanager.init.timeout` period.
- Issue 2821: Fixed "Traffic Router may choose wrong certificate when SNI names overlap"
- traffic_ops/app/bin/checks/ToDnssecRefresh.pl now requires "user" and "pass" parameters of an operations-level user! Update your scripts accordingly! This was necessary to move to an API endpoint with proper authentication, which may be safely exposed.
- Traffic Monitor UI updated to support HTTP or HTTPS traffic.
- Traffic Monitor health/stat time now includes full body download (like prior TM <=2.1 version)
- Modified Traffic Router logging format to include an additional field for DNS log entries, namely `rhi`. This defaults to '-' and is only used when EDNS0 client subnet extensions are enabled and a client subnet is present in the request. When enabled and a subnet is present, the subnet appears in the `chi` field and the resolver address is in the `rhi` field.
- Changed traffic_ops_ort.pl so that hdr_rw-&lt;ds&gt;.config files are compared with strict ordering and line duplication when detecting configuration changes.
- Traffic Ops (golang), Traffic Monitor, Traffic Stats are now compiled using Go version 1.11. Grove was already being compiled with this version which improves performance for TLS when RSA certificates are used.
- Fixed issue #3497: TO API clients that don't specify the latest minor version will overwrite/default any fields introduced in later versions
- Fixed permissions on DELETE /api/$version/deliveryservice_server/{dsid}/{serverid} endpoint
- Issue 3476: Traffic Router returns partial result for CLIENT_STEERING Delivery Services when Regional Geoblocking or Anonymous Blocking is enabled.
- Upgraded Traffic Portal to AngularJS 1.7.8
- Issue 3275: Improved the snapshot diff performance and experience.
- Issue 3550: Fixed TC golang client setting for cache control max age
- Issue #3605: Fixed Traffic Monitor custom ports in health polling URL.
- Issue 3587: Fixed Traffic Ops Golang reverse proxy and Riak logs to be consistent with the format of other error logs.
- Database migrations have been collapsed. Rollbacks to migrations that previously existed are no longer possible.
- Issue #3750: Fixed Grove access log fractional seconds.
- Issue #3646: Fixed Traffic Monitor Thresholds.
- Modified Traffic Router API to be available via HTTPS.
- Added fields to traffic_portal_properties.json to configure SSO through OAuth.
- Added field to cdn.conf to configure whitelisted URLs for Json Key Set URL returned from OAuth provider.
- Improved [profile comparison view in Traffic Portal](https://github.com/apache/trafficcontrol/blob/master/blueprints/profile-param-compare-manage.md).
- Issue #3871 - provides users with a specified role the ability to mark any delivery service request as complete.
- Fixed Traffic Ops Golang POST servers/id/deliveryservice continuing erroneously after a database error.
- Fixed Traffic Ops Golang POST servers/id/deliveryservice double-logging errors.
- Issue #4131 - The "Clone Delivery Service Assignments" menu item is hidden on a cache when the cache has zero delivery service assignments to clone.
- Traffic Portal - Turn off TLSv1
- Removed Traffic Portal dependency on Restangular
- Issue #1486 - Dashboard graph for bandwidth now displays units in the tooltip when hovering over a data point

### Deprecated/Removed
- Traffic Ops API Endpoints
  - /api/1.1/cachegroup_fallbacks
  - /api_capabilities `POST`

## [3.0.0] - 2018-10-30
### Added
- Removed MySQL-to-Postgres migration tools.  This tool is supported for 1.x to 2.x upgrades only and should not be used with 3.x.
- Backup Edge Cache group: If the matched group in the CZF is not available, this list of backup edge cache group configured via Traffic Ops API can be used as backup. In the event of all backup edge cache groups not available, GEO location can be optionally used as further backup. APIs detailed [here](http://traffic-control-cdn.readthedocs.io/en/latest/development/traffic_ops_api/v12/cachegroup_fallbacks.html)
- Traffic Ops Golang Proxy Endpoints
  - /api/1.4/users `(GET,POST,PUT)`
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
- Cachegroup Localization Methods: The Cachegroup API now supports an optional `localizationMethods` field which specifies the localization methods allowed for that cachegroup (currently 'DEEP_CZ', 'CZ', and 'GEO'). By default if this field is null/empty, all localization methods are enabled. After Traffic Router has localized a client, it will only route that client to cachegroups that have enabled the localization method used. For example, this can be used to prevent GEO-localized traffic (i.e. most likely from off-net/internet clients) to cachegroups that aren't optimal for internet traffic.
- Traffic Monitor Client Update: Traffic Monitor is updated to use the Traffic Ops v13 client.
- Removed previously deprecated `traffic_monitor_java`
- Added `infrastructure/cdn-in-a-box` for Apachecon 2018 demonstration
- The CacheURL Delivery service field is deprecated.  If you still need this functionality, you can create the configuration explicitly via the raw remap field.

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

[unreleased]: https://github.com/apache/trafficcontrol/compare/RELEASE-8.0.0...HEAD
[8.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-8.0.0...RELEASE-7.0.0
[7.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-7.0.0...RELEASE-6.0.0
[6.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-6.0.0...RELEASE-5.0.0
[5.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-4.1.0...RELEASE-5.0.0
[4.1.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-4.0.0...RELEASE-4.1.0
[4.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-3.0.0...RELEASE-4.0.0
[3.0.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-2.2.0...RELEASE-3.0.0
[2.2.0]: https://github.com/apache/trafficcontrol/compare/RELEASE-2.1.0...RELEASE-2.2.0
