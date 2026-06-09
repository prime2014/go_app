package comments

import (
	"time"
)

type Comments struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	BlogID    uint      `json:"blog_id" gorm:"not null"`           // Removed trailing semicolon
	AuthorID  uint      `json:"author_id" gorm:"not null"`         // Removed trailing semicolon
	Comment   string    `json:"comment" gorm:"type:text;not null"` // Changed comma to semicolon inside quotes
	Likes     uint      `json:"likes" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
