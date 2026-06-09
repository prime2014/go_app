package blogs

import (
	"comments"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type BlogStatus string

// 2. Declare your status constants
const (
	StatusDraft     BlogStatus = "draft"
	StatusPublished BlogStatus = "published"
	StatusArchived  BlogStatus = "archived"
)

type Blogs struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Title   string `json:"title" gorm:"size:100"`
	Slug    string `json:"slug" gorm:"size:100"`
	UserID  uint   `json:"user_id"`
	Article string `json:"article" gorm:"type:text;not null"`

	Status BlogStatus `json:"status" gorm:"type:varchar(20);default:'draft';not null;index"`

	Likes     int            `json:"likes" gorm:"default:0"`
	Views     int            `json:"views" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Images   []BlogImages        `json:"images" gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;"`
	Comments []comments.Comments `json:"comments" gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;"`
}

func (b *Blogs) BeforeCreate(tx *gorm.DB) (err error) {
	b.Slug = slug.Make(b.Title)
	return nil
}

type BlogImages struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Src    string `json:"src" gorm:"not null"`
	Width  int    `json:"width"`
	Height int    `json:"height"`

	BlogID uint `json:"blog_id"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
