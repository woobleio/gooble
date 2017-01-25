--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.5
-- Dumped by pg_dump version 9.5.5

-- Started on 2017-01-25 23:35:39 UTC

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
-- TOC entry 2230 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 194 (class 1255 OID 16386)
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

--
-- TOC entry 195 (class 1255 OID 16635)
-- Name: update_renew(); Type: FUNCTION; Schema: public; Owner: wooble
--

CREATE FUNCTION update_renew() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
NEW.renewed_at := current_date;
RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_renew() OWNER TO wooble;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 181 (class 1259 OID 16387)
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
    current_plan_id integer
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
-- TOC entry 2231 (class 0 OID 0)
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
    price numeric NOT NULL,
    thumb_url text,
    description text,
    is_unlisted boolean NOT NULL
);


ALTER TABLE creation OWNER TO wooble;

--
-- TOC entry 2232 (class 0 OID 0)
-- Dependencies: 183
-- Name: COLUMN creation.is_unlisted; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN creation.is_unlisted IS 'True when the creator chose to delete it. The creation still exists but isn''t listed anymore.';


--
-- TOC entry 184 (class 1259 OID 16410)
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
-- TOC entry 2233 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE creation_id_seq OWNED BY creation.id;


--
-- TOC entry 193 (class 1259 OID 16680)
-- Name: creation_purchase; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation_purchase (
    user_id integer NOT NULL,
    creation_id integer NOT NULL,
    price numeric NOT NULL,
    purchased_at date DEFAULT ('now'::text)::date NOT NULL
);


ALTER TABLE creation_purchase OWNER TO wooble;

--
-- TOC entry 185 (class 1259 OID 16412)
-- Name: engine; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE engine (
    name text NOT NULL,
    extension text NOT NULL,
    content_type text NOT NULL
);


ALTER TABLE engine OWNER TO wooble;

--
-- TOC entry 186 (class 1259 OID 16418)
-- Name: package; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package (
    id integer NOT NULL,
    user_id integer NOT NULL,
    title text NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    domains text[] NOT NULL,
    key text NOT NULL,
    source text
);


ALTER TABLE package OWNER TO wooble;

--
-- TOC entry 187 (class 1259 OID 16425)
-- Name: package_creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package_creation (
    package_id integer NOT NULL,
    creation_id integer NOT NULL
);


ALTER TABLE package_creation OWNER TO wooble;

--
-- TOC entry 188 (class 1259 OID 16428)
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
-- TOC entry 2234 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE package_id_seq OWNED BY package.id;


--
-- TOC entry 190 (class 1259 OID 16608)
-- Name: plan; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan (
    id integer NOT NULL,
    label text NOT NULL,
    price_per_month numeric DEFAULT 0 NOT NULL,
    price_per_year numeric DEFAULT 0 NOT NULL
);


ALTER TABLE plan OWNER TO wooble;

--
-- TOC entry 189 (class 1259 OID 16606)
-- Name: plan_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE plan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE plan_id_seq OWNER TO wooble;

--
-- TOC entry 2235 (class 0 OID 0)
-- Dependencies: 189
-- Name: plan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE plan_id_seq OWNED BY plan.id;


--
-- TOC entry 192 (class 1259 OID 16623)
-- Name: plan_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE plan_user (
    id integer NOT NULL,
    user_id integer NOT NULL,
    plan_id integer NOT NULL,
    nb_renew smallint NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    start_date date DEFAULT ('now'::text)::date NOT NULL,
    end_date date NOT NULL,
    renewed_at date
);


ALTER TABLE plan_user OWNER TO wooble;

--
-- TOC entry 2236 (class 0 OID 0)
-- Dependencies: 192
-- Name: TABLE plan_user; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON TABLE plan_user IS 'History of user plans';


--
-- TOC entry 2237 (class 0 OID 0)
-- Dependencies: 192
-- Name: COLUMN plan_user.nb_renew; Type: COMMENT; Schema: public; Owner: wooble
--

COMMENT ON COLUMN plan_user.nb_renew IS 'How many times the plan has been renewed';


--
-- TOC entry 191 (class 1259 OID 16621)
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
-- TOC entry 2238 (class 0 OID 0)
-- Dependencies: 191
-- Name: plan_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE plan_user_id_seq OWNED BY plan_user.id;


--
-- TOC entry 2031 (class 2604 OID 16430)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);


--
-- TOC entry 2038 (class 2604 OID 16431)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);


--
-- TOC entry 2040 (class 2604 OID 16432)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package ALTER COLUMN id SET DEFAULT nextval('package_id_seq'::regclass);


--
-- TOC entry 2041 (class 2604 OID 16611)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan ALTER COLUMN id SET DEFAULT nextval('plan_id_seq'::regclass);


--
-- TOC entry 2044 (class 2604 OID 16626)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user ALTER COLUMN id SET DEFAULT nextval('plan_user_id_seq'::regclass);


--
-- TOC entry 2210 (class 0 OID 16387)
-- Dependencies: 181
-- Data for Name: app_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY app_user (id, name, email, created_at, updated_at, is_creator, passwd, salt_key, current_plan_id) FROM stdin;
\.


--
-- TOC entry 2239 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('app_user_id_seq', 13, true);


--
-- TOC entry 2212 (class 0 OID 16398)
-- Dependencies: 183
-- Data for Name: creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation (id, title, creator_id, version, created_at, updated_at, has_document, has_script, has_style, engine, price, thumb_url, description, is_unlisted) FROM stdin;
\.


--
-- TOC entry 2240 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('creation_id_seq', 116, true);


--
-- TOC entry 2222 (class 0 OID 16680)
-- Dependencies: 193
-- Data for Name: creation_purchase; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation_purchase (user_id, creation_id, price, purchased_at) FROM stdin;
\.


--
-- TOC entry 2214 (class 0 OID 16412)
-- Dependencies: 185
-- Data for Name: engine; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY engine (name, extension, content_type) FROM stdin;
JSES5	.js	application/javascript
\.


--
-- TOC entry 2215 (class 0 OID 16418)
-- Dependencies: 186
-- Data for Name: package; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package (id, user_id, title, created_at, updated_at, domains, key, source) FROM stdin;
\.


--
-- TOC entry 2216 (class 0 OID 16425)
-- Dependencies: 187
-- Data for Name: package_creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package_creation (package_id, creation_id) FROM stdin;
\.


--
-- TOC entry 2241 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('package_id_seq', 35, true);


--
-- TOC entry 2219 (class 0 OID 16608)
-- Dependencies: 190
-- Data for Name: plan; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan (id, label, price_per_month, price_per_year) FROM stdin;
3	premium	20	230.33
\.


--
-- TOC entry 2242 (class 0 OID 0)
-- Dependencies: 189
-- Name: plan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('plan_id_seq', 3, true);


--
-- TOC entry 2221 (class 0 OID 16623)
-- Dependencies: 192
-- Data for Name: plan_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY plan_user (id, user_id, plan_id, nb_renew, created_at, start_date, end_date, renewed_at) FROM stdin;
\.


--
-- TOC entry 2243 (class 0 OID 0)
-- Dependencies: 191
-- Name: plan_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('plan_user_id_seq', 1, false);


--
-- TOC entry 2049 (class 2606 OID 16434)
-- Name: app_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2055 (class 2606 OID 16436)
-- Name: creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);


--
-- TOC entry 2078 (class 2606 OID 16687)
-- Name: creation_purchase_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT creation_purchase_pkey PRIMARY KEY (user_id, creation_id);


--
-- TOC entry 2051 (class 2606 OID 16438)
-- Name: email; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT email UNIQUE (email);


--
-- TOC entry 2061 (class 2606 OID 16440)
-- Name: engine_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY engine
    ADD CONSTRAINT engine_pkey PRIMARY KEY (name);


--
-- TOC entry 2053 (class 2606 OID 16442)
-- Name: name; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT name UNIQUE (name);


--
-- TOC entry 2070 (class 2606 OID 16444)
-- Name: package_creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_creation_pkey PRIMARY KEY (package_id, creation_id);


--
-- TOC entry 2064 (class 2606 OID 16446)
-- Name: package_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT package_pkey PRIMARY KEY (id);


--
-- TOC entry 2072 (class 2606 OID 16616)
-- Name: plan_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan
    ADD CONSTRAINT plan_pkey PRIMARY KEY (id);


--
-- TOC entry 2076 (class 2606 OID 16631)
-- Name: plan_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2059 (class 2606 OID 16448)
-- Name: title_creator; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT title_creator UNIQUE (title, creator_id);


--
-- TOC entry 2066 (class 2606 OID 16450)
-- Name: user_title; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT user_title UNIQUE (title, user_id);


--
-- TOC entry 2062 (class 1259 OID 16451)
-- Name: fki_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_app_user_id_fk ON package USING btree (user_id);


--
-- TOC entry 2067 (class 1259 OID 16452)
-- Name: fki_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_fk ON package_creation USING btree (creation_id);


--
-- TOC entry 2056 (class 1259 OID 16453)
-- Name: fki_engine_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_engine_fk ON creation USING btree (engine);


--
-- TOC entry 2057 (class 1259 OID 16454)
-- Name: fki_fk_app_user_id; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);


--
-- TOC entry 2068 (class 1259 OID 16455)
-- Name: fki_package_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_package_id_fk ON package_creation USING btree (package_id);


--
-- TOC entry 2073 (class 1259 OID 16674)
-- Name: fki_plan_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_app_user_id_fk ON plan_user USING btree (user_id);


--
-- TOC entry 2074 (class 1259 OID 16663)
-- Name: fki_plan_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_plan_id_fk ON plan_user USING btree (plan_id);


--
-- TOC entry 2079 (class 1259 OID 16700)
-- Name: fki_purchase_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_purchase_creation_id_fk ON creation_purchase USING btree (creation_id);


--
-- TOC entry 2080 (class 1259 OID 16694)
-- Name: fki_purchase_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_purchase_user_id_fk ON creation_purchase USING btree (user_id);


--
-- TOC entry 2093 (class 2620 OID 16457)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF version ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2092 (class 2620 OID 16458)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF name, email, is_creator ON app_user FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2094 (class 2620 OID 16459)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE ON package FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2095 (class 2620 OID 16639)
-- Name: update_renewed_at; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_renewed_at AFTER UPDATE OF renew ON plan_user FOR EACH ROW EXECUTE PROCEDURE update_renew();


--
-- TOC entry 2082 (class 2606 OID 16460)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES app_user(id);


--
-- TOC entry 2084 (class 2606 OID 16465)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2088 (class 2606 OID 16664)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2085 (class 2606 OID 16470)
-- Name: creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2083 (class 2606 OID 16475)
-- Name: engine_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT engine_fk FOREIGN KEY (engine) REFERENCES engine(name);


--
-- TOC entry 2086 (class 2606 OID 16485)
-- Name: package_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_id_fk FOREIGN KEY (package_id) REFERENCES package(id);


--
-- TOC entry 2089 (class 2606 OID 16669)
-- Name: plan_app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2087 (class 2606 OID 16658)
-- Name: plan_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY plan_user
    ADD CONSTRAINT plan_id_fk FOREIGN KEY (plan_id) REFERENCES plan(id);


--
-- TOC entry 2091 (class 2606 OID 16695)
-- Name: purchase_creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT purchase_creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2090 (class 2606 OID 16689)
-- Name: purchase_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation_purchase
    ADD CONSTRAINT purchase_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2081 (class 2606 OID 16675)
-- Name: user_plan_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT user_plan_id_fk FOREIGN KEY (id) REFERENCES plan_user(id);


--
-- TOC entry 2229 (class 0 OID 0)
-- Dependencies: 7
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-01-25 23:35:40 UTC

--
-- PostgreSQL database dump complete
--

