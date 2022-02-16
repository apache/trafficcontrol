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

CREATE TABLE IF NOT EXISTS public.server_config_update ( 
    server_id bigint NOT NULL PRIMARY KEY, 
    config_update_time TIMESTAMPTZ, 
    config_apply_time TIMESTAMPTZ, 
    revalidate_update_time TIMESTAMPTZ, 
    revalidate_apply_time TIMESTAMPTZ, 
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_server_id FOREIGN KEY(server_id) REFERENCES public.server(id) ON DELETE CASCADE
);

ALTER TABLE public.server DROP COLUMN IF EXISTS upd_pending, DROP COLUMN IF EXISTS reval_pending;