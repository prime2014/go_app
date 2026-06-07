package blogs

import "gorm.io/gorm"

type BlogServices struct {
	Db *gorm.DB
}

func (b *BlogServices) Create(blogDto BlogDto) (*Blogs, error) {
	blog := &Blogs{
		Title:   blogDto.Title,
		Article: blogDto.Article,
	}

	result := b.Db.Create(blog)

	if result.Error != nil {
		return &Blogs{}, result.Error
	}

	return blog, nil

}
