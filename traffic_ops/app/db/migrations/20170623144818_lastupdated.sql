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

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TRIGGER IF EXISTS on_update_current_timestamp ON asn;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON asn FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON cachegroup;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cachegroup FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON cachegroup_parameter;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cachegroup_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON cdn;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cdn FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON deliveryservice;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON deliveryservice_server;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON deliveryservice_tmuser;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_tmuser FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON division;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON division FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON federation;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON federation_deliveryservice;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_deliveryservice FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON federation_federation_resolver;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_federation_resolver FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON federation_resolver;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_resolver FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON federation_tmuser;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_tmuser FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON hwinfo;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON hwinfo FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON job;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON job_agent;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_agent FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON job_result;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_result FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON job_status;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON log;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON log FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON parameter;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON phys_location;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON phys_location FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON profile;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON profile_parameter;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON regex;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON regex FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON region;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON region FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON server;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON servercheck;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON servercheck FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON staticdnsentry;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON staticdnsentry FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON status;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON steering_target;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON steering_target FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON tm_user;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tm_user FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

DROP TRIGGER IF EXISTS on_update_current_timestamp ON type;
CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON type FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
