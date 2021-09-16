--
-- PostgreSQL database dump
--

-- Dumped from database version 13.2 (Debian 13.2-1.pgdg100+1)
-- Dumped by pg_dump version 13.2 (Debian 13.2-1.pgdg100+1)

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
-- Name: game; Type: TABLE; Schema: public; Owner: poc
--

CREATE TABLE public.game (
    id bigint NOT NULL,
    score integer DEFAULT 0 NOT NULL,
    max_times_through_deck integer DEFAULT 1000 NOT NULL
);


ALTER TABLE public.game OWNER TO poc;

--
-- Name: game_id_seq; Type: SEQUENCE; Schema: public; Owner: poc
--

CREATE SEQUENCE public.game_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.game_id_seq OWNER TO poc;

--
-- Name: game_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: poc
--

ALTER SEQUENCE public.game_id_seq OWNED BY public.game.id;


--
-- Name: history; Type: TABLE; Schema: public; Owner: poc
--

CREATE TABLE public.history (
    id bigint NOT NULL,
    game_id bigint NOT NULL,
    move_id bigint NOT NULL,
    move_number integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.history OWNER TO poc;

--
-- Name: history_id_seq; Type: SEQUENCE; Schema: public; Owner: poc
--

CREATE SEQUENCE public.history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.history_id_seq OWNER TO poc;

--
-- Name: history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: poc
--

ALTER SEQUENCE public.history_id_seq OWNED BY public.history.id;


--
-- Name: move; Type: TABLE; Schema: public; Owner: poc
--

CREATE TABLE public.move (
    id bigint NOT NULL,
    old_pile_num smallint DEFAULT 0 NOT NULL,
    old_pile_index smallint DEFAULT 0 NOT NULL,
    old_pile_position smallint DEFAULT 0 NOT NULL,
    new_pile_num smallint DEFAULT 0 NOT NULL,
    new_pile_index smallint DEFAULT 0 NOT NULL,
    new_pile_position smallint DEFAULT 0 NOT NULL
);


ALTER TABLE public.move OWNER TO poc;

--
-- Name: move_id_seq; Type: SEQUENCE; Schema: public; Owner: poc
--

CREATE SEQUENCE public.move_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.move_id_seq OWNER TO poc;

--
-- Name: move_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: poc
--

ALTER SEQUENCE public.move_id_seq OWNED BY public.move.id;


--
-- Name: pile_card; Type: TABLE; Schema: public; Owner: poc
--

CREATE TABLE public.pile_card (
    id bigint NOT NULL,
    pile_num smallint DEFAULT 0 NOT NULL,
    pile_index smallint DEFAULT 0 NOT NULL,
    suit smallint DEFAULT 0 NOT NULL,
    index smallint DEFAULT 0 NOT NULL,
    "position" integer DEFAULT 0 NOT NULL,
    game_id bigint NOT NULL
);


ALTER TABLE public.pile_card OWNER TO poc;

--
-- Name: pile_card_id_seq; Type: SEQUENCE; Schema: public; Owner: poc
--

CREATE SEQUENCE public.pile_card_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pile_card_id_seq OWNER TO poc;

--
-- Name: pile_card_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: poc
--

ALTER SEQUENCE public.pile_card_id_seq OWNED BY public.pile_card.id;


--
-- Name: game id; Type: DEFAULT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.game ALTER COLUMN id SET DEFAULT nextval('public.game_id_seq'::regclass);


--
-- Name: history id; Type: DEFAULT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.history ALTER COLUMN id SET DEFAULT nextval('public.history_id_seq'::regclass);


--
-- Name: move id; Type: DEFAULT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.move ALTER COLUMN id SET DEFAULT nextval('public.move_id_seq'::regclass);


--
-- Name: pile_card id; Type: DEFAULT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.pile_card ALTER COLUMN id SET DEFAULT nextval('public.pile_card_id_seq'::regclass);


--
-- Name: game game_pkey; Type: CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.game
    ADD CONSTRAINT game_pkey PRIMARY KEY (id);


--
-- Name: history history_game_move_uniqueness; Type: CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_game_move_uniqueness UNIQUE (game_id, move_id);


--
-- Name: history history_pkey; Type: CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_pkey PRIMARY KEY (id);


--
-- Name: move move_pkey; Type: CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.move
    ADD CONSTRAINT move_pkey PRIMARY KEY (id);


--
-- Name: pile_card pile_card_pkey; Type: CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.pile_card
    ADD CONSTRAINT pile_card_pkey PRIMARY KEY (id);


--
-- Name: history history_game_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id) ON DELETE CASCADE;


--
-- Name: history history_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.move(id) ON DELETE RESTRICT;


--
-- Name: pile_card pile_card_game_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: poc
--

ALTER TABLE ONLY public.pile_card
    ADD CONSTRAINT pile_card_game_id_fkey FOREIGN KEY (game_id) REFERENCES public.game(id);


--
-- PostgreSQL database dump complete
--

