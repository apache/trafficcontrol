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
# Delivery Services 'Active' field

## Problem Description
Setting a Delivery Service to "Inactive" actually only sets it to "not routed".
Currently, there is no way to create a Delivery Service (with assigned servers)
that will not be distributed to cache server configuration.

## Proposed Change
This blueprint proposes changing the `Active` property of Delivery Services
from a boolean to an enumerated string constant that can represent three
different "Activity States" for a Delivery Service.

## Data Model Impact
The proposed new type of the field is expressed below as a TypeScript
enumeration.
```typescript
/**
 * This defines what other components of ATC will consider a Delivery Service
 * "active".
 *
 * It's not an object exposed through the API in its own right, just a
 * specification of the allowed values.
*/
enum DeliveryServiceActiveState {
	/**
	 * A Delivery Service that is ”active” is one that is functionally
	 * in service, and fully capable of delivering content.
	 *
	 * This means that its configuration is deployed to Cache Servers and it is
	 * available for routing traffic.
	*/
	ACTIVE = 'ACTIVE',
	/**
	 * A Delivery Service that is ”inactive” is not available for
	 * routing and has not had its configuration distributed to its assigned
	 * Cache Servers.
	*/
	INACTIVE = 'INACTIVE',
	/**
	 * A Delivery Service that is ”primed” has had its configuration
	 * distributed to the various servers required to serve its content.
	 * However, the content itself is still inaccessible via routing.
	*/
	PRIMED = 'PRIMED'
}
```

We don't have a real data model for
<abbr title="Apache Traffic Control">ATC</abbr>, so what is proposed is that
everywhere a "Delivery Service" object is represented, the current "active"
field be re-typed to the above-described type. If there are representations
where this property is not expressed, it ought to be added for consistency's
sake.

## Impacted Components
While a relatively small change, this has some far-reaching repercussions. The
components affected are Traffic Portal (in time), Traffic Router, Traffic
Monitor, Traffic Ops, and <abbr title="Operational Readiness Test">ORT</abbr>.

### Traffic Portal Impact
None as yet, but when TP is made to use API version 3, it will need to account
for the new type of this field everywhere it parses Delivery Service objects
from Traffic Ops API responses, and/or sends requests.

Furthermore, once that is done, Traffic Portal's Delivery Service-based forms
will need to change from exposing the old boolean concept of "active" versus
"not active" to expressing the new enumerated values. Likewise all affected
tables will need to update to reflect the new values.

### Traffic Ops Impact

#### Database Impact
A new type must be declared as an enumeration of the above-described values.
The type of the `deliveryservices` table's `active` column will need to change
from `boolean` to the new enumerated type. The existing types must be coalesced
such that values which are currently `TRUE` become `'ACTIVE'`, and values which
are currently `FALSE` become `'PRIMED'`, thus preserving the existing behavior.

#### API Impact
The affected endpoints will be:

##### `/cdns/{{CDN Name}}/configs/monitoring`
Currently the Delivery Service information returned by this endpoint gives back
Delivery Services filtered such that - among other criteria - only those
that have `Active` values of `true` are returned. This must instead change to
all those that are **NOT** `INACTIVE` and an `active` field must be added to
the representations containing their "Activity State" value.

##### `/cdns/{{CDN Name}}/snapshot/new`
The Delivery Service representations returned by this endpoint are filtered
such that - among other criteria - only those that have `Active` values of
`true` are returned. This must instead change to all those that have `Active`
values of "ACTIVE".

No structural changes to the response payloads are strictly necessary, though
this author would strongly encourage whoever implements the changes add the
`active` field to the Delivery Service representations returned by this
endpoint to help move us toward a single, unified representations of objects,
and because it's a relatively small change to sneak in.

##### `/deliveryservices`
The `active` property of objects returned by this endpoints methods is
currently a boolean, and it must change to reflect the new enumerated string
constant type described above.

Furthermore, input objects must be parsed in accordance with the new type of
the `active` property. If an input is given with any value that is not one of
the enumerated string constants described above, the endpoint MUST respond with
a 400 Bad Request error response, accompanied by an `error`-level Alert that
explains why the request was rejected.

Finally, these changes must exist only in API v3. API v2 and lower must
coalesce "ACTIVE" to `true` and either other value to `false` in output and
perform the reverse operation on inputs to POST and PUT methods.

##### `/deliveryservices/{{ID}}`
Requires the same changes as `/deliveryservices`.

##### `/deliveryservices/{{ID}}/safe`
Requires the same changes as `/deliveryservices`, except that no changes are
necessary to the processing of its inputs.

##### `/servers/{{ID}}/deliveryservices`
Requires the same changes as `/deliveryservices`.

##### `/snapshot`
This endpoint currently creates Snapshots by selecting Delivery Services that
are filtered such that - among other criteria - only those that have `Active`
values of `true` are put into the newly created Snapshot. This must instead
change to all those that have `Active` values of "Active".

No structural changes to the response payloads are strictly necessary, though
this author would suggest that whoever implements the changes take this
opportunity to change the payload structure to conform to the API Guidelines
(success messages belong in `success`-level Alerts, not as response objects).

##### `/users/{{ID}}/deliveryservices`
This legacy endpoint must be updated to handle the new data as stored in the
database by coalescing "ACTIVE" to `true` and either other value to `false` in
its response payload objects. Unless somehow by the time this comes around
we've managed to get rid of API 1.x, which is the only version implementing
this endpoint.

#### Client Impact
Clients at or below API v2 should not need to change. The latest client version
will change what structures it returns, but need not change functionally.

And the Python client, of course, is unstructured and so is unaffected.

### ORT Impact
`atstccfg` will need to change handling of logic for generating Delivery
Service-based configuration to only consider Delivery Services that are
"ACTIVE" or "PRIMED" using the latest (API v3-based) client instead of the old
boolean-valued `active` flag.

### Traffic Monitor Impact
Traffic Monitor will need to update its polling logic to be able to parse
payloads from the `/cdns/{{CDN Name}}/configs/monitoring` endpoint that have
been modified as stated above. This should be done in a backward-compatible
manner, such that it can parse either that format or the old, APIv2 format.

## Documentation Impact
For every endpoint that has been modified, accompanying documentation changes
must be made. Further, the description of the `Active` property of a Delivery
Service in the Delivery Services overview section must change to reflect the
new type of that property.

## Testing Impact
All Traffic Ops API changes should be accompanied by associated client/API
integration tests, and Traffic Monitor's ability to parse the new payload
format in a backward-compatible manner should be accompanied by a corresponding
test.

When Traffic Portal is updated to make use of the newest API version its
associated UI changes should be accompanied by associated one or more tests.

## Performance Impact
This should positively impact performance in some cases where a Cache Server is
assigned to/within a Topology used by a Delivery Service that is "INACTIVE",
which was previously a concept that was incapable of being expressed by Traffic
Control. Thus, some configuration generation work may be avoided.

## Security Impact
> _"Are there any security risks to be aware of?"_

No, I don't believe so.

> _"What privilege level is required for these changes?"_

To make the changes? "Committer," I guess. If this question is about Traffic
Ops API Roles/Capabilities/Permissions, I suppose a better answer would be
"unchanged".

> _"Do these changes increase the attack surface (e.g. new untrusted input)?"_

No, not "new", just different. As far as "attack surface," though, that depends
and is further discussed two questions down.

> _"How will untrusted input be validated?"_

It's an enumerated type. So the input must be one of the enumerated string
constant values.

> _"If these changes are used maliciously or improperly, what could go wrong?"_

Typically I'm told "we trust operators" and since you need an operator role to
modify Delivery Services that would translate to an answer of "nothing that
matters". But to be more specific, you can now remove the ability of Cache
Servers to serve content for a Delivery Service, which can negatively impact
clients who are in the process of fetching object series or fragments from a
Cache Server. So use with caution.

> _"Will these changes adhere to multi-tenancy?"_

Tenancy rules unchanged, so yes.

> _"Will data be protected in transit (e.g. via HTTPS or TLS)?"_

The Traffic Ops API (at least the Go implementation) only serves content over
HTTPS as far as I'm aware, so yes.

> _"Will these changes require sensitive data that should be encrypted at rest?"_

No.

> _"Will these changes require handling of any secrets?"_

This feels like the same question as directly above it, but no.

> _"Will new SQL queries properly use parameter binding?"_

That this question needs to be asked is simultaneously baffling to me because I
don't know of any SQL library for any language that would tell you to do
anything else and horrifying to me because of what it implies about the past.

## Upgrade Impact
Traffic Monitor changes are backward-compatible, so no impact there. Traffic
Ops changes follow the API "version-promise" and so upgrades have no more or
less impact than any other breaking API change - if you keep using v2 you won't
even notice.

## Operations Impact
Operators that use raw API requests will need to be made aware of the changes
and their semantic meaning(s) - which should be covered handily by the
documentation. Operators using Traffic Portal will need the same information
when it is upgraded to use API v3 and provides the new UI for Delivery
Service's "Active" property.

### Developer Impact
Developers should be largely unaffected except for merge conflicts. Though it's
difficult to guess the full scope of what can and will be worked on by anyone.

## Alternatives
The one big alternative that was considered was making `active=false` mean "not
routed and not present in Cache Server configuration". However, this meant
losing the "in-between" state that "PRIMED" provides. So when a Delivery
Service is taken out of service, it would mean that the very next "queue
updates" that was done would remove it from Cache Server configuration even if
a Snapshot has not been taken. This results in clients being routed to Cache
Servers that cannot serve the content they requested.

By contrast, when setting a Delivery Service to "PRIMED", clients will no
longer be routed for that Delivery Service but even if updates are queued
before then clients using Cache Servers for ongoing content delivery will be
unaffected by the change. This allows operators to "drain" a Delivery Service
for a while before setting it to "INACTIVE" and removing it from Cache Server
configuration, thus minimizing client impact.

Another that has been brought up is changing the name of the current `active`
field to be `active-routed`, which would mean that Delivery Service will be
routed, and adding another field called `active-cached` which would mean that
the Delivery Service's configuration would be disseminated to the cache
servers. Thus the same amount of flexibility could be achieved without making a
breaking API change.

The trade-off is that it would allow Delivery Services to be in a state where
Traffic Router will direct clients to cache servers which are incapable of
properly serving them content.
