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

DROP TABLE IF EXISTS public.server_profile;

CREATE TRIGGER before_update_ip_address_trigger
    BEFORE UPDATE ON ip_address
    FOR EACH ROW WHEN (NEW.address <> OLD.address)
    EXECUTE PROCEDURE before_ip_address_table();

CREATE TRIGGER before_create_ip_address_trigger
    BEFORE INSERT ON ip_address
    FOR EACH ROW EXECUTE PROCEDURE before_ip_address_table();

CREATE TRIGGER before_update_server_trigger
    BEFORE UPDATE ON server
    FOR EACH ROW WHEN (NEW.profile <> OLD.profile)
    EXECUTE PROCEDURE before_server_table();

CREATE TRIGGER before_create_server_trigger
    BEFORE INSERT ON server
    FOR EACH ROW EXECUTE PROCEDURE before_server_table();
