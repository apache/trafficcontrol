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

ALTER TABLE public.deliveryservice
    ADD COLUMN required_capabilities TEXT[];

DO $$
DECLARE temprow RECORD;
BEGIN FOR temprow IN
SELECT deliveryservice_id, ARRAY_AGG(required_capability) AS required_capabilities FROM deliveryservices_required_capability drc GROUP BY drc.deliveryservice_id
    LOOP
UPDATE deliveryservice d SET required_capabilities = temprow.required_capabilities WHERE d.id = temprow.deliveryservice_id;
END LOOP;
END $$;

DROP TABLE public.deliveryservices_required_capability;
