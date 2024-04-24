-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    password character varying(60) COLLATE pg_catalog."default" NOT NULL,
    avatar text COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS user_token (
    user_id bigint NOT NULL, token text NOT NULL, PRIMARY KEY (user_id), FOREIGN KEY (user_id) REFERENCES public.users (user_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID
);

CREATE TABLE IF NOT EXISTS public.stories (
    stories_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY (
        INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1
    ), created_at timestamp
    with
        time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, creator bigint NOT NULL, CONSTRAINT stories_pkey PRIMARY KEY (stories_id)
);

CREATE TABLE IF NOT EXISTS public.banners
(
    banner_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    created_at time with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT banners_pkey PRIMARY KEY (banner_id)
);

CREATE TABLE IF NOT EXISTS public.story_banner (
    stories_id bigint NOT NULL, banner_id bigint NOT NULL, CONSTRAINT stories_banner_pkey PRIMARY KEY (banner_id, stories_id), CONSTRAINT stories_banner_banner_id_fkey FOREIGN KEY (banner_id) REFERENCES public.banners (banner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT stories_banner_stories_id_fkey FOREIGN KEY (stories_id) REFERENCES public.stories (stories_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS public.views (
    user_id bigint NOT NULL, banner_id bigint NOT NULL, CONSTRAINT views_pkey PRIMARY KEY (user_id, banner_id), CONSTRAINT views_banner_id_fkey FOREIGN KEY (banner_id) REFERENCES public.banners (banner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT views_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users (user_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- +goose Down
DROP TABLE users;
DROP TABLE stories;
DROP TABLE banners;
DROP TABLE story_banner;
DROP TABLE views;
