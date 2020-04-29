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

CREATE TABLE IF NOT EXISTS deleted_api_capability (
    id bigserial PRIMARY KEY,
    http_method http_method_t NOT NULL,
    route text NOT NULL,
    capability text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_asn (
    id bigint NOT NULL,
    asn bigint NOT NULL,
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_cachegroup (
    id bigint,
    name text NOT NULL,
    short_name text NOT NULL,
    parent_cachegroup_id bigint,
    secondary_parent_cachegroup_id bigint,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    fallback_to_closest boolean DEFAULT TRUE,
    coordinate bigint
);

CREATE TABLE IF NOT EXISTS deleted_cachegroup_fallbacks (
    primary_cg bigint NOT NULL,
    backup_cg bigint NOT NULL CHECK (primary_cg != backup_cg),
    set_order bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_cachegroup_localization_method (
    cachegroup bigint NOT NULL,
    method localization_method NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_cachegroup_parameter (
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_capability (
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_cdn (
    id bigint,
    name text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    dnssec_enabled boolean DEFAULT false NOT NULL,
    domain_name text NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_coordinate (
    id bigserial,
    name text NOT NULL,
    latitude numeric NOT NULL DEFAULT 0.0,
    longitude numeric NOT NULL DEFAULT 0.0,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice (
    id bigint,
    xml_id text NOT NULL,
    active boolean DEFAULT false NOT NULL,
    dscp bigint NOT NULL,
    signing_algorithm deliveryservice_signature_type,
    qstring_ignore smallint,
    geo_limit smallint DEFAULT '0'::smallint,
    http_bypass_fqdn text,
    dns_bypass_ip text,
    dns_bypass_ip6 text,
    dns_bypass_ttl bigint,
    type bigint NOT NULL,
    profile bigint,
    cdn_id bigint NOT NULL,
    ccr_dns_ttl bigint,
    global_max_mbps bigint,
    global_max_tps bigint,
    long_desc text,
    long_desc_1 text,
    long_desc_2 text,
    max_dns_answers bigint DEFAULT '5'::bigint,
    info_url text,
    miss_lat numeric,
    miss_long numeric,
    check_path text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    protocol smallint DEFAULT '0'::smallint,
    ssl_key_version bigint DEFAULT '0'::bigint,
    ipv6_routing_enabled boolean DEFAULT false,
    range_request_handling smallint DEFAULT '0'::smallint,
    edge_header_rewrite text,
    origin_shield text,
    mid_header_rewrite text,
    regex_remap text,
    cacheurl text,
    remap_text text,
    multi_site_origin boolean DEFAULT false,
    display_name text NOT NULL,
    tr_response_headers text,
    initial_dispersion bigint DEFAULT '1'::bigint,
    dns_bypass_cname text,
    tr_request_headers text,
    regional_geo_blocking boolean DEFAULT false NOT NULL,
    geo_provider smallint DEFAULT '0'::smallint,
    geo_limit_countries text,
    logs_enabled boolean DEFAULT false,
    multi_site_origin_algorithm smallint,
    geolimit_redirect_url text,
    tenant_id bigint NOT NULL,
    routing_name text NOT NULL DEFAULT 'cdn',
    deep_caching_type deep_caching_type NOT NULL DEFAULT 'NEVER',
    fq_pacing_rate bigint DEFAULT 0,
    anonymous_blocking_enabled boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice_regex (
    deliveryservice bigint NOT NULL,
    regex bigint NOT NULL,
    set_number bigint DEFAULT '0'::bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice_request (
    assignee_id bigint,
    author_id bigint NOT NULL,
    change_type change_types NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    id bigserial,
    last_edited_by_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    deliveryservice jsonb NOT NULL,
    status workflow_states NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice_request_comment (
    author_id bigint NOT NULL,
    deliveryservice_request_id bigint NOT NULL,
    id bigserial,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    value text NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice_server (
    deliveryservice bigint NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservice_tmuser (
    deliveryservice bigint NOT NULL,
    tm_user_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


CREATE TABLE IF NOT EXISTS deleted_division (
    id bigint NOT NULL,
    name text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_federation (
    id bigint NOT NULL,
    cname text NOT NULL,
    description text,
    ttl integer NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_federation_deliveryservice (
    federation bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_federation_federation_resolver (
    federation bigint NOT NULL,
    federation_resolver bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_federation_resolver (
    id bigint NOT NULL,
    ip_address text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_federation_tmuser (
    federation bigint NOT NULL,
    tm_user bigint NOT NULL,
    role bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_hwinfo (
    id bigint NOT NULL,
    serverid bigint NOT NULL,
    description text NOT NULL,
    val text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_job (
    id bigint NOT NULL,
    agent bigint,
    object_type text,
    object_name text,
    keyword text NOT NULL,
    parameters text,
    asset_url text NOT NULL,
    asset_type text NOT NULL,
    status bigint NOT NULL,
    start_time timestamp with time zone NOT NULL,
    entered_time timestamp with time zone NOT NULL,
    job_user bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    job_deliveryservice bigint
);

CREATE TABLE IF NOT EXISTS deleted_job_agent (
    id bigint NOT NULL,
    name text,
    description text,
    active integer DEFAULT 0 NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_job_status (
    id bigint NOT NULL,
    name text,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_log (
    id bigint NOT NULL,
    level text,
    message text NOT NULL,
    tm_user bigint NOT NULL,
    ticketnum text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_origin (
    id bigserial NOT NULL,
    name text NOT NULL,
    fqdn text NOT NULL,
    protocol origin_protocol NOT NULL DEFAULT 'http',
    is_primary boolean NOT NULL DEFAULT FALSE,
    port bigint, -- TODO: port numbers have a max of 65535 - this could be just an integer
    ip_address text, -- TODO: these should be inet type, not text
    ip6_address text,
    deliveryservice bigint NOT NULL,
    coordinate bigint,
    profile bigint,
    cachegroup bigint,
    tenant bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_parameter (
    id bigint NOT NULL,
    name text NOT NULL,
    config_file text,
    value text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    secure boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_phys_location (
    id bigint NOT NULL,
    name text NOT NULL,
    short_name text NOT NULL,
    address text NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    zip text NOT NULL,
    poc text,
    phone text,
    email text,
    comments text,
    region bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_profile (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    type profile_type NOT NULL,
    cdn bigint NOT NULL,
    routing_disabled boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS deleted_profile_parameter (
    profile bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_regex (
    id bigint NOT NULL,
    pattern text DEFAULT ''::text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_region (
    id bigint NOT NULL,
    name text NOT NULL,
    division bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_role (
    id bigint,
    name text NOT NULL,
    description text,
    priv_level bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_role_capability (
    role_id bigint NOT NULL,
    cap_name text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_server (
    id bigint NOT NULL,
    host_name text NOT NULL,
    domain_name text NOT NULL,
    tcp_port bigint,
    xmpp_id text,
    xmpp_passwd text,
    interface_name text NOT NULL,
    ip_address text NOT NULL,
    ip_netmask text NOT NULL,
    ip_gateway text NOT NULL,
    ip6_address text,
    ip6_gateway text,
    interface_mtu bigint DEFAULT '9000'::bigint NOT NULL,
    phys_location bigint NOT NULL,
    rack text,
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    type bigint NOT NULL,
    status bigint NOT NULL,
    offline_reason text,
    upd_pending boolean DEFAULT false NOT NULL,
    profile bigint NOT NULL,
    cdn_id bigint NOT NULL,
    mgmt_ip_address text,
    mgmt_ip_netmask text,
    mgmt_ip_gateway text,
    ilo_ip_address text,
    ilo_ip_netmask text,
    ilo_ip_gateway text,
    ilo_username text,
    ilo_password text,
    router_host_name text,
    router_port_name text,
    guid text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    https_port bigint,
    reval_pending boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS deleted_servercheck (
    id bigint NOT NULL,
    server bigint NOT NULL,
    aa bigint,
    ab bigint,
    ac bigint,
    ad bigint,
    ae bigint,
    af bigint,
    ag bigint,
    ah bigint,
    ai bigint,
    aj bigint,
    ak bigint,
    al bigint,
    am bigint,
    an bigint,
    ao bigint,
    ap bigint,
    aq bigint,
    ar bigint,
    bf bigint,
    at bigint,
    au bigint,
    av bigint,
    aw bigint,
    ax bigint,
    ay bigint,
    az bigint,
    ba bigint,
    bb bigint,
    bc bigint,
    bd bigint,
    be bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_snapshot (
    cdn text NOT NULL,
    content json NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_staticdnsentry (
    id bigint NOT NULL,
    host text NOT NULL,
    address text NOT NULL,
    type bigint NOT NULL,
    ttl bigint DEFAULT '3600'::bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    cachegroup bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_stats_summary (
    id bigint NOT NULL,
    cdn_name text DEFAULT 'all'::text NOT NULL,
    deliveryservice_name text NOT NULL,
    stat_name text NOT NULL,
    stat_value double precision NOT NULL,
    summary_time timestamp with time zone DEFAULT now() NOT NULL,
    stat_date date
);


CREATE TABLE IF NOT EXISTS deleted_status (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_steering_target (
    deliveryservice bigint NOT NULL,
    target bigint NOT NULL,
    value bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    type bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_tenant (
    id bigserial,
    name text NOT NULL,
    active boolean NOT NULL DEFAULT FALSE,
    parent_id bigint DEFAULT 1 CHECK (id != parent_id),
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_tm_user (
    id bigint NOT NULL,
    username text,
    public_ssh_key text,
    role bigint,
    uid bigint,
    gid bigint,
    local_passwd text,
    confirm_local_passwd text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    company text,
    email text,
    full_name text,
    new_user boolean DEFAULT false NOT NULL,
    address_line1 text,
    address_line2 text,
    city text,
    state_or_province text,
    phone_number text,
    postal_code text,
    country text,
    token text,
    registration_sent timestamp with time zone,
    tenant_id bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_to_extension (
    id bigint NOT NULL,
    name text NOT NULL,
    version text NOT NULL,
    info_url text NOT NULL,
    script_file text NOT NULL,
    isactive boolean DEFAULT false NOT NULL,
    additional_config_json text,
    description text,
    servercheck_short_name text,
    servercheck_column_name text,
    type bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_type (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    use_in_table text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_user_role (
    user_id bigint NOT NULL,
    role_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deleted_server_capability (
    name TEXT,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_server_server_capability (
    server_capability TEXT NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_deliveryservices_required_capability (
    required_capability TEXT NOT NULL,
    deliveryservice_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,

    PRIMARY KEY (deliveryservice_id, required_capability)
);
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS deleted_api_capability;
DROP TABLE IF EXISTS deleted_asn;
DROP TABLE IF EXISTS deleted_cachegroup;
DROP TABLE IF EXISTS deleted_cachegroup_fallbacks;
DROP TABLE IF EXISTS deleted_cachegroup_localization_method;
DROP TABLE IF EXISTS deleted_cachegroup_parameter;
DROP TABLE IF EXISTS deleted_capability;
DROP TABLE IF EXISTS deleted_cdn;
DROP TABLE IF EXISTS deleted_coordinate;
DROP TABLE IF EXISTS deleted_deliveryservice;
DROP TABLE IF EXISTS deleted_deliveryservice_regex;
DROP TABLE IF EXISTS deleted_deliveryservice_request;
DROP TABLE IF EXISTS deleted_deliveryservice_request_comment;
DROP TABLE IF EXISTS deleted_deliveryservice_server;
DROP TABLE IF EXISTS deleted_deliveryservice_tmuser;
DROP TABLE IF EXISTS deleted_division;
DROP TABLE IF EXISTS deleted_federation;
DROP TABLE IF EXISTS deleted_federation_deliveryservice;
DROP TABLE IF EXISTS deleted_federation_federation_resolver;
DROP TABLE IF EXISTS deleted_federation_resolver;
DROP TABLE IF EXISTS deleted_federation_tmuser;
DROP TABLE IF EXISTS deleted_hwinfo;
DROP TABLE IF EXISTS deleted_job;
DROP TABLE IF EXISTS deleted_job_agent;
DROP TABLE IF EXISTS deleted_job_status;
DROP TABLE IF EXISTS deleted_log;
DROP TABLE IF EXISTS deleted_origin;
DROP TABLE IF EXISTS deleted_parameter;
DROP TABLE IF EXISTS deleted_phys_location;
DROP TABLE IF EXISTS deleted_profile;
DROP TABLE IF EXISTS deleted_profile_parameter;
DROP TABLE IF EXISTS deleted_regex;
DROP TABLE IF EXISTS deleted_region;
DROP TABLE IF EXISTS deleted_role;
DROP TABLE IF EXISTS deleted_role_capability;
DROP TABLE IF EXISTS deleted_server;
DROP TABLE IF EXISTS deleted_servercheck;
DROP TABLE IF EXISTS deleted_snapshot;
DROP TABLE IF EXISTS deleted_staticdnsentry;
DROP TABLE IF EXISTS deleted_stats_summary;
DROP TABLE IF EXISTS deleted_status;
DROP TABLE IF EXISTS deleted_steering_target;
DROP TABLE IF EXISTS deleted_tenant;
DROP TABLE IF EXISTS deleted_tm_user;
DROP TABLE IF EXISTS deleted_to_extension;
DROP TABLE IF EXISTS deleted_type;
DROP TABLE IF EXISTS deleted_user_role;
DROP TABLE IF EXISTS deleted_server_capability;
DROP TABLE IF EXISTS deleted_server_server_capability;
DROP TABLE IF EXISTS deleted_deliveryservices_required_capability;