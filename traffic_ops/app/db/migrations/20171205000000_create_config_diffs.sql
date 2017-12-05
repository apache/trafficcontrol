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

--
-- Name: config_diffs; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE config_diffs (
    config_id bigserial NOT NULL PRIMARY KEY,
    server bigint NOT NULL REFERENCES server (id) ON UPDATE CASCADE ON DELETE CASCADE,
    config_name text NOT NULL,
    db_lines_missing text[],
    disk_lines_missing text[],
    last_checked timestamp without time zone NOT NULL
);

ALTER TABLE config_diffs OWNER to traffic_ops;