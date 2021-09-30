package swagger

type UserRequest struct {
	Email    string `json:"email" validate:"required" example:"test@test.test"`
	Password string `json:"password" validate:"required" example:"password"`
}

type UpdateUser struct {
	Email    string `json:"email" validate:"required" example:"test@test.test"`
	Password string `json:"password,omitempty" example:"password"`
}
