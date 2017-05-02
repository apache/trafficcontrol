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

-- THIS FILE INCLUDES DEMO DATA. IT BUILDS UPON THE RECORDS INSERTED IN SEEDS.SQL

-- parameters (demo)
insert into parameter (name, config_file, value) values ('tm.toolname', 'global', 'Traffic Ops') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.infourl', 'global', 'http://docs.cdnl.kabletown.net/traffic_control/html/') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.logourl', 'global', '/images/tc_logo.png') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('tm.instance_name', 'global', 'kabletown CDN') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation.polling.url', 'CRConfig.json', 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCity.dat.gz') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('geolocation6.polling.url', 'CRConfig.json', 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCityv6.dat.gz') ON CONFLICT (name, config_file, value) DO NOTHING;

insert into parameter (name, config_file, value) values ('tld.soa.admin', 'CRConfig.json', 'traffic_ops') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('health.polling.interval', 'rascal-config.txt', '8000') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threshold.loadavg', 'rascal.properties', '25.0') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('health.connection.timeout', 'rascal.properties', '25.0') ON CONFLICT (name, config_file, value) DO NOTHING;
insert into parameter (name, config_file, value) values ('health.threshold.availableBandwidthInKbps', 'rascal.properties', '>1750000') ON CONFLICT (name, config_file, value) DO NOTHING;

-- profile_parameters (demo)
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.url' and config_file = 'global' and value = 'https://tm.kabletown.net/') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.toolname' and config_file = 'global' and value = 'Traffic Ops') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.infourl' and config_file = 'global' and value = 'http://docs.cdnl.kabletown.net/traffic_control/html/') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.logourl' and config_file = 'global' and value = '/images/tc_logo.png') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.instance_name' and config_file = 'global' and value = 'kabletown CDN') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCity.dat.gz') ) ON CONFLICT (profile, parameter) DO NOTHING;
insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation6.polling.url' and config_file = 'CRConfig.json' and value = 'http://cdn-tools.cdnl.kabletown.net/cdn/MaxMind/GeoLiteCityv6.dat.gz') ) ON CONFLICT (profile, parameter) DO NOTHING;
