package blogs

type BlogDto struct {
	Title   string `json:"title" validate:"required"`
	Article string `json:"article" validate:"required"`
}

type EditBlogDto struct {
	// Title is optional, but required if Article is missing.
	// It cannot be an empty string if provided.
	Title *string `json:"title" validate:"omitempty,required_without=Article,min=1"`

	// Article is optional, but required if Title is missing.
	// It cannot be an empty string if provided.
	Article *string `json:"article" validate:"omitempty,required_without=Title,min=1"`
}
