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

CREATE TABLE IF NOT EXISTS cdni_capability_updates (
                                                 id bigserial NOT NULL,
                                                 request_type text NOT NULL,
                                                 ucdn text NOT NULL,
                                                 host text,
                                                 data json NOT NULL,
                                                 async_status_id bigint NOT NULL,
                                                 last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdni_capability_updates PRIMARY KEY (id),
    CONSTRAINT fk_cdni_capability_updates_async FOREIGN KEY (async_status_id) REFERENCES async_status(id) ON UPDATE CASCADE ON DELETE CASCADE
);
