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

CREATE TABLE IF NOT EXISTS deliveryservices_required_capability (
required_capability TEXT NOT NULL,
deliveryservice_id bigint NOT NULL,
last_updated timestamp with time zone DEFAULT now() NOT NULL,

PRIMARY KEY (deliveryservice_id, required_capability)
);

DO $$
DECLARE temprow RECORD;
BEGIN FOR temprow IN
select id as deliveryservice_id, unnest(required_capabilities) as required_capability from deliveryservice d group by d.id, d.required_capabilities
    LOOP
insert into deliveryservices_required_capability ("deliveryservice_id", "required_capability") values (temprow.deliveryservice_id, temprow.required_capability);
END LOOP;
END $$;

ALTER TABLE public.deliveryservice
DROP COLUMN required_capabilities;