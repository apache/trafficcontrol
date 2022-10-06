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
ADD COLUMN regional BOOLEAN NOT NULL DEFAULT FALSE;

/* Set `regional` to `false` if it does not exist */
UPDATE public.deliveryservice_request
SET deliveryservice = deliveryservice || '{"regional": false}'
WHERE deliveryservice->>'regional' IS NULL;

/* Set `regional` to `false` it does not exist and `original` is not null */
UPDATE public.deliveryservice_request
SET original = original || '{"regional": false}'
WHERE original->>'regional' IS NULL;
