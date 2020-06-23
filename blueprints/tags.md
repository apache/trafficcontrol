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
# Tags

## Problem Description
Tags are a fairly long-requested feature; to be able to associate a key word or
phrase with a Delivery Service, or a user. Today, some organizations are abusing
the overloaded "description" properties of a Delivery Service to provide
comma-separated key words and phrases, essentially inventing their own tags by
forcing semantics on a field that does not have them.

## Proposed Change
A new, incredibly simple object will be added to the data model: Tags. A Tag is
nothing more than a unique name, which can be associated with the following
objects:

- Cache Groups
- Delivery Services
- Origins as they appear in Traffic Ops API requests and responses
- Physical Locations as they appear in Traffic Ops API requests and responses
- Profiles
- Servers as they appear in Traffic Ops API requests and responses
- Tenants as they appear in Traffic Ops API requests and responses
- Topologies as they appear in Traffic Ops API requests and responses
- Users as they appear in Traffic Ops API requests and responses

## Component Impact
The only impacted components will be Traffic Ops and Traffic Portal, as other
components are strictly not allowed to derive meaning from an object's Tags.

### Traffic Ops Impact
A new object will be added, the Tag, which will appear through the API like the
following TypeScript interface.

```typescript
interface Tag {
	name: string;
}
```

Further, the aforementioned object types - which may be referred to as "taggable
objects" - will all be augmented with a new property: Tags, which will be a set
of Tag Names.

```typescript
interface Taggable {
	tags: Set<string>;
}
```

The Tags will be "shared" globally across all object types; that is, there are
not separate "server tags" and "Delivery Service tags", only Tags.

#### Database Specifics
The suggested method for storing Tag data is with one table for storing what
Tags exist, and separate tables for each Taggable object to be paired with zero
or more Tags.

The main Tag table ought to look like this:
```text
               Table "public.tag"
 Column | Type | Collation | Nullable | Default
--------+------+-----------+----------+---------
 name   | text |           | not null |
Indexes:
    "tag_pkey" PRIMARY KEY, btree (name)
```

and as an example of an object-to-Tag relation table, the server-to-Tag relation
table should resemble this:
```text
            Table "public.server_tag"
 Column |  Type  | Collation | Nullable | Default
--------+--------+-----------+----------+---------
 tag    | text   |           | not null |
 server | bigint |           | not null |
Indexes:
    "server_tag_pkey" PRIMARY KEY, btree (tag, server)
Foreign-key constraints:
    "server_tag_server_fkey" FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE
    "server_tag_tag_fkey" FOREIGN KEY (tag) REFERENCES tag(name) ON UPDATE CASCADE ON DELETE CASCADE
```

#### API Impact
The actual amount of work to add Tags to all Taggable objects is rather large,
considering the number of endpoints that need to be changed, but each change is
quite non-invasive, and can all be done atomically without breaking any
functionality between them.

The new Tags endpoint (`/tags`) shall support the HTTP request methods GET and
POST. A further endpoint, `/tags/\{\{Tag Name\}\}` will also be added that
supports GET, PUT and DELETE.

##### `/tags`

The GET method will retrieve representations of all Tag objects, and ought to
support the standard pagination methods.

The POST method will create new tags from a request payload. The payload may be
either a single representation of a Tag, or a set thereof. In any case, no Tag
may be created with a Name of a Tag that already exists. Also, the Name of a Tag
is subject to the restrictions that it not be empty and contain only
alphanumeric characters.

The PUT method will edit (change the name) of an existing Tag

#### Client Impact
The return data structures will change, but no code change to any Traffic Ops
client should be necessary for the addition of Tags to existing objects.
However, new methods/functions/procedures must be added to manipulate Tags
themselves.

### Traffic Portal Impact
A new UI section, complete with tables and forms for editing and creation, will
need to be added to accommodate the addition, editing and deletion of Tags.
Also, existing tables and forms for now-taggable objects will need to be updated
to accommodate the use of Tags. Each separate change can be done atomically, and
since the new collections need not be required, they can be done in non-invasive
fashion of arbitrary order (provided Tag manipulation itself is done first)
without disrupting any existing processes or data.

### Documentation Impact
Tags will need to be documented in the data model, both themselves and on the
modeled objects which now contain them. Also, API endpoints will need to updated
to reflect the new return structures, where applicable.

### Testing Impact
<!--
*How* will this impact testing?

What is the high-level test plan?
How should this be tested?
Can this be tested within the existing test frameworks?
How should the existing frameworks be enhanced in order to test this properly?
-->

### Performance Impact
<!--
*How* will this impact performance?

Are the changes expected to improve performance in any way?
Is there anything particularly CPU, network, or storage-intensive to be aware of?
What are the known bottlenecks to be aware of that may need to be addressed?
-->

### Security Impact
<!--
*How* will this impact overall security?

Are there any security risks to be aware of?
What privilege level is required for these changes?
Do these changes increase the attack surface (e.g. new untrusted input)?
How will untrusted input be validated?
If these changes are used maliciously or improperly, what could go wrong?
Will these changes adhere to multi-tenancy?
Will data be protected in transit (e.g. via HTTPS or TLS)?
Will these changes require sensitive data that should be encrypted at rest?
Will these changes require handling of any secrets?
Will new SQL queries properly use parameter binding?
-->

### Upgrade Impact
<!--
*How* will this impact the upgrade of an existing system?

Will a database migration be required?
Do the various components need to be upgraded in a specific order?
Will this affect the ability to rollback an upgrade?
Are there any special steps to be followed before an upgrade can be done?
Are there any special steps to be followed during the upgrade?
Are there any special steps to be followed after the upgrade is complete?
-->

### Operations Impact
<!--
*How* will this impact overall operation of the system?

Will the changes make it harder to operate the system?
Will the changes introduce new configuration that will need to be managed?
Can the changes be easily automated?
Do the changes have known limitations or risks that operators should be made aware of?
Will the changes introduce new steps to be followed for existing operations?
-->

### Developer Impact
<!--
*How* will this impact other developers?

Will it make it easier to set up a development environment?
Will it make the code easier to maintain?
What do other developers need to know about these changes?
Are the changes straightforward, or will new developer instructions be necessary?
-->

## Alternatives
<!--
What are some of the alternative solutions for this problem?
What are the pros and cons of each approach?
What design trade-offs were made and why?
-->

## Dependencies
<!--
Are there any significant new dependencies that will be required?
How were the dependencies assessed and chosen?
How will the new dependencies be managed?
Are the dependencies required at build-time, run-time, or both?
-->

## References
<!--
Include any references to external links here.
-->
