package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

func TestCategoryRepo_GetCategories(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go")
	mock.ExpectQuery("SELECT id, name FROM categories ORDER BY id DESC").WillReturnRows(rows)

	categories, err := repo.GetCategories()
	if err != nil {
		t.Fatalf("GetCategories error: %v", err)
	}
	if len(categories) != 1 || categories[0].Name != "Go" {
		t.Fatalf("unexpected categories: %+v", categories)
	}
}

func TestCategoryRepo_CreateCategory(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	mock.ExpectExec("INSERT INTO categories\\(name\\) VALUES\\(\\?\\)").WithArgs("Go").WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := repo.CreateCategory("Go")
	if err != nil {
		t.Fatalf("CreateCategory error: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected insert id 1, got %d", id)
	}
}

func TestCategoryRepo_UpdateCategory_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	mock.ExpectExec("UPDATE categories SET name = \\? WHERE id = \\?").WithArgs("Go", 99).WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.UpdateCategory(99, "Go")
	if err != db.ErrNoRowsAffected {
		t.Fatalf("expected ErrNoRowsAffected, got %v", err)
	}
}

func TestCategoryRepo_DeleteCategory(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	mock.ExpectExec("DELETE FROM categories WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.DeleteCategory(1)
	if err != nil {
		t.Fatalf("DeleteCategory error: %v", err)
	}
}

func TestCategoryRepo_GetCategoryByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go")
	mock.ExpectQuery("SELECT id, name FROM categories WHERE id = ?").WithArgs(1).WillReturnRows(rows)

	cat, err := repo.GetCategoryByID(1)
	if err != nil {
		t.Fatalf("GetCategoryByID error: %v", err)
	}
	if cat.ID != 1 || cat.Name != "Go" {
		t.Fatalf("unexpected category: %+v", cat)
	}
}

var _ = model.Category{}

func TestCategoryRepo_GetCategoryByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.CategoryRepoImpl

	mock.ExpectQuery("SELECT id, name FROM categories WHERE id = ?").WithArgs(99).WillReturnError(sql.ErrNoRows)
	_, err = repo.GetCategoryByID(99)
	if err != db.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
