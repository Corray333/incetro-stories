-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.banner_views (
    user_id bigint NOT NULL, 
    banner_id bigint NOT NULL, 
    CONSTRAINT views_pkey PRIMARY KEY (user_id, banner_id), 
    CONSTRAINT views_banner_id_fkey FOREIGN KEY (banner_id) 
        REFERENCES public.banners (banner_id) MATCH SIMPLE 
            ON UPDATE NO ACTION 
            ON DELETE NO ACTION, 
    CONSTRAINT views_user_id_fkey FOREIGN KEY (user_id) 
        REFERENCES public.users (user_id) MATCH SIMPLE 
        ON UPDATE NO ACTION 
        ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.banner_views;
-- +goose StatementEnd
