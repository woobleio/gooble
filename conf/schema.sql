--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.4

-- Started on 2016-11-16 01:26:29 UTC

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'SQL_ASCII';
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
-- TOC entry 2181 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 189 (class 1255 OID 16386)
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
-- TOC entry 181 (class 1259 OID 16387)
-- Name: app_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE app_user (
    id integer NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    created_at date DEFAULT ('now'::text)::date,
    updated_at date,
    is_creator boolean,
    passwd text NOT NULL
);


ALTER TABLE app_user OWNER TO wooble;

--
-- TOC entry 182 (class 1259 OID 16394)
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
-- TOC entry 2182 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE app_user_id_seq OWNED BY app_user.id;


--
-- TOC entry 183 (class 1259 OID 16396)
-- Name: creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation (
    id integer NOT NULL,
    title text DEFAULT 'toto'::bpchar NOT NULL,
    creator_id integer NOT NULL,
    version text DEFAULT 1.0 NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    has_document boolean DEFAULT false NOT NULL,
    has_script boolean DEFAULT false NOT NULL,
    has_style boolean DEFAULT false NOT NULL,
    engine text NOT NULL
);


ALTER TABLE creation OWNER TO wooble;

--
-- TOC entry 184 (class 1259 OID 16408)
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
-- TOC entry 2183 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE creation_id_seq OWNED BY creation.id;


--
-- TOC entry 185 (class 1259 OID 16410)
-- Name: engine; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE engine (
    name text NOT NULL,
    extension text NOT NULL,
    content_type text NOT NULL
);


ALTER TABLE engine OWNER TO wooble;

--
-- TOC entry 186 (class 1259 OID 16416)
-- Name: package; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package (
    id integer NOT NULL,
    user_id integer NOT NULL,
    title text NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date
);


ALTER TABLE package OWNER TO wooble;

--
-- TOC entry 187 (class 1259 OID 16422)
-- Name: package_creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE package_creation (
    package_id integer NOT NULL,
    creation_id integer NOT NULL
);


ALTER TABLE package_creation OWNER TO wooble;

--
-- TOC entry 188 (class 1259 OID 16425)
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
-- TOC entry 2184 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE package_id_seq OWNED BY package.id;


--
-- TOC entry 2011 (class 2604 OID 16427)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);


--
-- TOC entry 2018 (class 2604 OID 16428)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);


--
-- TOC entry 2019 (class 2604 OID 16429)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package ALTER COLUMN id SET DEFAULT nextval('package_id_seq'::regclass);


--
-- TOC entry 2166 (class 0 OID 16387)
-- Dependencies: 181
-- Data for Name: app_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY app_user (id, name, email, created_at, updated_at, is_creator, passwd) FROM stdin;
\.


--
-- TOC entry 2185 (class 0 OID 0)
-- Dependencies: 182
-- Name: app_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('app_user_id_seq', 5, true);


--
-- TOC entry 2168 (class 0 OID 16396)
-- Dependencies: 183
-- Data for Name: creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation (id, title, creator_id, version, created_at, updated_at, has_document, has_script, has_style, engine) FROM stdin;
\.


--
-- TOC entry 2186 (class 0 OID 0)
-- Dependencies: 184
-- Name: creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('creation_id_seq', 52, true);


--
-- TOC entry 2170 (class 0 OID 16410)
-- Dependencies: 185
-- Data for Name: engine; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY engine (name, extension, content_type) FROM stdin;
JSES5	.js	application/javascript
\.


--
-- TOC entry 2171 (class 0 OID 16416)
-- Dependencies: 186
-- Data for Name: package; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package (id, user_id, title, created_at, updated_at) FROM stdin;
\.


--
-- TOC entry 2172 (class 0 OID 16422)
-- Dependencies: 187
-- Data for Name: package_creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY package_creation (package_id, creation_id) FROM stdin;
\.


--
-- TOC entry 2187 (class 0 OID 0)
-- Dependencies: 188
-- Name: package_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('package_id_seq', 4, true);


--
-- TOC entry 2022 (class 2606 OID 16431)
-- Name: app_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2028 (class 2606 OID 16433)
-- Name: creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);


--
-- TOC entry 2024 (class 2606 OID 16435)
-- Name: email; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT email UNIQUE (email);


--
-- TOC entry 2034 (class 2606 OID 16437)
-- Name: engine_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY engine
    ADD CONSTRAINT engine_pkey PRIMARY KEY (name);


--
-- TOC entry 2026 (class 2606 OID 16439)
-- Name: name; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY app_user
    ADD CONSTRAINT name UNIQUE (name);


--
-- TOC entry 2043 (class 2606 OID 16441)
-- Name: package_creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT "package_creation_pkey" PRIMARY KEY (package_id, creation_id);


--
-- TOC entry 2037 (class 2606 OID 16443)
-- Name: package_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT package_pkey PRIMARY KEY (id);


--
-- TOC entry 2032 (class 2606 OID 16445)
-- Name: title_creator; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT title_creator UNIQUE (title, creator_id);


--
-- TOC entry 2039 (class 2606 OID 16447)
-- Name: user_title; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT user_title UNIQUE (title, user_id);


--
-- TOC entry 2035 (class 1259 OID 16448)
-- Name: fki_app_user_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_app_user_id_fk ON package USING btree (user_id);


--
-- TOC entry 2040 (class 1259 OID 16449)
-- Name: fki_creation_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_creation_id_fk ON package_creation USING btree (creation_id);


--
-- TOC entry 2029 (class 1259 OID 16450)
-- Name: fki_engine_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_engine_fk ON creation USING btree (engine);


--
-- TOC entry 2030 (class 1259 OID 16451)
-- Name: fki_fk_app_user_id; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);


--
-- TOC entry 2041 (class 1259 OID 16452)
-- Name: fki_package_id_fk; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_package_id_fk ON package_creation USING btree (package_id);


--
-- TOC entry 2050 (class 2620 OID 16453)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF version ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2049 (class 2620 OID 16454)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF name, email, is_creator ON app_user FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2051 (class 2620 OID 16455)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date AFTER UPDATE ON package FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2044 (class 2606 OID 16456)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES app_user(id);


--
-- TOC entry 2046 (class 2606 OID 16461)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (user_id) REFERENCES app_user(id);


--
-- TOC entry 2047 (class 2606 OID 16466)
-- Name: creation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT creation_id_fk FOREIGN KEY (creation_id) REFERENCES creation(id);


--
-- TOC entry 2045 (class 2606 OID 16471)
-- Name: engine_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT engine_fk FOREIGN KEY (engine) REFERENCES engine(name);


--
-- TOC entry 2048 (class 2606 OID 16476)
-- Name: package_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY package_creation
    ADD CONSTRAINT package_id_fk FOREIGN KEY (package_id) REFERENCES package(id);


--
-- TOC entry 2180 (class 0 OID 0)
-- Dependencies: 7
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2016-11-16 01:26:30 UTC

--
-- PostgreSQL database dump complete
--
