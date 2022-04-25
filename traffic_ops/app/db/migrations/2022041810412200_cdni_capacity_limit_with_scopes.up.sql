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

CREATE TABLE IF NOT EXISTS public.cdni_limits (
                                id bigserial NOT NULL,
                                limit_id text NOT NULL,
                                scope_type text,
                                scope_value text[],
                                limit_type text NOT NULL,
                                maximum_hard bigint NOT NULL,
                                maximum_soft bigint NOT NULL,
                                telemetry_id text NOT NULL,
                                telemetry_metric text NOT NULL,
                                capability_id bigint NOT NULL,
                                last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdni_limits PRIMARY KEY (id),
    CONSTRAINT fk_cdni_limits_telemetry FOREIGN KEY (telemetry_id) REFERENCES cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_cdni_limits_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE
);

INSERT INTO public.cdni_limits (limit_id,
                         scope_type,
                         scope_value,
                         limit_type,
                         maximum_hard,
                         maximum_soft,
                         telemetry_id,
                         telemetry_metric,
                         capability_id)
SELECT CONCAT('host_limit_', chl.limit_type, '_', chl.telemetry_metric),
       'published-host',
       ARRAY[chl.host],
       chl.limit_type,
       chl.maximum_hard,
       chl.maximum_soft,
       chl.telemetry_id,
       chl.telemetry_metric,
       chl.capability_id
FROM public.cdni_host_limits as chl;

INSERT INTO public.cdni_limits (limit_id,
                         scope_type,
                         scope_value,
                         limit_type,
                         maximum_hard,
                         maximum_soft,
                         telemetry_id,
                         telemetry_metric,
                         capability_id)
SELECT CONCAT('total_limit_', thl.limit_type, '_', thl.telemetry_metric),
       NULL,
       NULL,
       thl.limit_type,
       thl.maximum_hard,
       thl.maximum_soft,
       thl.telemetry_id,
       thl.telemetry_metric,
       thl.capability_id
FROM public.cdni_total_limits as thl;

DROP TABLE IF EXISTS public.cdni_total_limits;
DROP TABLE IF EXISTS public.cdni_host_limits;

ALTER TABLE public.cdni_telemetry ADD COLUMN configuration_url text DEFAULT '';
