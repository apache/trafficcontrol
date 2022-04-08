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

-- Add previously existing columns
ALTER TABLE public.server 
ADD COLUMN IF NOT EXISTS upd_pending bool NOT NULL DEFAULT false, 
ADD COLUMN IF NOT EXISTS reval_pending bool NOT NULL DEFAULT false;

-- Remove new columns
ALTER TABLE public.server 
DROP COLUMN IF EXISTS config_update_time, 
DROP COLUMN IF EXISTS config_apply_time,
DROP COLUMN IF EXISTS revalidate_update_time, 
DROP COLUMN IF EXISTS revalidate_apply_time;
