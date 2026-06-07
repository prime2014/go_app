package accounts

type SignUpDto struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
