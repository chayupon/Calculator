\connect calculator
SET TIME ZONE 'UTC';
DROP TABLE IF EXISTS public.history;
CREATE SEQUENCE history_sequence_seq;
CREATE TABLE public.history
(
    sequence smallint NOT NULL DEFAULT nextval('history_sequence_seq'::regclass),
    "time" timestamp with time zone NOT NULL,
    age numeric(10,0) ,
    -- operate character varying(1) COLLATE pg_catalog."default" NOT NULL,
    subprovince character varying(2000) COLLATE pg_catalog."default" ,
    -- result numeric(10,4) NOT NULL,
    province character varying(2000) COLLATE pg_catalog."default" ,
    CONSTRAINT history_pkey PRIMARY KEY (sequence)
);

ALTER TABLE public.history OWNER TO postgres;
