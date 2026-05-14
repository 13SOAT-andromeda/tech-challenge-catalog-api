package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/response"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	category, err := h.service.Create(req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "category not found")
		return
	}
	response.OK(c, category)
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.service.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, categories)
}

type updateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	category, err := h.service.Update(id, req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "category deleted"})
}
