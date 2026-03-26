package auth_dto

type SignUpDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LogInDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenDto struct {
	Role string `json:"role"`
}
