# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## [Unreleased]
### Added
- Per-DeliveryService Routing Names: you can now choose a Delivery Service's Routing Name (rather than a hardcoded "tr" or "edge" name). This might require a few pre-upgrade steps detailed [here](http://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops/migration_from_20_to_22.html#per-deliveryservice-routing-names)
- Golang Proxy Endpoints (R=REST endpoints for GET, POST, PUT, DELETE)
  - /api/1.3/asns (R)
  - /api/1.3/cdns (R)
  - /api/1.3/cdns/capacity
  - /api/1.3/cdns/configs
  - /api/1.3/cdns/dnsseckeys
  - /api/1.3/cdns/domain
  - /api/1.3/cdns/monitoring
  - /api/1.3/cdns/health
  - /api/1.3/cdns/routing
  - /api/1.3/deliveryservice_requests (R)
  - /api/1.3/divisions (R)
  - /api/1.3/hwinfos
  - /api/1.3/parameters (R)
  - /api/1.3/phys_locations (R)
  - /api/1.3/ping
  - /api/1.3/profiles (R)
  - /api/1.3/regions (R)
  - /api/1.3/servers (R)
  - /api/1.3/servers/checks
  - /api/1.3/servers/details
  - /api/1.3/servers/status
  - /api/1.3/servers/totals
  - /api/1.3/statuses (R)
  - /api/1.3/system/info
  - /api/1.3/types (R)

### Changed
- Reformatted this CHANGELOG file to the keep-a-changelog format

[Unreleased]: https://github.com/apache/incubator-trafficcontrol/compare/RELEASE-2.1.0...HEAD
