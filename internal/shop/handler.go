package shop

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
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(req.Name); err != nil {
		log.Error().Err(err).Str("module", "shop").Msg("Failed to Create Shop")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create shop"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shop created"})
}

func (h *Handler) List(c *gin.Context) {
	shops, err := h.service.List()
	if err != nil {
		log.Error().Err(err).Str("module", "shop").Msg("Failed to Get List Shops")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list shops"})
		return
	}
	c.JSON(http.StatusOK, shops)
}
