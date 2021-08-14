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

-- Revert

/**
 * Restore job_agent table with constraints and data.
 */

CREATE TABLE public.job_agent (
	id bigserial NOT NULL,
	"name" text NULL,
	description text NULL,
	active int4 NOT NULL DEFAULT 0,
	last_updated timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT idx_job_agent_id_primary PRIMARY KEY (id),
	CONSTRAINT job_agent_name_unique UNIQUE (name)
);
CREATE INDEX job_agent_last_updated_idx ON public.job_agent USING btree (last_updated DESC NULLS LAST);

-- Table Triggers

create trigger on_delete_current_timestamp after
delete
    on
    public.job_agent for each row execute function on_delete_current_timestamp_last_updated('job_agent');
create trigger on_update_current_timestamp before
update
    on
    public.job_agent for each row execute function on_update_current_timestamp_last_updated();

-- Permissions

ALTER TABLE public.job_agent OWNER TO traffic_ops;
GRANT ALL ON TABLE public.job_agent TO traffic_ops;

-- Data

INSERT INTO public.job_agent (name, description, active)
VALUES ('Example', 'Description of Purge Agent', 1) 
ON CONFLICT (name) DO NOTHING;

/**
 * Restore job_status table with constraints and data
 */

CREATE TABLE public.job_status (
	id bigserial NOT NULL,
	"name" text NULL,
	description text NULL,
	last_updated timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT idx_job_status_id_primary PRIMARY KEY (id),
	CONSTRAINT job_status_name_unique UNIQUE (name)
);
CREATE INDEX job_status_last_updated_idx ON public.job_status USING btree (last_updated DESC NULLS LAST);

-- Table Triggers

create trigger on_delete_current_timestamp after
delete
    on
    public.job_status for each row execute function on_delete_current_timestamp_last_updated('job_status');
create trigger on_update_current_timestamp before
update
    on
    public.job_status for each row execute function on_update_current_timestamp_last_updated();

-- Permissions

ALTER TABLE public.job_status OWNER TO traffic_ops;
GRANT ALL ON TABLE public.job_status TO traffic_ops;

-- Data

INSERT INTO job_status (name, description)
VALUES ('PENDING', 'Job is queued, but has not been picked up by any agents yet'),
('IN_PROGRESS', 'Job is being processed by agents'),
('COMPLETED', 'Job has finished'),
('CANCELLED', 'Job was cancelled')
ON CONFLICT (name) DO NOTHING;

/**
 * Restore job table with constraints and data
 */

-- Restore table
ALTER TABLE public.job 
ADD COLUMN IF NOT EXISTS agent int8 NULL,
ADD COLUMN IF NOT EXISTS object_type text NULL,
ADD COLUMN IF NOT EXISTS object_name text NULL,
ADD COLUMN IF NOT EXISTS keyword text NULL,
ADD COLUMN IF NOT EXISTS asset_type text NULL,
ADD COLUMN IF NOT EXISTS status int8 NULL;

ALTER TABLE public.job 
DROP COLUMN IF EXISTS invalidation_type CASCADE;

-- Restore indices and foreign key constraints
CREATE INDEX idx_job_fk_job_agent_id ON public.job USING btree (agent);
ALTER TABLE public.job ADD CONSTRAINT job_fk_job_agent_id FOREIGN KEY (agent) REFERENCES public.job_agent(id) ON DELETE CASCADE;

CREATE INDEX idx_job_fk_job_status_id ON public.job USING btree (status);
ALTER TABLE public.job ADD CONSTRAINT job_fk_job_status_id FOREIGN KEY (status) REFERENCES public.job_status(id);

-- Restore data
UPDATE public.job 
SET agent = 1;

UPDATE public.job 
SET status = 1;
ALTER TABLE public.job 
ALTER COLUMN status SET NOT NULL;

UPDATE public.job 
SET keyword = 'PURGE';
ALTER TABLE public.job 
ALTER COLUMN keyword SET NOT NULL;

UPDATE public.job 
SET asset_type = 'file';
ALTER TABLE public.job 
ALTER COLUMN asset_type SET NOT NULL;

ALTER TABLE public.job
RENAME COLUMN ttl_hr TO parameters;

ALTER TABLE public.job 
ALTER COLUMN parameters TYPE TEXT
USING CONCAT('TTL:', parameters, 'h');
