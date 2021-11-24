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

CREATE TABLE public.log_old (
	id bigint NOT NULL,
	"level" text,
	"message" text NOT NULL,
	tm_user bigint NOT NULL,
	ticketnum text,
	last_updated timestamp with time zone NOT NULL DEFAULT NOW(),
	CONSTRAINT idx_89634_primary PRIMARY KEY (id, tm_user)
);

CREATE SEQUENCE log_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

ALTER TABLE ONLY public.log_old
ALTER COLUMN id
SET DEFAULT nextval('log_id_seq'::regclass);

CREATE INDEX idx_89634_fk_log_1 ON public.log_old USING btree (tm_user);
CREATE INDEX idx_89634_idx_last_updated ON public.log_old USING btree (last_updated);

CREATE TRIGGER on_update_current_timestamp
BEFORE UPDATE ON public.log_old
FOR EACH ROW
	EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

CREATE TRIGGER on_delete_current_timestamp
BEFORE UPDATE ON public.log_old
FOR EACH ROW
	EXECUTE PROCEDURE on_delete_current_timestamp_last_updated();

ALTER TABLE ONLY public.log_old
ADD CONSTRAINT fk_log_1
FOREIGN KEY (tm_user) REFERENCES tm_user(id);

INSERT INTO public.log_old(
	"level",
	"message",
	tm_user,
	last_updated
)
SELECT 'APICHANGE', "message", "user", last_updated
FROM public.log
ORDER BY last_updated ASC;

DROP TABLE public.log CASCADE;

ALTER TABLE public.log_old RENAME TO log;
