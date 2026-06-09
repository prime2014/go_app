package comments

type CommentDto struct {
	Comment string `json:"comment" validate:"required"`
}
