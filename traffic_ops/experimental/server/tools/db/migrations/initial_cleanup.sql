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

alter table server rename column cdn_id to cdn;
alter table deliveryservice rename column cdn_id to cdn;
alter table deliveryservice_tmuser rename column tm_user_id to tm_user;
alter table hwinfo rename column serverid to server;
alter table job rename column agent to job_agent;
alter table job_result rename column agent to job_agent;
alter table job rename column job_user to tm_user;
alter table job rename column job_deliveryservice to deliveryservice;
alter table cachegroup rename column secondary_parent_cachegroup_id to secondary_parent_cachegroup;
alter table cachegroup rename column parent_cachegroup_id to parent_cachegroup;
