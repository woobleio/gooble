--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.4

-- Started on 2016-10-28 21:28:24 UTC

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
-- TOC entry 2152 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 187 (class 1255 OID 16509)
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
-- TOC entry 184 (class 1259 OID 16409)
-- Name: creation; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE creation (
    id integer NOT NULL,
    title text DEFAULT 'toto'::bpchar NOT NULL,
    creator_id integer NOT NULL,
    version text DEFAULT 1.0 NOT NULL,
    created_at date DEFAULT ('now'::text)::date NOT NULL,
    updated_at date,
    source_id integer NOT NULL
);


ALTER TABLE creation OWNER TO wooble;

--
-- TOC entry 183 (class 1259 OID 16407)
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
-- TOC entry 2153 (class 0 OID 0)
-- Dependencies: 183
-- Name: creation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE creation_id_seq OWNED BY creation.id;


--
-- TOC entry 186 (class 1259 OID 16481)
-- Name: source; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE source (
    id integer NOT NULL,
    host text NOT NULL
);


ALTER TABLE source OWNER TO wooble;

--
-- TOC entry 185 (class 1259 OID 16479)
-- Name: source_id_seq; Type: SEQUENCE; Schema: public; Owner: wooble
--

CREATE SEQUENCE source_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE source_id_seq OWNER TO wooble;

--
-- TOC entry 2154 (class 0 OID 0)
-- Dependencies: 185
-- Name: source_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE source_id_seq OWNED BY source.id;


--
-- TOC entry 182 (class 1259 OID 16390)
-- Name: app_user; Type: TABLE; Schema: public; Owner: wooble
--

CREATE TABLE "app_user" (
    id integer NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    created_at date DEFAULT ('now'::text)::date,
    updated_at date,
    is_creator boolean,
    passwd text NOT NULL
);


ALTER TABLE "app_user" OWNER TO wooble;

--
-- TOC entry 181 (class 1259 OID 16388)
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
-- TOC entry 2155 (class 0 OID 0)
-- Dependencies: 181
-- Name: app_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wooble
--

ALTER SEQUENCE app_user_id_seq OWNED BY "app_user".id;


--
-- TOC entry 2003 (class 2604 OID 16412)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);


--
-- TOC entry 2007 (class 2604 OID 16484)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY source ALTER COLUMN id SET DEFAULT nextval('source_id_seq'::regclass);


--
-- TOC entry 2001 (class 2604 OID 16393)
-- Name: id; Type: DEFAULT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY "app_user" ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);


--
-- TOC entry 2142 (class 0 OID 16409)
-- Dependencies: 184
-- Data for Name: creation; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY creation (id, title, creator_id, version, created_at, updated_at, source_id) FROM stdin;
\.


--
-- TOC entry 2156 (class 0 OID 0)
-- Dependencies: 183
-- Name: creation_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('creation_id_seq', 1, false);


--
-- TOC entry 2144 (class 0 OID 16481)
-- Dependencies: 186
-- Data for Name: source; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY source (id, host) FROM stdin;
\.


--
-- TOC entry 2157 (class 0 OID 0)
-- Dependencies: 185
-- Name: source_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('source_id_seq', 1, false);


--
-- TOC entry 2140 (class 0 OID 16390)
-- Dependencies: 182
-- Data for Name: app_user; Type: TABLE DATA; Schema: public; Owner: wooble
--

COPY "app_user" (id, name, email, created_at, updated_at, is_creator, passwd) FROM stdin;
\.


--
-- TOC entry 2158 (class 0 OID 0)
-- Dependencies: 181
-- Name: app_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wooble
--

SELECT pg_catalog.setval('app_user_id_seq', 1, false);


--
-- TOC entry 2015 (class 2606 OID 16417)
-- Name: creation_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);


--
-- TOC entry 2009 (class 2606 OID 16506)
-- Name: email; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY "app_user"
    ADD CONSTRAINT email UNIQUE (email);


--
-- TOC entry 2011 (class 2606 OID 16504)
-- Name: name; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY "app_user"
    ADD CONSTRAINT name UNIQUE (name);


--
-- TOC entry 2020 (class 2606 OID 16489)
-- Name: source_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY source
    ADD CONSTRAINT source_pkey PRIMARY KEY (id);


--
-- TOC entry 2018 (class 2606 OID 16502)
-- Name: title; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT title UNIQUE (title);


--
-- TOC entry 2013 (class 2606 OID 16398)
-- Name: app_user_pkey; Type: CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY "app_user"
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);


--
-- TOC entry 2016 (class 1259 OID 16495)
-- Name: fki_fk_app_user_id; Type: INDEX; Schema: public; Owner: wooble
--

CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);


--
-- TOC entry 2024 (class 2620 OID 16517)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF version ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2023 (class 2620 OID 16518)
-- Name: update_date; Type: TRIGGER; Schema: public; Owner: wooble
--

CREATE TRIGGER update_date BEFORE UPDATE OF name, email, is_creator ON "app_user" FOR EACH ROW EXECUTE PROCEDURE update_date();


--
-- TOC entry 2022 (class 2606 OID 16496)
-- Name: source_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT source_id_fk FOREIGN KEY (source_id) REFERENCES source(id);


--
-- TOC entry 2021 (class 2606 OID 16490)
-- Name: app_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wooble
--

ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES "app_user"(id);


--
-- TOC entry 2151 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2016-10-28 21:28:25 UTC

--
-- PostgreSQL database dump complete
--
