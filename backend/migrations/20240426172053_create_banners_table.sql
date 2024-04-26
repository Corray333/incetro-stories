-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.banners
(
    banner_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    created_at time with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT banners_pkey PRIMARY KEY (banner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.banners;
-- +goose StatementEnd
