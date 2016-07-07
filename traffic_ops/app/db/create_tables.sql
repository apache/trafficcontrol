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

-- Dumped from database version 9.5.3
-- Dumped by pg_dump version 9.5.3

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
-- Name: on_update_current_timestamp_last_updated(); Type: FUNCTION; Schema: public; Owner: jheitz200
--

CREATE FUNCTION on_update_current_timestamp_last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
   NEW.last_updated = now();
   RETURN NEW;
END;
$$;


ALTER FUNCTION public.on_update_current_timestamp_last_updated() OWNER TO jheitz200;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: cran; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE cran (
    id bigint NOT NULL,
    asn bigint NOT NULL,
    location bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE cran OWNER TO jheitz200;

--
-- Name: cran_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE cran_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cran_id_seq OWNER TO jheitz200;

--
-- Name: cran_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE cran_id_seq OWNED BY cran.id;


--
-- Name: deliveryservice; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE deliveryservice (
    id bigint NOT NULL,
    xml_id character varying(48) NOT NULL,
    active smallint NOT NULL,
    dscp bigint NOT NULL,
    signed boolean,
    qstring_ignore boolean,
    geo_limit boolean DEFAULT false,
    http_bypass_fqdn character varying(255),
    dns_bypass_ip character varying(45),
    dns_bypass_ip6 character varying(45),
    dns_bypass_ttl bigint,
    org_server_fqdn character varying(255),
    type bigint NOT NULL,
    profile bigint NOT NULL,
    ccr_dns_ttl bigint,
    global_max_mbps bigint,
    global_max_tps bigint,
    long_desc character varying(255),
    long_desc_1 character varying(255),
    long_desc_2 character varying(255),
    max_dns_answers bigint DEFAULT '0'::bigint,
    info_url character varying(255),
    miss_lat double precision,
    miss_long double precision,
    check_path character varying(255),
    header_rewrite bigint,
    last_updated timestamp with time zone DEFAULT now(),
    protocol smallint DEFAULT '0'::smallint
);


ALTER TABLE deliveryservice OWNER TO jheitz200;

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE deliveryservice_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE deliveryservice_id_seq OWNER TO jheitz200;

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE deliveryservice_id_seq OWNED BY deliveryservice.id;


--
-- Name: deliveryservice_regex; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE deliveryservice_regex (
    deliveryservice bigint NOT NULL,
    regex bigint NOT NULL,
    set_number bigint DEFAULT '0'::bigint
);


ALTER TABLE deliveryservice_regex OWNER TO jheitz200;

--
-- Name: deliveryservice_server; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE deliveryservice_server (
    deliveryservice bigint NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE deliveryservice_server OWNER TO jheitz200;

--
-- Name: deliveryservice_tmuser; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE deliveryservice_tmuser (
    deliveryservice bigint NOT NULL,
    tm_user_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE deliveryservice_tmuser OWNER TO jheitz200;

--
-- Name: division; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE division (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE division OWNER TO jheitz200;

--
-- Name: division_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE division_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE division_id_seq OWNER TO jheitz200;

--
-- Name: division_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE division_id_seq OWNED BY division.id;


--
-- Name: header_rewrite; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE header_rewrite (
    id bigint NOT NULL,
    hr_condition character varying(1024),
    action character varying(1024) NOT NULL
);


ALTER TABLE header_rewrite OWNER TO jheitz200;

--
-- Name: header_rewrite_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE header_rewrite_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE header_rewrite_id_seq OWNER TO jheitz200;

--
-- Name: header_rewrite_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE header_rewrite_id_seq OWNED BY header_rewrite.id;


--
-- Name: hwinfo; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE hwinfo (
    id bigint NOT NULL,
    serverid bigint NOT NULL,
    description character varying(256) NOT NULL,
    val character varying(256) NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE hwinfo OWNER TO jheitz200;

--
-- Name: hwinfo_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE hwinfo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE hwinfo_id_seq OWNER TO jheitz200;

--
-- Name: hwinfo_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE hwinfo_id_seq OWNED BY hwinfo.id;


--
-- Name: job; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE job (
    id bigint NOT NULL,
    agent bigint,
    object_type character varying(48),
    object_name character varying(256),
    keyword character varying(48) NOT NULL,
    parameters character varying(256),
    asset_url character varying(512) NOT NULL,
    asset_type character varying(48) NOT NULL,
    status bigint NOT NULL,
    start_time timestamp with time zone NOT NULL,
    entered_time timestamp with time zone NOT NULL,
    job_user bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE job OWNER TO jheitz200;

--
-- Name: job_agent; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE job_agent (
    id bigint NOT NULL,
    name character varying(128),
    description character varying(512),
    active integer DEFAULT 0 NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE job_agent OWNER TO jheitz200;

--
-- Name: job_agent_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE job_agent_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_agent_id_seq OWNER TO jheitz200;

--
-- Name: job_agent_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE job_agent_id_seq OWNED BY job_agent.id;


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_id_seq OWNER TO jheitz200;

--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE job_id_seq OWNED BY job.id;


--
-- Name: job_result; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE job_result (
    id bigint NOT NULL,
    job bigint NOT NULL,
    agent bigint NOT NULL,
    result character varying(48) NOT NULL,
    description character varying(512),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE job_result OWNER TO jheitz200;

--
-- Name: job_result_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE job_result_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_result_id_seq OWNER TO jheitz200;

--
-- Name: job_result_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE job_result_id_seq OWNED BY job_result.id;


--
-- Name: job_status; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE job_status (
    id bigint NOT NULL,
    name character varying(48),
    description character varying(256),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE job_status OWNER TO jheitz200;

--
-- Name: job_status_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE job_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_status_id_seq OWNER TO jheitz200;

--
-- Name: job_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE job_status_id_seq OWNED BY job_status.id;


--
-- Name: location; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE location (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    short_name character varying(255) NOT NULL,
    latitude double precision,
    longitude double precision,
    parent_location_id bigint,
    type bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE location OWNER TO jheitz200;

--
-- Name: location_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE location_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE location_id_seq OWNER TO jheitz200;

--
-- Name: location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE location_id_seq OWNED BY location.id;


--
-- Name: location_parameter; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE location_parameter (
    location bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE location_parameter OWNER TO jheitz200;

--
-- Name: log; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE log (
    id bigint NOT NULL,
    level character varying(45),
    message character varying(1024) NOT NULL,
    tm_user bigint NOT NULL,
    ticketnum character varying(64),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE log OWNER TO jheitz200;

--
-- Name: log_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE log_id_seq OWNER TO jheitz200;

--
-- Name: log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE log_id_seq OWNED BY log.id;


--
-- Name: parameter; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE parameter (
    id bigint NOT NULL,
    name character varying(1024) NOT NULL,
    config_file character varying(45) NOT NULL,
    value character varying(1024) NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE parameter OWNER TO jheitz200;

--
-- Name: parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE parameter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE parameter_id_seq OWNER TO jheitz200;

--
-- Name: parameter_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE parameter_id_seq OWNED BY parameter.id;


--
-- Name: phys_location; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE phys_location (
    id bigint NOT NULL,
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
    region bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE phys_location OWNER TO jheitz200;

--
-- Name: phys_location_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE phys_location_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE phys_location_id_seq OWNER TO jheitz200;

--
-- Name: phys_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE phys_location_id_seq OWNED BY phys_location.id;


--
-- Name: profile; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE profile (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(256),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE profile OWNER TO jheitz200;

--
-- Name: profile_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE profile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE profile_id_seq OWNER TO jheitz200;

--
-- Name: profile_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE profile_id_seq OWNED BY profile.id;


--
-- Name: profile_parameter; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE profile_parameter (
    profile bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE profile_parameter OWNER TO jheitz200;

--
-- Name: regex; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE regex (
    id bigint NOT NULL,
    pattern character varying(255) DEFAULT ''::character varying NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE regex OWNER TO jheitz200;

--
-- Name: regex_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE regex_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE regex_id_seq OWNER TO jheitz200;

--
-- Name: regex_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE regex_id_seq OWNED BY regex.id;


--
-- Name: region; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE region (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    division bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE region OWNER TO jheitz200;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE region_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE region_id_seq OWNER TO jheitz200;

--
-- Name: region_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE region_id_seq OWNED BY region.id;


--
-- Name: role; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE role (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(128),
    priv_level bigint NOT NULL
);


ALTER TABLE role OWNER TO jheitz200;

--
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE role_id_seq OWNER TO jheitz200;

--
-- Name: role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE role_id_seq OWNED BY role.id;


--
-- Name: server; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE server (
    id bigint NOT NULL,
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
    interface_mtu bigint DEFAULT '9000'::bigint NOT NULL,
    phys_location bigint NOT NULL,
    rack character varying(64),
    location bigint NOT NULL,
    type bigint NOT NULL,
    status bigint NOT NULL,
    profile bigint NOT NULL,
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
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE server OWNER TO jheitz200;

--
-- Name: server_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE server_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE server_id_seq OWNER TO jheitz200;

--
-- Name: server_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE server_id_seq OWNED BY server.id;


--
-- Name: serverstatus; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE serverstatus (
    id bigint NOT NULL,
    ilo_pingable boolean DEFAULT false NOT NULL,
    teng_pingable boolean DEFAULT false NOT NULL,
    fqdn_pingable boolean DEFAULT false,
    dscp boolean,
    firmware boolean,
    marvin boolean,
    ping6 boolean,
    upd_pending boolean,
    stats boolean,
    prox boolean,
    mtu boolean,
    ccr_online boolean,
    rascal boolean,
    chr bigint,
    cdu bigint,
    ort_errors bigint DEFAULT '-1'::bigint NOT NULL,
    mbps_out bigint DEFAULT '0'::bigint,
    clients_connected bigint DEFAULT '0'::bigint,
    server bigint NOT NULL,
    last_recycle_date timestamp with time zone,
    last_recycle_duration_hrs bigint DEFAULT '0'::bigint,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE serverstatus OWNER TO jheitz200;

--
-- Name: serverstatus_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE serverstatus_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE serverstatus_id_seq OWNER TO jheitz200;

--
-- Name: serverstatus_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE serverstatus_id_seq OWNED BY serverstatus.id;


--
-- Name: staticdnsentry; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE staticdnsentry (
    id bigint NOT NULL,
    host character varying(45) NOT NULL,
    address character varying(45) NOT NULL,
    type bigint NOT NULL,
    ttl bigint DEFAULT '3600'::bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    location bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE staticdnsentry OWNER TO jheitz200;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE staticdnsentry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE staticdnsentry_id_seq OWNER TO jheitz200;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE staticdnsentry_id_seq OWNED BY staticdnsentry.id;


--
-- Name: status; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE status (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(256),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE status OWNER TO jheitz200;

--
-- Name: status_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE status_id_seq OWNER TO jheitz200;

--
-- Name: status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE status_id_seq OWNED BY status.id;


--
-- Name: tm_user; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE tm_user (
    id bigint NOT NULL,
    username character varying(128),
    role bigint,
    uid bigint,
    gid bigint,
    local_passwd character varying(40),
    confirm_local_passwd character varying(40),
    last_updated timestamp with time zone DEFAULT now(),
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
    registration_sent timestamp with time zone DEFAULT '1998-12-31 17:00:00-07'::timestamp with time zone NOT NULL
);


ALTER TABLE tm_user OWNER TO jheitz200;

--
-- Name: tm_user_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE tm_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE tm_user_id_seq OWNER TO jheitz200;

--
-- Name: tm_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE tm_user_id_seq OWNED BY tm_user.id;


--
-- Name: to_extension; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE to_extension (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    version character varying(45) NOT NULL,
    info_url character varying(45) NOT NULL,
    script_file character varying(45) NOT NULL,
    isactive boolean NOT NULL,
    additional_config_json character varying(4096),
    description character varying(4096),
    servercheck_short_name character varying(8),
    servercheck_column_name character varying(10),
    type bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE to_extension OWNER TO jheitz200;

--
-- Name: to_extension_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE to_extension_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE to_extension_id_seq OWNER TO jheitz200;

--
-- Name: to_extension_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE to_extension_id_seq OWNED BY to_extension.id;


--
-- Name: type; Type: TABLE; Schema: public; Owner: jheitz200
--

CREATE TABLE type (
    id bigint NOT NULL,
    name character varying(45) NOT NULL,
    description character varying(45) NOT NULL,
    use_in_table character varying(45),
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE type OWNER TO jheitz200;

--
-- Name: type_id_seq; Type: SEQUENCE; Schema: public; Owner: jheitz200
--

CREATE SEQUENCE type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE type_id_seq OWNER TO jheitz200;

--
-- Name: type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: jheitz200
--

ALTER SEQUENCE type_id_seq OWNED BY type.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY cran ALTER COLUMN id SET DEFAULT nextval('cran_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice ALTER COLUMN id SET DEFAULT nextval('deliveryservice_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY division ALTER COLUMN id SET DEFAULT nextval('division_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY header_rewrite ALTER COLUMN id SET DEFAULT nextval('header_rewrite_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY hwinfo ALTER COLUMN id SET DEFAULT nextval('hwinfo_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job ALTER COLUMN id SET DEFAULT nextval('job_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_agent ALTER COLUMN id SET DEFAULT nextval('job_agent_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_result ALTER COLUMN id SET DEFAULT nextval('job_result_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_status ALTER COLUMN id SET DEFAULT nextval('job_status_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location ALTER COLUMN id SET DEFAULT nextval('location_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY log ALTER COLUMN id SET DEFAULT nextval('log_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY parameter ALTER COLUMN id SET DEFAULT nextval('parameter_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY phys_location ALTER COLUMN id SET DEFAULT nextval('phys_location_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY profile ALTER COLUMN id SET DEFAULT nextval('profile_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY regex ALTER COLUMN id SET DEFAULT nextval('regex_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY region ALTER COLUMN id SET DEFAULT nextval('region_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY role ALTER COLUMN id SET DEFAULT nextval('role_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server ALTER COLUMN id SET DEFAULT nextval('server_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY serverstatus ALTER COLUMN id SET DEFAULT nextval('serverstatus_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY staticdnsentry ALTER COLUMN id SET DEFAULT nextval('staticdnsentry_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY status ALTER COLUMN id SET DEFAULT nextval('status_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY tm_user ALTER COLUMN id SET DEFAULT nextval('tm_user_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY to_extension ALTER COLUMN id SET DEFAULT nextval('to_extension_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY type ALTER COLUMN id SET DEFAULT nextval('type_id_seq'::regclass);


--
-- Data for Name: cran; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY cran (id, asn, location, last_updated) FROM stdin;
\.


--
-- Name: cran_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('cran_id_seq', 1, true);


--
-- Data for Name: deliveryservice; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY deliveryservice (id, xml_id, active, dscp, signed, qstring_ignore, geo_limit, http_bypass_fqdn, dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, org_server_fqdn, type, profile, ccr_dns_ttl, global_max_mbps, global_max_tps, long_desc, long_desc_1, long_desc_2, max_dns_answers, info_url, miss_lat, miss_long, check_path, header_rewrite, last_updated, protocol) FROM stdin;
\.


--
-- Name: deliveryservice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('deliveryservice_id_seq', 1, true);


--
-- Data for Name: deliveryservice_regex; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY deliveryservice_regex (deliveryservice, regex, set_number) FROM stdin;
\.


--
-- Data for Name: deliveryservice_server; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY deliveryservice_server (deliveryservice, server, last_updated) FROM stdin;
\.


--
-- Data for Name: deliveryservice_tmuser; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY deliveryservice_tmuser (deliveryservice, tm_user_id, last_updated) FROM stdin;
\.


--
-- Data for Name: division; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY division (id, name, last_updated) FROM stdin;
\.


--
-- Name: division_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('division_id_seq', 1, true);


--
-- Data for Name: header_rewrite; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY header_rewrite (id, hr_condition, action) FROM stdin;
\.


--
-- Name: header_rewrite_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('header_rewrite_id_seq', 1, true);


--
-- Data for Name: hwinfo; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY hwinfo (id, serverid, description, val, last_updated) FROM stdin;
\.


--
-- Name: hwinfo_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('hwinfo_id_seq', 1, true);


--
-- Data for Name: job; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY job (id, agent, object_type, object_name, keyword, parameters, asset_url, asset_type, status, start_time, entered_time, job_user, last_updated) FROM stdin;
\.


--
-- Data for Name: job_agent; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY job_agent (id, name, description, active, last_updated) FROM stdin;
\.


--
-- Name: job_agent_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('job_agent_id_seq', 1, true);


--
-- Name: job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('job_id_seq', 1, true);


--
-- Data for Name: job_result; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY job_result (id, job, agent, result, description, last_updated) FROM stdin;
\.


--
-- Name: job_result_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('job_result_id_seq', 1, true);


--
-- Data for Name: job_status; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY job_status (id, name, description, last_updated) FROM stdin;
\.


--
-- Name: job_status_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('job_status_id_seq', 1, true);


--
-- Data for Name: location; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY location (id, name, short_name, latitude, longitude, parent_location_id, type, last_updated) FROM stdin;
\.


--
-- Name: location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('location_id_seq', 1, true);


--
-- Data for Name: location_parameter; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY location_parameter (location, parameter, last_updated) FROM stdin;
\.


--
-- Data for Name: log; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY log (id, level, message, tm_user, ticketnum, last_updated) FROM stdin;
\.


--
-- Name: log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('log_id_seq', 1, true);


--
-- Data for Name: parameter; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY parameter (id, name, config_file, value, last_updated) FROM stdin;
\.


--
-- Name: parameter_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('parameter_id_seq', 1, true);


--
-- Data for Name: phys_location; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY phys_location (id, name, short_name, address, city, state, zip, poc, phone, email, comments, region, last_updated) FROM stdin;
\.


--
-- Name: phys_location_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('phys_location_id_seq', 1, true);


--
-- Data for Name: profile; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY profile (id, name, description, last_updated) FROM stdin;
\.


--
-- Name: profile_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('profile_id_seq', 1, true);


--
-- Data for Name: profile_parameter; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY profile_parameter (profile, parameter, last_updated) FROM stdin;
\.


--
-- Data for Name: regex; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY regex (id, pattern, type, last_updated) FROM stdin;
\.


--
-- Name: regex_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('regex_id_seq', 1, true);


--
-- Data for Name: region; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY region (id, name, division, last_updated) FROM stdin;
\.


--
-- Name: region_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('region_id_seq', 1, true);


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY role (id, name, description, priv_level) FROM stdin;
\.


--
-- Name: role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('role_id_seq', 1, true);


--
-- Data for Name: server; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY server (id, host_name, domain_name, tcp_port, xmpp_id, xmpp_passwd, interface_name, ip_address, ip_netmask, ip_gateway, ip6_address, ip6_gateway, interface_mtu, phys_location, rack, location, type, status, profile, mgmt_ip_address, mgmt_ip_netmask, mgmt_ip_gateway, ilo_ip_address, ilo_ip_netmask, ilo_ip_gateway, ilo_username, ilo_password, router_host_name, router_port_name, last_updated) FROM stdin;
\.


--
-- Name: server_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('server_id_seq', 1, true);


--
-- Data for Name: serverstatus; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY serverstatus (id, ilo_pingable, teng_pingable, fqdn_pingable, dscp, firmware, marvin, ping6, upd_pending, stats, prox, mtu, ccr_online, rascal, chr, cdu, ort_errors, mbps_out, clients_connected, server, last_recycle_date, last_recycle_duration_hrs, last_updated) FROM stdin;
\.


--
-- Name: serverstatus_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('serverstatus_id_seq', 1, true);


--
-- Data for Name: staticdnsentry; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY staticdnsentry (id, host, address, type, ttl, deliveryservice, location, last_updated) FROM stdin;
\.


--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('staticdnsentry_id_seq', 1, true);


--
-- Data for Name: status; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY status (id, name, description, last_updated) FROM stdin;
\.


--
-- Name: status_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('status_id_seq', 1, true);


--
-- Data for Name: tm_user; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY tm_user (id, username, role, uid, gid, local_passwd, confirm_local_passwd, last_updated, company, email, full_name, new_user, address_line1, address_line2, city, state_or_province, phone_number, postal_code, country, local_user, token, registration_sent) FROM stdin;
\.


--
-- Name: tm_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('tm_user_id_seq', 1, true);


--
-- Data for Name: to_extension; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY to_extension (id, name, version, info_url, script_file, isactive, additional_config_json, description, servercheck_short_name, servercheck_column_name, type, last_updated) FROM stdin;
\.


--
-- Name: to_extension_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('to_extension_id_seq', 1, true);


--
-- Data for Name: type; Type: TABLE DATA; Schema: public; Owner: jheitz200
--

COPY type (id, name, description, use_in_table, last_updated) FROM stdin;
\.


--
-- Name: type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: jheitz200
--

SELECT pg_catalog.setval('type_id_seq', 1, true);


--
-- Name: idx_31799_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY cran
    ADD CONSTRAINT idx_31799_primary PRIMARY KEY (id, location);


--
-- Name: idx_31808_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT idx_31808_primary PRIMARY KEY (id, type);


--
-- Name: idx_31820_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT idx_31820_primary PRIMARY KEY (deliveryservice, regex);


--
-- Name: idx_31824_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT idx_31824_primary PRIMARY KEY (deliveryservice, server);


--
-- Name: idx_31829_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT idx_31829_primary PRIMARY KEY (deliveryservice, tm_user_id);


--
-- Name: idx_31836_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY division
    ADD CONSTRAINT idx_31836_primary PRIMARY KEY (id);


--
-- Name: idx_31844_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY header_rewrite
    ADD CONSTRAINT idx_31844_primary PRIMARY KEY (id);


--
-- Name: idx_31853_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY hwinfo
    ADD CONSTRAINT idx_31853_primary PRIMARY KEY (id);


--
-- Name: idx_31864_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job
    ADD CONSTRAINT idx_31864_primary PRIMARY KEY (id);


--
-- Name: idx_31875_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_agent
    ADD CONSTRAINT idx_31875_primary PRIMARY KEY (id);


--
-- Name: idx_31887_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_result
    ADD CONSTRAINT idx_31887_primary PRIMARY KEY (id);


--
-- Name: idx_31898_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_status
    ADD CONSTRAINT idx_31898_primary PRIMARY KEY (id);


--
-- Name: idx_31906_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location
    ADD CONSTRAINT idx_31906_primary PRIMARY KEY (id, type);


--
-- Name: idx_31912_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location_parameter
    ADD CONSTRAINT idx_31912_primary PRIMARY KEY (location, parameter);


--
-- Name: idx_31919_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY log
    ADD CONSTRAINT idx_31919_primary PRIMARY KEY (id, tm_user);


--
-- Name: idx_31930_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY parameter
    ADD CONSTRAINT idx_31930_primary PRIMARY KEY (id);


--
-- Name: idx_31941_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY phys_location
    ADD CONSTRAINT idx_31941_primary PRIMARY KEY (id);


--
-- Name: idx_31952_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY profile
    ADD CONSTRAINT idx_31952_primary PRIMARY KEY (id);


--
-- Name: idx_31958_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT idx_31958_primary PRIMARY KEY (profile, parameter);


--
-- Name: idx_31965_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY regex
    ADD CONSTRAINT idx_31965_primary PRIMARY KEY (id, type);


--
-- Name: idx_31974_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY region
    ADD CONSTRAINT idx_31974_primary PRIMARY KEY (id);


--
-- Name: idx_31982_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY role
    ADD CONSTRAINT idx_31982_primary PRIMARY KEY (id);


--
-- Name: idx_31988_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT idx_31988_primary PRIMARY KEY (id, location, type, status, profile);


--
-- Name: idx_32000_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY serverstatus
    ADD CONSTRAINT idx_32000_primary PRIMARY KEY (id, server);


--
-- Name: idx_32015_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT idx_32015_primary PRIMARY KEY (id);


--
-- Name: idx_32024_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY status
    ADD CONSTRAINT idx_32024_primary PRIMARY KEY (id);


--
-- Name: idx_32032_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY tm_user
    ADD CONSTRAINT idx_32032_primary PRIMARY KEY (id);


--
-- Name: idx_32046_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY to_extension
    ADD CONSTRAINT idx_32046_primary PRIMARY KEY (id);


--
-- Name: idx_32056_primary; Type: CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY type
    ADD CONSTRAINT idx_32056_primary PRIMARY KEY (id);


--
-- Name: idx_31799_cr_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31799_cr_id_unique ON cran USING btree (id);


--
-- Name: idx_31799_fk_cran_location1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31799_fk_cran_location1 ON cran USING btree (location);


--
-- Name: idx_31808_ds_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31808_ds_id_unique ON deliveryservice USING btree (id);


--
-- Name: idx_31808_ds_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31808_ds_name_unique ON deliveryservice USING btree (xml_id);


--
-- Name: idx_31808_fk_deliveryservice_header_rewrite1_idx; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31808_fk_deliveryservice_header_rewrite1_idx ON deliveryservice USING btree (header_rewrite);


--
-- Name: idx_31808_fk_deliveryservice_profile1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31808_fk_deliveryservice_profile1 ON deliveryservice USING btree (profile);


--
-- Name: idx_31808_fk_deliveryservice_type1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31808_fk_deliveryservice_type1 ON deliveryservice USING btree (type);


--
-- Name: idx_31820_fk_ds_to_regex_regex1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31820_fk_ds_to_regex_regex1 ON deliveryservice_regex USING btree (regex);


--
-- Name: idx_31824_fk_ds_to_cs_contentserver1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31824_fk_ds_to_cs_contentserver1 ON deliveryservice_server USING btree (server);


--
-- Name: idx_31829_fk_tm_userid; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31829_fk_tm_userid ON deliveryservice_tmuser USING btree (tm_user_id);


--
-- Name: idx_31836_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31836_name_unique ON division USING btree (name);


--
-- Name: idx_31853_fk_hwinfo1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31853_fk_hwinfo1 ON hwinfo USING btree (serverid);


--
-- Name: idx_31853_serverid; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31853_serverid ON hwinfo USING btree (serverid, description);


--
-- Name: idx_31864_fk_job_agent_id1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31864_fk_job_agent_id1 ON job USING btree (agent);


--
-- Name: idx_31864_fk_job_status_id1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31864_fk_job_status_id1 ON job USING btree (status);


--
-- Name: idx_31864_fk_job_user_id1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31864_fk_job_user_id1 ON job USING btree (job_user);


--
-- Name: idx_31887_fk_agent_id1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31887_fk_agent_id1 ON job_result USING btree (agent);


--
-- Name: idx_31887_fk_job_id1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31887_fk_job_id1 ON job_result USING btree (job);


--
-- Name: idx_31906_fk_location_1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31906_fk_location_1 ON location USING btree (parent_location_id);


--
-- Name: idx_31906_fk_location_type1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31906_fk_location_type1 ON location USING btree (type);


--
-- Name: idx_31906_lo_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31906_lo_id_unique ON location USING btree (id);


--
-- Name: idx_31906_loc_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31906_loc_name_unique ON location USING btree (name);


--
-- Name: idx_31906_loc_short_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31906_loc_short_unique ON location USING btree (short_name);


--
-- Name: idx_31912_fk_location; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31912_fk_location ON location_parameter USING btree (location);


--
-- Name: idx_31912_fk_parameter; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31912_fk_parameter ON location_parameter USING btree (parameter);


--
-- Name: idx_31919_fk_log_1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31919_fk_log_1 ON log USING btree (tm_user);


--
-- Name: idx_31941_fk_phys_location_region_idx; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31941_fk_phys_location_region_idx ON phys_location USING btree (region);


--
-- Name: idx_31941_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31941_name_unique ON phys_location USING btree (name);


--
-- Name: idx_31941_short_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31941_short_name_unique ON phys_location USING btree (short_name);


--
-- Name: idx_31952_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31952_name_unique ON profile USING btree (name);


--
-- Name: idx_31958_fk_atsprofile_atsparameters_atsparameters1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31958_fk_atsprofile_atsparameters_atsparameters1 ON profile_parameter USING btree (parameter);


--
-- Name: idx_31958_fk_atsprofile_atsparameters_atsprofile1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31958_fk_atsprofile_atsparameters_atsprofile1 ON profile_parameter USING btree (profile);


--
-- Name: idx_31965_fk_regex_type1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31965_fk_regex_type1 ON regex USING btree (type);


--
-- Name: idx_31965_re_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31965_re_id_unique ON regex USING btree (id);


--
-- Name: idx_31974_fk_region_division1_idx; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31974_fk_region_division1_idx ON region USING btree (division);


--
-- Name: idx_31974_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31974_name_unique ON region USING btree (name);


--
-- Name: idx_31988_cs_ip_address_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31988_cs_ip_address_unique ON server USING btree (ip_address);


--
-- Name: idx_31988_fk_contentserver_atsprofile1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31988_fk_contentserver_atsprofile1 ON server USING btree (profile);


--
-- Name: idx_31988_fk_contentserver_contentserverstatus1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31988_fk_contentserver_contentserverstatus1 ON server USING btree (status);


--
-- Name: idx_31988_fk_contentserver_contentservertype1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31988_fk_contentserver_contentservertype1 ON server USING btree (type);


--
-- Name: idx_31988_fk_contentserver_location; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31988_fk_contentserver_location ON server USING btree (location);


--
-- Name: idx_31988_fk_contentserver_phys_location1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_31988_fk_contentserver_phys_location1 ON server USING btree (phys_location);


--
-- Name: idx_31988_host_name; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31988_host_name ON server USING btree (host_name);


--
-- Name: idx_31988_ip6_address; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31988_ip6_address ON server USING btree (ip6_address);


--
-- Name: idx_31988_se_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_31988_se_id_unique ON server USING btree (id);


--
-- Name: idx_32000_fk_serverstatus_server1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32000_fk_serverstatus_server1 ON serverstatus USING btree (server);


--
-- Name: idx_32000_server; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32000_server ON serverstatus USING btree (server);


--
-- Name: idx_32000_ses_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32000_ses_id_unique ON serverstatus USING btree (id);


--
-- Name: idx_32015_combi_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32015_combi_unique ON staticdnsentry USING btree (host, address, deliveryservice, location);


--
-- Name: idx_32015_fk_staticdnsentry_ds; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32015_fk_staticdnsentry_ds ON staticdnsentry USING btree (deliveryservice);


--
-- Name: idx_32015_fk_staticdnsentry_location; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32015_fk_staticdnsentry_location ON staticdnsentry USING btree (location);


--
-- Name: idx_32015_fk_staticdnsentry_type; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32015_fk_staticdnsentry_type ON staticdnsentry USING btree (type);


--
-- Name: idx_32032_fk_user_1; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32032_fk_user_1 ON tm_user USING btree (role);


--
-- Name: idx_32032_username_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32032_username_unique ON tm_user USING btree (username);


--
-- Name: idx_32046_fk_ext_type_idx; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE INDEX idx_32046_fk_ext_type_idx ON to_extension USING btree (type);


--
-- Name: idx_32046_id_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32046_id_unique ON to_extension USING btree (id);


--
-- Name: idx_32056_name_unique; Type: INDEX; Schema: public; Owner: jheitz200
--

CREATE UNIQUE INDEX idx_32056_name_unique ON type USING btree (name);


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cran FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_tmuser FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON division FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON hwinfo FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_agent FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_result FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON location FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON location_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON log FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON phys_location FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON regex FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON region FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON serverstatus FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON staticdnsentry FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tm_user FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: jheitz200
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON type FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: fk_agent_id1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_result
    ADD CONSTRAINT fk_agent_id1 FOREIGN KEY (agent) REFERENCES job_agent(id) ON DELETE CASCADE;


--
-- Name: fk_atsprofile_atsparameters_atsparameters1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT fk_atsprofile_atsparameters_atsparameters1 FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;


--
-- Name: fk_atsprofile_atsparameters_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT fk_atsprofile_atsparameters_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id) ON DELETE CASCADE;


--
-- Name: fk_contentserver_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id);


--
-- Name: fk_contentserver_contentserverstatus1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_contentserverstatus1 FOREIGN KEY (status) REFERENCES status(id);


--
-- Name: fk_contentserver_contentservertype1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_contentservertype1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_contentserver_location; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_location FOREIGN KEY (location) REFERENCES location(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: fk_contentserver_phys_location1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_phys_location1 FOREIGN KEY (phys_location) REFERENCES phys_location(id);


--
-- Name: fk_cran_location1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY cran
    ADD CONSTRAINT fk_cran_location1 FOREIGN KEY (location) REFERENCES location(id);


--
-- Name: fk_deliveryservice_header_rewrite1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_deliveryservice_header_rewrite1 FOREIGN KEY (header_rewrite) REFERENCES header_rewrite(id);


--
-- Name: fk_deliveryservice_profile1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_deliveryservice_profile1 FOREIGN KEY (profile) REFERENCES profile(id);


--
-- Name: fk_deliveryservice_type1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_deliveryservice_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_ds_to_cs_contentserver1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT fk_ds_to_cs_contentserver1 FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_cs_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT fk_ds_to_cs_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_regex_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT fk_ds_to_regex_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_regex_regex1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT fk_ds_to_regex_regex1 FOREIGN KEY (regex) REFERENCES regex(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ext_type; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY to_extension
    ADD CONSTRAINT fk_ext_type FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_hwinfo1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY hwinfo
    ADD CONSTRAINT fk_hwinfo1 FOREIGN KEY (serverid) REFERENCES server(id) ON DELETE CASCADE;


--
-- Name: fk_job_agent_id1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_agent_id1 FOREIGN KEY (agent) REFERENCES job_agent(id) ON DELETE CASCADE;


--
-- Name: fk_job_id1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job_result
    ADD CONSTRAINT fk_job_id1 FOREIGN KEY (job) REFERENCES job(id) ON DELETE CASCADE;


--
-- Name: fk_job_status_id1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_status_id1 FOREIGN KEY (status) REFERENCES job_status(id);


--
-- Name: fk_job_user_id1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_user_id1 FOREIGN KEY (job_user) REFERENCES tm_user(id);


--
-- Name: fk_location; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location_parameter
    ADD CONSTRAINT fk_location FOREIGN KEY (location) REFERENCES location(id) ON DELETE CASCADE;


--
-- Name: fk_location_1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location
    ADD CONSTRAINT fk_location_1 FOREIGN KEY (parent_location_id) REFERENCES location(id);


--
-- Name: fk_location_type1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location
    ADD CONSTRAINT fk_location_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_log_1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY log
    ADD CONSTRAINT fk_log_1 FOREIGN KEY (tm_user) REFERENCES tm_user(id);


--
-- Name: fk_parameter; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY location_parameter
    ADD CONSTRAINT fk_parameter FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;


--
-- Name: fk_phys_location_region; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY phys_location
    ADD CONSTRAINT fk_phys_location_region FOREIGN KEY (region) REFERENCES region(id);


--
-- Name: fk_regex_type1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY regex
    ADD CONSTRAINT fk_regex_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_region_division1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY region
    ADD CONSTRAINT fk_region_division1 FOREIGN KEY (division) REFERENCES division(id);


--
-- Name: fk_serverstatus_server1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY serverstatus
    ADD CONSTRAINT fk_serverstatus_server1 FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE;


--
-- Name: fk_staticdnsentry_ds; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_ds FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id);


--
-- Name: fk_staticdnsentry_location; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_location FOREIGN KEY (location) REFERENCES location(id);


--
-- Name: fk_staticdnsentry_type; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_type FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_tm_user_ds; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT fk_tm_user_ds FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_tm_user_id; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT fk_tm_user_id FOREIGN KEY (tm_user_id) REFERENCES tm_user(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_user_1; Type: FK CONSTRAINT; Schema: public; Owner: jheitz200
--

ALTER TABLE ONLY tm_user
    ADD CONSTRAINT fk_user_1 FOREIGN KEY (role) REFERENCES role(id) ON DELETE SET NULL;


--
-- Name: public; Type: ACL; Schema: -; Owner: jheitz200
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM jheitz200;
GRANT ALL ON SCHEMA public TO jheitz200;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--
