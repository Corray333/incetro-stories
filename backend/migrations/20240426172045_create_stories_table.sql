-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.stories (
    story_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY (
        INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1
    ), 
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP), 
    creator bigint NOT NULL, 
    CONSTRAINT stories_pkey PRIMARY KEY (story_id)
);
CREATE INDEX IF NOT EXISTS stories_creator_idx ON public.stories (creator);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.stories;
DROP INDEX IF EXISTS stories_creator_idx;
-- +goose StatementEnd
