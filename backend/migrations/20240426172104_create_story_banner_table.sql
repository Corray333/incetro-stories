-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.story_banner (
    story_id bigint NOT NULL, banner_id bigint NOT NULL, CONSTRAINT stories_banner_pkey PRIMARY KEY (banner_id, story_id), CONSTRAINT stories_banner_banner_id_fkey FOREIGN KEY (banner_id) REFERENCES public.banners (banner_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT stories_banner_story_id_fkey FOREIGN KEY (story_id) REFERENCES public.stories (story_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.story_banner;
-- +goose StatementEnd
