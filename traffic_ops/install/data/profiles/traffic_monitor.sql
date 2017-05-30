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

INSERT INTO profile (name, description, type) VALUES ('TM_PROFILE','Traffic Monitor','TM_PROFILE') ON CONFLICT (name) DO NOTHING;
INSERT INTO parameter (name, config_file, value) VALUES ('hack.ttl','rascal-config.txt','30') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'hack.ttl' and config_file = 'rascal-config.txt' and value = '30') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.event-count','rascal-config.txt','200') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'health.event-count' and config_file = 'rascal-config.txt' and value = '200') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.polling.interval','rascal-config.txt','6000') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'health.polling.interval' and config_file = 'rascal-config.txt' and value = '6000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.threadPool','rascal-config.txt','4') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'health.threadPool' and config_file = 'rascal-config.txt' and value = '4') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.timepad','rascal-config.txt','0') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'health.timepad' and config_file = 'rascal-config.txt' and value = '0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('heartbeat.polling.interval','rascal-config.txt','2000') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'heartbeat.polling.interval' and config_file = 'rascal-config.txt' and value = '2000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','rascal-config.txt','/opt/traffic_monitor/conf') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'location' and config_file = 'rascal-config.txt' and value = '/opt/traffic_monitor/conf') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('peers.polling.interval','rascal-config.txt','1000') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'peers.polling.interval' and config_file = 'rascal-config.txt' and value = '1000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('tm.crConfig.polling.url','rascal-config.txt','https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'tm.crConfig.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.xml') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('tm.dataServer.polling.url','rascal-config.txt','https://${tmHostname}/dataserver/orderby/id') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'tm.dataServer.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/dataserver/orderby/id') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('tm.healthParams.polling.url','rascal-config.txt','https://${tmHostname}/health/${cdnName}') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'tm.healthParams.polling.url' and config_file = 'rascal-config.txt' and value = 'https://${tmHostname}/health/${cdnName}') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('tm.polling.interval','rascal-config.txt','60000') ON CONFLICT (name, config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'TM_PROFILE'), (select id from parameter where name = 'tm.polling.interval' and config_file = 'rascal-config.txt' and value = '60000') )  ON CONFLICT (profile, parameter) DO NOTHING;

