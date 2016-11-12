PGDMP     ,    :             
    t            wooble    9.5.4    9.5.4 #    c           0    0    ENCODING    ENCODING     #   SET client_encoding = 'SQL_ASCII';
                       false            d           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                       false            e           1262    16384    wooble    DATABASE     v   CREATE DATABASE wooble WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';
    DROP DATABASE wooble;
             postgres    false                        2615    2200    public    SCHEMA        CREATE SCHEMA public;
    DROP SCHEMA public;
             postgres    false            f           0    0    SCHEMA public    COMMENT     6   COMMENT ON SCHEMA public IS 'standard public schema';
                  postgres    false    7            g           0    0    public    ACL     �   REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;
                  postgres    false    7                        3079    12361    plpgsql 	   EXTENSION     ?   CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
    DROP EXTENSION plpgsql;
                  false            h           0    0    EXTENSION plpgsql    COMMENT     @   COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';
                       false    1            �            1255    16386    update_date()    FUNCTION     �   CREATE FUNCTION update_date() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
NEW.updated_at := current_date;
RETURN NEW;
END;

$$;
 $   DROP FUNCTION public.update_date();
       public       wooble    false    7    1            �            1259    16387    app_user    TABLE     �   CREATE TABLE app_user (
    id integer NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    created_at date DEFAULT ('now'::text)::date,
    updated_at date,
    is_creator boolean,
    passwd text NOT NULL
);
    DROP TABLE public.app_user;
       public         wooble    false    7            �            1259    16394    app_user_id_seq    SEQUENCE     q   CREATE SEQUENCE app_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.app_user_id_seq;
       public       wooble    false    7    181            i           0    0    app_user_id_seq    SEQUENCE OWNED BY     5   ALTER SEQUENCE app_user_id_seq OWNED BY app_user.id;
            public       wooble    false    182            �            1259    16396    creation    TABLE     �  CREATE TABLE creation (
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
    DROP TABLE public.creation;
       public         wooble    false    7            �            1259    16408    creation_id_seq    SEQUENCE     q   CREATE SEQUENCE creation_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.creation_id_seq;
       public       wooble    false    183    7            j           0    0    creation_id_seq    SEQUENCE OWNED BY     5   ALTER SEQUENCE creation_id_seq OWNED BY creation.id;
            public       wooble    false    184            �            1259    16410    engine    TABLE     m   CREATE TABLE engine (
    name text NOT NULL,
    extension text NOT NULL,
    content_type text NOT NULL
);
    DROP TABLE public.engine;
       public         wooble    false    7            �           2604    16424    id    DEFAULT     \   ALTER TABLE ONLY app_user ALTER COLUMN id SET DEFAULT nextval('app_user_id_seq'::regclass);
 :   ALTER TABLE public.app_user ALTER COLUMN id DROP DEFAULT;
       public       wooble    false    182    181            �           2604    16425    id    DEFAULT     \   ALTER TABLE ONLY creation ALTER COLUMN id SET DEFAULT nextval('creation_id_seq'::regclass);
 :   ALTER TABLE public.creation ALTER COLUMN id DROP DEFAULT;
       public       wooble    false    184    183            \          0    16387    app_user 
   TABLE DATA               X   COPY app_user (id, name, email, created_at, updated_at, is_creator, passwd) FROM stdin;
    public       wooble    false    181   �#       k           0    0    app_user_id_seq    SEQUENCE SET     6   SELECT pg_catalog.setval('app_user_id_seq', 1, true);
            public       wooble    false    182            ^          0    16396    creation 
   TABLE DATA               �   COPY creation (id, title, creator_id, version, created_at, updated_at, has_document, has_script, has_style, engine) FROM stdin;
    public       wooble    false    183   $       l           0    0    creation_id_seq    SEQUENCE SET     7   SELECT pg_catalog.setval('creation_id_seq', 26, true);
            public       wooble    false    184            `          0    16410    engine 
   TABLE DATA               8   COPY engine (name, extension, content_type) FROM stdin;
    public       wooble    false    185   5$       �           2606    16428    app_user_pkey 
   CONSTRAINT     M   ALTER TABLE ONLY app_user
    ADD CONSTRAINT app_user_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.app_user DROP CONSTRAINT app_user_pkey;
       public         wooble    false    181    181            �           2606    16430    creation_pkey 
   CONSTRAINT     M   ALTER TABLE ONLY creation
    ADD CONSTRAINT creation_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.creation DROP CONSTRAINT creation_pkey;
       public         wooble    false    183    183            �           2606    16432    email 
   CONSTRAINT     C   ALTER TABLE ONLY app_user
    ADD CONSTRAINT email UNIQUE (email);
 8   ALTER TABLE ONLY public.app_user DROP CONSTRAINT email;
       public         wooble    false    181    181            �           2606    16434    engine_pkey 
   CONSTRAINT     K   ALTER TABLE ONLY engine
    ADD CONSTRAINT engine_pkey PRIMARY KEY (name);
 <   ALTER TABLE ONLY public.engine DROP CONSTRAINT engine_pkey;
       public         wooble    false    185    185            �           2606    16438    name 
   CONSTRAINT     A   ALTER TABLE ONLY app_user
    ADD CONSTRAINT name UNIQUE (name);
 7   ALTER TABLE ONLY public.app_user DROP CONSTRAINT name;
       public         wooble    false    181    181            �           2606    16440    title 
   CONSTRAINT     C   ALTER TABLE ONLY creation
    ADD CONSTRAINT title UNIQUE (title);
 8   ALTER TABLE ONLY public.creation DROP CONSTRAINT title;
       public         wooble    false    183    183            �           1259    16441    fki_engine_fk    INDEX     =   CREATE INDEX fki_engine_fk ON creation USING btree (engine);
 !   DROP INDEX public.fki_engine_fk;
       public         wooble    false    183            �           1259    16442    fki_fk_app_user_id    INDEX     F   CREATE INDEX fki_fk_app_user_id ON creation USING btree (creator_id);
 &   DROP INDEX public.fki_fk_app_user_id;
       public         wooble    false    183            �           2620    16444    update_date    TRIGGER     n   CREATE TRIGGER update_date BEFORE UPDATE OF version ON creation FOR EACH ROW EXECUTE PROCEDURE update_date();
 -   DROP TRIGGER update_date ON public.creation;
       public       wooble    false    183    186    183            �           2620    16445    update_date    TRIGGER     ~   CREATE TRIGGER update_date BEFORE UPDATE OF name, email, is_creator ON app_user FOR EACH ROW EXECUTE PROCEDURE update_date();
 -   DROP TRIGGER update_date ON public.app_user;
       public       wooble    false    181    186    181    181    181            �           2606    16446    app_user_id_fk    FK CONSTRAINT     n   ALTER TABLE ONLY creation
    ADD CONSTRAINT app_user_id_fk FOREIGN KEY (creator_id) REFERENCES app_user(id);
 A   ALTER TABLE ONLY public.creation DROP CONSTRAINT app_user_id_fk;
       public       wooble    false    2009    181    183            �           2606    16451 	   engine_fk    FK CONSTRAINT     e   ALTER TABLE ONLY creation
    ADD CONSTRAINT engine_fk FOREIGN KEY (engine) REFERENCES engine(name);
 <   ALTER TABLE ONLY public.creation DROP CONSTRAINT engine_fk;
       public       wooble    false    183    2021    185            \      x������ � �      ^      x������ � �      `   .   x��
v6���*�L,(��LN,�����J,K,N.�,(����� ��4     