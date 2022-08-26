## What does this PR (Pull Request) do?

This PR makes the Go components of Traffic Control build using Go version {GO_VERSION} and updates the `golang.org/x/` dependencies.

See the Go {GO_VERSION} [release notes](https://golang.org/doc/devel/release.html#go{GO_MAJOR_VERSION}):

<!--
The release notes are licensed with the Go license.

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
-->
> {RELEASE_NOTES}

## Which Traffic Control components are affected by this PR?
- Traffic Control Cache Config (`t3c`, formerly ORT)
- Traffic Control Health Client (tc-health-client)
- Traffic Control Client <!-- Please specify which (Python, Go, or Java) -->
- Traffic Monitor
- Traffic Ops
- Traffic Stats
- Grove
- CDN in a Box - Enroller
- CI tests for Go components
- Build system - Go version in builder images

## What is the best way to verify this PR?
Run unit tests and API tests. Since this is only a patch-level version update, [the only changes]({MILESTONE_URL}?closed=1) were bugfixes. Breaking changes would be unexpected.

## The following criteria are ALL met by this PR
- [x] Existing tests are sufficient, no additional tests necessary
- [x] The documentation only mentions the major Go version, no documentation updates necessary.
- [x] The changelog already mentions updating to Go {GO_MAJOR_VERSION}, no additional changelog message necessary.
- [x] This PR **DOES NOT FIX A SERIOUS SECURITY VULNERABILITY** (see [the Apache Software Foundation's security guidelines](https://www.apache.org/security/) for details)

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
