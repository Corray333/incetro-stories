package storage

import (
	"github.com/Corray333/stories/internal/domains/story/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB) (*Storage, error) {

	_, err := db.Query(`
		CREATE TABLE IF NOT EXISTS public.stories
		(
			stories_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT stories_pkey PRIMARY KEY (stories_id)
		);
		CREATE TABLE IF NOT EXISTS public.banners
		(
			banner_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
			name text COLLATE pg_catalog."default",
			description text COLLATE pg_catalog."default",
			created_at time with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT banners_pkey PRIMARY KEY (banner_id)
		);
		CREATE TABLE IF NOT EXISTS public.stories_banner
		(
			stories_id bigint NOT NULL,
			banner_id bigint NOT NULL,
			CONSTRAINT stories_banner_pkey PRIMARY KEY (banner_id, stories_id),
			CONSTRAINT stories_banner_banner_id_fkey FOREIGN KEY (banner_id)
				REFERENCES public.banners (banner_id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION,
			CONSTRAINT stories_banner_stories_id_fkey FOREIGN KEY (stories_id)
				REFERENCES public.stories (stories_id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
		);
		CREATE TABLE IF NOT EXISTS public.views
		(
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
	`)
	return &Storage{db: db}, err
}

// InsertBanner inserts a new banner into the database and returns the id
func (s *Storage) InsertBanner(storyId string, banner types.Banner) (int64, error) {
	var id int64
	rows, err := s.db.Queryx(`
	INSERT INTO banners (name, description) VALUES ($1, $2) RETURNING banner_id;
`, banner.Name, banner.Description)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	_, err = s.db.Queryx(`
	INSERT INTO stories_banner (stories_id, banner_id) VALUES ($1, $2);
`, storyId, id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// InsertStory inserts a new story into the database and returns the id
func (s *Storage) InsertStory(story types.Story) (int64, error) {
	var id int64
	rows, err := s.db.Queryx(`
		INSERT INTO stories (created_at) VALUES (CURRENT_TIMESTAMP) RETURNING stories_id;
	`)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

// SelectStories returns all the stories from the database
func (s *Storage) SelectStories(filters string) ([]types.Story, error) {
	var stories []types.Story
	type row struct {
		StoriesID       int64  `db:"stories_id"`
		BannerID        int64  `db:"banner_id"`
		BannerName      string `db:"banner_name"`
		Description     string `db:"description"`
		StoryCreatedAt  string `db:"story_created_at"`
		BannerCreatedAt string `db:"banner_created_at"`
		UserID          int64  `db:"user_id"`
	}

	rows, err := s.db.Queryx(`
	SELECT stories.stories_id, banners.banner_id, banners.name AS banner_name, description, stories.created_at AS story_created_at, banners.created_at AS banner_created_at
	FROM stories_banner INNER JOIN banners ON banners.banner_id = stories_banner.banner_id 
	INNER JOIN stories ON stories.stories_id = stories_banner.stories_id ` + filters + ";")
	if err != nil {
		return nil, err
	}
	var r row
	counter := -1
	for rows.Next() {
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		var story types.Story
		story.ID = r.StoriesID
		story.CreatedAt = r.StoryCreatedAt
		// story.Name = r.StoriesName
		// story.UserID = r.UserID
		banner := types.Banner{
			ID:          r.BannerID,
			Name:        r.BannerName,
			Description: r.Description,
			CreatedAt:   r.BannerCreatedAt,
		}
		if counter >= 0 && story.ID == stories[counter].ID {
			stories[counter].Banners = append(stories[counter].Banners, banner)
		} else {
			story.Banners = append(story.Banners, banner)
			stories = append(stories, story)
			counter++
		}
	}
	return stories, nil
}

// InsertView inserts a new view into the database
func (s *Storage) InsertView(user_id int64, banner_id string) error {
	_, err := s.db.Queryx(`
		INSERT INTO views VALUES ($1, $2);
	`, user_id, banner_id)
	if err != nil {
		return err
	}
	return nil
}
