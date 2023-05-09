package ldap

import "encoding/asn1"

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

const (
	iTUTT              = 0
	data               = 9
	pSS                = 2342
	uCL                = 19200300
	pilot              = 100
	pilotAttributeType = 1
	uID                = 1
)

// OIDType is the Object Identifier value for UID used within LDAP
// LDAP OID reference: https://ldap.com/ldap-oid-reference-guide/
// 0.9.2342.19200300.100.1.1	uid	Attribute Type (see RFC 4519)
var OIDType = asn1.ObjectIdentifier{iTUTT, data, pSS, uCL, pilot, pilotAttributeType, uID}
