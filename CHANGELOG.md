# Changelog
All notable changes to this project will be documented in this file.
=======

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## [Unreleased]
### Added
- Per-DeliveryService Routing Names: you can now choose a Delivery Service's Routing Name (rather than a hardcoded "tr" or "edge" name). This might require a few pre-upgrade steps detailed [here](http://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops/migration_from_20_to_22.html#per-deliveryservice-routing-names)

- Backup Edge Cache group: Backup Edge group for a particular cache group can be configured using coverage zone files. With this, users will be able control traffic within portions of their network, there by avoiding choosing fall back cache groups using geo co-ordinates(getClosestAvailableCachegroup) This can be controlled using "backupZones" which contains configuration for backup cache groups parsed from Coverage zone file as explained [here] (http://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops/using.html#the-coverage-zone-file-and-asn-table)


### Changed
- Reformatted this CHANGELOG file to the keep-a-changelog format

-[Unreleased]: https://github.com/apache/incubator-trafficcontrol/compare/RELEASE-2.1.0...HEAD 

