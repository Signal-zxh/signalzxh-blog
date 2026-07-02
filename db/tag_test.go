package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Signal-zxh/signalzxh-blog/db"
)

func TestTagRepo_GetTags(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.TagRepoImpl

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go")
	mock.ExpectQuery("SELECT id, name FROM tags ORDER BY id DESC").WillReturnRows(rows)

	tags, err := repo.GetTags()
	if err != nil {
		t.Fatalf("GetTags error: %v", err)
	}
	if len(tags) != 1 || tags[0].Name != "Go" {
		t.Fatalf("unexpected tags: %+v", tags)
	}
}

func TestTagRepo_CreateTag(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.TagRepoImpl

	mock.ExpectExec("INSERT INTO tags\\(name\\) VALUES\\(\\?\\)").WithArgs("Go").WillReturnResult(sqlmock.NewResult(1, 1))
	id, err := repo.CreateTag("Go")
	if err != nil {
		t.Fatalf("CreateTag error: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected insert id 1, got %d", id)
	}
}

func TestTagRepo_UpdateTag_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.TagRepoImpl

	mock.ExpectExec("UPDATE tags SET name = \\? WHERE id = \\?").WithArgs("Go", 99).WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.UpdateTag(99, "Go")
	if err != db.ErrNoRowsAffected {
		t.Fatalf("expected ErrNoRowsAffected, got %v", err)
	}
}

func TestTagRepo_DeleteTag(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.TagRepoImpl

	mock.ExpectExec("DELETE FROM tags WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.DeleteTag(1)
	if err != nil {
		t.Fatalf("DeleteTag error: %v", err)
	}
}

func TestTagRepo_GetTagByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock db failed: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB
	repo := db.TagRepoImpl

	mock.ExpectQuery("SELECT id, name FROM tags WHERE id = ?").WithArgs(99).WillReturnError(sql.ErrNoRows)
	_, err = repo.GetTagByID(99)
	if err != db.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
