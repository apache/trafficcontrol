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
# Add REFETCH capability option for Content Invalidation

## Problem Description

Currently, within ATC, there is a concept of Invalidation Jobs. These Invalidation Jobs give a user the ability to queue an invalidation for a resource, primarily based on regular expressions. The invalidation is gathered and treated as though there was a cache **STALE**, allowing the cache to query the origin server to **REFRESH** the resource. However, should the cache policy still be incorrect or misconfigured, the resource could be updated on the origin server, but the cache will still receive a 304 - Not Modified HTTP status response.

## Proposed Change

 To address this potential conflict, a proposal to add **REFETCH** as an option for Invalidation Jobs. This will then be treated by caches as a **MISS** (rather than a **STALE**), thusly allowing the cache to retrieve the resource regardless of cache policies. The original **REFRESH**/**STALE** will be the default option, where **REFETCH**/**MISS** will be the addition.

### Traffic Portal Impact

##### Create and Update
Traffic Portal will need to update the Invalidation Job to account for the different options. When creating an Invalidation Job both options will need to be present (Perhaps a radio button? Default will be the original **REFRESH**).

Tooltips should be added to ensure an understanding of this feature at a high level.

##### Read
When displaying the information, the **Invalidation Requests** table current shows the `Parameters` field, so the display will be straight forward with no manipulation.

However, we derived and calculate the expiration field based on the TTL. This code will need to be modified to account for the additional information contained in the `Parameters` field.

### Traffic Ops Impact

Both the API and the database schema will likely be updated, which in turn will result in changes downstream (such as T3C/ORT, clients) as well.

#### REST API Impact

No new endpoints will be required. However the current invalidation job will now include an optional field during `Create`. Invalidation jobs are added by submitting a POST to the jobs endpoint. 

**POST** /api/4.0/jobs

##### Current Request

Body:
```
{"startTime":"2021-06-02T15:23:21.348Z","deliveryService":11,"regex":"/path/.*\\.jpeg","ttl":24}
```

Which is mapped to a go `struct` in the `go-tc` lib.
```go
type InvalidationJobInput struct {
	DeliveryService *interface{} `json:"deliveryService"`
	Regex *string `json:"regex"`
	StartTime *Time `json:"startTime"`
	TTL *interface{} `json:"ttl"`
	dsid *uint
	ttl  *time.Duration
}
```

##### Proposed

Add an "InvalidationType" to signify a specific type of invalidation request. The InvalidationType is an optional field and will not break backwards compatibility with previous API versions. If the field is included, it _must_ be either "refetch" or "refresh".

Body:
```
{"startTime":"2021-06-02T15:23:21.348Z","deliveryService":11,"regex":"/path/.*\\.jpeg","ttl":24,"invalidationType":"refresh"}
```

This struct now contains the `InvalidationType *string` field.
```go
type InvalidationJobInput struct {
	DeliveryService *interface{} `json:"deliveryService"`
	Regex *string `json:"regex"`
	InvalidationType *string `json:"invalidationType,omitempty"`
	StartTime *Time `json:"startTime"`
	TTL *interface{} `json:"ttl"`
	dsid *uint
	ttl  *time.Duration
}
```

##### Parsing the value

Since the field is optional and existing functionality only signifies a **REFRESH**/**STALE** capability, if the field is omitted, empty, malformed, or in any way _not_ `refetch` then it will be treated as `refresh`.

##### Response

The response will be modified, then, to return this new value as well.

Sample current response:
```
{"alerts":[{"text":"Invalidation request created for http://amc-linear-origin.local.tld/path/.*\\.jpeg, start:2021-06-02 15:23:21.348 +0000 UTC end 2021-06-03 15:23:21.348 +0000 UTC","level":"success"}],"response":{"assetUrl":"http://amc-linear-origin.local.tld/path/.*\\.jpeg","createdBy":"admin","deliveryService":"amc-live","id":1,"keyword":"PURGE","parameters":"TTL:24h""startTime":"2021-06-02 09:23:21-06"}}
```

Sample new response (includes the `invalidationType` on parameters field):
```
{"alerts":[{"text":"Invalidation request created for http://amc-linear-origin.local.tld/path/.*\\.jpeg, start:2021-06-02 15:23:21.348 +0000 UTC end 2021-06-03 15:23:21.348 +0000 UTC","level":"success"}],"response":{"assetUrl":"http://amc-linear-origin.local.tld/path/.*\\.jpeg","createdBy":"admin","deliveryService":"amc-live","id":1,"keyword":"PURGE","parameters":"TTL:24h,invalidationType:refresh","startTime":"2021-06-02 09:23:21-06"}}
```

___

> Note: There are still 1.x routes that reference `UserInvalidationJob`, such as 
		`user/current/jobs(/|\.json/?)?$`
		`user/current/jobs(/|\.json/?)?$`
		These routes are currently deprecated and the corresponding `structs` will be removed in a future release as well.

#### Client Impact

Likewise with Traffic Portal, the `go` clients will need to be updated to provide this additional functionality. Since an additional field has been added to `InvalidationJobInput` in `go-tc` lib, this can be set by the client as well. If left unset, it will default to "refresh".

#### Data Model / Database Impact


The current column `parameters` will now contain a cskv (comma separated key value) string. Currently it only stores the `TTL` for the invalidation request:
```
TTL:48h
```

Moving forward, this column will also contain the type of cache invalidation. For instance, the string may read:
```
TTL:48h,invalidationType:refetch
```

If there is no `invalidationType` in the cskv, it is assumed to be **REFRESH**/**STALE** as it's default value keeping with the current implementation. Otherwise the `invalidationType` will be either `refetch` or `refresh`, defaulting to `refresh` during validation.

*OPTIONAL, BUT RECOMMENDED*: As part of this effort, the _Boy Scout Rule_ will be applied ("Always leave the campground cleaner than you found it."). The `agent`, `status`, `asset_type`, `object_type`, and `object_name` columns will be removed. `agent` and `status` are currently hardcoded to the value of 1. Similarly, `asset_type` is never a anything besides "file". `object_type` and `object_name` are not used at all. This would require a DB migration with Goose as well. If this were not implemented, no migration is necessary.

> The removal of these columns will impact the `UserInvalidationJob` struct. Even though the only endpoint utilizing this struct (v1.1) is deprecated, it is still in use.

### ORT/T3C Impact

The `regexrevalidatedotconfig.go` in _lib/tc-atscfg_ currently has a function called `filterJobs` that is responsible for parsing the _parameters_ information from the job. Right now, that is only parsing the **TTL**, however it would then also need to parse the optional, extra field of **invalidationType**. 

```go
type revalJob struct {
	AssetURL string
	PurgeEnd time.Time
	Type     string // MISS or STALE (default)
}
```

As we can see above, if the value **invalidationType** is missing or equal to **refresh**, then the `revalJob` struct's `Type` field will be set to **STALE**. Otherwise, if the **invalidationType** is **refetch** then the `Type` field will be set to **MISS**. The struct is already prepared for this information, it needs only to be parsed by during the `filterJobs` function call.


### Traffic Monitor Impact

N/A - No changes

### Traffic Router Impact

N/A - No changes

### Traffic Stats Impact

N/A - No changes

### Traffic Vault Impact

N/A - No changes

### Documentation Impact

This information will need to be added to the current **Forcing Content Invalidation** section under the **_Quick How To Guides_** section under the **_Administrator's Guide_**
[Content Invalidation](https://traffic-control-cdn.readthedocs.io/en/latest/admin/quick_howto/content_invalidation.html)

Additionally, the **_Traffic Ops API_** `jobs` routes will need to be updated with the changes, API V1-V4.
[V4 Jobs](https://traffic-control-cdn.readthedocs.io/en/latest/admin/quick_howto/content_invalidation.html)

### Testing Impact

##### Unit tests

For ORT/T3C there are already unit tests that will need to be updated to account for this contingency.

Update:
```
github.com/apache/trafficcontrol/lib/go-atscfg/regexrevalidatedotconfig_test.go
```

There are no unit tests for invalidation jobs in `traffic_ops_golang`. This provides an opportunity to create unit tests to validate current and new functionality.

Add:
```
github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/invalidationjobs/invalidationjobs.go
```

##### Integration/E2E Tests:
There are already existing integration tests for the various APIs (v1-v4) for Traffic Ops. Each will need to have this optional functionality tested as well.
```
github.com/apache/trafficcontrol/traffic_ops/testing/api/v1/jobs_test.go
github.com/apache/trafficcontrol/traffic_ops/testing/api/v2/jobs_test.go
github.com/apache/trafficcontrol/traffic_ops/testing/api/v3/jobs_test.go
github.com/apache/trafficcontrol/traffic_ops/testing/api/v4/jobs_test.go
```

### Performance Impact

There will be no performance impact for Traffic Control.

> Note: This functionality may create a performance impact on caches that implement a REFETCH/MISS manual override based on a regex.

### Security Impact

The validation in Traffic Ops of the `invalidationType` field will be such that it can only be explicitly set to **refetch**. Any other value (missing, malformed, wrong data type, etc.) will result in either a 400 level error or a default to **refresh**. No other permissions are modified.

> Current permissions require `PrivLevelPortal` to create, update, delete. For read, only `PrivLevelReadOnly` is needed.

### Upgrade Impact

Unless the database schema is changed, there will be no required migration. Since the field is optional and defaults to a value, any previous data will remain unaltered.

Those utilizing the clients will need to update to be able to utilize the new type of invalidation job.

Once an upgrade is complete, manual test can verify the changes were done correctly (for example, either through code utilizing a client or through the traffic portal interface)

### Operations Impact

Operators should be made aware in the documentation of the potential performance hit the cache might experience by using **reFETCH** resulting in a **MISS** over **reFRESH** resulting in **STALE** (Default).

### Developer Impact

There will be no impact for developers moving forward. If the database columns were removed/cleaned up it may lighten a slight cognitive load since the fields aren't representative or used in the non-deprecated implementation.

## Alternatives

Utilizing the currently existing `parameters` field will require some code modifications, but it will be minimally invasive. Another field could be added and become an additional column to the database schema. It would then be returned in the response object from the API, however the primary concern is that the generic `jobs` table to now to be specifically aware of invalidation jobs (it appears to have been written generically originally) and this may result in unintended consequences for clients.

## Dependencies

N/A - No changes
