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

Currently, within ATC, there is a concept of Invalidation Jobs. These Invalidation Jobs give a user the ability to queue an invalidation for a resource, primarily based on regular expressions. The invalidation is gathered and treated as though there was a cache **STALE**, allowing the CDN to query the origin server to **REFRESH** the resource. However, should the cache policy still be incorrect or misconfigured, the resource could be updated on the origin server, but the CDN will still receive a 304 - Not Modified HTTP status response.

It should be noted that this problem originally arose from a misconfigured origin, however this applies to any parent cache within the topology.

## Proposed Change

To address this potential conflict, a proposal to add **REFETCH** as an option for Invalidation Jobs. This will then be treated by caches as a **MISS** (rather than a **STALE**), thusly allowing the cache to retrieve the resource regardless of cache policies. The original **REFRESH**/**STALE** will be the default option, where **REFETCH**/**MISS** will be the addition.

### Traffic Portal Impact

##### Create and Update
Traffic Portal will need to update the Invalidation Job to account for the different options. When creating an Invalidation Job both options will need to be present (Perhaps a radio button? Default will be the original **REFRESH**).

Tooltips should be added to ensure an understanding of this feature at a high level.

##### Read
When displaying the information, the **Invalidation Requests** table current shows the `Parameters` field, which is only the TTL in hours in the format `TTL:%dh`. The `Parameters` field will be changed to a `TTL` field. This can be displayed directly on the table, however it should be made clear that this is in hours and no other time value.

Additionally, since we derive and calculate the expiration field based on the TTL some code will need to be modified to account for the change in field name and value.

Speaking of Time, time values returned by the server will be formatted to follow RFC3339 per API guidelines.

### Traffic Ops Impact

Both the API and the database schema will likely be updated, which in turn will result in changes downstream (such as T3C/ORT, clients) as well.

#### REST API Impact

No new endpoints will be required. However the current invalidation job will now include an optional field during `Create`. Invalidation jobs are added by submitting a POST to the jobs endpoint.

Globally, a new parameter will be added (recommend something akin to `refetch_enabled`, which defaults to **false**), that will be used to validate any **POSTS** calls to the _/jobs_ endpoint. This will be an initial check done to perform that the CDN is configured to process and return `refetch` jobs. If a `refetch` is submitted and the value of the parameter is **true**, the **POST** will succeed. If the value is `refetch` and the parameter is set to **false**, the **POST** will fail with a 400 - Bad Request. If the value is `refresh`, the parameter need not be checked at all and the **POST** will succeed. 

**POST** /api/4.0/jobs

##### Current Request

Body:
```json
{
	"startTime":"2021-06-02T15:23:21.348Z",
	"deliveryService":11,
	"regex":"/path/.*\\.jpeg",
	"ttl":24
}
```

Which is mapped to a go `struct` in the `go-tc` lib.
```go
type InvalidationJobInput struct {
	DeliveryService *interface{} `json:"deliveryService"`
	Regex           *string      `json:"regex"`
	StartTime       *Time        `json:"startTime"`
	TTL             *interface{} `json:"ttl"`
	dsid            *uint
	ttl             *time.Duration
}
```

##### Proposed

Add an "InvalidationType" to signify a specific type of invalidation request. If the field is included, it _must_ be either "refetch" or "refresh".

Body:
```json
{
	"startTime":"2021-06-02T15:23:21.348Z",
	"deliveryService":11,
	"regex":"/path/.*\\.jpeg",
	"ttl":24,
	"invalidationType":"refresh"
}
```

This struct now contains the `InvalidationType *string` field. Additionally, the `DeliveryService` and `TTL` are no longer empty interfaces. Also `DeliveryService` and `Regex` are no longer optional fields. 

```go
type InvalidationJobInput struct {
	DeliveryService  string  `json:"deliveryService"`
	Regex            string  `json:"regex"`
	InvalidationType *string `json:"invalidationType"`
	StartTime        *Time   `json:"startTime"`
	TTL              *uint   `json:"ttl"`
}
```

##### Parsing the value

The value can only pass validation if it is either explicitly `refresh` or `refetch`. Any other value (including missing/omitted) will be treated as 400 - Bad Content.

##### Response

The response will be modified, then, to return this new value as well. Additionally, _Time_ values will be formatted to RFC3339 to follow API guidelines. This, plus some database changes, will require changes to the struct used for reading the values from the DB.

The current struct:

```go
type InvalidationJob struct {
	AssetURL        *string `json:"assetUrl"`
	CreatedBy       *string `json:"createdBy"`
	DeliveryService *string `json:"deliveryService"`
	ID              *uint64 `json:"id"`
	Keyword         *string `json:"keyword"`
	Parameters      *string `json:"parameters"`
	StartTime       *Time   `json:"startTime"`
}
```

Will be changed to:

```go
type InvalidationJob struct {
	AssetURL         *string `json:"assetUrl"`
	CreatedBy        *string `json:"createdBy"`
	DeliveryService  *string `json:"deliveryService"`
	ID               *uint64 `json:"id"`
	Keyword          *string `json:"keyword"`
	TTL              *int    `json:"ttl"`
	InvalidationType *string `json:"invalidationType"`
	StartTime        *Time   `json:"startTime"`
}
```

Sample current response:
```json
{
	"alerts":[
		{
			"text":"Invalidation request created for http://amc-linear-origin.local.tld/path/.*\\.jpeg, start:2021-06-02 15:23:21.348 +0000 UTC end 2021-06-03 15:23:21.348 +0000 UTC",
			"level":"success"
		}
	],
	"response":{
		"assetUrl":"http://amc-linear-origin.local.tld/path/.*\\.jpeg",
		"createdBy":"admin",
		"deliveryService":"amc-live",
		"id":1,
		"keyword":"PURGE",
		"parameters":"TTL:24h",
		"startTime":"2021-06-02 09:23:21-06"
	}
}
```

Sample new response (includes the `invalidationType` on parameters field, an updated `alert text` field, and a RFC3339 formated startTime):
```json
{
	"alerts":[
		{
			"text":"Invalidation (refresh) request created for http://amc-linear-origin.local.tld/path/.*\\.jpeg, start:2021-06-02 15:23:21.348 +0000 UTC end 2021-06-03 15:23:21.348 +0000 UTC",
			"level":"success"
		}
	],
	"response":{
		"assetUrl":"http://amc-linear-origin.local.tld/path/.*\\.jpeg",
		"createdBy":"admin",
		"deliveryService":"amc-live",
		"id":1,
		"keyword":"PURGE",
		"ttl":24,
		"invalidationType":"refresh",
		"startTime":"2021-06-02T09:23:21-06Z07:00"
	}
}
```

___

> Note: There are still 1.x routes that reference `UserInvalidationJob`, such as 
		`user/current/jobs(/|\.json/?)?$`
		`user/current/jobs(/|\.json/?)?$`
		These routes are currently deprecated and the corresponding `structs` will be removed in a future release as well.

#### Client Impact

Likewise with Traffic Portal, the `go` clients will need to be updated to provide this additional functionality. Since an additional field has been added to `InvalidationJobInput` in `go-tc` lib, this can be set by the client as well.

There is also

#### Data Model / Database Impact

When referring to _Jobs_, these are relegated to three database tables in the TO DB; _jobs_, _job\_agent_, and _job\_status_.

> The _jobs_ concept appears to have been intended to be generic and flexible, however it's only ever been implemented to record invalidation jobs.

The current column `parameters` will be converted to `ttl` and contain an INT datatype representing the TTL in hours. 
```
ttl:48h
```

I propose adding an additional column, _invalidation\_type_. This column will be non-nullable. Default value will be `refresh` (rather than a nullable `NULL` value).

```
invalidation_type: refresh
```

The jobs table will then look more like:

| id | agent | object\_type | object\_name | invalidation\_type | keyword | ttl | asset\_url | asset\_type | status | start\_time | entered\_time | job\_user | last\_updated | job\_deliveryservice |
|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|


*OPTIONAL, BUT RECOMMENDED*: As part of this effort, the _Boy Scout Rule_ will be applied ("Always leave the campground cleaner than you found it."). The `agent`, `status`, `asset_type`, `keyword`, `object_type`, and `object_name` columns will be removed. `agent` and `status` are currently hardcoded to the value of 1 and don't appear to be accessed beyond the INSERT. Similarly, `asset_type` is always "file" and `keyword` is always "PURGE". `object_type` and `object_name` are not used at all. Also, if `status` is removed, the _job\_status_ table can also be removed. Likewise with the _job\_agent_ table if the `job_agent` column is removed. If this were the case, it will also impact the REST API as well since the `keyword` field will not longer be returned. Finally, for clarity, the _jobs_ table could be renamed to _invalidationjob_ since that is the only functionality provided. In the future, should we need to expand the concept of _jobs_ we would still be able to do so.

> The removal of these columns will impact the `UserInvalidationJob` struct. Even though the only endpoint utilizing this struct (v1.1) is deprecated, it is still in use.

### ORT/T3C Impact

The changes required to implement the final functionality have already been implemented within ORT/T3C. The generation of the invalidation regex jobs previously resembled:

```
# regex purgeExpiryTime
refreshasset 1623861151
```

They now include an optional third field that is either **STALE** or **MISS**.

```
# regex purgeExpiryTime (optional)type
refreshasset 1623861151
refreshasset 1623861151 STALE
refreshasset 1623861151 MISS
```

The `regexrevalidatedotconfig.go` in _lib/tc-atscfg_ currently has a function called `filterJobs` that is responsible for parsing the _parameters_ information from the job. Right now, that is parsing the **TTL**. However, this will need to account for the new specific TTL field as well as the new InvalidationType field.

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

The validation in Traffic Ops of the `invalidationType` field will be such that it can only be explicitly set to **refresh** or **refetch**. Any other value (missing, malformed, wrong data type, etc.) will result in either a 400 level error. No other permissions are modified.

> Current permissions require `PrivLevelPortal` to create, update, delete. For read, only `PrivLevelReadOnly` is needed.

### Upgrade Impact

The API for v4 and database schemas will change, however with the addition of the Parameter to ensure the feature is checked for safe guards, there is minimal impact on the upgrades. Clients will continue to create invalidation jobs (`refresh`) as they did before without impacting the caches downstream. However the caches downstream must implement the proper regexrevalidation plugin before enabling the feature via Parameters.

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
