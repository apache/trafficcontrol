v2.2.0 [unreleased]
-------------------

### Upgrading

#### Per-DeliveryService Routing Names
A new Delivery Services feature has been added that might require a few pre-upgrade steps: Per-DeliveryService Routing Names. Before this release, DNS Delivery Services were hardcoded to use the name "edge", i.e. "edge.myds.mycdn.com", and HTTP Delivery Services use the name "tr" (or previously "ccr"), i.e. "tr.myds.mycdn.com". As of 2.2, Routing Names will default to "cdn" if left unspecified and can be set to any arbitrary non-dotted hostname.

Pre-2.2 the HTTP Routing Name is configurable via the `http.routing.name` option in in the Traffic Router http.properties config file. If your CDN uses that option to change the name from "tr" to something else, then you will need to perform the following steps for *each* CDN affected:
1. In Traffic Ops, create the following profile parameter (double-check for typos, trailing spaces, etc):

   **name:** upgrade_http_routing_name  
   **config file:** temp  
   **value:** whatever value is used for the affected CDN's http.routing.name

2. Add this parameter to a **single** profile in the affected CDN

With those profile parameters in place Traffic Ops can be safely upgraded to 2.2. Before taking a post-upgrade snapshot, make sure to check your Delivery Service example URLs for unexpected Routing Name changes. Once Traffic Ops has been upgraded to 2.2 and a post-upgrade snapshot has been taken, your Traffic Routers can be upgraded to 2.2 (Traffic Routers must be upgraded *after* Traffic Ops so that they can work with custom per-DeliveryService Routing Names).
