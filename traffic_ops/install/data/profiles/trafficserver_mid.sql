INSERT INTO profile (name, description, type) VALUES ('MID_PROFILE','Mid Cache','ATS_PROFILE') ON CONFLICT (name) DO NOTHING;
INSERT INTO parameter (name, config_file, value) VALUES ('astats_over_http','package','1.2-8.el6.x86_64') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'astats_over_http' and config_file = 'package' and value = '1.2-8.el6.x86_64') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('astats_over_http.so','plugin.config','') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'astats_over_http.so' and config_file = 'plugin.config' and value = '') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('coalesce_masklen_v4','ip_allow.config','16') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'coalesce_masklen_v4' and config_file = 'ip_allow.config' and value = '16') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('coalesce_masklen_v6','ip_allow.config','40') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'coalesce_masklen_v6' and config_file = 'ip_allow.config' and value = '40') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('coalesce_number_v4','ip_allow.config','5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'coalesce_number_v4' and config_file = 'ip_allow.config' and value = '5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('coalesce_number_v6','ip_allow.config','5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'coalesce_number_v6' and config_file = 'ip_allow.config' and value = '5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.accept_threads','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.accept_threads' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.admin.admin_user','records.config','STRING admin') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.admin.admin_user' and config_file = 'records.config' and value = 'STRING admin') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.admin.autoconf_port','records.config','INT 8083') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.admin.autoconf_port' and config_file = 'records.config' and value = 'INT 8083') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.admin.number_config_bak','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.admin.number_config_bak' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.admin.user_id','records.config','STRING ats') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.admin.user_id' and config_file = 'records.config' and value = 'STRING ats') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.alarm.abs_path','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.alarm.abs_path' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.alarm.bin','records.config','STRING example_alarm_bin.sh') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.alarm.bin' and config_file = 'records.config' and value = 'STRING example_alarm_bin.sh') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.alarm_email','records.config','STRING ats') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.alarm_email' and config_file = 'records.config' and value = 'STRING ats') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.allocator.debug_filter','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.allocator.debug_filter' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.allocator.enable_reclaim','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.allocator.enable_reclaim' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.allocator.max_overage','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.allocator.max_overage' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.allocator.thread_freelist_size','records.config','INT 1024') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.allocator.thread_freelist_size' and config_file = 'records.config' and value = 'INT 1024') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.body_factory.enable_customizations','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.enable_customizations' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.body_factory.enable_logging','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.enable_logging' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.body_factory.response_suppression_mode','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.response_suppression_mode' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.body_factory.template_sets_dir','records.config','STRING etc/trafficserver/body_factory') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.body_factory.template_sets_dir' and config_file = 'records.config' and value = 'STRING etc/trafficserver/body_factory') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.control.filename','records.config','STRING cache.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.control.filename' and config_file = 'records.config' and value = 'STRING cache.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.enable_read_while_writer','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.enable_read_while_writer' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.hosting_filename','records.config','STRING hosting.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.hosting_filename' and config_file = 'records.config' and value = 'STRING hosting.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.http.compatibility.4-2-0-fixup' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ip_allow.filename','records.config','STRING ip_allow.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ip_allow.filename' and config_file = 'records.config' and value = 'STRING ip_allow.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.limits.http.max_alts','records.config','INT 5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.limits.http.max_alts' and config_file = 'records.config' and value = 'INT 5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.max_doc_size','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.max_doc_size' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.min_average_object_size','records.config','INT 131072') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.min_average_object_size' and config_file = 'records.config' and value = 'INT 131072') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.mutex_retry_delay','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.mutex_retry_delay' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.permit.pinning','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.permit.pinning' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ram_cache.algorithm','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.algorithm' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ram_cache.compress','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.compress' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ram_cache_cutoff','records.config','INT 1073741824') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache_cutoff' and config_file = 'records.config' and value = 'INT 1073741824') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ram_cache.size','records.config','INT 34359738368') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.size' and config_file = 'records.config' and value = 'INT 34359738368') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.ram_cache.use_seen_filter','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.ram_cache.use_seen_filter' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.target_fragment_size','records.config','INT 1048576') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.target_fragment_size' and config_file = 'records.config' and value = 'INT 1048576') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cache.threads_per_disk','records.config','INT 8') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cache.threads_per_disk' and config_file = 'records.config' and value = 'INT 8') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.cluster_configuration ','records.config','STRING cluster.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.cluster_configuration ' and config_file = 'records.config' and value = 'STRING cluster.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.cluster_port','records.config','INT 8086') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.cluster_port' and config_file = 'records.config' and value = 'INT 8086') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.ethernet_interface','records.config','STRING lo') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.ethernet_interface' and config_file = 'records.config' and value = 'STRING lo') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.log_bogus_mc_msgs','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.log_bogus_mc_msgs' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.mc_group_addr','records.config','STRING 224.0.1.37') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mc_group_addr' and config_file = 'records.config' and value = 'STRING 224.0.1.37') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.mcport','records.config','INT 8089') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mcport' and config_file = 'records.config' and value = 'INT 8089') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.mc_ttl','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.mc_ttl' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.cluster.rsport','records.config','INT 8088') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.cluster.rsport' and config_file = 'records.config' and value = 'INT 8088') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.config_dir','records.config','STRING etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.config_dir' and config_file = 'records.config' and value = 'STRING etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.core_limit','records.config','INT -1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.core_limit' and config_file = 'records.config' and value = 'INT -1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.diags.debug.enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.diags.debug.enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.diags.debug.tags','records.config','STRING http.*|dns.*') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.diags.debug.tags' and config_file = 'records.config' and value = 'STRING http.*|dns.*') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.diags.show_location','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.diags.show_location' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.lookup_timeout','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.lookup_timeout' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.max_dns_in_flight','records.config','INT 2048') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.max_dns_in_flight' and config_file = 'records.config' and value = 'INT 2048') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.nameservers','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.nameservers' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.resolv_conf','records.config','STRING /etc/resolv.conf') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.resolv_conf' and config_file = 'records.config' and value = 'STRING /etc/resolv.conf') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.round_robin_nameservers','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.round_robin_nameservers' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.search_default_domains','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.search_default_domains' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.splitDNS.enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.splitDNS.enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.url_expansions','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.url_expansions' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dns.validate_query_name','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dns.validate_query_name' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.dump_mem_info_frequency','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.dump_mem_info_frequency' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.env_prep','records.config','STRING example_prep.sh') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.env_prep' and config_file = 'records.config' and value = 'STRING example_prep.sh') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.exec_thread.affinity','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.affinity' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.exec_thread.autoconfig','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.autoconfig' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.exec_thread.autoconfig.scale','records.config','FLOAT 1.5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.autoconfig.scale' and config_file = 'records.config' and value = 'FLOAT 1.5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.exec_thread.limit','records.config','INT 32') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.exec_thread.limit' and config_file = 'records.config' and value = 'INT 32') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.header.parse.no_host_url_redirect','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.header.parse.no_host_url_redirect' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.serve_stale_for','records.config','INT 6') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.serve_stale_for' and config_file = 'records.config' and value = 'INT 6') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.size','records.config','INT 120000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.size' and config_file = 'records.config' and value = 'INT 120000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.storage_size','records.config','INT 33554432') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.storage_size' and config_file = 'records.config' and value = 'INT 33554432') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.strict_round_robin','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.strict_round_robin' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.timeout','records.config','INT 1440') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.timeout' and config_file = 'records.config' and value = 'INT 1440') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.hostdb.ttl_mode','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.hostdb.ttl_mode' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.accept_no_activity_timeout','records.config','INT 120') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.accept_no_activity_timeout' and config_file = 'records.config' and value = 'INT 120') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_insert_client_ip','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_insert_client_ip' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_other_header_list','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_other_header_list' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_remove_client_ip','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_client_ip' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_remove_cookie','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_cookie' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_remove_from','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_from' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_remove_referer','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_referer' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.anonymize_remove_user_agent','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.anonymize_remove_user_agent' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.background_fill_active_timeout','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.background_fill_active_timeout' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.background_fill_completed_threshold','records.config','FLOAT 0.0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.background_fill_completed_threshold' and config_file = 'records.config' and value = 'FLOAT 0.0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.allow_empty_doc','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.allow_empty_doc' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.cache_responses_to_cookies','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.cache_responses_to_cookies' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.cache_urls_that_look_dynamic' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.enable_default_vary_headers','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.enable_default_vary_headers' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.fuzz.probability','records.config','FLOAT 0.005') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.fuzz.probability' and config_file = 'records.config' and value = 'FLOAT 0.005') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.fuzz.time','records.config','INT 240') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.fuzz.time' and config_file = 'records.config' and value = 'INT 240') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.guaranteed_max_lifetime','records.config','INT 2592000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.guaranteed_max_lifetime' and config_file = 'records.config' and value = 'INT 2592000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.heuristic_lm_factor','records.config','FLOAT 0.10') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_lm_factor' and config_file = 'records.config' and value = 'FLOAT 0.10') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.heuristic_max_lifetime','records.config','INT 86400') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_max_lifetime' and config_file = 'records.config' and value = 'INT 86400') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.heuristic_min_lifetime','records.config','INT 3600') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.heuristic_min_lifetime' and config_file = 'records.config' and value = 'INT 3600') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.http','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.http' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ignore_accept_encoding_mismatch','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_accept_encoding_mismatch' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ignore_authentication','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_authentication' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ignore_client_cc_max_age','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_client_cc_max_age' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ignore_client_no_cache','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_client_no_cache' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ignore_server_no_cache','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ignore_server_no_cache' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.ims_on_client_no_cache','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.ims_on_client_no_cache' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.max_stale_age','records.config','INT 604800') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.max_stale_age' and config_file = 'records.config' and value = 'INT 604800') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.range.lookup','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.range.lookup' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.required_headers','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.required_headers' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.vary_default_images','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_images' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.vary_default_other','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_other' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.vary_default_text','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.vary_default_text' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests','records.config','INT -1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.when_to_add_no_cache_to_msie_requests' and config_file = 'records.config' and value = 'INT -1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.cache.when_to_revalidate','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.cache.when_to_revalidate' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.chunking_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.chunking_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.congestion_control.enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.congestion_control.enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.connect_attempts_max_retries','records.config','INT 6') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_max_retries' and config_file = 'records.config' and value = 'INT 6') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.connect_attempts_max_retries_dead_server','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_max_retries_dead_server' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.connect_attempts_rr_retries','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_rr_retries' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.connect_attempts_timeout','records.config','INT 10') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 10') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.connect_ports','records.config','STRING 443 563') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.connect_ports' and config_file = 'records.config' and value = 'STRING 443 563') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.down_server.abort_threshold','records.config','INT 10') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.down_server.abort_threshold' and config_file = 'records.config' and value = 'INT 10') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.down_server.cache_time','records.config','INT 300') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.down_server.cache_time' and config_file = 'records.config' and value = 'INT 300') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.enable_http_stats','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.enable_http_stats' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.enable_url_expandomatic','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.enable_url_expandomatic' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.forward.proxy_auth_to_parent','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.forward.proxy_auth_to_parent' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.insert_age_in_response','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_age_in_response' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.insert_request_via_str','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_request_via_str' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.insert_response_via_str','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_response_via_str' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.insert_squid_x_forwarded_for','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.insert_squid_x_forwarded_for' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.keep_alive_enabled_in','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_enabled_in' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.keep_alive_enabled_out','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_enabled_out' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.keep_alive_no_activity_timeout_in','records.config','INT 115') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_in' and config_file = 'records.config' and value = 'INT 115') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.keep_alive_no_activity_timeout_out','records.config','INT 120') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.keep_alive_no_activity_timeout_out' and config_file = 'records.config' and value = 'INT 120') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.negative_caching_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.negative_caching_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.negative_caching_lifetime','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.negative_caching_lifetime' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.no_dns_just_forward_to_parent','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.no_dns_just_forward_to_parent' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.normalize_ae_gzip','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.normalize_ae_gzip' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.origin_server_pipeline','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.origin_server_pipeline' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.fail_threshold','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.fail_threshold' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.file','records.config','STRING parent.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.file' and config_file = 'records.config' and value = 'STRING parent.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.per_parent_connect_attempts' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.retry_time','records.config','INT 60') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.retry_time' and config_file = 'records.config' and value = 'INT 60') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy_routing_enable','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy_routing_enable' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.parent_proxy.total_connect_attempts','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.parent_proxy.total_connect_attempts' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.post_connect_attempts_timeout','records.config','INT 1800') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.post_connect_attempts_timeout' and config_file = 'records.config' and value = 'INT 1800') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.push_method_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.push_method_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.referer_default_redirect','records.config','STRING http://www.example.com/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_default_redirect' and config_file = 'records.config' and value = 'STRING http://www.example.com/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.referer_filter','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_filter' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.referer_format_redirect','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.referer_format_redirect' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.response_server_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.response_server_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.send_http11_requests','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.send_http11_requests' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.server_ports','records.config','STRING 80 80:ipv6') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.server_ports' and config_file = 'records.config' and value = 'STRING 80 80:ipv6') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.share_server_sessions','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.share_server_sessions' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.slow.log.threshold','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.slow.log.threshold' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.transaction_active_timeout_in','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_active_timeout_in' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.transaction_active_timeout_out','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_active_timeout_out' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.transaction_no_activity_timeout_in','records.config','INT 30') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_no_activity_timeout_in' and config_file = 'records.config' and value = 'INT 30') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.transaction_no_activity_timeout_out','records.config','INT 30') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.transaction_no_activity_timeout_out' and config_file = 'records.config' and value = 'INT 30') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.uncacheable_requests_bypass_parent','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.uncacheable_requests_bypass_parent' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.http.user_agent_pipeline','records.config','INT 8') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.http.user_agent_pipeline' and config_file = 'records.config' and value = 'INT 8') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.icp.enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.icp.enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.icp.icp_interface','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.icp.icp_interface' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.icp.icp_port','records.config','INT 3130') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.icp.icp_port' and config_file = 'records.config' and value = 'INT 3130') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.icp.multicast_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.icp.multicast_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.icp.query_timeout','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.icp.query_timeout' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.auto_delete_rolled_files','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.auto_delete_rolled_files' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.collation_host','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_host' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.collation_host_tagged','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_host_tagged' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.collation_port','records.config','INT 8085') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_port' and config_file = 'records.config' and value = 'INT 8085') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.collation_retry_sec','records.config','INT 5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_retry_sec' and config_file = 'records.config' and value = 'INT 5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.collation_secret','records.config','STRING foobar') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.collation_secret' and config_file = 'records.config' and value = 'STRING foobar') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.common_log_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.common_log_header','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_header' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.common_log_is_ascii','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_is_ascii' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.common_log_name','records.config','STRING common') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.common_log_name' and config_file = 'records.config' and value = 'STRING common') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.custom_logs_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.custom_logs_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended2_log_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended2_log_header','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_header' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended2_log_is_ascii','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_is_ascii' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended2_log_name','records.config','STRING extended2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended2_log_name' and config_file = 'records.config' and value = 'STRING extended2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended_log_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended_log_header','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_header' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended_log_is_ascii','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_is_ascii' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.extended_log_name','records.config','STRING extended') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.extended_log_name' and config_file = 'records.config' and value = 'STRING extended') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.hostname','records.config','STRING localhost') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.hostname' and config_file = 'records.config' and value = 'STRING localhost') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.logfile_dir','records.config','STRING var/log/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.logfile_dir' and config_file = 'records.config' and value = 'STRING var/log/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.logfile_perm','records.config','STRING rw-r--r--') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.logfile_perm' and config_file = 'records.config' and value = 'STRING rw-r--r--') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.logging_enabled','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.logging_enabled' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.max_secs_per_buffer','records.config','INT 5') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.max_secs_per_buffer' and config_file = 'records.config' and value = 'INT 5') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.max_space_mb_for_logs','records.config','INT 25000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_for_logs' and config_file = 'records.config' and value = 'INT 25000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.max_space_mb_for_orphan_logs','records.config','INT 25') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_for_orphan_logs' and config_file = 'records.config' and value = 'INT 25') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.max_space_mb_headroom','records.config','INT 1000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.max_space_mb_headroom' and config_file = 'records.config' and value = 'INT 1000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.rolling_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.rolling_interval_sec','records.config','INT 86400') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_interval_sec' and config_file = 'records.config' and value = 'INT 86400') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.rolling_offset_hr','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_offset_hr' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.rolling_size_mb','records.config','INT 10') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.rolling_size_mb' and config_file = 'records.config' and value = 'INT 10') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.sampling_frequency','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.sampling_frequency' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.separate_host_logs','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.separate_host_logs' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.separate_icp_logs','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.separate_icp_logs' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.squid_log_enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.squid_log_header','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_header' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.squid_log_is_ascii','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_is_ascii' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.squid_log_name','records.config','STRING squid') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.squid_log_name' and config_file = 'records.config' and value = 'STRING squid') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.log.xml_config_file','records.config','STRING logs_xml.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.log.xml_config_file' and config_file = 'records.config' and value = 'STRING logs_xml.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.mlock_enabled','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.mlock_enabled' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.connections_throttle','records.config','INT 500000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.connections_throttle' and config_file = 'records.config' and value = 'INT 500000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.default_inactivity_timeout','records.config','INT 180') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.default_inactivity_timeout' and config_file = 'records.config' and value = 'INT 180') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.defer_accept','records.config','INT 45') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.defer_accept' and config_file = 'records.config' and value = 'INT 45') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.sock_recv_buffer_size_in','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_recv_buffer_size_in' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.sock_recv_buffer_size_out','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_recv_buffer_size_out' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.sock_send_buffer_size_in','records.config','INT 262144') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_send_buffer_size_in' and config_file = 'records.config' and value = 'INT 262144') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.net.sock_send_buffer_size_out','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.net.sock_send_buffer_size_out' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.output.logfile','records.config','STRING traffic.out') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.output.logfile' and config_file = 'records.config' and value = 'STRING traffic.out') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.process_manager.mgmt_port','records.config','INT 8084') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.process_manager.mgmt_port' and config_file = 'records.config' and value = 'INT 8084') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.proxy_binary_opts','records.config','STRING -M') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.proxy_binary_opts' and config_file = 'records.config' and value = 'STRING -M') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.proxy_name','records.config','STRING __HOSTNAME__') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.proxy_name' and config_file = 'records.config' and value = 'STRING __HOSTNAME__') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.reverse_proxy.enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.reverse_proxy.enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.snapshot_dir','records.config','STRING snapshots') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.snapshot_dir' and config_file = 'records.config' and value = 'STRING snapshots') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.CA.cert.filename','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.CA.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.CA.cert.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.CA.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.CA.cert.filename','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.CA.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.CA.cert.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.CA.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.cert.filename','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.cert.filename' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.certification_level','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.certification_level' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.cert.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.private_key.filename','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.private_key.filename' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.private_key.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.private_key.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.client.verify.server','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.client.verify.server' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.compression','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.compression' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.number.threads','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.number.threads' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.cert_chain.filename','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cert_chain.filename' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.cert.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cert.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.cipher_suite','records.config','STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.cipher_suite' and config_file = 'records.config' and value = 'STRING ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384:AES128-GCM-SHA256:AES256-GCM-SHA384:ECDHE-RSA-RC4-SHA:ECDHE-RSA-AES128-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA:RC4-MD5:AES128-SHA:AES256-SHA:DES-CBC3-SHA!SRP:!DSS:!PSK:!aNULL:!eNULL:!SSLv2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.honor_cipher_order','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.honor_cipher_order' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.multicert.filename','records.config','STRING ssl_multicert.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.multicert.filename' and config_file = 'records.config' and value = 'STRING ssl_multicert.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.server.private_key.path','records.config','STRING etc/trafficserver/ssl') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.server.private_key.path' and config_file = 'records.config' and value = 'STRING etc/trafficserver/ssl') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.SSLv2','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.SSLv2' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.SSLv3','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.SSLv3' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.ssl.TLSv1','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.ssl.TLSv1' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.stack_dump_enabled','records.config','INT 1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.stack_dump_enabled' and config_file = 'records.config' and value = 'INT 1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.syslog_facility','records.config','STRING LOG_DAEMON') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.syslog_facility' and config_file = 'records.config' and value = 'STRING LOG_DAEMON') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.system.mmap_max','records.config','INT 2097152') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.system.mmap_max' and config_file = 'records.config' and value = 'INT 2097152') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.task_threads','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.task_threads' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.temp_dir','records.config','STRING /tmp') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.temp_dir' and config_file = 'records.config' and value = 'STRING /tmp') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.update.concurrent_updates','records.config','INT 100') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.update.concurrent_updates' and config_file = 'records.config' and value = 'INT 100') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.update.enabled','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.update.enabled' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.update.force','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.update.force' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.update.retry_count','records.config','INT 10') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.update.retry_count' and config_file = 'records.config' and value = 'INT 10') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.update.retry_interval','records.config','INT 2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.update.retry_interval' and config_file = 'records.config' and value = 'INT 2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.url_remap.default_to_server_pac','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.default_to_server_pac' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.url_remap.default_to_server_pac_port','records.config','INT -1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.default_to_server_pac_port' and config_file = 'records.config' and value = 'INT -1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.url_remap.filename','records.config','STRING remap.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.filename' and config_file = 'records.config' and value = 'STRING remap.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.url_remap.pristine_host_hdr','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.pristine_host_hdr' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('CONFIG proxy.config.url_remap.remap_required','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'CONFIG proxy.config.url_remap.remap_required' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('Disk_Volume','storage.config','1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'Disk_Volume' and config_file = 'storage.config' and value = '1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('Drive_Letters','storage.config','b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'Drive_Letters' and config_file = 'storage.config' and value = 'b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('Drive_Prefix','storage.config','/dev/sd') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'Drive_Prefix' and config_file = 'storage.config' and value = '/dev/sd') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.connection.timeout','rascal.properties','2000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'health.connection.timeout' and config_file = 'rascal.properties' and value = '2000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.polling.url','rascal.properties','http://${hostname}/_astats?application=&inf.name=${interface_name}') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'health.polling.url' and config_file = 'rascal.properties' and value = 'http://${hostname}/_astats?application=&inf.name=${interface_name}') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.threshold.availableBandwidthInKbps','rascal.properties','>1750000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'health.threshold.availableBandwidthInKbps' and config_file = 'rascal.properties' and value = '>1750000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.threshold.loadavg','rascal.properties','25.0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'health.threshold.loadavg' and config_file = 'rascal.properties' and value = '25.0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('health.threshold.queryTime','rascal.properties','1000') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'health.threshold.queryTime' and config_file = 'rascal.properties' and value = '1000') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('history.count','rascal.properties','30') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'history.count' and config_file = 'rascal.properties' and value = '30') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LOCAL proxy.config.cache.interim.storage','records.config','STRING NULL') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LOCAL proxy.config.cache.interim.storage' and config_file = 'records.config' and value = 'STRING NULL') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LOCAL proxy.local.cluster.type','records.config','INT 3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LOCAL proxy.local.cluster.type' and config_file = 'records.config' and value = 'INT 3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LOCAL proxy.local.log.collation_mode','records.config','INT 0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LOCAL proxy.local.log.collation_mode' and config_file = 'records.config' and value = 'INT 0') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','remap.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'remap.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','regex_revalidate.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'regex_revalidate.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','astats.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'astats.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','12M_facts','/opt/ort') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = '12M_facts' and value = '/opt/ort') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','logs_xml.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'logs_xml.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_qstring.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_qstring.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_cedexis-perf-tune.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_cedexis-perf-tune.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_cedexis-perf-ext.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_cedexis-perf-ext.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_cim-sfv.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_cim-sfv.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_steam-dns-ext.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_steam-dns-ext.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_steam-dns.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_steam-dns.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cacheurl_yahoo-nfl-dns.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cacheurl_yahoo-nfl-dns.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','ip_allow.config','/opt/trafficserver/etc/trafficserver') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'ip_allow.config' and value = '/opt/trafficserver/etc/trafficserver') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','50-ats.rules','/etc/udev/rules.d/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = '50-ats.rules' and value = '/etc/udev/rules.d/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','volume.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'volume.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','storage.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'storage.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','records.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'records.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','plugin.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'plugin.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','parent.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'parent.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','hosting.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'hosting.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('location','cache.config','/opt/trafficserver/etc/trafficserver/') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'location' and config_file = 'cache.config' and value = '/opt/trafficserver/etc/trafficserver/') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogFormat.Format','logs_xml.config','%<cqtq> chi=%<chi> phn=%<phn> php=%<php> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> pqsn=%<pqsn> uas="%<{User-Agent}cqh>" ') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogFormat.Format' and config_file = 'logs_xml.config' and value = '%<cqtq> chi=%<chi> phn=%<phn> php=%<php> shn=%<shn> url=%<cquuc> cqhm=%<cqhm> cqhv=%<cqhv> pssc=%<pssc> ttms=%<ttms> b=%<pscl> sssc=%<sssc> sscl=%<sscl> cfsc=%<cfsc> pfsc=%<pfsc> crc=%<crc> phr=%<phr> pqsn=%<pqsn> uas="%<{User-Agent}cqh>" xmt="%<{X-MoneyTrace}cqh>" ') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogFormat.Name','logs_xml.config','custom_ats_2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogFormat.Name' and config_file = 'logs_xml.config' and value = 'custom_ats_2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.Filename','logs_xml.config','custom_ats_2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.Filename' and config_file = 'logs_xml.config' and value = 'custom_ats_2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.Format','logs_xml.config','custom_ats_2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.Format' and config_file = 'logs_xml.config' and value = 'custom_ats_2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.RollingEnabled','logs_xml.config','3') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.RollingEnabled' and config_file = 'logs_xml.config' and value = '3') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.RollingIntervalSec','logs_xml.config','86400') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.RollingIntervalSec' and config_file = 'logs_xml.config' and value = '86400') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.RollingOffsetHr','logs_xml.config','11') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.RollingOffsetHr' and config_file = 'logs_xml.config' and value = '11') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('LogObject.RollingSizeMb','logs_xml.config','1024') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'LogObject.RollingSizeMb' and config_file = 'logs_xml.config' and value = '1024') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('path','astats.config','_astats') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'path' and config_file = 'astats.config' and value = '_astats') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('RAM_Drive_Letters','storage.config','0,1,2,3,4,5,6,7') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'RAM_Drive_Letters' and config_file = 'storage.config' and value = '0,1,2,3,4,5,6,7') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('RAM_Drive_Prefix','storage.config','/dev/ram') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'RAM_Drive_Prefix' and config_file = 'storage.config' and value = '/dev/ram') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('RAM_Volume','storage.config','2') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'RAM_Volume' and config_file = 'storage.config' and value = '2') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('record_types','astats.config','122') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'record_types' and config_file = 'astats.config' and value = '122') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('regex_revalidate.so','plugin.config','--config regex_revalidate.config') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'regex_revalidate.so' and config_file = 'plugin.config' and value = '--config regex_revalidate.config') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('trafficserver','chkconfig','0:off	1:off	2:on	3:on	4:on	5:on	6:off') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'trafficserver' and config_file = 'chkconfig' and value = '0:off	1:off	2:on	3:on	4:on	5:on	6:off') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('trafficserver','package','5.3.2-744.31aba39.el6.x86_64') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'trafficserver' and config_file = 'package' and value = '5.3.2-744.31aba39.el6.x86_64') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('weight','parent.config','0.1') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'weight' and config_file = 'parent.config' and value = '0.1') )  ON CONFLICT (profile, parameter) DO NOTHING;

INSERT INTO parameter (name, config_file, value) VALUES ('weight','CRConfig.json','1.0') ON CONFLICT (name,config_file, value) DO NOTHING;
INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = 'MID_PROFILE'), (select id from parameter where name = 'weight' and config_file = 'CRConfig.json' and value = '1.0') )  ON CONFLICT (profile, parameter) DO NOTHING;

