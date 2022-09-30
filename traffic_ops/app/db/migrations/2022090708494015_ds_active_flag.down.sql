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

ALTER TABLE public.deliveryservice
ADD COLUMN active_flag boolean DEFAULT FALSE NOT NULL;

UPDATE public.deliveryservice
SET active_flag = FALSE
WHERE active = 'PRIMED' OR active = 'INACTIVE';

UPDATE public.deliveryservice
SET active_flag = TRUE
WHERE active = 'ACTIVE';

ALTER TABLE public.deliveryservice DROP COLUMN active;
ALTER TABLE public.deliveryservice RENAME COLUMN active_flag TO active;
DROP TYPE public.ds_active_state;

UPDATE public.deliveryservice_request
SET
	deliveryservice = jsonb_set(deliveryservice, '{active}', 'true')
WHERE
	deliveryservice IS NOT NULL
	AND
	deliveryservice ? 'active'
	AND
	deliveryservice ->> 'active' = 'ACTIVE';
UPDATE public.deliveryservice_request
SET
	deliveryservice = jsonb_set(deliveryservice, '{active}', 'false')
WHERE
	deliveryservice IS NOT NULL
	AND
	deliveryservice ? 'active'
	AND (
		deliveryservice ->> 'active' = 'PRIMED'
		OR
		deliveryservice ->> 'active' = 'INACTIVE'
	);
UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{active}', 'true')
WHERE
	original IS NOT NULL
	AND
	original ? 'active'
	AND
	original ->> 'active' = 'ACTIVE';
UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{active}', 'false')
WHERE
	original IS NOT NULL
	AND
	original ? 'active'
	AND (
		original ->> 'active' = 'PRIMED'
		OR
		original ->> 'active' = 'INACTIVE'
	);
