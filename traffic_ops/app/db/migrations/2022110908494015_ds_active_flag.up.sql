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

CREATE TYPE public.ds_active_state AS ENUM (
	'ACTIVE',
	'INACTIVE',
	'PRIMED'
);

ALTER TABLE public.deliveryservice
ADD COLUMN active_state ds_active_state NOT NULL DEFAULT 'INACTIVE';

UPDATE public.deliveryservice SET active_state = 'ACTIVE' WHERE active IS TRUE;
UPDATE public.deliveryservice SET active_state = 'PRIMED' WHERE active IS FALSE;

ALTER TABLE public.deliveryservice DROP COLUMN active;
ALTER TABLE public.deliveryservice RENAME COLUMN active_state TO active;

UPDATE public.deliveryservice_request
SET
	deliveryservice = deliveryservice || '{"active": "ACTIVE"}'
WHERE
	deliveryservice ->> 'active' = 'true';

UPDATE public.deliveryservice_request
SET
	deliveryservice = deliveryservice || '{"active": "PRIMED"}'
WHERE
	deliveryservice ->> 'active' = 'false';

UPDATE public.deliveryservice_request
SET
	original = original || '{"active": "ACTIVE"}'
WHERE
	original ->> 'active' = 'true';

UPDATE public.deliveryservice_request
SET
	original = original || '{"active": "PRIMED"}'
WHERE
	original ->> 'active' = 'false';

UPDATE public.deliveryservice_request
SET
	original = original || CAST('{"lastUpdated": "' || replace(replace(original ->> 'lastUpdated', ' ', 'T'), '+00', 'Z') || '"}' AS jsonb)
WHERE
	original ->> 'lastUpdated' IS NOT NULL;

UPDATE public.deliveryservice_request
SET
	deliveryservice = deliveryservice || CAST('{"lastUpdated": "' || replace(replace(deliveryservice ->> 'lastUpdated', ' ', 'T'), '+00', 'Z') || '"}' AS jsonb)
WHERE
	deliveryservice ->> 'lastUpdated' IS NOT NULL;
