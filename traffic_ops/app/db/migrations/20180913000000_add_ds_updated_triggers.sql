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

-- these triggers update the deliveryservice timestamp when a ds_regex, ds_server, or staticdnsentry is deleted. This is necessary to update the DS timestamp in the CRConfig, to notify the Traffic Router that the DS needs updated. Updates and inserts can be queried to aggregate their times, but deletes can't otherwise be detected.
-- these triggers MUST be removed when the Self Service Timestamp plan is implemented, to avoid mangling the deliveryservice history. See https://cwiki.apache.org/confluence/display/TC/Traffic+Control++Self+Service+Proposal+for+Change+Integrity

create or replace function delete_update_dsf() returns trigger as $$
begin
  update deliveryservice set last_updated=now() where id = old.deliveryservice;
  return old;
end;
$$ language plpgsql;
create trigger dss_delete_update_ds after delete on deliveryservice_server for each row execute procedure delete_update_dsf();
create trigger sde_delete_update_ds after delete on staticdnsentry for each row execute procedure delete_update_dsf();
create trigger dsr_delete_update_ds after delete on deliveryservice_regex for each row execute procedure delete_update_dsf();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TRIGGER IF EXISTS dss_delete_update_ds ON deliveryservice_server;
DROP TRIGGER IF EXISTS sde_delete_update_ds ON staticdnsentry;
DROP TRIGGER IF EXISTS dsr_delete_update_ds ON deliveryservice_regex;
DROP FUNCTION IF EXISTS delete_update_dsf();
