package comments

import "gorm.io/gorm"

type CommentService struct {
	Db *gorm.DB
}

func (c *CommentService) CreateComment(commentDto CommentDto, blogID, authorID uint) (*Comments, error) {

	comment := &Comments{
		BlogID:   blogID,
		AuthorID: authorID,
		Comment:  commentDto.Comment,
	}

	result := c.Db.Create(comment)

	if result.Error != nil {
		return &Comments{}, result.Error
	}

	return comment, nil
}
