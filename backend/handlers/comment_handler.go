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

type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{DB: config.DB}
}

// AddComment - добавление комментария к цитате
func (h *CommentHandler) AddComment(c *gin.Context) {
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

	var input models.CommentCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование цитаты
	var quote models.Quote
	if err := h.DB.First(&quote, quoteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	comment := models.Comment{
		Content: input.Content,
		QuoteID: uint(quoteID),
		UserID:  &userIDUint,
	}

	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	h.DB.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusCreated, comment)
}

// GetComments - получение комментариев для цитаты
func (h *CommentHandler) GetComments(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
		return
	}

	// Проверяем существование цитаты
	var quote models.Quote
	if err := h.DB.First(&quote, quoteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	var comments []models.Comment
	if err := h.DB.Preload("User").
		Where("quote_id = ?", quoteID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// LikeComment - лайк комментария
func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var existingLike models.CommentLike
	result := h.DB.Where("comment_id = ? AND user_id = ?", commentID, userIDUint).First(&existingLike)

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var comment models.Comment
		if err := tx.First(&comment, commentID).Error; err != nil {
			return err
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Создаем новый лайк
			like := models.CommentLike{
				CommentID: uint(commentID),
				UserID:    userIDUint,
			}
			if err := tx.Create(&like).Error; err != nil {
				return err
			}
			comment.LikesCount++
		} else {
			// Удаляем существующий лайк
			if err := tx.Delete(&existingLike).Error; err != nil {
				return err
			}
			comment.LikesCount--
		}

		return tx.Save(&comment).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like updated successfully"})
}

// DeleteComment - удаление комментария
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var comment models.Comment
	if err := h.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Проверяем владельца
	if comment.UserID == nil || *comment.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	if err := h.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// UpdateComment - обновление комментария
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint := userID.(uint)

	var comment models.Comment
	if err := h.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Проверяем владельца
	if comment.UserID == nil || *comment.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own comments"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=1,max=500"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.Content = input.Content
	if err := h.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
