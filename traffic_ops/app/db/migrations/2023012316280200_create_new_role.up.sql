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

INSERT INTO public."role" (name, description, priv_level)
VALUES('trouter', 'Limited role for Traffic Router calls to Traffic Ops', 10);

INSERT INTO public.role_capability (role_id, cap_name)
    VALUES (
        (SELECT id FROM role WHERE name='trouter'),
        UNNEST(ARRAY[
            'CDN:READ',
            'DELIVERY-SERVICE:READ',
            'DNS-SEC:READ',
            'STEERING:READ',
            'FEDERATION-RESOLVER:READ',
            'DS-SECURITY-KEY:READ']
        )
    );
