/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership.  The ASF
 * licenses this file to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

INSERT INTO public."role"
(id, name, description, priv_level, last_updated)
VALUES(236, 'trouter', 'Role utilized by the trouter user', 30, '2023-01-11 18:18:11.800');

INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'CDN:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'DELIVERY-SERVICE:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'DNS-SEC:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'FEDERATION:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'STEERING:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'FEDERATION-RESOLVER:READ', '2023-01-23 06:00:00.000');
INSERT INTO public.role_capability
(role_id, cap_name, last_updated)
VALUES(236, 'DS-SECURITY-KEY:READ', '2023-01-23 06:00:00.000');

UPDATE public.tm_user
SET username='trouter1', public_ssh_key=NULL, "role"=236, uid=NULL, gid=NULL, local_passwd='SCRYPT:16384:8:1:0/DsktE2F1KUJAwIcH3YHuuaOzQboEhu9tNruRBECD/h6yTGIVH5SoCr+VlhSttLbn82hB2SQucrkE1Uyw4svA==:Ww9mFRbIZAW3a1Ah6zI6WV0MVNsWmiffhyT3zEJ2R5/2IQqrwO067WRArcRmcGzzJdnYHqRynLzhREPiC33+xg==', confirm_local_passwd=NULL, last_updated='2023-01-12 02:10:15.053', company=NULL, email='cdn_ops2@comcast.com', full_name='Traffic Router1 Service Account', new_user=false, address_line1=NULL, address_line2=NULL, city=NULL, state_or_province=NULL, phone_number=NULL, postal_code=NULL, country=NULL, "token"=NULL, registration_sent=NULL, tenant_id=1, last_authenticated='2023-01-12 02:10:15.053', ucdn=''
WHERE id=1136;
