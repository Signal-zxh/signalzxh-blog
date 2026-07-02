package handler

import (
	"net/http"
	"strconv"

	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// @Summary 获取分类列表
// @Description 获取所有分类
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Category}
// @Failure 500 {object} model.Response
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("获取分类失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(categories))
}

// @Summary 获取单个分类
// @Description 根据ID获取分类
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} model.Response{data=model.Category}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的分类ID"))
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("分类不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("获取分类失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(category))
}

// @Summary 创建分类
// @Description 创建新分类
// @Tags categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name body string true "分类名称"
// @Success 200 {object} model.Response{data=int64}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("请求参数错误"))
		return
	}

	id, err := h.service.CreateCategory(req.Name)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("分类名称无效"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("创建分类失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(id))
}

// @Summary 更新分类
// @Description 更新分类信息
// @Tags categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Param name body string true "分类名称"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的分类ID"))
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("请求参数错误"))
		return
	}

	err = h.service.UpdateCategory(id, req.Name)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("参数无效"))
			return
		}
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("分类不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("更新分类失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(nil))
}

// @Summary 删除分类
// @Description 删除分类
// @Tags categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的分类ID"))
		return
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("参数无效"))
			return
		}
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("分类不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("删除分类失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(nil))
}
