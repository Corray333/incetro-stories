package storage

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
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
func (s *Storage) InsertBanner(project_id string, storyId string, uid int, banner types.Banner, file multipart.File, fileName string) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var story_id int
	if storyId == "" {
		var story types.Story
		story.Creator = uid

		rows, err := tx.Queryx(`
			INSERT INTO stories (creator, project_id) VALUES ($1, $2) RETURNING story_id;
		`, story.Creator, project_id)
		if err != nil {
			return err
		}
		for rows.Next() {
			err := rows.Scan(&story_id)
			if err != nil {
				return err
			}
		}
	} else {
		// Check if the story exists
		story_id, _ = strconv.Atoi(storyId)
		rows, err := tx.Queryx(`SELECT story_id FROM stories JOIN project_user on stories.project_id = project_user.project_id WHERE story_id = $1 AND creator = $2;`, story_id, uid)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return fmt.Errorf("story does not exist or user does not have access to it")
		}
	}

	var banner_id int
	rows, err := tx.Queryx(`INSERT INTO banners VALUES (DEFAULT, DEFAULT) RETURNING banner_id;`)
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&banner_id)
		if err != nil {
			return err
		}
	}
	defer rows.Close()
	_, err = tx.Exec(`
		INSERT INTO story_banner (story_id, banner_id) VALUES ($1, $2);
	`, story_id, banner_id)
	if err != nil {
		return err
	}

	for _, lang := range banner.Langs {
		_, err := tx.Exec(`
			INSERT INTO banner_lang VALUES ($1, $2, $3, $4);
		`, banner_id, lang.Lang, lang.Title, lang.Description)
		if err != nil {
			return err
		}
	}

	randomStr := ""
	for i := 0; i < 10; i++ {
		randomStr += strconv.Itoa(rand.Intn(10))
	}

	fileName = strconv.Itoa(banner_id) + "-" + randomStr + filepath.Ext(fileName)

	newFile, _ := os.Create("../files/images/banners/banner" + fileName)
	defer newFile.Close()
	if _, err := io.Copy(newFile, file); err != nil {
		return err
	}

	// Update the media url of the banner
	// TODO: replace with url from .env
	if _, err = tx.Exec(`UPDATE banners SET media_url = $1 WHERE banner_id = $2;`, "http://localhost:3001/images/banners/banner"+fileName, banner_id); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// SelectStories returns all the stories from the database
func (s *Storage) SelectStories(project_id, story_id, banner_id, creator, offset, lang string) ([]types.Story, error) {
	var stories []types.Story
	type row struct {
		StoriesID       int    `db:"story_id"`
		BannerID        int    `db:"banner_id"`
		BannerTitle     string `db:"banner_title"`
		Description     string `db:"description"`
		StoryCreatedAt  string `db:"story_created_at"`
		BannerCreatedAt string `db:"banner_created_at"`
		MediaURL        string `db:"media_url"`
		Creator         int    `db:"creator"`
		UserID          int    `db:"user_id"`
		Lang            string `db:"lang"`
	}

	off, _ := strconv.Atoi(offset)
	where := squirrel.Eq{}
	where["project_id"] = project_id
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

	query := squirrel.Select("stories.story_id, banners.banner_id, banner_lang.title AS banner_title, banner_lang.description, lang, stories.created_at AS story_created_at, banners.created_at AS banner_created_at, media_url, stories.creator").
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
			MediaURL:  r.MediaURL,
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
