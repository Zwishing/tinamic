package user

// SignUp struct to describe register a new user.
type SignUp struct {
	Name     string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
	UserRole string `json:"user_role" validate:"required,lte=25"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Name         string `json:"email" validate:"required,email,lte=255"`
	PasswordHash string `json:"password" validate:"required,lte=255"`
}
