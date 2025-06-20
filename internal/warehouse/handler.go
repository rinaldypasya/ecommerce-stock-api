package warehouse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name   string `json:"name"`
		ShopID uint   `json:"shop_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(req.Name, req.ShopID); err != nil {
		log.Error().Err(err).Str("module", "user").Msg("Failed to Create Warehouse")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create warehouse"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "warehouse created"})
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	var req struct {
		ID     uint `json:"id"`
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateStatus(req.ID, req.Active); err != nil {
		log.Error().Err(err).Str("module", "user").Msg("Failed to Update Status of Warehouse")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *Handler) List(c *gin.Context) {
	whs, err := h.service.List()
	if err != nil {
		log.Error().Err(err).Str("module", "user").Msg("Failed to Get List Warehouses")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list warehouses"})
		return
	}
	c.JSON(http.StatusOK, whs)
}

func (h *Handler) TransferStock(c *gin.Context) {
	var req struct {
		ProductID       uint `json:"product_id"`
		FromWarehouseID uint `json:"from_warehouse_id"`
		ToWarehouseID   uint `json:"to_warehouse_id"`
		Quantity        int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TransferStock(req.ProductID, req.FromWarehouseID, req.ToWarehouseID, req.Quantity)
	if err != nil {
		log.Error().Err(err).Str("module", "warehouse").Msg("Failed to Transfer Stock")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock transferred successfully"})
}
