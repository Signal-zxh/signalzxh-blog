package db

import (
	"database/sql"
	"errors"

	"github.com/Signal-zxh/signalzxh-blog/model"
)

var ErrNoRowsAffected = errors.New("no rows affected")
var ErrNotFound = errors.New("not found")

type PostRepo interface {
	GetPosts() ([]model.Post, error)
	GetPostsByPage(page, pageSize int) ([]model.Post, error)
	GetPostsCount() (int, error)
	CreatePost(post model.Post) (int64, error)
	UpdatePost(post model.Post) error
	DeletePost(id int) error
	GetPostByID(id int) (model.Post, error)
	GetPostsByCategory(categoryID int, page, pageSize int) ([]model.Post, error)
	GetPostsByCategoryCount(categoryID int) (int, error)
	GetPostsWithCategoryTag(id int) (model.PostWithCategoryTag, error)
	GetPostsWithCategoryTagByPage(page, pageSize int) ([]model.PostWithCategoryTag, error)
}

type postRepo struct{}

var PostRepoImpl = &postRepo{}

func (r *postRepo) GetPosts() ([]model.Post, error) {
	rows, err := DB.Query("SELECT id, title, content, user_id, category_id FROM posts ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) GetPostsByPage(page, pageSize int) ([]model.Post, error) {
	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		"SELECT id, title, content, user_id, category_id FROM posts ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) GetPostsCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *postRepo) CreatePost(post model.Post) (int64, error) {
	res, err := DB.Exec(
		"INSERT INTO posts(title, content, user_id, category_id) VALUES(?, ?, ?, ?)",
		post.Title, post.Content, post.UserID, post.CategoryID,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *postRepo) UpdatePost(post model.Post) error {
	res, err := DB.Exec("UPDATE posts SET title = ?, content = ?, category_id = ? WHERE id = ?", post.Title, post.Content, post.CategoryID, post.ID)
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

func (r *postRepo) DeletePost(id int) error {
	res, err := DB.Exec("DELETE FROM posts WHERE id = ?", id)
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

func (r *postRepo) GetPostByID(id int) (model.Post, error) {
	row := DB.QueryRow("SELECT id, title, content, user_id, category_id FROM posts WHERE id = ?", id)

	var post model.Post

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}

	return post, nil
}

func (r *postRepo) GetPostsByCategory(categoryID int, page, pageSize int) ([]model.Post, error) {
	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		"SELECT id, title, content, user_id, category_id FROM posts WHERE category_id = ? ORDER BY id DESC LIMIT ? OFFSET ?",
		categoryID, pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostsByCategoryCount(categoryID int) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE category_id = ?", categoryID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *postRepo) GetPostsWithCategoryTag(id int) (model.PostWithCategoryTag, error) {
	var post model.PostWithCategoryTag

	row := DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.user_id, COALESCE(c.name, '') as category
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?`, id)

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.PostWithCategoryTag{}, ErrNotFound
		}
		return model.PostWithCategoryTag{}, err
	}

	tags, err := TagRepoImpl.GetTagsByPostID(id)
	if err != nil {
		return model.PostWithCategoryTag{}, err
	}

	post.Tags = make([]string, len(tags))
	for i, t := range tags {
		post.Tags[i] = t.Name
	}

	return post, nil
}

func (r *postRepo) GetPostsWithCategoryTagByPage(page, pageSize int) ([]model.PostWithCategoryTag, error) {
	offset := (page - 1) * pageSize

	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.user_id, COALESCE(c.name, '') as category
		FROM posts p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.PostWithCategoryTag
	for rows.Next() {
		var post model.PostWithCategoryTag
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Category); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range posts {
		tags, err := TagRepoImpl.GetTagsByPostID(posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].Tags = make([]string, len(tags))
		for j, t := range tags {
			posts[i].Tags[j] = t.Name
		}
	}

	return posts, nil
}
