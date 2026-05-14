package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/response"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.Create(&product); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "product not found")
		return
	}

	response.OK(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.service.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, products)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	product.ID = id

	if err := h.service.Update(&product); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, nil)
}
