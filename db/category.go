package db

import (
	"database/sql"

	"github.com/Signal-zxh/signalzxh-blog/model"
)

type CategoryRepo interface {
	GetCategories() ([]model.Category, error)
	GetCategoryByID(id int) (model.Category, error)
	CreateCategory(name string) (int64, error)
	UpdateCategory(id int, name string) error
	DeleteCategory(id int) error
}

type categoryRepo struct{}

var CategoryRepoImpl = &categoryRepo{}

func (r *categoryRepo) GetCategories() ([]model.Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) GetCategoryByID(id int) (model.Category, error) {
	row := DB.QueryRow("SELECT id, name FROM categories WHERE id = ?", id)
	var c model.Category
	err := row.Scan(&c.ID, &c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Category{}, ErrNotFound
		}
		return model.Category{}, err
	}
	return c, nil
}

func (r *categoryRepo) CreateCategory(name string) (int64, error) {
	res, err := DB.Exec("INSERT INTO categories(name) VALUES(?)", name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *categoryRepo) UpdateCategory(id int, name string) error {
	res, err := DB.Exec("UPDATE categories SET name = ? WHERE id = ?", name, id)
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

func (r *categoryRepo) DeleteCategory(id int) error {
	res, err := DB.Exec("DELETE FROM categories WHERE id = ?", id)
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
