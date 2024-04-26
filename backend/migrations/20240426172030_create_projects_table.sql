-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects
(
    project_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    created_at time with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator bigint NOT NULL,
    cover text COLLATE pg_catalog."default",
    CONSTRAINT projects_pkey PRIMARY KEY (project_id),
    CONSTRAINT projects_creator_fkey FOREIGN KEY (creator)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
