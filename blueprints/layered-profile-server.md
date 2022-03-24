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

Layered Profiles allow assigning multiple profiles to Servers. If multiple profiles have a parameter with the same name and config file, the parameter from the last profile in the ordering is applied.

Layered Profiles is exactly as powerful as the existing profiles, it doesn't enable any new things. It makes profiles much easier to manage.

With Layered Profiles, hundreds of profiles become a few dozen, each representing a logical group. For example, a server might have the profiles, in order:
1. EDGE
2. AMIGA_123
3. CDN_FOO
4. DATACENTER_ABC
5. ATS_714
6. Custom_Hack_42

### Traffic Portal Impact

1. A UI to view all parameters currently applied to a server, as well as the profile each parameter came from. A new page, linked from Server pages. For eg:

| Name                                             |Config File               | Value                                | Profile   |
|--------------------------------------------------|--------------------------|--------------------------------------|-----------|
| location                                         | url_sig_myds.config      | /opt/trafficserver/etc/trafficserver | EDGE      |
| Drive_Prefix                                     | storage.config           | /dev/sd                              | AMIGA_123 |
| error_url                                        | url_sig_myotherds.config | 403                                  | EDGE      |
| CONFIG proxy.config.exec_thread.autoconfig.scale | records.config           | FLOAT 1.5                            | ATS_714   |

2. A UI change to add, remove, and reorder (sortable list) profiles for Servers, on the existing Server pages.
3. Filtering based on Profiles will also need to be updated to take into account the plurality of Profiles.

### Traffic Ops Impact

- Add backend logic to the below-mentioned existing API endpoints, in order to show/create/update/delete a sorted list of profiles for server.

- `/servers/  GET, POST`
- `/servers/{id} PUT, DELETE`
- `/servers/details`
- `/deliveryservices/{{ID}}/servers`
- `/deliveryservices/{{ID}}/servers/eligible`

#### REST API Impact

- No new endpoints are required

**Existing Endpoints**

- Modify JSON request and response for existing servers endpoints.
- **_Note_**: All fields in the response structure for API 4.0 except `profileId`, `profileDesc` and `profile`, remain same. 
  - In API 4.0 `profile` changes to `profiles` and `profileId` and `profileDesc` are no longer part of response structure

#### API 4.0 GET
- JSON **response** with the proposed change will look as follows:
  `/servers?id=5`
```JSON
{
  "response": [{
    "id": 5,
    "profileNames": ["MID", "AMIGA_123", "CDN_FOO"]
  }]
}
```

`/servers`
```JSON
{
  "response": [{
    "profileNames": ["MID", "AMIGA_123", "CDN_FOO"]
  }]
}
```

#### API 4.0 POST/PUT
JSON **request** with the proposed change will look as follows:

`POST /servers `
```JSON
{
    "cachegroupId": 6,
    "cdnId": 2,
    "profileNames": ["MID", "AMIGA_123", "CDN_FOO"]
}
```

`PUT /servers/5`
```JSON
{
    "cachegroupId": 6,
    "cdnId": 2,
    "profileNames": ["MID", "AMIGA_123", "CDN_FOO"]
}
```

return following **response**
```JSON
{ "alerts": [
  {
    "text": "Server created/updated",
    "level": "success"
  }
],
  "response": {
    "cachegroup": "CDN_in_a_Box_Mid",
    "cachegroupId": 6,
    "cdnId": 2,
    "id": 5,
    "profileNames": ["MID", "AMIGA_123", "CDN_FOO"]
  }
}
```

The following table describes the top level `layered_profile` object for servers:

| field           | type                 | optionality | description                                              |
| ----------------| ---------------------| ----------- | ---------------------------------------------------------|
| server          | bigint               | required    | the server id associated with a given profile            |
| profile_name    | text                 | required    | the profile name associated with a server                |
| order           | bigint               | required    | the order in which a profile is applied to a server      |

**API constraints**
- In API 4.0, 
  - GET `/servers` object, `profiles` field is an array, 
  - and a PUT or POST `/servers` allows multiple profiles to be assigned, and displayed in the UI, etc.
    The UI may or may not display a list (which can only add 1), but the client implements handling multiple, and the API is documented to potentially return multiple and how their parameters must be applied.
- Add `Profiles` key to API endpoints for Server objects

#### Client Impact
- Existing Go client methods will be updated for the `/servers` endpoints in order to write TO API tests for the exising endpoints.
- All functions which get Parameters on Servers must be changed to get the Parameters of all assigned Profiles in order.

#### Data Model / Database Impact

- A new Database Table `server_profiles` as described below, will be created:
```text
            Table "traffic_ops.server_profiles"
     Column    |  Type                    | Collation | Nullable | Default
---------------+--------------------------+-----------+----------+--------
 server        | bigint                   |           | not null |
 profile_name  | text                     |           | not null |
 order         | bigint                   |           | not null |
Indexes:
    "pk_server_profile" PRIMARY KEY(profile_name, server)
Foreign-key constraints:
    "fk_server_id" FOREIGN KEY (server) REFERENCES public.server(id) ON DELETE CASCADE ON UPDATE CASCADE
    "fk_server_profile_name_profile" FOREIGN KEY (profile_name) REFERENCES public.profile(name) ON DELETE RESTRICT ON UPDATE CASCADE,
```

All profiles assigned to a given server will have the same values of:
- type
- cdn
If any of those differ within the same server, it's probably a mistake.

### ORT Impact
New set of profiles will be created and the profile-parameter relationship will change.

### Traffic Monitor Impact
No impact

### Traffic Router Impact
These changes will affect the Snapshot generation (both crconfig and monitoring). Even though that is more of a TO impact.
Reason being that Snapshots and Monitoring Configurations for a CDN include Profile and Parameter information for the servers, Traffic Monitors, and Traffic Routersz.

### Traffic Stats Impact
No impact

### Traffic Vault Impact
No impact

### Documentation Impact
All existing endpoints will need to be updated, along with the documentation explaining `layered_profile`.

### Testing Impact
Client/API integration tests should be updated to verify the `layered_profile` functionality on existing API `/servers` endpoints.

### Performance Impact
We do not anticipate any performance impact with layered_profile approach.

### Security Impact
We do not anticipate any impact on security.

### Upgrade Impact
- A Database Migration to:
  - insert existing server profiles along with their order into the new table(server_profiles)

### Operations Impact
The profile-parameter relationship will change based on new sets of profile and Operations will have to learn on how to assign profile order to a server.

### Developer Impact
Developers will most likely need to use layered_profile, so they'll need to be familiar with the process
of adding, sorting, deleting, debugging and working with layered_profiles.
When searching for any parameters assigned to a profile for a given server, one will need to look up all profiles assigned to the server and then process all the parameters in order to get the "final" view of the profile.

## Alternatives
None, except to keep using existing Profiles.

## Dependencies
None

## References
None
