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

-- remove from role_capability table
DELETE FROM public.role_capability WHERE cap_name='ASN:CREATE';
DELETE FROM public.role_capability WHERE cap_name='ASN:DELETE';
DELETE FROM public.role_capability WHERE cap_name='ASN:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='CACHE-GROUP:CREATE';
DELETE FROM public.role_capability WHERE cap_name='CACHE-GROUP:DELETE';
DELETE FROM public.role_capability WHERE cap_name='CACHE-GROUP:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='CDN-LOCK:CREATE';
DELETE FROM public.role_capability WHERE cap_name='CDN-LOCK:DELETE';
DELETE FROM public.role_capability WHERE cap_name='CDN-SNAPSHOT:CREATE';
DELETE FROM public.role_capability WHERE cap_name='CDN:CREATE';
DELETE FROM public.role_capability WHERE cap_name='CDN:DELETE';
DELETE FROM public.role_capability WHERE cap_name='CDN:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='COORDINATE:CREATE';
DELETE FROM public.role_capability WHERE cap_name='COORDINATE:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='COORDINATE:DELETE';
DELETE FROM public.role_capability WHERE cap_name='DELIVERY-SERVICE-SAFE:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='DELIVERY-SERVICE:CREATE';
DELETE FROM public.role_capability WHERE cap_name='DELIVERY-SERVICE:DELETE';
DELETE FROM public.role_capability WHERE cap_name='DIVISION:CREATE';
DELETE FROM public.role_capability WHERE cap_name='DIVISION:DELETE';
DELETE FROM public.role_capability WHERE cap_name='DIVISION:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='DNS-SEC:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='ISO:GENERATE';
DELETE FROM public.role_capability WHERE cap_name='ORIGIN:CREATE';
DELETE FROM public.role_capability WHERE cap_name='ORIGIN:DELETE';
DELETE FROM public.role_capability WHERE cap_name='ORIGIN:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='PARAMETER:CREATE';
DELETE FROM public.role_capability WHERE cap_name='PARAMETER:DELETE';
DELETE FROM public.role_capability WHERE cap_name='PARAMETER:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='PHYSICAL-LOCATION:CREATE';
DELETE FROM public.role_capability WHERE cap_name='PHYSICAL-LOCATION:DELETE';
DELETE FROM public.role_capability WHERE cap_name='PHYSICAL-LOCATION:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='PROFILE:CREATE';
DELETE FROM public.role_capability WHERE cap_name='PROFILE:DELETE';
DELETE FROM public.role_capability WHERE cap_name='PROFILE:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='REGION:CREATE';
DELETE FROM public.role_capability WHERE cap_name='REGION:DELETE';
DELETE FROM public.role_capability WHERE cap_name='REGION:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CAPABILITY:CREATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CAPABILITY:DELETE';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CAPABILITY:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER:CREATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER:DELETE';
DELETE FROM public.role_capability WHERE cap_name='SERVER:QUEUE';
DELETE FROM public.role_capability WHERE cap_name='SERVER:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='SERVICE-CATEGORY:CREATE';
DELETE FROM public.role_capability WHERE cap_name='SERVICE-CATEGORY:DELETE';
DELETE FROM public.role_capability WHERE cap_name='SERVICE-CATEGORY:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='STATIC-DN:CREATE';
DELETE FROM public.role_capability WHERE cap_name='STATIC-DN:DELETE';
DELETE FROM public.role_capability WHERE cap_name='STATIC-DN:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='STATUS:CREATE';
DELETE FROM public.role_capability WHERE cap_name='STATUS:DELETE';
DELETE FROM public.role_capability WHERE cap_name='STATUS:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='TENANT:CREATE';
DELETE FROM public.role_capability WHERE cap_name='TENANT:DELETE';
DELETE FROM public.role_capability WHERE cap_name='TENANT:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='TOPOLOGY:CREATE';
DELETE FROM public.role_capability WHERE cap_name='TOPOLOGY:DELETE';
DELETE FROM public.role_capability WHERE cap_name='TOPOLOGY:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='TYPE:CREATE';
DELETE FROM public.role_capability WHERE cap_name='TYPE:DELETE';
DELETE FROM public.role_capability WHERE cap_name='TYPE:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='USER:CREATE';
DELETE FROM public.role_capability WHERE cap_name='USER:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CHECK:CREATE';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CHECK:DELETE';
DELETE FROM public.role_capability WHERE cap_name='STAT:CREATE';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION:CREATE';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION:DELETE';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION-RESOLVER:CREATE';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION-RESOLVER:DELETE';
DELETE FROM public.role_capability WHERE cap_name='DELIVERY-SERVICE:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='JOB:CREATE';
DELETE FROM public.role_capability WHERE cap_name='JOB:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='JOB:DELETE';
DELETE FROM public.role_capability WHERE cap_name='DS-REQUEST:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='DS-REQUEST:CREATE';
DELETE FROM public.role_capability WHERE cap_name='DS-REQUEST:DELETE';
DELETE FROM public.role_capability WHERE cap_name='STEERING:CREATE';
DELETE FROM public.role_capability WHERE cap_name='STEERING:UPDATE';
DELETE FROM public.role_capability WHERE cap_name='STEERING:DELETE';
DELETE FROM public.role_capability WHERE cap_name='ASN:READ';
DELETE FROM public.role_capability WHERE cap_name='ASYNC-STATUS:READ';
DELETE FROM public.role_capability WHERE cap_name='CACHE-GROUP:READ';
DELETE FROM public.role_capability WHERE cap_name='CAPABILITY:READ';
DELETE FROM public.role_capability WHERE cap_name='CDN-SNAPSHOT:READ';
DELETE FROM public.role_capability WHERE cap_name='CDN:READ';
DELETE FROM public.role_capability WHERE cap_name='COORDINATE:READ';
DELETE FROM public.role_capability WHERE cap_name='DELIVERY-SERVICE:READ';
DELETE FROM public.role_capability WHERE cap_name='DIVISION:READ';
DELETE FROM public.role_capability WHERE cap_name='DS-REQUEST:READ';
DELETE FROM public.role_capability WHERE cap_name='DS-SECURITY-KEY:READ';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION:READ';
DELETE FROM public.role_capability WHERE cap_name='FEDERATION-RESOLVER:READ';
DELETE FROM public.role_capability WHERE cap_name='ISO:READ';
DELETE FROM public.role_capability WHERE cap_name='JOB:READ';
DELETE FROM public.role_capability WHERE cap_name='LOG:READ';
DELETE FROM public.role_capability WHERE cap_name='MONITOR-CONFIG:READ';
DELETE FROM public.role_capability WHERE cap_name='ORIGIN:READ';
DELETE FROM public.role_capability WHERE cap_name='PARAMETER:READ';
DELETE FROM public.role_capability WHERE cap_name='PHYSICAL-LOCATION:READ';
DELETE FROM public.role_capability WHERE cap_name='PLUGIN-READ';
DELETE FROM public.role_capability WHERE cap_name='PROFILE:READ';
DELETE FROM public.role_capability WHERE cap_name='REGION:READ';
DELETE FROM public.role_capability WHERE cap_name='ROLE:READ';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CAPABILITY:READ';
DELETE FROM public.role_capability WHERE cap_name='SERVER:READ';
DELETE FROM public.role_capability WHERE cap_name='SERVICE-CATEGORY:READ';
DELETE FROM public.role_capability WHERE cap_name='STATIC-DN:READ';
DELETE FROM public.role_capability WHERE cap_name='STATUS:READ';
DELETE FROM public.role_capability WHERE cap_name='SERVER-CHECK:READ';
DELETE FROM public.role_capability WHERE cap_name='STEERING:READ';
DELETE FROM public.role_capability WHERE cap_name='STAT:READ';
DELETE FROM public.role_capability WHERE cap_name='TENANT:READ';
DELETE FROM public.role_capability WHERE cap_name='TOPOLOGY:READ';
DELETE FROM public.role_capability WHERE cap_name='TRAFFIC-VAULT:READ';
DELETE FROM public.role_capability WHERE cap_name='TYPE:READ';
DELETE FROM public.role_capability WHERE cap_name='USER:READ';
