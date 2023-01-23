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

UPDATE public.tm_user
SET username='trouter1', public_ssh_key=NULL, "role"=4, uid=NULL, gid=NULL, local_passwd='SCRYPT:16384:8:1:0/DsktE2F1KUJAwIcH3YHuuaOzQboEhu9tNruRBECD/h6yTGIVH5SoCr+VlhSttLbn82hB2SQucrkE1Uyw4svA==:Ww9mFRbIZAW3a1Ah6zI6WV0MVNsWmiffhyT3zEJ2R5/2IQqrwO067WRArcRmcGzzJdnYHqRynLzhREPiC33+xg==', confirm_local_passwd=NULL, last_updated='2023-01-12 02:10:15.053', company=NULL, email='cdn_ops2@comcast.com', full_name='Traffic Router1 Service Account', new_user=false, address_line1=NULL, address_line2=NULL, city=NULL, state_or_province=NULL, phone_number=NULL, postal_code=NULL, country=NULL, "token"=NULL, registration_sent=NULL, tenant_id=1, last_authenticated='2023-01-12 02:10:15.053', ucdn=''
WHERE id=1136;

DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='CDN:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='DELIVERY-SERVICE:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='DNS-SEC:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='FEDERATION:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='STEERING:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='FEDERATION-RESOLVER:READ';
DELETE FROM public.role_capability
WHERE role_id=236 AND cap_name='DS-SECURITY-KEY:READ';

DELETE FROM public."role"
WHERE id=236;