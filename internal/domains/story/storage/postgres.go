package storage

import (
	"time"

	"github.com/Corray333/univer_cs/internal/domains/story/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB) (Storage, error) {

	_, err := db.Query(`
		CREATE TABLE IF NOT EXISTS public.stories
		(
			stories_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			creator bigint NOT NULL,
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
	return Storage{db: db}, err
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
		INSERT INTO stories (created_at, creator) VALUES (CURRENT_TIMESTAMP, $1) RETURNING stories_id;
	`, story.Creator)
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
		Creator         int64  `db:"creator"`
		UserID          int64  `db:"user_id"`
	}

	rows, err := s.db.Queryx(`
	SELECT stories.stories_id, banners.banner_id, banners.name AS banner_name, description, stories.created_at AS story_created_at, banners.created_at AS banner_created_at, stories.creator
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
		story.Creator = r.Creator
		story.CreatedAt = r.StoryCreatedAt
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

// UpdateBannerMedia updates the media url of the banner
func (s *Storage) UpdateBannerMedia(bannerId string, mediaURL string) error {
	_, err := s.db.Queryx(`
		UPDATE banners SET media = $1 WHERE banner_id = $2;
	`, mediaURL, bannerId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBannerName updates the name of the banner
func (s *Storage) UpdateBannerName(bannerId string, name string) error {
	_, err := s.db.Queryx(`
		UPDATE banners SET name = $1 WHERE banner_id = $2;
	`, name, bannerId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBannerDescription updates the description of the banner
func (s *Storage) UpdateBannerDescription(bannerId string, description string) error {
	_, err := s.db.Queryx(`
		UPDATE banners SET description = $1 WHERE banner_id = $2;
	`, description, bannerId)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBannerTimestamp updates the timestamp of the banner
func (s *Storage) UpdateBannerTimestamp(bannerId string, timestamp time.Time) error {
	_, err := s.db.Queryx(`
		UPDATE stories SET expires_at = $1 WHERE banner_id = $2;
	`, timestamp.String(), bannerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateBanner(banner types.Banner) error {
	_, err := s.db.Queryx(`
		UPDATE banners SET name = $1, description = $2 WHERE banner_id = $3;
	`, banner.Name, banner.Description, banner.ID)
	if err != nil {
		return err
	}
	return nil
}
