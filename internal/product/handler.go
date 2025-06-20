package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ListProducts(c *gin.Context) {
	products, err := h.service.GetProductList()
	if err != nil {
		log.Error().Err(err).Str("module", "product").Msg("Failed to Get Product List")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
