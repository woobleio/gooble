--
-- PostgreSQL database dump
--

-- CREATE TABLE preview_position (
--   position_id text NOT NULL PRIMARY KEY,
--   style_source text NOT NULL
-- )

-- ALTER TABLE creation ADD COLUMN preview_position_id text
-- ALTER TABLE ONLY creation ADD CONSTRAINT preview_position_fk FOREIGN KEY (preview_position_id) REFERENCES preview_position(position_id);

-- INSERT INTO preview_position (position_id, style_source) VALUES
--   ('TL', '')
--   ('TM', 'body { display: inline-block; width: auto; position: absolute; left: 50%; transform: translateX(-50%);}')
--   ('TR', 'body { float: right; }')
--   ('ML', 'body { display: inline-block; width: auto; position: absolute; top: 50%; transform: translateY(-50%);}')
--   ('M', 'body { display: inline-block; width: auto; position: absolute; left: 50%; top: 50%; transform: translate(-50%,-50%);}')
--   ('MR', 'body { display: inline-block; width: auto; position: absolute; right: 0; top: 50%; transform: translateY(-50%);}')
--   ('BL', 'body { display: inline-block; width: auto; position: absolute; bottom: 0;}')
--   ('BM', 'body { display: inline-block; width: auto; position: absolute; bottom: 0; right: 50%; transform: translateX(+50%);}')
--   ('BR', 'body { display: inline-block; width: auto; position: absolute; bottom: 0; right: 0;}');

-- Dumped from database version 9.5.6
-- Dumped by pg_dump version 9.5.6

-- Started on 2017-12-20 12:25:41 UTC

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12361)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2265 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 564 (class 1247 OID 16387)
-- Name: enum_creation_states; Type: TYPE; Schema: public; Owner: wooble
--

CREATE TYPE enum_creation_states AS ENUM (
    'draft',
    'public',
    'private',
    'delete'
);


ALTER TYPE enum_creation_states OWNER TO wooble;

--
-- TOC entry 197 (class 1255 OID 16395)
-- Name: update_date(); Type: FUNCTION; Schema: public; Owner: wooble
--

CREATE FUNCTION update_date() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
NEW.updated_at := current_date;
RETURN NEW;
END;

$$;


ALTER FUNCTION public.update_date() OWNER TO wooble;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 181 (class 1259 OID 16396)
-- Name: app_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE app_user (
    id integer NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    created_at date DEFAULT ('now'::text)::date,
    updated_at date,
    is_creator boolean DEFAULT false,
    passwd text,
    salt_key text NOT NULL,
    customer_id text NOT NULL,
    fund integer DEFAULT 0 NOT NULL,
    deleted_at date,
    account_id text,
    pic_path text,
    codepen_name text,
    dribbble_name text,
    github_name text,
    twitter_name text,
    website text,
    fullname text NOT NULL,
    is_vip boolean DEFAULT false,
    is_active boolean DEFAULT false
);


ALTER TABLE app_user OWNER TO wooble;

--
-- TOC entry 182 (class 1259 OID 16405)
-- Name: app_user_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE app_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE app_user_id_seq OWNER TO wooble;

--
-- TOC entry 2266 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE app_user_id_seq OWNED BY app_user.id;


--
-- TOC entry 183 (class 1259 OID 16407)
-- Name: creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation (
    id integer NOT NULL,
    title text DEFAULT 'unknown'::bpchar NOT NULL,
    creator_id integer NOT NULL,
    updated_at date,
    engine text NOT NULL,
    thumb_path text,
    description text,
    alias text DEFAULT 'woobly'::text,
    state enum_creation_states DEFAULT 'draft'::enum_creation_states NOT NULL,
    old_creator_id integer,
    versions integer[],
    created_at timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone,
    is_featured boolean DEFAULT false,
    is_thumb_preview boolean DEFAULT true NOT NULL,
    preview_position_id text DEFAULT 'TL'::text
);


ALTER TABLE creation OWNER TO wooble;

--
-- TOC entry 193 (class 1259 OID 74081)
-- Name: creation_function; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation_function (
    creation_id integer NOT NULL,
    version integer NOT NULL,
    call text NOT NULL,
    detail text
);


ALTER TABLE creation_function OWNER TO wooble;

--
-- TOC entry 184 (class 1259 OID 16418)
-- Name: creation_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE creation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE creation_id_seq OWNER TO wooble;

--
-- TOC entry 2267 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE creation_id_seq OWNED BY creation.id;


--
-- TOC entry 192 (class 1259 OID 33116)
-- Name: creation_param; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation_param (
    creation_id integer NOT NULL,
    field text NOT NULL,
    value text,
    version integer NOT NULL
);


ALTER TABLE creation_param OWNER TO wooble;

--
-- TOC entry 185 (class 1259 OID 16427)
-- Name: engine; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE engine (
    name text NOT NULL,
    extension text NOT NULL,
    content_type text NOT NULL
);


ALTER TABLE engine OWNER TO wooble;

--
-- TOC entry 186 (class 1259 OID 16433)
-- Name: package; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package (
    id integer NOT NULL,
    user_id integer NOT NULL,
    title text NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    source text,
    referer text,
    build_required boolean DEFAULT true,
    built_at timestamp with time zone,
    nb_build integer
);


ALTER TABLE package OWNER TO wooble;

--
-- TOC entry 187 (class 1259 OID 16440)
-- Name: package_creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package_creation (
    package_id integer NOT NULL,
    creation_id integer NOT NULL,
    alias text,
    version integer,
    id integer NOT NULL
);


ALTER TABLE package_creation OWNER TO wooble;

--
-- TOC entry 195 (class 1259 OID 74126)
-- Name: package_creation_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE package_creation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE package_creation_id_seq OWNER TO wooble;

--
-- TOC entry 2268 (class 0 OID 0)
-- Dependencies: 195
-- Name: package_creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE package_creation_id_seq OWNED BY package_creation.id;


--
-- TOC entry 194 (class 1259 OID 74102)
-- Name: package_creation_param; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package_creation_param (
    field text NOT NULL,
    value text,
    package_creation_id integer NOT NULL
);


ALTER TABLE package_creation_param OWNER TO wooble;

--
-- TOC entry 188 (class 1259 OID 16447)
-- Name: package_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE package_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE package_id_seq OWNER TO wooble;

--
-- TOC entry 2269 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE package_id_seq OWNED BY package.id;


--
-- TOC entry 189 (class 1259 OID 16449)
-- Name: plan; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan (
    label text NOT NULL,
    price_per_month numeric DEFAULT 0 NOT NULL,
    price_per_year numeric DEFAULT 0 NOT NULL,
    nb_pkg integer DEFAULT 1,
    nb_crea integer DEFAULT 1,
    level smallint
);


ALTER TABLE plan OWNER TO wooble;

--
-- TOC entry 2270 (class 0 OID 0)
-- Dependencies: 189
-- Name: COLUMN plan.level; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN plan.level IS 'Specify at which level is the plan';


--
-- TOC entry 190 (class 1259 OID 16459)
-- Name: plan_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan_user (
    id integer NOT NULL,
    user_id integer NOT NULL,
    nb_renew smallint NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    start_date date DEFAULT ('now'::text)::date NOT NULL,
    end_date date,
    plan_label text NOT NULL,
    unsub_date date
);


ALTER TABLE plan_user OWNER TO wooble;

--
-- TOC entry 2271 (class 0 OID 0)
-- Dependencies: 190
-- Name: TABLE plan_user; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON TABLE plan_user IS 'History of user plans';


--
-- TOC entry 2272 (class 0 OID 0)
-- Dependencies: 190
-- Name: COLUMN plan_user.nb_renew; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN plan_user.nb_renew IS 'How many times the plan has been renewed';


--
-- TOC entry 191 (class 1259 OID 16467)
-- Name: plan_user_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE plan_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE plan_user_id_seq OWNER TO wooble;

--
-- TOC entry 2273 (class 0 OID 0)
-- Dependencies: 191
-- Name: plan_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE plan_user_id_seq OWNED BY plan_user.id;


--
-- TOC entry 196 (class 1259 OID 82294)
-- Name: preview_position; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE preview_position (
    position_id text NOT NULL,
    style_source text NOT NULL
);


ALTER TABLE preview_position OWNER TO wooble;

--
-- TOC entry 2051 (class 2604 OID 16469)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);


--
-- TOC entry 2057 (class 2604 OID 16470)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);


--
-- TOC entry 2063 (class 2604 OID 16471)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package ALTER COLUMN id SET DEFAULT nextval('package_id_seq'::regclass);


--
-- TOC entry 2065 (class 2604 OID 74128)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation ALTER COLUMN id SET DEFAULT nextval('package_creation_id_seq'::regclass);


--
-- TOC entry 2072 (class 2604 OID 16472)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user ALTER COLUMN id SET DEFAULT nextval('plan_user_id_seq'::regclass);


--
-- TOC entry 2242 (class 0 OID 16396)
-- Dependencies: 181
-- Data for Name: app_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY app_user (id, name, email, created_at, updated_at, is_creator, passwd, salt_key, customer_id, fund, deleted_at, account_id, pic_path, codepen_name, dribbble_name, github_name, twitter_name, website, fullname, is_vip, is_active) FROM stdin;
\.


--
-- TOC entry 2274 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('app_user_id_seq', 145, true);


--
-- TOC entry 2244 (class 0 OID 16407)
-- Dependencies: 183
-- Data for Name: creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation (id, title, creator_id, updated_at, engine, thumb_path, description, alias, state, old_creator_id, versions, created_at, is_featured, is_thumb_preview, preview_position_id) FROM stdin;
\.


--
-- TOC entry 2254 (class 0 OID 74081)
-- Dependencies: 193
-- Data for Name: creation_function; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation_function (creation_id, version, call, detail) FROM stdin;
\.


--
-- TOC entry 2275 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('creation_id_seq', 407, true);


--
-- TOC entry 2253 (class 0 OID 33116)
-- Dependencies: 192
-- Data for Name: creation_param; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation_param (creation_id, field, value, version) FROM stdin;
\.


--
-- TOC entry 2246 (class 0 OID 16427)
-- Dependencies: 185
-- Data for Name: engine; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY engine (name, extension, content_type) FROM stdin;
JS	.js	application/javascript
\.


--
-- TOC entry 2247 (class 0 OID 16433)
-- Dependencies: 186
-- Data for Name: package; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package (id, user_id, title, created_at, updated_at, source, referer, build_required, built_at, nb_build) FROM stdin;
\.


--
-- TOC entry 2248 (class 0 OID 16440)
-- Dependencies: 187
-- Data for Name: package_creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package_creation (package_id, creation_id, alias, version, id) FROM stdin;
\.


--
-- TOC entry 2276 (class 0 OID 0)
-- Dependencies: 195
-- Name: package_creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('package_creation_id_seq', 32, true);


--
-- TOC entry 2255 (class 0 OID 74102)
-- Dependencies: 194
-- Data for Name: package_creation_param; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package_creation_param (field, value, package_creation_id) FROM stdin;
\.


--
-- TOC entry 2277 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('package_id_seq', 139, true);


--
-- TOC entry 2250 (class 0 OID 16449)
-- Dependencies: 189
-- Data for Name: plan; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan (label, price_per_month, price_per_year, nb_pkg, nb_crea, level) FROM stdin;
Visitor	0	0	2	3	0
Architect	35	800	0	0	2
Woobler	1300	25000	5	8	1
\.


--
-- TOC entry 2251 (class 0 OID 16459)
-- Dependencies: 190
-- Data for Name: plan_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan_user (id, user_id, nb_renew, created_at, start_date, end_date, plan_label, unsub_date) FROM stdin;
\.


--
-- TOC entry 2278 (class 0 OID 0)
-- Dependencies: 191
-- Name: plan_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('plan_user_id_seq', 95, true);


--
-- TOC entry 2257 (class 0 OID 82294)
-- Dependencies: 196
-- Data for Name: preview_position; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY preview_position (position_id, style_source) FROM stdin;
TL
TM	body { display: inline-block; width: auto; position: absolute; left: 50%; transform: translateX(-50%);}
TR	body { float: right; }
ML	body { display: inline-block; width: auto; position: absolute; top: 50%; transform: translateY(-50%);}
M	body { display: inline-block; width: auto; position: absolute; left: 50%; top: 50%; transform: translate(-50%,-50%);}
MR	body { display: inline-block; width: auto; position: absolute; right: 0; top: 50%; transform: translateY(-50%);}
BL	body { display: inline-block; width: auto; position: absolute; bottom: 0;}
BM	body { display: inline-block; width: auto; position: absolute; bottom: 0; right: 50%; transform: translateX(+50%);}
BR	body { display: inline-block; width: auto; position: absolute; bottom: 0; right: 0;}
\.


--
-- TOC entry 2074 (class 2606 OID 16474)
-- Name: app_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2107 (class 2606 OID 74096)
-- Name: creation_function_pk; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_function
    ADD CONSTRAINT creation_function_pk PRIMARY KEY (creation_id, version, call);


--
-- TOC entry 2104 (class 2606 OID 41309)
-- Name: creation_param_pk; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_param
    ADD CONSTRAINT creation_param_pk PRIMARY KEY (creation_id, field, version);


--
-- TOC entry 2082 (class 2606 OID 16476)
-- Name: creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);


--
-- TOC entry 2076 (class 2606 OID 16480)
-- Name: customer_id; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT customer_id UNIQUE (customer_id);


--
-- TOC entry 2078 (class 2606 OID 16482)
-- Name: email; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT email UNIQUE (email);


--
-- TOC entry 2087 (class 2606 OID 16484)
-- Name: engine_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY engine
    ADD CONSTRAINT engine_pkey PRIMARY KEY (name);


--
-- TOC entry 2098 (class 2606 OID 16486)
-- Name: label_pk; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan
    ADD CONSTRAINT label_pk PRIMARY KEY (label);


--
-- TOC entry 2080 (class 2606 OID 16488)
-- Name: name; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT name UNIQUE (name);


--
-- TOC entry 2094 (class 2606 OID 74139)
-- Name: package_creation_id_uq; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_creation_id_uq UNIQUE (id);


--
-- TOC entry 2111 (class 2606 OID 74125)
-- Name: package_creation_param_pk; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation_param
    ADD CONSTRAINT package_creation_param_pk PRIMARY KEY (field, package_creation_id);


--
-- TOC entry 2096 (class 2606 OID 16490)
-- Name: package_creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_creation_pkey PRIMARY KEY (package_id, creation_id);


--
-- TOC entry 2090 (class 2606 OID 16492)
-- Name: package_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT package_pkey PRIMARY KEY (id);


--
-- TOC entry 2102 (class 2606 OID 16494)
-- Name: plan_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2113 (class 2606 OID 82301)
-- Name: preview_position_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY preview_position
    ADD CONSTRAINT preview_position_pkey PRIMARY KEY (position_id);


--
-- TOC entry 2088 (class 1259 OID 16495)
-- Name: fki_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_app_user_id_fk ON package USING btree (user_id);


--
-- TOC entry 2091 (class 1259 OID 16496)
-- Name: fki_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_fk ON package_creation USING btree (creation_id);


--
-- TOC entry 2108 (class 1259 OID 74094)
-- Name: fki_creation_id_function_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_function_fk ON creation_function USING btree (creation_id);


--
-- TOC entry 2105 (class 1259 OID 33129)
-- Name: fki_creation_id_param_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_param_fk ON creation_param USING btree (creation_id);


--
-- TOC entry 2083 (class 1259 OID 16497)
-- Name: fki_engine_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_engine_fk ON creation USING btree (engine);


--
-- TOC entry 2084 (class 1259 OID 16498)
-- Name: fki_fk_app_user_id; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);


--
-- TOC entry 2109 (class 1259 OID 74145)
-- Name: fki_package_creation_param_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_package_creation_param_fk ON package_creation_param USING btree (package_creation_id);


--
-- TOC entry 2092 (class 1259 OID 16499)
-- Name: fki_package_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_package_id_fk ON package_creation USING btree (package_id);


--
-- TOC entry 2099 (class 1259 OID 16500)
-- Name: fki_plan_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_app_user_id_fk ON plan_user USING btree (user_id);


--
-- TOC entry 2100 (class 1259 OID 16501)
-- Name: fki_plan_label_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_label_fk ON plan_user USING btree (plan_label);


--
-- TOC entry 2085 (class 1259 OID 82317)
-- Name: fki_preview_position_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_preview_position_fk ON creation USING btree (preview_position_id);


--
-- TOC entry 2125 (class 2620 OID 16504)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE OF name, email, is_creator ON app_user FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2126 (class 2620 OID 16505)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE OF title, description, versions ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2127 (class 2620 OID 16506)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE ON package FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2114 (class 2606 OID 16507)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES app_user(id);


--
-- TOC entry 2117 (class 2606 OID 16512)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2118 (class 2606 OID 16517)
-- Name: creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2123 (class 2606 OID 74097)
-- Name: creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_function
    ADD CONSTRAINT creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id) ON DELETE CASCADE;


--
-- TOC entry 2122 (class 2606 OID 33130)
-- Name: creation_id_param_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_param
    ADD CONSTRAINT creation_id_param_fk FOREIGN KEY (creation_id) REFERENCES creation(id) ON DELETE CASCADE;


--
-- TOC entry 2115 (class 2606 OID 16522)
-- Name: engine_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT engine_fk FOREIGN KEY (engine) REFERENCES engine(name);


--
-- TOC entry 2124 (class 2606 OID 74146)
-- Name: package_creation_param_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation_param
    ADD CONSTRAINT package_creation_param_fk FOREIGN KEY (package_creation_id) REFERENCES package_creation(id) ON DELETE CASCADE;


--
-- TOC entry 2119 (class 2606 OID 16527)
-- Name: package_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_id_fk FOREIGN KEY (package_id) REFERENCES package(id) ON DELETE CASCADE;


--
-- TOC entry 2120 (class 2606 OID 16532)
-- Name: plan_app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2121 (class 2606 OID 16537)
-- Name: plan_label_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_label_fk FOREIGN KEY (plan_label) REFERENCES plan(label);


--
-- TOC entry 2116 (class 2606 OID 82312)
-- Name: preview_position_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT preview_position_fk FOREIGN KEY (preview_position_id) REFERENCES preview_position(position_id);


--
-- TOC entry 2264 (class 0 OID 0)
-- Dependencies: 7
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-12-20 12:25:42 UTC

--
-- PostgreSQL database dump complete
--
