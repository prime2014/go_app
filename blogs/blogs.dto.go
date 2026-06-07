package blogs

type BlogDto struct {
	Title   string `json:"title" validate:"required"`
	Article string `json:"article" validate:"required"`
}
