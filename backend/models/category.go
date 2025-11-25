package models

type Category struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Quotes      []Quote `gorm:"foreignKey:CategoryID" json:"quotes,omitempty"`
}
