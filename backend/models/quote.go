package models

import "time"

type Quote struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Content       string      `gorm:"type:text;not null" json:"content" binding:"required,min=1,max=1000"`
	Author        string      `gorm:"size:100" json:"author" binding:"required,min=1,max=100"`
	UserID        *uint       `json:"user_id"`
	User          User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID    *uint       `json:"category_id" binding:"required"`
	Category      Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	LikesCount    int         `gorm:"default:0" json:"likes_count"`
	DislikesCount int         `gorm:"default:0" json:"dislikes_count"`
	Comments      []Comment   `gorm:"foreignKey:QuoteID;constraint:OnDelete:CASCADE;" json:"comments,omitempty"`
	QuoteLikes    []QuoteLike `gorm:"foreignKey:QuoteID;constraint:OnDelete:CASCADE;" json:"quote_likes,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// QuoteCreateRequest для валидации при создании
type QuoteCreateRequest struct {
	Content    string `json:"content" binding:"required,min=1,max=1000"`
	Author     string `json:"author" binding:"required,min=1,max=100"`
	CategoryID uint   `json:"category_id" binding:"required"`
}

type QuoteUpdateRequest struct {
	Content    string `json:"content" binding:"omitempty,min=1,max=1000"`
	Author     string `json:"author" binding:"omitempty,min=1,max=100"`
	CategoryID uint   `json:"category_id" binding:"omitempty"`
}
