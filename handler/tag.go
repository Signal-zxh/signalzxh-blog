package handler

import (
	"net/http"
	"strconv"

	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service service.TagService
}

func NewTagHandler(s service.TagService) *TagHandler {
	return &TagHandler{service: s}
}

// @Summary 获取标签列表
// @Description 获取所有标签
// @Tags tags
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=[]model.Tag}
// @Failure 500 {object} model.Response
// @Router /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.service.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("获取标签失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(tags))
}

// @Summary 获取单个标签
// @Description 根据ID获取标签
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} model.Response{data=model.Tag}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的标签ID"))
		return
	}

	tag, err := h.service.GetTagByID(id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("标签不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("获取标签失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(tag))
}

// @Summary 创建标签
// @Description 创建新标签
// @Tags tags
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name body string true "标签名称"
// @Success 200 {object} model.Response{data=int64}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("请求参数错误"))
		return
	}

	id, err := h.service.CreateTag(req.Name)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("标签名称无效"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("创建标签失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(id))
}

// @Summary 更新标签
// @Description 更新标签信息
// @Tags tags
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "标签ID"
// @Param name body string true "标签名称"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的标签ID"))
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("请求参数错误"))
		return
	}

	err = h.service.UpdateTag(id, req.Name)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("参数无效"))
			return
		}
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("标签不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("更新标签失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(nil))
}

// @Summary 删除标签
// @Description 删除标签
// @Tags tags
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "标签ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("无效的标签ID"))
		return
	}

	err = h.service.DeleteTag(id)
	if err != nil {
		if err == service.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, model.Fail("参数无效"))
			return
		}
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, model.Fail("标签不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.Fail("删除标签失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(nil))
}
