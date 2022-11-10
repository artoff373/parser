--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.5

-- Started on 2022-11-07 13:36:21

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
-- TOC entry 3329 (class 1262 OID 16565)
-- Name: Search; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "Search" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'Russian_Russia.1251';


ALTER DATABASE "Search" OWNER TO postgres;

\connect "Search"

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
-- TOC entry 4 (class 2615 OID 16566)
-- Name: Search; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA "Search";


ALTER SCHEMA "Search" OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 16591)
-- Name: Posts; Type: TABLE; Schema: Search; Owner: postgres
--

CREATE TABLE "Search"."Posts" (
    id integer NOT NULL,
    title text,
    text text,
    relev real,
    url text,
    is_in_report boolean,
    fresh boolean,
    profile_id integer,
    source_id integer,
    dictionary text[],
    pub_date text,
    search_date text
);


ALTER TABLE "Search"."Posts" OWNER TO postgres;

--
-- TOC entry 3330 (class 0 OID 0)
-- Dependencies: 215
-- Name: TABLE "Posts"; Type: COMMENT; Schema: Search; Owner: postgres
--

COMMENT ON TABLE "Search"."Posts" IS 'found news';


--
-- TOC entry 214 (class 1259 OID 16590)
-- Name: Posts_id_seq; Type: SEQUENCE; Schema: Search; Owner: postgres
--

CREATE SEQUENCE "Search"."Posts_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE "Search"."Posts_id_seq" OWNER TO postgres;

--
-- TOC entry 3331 (class 0 OID 0)
-- Dependencies: 214
-- Name: Posts_id_seq; Type: SEQUENCE OWNED BY; Schema: Search; Owner: postgres
--

ALTER SEQUENCE "Search"."Posts_id_seq" OWNED BY "Search"."Posts".id;


--
-- TOC entry 211 (class 1259 OID 16568)
-- Name: Profile; Type: TABLE; Schema: Search; Owner: postgres
--

CREATE TABLE "Search"."Profile" (
    id integer NOT NULL,
    name text,
    keys text[],
    last_search text
);


ALTER TABLE "Search"."Profile" OWNER TO postgres;

--
-- TOC entry 3332 (class 0 OID 0)
-- Dependencies: 211
-- Name: TABLE "Profile"; Type: COMMENT; Schema: Search; Owner: postgres
--

COMMENT ON TABLE "Search"."Profile" IS 'Search profiles';


--
-- TOC entry 210 (class 1259 OID 16567)
-- Name: Profile_id_seq; Type: SEQUENCE; Schema: Search; Owner: postgres
--

CREATE SEQUENCE "Search"."Profile_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE "Search"."Profile_id_seq" OWNER TO postgres;

--
-- TOC entry 3333 (class 0 OID 0)
-- Dependencies: 210
-- Name: Profile_id_seq; Type: SEQUENCE OWNED BY; Schema: Search; Owner: postgres
--

ALTER SEQUENCE "Search"."Profile_id_seq" OWNED BY "Search"."Profile".id;


--
-- TOC entry 213 (class 1259 OID 16577)
-- Name: Sources; Type: TABLE; Schema: Search; Owner: postgres
--

CREATE TABLE "Search"."Sources" (
    id integer NOT NULL,
    name text,
    url text,
    selector text,
    profile_id integer
);


ALTER TABLE "Search"."Sources" OWNER TO postgres;

--
-- TOC entry 3334 (class 0 OID 0)
-- Dependencies: 213
-- Name: TABLE "Sources"; Type: COMMENT; Schema: Search; Owner: postgres
--

COMMENT ON TABLE "Search"."Sources" IS 'Sources for searching news';


--
-- TOC entry 212 (class 1259 OID 16576)
-- Name: Sources_id_seq; Type: SEQUENCE; Schema: Search; Owner: postgres
--

CREATE SEQUENCE "Search"."Sources_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE "Search"."Sources_id_seq" OWNER TO postgres;

--
-- TOC entry 3335 (class 0 OID 0)
-- Dependencies: 212
-- Name: Sources_id_seq; Type: SEQUENCE OWNED BY; Schema: Search; Owner: postgres
--

ALTER SEQUENCE "Search"."Sources_id_seq" OWNED BY "Search"."Sources".id;


--
-- TOC entry 3177 (class 2604 OID 18165)
-- Name: Posts id; Type: DEFAULT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Posts" ALTER COLUMN id SET DEFAULT nextval('"Search"."Posts_id_seq"'::regclass);


--
-- TOC entry 3175 (class 2604 OID 18166)
-- Name: Profile id; Type: DEFAULT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Profile" ALTER COLUMN id SET DEFAULT nextval('"Search"."Profile_id_seq"'::regclass);


--
-- TOC entry 3176 (class 2604 OID 18167)
-- Name: Sources id; Type: DEFAULT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Sources" ALTER COLUMN id SET DEFAULT nextval('"Search"."Sources_id_seq"'::regclass);


--
-- TOC entry 3183 (class 2606 OID 16598)
-- Name: Posts Posts_pkey; Type: CONSTRAINT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Posts"
    ADD CONSTRAINT "Posts_pkey" PRIMARY KEY (id);


--
-- TOC entry 3179 (class 2606 OID 16575)
-- Name: Profile Profile_pkey; Type: CONSTRAINT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Profile"
    ADD CONSTRAINT "Profile_pkey" PRIMARY KEY (id);


--
-- TOC entry 3181 (class 2606 OID 16584)
-- Name: Sources Sources_pkey; Type: CONSTRAINT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Sources"
    ADD CONSTRAINT "Sources_pkey" PRIMARY KEY (id);


--
-- TOC entry 3184 (class 2606 OID 16585)
-- Name: Sources profile link; Type: FK CONSTRAINT; Schema: Search; Owner: postgres
--

ALTER TABLE ONLY "Search"."Sources"
    ADD CONSTRAINT "profile link" FOREIGN KEY (profile_id) REFERENCES "Search"."Profile"(id);


-- Completed on 2022-11-07 13:36:21

--
-- PostgreSQL database dump complete
--

