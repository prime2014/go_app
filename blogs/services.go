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

func (b *BlogServices) Edit(blogDto EditBlogDto, userID uint, blogID uint) (*Blogs, error) {
	var blog Blogs
	result := b.Db.First(&blog, blogID)

	if result.Error != nil {
		return &Blogs{}, result.Error
	}

	if blogDto.Title != nil {
		blog.Title = *blogDto.Title
	}

	if blogDto.Article != nil {
		blog.Article = *blogDto.Article
	}

	b.Db.Save(&blog)

	return &blog, nil
}
