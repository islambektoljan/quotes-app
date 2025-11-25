package models

import "time"

type Comment struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Content      string        `gorm:"type:text;not null" json:"content" binding:"required,min=1,max=500"`
	QuoteID      uint          `gorm:"not null" json:"quote_id"`
	UserID       *uint         `json:"user_id"`
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LikesCount   int           `gorm:"default:0" json:"likes_count"`
	CommentLikes []CommentLike `gorm:"foreignKey:CommentID" json:"comment_likes,omitempty"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
}

type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}
