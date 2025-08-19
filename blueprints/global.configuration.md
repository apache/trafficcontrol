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
# Global Configuration

## Problem Description
Currently, a lot of global configuration for Traffic Ops is handled by a set of
Parameters which may or may not exist at any given time, and may or may not have
a valid value for the configuration the represent, *and* may or may not have the
Config File value "global", **and** may or may not be assigned to the a Profile
named "GLOBAL". This is annoying, lacks validation, doesn't allow users to be
aware of the current/default configuration (or even which value is used as in
many cases the actual Parameter value used is undefined).

## Proposed Change
The Global Configuration singleton object will contain information pertaining to
the configuration of the Traffic Control system as a whole - most notably it
stores Traffic Ops API access information and provides information for
geographically locating clients of Apache Traffic Control CDNs.

## Data Model Impact
<a name="sec:data-model"></a>
```typescript
interface GlobalConfiguration {
	/**
	 * The default latitude to use for CDN clients when geo-location fails, on
	 * the interval [-90, 90].
	 */
	defaultGeoMissLatitude: number;
	/**
	 * The default longitude to use for CDN clients when geo-location fails, on
	 * the interval [-180, 180].
	 */
	defaultGeoMissLongitude: number;
	/**
	 * The location of an IPv4 geolocation database (replaces
	 * `geolocation.polling.url` and `alt.geolocation.polling.url`).
	 */
	geolocationIPv4PollingURL: URL;
	/**
	 * The location of an IPv6 geolocation database (replaces
	 * `geolocation6.polling.url`). If `null`, `geolocationIPv4PollingURL` is
	 * assumed to be the location of a database that's either capable of
	 * localizing both IPv4 and IPv6 addresses, or only IPv4 addresses but no
	 * IPv6 localization is needed/desired.
	 */
	geolocationIPv6PollingURL: URL | null;
	/**
	 * A URL at which information about the ATC system may be found (replaces
	 * `tm.info_url`).
	 */
	informationURL: URL | null;
	/**
	 * The name of the ATC instance (replaces `tm.instance_name`).
	 */
	instanceName: string;
	/**
	 * The maximum duration - in days - of content invalidation jobs. `null`
	 * means "no limit" (replaces `maxRevalDurationDays`).
	 */
	maxRevalidationDays: bigint | null; // >= 0
	/**
	 * A URL for a reverse proxy to the Traffic Ops instance(s) (replaces
	 * `tm.rev_proxy.url`).
	 */
	reverseProxyURL: URL & {search: "", pathname: "/"} | null;
	/**
	 * The name of the tool that serves the Traffic Ops API (replaces
	 * `tm.toolname`).
	 */
	toolName: string;
	/**
	 * The URL of a Proxy that should be used instead of making requests
	 * directly to Traffic Monitors when serving Traffic Ops API endpoints that
	 * expose Traffic Monitor data (replaces `tm.traffic_mon_fwd_proxy`).
	 */
	trafficMonitorProxy: URL & {search: "", pathname: "/"} | null;
	/**
	 * The Name of the Status of Traffic Monitor that will be used as the source
	 * of truth for the availability of cache servers (replaces
	 * `tm_query_status_override`).
	 *
	 * @default "ONLINE"
	 */
	trafficMonitorQueryStatus: string;
	/**
	 * The URL of a Proxy that should be used instead of making requests
	 * directly to Traffic Routers when serving Traffic Ops API endpoints that
	 * expose Traffic Monitor data (replaces `tm.traffic_rtr_fwd_proxy`).
	 */
	trafficRouterProxy: URL & {search: "", pathname: "/"} | null;
	/**
	 * The canonical URL used to make requests to the Traffic Ops API (replaces
	 * `tm.url`).
	 */
	trafficOpsURL: URL & {search: "", pathname: "/"};
}
```

`tm.toolname` and `tm.info_url` don't appear to actually be used by anything. So
I'm not sure if the `toolName` and `informationURL` properties are actually
necessary; hopefully the review of this blueprint will shed some light on that.
Furthermore, it seems that `geolocation6.polling.url` may no longer be used.

## Component Impact
This change can have potentially far-reaching effects. This will affect T3C,
Traffic Portal, and Traffic Ops, primarily, but through CDN Snapshots it'll also
indirectly affect Traffic Router and Traffic Monitor is technically unaffected
but access to Traffic Monitor through the Traffic Ops API is.

### Traffic Portal Impact
Traffic Portal will need a new UI view for setting configuration values. This
should be very straightforward. A sample is provided below, as it would appear
in TPv2. The code is available in [Appendix A](#appendix).

![](img/global.configuration.example.png "Mockup of proposed Global Configuration controls")


### Traffic Ops Impact
Traffic Ops will need to not only add a new endpoint for interacting with the
Global Configuration, but it'll also need to rework existing endpoints to use it
rather than the old Parameters - mostly that'll affect Content Invalidation Job
creation and CDN Snapshots.

#### Database Impact
To store the new data, the following schema is suggested:

```
                                                     Table "public.global_configuration"
            Column            |  Type   | Collation | Nullable |                                   Default
------------------------------+---------+-----------+----------+------------------------------------------------------------------------------
 default_geo_miss_latitude    | numeric |           | not null | 0.0
 default_geo_miss_longitude   | numeric |           | not null | 0.0
 geolocation_ipv4_polling_url | text    |           | not null | 'file:///opt/traffic_router/etc/traffic_router/Geo_lite2-City.mmdb.gz'::text
 geolocation_ipv6_polling_url | text    |           |          |
 information_url              | text    |           |          |
 instance_name                | text    |           | not null | 'Traffic Ops'::text
 max_revalidation_days        | bigint  |           |          |
 reverse_proxy_url            | text    |           |          |
 tool_name                    | text    |           | not null | 'Traffic Ops'::text
 traffic_monitor_proxy        | text    |           |          |
 traffic_monitor_query_status | text    |           | not null | 'ONLINE'::text
 traffic_router_proxy         | text    |           |          |
 traffic_ops_url              | text    |           | not null | 'https://trafficops.infra.ciab.test/'::text
Indexes:
    "global_configuration_bool_idx" UNIQUE, btree ((true))
Check constraints:
    "valid_latitude" CHECK (default_geo_miss_latitude <= 90::numeric AND default_geo_miss_latitude >= '-90'::integer::numeric)
    "valid_longitude" CHECK (default_geo_miss_longitude <= 180::numeric AND default_geo_miss_longitude >= '-180'::integer::numeric)
    "valid_max_revalidation_days" CHECK (max_revalidation_days >= 0)
```

Note the unique index on a constant expression - `true` - which prevents
inserting a second row.

#### API Impact
The API will need a new endpoint for dealing with the Global Configuration.

##### `/configuration`
Since the Global Configuration exists ab initio and always exists, the endpoint
need only support GET and PUT.

###### `GET`

Authentication Required: Yes

Required Permissions:    `CONFIGURATION:READ`

Response Type:           Object


Response Structure - refer to <a href="#sec:data-model">Data Model Impact</a>.

###### `PUT`

Authentication Required: Yes

Required Permissions:    `CONFIGURATION:READ`, `CONFIGURATION:UPDATE`

Response Type:           Object


Request Structure - refer to <a href="#sec:data-model">Data Model Impact</a>.

Response Structure - refer to <a href="#sec:data-model">Data Model Impact</a>.

##### `/caches/stats`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "caches stats". It will instead need to use
the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

##### `/cdns/capacity`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "cdns capacity". It will instead need to use
the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

##### `/cdns/health`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "cdns health". It will instead need to use
the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

#### `/cdns/routing`
When choosing a Traffic Router from which to retrieve "cdns routing"
information, this endpoint will randomly choose the value of a Parameter with
the Config File "global" and the Name `tm.traffic_rtr_fwd_proxy` to use as a
proxy for requests that would normally go directly to some Traffic Router. It
must now instead used the `trafficRouterProxy` Global Configuration setting.

##### `/cdns/{{Name}}/health`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "cdn health". It will instead need to use
the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

##### `/cdns/{{Name}}/snapshot`
Currently, this endpoint uses some randomly selected Parameter with the Config
File "CRConfig.json" to set the `edge.http.routing`, `edge.dns.routing`,
`dns.consistent.routing`, `dnssec.enabled`, `dnssec.zone.diffing.enabled`,
`dnssec.rrsig.cache.enabled`, **and** `strip.special.query.params` properties of
the Snapshot's `config` property from Parameters of the same Name. These need
not be assigned to any particular Profile, but must be assigned to some Profile
within the CDN for which the Snapshot is generated (which is to say most likely
not the "GLOBAL" Profile as that is typically in the "ALL" CDN which will break
in various undocumented scenarios if you try to actually use it as a CDN). All
of these settings ought to be set on Traffic Router Profiles (and really ought
to be first-class properties thereof so they can be properly validated and given
proper, visible default values, but that's beyond this blueprint's scope). The
recommendation this blueprint will make is that the endpoint be changed to
consider Traffic Router Profiles, and/or Traffic Routers be changed to pick up
the value from its own Profile's Parameter list instead of this quasi-global
configuration, but technically no changes are necessary to implement Global
Configuration because these are merely quasi-global Parameter settings.

This endpoint *also* uses randomly selected values for the truly "global"
Parameters `geolocation.polling.url` and `alt.geolocation.polling.url`. The
information used to fill in the former should be sourced now from the Global
Configuration, while the latter appears to be undocumented and support for it
should just be removed entirely.

Finally, this endpoint creates Snapshot structures by filling in the `tm_host`
property of the Snapshot's `stats` property by randomly selecting the Value of a
Parameter with the Name `tm.url`. This should instead use the `trafficOpsURL`
Global Configuration setting.

##### `/cdns/{{Name}}/snapshot/new`
Currently, this endpoint uses some randomly selected Parameter with the Config
File "CRConfig.json" to set the `edge.http.routing`, `edge.dns.routing`,
`dns.consistent.routing`, `dnssec.enabled`, `dnssec.zone.diffing.enabled`,
`dnssec.rrsig.cache.enabled`, **and** `strip.special.query.params` properties of
the Snapshot's `config` property from Parameters of the same Name. These need
not be assigned to any particular Profile, but must be assigned to some Profile
within the CDN for which the Snapshot is generated (which is to say most likely
not the "GLOBAL" Profile as that is typically in the "ALL" CDN which will break
in various undocumented scenarios if you try to actually use it as a CDN). All
of these settings ought to be set on Traffic Router Profiles (and really ought
to be first-class properties thereof so they can be properly validated and given
proper, visible default values, but that's beyond this blueprint's scope). The
recommendation this blueprint will make is that the endpoint be changed to
consider Traffic Router Profiles, and/or Traffic Routers be changed to pick up
the value from its own Profile's Parameter list instead of this quasi-global
configuration, but technically no changes are necessary to implement Global
Configuration because these are merely quasi-global Parameter settings.

This endpoint *also* uses randomly selected values for the truly "global"
Parameters `geolocation.polling.url` and `alt.geolocation.polling.url`. The
information used to fill in the former should be sourced now from the Global
Configuration, while the latter appears to be undocumented and support for it
should just be removed entirely.

Finally, this endpoint creates Snapshot structures by filling in the `tm_host`
property of the Snapshot's `stats` property by randomly selecting the Value of a
Parameter with the Name `tm.url`. This should instead use the `trafficOpsURL`
Global Configuration setting.

##### `/deliveryservices/{{ID}}/capacity`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "Delivery Service capacity". It will instead
need to use the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

##### `/deliveryservices/{{ID}}/health`
Currently uses the "global" Parameter `tm_query_status_override` to determine
which Traffic Monitors to query for "Delivery Service health". It will instead
need to use the `trafficMonitorQueryStatus` setting of the Global Configuration.

Also uses the "global" Parameter `tm.traffic_mon_fwd_proxy` to make requests to
a forward proxy for Traffic Monitor services instead of actual Traffic Monitor
URLs, when configured. It will instead need to use the `trafficMonitorProxy`
Global Configuration setting.

#### `/deliveryservices/{{ID}}/routing`
When choosing a Traffic Router from which to retrieve "Delivery Service routing"
information, this endpoint will randomly choose the value of a Parameter with
the Config File "global" and the Name `tm.traffic_rtr_fwd_proxy` to use as a
proxy for requests that would normally go directly to some Traffic Router. It
must now instead used the `trafficRouterProxy` Global Configuration setting.

##### `/jobs`
Job creation (`POST` method) checks the "global" Parameter
`maxRevalDurationDays` to determine if a submitted Content Invalidation Job is
valid with respect to its duration. It will instead need to use the
`maxRevalDurationDays` Global Configuration property.

Job creation (`POST` method) checks the "global" Parameter `use_reval_pending`
to determine how to set revalidation flags. In the case that this Parameter
exists and the randomly chosen Value is exactly `"0"`, it will cause a CDN-wide
"Queue Updates". This is probably never actually desirable. The endpoint should
instead always use "Rapid Revalidation" and simply deprecated and eventually
remove its behavior's dependence on this "global" Parameter.

When a client attempts to create a Content Invalidation Job (i.e. makes a `POST`
request) with the "REFETCH" type, this endpoint checks for the value of a
Parameter with the name `refetch_enabled` and the Config File "global". If more
than one exists, the value used is chosen essentially at random. Currently, the
default behavior when no such Parameter exists is to act as though it were
"false". Now that the Parameter has been introduced in a major release, the
behavior should be changed to instead act as though it were "true" by default
and deprecate the use of the Parameter to disable it, as in the future REFETCH
will be guaranteed to be supported by any supported T3C version.

##### `/servers/{{Host Name}}/update-status`
This endpoint uses a Parameter with the name `use_reval_pending` to set the
property of the same name on the array response elements it returns. This
property should be deprecated (and in fact because host names are not unique we
really ought to scrap this entire endpoint as it can't possibly work in a
well-defined way).

##### `/snapshot`
Currently, this endpoint uses some randomly selected Parameter with the Config
File "CRConfig.json" to set the `edge.http.routing`, `edge.dns.routing`,
`dns.consistent.routing`, `dnssec.enabled`, `dnssec.zone.diffing.enabled`,
`dnssec.rrsig.cache.enabled`, **and** `strip.special.query.params` properties of
the Snapshot's `config` property from Parameters of the same Name. These need
not be assigned to any particular Profile, but must be assigned to some Profile
within the CDN for which the Snapshot is generated (which is to say most likely
not the "GLOBAL" Profile as that is typically in the "ALL" CDN which will break
in various undocumented scenarios if you try to actually use it as a CDN). All
of these settings ought to be set on Traffic Router Profiles (and really ought
to be first-class properties thereof so they can be properly validated and given
proper, visible default values, but that's beyond this blueprint's scope). The
recommendation this blueprint will make is that the endpoint be changed to
consider Traffic Router Profiles, and/or Traffic Routers be changed to pick up
the value from its own Profile's Parameter list instead of this quasi-global
configuration, but technically no changes are necessary to implement Global
Configuration because these are merely quasi-global Parameter settings.

This endpoint *also* uses randomly selected values for the truly "global"
Parameters `geolocation.polling.url` and `alt.geolocation.polling.url`. The
information used to fill in the former should be sourced now from the Global
Configuration, while the latter appears to be undocumented and support for it
should just be removed entirely.

Finally, this endpoint creates Snapshot structures by filling in the `tm_host`
property of the Snapshot's `stats` property by randomly selecting the Value of a
Parameter with the Name `tm.url`. This should instead use the `trafficOpsURL`
Global Configuration setting.

##### `/system/info`
This endpoint currently does nothing more than list the Names and Values of all
Parameters with the Config File value of exactly `"global"`. These Parameters
may or may not be assigned to some Profile which may or may not have the name
"GLOBAL", and if there are multiple Parameters in that Config File with the same
Name, the Value shown by this endpoint is not defined.

This endpoint should be deprecated and removed as soon as possible once Global
Configuration is introduced.

##### `/user/reset_password`
This endpoint uses the Value of a randomly selected Parameter with the Name
`tm.instance_name` and the Config File `"global"` to be inserted into emails for
password resets. It must instead use the `instanceName` Global Configuration
setting. This will also alleviate the problem where users can never reset their
passwords if that Parameter doesn't exist (results in a `500 Internal Server
Error` response).

##### `/users/register`
This endpoint uses the Value of a randomly selected Parameter with the Name
`tm.instance_name` and the Config File `"global"` to be inserted into emails for
user registration. It must instead use the `instanceName` Global Configuration
setting. This will also alleviate the problem where users can never be
registered if that Parameter doesn't exist (results in a `500 Internal Server
Error` response).

#### Client Impact
The Go and Python clients will need methods added for interacting with the
Global Configuration.

#### Other
##### Check Script(s) Impact
The `ToORTCheck.pl` script currently checks for a randomly chosen `tm.url`
"global" Parameter to find out what URL to use for requests. Besides the request
used to find that out, of course. Frankly I don't even know if that script still
works, but if it does and people care about it it'll need to be update to
instead use the `trafficOpsURL` Global Configuration Setting instead, and fall
back on the old Parameter way in the event that it's `null` (but that behavior
should be deprecated and removed as soon as possible).

##### Postinstall Impact
Postinstall sets a number of "global" Parameter values, and must be changed to
instead set Global Configuration values.

### T3C Impact
T3C currently uses a Parameter named `tm.rev_proxy.url` assigned to a Profile
named "GLOBAL" (which may or may not exist and may or may not have multiple
`tm.rev_proxy.url` Parameters assigned under arbitrary Config Files with
possibly conflicting Values) to determine if it should use a proxy for requests
to Traffic Ops. It will instead need to use the `reverseProxyURL` property of
the Global Configuration.

It also checks a Parameter with the Name `maxRevalDurationDays` chosen at random
from any such Parameters assigned to any existing Profile with the Name "GLOBAL"
as a cap on Content Invalidation Jobs. It will instead need to check the Global
Configuration setting of the same name.

T3C also uses a randomly selected "global" Parameter named `tm.url` to determine
the URL to use when making requests to Traffic Ops (assuming there is no reverse
proxy configured for use instead) and that will need to change to the
`trafficOpsURL` Global Configuration setting.

If, however, the value of any property/setting is `null`, it should fall back
(for now) on the old Parameter, during a deprecation period for said Parameter.

Lastly, T3C will randomly select from the available "global" Parameters named
`use_reval_pending` to check whether to wait for parents to apply updates before
applying updates itself. This behavior should be deprecated and eventually T3C
will only apply "Revalidations" "Rapidly".

### Traffic Monitor Impact
Traffic Monitor itself doesn't need to care about the Global Configuration, it
is only impacted transitively by the changes to endpoints that request Traffic
Monitor data.

### Traffic Router Impact
Traffic Router itself doesn't need to care about the Global Configuration, it
is only impacted transitively by the changes to endpoints that request Traffic
Router data.

### Traffic Stats Impact
None.

### Traffic Vault Impact
None.

## Documentation Impact
Besides API documentations for new and newly deprecated routes, an overview
section should be added for Global Configuration and all of the myriad of
references to "global" Parameters should be removed/reworked into links to the
configuration's docs.

## Testing Impact
Tests for each component will need to change to reflect their new behavior.
Mostly I expect T3C tests will need updating to account for the new setup
process required to put equivalent data into the system, and Traffic Ops tests
will likewise need to change their database mocking to reflect the new object.

## Performance Impact
No performance impact is expected.

## Security Impact
No security impact is expected.

## Upgrade Impact
The database migration should set the values of the configuration settings
according to the values of the Parameters they are replacing. So on upgrade the
behavior should be unchanged by design.

## Operations Impact
Operators will need to be aware of the new configuration options and be informed
that they should no longer use the old "global" Parameters. Ideally with the
Traffic Portal UI controls this should be a very smooth transition, putting all
such information in a single place with definite values and hints for input.

## Developer Impact
Developers will need to know to use the Global Configuration instead of "global"
Parameters in future applications (not that I'm anticipating any in general).

## Alternatives
One alternative would be to move these Parameters to the Profile of Traffic Ops
itself, but that has two main disadvantages:

1. Parameter Values are unvalidated
2. What do you do when Parameter Values for multiple Traffic Ops instances
conflict?

To say nothing of the fact that this doesn't solve the inherent randomness of
duplicate Parameters being possible.

## Dependencies
None.

## References
N/A

## Appendix A - Code samples
<a name="appendix"></a>
### Traffic Portal v2 Mock Markup
```html
<style type="text/css">
	mat-card {
		max-width: 500px;
		margin: auto;
	}
	form {
		display: grid;
		grid-template-columns: 1fr;
	}
	footer {
		display: inline-flex;
		justify-content: space-between;
		padding-top: 16px;
	}
</style>

<mat-card>
<form>
<mat-form-field>
	<mat-label>Default Geo Miss Latitude</mat-label>
	<input matInput name="defaultGeoMissLatitude" type="number" required min="-90" max="90"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Default Geo Miss Longitude</mat-label>
	<input matInput name="defaultGeoMissLongitude" type="number" required min="-180" max="180"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Geolocation IPv4Polling URL</mat-label>
	<input matInput name="geolocationIPv4PollingURL" type="url" required/>
</mat-form-field>
<mat-form-field>
	<mat-label>Geolocation IPv6Polling URL</mat-label>
	<input matInput name="geolocationIPv6PollingURL" type="url"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Information URL</mat-label>
	<input matInput name="informationURL" type="url"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Instance Name</mat-label>
	<input matInput name="instanceName" type="text" required/>
</mat-form-field>
<mat-form-field>
	<mat-label>Max Revalidation Days</mat-label>
	<input matInput name="maxRevalidationDays" type="number" min="0"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Reverse Proxy URL</mat-label>
	<input matInput name="reverseProxyURL" type="url"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Tool Name</mat-label>
	<input matInput name="toolName" type="string" required/>
</mat-form-field>
<mat-form-field>
	<mat-label>Traffic Monitor Proxy</mat-label>
	<input matInput name="trafficMonitorProxy" type="url"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Traffic Monitor Query Status</mat-label>
	<input matInput name="trafficMonitorQueryStatus" type="string" required/>
</mat-form-field>
<mat-form-field>
	<mat-label>Traffic Router Proxy</mat-label>
	<input matInput name="trafficRouterProxy" type="url"/>
</mat-form-field>
<mat-form-field>
	<mat-label>Traffic Ops URL</mat-label>
	<input matInput name="trafficOpsURL" type="url" required/>
</mat-form-field>
<mat-divider></mat-divider>
<footer>
	<button mat-raised-button color="warn" type="reset">Cancel</button>
	<button mat-raised-button color="primary">Save</button>
</footer>
</form>
</mat-card>
```
