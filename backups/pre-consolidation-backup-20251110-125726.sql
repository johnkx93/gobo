--
-- PostgreSQL database dump
--

\restrict rhr3LFJmf1heVyUtbL6deJ8y0WLHslUbwaleBYdJNBCIwl220c2UJ9mxegsaxcZ

-- Dumped from database version 16.10
-- Dumped by pg_dump version 16.10

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

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: audit_action; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.audit_action AS ENUM (
    'CREATE',
    'UPDATE',
    'DELETE'
);


ALTER TYPE public.audit_action OWNER TO postgres;

--
-- Name: update_admins_updated_at(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_admins_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_admins_updated_at() OWNER TO postgres;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admins (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    role character varying(50) DEFAULT 'admin'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.admins OWNER TO postgres;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
)
PARTITION BY RANGE (created_at);


ALTER TABLE public.audit_logs OWNER TO postgres;

--
-- Name: audit_logs_2025_11; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs_2025_11 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.audit_logs_2025_11 OWNER TO postgres;

--
-- Name: audit_logs_2025_12; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs_2025_12 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.audit_logs_2025_12 OWNER TO postgres;

--
-- Name: audit_logs_2026_01; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs_2026_01 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.audit_logs_2026_01 OWNER TO postgres;

--
-- Name: audit_logs_2026_02; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs_2026_02 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.audit_logs_2026_02 OWNER TO postgres;

--
-- Name: audit_logs_default; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs_default (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    action public.audit_action NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id uuid NOT NULL,
    old_data jsonb,
    new_data jsonb,
    request_id character varying(100),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.audit_logs_default OWNER TO postgres;

--
-- Name: error_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
)
PARTITION BY RANGE (created_at);


ALTER TABLE public.error_logs OWNER TO postgres;

--
-- Name: error_logs_2025_11; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs_2025_11 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.error_logs_2025_11 OWNER TO postgres;

--
-- Name: error_logs_2025_12; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs_2025_12 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.error_logs_2025_12 OWNER TO postgres;

--
-- Name: error_logs_2026_01; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs_2026_01 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.error_logs_2026_01 OWNER TO postgres;

--
-- Name: error_logs_2026_02; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs_2026_02 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.error_logs_2026_02 OWNER TO postgres;

--
-- Name: error_logs_default; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.error_logs_default (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    request_id character varying(100),
    error_type character varying(100) NOT NULL,
    error_message text NOT NULL,
    stack_trace text,
    request_path character varying(255),
    request_method character varying(10),
    ip_address character varying(45),
    user_agent text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.error_logs_default OWNER TO postgres;

--
-- Name: menu_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.menu_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    parent_id uuid,
    code character varying(100) NOT NULL,
    label character varying(255) NOT NULL,
    icon character varying(50),
    path character varying(255),
    permission_id uuid,
    order_index integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.menu_items OWNER TO postgres;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    order_number character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    total_amount numeric(10,2) DEFAULT 0.00 NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
)
PARTITION BY RANGE (created_at);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_2025; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders_2025 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    order_number character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    total_amount numeric(10,2) DEFAULT 0.00 NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.orders_2025 OWNER TO postgres;

--
-- Name: orders_2026; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders_2026 (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    order_number character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    total_amount numeric(10,2) DEFAULT 0.00 NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.orders_2026 OWNER TO postgres;

--
-- Name: orders_default; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders_default (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    order_number character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    total_amount numeric(10,2) DEFAULT 0.00 NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.orders_default OWNER TO postgres;

--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    code character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    category character varying(50) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.permissions OWNER TO postgres;

--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_permissions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    role character varying(50) NOT NULL,
    permission_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.role_permissions OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: audit_logs_2025_11; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ATTACH PARTITION public.audit_logs_2025_11 FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');


--
-- Name: audit_logs_2025_12; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ATTACH PARTITION public.audit_logs_2025_12 FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');


--
-- Name: audit_logs_2026_01; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ATTACH PARTITION public.audit_logs_2026_01 FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');


--
-- Name: audit_logs_2026_02; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ATTACH PARTITION public.audit_logs_2026_02 FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');


--
-- Name: audit_logs_default; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ATTACH PARTITION public.audit_logs_default DEFAULT;


--
-- Name: error_logs_2025_11; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs ATTACH PARTITION public.error_logs_2025_11 FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');


--
-- Name: error_logs_2025_12; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs ATTACH PARTITION public.error_logs_2025_12 FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');


--
-- Name: error_logs_2026_01; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs ATTACH PARTITION public.error_logs_2026_01 FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');


--
-- Name: error_logs_2026_02; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs ATTACH PARTITION public.error_logs_2026_02 FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');


--
-- Name: error_logs_default; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs ATTACH PARTITION public.error_logs_default DEFAULT;


--
-- Name: orders_2025; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ATTACH PARTITION public.orders_2025 FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');


--
-- Name: orders_2026; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ATTACH PARTITION public.orders_2026 FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2027-01-01 00:00:00+00');


--
-- Name: orders_default; Type: TABLE ATTACH; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ATTACH PARTITION public.orders_default DEFAULT;


--
-- Data for Name: admins; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.admins (id, email, username, password_hash, first_name, last_name, role, is_active, created_at, updated_at) FROM stdin;
2b46042d-73e5-48dc-8a59-410220ff524d	admin@example.com	superadmin	$2a$10$rZ0qV3wF8Xr0mE5yB.qI0.WxKGZGQXZJYH0wLZvMZN5Q3wYX0Q0Qq	Super	Admin	super_admin	t	2025-11-10 03:18:46.233583+00	2025-11-10 03:18:46.233583+00
10db0711-722f-44b9-81c4-2448c3adc21f	admin0@investorinnovative.org	admin_The Notorious D.O.G.4510	$2a$10$AwPmI0g2IaATLf05grE0xeQPMeprBLIkLlpnk1GziPTg.KJ01Z7jC	Roman	Spencer	super_admin	t	2025-11-10 03:24:51.423308+00	2025-11-10 03:24:51.423308+00
ca3a2915-aba4-42b1-8538-dd60628aa165	admin1@seniorfunctionalities.com	admin_PantherShakeer106911	$2a$10$AwPmI0g2IaATLf05grE0xeQPMeprBLIkLlpnk1GziPTg.KJ01Z7jC	Vincenzo	McCullough	moderator	t	2025-11-10 03:24:51.425184+00	2025-11-10 03:24:51.425184+00
b2b0dfaa-0c07-46cc-a368-642144be011f	admin2@customerefficient.name	admin_muskrat_922	$2a$10$AwPmI0g2IaATLf05grE0xeQPMeprBLIkLlpnk1GziPTg.KJ01Z7jC	Oral	Considine	moderator	t	2025-11-10 03:24:51.426085+00	2025-11-10 03:24:51.426085+00
\.


--
-- Data for Name: audit_logs_2025_11; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audit_logs_2025_11 (id, user_id, action, entity_type, entity_id, old_data, new_data, request_id, ip_address, user_agent, metadata, created_at) FROM stdin;
fd87af64-c76e-44fd-8f42-b2c211204c1f	\N	CREATE	users	3771ce44-478c-45b5-aca4-06e018883f9a	\N	{"id": "3771ce44-478c-45b5-aca4-06e018883f9a", "email": "asd@dsa20.com", "username": "mmm4", "last_name": null, "created_at": "2025-11-07T16:05:52.517015+08:00", "first_name": null, "updated_at": "2025-11-07T16:05:52.517015+08:00"}	133bd5ba-665f-43a5-82c4-253f6dffe71d	127.0.0.1	RapidAPI/4.4.3 (Macintosh; OS X/26.1.0) GCDHTTPRequest	\N	2025-11-07 08:05:52.518928+00
\.


--
-- Data for Name: audit_logs_2025_12; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audit_logs_2025_12 (id, user_id, action, entity_type, entity_id, old_data, new_data, request_id, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: audit_logs_2026_01; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audit_logs_2026_01 (id, user_id, action, entity_type, entity_id, old_data, new_data, request_id, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: audit_logs_2026_02; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audit_logs_2026_02 (id, user_id, action, entity_type, entity_id, old_data, new_data, request_id, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: audit_logs_default; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.audit_logs_default (id, user_id, action, entity_type, entity_id, old_data, new_data, request_id, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: error_logs_2025_11; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.error_logs_2025_11 (id, user_id, request_id, error_type, error_message, stack_trace, request_path, request_method, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: error_logs_2025_12; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.error_logs_2025_12 (id, user_id, request_id, error_type, error_message, stack_trace, request_path, request_method, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: error_logs_2026_01; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.error_logs_2026_01 (id, user_id, request_id, error_type, error_message, stack_trace, request_path, request_method, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: error_logs_2026_02; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.error_logs_2026_02 (id, user_id, request_id, error_type, error_message, stack_trace, request_path, request_method, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: error_logs_default; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.error_logs_default (id, user_id, request_id, error_type, error_message, stack_trace, request_path, request_method, ip_address, user_agent, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: menu_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.menu_items (id, parent_id, code, label, icon, path, permission_id, order_index, is_active, created_at, updated_at) FROM stdin;
168f960e-c970-48c1-a538-0a5261cbfe85	\N	users	User Management	users	\N	\N	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
723c2008-4de5-4a64-9c21-9a7c5e50f23a	\N	orders	Order Management	shopping-cart	\N	\N	2	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
f868b554-5136-4f83-9e93-454c64e6c5f0	\N	admins	Admin Management	shield	\N	0d67a644-1532-4b0b-a969-26e97f79bd31	3	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
1d9949e5-ec45-43bf-b145-22d8766449e3	\N	settings	Settings	settings	\N	\N	4	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
897ba22e-df7e-4cdb-8edd-24c51c8b7082	\N	analytics	Analytics	chart-bar	\N	\N	5	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
26369b41-0e1d-4bf8-b2f4-a61785675777	168f960e-c970-48c1-a538-0a5261cbfe85	users-create	Create User	\N	/admin/users/create	a819da74-fc0c-4957-854a-ebd26d380964	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
7482e01c-915c-42e6-9f62-380d9f541c17	168f960e-c970-48c1-a538-0a5261cbfe85	users-list	User List	\N	/admin/users	6c349a37-ebf2-4027-aec5-9addba8c68b2	2	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
d6909bc4-ab45-46e5-8c5a-d49b01e538a2	723c2008-4de5-4a64-9c21-9a7c5e50f23a	orders-list	Order List	\N	/admin/orders	39b03ece-3e5b-4f74-8890-4f27dfc22cbb	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
039728bf-702d-4f0a-9c17-963f82e1fb10	723c2008-4de5-4a64-9c21-9a7c5e50f23a	orders-update	Update Orders	\N	/admin/orders/bulk-update	c2b06b19-861d-4dd7-ba19-a5c35c8a9454	2	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
3560e89f-d047-4be5-a61c-24ddd2b564dc	f868b554-5136-4f83-9e93-454c64e6c5f0	admins-list	Admin List	\N	/admin/admins	f821ee73-995f-4f98-8d98-834143907b52	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
6d81cfa0-c1b1-40a4-9b6d-ba22c3653e71	1d9949e5-ec45-43bf-b145-22d8766449e3	settings-general	General Settings	\N	/admin/settings/general	4cc40686-86b1-4d38-9d99-a3f1a0cf20c3	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
4ab5cade-7d2e-453d-bac2-308995c794de	1d9949e5-ec45-43bf-b145-22d8766449e3	settings-security	Security	\N	/admin/settings/security	2f7e9085-c8ce-4de0-8b71-47a2d130c2d1	2	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
e1b863bd-e1b7-4b7f-8e1e-9fe22d98d39a	897ba22e-df7e-4cdb-8edd-24c51c8b7082	analytics-dashboard	Dashboard	\N	/admin/analytics/dashboard	2765ad95-382a-47c8-90c8-ad680aabc645	1	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
1de639d8-e41e-4be5-800d-79002f87dbe9	897ba22e-df7e-4cdb-8edd-24c51c8b7082	analytics-reports	Reports	\N	/admin/analytics/reports	01f402b6-ee4f-4060-bb82-62280927e8ca	2	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
\.


--
-- Data for Name: orders_2025; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders_2025 (id, user_id, order_number, status, total_amount, notes, created_at, updated_at) FROM stdin;
bc9cc92a-5cf3-4ed3-a874-90ccf2fcbecd	7fe9c1b5-b5d2-4674-b4f8-d20b5bf07ce4	ORD-20251105-269624	shipped	778.98	\N	2025-11-05 08:34:19.367977+00	2025-11-05 08:34:19.367977+00
b27b110d-261e-4de1-8211-121dc07bda77	7fe9c1b5-b5d2-4674-b4f8-d20b5bf07ce4	ORD-20251105-925921	pending	954.49	\N	2025-11-05 08:34:19.370645+00	2025-11-05 08:34:19.370645+00
b6dc3539-3c55-4965-976b-08dd7dc1b320	4a5bbae2-90e9-4e15-8607-7aee1ca1fad7	ORD-20251105-217731	shipped	507.42	Evenings in Austin invite quieter woman.	2025-11-05 08:34:19.371243+00	2025-11-05 08:34:19.371243+00
febc0694-1261-40cf-a0e0-8b6218ed55d4	5538dfac-492d-45c7-98f5-8f901759f313	ORD-20251105-252514	cancelled	52.85	\N	2025-11-05 08:34:19.371746+00	2025-11-05 08:34:19.371746+00
c07f2cf1-aab7-493e-a6df-caf9acdab0af	dd6a7abb-5a19-42c4-b3fb-0b0f94010e1e	ORD-20251105-152631	shipped	691.90	Before launch, you fight now.	2025-11-05 08:34:19.372198+00	2025-11-05 08:34:19.372198+00
eb57d6f7-fb7d-4f10-9f31-613135fc98cf	bfbd9bac-abcd-4cec-89e1-df20a9a99238	ORD-20251105-945829	completed	539.39	Continuously measure the number and embarrass the outliers.	2025-11-05 08:34:19.372686+00	2025-11-05 08:34:19.372686+00
ba23857f-338e-4eff-a895-00cea016dd52	9e280b82-1658-46b2-9ece-2064f12e9eb7	ORD-20251105-132218	cancelled	372.74	\N	2025-11-05 08:34:19.373139+00	2025-11-05 08:34:19.373139+00
cfdd6c8a-5b8e-48d8-a221-defff6a57fb1	eaa88a4d-a54f-42ee-8643-82362a1ab2b0	ORD-20251105-168357	pending	823.37	\N	2025-11-05 08:34:19.373487+00	2025-11-05 08:34:19.373487+00
7571bc16-ee83-4b54-b31a-4dbbea9885cd	adc3f01c-9b7d-4fd6-a103-7bf890fc2b22	ORD-20251105-624769	completed	482.67	\N	2025-11-05 08:34:19.373908+00	2025-11-05 08:34:19.373908+00
01393f78-7487-464e-8595-dc78e4f46d7e	5538dfac-492d-45c7-98f5-8f901759f313	ORD-20251105-373072	pending	98.37	\N	2025-11-05 08:34:19.374337+00	2025-11-05 08:34:19.374337+00
2d1b672c-10f0-4017-be08-620ce30211aa	4a5bbae2-90e9-4e15-8607-7aee1ca1fad7	ORD-20251105-513800	pending	873.82	\N	2025-11-05 08:34:19.374954+00	2025-11-05 08:34:19.374954+00
42a9b567-d93c-486d-842f-ba1a85336e0c	adc3f01c-9b7d-4fd6-a103-7bf890fc2b22	ORD-20251105-288532	cancelled	735.48	\N	2025-11-05 08:34:19.37559+00	2025-11-05 08:34:19.37559+00
858602ef-dea1-4f30-a9e6-634252b299ad	bfbd9bac-abcd-4cec-89e1-df20a9a99238	ORD-20251105-182818	cancelled	43.09	\N	2025-11-05 08:34:19.37611+00	2025-11-05 08:34:19.37611+00
9798d929-8829-4db9-8e88-9dba7abca569	5538dfac-492d-45c7-98f5-8f901759f313	ORD-20251105-420487	completed	499.14	\N	2025-11-05 08:34:19.376672+00	2025-11-05 08:34:19.376672+00
81db6c40-ef97-48d0-a519-01910c22dbb5	4825d06b-421a-4a86-ab8d-303118a024dc	ORD-20251105-658511	cancelled	480.92	\N	2025-11-05 08:34:19.377267+00	2025-11-05 08:34:19.377267+00
d387f543-c7af-44f6-9671-7944e3ada3de	adc3f01c-9b7d-4fd6-a103-7bf890fc2b22	ORD-20251105-785073	completed	281.05	\N	2025-11-05 08:34:19.377626+00	2025-11-05 08:34:19.377626+00
d43ebd8e-72dd-409f-84c2-6828320866c4	4825d06b-421a-4a86-ab8d-303118a024dc	ORD-20251105-489474	cancelled	478.59	Onward to better time!	2025-11-05 08:34:19.377941+00	2025-11-05 08:34:19.377941+00
53f149c0-dfdb-43d3-a6e4-8027a7656198	eaa88a4d-a54f-42ee-8643-82362a1ab2b0	ORD-20251105-867309	pending	420.06	\N	2025-11-05 08:34:19.378256+00	2025-11-05 08:34:19.378256+00
9d5497ad-69bb-4845-8b28-45f6f8487d55	5538dfac-492d-45c7-98f5-8f901759f313	ORD-20251105-166419	cancelled	563.35	\N	2025-11-05 08:34:19.378583+00	2025-11-05 08:34:19.378583+00
fe7fc9e2-4a4e-41ad-996a-71c12640dabf	a9dcf66a-cacd-4a2d-bd5d-7cb3105611e0	ORD-20251105-747119	shipped	136.60	Practice place drills regularly.	2025-11-05 08:34:19.378901+00	2025-11-05 08:34:19.378901+00
0b94bcf3-acb2-4166-b2f9-d54088f3a207	7a64aa6f-0c20-4e9c-b563-0cd24e00c5c7	ORD-20251110-828630	pending	761.06	\N	2025-11-10 03:24:51.483329+00	2025-11-10 03:24:51.483329+00
a16f7b76-5179-4feb-8ff3-df8ec5cade83	d4ec9cf1-fc39-47fd-97d3-22e59f6b2f21	ORD-20251110-352500	completed	238.31	Deliberately surprise the year.	2025-11-10 03:24:51.485157+00	2025-11-10 03:24:51.485157+00
7fce89bb-1a9c-4767-ac93-6d7c1e4e7cec	78784647-3704-4eab-9762-42e8182cd55a	ORD-20251110-171749	shipped	717.93	Instrument the hand for observability.	2025-11-10 03:24:51.485733+00	2025-11-10 03:24:51.485733+00
df80c36b-234e-420f-9d39-e9741f64c120	d4ec9cf1-fc39-47fd-97d3-22e59f6b2f21	ORD-20251110-704237	cancelled	500.79	\N	2025-11-10 03:24:51.486251+00	2025-11-10 03:24:51.486251+00
738c1032-5246-47e1-811c-493c73dbe3d0	7a64aa6f-0c20-4e9c-b563-0cd24e00c5c7	ORD-20251110-255303	completed	764.51	\N	2025-11-10 03:24:51.486837+00	2025-11-10 03:24:51.486837+00
26cf24d1-3db0-4764-bc9b-a5eccaadbb01	ef425dda-42d7-4b1c-b36e-136e96542398	ORD-20251110-517091	completed	569.94	Quietly harden the case last.	2025-11-10 03:24:51.487461+00	2025-11-10 03:24:51.487461+00
a07e2e79-7b94-422d-b967-13783a526a8b	d4ec9cf1-fc39-47fd-97d3-22e59f6b2f21	ORD-20251110-229829	pending	913.52	Invite review for the company in Las Vegas.	2025-11-10 03:24:51.487978+00	2025-11-10 03:24:51.487978+00
5d9cd219-b5c2-4ebe-b550-85da00e5e2e0	1a85d898-3c87-417d-bce4-9a7653ac5dbf	ORD-20251110-801996	cancelled	333.48	Track day over time monthly.	2025-11-10 03:24:51.48857+00	2025-11-10 03:24:51.48857+00
adb86631-5767-424c-9adb-252804c0bc89	78784647-3704-4eab-9762-42e8182cd55a	ORD-20251110-988342	completed	868.01	\N	2025-11-10 03:24:51.489126+00	2025-11-10 03:24:51.489126+00
99c7be4f-e848-42d1-8b11-0a7527e95c90	d4ec9cf1-fc39-47fd-97d3-22e59f6b2f21	ORD-20251110-988917	cancelled	483.89	\N	2025-11-10 03:24:51.48969+00	2025-11-10 03:24:51.48969+00
ef820209-afb8-4849-b62d-ad8fcf2984a8	1a85d898-3c87-417d-bce4-9a7653ac5dbf	ORD-20251110-707714	completed	892.28	\N	2025-11-10 03:24:51.490171+00	2025-11-10 03:24:51.490171+00
f40df804-398f-449b-8547-1e5442cb95ce	78784647-3704-4eab-9762-42e8182cd55a	ORD-20251110-722998	cancelled	954.50	\N	2025-11-10 03:24:51.490647+00	2025-11-10 03:24:51.490647+00
e19e7ae8-b2b2-4427-a5a2-6915944372ae	78784647-3704-4eab-9762-42e8182cd55a	ORD-20251110-816434	cancelled	415.37	According to the hand, align expectations.	2025-11-10 03:24:51.491091+00	2025-11-10 03:24:51.491091+00
fbc747ef-ed0e-4793-be89-51c5910b1efd	68f85432-2588-4de0-9507-5bff96c3ba54	ORD-20251110-300444	shipped	926.46	Create a fallback for case.	2025-11-10 03:24:51.491525+00	2025-11-10 03:24:51.491525+00
afd3e44f-2b24-49a7-92e4-3c8a2eac71b3	1a85d898-3c87-417d-bce4-9a7653ac5dbf	ORD-20251110-877723	shipped	472.62	\N	2025-11-10 03:24:51.491968+00	2025-11-10 03:24:51.491968+00
37e12b98-b1c9-45a4-bb7e-31d87ba1cc74	ef425dda-42d7-4b1c-b36e-136e96542398	ORD-20251110-824970	shipped	705.76	\N	2025-11-10 03:24:51.492413+00	2025-11-10 03:24:51.492413+00
2df22e14-786c-4234-8226-9539b0e748c1	68f85432-2588-4de0-9507-5bff96c3ba54	ORD-20251110-661937	pending	404.81	\N	2025-11-10 03:24:51.492883+00	2025-11-10 03:24:51.492883+00
660e6d10-504e-4761-adf9-0e8aa50168b5	68f85432-2588-4de0-9507-5bff96c3ba54	ORD-20251110-964280	cancelled	133.59	Automate hand recovery gracefully.	2025-11-10 03:24:51.493379+00	2025-11-10 03:24:51.493379+00
675b8299-a94b-4973-a6da-5691d76106ae	1a85d898-3c87-417d-bce4-9a7653ac5dbf	ORD-20251110-657069	pending	509.66	Neither pout when the case spikes.	2025-11-10 03:24:51.493863+00	2025-11-10 03:24:51.493863+00
b5a17b63-86fe-4e30-952f-6089b5d015d9	3815d1a5-5ea7-4fc4-b0ea-5c7a113de884	ORD-20251110-769653	completed	908.89	Across the hand, we insert a smaller part.	2025-11-10 03:24:51.494764+00	2025-11-10 03:24:51.494764+00
\.


--
-- Data for Name: orders_2026; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders_2026 (id, user_id, order_number, status, total_amount, notes, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: orders_default; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders_default (id, user_id, order_number, status, total_amount, notes, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.permissions (id, code, name, description, category, is_active, created_at, updated_at) FROM stdin;
a819da74-fc0c-4957-854a-ebd26d380964	users.create	Create Users	Ability to create new users	users	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
6c349a37-ebf2-4027-aec5-9addba8c68b2	users.read	Read Users	Ability to view user information	users	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
369a23c2-b927-4701-92d8-df5a6cda735b	users.update	Update Users	Ability to update user information	users	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
cad203f1-40c1-4c46-83cc-1e60cf70d872	users.delete	Delete Users	Ability to delete users	users	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
77ad3e8d-e54a-4f51-8f4b-ef95ebb58bb9	orders.create	Create Orders	Ability to create new orders	orders	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
39b03ece-3e5b-4f74-8890-4f27dfc22cbb	orders.read	Read Orders	Ability to view order information	orders	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
c2b06b19-861d-4dd7-ba19-a5c35c8a9454	orders.update	Update Orders	Ability to update order information	orders	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
204a06ad-fc7c-4a32-985e-8cba45792c33	orders.delete	Delete Orders	Ability to delete orders	orders	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
34870755-ce90-440f-b8b9-fc7d7f721aa4	admins.create	Create Admins	Ability to create new admin accounts	admins	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
f821ee73-995f-4f98-8d98-834143907b52	admins.read	Read Admins	Ability to view admin information	admins	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
5a1fb116-1313-4ace-8d16-776f8d89f060	admins.update	Update Admins	Ability to update admin information	admins	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
9b2a7e60-3aeb-446e-9902-8cb6e24a003e	admins.delete	Delete Admins	Ability to delete admin accounts	admins	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
0d67a644-1532-4b0b-a969-26e97f79bd31	admins.manage	Manage Admins	Full admin management access	admins	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
4cc40686-86b1-4d38-9d99-a3f1a0cf20c3	settings.general	General Settings	Ability to manage general settings	settings	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
2f7e9085-c8ce-4de0-8b71-47a2d130c2d1	settings.security	Security Settings	Ability to manage security settings	settings	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
2765ad95-382a-47c8-90c8-ad680aabc645	analytics.dashboard	Analytics Dashboard	Ability to view analytics dashboard	analytics	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
01f402b6-ee4f-4060-bb82-62280927e8ca	analytics.reports	Analytics Reports	Ability to view and generate reports	analytics	t	2025-11-10 04:19:16.195758+00	2025-11-10 04:19:16.195758+00
\.


--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.role_permissions (id, role, permission_id, created_at) FROM stdin;
8fbea99f-a684-4db2-b119-16879f4a40a5	super_admin	a819da74-fc0c-4957-854a-ebd26d380964	2025-11-10 04:19:16.195758+00
88ec275b-04f5-4395-8061-0ca98cfb8c1d	super_admin	6c349a37-ebf2-4027-aec5-9addba8c68b2	2025-11-10 04:19:16.195758+00
1acf2e1b-8cad-4657-976d-0923a2166527	super_admin	369a23c2-b927-4701-92d8-df5a6cda735b	2025-11-10 04:19:16.195758+00
d5ae8bac-cd29-4bb9-9825-88760f339a5c	super_admin	cad203f1-40c1-4c46-83cc-1e60cf70d872	2025-11-10 04:19:16.195758+00
79fd7361-8e3c-4e5d-8465-df5867d25b8c	super_admin	77ad3e8d-e54a-4f51-8f4b-ef95ebb58bb9	2025-11-10 04:19:16.195758+00
ee2c45af-9754-4dab-84fd-39ae6b94c829	super_admin	39b03ece-3e5b-4f74-8890-4f27dfc22cbb	2025-11-10 04:19:16.195758+00
3273bd33-bf18-46d9-b073-67c2b7968121	super_admin	c2b06b19-861d-4dd7-ba19-a5c35c8a9454	2025-11-10 04:19:16.195758+00
7ad9b076-ea73-4131-b766-8b4101b1a31a	super_admin	204a06ad-fc7c-4a32-985e-8cba45792c33	2025-11-10 04:19:16.195758+00
b8468c80-7992-4039-9aa9-cad23ea4fe24	super_admin	34870755-ce90-440f-b8b9-fc7d7f721aa4	2025-11-10 04:19:16.195758+00
33b9c520-2e7b-4e32-ac27-7888ec48e74a	super_admin	f821ee73-995f-4f98-8d98-834143907b52	2025-11-10 04:19:16.195758+00
89ec9bae-753b-421b-b85a-59ea0f81af6d	super_admin	5a1fb116-1313-4ace-8d16-776f8d89f060	2025-11-10 04:19:16.195758+00
155fdf90-35e5-47e3-9f8d-dabcb52f0e47	super_admin	9b2a7e60-3aeb-446e-9902-8cb6e24a003e	2025-11-10 04:19:16.195758+00
81d5dc8c-70e5-48c1-91a6-fa94cb6b525b	super_admin	0d67a644-1532-4b0b-a969-26e97f79bd31	2025-11-10 04:19:16.195758+00
ea1546e9-4243-4940-850f-212eca3d5b81	super_admin	4cc40686-86b1-4d38-9d99-a3f1a0cf20c3	2025-11-10 04:19:16.195758+00
14f3417e-551e-4059-a8f5-9204582a7005	super_admin	2f7e9085-c8ce-4de0-8b71-47a2d130c2d1	2025-11-10 04:19:16.195758+00
6bbccbb7-c671-4e35-8843-42321fe06601	super_admin	2765ad95-382a-47c8-90c8-ad680aabc645	2025-11-10 04:19:16.195758+00
207cd83c-1d77-4303-ab55-042a84ecc188	super_admin	01f402b6-ee4f-4060-bb82-62280927e8ca	2025-11-10 04:19:16.195758+00
fb40faa7-4956-49ee-8845-0f6b45dd9eca	admin	a819da74-fc0c-4957-854a-ebd26d380964	2025-11-10 04:19:16.195758+00
12123c79-30d4-4626-815c-c62f949da85f	admin	6c349a37-ebf2-4027-aec5-9addba8c68b2	2025-11-10 04:19:16.195758+00
7f159b4e-a649-4ea2-b514-aa3c40391e36	admin	369a23c2-b927-4701-92d8-df5a6cda735b	2025-11-10 04:19:16.195758+00
949397b5-3aee-45ae-856b-f3ff826df348	admin	77ad3e8d-e54a-4f51-8f4b-ef95ebb58bb9	2025-11-10 04:19:16.195758+00
c0e85063-2046-4042-989a-8df1f3ffb5d9	admin	39b03ece-3e5b-4f74-8890-4f27dfc22cbb	2025-11-10 04:19:16.195758+00
de614610-f44a-413a-9493-8aa0e3247d60	admin	c2b06b19-861d-4dd7-ba19-a5c35c8a9454	2025-11-10 04:19:16.195758+00
90d93c37-8e63-4182-92ce-d6bd1793c81c	admin	4cc40686-86b1-4d38-9d99-a3f1a0cf20c3	2025-11-10 04:19:16.195758+00
ca3ceb0d-f8dd-4164-bb07-c6beed35bb6b	admin	2765ad95-382a-47c8-90c8-ad680aabc645	2025-11-10 04:19:16.195758+00
beafa001-b9ef-46f8-9b1d-336937ac23fa	moderator	6c349a37-ebf2-4027-aec5-9addba8c68b2	2025-11-10 04:19:16.195758+00
7df646c0-c43d-462a-ab4f-8b1445e53520	moderator	39b03ece-3e5b-4f74-8890-4f27dfc22cbb	2025-11-10 04:19:16.195758+00
9c01c353-7353-43bb-a312-977ab02332b1	moderator	2765ad95-382a-47c8-90c8-ad680aabc645	2025-11-10 04:19:16.195758+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
10	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, username, password_hash, first_name, last_name, created_at, updated_at) FROM stdin;
42a08b8c-bd12-4f93-9b5c-fc0329ece04f	loyalabshire@oberbrunner.info	Oliveclam0	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Cara	Dach	2025-11-05 08:34:06.243711+00	2025-11-05 08:34:06.243711+00
341a3946-cede-4238-9ff3-e8a639053ee5	andreschiller@romaguera.info	RobertRosyBrown1	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Beverly	Rogahn	2025-11-05 08:34:06.247778+00	2025-11-05 08:34:06.247778+00
d38318f9-54e7-46f4-8c15-cd27302ad921	maeganfunk@heller.net	Mossiesome2	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Sigmund	Eichmann	2025-11-05 08:34:06.248553+00	2025-11-05 08:34:06.248553+00
d40f6d7e-87dc-4c1b-8e56-0595ccafbf0f	justinaschinner@barton.name	anythingAbel3	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Mara	Botsford	2025-11-05 08:34:06.249048+00	2025-11-05 08:34:06.249048+00
1ed1a406-aff0-4db8-a0a2-2f58b9d973ba	adolfkohler@kihn.org	trip4794	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Josianne	Spencer	2025-11-05 08:34:06.24956+00	2025-11-05 08:34:06.24956+00
3ef1cc3b-839d-4e53-bb29-ad85b983d3ed	tobymcclure@carter.io	theFisher5	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Jarrell	Bosco	2025-11-05 08:34:06.250068+00	2025-11-05 08:34:06.250068+00
98292ef6-0cdf-4ad1-8cbe-fbe4d46d1645	sylvesterlang@beatty.net	NeilMayert6	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Ollie	Haag	2025-11-05 08:34:06.2505+00	2025-11-05 08:34:06.2505+00
0e27bbf4-00c3-4279-bcfd-6fa3959f1b95	marquesokuneva@baumbach.com	NuttySpider067	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Ora	Olson	2025-11-05 08:34:06.25097+00	2025-11-05 08:34:06.25097+00
1ff11ab9-0e81-41aa-a76a-9421153718f9	dameonzemlak@walker.io	Terrell79538	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Emmanuelle	Trantow	2025-11-05 08:34:06.251377+00	2025-11-05 08:34:06.251377+00
0ee653d3-552b-424e-bf71-d14b206306b5	kamrynmohr@mante.com	Representative8309	$2a$10$9KUDWiM5uiGtgXurjPt4EOf91Cup2dlrAjHcl07Ko7FqI1O.lOtlG	Americo	Auer	2025-11-05 08:34:06.251756+00	2025-11-05 08:34:06.251756+00
a9dcf66a-cacd-4a2d-bd5d-7cb3105611e0	yolandahickle@conroy.info	Representative8200	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Adaline	Raynor	2025-11-05 08:34:19.359976+00	2025-11-05 08:34:19.359976+00
dd6a7abb-5a19-42c4-b3fb-0b0f94010e1e	rosalynbins@torp.org	SpottedWolf1	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Raymundo	Friesen	2025-11-05 08:34:19.362548+00	2025-11-05 08:34:19.362548+00
bfbd9bac-abcd-4cec-89e1-df20a9a99238	neilturner@douglas.io	mrBuckridge2	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Dorcas	Quitzon	2025-11-05 08:34:19.363067+00	2025-11-05 08:34:19.363067+00
4825d06b-421a-4a86-ab8d-303118a024dc	sophierosenbaum@krajcik.org	Julie.413	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Mark	Quigley	2025-11-05 08:34:19.363561+00	2025-11-05 08:34:19.363561+00
4a5bbae2-90e9-4e15-8607-7aee1ca1fad7	roscoefahey@kunde.biz	ForestGreen_bevy4	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Miguel	Doyle	2025-11-05 08:34:19.364065+00	2025-11-05 08:34:19.364065+00
5538dfac-492d-45c7-98f5-8f901759f313	armanddaniel@bernier.org	DarkSalmoncat5	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Lane	Leuschke	2025-11-05 08:34:19.364564+00	2025-11-05 08:34:19.364564+00
7fe9c1b5-b5d2-4674-b4f8-d20b5bf07ce4	koreycartwright@schamberger.biz	Magenta_troop6	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Brennan	Hand	2025-11-05 08:34:19.365049+00	2025-11-05 08:34:19.365049+00
eaa88a4d-a54f-42ee-8643-82362a1ab2b0	owengrant@franecki.org	reindeer.837	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Jacques	Little	2025-11-05 08:34:19.365471+00	2025-11-05 08:34:19.365471+00
9e280b82-1658-46b2-9ece-2064f12e9eb7	michaelakunze@barrows.info	CoraSaddleBrown8	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Victor	Grady	2025-11-05 08:34:19.365936+00	2025-11-05 08:34:19.365936+00
adc3f01c-9b7d-4fd6-a103-7bf890fc2b22	samarapadberg@funk.com	LoganPurple9	$2a$10$W222OHHlCfqwWztGZBe2QuiOgHyWTQqKdOM0K7i5P1a4NxzLU/XMy	Juwan	Mayer	2025-11-05 08:34:19.366487+00	2025-11-05 08:34:19.366487+00
af5a86ac-b114-4c6b-a940-e05d0703362f	email@email.com	username	$2a$10$I7YTMrZLXRY5ZCB3nln.ee.mmvfcVLu.2qAa5jomM5yOBiWwVsyj.	\N	\N	2025-11-06 07:30:26.694461+00	2025-11-06 07:30:26.694461+00
bb8c254f-b0e7-4984-aab0-26a1f160d2a2	test@test.com	username1	$2a$10$TRt6NzOF.RAO2u0aCtglm.7HvkX6kqN4MoJEgp9XH6bzzvPaplA/u	\N	\N	2025-11-06 07:34:33.716204+00	2025-11-06 07:34:33.716204+00
16da6b0b-ae6c-4bf7-86d6-fe9d788d9497	test@test2.com	asd123	$2a$10$S1cUirMO2ioln7NwXZ7FwOJkgHl7w1HYiPkoc5aSbkfO8FMfor2em	\N	\N	2025-11-06 07:48:58.536896+00	2025-11-06 07:48:58.536896+00
084f6366-d88f-4879-885a-3540859ce456	test@test22.com	asd1232	$2a$10$Yi8Knfqsw71ird53eoi0K.B4jkbuxtFGp9R3bbnxVPRCHbIoWRmxW	\N	\N	2025-11-06 07:49:57.363116+00	2025-11-06 07:49:57.363116+00
68e36490-c7f8-4918-913c-b0e4d9e83226	test@test222.com	asd12321	$2a$10$y.sEn6z4QCJqz2bA1QiPiuoenroEXW/uO6fgdsoglGHps3J9hy/mK	\N	\N	2025-11-06 07:52:58.560361+00	2025-11-06 07:52:58.560361+00
c52ab351-85e1-4f7d-bb73-773b2fe222da	asd@dsa.com	mmm	$2a$10$aTL4eBZFieB7yCpIbNUEXe2d8iZNXpUniT31J4rxyB7JgeCoA4/9m	\N	\N	2025-11-06 08:49:01.784572+00	2025-11-06 08:49:01.784572+00
110ecfd3-614f-4e18-ba3e-3e3536a9256a	asd@dsa1.com	mmm1	$2a$10$C.ZAnUd2/71Xu.e.GZopW./Qp/o3ABBmjJhj6Es50zgWZC6S.Uyqu	\N	\N	2025-11-07 07:54:37.763688+00	2025-11-07 07:54:37.763688+00
ee5787b1-b0ac-4119-9cd3-16ee1040992b	asd@dsa12.com	mmm2	$2a$10$p.zCNGrtGING365GJbZ/SeoepHvSRwFG6dohM7QPY0uLi4fXWZHtO	\N	\N	2025-11-07 08:02:21.761833+00	2025-11-07 08:02:21.761833+00
b27bf955-43e1-4e8d-a3ec-f0cae6c17d37	asd@dsa10.com	mmm3	$2a$10$m47gGN72dSmWBqv.G46D4ufKqpKy3dS.pp8HRYHKVsgqmRyGk2zJe	\N	\N	2025-11-07 08:05:14.030644+00	2025-11-07 08:05:14.030644+00
3771ce44-478c-45b5-aca4-06e018883f9a	asd@dsa20.com	mmm4	$2a$10$zwjErQZW/g/Vk/2t3KOkx.JMsNYt9oESoro2ggMAU7rpZcuylq9Re	\N	\N	2025-11-07 08:05:52.517015+00	2025-11-07 08:05:52.517015+00
78784647-3704-4eab-9762-42e8182cd55a	shaniabashirian@heller.com	Beans1220	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Norval	Thompson	2025-11-10 03:24:51.474154+00	2025-11-10 03:24:51.474154+00
d4ec9cf1-fc39-47fd-97d3-22e59f6b2f21	allancarroll@ondricka.biz	CarefulRefrigerator1	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Rowan	Bradtke	2025-11-10 03:24:51.475612+00	2025-11-10 03:24:51.475612+00
ef425dda-42d7-4b1c-b36e-136e96542398	dennismaggio@jerde.com	friendlydinosaur2	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Bernadette	Collins	2025-11-10 03:24:51.476511+00	2025-11-10 03:24:51.476511+00
949d5e16-4383-4975-b5fd-5b0ed6a2a1a9	wilsonfay@quigley.net	RosemarieBurlyWood3	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Kiley	Legros	2025-11-10 03:24:51.477428+00	2025-11-10 03:24:51.477428+00
7a64aa6f-0c20-4e9c-b563-0cd24e00c5c7	maudetoy@ullrich.net	Ari0034	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Kelley	Murazik	2025-11-10 03:24:51.478244+00	2025-11-10 03:24:51.478244+00
7a6fd6cc-c579-4704-a9d1-7eaaa5ac5f28	catherinestiedemann@tromp.name	Kulas_IwJ5	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Alba	Watsica	2025-11-10 03:24:51.479075+00	2025-11-10 03:24:51.479075+00
68f85432-2588-4de0-9507-5bff96c3ba54	ericherdman@pfeffer.name	LightCoral_riches6	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Jazmin	Hyatt	2025-11-10 03:24:51.479744+00	2025-11-10 03:24:51.479744+00
55b0fd16-bcb6-4e0f-817b-5bda4737bce0	genebeer@von.io	Damian7367	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Selmer	Watsica	2025-11-10 03:24:51.480421+00	2025-11-10 03:24:51.480421+00
3815d1a5-5ea7-4fc4-b0ea-5c7a113de884	nyasiaabbott@armstrong.biz	Planner728	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Skylar	Nitzsche	2025-11-10 03:24:51.481088+00	2025-11-10 03:24:51.481088+00
1a85d898-3c87-417d-bce4-9a7653ac5dbf	catharineswift@stroman.org	OReilly69569	$2a$10$ZTB4vuG8c98zNV0Qm98CvO14vH/Szt9tDkEYKA4D6VEGRXlAQkyUC	Camille	Carroll	2025-11-10 03:24:51.482242+00	2025-11-10 03:24:51.482242+00
\.


--
-- Name: admins admins_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_email_key UNIQUE (email);


--
-- Name: admins admins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_pkey PRIMARY KEY (id);


--
-- Name: admins admins_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_username_key UNIQUE (username);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id, created_at);


--
-- Name: audit_logs_2025_11 audit_logs_2025_11_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs_2025_11
    ADD CONSTRAINT audit_logs_2025_11_pkey PRIMARY KEY (id, created_at);


--
-- Name: audit_logs_2025_12 audit_logs_2025_12_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs_2025_12
    ADD CONSTRAINT audit_logs_2025_12_pkey PRIMARY KEY (id, created_at);


--
-- Name: audit_logs_2026_01 audit_logs_2026_01_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs_2026_01
    ADD CONSTRAINT audit_logs_2026_01_pkey PRIMARY KEY (id, created_at);


--
-- Name: audit_logs_2026_02 audit_logs_2026_02_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs_2026_02
    ADD CONSTRAINT audit_logs_2026_02_pkey PRIMARY KEY (id, created_at);


--
-- Name: audit_logs_default audit_logs_default_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs_default
    ADD CONSTRAINT audit_logs_default_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs error_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs
    ADD CONSTRAINT error_logs_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs_2025_11 error_logs_2025_11_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs_2025_11
    ADD CONSTRAINT error_logs_2025_11_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs_2025_12 error_logs_2025_12_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs_2025_12
    ADD CONSTRAINT error_logs_2025_12_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs_2026_01 error_logs_2026_01_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs_2026_01
    ADD CONSTRAINT error_logs_2026_01_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs_2026_02 error_logs_2026_02_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs_2026_02
    ADD CONSTRAINT error_logs_2026_02_pkey PRIMARY KEY (id, created_at);


--
-- Name: error_logs_default error_logs_default_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.error_logs_default
    ADD CONSTRAINT error_logs_default_pkey PRIMARY KEY (id, created_at);


--
-- Name: menu_items menu_items_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_code_key UNIQUE (code);


--
-- Name: menu_items menu_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id, created_at);


--
-- Name: orders_2025 orders_2025_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders_2025
    ADD CONSTRAINT orders_2025_pkey PRIMARY KEY (id, created_at);


--
-- Name: orders_2026 orders_2026_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders_2026
    ADD CONSTRAINT orders_2026_pkey PRIMARY KEY (id, created_at);


--
-- Name: orders_default orders_default_pkey1; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders_default
    ADD CONSTRAINT orders_default_pkey1 PRIMARY KEY (id, created_at);


--
-- Name: permissions permissions_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_code_key UNIQUE (code);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_role_permission_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_permission_id_key UNIQUE (role, permission_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_created_at ON ONLY public.audit_logs USING btree (created_at DESC);


--
-- Name: audit_logs_2025_11_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_11_created_at_idx ON public.audit_logs_2025_11 USING btree (created_at DESC);


--
-- Name: idx_audit_logs_composite; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_composite ON ONLY public.audit_logs USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: audit_logs_2025_11_entity_type_entity_id_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_11_entity_type_entity_id_created_at_idx ON public.audit_logs_2025_11 USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: idx_audit_logs_entity; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_entity ON ONLY public.audit_logs USING btree (entity_type, entity_id);


--
-- Name: audit_logs_2025_11_entity_type_entity_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_11_entity_type_entity_id_idx ON public.audit_logs_2025_11 USING btree (entity_type, entity_id);


--
-- Name: idx_audit_logs_request_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_request_id ON ONLY public.audit_logs USING btree (request_id);


--
-- Name: audit_logs_2025_11_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_11_request_id_idx ON public.audit_logs_2025_11 USING btree (request_id);


--
-- Name: idx_audit_logs_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_user_id ON ONLY public.audit_logs USING btree (user_id);


--
-- Name: audit_logs_2025_11_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_11_user_id_idx ON public.audit_logs_2025_11 USING btree (user_id);


--
-- Name: audit_logs_2025_12_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_12_created_at_idx ON public.audit_logs_2025_12 USING btree (created_at DESC);


--
-- Name: audit_logs_2025_12_entity_type_entity_id_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_12_entity_type_entity_id_created_at_idx ON public.audit_logs_2025_12 USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: audit_logs_2025_12_entity_type_entity_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_12_entity_type_entity_id_idx ON public.audit_logs_2025_12 USING btree (entity_type, entity_id);


--
-- Name: audit_logs_2025_12_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_12_request_id_idx ON public.audit_logs_2025_12 USING btree (request_id);


--
-- Name: audit_logs_2025_12_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2025_12_user_id_idx ON public.audit_logs_2025_12 USING btree (user_id);


--
-- Name: audit_logs_2026_01_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_01_created_at_idx ON public.audit_logs_2026_01 USING btree (created_at DESC);


--
-- Name: audit_logs_2026_01_entity_type_entity_id_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_01_entity_type_entity_id_created_at_idx ON public.audit_logs_2026_01 USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: audit_logs_2026_01_entity_type_entity_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_01_entity_type_entity_id_idx ON public.audit_logs_2026_01 USING btree (entity_type, entity_id);


--
-- Name: audit_logs_2026_01_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_01_request_id_idx ON public.audit_logs_2026_01 USING btree (request_id);


--
-- Name: audit_logs_2026_01_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_01_user_id_idx ON public.audit_logs_2026_01 USING btree (user_id);


--
-- Name: audit_logs_2026_02_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_02_created_at_idx ON public.audit_logs_2026_02 USING btree (created_at DESC);


--
-- Name: audit_logs_2026_02_entity_type_entity_id_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_02_entity_type_entity_id_created_at_idx ON public.audit_logs_2026_02 USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: audit_logs_2026_02_entity_type_entity_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_02_entity_type_entity_id_idx ON public.audit_logs_2026_02 USING btree (entity_type, entity_id);


--
-- Name: audit_logs_2026_02_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_02_request_id_idx ON public.audit_logs_2026_02 USING btree (request_id);


--
-- Name: audit_logs_2026_02_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_2026_02_user_id_idx ON public.audit_logs_2026_02 USING btree (user_id);


--
-- Name: audit_logs_default_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_default_created_at_idx ON public.audit_logs_default USING btree (created_at DESC);


--
-- Name: audit_logs_default_entity_type_entity_id_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_default_entity_type_entity_id_created_at_idx ON public.audit_logs_default USING btree (entity_type, entity_id, created_at DESC);


--
-- Name: audit_logs_default_entity_type_entity_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_default_entity_type_entity_id_idx ON public.audit_logs_default USING btree (entity_type, entity_id);


--
-- Name: audit_logs_default_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_default_request_id_idx ON public.audit_logs_default USING btree (request_id);


--
-- Name: audit_logs_default_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_logs_default_user_id_idx ON public.audit_logs_default USING btree (user_id);


--
-- Name: idx_error_logs_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_error_logs_created_at ON ONLY public.error_logs USING btree (created_at DESC);


--
-- Name: error_logs_2025_11_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_11_created_at_idx ON public.error_logs_2025_11 USING btree (created_at DESC);


--
-- Name: idx_error_logs_error_type; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_error_logs_error_type ON ONLY public.error_logs USING btree (error_type);


--
-- Name: error_logs_2025_11_error_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_11_error_type_idx ON public.error_logs_2025_11 USING btree (error_type);


--
-- Name: idx_error_logs_request_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_error_logs_request_id ON ONLY public.error_logs USING btree (request_id);


--
-- Name: error_logs_2025_11_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_11_request_id_idx ON public.error_logs_2025_11 USING btree (request_id);


--
-- Name: idx_error_logs_request_path; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_error_logs_request_path ON ONLY public.error_logs USING btree (request_path);


--
-- Name: error_logs_2025_11_request_path_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_11_request_path_idx ON public.error_logs_2025_11 USING btree (request_path);


--
-- Name: idx_error_logs_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_error_logs_user_id ON ONLY public.error_logs USING btree (user_id);


--
-- Name: error_logs_2025_11_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_11_user_id_idx ON public.error_logs_2025_11 USING btree (user_id);


--
-- Name: error_logs_2025_12_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_12_created_at_idx ON public.error_logs_2025_12 USING btree (created_at DESC);


--
-- Name: error_logs_2025_12_error_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_12_error_type_idx ON public.error_logs_2025_12 USING btree (error_type);


--
-- Name: error_logs_2025_12_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_12_request_id_idx ON public.error_logs_2025_12 USING btree (request_id);


--
-- Name: error_logs_2025_12_request_path_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_12_request_path_idx ON public.error_logs_2025_12 USING btree (request_path);


--
-- Name: error_logs_2025_12_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2025_12_user_id_idx ON public.error_logs_2025_12 USING btree (user_id);


--
-- Name: error_logs_2026_01_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_01_created_at_idx ON public.error_logs_2026_01 USING btree (created_at DESC);


--
-- Name: error_logs_2026_01_error_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_01_error_type_idx ON public.error_logs_2026_01 USING btree (error_type);


--
-- Name: error_logs_2026_01_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_01_request_id_idx ON public.error_logs_2026_01 USING btree (request_id);


--
-- Name: error_logs_2026_01_request_path_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_01_request_path_idx ON public.error_logs_2026_01 USING btree (request_path);


--
-- Name: error_logs_2026_01_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_01_user_id_idx ON public.error_logs_2026_01 USING btree (user_id);


--
-- Name: error_logs_2026_02_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_02_created_at_idx ON public.error_logs_2026_02 USING btree (created_at DESC);


--
-- Name: error_logs_2026_02_error_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_02_error_type_idx ON public.error_logs_2026_02 USING btree (error_type);


--
-- Name: error_logs_2026_02_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_02_request_id_idx ON public.error_logs_2026_02 USING btree (request_id);


--
-- Name: error_logs_2026_02_request_path_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_02_request_path_idx ON public.error_logs_2026_02 USING btree (request_path);


--
-- Name: error_logs_2026_02_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_2026_02_user_id_idx ON public.error_logs_2026_02 USING btree (user_id);


--
-- Name: error_logs_default_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_default_created_at_idx ON public.error_logs_default USING btree (created_at DESC);


--
-- Name: error_logs_default_error_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_default_error_type_idx ON public.error_logs_default USING btree (error_type);


--
-- Name: error_logs_default_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_default_request_id_idx ON public.error_logs_default USING btree (request_id);


--
-- Name: error_logs_default_request_path_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_default_request_path_idx ON public.error_logs_default USING btree (request_path);


--
-- Name: error_logs_default_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX error_logs_default_user_id_idx ON public.error_logs_default USING btree (user_id);


--
-- Name: idx_admins_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_admins_email ON public.admins USING btree (email);


--
-- Name: idx_admins_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_admins_role ON public.admins USING btree (role);


--
-- Name: idx_admins_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_admins_username ON public.admins USING btree (username);


--
-- Name: idx_menu_items_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_menu_items_is_active ON public.menu_items USING btree (is_active);


--
-- Name: idx_menu_items_parent_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_menu_items_parent_id ON public.menu_items USING btree (parent_id);


--
-- Name: idx_menu_items_permission_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_menu_items_permission_id ON public.menu_items USING btree (permission_id);


--
-- Name: idx_permissions_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_permissions_code ON public.permissions USING btree (code);


--
-- Name: idx_permissions_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_permissions_is_active ON public.permissions USING btree (is_active);


--
-- Name: idx_role_permissions_permission_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_role_permissions_permission_id ON public.role_permissions USING btree (permission_id);


--
-- Name: idx_role_permissions_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_role_permissions_role ON public.role_permissions USING btree (role);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: audit_logs_2025_11_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_created_at ATTACH PARTITION public.audit_logs_2025_11_created_at_idx;


--
-- Name: audit_logs_2025_11_entity_type_entity_id_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_composite ATTACH PARTITION public.audit_logs_2025_11_entity_type_entity_id_created_at_idx;


--
-- Name: audit_logs_2025_11_entity_type_entity_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_entity ATTACH PARTITION public.audit_logs_2025_11_entity_type_entity_id_idx;


--
-- Name: audit_logs_2025_11_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.audit_logs_pkey ATTACH PARTITION public.audit_logs_2025_11_pkey;


--
-- Name: audit_logs_2025_11_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_request_id ATTACH PARTITION public.audit_logs_2025_11_request_id_idx;


--
-- Name: audit_logs_2025_11_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_user_id ATTACH PARTITION public.audit_logs_2025_11_user_id_idx;


--
-- Name: audit_logs_2025_12_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_created_at ATTACH PARTITION public.audit_logs_2025_12_created_at_idx;


--
-- Name: audit_logs_2025_12_entity_type_entity_id_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_composite ATTACH PARTITION public.audit_logs_2025_12_entity_type_entity_id_created_at_idx;


--
-- Name: audit_logs_2025_12_entity_type_entity_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_entity ATTACH PARTITION public.audit_logs_2025_12_entity_type_entity_id_idx;


--
-- Name: audit_logs_2025_12_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.audit_logs_pkey ATTACH PARTITION public.audit_logs_2025_12_pkey;


--
-- Name: audit_logs_2025_12_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_request_id ATTACH PARTITION public.audit_logs_2025_12_request_id_idx;


--
-- Name: audit_logs_2025_12_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_user_id ATTACH PARTITION public.audit_logs_2025_12_user_id_idx;


--
-- Name: audit_logs_2026_01_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_created_at ATTACH PARTITION public.audit_logs_2026_01_created_at_idx;


--
-- Name: audit_logs_2026_01_entity_type_entity_id_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_composite ATTACH PARTITION public.audit_logs_2026_01_entity_type_entity_id_created_at_idx;


--
-- Name: audit_logs_2026_01_entity_type_entity_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_entity ATTACH PARTITION public.audit_logs_2026_01_entity_type_entity_id_idx;


--
-- Name: audit_logs_2026_01_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.audit_logs_pkey ATTACH PARTITION public.audit_logs_2026_01_pkey;


--
-- Name: audit_logs_2026_01_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_request_id ATTACH PARTITION public.audit_logs_2026_01_request_id_idx;


--
-- Name: audit_logs_2026_01_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_user_id ATTACH PARTITION public.audit_logs_2026_01_user_id_idx;


--
-- Name: audit_logs_2026_02_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_created_at ATTACH PARTITION public.audit_logs_2026_02_created_at_idx;


--
-- Name: audit_logs_2026_02_entity_type_entity_id_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_composite ATTACH PARTITION public.audit_logs_2026_02_entity_type_entity_id_created_at_idx;


--
-- Name: audit_logs_2026_02_entity_type_entity_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_entity ATTACH PARTITION public.audit_logs_2026_02_entity_type_entity_id_idx;


--
-- Name: audit_logs_2026_02_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.audit_logs_pkey ATTACH PARTITION public.audit_logs_2026_02_pkey;


--
-- Name: audit_logs_2026_02_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_request_id ATTACH PARTITION public.audit_logs_2026_02_request_id_idx;


--
-- Name: audit_logs_2026_02_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_user_id ATTACH PARTITION public.audit_logs_2026_02_user_id_idx;


--
-- Name: audit_logs_default_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_created_at ATTACH PARTITION public.audit_logs_default_created_at_idx;


--
-- Name: audit_logs_default_entity_type_entity_id_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_composite ATTACH PARTITION public.audit_logs_default_entity_type_entity_id_created_at_idx;


--
-- Name: audit_logs_default_entity_type_entity_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_entity ATTACH PARTITION public.audit_logs_default_entity_type_entity_id_idx;


--
-- Name: audit_logs_default_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.audit_logs_pkey ATTACH PARTITION public.audit_logs_default_pkey;


--
-- Name: audit_logs_default_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_request_id ATTACH PARTITION public.audit_logs_default_request_id_idx;


--
-- Name: audit_logs_default_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_audit_logs_user_id ATTACH PARTITION public.audit_logs_default_user_id_idx;


--
-- Name: error_logs_2025_11_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_created_at ATTACH PARTITION public.error_logs_2025_11_created_at_idx;


--
-- Name: error_logs_2025_11_error_type_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_error_type ATTACH PARTITION public.error_logs_2025_11_error_type_idx;


--
-- Name: error_logs_2025_11_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.error_logs_pkey ATTACH PARTITION public.error_logs_2025_11_pkey;


--
-- Name: error_logs_2025_11_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_id ATTACH PARTITION public.error_logs_2025_11_request_id_idx;


--
-- Name: error_logs_2025_11_request_path_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_path ATTACH PARTITION public.error_logs_2025_11_request_path_idx;


--
-- Name: error_logs_2025_11_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_user_id ATTACH PARTITION public.error_logs_2025_11_user_id_idx;


--
-- Name: error_logs_2025_12_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_created_at ATTACH PARTITION public.error_logs_2025_12_created_at_idx;


--
-- Name: error_logs_2025_12_error_type_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_error_type ATTACH PARTITION public.error_logs_2025_12_error_type_idx;


--
-- Name: error_logs_2025_12_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.error_logs_pkey ATTACH PARTITION public.error_logs_2025_12_pkey;


--
-- Name: error_logs_2025_12_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_id ATTACH PARTITION public.error_logs_2025_12_request_id_idx;


--
-- Name: error_logs_2025_12_request_path_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_path ATTACH PARTITION public.error_logs_2025_12_request_path_idx;


--
-- Name: error_logs_2025_12_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_user_id ATTACH PARTITION public.error_logs_2025_12_user_id_idx;


--
-- Name: error_logs_2026_01_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_created_at ATTACH PARTITION public.error_logs_2026_01_created_at_idx;


--
-- Name: error_logs_2026_01_error_type_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_error_type ATTACH PARTITION public.error_logs_2026_01_error_type_idx;


--
-- Name: error_logs_2026_01_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.error_logs_pkey ATTACH PARTITION public.error_logs_2026_01_pkey;


--
-- Name: error_logs_2026_01_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_id ATTACH PARTITION public.error_logs_2026_01_request_id_idx;


--
-- Name: error_logs_2026_01_request_path_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_path ATTACH PARTITION public.error_logs_2026_01_request_path_idx;


--
-- Name: error_logs_2026_01_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_user_id ATTACH PARTITION public.error_logs_2026_01_user_id_idx;


--
-- Name: error_logs_2026_02_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_created_at ATTACH PARTITION public.error_logs_2026_02_created_at_idx;


--
-- Name: error_logs_2026_02_error_type_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_error_type ATTACH PARTITION public.error_logs_2026_02_error_type_idx;


--
-- Name: error_logs_2026_02_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.error_logs_pkey ATTACH PARTITION public.error_logs_2026_02_pkey;


--
-- Name: error_logs_2026_02_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_id ATTACH PARTITION public.error_logs_2026_02_request_id_idx;


--
-- Name: error_logs_2026_02_request_path_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_path ATTACH PARTITION public.error_logs_2026_02_request_path_idx;


--
-- Name: error_logs_2026_02_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_user_id ATTACH PARTITION public.error_logs_2026_02_user_id_idx;


--
-- Name: error_logs_default_created_at_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_created_at ATTACH PARTITION public.error_logs_default_created_at_idx;


--
-- Name: error_logs_default_error_type_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_error_type ATTACH PARTITION public.error_logs_default_error_type_idx;


--
-- Name: error_logs_default_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.error_logs_pkey ATTACH PARTITION public.error_logs_default_pkey;


--
-- Name: error_logs_default_request_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_id ATTACH PARTITION public.error_logs_default_request_id_idx;


--
-- Name: error_logs_default_request_path_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_request_path ATTACH PARTITION public.error_logs_default_request_path_idx;


--
-- Name: error_logs_default_user_id_idx; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.idx_error_logs_user_id ATTACH PARTITION public.error_logs_default_user_id_idx;


--
-- Name: orders_2025_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.orders_pkey ATTACH PARTITION public.orders_2025_pkey;


--
-- Name: orders_2026_pkey; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.orders_pkey ATTACH PARTITION public.orders_2026_pkey;


--
-- Name: orders_default_pkey1; Type: INDEX ATTACH; Schema: public; Owner: postgres
--

ALTER INDEX public.orders_pkey ATTACH PARTITION public.orders_default_pkey1;


--
-- Name: admins trigger_update_admins_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_admins_updated_at BEFORE UPDATE ON public.admins FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: menu_items trigger_update_menu_items_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_menu_items_updated_at BEFORE UPDATE ON public.menu_items FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: orders trigger_update_orders_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_orders_updated_at BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: permissions trigger_update_permissions_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_permissions_updated_at BEFORE UPDATE ON public.permissions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users trigger_update_users_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: error_logs error_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE public.error_logs
    ADD CONSTRAINT error_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: menu_items menu_items_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.menu_items(id) ON DELETE CASCADE;


--
-- Name: menu_items menu_items_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE SET NULL;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict rhr3LFJmf1heVyUtbL6deJ8y0WLHslUbwaleBYdJNBCIwl220c2UJ9mxegsaxcZ

