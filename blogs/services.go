package blogs

import "gorm.io/gorm"

type BlogServices struct {
	Db *gorm.DB
}

func (b *BlogServices) Create(blogDto BlogDto, userID uint) (*Blogs, error) {
	blog := &Blogs{
		Title:   blogDto.Title,
		Article: blogDto.Article,
		UserID:  userID,
	}

	result := b.Db.Create(blog)

	if result.Error != nil {
		return &Blogs{}, result.Error
	}

	return blog, nil

}
