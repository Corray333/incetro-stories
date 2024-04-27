package storage

import (
	"strconv"

	"github.com/Corray333/univer_cs/internal/domains/story/types"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB) *Storage {

	return &Storage{db: db}
}

// InsertBanner inserts a new banner into the database and returns the id
func (s *Storage) InsertBanner(storyId string, uid int, banners []types.Banner) (int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	var story_id int
	if storyId == "" {
		var story types.Story
		story.Creator = uid

		rows, err := tx.Queryx(`
			INSERT INTO stories (creator) VALUES ($1) RETURNING story_id;
		`, story.Creator)
		if err != nil {
			return -1, err
		}
		for rows.Next() {
			err := rows.Scan(&story_id)
			if err != nil {
				return -1, err
			}
		}
	} else {
		// Check if the story exists
		story_id, _ = strconv.Atoi(storyId)
	}

	var banner_id int
	rows, err := tx.Queryx(`INSERT INTO banners VALUES (DEFAULT, DEFAULT) RETURNING banner_id;`)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err := rows.Scan(&banner_id)
		if err != nil {
			return -1, err
		}
	}
	defer rows.Close()
	_, err = tx.Exec(`
		INSERT INTO story_banner (story_id, banner_id) VALUES ($1, $2);
	`, story_id, banner_id)
	if err != nil {
		return -1, err
	}

	for _, banner := range banners {
		for _, lang := range banner.Langs {
			_, err := tx.Exec(`
				INSERT INTO banner_lang VALUES ($1, $2, $3, $4);
			`, banner_id, lang.Lang, lang.Title, lang.Description)
			if err != nil {
				return -1, err
			}
		}
	}
	tx.Commit()
	return banner_id, nil
}

// SelectStories returns all the stories from the database
func (s *Storage) SelectStories(story_id, banner_id, creator, offset, lang string) ([]types.Story, error) {
	var stories []types.Story
	type row struct {
		StoriesID       int    `db:"story_id"`
		BannerID        int    `db:"banner_id"`
		BannerTitle     string `db:"banner_title"`
		Description     string `db:"description"`
		StoryCreatedAt  string `db:"story_created_at"`
		BannerCreatedAt string `db:"banner_created_at"`
		Creator         int    `db:"creator"`
		UserID          int    `db:"user_id"`
		Lang            string `db:"lang"`
	}

	off, _ := strconv.Atoi(offset)
	where := squirrel.Eq{}
	if creator != "" {
		where["creator"] = creator
	}
	if banner_id != "" {
		where["banners.banner_id"] = banner_id
	}
	if story_id != "" {
		where["stories.story_id"] = story_id
	}
	if lang != "" && lang != "all" {
		where["lang"] = lang
	}

	query := squirrel.Select("stories.story_id, banners.banner_id, banner_lang.title AS banner_title, banner_lang.description, lang, stories.created_at AS story_created_at, banners.created_at AS banner_created_at, stories.creator").
		From("story_banner").
		Join("banners on banners.banner_id = story_banner.banner_id ").
		Join("stories ON stories.story_id = story_banner.story_id").
		Join("banner_lang ON banner_lang.banner_id = banners.banner_id").
		Where(where).OrderBy("stories.story_id", "banners.banner_id").Offset(uint64(off))
	sql, args, _ := query.PlaceholderFormat(squirrel.Dollar).ToSql()

	rows, err := s.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	var r row
	counterStory := -1
	for rows.Next() {
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		var story types.Story
		story.ID = r.StoriesID
		story.Creator = r.Creator
		story.CreatedAt = r.StoryCreatedAt

		lang := types.BannerLang{Lang: r.Lang, Title: r.BannerTitle, Description: r.Description}

		banner := types.Banner{
			ID:        r.BannerID,
			Langs:     []types.BannerLang{lang},
			CreatedAt: r.BannerCreatedAt,
		}

		if counterStory >= 0 && story.ID == stories[counterStory].ID {
			if banner.ID == stories[counterStory].Banners[len(stories[counterStory].Banners)-1].ID {
				stories[counterStory].Banners[len(stories[counterStory].Banners)-1].Langs = append(stories[counterStory].Banners[len(stories[counterStory].Banners)-1].Langs, lang)
			} else {
				stories[counterStory].Banners = append(stories[counterStory].Banners, banner)
			}
		} else {
			story.Banners = append(story.Banners, banner)
			stories = append(stories, story)
			counterStory++
		}
	}
	return stories, nil
}

// InsertView inserts a new view into the database
func (s *Storage) InsertView(user_id int, banner_id string) error {
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

// func (s *Storage) UpdateBanner(banner types.Banner) error {
// 	_, err := s.db.Queryx(`
// 		UPDATE banners SET title = $1, description = $2 WHERE banner_id = $3;
// 	`, banner.Title, banner.Description, banner.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
