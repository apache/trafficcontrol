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

UPDATE public.deliveryservice_request
SET
    deliveryservice =
        CASE
            WHEN deliveryservice -> 'geoLimitCountries' = '""' THEN jsonb_set(deliveryservice, '{geoLimitCountries}', '[]')
            ELSE jsonb_set(deliveryservice, '{geoLimitCountries}', ('[' || replace((deliveryservice->'geoLimitCountries')::text, ',', '","') || ']')::jsonb)
        END
WHERE
    deliveryservice IS NOT NULL;

UPDATE public.deliveryservice_request
SET
    original =
        CASE
            WHEN original -> 'geoLimitCountries' = '""' THEN jsonb_set(original, '{geoLimitCountries}', '[]')
            ELSE jsonb_set(original, '{geoLimitCountries}', ('[' || replace((original->'geoLimitCountries')::text, ',', '","') || ']')::jsonb)
        END
WHERE
    original IS NOT NULL;
