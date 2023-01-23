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
DELETE FROM public.role_capability rc WHERE rc.cap_name='DELIVERY-SERVICE-SAFE:UPDATE' AND rc.role_id IN (SELECT id FROM public.role r WHERE r.priv_level < 20);

INSERT INTO public.role_capability (role_id, cap_name)
SELECT id, perm
FROM public.role
         CROSS JOIN ( VALUES
                          ('DELIVERY-SERVICE:UPDATE')
) AS perms(perm)
WHERE "priv_level" < 20 AND "priv_level" > 0
    ON CONFLICT DO NOTHING;




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

