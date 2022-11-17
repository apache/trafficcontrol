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
	deliveryservice = jsonb_set(deliveryservice, '{active}', '"ACTIVE"')
WHERE
	deliveryservice IS NOT NULL
	AND
	deliveryservice ? 'active'
	AND
	(deliveryservice -> 'active')::boolean IS TRUE;
UPDATE public.deliveryservice_request
SET
	deliveryservice = jsonb_set(deliveryservice, '{active}', '"PRIMED"')
WHERE
	deliveryservice IS NOT NULL
	AND
	deliveryservice ? 'active'
	AND
	(deliveryservice -> 'active')::boolean IS FALSE;
UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{active}', '"ACTIVE"')
WHERE
	original IS NOT NULL
	AND
	original ? 'active'
	AND
	(original -> 'active')::boolean IS TRUE;
UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{active}', '"PRIMED"')
WHERE
	original IS NOT NULL
	AND
	original ? 'active'
	AND
	(original -> 'active')::boolean IS FALSE;
