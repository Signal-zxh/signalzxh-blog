package agent

import (
	"encoding/json"
	"strconv"

	"github.com/Signal-zxh/signalzxh-blog/service"
)

type ToolService struct {
	postService service.PostService
}

func NewToolService(postService service.PostService) *ToolService {
	return &ToolService{postService: postService}
}

func (s *ToolService) GetPosts() string {
	posts, err := s.postService.GetPosts()
	if err != nil {
		return "获取文章列表失败: " + err.Error()
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return "序列化文章列表失败: " + err.Error()
	}

	return string(data)
}

func (s *ToolService) GetPostByID(id string) string {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return "无效的文章ID: " + id
	}

	post, err := s.postService.GetPostByID(postID)
	if err != nil {
		return "获取文章失败: " + err.Error()
	}

	data, err := json.Marshal(post)
	if err != nil {
		return "序列化文章失败: " + err.Error()
	}

	return string(data)
}
