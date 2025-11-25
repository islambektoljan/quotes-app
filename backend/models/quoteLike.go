package models

import "time"

type QuoteLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	QuoteID   uint      `gorm:"not null" json:"quote_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Type      string    `gorm:"size:10;check:type IN ('like', 'dislike')" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
