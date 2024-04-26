-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner_lang (
    banner_id INT NOT NULL,
    lang VARCHAR(10) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (banner_id, lang),
    FOREIGN KEY (banner_id) REFERENCES public.banners (banner_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.banner_lang;
-- +goose StatementEnd
