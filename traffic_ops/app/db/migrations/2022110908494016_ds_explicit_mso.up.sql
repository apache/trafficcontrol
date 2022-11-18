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

UPDATE public.deliveryservice
SET multi_site_origin = FALSE
WHERE multi_site_origin IS NULL;

ALTER TABLE public.deliveryservice
ALTER COLUMN multi_site_origin
SET NOT NULL;

UPDATE public.deliveryservice_request
SET
	deliveryservice = jsonb_set(deliveryservice, '{multiSiteOrigin}', 'false')
WHERE
	deliveryservice IS NOT NULL
	AND
	(
		NOT (deliveryservice ? 'multiSiteOrigin')
		OR
		jsonb_typeof(deliveryservice -> 'multiSiteOrigin') = 'null'
	);

UPDATE public.deliveryservice_request
SET
	original = jsonb_set(original, '{multiSiteOrigin}', 'false')
WHERE
	original IS NOT NULL
	AND
	(
		NOT (original ? 'multiSiteOrigin')
		OR
		jsonb_typeof(original -> 'multiSiteOrigin') = 'null'
	);
