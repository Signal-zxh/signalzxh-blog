package db

import (
	"database/sql"
	"log"

	"github.com/Signal-zxh/signalzxh-blog/model"
)

type TagRepo interface {
	GetTags() ([]model.Tag, error)
	GetTagByID(id int) (model.Tag, error)
	CreateTag(name string) (int64, error)
	UpdateTag(id int, name string) error
	DeleteTag(id int) error
	GetOrCreateTag(name string) (int64, error)
	GetTagsByPostID(postID int) ([]model.Tag, error)
	AddTagsToPost(postID int, tagIDs []int) error
	RemoveTagsFromPost(postID int) error
}

type tagRepo struct{}

var TagRepoImpl = &tagRepo{}

func (r *tagRepo) GetTags() ([]model.Tag, error) {
	rows, err := DB.Query("SELECT id, name FROM tags ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepo) GetTagByID(id int) (model.Tag, error) {
	row := DB.QueryRow("SELECT id, name FROM tags WHERE id = ?", id)
	var t model.Tag
	err := row.Scan(&t.ID, &t.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Tag{}, ErrNotFound
		}
		return model.Tag{}, err
	}
	return t, nil
}

func (r *tagRepo) CreateTag(name string) (int64, error) {
	res, err := DB.Exec("INSERT INTO tags(name) VALUES(?)", name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *tagRepo) UpdateTag(id int, name string) error {
	res, err := DB.Exec("UPDATE tags SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (r *tagRepo) DeleteTag(id int) error {
	res, err := DB.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (r *tagRepo) GetOrCreateTag(name string) (int64, error) {
	row := DB.QueryRow("SELECT id FROM tags WHERE name = ?", name)
	var id int64
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return r.CreateTag(name)
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *tagRepo) GetTagsByPostID(postID int) ([]model.Tag, error) {
	rows, err := DB.Query(`
		SELECT t.id, t.name 
		FROM tags t 
		JOIN post_tags pt ON t.id = pt.tag_id 
		WHERE pt.post_id = ? 
		ORDER BY t.id DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepo) AddTagsToPost(postID int, tagIDs []int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	if err := r.RemoveTagsFromPost(postID); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", rbErr)
		}
		return err
	}

	for _, tagID := range tagIDs {
		_, err := tx.Exec("INSERT INTO post_tags(post_id, tag_id) VALUES(?, ?)", postID, tagID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
				log.Printf("tx rollback error: %v", rbErr)
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *tagRepo) RemoveTagsFromPost(postID int) error {
	_, err := DB.Exec("DELETE FROM post_tags WHERE post_id = ?", postID)
	return err
}
