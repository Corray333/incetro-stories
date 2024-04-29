-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.banners
(
    banner_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    created_at bigint NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP),
    media_url text NOT NULL DEFAULT '',
    CONSTRAINT banners_pkey PRIMARY KEY (banner_id)
);
CREATE INDEX IF NOT EXISTS banners_created_at_idx ON public.banners (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.banners;
DROP INDEX IF EXISTS public.banners_created_at_idx;
-- +goose StatementEnd
