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

# API Guidelines

## Problem Description
The Traffic Ops API duplicates a lot of functionality, is at times
self-inconsistent, and can be difficult to work with.

## Proposed Changes
Many manifestations of the above problems could be avoided entirely by adhering
to a set of guidelines for how our API ought to behave. That would certainly
make it more consistent, and would make the intended behavior much easier to
discern. Herein lie a proposed set of API guidelines that have been discussed
by the Traffic Ops working group over the past few months.

These rules aren't meant to necessarily be ironclad, but there should be a very
convincing argument accompanying endpoints that do not follow them.

### Response Bodies
All valid API responses will be in the form of some serialized object. The main
data that represents the result of the client's request MUST appear in the
`response` property of that object. If a warning, error message, success
message, or informational message is to be issued by the server, then they MUST
appear in the `alerts` property of the response. Some endpoints may return
ancillary statistics such as the total number of objects when pagination occurs,
which should be placed in the `summary` property of the response.

#### Response
The `response` property of a serialized response object MUST only contain object
representations as requested by the client. In particular, it MUST NOT contain
admonitions, success messages, informative messages, or statistic summaries
beyond the scope requested by the client.

Equally unacceptable API responses are shown below.

```json
{
	"response": "Thing was successfully created."
}
```

```json
{
	"response": {"foo": "bar"},
	"someOtherField": {"someOtherObject"}
}
```

When requests are serviced by Traffic Ops that pass data asking that the
returned object list be filtered, the response property MUST be a filtered array
of those objects (assuming the request may be successfully serviced). This is
true even if filtering is being done according to a uniquely identifying
property - e.g. a numeric ID. The `response` field of an object returned in
response to a request to create, update, or delete one or more resources may be
either a single object representation or an array thereof according to the
number of objects created, updated, or deleted. However, if a handler is
*capable* of creating, updating, or deleting more than one object at a time, the
`response` field SHOULD consistently be represented as an array - even if its
length would only be 1.

The proper value of an empty collection is an empty collection. If a Foo can
have zero or more Bars, then the representation of a Foo with no Bars MUST be an
empty array/list/set, and in particular MUST NOT be either missing from the
representation or represented as the "Null" value of the representation format.
That is, if Foos have no other property than their Bars, then a Foo with no Bars
may be represented in JSON encoding as `{"bars":[]}`, but not as `{"bars":null}`
or `{}`. Similarly, an empty string field is properly represented as an empty
string - e.g. `{"bar":""}` not `{"bar":null}` or `{}` - and the "zero-value" of
numbers is zero itself - e.g. `{"bar":0}` not `{"bar":null}` or `{}`. Note that
"null" values *are allowed* when appropriate, but "null" values represent the
*absence* of a value rather than the "zero-value" of a property. If a property
is *missing* from an object representation it indicates the absence of that
property, and because of that there must be a *very convincing* argument if and
when that is the case.

As a special case, endpoints that report statistics including minimums,
maximums and arithmetic means of data sets MUST use the property names `min`,
`max`, and `mean`, respectively, to express those concepts. These SHOULD be
properties of `response` directly whenever that makes sense.

#### Alerts
Alerts should be presented as an array containing objects which each conform to
the object definition laid out by
[the ATC library's Alert structure](https://pkg.go.dev/github.com/apache/trafficcontrol/lib/go-tc#Alert).
The allowable `level`s of an Alert are:

- `error` - This level MUST be used to indicate conditions that caused a request
to fail. Because of this, this level MUST NOT appear in the `alerts` array of
responses with any HTTP response code less than 400. Details of server workings
and/or failing components MUST NOT be exposed in this message, which should
otherwise be as descriptive as possible.
- `info` - This level SHOULD be used to convey supplementary information to a
user that is not directly the result of their request. This SHOULD NOT carry
information indicating whether or not the request succeeded and why/why not, as
that is best left to the `error` and `success` levels.
- `success` - This level MUST be used to convey success messages to the client.
In general, it is expected that the message will be directly displayed to th
user by the client, and therefore can be used to add human-friendly details
about a request beyond the response payload. This level MUST NOT appear in the
`alerts` array of responses with an HTTP response code that is not between 200
and 399 (inclusive).
- `warning` - This level is used to warn clients of potentially dangerous
conditions when said conditions have not caused a request to fail. The best
example of this is deprecation warnings, which should appear on all API routes
that have been deprecated.

#### Summary
The `summary` object is used to provide summary statistics about object
collections. In general the use of `summary` is left to be defined by API
endpoints (subject to some restrictions). However, its use is not appropriate in
cases where the user is specifically requesting summary statistics, but should
rather be used to provide supporting information - pre-calculated - about a set
of objects or data that the client *has* requested.

Endpoints MUST use the following, reserved properties of `summary` for their
described purposes (when use of `summary` is appropriate) rather than defining
new `summary` or `response` properties to suit the same purpose:

- `count` - Count contains an unsigned integer that defines the total number of
results that could possibly be returned given the non-pagination query
parameters supplied by the client.

### HTTP Request Methods
[RFC 7231 Section 4](https://tools.ietf.org/html/rfc7231#section-4) defines the
semantics of HTTP/1.1 request methods. Authors should conform to that set of
standards whenever possible, but for convenience the methods recognized by
Traffic Ops and their meanings in that context are herein defined.

#### GET
HTTP GET requests are issued by clients who want some data in response. In the
context of Traffic Ops, this generally means a serialized representation of some
object. GET requests MUST NOT alter the state of the server. An example of an
API endpoint created in API version 1 that violates this restriction is
`cdns/name/{{name}}/dnsseckeys/delete`.

This is the standard method to be used by all read-only operations, and as such
handlers for this method should generally be accessible to users with the
"read-only" Role.

All endpoints dealing with the manipulation or fetching representations of
"Traffic Control Objects" MUST support this method.

#### POST
POST requests ask the server to process some provided data. Most commonly, in
Traffic Ops, this means creating an object based on the serialization of said
object contained in the request body, but it can also be used virtually whenever
no other method is appropriate. When an object *is* created, the response body
MUST contain a representation of the newly created object. POST requests do not
need to be *idempotent*, unlike PUT requests.

#### PUT
PUT is used to replace existing data with new data that is provided in the
request body.
[RFC 2616 Section 9.1.2](https://tools.ietf.org/html/rfc2616#section-9.1.2)
lists PUT as an "idempotent" request method, which means that subsequent
identical requests should ensure the same state is maintained on the server.
What this means is that a client that PUTs an object representation to Traffic
Ops expects that if they then GET a representation of that object, do the same
PUT again and GET another representation, the two retrieved representations
should be identical. Effectively, the `lastUpdated` field that is common to
objects in the Traffic Ops API violates this, but the other properties of
objects - which can actually be defined - generally obey this restriction. In
general, fulfilling this restriction means that handlers will need to require
the entirety of an object be defined in the request body.

When an object is replaced, the response body MUST contain a representation of
the object after replacement.
While RFC 2616 states that servers MAY create objects for the passed
representations if they do not already exist, Traffic Ops API endpoint authors
MUST instead use POST handlers for object creation.

All endpoints that support the PUT request method MUST also support the
If-Unmodified-Since HTTP header.

#### PATCH
At the time of this writing, no Traffic Ops API endpoints handle the PATCH
request method. PATCH requests that the server's stored data be mutated in some
way using data provided in the request body. Unlike PUT, PATCH is not
*idempotent*, which essentially means that it can be used to change only part of
a stored object. When an object is modified, the response body MUST contain a
representation of the object after modification, and that representation SHOULD
fully describe the modified object, even the parts that were not modified.

Handlers that implement PATCH in the Traffic Ops API MUST use conditional
requests to ensure that race conditions are not a problem, specifically they
MUST support using ETag and If-Match, and SHOULD also support
If-Not-Modified-Since.

Clients SHOULD use PATCH requests rather than PUT requests for modifying
existing resources whenever it is supported.

#### DELETE
DELETE destroys an object stored on the server. Typically the request will
contain identifying information for the object(s) to be destroyed either in the
request URI or in the request's body. Traffic Ops API endpoint authors MUST use
this request method whenever an object identified by the request URI is being
destroyed. When such deletion successfully occurs, the response body MUST
contain a representation of the destroyed object.

### HTTP Response Codes
Proper use of HTTP response codes can significantly improve user interfaces
built on top of the API. What follows is a (non-exhaustive) set of response
codes and their appropriate use in the context of Traffic Ops. For more complete
information, refer to
[the Mozilla Developer Network's HTTP Response Code list](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status).

#### 200 OK
This indicates the request succeeded, with no additional semantics. This MUST be
the exact response status code of successful GET requests. This is also the
default "success" response code for any other request.

#### 201 Created
This indicates that a resource was successfully created on the server. This MUST
be the response status code of POST requests that create a new object or objects
on the server, and in that case the response SHOULD also include a Location
header that provides a URI where a representation of the newly created object
may be requested.

#### 202 Accepted
`202 Accepted` MUST be used when the server is performing some task
asynchronously (e.g. refreshing DNSSEC keys) but the status of that task
cannot be ascertained at the current time. Ideally in this case, when the task
completes - either successfully or by failing - the Traffic Ops changelog will
be updated to indicate that status, along with information to uniquely identify
the task (e.g. username and date/time when the task started).

Endpoints that create asynchronous jobs SHOULD provide a URI to which the
client may send GET requests to obtain a representation of the job's current
state in the Location HTTP header. They MAY also provide an `info`-level
Alert that provides the same or similar information in a more human-friendly
manner.

The responses to such GET requests are subject to the same restrictions as any
other API endpoint, but have the added restriction that the `response` objects
sent MUST have the `status` property, which is a string limited to one of the
following values and having the associated semantics:

- PENDING - This means the job has been started but is not yet completed.
- SUCCEEDED - This means that the asynchronous job has completed and encountered
no errors.
- FAILED - The task encountered errors and was unable to continue, and thus has
been terminated.

Note that the response code of the response carrying this information MUST NOT
depend on the value of `status`. In particular, a response that successfully
reports the status of a FAILED asynchronous task is still successfully servicing
a client's GET request, and therefore MUST have the `200 OK` response status
code. However, a response encoding a FAILED `status` MUST be accompanied by one
or more `error`-level Alerts that explain (to the greatest degree of detail
allowable securely) why the job failed.

These responses MUST also include the `startTime` and `endTime` properties which
indicate, respectively, the time at which the asynchronous job started and the
time at which it concluded. A job that has not started MUST have a Null-valued
`startTime` and likewise a job that has yet to conclude MUST have a Null-value
`endTime`.

#### 400 Bad Request
In general this is used when there's something syntactically wrong with the
client's request. For example, Traffic Ops MUST respond with this code when the
request body was improperly encoded. In most cases, this is also the proper
response code when the client submits data that is not semantically correct. For
example, dates/times represented as timestamp strings in an unsupported format should trigger this response code.

This is also the default "client failure" response code for any other request.

The response body MUST include an entry in the `alerts` array that describes to
the client what was wrong with the request.

#### 401 Unauthorized
This MUST be the response code when a client without valid authorization
information in the HTTP headers requests a resource which cannot be accessed
without first authorizing. Which should be everything except `/ping` and
endpoints that provide authorization.

#### 403 Forbidden
This MUST be used whenever the client is logged-in, but still does not have
access to the resource they are requesting. It MUST also be used when they have
some access to the resource, but not with the specific request method they used.
This can pertain to restricted access on the basis of Role, User Permissions, as
well as Tenancy.

The response body MUST NOT disclose any information regarding why the user was
denied access.

#### 404 Not Found
This MUST be the returned status code when the client requests a path that does
not exist on the server. Note that a _path_ does not include a _query string_;
in the URL `http://example.test/some/path?query#frag` the _path_ consists of
only `/some/path`.

#### 409 Conflict
This SHOULD be used when the request cannot be completed because the current
state of the server is fundamentally incompatible with the request. For example,
creating a new user with an email that is already in use should result in this
response.

Additionally, this MAY be used instead of `404 Not Found` when the client is
requesting a link between an object identified by the request URI and some
other object (e.g. when assigning a server to a Delivery Service) when the other
object does not exist. If the request URI identifies an object that does not
exist, the response MUST use `404 Not Found` instead.

This is also the proper response status code when the conditions of a request
cannot be met, e.g. when a client submits a PATCH request for a resource with an
If-Match header that does not match the stored object's ETag.

The request body MUST indicate what the conflict is that prevented the request
from being fulfilled via one or more `error`-level alerts.

#### 500 Internal Server Error
When the Traffic Ops server encounters some error - through no fault of the
client or their request - that renders it incapable of servicing the client's
request, it MUST return this status code if no other code is more appropriate.
The response body in this case SHOULD indicate that an error occurred, but
MUST NOT divulge details about what data was being processed, what (if any)
other components are not functioning properly, or what process failed. Generally
it is advisable that the resultant \code{alerts} array entry just say "Internal
Server Error" and nothing else.

#### 501 Not Implemented
This is the response code used when the client requests an API version not
implemented by the server. It SHOULD NOT be used in any other case.

#### 502 Bad Gateway
This code indicates that some other service on which the endpoint's processes
depend has given back improper data or an error response. It MAY be used (with
caution) by plugin developers, but SHOULD NOT be used by authors of proper API
endpoints, as that divulges information about failing connected systems and
potentially gives an attacker information about Traffic Control's weak points.
API endpoint authors should instead use `500 Internal Server Error`.

#### 504 Gateway Timeout
This code indicates that a connection timeout occurred when attempting to
contact some other service on which the endpoint's processes depend. It MAY be
used (with caution) by plugin developers, but SHOULD NOT be used by authors of
proper API endpoints, as that divulges information about failing connected
systems and potentially gives an attacker information about Traffic Control's
weak points. API endpoint authors should instead use `500 Internal Server
Error`.

### Documentation
All endpoints must be properly documented. For guidelines for writing API
documentation, refer to
[that section of the official ATC documentation](https://traffic-control-cdn.readthedocs.io/en/latest/development/documentation_guidelines.html#documenting-api-routes).

### Passing Request Data
Request data may be passed in the request body or as a
`application/x-www-form-urlencoded`-encoded query string in the request URI, or
as a part of the request path. Request data MUST NOT be passed through a portion
of the request path unless it uniquely identifies a resource with which the
client may interact. For example, `/foos/{{ID}}}` is an acceptable path for
dealing with the particular "Foo" object that has some identifier `ID`, but
`logs/{{Number of Days}}/days` is unacceptable because reasonable default
behavior can be provided if no number of days is given in the query string
parameters, and that doesn't help uniquely identify a resource. Request path
parameters should use double "curly-braces" (<kbd>{</kbd> and <kbd>}</kbd>) to
call out variable components of the request path in documentation and
references. Request path parameters MUST NOT be used for data that is optional
to the request (somewhat obviously). Note that all endpoints dealing with the
manipulation of "Traffic Control Objects"  MUST support the GET HTTP request
method.

When accepting data in the request body of requests, the endpoint MUST properly
document the object representations (properties and their types) it accepts and
MUST reject semantically invalid data with a `400 Bad Request`.  For example,
if an endpoint specifies it accepts a representation of a Foo object, assuming
Foo objects possess only the Bar property which is an arbitrary string, then
the endpoint MUST accept `{"bar": "testquest"}` as semantically valid (The data
may be rejected for other reasons, e.g. if a Foo with such a Bar property
already exists and Bars must be unique among all Foos) and MUST reject
`{"bar": "testquest", "someOtherProperty": 10}` as semantically invalid. This
is in contrast to the API's behavior at the time of this writing, which
silently ignores unrecognized properties of request body objects.

The decision to pass data in the request body or query string is mainly up to
the author, but some helpful tips:

- GET and DELETE requests do not typically provide request bodies.
- Query parameters should nearly always be optional. If data is required by an
endpoint, consider requiring it in the request body. If the data identifies a
resource, it ought to be a path parameter.
- Request body data often represents objects that are being created or updated.
If an object is being created or updated, it ought to be defined in the request
body, and if any additional data is (possibly optionally) required then it ought
to be passed in the query string to separate it from the object definition.
- The following query parameters are reserved for special use by Traffic Ops
endpoint handlers, and may not be used for any purpose other than their
prescribed functions.

	- `limit`
	- `newerThan`
	- `offset`
	- `olderThan`
	- `orderby`
	- `page`
	- `sortOrder`

### Duplicate Endpoints
No two endpoints should serve the same purpose. While it's fine to overlap a
bit, an endpoint like `/foo_bars` should not exist solely to edit the Bars
property of Foo objects (which can ostensibly be edited just fine on the object
itself), for example. Ideally, there should be exactly one way to accomplish
something through the API.

A caveat, though, is object relationships. For example, a Delivery Service has
zero or more Cache Servers assigned to it, and in turn Cache Servers may be
assigned to zero or more Delivery Services (a "has-and-belongs-to-many"
relationship). Thus it is permissible to be able to edit the Delivery Services
property of a Cache Server using the `/cache_servers` API endpoint as well as to
be able to edit the Cache Servers property of a Delivery Service using the
`/delivery_services` API endpoint - though they arguably provide equivalent
functionality in that way (although at the time of this writing the former
endpoint doesn't exist and the latter doesn't offer that functionality - this is
just an example).

### Date/Time Format
Dates MUST be represented in either
[RFC3339](https://tools.ietf.org/html/rfc3339) (with or without nanosecond
precision) or as integers indicating the number of nanoseconds past the Unix
epoch at which the date/time occurs. In either case, Dates included in responses
from Traffic Ops MUST be in UTC. Wherever date/times are accepted as input,
Traffic Ops API endpoints MUST accept either format and SHOULD NOT accept
anything else.

Traffic Ops endpoints MUST return dates and times in
[RFC3339](https://tools.ietf.org/html/rfc3339) format with nanosecond precision.
Endpoints MAY provide ways for the client to specify alternate representations,
but these SHOULD be restricted to only Unix epoch timestamps in nanoseconds.

### Age Filtering
Whenever object age is a property of that object (which is quite often in the
form of `lastUpdated`), Traffic Ops endpoint handlers that respond to requests
for object representations (i.e. GET requests) SHOULD support filtering by age.
If age filtering is implemented, it MUST be made available using the query
parameters in the table below.

| Parameter   | Meaning                                                                                                                                                                                                                                                                    |
|-------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `newerThan` | A timestamp to be used as the lower limit on an object's age. Objects older than this MUST NOT appear in the response body. That is, the response will be the set of all objects in the collection with a modification date that is greater than *or equal to* this value. |
| `olderThan` | A timestamp to be used as the upper limit an object's age. Objects newer than this MUST NOT appear in the response body. That is, the response will be the set of all objects in the collection with a modification date that is less than *or equal to* this value.       |

The format of these timestamps - in accordance with the above Date/Time Format
of this blueprint - MUST be accepted as Unix epoch timestamps in nanoseconds,
AND in the form of RFC3339 date/time strings.

Endpoints MAY return errors when a client request gives these parameters
improper or invalid values, but MUST at least provide a warning. When ambiguity
or errors in age filtering controls render age filtering impossible, the handler
MUST NOT perform age filtering.

### Tenancy
When a client requests access to a set of stored objects that are "tenantable"
inevitably some of them will be inaccessible to the user on the basis of their
Tenant. Traffic Ops endpoint handlers that respond to requests for such object
representations (i.e. GET requests) MUST filter their results implicitly
according to the requesting Tenant's access. Any request that would modify,
create, or destroy an object to which the requesting Tenant does not have access
MUST NOT be fulfilled by the server (obviously) and in that case the response
status code MUST be `403 Forbidden`. Furthermore, if a request for a
representation of a Tenant-inaccessible object is made *explicitly for said*
*object* (e.g. `GET /foos/{{ID}}` rather than `GET /foos?id={{ID}}`) the
response status code MUST be `403 Forbidden`.

### Naming Conventions
The names of properties of objects as they appear in said objects'
serializations ought to conform to "camelCase" naming. Initialisms,
abbreviations, and acronyms that appear in property names should be capitalized
unless they are at the very beginning of the name. For example,
`myIPAddress` and `someProperty` are both well-formed property names, while
`IPAddress`, `someproperty` and `SomeProperty` are not.

Query string parameters MUST also follow "camelCase" naming.

API endpoints themselves should have a name that conveys their purpose. For
example, `/cdns` is an endpoint that deals with manipulating, creating,
destroying, or retrieving representations of CDNs. Request paths MUST use
"snake_case" to separate words whenever necessary, and MUST never include the
action being performed by the handler; instead that is decided by the request
*method*. For example, `/myObject/delete` is a poor request path name for both
of those reasons. Furthermore, when an endpoint deals with an object type of
which there are typically multiple, the request path should be plural, e.g.
`/cdns` is better than `/cdn`.

API endpoints MAY support trailing slashes (<kbd>/</kbd>) in the request path,
but MUST NOT include suffixes that indicate a particular encoding ("file
extensions"); that's what the Content-Type header is for. For example, in
API version 1.x, `/foos` and `/foos.json` are both equally valid ways
to access the `/foos` endpoint handlers - this is no longer allowed!

### Relationships as Objects
Relationships SHOULD NOT be represented through the API as objects in their own
right. For example, instead of an endpoint like `/delivery_service_servers` used
to manipulate assignments of Cache Servers to Delivery Services, a Delivery
Service itself has Servers as a property. Thus assignments are manipulated by
manipulating that property. So the only endpoints necessary for fully defining
and dealing with such relationships are `/delivery_services` and `/servers`.

### Change Logging
All manipulations of objects (i.e. any operation that is not merely "reading"
data) MUST add a Change Log entry indicating what was changed.

## Component Impact
The intent here is to establish guidelines *moving forward* - no breaking
changes should be made inside a major version and this is **not** a proposal to
create a new major version. As such, impact besides documentation is expected to
be small.

### Traffic Portal Impact
Hopefully none will be necessary, but it may be beneficial to utilize new API
features, which should be evaluated as they are implemented.

### Traffic Ops Impact
The Traffic Ops API doesn't *need* to change at all, but endpoints may be
updated to support Age Filtering, If-Unmodified-Since, and PATCH/ETag. Endpoints
that don't conform to standards could be deprecated as new handlers that *do*
conform are written. Changes such as proper HTTP response codes may be
implemented immediately on existing endpoints, as those are undocumented.

### Client Impact
Clients that currently consider `200 OK` to be the only successful response code
should update that assumption to consider anything on the interval [200, 300) to
be a successful code, as per RFC standards. But they don't have to until that's
actually implemented in Traffic Ops.

### Documentation Impact
The informational content blueprint should be converted to reST and added to the
official documentation.

## Developer Impact
Developers must now be aware of and abide by the guidelines herein outlined from
the point at which it is approved onward.
