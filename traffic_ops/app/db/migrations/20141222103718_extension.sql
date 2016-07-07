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

update serverstatus set last_recycle_date = NULL;

-- migrate the current table
alter table serverstatus 
  ALTER COLUMN ilo_pingable TYPE integer,
  ALTER COLUMN ilo_pingable DROP NOT NULL,
  ALTER COLUMN ilo_pingable SET DEFAULT NULL,
  RENAME COLUMN ilo_pingable TO aa,

  ALTER COLUMN teng_pingable TYPE integer,
  ALTER COLUMN teng_pingable DROP NOT NULL,
  ALTER COLUMN teng_pingable SET DEFAULT NULL,
  RENAME COLUMN teng_pingable TO ab,

  ALTER COLUMN fqdn_pingable TYPE integer,
  ALTER COLUMN fqdn_pingable DROP NOT NULL,
  ALTER COLUMN fqdn_pingable SET DEFAULT NULL,
  RENAME COLUMN fqdn_pingable TO ab,

  ALTER COLUMN dscp TYPE integer,
  ALTER COLUMN dscp DROP NOT NULL,
  ALTER COLUMN dscp SET DEFAULT NULL,
  RENAME COLUMN dscp TO ad,

  ALTER COLUMN firmware TYPE integer,
  ALTER COLUMN firmware DROP NOT NULL,
  ALTER COLUMN firmware SET DEFAULT NULL,
  RENAME COLUMN firmware TO ae,

  ALTER COLUMN marvin TYPE integer,
  ALTER COLUMN marvin DROP NOT NULL,
  ALTER COLUMN marvin SET DEFAULT NULL,
  RENAME COLUMN marvin TO af,

  ALTER COLUMN ping6 TYPE integer,
  ALTER COLUMN ping6 DROP NOT NULL,
  ALTER COLUMN ping6 SET DEFAULT NULL,
  RENAME COLUMN ping6 TO ag,

  ALTER COLUMN upd_pending TYPE integer,
  ALTER COLUMN upd_pending DROP NOT NULL,
  ALTER COLUMN upd_pending SET DEFAULT NULL,
  RENAME COLUMN upd_pending TO ah,

  ALTER COLUMN stats TYPE integer,
  ALTER COLUMN stats DROP NOT NULL,
  ALTER COLUMN stats SET DEFAULT NULL,
  RENAME COLUMN stats TO ai,

  ALTER COLUMN prox TYPE integer,
  ALTER COLUMN prox DROP NOT NULL,
  ALTER COLUMN prox SET DEFAULT NULL,
  RENAME COLUMN prox TO aj,

  ALTER COLUMN mtu TYPE integer,
  ALTER COLUMN mtu DROP NOT NULL,
  ALTER COLUMN mtu SET DEFAULT NULL,
  RENAME COLUMN mtu TO ak,

  ALTER COLUMN ccr_online TYPE integer,
  ALTER COLUMN ccr_online DROP NOT NULL,
  ALTER COLUMN ccr_online SET DEFAULT NULL,
  RENAME COLUMN ccr_online TO al,

  ALTER COLUMN rascal TYPE integer,
  ALTER COLUMN rascal DROP NOT NULL,
  ALTER COLUMN rascal SET DEFAULT NULL,
  RENAME COLUMN rascal TO am,

  ALTER COLUMN chr TYPE integer,
  ALTER COLUMN chr DROP NOT NULL,
  ALTER COLUMN chr SET DEFAULT NULL,
  RENAME COLUMN chr TO an,

  ALTER COLUMN cdu TYPE integer,
  ALTER COLUMN cdu DROP NOT NULL,
  ALTER COLUMN cdu SET DEFAULT NULL,
  RENAME COLUMN cdu TO ao,

  ALTER COLUMN ort_errors TYPE integer,
  ALTER COLUMN ort_errors DROP NOT NULL,
  ALTER COLUMN ort_errors SET DEFAULT NULL,
  RENAME COLUMN ort_errors TO ap,

  ALTER COLUMN mbps_out TYPE integer,
  ALTER COLUMN mbps_out DROP NOT NULL,
  ALTER COLUMN mbps_out SET DEFAULT NULL,
  RENAME COLUMN mbps_out TO aq,

  ALTER COLUMN clients_connected TYPE integer,
  ALTER COLUMN clients_connected DROP NOT NULL,
  ALTER COLUMN clients_connected SET DEFAULT NULL,
  RENAME COLUMN clients_connected TO ar,

  ALTER COLUMN last_recycle_date TYPE integer,
  ALTER COLUMN last_recycle_date DROP NOT NULL,
  ALTER COLUMN last_recycle_date SET DEFAULT NULL,
  RENAME COLUMN last_recycle_date TO as,

  ALTER COLUMN last_recycle_duration_hrs TYPE integer,
  ALTER COLUMN last_recycle_duration_hrs DROP NOT NULL,
  ALTER COLUMN last_recycle_duration_hrs SET DEFAULT NULL,
  RENAME COLUMN last_recycle_duration_hrs TO at;
