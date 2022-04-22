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

/* Downgrading this will lose some data because the specification has changed.
 * Previously, only total and host limits were set. Now with the scope field, it is not as limited so downgrading
 * cannot put the data back in the old tables since it wont necessarily fit into those buckets.
*/
CREATE TABLE IF NOT EXISTS cdni_total_limits (
                                                 limit_type text NOT NULL,
                                                 maximum_hard bigint NOT NULL,
                                                 maximum_soft bigint NOT NULL,
                                                 telemetry_id text NOT NULL,
                                                 telemetry_metric text NOT NULL,
                                                 capability_id bigint NOT NULL,
                                                 last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdni_total_limits PRIMARY KEY (capability_id, telemetry_id),
    CONSTRAINT fk_cdni_total_limits_telemetry FOREIGN KEY (telemetry_id) REFERENCES cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_cdni_total_limits_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS cdni_host_limits (
                                                limit_type text NOT NULL,
                                                maximum_hard bigint NOT NULL,
                                                maximum_soft bigint NOT NULL,
                                                telemetry_id text NOT NULL,
                                                telemetry_metric text NOT NULL,
                                                capability_id bigint NOT NULL,
                                                host text NOT NULL,
                                                last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdni_host_limits PRIMARY KEY (capability_id, telemetry_id, host),
    CONSTRAINT fk_cdni_host_limits_telemetry FOREIGN KEY (telemetry_id) REFERENCES cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_cdni_total_limits_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE
);

DROP TABLE IF EXISTS cdni_limits;

ALTER TABLE cdni_telemetry DROP COLUMN IF EXISTS configuration_url;
