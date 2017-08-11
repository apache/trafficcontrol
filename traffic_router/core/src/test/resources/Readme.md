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


The Certificate Authority used for testing and installed in the trust store temporarily during RouterTest
was created using the guidelines found at this website

https://jamielinux.com/docs/openssl-certificate-authority

Keys were created using the same mechanism that happens when one clicks 'generate ssl keys' for a delivery service in traffic ops

openssl req -nodes -newkey rsa:2048 -keyout private/https-only-test.key -out csr/https-only-test.csr

The certificate is then signed like:
openssl ca -config ../intermediate/openssl.cnf -extensions server_cert -days 7000 -notext -md sha256 -in csr/https-only-test.csr -out certs/https-only-test.crt

then the encoded data in sslkeys.json and sslkeys-missing-1.json were like the following

for the 'crt' attribute
cat primary/certs/http-to-https-test.crt intermediate/certs/intermediate.crt | base64 -b 76

for the 'key' attribute
cat private/https-nocert.key | base64 -b 76

The Root CA was put into a java keystore file like:

keytool -import -trustcacerts -alias root -file root/certs/root.crt -keystore keystore.jks

And this file is now under resources and loaded by RouterTestd