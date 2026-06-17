package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct{}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(posts))
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.Post

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	id, err := service.CreatePost(req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"id":    id,
		"title": req.Title,
	}))
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	var req model.Post

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	err = service.UpdatePost(id, req.Title)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
		} else {
			c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"message": "updated successfully",
		"id":      id,
		"title":   req.Title,
	}))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	err = service.DeletePost(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
		} else {
			c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"message": "deleted successfully",
	}))
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid id"))
		return
	}

	post, err := service.GetPostByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, model.Fail("post not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(post))
}
