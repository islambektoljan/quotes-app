package handlers

import (
	"errors"
	"net/http"
	"quotes-app/config"
	"quotes-app/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuoteHandler struct {
	DB *gorm.DB
}

func NewQuoteHandler() *QuoteHandler {
	return &QuoteHandler{DB: config.DB}
}

// GetQuotes - получение цитат с фильтрацией и пагинацией
func (h *QuoteHandler) GetQuotes(c *gin.Context) {
	var quotes []models.Quote

	query := h.DB.Preload("User").Preload("Category")

	// Фильтрация по категории
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			query = query.Where("category_id = ?", id)
		}
	}

	// Фильтрация по автору цитаты
	if author := c.Query("author"); author != "" {
		query = query.Where("author ILIKE ?", "%"+author+"%")
	}

	// Поиск по содержанию
	if content := c.Query("content"); content != "" {
		query = query.Where("content ILIKE ?", "%"+content+"%")
	}

	// Сортировка
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	query = query.Order(sort + " " + order)

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Model(&models.Quote{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quotes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quotes": quotes,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (int(total) + limit - 1) / limit,
		},
	})
}

// GetQuoteByID - получение цитаты по ID
func (h *QuoteHandler) GetQuoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	var quote models.Quote
	if err := h.DB.Preload("User").Preload("Category").
		Preload("Comments").Preload("Comments.User").
		First(&quote, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quote"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// CreateQuote - создание цитаты
func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var input models.QuoteCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	// Проверяем существование категории
	var category models.Category
	if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	quote := models.Quote{
		Content:    input.Content,
		Author:     input.Author,
		CategoryID: &input.CategoryID,
		UserID:     &userIDUint, // Теперь правильно
	}

	if err := h.DB.Create(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quote"})
		return
	}

	h.DB.Preload("User").Preload("Category").First(&quote, quote.ID)
	c.JSON(http.StatusCreated, quote)
}

// UpdateQuote - обновление цитаты
func (h *QuoteHandler) UpdateQuote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var quote models.Quote
	if err := h.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	// Проверяем владельца
	if quote.UserID == nil || *quote.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own quotes"})
		return
	}

	var input models.QuoteUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем только переданные поля
	updates := make(map[string]interface{})
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Author != "" {
		updates["author"] = input.Author
	}
	if input.CategoryID != 0 {
		// Проверяем существование категории
		var category models.Category
		if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		updates["category_id"] = input.CategoryID
	}

	if err := h.DB.Model(&quote).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote"})
		return
	}

	h.DB.Preload("User").Preload("Category").First(&quote, quote.ID)
	c.JSON(http.StatusOK, quote)
}

// DeleteQuote - удаление цитаты
func (h *QuoteHandler) DeleteQuote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var quote models.Quote
	if err := h.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	// Проверяем владельца
	if quote.UserID == nil || *quote.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own quotes"})
		return
	}

	if err := h.DB.Delete(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
}

// LikeQuote - лайк цитаты
func (h *QuoteHandler) LikeQuote(c *gin.Context) {
	h.handleQuoteReaction(c, "like")
}

// DislikeQuote - дизлайк цитаты
func (h *QuoteHandler) DislikeQuote(c *gin.Context) {
	h.handleQuoteReaction(c, "dislike")
}

func (h *QuoteHandler) handleQuoteReaction(c *gin.Context, reactionType string) {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var existingLike models.QuoteLike
	result := h.DB.Where("quote_id = ? AND user_id = ?", quoteID, userIDUint).First(&existingLike)

	// Начинаем транзакцию
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var quote models.Quote
		if err := tx.First(&quote, quoteID).Error; err != nil {
			return err
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Создаем новую реакцию
			like := models.QuoteLike{
				QuoteID: uint(quoteID),
				UserID:  userIDUint,
				Type:    reactionType,
			}
			if err := tx.Create(&like).Error; err != nil {
				return err
			}

			// Обновляем счетчики
			if reactionType == "like" {
				quote.LikesCount++
			} else {
				quote.DislikesCount++
			}
		} else {
			if existingLike.Type == reactionType {
				// Удаляем существующую реакцию
				if err := tx.Delete(&existingLike).Error; err != nil {
					return err
				}
				// Уменьшаем счетчик
				if reactionType == "like" {
					quote.LikesCount--
				} else {
					quote.DislikesCount--
				}
			} else {
				// Меняем реакцию
				existingLike.Type = reactionType
				if err := tx.Save(&existingLike).Error; err != nil {
					return err
				}
				// Обновляем счетчики
				if reactionType == "like" {
					quote.LikesCount++
					quote.DislikesCount--
				} else {
					quote.LikesCount--
					quote.DislikesCount++
				}
			}
		}

		return tx.Save(&quote).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction updated successfully"})
}
