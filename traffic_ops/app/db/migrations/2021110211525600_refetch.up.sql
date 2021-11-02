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

-- Migrate

/**
 * Alter job_agent table and data
 */

DROP TABLE IF EXISTS public.job_agent CASCADE;

/**
 * Alter job_status table and data
 */

DROP TABLE IF EXISTS public.job_status CASCADE;

/**
 * Alter job table and data
 */

ALTER TABLE public.job
DROP COLUMN IF EXISTS agent CASCADE,
DROP COLUMN IF EXISTS object_type CASCADE,
DROP COLUMN IF EXISTS object_name CASCADE,
DROP COLUMN IF EXISTS keyword CASCADE,
DROP COLUMN IF EXISTS asset_type CASCADE,
DROP COLUMN IF EXISTS status CASCADE;

ALTER TABLE public.job 
ADD COLUMN IF NOT EXISTS invalidation_type text NOT NULL DEFAULT 'REFRESH';

/*
 * If the asset_url contains the temporary fix for refetch
 * (adding ##REFETCH## to the end of the url) then assign the
 * correct invalidation_type.
 */
UPDATE public.job
SET invalidation_type = 'REFETCH'
WHERE asset_url LIKE '%##REFETCH##%';

ALTER TABLE public.job
RENAME COLUMN parameters TO ttl_hr;

ALTER TABLE public.job 
ALTER COLUMN ttl_hr TYPE integer
USING CAST(substring(ttl_hr,'TTL:([0-9]+)h') AS integer);
