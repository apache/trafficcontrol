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
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.1
-- Dumped by pg_dump version 9.5.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: traffic_ops; Type: DATABASE; Schema: -; Owner: touser
--

CREATE DATABASE traffic_ops WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


ALTER DATABASE traffic_ops OWNER TO touser;

\connect traffic_ops_debug

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: asns; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE asns (
    asn integer NOT NULL,
    cachegroup text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE asns OWNER TO touser;

--
-- Name: asns_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW asns AS
 SELECT asns.asn,
    asns.cachegroup,
    asns.created_at,
    pg_xact_commit_timestamp(asns.xmin) AS last_updated
   FROM asns;


ALTER TABLE asns OWNER TO touser;


CREATE TABLE cachegroups_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE cachegroups_types OWNER TO touser;

--
-- Name: cachegroups; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cachegroups (
    name text NOT NULL,
    description text NOT NULL,
    latitude numeric,
    longitude numeric,
    parent_cachegroup text,
    type text NOT NULL REFERENCES cachegroups_types (name),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cachegroups OWNER TO touser;

--
-- Name: cachegroups_parameters; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cachegroups_parameters (
    cachegroup text NOT NULL,
    parameter_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cachegroups_parameters OWNER TO touser;

--
-- Name: cdns; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cdns (
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cdns OWNER TO touser;

--
-- Name: parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE parameter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE parameter_id_seq OWNER TO touser;

--
-- Name: parameters; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE parameters (
    id bigint DEFAULT nextval('parameter_id_seq'::regclass) NOT NULL,
    name text NOT NULL,
    config_file text NOT NULL,
    value text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE parameters OWNER TO touser;

--
-- Name: profiles; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE profiles (
    name text NOT NULL,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE profiles OWNER TO touser;

--
-- Name: profiles_parameters; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE profiles_parameters (
    profile text NOT NULL,
    parameter_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE profiles_parameters OWNER TO touser;

ALTER TABLE ONLY profiles_parameters
    ADD CONSTRAINT profiles_parameters_profile_parameter_id_pkey PRIMARY KEY (profile, parameter_id);

CREATE TABLE servers_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE servers_types OWNER TO touser;

--
-- Name: servers; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE servers (
    host_name text NOT NULL,
    domain_name text NOT NULL,
    tcp_port integer,
    xmpp_id text,
    xmpp_passwd text,
    interface_name text NOT NULL,
    ip_address inet,
    ip_gateway inet,
    ip6_address inet,
    ip6_gateway inet,
    interface_mtu integer DEFAULT 9000 NOT NULL,
    phys_location text NOT NULL,
    rack text,
    cachegroup text NOT NULL,
    type text NOT NULL REFERENCES servers_types (name),
    status text NOT NULL,
    upd_pending boolean DEFAULT false NOT NULL,
    profile text NOT NULL,
    cdn text NOT NULL,
    mgmt_ip_address inet,
    mgmt_ip_gateway inet,
    ilo_ip_address inet,
    ilo_ip_gateway inet,
    ilo_username text,
    ilo_password text,
    router_host_name text,
    router_port_name text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE servers OWNER TO touser;

ALTER TABLE ONLY servers
    ADD CONSTRAINT servers_host_name_tcp_port_pkey PRIMARY KEY (host_name, tcp_port);

ALTER TABLE ONLY servers
    ADD CONSTRAINT servers_profile_profiles_name_fkey FOREIGN KEY (profile) REFERENCES profiles(name);

--
-- Name: content_routers_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW content_routers AS
 SELECT servers.ip_address AS ip,
    servers.ip6_address AS ip6,
    servers.profile,
    servers.cachegroup AS location,
    servers.status,
    servers.tcp_port AS port,
    servers.host_name,
    concat(servers.host_name, '.', servers.domain_name) AS fqdn,
    parameters.value AS apiport,
    servers.cdn
   FROM (((servers
     JOIN profiles ON ((profiles.name = servers.profile)))
     JOIN profiles_parameters ON ((profiles_parameters.profile = profiles.name)))
     JOIN parameters ON ((parameters.id = profiles_parameters.parameter_id)))
  WHERE ((servers.type = 'CCR'::text) AND (parameters.name = 'api.port'::text));


ALTER TABLE content_routers OWNER TO touser;

--
-- Name: content_servers_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW content_servers AS
 SELECT DISTINCT servers.host_name,
    servers.profile,
    servers.type,
    servers.cachegroup AS location,
    servers.ip_address AS ip,
    servers.cdn,
    servers.status,
    servers.cachegroup AS cache_group,
    servers.ip6_address AS ip6,
    servers.tcp_port AS port,
    concat(servers.host_name, '.', servers.domain_name) AS fqdn,
    servers.interface_name,
    parameters.value AS hash_count
   FROM (((servers
     JOIN profiles ON ((profiles.name = servers.profile)))
     JOIN profiles_parameters ON ((profiles_parameters.profile = profiles.name)))
     JOIN parameters ON ((parameters.id = profiles_parameters.parameter_id)))
  WHERE ((parameters.name = 'weight'::text) AND (servers.status = ANY (ARRAY['REPORTED'::text, 'ONLINE'::text])) AND (servers.type = 'EDGE'::text));


ALTER TABLE content_servers OWNER TO touser;


CREATE TABLE deliveryservices_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE deliveryservices_types OWNER TO touser;

--
-- Name: deliveryservices; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservices (
    name text NOT NULL,
    display_name text NOT NULL,
    description text NOT NULL,
    cdn text NOT NULL,
    domain text NOT NULL,
    active boolean NOT NULL,
    dscp smallint NOT NULL,
    signed boolean NOT NULL,
    qstring_ignore boolean NOT NULL,
    geo_limit boolean NOT NULL,
    http_bypass_fqdn text,
    dns_bypass_ip inet,
    dns_bypass_ip6 inet,
    dns_bypass_ttl integer,
    org_server_fqdn text,
    type text NOT NULL REFERENCES deliveryservices_types (name),
    profile text NOT NULL,
    dns_ttl integer,
    global_max_mbps integer,
    global_max_tps integer,
    max_dns_answers integer DEFAULT 0,
    info_url text,
    miss_lat numeric,
    miss_long numeric,
    check_path text,
    protocol smallint DEFAULT 0,
    ssl_key_version bigint DEFAULT 0,
    ipv6_routing_enabled boolean NOT NULL,
    range_request_handling smallint DEFAULT 0,
    edge_header_rewrite text,
    origin_shield text,
    mid_header_rewrite text,
    regex_remap text,
    cacheurl text,
    remap_text text,
    multi_site_origin boolean,
    tr_response_headers text,
    initial_dispersion integer DEFAULT 1 NOT NULL,
    dns_bypass_cname text,
    tr_request_headers text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE deliveryservices OWNER TO touser;

ALTER TABLE ONLY deliveryservices
    ADD CONSTRAINT deliveryservices_name_pkey PRIMARY KEY (name);

--
-- Name: deliveryservices_regexes; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservices_regexes (
    deliveryservice text NOT NULL,
    regex_id bigint NOT NULL,
    set_number integer DEFAULT 0
);


ALTER TABLE deliveryservices_regexes OWNER TO touser;

--
-- Name: deliveryservices_servers; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservices_servers (
    deliveryservice name NOT NULL,
    server text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE deliveryservices_servers OWNER TO touser;

--
-- Name: regexes_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE regexes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE regexes_id_seq OWNER TO touser;


CREATE TABLE regexes_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE regexes_types OWNER TO touser;

--
-- Name: regexes; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE regexes (
    id bigint DEFAULT nextval('regexes_id_seq'::regclass) NOT NULL,
    pattern text NOT NULL,
    type text NOT NULL REFERENCES regexes_types (name),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE regexes OWNER TO touser;

ALTER TABLE ONLY regexes
    ADD CONSTRAINT regexes_id_pkey PRIMARY KEY (id);

--
-- Name: cr_deliveryservice_server_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW cr_deliveryservice_server AS
 SELECT DISTINCT regexes.pattern,
    deliveryservices.name,
    servers.cdn,
    servers.host_name AS server_name
   FROM ((((deliveryservices
     JOIN deliveryservices_regexes ON ((deliveryservices_regexes.deliveryservice = deliveryservices.name)))
     JOIN regexes ON ((regexes.id = deliveryservices_regexes.regex_id)))
     JOIN deliveryservices_servers ON ((deliveryservices.name = (deliveryservices_servers.deliveryservice)::text)))
     JOIN servers ON ((servers.host_name = deliveryservices_servers.server)))
  WHERE (deliveryservices.type <> 'ANY_MAP'::text);


ALTER TABLE cr_deliveryservice_server OWNER TO touser;

--
-- Name: staticdnsentries_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE staticdnsentries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE staticdnsentries_id_seq OWNER TO touser;


CREATE TABLE staticdnsentries_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE staticdnsentries_types OWNER TO touser;

--
-- Name: staticdnsentries; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE staticdnsentries (
    id integer DEFAULT nextval('staticdnsentries_id_seq'::regclass) NOT NULL,
    name character varying(63) NOT NULL,
    type text NOT NULL REFERENCES staticdnsentries_types (name),
    class character varying(2) NOT NULL,
    ttl bigint DEFAULT 3600 NOT NULL,
    rdata character varying(255) NOT NULL,
    deliveryservice text NOT NULL,
    cachegroup text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE staticdnsentries OWNER TO touser;

ALTER TABLE ONLY staticdnsentries
    ADD CONSTRAINT staticdnsentries_id_pkey PRIMARY KEY (id);

--
-- Name: crconfig_ds_data_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW crconfig_ds_data AS
 SELECT deliveryservices.name,
    deliveryservices.profile,
    deliveryservices.dns_ttl,
    deliveryservices.global_max_mbps,
    deliveryservices.global_max_tps,
    deliveryservices.max_dns_answers,
    deliveryservices.miss_lat,
    deliveryservices.miss_long,
    deliveryservices.protocol,
    deliveryservices.ipv6_routing_enabled,
    deliveryservices.tr_request_headers,
    deliveryservices.tr_response_headers,
    deliveryservices.initial_dispersion,
    deliveryservices.dns_bypass_cname,
    deliveryservices.dns_bypass_ip,
    deliveryservices.dns_bypass_ip6,
    deliveryservices.dns_bypass_ttl,
    deliveryservices.geo_limit,
    deliveryservices.cdn,
    regexes.pattern AS match_pattern,
    regextypes.name AS match_type,
    deliveryservices_regexes.set_number,
    staticdnsentries.name AS sdns_host,
    staticdnsentries.rdata AS sdns_address,
    staticdnsentries.ttl AS sdns_ttl,
    sdnstypes.name AS sdns_type
   FROM (((((deliveryservices
     LEFT JOIN staticdnsentries ON ((deliveryservices.name = staticdnsentries.deliveryservice)))
     JOIN deliveryservices_regexes ON ((deliveryservices_regexes.deliveryservice = deliveryservices.name)))
     JOIN regexes ON ((regexes.id = deliveryservices_regexes.regex_id)))
     JOIN regexes_types regextypes ON ((regextypes.name = regexes.type)))
     LEFT JOIN staticdnsentries_types sdnstypes ON ((sdnstypes.name = (staticdnsentries.type)::text)));


ALTER TABLE crconfig_ds_data OWNER TO touser;

--
-- Name: crconfig_params_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW crconfig_params AS
 SELECT DISTINCT servers.cdn,
    servers.profile,
    servers.type AS stype,
    parameters.name AS pname,
    parameters.config_file AS cfile,
    parameters.value AS pvalue
   FROM (((servers
     JOIN profiles ON ((profiles.name = servers.profile)))
     JOIN profiles_parameters ON ((profiles_parameters.profile = servers.profile)))
     JOIN parameters ON ((parameters.id = profiles_parameters.parameter_id)))
  WHERE ((servers.type = ANY (ARRAY['EDGE'::text, 'MID'::text, 'CCR'::text])) AND (parameters.config_file = 'CRConfig.json'::text));


ALTER TABLE crconfig_params OWNER TO touser;

--
-- Name: deliveryservices_users; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservices_users (
    deliveryservice text NOT NULL,
    username text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE deliveryservices_users OWNER TO touser;

--
-- Name: divisions; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE divisions (
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE divisions OWNER TO touser;

--
-- Name: domains; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE domains (
    name text NOT NULL,
    cdn text NOT NULL,
    dnssec boolean NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE domains OWNER TO touser;


CREATE TABLE extensions_types (
    name text PRIMARY KEY,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);
ALTER TABLE extensions_types OWNER TO touser;

--
-- Name: extensions; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE extensions (
    name text NOT NULL,
    short_name text,
    description text,
    version text NOT NULL,
    info_url text NOT NULL,
    script_file text NOT NULL,
    active boolean NOT NULL,
    additional_config_json text,
    type text NOT NULL REFERENCES extensions_types (name),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE extensions OWNER TO touser;

ALTER TABLE ONLY extensions
    ADD CONSTRAINT extensions_name_pkey PRIMARY KEY (name);

--
-- Name: federation_resolvers_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE federation_resolvers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federation_resolvers_id_seq OWNER TO touser;

--
-- Name: federation_resolvers; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_resolvers (
    id bigint DEFAULT nextval('federation_resolvers_id_seq'::regclass) NOT NULL,
    ip_address inet NOT NULL,
    type text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE federation_resolvers OWNER TO touser;

--
-- Name: federation_users; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_users (
    federation_id bigint NOT NULL,
    username text NOT NULL,
    role text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE federation_users OWNER TO touser;

--
-- Name: federations_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE federations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federations_id_seq OWNER TO touser;

--
-- Name: federations; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federations (
    id bigint DEFAULT nextval('federations_id_seq'::regclass) NOT NULL,
    cname text NOT NULL,
    description text,
    ttl integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE federations OWNER TO touser;

--
-- Name: federations_deliveryservices; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federations_deliveryservices (
    federation_id bigint NOT NULL,
    deliveryservice text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE federations_deliveryservices OWNER TO touser;

--
-- Name: federations_federation_resolvers; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federations_federation_resolvers (
    federation_id bigint NOT NULL,
    federation_resolver integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE federations_federation_resolvers OWNER TO touser;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goose_db_version_id_seq OWNER TO touser;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE goose_db_version (
    id integer DEFAULT nextval('goose_db_version_id_seq'::regclass) NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE goose_db_version OWNER TO touser;

--
-- Name: log_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE log_id_seq OWNER TO touser;

--
-- Name: log; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE log (
    id integer DEFAULT nextval('log_id_seq'::regclass) NOT NULL,
    level text,
    message text NOT NULL,
    username text NOT NULL,
    ticketnum text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE log OWNER TO touser;

ALTER TABLE ONLY log
    ADD CONSTRAINT log_id_pkey PRIMARY KEY (id);

--
-- Name: monitors_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW monitors AS
 SELECT servers.ip_address AS ip,
    servers.ip6_address AS ip6,
    servers.profile,
    servers.cachegroup AS location,
    servers.status,
    servers.tcp_port AS port,
    concat(servers.host_name, '.', servers.domain_name) AS fqdn,
    servers.cdn,
    servers.host_name
   FROM servers
  WHERE (servers.type = 'RASCAL'::text);


ALTER TABLE monitors OWNER TO touser;

--
-- Name: phys_locations; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE phys_locations (
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
    region text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE phys_locations OWNER TO touser;

ALTER TABLE ONLY phys_locations
    ADD CONSTRAINT phys_locations_name_pkey PRIMARY KEY (name);

--
-- Name: profiles_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW profiles AS
 SELECT profiles.name,
    profiles.description,
    profiles.created_at,
    pg_xact_commit_timestamp(profiles.xmin) AS last_updated
   FROM profiles;


ALTER TABLE profiles OWNER TO touser;

--
-- Name: regions; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE regions (
    name text NOT NULL,
    division text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE regions OWNER TO touser;

ALTER TABLE ONLY regions
    ADD CONSTRAINT regions_name_pkey PRIMARY KEY (name);

--
-- Name: regions_v; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW regions AS
 SELECT regions.name,
    regions.division,
    regions.created_at,
    pg_xact_commit_timestamp(regions.xmin) AS last_updated
   FROM regions;


ALTER TABLE regions OWNER TO touser;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE roles (
    name text NOT NULL,
    description text,
    priv_level integer NOT NULL
);


ALTER TABLE roles OWNER TO touser;

ALTER TABLE ONLY roles
    ADD CONSTRAINT roles_name_pkey PRIMARY KEY (name);

--
-- Name: stats_summary; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE stats_summary (
    cdn_name text NOT NULL,
    deliveryservice text NOT NULL,
    stat_name text NOT NULL,
    stat_value numeric NOT NULL,
    stat_date date NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE stats_summary OWNER TO touser;

ALTER TABLE ONLY stats_summary
    ADD CONSTRAINT stats_summary_cdn_name_deliveryservice_stat_name_stat_date_pkey PRIMARY KEY (cdn_name, deliveryservice, stat_name, stat_date);

--
-- Name: statuses; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE statuses (
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE statuses OWNER TO touser;

ALTER TABLE ONLY statuses
    ADD CONSTRAINT statuses_name_pkey PRIMARY KEY (name);

--
-- Name: users; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE users (
    username text NOT NULL,
    role text,
    email text,
    full_name text,
    ssh_pub_key text,
    local_password text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE users OWNER TO touser;

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_pkey PRIMARY KEY (username);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_role_roles_name_fkey FOREIGN KEY (role) REFERENCES roles(name);

--
-- Name: asns_asn_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY asns
    ADD CONSTRAINT asns_asn_pkey PRIMARY KEY (asn);


--
-- Name: cachegroup_name_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroups
    ADD CONSTRAINT cachegroup_name_pkey PRIMARY KEY (name);


--
-- Name: cachegroups_parameters_cachegroup_parameter_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroups_parameters
    ADD CONSTRAINT cachegroups_parameters_cachegroup_parameter_id_pkey PRIMARY KEY (cachegroup, parameter_id);


--
-- Name: cdns_name_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cdns
    ADD CONSTRAINT cdns_name_pkey PRIMARY KEY (name);


--
-- Name: deliveryservices_name_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices
    ADD CONSTRAINT deliveryservices_name_pkey PRIMARY KEY (name);


--
-- Name: deliveryservices_regexes_deliveryservice_regex_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices_regexes
    ADD CONSTRAINT deliveryservices_regexes_deliveryservice_regex_id_pkey PRIMARY KEY (deliveryservice, regex_id);


--
-- Name: deliveryservices_servers_deliveryservice_server_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices_servers
    ADD CONSTRAINT deliveryservices_servers_deliveryservice_server_pkey PRIMARY KEY (deliveryservice, server);


--
-- Name: deliveryservices_users_deliveryservice_username_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices_users
    ADD CONSTRAINT deliveryservices_users_deliveryservice_username_pkey PRIMARY KEY (deliveryservice, username);


--
-- Name: divisions_name_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY divisions
    ADD CONSTRAINT divisions_name_pkey PRIMARY KEY (name);


--
-- Name: domains_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY domains
    ADD CONSTRAINT domains_pkey PRIMARY KEY (name);


--
-- Name: federation_resolvers_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_resolvers
    ADD CONSTRAINT federation_resolvers_id_pkey PRIMARY KEY (id);


--
-- Name: federation_users_federation_username_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_users
    ADD CONSTRAINT federation_users_federation_username_pkey PRIMARY KEY (federation_id, username);


--
-- Name: federations_deliveryservices_federation_id_deliveryservice_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federations_deliveryservices
    ADD CONSTRAINT federations_deliveryservices_federation_id_deliveryservice_pkey PRIMARY KEY (federation_id, deliveryservice);


--
-- Name: federations_federation_resolvers_federation_id_federation_resol; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federations_federation_resolvers
    ADD CONSTRAINT federations_federation_resolvers_federation_id_federation_resol PRIMARY KEY (federation_id, federation_resolver);


--
-- Name: federations_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federations
    ADD CONSTRAINT federations_id_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY goose_db_version
    ADD CONSTRAINT goose_db_version_id_pkey PRIMARY KEY (id);


--
-- Name: parameters_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY parameters
    ADD CONSTRAINT parameters_id_pkey PRIMARY KEY (id);


--
-- Name: profiles_name_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY profiles
    ADD CONSTRAINT profiles_name_pkey PRIMARY KEY (name);

--
-- Name: cachegroups_short_name; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX cachegroups_short_name ON cachegroups USING btree (description);


--
-- Name: federation_resolvers_ip_address; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX federation_resolvers_ip_address ON federation_resolvers USING btree (ip_address);


--
-- Name: parameters_name_config_file_value_idx; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX parameters_name_config_file_value_idx ON parameters USING btree (name, config_file, value);


--
-- Name: asns_cchegroup_cachegroups_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY asns
    ADD CONSTRAINT asns_cchegroup_cachegroups_name_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroups(name);

--
-- Name: cachegroups_parameters_cachegroup_cachegroups_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroups_parameters
    ADD CONSTRAINT cachegroups_parameters_cachegroup_cachegroups_name_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroups(name);


--
-- Name: cachegroups_parameters_parameter_id_parameters_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroups_parameters
    ADD CONSTRAINT cachegroups_parameters_parameter_id_parameters_id_fkey FOREIGN KEY (parameter_id) REFERENCES parameters(id);


--
-- Name: cachegroups_parent_cachegroup_cachegroups_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroups
    ADD CONSTRAINT cachegroups_parent_cachegroup_cachegroups_name_fkey FOREIGN KEY (parent_cachegroup) REFERENCES cachegroups(name);

--
-- Name: deliveryservices_cdn_cdns_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices
    ADD CONSTRAINT deliveryservices_cdn_cdns_name_fkey FOREIGN KEY (cdn) REFERENCES cdns(name);


--
-- Name: deliveryservices_domain_domains_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices
    ADD CONSTRAINT deliveryservices_domain_domains_name_fkey FOREIGN KEY (domain) REFERENCES domains(name);


--
-- Name: deliveryservices_profile_profiles_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservices
    ADD CONSTRAINT deliveryservices_profile_profiles_name_fkey FOREIGN KEY (profile) REFERENCES profiles(name);

--
-- Name: domains_cdn_fkey; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY domains
    ADD CONSTRAINT domains_cdn_fkey FOREIGN KEY (cdn) REFERENCES cdns(name);

CREATE TABLE crconfig_snapshots (
    cdn text NOT NULL REFERENCES cdns (name),
    snapshot text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (cdn, created_at)
);
ALTER TABLE crconfig_snapshots OWNER TO touser;

-- \todo add port to deliveryservices_servers and join on it. Because server's PK is a compound key, host_name and port.
CREATE OR REPLACE VIEW csconfig_remap as
select deliveryservices.*, regexes.pattern as r_pattern, servers.host_name as server_name
from servers
join deliveryservices_servers on deliveryservices_servers.server = servers.host_name
join deliveryservices on deliveryservices.name = deliveryservices_servers.deliveryservice
join deliveryservices_regexes on deliveryservices_regexes.deliveryservice = deliveryservices.name
join regexes on regexes.id = deliveryservices_regexes.regex_id;

--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;

--
-- PostgreSQL database dump complete
--

