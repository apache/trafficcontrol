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


A simple authentication server written in go that authenticates user agains the `tm_user` table and returns a jwt representing the user, incl. its API access capabilities, derived from the user's role.

* To run:
`go run auth.go auth.config my-secret`
`secret` is used for jwt signing

* To login:
`curl --insecure -X POST -Lkvs --header "Content-Type:application/json" https://localhost:9004/login -d'{"username":"username", "password":"password"}'`
