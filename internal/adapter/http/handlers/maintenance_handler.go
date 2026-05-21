package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/adapter/http/response"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/application/services"
	"github.com/juliovaz/tech-challenge-catalog-api/internal/domain"
)

type MaintenanceHandler struct {
	service *services.MaintenanceService
}

func NewMaintenanceHandler(service *services.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

func (h *MaintenanceHandler) Create(c *gin.Context) {
	var maintenance domain.Maintenance
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.Create(&maintenance); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, maintenance)
}

func (h *MaintenanceHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	maintenance, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "maintenance not found")
		return
	}

	response.OK(c, maintenance)
}

func (h *MaintenanceHandler) List(c *gin.Context) {
	maintenances, err := h.service.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, maintenances)
}

func (h *MaintenanceHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var maintenance domain.Maintenance
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	maintenance.ID = id

	if err := h.service.Update(&maintenance); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, maintenance)
}

func (h *MaintenanceHandler) Delete(c *gin.Context) {
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

func (h *MaintenanceHandler) CheckBatch(c *gin.Context) {
	raw := c.Query("ids")
	if raw == "" {
		response.BadRequest(c, "ids query param is required")
		return
	}

	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	maintenances, err := h.service.GetBatch(ids)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	type maintenanceItem struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		PriceCents int64  `json:"price_cents"`
	}
	items := make([]maintenanceItem, len(maintenances))
	for i, m := range maintenances {
		items[i] = maintenanceItem{
			ID:         m.SerialID,
			Name:       m.Description,
			PriceCents: int64(m.BasePrice * 100),
		}
	}

	c.JSON(http.StatusOK, items)
}
