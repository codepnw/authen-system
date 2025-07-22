package user

type CreateUserRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}
