/*

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

CREATE SEQUENCE IF NOT EXISTS async_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS async_status (
    id bigint NOT NULL DEFAULT nextval('async_status_id_seq'::regclass),
    status TEXT NOT NULL,
    message TEXT,
    start_time timestamp with time zone DEFAULT now() NOT NULL,
    end_time timestamp with time zone,

    PRIMARY KEY (id)
);

ALTER SEQUENCE async_status_id_seq OWNED BY async_status.id;
