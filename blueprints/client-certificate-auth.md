<!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
<Paste>"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->

# Client Certificate Authentication

## Problem Description

Passwords currently serve as a powerful, but inflexible, way for clients to
authenticate with Traffic Ops. However, an account can only have one password
at a time, which means that service accounts tend to have a single shared
password. This is undesirable from a security perspective.

Client certificates are a flexible tool that allows operators to assign
independent credentials to users of a service account. These credentials can
have varying expirations and be issued in accordance with the operator's
security policy.

## Proposed Change

When Traffic Portal or Traffic Ops receives a request, in addition to accepting
a valid token, it will accept a valid TLS certificate chain. Traffic Ops will
have a list of root certificates that it accepts.

Issuance and management of these certificates will occur outside of Traffic
Control and is the responsibility of the operator.

### Traffic Portal Impact

Traffic Portal will pull user information from requests before they are
proxied. Connections between TP and TO will be signed with a client certificate
to protect the connection.

When TP proxies a connection to TO, it validates the TLS certificate against
it's list of authorized roots. Iff it is valid, it includes the Subject of the
certificate verbatim in a request header named `Client-Cert-Subject`.
Otherwise, it strips any header named `Client-Cert-Subject` from the request.

Whether the certificate is valid or not, it removes any existing
`Client-Cert-Public-Key` header and adds a `Client-Cert-Public-Key` with the
base64-encoded public key of the client certificate.

TP will also add a Forwarded header.

TP will have a list of valid roots and a client certificate with private key in
it's configuration. The client certificate should have an empty value in the
UID Relative Distinguished Name to indicate that it is considered a trusted
proxy. TP will use this client certificate for all requests to Traffic Ops.

An invalid certificate should be processed as if it were not present, except
for the fact that its public key will be copied into the
`Client-Cert-Public-Key` header.

#### Certificate Expiration Notification

When an authorized root is within 30 days of expiration, Traffic Portal should
put a banner up for every administrative user on login indicating how long
remains until expiration.

### Traffic Ops Impact

Traffic Ops will need to accept certificates as proof of identity as well in
addition to tokens.

A certificate provides proof of identity for a given user iff:
  - the user id matches the value of the Relative Distinguished Name with the
    attribute UID in the Subject of the certificate and it chains through valid
    intermediates to an authorized root; or
  - the value of the UID RDN in the Subject of the certificate is empty, it
    chains through valid intermediates to an authorized root, and the user id
    matches the value of the single UID RDN of the `Client-Cert-Subject` header.

A certificate may not contain multiple UID RDNs. If it does, it is invalid.

In all log messages that contain a user's identity, the CN of the Subject of
the certificate and the IP address of connection should be logged as well. If
the user id is matched in the `Client-Cert-Subject` on account of an empty UID
in the Subject of the certificate, the CN of the Client-Cert-Subject and the IP
address from the Forwarded header should be logged as well.

An invalid certificate should be processed as though it were not present. So a
valid security token would still be accepted and a user that presents neither
would be redirected to the login page.

If Traffic Ops issues a token as a result of a login with valid credentials, it
should include the public key of the client certificate (either directly from
the certificate or from the `Client-Cert-Public-Key` header if presented with a
valid certificate with an empty UID) as a field in the encoded JSON. When a
token is validated, the public key in it should be compared to the public key
from either the certificate or the header as appropriate and and the token
rejected if they do not match.

#### REST API Impact

As part of Traffic Ops, the REST API will accept valid certificates as proof of
identity.

#### Client Impact

The official client libraries should expose an interface that allows a client
certificate to be provided so they can use the new functionality.

#### Data Model Impact

None. The relevant data is stored only in config files.

### t3c Impact

t3c needs to optionally accept a client certificate in place of existing
credentials.

### Traffic Monitor Impact

Traffic Monitor needs to optionally accept a client certificate in place of
existing credentials. This certificate needs to be able to be modified without
restarting Traffic Monitor.

### Traffic Router Impact

Traffic Router needs to optionally accept a client certificate in place of
existing credentials. This certificate needs to be able to be modified without
restarting Traffic Router.

### Traffic Stats Impact

Traffic Stats needs to optionally accept a client certificate in place of
existing credentials. This certificate needs to be able to be modified without
restarting Traffic Stats.

### Documentation Impact

The documentation should include a broad overview of how client certificates
can be used to authenticate. It should also detail the changes in config files
and tool parameters.

### Testing Impact

The unit tests need to generate certificates on-the-fly as part of test setup
so that tests do not periodically expire with baked in certificates. The
authentication functions should be part of the unit tests and the entire system
should be tested end-to-end with certificates as part of the integration tests.

### Performance Impact

There should be no serious performance impact of this beyond that being used
for TLS itself. The basic costs of TLS are very slightly more than unencrypted
transport, but this feature is entirely unavailable without TLS.

### Security Impact

The security impact of this feature will depend greatly on the operational
practices of the operator. Private root keys need to be properly maintained and
certificates should be rotated frequently. If a key or intermediate is
compromised, the root should be rotated and all keys reissued immediately.

Root certificates are global and do not apply to tenants in particular. Any
root cert can provide proof of identity for any user of any tenant.

Someone with administrative credentials could easily create a backdoor with a
long-lived cert that allows them to authenticate as any user. Operators will
need to keep tight control of their private keys. This is not worse than the
existing situation, though, because someone with administrative access to a
Traffic Ops system already has such access.

If Traffic Ops is serving unencrypted traffic, this feature will not
meaningfully improve security.

### Upgrade Impact

An upgraded system will initially contain an empty list of authorized roots.
This means that no client certificates will be valid and since invalid
certificates are ignored, the login behaviour will be identical.

#### Downgrade Impact

A system that is downgraded will lose the ability to understand client
certificates. Any clients that were relying on them will naturally lose their
ability to authenticate.

### Operations Impact

There is no operational impact for those who do not wish to use Client Certificates.

Using Client Certificates, however, requires an operator to have a robust and
secure Public Key Infrastructure. This is a significant endeavour that will
require operators to either manually manage a considerable number of keys and
certificates or to develop tooling that manages it automatically.

The PKI properly belongs entirely outside of Traffic Control, however, because
operators are likely to have their own requirements and existing tooling around
how this is handled.

If an operator uses a caching proxy to access Traffic Ops, it will need to be
issued an empty-UID client certificate in order to proxy authentication. It
will also need to validate presented certs against the authorized roots, copy
the Subject into the `Client-Cert-Subject` header, add the Forwarded header,
and include the UID in the cache key.

### Developer Impact

Developers with test systems may find it's easier to develop on test systems
using certificates instead of passwords. It will present an additional option
at least. Existing workflows will work just fine, though.

### Alternatives

The main alternative is the one already in use: passwords. We could accomplish
some of what we are attempting to do here by allowing users to have an
arbitrary number of valid passwords instead of one. Given the fact that
passwords are sensitive, it would be very, very difficult to rotate passwords
on an ongoing basis or in the event that one was compromised.

#### Certificate Revocation Lists

A CRL could potentially allow intermediates to be revoked without revoking the
root. The complexity is considerable and rotating roots should be reasonably
doable if necessary.

#### Tenancy

Adding tenancy features to the certificates themselves drastically increases
the complexity and requires intermediaries like Traffic Portal to implement the
tenancy business logic directly. This isn't worth it for the flexibility that
tenant-based intermediaries brings.

