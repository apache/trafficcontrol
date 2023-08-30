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
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cdni_capabilities; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE public.cdni_capabilities (
    id bigint NOT NULL,
    type text NOT NULL,
    ucdn text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.cdni_capabilities OWNER TO traffic_ops;

--
-- Name: cdni_capabilities_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE public.cdni_capabilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cdni_capabilities_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_capabilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE public.cdni_capabilities_id_seq OWNED BY public.cdni_capabilities.id;


--
-- Name: cdni_footprints; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE public.cdni_footprints (
    id bigint NOT NULL,
    footprint_type text NOT NULL,
    footprint_value text[] NOT NULL,
    ucdn text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.cdni_footprints OWNER TO traffic_ops;

--
-- Name: cdni_footprints_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE public.cdni_footprints_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cdni_footprints_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_footprints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE public.cdni_footprints_id_seq OWNED BY public.cdni_footprints.id;


--
-- Name: cdni_limits; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE public.cdni_limits (
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


ALTER TABLE public.cdni_limits OWNER TO traffic_ops;

--
-- Name: cdni_limits_id_seq; Type: SEQUENCE; Schema: public; Owner: traffic_ops
--

CREATE SEQUENCE public.cdni_limits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cdni_limits_id_seq OWNER TO traffic_ops;

--
-- Name: cdni_limits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: traffic_ops
--

ALTER SEQUENCE public.cdni_limits_id_seq OWNED BY public.cdni_limits.id;


--
-- Name: cdni_telemetry; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE public.cdni_telemetry (
    id text NOT NULL,
    type text NOT NULL,
    capability_id bigint NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL,
    configuration_url text DEFAULT ''::text
);


ALTER TABLE public.cdni_telemetry OWNER TO traffic_ops;

--
-- Name: cdni_telemetry_metrics; Type: TABLE; Schema: public; Owner: traffic_ops
--

CREATE TABLE public.cdni_telemetry_metrics (
    name text NOT NULL,
    time_granularity bigint NOT NULL,
    data_percentile bigint NOT NULL,
    latency integer NOT NULL,
    telemetry_id text NOT NULL,
    last_updated timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.cdni_telemetry_metrics OWNER TO traffic_ops;

--
-- Name: cdni_capabilities id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_capabilities ALTER COLUMN id SET DEFAULT nextval('public.cdni_capabilities_id_seq'::regclass);


--
-- Name: cdni_footprints id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_footprints ALTER COLUMN id SET DEFAULT nextval('public.cdni_footprints_id_seq'::regclass);


--
-- Name: cdni_limits id; Type: DEFAULT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_limits ALTER COLUMN id SET DEFAULT nextval('public.cdni_limits_id_seq'::regclass);


--
-- Data for Name: cdni_capabilities; Type: TABLE DATA; Schema: public; Owner: traffic_ops
--

COPY public.cdni_capabilities (id, type, ucdn, last_updated) FROM stdin;
\.


--
-- Data for Name: cdni_footprints; Type: TABLE DATA; Schema: public; Owner: traffic_ops
--

COPY public.cdni_footprints (id, footprint_type, footprint_value, ucdn, capability_id, last_updated) FROM stdin;
\.


--
-- Data for Name: cdni_limits; Type: TABLE DATA; Schema: public; Owner: traffic_ops
--

COPY public.cdni_limits (id, limit_id, scope_type, scope_value, limit_type, maximum_hard, maximum_soft, telemetry_id, telemetry_metric, capability_id, last_updated) FROM stdin;
\.


--
-- Data for Name: cdni_telemetry; Type: TABLE DATA; Schema: public; Owner: traffic_ops
--

COPY public.cdni_telemetry (id, type, capability_id, last_updated, configuration_url) FROM stdin;
\.


--
-- Data for Name: cdni_telemetry_metrics; Type: TABLE DATA; Schema: public; Owner: traffic_ops
--

COPY public.cdni_telemetry_metrics (name, time_granularity, data_percentile, latency, telemetry_id, last_updated) FROM stdin;
\.


--
-- Name: cdni_capabilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: traffic_ops
--

SELECT pg_catalog.setval('public.cdni_capabilities_id_seq', 1, false);


--
-- Name: cdni_footprints_id_seq; Type: SEQUENCE SET; Schema: public; Owner: traffic_ops
--

SELECT pg_catalog.setval('public.cdni_footprints_id_seq', 1, false);


--
-- Name: cdni_limits_id_seq; Type: SEQUENCE SET; Schema: public; Owner: traffic_ops
--

SELECT pg_catalog.setval('public.cdni_limits_id_seq', 1, false);


--
-- Name: cdni_capabilities pk_cdni_capabilities; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_capabilities
    ADD CONSTRAINT pk_cdni_capabilities PRIMARY KEY (id);


--
-- Name: cdni_footprints pk_cdni_footprints; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_footprints
    ADD CONSTRAINT pk_cdni_footprints PRIMARY KEY (id);


--
-- Name: cdni_limits pk_cdni_limits; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_limits
    ADD CONSTRAINT pk_cdni_limits PRIMARY KEY (id);


--
-- Name: cdni_telemetry pk_cdni_telemetry; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_telemetry
    ADD CONSTRAINT pk_cdni_telemetry PRIMARY KEY (id);


--
-- Name: cdni_telemetry_metrics pk_cdni_telemetry_metrics; Type: CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_telemetry_metrics
    ADD CONSTRAINT pk_cdni_telemetry_metrics PRIMARY KEY (name);


--
-- Name: cdni_footprints fk_cdni_footprint_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_footprints
    ADD CONSTRAINT fk_cdni_footprint_capabilities FOREIGN KEY (capability_id) REFERENCES public.cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cdni_limits fk_cdni_limits_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_limits
    ADD CONSTRAINT fk_cdni_limits_capabilities FOREIGN KEY (capability_id) REFERENCES public.cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cdni_limits fk_cdni_limits_telemetry; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_limits
    ADD CONSTRAINT fk_cdni_limits_telemetry FOREIGN KEY (telemetry_id) REFERENCES public.cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cdni_telemetry fk_cdni_telemetry_capabilities; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_telemetry
    ADD CONSTRAINT fk_cdni_telemetry_capabilities FOREIGN KEY (capability_id) REFERENCES public.cdni_capabilities(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: cdni_telemetry_metrics fk_cdni_telemetry_metrics_telemetry; Type: FK CONSTRAINT; Schema: public; Owner: traffic_ops
--

ALTER TABLE ONLY public.cdni_telemetry_metrics
    ADD CONSTRAINT fk_cdni_telemetry_metrics_telemetry FOREIGN KEY (telemetry_id) REFERENCES public.cdni_telemetry(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

