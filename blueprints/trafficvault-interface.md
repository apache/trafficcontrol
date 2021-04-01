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

# Traffic Vault Interface for Secret Data Store

## Problem Description

Currently, Riak is the only supported data store for *secrets* and their
related data (e.g. public and private keys, certificates, URI signing keys,
etc) in Traffic Ops, and Riak is largely unmaintained now and has become a
security risk. Traffic Ops needs more options for secret storage backends,
such as PostgreSQL, which is already used as the Traffic Ops database today.
However, it may be desirable to use other backends for Traffic Vault, such as
HashiCorp Vault, so the implementation should provide a clear interface for
implementations of different backends, with the ability to configure which
backend Traffic Ops should use.

Having this interface-based implementation will better separate the concerns of
reading/writing data to Riak or some other secure data store and everything
outside of that scope. Adding support for a new secure data store would then be
mainly just implementing the read/write integration between TO and the new data
store, without having to worry about everything else in between the TO API and
reading/writing to the data store, such as authentication/authorization,
request/response (de)serialization, input validation, HTTP handling, etc.

## Proposed Change

Within Traffic Ops, a new `TrafficVault` interface will be defined with the
minimum set of methods to fulfill the existing capabilities of Traffic Vault
(storing/reading sslkeys, DNSSEC keys, etc.). The existing Riak-based
implementation will be refactored into a concrete implementation of the
`TrafficVault` interface so that no external functionality will be lost when
using the Riak implementation.

Once the existing Riak implementation has been refactored to implement the
`TrafficVault` interface, a *new* PostgreSQL-based implementation of the
`TrafficVault` interface will be added. This implementation will make
PostgreSQL the primary replacement for Riak (as discussed on the mailing list),
but the community is still open to the idea of other implementations as well.

New configuration options will be added to `cdn.conf` which will allow the
choice of Traffic Vault backend as well as the location of that backend's
configuration file. Backend implementations may each have their own unique
configuration file which they are responsible for parsing. If unspecified,
Traffic Vault will default to the existing Riak backend and config.

In order to *swap* `TrafficVault` backends, e.g. from Riak to PostgreSQL, an
operator would have to stand up PostgreSQL (or add a new database to the
existing PostgreSQL instance used for Traffic Ops), copy all the data from Riak
to PostgreSQL, stop Traffic Ops, reconfigure the `TrafficVault` backend from
Riak to PostgreSQL, add the necessary PostgreSQL-specific configuration, start
Traffic Ops, and verify that all the "secure" TO API endpoints return and write
the same data as when the Riak backend is configured to use.

### Traffic Portal Impact

n/a

### Traffic Ops Impact

#### REST API Impact

There are several endpoints that the TO API currently uses Riak for:
- GET    /api/$version/cdns/name/:name/sslkeys
- GET    /api/$version/deliveryservices/xmlId/#xmlid/sslkeys
- GET    /api/$version/deliveryservices/hostname/#hostname/sslkeys
- POST   /api/$version/deliveryservices/sslkeys/generate
- POST   /api/$version/deliveryservices/sslkeys/add
- GET    /api/$version/deliveryservices/xmlId/:xmlid/sslkeys/delete
- POST   /api/$version/deliveryservices/xmlId/:xmlId/urlkeys/generate
- POST   /api/$version/deliveryservices/xmlId/:xmlId/urlkeys/copyFromXmlId/:copyFromXmlId
- GET    /api/$version/deliveryservices/xmlId/:xmlId/urlkeys
- GET    /api/$version/deliveryservices/:id/urlkeys
- GET    /api/$version/cdns/name/:name/dnsseckeys
- POST   /api/$version/cdns/dnsseckeys/generate
- GET    /api/$version/cdns/dnsseckeys/refresh
- GET    /api/$version/cdns/name/:name/dnsseckeys/delete
- POST   /api/$version/cdns/:name/dnsseckeys/ksk/generate
- GET    /api/$version/deliveryservices/:xmlID/urisignkeys
- POST   /api/$version/deliveryservices/:xmlID/urisignkeys
- PUT    /api/$version/deliveryservices/:xmlID/urisignkeys
- DELETE /api/$version/deliveryservices/:xmlID/urisignkeys
- GET    /api/$version/vault/bucket/:bucket/key/:key/values
- GET    /api/$version/vault/ping
- PUT    /api/$version/snapshot (certs from deleted DSes are deleted from Riak)
- PUT    /api/$version/deliveryservice/:id (creating DNSSEC keys)

The above endpoints will be updated to use the `TrafficVault` interface, of
which the concrete implementation will handle reading from and writing to the
secure data store backend that has been configured for Traffic Ops. The format
of the API requests and responses should remain exactly the same as they are
currently, so from a TO API client perspective there should be no change. The
only difference will be which backend TO is using as its secure data store.

#### Client Impact

Currently, not all of the aforementioned API endpoints have support in the
Traffic Ops Go client. As part of implementing the tests for this blueprint,
each of the aforementioned API endpoints will have corresponding methods added
to the Traffic Ops Go client (if the methods do not exist already).

#### Data Model / Database Impact

The TO internal data model shouldn't require much change, but new structs may
be required for marshalling/unmarshalling between TO and any new `TrafficVault`
backends. The lib/go-tc structs should be unaffected by this change.

This change should not require any new changes to the database schema, but it
might be required to seed new server/profile types for adding new servers into
TO that correspond to a `TrafficVault` backend (if applicable), similar to the
server type `RIAK` and profile type `RIAK_PROFILE` for the existing Riak
backend.

### ORT Impact

n/a

### Traffic Monitor Impact

n/a

### Traffic Router Impact

n/a

### Traffic Stats Impact

n/a

### Traffic Vault Impact

The existing Traffic Vault Riak backend will remain the default `TrafficVault`
backend if no backend has been chosen in the TO configuration, and the existing
format of Riak requests and responses should not change due to this blueprint.
PostgreSQL will be added alongside Riak as a supported `TrafficVault` backend
implementation.

No *new* data requirements will be added to the system as part of this change,
and existing data in Riak will be unchanged.

### Documentation Impact

Any new Traffic Vault backends should be documented along with their required
configuration. Existing Riak-related documentation should be reorganized into
backend-specific places -- i.e. Traffic Vault should have an overview, with a
listing of supported backends (including Riak). Each supported backend would
then have its own page of documentation describing setup, configuration, etc.
of that specific backend.

### Testing Impact

As part of implementing this blueprint, the optional "system" tests in the TO
API testing framework will be augmented to test the endpoints that depend on
the `TrafficVault` backend. To help facilitate running these tests in
cdn-in-a-box, support for the PostgreSQL `TrafficVault` backend will be added
to cdn-in-a-box, and PostgreSQL will become the new default backend for
`TrafficVault` within cdn-in-a-box. However, this can remain configurable so
that developers could still choose to run Riak as the `TrafficVault` backend in
cdn-in-a-box.

### Performance Impact

n/a

### Security Impact

The new refactored Riak implementation of the `TrafficVault` interface should
not introduce any changes in the traffic between TO and Riak -- traffic that is
encrypted today will remain encrypted after this blueprint is implemented.

The new PostgreSQL implementation will take advantage of the
[pgcrypto](https://www.postgresql.org/docs/current/pgcrypto.html) module in
order to perform encryption at rest, and data will be encrypted in transit as
long as ssl is enabled in postgres and used by the client (per the
administrator's discretion).

### Upgrade Impact

Traffic Ops should be able to be upgraded without any new required
configuration changes. By default (if no specific `TrafficVault` backend has
been specified in the configuration), Traffic Ops will assume that Riak is the
`TrafficVault` backend. If a specific `TrafficVault` backend has not been
specified and Riak is not enabled, the API endpoints that require a
`TrafficVault` backend will return an error similar to how they return an error
today when Riak is not enabled.

Once Traffic Ops has been upgraded to a version that supports new
`TrafficVault` backends, a Traffic Ops administrator would be free to setup
PostgreSQL and enable it as the `TrafficVault` backend. Until such time,
pre-existing Riak-based installations should continue to work without further
changes.

### Operations Impact

By default, operators should be able to ignore the new `TrafficVault`
configuration and continue using their existing Riak configuration if they
choose to do so (at their own risk). However, it will be strongly recommended
to migrate to the PostgreSQL backend once ready. Operators may be able to lean
on their existing PostgreSQL automation/support in order to setup a new
database for Traffic Vault, and ATC will provide tools to help administer the
PostgreSQL Traffic Vault database (the existing `db/admin` tool may be
augmented to handle this). If new fields are added to the APIs that require
Traffic Vault, operators may need to run provided migrations for the Traffic
Vault database, similar to running migrations for the Traffic Ops database.

### Developer Impact

Currently, Riak is the only supported Traffic Vault backend, and configuring
Riak and running it locally for a development environment is non trivial. By
contrast, it is fairly easy to stand up a PostgreSQL databse locally for
development/testing purposes. By adding support for PostgreSQL as a Traffic
Vault backend, it will be easier to develop/test endpoints that require Traffic
Vault, especially considering that developers already have a PostgreSQL
database for developing/testing Traffic Ops.

Developers should be made aware of the new `TrafficVault` interface and the
fact that there should be a clean separation between the TO API implementation
of the business logic and the integration with a particular `TrafficVault`
backend.  Simply put, the business logic should only depend on the
`TrafficVault` interface, never on a concrete implementation of that interface
(such as Riak). This interface and its associated methods should be well
documented in the code itself, and a `README` about the `TrafficVault`
interface might be prudent as well, so that developers will have an easier time
implementing support for new `TrafficVault` backends.

## Alternatives

Some alternatives might include:
- rather than abstract the `TrafficVault` implementations behind an interface
  and configuring which one to use, just swap in support for a different
  backend
  - Pros:
    - by being more direct, might be easier to implement
  - Cons:
    - an upgrade to the latest version of TO would require swapping in the new
      data store immediately
    - stuck with one choice of backend
    - more difficult to swap the implementation with a new data store
      integration
- use the Traffic Ops API plugin system
  - Pros:
    - able to override a subset of the routes that require Traffic Vault, in
      case some features (e.g. DNSSEC) are unused
    - cleaner separation for proprietary code/backends
  - Cons:
    - rather than focusing on just the data store integration, overriding
      plugins are also responsible for everything else, including business
      logic, HTTP handling, validation, etc.

## Dependencies

As Traffic Ops already requires a PostgreSQL database, introducing PostgreSQL
support as a Traffic Vault backend does not introduce a new dependency to the
project. Additionally, it allows operators to *remove* Riak as a dependency.

## References

n/a
