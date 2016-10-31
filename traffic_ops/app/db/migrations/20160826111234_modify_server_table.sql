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
set @exist1 := (select count(*) from information_schema.statistics where table_name = 'server' and index_name = 'cs_ip_address_UNIQUE' and table_schema = database());
set @sqlstmt1 := if( @exist1 > 0, 'alter table server drop index cs_ip_address_UNIQUE', 'alter table server');
PREPARE stmt1 FROM @sqlstmt1;
EXECUTE stmt1;

set @exist2 := (select count(*) from information_schema.statistics where table_name = 'server' and index_name = 'ip6_address' and table_schema = database());
set @sqlstmt2 := if( @exist2 > 0, 'alter table server drop index ip6_address', 'alter table server');
PREPARE stmt2 FROM @sqlstmt2;
EXECUTE stmt2;

set @exist3 := (select count(*) from information_schema.statistics where table_name = 'server' and index_name = 'host_name' and table_schema = database());
set @sqlstmt3 := if( @exist3 > 0, 'alter table server drop index host_name', 'alter table server');
PREPARE stmt3 FROM @sqlstmt3;
EXECUTE stmt3;

alter table server modify host_name varchar(63) not null;
alter table server modify domain_name varchar(63) not null;

alter table server add unique key `ip_profile` (ip_address, profile);
alter table server add unique key `ip6_profile` (ip6_address, profile);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
create unique index cs_ip_address_UNIQUE on server (ip_address);
create unique index host_name on server (host_name);
create unique index ip6_address on server (ip6_address);
alter table server modify host_name varchar(45) not null;
alter table server modify domain_name varchar(45) not null;

set @exist4 := (select count(*) from information_schema.statistics where table_name = 'server' and index_name = 'ip_profile' and table_schema = database());
set @sqlstmt4 := if( @exist4 > 0, 'alter table server drop index ip_profile', 'alter table server');
PREPARE stmt4 FROM @sqlstmt4;
EXECUTE stmt4;

set @exist5 := (select count(*) from information_schema.statistics where table_name = 'server' and index_name = 'ip6_profile' and table_schema = database());
set @sqlstmt5 := if( @exist5 > 0, 'alter table server drop index ip6_profile', 'alter table server');
PREPARE stmt5 FROM @sqlstmt5;
EXECUTE stmt5;



