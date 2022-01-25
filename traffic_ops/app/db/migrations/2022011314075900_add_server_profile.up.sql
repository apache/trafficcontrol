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

CREATE TABLE IF NOT EXISTS public.server_profile (
                                                     server bigint NOT NULL,
                                                     profile_name text[] NOT NULL,
                                                     order bigint[] NOT NULL

    CONSTRAINT pk_server_profile PRIMARY KEY(profile_name, server, order)
    CONSTRAINT fk_server_id FOREIGN KEY (server) REFERENCES server(id)
    CONSTRAINT fk_profile_name FOREIGN KEY (profile_name) REFERENCES profile(name)
    );

/*

INSERT INTO public.server_profile(server_id, profile_name, order)
SELECT id, profile FROM server

*/