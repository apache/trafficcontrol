--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.0
-- Dumped by pg_dump version 9.5.0

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

\connect traffic_ops

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

--
-- Name: asn_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE asn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE asn_id_seq OWNER TO touser;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: asn; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE asn (
    id integer DEFAULT nextval('asn_id_seq'::regclass) NOT NULL,
    asn integer NOT NULL,
    cachegroup integer DEFAULT 0 NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE asn OWNER TO touser;

--
-- Name: cachegroup_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE cachegroup_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cachegroup_id_seq OWNER TO touser;

--
-- Name: cachegroup; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cachegroup (
    id integer DEFAULT nextval('cachegroup_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    short_name character varying(255) NOT NULL,
    latitude double precision,
    longitude double precision,
    parent_cachegroup_id integer,
    type integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE cachegroup OWNER TO touser;

--
-- Name: api_asns; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW api_asns AS
 SELECT a.last_updated AS "lastUpdated",
    a.id,
    a.asn,
    c.name AS cachegroup
   FROM asn a,
    cachegroup c
  WHERE (a.id = c.id);


ALTER TABLE api_asns OWNER TO touser;

--
-- Name: profile_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE profile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE profile_id_seq OWNER TO touser;

--
-- Name: profile; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE profile (
    id integer DEFAULT nextval('profile_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(256),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE profile OWNER TO touser;

--
-- Name: api_profiles; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW api_profiles AS
 SELECT profile.last_updated AS "lastUpdated",
    profile.name,
    profile.id,
    profile.description
   FROM profile;


ALTER TABLE api_profiles OWNER TO touser;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE region_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE region_id_seq OWNER TO touser;

--
-- Name: region; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE region (
    id integer DEFAULT nextval('region_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    division integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE region OWNER TO touser;

--
-- Name: api_region; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW api_region AS
 SELECT region.id,
    region.name,
    region.division AS division_id
   FROM region;


ALTER TABLE api_region OWNER TO touser;

--
-- Name: cachegroup_parameter; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cachegroup_parameter (
    cachegroup integer DEFAULT 0 NOT NULL,
    parameter integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE cachegroup_parameter OWNER TO touser;

--
-- Name: cdn_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE cdn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cdn_id_seq OWNER TO touser;

--
-- Name: cdn; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE cdn (
    id integer DEFAULT nextval('cdn_id_seq'::regclass) NOT NULL,
    name character varying(127),
    last_updated timestamp without time zone DEFAULT now() NOT NULL,
    dnssec_enabled smallint DEFAULT 0
);


ALTER TABLE cdn OWNER TO touser;

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
-- Name: parameter; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE parameter (
    id integer DEFAULT nextval('parameter_id_seq'::regclass) NOT NULL,
    name character varying(1024) NOT NULL,
    config_file character varying(45) NOT NULL,
    value character varying(1024) NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE parameter OWNER TO touser;

--
-- Name: profile_parameter; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE profile_parameter (
    profile integer NOT NULL,
    parameter integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE profile_parameter OWNER TO touser;

--
-- Name: server_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE server_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE server_id_seq OWNER TO touser;

--
-- Name: server; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE server (
    id integer DEFAULT nextval('server_id_seq'::regclass) NOT NULL,
    host_name character varying(45) NOT NULL,
    domain_name character varying(45) NOT NULL,
    tcp_port bigint,
    xmpp_id character varying(256),
    xmpp_passwd character varying(45),
    interface_name character varying(45) NOT NULL,
    ip_address character varying(45) NOT NULL,
    ip_netmask character varying(45) NOT NULL,
    ip_gateway character varying(45) NOT NULL,
    ip6_address character varying(50),
    ip6_gateway character varying(50),
    interface_mtu integer DEFAULT 9000 NOT NULL,
    phys_location integer NOT NULL,
    rack character varying(64),
    cachegroup integer DEFAULT 0 NOT NULL,
    type integer NOT NULL,
    status integer NOT NULL,
    upd_pending boolean DEFAULT false NOT NULL,
    profile integer NOT NULL,
    cdn_id integer,
    mgmt_ip_address character varying(45),
    mgmt_ip_netmask character varying(45),
    mgmt_ip_gateway character varying(45),
    ilo_ip_address character varying(45),
    ilo_ip_netmask character varying(45),
    ilo_ip_gateway character varying(45),
    ilo_username character varying(45),
    ilo_password character varying(45),
    router_host_name character varying(256),
    router_port_name character varying(256),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE server OWNER TO touser;

--
-- Name: status_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE status_id_seq OWNER TO touser;

--
-- Name: status; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE status (
    id integer DEFAULT nextval('status_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(256),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE status OWNER TO touser;

--
-- Name: type_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE type_id_seq OWNER TO touser;

--
-- Name: type; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE type (
    id integer DEFAULT nextval('type_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(256),
    use_in_table character varying(45),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE type OWNER TO touser;

--
-- Name: content_routers; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW content_routers AS
 SELECT server.ip_address AS ip,
    server.ip6_address AS ip6,
    profile.name AS profile,
    cachegroup.name AS location,
    status.name AS status,
    server.tcp_port AS port,
    server.host_name,
    concat(server.host_name, '.', server.domain_name) AS fqdn,
    parameter.value AS apiport,
    cdn.name AS cdnname
   FROM (((((((server
     JOIN profile ON ((profile.id = server.profile)))
     JOIN profile_parameter ON ((profile_parameter.profile = profile.id)))
     JOIN parameter ON ((parameter.id = profile_parameter.parameter)))
     JOIN cachegroup ON ((cachegroup.id = server.cachegroup)))
     JOIN status ON ((status.id = server.status)))
     JOIN cdn ON ((cdn.id = server.cdn_id)))
     JOIN type ON ((type.id = server.type)))
  WHERE (((type.name)::text = 'CCR'::text) AND ((parameter.name)::text = 'api.port'::text));


ALTER TABLE content_routers OWNER TO touser;

--
-- Name: content_servers; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW content_servers AS
 SELECT DISTINCT server.host_name,
    profile.name AS profile,
    type.name AS type,
    cachegroup.name AS location_id,
    server.ip_address AS ip,
    cdn.name AS cdnname,
    status.name AS status,
    cachegroup.name AS cache_group,
    server.ip6_address AS ip6,
    server.tcp_port AS port,
    concat(server.host_name, '.', server.domain_name) AS fqdn,
    server.interface_name,
    parameter.value AS hash_count
   FROM (((((((server
     JOIN profile ON ((profile.id = server.profile)))
     JOIN profile_parameter ON ((profile_parameter.profile = profile.id)))
     JOIN parameter ON ((parameter.id = profile_parameter.parameter)))
     JOIN cachegroup ON ((cachegroup.id = server.cachegroup)))
     JOIN type ON ((type.id = server.type)))
     JOIN status ON ((status.id = server.status)))
     JOIN cdn ON (((cdn.id = server.cdn_id) AND ((parameter.name)::text = 'weight'::text) AND (server.status IN ( SELECT status_1.id
           FROM status status_1
          WHERE (((status_1.name)::text = 'REPORTED'::text) OR ((status_1.name)::text = 'ONLINE'::text)))) AND (server.type = ( SELECT type_1.id
           FROM type type_1
          WHERE ((type_1.name)::text = 'EDGE'::text))))));


ALTER TABLE content_servers OWNER TO touser;

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE deliveryservice_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE deliveryservice_id_seq OWNER TO touser;

--
-- Name: deliveryservice; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservice (
    id integer DEFAULT nextval('deliveryservice_id_seq'::regclass) NOT NULL,
    xml_id character varying(48) NOT NULL,
    active smallint NOT NULL,
    dscp integer NOT NULL,
    signed boolean,
    qstring_ignore boolean,
    geo_limit boolean DEFAULT false,
    http_bypass_fqdn character varying(255),
    dns_bypass_ip character varying(45),
    dns_bypass_ip6 character varying(45),
    dns_bypass_ttl integer,
    org_server_fqdn character varying(255),
    type integer NOT NULL,
    profile integer NOT NULL,
    cdn_id integer,
    ccr_dns_ttl integer,
    global_max_mbps integer,
    global_max_tps integer,
    long_desc character varying(1024),
    long_desc_1 character varying(1024),
    long_desc_2 character varying(1024),
    max_dns_answers integer DEFAULT 0,
    info_url character varying(255),
    miss_lat double precision,
    miss_long double precision,
    check_path character varying(255),
    last_updated timestamp without time zone DEFAULT now(),
    protocol smallint DEFAULT 0,
    ssl_key_version integer DEFAULT 0,
    ipv6_routing_enabled smallint,
    range_request_handling smallint DEFAULT 0,
    edge_header_rewrite character varying(2048),
    origin_shield character varying(1024),
    mid_header_rewrite character varying(2048),
    regex_remap character varying(1024),
    cacheurl character varying(1024),
    remap_text character varying(2048),
    multi_site_origin boolean,
    display_name character varying(48) NOT NULL,
    tr_response_headers character varying(1024),
    initial_dispersion integer DEFAULT 1,
    dns_bypass_cname character varying(255),
    tr_request_headers character varying(1024)
);


ALTER TABLE deliveryservice OWNER TO touser;

--
-- Name: deliveryservice_regex; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservice_regex (
    deliveryservice integer NOT NULL,
    regex integer NOT NULL,
    set_number integer DEFAULT 0
);


ALTER TABLE deliveryservice_regex OWNER TO touser;

--
-- Name: deliveryservice_server; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservice_server (
    deliveryservice integer NOT NULL,
    server integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE deliveryservice_server OWNER TO touser;

--
-- Name: regex_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE regex_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE regex_id_seq OWNER TO touser;

--
-- Name: regex; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE regex (
    id integer DEFAULT nextval('regex_id_seq'::regclass) NOT NULL,
    pattern character varying(255) DEFAULT ''::character varying NOT NULL,
    type integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE regex OWNER TO touser;

--
-- Name: cr_deliveryservice_server; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW cr_deliveryservice_server AS
 SELECT DISTINCT regex.pattern,
    deliveryservice.xml_id,
    deliveryservice.id AS ds_id,
    server.id AS srv_id,
    cdn.name AS cdnname,
    server.host_name AS server_name
   FROM (((((deliveryservice
     JOIN deliveryservice_regex ON ((deliveryservice_regex.deliveryservice = deliveryservice.id)))
     JOIN regex ON ((regex.id = deliveryservice_regex.regex)))
     JOIN deliveryservice_server ON ((deliveryservice.id = deliveryservice_server.deliveryservice)))
     JOIN server ON ((server.id = deliveryservice_server.server)))
     JOIN cdn ON ((cdn.id = server.cdn_id)))
  WHERE (deliveryservice.type <> ( SELECT type.id
           FROM type
          WHERE ((type.name)::text = 'ANY_MAP'::text)));


ALTER TABLE cr_deliveryservice_server OWNER TO touser;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE staticdnsentry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE staticdnsentry_id_seq OWNER TO touser;

--
-- Name: staticdnsentry; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE staticdnsentry (
    id integer DEFAULT nextval('staticdnsentry_id_seq'::regclass) NOT NULL,
    host character varying(45) NOT NULL,
    address character varying(45) NOT NULL,
    type integer NOT NULL,
    ttl integer DEFAULT 3600 NOT NULL,
    deliveryservice integer NOT NULL,
    cachegroup integer,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE staticdnsentry OWNER TO touser;

--
-- Name: crconfig_ds_data; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW crconfig_ds_data AS
 SELECT deliveryservice.xml_id,
    deliveryservice.profile,
    deliveryservice.ccr_dns_ttl,
    deliveryservice.global_max_mbps,
    deliveryservice.global_max_tps,
    deliveryservice.max_dns_answers,
    deliveryservice.miss_lat,
    deliveryservice.miss_long,
    protocoltype.name AS protocol,
    deliveryservice.ipv6_routing_enabled,
    deliveryservice.tr_request_headers,
    deliveryservice.tr_response_headers,
    deliveryservice.initial_dispersion,
    deliveryservice.dns_bypass_cname,
    deliveryservice.dns_bypass_ip,
    deliveryservice.dns_bypass_ip6,
    deliveryservice.dns_bypass_ttl,
    deliveryservice.geo_limit,
    cdn.name AS cdn_name,
    regex.pattern AS match_pattern,
    regextype.name AS match_type,
    deliveryservice_regex.set_number,
    staticdnsentry.host AS sdns_host,
    staticdnsentry.address AS sdns_address,
    staticdnsentry.ttl AS sdns_ttl,
    sdnstype.name AS sdns_type
   FROM (((((((deliveryservice
     JOIN cdn ON ((cdn.id = deliveryservice.cdn_id)))
     LEFT JOIN staticdnsentry ON ((deliveryservice.id = staticdnsentry.deliveryservice)))
     JOIN deliveryservice_regex ON ((deliveryservice_regex.deliveryservice = deliveryservice.id)))
     JOIN regex ON ((regex.id = deliveryservice_regex.regex)))
     JOIN type protocoltype ON ((protocoltype.id = deliveryservice.type)))
     JOIN type regextype ON ((regextype.id = regex.type)))
     LEFT JOIN type sdnstype ON ((sdnstype.id = staticdnsentry.type)));


ALTER TABLE crconfig_ds_data OWNER TO touser;

--
-- Name: crconfig_params; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW crconfig_params AS
 SELECT DISTINCT cdn.name AS cdn_name,
    cdn.id AS cdn_id,
    server.profile AS profile_id,
    server.type AS stype,
    parameter.name AS pname,
    parameter.config_file AS cfile,
    parameter.value AS pvalue
   FROM ((((server
     JOIN cdn ON ((cdn.id = server.cdn_id)))
     JOIN profile ON ((profile.id = server.profile)))
     JOIN profile_parameter ON ((profile_parameter.profile = server.profile)))
     JOIN parameter ON ((parameter.id = profile_parameter.parameter)))
  WHERE ((server.type IN ( SELECT type.id
           FROM type
          WHERE ((type.name)::text = ANY ((ARRAY['EDGE'::character varying, 'MID'::character varying, 'CCR'::character varying])::text[])))) AND ((parameter.config_file)::text = 'CRConfig.json'::text));


ALTER TABLE crconfig_params OWNER TO touser;

--
-- Name: deliveryservice_tmuser; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE deliveryservice_tmuser (
    deliveryservice integer NOT NULL,
    tm_user_id integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE deliveryservice_tmuser OWNER TO touser;

--
-- Name: division_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE division_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE division_id_seq OWNER TO touser;

--
-- Name: division; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE division (
    id integer DEFAULT nextval('division_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE division OWNER TO touser;

--
-- Name: divisions; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW divisions AS
 SELECT division.id,
    division.name,
    division.last_updated
   FROM division;


ALTER TABLE divisions OWNER TO touser;

--
-- Name: federation_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE federation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federation_id_seq OWNER TO touser;

--
-- Name: federation; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation (
    id integer DEFAULT nextval('federation_id_seq'::regclass) NOT NULL,
    cname character varying(1024) NOT NULL,
    description character varying(1024),
    ttl integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE federation OWNER TO touser;

--
-- Name: federation_deliveryservice; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_deliveryservice (
    federation integer NOT NULL,
    deliveryservice integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE federation_deliveryservice OWNER TO touser;

--
-- Name: federation_federation_resolver; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_federation_resolver (
    federation integer NOT NULL,
    federation_resolver integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE federation_federation_resolver OWNER TO touser;

--
-- Name: federation_resolver_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE federation_resolver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federation_resolver_id_seq OWNER TO touser;

--
-- Name: federation_resolver; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_resolver (
    id integer DEFAULT nextval('federation_resolver_id_seq'::regclass) NOT NULL,
    ip_address character varying(50) NOT NULL,
    type integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE federation_resolver OWNER TO touser;

--
-- Name: federation_tmuser; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE federation_tmuser (
    federation integer NOT NULL,
    tm_user integer NOT NULL,
    role integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE federation_tmuser OWNER TO touser;

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
-- Name: hwinfo_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE hwinfo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE hwinfo_id_seq OWNER TO touser;

--
-- Name: hwinfo; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE hwinfo (
    id integer DEFAULT nextval('hwinfo_id_seq'::regclass) NOT NULL,
    serverid integer NOT NULL,
    description character varying(256) NOT NULL,
    val character varying(256) NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE hwinfo OWNER TO touser;

--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_id_seq OWNER TO touser;

--
-- Name: job; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE job (
    id integer DEFAULT nextval('job_id_seq'::regclass) NOT NULL,
    agent integer,
    object_type character varying(48),
    object_name character varying(256),
    keyword character varying(48) NOT NULL,
    parameters character varying(256),
    asset_url character varying(512) NOT NULL,
    asset_type character varying(48) NOT NULL,
    status integer NOT NULL,
    start_time timestamp without time zone NOT NULL,
    entered_time timestamp without time zone NOT NULL,
    job_user integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now(),
    job_deliveryservice integer
);


ALTER TABLE job OWNER TO touser;

--
-- Name: job_agent_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE job_agent_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_agent_id_seq OWNER TO touser;

--
-- Name: job_agent; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE job_agent (
    id integer DEFAULT nextval('job_agent_id_seq'::regclass) NOT NULL,
    name character varying(128),
    description character varying(512),
    active integer DEFAULT 0 NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE job_agent OWNER TO touser;

--
-- Name: job_result_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE job_result_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_result_id_seq OWNER TO touser;

--
-- Name: job_result; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE job_result (
    id integer DEFAULT nextval('job_result_id_seq'::regclass) NOT NULL,
    job integer NOT NULL,
    agent integer NOT NULL,
    result character varying(48) NOT NULL,
    description character varying(512),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE job_result OWNER TO touser;

--
-- Name: job_status_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE job_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_status_id_seq OWNER TO touser;

--
-- Name: job_status; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE job_status (
    id integer DEFAULT nextval('job_status_id_seq'::regclass) NOT NULL,
    name character varying(48),
    description character varying(256),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE job_status OWNER TO touser;

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
    level character varying(45),
    message character varying(1024) NOT NULL,
    tm_user integer NOT NULL,
    ticketnum character varying(64),
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE log OWNER TO touser;

--
-- Name: monitors; Type: VIEW; Schema: public; Owner: touser
--

CREATE VIEW monitors AS
 SELECT server.ip_address AS ip,
    server.ip6_address AS ip6,
    profile.name AS profile,
    cachegroup.name AS location,
    status.name AS status,
    server.tcp_port AS port,
    concat(server.host_name, '.', server.domain_name) AS fqdn,
    cdn.name AS cdnname,
    server.host_name
   FROM (((((server
     JOIN profile ON ((profile.id = server.profile)))
     JOIN cachegroup ON ((cachegroup.id = server.cachegroup)))
     JOIN status ON ((status.id = server.status)))
     JOIN cdn ON ((cdn.id = server.cdn_id)))
     JOIN type ON ((type.id = server.type)))
  WHERE ((type.name)::text = 'RASCAL'::text);


ALTER TABLE monitors OWNER TO touser;

--
-- Name: phys_location_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE phys_location_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE phys_location_id_seq OWNER TO touser;

--
-- Name: phys_location; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE phys_location (
    id integer DEFAULT nextval('phys_location_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    short_name character varying(12) NOT NULL,
    address character varying(128) NOT NULL,
    city character varying(128) NOT NULL,
    state character varying(2) NOT NULL,
    zip character varying(5) NOT NULL,
    poc character varying(128),
    phone character varying(45),
    email character varying(128),
    comments character varying(256),
    region integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE phys_location OWNER TO touser;

--
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE role_id_seq OWNER TO touser;

--
-- Name: role; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE role (
    id integer DEFAULT nextval('role_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(128),
    priv_level integer NOT NULL
);


ALTER TABLE role OWNER TO touser;

--
-- Name: servercheck_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE servercheck_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE servercheck_id_seq OWNER TO touser;

--
-- Name: servercheck; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE servercheck (
    id integer DEFAULT nextval('servercheck_id_seq'::regclass) NOT NULL,
    server integer NOT NULL,
    aa integer,
    ab integer,
    ac integer,
    ad integer,
    ae integer,
    af integer,
    ag integer,
    ah integer,
    ai integer,
    aj integer,
    ak integer,
    al integer,
    am integer,
    an integer,
    ao integer,
    ap integer,
    aq integer,
    ar integer,
    "as" integer,
    at integer,
    au integer,
    av integer,
    aw integer,
    ax integer,
    ay integer,
    az integer,
    ba integer,
    bb integer,
    bc integer,
    bd integer,
    be integer,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE servercheck OWNER TO touser;

--
-- Name: stats_summary_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE stats_summary_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE stats_summary_id_seq OWNER TO touser;

--
-- Name: stats_summary; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE stats_summary (
    id integer DEFAULT nextval('stats_summary_id_seq'::regclass) NOT NULL,
    cdn_name character varying(255) DEFAULT 'all'::character varying NOT NULL,
    deliveryservice_name character varying(255) NOT NULL,
    stat_name character varying(255) NOT NULL,
    stat_value real NOT NULL,
    summary_time timestamp without time zone DEFAULT now() NOT NULL,
    stat_date date
);


ALTER TABLE stats_summary OWNER TO touser;

--
-- Name: tm_user_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE tm_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE tm_user_id_seq OWNER TO touser;

--
-- Name: tm_user; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE tm_user (
    id integer DEFAULT nextval('tm_user_id_seq'::regclass) NOT NULL,
    username character varying(128),
    role integer,
    uid integer,
    gid integer,
    local_passwd character varying(40),
    confirm_local_passwd character varying(40),
    last_updated timestamp without time zone DEFAULT now(),
    company character varying(256),
    email character varying(128),
    full_name character varying(256),
    new_user boolean DEFAULT true NOT NULL,
    address_line1 character varying(256),
    address_line2 character varying(256),
    city character varying(128),
    state_or_province character varying(128),
    phone_number character varying(25),
    postal_code character varying(11),
    country character varying(256),
    local_user boolean DEFAULT false NOT NULL,
    token character varying(50),
    registration_sent timestamp without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone NOT NULL
);


ALTER TABLE tm_user OWNER TO touser;

--
-- Name: to_extension_id_seq; Type: SEQUENCE; Schema: public; Owner: touser
--

CREATE SEQUENCE to_extension_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE to_extension_id_seq OWNER TO touser;

--
-- Name: to_extension; Type: TABLE; Schema: public; Owner: touser
--

CREATE TABLE to_extension (
    id integer DEFAULT nextval('to_extension_id_seq'::regclass) NOT NULL,
    name character varying(45) NOT NULL,
    version character varying(45) NOT NULL,
    info_url character varying(45) NOT NULL,
    script_file character varying(45) NOT NULL,
    isactive boolean NOT NULL,
    additional_config_json character varying(4096),
    description character varying(4096),
    servercheck_short_name character varying(8),
    servercheck_column_name character varying(10),
    type integer NOT NULL,
    last_updated timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE to_extension OWNER TO touser;

--
-- Name: asn_id_cachegroup_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY asn
    ADD CONSTRAINT asn_id_cachegroup_pkey PRIMARY KEY (id, cachegroup);


--
-- Name: cachegroup_id_type_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT cachegroup_id_type_pkey PRIMARY KEY (id, type);


--
-- Name: cachegroup_parameter_cachegroup_parameter_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cachegroup_parameter
    ADD CONSTRAINT cachegroup_parameter_cachegroup_parameter_pkey PRIMARY KEY (cachegroup, parameter);


--
-- Name: cdn_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY cdn
    ADD CONSTRAINT cdn_id_pkey PRIMARY KEY (id);


--
-- Name: deliveryservice_id_type_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT deliveryservice_id_type_pkey PRIMARY KEY (id, type);


--
-- Name: deliveryservice_regex_deliveryservice_regex_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT deliveryservice_regex_deliveryservice_regex_pkey PRIMARY KEY (deliveryservice, regex);


--
-- Name: deliveryservice_server_deliveryservice_server_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT deliveryservice_server_deliveryservice_server_pkey PRIMARY KEY (deliveryservice, server);


--
-- Name: deliveryservice_tmuser_deliveryservice_tm_user_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT deliveryservice_tmuser_deliveryservice_tm_user_id_pkey PRIMARY KEY (deliveryservice, tm_user_id);


--
-- Name: division_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY division
    ADD CONSTRAINT division_id_pkey PRIMARY KEY (id);


--
-- Name: federation_deliveryservice_federation_deliveryservice_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_deliveryservice
    ADD CONSTRAINT federation_deliveryservice_federation_deliveryservice_pkey PRIMARY KEY (federation, deliveryservice);


--
-- Name: federation_federation_resolver_federation_federation_resolver_p; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_federation_resolver
    ADD CONSTRAINT federation_federation_resolver_federation_federation_resolver_p PRIMARY KEY (federation, federation_resolver);


--
-- Name: federation_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation
    ADD CONSTRAINT federation_id_pkey PRIMARY KEY (id);


--
-- Name: federation_resolver_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_resolver
    ADD CONSTRAINT federation_resolver_id_pkey PRIMARY KEY (id);


--
-- Name: federation_tmuser_federation_tm_user_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY federation_tmuser
    ADD CONSTRAINT federation_tmuser_federation_tm_user_pkey PRIMARY KEY (federation, tm_user);


--
-- Name: goose_db_version_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY goose_db_version
    ADD CONSTRAINT goose_db_version_id_pkey PRIMARY KEY (id);


--
-- Name: hwinfo_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY hwinfo
    ADD CONSTRAINT hwinfo_id_pkey PRIMARY KEY (id);


--
-- Name: job_id_pkey; Type: CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY job
    ADD CONSTRAINT job_id_pkey PRIMARY KEY (id);


--
-- Name: asn_cachegroup; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX asn_cachegroup ON asn USING btree (cachegroup);


--
-- Name: asn_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX asn_id ON asn USING btree (id);


--
-- Name: cachegroup_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX cachegroup_id ON cachegroup USING btree (id);


--
-- Name: cachegroup_name; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX cachegroup_name ON cachegroup USING btree (name);


--
-- Name: cachegroup_parameter_parameter; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX cachegroup_parameter_parameter ON cachegroup_parameter USING btree (parameter);


--
-- Name: cachegroup_parent_cachegroup_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX cachegroup_parent_cachegroup_id ON cachegroup USING btree (parent_cachegroup_id);


--
-- Name: cachegroup_short_name; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX cachegroup_short_name ON cachegroup USING btree (short_name);


--
-- Name: cachegroup_type; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX cachegroup_type ON cachegroup USING btree (type);


--
-- Name: cdn_name; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX cdn_name ON cdn USING btree (name);


--
-- Name: deliveryservice_cdn_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_cdn_id ON deliveryservice USING btree (cdn_id);


--
-- Name: deliveryservice_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX deliveryservice_id ON deliveryservice USING btree (id);


--
-- Name: deliveryservice_profile; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_profile ON deliveryservice USING btree (profile);


--
-- Name: deliveryservice_regex_regex; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_regex_regex ON deliveryservice_regex USING btree (regex);


--
-- Name: deliveryservice_server_server; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_server_server ON deliveryservice_server USING btree (server);


--
-- Name: deliveryservice_tmuser_tm_user_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_tmuser_tm_user_id ON deliveryservice_tmuser USING btree (tm_user_id);


--
-- Name: deliveryservice_type; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX deliveryservice_type ON deliveryservice USING btree (type);


--
-- Name: deliveryservice_xml_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX deliveryservice_xml_id ON deliveryservice USING btree (xml_id);


--
-- Name: division_name; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX division_name ON division USING btree (name);


--
-- Name: federation_deliveryservice_deliveryservice; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_deliveryservice_deliveryservice ON federation_deliveryservice USING btree (deliveryservice);


--
-- Name: federation_federation_resolver_federation; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_federation_resolver_federation ON federation_federation_resolver USING btree (federation);


--
-- Name: federation_federation_resolver_federation_resolver; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_federation_resolver_federation_resolver ON federation_federation_resolver USING btree (federation_resolver);


--
-- Name: federation_resolver_ip_address; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX federation_resolver_ip_address ON federation_resolver USING btree (ip_address);


--
-- Name: federation_resolver_type; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_resolver_type ON federation_resolver USING btree (type);


--
-- Name: federation_tmuser_federation; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_tmuser_federation ON federation_tmuser USING btree (federation);


--
-- Name: federation_tmuser_role; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_tmuser_role ON federation_tmuser USING btree (role);


--
-- Name: federation_tmuser_tm_user; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX federation_tmuser_tm_user ON federation_tmuser USING btree (tm_user);


--
-- Name: goose_db_version_id; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX goose_db_version_id ON goose_db_version USING btree (id);


--
-- Name: hwinfo_serverid; Type: INDEX; Schema: public; Owner: touser
--

CREATE INDEX hwinfo_serverid ON hwinfo USING btree (serverid);


--
-- Name: hwinfo_serverid_description; Type: INDEX; Schema: public; Owner: touser
--

CREATE UNIQUE INDEX hwinfo_serverid_description ON hwinfo USING btree (serverid, description);


--
-- Name: divisionfk; Type: FK CONSTRAINT; Schema: public; Owner: touser
--

ALTER TABLE ONLY region
    ADD CONSTRAINT divisionfk FOREIGN KEY (division) REFERENCES division(id);


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

