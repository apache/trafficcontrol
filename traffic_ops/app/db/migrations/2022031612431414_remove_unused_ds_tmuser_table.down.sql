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

CREATE TABLE public.deliveryservice_tmuser (
	deliveryservice bigint NOT NULL,
	tm_user_id bigint NOT NULL,
	last_updated timestamp with time zone NOT NULL DEFAULT now(),
	CONSTRAINT idx_89525_primary PRIMARY KEY (deliveryservice, tm_user_id)
);

CREATE INDEX idx_89525_fk_tm_userid
ON public.deliveryservice_tmuser
USING btree (tm_user_id);

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE ON public.deliveryservice_tmuser
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('public.deliveryservice_tmuser');

CREATE TRIGGER on_update_current_timestamp
BEFORE UPDATE ON public.deliveryservice_tmuser
FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated('public.deliveryservice_tmuser');

ALTER TABLE ONLY public.deliveryservice_tmuser
ADD CONSTRAINT fk_tm_user_ds
FOREIGN KEY (deliveryservice)
REFERENCES deliveryservice(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

ALTER TABLE ONLY public.deliveryservice_tmuser
ADD CONSTRAINT fk_tm_user_id
FOREIGN KEY (tm_user_id)
REFERENCES tm_user(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
