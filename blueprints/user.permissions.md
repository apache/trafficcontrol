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
# Enforced User Permissions

## Problem Description
Currently, the Permissions (currently "Capabilities", henceforth referred to as
"Permissions" to avoid confusion with Server Capabilities) afforded to a user
are defined by the "Privilege Level" of their Role. It essentially defines a
scale on the interval [0,30], where a higher number allows more things. This
isn't as scalable or configurable as it could be, and involves granting more
permissions to a user than they actually need just because every Privilege
Level encapsulates all of the permissions of every Privilege Level below it,
making disjoint Permission sets for disparate Roles logically impossible.

## Proposed Change
Instead of a system of rigid supersets of Permissions, this blueprint proposes
that the system of user "Capabilities" be expanded upon, renamed to Permissions
(again, to avoid confusion) and enforced by the API.

Under this system, each endpoint will declare the Permissions it requires for
specific operations and check that the requesting user has those required, if
any.

This information is _not_ to be stored in the database nor exposed through the
API, but documented. This change is discussed in greater detail in the
"Alternatives" section.

The enforcement of these Permissions is proposed to be configurable via an
entry in the Traffic Ops configuration file. Thus, for one major version the
API will need to support Permissions-based authorization as well as Privilege
Level-based authorization.

## Data Model Impact
"Capabilities" as they are known today will be removed from the data model.
Also, the current model of a Role will need to be modified. The current model
is as follows:

```typescript
interface Role {
	capabilities: Array<string> | null;
	description: string | null;
	id?: number | null;
	name: string | null;
	privLevel: number | null;
}
```

The proposed new model is shown below.

```typescript
interface Role {
	description: string;
	lastUpdated?: Date;
	name: string;
	permissions: Array<string>;
}
```

Also, because the ID is no longer exposed via the API, the model for users will
need to change to only reference Role name, not ID. So it'll change from...

```typescript
interface User {
	addressLine1: string | null;
	addressLine2: string | null;
	city: string | null;
	company: string | null;
	confirmLocalPassword?: string | null;
	country: string | null;
	email: string | null;
	fullName: string | null;
	gid: number | null;
	id?: number | null;
	lastUpdated?: string | null;
	localPassword?: string | null;
	newUser: boolean | null;
	phoneNumber: string | null;
	postalCode: string | null;
	publicSSHKey: string | null;
	registrationSent: string | null;
	role: number | null;
	roleName?: string | null;
	stateOrProvince: string | null;
	tenant?: string | null;
	tenantID: string | null;
	uid: number | null;
	username: string | null;
}
```

... to ...

```typescript
interface User {
	addressLine1: string | null;
	addressLine2: string | null;
	city: string | null;
	company: string | null;
	confirmLocalPassword?: string | null;
	country: string | null;
	email: string | null;
	fullName: string;
	gid: number | null;
	id?: number;
	lastUpdated?: Date;
	localPassword?: string | null;
	newUser: boolean;
	phoneNumber: string | null;
	postalCode: string | null;
	publicSSHKey: string | null;
	registrationSent: string | null;
	role: string;
	stateOrProvince: string | null;
	tenant?: string;
	tenantID: number;
	uid: number | null;
	username: string;
}
```

... replacing the required `role` field that was a Role ID with one that is
the Role's name (and therefore removing `roleName` as it's now redundant).

(Note that `lastUpdated` changed from `string | null` to `Date`, which
represents the transition to using proper RFC3339 formatting, and also some
fields which cannot be `null` for a valid user are no longer `| null`. Those
aren't changes necessary to implement Permissions, but can be done easily at
the same time as necessary changes to User.)

## Component Impact
This change primarily concerns Traffic Ops and Traffic Portal, though any
client of the API will need to be aware of the changes to its authorization
system.

### Traffic Portal Impact
There is a "Capabilities" view currently in Traffic Portal that provides a
table of user Permissions, though it is not linked to the side navigation bar.
This should be removed.

The Roles' details pages will need to be augmented with controls for adding and
removing Permissions. The Permissions available for adding to a Role should
include an autocompletion of those known by Traffic Portal to exist, but
shouldn't restrict the addition of Permissions not known to exist to accomodate
running older versions of Traffic Portal with newer versions of Traffic Ops, as
well as any plugins unbeknownst to Traffic Portal that add Traffic Ops API
endpoints and may declare their own Permissions that do not exist in a vanilla
Traffic Control environment.

Also, the Roles table will be augmented with the ability to filter Roles by the
presence or absence of a Permission, although the full itemization of the
Permissions afforded to a Role would, in most cases, be far too large to
comfortably display in the table itself.

Under-the-hood, API services will need to change from using Role IDs to Role
names.

### Traffic Ops Impact
Traffic Ops will reqire changes to its Roles and Permissions API, chiefly, and
the database tables that back them. Enforcement of API Permissions, though,
will also have more wide-ranging impact that will require changes to very
nearly every API endpoint.

Additionally, a configuation file field will be added named `usePermissions`
that is an optional boolean which, when present and `true` causes Traffic Ops
to use Permissions rather than Role Privilege Level to determine a user's
authorization for a given operation.

#### API Impact
Structurally, the only necessary changes are to the `/roles` and `/users`
endpoints, which will need to be updated to output structures consistent with
the changes outlined in the Data Model Impact section. The `/capabilities` and
`/api_capabilities` endpoints will be removed.

The more pervasive changes will be to all authenticated API endpoints which
shall be updated to respect Permissions given the correct configuration
setting.

The exact Permissions that need to exist for each endpoint are best left to
debate on the changeset that implements them.

When a user attempts an operation for which they do not have sufficient
Permissions, the API MUST respond with a `403 Forbidden` response containing an
error-level Alert that describes what operation is not permitted and what
Permission the user is missing that would allow them to proceed.

Also, the `admin` Role that comes with new Traffic Ops installations currently
has two "Capabilities" - `all-read` and `all-write`; this intention should be
replicated, but not necessarily the implementation. The Role named `admin` is to
be treated specially in that it has *all* Permissions - it should, therefore, be
impossible to delete the Role or to update its Permissions or Name. This writer
suggests that for display purposes the Permissions used to communicate this
ability be reduced to simply `all`, and that rather than checking for this
Permission endpoints would only consider the Name of the Role. It's not terribly
important that users be disallowed to alter the `admin` Role's Description, but
it may be easier to implement the other restrictions in a way that does disallow
that.

#### Client Impact
As Permissions are defined by the Role of the authenticated user, the only
client changes are necessary beyond those necessitated by the removal of two
API endpoints and the renaming and restructuring of a third are to stop
filling in a Role ID on user creation/update from a Role name if an ID was
unset (this will now be handled by the API itself, which will save a full
HTTP request on each user creation/modification).

#### Database Impact
The new model for a Role does not allow a `null` description; a simple
migration that coalesces existing `NULL` values to an empty string and adds a
check constraint should be all that's actually required. No other immediate
changes should be made - including dropping the now-unused `api_capability`
table -, since old API versions will still need access to the deprecated
columns. However, the foreign key constraint on the `role_capability` table
that links a "Capability" name to a row in the `capability` table should be
dropped, as that is no longer the source of truth for valid Permissions.

Roles that exist already should be ensured via database migration to have
whatever Permissions are necessary to use the API endpoints that their
"Privilege Level" already allows them to use. For privilege levels that fall
between values needed for any API endpoints, the only safe thing to do is "round
down" and only give them the Permissions they would need to have the same level
of API access as the next-lowest Privilege Level that **is** explicitly required
by some method and/or endpoint.

If the `admin` Role does not exist, the database migrations for this change
should create it.

## Documentation Impact
The new configuration option will need to be documented, documentation for
removed API endpoints will itself need to be removed, documentation for the
`/roles` and `/users` endpoints will need to be updated to reflect the new
request and response structures, and as Permissions are implemented on each
endpoint the Permissions it requires for various actions will need to be
defined.

The special `admin` Role must also be documented somewhere (most appropriately
in an overview section on Roles and Permissions).

## Testing Impact
The most significant testing changes will need to be made to the Go Traffic Ops
client integration tests, which should verify each endpoint's proper
Permissions-based authorization. Traffic Portal functionality will also require
the appropriate end-to-end tests.

## Performance Impact
No significant performance impact is expected, since the Permissions of a user
are already queried at the time of servicing a request by every authenticated
endpoint. Some negligible, constant time will need to be spend determining how
to authorize an authenticated user.

## Security Impact
Careful consideration must be given to the design and implementation of each
Permission. For example, this author believes a Permission named
`do-secrets-things` that allows a Role unrestricted read and write access to
any and all DNSSEC, SSL, URL Signature, and URI signing keys would be a poor
design from a security standpoint. Permissions should be broad enough to
encompass a single, well-defined purpose, and no more. In many cases, though,
the existing "Capabilities" concepts will be good enough to build from (e.g.
`delivery-services-write`, `cdns-read`).

## Upgrade Impact
There will be a database migration to run, but since the default configuration
will be to ignore Permissions and just use Privilege Level, there isn't much in
the way of upgrade impact immediately. The bigger step is getting ready for the
_next_ upgrade, when Permissions stop being optional.

## Operations Impact
Before enabling Permissions, operators will need to ensure that all Roles have
the appropriate Permissions to accomplish their necessary tasks. They will also
need to be aware of the Permissions required to accomplish any given task, or
at least where in the documentation to look for guidance in that regard.

### Developer Impact
From now on, new endpoints will need to be designed with their Permissions in
mind. For many endpoints this will probably be trivial, but must always be at
least considered.

## Alternatives
This section lists possible alternatives to the proposed system that were
considered.

### The current "Capabilities" system
The current system of "Capabilities", if enforced, was a possible alternative
to the system herein described. However, that system lacks the ability to
express a user's permission to do something beyond a combination of HTTP
request method and path. For example, the CDN Locks blueprint (#5834) proposes
a system by which a user may delete a lock that they created, but also states
that "`admin` users" (Privilege Level 30 or above) can delete any lock created
by any user. However, these two permissions use the same request method and
path, meaning that under the current system it would be impossible to give
admin users the ability to do that overriding deletion without also giving it
to everyone else. The system described here is far more flexible. It could even
eliminate the need for the `/deliveryservices/safe` API endpoint, which can be
expressed instead as two separate Permissions: one that allows making changes
to the "safe" Delivery Service fields and one that allows all others.

### Configuration file-based Role definitions
Another alternative would have been to hard-code user Permissions and possibly
also available Roles into a configuration file for Traffic Ops, instead of
exposing the Role creation and Permissions assignment via the Traffic Ops API
and storing that in the database. That approach does pose a problem for
automatic migration on upgrade, since it's beyond the scope of a database
migration. The new mappings would need to be constructed from the old Roles in
some way, either by hand or with some tool created and provided by ATC. The
reverse situation, a downgrade, is more or less impossible to handle without
manual intervention.

From a security standpoint, it could be argued that an on-disk file defining
Roles and Permissions provides fewer attack vectors through which a privilege
escalation could be performed.

From a performance standpoint, it would also require less data to be queried
from the DB with every single request, because the Roles-to-Permissions mappings
would be stored in memory. Given how many API endpoints currently exist, that
could save quite a lot of data traversal, especially for administrative-type
users which will have most (if not all) of the available Permissions. However,
this runs the risk of database state drifting out-of-sync with the configured
Roles; for example if a user has the "Foo" Role and then that Role is deleted
from the configuration file which is then reloaded, how does the API handle this
now non-existent Role?

From a developer standpoint, that is less code to implement and maintain,
allowing us to finish the project faster and require less maintenance. Although,
some fairly substantial API changes are still entailed, and changes to the Roles
and Permissions-related endpoints would still need to be made so that e.g.
Traffic Portal can decide what functionality to offer a user. This is also only
less work if a tool is not made to handle migration.

Additionally, testing this implementation would prove difficult in an end-to-end
manner, as the current end-to-end testing using the various Go Traffic Ops
clients does not have access to the process memory of Traffic Ops itself, so
either end-to-end testing of Roles and Permissions as an access control measure
would need to be abandoned, or a new testing method and/or framework would need
to be adopted.

From an operator standpoint, if a user is requesting more Permission, using the
API would be more convenient than say, making a change to an Ansible file
template. However, users needing more Permission should be relatively rare, and
editing an Ansible file template is pretty straightforward anyways. However,
users needing _less_ Permissions is much more time-sensitive than needing
_more_, and operators would need to perform configuration deployment - whatever
that entails for their organization (possibly including restarting Traffic Ops
although that's not necessarily required for loading changes to the specific
Roles and Permissions configuration file) rather than using the user-friendly
interface of Traffic Portal.
