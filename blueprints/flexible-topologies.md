<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

# Flexible Topologies

## Problem Description

Today, a Traffic Control CDN is limited to 2 tiers -- *EDGE* and *MID* -- with the option to skip the *MID* tier for certain Delivery Service types (e.g. `HTTP_LIVE` and `HTTP_NO_CACHE`). In addition, a CDN is limited to one global parent hierarchy, which is defined via the `parent_cachegroup` and `secondary_parent_cachegroup` fields of cachegroups. Both of these problems limit a CDN's ability to scale with increased demand and changing usage patterns, and providing the ability to add more tiers to a CDN helps it keep up with that growth. A Topology that works well for one set of Delivery Services might not be ideal for a different set of Delivery Services, and a CDN needs the flexibility to provide the best Topology for any given Delivery Service -- with any number of tiers and custom caching hierarchies.

## Proposed Change

Traffic Control will provide the ability to define one or more Topologies, and a Topology can have any number of Delivery Services assigned to it. A Topology will be composed of Cachegroups along with their primary/secondary parent relationships to other Cachegroups as defined by the Topology.

If a Delivery Service is assigned to a Topology, any `deliveryservice_server` assignments it has to `EDGE` caches will be ignored, because it will be assigned to all caches in the Delivery Service's CDN (filtered by server capabilities) that belong to the Topology's cachegroups. Ideally, this feature will obsolete legacy `deliveryservice_server` assignments, since Topologies negate the need to assign Delivery Services to individual `EDGE` caches. Nonetheless, legacy `deliveryservice_server` assignments will be supported alongside Topology-based Delivery Services for some time until all Delivery Services have been migrated to Topologies.

### Traffic Portal Impact

Traffic Portal will need new pages for:
- creating and viewing Topologies (and since TP sidebar menu currently has "Topology" as an item already, that may need to be renamed if "Topologies" is going to be a sub-menu item of that)

Existing Traffic Portal pages will need to be updated:
- all the delivery services and delivery service requests views, to include the new `topology` field
- delivery service server assignment views should prohibit assigning `EDGE` servers to a delivery service that has a Topology assigned already (`ORIGIN` servers may still need to be assignable for MSO purposes)
- the CDN snapshot view, in order to account for a new top-level section (`topologies`) in the `CRConfig`

Since Delivery Services will no longer be constrained to one global Topology as they are today, it would be extremely useful to be able to visualize a Delivery Service's Topology like a tree, where each node in the tree is a cachegroup, and the edges between nodes are the primary/secondary parent relationships between them. Clicking on a particular node would show all the servers in that cachegroup that could serve a request for the Delivery Service. This visualization will most likely be different from the Topology form for creating a Topology and does not necessarily need to be provided by Traffic Portal.

### Traffic Ops Impact

Traffic Ops will provide the ability to create Topologies, composed of cachegroups and parent relationships, which will be assignable to one or more Delivery Services.

#### REST API Impact

The following is the JSON representation of a `Topology` object:

```JSON
{
    "name": "foo",
    "description": "a foo topology",
    "nodes": [
        {
            "cachegroup": "child-cachegroup",
            "parents": [1, 2]
        },
        {
            "cachegroup": "parent-cachegroup",
            "parents": []
        },
        {
            "cachegroup": "secondary-parent-cachegroup",
            "parents": []
        }
    ]
}
```

The following table describes the top-level `Topology` object:

| field       | type                        | optionality | description                                                         |
| ----------- | --------------------------- | ----------- | ------------------------------------------------------------------- |
| name        | string                      | required    | a unique name for identifying this Topology                         |
| description | string                      | required    | the description of this Topology                                    |
| nodes       | array of `node` sub-objects | required    | the set of `nodes` in this topology, similar to an *adjacency list* |

The following table describes the `node` sub-object:

| field      | type              | optionality | description                                                                                                                                                                                                 |
| ---------- | ----------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| cachegroup | string            | required    | the `short_name` of a cachegroup this node maps to in the Topology                                                                                                                                          |
| parents    | array of integers | required    | zero-based indexes to other nodes in the Topology's `nodes` array, where the 1st element is for the *primary* parent relationship and the 2nd element is for the *secondary* parent relationship, and so on |

API constraints:
- a Topology's `name` must consist of alphanumeric or hyphen characters
- a Topology must have at least 1 `node`; otherwise, it is useless
- there cannot be multiple `nodes` for the same cachegroup in a Topology
- `parents` must have 0, 1 or 2 elements, cannot contain duplicates, cannot contain the index of its own `node`, and cannot contain the index of `nodes` whose cachegroup is of type `EDGE_LOC`
- leaf `nodes` must be cachegroups of type `EDGE_LOC`
- all `nodes` in the Topology must be reachable -- i.e. a `node` is either a leaf (which would be an `EDGE_LOC`) or is a parent of at least one other node
- a Topology cannot contain a cycle (through any combination of primary/secondary parent relationships)
- a Topology cannot be deleted if one or more Delivery Services are still assigned to it
- a Cachegroup cannot be deleted if it is currently being used in a Topology
- a Topology cannot have `STEERING` or `CLIENT_STEERING` delivery services assigned to it (because those types are not assigned to caches -- their _targets_ are)

The following new endpoints will be required:

##### `GET /topologies`

response JSON:
```JSON
{ "response": [
    {
        "name": "foo",
        "description": "a foo topology",
        "nodes": [
            {
                "cachegroup": "child-cachegroup",
                "parents": [1, 2]
            },
            {
                "cachegroup": "parent-cachegroup",
                "parents": []
            },
            {
                "cachegroup": "secondary-parent-cachegroup",
                "parents": []
            }
        ]
    }
]}
```

##### `POST /topologies`

request JSON:
```JSON
{
    "name": "foo",
    "description": "a foo topology",
    "nodes": [
        {
            "cachegroup": "child-cachegroup",
            "parents": [1, 2]
        },
        {
            "cachegroup": "parent-cachegroup",
            "parents": []
        },
        {
            "cachegroup": "secondary-parent-cachegroup",
            "parents": []
        }
    ]
}
```

response JSON:
```JSON
{
    "alerts": [
        {
            "text": "topology was created successfully",
            "level": "success"
        }
    ],
    "response": {
        "name": "foo",
        "description": "a foo topology",
        "nodes": [
            {
                "cachegroup": "child-cachegroup",
                "parents": [1, 2]
            },
            {
                "cachegroup": "parent-cachegroup",
                "parents": []
            },
            {
                "cachegroup": "secondary-parent-cachegroup",
                "parents": []
            }
        ]
    }
}

```

##### `PUT /topologies?name=foo`

request JSON:
```JSON
{
    "name": "foo",
    "description": "a foo topology",
    "nodes": [
        {
            "cachegroup": "child-cachegroup",
            "parents": [1, 2]
        },
        {
            "cachegroup": "parent-cachegroup",
            "parents": []
        },
        {
            "cachegroup": "secondary-parent-cachegroup",
            "parents": []
        }
    ]
}
```

response JSON:
```JSON
{
    "alerts": [
        {
            "text": "topology was updated successfully",
            "level": "success"
        }
    ],
    "response": {
        "name": "foo",
        "description": "a foo topology",
        "nodes": [
            {
                "cachegroup": "child-cachegroup",
                "parents": [1, 2]
            },
            {
                "cachegroup": "parent-cachegroup",
                "parents": []
            },
            {
                "cachegroup": "secondary-parent-cachegroup",
                "parents": []
            }
        ]
    }
}
```

##### `DELETE /topologies?name=foo`

response JSON:
```JSON
{
    "alerts": [
        {
            "text": "topology was deleted successfully",
            "level": "success"
        }
    ]
}
```

##### `/deliveryservices` endpoints

All relevant Delivery Service APIs will have their JSON request and response objects modified to include a new `topology` field which references the name of the topology it's assigned to:
```JSON
{
    ...
    "topology": "foo"
}
```

The `GET /deliveryservices` endpoint should be updated to support `?topology=foo` as a query parameter to retrieve all the Delivery Services that are assigned to a given Topology.

##### `/cachegroups` endpoints

`GET /cachegroups` should be updated to support `?topology=foo` as a query parameter to retrieve all the Cachegroups that are used in a given Topology.

`DELETE /cachegroups` should be updated to return a useful error when trying to delete a Cachegroup that is currently used in a Topology.

`PUT /cachegroups` should update all foreign key references if `short_name` is updated. This should be done automatically in the database via the FK references.

##### The various `/snapshot` endpoints

The various `/snapshot` endpoints will need to be updated to include new Topologies data along with their associations to Delivery Services in the `CRConfig.json` snapshot. The data should only include the `EDGE_LOC` cachegroups of the Topologies, because those are all Traffic Router needs.

##### Various endpoints that are affected by cachegroup parentage or deliveryservice-server assignment

API endpoints that do things such as the following may need to be updated to take Topology-based Delivery Service assignment and parentage into account:
- assign a Delivery Service to a server (or vice versa). It should be prohibited to assign an `EDGE` server to a Delivery Service (or vice versa) that has a Topology assigned already.
- perform an operation on "child" cachegroups -- like queueing updates on "child" caches when changing the status of a "parent" cache. Child caches via both Topologies and legacy Cachegroup parentage would all need to be updated in that case.

#### Client Impact

New Go client methods will be added for the `/topologies` endpoints in order to write TO API tests for the new endpoints. The `/deliveryservices` client methods won't need modified as the new `Topology` field will simply be added to the `DeliveryService` struct. New client methods for the Python client will also be added for each of the new `/topologies` endpoints.

#### Data Model Impact

##### Go structs

New structs will be added for the `/topologies` endpoints, mapping directly to the JSON request bodies in the REST API Impact section:
- `Topology`: for the top-level Topology object, which includes the `nodes` array
- `TopologyNode`: for objects in a Topology's `nodes` array

The `DeliveryService` struct will be updated with a new `Topology` field which is the name of the Topology the Delivery Service is assigned to.

##### Traffic Ops Database

A new `topology` table will be created:

| column      | type | modifiers    |
| ----------- | ---- | ------------ |
| name        | text | not null, PK |
| description | text |              |

A new `topology_cachegroup` table will be created to model the association of cachegroups to topologies:

| column     | type | modifiers                                       |
| ---------- | ---- | ----------------------------------------------- |
| id         | int  | not null, PK                                    |
| topology   | text | not null, FK: references topology(name)         |
| cachegroup | text | not null, FK: references cachegroup(short_name) |

**Constraints**:
- unique (topology, cachegroup) -- a cachegroup can only be in a Topology once.

A new `topology_cachegroup_parents` table will be created to model the parent relationships of cachegroups within a topology:

| column | type | modifiers                                        |
| ------ | ---- | ------------------------------------------------ |
| child  | int  | not null, FK: references topology_cachegroup(id) |
| parent | int  | not null, FK: references topology_cachegroup(id) |
| rank   | int  | not null                                         |

**Constraints**:
- unique (child, rank) -- within a Topology, a cachegroup can only have one primary parent, one secondary parent, and so on.
- unique (child, parent) -- within a Topology, a cachegroup cannot relate to another cachegroup more than once.
- check (rank is either 1 or 2) -- a cachegroup can only have primary and secondary parents currently.

The `deliveryservice` table will be updated to add a new column for the topology it is associated to:

| column   | type | modifiers                     |
| -------- | ---- | ----------------------------- |
| topology | text | FK: references topology(name) |

### ORT Impact

`atstccfg` will need to be updated to request the Topologies from Traffic Ops and use that data to determine the following for config generation:
- what delivery services are assigned to a cache via Topologies -- in addition to legacy `deliveryservice_server` assignments used today -- still taking Server Capabilities into account
- if a delivery service is assigned to a cache via a Topology, the parent and secondary parent cachegroups for that delivery service are determined via the Topology it's assigned to. Otherwise, the parents are determined by the server's cachegroup as they are today.

Since new Topologies can be more than 2 tiers (`EDGE` -> `MID`), `atstccfg` may need to break some assumptions about the current 2-tier hierarchy in order to work with an arbitrary number of "forward proxy tiers" -- e.g. `EDGE` -> `MID` -> `MID` -> `ORIGIN`. Basically, `MID` caches need to be able to forward requests to other `MID` caches -- they can no longer assume that their parents are always origins.

#### ATS config-related fields on Delivery Services

Edge header rewrite rules and raw remap rules should still be applied only to the `EDGE` tier of Topologies. Mid header rewrite rules should probably only be applied at the *first* `MID` tier if there are multiple `MID` tiers in a Topology, because it might not be safe to assume that all mid header rewrite rules could be applied safely through multiple `MID` tiers.

### Traffic Monitor Impact

There should be little (if any) impact to Traffic Monitor for this feature.

### Traffic Router Impact

Traffic Router will need to be made aware of Topologies and their associations to Delivery Services via additions to the CRConfig. No new TR profile parameters should be required to enable Topology-based routing since Topologies are configurable on a per-delivery-service basis.

The CRConfig will need a new top-level field for `topologies`, which will be a map of Topology names to arrays of cachegroup names that make up the "edge tier" of that Topology. The `deliveryServices` section will add an optional `topology` field to each delivery service that is assigned to a Topology. Delivery services that have a Topology assigned will not be referenced explicitly by any `contentServer` objects. Traffic Router will use the Topology information to determine which "edge" cachegroups can be routed to for a particular delivery service. Taking Server Capabilities into account, the CRConfig modifications would ideally be designed in a way that doesn't require Traffic Router to be aware of what Server Capabilities edges have.

Since Topologies will be optional, new Traffic Routers should remain backwards-compatible with old CRConfigs, and old Traffic Routers should remain forwards-compatible with new CRConfigs because Traffic Router ignores unknown fields by default.

### Traffic Stats Impact

There should be little (if any) impact to Traffic Stats for this feature.

### Traffic Vault Impact

This feature should not require any changes to Traffic Vault or its related APIs in Traffic Ops.

### Documentation Impact

The Traffic Ops API reference will be updated to include the new `/topologies` API endpoints as well as all of the relevant `deliveryservices` endpoints. It may be useful to include a new "Topologies Overview/How-To" section in the docs describing how to create and use custom Topologies.

### Testing Impact

For Traffic Ops, new API tests will be written for the new `/topologies` API endpoints, and existing API tests for the `/deliveryservices` endpoints will be updated to test the Topology association.

Traffic Router unit and/or integration tests will be added to test the new functionality of Topology-based delivery services.

Unit tests for `atstccfg` will be added in order to validate ATS config generation for Topology-based delivery services.

Automated end-to-end environment creation (such as CDN-in-a-box) should be updated to include the creation of arbitrary Topologies that would be assigned to one or more delivery services. Those delivery services should be tested by an end user HTTP client to verify the basic functionality of their Topologies. Additionally, the request/data flows should be observed to verify that they match up with the given Topology as expected.

### Performance Impact

There should be no visible impact to performance from an application perspective -- this feature does not introduce anything particularly CPU, network, or storage-intensive. However, this feature will allow the tuning of Topologies in the CDN, which *will* affect end-to-end performance of the CDN in terms of things like latency, cache hit ratio, cache efficiency, etc. CDN architects will be able to make certain trade-offs in their Topologies until their desired end-to-end performance characteristics are met.

Additionally, the more Delivery Services that are migrated from legacy `deliveryservice_server` assignments to Topologies, the smaller the size of the `CRConfig` will get. Currently, `deliveryservice_server` assignments are responsible for the most of the growth in `CRConfig` size due to the natural addition of both Delivery Services and Servers to a CDN over time, and migrating to Topology-based Delivery Services should noticeably reduce the size and growth of the `CRConfig` over time.

### Security Impact

Probably the biggest impact to security will be the custom Topologies themselves, specifically in terms of breaking the assumptions of the existing 2-tier CDN architecture. There will no longer be a single, global parent hierarchy, so things like firewalls/ACLs that allow caches to communicate with each other may need to be updated to account for custom parent hierarchies.

Creating or updating Topologies should be restricted to users with the `operations` role or higher, and Topologies are CDN-wide -- they do not belong to a particular tenant.

### Upgrade Impact

This feature will require a database migration to create new tables and add a new field to the `deliveryservice` table, but existing data does not need to be modified or migrated. Therefore, rolling back the database migration will not cause any data loss until Topologies are actually created and assigned. Additionally, since this feature will not remove any existing tables or columns, the new database schema should be backwards-compatible with the previous version of Traffic Ops. However, this blueprint cannot make any claim as to the backwards-compatibility of schema changes required by *other* features, which would affect the overall backwards-compatibility of the entire release.

This feature will not require components to be upgraded in a specific order, and no special manual steps will be required before, during, or after the upgrade is done.

### Operations Impact

One of the bigger day-to-day operational impacts this feature will have is in the assignment of Delivery Services to Topologies instead of to individual `EDGE` caches. In theory, if the number of unique Topologies is kept as low as possible, it should be easier to choose a specific Topology for a Delivery Service than to assign it to individual caches (unless the default approach is to just assign delivery services to *all* caches). This could be made easier by defining a kind of *default* Topology which is assigned to until it is determined that the default Topology does not meet the Delivery Service's requirements.

Until legacy `deliveryservice_server` assignments are fully removed in favor of Topology-based assignments, there will be extra operational overhead due to having two different ways to assign Delivery Services, but this overhead should be temporary as CDN operators should migrate all Delivery Services to Topologies as soon as possible. That one-time manual migration should also be considered an operational impact and should be scheduled to take place sometime after the upgrade is completed.

Once all Delivery Services have been migrated to Topologies, CDN operators will no longer have to "clone Delivery Service assignments" from a nearby edge cache in the same cachegroup when adding a new edge cache into the CDN. Since Topologies are composed of cachegroups, simply adding a server into a particular cachegroup will give it all the Delivery Service assignments from all the Topologies the cachegroup is used in.

Troubleshooting issues with Delivery Services on the CDN may become slightly more difficult due to the fact that you will no longer be able to assume that any given Delivery Service will follow the same global hierarchy. Unique Topologies will provide unique paths through the CDN from client to origin, so an operator now needs to look up what the path *should be* based on the Delivery Service's Topology in order to triage issues along the delivery path. This may necessitate the development of Topology-based troubleshooting tools and visualizations.

#### Third Party Logging/Monitoring/Analytics

Any external 3rd party logging, analytics, or monitoring tools presuming a 2-tiered CDN architecture may need to be updated to deal with more dynamic and complex N-tiered Topologies, and it may become more difficult to model them accurately. For instance, each "tier" may need to be grouped or modeled differently, and relying on the server type `MID` to mean "the mid tier" will no longer be accurate.

### Developer Impact

Developers primarily need to be aware that they will no longer be able to make the assumption that a CDN is limited to the traditional 2 tiers (`EDGE` and `MID`). There will be any number of tiers, and *cache parentage* is no longer driven by just a server's `cachegroup.parent_cachegroup` and `cachegroup.secondary_parent_cachegroup`. Topologies, which are per Delivery Service, will drive the parentage when used. Additionally, both Topology-based server assignment and legacy `deliveryservice_server` assignment will drive what Delivery Services are "assigned" to a cache for the time being. So, once this feature has been completed, developers will need to keep these changes in mind so that these old assumptions are not mistakenly brought back into the codebase.

## Alternatives

Some alternative designs for the API have been discussed, but they depend on whether or not we'd like to support having the same cachegroup in a Topology more than once (for example, to have an `EDGE_LOC` cachegroup be its own parent but using a different caching medium). That will be prohibited in the implementation of this blueprint, but the current design allows for that functionality to be added in the future without requiring further database migrations or breaking changes to the API.

## Dependencies

This feature requires no known external dependencies at this time.

## References

n/a