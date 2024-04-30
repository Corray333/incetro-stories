package storage

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/Corray333/univer_cs/internal/domains/project/types"
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

func (s Storage) InsertProject(uid int, cover multipart.File, project types.Project) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var project_id int
	rows, err := tx.Queryx("INSERT INTO projects (name, description, creator) VALUES ($1, $2, $3) RETURNING project_id", project.Name, project.Description, uid)
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(&project_id)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	randomStr := ""
	for i := 0; i < 10; i++ {
		randomStr += strconv.Itoa(rand.Intn(10))
	}

	fileName := strconv.Itoa(project_id) + "-" + randomStr + ".png"

	newFile, err := os.Create("../files/images/projects/project" + fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err = io.Copy(newFile, cover); err != nil {
		return err
	}

	if _, err = tx.Exec("UPDATE projects SET cover = $1 WHERE project_id = $2", "images/projects/project"+fileName, project_id); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO project_user VALUES ($1, $2)", project_id, uid); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetProjects(project_id string, uid int) ([]types.Project, error) {
	var projects []types.Project
	if project_id != "" {
		err := s.db.Select(&projects, "SELECT projects.project_id, projects.name, projects.description, projects.cover FROM projects INNER JOIN project_user ON projects.project_id = project_user.project_id WHERE project_user.user_id = $1 AND projects.project_id = $2", uid, project_id)
		if err != nil {
			return nil, err
		}
		return projects, nil
	}
	err := s.db.Select(&projects, "SELECT projects.project_id, projects.name, projects.description, projects.cover FROM projects INNER JOIN project_user ON projects.project_id = project_user.project_id WHERE project_user.user_id = $1", uid)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
