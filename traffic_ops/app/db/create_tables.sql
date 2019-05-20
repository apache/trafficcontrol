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

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.5

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


SET search_path = public, pg_catalog;

--
-- Name: on_update_current_timestamp_last_updated(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE FUNCTION on_update_current_timestamp_last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.last_updated = now();
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.on_update_current_timestamp_last_updated() OWNER TO traffic_ops;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: change_types; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE change_types AS ENUM (
    'create',
    'update',
    'delete'
);

--
-- Name: deep_caching_type; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE deep_caching_type AS ENUM (
    'NEVER',
    'ALWAYS'
);

--
-- Name: http_method_t; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE http_method_t AS ENUM (
    'GET',
    'PUT',
    'POST',
    'PATCH',
    'DELETE'
);

--
-- Name: localization_method; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE localization_method AS ENUM (
    'CZ',
    'DEEP_CZ',
    'GEO'
);

--
-- Name: origin_protocol; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE origin_protocol AS ENUM (
    'http',
    'https'
);

--
-- Name: profile_type; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE profile_type AS ENUM (
    'ATS_PROFILE',
    'TR_PROFILE',
    'TM_PROFILE',
    'TS_PROFILE',
    'TP_PROFILE',
    'INFLUXDB_PROFILE',
    'RIAK_PROFILE',
    'SPLUNK_PROFILE',
    'DS_PROFILE',
    'ORG_PROFILE',
    'KAFKA_PROFILE',
    'LOGSTASH_PROFILE',
    'ES_PROFILE',
    'UNK_PROFILE',
    'GROVE_PROFILE'
);

--
-- Name: workflow_states; Type: TYPE; Schema: public; Owner: traffic_ops
--

CREATE TYPE workflow_states AS ENUM (
    'draft',
    'submitted',
    'rejected',
    'pending',
    'complete'
);

--
-- Name: deliveryservice_signature_type; Type: DOMAIN; Schema: public; Owner: traffic_ops
--

CREATE DOMAIN deliveryservice_signature_type AS text CHECK (VALUE IN ('url_sig', 'uri_signing'));

--
-- Name: api_capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE api_capability (
    id bigserial NOT NULL,
    http_method http_method_t,
    route text,
    capability text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

ALTER TABLE api_capability OWNER TO traffic_ops;

--
-- Name: asn; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE asn (
    id bigint NOT NULL,
    asn bigint NOT NULL,
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE asn OWNER TO traffic_ops;

--
-- Name: asn_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE asn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE asn_id_seq OWNER TO traffic_ops;

--
-- Name: asn_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE asn_id_seq OWNED BY asn.id;

--
-- Name: cachegroup; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cachegroup (
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


ALTER TABLE cachegroup OWNER TO traffic_ops;

--
-- Name: cachegroup_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE cachegroup_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cachegroup_id_seq OWNER TO traffic_ops;

--
-- Name: cachegroup_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cachegroup_id_seq OWNED BY cachegroup.id;

--
-- Name: cachegroup_fallbacks; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cachegroup_fallbacks (
    primary_cg bigint,
    backup_cg bigint CHECK (primary_cg != backup_cg),
    set_order bigint NOT NULL
);

ALTER TABLE cachegroup_fallbacks OWNER TO traffic_ops;

--
-- Name: cachegroup_localization_method; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cachegroup_localization_method (
    cachegroup bigint,
    method localization_method
);

--
-- Name: cachegroup_parameter; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cachegroup_parameter (
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE cachegroup_parameter OWNER TO traffic_ops;

--
-- Name: capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE capability (
    name text,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

ALTER TABLE capability OWNER TO traffic_ops;

--
-- Name: cdn; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdn (
    id bigint,
    name text UNIQUE NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    dnssec_enabled boolean DEFAULT false NOT NULL,
    domain_name text UNIQUE NOT NULL
);


ALTER TABLE cdn OWNER TO traffic_ops;

--
-- Name: cdn_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE cdn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cdn_id_seq OWNER TO traffic_ops;

--
-- Name: cdn_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cdn_id_seq OWNED BY cdn.id;

--
-- Name: coordinate; Type: TABLE; Schema: public: Owner: traffic_ops
--

CREATE TABLE coordinate (
    id bigserial UNIQUE NOT NULL,
    name text,
    latitude numeric NOT NULL DEFAULT 0.0,
    longitude numeric NOT NULL DEFAULT 0.0,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

--
-- Name: deliveryservice; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice (
    id bigint,
    xml_id text NOT NULL,
    active boolean DEFAULT false NOT NULL,
    dscp bigint NOT NULL,
    signing_algorithm deliveryservice_signature_type DEFAULT NULL,
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
    max_dns_answers bigint DEFAULT 5,
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
    routing_name text NOT NULL DEFAULT 'cdn' CHECK (length(routing_name) > 0),
    deep_caching_type deep_caching_type NOT NULL DEFAULT 'NEVER',
    fq_pacing_rate bigint DEFAULT 0,
    anonymous_blocking_enabled boolean NOT NULL DEFAULT FALSE,
    consistent_hash_regex text
);


ALTER TABLE deliveryservice OWNER TO traffic_ops;

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE deliveryservice_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE deliveryservice_id_seq OWNER TO traffic_ops;

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE deliveryservice_id_seq OWNED BY deliveryservice.id;


--
-- Name: deliveryservice_regex; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice_regex (
    deliveryservice bigint NOT NULL,
    regex bigint NOT NULL,
    set_number bigint DEFAULT '0'::bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE deliveryservice_regex OWNER TO traffic_ops;

--
-- Name: deliveryservice_request; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice_request (
    assignee_id bigint,
    author_id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    id bigserial,
    last_edited_by bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    deliveryservice jsonb NOT NULL,
    status workflow_states NOT NULL
);

ALTER TABLE deliveryservice_request OWNER TO traffic_ops;

--
-- Name: deliveryservice_request_comment; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice_request_comment (
    author_id bigint NOT NULL,
    deliveryservice_request_id bigint NOT NULL,
    id bigserial,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    value text NOT NULL
);

--
-- Name: deliveryservice_server; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice_server (
    deliveryservice bigint NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now()
);


ALTER TABLE deliveryservice_server OWNER TO traffic_ops;

--
-- Name: deliveryservice_tmuser; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE deliveryservice_tmuser (
    deliveryservice bigint NOT NULL,
    tm_user_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE deliveryservice_tmuser OWNER TO traffic_ops;

--
-- Name: division; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE division (
    id bigint NOT NULL,
    name text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE division OWNER TO traffic_ops;

--
-- Name: division_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE division_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE division_id_seq OWNER TO traffic_ops;

--
-- Name: division_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE division_id_seq OWNED BY division.id;


--
-- Name: federation; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE federation (
    id bigint NOT NULL,
    cname text NOT NULL,
    description text,
    ttl integer NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE federation OWNER TO traffic_ops;

--
-- Name: federation_deliveryservice; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE federation_deliveryservice (
    federation bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE federation_deliveryservice OWNER TO traffic_ops;

--
-- Name: federation_federation_resolver; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE federation_federation_resolver (
    federation bigint NOT NULL,
    federation_resolver bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE federation_federation_resolver OWNER TO traffic_ops;

--
-- Name: federation_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE federation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federation_id_seq OWNER TO traffic_ops;

--
-- Name: federation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE federation_id_seq OWNED BY federation.id;


--
-- Name: federation_resolver; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE federation_resolver (
    id bigint NOT NULL,
    ip_address text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE federation_resolver OWNER TO traffic_ops;

--
-- Name: federation_resolver_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE federation_resolver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE federation_resolver_id_seq OWNER TO traffic_ops;

--
-- Name: federation_resolver_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE federation_resolver_id_seq OWNED BY federation_resolver.id;


--
-- Name: federation_tmuser; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE federation_tmuser (
    federation bigint NOT NULL,
    tm_user bigint NOT NULL,
    role bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE federation_tmuser OWNER TO traffic_ops;

--
-- Name: hwinfo; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE hwinfo (
    id bigint NOT NULL,
    serverid bigint NOT NULL,
    description text NOT NULL,
    val text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE hwinfo OWNER TO traffic_ops;

--
-- Name: hwinfo_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE hwinfo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE hwinfo_id_seq OWNER TO traffic_ops;

--
-- Name: hwinfo_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE hwinfo_id_seq OWNED BY hwinfo.id;


--
-- Name: job; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE job (
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


ALTER TABLE job OWNER TO traffic_ops;

--
-- Name: job_agent; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE job_agent (
    id bigint NOT NULL,
    name text UNIQUE,
    description text,
    active integer DEFAULT 0 NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE job_agent OWNER TO traffic_ops;

--
-- Name: job_agent_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE job_agent_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_agent_id_seq OWNER TO traffic_ops;

--
-- Name: job_agent_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE job_agent_id_seq OWNED BY job_agent.id;


--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE job_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_id_seq OWNER TO traffic_ops;

--
-- Name: job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE job_id_seq OWNED BY job.id;

--
-- Name: job_status; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE job_status (
    id bigint NOT NULL,
    name text UNIQUE,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE job_status OWNER TO traffic_ops;

--
-- Name: job_status_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE job_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE job_status_id_seq OWNER TO traffic_ops;

--
-- Name: job_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE job_status_id_seq OWNED BY job_status.id;


--
-- Name: log; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE log (
    id bigint NOT NULL,
    level text,
    message text NOT NULL,
    tm_user bigint NOT NULL,
    ticketnum text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE log OWNER TO traffic_ops;

--
-- Name: log_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE log_id_seq OWNER TO traffic_ops;

--
-- Name: log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE log_id_seq OWNED BY log.id;

--
-- Name: origin; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE origin (
    id bigserial UNIQUE NOT NULL,
    name text,
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

ALTER TABLE origin OWNER TO traffic_ops;

--
-- Name: parameter; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE parameter (
    id bigint NOT NULL,
    name text NOT NULL,
    config_file text,
    value text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    secure boolean DEFAULT false NOT NULL
);


ALTER TABLE parameter OWNER TO traffic_ops;

--
-- Name: parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE parameter_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE parameter_id_seq OWNER TO traffic_ops;

--
-- Name: parameter_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE parameter_id_seq OWNED BY parameter.id;


--
-- Name: phys_location; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE phys_location (
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


ALTER TABLE phys_location OWNER TO traffic_ops;

--
-- Name: phys_location_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE phys_location_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE phys_location_id_seq OWNER TO traffic_ops;

--
-- Name: phys_location_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE phys_location_id_seq OWNED BY phys_location.id;

--
-- Name: profile; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE profile (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    type profile_type NOT NULL,
    cdn bigint NOT NULL,
    routing_disabled boolean NOT NULL DEFAULT FALSE
);


ALTER TABLE profile OWNER TO traffic_ops;

--
-- Name: profile_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE profile_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE profile_id_seq OWNER TO traffic_ops;

--
-- Name: profile_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE profile_id_seq OWNED BY profile.id;


--
-- Name: profile_parameter; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE profile_parameter (
    profile bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE profile_parameter OWNER TO traffic_ops;

--
-- Name: regex; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE regex (
    id bigint NOT NULL,
    pattern text DEFAULT ''::text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE regex OWNER TO traffic_ops;

--
-- Name: regex_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE regex_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE regex_id_seq OWNER TO traffic_ops;

--
-- Name: regex_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE regex_id_seq OWNED BY regex.id;


--
-- Name: region; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE region (
    id bigint NOT NULL,
    name text NOT NULL,
    division bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE region OWNER TO traffic_ops;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE region_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE region_id_seq OWNER TO traffic_ops;

--
-- Name: region_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE region_id_seq OWNED BY region.id;


--
-- Name: role; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE role (
    id bigint,
    name text UNIQUE NOT NULL,
    description text,
    priv_level bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE role OWNER TO traffic_ops;

--
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE role_id_seq OWNER TO traffic_ops;

--
-- Name: role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE role_id_seq OWNED BY role.id;

--
-- Name: role_capability; Type TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE role_capability (
    role_id bigint,
    cap_name text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

ALTER TABLE role_capability OWNER TO traffic_ops;

--
-- Name: server; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE server (
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


ALTER TABLE server OWNER TO traffic_ops;

--
-- Name: server_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE server_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE server_id_seq OWNER TO traffic_ops;

--
-- Name: server_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE server_id_seq OWNED BY server.id;


--
-- Name: servercheck; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE servercheck (
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


ALTER TABLE servercheck OWNER TO traffic_ops;

--
-- Name: servercheck_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE servercheck_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE servercheck_id_seq OWNER TO traffic_ops;

--
-- Name: servercheck_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE servercheck_id_seq OWNED BY servercheck.id;

--
-- Name: snapshot; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE snapshot (
    cdn text NOT NULL,
    crconfig json NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    monitoring json NOT NULL
);

ALTER TABLE snapshot OWNER TO traffic_ops;

--
-- Name: staticdnsentry; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE staticdnsentry (
    id bigint NOT NULL,
    host text NOT NULL,
    address text NOT NULL,
    type bigint NOT NULL,
    ttl bigint DEFAULT '3600'::bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    cachegroup bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE staticdnsentry OWNER TO traffic_ops;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE staticdnsentry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE staticdnsentry_id_seq OWNER TO traffic_ops;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE staticdnsentry_id_seq OWNED BY staticdnsentry.id;


--
-- Name: stats_summary; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE stats_summary (
    id bigint NOT NULL,
    cdn_name text DEFAULT 'all'::text NOT NULL,
    deliveryservice_name text NOT NULL,
    stat_name text NOT NULL,
    stat_value double precision NOT NULL,
    summary_time timestamp with time zone DEFAULT now() NOT NULL,
    stat_date date
);


ALTER TABLE stats_summary OWNER TO traffic_ops;

--
-- Name: stats_summary_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE stats_summary_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE stats_summary_id_seq OWNER TO traffic_ops;

--
-- Name: stats_summary_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE stats_summary_id_seq OWNED BY stats_summary.id;


--
-- Name: status; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE status (
    id bigint NOT NULL,
    name text UNIQUE NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE status OWNER TO traffic_ops;

--
-- Name: status_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE status_id_seq OWNER TO traffic_ops;

--
-- Name: status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE status_id_seq OWNED BY status.id;


--
-- Name: steering_target; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE steering_target (
    deliveryservice bigint NOT NULL,
    target bigint NOT NULL,
    value bigint NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE steering_target OWNER TO traffic_ops;

--
-- Name: tenant; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE tenant (
    id bigserial,
    name text UNIQUE NOT NULL,
    active boolean NOT NULL DEFAULT FALSE,
    parent_id bigint DEFAULT 1 CHECK (id != parent_id),
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

--
-- Name: tm_user; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE tm_user (
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


ALTER TABLE tm_user OWNER TO traffic_ops;

--
-- Name: tm_user_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE tm_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE tm_user_id_seq OWNER TO traffic_ops;

--
-- Name: tm_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE tm_user_id_seq OWNED BY tm_user.id;


--
-- Name: to_extension; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE to_extension (
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


ALTER TABLE to_extension OWNER TO traffic_ops;

--
-- Name: to_extension_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE to_extension_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE to_extension_id_seq OWNER TO traffic_ops;

--
-- Name: to_extension_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE to_extension_id_seq OWNED BY to_extension.id;


--
-- Name: type; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE type (
    id bigint NOT NULL,
    name text UNIQUE NOT NULL,
    description text,
    use_in_table text,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);


ALTER TABLE type OWNER TO traffic_ops;

--
-- Name: type_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE type_id_seq OWNER TO traffic_ops;

--
-- Name: type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE type_id_seq OWNED BY type.id;

--
-- Name: user_role; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE user_role (
    user_id bigint NOT NULL,
    role_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

ALTER TABLE user_role OWNER TO traffic_ops;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY asn ALTER COLUMN id SET DEFAULT nextval('asn_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup ALTER COLUMN id SET DEFAULT nextval('cachegroup_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdn ALTER COLUMN id SET DEFAULT nextval('cdn_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice ALTER COLUMN id SET DEFAULT nextval('deliveryservice_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY division ALTER COLUMN id SET DEFAULT nextval('division_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation ALTER COLUMN id SET DEFAULT nextval('federation_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_resolver ALTER COLUMN id SET DEFAULT nextval('federation_resolver_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY hwinfo ALTER COLUMN id SET DEFAULT nextval('hwinfo_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job ALTER COLUMN id SET DEFAULT nextval('job_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job_agent ALTER COLUMN id SET DEFAULT nextval('job_agent_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job_status ALTER COLUMN id SET DEFAULT nextval('job_status_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY log ALTER COLUMN id SET DEFAULT nextval('log_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY parameter ALTER COLUMN id SET DEFAULT nextval('parameter_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY phys_location ALTER COLUMN id SET DEFAULT nextval('phys_location_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile ALTER COLUMN id SET DEFAULT nextval('profile_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY regex ALTER COLUMN id SET DEFAULT nextval('regex_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY region ALTER COLUMN id SET DEFAULT nextval('region_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY role ALTER COLUMN id SET DEFAULT nextval('role_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server ALTER COLUMN id SET DEFAULT nextval('server_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY servercheck ALTER COLUMN id SET DEFAULT nextval('servercheck_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY staticdnsentry ALTER COLUMN id SET DEFAULT nextval('staticdnsentry_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY stats_summary ALTER COLUMN id SET DEFAULT nextval('stats_summary_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY status ALTER COLUMN id SET DEFAULT nextval('status_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tm_user ALTER COLUMN id SET DEFAULT nextval('tm_user_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY to_extension ALTER COLUMN id SET DEFAULT nextval('to_extension_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY type ALTER COLUMN id SET DEFAULT nextval('type_id_seq'::regclass);

--
-- Name: pk_api_capability; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY api_capability
    ADD CONSTRAINT pk_api_capability PRIMARY KEY (http_method, route, capability);

--
-- Name: idx_89468_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY asn
    ADD CONSTRAINT idx_89468_primary PRIMARY KEY (id, cachegroup);


--
-- Name: idx_89476_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT idx_89476_primary PRIMARY KEY (id, type);

--
-- Name: pk_cachegroup_fallbacks; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_fallbacks
    ADD CONSTRAINT pk_cachegroup_fallbacks PRIMARY KEY (primary_cg, backup_cg);

--
-- Name: idx_89484_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_parameter
    ADD CONSTRAINT idx_89484_primary PRIMARY KEY (cachegroup, parameter);

--
-- Name: pk_capability; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY capability
    ADD CONSTRAINT pk_capability PRIMARY KEY (name);

--
-- Name: idx_89491_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdn
    ADD CONSTRAINT idx_89491_primary PRIMARY KEY (id);

--
-- Name: pk_coordinate; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY coordinate
    ADD CONSTRAINT pk_coordinate PRIMARY KEY (name);

--
-- Name: idx_89502_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT idx_89502_primary PRIMARY KEY (id, type);


--
-- Name: idx_89517_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT idx_89517_primary PRIMARY KEY (deliveryservice, regex);

--
-- Name: pk_deliveryservice_request; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request
    ADD CONSTRAINT pk_deliveryservice_request PRIMARY KEY (id);

--
-- Name: pk_deliveryservice_request_comment; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request_comment
    ADD CONSTRAINT pk_deliveryservice_request_comment PRIMARY KEY (id);

--
-- Name: idx_89521_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT idx_89521_primary PRIMARY KEY (deliveryservice, server);


--
-- Name: idx_89525_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT idx_89525_primary PRIMARY KEY (deliveryservice, tm_user_id);


--
-- Name: idx_89531_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY division
    ADD CONSTRAINT idx_89531_primary PRIMARY KEY (id);


--
-- Name: idx_89541_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation
    ADD CONSTRAINT idx_89541_primary PRIMARY KEY (id);


--
-- Name: idx_89549_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_deliveryservice
    ADD CONSTRAINT idx_89549_primary PRIMARY KEY (federation, deliveryservice);


--
-- Name: idx_89553_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_federation_resolver
    ADD CONSTRAINT idx_89553_primary PRIMARY KEY (federation, federation_resolver);


--
-- Name: idx_89559_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_resolver
    ADD CONSTRAINT idx_89559_primary PRIMARY KEY (id);


--
-- Name: idx_89567_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_tmuser
    ADD CONSTRAINT idx_89567_primary PRIMARY KEY (federation, tm_user);


--
-- Name: idx_89583_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY hwinfo
    ADD CONSTRAINT idx_89583_primary PRIMARY KEY (id);


--
-- Name: idx_89593_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job
    ADD CONSTRAINT idx_89593_primary PRIMARY KEY (id);


--
-- Name: idx_89603_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job_agent
    ADD CONSTRAINT idx_89603_primary PRIMARY KEY (id);


--
-- Name: idx_89624_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job_status
    ADD CONSTRAINT idx_89624_primary PRIMARY KEY (id);

--
-- Name: pk_cachegroup_localization_method; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_localization_method
    ADD CONSTRAINT pk_cachegroup_localization_method PRIMARY KEY (cachegroup, method);

--
-- Name: idx_89634_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY log
    ADD CONSTRAINT idx_89634_primary PRIMARY KEY (id, tm_user);

--
-- Name: pk_origin; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT pk_origin PRIMARY KEY (name);

--
-- Name: idx_89644_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY parameter
    ADD CONSTRAINT idx_89644_primary PRIMARY KEY (id);


--
-- Name: idx_89655_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY phys_location
    ADD CONSTRAINT idx_89655_primary PRIMARY KEY (id);


--
-- Name: idx_89665_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile
    ADD CONSTRAINT idx_89665_primary PRIMARY KEY (id);


--
-- Name: idx_89673_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT idx_89673_primary PRIMARY KEY (profile, parameter);


--
-- Name: idx_89679_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY regex
    ADD CONSTRAINT idx_89679_primary PRIMARY KEY (id, type);


--
-- Name: idx_89690_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY region
    ADD CONSTRAINT idx_89690_primary PRIMARY KEY (id);


--
-- Name: idx_89700_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY role
    ADD CONSTRAINT idx_89700_primary PRIMARY KEY (id);

--
-- Name: pk_role_capability; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY role_capability
    ADD CONSTRAINT pk_role_capability PRIMARY KEY (role_id, cap_name);

--
-- Name: idx_89709_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT idx_89709_primary PRIMARY KEY (id, cachegroup, type, status, profile);


--
-- Name: idx_89722_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY servercheck
    ADD CONSTRAINT idx_89722_primary PRIMARY KEY (id, server);

--
-- Name: pk_snapshot; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY snapshot
    ADD CONSTRAINT pk_snapshot PRIMARY KEY (cdn);

--
-- Name: idx_89729_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT idx_89729_primary PRIMARY KEY (id);


--
-- Name: idx_89740_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY stats_summary
    ADD CONSTRAINT idx_89740_primary PRIMARY KEY (id);


--
-- Name: idx_89751_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY status
    ADD CONSTRAINT idx_89751_primary PRIMARY KEY (id);


--
-- Name: idx_89759_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY steering_target
    ADD CONSTRAINT idx_89759_primary PRIMARY KEY (deliveryservice, target);

--
-- Name: pk_tenant; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tenant
    ADD CONSTRAINT pk_tenant PRIMARY KEY (id);

--
-- Name: idx_89765_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tm_user
    ADD CONSTRAINT idx_89765_primary PRIMARY KEY (id);


--
-- Name: idx_89776_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY to_extension
    ADD CONSTRAINT idx_89776_primary PRIMARY KEY (id);


--
-- Name: idx_89786_primary; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY type
    ADD CONSTRAINT idx_89786_primary PRIMARY KEY (id);

--
-- Name: pk_user_role; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY user_role
    ADD CONSTRAINT pk_user_role PRIMARY KEY (user_id, role_id);

--
-- Name: idx_89468_cr_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89468_cr_id_unique ON asn USING btree (id);


--
-- Name: idx_89468_fk_cran_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89468_fk_cran_cachegroup1 ON asn USING btree (cachegroup);


--
-- Name: idx_89476_cg_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89476_cg_name_unique ON cachegroup USING btree (name);


--
-- Name: idx_89476_cg_short_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89476_cg_short_unique ON cachegroup USING btree (short_name);


--
-- Name: idx_89476_fk_cg_1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89476_fk_cg_1 ON cachegroup USING btree (parent_cachegroup_id);


--
-- Name: idx_89476_fk_cg_secondary; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89476_fk_cg_secondary ON cachegroup USING btree (secondary_parent_cachegroup_id);


--
-- Name: idx_89476_fk_cg_type1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89476_fk_cg_type1 ON cachegroup USING btree (type);


--
-- Name: idx_89476_lo_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89476_lo_id_unique ON cachegroup USING btree (id);


--
-- Name: idx_89484_fk_parameter; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89484_fk_parameter ON cachegroup_parameter USING btree (parameter);


--
-- Name: idx_89491_cdn_cdn_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89491_cdn_cdn_unique ON cdn USING btree (name);

--
-- Name: idx_89502_ds_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89502_ds_id_unique ON deliveryservice USING btree (id);

--
-- Name: idx_k_deliveryservice_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_k_deliveryservice_tenant_idx ON deliveryservice USING btree (tenant_id);

--
-- Name: idx_89502_ds_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89502_ds_name_unique ON deliveryservice USING btree (xml_id);


--
-- Name: idx_89502_fk_cdn1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89502_fk_cdn1 ON deliveryservice USING btree (cdn_id);


--
-- Name: idx_89502_fk_deliveryservice_profile1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89502_fk_deliveryservice_profile1 ON deliveryservice USING btree (profile);


--
-- Name: idx_89502_fk_deliveryservice_type1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89502_fk_deliveryservice_type1 ON deliveryservice USING btree (type);


--
-- Name: idx_89517_fk_ds_to_regex_regex1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89517_fk_ds_to_regex_regex1 ON deliveryservice_regex USING btree (regex);


--
-- Name: idx_89521_fk_ds_to_cs_contentserver1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89521_fk_ds_to_cs_contentserver1 ON deliveryservice_server USING btree (server);


--
-- Name: idx_89525_fk_tm_userid; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89525_fk_tm_userid ON deliveryservice_tmuser USING btree (tm_user_id);


--
-- Name: idx_89531_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89531_name_unique ON division USING btree (name);


--
-- Name: idx_89549_fk_fed_to_ds1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89549_fk_fed_to_ds1 ON federation_deliveryservice USING btree (deliveryservice);


--
-- Name: idx_89553_fk_federation_federation_resolver; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89553_fk_federation_federation_resolver ON federation_federation_resolver USING btree (federation);


--
-- Name: idx_89553_fk_federation_resolver_to_fed1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89553_fk_federation_resolver_to_fed1 ON federation_federation_resolver USING btree (federation_resolver);


--
-- Name: idx_89559_federation_resolver_ip_address; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89559_federation_resolver_ip_address ON federation_resolver USING btree (ip_address);


--
-- Name: idx_89559_fk_federation_mapping_type; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89559_fk_federation_mapping_type ON federation_resolver USING btree (type);


--
-- Name: idx_89567_fk_federation_federation_resolver; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89567_fk_federation_federation_resolver ON federation_tmuser USING btree (federation);


--
-- Name: idx_89567_fk_federation_tmuser_role; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89567_fk_federation_tmuser_role ON federation_tmuser USING btree (role);


--
-- Name: idx_89567_fk_federation_tmuser_tmuser; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89567_fk_federation_tmuser_tmuser ON federation_tmuser USING btree (tm_user);


--
-- Name: idx_89583_fk_hwinfo1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89583_fk_hwinfo1 ON hwinfo USING btree (serverid);


--
-- Name: idx_89583_serverid; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89583_serverid ON hwinfo USING btree (serverid, description);


--
-- Name: idx_89593_fk_job_agent_id1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89593_fk_job_agent_id1 ON job USING btree (agent);


--
-- Name: idx_89593_fk_job_deliveryservice1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89593_fk_job_deliveryservice1 ON job USING btree (job_deliveryservice);


--
-- Name: idx_89593_fk_job_status_id1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89593_fk_job_status_id1 ON job USING btree (status);


--
-- Name: idx_89593_fk_job_user_id1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89593_fk_job_user_id1 ON job USING btree (job_user);


--
-- Name: idx_89634_fk_log_1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89634_fk_log_1 ON log USING btree (tm_user);


--
-- Name: idx_89634_idx_last_updated; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89634_idx_last_updated ON log USING btree (last_updated);


--
-- Name: idx_89644_parameter_name_value_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89644_parameter_name_value_idx ON parameter USING btree (name, value);


--
-- Name: idx_89655_fk_phys_location_region_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89655_fk_phys_location_region_idx ON phys_location USING btree (region);


--
-- Name: idx_89655_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89655_name_unique ON phys_location USING btree (name);


--
-- Name: idx_89655_short_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89655_short_name_unique ON phys_location USING btree (short_name);


--
-- Name: idx_89665_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89665_name_unique ON profile USING btree (name);

--
-- Name: idx_181818_fk_cdn1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_181818_fk_cdn1 ON profile USING btree (id);


--
-- Name: idx_89673_fk_atsprofile_atsparameters_atsparameters1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89673_fk_atsprofile_atsparameters_atsparameters1 ON profile_parameter USING btree (parameter);


--
-- Name: idx_89673_fk_atsprofile_atsparameters_atsprofile1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89673_fk_atsprofile_atsparameters_atsprofile1 ON profile_parameter USING btree (profile);


--
-- Name: idx_89679_fk_regex_type1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89679_fk_regex_type1 ON regex USING btree (type);


--
-- Name: idx_89679_re_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89679_re_id_unique ON regex USING btree (id);


--
-- Name: idx_89690_fk_region_division1_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89690_fk_region_division1_idx ON region USING btree (division);


--
-- Name: idx_89690_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89690_name_unique ON region USING btree (name);


--
-- Name: idx_89709_fk_cdn2; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_cdn2 ON server USING btree (cdn_id);


--
-- Name: idx_89709_fk_contentserver_atsprofile1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_contentserver_atsprofile1 ON server USING btree (profile);


--
-- Name: idx_89709_fk_contentserver_contentserverstatus1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_contentserver_contentserverstatus1 ON server USING btree (status);


--
-- Name: idx_89709_fk_contentserver_contentservertype1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_contentserver_contentservertype1 ON server USING btree (type);


--
-- Name: idx_89709_fk_contentserver_phys_location1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_contentserver_phys_location1 ON server USING btree (phys_location);


--
-- Name: idx_89709_fk_server_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89709_fk_server_cachegroup1 ON server USING btree (cachegroup);


--
-- Name: idx_89709_ip6_profile; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89709_ip6_profile ON server USING btree (ip6_address, profile);


--
-- Name: idx_89709_ip_profile; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89709_ip_profile ON server USING btree (ip_address, profile);


--
-- Name: idx_89709_se_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89709_se_id_unique ON server USING btree (id);


--
-- Name: idx_89722_fk_serverstatus_server1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89722_fk_serverstatus_server1 ON servercheck USING btree (server);


--
-- Name: idx_89722_server; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89722_server ON servercheck USING btree (server);


--
-- Name: idx_89722_ses_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89722_ses_id_unique ON servercheck USING btree (id);


--
-- Name: idx_89729_combi_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89729_combi_unique ON staticdnsentry USING btree (host, address, deliveryservice, cachegroup);


--
-- Name: idx_89729_fk_staticdnsentry_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89729_fk_staticdnsentry_cachegroup1 ON staticdnsentry USING btree (cachegroup);


--
-- Name: idx_89729_fk_staticdnsentry_ds; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89729_fk_staticdnsentry_ds ON staticdnsentry USING btree (deliveryservice);


--
-- Name: idx_89729_fk_staticdnsentry_type; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89729_fk_staticdnsentry_type ON staticdnsentry USING btree (type);

--
-- Name: idx_k_tenant_parent_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_k_tenant_parent_tenant_idx ON tenant USING btree (parent_id);

--
-- Name: idx_89765_fk_user_1; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89765_fk_user_1 ON tm_user USING btree (role);

--
-- Name: idx_k_tm_user_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_k_tm_user_tenant_idx ON tm_user USING btree (tenant_id);

--
-- Name: idx_89765_tmuser_email_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89765_tmuser_email_unique ON tm_user USING btree (email);


--
-- Name: idx_89765_username_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89765_username_unique ON tm_user USING btree (username);


--
-- Name: idx_89776_fk_ext_type_idx; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX idx_89776_fk_ext_type_idx ON to_extension USING btree (type);


--
-- Name: idx_89776_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX idx_89776_id_unique ON to_extension USING btree (id);

--
-- Name: cachegroup_coordinate_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX cachegroup_coordinate_fkey ON cachegroup USING btree (coordinate);

--
-- Name: cachegroup_localization_method_cachegroup_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX cachegroup_localization_method_cachegroup_fkey ON cachegroup_localization_method USING btree (cachegroup);

--
-- Name: origin_is_primary_deliveryservice_constraint; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE UNIQUE INDEX origin_is_primary_deliveryservice_constraint ON origin (is_primary, deliveryservice) WHERE is_primary;

--
-- Name: origin_deliveryservice_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX origin_deliveryservice_fkey ON origin USING btree (deliveryservice);

--
-- Name: origin_coordinate_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX origin_coordinate_fkey ON origin USING btree (coordinate);

--
-- Name: origin_profile_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX origin_profile_fkey ON origin USING btree (profile);

--
-- Name: origin_cachegroup_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX origin_cachegroup_fkey ON origin USING btree (cachegroup);

--
-- Name: origin_tenant_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE INDEX origin_tenant_fkey ON origin USING btree (tenant);


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON api_capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: idx_89786_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON asn FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cachegroup FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cachegroup_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON cdn FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON coordinate FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_regex FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_request FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_request_comment FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON deliveryservice_tmuser FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON division FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_deliveryservice FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_federation_resolver FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_resolver FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON federation_tmuser FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON hwinfo FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_agent FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON job_status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON log FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON origin FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON phys_location FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON profile_parameter FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON regex FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON region FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON role FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON role_capability FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON server FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON servercheck FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON snapshot FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON staticdnsentry FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON status FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON steering_target FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tenant FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON tm_user FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();


--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON type FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: on_update_current_timestamp; Type: TRIGGER; Schema: public; Owner: traffic_ops
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON user_role FOR EACH ROW EXECUTE PROCEDURE on_update_current_timestamp_last_updated();

--
-- Name: fk_atsprofile_atsparameters_atsparameters1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT fk_atsprofile_atsparameters_atsparameters1 FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;


--
-- Name: fk_atsprofile_atsparameters_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile_parameter
    ADD CONSTRAINT fk_atsprofile_atsparameters_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id) ON DELETE CASCADE;


--
-- Name: fk_cdn1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_cdn1 FOREIGN KEY (cdn_id) REFERENCES cdn(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: fk_cdn2; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_cdn2 FOREIGN KEY (cdn_id) REFERENCES cdn(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: fk_cg_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT fk_cg_1 FOREIGN KEY (parent_cachegroup_id) REFERENCES cachegroup(id);


--
-- Name: fk_cg_param_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_parameter
    ADD CONSTRAINT fk_cg_param_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id) ON DELETE CASCADE;


--
-- Name: fk_cg_secondary; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT fk_cg_secondary FOREIGN KEY (secondary_parent_cachegroup_id) REFERENCES cachegroup(id);


--
-- Name: fk_cg_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT fk_cg_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_contentserver_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id);


--
-- Name: fk_contentserver_contentserverstatus1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_contentserverstatus1 FOREIGN KEY (status) REFERENCES status(id);


--
-- Name: fk_contentserver_contentservertype1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_contentservertype1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_contentserver_phys_location1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_contentserver_phys_location1 FOREIGN KEY (phys_location) REFERENCES phys_location(id);


--
-- Name: fk_cran_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY asn
    ADD CONSTRAINT fk_cran_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id);


--
-- Name: fk_deliveryservice_profile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_deliveryservice_profile1 FOREIGN KEY (profile) REFERENCES profile(id);


--
-- Name: fk_deliveryservice_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_deliveryservice_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_ds_to_cs_contentserver1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT fk_ds_to_cs_contentserver1 FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_cs_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_server
    ADD CONSTRAINT fk_ds_to_cs_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_regex_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT fk_ds_to_regex_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ds_to_regex_regex1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_regex
    ADD CONSTRAINT fk_ds_to_regex_regex1 FOREIGN KEY (regex) REFERENCES regex(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_ext_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY to_extension
    ADD CONSTRAINT fk_ext_type FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_federation_federation_resolver1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_federation_resolver
    ADD CONSTRAINT fk_federation_federation_resolver1 FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_mapping_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_resolver
    ADD CONSTRAINT fk_federation_mapping_type FOREIGN KEY (type) REFERENCES type(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_resolver_to_fed1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_federation_resolver
    ADD CONSTRAINT fk_federation_resolver_to_fed1 FOREIGN KEY (federation_resolver) REFERENCES federation_resolver(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_tmuser_federation; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_tmuser
    ADD CONSTRAINT fk_federation_tmuser_federation FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_tmuser_role; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_tmuser
    ADD CONSTRAINT fk_federation_tmuser_role FOREIGN KEY (role) REFERENCES role(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_tmuser_tmuser; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_tmuser
    ADD CONSTRAINT fk_federation_tmuser_tmuser FOREIGN KEY (tm_user) REFERENCES tm_user(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_to_ds1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_deliveryservice
    ADD CONSTRAINT fk_federation_to_ds1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_federation_to_fed1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY federation_deliveryservice
    ADD CONSTRAINT fk_federation_to_fed1 FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_hwinfo1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY hwinfo
    ADD CONSTRAINT fk_hwinfo1 FOREIGN KEY (serverid) REFERENCES server(id) ON DELETE CASCADE;


--
-- Name: fk_job_agent_id1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_agent_id1 FOREIGN KEY (agent) REFERENCES job_agent(id) ON DELETE CASCADE;


--
-- Name: fk_job_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_deliveryservice1 FOREIGN KEY (job_deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE;


--
-- Name: fk_job_status_id1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_status_id1 FOREIGN KEY (status) REFERENCES job_status(id);


--
-- Name: fk_job_user_id1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY job
    ADD CONSTRAINT fk_job_user_id1 FOREIGN KEY (job_user) REFERENCES tm_user(id);


--
-- Name: fk_log_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY log
    ADD CONSTRAINT fk_log_1 FOREIGN KEY (tm_user) REFERENCES tm_user(id);


--
-- Name: fk_parameter; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_parameter
    ADD CONSTRAINT fk_parameter FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;


--
-- Name: fk_phys_location_region; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY phys_location
    ADD CONSTRAINT fk_phys_location_region FOREIGN KEY (region) REFERENCES region(id);


--
-- Name: fk_regex_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY regex
    ADD CONSTRAINT fk_regex_type1 FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_region_division1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY region
    ADD CONSTRAINT fk_region_division1 FOREIGN KEY (division) REFERENCES division(id);


--
-- Name: fk_server_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server
    ADD CONSTRAINT fk_server_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: fk_serverstatus_server1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY servercheck
    ADD CONSTRAINT fk_serverstatus_server1 FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE;


--
-- Name: fk_staticdnsentry_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id);


--
-- Name: fk_staticdnsentry_ds; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_ds FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE;


--
-- Name: fk_staticdnsentry_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY staticdnsentry
    ADD CONSTRAINT fk_staticdnsentry_type FOREIGN KEY (type) REFERENCES type(id);


--
-- Name: fk_steering_target_delivery_service; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY steering_target
    ADD CONSTRAINT fk_steering_target_delivery_service FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_steering_target_target; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY steering_target
    ADD CONSTRAINT fk_steering_target_target FOREIGN KEY (target) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_tm_user_ds; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT fk_tm_user_ds FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_tm_user_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_tmuser
    ADD CONSTRAINT fk_tm_user_id FOREIGN KEY (tm_user_id) REFERENCES tm_user(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fk_user_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tm_user
    ADD CONSTRAINT fk_user_1 FOREIGN KEY (role) REFERENCES role(id) ON DELETE SET NULL;


--
-- Name: fk_capability; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY api_capability
    ADD CONSTRAINT fk_capability FOREIGN KEY (capability) REFERENCES capability (name) ON DELETE RESTRICT;

--
-- Name: steering_target_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY steering_target
    ADD CONSTRAINT steering_target_type_fkey FOREIGN KEY (type) REFERENCES type (id);

--
-- Name: fk_tenantid; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice
    ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL;

--
-- Name: fk_author; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request
    ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES tm_user(id) ON DELETE CASCADE;

--
-- Name: fk_assignee; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request
    ADD CONSTRAINT fk_assignee FOREIGN KEY (assignee_id) REFERENCES tm_user(id) ON DELETE SET NULL;

--
-- Name: fk_last_edited_by; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request
    ADD CONSTRAINT fk_last_edited_by FOREIGN KEY (last_edited_by) REFERENCES tm_user (id) ON DELETE CASCADE;

--
-- Name: fk_author; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request_comment
    ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES tm_user (id) ON DELETE CASCADE;

--
-- Name: fk_profile; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT fk_profile FOREIGN KEY (profile) REFERENCES profile (id) ON DELETE RESTRICT;

--
-- Name: fk_deliveryservice; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT fk_deliveryservice FOREIGN KEY (deliveryservice) REFERENCES deliveryservice (id) ON DELETE CASCADE;

--
-- Name: fk_coordinate; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT fk_coordinate FOREIGN KEY (coordinate) REFERENCES coordinate (id) ON DELETE RESTRICT;

--
-- Name: fk_cachegroup; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT fk_cachegroup FOREIGN KEY (cachegroup) REFERENCES cachegroup (id) ON DELETE RESTRICT;

--
-- Name: fk_tenant; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY origin
    ADD CONSTRAINT fk_tenant FOREIGN KEY (tenant) REFERENCES tenant (id) ON DELETE RESTRICT;

--
-- Name: fk_coordinate; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup
    ADD CONSTRAINT fk_coordinate FOREIGN KEY (coordinate) REFERENCES coordinate (id);

--
-- Name: fk_primary_cg; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_fallbacks
    ADD CONSTRAINT fk_primary_cg FOREIGN KEY (primary_cg) REFERENCES cachegroup (id);

--
-- Name: fk_backup_cg; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_fallbacks
    ADD CONSTRAINT fk_backup_cg FOREIGN KEY (backup_cg) REFERENCES cachegroup (id);

--
-- Name: fk_cachegroup; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cachegroup_localization_method
    ADD CONSTRAINT fk_cachegroup FOREIGN KEY (cachegroup) REFERENCES cachegroup (id) ON DELETE CASCADE;

--
-- Name: fk_deliveryservice_request; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY deliveryservice_request_comment
    ADD CONSTRAINT fk_deliveryservice_request FOREIGN KEY (deliveryservice_request_id) REFERENCES deliveryservice_request (id) ON DELETE CASCADE;

--
-- Name: fk_cdn; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY profile
    ADD CONSTRAINT fk_cdn FOREIGN KEY (cdn) REFERENCES cdn (id) MATCH SIMPLE ON UPDATE RESTRICT ON DELETE RESTRICT;

--
-- Name: fk_role_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY role_capability
    ADD CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE CASCADE;

--
-- Name: fk_cap_name; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY role_capability
    ADD CONSTRAINT fk_cap_name FOREIGN KEY (cap_name) REFERENCES capability (name) ON DELETE RESTRICT;

--
-- Name: fk_cdn; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY snapshot
    ADD CONSTRAINT fk_cdn FOREIGN KEY (cdn) REFERENCES cdn (name) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: fk_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tenant
    ADD CONSTRAINT fk_parent_id FOREIGN KEY (parent_id) REFERENCES tenant (id);

--
-- Name: fk_tenant_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY tm_user
    ADD CONSTRAINT fk_tenant_id FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL;

--
-- Name: fk_tm_user; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY user_role
    ADD CONSTRAINT fk_tm_user FOREIGN KEY (user_id) REFERENCES tm_user (id) ON DELETE CASCADE;

--
-- Name: fk_role; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY user_role
    ADD CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE RESTRICT;


--
-- Name: public; Type: ACL; Schema: -; Owner: traffic_ops
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM traffic_ops;
GRANT ALL ON SCHEMA public TO traffic_ops;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

