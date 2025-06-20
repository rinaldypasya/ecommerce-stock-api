package user

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

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Register(req.Email, req.Phone, req.Password); err != nil {
		log.Error().Err(err).Str("module", "user").Msg("Failed to Register")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"` // email or phone
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Login(req.Identifier, req.Password)
	if err != nil {
		log.Error().Err(err).Str("module", "user").Msg("Failed to Login")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
