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

ALTER TABLE public.cdn_lock ADD CONSTRAINT cdn_lock_cdn_username_unique UNIQUE (username, cdn);
CREATE TABLE IF NOT EXISTS public.cdn_lock_user (
                                                     owner text NOT NULL,
                                                     cdn text NOT NULL,
                                                     username text NOT NULL,
    CONSTRAINT pk_cdn_lock_user PRIMARY KEY (owner, cdn, username),
    CONSTRAINT fk_shared_username FOREIGN KEY (username) REFERENCES public.tm_user(username),
    CONSTRAINT fk_owner FOREIGN KEY (owner, cdn) REFERENCES public.cdn_lock(username, cdn) ON DELETE CASCADE
    );
