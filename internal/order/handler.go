package order

import (
	"ecommerce-stock-api/models"
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

func (h *Handler) Checkout(c *gin.Context) {
	var req struct {
		UserID uint               `json:"user_id"`
		Items  []models.OrderItem `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Checkout(req.UserID, req.Items); err != nil {
		log.Error().Err(err).Str("module", "order").Msg("Failed to Checkout")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order created and stock reserved"})
}
