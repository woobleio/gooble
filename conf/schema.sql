--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.5
-- Dumped by pg_dump version 9.5.5

-- Started on 2017-02-13 18:01:03 UTC

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
-- TOC entry 2227 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 193 (class 1255 OID 16386)
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
-- TOC entry 181 (class 1259 OID 16388)
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
    total_due integer NOT NULL DEFAULT 0
);


ALTER TABLE app_user OWNER TO wooble;

--
-- TOC entry 182 (class 1259 OID 16396)
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
-- TOC entry 2228 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE app_user_id_seq OWNED BY app_user.id;


--
-- TOC entry 183 (class 1259 OID 16398)
-- Name: creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation (
    id integer NOT NULL,
    title text DEFAULT 'unknown'::bpchar NOT NULL,
    creator_id integer NOT NULL,
    version text DEFAULT 1.0 NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    has_document boolean DEFAULT false NOT NULL,
    has_script boolean DEFAULT false NOT NULL,
    has_style boolean DEFAULT false NOT NULL,
    engine text NOT NULL,
    price integer DEFAULT 0 NOT NULL,
    thumb_url text,
    description text,
    is_unlisted boolean DEFAULT false NOT NULL
);


ALTER TABLE creation OWNER TO wooble;

--
-- TOC entry 2229 (class 0 OID 0)
-- Dependencies: 183
-- Name: COLUMN creation.is_unlisted; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN creation.is_unlisted IS 'True when the creator chose to delete it. The creation still exists but isn''t listed anymore.';


--
-- TOC entry 184 (class 1259 OID 16412)
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
-- TOC entry 2230 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE creation_id_seq OWNED BY creation.id;


--
-- TOC entry 185 (class 1259 OID 16414)
-- Name: creation_purchase; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation_purchase (
    user_id integer NOT NULL,
    creation_id integer NOT NULL,
    total integer NOT NULL,
    purchased_at date DEFAULT ('now'::text)::date NOT NULL,
    charge_id text NOT NULL
);


ALTER TABLE creation_purchase OWNER TO wooble;

--
-- TOC entry 186 (class 1259 OID 16421)
-- Name: engine; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE engine (
    name text NOT NULL,
    extension text NOT NULL,
    content_type text NOT NULL
);


ALTER TABLE engine OWNER TO wooble;

--
-- TOC entry 187 (class 1259 OID 16427)
-- Name: package; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package (
    id integer NOT NULL,
    user_id integer NOT NULL,
    title text NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    domains text[] NOT NULL,
    source text
);


ALTER TABLE package OWNER TO wooble;

--
-- TOC entry 188 (class 1259 OID 16434)
-- Name: package_creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package_creation (
    package_id integer NOT NULL,
    creation_id integer NOT NULL,
    alias text
);


ALTER TABLE package_creation OWNER TO wooble;

--
-- TOC entry 189 (class 1259 OID 16440)
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
-- TOC entry 2231 (class 0 OID 0)
-- Dependencies: 189
-- Name: package_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE package_id_seq OWNED BY package.id;


--
-- TOC entry 190 (class 1259 OID 16442)
-- Name: plan; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan (
    label text NOT NULL,
    price_per_month numeric DEFAULT 0 NOT NULL,
    price_per_year numeric DEFAULT 0 NOT NULL,
    nb_pkg integer DEFAULT 1,
    nb_crea integer DEFAULT 1,
    nb_domains integer DEFAULT 1
);


ALTER TABLE plan OWNER TO wooble;

--
-- TOC entry 191 (class 1259 OID 16452)
-- Name: plan_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan_user (
    id integer NOT NULL,
    user_id integer NOT NULL,
    nb_renew smallint NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    start_date date DEFAULT ('now'::text)::date NOT NULL,
    end_date date NOT NULL,
    plan_label text NOT NULL
);


ALTER TABLE plan_user OWNER TO wooble;

--
-- TOC entry 2232 (class 0 OID 0)
-- Dependencies: 191
-- Name: TABLE plan_user; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON TABLE plan_user IS 'History of user plans';


--
-- TOC entry 2233 (class 0 OID 0)
-- Dependencies: 191
-- Name: COLUMN plan_user.nb_renew; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN plan_user.nb_renew IS 'How many times the plan has been renewed';


--
-- TOC entry 192 (class 1259 OID 16457)
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
-- TOC entry 2234 (class 0 OID 0)
-- Dependencies: 192
-- Name: plan_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE plan_user_id_seq OWNED BY plan_user.id;


--
-- TOC entry 2030 (class 2604 OID 16459)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);


--
-- TOC entry 2038 (class 2604 OID 16460)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);


--
-- TOC entry 2042 (class 2604 OID 16461)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package ALTER COLUMN id SET DEFAULT nextval('package_id_seq'::regclass);


--
-- TOC entry 2050 (class 2604 OID 16463)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user ALTER COLUMN id SET DEFAULT nextval('plan_user_id_seq'::regclass);


--
-- TOC entry 2208 (class 0 OID 16388)
-- Dependencies: 181
-- Data for Name: app_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY app_user (id, name, email, created_at, updated_at, is_creator, passwd, salt_key, customer_id, total_due) FROM stdin;
\.


--
-- TOC entry 2235 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('app_user_id_seq', 47, true);


--
-- TOC entry 2210 (class 0 OID 16398)
-- Dependencies: 183
-- Data for Name: creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation (id, title, creator_id, version, created_at, updated_at, has_document, has_script, has_style, engine, price, thumb_url, description, is_unlisted) FROM stdin;
\.


--
-- TOC entry 2236 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('creation_id_seq', 131, true);


--
-- TOC entry 2212 (class 0 OID 16414)
-- Dependencies: 185
-- Data for Name: creation_purchase; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation_purchase (user_id, creation_id, total, purchased_at, charge_id) FROM stdin;
\.


--
-- TOC entry 2213 (class 0 OID 16421)
-- Dependencies: 186
-- Data for Name: engine; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY engine (name, extension, content_type) FROM stdin;
JSES5	.js	application/javascript
\.


--
-- TOC entry 2214 (class 0 OID 16427)
-- Dependencies: 187
-- Data for Name: package; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package (id, user_id, title, created_at, updated_at, domains, source) FROM stdin;
\.


--
-- TOC entry 2215 (class 0 OID 16434)
-- Dependencies: 188
-- Data for Name: package_creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package_creation (package_id, creation_id, alias) FROM stdin;
\.


--
-- TOC entry 2237 (class 0 OID 0)
-- Dependencies: 189
-- Name: package_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('package_id_seq', 44, true);


--
-- TOC entry 2217 (class 0 OID 16442)
-- Dependencies: 190
-- Data for Name: plan; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan (label, price_per_month, price_per_year, nb_pkg, nb_crea, nb_domains) FROM stdin;
free	0	0	1	10	1
premium	2000	21525	2	20	2
pro	10000	110000	0	0	0
\.


--
-- TOC entry 2218 (class 0 OID 16452)
-- Dependencies: 191
-- Data for Name: plan_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan_user (id, user_id, nb_renew, created_at, start_date, end_date, plan_label) FROM stdin;
\.


--
-- TOC entry 2238 (class 0 OID 0)
-- Dependencies: 192
-- Name: plan_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('plan_user_id_seq', 15, true);


--
-- TOC entry 2052 (class 2606 OID 16465)
-- Name: app_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2060 (class 2606 OID 16467)
-- Name: creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);


--
-- TOC entry 2064 (class 2606 OID 16469)
-- Name: creation_purchase_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT creation_purchase_pkey PRIMARY KEY (user_id, creation_id);


--
-- TOC entry 2054 (class 2606 OID 16766)
-- Name: customer_id; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT customer_id UNIQUE (customer_id);


--
-- TOC entry 2056 (class 2606 OID 16471)
-- Name: email; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT email UNIQUE (email);


--
-- TOC entry 2068 (class 2606 OID 16473)
-- Name: engine_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY engine
    ADD CONSTRAINT engine_pkey PRIMARY KEY (name);


--
-- TOC entry 2077 (class 2606 OID 16728)
-- Name: label_pk; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan
    ADD CONSTRAINT label_pk PRIMARY KEY (label);


--
-- TOC entry 2058 (class 2606 OID 16475)
-- Name: name; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT name UNIQUE (name);


--
-- TOC entry 2075 (class 2606 OID 16477)
-- Name: package_creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_creation_pkey PRIMARY KEY (package_id, creation_id);


--
-- TOC entry 2071 (class 2606 OID 16479)
-- Name: package_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT package_pkey PRIMARY KEY (id);


--
-- TOC entry 2081 (class 2606 OID 16483)
-- Name: plan_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2069 (class 1259 OID 16484)
-- Name: fki_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_app_user_id_fk ON package USING btree (user_id);


--
-- TOC entry 2072 (class 1259 OID 16485)
-- Name: fki_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_fk ON package_creation USING btree (creation_id);


--
-- TOC entry 2061 (class 1259 OID 16486)
-- Name: fki_engine_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_engine_fk ON creation USING btree (engine);


--
-- TOC entry 2062 (class 1259 OID 16487)
-- Name: fki_fk_app_user_id; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);


--
-- TOC entry 2073 (class 1259 OID 16488)
-- Name: fki_package_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_package_id_fk ON package_creation USING btree (package_id);


--
-- TOC entry 2078 (class 1259 OID 16489)
-- Name: fki_plan_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_app_user_id_fk ON plan_user USING btree (user_id);


--
-- TOC entry 2079 (class 1259 OID 16758)
-- Name: fki_plan_label_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_label_fk ON plan_user USING btree (plan_label);


--
-- TOC entry 2065 (class 1259 OID 16491)
-- Name: fki_purchase_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_purchase_creation_id_fk ON creation_purchase USING btree (creation_id);


--
-- TOC entry 2066 (class 1259 OID 16492)
-- Name: fki_purchase_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_purchase_user_id_fk ON creation_purchase USING btree (user_id);


--
-- TOC entry 2092 (class 2620 OID 16494)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE OF version ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2091 (class 2620 OID 16495)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE OF name, email, is_creator ON app_user FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2093 (class 2620 OID 16496)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE ON package FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2082 (class 2606 OID 16498)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES app_user(id);


--
-- TOC entry 2086 (class 2606 OID 16503)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2087 (class 2606 OID 16513)
-- Name: creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2083 (class 2606 OID 16518)
-- Name: engine_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT engine_fk FOREIGN KEY (engine) REFERENCES engine(name);


--
-- TOC entry 2088 (class 2606 OID 16523)
-- Name: package_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_id_fk FOREIGN KEY (package_id) REFERENCES package(id);


--
-- TOC entry 2089 (class 2606 OID 16528)
-- Name: plan_app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2090 (class 2606 OID 16753)
-- Name: plan_label_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_label_fk FOREIGN KEY (plan_label) REFERENCES plan(label);


--
-- TOC entry 2084 (class 2606 OID 16538)
-- Name: purchase_creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT purchase_creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2085 (class 2606 OID 16543)
-- Name: purchase_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT purchase_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2226 (class 0 OID 0)
-- Dependencies: 7
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-02-13 18:01:04 UTC

--
-- PostgreSQL database dump complete
--

