package handlers

import (
	"net/http"
	"quotes-app/config"
	"quotes-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{DB: config.DB}
}

// GetCategories - получение всех категорий
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := h.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
