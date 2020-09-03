\connect calculator
SET TIME ZONE 'UTC';
DROP TABLE IF EXISTS public.history;
CREATE SEQUENCE history_sequence_seq;
CREATE TABLE public.history
(
    sequence smallint NOT NULL DEFAULT nextval('history_sequence_seq'::regclass),
    "time" timestamp with time zone NOT NULL,
    input1 numeric(10,4) NOT NULL,
    operate character varying(1) COLLATE pg_catalog."default" NOT NULL,
    input2 numeric(10,4) NOT NULL,
    result numeric(10,4) NOT NULL,
    errordescripe character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT history_pkey PRIMARY KEY (sequence)
);

ALTER TABLE public.history OWNER TO postgres;
