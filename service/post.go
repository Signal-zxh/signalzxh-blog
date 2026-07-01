package service

import (
	"errors"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service/cache"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type PostService interface {
	GetPostByID(id int) (model.Post, error)
	GetPosts() ([]model.Post, error)
	GetPostsByPage(page, pageSize int) ([]model.Post, int, error)
	CreatePost(title, content string, userID int) (int64, error)
	UpdatePost(id int, title, content string) error
	DeletePost(id int) error
	GetPostWithCategoryTag(id int) (model.PostWithCategoryTag, error)
	GetPostsWithCategoryTagByPage(page, pageSize int) ([]model.PostWithCategoryTag, int, error)
	GetPostsByCategory(categoryID, page, pageSize int) ([]model.Post, int, error)
	CreatePostWithCategoryTag(title, content string, userID, categoryID int, tagNames []string) (int64, error)
	UpdatePostWithCategoryTag(id, categoryID int, title, content string, tagNames []string) error
}

type postService struct {
	repo  db.PostRepo
	cache cache.PostCache
	tagDB db.TagRepo
}

func NewPostService(repo db.PostRepo, c cache.PostCache) PostService {
	return &postService{repo: repo, cache: c, tagDB: db.TagRepoImpl}
}

func (s *postService) GetPostByID(id int) (model.Post, error) {
	post, found, err := s.cache.GetPostByID(id)
	if err == nil && found {
		if post.ID == 0 {
			return model.Post{}, ErrNotFound
		}
		return post, nil
	}

	post, err = s.repo.GetPostByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			_ = s.cache.SetNilPost(id, 1*time.Minute) //nolint:errcheck
			return model.Post{}, ErrNotFound
		}
		return model.Post{}, err
	}

	_ = s.cache.SetPost(post, 10*time.Minute) //nolint:errcheck
	return post, nil
}

func (s *postService) GetPosts() ([]model.Post, error) {
	posts, found, err := s.cache.GetPosts()
	if err == nil && found {
		return posts, nil
	}

	posts, err = s.repo.GetPosts()
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetPosts(posts, 10*time.Minute) //nolint:errcheck
	return posts, nil
}

func (s *postService) GetPostsByPage(page, pageSize int) ([]model.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	posts, found, err := s.cache.GetPostsByPage(page, pageSize)
	if err == nil && found {
		count, err := s.repo.GetPostsCount()
		if err != nil {
			return posts, 0, err
		}
		return posts, count, nil
	}

	posts, err = s.repo.GetPostsByPage(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetPostsCount()
	if err != nil {
		return nil, 0, err
	}

	_ = s.cache.SetPostsByPage(posts, page, pageSize, 10*time.Minute) //nolint:errcheck
	return posts, count, nil
}

func (s *postService) CreatePost(title, content string, userID int) (int64, error) {
	if title == "" || len(title) > 255 {
		return 0, ErrInvalidInput
	}

	post := model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	id, err := s.repo.CreatePost(post)
	if err != nil {
		return 0, err
	}

	_ = s.cache.InvalidatePosts() //nolint:errcheck
	return id, nil
}

func (s *postService) UpdatePost(id int, title, content string) error {
	if title == "" || len(title) > 255 {
		return ErrInvalidInput
	}

	post := model.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}

	err := s.repo.UpdatePost(post)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}

	_ = s.cache.InvalidatePost(id) //nolint:errcheck
	return nil
}

func (s *postService) DeletePost(id int) error {
	err := s.repo.DeletePost(id)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}

	_ = s.cache.InvalidatePost(id) //nolint:errcheck
	_ = s.cache.InvalidatePosts()  //nolint:errcheck
	return nil
}

func (s *postService) GetPostWithCategoryTag(id int) (model.PostWithCategoryTag, error) {
	if id <= 0 {
		return model.PostWithCategoryTag{}, ErrInvalidInput
	}
	post, err := s.repo.GetPostsWithCategoryTag(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return model.PostWithCategoryTag{}, ErrNotFound
		}
		return model.PostWithCategoryTag{}, err
	}
	return post, nil
}

func (s *postService) GetPostsWithCategoryTagByPage(page, pageSize int) ([]model.PostWithCategoryTag, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	posts, err := s.repo.GetPostsWithCategoryTagByPage(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetPostsCount()
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) GetPostsByCategory(categoryID, page, pageSize int) ([]model.Post, int, error) {
	if categoryID <= 0 {
		return nil, 0, ErrInvalidInput
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	posts, err := s.repo.GetPostsByCategory(categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetPostsByCategoryCount(categoryID)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *postService) CreatePostWithCategoryTag(title, content string, userID, categoryID int, tagNames []string) (int64, error) {
	if title == "" || len(title) > 255 {
		return 0, ErrInvalidInput
	}

	post := model.Post{
		Title:      title,
		Content:    content,
		UserID:     userID,
		CategoryID: categoryID,
	}

	id, err := s.repo.CreatePost(post)
	if err != nil {
		return 0, err
	}

	if len(tagNames) > 0 {
		var tagIDs []int
		for _, name := range tagNames {
			tagID, err := s.tagDB.GetOrCreateTag(name)
			if err != nil {
				return 0, err
			}
			tagIDs = append(tagIDs, int(tagID))
		}
		if err := s.tagDB.AddTagsToPost(int(id), tagIDs); err != nil {
			return 0, err
		}
	}

	_ = s.cache.InvalidatePosts() //nolint:errcheck
	return id, nil
}

func (s *postService) UpdatePostWithCategoryTag(id, categoryID int, title, content string, tagNames []string) error {
	if id <= 0 || title == "" || len(title) > 255 {
		return ErrInvalidInput
	}

	post := model.Post{
		ID:         id,
		Title:      title,
		Content:    content,
		CategoryID: categoryID,
	}

	err := s.repo.UpdatePost(post)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}

	if len(tagNames) > 0 {
		var tagIDs []int
		for _, name := range tagNames {
			tagID, err := s.tagDB.GetOrCreateTag(name)
			if err != nil {
				return err
			}
			tagIDs = append(tagIDs, int(tagID))
		}
		if err := s.tagDB.AddTagsToPost(id, tagIDs); err != nil {
			return err
		}
	} else {
		_ = s.tagDB.RemoveTagsFromPost(id) //nolint:errcheck
	}

	_ = s.cache.InvalidatePost(id) //nolint:errcheck
	_ = s.cache.InvalidatePosts()  //nolint:errcheck
	return nil
}
