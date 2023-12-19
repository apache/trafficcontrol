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

-- Dumped from database version 13.12
-- Dumped by pg_dump version 13.12

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

CREATE OR REPLACE FUNCTION on_update_current_timestamp_last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.last_updated = now();
  RETURN NEW;
END;
$$;

ALTER FUNCTION public.on_update_current_timestamp_last_updated() OWNER TO traffic_ops;

--
-- Name: before_server_table(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE OR REPLACE FUNCTION before_server_table()
    RETURNS TRIGGER AS
$$
DECLARE
    server_count BIGINT;
BEGIN
    WITH server_ips AS (
        SELECT s.id, i.name, ip.address, s.profile
        FROM server s
                JOIN interface i on i.server = s.ID
                JOIN ip_address ip on ip.Server = s.ID and ip.interface = i.name
        WHERE i.monitor = true
    )
    SELECT count(*)
    INTO server_count
    FROM server_ips sip
             JOIN server_ips sip2 on sip.id <> sip2.id
    WHERE sip.id = NEW.id
      AND sip2.address = sip.address
      AND sip2.profile = sip.profile;

    IF server_count > 0 THEN
        RAISE EXCEPTION 'Server [id:%] does not have a unique ip_address over the profile [id:%], [%] conflicts',
            NEW.id,
            NEW.profile,
            server_count;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
ALTER FUNCTION public.before_server_table() OWNER TO traffic_ops;

--
-- Name: before_ip_address_table(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE OR REPLACE FUNCTION before_ip_address_table()
    RETURNS TRIGGER
AS
$$
DECLARE
    server_count   BIGINT;
    server_id      BIGINT;
    server_profile BIGINT;
BEGIN
    WITH server_ips AS (
        SELECT s.id as sid, ip.interface, i.name, ip.address, s.profile, ip.server
        FROM server s
                 JOIN interface i
                      on i.server = s.ID
                 JOIN ip_address ip
                      on ip.Server = s.ID and ip.interface = i.name
        WHERE ip.service_address = true
    )
    SELECT count(distinct(sip.sid)), sip.sid, sip.profile
    INTO server_count, server_id, server_profile
    FROM server_ips sip
    WHERE (sip.server <> NEW.server AND (SELECT host(sip.address)) = (SELECT host(NEW.address)) AND sip.profile = (SELECT profile from server s WHERE s.id = NEW.server))
    GROUP BY sip.sid, sip.profile;

    IF server_count > 0 THEN
        RAISE EXCEPTION 'ip_address is not unique across the server [id:%] profile [id:%], [%] conflicts',
            server_id,
            server_profile,
            server_count;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;
ALTER FUNCTION public.before_ip_address_table() OWNER TO traffic_ops;

--
-- Name: on_delete_current_timestamp_last_updated(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE OR REPLACE FUNCTION on_delete_current_timestamp_last_updated()
    RETURNS trigger
AS $$
BEGIN
  update last_deleted set last_updated = now() where table_name = TG_ARGV[0];
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

ALTER FUNCTION on_delete_current_timestamp_last_updated() OWNER TO traffic_ops;

--
-- Name: update_ds_timestamp_on_insert(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE OR REPLACE FUNCTION update_ds_timestamp_on_insert()
    RETURNS trigger
    AS $$
BEGIN
    UPDATE deliveryservice
    SET last_updated=now()
    WHERE id IN (
        SELECT deliveryservice
        FROM CAST(NEW AS deliveryservice_tls_version)
    );
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

--
-- Name: update_ds_timestamp_on_delete(); Type: FUNCTION; Schema: public; Owner: traffic_ops
--

CREATE OR REPLACE FUNCTION update_ds_timestamp_on_delete()
    RETURNS trigger
    AS $$
BEGIN
    UPDATE deliveryservice
    SET last_updated=now()
    WHERE id IN (
        SELECT deliveryservice
        FROM CAST(OLD AS deliveryservice_tls_version)
    );
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

SET default_tablespace = '';

SET default_with_oids = false;

DO $$
BEGIN
IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'change_types') THEN
    --
    -- Name: change_types; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE change_types AS ENUM (
        'create',
        'update',
        'delete'
    );
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'deep_caching_type') THEN
    --
    -- Name: deep_caching_type; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE deep_caching_type AS ENUM (
        'NEVER',
        'ALWAYS'
    );
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'http_method_t') THEN
    --
    -- Name: http_method_t; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE http_method_t AS ENUM (
        'GET',
        'POST',
        'PUT',
        'PATCH',
        'DELETE'
    );
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'origin_protocol') THEN
    --
    -- Name: localization_method; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE localization_method AS ENUM (
        'CZ',
        'DEEP_CZ',
        'GEO'
    );
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'origin_protocol') THEN
    --
    -- Name: origin_protocol; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE origin_protocol AS ENUM (
        'http',
        'https'
    );
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'profile_type') THEN
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
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'workflow_states') THEN
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
END IF;

IF NOT EXISTS(SELECT FROM pg_type WHERE typname = 'server_ip_address') THEN
    --
    -- Name: server_ip_address; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE server_ip_address AS (address inet, gateway inet, service_address boolean);
END IF;

IF NOT EXISTS(SELECT FROM pg_type WHERE typname = 'server_interface') THEN
    --
    -- Name: server_interface; Type: TYPE; Schema: public; Owner: traffic_ops
    --

    CREATE TYPE server_interface AS (ip_addresses server_ip_address ARRAY, max_bandwidth bigint, monitor boolean, mtu bigint, name text);
END IF;

IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'deliveryservice_signature_type') THEN
    --
    -- Name: deliveryservice_signature_type; Type: DOMAIN; Schema: public; Owner: traffic_ops
    --

    CREATE DOMAIN deliveryservice_signature_type AS text CHECK (VALUE IN ('url_sig', 'uri_signing'));
END IF;
END$$;

--
-- Name: acme_account; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS acme_account (
    email text NOT NULL,
    private_key text NOT NULL,
    provider text NOT NULL,
    uri text NOT NULL,
    PRIMARY KEY (email, provider)
);

--
-- Name: api_capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS api_capability (
    id bigserial PRIMARY KEY,
    http_method http_method_t NOT NULL,
    route text NOT NULL,
    capability text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    UNIQUE (http_method, route, capability)
);

ALTER TABLE api_capability OWNER TO traffic_ops;

--
-- Name: asn; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS asn (
    id bigint NOT NULL,
    asn bigint NOT NULL,
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89468_primary PRIMARY KEY (id, cachegroup)
);

ALTER TABLE asn OWNER TO traffic_ops;

--
-- Name: asn_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS asn_id_seq
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
-- Name: async_status; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS async_status (
    id bigint NOT NULL,
    status TEXT NOT NULL,
    message TEXT,
    start_time timestamp with time zone DEFAULT now() NOT NULL,
    end_time timestamp with time zone,

    CONSTRAINT async_status_pkey PRIMARY KEY (id)
);

ALTER TABLE async_status OWNER TO traffic_ops;

--
-- Name: async_status_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS async_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE async_status_id_seq OWNER TO traffic_ops;

--
-- Name: async_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE async_status_id_seq OWNED BY async_status.id;

--
-- Name: cachegroup; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cachegroup (
    id bigint,
    name text NOT NULL,
    short_name text NOT NULL,
    parent_cachegroup_id bigint,
    secondary_parent_cachegroup_id bigint,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    fallback_to_closest boolean DEFAULT TRUE,
    coordinate bigint,
    CONSTRAINT idx_89476_primary PRIMARY KEY (id, type)
);

ALTER TABLE cachegroup OWNER TO traffic_ops;

--
-- Name: cachegroup_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS cachegroup_id_seq
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

CREATE TABLE IF NOT EXISTS cachegroup_fallbacks (
    primary_cg bigint NOT NULL,
    backup_cg bigint NOT NULL CHECK (primary_cg != backup_cg),
    set_order bigint NOT NULL,
    UNIQUE (primary_cg, backup_cg),
    UNIQUE (primary_cg, set_order)
);

ALTER TABLE cachegroup_fallbacks OWNER TO traffic_ops;

--
-- Name: cachegroup_localization_method; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cachegroup_localization_method (
    cachegroup bigint NOT NULL,
    method localization_method NOT NULL,
    UNIQUE (cachegroup, method)
);

--
-- Name: cachegroup_parameter; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cachegroup_parameter (
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89484_primary PRIMARY KEY (cachegroup, parameter)
);

ALTER TABLE cachegroup_parameter OWNER TO traffic_ops;

--
-- Name: capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS capability (
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT capability_pkey PRIMARY KEY (name)
);

ALTER TABLE capability OWNER TO traffic_ops;

--
-- Name: cdn; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdn (
    id bigint,
    name text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    dnssec_enabled boolean DEFAULT false NOT NULL,
    domain_name text NOT NULL,
    CONSTRAINT cdn_domain_name_unique UNIQUE (domain_name),
    CONSTRAINT idx_89491_primary PRIMARY KEY (id)
);

ALTER TABLE cdn OWNER TO traffic_ops;

--
-- Name: cdn_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS cdn_id_seq
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
-- Name: cdn_lock; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdn_lock (
    username text NOT NULL,
    cdn text NOT NULL,
    message text,
    soft boolean NOT NULL DEFAULT TRUE,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT pk_cdn_lock PRIMARY KEY ("cdn")
);

ALTER TABLE cdn_lock OWNER TO traffic_ops;

--
-- Name: cdn_notification; Type: TABLE; Schema: public; Owner: traffic_ops

CREATE TABLE cdn_notification (
    id bigint NOT NULL,
    cdn text NOT NULL,
    "user" text NOT NULL,
    notification text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT cdn_notification_pkey PRIMARY KEY (id)
);

ALTER TABLE cdn_notification OWNER TO traffic_ops;

--
-- Name: cdn_notification_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS cdn_notification_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE cdn_notification_id_seq OWNER TO traffic_ops;

--
-- Name: cdn_notification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cdn_notification_id_seq OWNED BY cdn_notification.id;

--
-- Name: coordinate; Type: TABLE; Schema: public: Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS coordinate (
    id bigserial,
    name text UNIQUE NOT NULL,
    latitude numeric NOT NULL DEFAULT 0.0,
    longitude numeric NOT NULL DEFAULT 0.0,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT coordinate_pkey PRIMARY KEY (id)
);

--
-- Name: deliveryservice; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice (
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
    geolimit_redirect_url text,
    tenant_id bigint NOT NULL,
    routing_name text NOT NULL DEFAULT 'cdn',
    deep_caching_type deep_caching_type NOT NULL DEFAULT 'NEVER',
    fq_pacing_rate bigint DEFAULT 0,
    anonymous_blocking_enabled boolean NOT NULL DEFAULT FALSE,
    consistent_hash_regex text,
    max_origin_connections bigint NOT NULL DEFAULT 0 CHECK (max_origin_connections >= 0),
    ecs_enabled boolean NOT NULL DEFAULT false,
    range_slice_block_size integer CHECK (range_slice_block_size >= 262144 AND range_slice_block_size <= 33554432) DEFAULT NULL,
    topology text,
    first_header_rewrite text,
    inner_header_rewrite text,
    last_header_rewrite text,
    service_category text,
    max_request_header_bytes int NOT NULL DEFAULT 0,
    CONSTRAINT routing_name_not_empty CHECK ((length(routing_name) > 0)),
    CONSTRAINT idx_89502_primary PRIMARY KEY (id, type)
);

ALTER TABLE deliveryservice OWNER TO traffic_ops;

--
-- Name: deliveryservices_required_capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservices_required_capability (
    required_capability TEXT NOT NULL,
    deliveryservice_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,

    PRIMARY KEY (deliveryservice_id, required_capability)
);

--
-- Name: deliveryservice_consistent_hash_query_param; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice_consistent_hash_query_param (
    name TEXT NOT NULL,
    deliveryservice_id bigint NOT NULL,
    CONSTRAINT name_empty CHECK (length(name) > 0),
    CONSTRAINT name_reserved CHECK (name NOT IN ('format','trred')),
    PRIMARY KEY (name, deliveryservice_id)
);

--
-- Name: deliveryservice_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS deliveryservice_id_seq
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

CREATE TABLE IF NOT EXISTS deliveryservice_regex (
    deliveryservice bigint NOT NULL,
    regex bigint NOT NULL,
    set_number bigint DEFAULT '0'::bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89517_primary PRIMARY KEY (deliveryservice, regex)
);

ALTER TABLE deliveryservice_regex OWNER TO traffic_ops;

--
-- Name: deliveryservice_request; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice_request (
    assignee_id bigint,
    author_id bigint NOT NULL,
    change_type change_types NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    id bigserial,
    last_edited_by_id bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    deliveryservice jsonb DEFAULT NULL,
    status workflow_states NOT NULL,
    original jsonb DEFAULT NULL,
    CONSTRAINT deliveryservice_request_pkey PRIMARY KEY (id),
    CONSTRAINT appropriate_requested_and_original_for_change_type CHECK (
        (change_type = 'delete' AND original IS NOT NULL AND deliveryservice IS NULL)
        OR
        (change_type = 'create' AND original IS NULL AND deliveryservice IS NOT NULL)
        OR (
            change_type = 'update' AND
            deliveryservice IS NOT NULL AND
            (
                (
                    (status = 'complete' OR status = 'rejected' OR status = 'pending')
                    AND
                    original IS NOT NULL
                )
                OR
                (
                    (status = 'draft' OR status = 'submitted')
                    AND
                    original IS NULL
                )
            )
        )
    )
);

ALTER TABLE deliveryservice_request OWNER TO traffic_ops;

--
-- Name: deliveryservice_request_comment; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice_request_comment (
    author_id bigint NOT NULL,
    deliveryservice_request_id bigint NOT NULL,
    id bigserial,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    value text NOT NULL,
    CONSTRAINT deliveryservice_request_comment_pkey PRIMARY KEY (id)
);

--
-- Name: deliveryservice_server; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice_server (
    deliveryservice bigint NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT idx_89521_primary PRIMARY KEY (deliveryservice, server)
);

ALTER TABLE deliveryservice_server OWNER TO traffic_ops;

--
-- Name: deliveryservice_tls_version; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS deliveryservice_tls_version (
    deliveryservice bigint NOT NULL,
    tls_version text NOT NULL,
    CONSTRAINT deliveryservice_tls_version_pkey PRIMARY KEY (deliveryservice, tls_version),
    CONSTRAINT deliveryservice_tls_version_tls_version_check CHECK (tls_version <> '')
);

ALTER TABLE deliveryservice_tls_version OWNER TO traffic_ops;

--
-- Name: update_ds_timestamp_on_tls_version_insertion; Type: TRIGGER; Schema: public; Owner: traffic_ops
--
DROP TRIGGER IF EXISTS update_ds_timestamp_on_tls_version_insertion on deliveryservice_tls_version;
CREATE TRIGGER update_ds_timestamp_on_tls_version_insertion
    AFTER INSERT ON deliveryservice_tls_version
    FOR EACH ROW EXECUTE PROCEDURE update_ds_timestamp_on_insert();

--
-- Name: update_ds_timestamp_on_tls_version_delete; Type: TRIGGER; Schema: public; Owner: traffic_ops
--
DROP TRIGGER IF EXISTS update_ds_timestamp_on_tls_version_delete on deliveryservice_tls_version;
CREATE TRIGGER update_ds_timestamp_on_tls_version_delete
    AFTER DELETE ON deliveryservice_tls_version
    FOR EACH ROW EXECUTE PROCEDURE update_ds_timestamp_on_delete();

--
-- Name: division; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS division (
    id bigint NOT NULL,
    name text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89531_primary PRIMARY KEY (id)
);

ALTER TABLE division OWNER TO traffic_ops;

--
-- Name: division_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS division_id_seq
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
-- Name: dnschallenges; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS dnschallenges (
    fqdn text NOT NULL,
    record text NOT NULL,
    xml_id text NOT NULL
);

--
-- Name: federation; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS federation (
    id bigint NOT NULL,
    cname text NOT NULL,
    description text,
    ttl integer NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89541_primary PRIMARY KEY (id)
);

ALTER TABLE federation OWNER TO traffic_ops;

--
-- Name: federation_deliveryservice; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS federation_deliveryservice (
    federation bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89549_primary PRIMARY KEY (federation, deliveryservice)
);

ALTER TABLE federation_deliveryservice OWNER TO traffic_ops;

--
-- Name: federation_federation_resolver; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS federation_federation_resolver (
    federation bigint NOT NULL,
    federation_resolver bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89553_primary PRIMARY KEY (federation, federation_resolver)
);

ALTER TABLE federation_federation_resolver OWNER TO traffic_ops;

--
-- Name: federation_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS federation_id_seq
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

CREATE TABLE IF NOT EXISTS federation_resolver (
    id bigint NOT NULL,
    ip_address text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89559_primary PRIMARY KEY (id)
);


ALTER TABLE federation_resolver OWNER TO traffic_ops;

--
-- Name: federation_resolver_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS federation_resolver_id_seq
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

CREATE TABLE IF NOT EXISTS federation_tmuser (
    federation bigint NOT NULL,
    tm_user bigint NOT NULL,
    role bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89567_primary PRIMARY KEY (federation, tm_user)
);

ALTER TABLE federation_tmuser OWNER TO traffic_ops;

--
-- Name: hwinfo; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS hwinfo (
    id bigint NOT NULL,
    serverid bigint NOT NULL,
    description text NOT NULL,
    val text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89583_primary PRIMARY KEY (id)
);

ALTER TABLE hwinfo OWNER TO traffic_ops;

--
-- Name: hwinfo_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS hwinfo_id_seq
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
-- Name: interface; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS interface (
    max_bandwidth bigint DEFAULT NULL CHECK (max_bandwidth IS NULL OR max_bandwidth >= 0),
    monitor boolean NOT NULL,
    mtu bigint DEFAULT 1500,
    name text NOT NULL CHECK (name != ''),
    server bigint NOT NULL,
    router_host_name text NOT NULL DEFAULT '',
    router_port_name text NOT NULL DEFAULT '',
    PRIMARY KEY (name, server)
);

ALTER TABLE interface
ADD CONSTRAINT interface_mtu_check
CHECK (((mtu IS NULL) OR (mtu >= 1280)));

--
-- Name: ip_address; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS ip_address (
    address inet NOT NULL,
    gateway inet CHECK (
        gateway IS NULL OR (
            family(gateway) = 4 AND
            masklen(gateway) = 32
        ) OR (
            family(gateway) = 6 AND
            masklen(gateway) = 128
        )
    ),
    interface text NOT NULL,
    server bigint NOT NULL,
    service_address boolean NOT NULL DEFAULT FALSE,
    PRIMARY KEY (address, interface, server)
);

--
-- Name: job; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS job (
    id bigint NOT NULL,
    ttl_hr integer,
    asset_url text NOT NULL,
    start_time timestamp with time zone NOT NULL,
    entered_time timestamp with time zone NOT NULL,
    job_user bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    job_deliveryservice bigint,
    invalidation_type text NOT NULL DEFAULT 'REFRESH',
    CONSTRAINT idx_89593_primary PRIMARY KEY (id)
);

ALTER TABLE job OWNER TO traffic_ops;

--
-- Name: job_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS job_id_seq
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

CREATE TABLE IF NOT EXISTS last_deleted (
    table_name text NOT NULL PRIMARY KEY,
    last_updated timestamp with time zone NOT NULL DEFAULT now()
);

--
-- Name: log; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS log (
    id bigint NOT NULL,
    level text,
    message text NOT NULL,
    tm_user bigint NOT NULL,
    ticketnum text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89634_primary PRIMARY KEY (id, tm_user)
);

ALTER TABLE log OWNER TO traffic_ops;

--
-- Name: log_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS log_id_seq
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

CREATE TABLE IF NOT EXISTS origin (
    id bigserial NOT NULL,
    name text UNIQUE NOT NULL,
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
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT origin_pkey PRIMARY KEY (id)
);

ALTER TABLE origin OWNER TO traffic_ops;

--
-- Name: parameter; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS parameter (
    id bigint NOT NULL,
    name text NOT NULL,
    config_file text,
    value text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    secure boolean DEFAULT false NOT NULL,
    CONSTRAINT unique_param UNIQUE (name, config_file, value),
    CONSTRAINT idx_89644_primary PRIMARY KEY (id)
);


ALTER TABLE parameter OWNER TO traffic_ops;

--
-- Name: parameter_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS parameter_id_seq
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

CREATE TABLE IF NOT EXISTS phys_location (
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
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89655_primary PRIMARY KEY (id)
);

ALTER TABLE phys_location OWNER TO traffic_ops;

--
-- Name: phys_location_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS phys_location_id_seq
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

CREATE TABLE IF NOT EXISTS profile (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    type profile_type NOT NULL,
    cdn bigint NOT NULL,
    routing_disabled boolean NOT NULL DEFAULT FALSE,
    CONSTRAINT idx_89665_primary PRIMARY KEY (id)
);

ALTER TABLE profile OWNER TO traffic_ops;

--
-- Name: profile_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS profile_id_seq
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

CREATE TABLE IF NOT EXISTS profile_parameter (
    profile bigint NOT NULL,
    parameter bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89673_primary PRIMARY KEY (profile, parameter)
);


ALTER TABLE profile_parameter OWNER TO traffic_ops;

--
-- Name: regex; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS regex (
    id bigint NOT NULL,
    pattern text DEFAULT ''::text NOT NULL,
    type bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89679_primary PRIMARY KEY (id, type)
);

ALTER TABLE regex OWNER TO traffic_ops;

--
-- Name: regex_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS regex_id_seq
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

CREATE TABLE IF NOT EXISTS region (
    id bigint NOT NULL,
    name text NOT NULL,
    division bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89690_primary PRIMARY KEY (id)
);

ALTER TABLE region OWNER TO traffic_ops;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS region_id_seq
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

CREATE TABLE IF NOT EXISTS role (
    id bigint,
    name text NOT NULL,
    description text NOT NULL,
    priv_level bigint NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT role_name_unique UNIQUE (name),
    CONSTRAINT idx_89700_primary PRIMARY KEY (id)
);

ALTER TABLE role OWNER TO traffic_ops;

--
-- Name: role_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS role_id_seq
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

CREATE TABLE IF NOT EXISTS role_capability (
    role_id bigint NOT NULL,
    cap_name text NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    UNIQUE (role_id, cap_name)
);

ALTER TABLE role_capability OWNER TO traffic_ops;

--
-- Name: server; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS server (
    id bigint NOT NULL,
    host_name text NOT NULL,
    domain_name text NOT NULL,
    tcp_port bigint,
    xmpp_id text,
    xmpp_passwd text,
    phys_location bigint NOT NULL,
    rack text,
    cachegroup bigint DEFAULT '0'::bigint NOT NULL,
    type bigint NOT NULL,
    status bigint NOT NULL,
    offline_reason text,
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
    guid text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    https_port bigint,
    status_last_updated timestamp with time zone,
    config_update_time timestamp with time zone NOT NULL DEFAULT TIMESTAMP 'epoch',
    config_apply_time timestamp with time zone NOT NULL DEFAULT TIMESTAMP 'epoch',
    revalidate_update_time timestamp with time zone NOT NULL DEFAULT TIMESTAMP 'epoch',
    revalidate_apply_time timestamp with time zone NOT NULL DEFAULT TIMESTAMP 'epoch',
    CONSTRAINT idx_89709_primary PRIMARY KEY (id)
);

ALTER TABLE server OWNER TO traffic_ops;

--
-- Name: server_capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS server_capability (
    name TEXT PRIMARY KEY,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT name_empty CHECK (length(name) > 0)
);

--
-- Name: server_server_capability; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS server_server_capability (
    server_capability TEXT NOT NULL,
    server bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,

    PRIMARY KEY (server, server_capability)
);

--
-- Name: service_category; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS service_category (
    name TEXT PRIMARY KEY CHECK (name <> ''),
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

--
-- Name: server_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS server_id_seq
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

CREATE TABLE IF NOT EXISTS servercheck (
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
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89722_primary PRIMARY KEY (id, server)
);

ALTER TABLE servercheck OWNER TO traffic_ops;

--
-- Name: servercheck_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS servercheck_id_seq
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

CREATE TABLE IF NOT EXISTS snapshot (
    cdn text NOT NULL,
    crconfig json NOT NULL,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    monitoring json NOT NULL,
    CONSTRAINT snapshot_pkey PRIMARY KEY (cdn)
);

ALTER TABLE snapshot OWNER TO traffic_ops;

--
-- Name: staticdnsentry; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS staticdnsentry (
    id bigint NOT NULL,
    host text NOT NULL,
    address text NOT NULL,
    type bigint NOT NULL,
    ttl bigint DEFAULT '3600'::bigint NOT NULL,
    deliveryservice bigint NOT NULL,
    cachegroup bigint,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT idx_89729_primary PRIMARY KEY (id)
);

ALTER TABLE staticdnsentry OWNER TO traffic_ops;

--
-- Name: staticdnsentry_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS staticdnsentry_id_seq
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

CREATE TABLE IF NOT EXISTS stats_summary (
    id bigint NOT NULL,
    cdn_name text DEFAULT 'all'::text NOT NULL,
    deliveryservice_name text NOT NULL,
    stat_name text NOT NULL,
    stat_value double precision NOT NULL,
    summary_time timestamp with time zone DEFAULT now() NOT NULL,
    stat_date date,
    CONSTRAINT idx_89740_primary PRIMARY KEY (id)
);

ALTER TABLE stats_summary OWNER TO traffic_ops;

--
-- Name: stats_summary_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS stats_summary_id_seq
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

CREATE TABLE IF NOT EXISTS status (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT status_name_unique UNIQUE (name),
    CONSTRAINT idx_89751_primary PRIMARY KEY (id)
);

ALTER TABLE status OWNER TO traffic_ops;

--
-- Name: status_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS status_id_seq
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

CREATE TABLE IF NOT EXISTS steering_target (
    deliveryservice bigint NOT NULL,
    target bigint NOT NULL,
    value bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    type bigint NOT NULL,
    CONSTRAINT idx_89759_primary PRIMARY KEY (deliveryservice, target)
);

ALTER TABLE steering_target OWNER TO traffic_ops;

--
-- Name: tenant; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS tenant (
    id bigserial,
    name text UNIQUE NOT NULL,
    active boolean NOT NULL DEFAULT FALSE,
    parent_id bigint DEFAULT 1 CHECK (id != parent_id),
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT tenant_pkey PRIMARY KEY (id)
);

--
-- Name: tm_user; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS tm_user (
    id bigint NOT NULL,
    username text NOT NULL,
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
    tenant_id bigint NOT NULL,
    last_authenticated timestamp with time zone,
    ucdn text NOT NULL DEFAULT '',
    CONSTRAINT idx_89765_primary PRIMARY KEY (id)
);

ALTER TABLE tm_user OWNER TO traffic_ops;

--
-- Name: tm_user_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS tm_user_id_seq
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

CREATE TABLE IF NOT EXISTS to_extension (
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
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT idx_89776_primary PRIMARY KEY (id)
);

ALTER TABLE to_extension OWNER TO traffic_ops;

--
-- Name: to_extension_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS to_extension_id_seq
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
-- Name: toplogy; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS topology (
    name text PRIMARY KEY,
    description text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

--
-- Name: topology_cachegroup; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS topology_cachegroup (
    id BIGSERIAL PRIMARY KEY,
    topology text NOT NULL,
    cachegroup text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT unique_topology_cachegroup UNIQUE (topology, cachegroup)
);

--
-- Name: topology_cachegroup_parents; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS topology_cachegroup_parents (
    child bigint NOT NULL,
    parent bigint NOT NULL,
    rank integer NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT topology_cachegroup_parents_rank_check CHECK (rank = 1 OR rank = 2),
    CONSTRAINT unique_child_rank UNIQUE (child, rank),
    CONSTRAINT unique_child_parent UNIQUE (child, parent)
);

--
-- Name: type; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS type (
    id bigint NOT NULL,
    name text NOT NULL,
    description text,
    use_in_table text,
    last_updated timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT type_name_unique UNIQUE(name),
    CONSTRAINT idx_89786_primary PRIMARY KEY (id)
);

ALTER TABLE type OWNER TO traffic_ops;

--
-- Name: type_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE IF NOT EXISTS type_id_seq
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

DO $$ BEGIN
IF NOT EXISTS (SELECT FROM information_schema.tables
    WHERE table_name = 'profile_type_values'
    AND table_type = 'VIEW') THEN
    --
    -- Name: profile_type_values; Type: VIEW; Schema: public; Owner: traffic_ops
    --

    CREATE VIEW profile_type_values AS
        SELECT unnest(enum_range(NULL::profile_type)) AS VALUE
            ORDER BY (unnest(enum_range(NULL::profile_type)));

    ALTER TABLE profile_type_values OWNER TO traffic_ops;
END IF;
END$$;

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY asn ALTER COLUMN id SET DEFAULT nextval('asn_id_seq'::regclass);

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY async_status ALTER COLUMN id SET DEFAULT nextval('async_status_id_seq'::regclass);

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

ALTER TABLE ONLY cdn_notification ALTER COLUMN id SET DEFAULT nextval('cdn_notification_id_seq'::regclass);

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

DO $$ BEGIN
IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'asn' AND column_name = 'id') THEN
    --
    -- Name: idx_89468_cr_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89468_cr_id_unique ON asn USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'asn' AND column_name = 'cachegroup') THEN
    --
    -- Name: idx_89468_fk_cran_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89468_fk_cran_cachegroup1 ON asn USING btree (cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'name') THEN
    --
    -- Name: idx_89476_cg_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89476_cg_name_unique ON cachegroup USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'short_name') THEN
    --
    -- Name: idx_89476_cg_short_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89476_cg_short_unique ON cachegroup USING btree (short_name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'parent_cachegroup_id') THEN
    --
    -- Name: idx_89476_fk_cg_1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89476_fk_cg_1 ON cachegroup USING btree (parent_cachegroup_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'secondary_parent_cachegroup_id') THEN
    --
    -- Name: idx_89476_fk_cg_secondary; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89476_fk_cg_secondary ON cachegroup USING btree (secondary_parent_cachegroup_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'type') THEN
    --
    -- Name: idx_89476_fk_cg_type1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89476_fk_cg_type1 ON cachegroup USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'id') THEN
    --
    -- Name: idx_89476_lo_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89476_lo_id_unique ON cachegroup USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup_parameter' AND column_name = 'parameter') THEN
    --
    -- Name: idx_89484_fk_parameter; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89484_fk_parameter ON cachegroup_parameter USING btree (parameter);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cdn' AND column_name = 'name') THEN
    --
    -- Name: idx_89491_cdn_cdn_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89491_cdn_cdn_unique ON cdn USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'id') THEN
    --
    -- Name: idx_89502_ds_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89502_ds_id_unique ON deliveryservice USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'tenant_id') THEN
    --
    -- Name: idx_k_deliveryservice_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_k_deliveryservice_tenant_idx ON deliveryservice USING btree (tenant_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'xml_id') THEN
    --
    -- Name: idx_89502_ds_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89502_ds_name_unique ON deliveryservice USING btree (xml_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'cdn_id') THEN
    --
    -- Name: idx_89502_fk_cdn1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89502_fk_cdn1 ON deliveryservice USING btree (cdn_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'profile') THEN
    --
    -- Name: idx_89502_fk_deliveryservice_profile1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89502_fk_deliveryservice_profile1 ON deliveryservice USING btree (profile);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'type') THEN
    --
    -- Name: idx_89502_fk_deliveryservice_type1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89502_fk_deliveryservice_type1 ON deliveryservice USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice_regex' AND column_name = 'regex') THEN
    --
    -- Name: idx_89517_fk_ds_to_regex_regex1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89517_fk_ds_to_regex_regex1 ON deliveryservice_regex USING btree (regex);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice_server' AND column_name = 'server') THEN
    --
    -- Name: idx_89521_fk_ds_to_cs_contentserver1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89521_fk_ds_to_cs_contentserver1 ON deliveryservice_server USING btree (server);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'division' AND column_name = 'name') THEN
    --
    -- Name: idx_89531_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89531_name_unique ON division USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_deliveryservice' AND column_name = 'deliveryservice') THEN
    --
    -- Name: idx_89549_fk_fed_to_ds1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89549_fk_fed_to_ds1 ON federation_deliveryservice USING btree (deliveryservice);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_federation_resolver' AND column_name = 'federation') THEN
    --
    -- Name: idx_89553_fk_federation_federation_resolver; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89553_fk_federation_federation_resolver ON federation_federation_resolver USING btree (federation);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_federation_resolver' AND column_name = 'federation_resolver') THEN
    --
    -- Name: idx_89553_fk_federation_resolver_to_fed1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89553_fk_federation_resolver_to_fed1 ON federation_federation_resolver USING btree (federation_resolver);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_resolver' AND column_name = 'ip_address') THEN
    --
    -- Name: idx_89559_federation_resolver_ip_address; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89559_federation_resolver_ip_address ON federation_resolver USING btree (ip_address);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_resolver' AND column_name = 'type') THEN
    --
    -- Name: idx_89559_fk_federation_mapping_type; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89559_fk_federation_mapping_type ON federation_resolver USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_tmuser' AND column_name = 'federation') THEN
    --
    -- Name: idx_89567_fk_federation_federation_resolver; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89567_fk_federation_federation_resolver ON federation_tmuser USING btree (federation);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_tmuser' AND column_name = 'role') THEN
    --
    -- Name: idx_89567_fk_federation_tmuser_role; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89567_fk_federation_tmuser_role ON federation_tmuser USING btree (role);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'federation_tmuser' AND column_name = 'tm_user') THEN
    --
    -- Name: idx_89567_fk_federation_tmuser_tmuser; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89567_fk_federation_tmuser_tmuser ON federation_tmuser USING btree (tm_user);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'hwinfo' AND column_name = 'serverid') THEN
    --
    -- Name: idx_89583_fk_hwinfo1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89583_fk_hwinfo1 ON hwinfo USING btree (serverid);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'hwinfo' AND column_name = 'serverid')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'hwinfo' AND column_name = 'description') THEN
    --
    -- Name: idx_89583_serverid; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89583_serverid ON hwinfo USING btree (serverid, description);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'job' AND column_name = 'job_deliveryservice') THEN
    --
    -- Name: idx_89593_fk_job_deliveryservice1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89593_fk_job_deliveryservice1 ON job USING btree (job_deliveryservice);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'job' AND column_name = 'job_user') THEN
    --
    -- Name: idx_89593_fk_job_user_id1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89593_fk_job_user_id1 ON job USING btree (job_user);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'job' AND column_name = 'start_time') THEN
    --
    -- Name: job_start_time_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS job_start_time_idx ON job (start_time DESC NULLS LAST);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'log' AND column_name = 'tm_user') THEN
    --
    -- Name: idx_89634_fk_log_1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89634_fk_log_1 ON log USING btree (tm_user);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'log' AND column_name = 'last_updated') THEN
    --
    -- Name: idx_89634_idx_last_updated; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89634_idx_last_updated ON log USING btree (last_updated);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'parameter' AND column_name = 'name')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'parameter' AND column_name = 'value') THEN
    --
    -- Name: idx_89644_parameter_name_value_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89644_parameter_name_value_idx ON parameter USING btree (name, value);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'phys_location' AND column_name = 'region') THEN
    --
    -- Name: idx_89655_fk_phys_location_region_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89655_fk_phys_location_region_idx ON phys_location USING btree (region);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'phys_location' AND column_name = 'name') THEN
    --
    -- Name: idx_89655_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89655_name_unique ON phys_location USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'phys_location' AND column_name = 'short_name') THEN
    --
    -- Name: idx_89655_short_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89655_short_name_unique ON phys_location USING btree (short_name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'profile' AND column_name = 'name') THEN
    --
    -- Name: idx_89665_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89665_name_unique ON profile USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'profile' AND column_name = 'cdn') THEN
    --
    -- Name: idx_181818_fk_cdn1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_181818_fk_cdn1 ON profile USING btree (cdn);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'profile_parameter' AND column_name = 'parameter') THEN
    --
    -- Name: idx_89673_fk_atsprofile_atsparameters_atsparameters1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89673_fk_atsprofile_atsparameters_atsparameters1 ON profile_parameter USING btree (parameter);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'profile_parameter' AND column_name = 'profile') THEN
    --
    -- Name: idx_89673_fk_atsprofile_atsparameters_atsprofile1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89673_fk_atsprofile_atsparameters_atsprofile1 ON profile_parameter USING btree (profile);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'regex' AND column_name = 'type') THEN
    --
    -- Name: idx_89679_fk_regex_type1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89679_fk_regex_type1 ON regex USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'regex' AND column_name = 'id') THEN
    --
    -- Name: idx_89679_re_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89679_re_id_unique ON regex USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'region' AND column_name = 'division') THEN
    --
    -- Name: idx_89690_fk_region_division1_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89690_fk_region_division1_idx ON region USING btree (division);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'region' AND column_name = 'name') THEN
    --
    -- Name: idx_89690_name_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89690_name_unique ON region USING btree (name);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'cdn_id') THEN
    --
    -- Name: idx_89709_fk_cdn2; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_cdn2 ON server USING btree (cdn_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'profile') THEN
    --
    -- Name: idx_89709_fk_contentserver_atsprofile1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_contentserver_atsprofile1 ON server USING btree (profile);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'status') THEN
    --
    -- Name: idx_89709_fk_contentserver_contentserverstatus1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_contentserver_contentserverstatus1 ON server USING btree (status);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'type') THEN
    --
    -- Name: idx_89709_fk_contentserver_contentservertype1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_contentserver_contentservertype1 ON server USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'phys_location') THEN
    --
    -- Name: idx_89709_fk_contentserver_phys_location1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_contentserver_phys_location1 ON server USING btree (phys_location);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'cachegroup') THEN
    --
    -- Name: idx_89709_fk_server_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89709_fk_server_cachegroup1 ON server USING btree (cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'server' AND column_name = 'id') THEN
    --
    -- Name: idx_89709_se_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89709_se_id_unique ON server USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'servercheck' AND column_name = 'server') THEN
    --
    -- Name: idx_89722_fk_serverstatus_server1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89722_fk_serverstatus_server1 ON servercheck USING btree (server);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'servercheck' AND column_name = 'server') THEN
    --
    -- Name: idx_89722_server; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89722_server ON servercheck USING btree (server);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'servercheck' AND column_name = 'id') THEN
    --
    -- Name: idx_89722_ses_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89722_ses_id_unique ON servercheck USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'host')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'address')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'deliveryservice')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'cachegroup') THEN
    --
    -- Name: idx_89729_combi_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89729_combi_unique ON staticdnsentry USING btree (host, address, deliveryservice, cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'cachegroup') THEN
    --
    -- Name: idx_89729_fk_staticdnsentry_cachegroup1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89729_fk_staticdnsentry_cachegroup1 ON staticdnsentry USING btree (cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'deliveryservice') THEN
    --
    -- Name: idx_89729_fk_staticdnsentry_ds; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89729_fk_staticdnsentry_ds ON staticdnsentry USING btree (deliveryservice);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'staticdnsentry' AND column_name = 'type') THEN
    --
    -- Name: idx_89729_fk_staticdnsentry_type; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89729_fk_staticdnsentry_type ON staticdnsentry USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'topology_cachegroup' AND column_name = 'cachegroup') THEN
    --
    -- Name: topology_cachegroup_cachegroup_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS topology_cachegroup_cachegroup_fkey ON topology_cachegroup USING btree (cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'topology_cachegroup' AND column_name = 'topology') THEN
    --
    -- Name: topology_cachegroup_topology_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS topology_cachegroup_topology_fkey ON topology_cachegroup USING btree (topology);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'topology_cachegroup_parents' AND column_name = 'child') THEN
    --
    -- Name: topology_cachegroup_parents_child_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS topology_cachegroup_parents_child_fkey ON topology_cachegroup_parents USING btree (child);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'topology_cachegroup_parents' AND column_name = 'parent') THEN
    --
    -- Name: topology_cachegroup_parents_parents_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS topology_cachegroup_parents_parents_fkey ON topology_cachegroup_parents USING btree (parent);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'deliveryservice' AND column_name = 'topology') THEN
    --
    -- Name: deliveryservice_topology_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS deliveryservice_topology_fkey ON deliveryservice USING btree (topology);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'tenant' AND column_name = 'parent_id') THEN
    --
    -- Name: idx_k_tenant_parent_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_k_tenant_parent_tenant_idx ON tenant USING btree (parent_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'tm_user' AND column_name = 'role') THEN
    --
    -- Name: idx_89765_fk_user_1; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89765_fk_user_1 ON tm_user USING btree (role);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'tm_user' AND column_name = 'tenant_id') THEN
    --
    -- Name: idx_k_tm_user_tenant_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_k_tm_user_tenant_idx ON tm_user USING btree (tenant_id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'tm_user' AND column_name = 'email') THEN
    --
    -- Name: idx_89765_tmuser_email_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89765_tmuser_email_unique ON tm_user USING btree (email);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'tm_user' AND column_name = 'username') THEN
    --
    -- Name: idx_89765_username_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89765_username_unique ON tm_user USING btree (username);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'to_extension' AND column_name = 'type') THEN
    --
    -- Name: idx_89776_fk_ext_type_idx; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS idx_89776_fk_ext_type_idx ON to_extension USING btree (type);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'to_extension' AND column_name = 'id') THEN
    --
    -- Name: idx_89776_id_unique; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS idx_89776_id_unique ON to_extension USING btree (id);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup' AND column_name = 'coordinate') THEN
    --
    -- Name: cachegroup_coordinate_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS cachegroup_coordinate_fkey ON cachegroup USING btree (coordinate);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'cachegroup_localization_method' AND column_name = 'cachegroup') THEN
    --
    -- Name: cachegroup_localization_method_cachegroup_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS cachegroup_localization_method_cachegroup_fkey ON cachegroup_localization_method USING btree(cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'is_primary')
    AND EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'deliveryservice') THEN
    --
    -- Name: origin_is_primary_deliveryservice_constraint; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE UNIQUE INDEX IF NOT EXISTS origin_is_primary_deliveryservice_constraint ON origin (is_primary, deliveryservice) WHERE is_primary;
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'deliveryservice') THEN
    --
    -- Name: origin_deliveryservice_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS origin_deliveryservice_fkey ON origin USING btree (deliveryservice);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'coordinate') THEN
    --
    -- Name: origin_coordinate_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS origin_coordinate_fkey ON origin USING btree (coordinate);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'profile') THEN
    --
    -- Name: origin_profile_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS origin_profile_fkey ON origin USING btree (profile);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'cachegroup') THEN
    --
    -- Name: origin_cachegroup_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS origin_cachegroup_fkey ON origin USING btree (cachegroup);
END IF;

IF EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'origin' AND column_name = 'tenant') THEN
    --
    -- Name: origin_tenant_fkey; Type: INDEX; Schema: public; Owner: traffic_ops
    --

    CREATE INDEX IF NOT EXISTS origin_tenant_fkey ON origin USING btree (tenant);
END IF;
END$$;


DO $$
    DECLARE table_names VARCHAR[] := CAST(ARRAY[
        'api_capability',
        'asn',
        'cachegroup',
        'cachegroup_fallbacks',
        'cachegroup_localization_method',
        'cachegroup_parameter',
        'capability',
        'cdn',
        'coordinate',
        'deliveryservice',
        'deliveryservice_regex',
        'deliveryservice_request',
        'deliveryservice_request_comment',
        'deliveryservice_server',
        'deliveryservices_required_capability',
        'division',
        'federation',
        'federation_deliveryservice',
        'federation_federation_resolver',
        'federation_resolver',
        'federation_tmuser',
        'hwinfo',
        'job',
        'log',
        'origin',
        'parameter',
        'phys_location',
        'profile',
        'profile_parameter',
        'regex',
        'region',
        'role',
        'role_capability',
        'server',
        'server_capability',
        'server_server_capability',
        'servercheck',
        'service_category',
        'snapshot',
        'staticdnsentry',
        'stats_summary',
        'status',
        'steering_target',
        'tenant',
        'tm_user',
        'to_extension',
        'topology',
        'topology_cachegroup',
        'topology_cachegroup_parents',
        'type'
    ] AS VARCHAR[]);
    table_name TEXT;
    trigger_name TEXT := 'on_delete_current_timestamp';
    trigger_exists BOOLEAN;

BEGIN
    FOREACH table_name IN ARRAY table_names
        LOOP
            EXECUTE FORMAT('SELECT EXISTS (
                SELECT
                FROM pg_catalog.pg_trigger
                WHERE tgname = ''%s''
                AND tgrelid = CAST(''%s'' AS REGCLASS))
                ',
                QUOTE_IDENT(trigger_name),
                QUOTE_IDENT(table_name)) INTO trigger_exists;
            IF NOT trigger_exists
            THEN
                EXECUTE FORMAT('
                    CREATE TRIGGER %s
                    AFTER DELETE ON %s
                    FOR EACH ROW
                        EXECUTE PROCEDURE %s_last_updated(''%s'');
                    ',
                    QUOTE_IDENT(trigger_name),
                    QUOTE_IDENT(table_name),
                    QUOTE_IDENT(trigger_name),
                    QUOTE_IDENT(table_name)
                    );
            END IF;
            IF table_name = 'topology' THEN
                EXECUTE FORMAT('
                        CREATE INDEX IF NOT EXISTS %s_last_updated_idx
                               ON %s (last_updated DESC NULLS LAST);
                        ',
                        QUOTE_IDENT(table_name),
                        QUOTE_IDENT(table_name)
                    );
            ELSIF table_name = 'phys_location' THEN
            EXECUTE FORMAT('
                        CREATE INDEX IF NOT EXISTS phys_location_last_updated_idx
                               ON %s (last_updated DESC NULLS LAST);
                        ',
                        QUOTE_IDENT(table_name)
                );

            ELSIF NOT (table_name = 'stats_summary' OR table_name = 'cachegroup_fallbacks' OR table_name = 'cachegroup_localization_method')THEN
                EXECUTE FORMAT('
                        CREATE INDEX IF NOT EXISTS %s_last_updated_idx
                               ON %s (last_updated DESC NULLS LAST);
                        ',
                        QUOTE_IDENT(table_name),
                        QUOTE_IDENT(table_name)
                        );
            END IF;
        END LOOP;
END
$$;

--
-- Add on_update_current_timestamp TRIGGER to all tables
--
DO $$
DECLARE
    table_names VARCHAR[] := CAST(ARRAY[
        'api_capability',
        'asn',
        'cachegroup',
        'cachegroup_parameter',
        'capability',
        'cdn',
        'cdn_lock',
        'cdn_notification',
        'coordinate',
        'deliveryservice',
        'deliveryservice_regex',
        'deliveryservice_request',
        'deliveryservice_request_comment',
        'deliveryservice_server',
        'division',
        'federation',
        'federation_deliveryservice',
        'federation_federation_resolver',
        'federation_resolver',
        'federation_tmuser',
        'hwinfo',
        'job',
        'log',
        'origin',
        'parameter',
        'phys_location',
        'profile',
        'profile_parameter',
        'regex',
        'region',
        'role',
        'role_capability',
        'server',
        'servercheck',
        'snapshot',
        'staticdnsentry',
        'status',
        'steering_target',
        'tenant',
        'tm_user',
        'topology',
        'topology_cachegroup',
        'topology_cachegroup_parents',
        'type'
    ] AS VARCHAR[]);
    table_name TEXT;
    trigger_name TEXT := 'on_update_current_timestamp';
    trigger_exists BOOLEAN;
BEGIN
    FOREACH table_name IN ARRAY table_names
    LOOP
        EXECUTE FORMAT('SELECT EXISTS (
            SELECT
            FROM pg_catalog.pg_trigger
            WHERE tgname = ''%s''
            AND tgrelid = CAST(''%s'' AS REGCLASS))
            ',
            QUOTE_IDENT(trigger_name),
            QUOTE_IDENT(table_name)) INTO trigger_exists;
        IF NOT trigger_exists
        THEN
            EXECUTE FORMAT('
                    CREATE TRIGGER %s
                    BEFORE UPDATE ON %s
                    FOR EACH ROW
                        EXECUTE PROCEDURE %s_last_updated();
                ',
                QUOTE_IDENT(trigger_name),
                QUOTE_IDENT(table_name),
                QUOTE_IDENT(trigger_name)
            );
        END IF;
    END LOOP;
END$$;

-- New code block to deallocate table_name variable to avoid identifier collision
DO $$ BEGIN
IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_notification_cdn' AND table_name = 'cdn_notification') THEN
    --
    -- Name: fk_notification_cdn; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cdn_notification
        ADD CONSTRAINT fk_notification_cdn FOREIGN KEY (cdn) REFERENCES cdn(name) ON DELETE CASCADE ON UPDATE CASCADE;

END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_notification_user' AND table_name = 'cdn_notification') THEN
    --
    -- Name: fk_notification_user; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cdn_notification
        ADD CONSTRAINT fk_notification_user FOREIGN KEY ("user") REFERENCES tm_user(username) ON DELETE CASCADE ON UPDATE CASCADE;

END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_atsprofile_atsparameters_atsparameters1' AND table_name = 'profile_parameter') THEN
    --
    -- Name: fk_atsprofile_atsparameters_atsparameters1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY profile_parameter
        ADD CONSTRAINT fk_atsprofile_atsparameters_atsparameters1 FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'deliveryservice_service_category_fkey' AND table_name = 'deliveryservice') THEN
    --
    -- Name: deliveryservice_service_category_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT deliveryservice_service_category_fkey FOREIGN KEY (service_category) REFERENCES service_category(name) ON UPDATE CASCADE;
END IF;


IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'deliveryservice_tls_version_deliveryservice_fkey' AND table_name = 'deliveryservice_tls_version') THEN
    --
    -- Name: deliveryservice_tls_version_deliveryservice_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_tls_version
        ADD CONSTRAINT deliveryservice_tls_version_deliveryservice_fkey FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;


IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'ip_address_server_fkey' AND table_name = 'ip_address') THEN
    --
    -- Name: ip_address_server_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY ip_address
        ADD CONSTRAINT ip_address_server_fkey FOREIGN KEY (interface, server) REFERENCES interface(name, server) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'ip_address_interface_fkey' AND table_name = 'ip_address') THEN
    --
    -- Name: ip_address_interface_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY ip_address
        ADD CONSTRAINT ip_address_interface_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'topology_cachegroup_cachegroup_fkey' AND table_name = 'topology_cachegroup') THEN
    --
    -- Name: topology_cachegroup_cachegroup_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY topology_cachegroup
        ADD CONSTRAINT topology_cachegroup_cachegroup_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroup(name) ON UPDATE CASCADE ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'topology_cachegroup_topology_fkey' AND table_name = 'topology_cachegroup') THEN
    --
    -- Name: topology_cachegroup_topology_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY topology_cachegroup
        ADD CONSTRAINT topology_cachegroup_topology_fkey FOREIGN KEY (topology) REFERENCES topology(name) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'topology_cachegroup_parents_child_fkey' AND table_name = 'topology_cachegroup_parents') THEN
    --
    -- Name: topology_cachegroup_parents_child_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY topology_cachegroup_parents
        ADD CONSTRAINT topology_cachegroup_parents_child_fkey FOREIGN KEY (child) REFERENCES topology_cachegroup(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'topology_cachegroup_parents_parent_fkey' AND table_name = 'topology_cachegroup_parents') THEN
    --
    -- Name: topology_cachegroup_parents_parent_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY topology_cachegroup_parents
        ADD CONSTRAINT topology_cachegroup_parents_parent_fkey FOREIGN KEY (parent) REFERENCES topology_cachegroup(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'fk_deliveryservice_id' AND table_name = 'deliveryservices_required_capability') THEN
    --
    -- Name: fk_deliveryservice_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservices_required_capability
        ADD CONSTRAINT fk_deliveryservice_id FOREIGN KEY (deliveryservice_id) REFERENCES deliveryservice(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'fk_required_capability' AND table_name = 'deliveryservices_required_capability') THEN
    --
    -- Name: fk_required_capability; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservices_required_capability
        ADD CONSTRAINT fk_required_capability FOREIGN KEY (required_capability) REFERENCES server_capability(name) ON UPDATE CASCADE ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'deliveryservice_topology_fkey' AND table_name = 'deliveryservice') THEN
    --
    -- Name: deliveryservice_topology_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT deliveryservice_topology_fkey FOREIGN KEY (topology) REFERENCES topology (name) ON UPDATE CASCADE ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'fk_server' AND table_name = 'server_server_capability') THEN
    --
    -- Name: fk_server; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server_server_capability
        ADD CONSTRAINT fk_server FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'fk_server_capability' AND table_name = 'server_server_capability') THEN
    --
    -- Name: fk_server_capability; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server_server_capability
        ADD CONSTRAINT fk_server_capability FOREIGN KEY (server_capability) REFERENCES server_capability(name) ON UPDATE CASCADE ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT  FROM information_schema.table_constraints WHERE constraint_name = 'fk_deliveryservice' AND table_name = 'deliveryservice_consistent_hash_query_param') THEN
    --
    -- Name: fk_deliveryservice; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE deliveryservice_consistent_hash_query_param
        ADD CONSTRAINT fk_deliveryservice FOREIGN KEY (deliveryservice_id) REFERENCES deliveryservice(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_atsprofile_atsparameters_atsprofile1' AND table_name = 'profile_parameter') THEN
    --
    -- Name: fk_atsprofile_atsparameters_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY profile_parameter
        ADD CONSTRAINT fk_atsprofile_atsparameters_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cdn1' AND table_name = 'deliveryservice') THEN
    --
    -- Name: fk_cdn1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT fk_cdn1 FOREIGN KEY (cdn_id) REFERENCES cdn(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cdn2' AND table_name = 'server') THEN
    --
    -- Name: fk_cdn2; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_cdn2 FOREIGN KEY (cdn_id) REFERENCES cdn(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cg_1' AND table_name = 'cachegroup') THEN
    --
    -- Name: fk_cg_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup
        ADD CONSTRAINT fk_cg_1 FOREIGN KEY (parent_cachegroup_id) REFERENCES cachegroup(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cg_param_cachegroup1' AND table_name = 'cachegroup_parameter') THEN
    --
    -- Name: fk_cg_param_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup_parameter
        ADD CONSTRAINT fk_cg_param_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cg_secondary' AND table_name = 'cachegroup') THEN
    --
    -- Name: fk_cg_secondary; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup
        ADD CONSTRAINT fk_cg_secondary FOREIGN KEY (secondary_parent_cachegroup_id) REFERENCES cachegroup(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cg_type1' AND table_name = 'cachegroup') THEN
    --
    -- Name: fk_cg_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup
        ADD CONSTRAINT fk_cg_type1 FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_contentserver_atsprofile1' AND table_name = 'server') THEN
    --
    -- Name: fk_contentserver_atsprofile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_contentserver_atsprofile1 FOREIGN KEY (profile) REFERENCES profile(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_contentserver_contentserverstatus1' AND table_name = 'server') THEN
    --
    -- Name: fk_contentserver_contentserverstatus1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_contentserver_contentserverstatus1 FOREIGN KEY (status) REFERENCES status(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_contentserver_contentservertype1' AND table_name = 'server') THEN
    --
    -- Name: fk_contentserver_contentservertype1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_contentserver_contentservertype1 FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_contentserver_phys_location1' AND table_name = 'server') THEN
    --
    -- Name: fk_contentserver_phys_location1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_contentserver_phys_location1 FOREIGN KEY (phys_location) REFERENCES phys_location(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cran_cachegroup1' AND table_name = 'asn') THEN
    --
    -- Name: fk_cran_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY asn
        ADD CONSTRAINT fk_cran_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_deliveryservice_profile1' AND table_name = 'deliveryservice') THEN
    --
    -- Name: fk_deliveryservice_profile1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT fk_deliveryservice_profile1 FOREIGN KEY (profile) REFERENCES profile(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_deliveryservice_type1' AND table_name = 'deliveryservice') THEN
    --
    -- Name: fk_deliveryservice_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT fk_deliveryservice_type1 FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_ds_to_cs_contentserver1' AND table_name = 'deliveryservice_server') THEN
    --
    -- Name: fk_ds_to_cs_contentserver1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_server
        ADD CONSTRAINT fk_ds_to_cs_contentserver1 FOREIGN KEY (server) REFERENCES server(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_ds_to_cs_deliveryservice1' AND table_name = 'deliveryservice_server') THEN
    --
    -- Name: fk_ds_to_cs_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_server
        ADD CONSTRAINT fk_ds_to_cs_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_ds_to_regex_deliveryservice1' AND table_name = 'deliveryservice_regex') THEN
    --
    -- Name: fk_ds_to_regex_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_regex
        ADD CONSTRAINT fk_ds_to_regex_deliveryservice1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_ds_to_regex_regex1' AND table_name = 'deliveryservice_regex') THEN
    --
    -- Name: fk_ds_to_regex_regex1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_regex
        ADD CONSTRAINT fk_ds_to_regex_regex1 FOREIGN KEY (regex) REFERENCES regex(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_ext_type' AND table_name = 'to_extension') THEN
    --
    -- Name: fk_ext_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY to_extension
        ADD CONSTRAINT fk_ext_type FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_federation_resolver1' AND table_name = 'federation_federation_resolver') THEN
    --
    -- Name: fk_federation_federation_resolver1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_federation_resolver
        ADD CONSTRAINT fk_federation_federation_resolver1 FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_mapping_type' AND table_name = 'federation_resolver') THEN
    --
    -- Name: fk_federation_mapping_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_resolver
        ADD CONSTRAINT fk_federation_mapping_type FOREIGN KEY (type) REFERENCES type(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_resolver_to_fed1' AND table_name = 'federation_federation_resolver') THEN
    --
    -- Name: fk_federation_resolver_to_fed1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_federation_resolver
        ADD CONSTRAINT fk_federation_resolver_to_fed1 FOREIGN KEY (federation_resolver) REFERENCES federation_resolver(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_tmuser_federation' AND table_name = 'federation_tmuser') THEN
    --
    -- Name: fk_federation_tmuser_federation; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_tmuser
        ADD CONSTRAINT fk_federation_tmuser_federation FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_tmuser_role' AND table_name = 'federation_tmuser') THEN
    --
    -- Name: fk_federation_tmuser_role; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_tmuser
        ADD CONSTRAINT fk_federation_tmuser_role FOREIGN KEY (role) REFERENCES role(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_tmuser_tmuser' AND table_name = 'federation_tmuser') THEN
    --
    -- Name: fk_federation_tmuser_tmuser; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_tmuser
        ADD CONSTRAINT fk_federation_tmuser_tmuser FOREIGN KEY (tm_user) REFERENCES tm_user(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_to_ds1' AND table_name = 'federation_deliveryservice') THEN
    --
    -- Name: fk_federation_to_ds1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_deliveryservice
        ADD CONSTRAINT fk_federation_to_ds1 FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_federation_to_fed1' AND table_name = 'federation_deliveryservice') THEN
    --
    -- Name: fk_federation_to_fed1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY federation_deliveryservice
        ADD CONSTRAINT fk_federation_to_fed1 FOREIGN KEY (federation) REFERENCES federation(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_hwinfo1' AND table_name = 'hwinfo') THEN
    --
    -- Name: fk_hwinfo1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY hwinfo
        ADD CONSTRAINT fk_hwinfo1 FOREIGN KEY (serverid) REFERENCES server(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'interface_server_fkey' AND table_name = 'interface') THEN
    --
    -- Name: interface_server_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE interface
        ADD CONSTRAINT interface_server_fkey FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_job_deliveryservice1' AND table_name = 'job') THEN
    --
    -- Name: fk_job_deliveryservice1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY job
        ADD CONSTRAINT fk_job_deliveryservice1 FOREIGN KEY (job_deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_job_user_id1' AND table_name = 'job') THEN
    --
    -- Name: fk_job_user_id1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY job
        ADD CONSTRAINT fk_job_user_id1 FOREIGN KEY (job_user) REFERENCES tm_user(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_log_1' AND table_name = 'log') THEN
    --
    -- Name: fk_log_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY log
        ADD CONSTRAINT fk_log_1 FOREIGN KEY (tm_user) REFERENCES tm_user(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_parameter' AND table_name = 'cachegroup_parameter') THEN
    --
    -- Name: fk_parameter; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup_parameter
        ADD CONSTRAINT fk_parameter FOREIGN KEY (parameter) REFERENCES parameter(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_phys_location_region' AND table_name = 'phys_location') THEN
    --
    -- Name: fk_phys_location_region; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY phys_location
        ADD CONSTRAINT fk_phys_location_region FOREIGN KEY (region) REFERENCES region(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_regex_type1' AND table_name = 'regex') THEN
    --
    -- Name: fk_regex_type1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY regex
        ADD CONSTRAINT fk_regex_type1 FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_region_division1' AND table_name = 'region') THEN
    --
    -- Name: fk_region_division1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY region
        ADD CONSTRAINT fk_region_division1 FOREIGN KEY (division) REFERENCES division(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_server_cachegroup1' AND table_name = 'server') THEN
    --
    -- Name: fk_server_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY server
        ADD CONSTRAINT fk_server_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_serverstatus_server1' AND table_name = 'servercheck') THEN
    --
    -- Name: fk_serverstatus_server1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY servercheck
        ADD CONSTRAINT fk_serverstatus_server1 FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_staticdnsentry_cachegroup1' AND table_name = 'staticdnsentry') THEN
    --
    -- Name: fk_staticdnsentry_cachegroup1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY staticdnsentry
        ADD CONSTRAINT fk_staticdnsentry_cachegroup1 FOREIGN KEY (cachegroup) REFERENCES cachegroup(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_staticdnsentry_ds' AND table_name = 'staticdnsentry') THEN
    --
    -- Name: fk_staticdnsentry_ds; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY staticdnsentry
        ADD CONSTRAINT fk_staticdnsentry_ds FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_staticdnsentry_type' AND table_name = 'staticdnsentry') THEN
    --
    -- Name: fk_staticdnsentry_type; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY staticdnsentry
        ADD CONSTRAINT fk_staticdnsentry_type FOREIGN KEY (type) REFERENCES type(id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_steering_target_delivery_service' AND table_name = 'steering_target') THEN
    --
    -- Name: fk_steering_target_delivery_service; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY steering_target
        ADD CONSTRAINT fk_steering_target_delivery_service FOREIGN KEY (deliveryservice) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_steering_target_target' AND table_name = 'steering_target') THEN
    --
    -- Name: fk_steering_target_target; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY steering_target
        ADD CONSTRAINT fk_steering_target_target FOREIGN KEY (target) REFERENCES deliveryservice(id) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_user_1' AND table_name = 'tm_user') THEN
    --
    -- Name: fk_user_1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY tm_user
        ADD CONSTRAINT fk_user_1 FOREIGN KEY (role) REFERENCES role(id) ON DELETE SET NULL;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_capability' AND table_name = 'api_capability') THEN
    --
    -- Name: fk_capability; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY api_capability
        ADD CONSTRAINT fk_capability FOREIGN KEY (capability) REFERENCES capability (name) ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'steering_target_type_fkey' AND table_name = 'steering_target') THEN
    --
    -- Name: steering_target_type_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY steering_target
        ADD CONSTRAINT steering_target_type_fkey FOREIGN KEY (type) REFERENCES type (id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_tenantid' AND table_name = 'deliveryservice') THEN
    --
    -- Name: fk_tenantid; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice
        ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_author' AND table_name = 'deliveryservice_request') THEN
    --
    -- Name: fk_author; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_request
        ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES tm_user(id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_assignee' AND table_name = 'deliveryservice_request') THEN
    --
    -- Name: fk_assignee; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_request
        ADD CONSTRAINT fk_assignee FOREIGN KEY (assignee_id) REFERENCES tm_user(id) ON DELETE SET NULL;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_last_edited_by' AND table_name = 'deliveryservice_request') THEN
    --
    -- Name: fk_last_edited_by; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_request
        ADD CONSTRAINT fk_last_edited_by FOREIGN KEY (last_edited_by_id) REFERENCES tm_user (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_author' AND table_name = 'deliveryservice_request_comment') THEN
    --
    -- Name: fk_author; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_request_comment
        ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES tm_user (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'origin_profile_fkey' AND table_name = 'origin') THEN
    --
    -- Name: origin_profile_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY origin
        ADD CONSTRAINT origin_profile_fkey FOREIGN KEY (profile) REFERENCES profile (id) ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'origin_deliveryservice_fkey' AND table_name = 'origin') THEN
    --
    -- Name: origin_deliveryservice_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY origin
        ADD CONSTRAINT origin_deliveryservice_fkey FOREIGN KEY (deliveryservice) REFERENCES deliveryservice (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'origin_coordinate_fkey' AND table_name = 'origin') THEN
    --
    -- Name: origin_coordinate_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY origin
        ADD CONSTRAINT origin_coordinate_fkey FOREIGN KEY (coordinate) REFERENCES coordinate (id) ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'origin_cachegroup_fkey' AND table_name = 'origin') THEN
    --
    -- Name: origin_cachegroup_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY origin
        ADD CONSTRAINT origin_cachegroup_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroup (id) ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'origin_tenant_fkey' AND table_name = 'origin') THEN
    --
    -- Name: origin_tenant_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY origin
        ADD CONSTRAINT origin_tenant_fkey FOREIGN KEY (tenant) REFERENCES tenant (id) ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fkey_lock_cdn' AND table_name = 'cdn_lock') THEN
    --
    -- Name: fk_lock_cdn; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cdn_lock
        ADD CONSTRAINT fk_lock_cdn FOREIGN KEY ("cdn") REFERENCES cdn(name) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fkey_lock_username' AND table_name = 'cdn_lock') THEN
    --
    -- Name: fk_lock_username; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cdn_lock
        ADD CONSTRAINT fk_lock_username FOREIGN KEY ("username") REFERENCES tm_user(username) ON DELETE CASCADE ON UPDATE CASCADE;
END IF;


IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'cachegroup_coordinate_fkey' AND table_name = 'cachegroup') THEN
    --
    -- Name: cachegroup_coordinate_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup
        ADD CONSTRAINT cachegroup_coordinate_fkey FOREIGN KEY (coordinate) REFERENCES coordinate (id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_primary_cg' AND table_name = 'cachegroup_fallbacks') THEN
    --
    -- Name: fk_primary_cg; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup_fallbacks
        ADD CONSTRAINT fk_primary_cg FOREIGN KEY (primary_cg) REFERENCES cachegroup (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_backup_cg' AND table_name = 'cachegroup_fallbacks') THEN
    --
    -- Name: fk_backup_cg; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup_fallbacks
        ADD CONSTRAINT fk_backup_cg FOREIGN KEY (backup_cg) REFERENCES cachegroup (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'cachegroup_localization_method_cachegroup_fkey' AND table_name = 'cachegroup_localization_method') THEN
    --
    -- Name: cachegroup_localization_method_cachegroup_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY cachegroup_localization_method
        ADD CONSTRAINT cachegroup_localization_method_cachegroup_fkey FOREIGN KEY (cachegroup) REFERENCES cachegroup (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_deliveryservice_request' AND table_name = 'deliveryservice_request_comment') THEN
    --
    -- Name: fk_deliveryservice_request; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY deliveryservice_request_comment
        ADD CONSTRAINT fk_deliveryservice_request FOREIGN KEY (deliveryservice_request_id) REFERENCES deliveryservice_request (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_cdn1' AND table_name = 'profile') THEN
    --
    -- Name: fk_cdn1; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY profile
        ADD CONSTRAINT fk_cdn1 FOREIGN KEY (cdn) REFERENCES cdn (id) ON UPDATE RESTRICT ON DELETE RESTRICT;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_role_id' AND table_name = 'role_capability') THEN
    --
    -- Name: fk_role_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY role_capability
        ADD CONSTRAINT fk_role_id FOREIGN KEY (role_id) REFERENCES role (id) ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'snapshot_cdn_fkey' AND table_name = 'snapshot') THEN
    --
    -- Name: snapshot_cdn_fkey; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY snapshot
        ADD CONSTRAINT snapshot_cdn_fkey FOREIGN KEY (cdn) REFERENCES cdn (name) ON UPDATE CASCADE ON DELETE CASCADE;
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_parentid' AND table_name = 'tenant') THEN
    --
    -- Name: fk_parentid; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY tenant
        ADD CONSTRAINT fk_parentid FOREIGN KEY (parent_id) REFERENCES tenant (id);
END IF;

IF NOT EXISTS (SELECT FROM information_schema.table_constraints WHERE constraint_name = 'fk_tenantid' AND table_name = 'tm_user') THEN
    --
    -- Name: fk_tenantid; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
    --

    ALTER TABLE ONLY tm_user
        ADD CONSTRAINT fk_tenantid FOREIGN KEY (tenant_id) REFERENCES tenant (id) MATCH FULL;
END IF;
END$$;


--
-- Name: public; Type: ACL; Schema: -; Owner: traffic_ops
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM traffic_ops;
GRANT ALL ON SCHEMA public TO traffic_ops;
GRANT ALL ON SCHEMA public TO PUBLIC;

--
-- Name: cdni_capabilities; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdni_capabilities (
    id bigint NOT NULL,
    type text NOT NULL,
    ucdn text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_capabilities OWNER TO traffic_ops;

--
-- Name: cdni_capabilities_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE cdni_capabilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE cdni_capabilities_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_capabilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cdni_capabilities_id_seq OWNED BY cdni_capabilities.id;

--
-- Name: cdni_footprints; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdni_footprints (
    id bigint NOT NULL,
    footprint_type text NOT NULL,
    footprint_value text[] NOT NULL,
    ucdn text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_footprints OWNER TO traffic_ops;

--
-- Name: cdni_footprints_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE cdni_footprints_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE cdni_footprints_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_footprints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cdni_footprints_id_seq OWNED BY cdni_footprints.id;

--
-- Name: cdni_limits; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdni_limits (
    id bigint NOT NULL,
    limit_id text NOT NULL,
    scope_type text,
    scope_value text[],
    limit_type text NOT NULL,
    maximum_hard bigint NOT NULL,
    maximum_soft bigint NOT NULL,
    telemetry_id text NOT NULL,
    telemetry_metric text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_limits OWNER TO traffic_ops;

--
-- Name: cdni_limits_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE cdni_limits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE cdni_limits_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_limits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE cdni_limits_id_seq OWNED BY cdni_limits.id;

--
-- Name: cdni_capabilities; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdni_capabilities (
    id bigserial NOT NULL,
    type text NOT NULL,
    ucdn text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_capabilities OWNER TO traffic_ops;

--
-- Name: cdni_footprints; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdni_footprints (
    id bigserial NOT NULL,
    footprint_type text NOT NULL,
    footprint_value text[] NOT NULL,
    ucdn text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_footprints OWNER TO traffic_ops;

--
-- Name: cdni_telemetry; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdni_telemetry (
    id text NOT NULL,
    type text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    configuration_url text DEFAULT ''::text
);

ALTER TABLE cdni_telemetry OWNER TO traffic_ops;

--
-- Name: cdni_telemetry_metrics; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE cdni_telemetry_metrics (
    name text NOT NULL,
    time_granularity bigint NOT NULL,
    data_percentile bigint NOT NULL,
    latency integer NOT NULL,
    telemetry_id text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_telemetry_metrics OWNER TO traffic_ops;

--
-- Name: cdni_limits; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdni_limits (
    id bigserial NOT NULL,
    limit_id text NOT NULL,
    scope_type text,
    scope_value text[],
    limit_type text NOT NULL,
    maximum_hard bigint NOT NULL,
    maximum_soft bigint NOT NULL,
    telemetry_id text NOT NULL,
    telemetry_metric text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_limits OWNER TO traffic_ops;

--
-- Name: cdn_lock_user; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdn_lock_user (
    owner text NOT NULL,
    cdn text NOT NULL,
    username text NOT NULL
);

ALTER TABLE cdn_lock_user OWNER TO traffic_ops;

--
-- Name: server_profile; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS server_profile (
    server bigint NOT NULL,
    profile_name text NOT NULL,
    priority int NOT NULL CHECK (priority >= 0)
);

ALTER TABLE server_profile OWNER TO traffic_ops;

--
-- Name: cdni_capability_updates; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE IF NOT EXISTS cdni_capability_updates (
    id bigserial NOT NULL,
    request_type text NOT NULL,
    ucdn text NOT NULL,
    host text,
    data json NOT NULL,
    async_status_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE cdni_capability_updates OWNER TO traffic_ops;

--
-- Name: cdni_capabilities id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_capabilities ALTER COLUMN id SET DEFAULT nextval('cdni_capabilities_id_seq'::regclass);

--
-- Name: cdni_footprints id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_footprints ALTER COLUMN id SET DEFAULT nextval('cdni_footprints_id_seq'::regclass);

--
-- Name: cdni_limits id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_limits ALTER COLUMN id SET DEFAULT nextval('cdni_limits_id_seq'::regclass);

--
-- Name: cdn_lock cdn_lock_cdn_username_unique; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE cdn_lock ADD CONSTRAINT cdn_lock_cdn_username_unique UNIQUE (username, cdn);

--
-- Name: cdni_capabilities pk_cdni_capabilities; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_capabilities
    ADD CONSTRAINT pk_cdni_capabilities PRIMARY KEY (id);

--
-- Name: cdni_footprints pk_cdni_footprints; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_footprints
    ADD CONSTRAINT pk_cdni_footprints PRIMARY KEY (id);

--
-- Name: cdni_limits pk_cdni_limits; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_limits
    ADD CONSTRAINT pk_cdni_limits PRIMARY KEY (id);

--
-- Name: cdni_telemetry pk_cdni_telemetry; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_telemetry
    ADD CONSTRAINT pk_cdni_telemetry PRIMARY KEY (id);

--
-- Name: cdni_telemetry_metrics pk_cdni_telemetry_metrics; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_telemetry_metrics
    ADD CONSTRAINT pk_cdni_telemetry_metrics PRIMARY KEY (name);

--
-- Name: cdn_lock_user pk_cdn_lock_user; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdn_lock_user
    ADD CONSTRAINT pk_cdn_lock_user PRIMARY KEY (owner, cdn, username);

--
-- Name: server_profile pk_server_profile; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server_profile
    ADD CONSTRAINT pk_server_profile PRIMARY KEY (profile_name, server);

--
-- Name: cdni_capability_updates pk_cdni_capability_updates; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_capability_updates
    ADD CONSTRAINT pk_cdni_capability_updates PRIMARY KEY (id);

--
-- Name: cdni_footprints fk_cdni_footprint_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_footprints
    ADD CONSTRAINT fk_cdni_footprint_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: cdni_limits fk_cdni_limits_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_limits
    ADD CONSTRAINT fk_cdni_limits_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: cdni_limits fk_cdni_limits_telemetry; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_limits
    ADD CONSTRAINT fk_cdni_limits_telemetry FOREIGN KEY (telemetry_id) REFERENCES cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: cdni_telemetry fk_cdni_telemetry_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_telemetry
    ADD CONSTRAINT fk_cdni_telemetry_capabilities FOREIGN KEY (capability_id) REFERENCES cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: cdni_telemetry_metrics fk_cdni_telemetry_metrics_telemetry; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_telemetry_metrics
    ADD CONSTRAINT fk_cdni_telemetry_metrics_telemetry FOREIGN KEY (telemetry_id) REFERENCES cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: cdn_lock_user fk_shared_username; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdn_lock_user
    ADD CONSTRAINT fk_shared_username FOREIGN KEY (username) REFERENCES tm_user(username);

--
-- Name: cdn_lock_user fk_owner; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdn_lock_user
    ADD CONSTRAINT fk_owner FOREIGN KEY (owner, cdn) REFERENCES cdn_lock(username, cdn) ON DELETE CASCADE;

--
-- Name: server_profile fk_server_id; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server_profile
    ADD CONSTRAINT fk_server_id FOREIGN KEY (server) REFERENCES server(id) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Name: server_profile fk_server_profile_name_profile; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY server_profile
    ADD CONSTRAINT fk_server_profile_name_profile FOREIGN KEY (profile_name) REFERENCES profile(name) ON UPDATE CASCADE ON DELETE RESTRICT;

--
-- Name: cdni_capability_updates fk_cdni_capability_updates_async; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY cdni_capability_updates
    ADD CONSTRAINT fk_cdni_capability_updates_async FOREIGN KEY (async_status_id) REFERENCES async_status(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--
