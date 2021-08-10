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
# Layered Profile

## Problem Description

Profiles are unwieldy and dangerous for Operations.

Currently, we have countless Profiles, in broad categories. For example, "EDGE_123_FOO_ABC_ATS_714-RC0-42" might represent servers which are
1. Edges
2. Amiga 123 machines
3. In the Foo CDN
4. In the ABC datacenter
5. Running ATS 7.14 RC0
6. Who knows what 42 means?

Suppose we have Amiga 456 machines, at DEF datacenters, and a Bar CDN, and some servers are running ATS 7.15, 8.0, and 8.1. We make profiles for each server we have fitting all those categories:

EDGE_123_FOO_ABC_ATS_714-RC0-27, EDGE_456_FOO_ABC_ATS_714-RC0-27, EDGE_123_BAR_ABC_ATS_714-RC0-42, EDGE_456_FOO_DEF_ATS_714-RC1-27, EDGE_123_FOO_DEF_ATS_800-RC4-29

- The number of Profiles quickly becomes a combinatorial explosion
- It's nearly impossible for Ops engineers to figure out what Parameters are assigned to all Profiles of a given CDN, machine, or datacenter.
- What about that one random Parameter on a Cloned profile that was changed a year ago? Is it still in place? Should it be? Did it need to be applied to all new Profiles cloned from this one?

## Proposed Change

Layered Profiles allow assigning multiple Profiles, in order, to both Delivery Services and Servers. If multiple Profiles have a Parameter with the same Name and Config File, the Parameter from the last Profile in the ordering is applied.

Layered Profiles is exactly as powerful as the existing Profiles, it doesn't enable any new things. It makes profiles much easier to manage.

With Layered Profiles, hundreds of Profiles become a few dozen, each representing a logical group. For example, a server might have the Profiles, in order:
1. EDGE
2. AMIGA_123
3. CDN_FOO
4. DATACENTER_ABC
5. ATS_714
6. Custom_Hack_42

### Traffic Portal Impact

1. A UI to view all parameters currently applied to a delivery service/server, as well as the profile each parameter came from. A new page, linked from DS and Server pages. For eg:

| Name                                             |Config File               | Value                                | Profile   |
|--------------------------------------------------|--------------------------|--------------------------------------|-----------|
| location                                         | url_sig_myds.config      | /opt/trafficserver/etc/trafficserver | EDGE      |
| Drive_Prefix                                     | storage.config           | /dev/sd                              | AMIGA_123 |
| error_url                                        | url_sig_myotherds.config | 403                                  | EDGE      |
| CONFIG proxy.config.exec_thread.autoconfig.scale | records.config           | FLOAT 1.5                            | ATS_714   |

2. A UI change to add, remove, and reorder (sortable list) profiles for both Delivery Services and Servers, on the existing DS and Server pages.

### Traffic Ops Impact

- Traffic Ops will need to add the logic to the below-mentioned existing API endpoints, in order to show/create/update/delete a sorted list of profiles for delivery services and server.
- `/deliveryservices/ GET, POST`
- `/deliveryservices/{id} PUT, DELETE`
- `/server/  GET, POST`
- `/servers/{id} PUT, DELETE`

#### REST API Impact

- No new endpoints are required

**Existing Endpoints**

- Modify JSON request and response for existing delivery services and servers endpoints.
- JSON **response** with the proposed change will look as follows:

`/deliveryservices?id=100`
```JSON
{ 
  "id": 100,
  ⋮
  "requested": {
    "profileIds": [
      1234,
      5678,
      9012
    ],
    "profileNames": [
      "EDGE",
      "ATS8",
      "TOP"
    ]
    ⋮
  }
}
```

`/servers?id=5`
```JSON
{ 
  "id": 5,
  ⋮
  "requested": {
    "profileIds": [
      1357,
      2468,
      2356
    ],
    "profileNames": [
      "MID",
      "TOP",
      "PRIMARY"
    ],
    ⋮
  }
}
```

JSON **request** with the proposed change will look as follows:

`/deliveryservices`
```JSON
{
  "requested": {
    "profileIds": [
      1234,
      5678,
      9012
    ],
    "profileNames": [
      "EDGE",
      "ATS8",
      "TOP"
    ],
    ⋮
  },
  ⋮
}
```

`/servers`
```JSON
{ 
  "requested": {
    "profileIds": [
      1357, 
      2468, 
      2356
    ],
    "profileNames": [
      "MID", 
      "TOP", 
      "PRIMARY"
    ],
    ⋮
  },
  ⋮
}
```

The following table describes the top level `layered_profile` object for delivery services:

| field              | type          | optionality | description                                                    |
| ------------------ | --------------| ----------- | ---------------------------------------------------------------|
| deliveryservice    | bigint        | required    | the delivery service id associated with a given profile        |
| profile            | bigint        | required    | the profile id associated with a delivery service              |
| order              | bigint        | required    | the order in which a profile is applied to a delivery service  |
| lastUpdated        | bigint        | required    | the last time this delivery service was updated                |

The following table describes the top level `layered_profile` object for servers:

| field       | type                 | optionality | description                                                    |
| ----------- | ---------------------| ----------- | ---------------------------------------------------------------|
| server      | bigint               | required    | the server id associated with a given profile                  |
| profile     | bigint               | required    | the profile id associated with a server                        |
| order       | bigint               | required    | the order in which a profile is applied to a server            |
| lastUpdated | bigint               | required    | the last time this server was updated                          |

**API constraints**
- TO REST APIs will not be backward compatibility and this feature will be a major version (v5.0) change.
- Add `Profiles` key to API endpoints for Server and DeliveryService objects (and ProfileNames, ProfileIDs)

#### Client Impact
- Existing Go client methods will be updated for the `/deliveryservices` and `/servers` endpoints in order to write TO API tests for the exising endpoints.
- All functions which get Parameters on Delivery Services or Servers must be changed to get the Parameters of all assigned Profiles in order.

#### Data Model / Database Impact

- A new Database Table `server_profiles` as described below, will be created:
- A new Database Table `deliveryservice_profiles` as described below, will be created:
```text
            Table "traffic_ops.deliveryservice_profiles"
     Column     |  Type                    | Collation | Nullable | Default
----------------+--------------------------+-----------+----------+--------
deliveryservice | bigint                   |           | not null |
profile         | bigint                   |           | not null |
order           | bigint                   |           | not null |
last_updated    | timestamp with time zone |           | not null | now()
Indexes:
"pk_deliveryservice_profile" PRIMARY KEY(profile)
Foreign-key constraints:
"fk_deliveryservice" FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id)
```

```text
            Table "traffic_ops.server_profiles"
     Column    |  Type                    | Collation | Nullable | Default
---------------+--------------------------+-----------+----------+--------
 server        | bigint                   |           | not null |
 profile       | bigint                   |           | not null |
 order         | bigint                   |           | not null |
 last_updated  | timestamp with time zone |           | not null | now() 
Indexes:
    "pk_server_profile" PRIMARY KEY(profile)
Foreign-key constraints:
    "fk_server" FOREIGN KEY (server) REFERENCES server(id)
```

### ORT Impact
New set of profiles will be created and the profile-parameter relationship will change.

### Traffic Monitor Impact
No impact

### Traffic Router Impact
No impact

### Traffic Stats Impact
No impact

### Traffic Vault Impact
No impact

### Documentation Impact
All existing endpoints will need to be updated, along with the documentation explaining `layered_profile`.

### Testing Impact
Client/API integration tests should be updated to verify the `layered_profile` functionality on existing API `/deliveryservices`, `/servers` endpoints.

### Performance Impact
We do not anticipate any performance impact with layered_profile approach.

### Security Impact
We do not anticipate any impact on security.

### Upgrade Impact
- A Database Migration to:
  - drop profile column in existing deliveryservice and server table
  - insert existing server and DS profiles along with their order into the new tables(deliveryservice_profiles, server_profiles)
- The new capability can just be added to the `seeds.sql` file.

### Operations Impact
The profile-parameter relationship will change based on new sets of profile and Operations will have to learn on how to assign profile order to a delivery service and/or server.

### Developer Impact
Developers will most likely need to use layered_profile, so they'll need to be familiar with the process
of adding, sorting, deleting, debugging and working with layered_profiles.

## Alternatives
None, except to keep using existing Profiles.

## Dependencies
None

## References
None
